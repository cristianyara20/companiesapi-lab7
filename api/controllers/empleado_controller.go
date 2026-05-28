package controllers

import (
	"companies-api/application/dtos"
	"companies-api/application/services"
	"net/http"
	"strconv"

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

// GetAll → GET /api/empleados?pagina=&tamano=&orden=&dir=&buscar=
// Listado paginado, filtrado y ordenado
func (c *EmpleadoController) GetAll(ctx *gin.Context) {
	paginaStr := ctx.DefaultQuery("pagina", "1")
	tamanoStr := ctx.DefaultQuery("tamano", "10")
	orden := ctx.DefaultQuery("orden", "id")
	dir := ctx.DefaultQuery("dir", "asc")
	buscar := ctx.Query("buscar")

	pagina, errPag := strconv.Atoi(paginaStr)
	tamano, errTam := strconv.Atoi(tamanoStr)

	if errPag != nil || errTam != nil || pagina < 1 || tamano < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Parámetros de paginación inválidos"})
		return
	}

	data, total, err := c.service.GetPaged(ctx.Request.Context(), pagina, tamano, orden, dir, buscar)
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

// GetById → GET /api/empleados/:id
func (c *EmpleadoController) GetById(ctx *gin.Context) {
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

// Create → POST /api/empleados
func (c *EmpleadoController) Create(ctx *gin.Context) {
	var dto dtos.CreateEmpleadoDTO
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

// Update → PUT /api/empleados/:id
// Reemplazo completo (se actualizan todos los campos que vengan)
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
	data, err := c.service.Update(ctx.Request.Context(), id, dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// Patch → PATCH /api/empleados/:id
// Actualización parcial
func (c *EmpleadoController) Patch(ctx *gin.Context) {
	id, err := parseID(ctx)
	if err != nil {
		return
	}
	var dto dtos.UpdateEmpleadoDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := c.service.PatchPartial(ctx.Request.Context(), id, dto)
	if err != nil {
		handleAPIError(ctx, err)
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
	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		handleAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// CreateRange → POST /api/empleados/lote
// Creación masiva (Bulk)
func (c *EmpleadoController) CreateRange(ctx *gin.Context) {
	var list []dtos.CreateEmpleadoDTO
	if err := ctx.ShouldBindJSON(&list); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := c.service.CreateRange(ctx.Request.Context(), list)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":   "✅ Empleados creados exitosamente en lote",
		"empleados": data,
	})
}

// DeleteRange → DELETE /api/empleados/lote
// Eliminación múltiple
func (c *EmpleadoController) DeleteRange(ctx *gin.Context) {
	var req struct {
		Ids []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.DeleteRange(ctx.Request.Context(), req.Ids)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
