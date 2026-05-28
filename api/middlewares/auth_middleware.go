package middlewares

import (
	"companies-api/application/services"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware extrae y valida el JWT token de la cabecera Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "cabecera Authorization ausente"})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "formato de token inválido, debe ser Bearer <token>"})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "super_secreto_desarrollo_12345"
		}

		claims := &services.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma inesperado")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token ausente, inválido o expirado"})
			ctx.Abort()
			return
		}

		// Establecer la información del usuario en el contexto de Gin
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_correo", claims.Correo)
		ctx.Set("user_rol", claims.Rol)
		if claims.CompaniaID != nil {
			ctx.Set("user_compania_id", *claims.CompaniaID)
		} else {
			ctx.Set("user_compania_id", nil)
		}

		ctx.Next()
	}
}

// RequireRole verifica que el usuario autenticado posea uno de los roles permitidos
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolVal, exists := ctx.Get("user_rol")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			ctx.Abort()
			return
		}

		userRol := rolVal.(string)
		isAllowed := false
		for _, role := range allowedRoles {
			if userRol == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado: rol insuficiente"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// EsPropietarioDeCompania implementa la política de propiedad (ownership)
// Un USUARIO solo puede actualizar/eliminar empleados de su misma compañía.
// El ADMIN está exento de esta restricción.
func EsPropietarioDeCompania(empleadoService *services.EmpleadoService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rolVal, exists := ctx.Get("user_rol")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			ctx.Abort()
			return
		}

		rol := rolVal.(string)
		if rol == "ADMIN" {
			// El administrador tiene acceso global, omitir validación de propiedad
			ctx.Next()
			return
		}

		// Obtener CompaniaID del token del usuario
		compVal, exists := ctx.Get("user_compania_id")
		if !exists || compVal == nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado: el usuario no pertenece a ninguna compañía"})
			ctx.Abort()
			return
		}
		userCompanyID := compVal.(uint)

		// Obtener ID del empleado desde el parámetro de la ruta (:id)
		empIDStr := ctx.Param("id")
		if empIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID del empleado ausente"})
			ctx.Abort()
			return
		}

		empID, err := strconv.ParseUint(empIDStr, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID del empleado inválido"})
			ctx.Abort()
			return
		}

		// Consultar el empleado para verificar la propiedad
		emp, err := empleadoService.GetById(ctx.Request.Context(), uint(empID))
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "empleado no encontrado"})
			ctx.Abort()
			return
		}

		// Validar si el empleado pertenece a la misma compañía del usuario
		if emp.CompaniaID != userCompanyID {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "acceso denegado: no es propietario de esta compañía"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
