package di

import (
	"context"

	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// Container holds all application dependencies
// Simple: apenas armazena, não cria (SRP)
type Container struct {
	clienteRepo      repository.ClienteRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	authenticateUC   *usecase.AuthenticateUseCase
	closeFunc        func() error
}

// Getters (encapsulamento)

func (c *Container) ClienteRepository() repository.ClienteRepository {
	return c.clienteRepo
}

func (c *Container) JWTService() service.JWTService {
	return c.jwtService
}

func (c *Container) ValidatorService() service.ValidatorService {
	return c.validatorService
}

func (c *Container) AuthenticateUseCase() *usecase.AuthenticateUseCase {
	return c.authenticateUC
}

// Close fecha recursos (lifecycle management)
func (c *Container) Close(ctx context.Context) error {
	if c.closeFunc != nil {
		return c.closeFunc()
	}
	return nil
}
