package dtos

// CreateEmpleadoDTO datos requeridos para crear un empleado
// binding:"email" → valida formato de correo
// binding:"gt=0"  → salario debe ser mayor a 0
type CreateEmpleadoDTO struct {
	Nombre     string  `json:"nombre"      binding:"required"`
	Apellido   string  `json:"apellido"    binding:"required"`
	Correo     string  `json:"correo"      binding:"required,email"`
	Cargo      string  `json:"cargo"       binding:"required"`
	Salario    float64 `json:"salario"     binding:"required,gt=0"`
	CompaniaID uint    `json:"compania_id" binding:"required"`
}

// UpdateEmpleadoDTO todos los campos son opcionales
type UpdateEmpleadoDTO struct {
	Nombre     string  `json:"nombre"`
	Apellido   string  `json:"apellido"`
	Correo     string  `json:"correo"`
	Cargo      string  `json:"cargo"`
	Salario    float64 `json:"salario"`
	CompaniaID uint    `json:"compania_id"`
}
