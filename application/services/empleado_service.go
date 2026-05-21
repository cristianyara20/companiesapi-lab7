package services

import (
	"companies-api/application/dtos"
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
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

func (s *EmpleadoService) GetAll() ([]entities.Empleado, error) {
	return s.uow.Empleados().GetAll()
}

func (s *EmpleadoService) GetById(id uint) (*entities.Empleado, error) {
	e, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return nil, errors.New("empleado no encontrado")
	}
	return e, nil
}

func (s *EmpleadoService) Create(dto dtos.CreateEmpleadoDTO) (*entities.Empleado, error) {
	s.logger.Info("Creando empleado", zap.String("nombre", dto.Nombre))

	// Verificar que la compañía existe antes de crear el empleado
	if _, err := s.uow.Companias().GetById(dto.CompaniaID); err != nil {
		return nil, errors.New("la compañía especificada no existe")
	}

	emp := &entities.Empleado{
		Nombre:     dto.Nombre,
		Apellido:   dto.Apellido,
		Correo:     dto.Correo,
		Cargo:      dto.Cargo,
		Salario:    dto.Salario,
		CompaniaID: dto.CompaniaID,
	}
	if err := s.uow.Empleados().Create(emp); err != nil {
		s.logger.Error("Error al crear empleado", zap.Error(err))
		return nil, err
	}
	s.logger.Info("✅ Empleado creado", zap.Uint("id", emp.ID))
	return emp, nil
}

func (s *EmpleadoService) Update(id uint, dto dtos.UpdateEmpleadoDTO) (*entities.Empleado, error) {
	emp, err := s.uow.Empleados().GetById(id)
	if err != nil {
		return nil, errors.New("empleado no encontrado")
	}
	if dto.Nombre != "" {
		emp.Nombre = dto.Nombre
	}
	if dto.Apellido != "" {
		emp.Apellido = dto.Apellido
	}
	if dto.Correo != "" {
		emp.Correo = dto.Correo
	}
	if dto.Cargo != "" {
		emp.Cargo = dto.Cargo
	}
	if dto.Salario > 0 {
		emp.Salario = dto.Salario
	}
	if dto.CompaniaID > 0 {
		emp.CompaniaID = dto.CompaniaID
	}
	if err := s.uow.Empleados().Update(emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *EmpleadoService) Delete(id uint) error {
	if _, err := s.uow.Empleados().GetById(id); err != nil {
		return errors.New("empleado no encontrado")
	}
	return s.uow.Empleados().Delete(id)
}
