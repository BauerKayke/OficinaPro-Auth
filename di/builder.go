package di

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// builder encapsula a lógica de construção do container.
// Segue o padrão Builder para construção passo a passo.
type builder struct {
	ctx              context.Context
	cfg              *config.Config
	db               *gorm.DB
	closeDB          func() error
	telemetryService service.TelemetryService
	closeTelemetry   func()
	clienteRepo      repository.ClienteRepository
	usuarioRepo      repository.UsuarioRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	authenticateUC   *usecase.AuthenticateUseCase
	authorizeUC      *usecase.AuthorizeUseCase
}

// newBuilder cria um novo builder para construção do container.
func newBuilder(ctx context.Context, cfg *config.Config) *builder {
	return &builder{
		ctx: ctx,
		cfg: cfg,
	}
}

// build executa a construção do container na ordem correta de dependências.
func (b *builder) build() (*Container, error) {
	// 1. Infrastructure (database)
	if err := b.setupDatabase(); err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	// 2. Telemetry (observabilidade)
	if err := b.setupTelemetry(); err != nil {
		b.cleanup()
		return nil, fmt.Errorf("telemetry: %w", err)
	}

	// 3. Repositories
	b.setupRepositories()

	// 4. Services
	b.setupServices()

	// 5. Use Cases
	b.setupUseCases()

	return b.toContainer(), nil
}

// setupDatabase configura conexão com banco de dados.
func (b *builder) setupDatabase() error {
	db, closeDB, err := NewDatabaseConnection(b.cfg)
	if err != nil {
		return err
	}

	b.db = db
	b.closeDB = closeDB
	return nil
}

// setupTelemetry configura serviço de observabilidade.
func (b *builder) setupTelemetry() error {
	telemetryService, closeTelemetry, err := NewTelemetryService(b.ctx, b.cfg)
	if err != nil {
		return err
	}

	b.telemetryService = telemetryService
	b.closeTelemetry = closeTelemetry
	return nil
}

// setupRepositories cria e configura repositórios.
func (b *builder) setupRepositories() {
	b.clienteRepo, b.usuarioRepo = NewRepositories(b.db)
}

// setupServices cria e configura serviços de domínio.
func (b *builder) setupServices() {
	b.jwtService, b.validatorService = NewServices(b.cfg)
}

// setupUseCases cria e configura casos de uso.
func (b *builder) setupUseCases() {
	b.authenticateUC = NewAuthenticateUseCase(
		b.cfg,
		b.usuarioRepo,
		b.jwtService,
		b.telemetryService,
	)
	
	// Novo UseCase de Autorização
	b.authorizeUC = usecase.NewAuthorizeUseCase(b.jwtService)
}

// cleanup libera recursos em caso de erro durante construção.
func (b *builder) cleanup() {
	if b.closeTelemetry != nil {
		b.closeTelemetry()
	}
	if b.closeDB != nil {
		_ = b.closeDB()
	}
}

// toContainer converte builder para container final.
func (b *builder) toContainer() *Container {
	closeAll := func() error {
		if b.closeTelemetry != nil {
			b.closeTelemetry()
		}
		if b.closeDB != nil {
			return b.closeDB()
		}
		return nil
	}

	return &Container{
		clienteRepo:      b.clienteRepo,
		usuarioRepo:      b.usuarioRepo,
		jwtService:       b.jwtService,
		validatorService: b.validatorService,
		telemetryService: b.telemetryService,
		authenticateUC:   b.authenticateUC,
		authorizeUC:      b.authorizeUC,
		closeFunc:        closeAll,
	}
}
