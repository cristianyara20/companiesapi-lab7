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

type EmpleadoService struct {
	uow    interfaces.UnitOfWork
	logger *zap.Logger
}

func NewEmpleadoService(uow interfaces.UnitOfWork, logger *zap.Logger) *EmpleadoService {
	return &EmpleadoService{uow: uow, logger: logger}
}

func (s *EmpleadoService) GetAll(ctx context.Context) ([]entities.Empleado, error) {
	return s.uow.Empleados().GetAll(ctx)
}

func (s *EmpleadoService) GetById(ctx context.Context, id uint) (*entities.Empleado, error) {
	e, err := s.uow.Empleados().GetById(ctx, id)
	if err != nil {
		return nil, errors.New("empleado no encontrado")
	}
	return e, nil
}

func (s *EmpleadoService) Create(ctx context.Context, dto dtos.CreateEmpleadoDTO) (*entities.Empleado, error) {
	s.logger.Info("Creando empleado", zap.String("nombre", dto.Nombre))

	// 1. Validación estructural (Capa de Aplicación)
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	// 2. Validación de Negocio: Verificar que la compañía existe
	if _, err := s.uow.Companias().GetById(ctx, dto.CompaniaID); err != nil {
		return nil, &validation.ValidationError{
			Mensaje: "Error de validación de negocio",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "compania_id", Detalle: "La compañía especificada no existe"},
			},
		}
	}

	// 3. Validación de Negocio: Correo único
	existing, err := s.uow.Empleados().FindByCondition(ctx, "correo = ?", dto.Correo)
	if err == nil && len(existing) > 0 {
		return nil, &validation.ValidationError{
			Mensaje: "Error de validación de negocio",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "correo", Detalle: "El correo electrónico ya se encuentra registrado por otro empleado"},
			},
		}
	}

	emp := &entities.Empleado{
		Nombre:     dto.Nombre,
		Apellido:   dto.Apellido,
		Correo:     dto.Correo,
		Cargo:      dto.Cargo,
		Salario:    dto.Salario,
		CompaniaID: dto.CompaniaID,
	}
	if err := s.uow.Empleados().Create(ctx, emp); err != nil {
		s.logger.Error("Error al crear empleado", zap.Error(err))
		return nil, err
	}
	s.logger.Info("✅ Empleado creado", zap.Uint("id", emp.ID))
	return emp, nil
}

func (s *EmpleadoService) Update(ctx context.Context, id uint, dto dtos.UpdateEmpleadoDTO) (*entities.Empleado, error) {
	// 1. Validación estructural
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	emp, err := s.uow.Empleados().GetById(ctx, id)
	if err != nil {
		return nil, errors.New("empleado no encontrado")
	}

	// 2. Validaciones de Negocio para campos provistos
	if dto.CompaniaID != nil {
		if _, err := s.uow.Companias().GetById(ctx, *dto.CompaniaID); err != nil {
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "compania_id", Detalle: "La compañía especificada no existe"},
				},
			}
		}
		emp.CompaniaID = *dto.CompaniaID
	}

	if dto.Correo != nil {
		existing, err := s.uow.Empleados().FindByCondition(ctx, "correo = ? AND id != ?", *dto.Correo, id)
		if err == nil && len(existing) > 0 {
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "correo", Detalle: "El correo electrónico ya se encuentra registrado por otro empleado"},
				},
			}
		}
		emp.Correo = *dto.Correo
	}

	if dto.Nombre != nil {
		emp.Nombre = *dto.Nombre
	}
	if dto.Apellido != nil {
		emp.Apellido = *dto.Apellido
	}
	if dto.Cargo != nil {
		emp.Cargo = *dto.Cargo
	}
	if dto.Salario != nil {
		emp.Salario = *dto.Salario
	}

	if err := s.uow.Empleados().Update(ctx, emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *EmpleadoService) Delete(ctx context.Context, id uint) error {
	if _, err := s.uow.Empleados().GetById(ctx, id); err != nil {
		return errors.New("empleado no encontrado")
	}
	return s.uow.Empleados().Delete(ctx, id)
}

// CreateRange — Creación masiva (Bulk Create) en una sola transacción del Unit of Work
func (s *EmpleadoService) CreateRange(ctx context.Context, list []dtos.CreateEmpleadoDTO) ([]entities.Empleado, error) {
	s.logger.Info("Creando empleados en lote", zap.Int("cantidad", len(list)))

	if len(list) == 0 {
		return nil, &validation.ValidationError{
			Mensaje: "Error de validación",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "empleados", Detalle: "Debe proveer al menos un empleado"},
			},
		}
	}

	// 1. Iniciar transacción del Unit of Work
	if err := s.uow.BeginTransaction(); err != nil {
		return nil, err
	}

	var empleadosToCreate []entities.Empleado

	for i, empDTO := range list {
		// Validar DTO estructuralmente
		if err := validation.ValidateStruct(empDTO); err != nil {
			s.uow.Rollback()
			return nil, err
		}

		// Validar que la compañía existe
		if _, err := s.uow.Companias().GetById(ctx, empDTO.CompaniaID); err != nil {
			s.uow.Rollback()
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "compania_id", Detalle: "La compañía especificada en el índice " + string(rune(i)) + " no existe"},
				},
			}
		}

		// Validar correo único en DB
		exists, err := s.uow.Empleados().FindByCondition(ctx, "correo = ?", empDTO.Correo)
		if err == nil && len(exists) > 0 {
			s.uow.Rollback()
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "correo", Detalle: "El correo [" + empDTO.Correo + "] ya está registrado por otro empleado. Transacción cancelada."},
				},
			}
		}

		empleadosToCreate = append(empleadosToCreate, entities.Empleado{
			Nombre:     empDTO.Nombre,
			Apellido:   empDTO.Apellido,
			Correo:     empDTO.Correo,
			Cargo:      empDTO.Cargo,
			Salario:    empDTO.Salario,
			CompaniaID: empDTO.CompaniaID,
		})
	}

	// 2. Crear rango en lote
	if err := s.uow.Empleados().CreateRange(ctx, empleadosToCreate); err != nil {
		s.uow.Rollback()
		return nil, err
	}

	// 3. Confirmar transacción
	if err := s.uow.Commit(); err != nil {
		s.uow.Rollback()
		return nil, err
	}

	return empleadosToCreate, nil
}

// DeleteRange — Eliminación múltiple en una sola transacción del Unit of Work
func (s *EmpleadoService) DeleteRange(ctx context.Context, ids []uint) error {
	s.logger.Info("Eliminando empleados en lote", zap.Int("cantidad", len(ids)))

	if len(ids) == 0 {
		return &validation.ValidationError{
			Mensaje: "Error de validación",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "ids", Detalle: "Debe proveer al menos un ID"},
			},
		}
	}

	if err := s.uow.BeginTransaction(); err != nil {
		return err
	}

	// Verificar existencia de cada ID antes de eliminar
	for _, id := range ids {
		if _, err := s.uow.Empleados().GetById(ctx, id); err != nil {
			s.uow.Rollback()
			return errors.New("empleado con ID " + string(rune(id)) + " no existe")
		}
	}

	if err := s.uow.Empleados().DeleteRange(ctx, ids); err != nil {
		s.uow.Rollback()
		return err
	}

	return s.uow.Commit()
}

// GetPaged — Listado paginado, filtrado y ordenado
func (s *EmpleadoService) GetPaged(ctx context.Context, pagina int, tamano int, orden string, dir string, buscar string) ([]entities.Empleado, int64, error) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 {
		tamano = 10
	}
	return s.uow.Empleados().GetPaged(ctx, pagina, tamano, orden, dir, buscar)
}

// GetPagedByCompaniaID — Listado paginado de empleados de una compañía
func (s *EmpleadoService) GetPagedByCompaniaID(ctx context.Context, companiaID uint, pagina int, tamano int) ([]entities.Empleado, int64, error) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 {
		tamano = 10
	}
	return s.uow.Empleados().GetPagedByCompaniaID(ctx, companiaID, pagina, tamano)
}

// PatchPartial — Actualización parcial de empleado (Capa de Aplicación)
func (s *EmpleadoService) PatchPartial(ctx context.Context, id uint, dto dtos.UpdateEmpleadoDTO) (*entities.Empleado, error) {
	s.logger.Info("Actualizando parcialmente empleado", zap.Uint("id", id))

	// 1. Validación estructural
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	// Verificar existencia
	if _, err := s.uow.Empleados().GetById(ctx, id); err != nil {
		return nil, errors.New("empleado no encontrado")
	}

	// Construir mapa de cambios
	changes := make(map[string]interface{})

	if dto.Nombre != nil {
		changes["nombre"] = *dto.Nombre
	}
	if dto.Apellido != nil {
		changes["apellido"] = *dto.Apellido
	}
	if dto.Cargo != nil {
		changes["cargo"] = *dto.Cargo
	}
	if dto.Salario != nil {
		changes["salario"] = *dto.Salario
	}

	if dto.CompaniaID != nil {
		// Validar que la compañía existe
		if _, err := s.uow.Companias().GetById(ctx, *dto.CompaniaID); err != nil {
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "compania_id", Detalle: "La compañía especificada no existe"},
				},
			}
		}
		changes["compania_id"] = *dto.CompaniaID
	}

	if dto.Correo != nil {
		// Validar correo único
		existing, err := s.uow.Empleados().FindByCondition(ctx, "correo = ? AND id != ?", *dto.Correo, id)
		if err == nil && len(existing) > 0 {
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "correo", Detalle: "El correo electrónico ya se encuentra registrado por otro empleado"},
				},
			}
		}
		changes["correo"] = *dto.Correo
	}

	if len(changes) == 0 {
		return s.uow.Empleados().GetById(ctx, id)
	}

	if err := s.uow.Empleados().PatchPartial(ctx, id, changes); err != nil {
		return nil, err
	}

	return s.uow.Empleados().GetById(ctx, id)
}
