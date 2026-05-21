package dtos

// CreateCompaniaDTO datos requeridos para crear una compañía
// binding:"required" → Gin valida que el campo no esté vacío
type CreateCompaniaDTO struct {
	Nombre    string `json:"nombre"    binding:"required"`
	Direccion string `json:"direccion" binding:"required"`
	Telefono  string `json:"telefono"  binding:"required"`
}

// UpdateCompaniaDTO todos los campos son opcionales en un PUT
type UpdateCompaniaDTO struct {
	Nombre    string `json:"nombre"`
	Direccion string `json:"direccion"`
	Telefono  string `json:"telefono"`
}

// CreateCompaniaConEmpleadosDTO para el endpoint transaccional
// Crea compañía + empleados en una sola transacción
type CreateCompaniaConEmpleadosDTO struct {
	Nombre    string              `json:"nombre"    binding:"required"`
	Direccion string              `json:"direccion" binding:"required"`
	Telefono  string              `json:"telefono"  binding:"required"`
	Empleados []CreateEmpleadoDTO `json:"empleados" binding:"required,min=1"`
}
