package entities

// Empleado representa la tabla "empleados" en PostgreSQL
type Empleado struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"      json:"id"`
	Nombre     string    `gorm:"size:100;not null"             json:"nombre"`
	Apellido   string    `gorm:"size:100;not null"             json:"apellido"`
	Correo     string    `gorm:"size:150;not null;uniqueIndex" json:"correo"`
	Cargo      string    `gorm:"size:100;not null"             json:"cargo"`
	Salario    float64   `gorm:"not null"                      json:"salario"`
	// Llave foránea hacia companias.id
	CompaniaID uint      `gorm:"not null;index"                json:"compania_id"`
	// Relación navegación (equivale a virtual Compania en EF Core)
	Compania   *Compania `gorm:"foreignKey:CompaniaID"         json:"compania,omitempty"`
}