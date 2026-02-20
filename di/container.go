package di

import (
	"context"
	"fmt"

	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// Container armazena todas as dependências da aplicação.
// Fornece acesso controlado via getters (encapsulamento).
// Responsável por gerenciar ciclo de vida dos recursos (Close).
type Container struct {
	clienteRepo      repository.ClienteRepository
	usuarioRepo      repository.UsuarioRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	telemetryService service.TelemetryService
	authenticateUC   *usecase.AuthenticateUseCase
	authorizeUC      *usecase.AuthorizeUseCase
	closeFunc        func() error
}

// NewContainer cria e configura o container DI com todas as dependências.
// Orquestra a inicialização seguindo ordem de dependências.
func NewContainer(ctx context.Context) (*Container, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	builder := newBuilder(ctx, cfg)
	return builder.build()
}

// ClienteRepository retorna o repositório de clientes configurado.
func (c *Container) ClienteRepository() repository.ClienteRepository {
	return c.clienteRepo
}

// UsuarioRepository retorna o repositório de usuários configurado.
func (c *Container) UsuarioRepository() repository.UsuarioRepository {
	return c.usuarioRepo
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

// AuthorizeUseCase retorna o caso de uso de autorização configurado.
func (c *Container) AuthorizeUseCase() *usecase.AuthorizeUseCase {
	return c.authorizeUC
}

// Close libera recursos do container (conexões de banco, telemetry, etc).
// Deve ser chamado ao encerrar a aplicação para evitar vazamento de recursos.
func (c *Container) Close(ctx context.Context) error {
	if c.telemetryService != nil {
		_ = c.telemetryService.Shutdown(ctx)
	}

	if c.closeFunc != nil {
		return c.closeFunc()
	}
	return nil
}
