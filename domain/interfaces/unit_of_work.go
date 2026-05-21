package interfaces

// UnitOfWork coordina repositorios y maneja transacciones
// Equivale al DbContext de EF Core pero implementado manualmente en Go
// Flujo: BeginTransaction → operaciones → Commit (o Rollback si hay error)
type UnitOfWork interface {
	Companias() CompaniaRepository  // acceso al repositorio de compañías
	Empleados() EmpleadoRepository  // acceso al repositorio de empleados
	BeginTransaction() error         // inicia una transacción en PostgreSQL
	Commit() error                   // confirma todos los cambios
	Rollback() error                 // revierte todos los cambios
}