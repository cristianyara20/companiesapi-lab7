package controllers

import (
	"companies-api/application/validation"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleAPIError centraliza el formateo y envío de respuestas de error a nivel de HTTP
func handleAPIError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	// 1. Error de validación estructural o de negocio (Capa de Aplicación) -> 422 Unprocessable Entity
	var valErr *validation.ValidationError
	if errors.As(err, &valErr) {
		ctx.JSON(http.StatusUnprocessableEntity, valErr)
		return
	}

	// 2. Error de recurso no encontrado -> 404 Not Found
	errMsg := err.Error()
	if errMsg == "compañía no encontrada" || errMsg == "empleado no encontrado" || errMsg == "usuario no encontrado" || errMsg == "la compañía especificada no existe" {
		ctx.JSON(http.StatusNotFound, gin.H{"error": errMsg})
		return
	}

	// 3. Error de credenciales inválidas -> 401 Unauthorized
	if errMsg == "credenciales inválidas" || errMsg == "token inválido" || errMsg == "token ausente o expirado" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
		return
	}

	// 4. Error de autorización -> 403 Forbidden
	if errMsg == "acceso denegado: rol insuficiente" || errMsg == "acceso denegado: no es propietario de esta compañía" {
		ctx.JSON(http.StatusForbidden, gin.H{"error": errMsg})
		return
	}

	// 5. Error genérico del servidor -> 500 Internal Server Error
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
}
