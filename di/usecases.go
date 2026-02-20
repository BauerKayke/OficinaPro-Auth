package di

import (
	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// NewAuthenticateUseCase cria e configura o caso de uso de autenticação.
// Injeta todas as dependências necessárias.
func NewAuthenticateUseCase(
	cfg *config.Config,
	usuarioRepo repository.UsuarioRepository,
	jwtService service.JWTService,
	telemetryService service.TelemetryService,
) *usecase.AuthenticateUseCase {
	return usecase.NewAuthenticateUseCase(
		usuarioRepo,
		jwtService,
		telemetryService,
		cfg.JWT.Expiration,
	)
}
