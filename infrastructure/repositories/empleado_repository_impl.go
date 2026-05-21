package repositories

import (
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"

	"gorm.io/gorm"
)

type empleadoRepositoryImpl struct {
	db *gorm.DB
}

func NewEmpleadoRepository(db *gorm.DB) interfaces.EmpleadoRepository {
	return &empleadoRepositoryImpl{db: db}
}

func (r *empleadoRepositoryImpl) GetAll() ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.Find(&list).Error
}

func (r *empleadoRepositoryImpl) GetById(id uint) (*entities.Empleado, error) {
	var e entities.Empleado
	if err := r.db.First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *empleadoRepositoryImpl) Create(e *entities.Empleado) error {
	return r.db.Create(e).Error
}

func (r *empleadoRepositoryImpl) Update(e *entities.Empleado) error {
	return r.db.Save(e).Error
}

func (r *empleadoRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Empleado{}, id).Error
}

func (r *empleadoRepositoryImpl) FindByCondition(cond string, args ...interface{}) ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.Where(cond, args...).Find(&list).Error
}

func (r *empleadoRepositoryImpl) GetByCompaniaID(cid uint) ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.Where("compania_id = ?", cid).Find(&list).Error
}
