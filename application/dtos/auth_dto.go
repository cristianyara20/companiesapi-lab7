package dtos

// RegisterDTO datos requeridos para registrar un nuevo usuario
type RegisterDTO struct {
	Nombre     string `json:"nombre"      validate:"required,min=3,max=100"`
	Correo     string `json:"correo"      validate:"required,email"`
	Contrasena string `json:"contrasena"  validate:"required,min=6"`
	Rol        string `json:"rol"         validate:"required,oneof=ADMIN USUARIO"`
	CompaniaID *uint  `json:"compania_id" validate:"omitempty"`
}

// LoginDTO datos requeridos para iniciar sesión
type LoginDTO struct {
	Correo     string `json:"correo"     validate:"required,email"`
	Contrasena string `json:"contrasena" validate:"required"`
}
