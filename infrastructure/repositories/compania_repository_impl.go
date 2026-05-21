package repositories

import (
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"

	"gorm.io/gorm"
)

// companiaRepositoryImpl es la implementación concreta del repositorio
// Vive en Infrastructure, no en Domain (Onion Architecture)
type companiaRepositoryImpl struct {
	db *gorm.DB
}

// NewCompaniaRepository crea el repositorio con el db inyectado
// El db puede ser la conexión normal o una transacción activa
func NewCompaniaRepository(db *gorm.DB) interfaces.CompaniaRepository {
	return &companiaRepositoryImpl{db: db}
}

func (r *companiaRepositoryImpl) GetAll() ([]entities.Compania, error) {
	var list []entities.Compania
	result := r.db.Find(&list)
	return list, result.Error
}

func (r *companiaRepositoryImpl) GetById(id uint) (*entities.Compania, error) {
	var c entities.Compania
	if err := r.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *companiaRepositoryImpl) Create(c *entities.Compania) error {
	// IMPORTANTE: NO llama a Commit aquí
	// El Commit es responsabilidad exclusiva del Unit of Work
	return r.db.Create(c).Error
}

func (r *companiaRepositoryImpl) Update(c *entities.Compania) error {
	return r.db.Save(c).Error
}

func (r *companiaRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&entities.Compania{}, id).Error
}

func (r *companiaRepositoryImpl) FindByCondition(cond string, args ...interface{}) ([]entities.Compania, error) {
	var list []entities.Compania
	result := r.db.Where(cond, args...).Find(&list)
	return list, result.Error
}

// GetWithEmpleados carga la compañía con sus empleados
// Equivale a .Include(c => c.Empleados) en EF Core
func (r *companiaRepositoryImpl) GetWithEmpleados(id uint) (*entities.Compania, error) {
	var c entities.Compania
	if err := r.db.Preload("Empleados").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}
