package unitofwork

import (
	"companies-api/domain/interfaces"
	"companies-api/infrastructure/repositories"
	"errors"

	"gorm.io/gorm"
)

// UnitOfWorkImpl coordina todos los repositorios bajo una misma transacción
type UnitOfWorkImpl struct {
	db        *gorm.DB
	tx        *gorm.DB
	companias interfaces.CompaniaRepository
	empleados interfaces.EmpleadoRepository
	usuarios  interfaces.UsuarioRepository
}

// NewUnitOfWork crea una nueva instancia del Unit of Work
func NewUnitOfWork(db *gorm.DB) interfaces.UnitOfWork {
	return &UnitOfWorkImpl{
		db: db,
		tx: db, // por defecto usa la conexión sin transacción
	}
}

// Companias devuelve el repositorio de compañías usando la sesión activa
func (u *UnitOfWorkImpl) Companias() interfaces.CompaniaRepository {
	if u.companias == nil {
		u.companias = repositories.NewCompaniaRepository(u.tx)
	}
	return u.companias
}

// Empleados devuelve el repositorio de empleados usando la sesión activa
func (u *UnitOfWorkImpl) Empleados() interfaces.EmpleadoRepository {
	if u.empleados == nil {
		u.empleados = repositories.NewEmpleadoRepository(u.tx)
	}
	return u.empleados
}

// Usuarios devuelve el repositorio de usuarios usando la sesión activa
func (u *UnitOfWorkImpl) Usuarios() interfaces.UsuarioRepository {
	if u.usuarios == nil {
		u.usuarios = repositories.NewUsuarioRepository(u.tx)
	}
	return u.usuarios
}

// BeginTransaction inicia una transacción real en PostgreSQL
func (u *UnitOfWorkImpl) BeginTransaction() error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	u.tx = tx         // reemplaza la sesión activa por la transacción
	u.companias = nil // reinicia repos para que usen el nuevo tx
	u.empleados = nil
	u.usuarios = nil
	return nil
}

// Commit confirma todos los cambios de la transacción en PostgreSQL
func (u *UnitOfWorkImpl) Commit() error {
	if u.tx == u.db {
		return errors.New("no hay transacción activa para confirmar")
	}
	err := u.tx.Commit().Error
	u.tx = u.db // restaura la sesión a la conexión original
	u.companias = nil
	u.empleados = nil
	u.usuarios = nil
	return err
}

// Rollback revierte TODOS los cambios hechos desde BeginTransaction
func (u *UnitOfWorkImpl) Rollback() error {
	if u.tx == u.db {
		return nil // no hay transacción activa, nada que revertir
	}
	err := u.tx.Rollback().Error
	u.tx = u.db
	u.companias = nil
	u.empleados = nil
	u.usuarios = nil
	return err
}
