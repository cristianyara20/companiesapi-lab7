package dtos

// CreateCompaniaDTO datos requeridos para crear una compañía
type CreateCompaniaDTO struct {
	Nombre    string `json:"nombre"    validate:"required,min=3,max=100"`
	Direccion string `json:"direccion" validate:"omitempty"`
	Telefono  string `json:"telefono"  validate:"required,numeric,min=7,max=15"`
}

// UpdateCompaniaDTO campos opcionales para actualizar una compañía
type UpdateCompaniaDTO struct {
	Nombre    *string `json:"nombre"    validate:"omitempty,min=3,max=100"`
	Direccion *string `json:"direccion" validate:"omitempty"`
	Telefono  *string `json:"telefono"  validate:"omitempty,numeric,min=7,max=15"`
}

// CreateCompaniaConEmpleadosDTO para el endpoint transaccional
type CreateCompaniaConEmpleadosDTO struct {
	Nombre    string              `json:"nombre"    validate:"required,min=3,max=100"`
	Direccion string              `json:"direccion" validate:"required"`
	Telefono  string              `json:"telefono"  validate:"required,numeric,min=7,max=15"`
	Empleados []CreateEmpleadoDTO `json:"empleados" validate:"required,min=1,dive"`
}
