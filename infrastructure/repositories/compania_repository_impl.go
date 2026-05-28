package repositories

import (
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"context"

	"gorm.io/gorm"
)

// companiaRepositoryImpl es la implementación concreta del repositorio de compañías
type companiaRepositoryImpl struct {
	db *gorm.DB
}

// NewCompaniaRepository crea el repositorio con el db inyectado
func NewCompaniaRepository(db *gorm.DB) interfaces.CompaniaRepository {
	return &companiaRepositoryImpl{db: db}
}

func (r *companiaRepositoryImpl) GetAll(ctx context.Context) ([]entities.Compania, error) {
	var list []entities.Compania
	result := r.db.WithContext(ctx).Find(&list)
	return list, result.Error
}

func (r *companiaRepositoryImpl) GetById(ctx context.Context, id uint) (*entities.Compania, error) {
	var c entities.Compania
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *companiaRepositoryImpl) Create(ctx context.Context, c *entities.Compania) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *companiaRepositoryImpl) Update(ctx context.Context, c *entities.Compania) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *companiaRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Compania{}, id).Error
}

func (r *companiaRepositoryImpl) FindByCondition(ctx context.Context, cond string, args ...interface{}) ([]entities.Compania, error) {
	var list []entities.Compania
	result := r.db.WithContext(ctx).Where(cond, args...).Find(&list)
	return list, result.Error
}

func (r *companiaRepositoryImpl) GetWithEmpleados(ctx context.Context, id uint) (*entities.Compania, error) {
	var c entities.Compania
	if err := r.db.WithContext(ctx).Preload("Empleados").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
