package di

import (
	"context"

	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// Container armazena todas as dependências da aplicação.
// Fornece acesso controlado via getters (encapsulamento).
// Responsável por gerenciar ciclo de vida dos recursos (Close).
type Container struct {
	clienteRepo      repository.ClienteRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	telemetryService service.TelemetryService
	authenticateUC   *usecase.AuthenticateUseCase
	closeFunc        func() error
}

// ClienteRepository retorna o repositório de clientes configurado.
func (c *Container) ClienteRepository() repository.ClienteRepository {
	return c.clienteRepo
}

// JWTService retorna o serviço JWT configurado.
func (c *Container) JWTService() service.JWTService {
	return c.jwtService
}

// ValidatorService retorna o serviço de validação configurado.
func (c *Container) ValidatorService() service.ValidatorService {
	return c.validatorService
}

// TelemetryService retorna o serviço de telemetria configurado.
func (c *Container) TelemetryService() service.TelemetryService {
	return c.telemetryService
}

// AuthenticateUseCase retorna o caso de uso de autenticação configurado.
func (c *Container) AuthenticateUseCase() *usecase.AuthenticateUseCase {
	return c.authenticateUC
}

// Close libera recursos do container (conexões de banco, telemetry, etc).
// Deve ser chamado ao encerrar a aplicação para evitar vazamento de recursos.
func (c *Container) Close(ctx context.Context) error {
	// Finalizar telemetry primeiro (flush de métricas/traces)
	if c.telemetryService != nil {
		if err := c.telemetryService.Shutdown(ctx); err != nil {
			// Log erro mas continua cleanup
		}
	}

	// Fechar outras conexões
	if c.closeFunc != nil {
		return c.closeFunc()
	}
	return nil
}
