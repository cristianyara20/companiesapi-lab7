package services

import (
	"companies-api/application/dtos"
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"errors"

	"go.uber.org/zap"
)

// CompaniaService contiene la lógica de negocio de compañías
// REGLA IMPORTANTE: nunca accede directamente al ORM
// Siempre usa el UnitOfWork para acceder a los repositorios
type CompaniaService struct {
	uow    interfaces.UnitOfWork
	logger *zap.Logger
}

func NewCompaniaService(uow interfaces.UnitOfWork, logger *zap.Logger) *CompaniaService {
	return &CompaniaService{uow: uow, logger: logger}
}

func (s *CompaniaService) GetAll() ([]entities.Compania, error) {
	return s.uow.Companias().GetAll()
}

func (s *CompaniaService) GetById(id uint) (*entities.Compania, error) {
	c, err := s.uow.Companias().GetById(id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	return c, nil
}

func (s *CompaniaService) Create(dto dtos.CreateCompaniaDTO) (*entities.Compania, error) {
	s.logger.Info("Creando compañía", zap.String("nombre", dto.Nombre))

	c := &entities.Compania{
		Nombre:    dto.Nombre,
		Direccion: dto.Direccion,
		Telefono:  dto.Telefono,
	}
	if err := s.uow.Companias().Create(c); err != nil {
		s.logger.Error("Error al crear compañía", zap.Error(err))
		return nil, err
	}
	s.logger.Info("✅ Compañía creada", zap.Uint("id", c.ID))
	return c, nil
}

func (s *CompaniaService) Update(id uint, dto dtos.UpdateCompaniaDTO) (*entities.Compania, error) {
	c, err := s.uow.Companias().GetById(id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	if dto.Nombre != "" {
		c.Nombre = dto.Nombre
	}
	if dto.Direccion != "" {
		c.Direccion = dto.Direccion
	}
	if dto.Telefono != "" {
		c.Telefono = dto.Telefono
	}
	if err := s.uow.Companias().Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CompaniaService) Delete(id uint) error {
	if _, err := s.uow.Companias().GetById(id); err != nil {
		return errors.New("compañía no encontrada")
	}
	return s.uow.Companias().Delete(id)
}

func (s *CompaniaService) GetEmpleados(id uint) (*entities.Compania, error) {
	c, err := s.uow.Companias().GetWithEmpleados(id)
	if err != nil {
		return nil, errors.New("compañía no encontrada")
	}
	return c, nil
}

// CreateConEmpleados — CASO TRANSACCIONAL OBLIGATORIO
// Crea una compañía y todos sus empleados en UNA sola transacción.
// Si falla CUALQUIER empleado → Rollback → no se guarda NADA.
func (s *CompaniaService) CreateConEmpleados(dto dtos.CreateCompaniaConEmpleadosDTO) (*entities.Compania, error) {
	s.logger.Info("🔄 INICIO de transacción",
		zap.String("compania", dto.Nombre),
		zap.Int("empleados", len(dto.Empleados)),
	)

	// PASO 1: iniciar transacción
	if err := s.uow.BeginTransaction(); err != nil {
		s.logger.Error("Error iniciando transacción", zap.Error(err))
		return nil, err
	}

	// PASO 2: crear compañía dentro de la transacción
	compania := &entities.Compania{
		Nombre:    dto.Nombre,
		Direccion: dto.Direccion,
		Telefono:  dto.Telefono,
	}
	if err := s.uow.Companias().Create(compania); err != nil {
		s.logger.Error("❌ Error en compañía → ROLLBACK", zap.Error(err))
		s.uow.Rollback()
		return nil, err
	}
	s.logger.Info("Compañía lista en transacción", zap.Uint("id", compania.ID))

	// PASO 3: crear cada empleado dentro de la MISMA transacción
	for _, empDTO := range dto.Empleados {
		emp := &entities.Empleado{
			Nombre:     empDTO.Nombre,
			Apellido:   empDTO.Apellido,
			Correo:     empDTO.Correo,
			Cargo:      empDTO.Cargo,
			Salario:    empDTO.Salario,
			CompaniaID: compania.ID,
		}
		if err := s.uow.Empleados().Create(emp); err != nil {
			s.logger.Error("❌ Error en empleado → ROLLBACK TOTAL",
				zap.String("correo", empDTO.Correo),
				zap.Error(err),
			)
			// Revierte la compañía Y todos los empleados ya creados
			s.uow.Rollback()
			return nil, errors.New("error en empleado [" + empDTO.Correo + "]: toda la operación fue revertida")
		}
		s.logger.Info("Empleado listo en transacción", zap.String("nombre", empDTO.Nombre))
	}

	// PASO 4: confirmar todo junto
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
