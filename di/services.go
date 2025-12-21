package di

import (
	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/infrastructure/jwt"
	"github.com/oficinapro/auth-service/internal/infrastructure/validator"
)

// NewServices cria e configura serviços de domínio.
// Retorna JWT service e validator service configurados.
func NewServices(cfg *config.Config) (service.JWTService, service.ValidatorService) {
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)
	validatorService := validator.NewCPFValidator()
	return jwtService, validatorService
}
