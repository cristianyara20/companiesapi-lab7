package services

import (
	"companies-api/application/dtos"
	"companies-api/application/validation"
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"context"
	"errors"

	"go.uber.org/zap"
)

// CompaniaService contiene la lógica de negocio de compañías
type CompaniaService struct {
	uow    interfaces.UnitOfWork
	logger *zap.Logger
}

func NewCompaniaService(uow interfaces.UnitOfWork, logger *zap.Logger) *CompaniaService {
	return &CompaniaService{uow: uow, logger: logger}
}

func (s *CompaniaService) GetAll(ctx context.Context) ([]entities.Compania, error) {
	return s.uow.Companias().GetAll(ctx)
}

func (s *CompaniaService) GetById(ctx context.Context, id uint) (*entities.Compania, error) {
	c, err := s.uow.Companias().GetById(ctx, id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	return c, nil
}

func (s *CompaniaService) Create(ctx context.Context, dto dtos.CreateCompaniaDTO) (*entities.Compania, error) {
	s.logger.Info("Creando compañía", zap.String("nombre", dto.Nombre))

	// 1. Validación estructural (Capa de Aplicación)
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	c := &entities.Compania{
		Nombre:    dto.Nombre,
		Direccion: dto.Direccion,
		Telefono:  dto.Telefono,
	}
	if err := s.uow.Companias().Create(ctx, c); err != nil {
		s.logger.Error("Error al crear compañía", zap.Error(err))
		return nil, err
	}
	s.logger.Info("✅ Compañía creada", zap.Uint("id", c.ID))
	return c, nil
}

func (s *CompaniaService) Update(ctx context.Context, id uint, dto dtos.UpdateCompaniaDTO) (*entities.Compania, error) {
	// 1. Validación estructural (Capa de Aplicación)
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	c, err := s.uow.Companias().GetById(ctx, id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	
	if dto.Nombre != nil {
		c.Nombre = *dto.Nombre
	}
	if dto.Direccion != nil {
		c.Direccion = *dto.Direccion
	}
	if dto.Telefono != nil {
		c.Telefono = *dto.Telefono
	}
	
	if err := s.uow.Companias().Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CompaniaService) Delete(ctx context.Context, id uint) error {
	if _, err := s.uow.Companias().GetById(ctx, id); err != nil {
		return errors.New("compañía no encontrada")
	}
	return s.uow.Companias().Delete(ctx, id)
}

func (s *CompaniaService) GetEmpleados(ctx context.Context, id uint) (*entities.Compania, error) {
	c, err := s.uow.Companias().GetWithEmpleados(ctx, id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	return c, nil
}

func (s *CompaniaService) GetEmpleadosPaged(ctx context.Context, id uint, pagina, tamano int) ([]entities.Empleado, int64, error) {
	if _, err := s.uow.Companias().GetById(ctx, id); err != nil {
		return nil, 0, errors.New("compañía no encontrada")
	}
	return s.uow.Empleados().GetPagedByCompaniaID(ctx, id, pagina, tamano)
}

// CreateConEmpleados — CASO TRANSACCIONAL OBLIGATORIO (con Contexto y Validación)
func (s *CompaniaService) CreateConEmpleados(ctx context.Context, dto dtos.CreateCompaniaConEmpleadosDTO) (*entities.Compania, error) {
	s.logger.Info("🔄 INICIO de transacción",
		zap.String("compania", dto.Nombre),
		zap.Int("empleados", len(dto.Empleados)),
	)

	// 1. Validación estructural de toda la petición (incluidos los empleados anidados)
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	// 2. Iniciar transacción del Unit of Work
	if err := s.uow.BeginTransaction(); err != nil {
		s.logger.Error("Error iniciando transacción", zap.Error(err))
		return nil, err
	}

	// 3. Crear compañía dentro de la transacción
	compania := &entities.Compania{
		Nombre:    dto.Nombre,
		Direccion: dto.Direccion,
		Telefono:  dto.Telefono,
	}
	if err := s.uow.Companias().Create(ctx, compania); err != nil {
		s.logger.Error("❌ Error en compañía → ROLLBACK", zap.Error(err))
		s.uow.Rollback()
		return nil, err
	}
	s.logger.Info("Compañía lista en transacción", zap.Uint("id", compania.ID))

	// 4. Crear cada empleado dentro de la misma transacción validando reglas de negocio
	for _, empDTO := range dto.Empleados {
		// Validar Regla de Negocio: Correo Único en DB
		exists, err := s.uow.Empleados().FindByCondition(ctx, "correo = ?", empDTO.Correo)
		if err == nil && len(exists) > 0 {
			s.logger.Error("❌ Correo duplicado detectado en lote → ROLLBACK", zap.String("correo", empDTO.Correo))
			s.uow.Rollback()
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "correo", Detalle: "El correo [" + empDTO.Correo + "] ya está registrado por otro empleado. Transacción cancelada."},
				},
			}
		}

		emp := &entities.Empleado{
			Nombre:     empDTO.Nombre,
			Apellido:   empDTO.Apellido,
			Correo:     empDTO.Correo,
			Cargo:      empDTO.Cargo,
			Salario:    empDTO.Salario,
			CompaniaID: compania.ID,
		}
		if err := s.uow.Empleados().Create(ctx, emp); err != nil {
			s.logger.Error("❌ Error en empleado → ROLLBACK TOTAL",
				zap.String("correo", empDTO.Correo),
				zap.Error(err),
			)
			s.uow.Rollback()
			return nil, errors.New("error en empleado [" + empDTO.Correo + "]: toda la operación fue revertida")
		}
		s.logger.Info("Empleado listo en transacción", zap.String("nombre", empDTO.Nombre))
	}

	// 5. Confirmar todo junto
	if err := s.uow.Commit(); err != nil {
		s.logger.Error("❌ Error en Commit → ROLLBACK", zap.Error(err))
		s.uow.Rollback()
		return nil, err
	}

	s.logger.Info("✅ COMMIT exitoso",
		zap.Uint("compania_id", compania.ID),
		zap.Int("empleados_creados", len(dto.Empleados)),
	)
	return compania, nil
}
