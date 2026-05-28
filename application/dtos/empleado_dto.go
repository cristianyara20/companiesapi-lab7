package dtos

// CreateEmpleadoDTO datos requeridos para crear un empleado
type CreateEmpleadoDTO struct {
	Nombre     string  `json:"nombre"      validate:"required,min=1"`
	Apellido   string  `json:"apellido"    validate:"required,min=1"`
	Correo     string  `json:"correo"      validate:"required,email"`
	Cargo      string  `json:"cargo"       validate:"required"`
	Salario    float64 `json:"salario"     validate:"required,gt=0"`
	CompaniaID uint    `json:"compania_id" validate:"required"`
}

// UpdateEmpleadoDTO campos opcionales para actualizar un empleado
type UpdateEmpleadoDTO struct {
	Nombre     *string  `json:"nombre"      validate:"omitempty,min=1"`
	Apellido   *string  `json:"apellido"    validate:"omitempty,min=1"`
	Correo     *string  `json:"correo"      validate:"omitempty,email"`
	Cargo      *string  `json:"cargo"       validate:"omitempty"`
	Salario    *float64 `json:"salario"     validate:"omitempty,gt=0"`
	CompaniaID *uint    `json:"compania_id" validate:"omitempty"`
}
