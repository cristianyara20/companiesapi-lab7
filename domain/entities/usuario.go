package entities

import "time"

// Usuario representa la tabla "usuarios" en PostgreSQL
// Almacena credenciales con contraseña hasheada (bcrypt) y rol de autorización
type Usuario struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"      json:"id"`
	Nombre         string    `gorm:"size:100;not null"             json:"nombre"`
	Correo         string    `gorm:"size:150;not null;uniqueIndex" json:"correo"`
	ContrasenaHash string    `gorm:"size:255;not null"             json:"-"`
	Rol            string    `gorm:"size:20;not null;default:'USUARIO'" json:"rol"`
	CompaniaID     *uint     `gorm:"index"                         json:"compania_id,omitempty"`
	FechaCreacion  time.Time `gorm:"autoCreateTime"                json:"fecha_creacion"`
	// Relación navegación (opcional)
	Compania *Compania `gorm:"foreignKey:CompaniaID" json:"compania,omitempty"`
}
