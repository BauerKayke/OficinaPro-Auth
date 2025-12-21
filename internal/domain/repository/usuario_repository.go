// Package repository define interfaces de acesso a dados (Repository Pattern).
// Repositories abstraem a persistência seguindo princípio de Dependency Inversion.
package repository

import (
	"context"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// UsuarioRepository define contrato para acesso a dados de usuários.
// Interface permite múltiplas implementações (Postgres, MySQL, Mock, etc).
type UsuarioRepository interface {
	// FindByEmail busca usuário por email único.
	// Retorna entity.ErrUsuarioNotFound se não encontrado.
	FindByEmail(ctx context.Context, email string) (*entity.Usuario, error)

	// FindByID busca usuário por ID.
	// Retorna entity.ErrUsuarioNotFound se não encontrado.
	FindByID(ctx context.Context, id int64) (*entity.Usuario, error)

	// Create cria novo usuário no banco.
	// Retorna entity.ErrEmailAlreadyExists se email já existe.
	Create(ctx context.Context, usuario *entity.Usuario) error

	// Update atualiza usuário existente.
	// Retorna entity.ErrUsuarioNotFound se não encontrado.
	Update(ctx context.Context, usuario *entity.Usuario) error

	// Deactivate marca usuário como inativo (soft delete).
	Deactivate(ctx context.Context, id int64) error
}
