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
type CompaniaController struct {
	service *services.CompaniaService
	logger  *zap.Logger
}

func NewCompaniaController(s *services.CompaniaService, l *zap.Logger) *CompaniaController {
	return &CompaniaController{service: s, logger: l}
}

// GetAll → GET /api/companias
func (c *CompaniaController) GetAll(ctx *gin.Context) {
	data, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetById → GET /api/companias/:id
func (c *CompaniaController) GetById(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	data, err := c.service.GetById(ctx.Request.Context(), id)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Create → POST /api/companias
func (c *CompaniaController) Create(ctx *gin.Context) {
	var dto dtos.CreateCompaniaDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.Create(ctx.Request.Context(), dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, data)
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
	data, err := c.service.Update(ctx.Request.Context(), id, dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Delete → DELETE /api/companias/:id
func (c *CompaniaController) Delete(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// GetEmpleados → GET /api/companias/:id/empleados?pagina=&tamano=
func (c *CompaniaController) GetEmpleados(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}

	paginaStr := ctx.DefaultQuery("pagina", "1")
	tamanoStr := ctx.DefaultQuery("tamano", "10")

	pagina, errPag := strconv.Atoi(paginaStr)
	tamano, errTam := strconv.Atoi(tamanoStr)

	if errPag != nil || errTam != nil || pagina < 1 || tamano < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de paginación inválidos"})
		return
	}

	data, total, err := c.service.GetEmpleadosPaged(ctx.Request.Context(), id, pagina, tamano)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	totalPages := int((total + int64(tamano) - 1) / int64(tamano))
	if totalPages == 0 {
		totalPages = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"datos":        data,
		"pagina":       pagina,
		"tamano":       tamano,
		"total":        total,
		"totalPaginas": totalPages,
	})
}

// CreateConEmpleados → POST /api/companias/con-empleados
func (c *CompaniaController) CreateConEmpleados(ctx *gin.Context) {
	var dto dtos.CreateCompaniaConEmpleadosDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.CreateConEmpleados(ctx.Request.Context(), dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "✅ Compañía y empleados creados en una sola transacción",
		"compania": data,
	})
}

// parseID extrae el parámetro :id de la URL y lo convierte a uint
func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "El ID debe ser un número entero positivo"})
	}
	return uint(id), err
}
