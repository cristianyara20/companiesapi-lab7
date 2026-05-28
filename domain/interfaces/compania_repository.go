package interfaces

import (
	"companies-api/domain/entities"
	"context"
)

// CompaniaRepository define el contrato del repositorio de compañías
type CompaniaRepository interface {
	GetAll(ctx context.Context) ([]entities.Compania, error)
	GetById(ctx context.Context, id uint) (*entities.Compania, error)
	Create(ctx context.Context, compania *entities.Compania) error
	Update(ctx context.Context, compania *entities.Compania) error
	Delete(ctx context.Context, id uint) error
	FindByCondition(ctx context.Context, condition string, args ...interface{}) ([]entities.Compania, error)
	GetWithEmpleados(ctx context.Context, id uint) (*entities.Compania, error)
}