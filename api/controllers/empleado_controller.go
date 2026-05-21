package controllers

import (
	"companies-api/application/dtos"
	"companies-api/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EmpleadoController struct {
	service *services.EmpleadoService
	logger  *zap.Logger
}

func NewEmpleadoController(s *services.EmpleadoService, l *zap.Logger) *EmpleadoController {
	return &EmpleadoController{service: s, logger: l}
}

// GetAll → GET /api/empleados
func (c *EmpleadoController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetById → GET /api/empleados/:id
func (c *EmpleadoController) GetById(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	data, err := c.service.GetById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Create → POST /api/empleados
func (c *EmpleadoController) Create(ctx *gin.Context) {
	var dto dtos.CreateEmpleadoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.Create(dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, data)
}

// Update → PUT /api/empleados/:id
func (c *EmpleadoController) Update(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	var dto dtos.UpdateEmpleadoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.Update(id, dto)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Delete → DELETE /api/empleados/:id
func (c *EmpleadoController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	if err := c.service.Delete(id); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
