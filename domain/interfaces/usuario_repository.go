package interfaces

import (
	"companies-api/domain/entities"
	"context"
)

// UsuarioRepository define el contrato del repositorio de usuarios
type UsuarioRepository interface {
	Create(ctx context.Context, usuario *entities.Usuario) error
	GetByCorreo(ctx context.Context, correo string) (*entities.Usuario, error)
	GetById(ctx context.Context, id uint) (*entities.Usuario, error)
}
