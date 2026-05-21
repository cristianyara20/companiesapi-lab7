package entities

import "time"

// Compania representa la tabla "companias" en PostgreSQL
// Los tags `gorm:"..."` configuran columnas, igual que Fluent API en EF Core
// Los tags `json:"..."` configuran cómo se serializa en JSON
type Compania struct {
	ID            uint       `gorm:"primaryKey;autoIncrement"  json:"id"`
	Nombre        string     `gorm:"size:100;not null"         json:"nombre"`
	Direccion     string     `gorm:"size:200;not null"         json:"direccion"`
	Telefono      string     `gorm:"size:20;not null"          json:"telefono"`
	FechaCreacion time.Time  `gorm:"autoCreateTime"            json:"fecha_creacion"`
	// Relación 1:N — una compañía tiene muchos empleados
	// Equivale a ICollection<Empleado> en EF Core
	Empleados []Empleado `gorm:"foreignKey:CompaniaID" json:"empleados,omitempty"`
}