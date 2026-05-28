package interfaces

import (
	"companies-api/domain/entities"
	"context"
)

// EmpleadoRepository define el contrato del repositorio de empleados
type EmpleadoRepository interface {
	GetAll(ctx context.Context) ([]entities.Empleado, error)
	GetById(ctx context.Context, id uint) (*entities.Empleado, error)
	Create(ctx context.Context, empleado *entities.Empleado) error
	Update(ctx context.Context, empleado *entities.Empleado) error
	Delete(ctx context.Context, id uint) error
	FindByCondition(ctx context.Context, condition string, args ...interface{}) ([]entities.Empleado, error)
	GetByCompaniaID(ctx context.Context, companiaID uint) ([]entities.Empleado, error)
	
	// Operaciones de colección requeridas (Parte II)
	CreateRange(ctx context.Context, empleados []entities.Empleado) error
	DeleteRange(ctx context.Context, ids []uint) error
	GetPaged(ctx context.Context, pagina int, tamano int, orden string, dir string, buscar string) ([]entities.Empleado, int64, error)
	GetPagedByCompaniaID(ctx context.Context, companiaID uint, pagina int, tamano int) ([]entities.Empleado, int64, error)
	PatchPartial(ctx context.Context, id uint, cambios map[string]interface{}) error
}