package interfaces

import "companies-api/domain/entities"

// CompaniaRepository define el contrato del repositorio
// Equivale a la interfaz ICompaniaRepository en ASP.NET Core
// La implementación concreta estará en infrastructure/repositories
type CompaniaRepository interface {
	GetAll() ([]entities.Compania, error)
	GetById(id uint) (*entities.Compania, error)
	Create(compania *entities.Compania) error
	Update(compania *entities.Compania) error
	Delete(id uint) error
	FindByCondition(condition string, args ...interface{}) ([]entities.Compania, error)
	GetWithEmpleados(id uint) (*entities.Compania, error)
}