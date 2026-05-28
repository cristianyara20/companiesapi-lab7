package controllers

import (
	"companies-api/application/dtos"
	"companies-api/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthController struct {
	service *services.AuthService
	logger  *zap.Logger
}

func NewAuthController(s *services.AuthService, l *zap.Logger) *AuthController {
	return &AuthController{service: s, logger: l}
}

// Registro → POST /api/auth/registro
func (ctrl *AuthController) Registro(ctx *gin.Context) {
	var dto dtos.RegisterDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.service.Register(ctx.Request.Context(), dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "✅ Usuario registrado exitosamente",
		"usuario": gin.H{
			"id":             user.ID,
			"nombre":         user.Nombre,
			"correo":         user.Correo,
			"rol":            user.Rol,
			"compania_id":    user.CompaniaID,
			"fecha_creacion": user.FechaCreacion,
		},
	})
}

// Login → POST /api/auth/login
func (ctrl *AuthController) Login(ctx *gin.Context) {
	var dto dtos.LoginDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := ctrl.service.Login(ctx.Request.Context(), dto)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"usuario": gin.H{
			"id":     user.ID,
			"nombre": user.Nombre,
			"correo": user.Correo,
			"rol":    user.Rol,
		},
	})
}

// Perfil → GET /api/auth/perfil
func (ctrl *AuthController) Perfil(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userID := userIDVal.(uint)
	user, err := ctrl.service.GetPerfil(ctx.Request.Context(), userID)
	if err != nil {
		handleAPIError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":             user.ID,
		"nombre":         user.Nombre,
		"correo":         user.Correo,
		"rol":            user.Rol,
		"compania_id":    user.CompaniaID,
		"fecha_creacion": user.FechaCreacion,
	})
}
