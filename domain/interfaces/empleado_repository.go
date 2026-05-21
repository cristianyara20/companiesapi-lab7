package interfaces

import "companies-api/domain/entities"

// EmpleadoRepository define el contrato del repositorio de empleados
type EmpleadoRepository interface {
	GetAll() ([]entities.Empleado, error)
	GetById(id uint) (*entities.Empleado, error)
	Create(empleado *entities.Empleado) error
	Update(empleado *entities.Empleado) error
	Delete(id uint) error
	FindByCondition(condition string, args ...interface{}) ([]entities.Empleado, error)
	GetByCompaniaID(companiaID uint) ([]entities.Empleado, error)
}