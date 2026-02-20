package di

import (
	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/infrastructure/database"
)

// NewRepositories cria e configura repositórios da aplicação.
// Retorna ClienteRepository e UsuarioRepository configurados.
func NewRepositories(db *gorm.DB) (repository.ClienteRepository, repository.UsuarioRepository) {
	clienteRepo := database.NewGormClienteRepository(db)
	usuarioRepo := database.NewGormUsuarioRepository(db)
	return clienteRepo, usuarioRepo
}
