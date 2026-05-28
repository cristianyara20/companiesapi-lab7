package repositories

import (
	"companies-api/domain/entities"
	"companies-api/domain/interfaces"
	"context"

	"gorm.io/gorm"
)

type usuarioRepositoryImpl struct {
	db *gorm.DB
}

// NewUsuarioRepository crea el repositorio de usuarios con el db inyectado
func NewUsuarioRepository(db *gorm.DB) interfaces.UsuarioRepository {
	return &usuarioRepositoryImpl{db: db}
}

func (r *usuarioRepositoryImpl) Create(ctx context.Context, u *entities.Usuario) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *usuarioRepositoryImpl) GetByCorreo(ctx context.Context, correo string) (*entities.Usuario, error) {
	var u entities.Usuario
	if err := r.db.WithContext(ctx).Where("correo = ?", correo).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *usuarioRepositoryImpl) GetById(ctx context.Context, id uint) (*entities.Usuario, error) {
	var u entities.Usuario
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
