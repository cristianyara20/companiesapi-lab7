package repositories

import (
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"context"
	"strings"

	"gorm.io/gorm"
)

type empleadoRepositoryImpl struct {
	db *gorm.DB
}

// NewEmpleadoRepository crea el repositorio de empleados con el db inyectado
func NewEmpleadoRepository(db *gorm.DB) interfaces.EmpleadoRepository {
	return &empleadoRepositoryImpl{db: db}
}

func (r *empleadoRepositoryImpl) GetAll(ctx context.Context) ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.WithContext(ctx).Find(&list).Error
}

func (r *empleadoRepositoryImpl) GetById(ctx context.Context, id uint) (*entities.Empleado, error) {
	var e entities.Empleado
	if err := r.db.WithContext(ctx).First(&e, id).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *empleadoRepositoryImpl) Create(ctx context.Context, e *entities.Empleado) error {
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *empleadoRepositoryImpl) Update(ctx context.Context, e *entities.Empleado) error {
	return r.db.WithContext(ctx).Save(e).Error
}

func (r *empleadoRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Empleado{}, id).Error
}

func (r *empleadoRepositoryImpl) FindByCondition(ctx context.Context, cond string, args ...interface{}) ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.WithContext(ctx).Where(cond, args...).Find(&list).Error
}

func (r *empleadoRepositoryImpl) GetByCompaniaID(ctx context.Context, cid uint) ([]entities.Empleado, error) {
	var list []entities.Empleado
	return list, r.db.WithContext(ctx).Where("compania_id = ?", cid).Find(&list).Error
}

// CreateRange realiza inserción en lote (bulk insert)
func (r *empleadoRepositoryImpl) CreateRange(ctx context.Context, empleados []entities.Empleado) error {
	return r.db.WithContext(ctx).Create(&empleados).Error
}

// DeleteRange realiza eliminación múltiple por IDs
func (r *empleadoRepositoryImpl) DeleteRange(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Empleado{}, ids).Error
}

// GetPaged devuelve una página filtrada y ordenada de empleados junto con el total
func (r *empleadoRepositoryImpl) GetPaged(ctx context.Context, pagina int, tamano int, orden string, dir string, buscar string) ([]entities.Empleado, int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Empleado{})

	// 1. Filtrado / Búsqueda
	if buscar != "" {
		buscarLower := "%" + strings.ToLower(buscar) + "%"
		query = query.Where("LOWER(nombre) LIKE ? OR LOWER(apellido) LIKE ? OR LOWER(correo) LIKE ?", buscarLower, buscarLower, buscarLower)
	}

	// 2. Conteo de Total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 3. Ordenamiento Dinámico Sanitizado (para prevenir SQL Injection)
	allowedColumns := map[string]bool{
		"id":          true,
		"nombre":      true,
		"apellido":    true,
		"correo":      true,
		"cargo":       true,
		"salario":     true,
		"compania_id": true,
	}

	ordenCol := "id"
	if allowedColumns[strings.ToLower(orden)] {
		ordenCol = strings.ToLower(orden)
	}

	ordenDir := "asc"
	if strings.ToLower(dir) == "desc" {
		ordenDir = "desc"
	}

	orderBy := ordenCol + " " + ordenDir

	// 4. Paginación y Consulta de Datos
	var list []entities.Empleado
	offset := (pagina - 1) * tamano
	err := query.Order(orderBy).Offset(offset).Limit(tamano).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// GetPagedByCompaniaID obtiene empleados paginados para una compañía específica
func (r *empleadoRepositoryImpl) GetPagedByCompaniaID(ctx context.Context, companiaID uint, pagina int, tamano int) ([]entities.Empleado, int64, error) {
	query := r.db.WithContext(ctx).Model(&entities.Empleado{}).Where("compania_id = ?", companiaID)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []entities.Empleado
	offset := (pagina - 1) * tamano
	err := query.Offset(offset).Limit(tamano).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// PatchPartial actualiza parcialmente un empleado utilizando un mapa de cambios
func (r *empleadoRepositoryImpl) PatchPartial(ctx context.Context, id uint, cambios map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&entities.Empleado{}).Where("id = ?", id).Updates(cambios).Error
}
