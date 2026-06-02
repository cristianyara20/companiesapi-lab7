package services

import (
	"companies-api/application/dtos"
	"companies-api/application/validation"
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// JWTClaims define los claims del token JWT.
// Se expone aquí para que el middleware pueda importarlo sin ciclos.
type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CompanyID *uint  `json:"company_id,omitempty"`
	City      string `json:"city"`
	jwt.RegisteredClaims
}

// AuthService contiene la lógica de autenticación y autorización
type AuthService struct {
	uow    interfaces.UnitOfWork
	logger *zap.Logger
}

func NewAuthService(uow interfaces.UnitOfWork, logger *zap.Logger) *AuthService {
	return &AuthService{uow: uow, logger: logger}
}

// Register registra un nuevo usuario con contraseña hasheada
func (s *AuthService) Register(ctx context.Context, dto dtos.RegisterDTO) (*entities.Usuario, error) {
	s.logger.Info("Registrando nuevo usuario", zap.String("correo", dto.Correo))

	// 1. Validación estructural del DTO (Capa de Aplicación)
	if err := validation.ValidateStruct(dto); err != nil {
		return nil, err
	}

	// 2. Validación de negocio: correo único
	existing, err := s.uow.Usuarios().GetByCorreo(ctx, dto.Correo)
	if err == nil && existing != nil {
		return nil, &validation.ValidationError{
			Mensaje: "Error de validación de negocio",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "correo", Detalle: "El correo ya está registrado por otro usuario"},
			},
		}
	}

	// 3. Si el rol es USUARIO, verificar que la compañía existe
	if dto.Rol == "USUARIO" && dto.CompaniaID != nil {
		if _, err := s.uow.Companias().GetById(ctx, *dto.CompaniaID); err != nil {
			return nil, &validation.ValidationError{
				Mensaje: "Error de validación de negocio",
				Errores: []validation.ValidationErrorDetail{
					{Campo: "compania_id", Detalle: "La compañía especificada no existe"},
				},
			}
		}
	}

	// 4. Hash de la contraseña con bcrypt (contraseña nunca en texto plano)
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Error generando hash de contraseña", zap.Error(err))
		return nil, errors.New("error al procesar la contraseña")
	}

	u := &entities.Usuario{
		Nombre:         dto.Nombre,
		Correo:         dto.Correo,
		ContrasenaHash: string(hash),
		Rol:            dto.Rol,
		CompaniaID:     dto.CompaniaID,
	}

	if err := s.uow.Usuarios().Create(ctx, u); err != nil {
		s.logger.Error("Error al registrar usuario", zap.Error(err))
		return nil, err
	}

	s.logger.Info("✅ Usuario registrado", zap.Uint("id", u.ID), zap.String("rol", u.Rol))
	return u, nil
}

// Login valida credenciales y devuelve un token JWT firmado
func (s *AuthService) Login(ctx context.Context, dto dtos.LoginDTO) (string, *entities.Usuario, error) {
	// 1. Validación estructural
	if err := validation.ValidateStruct(dto); err != nil {
		return "", nil, err
	}

	// 2. Buscar usuario por correo
	u, err := s.uow.Usuarios().GetByCorreo(ctx, dto.Correo)
	if err != nil {
		return "", nil, &validation.ValidationError{
			Mensaje: "Credenciales inválidas",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "correo", Detalle: "El correo o la contraseña son incorrectos"},
			},
		}
	}

	// 3. Comparar contraseña con bcrypt (nunca texto plano)
	if err := bcrypt.CompareHashAndPassword([]byte(u.ContrasenaHash), []byte(dto.Contrasena)); err != nil {
		return "", nil, &validation.ValidationError{
			Mensaje: "Credenciales inválidas",
			Errores: []validation.ValidationErrorDetail{
				{Campo: "contrasena", Detalle: "El correo o la contraseña son incorrectos"},
			},
		}
	}

	// 4. Generar JWT — la clave viene del entorno, NUNCA del código fuente
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback solo para desarrollo (no exponer en producción)
		secret = "super_secreto_desarrollo_12345"
		s.logger.Warn("⚠️  JWT_SECRET no configurado — usando valor de desarrollo")
	}

	city := "Global"
	if u.CompaniaID != nil {
		switch *u.CompaniaID {
		case 1:
			city = "Bogotá"
		case 2:
			city = "Medellín"
		case 3:
			city = "Cali"
		}
	}

	claims := &JWTClaims{
		UserID:    u.ID,
		Email:     u.Correo,
		Role:      u.Rol,
		CompanyID: u.CompaniaID,
		City:      city,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "companies-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		s.logger.Error("Error firmando token JWT", zap.Error(err))
		return "", nil, errors.New("error al generar el token de acceso")
	}

	s.logger.Info("✅ Login exitoso", zap.String("correo", u.Correo), zap.String("rol", u.Rol))
	return tokenString, u, nil
}

// GetPerfil devuelve el perfil del usuario autenticado por su ID
func (s *AuthService) GetPerfil(ctx context.Context, userID uint) (*entities.Usuario, error) {
	u, err := s.uow.Usuarios().GetById(ctx, userID)
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}
