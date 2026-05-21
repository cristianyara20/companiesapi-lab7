package unitofwork

import (
	"companies-api/domain/interfaces"
	"companies-api/infrastructure/repositories"
	"errors"

	"gorm.io/gorm"
)

// UnitOfWorkImpl coordina todos los repositorios bajo una misma transacción
//
// ¿Cómo funciona?
//   - u.db  = conexión original a PostgreSQL (nunca cambia)
//   - u.tx  = sesión activa (empieza como u.db, cambia a transacción con BeginTransaction)
//   - Los repositorios siempre usan u.tx, así que automáticamente participan
//     en la transacción activa sin saberlo
//
// Equivalente conceptual en EF Core:
//   - db.Begin()    → dbContext.Database.BeginTransactionAsync()
//   - tx.Commit()   → transaction.CommitAsync()
//   - tx.Rollback() → transaction.RollbackAsync()
type UnitOfWorkImpl struct {
	db        *gorm.DB
	tx        *gorm.DB
	companias interfaces.CompaniaRepository
	empleados interfaces.EmpleadoRepository
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

// BeginTransaction inicia una transacción real en PostgreSQL
// Todos los repositorios creados después de esta llamada usarán la misma tx
func (u *UnitOfWorkImpl) BeginTransaction() error {
	tx := u.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	u.tx = tx         // reemplaza la sesión activa por la transacción
	u.companias = nil // reinicia repos para que usen el nuevo tx
	u.empleados = nil
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
	return err
}
