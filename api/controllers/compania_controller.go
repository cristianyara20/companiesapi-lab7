package controllers

import (
	"companies-api/application/dtos"
	"companies-api/application/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CompaniaController maneja las peticiones HTTP de compañías
// Equivale al [ApiController] de ASP.NET Core
// REGLA: el controller NO accede al ORM ni a los repositorios directamente
type CompaniaController struct {
	service *services.CompaniaService
	logger  *zap.Logger
}

func NewCompaniaController(s *services.CompaniaService, l *zap.Logger) *CompaniaController {
	return &CompaniaController{service: s, logger: l}
}

// GetAll → GET /api/companias
func (c *CompaniaController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data) // 200
}

// GetById → GET /api/companias/:id
func (c *CompaniaController) GetById(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	data, err := c.service.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 404
		return
	}
	ctx.JSON(http.StatusOK, data) // 200
}

// Create → POST /api/companias
func (c *CompaniaController) Create(ctx *gin.Context) {
	var dto dtos.CreateCompaniaDTO
	// ShouldBindJSON valida el body y los campos requeridos
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // 400
		return
	}
	data, err := c.service.Create(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // 500
		return
	}
	ctx.JSON(http.StatusCreated, data) // 201
}

// Update → PUT /api/companias/:id
func (c *CompaniaController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	var dto dtos.UpdateCompaniaDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.Update(id, dto)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data) // 200
}

// Delete → DELETE /api/companias/:id
func (c *CompaniaController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil) // 204
}

// GetEmpleados → GET /api/companias/:id/empleados
func (c *CompaniaController) GetEmpleados(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	data, err := c.service.GetEmpleados(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data.Empleados) // 200
}

// CreateConEmpleados → POST /api/companias/con-empleados
// Endpoint transaccional obligatorio
func (c *CompaniaController) CreateConEmpleados(ctx *gin.Context) {
	var dto dtos.CreateCompaniaConEmpleadosDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.CreateConEmpleados(dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "✅ Compañía y empleados creados en una sola transacción",
		"compania": data,
	})
}

// parseID extrae el parámetro :id de la URL y lo convierte a uint
// Equivale a [FromRoute] int id en ASP.NET Core
func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero positivo"})
	}
	return uint(id), err
}
