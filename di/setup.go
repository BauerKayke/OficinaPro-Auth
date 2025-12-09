package di

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/infrastructure/database"
	"github.com/oficinapro/auth-service/internal/infrastructure/jwt"
	"github.com/oficinapro/auth-service/internal/infrastructure/telemetry"
	"github.com/oficinapro/auth-service/internal/infrastructure/validator"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// NewContainer cria e configura o container DI com todas as dependências.
// Segue o princípio KISS: cria tudo de forma clara e explícita.
// Retorna erro se falhar ao configurar qualquer dependência crítica.
func NewContainer(ctx context.Context) (*Container, error) {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	// 2. Infrastructure (database)
	db, closeDB, err := setupDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}

	// 3. Telemetry (observabilidade)
	telemetryService, closeTelemetry, err := setupTelemetry(ctx, cfg)
	if err != nil {
		closeDB()
		return nil, fmt.Errorf("telemetry: %w", err)
	}

	// 4. Repositories
	clienteRepo := setupRepositories(db)

	// 5. Services
	jwtService, validatorService := setupServices(cfg)

	// 6. Use Cases (com telemetry injetado)
	authenticateUC := setupUseCases(cfg, clienteRepo, jwtService, validatorService, telemetryService)

	// 7. Montar container com cleanup function
	closeAll := func() error {
		closeTelemetry()
		return closeDB()
	}

	return &Container{
		clienteRepo:      clienteRepo,
		jwtService:       jwtService,
		validatorService: validatorService,
		telemetryService: telemetryService,
		authenticateUC:   authenticateUC,
		closeFunc:        closeAll,
	}, nil
}

// setupDatabase configura conexão com banco de dados PostgreSQL.
// Retorna conexão GORM, função de cleanup e erro caso falhe.
func setupDatabase(cfg *config.Config) (*gorm.DB, func() error, error) {
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Database:        cfg.Database.Database,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		SSLMode:         cfg.Database.SSLMode,
		MaxConnections:  cfg.Database.MaxConnections,
		MaxIdleConns:    cfg.Database.MaxIdleConnections,
		ConnMaxLifetime: cfg.Database.ConnectionTimeout,
	}

	db, err := database.NewGormConnection(dbConfig)
	if err != nil {
		return nil, nil, err
	}

	closeFunc := func() error {
		return database.CloseGorm(db)
	}

	return db, closeFunc, nil
}

// setupTelemetry configura serviço de observabilidade (OpenTelemetry + New Relic).
// Retorna telemetry service, função de cleanup e erro caso falhe.
func setupTelemetry(ctx context.Context, cfg *config.Config) (service.TelemetryService, func(), error) {
	// Se telemetry estiver desabilitado, retorna NoOp
	if !cfg.Telemetry.Enabled {
		return telemetry.NewNoOpTelemetryService(), func() {}, nil
	}

	// Criar configuração de telemetry
	telemetryCfg := telemetry.Config{
		ServiceName:      cfg.Telemetry.ServiceName,
		ServiceVersion:   cfg.Telemetry.ServiceVersion,
		Environment:      cfg.App.Environment,
		NewRelicKey:      cfg.Telemetry.NewRelicKey,
		NewRelicEndpoint: cfg.Telemetry.NewRelicEndpoint,
		SampleRate:       cfg.Telemetry.SampleRate,
	}

	// Criar serviço OpenTelemetry
	otelService, err := telemetry.NewOTelService(ctx, telemetryCfg)
	if err != nil {
		return nil, nil, err
	}

	// Função de cleanup
	cleanup := func() {
		_ = otelService.Shutdown(context.Background())
	}

	return otelService, cleanup, nil
}

// setupRepositories cria e configura repositórios da aplicação.
func setupRepositories(db *gorm.DB) repository.ClienteRepository {
	return database.NewGormClienteRepository(db)
}

// setupServices cria e configura serviços de domínio.
// Retorna JWT service e validator service configurados.
func setupServices(cfg *config.Config) (service.JWTService, service.ValidatorService) {
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)
	validatorService := validator.NewCPFValidator()
	return jwtService, validatorService
}

// setupUseCases cria e configura casos de uso da aplicação.
// Injeta todas as dependências necessárias via parâmetros, incluindo telemetry.
func setupUseCases(
	cfg *config.Config,
	clienteRepo repository.ClienteRepository,
	jwtService service.JWTService,
	validatorService service.ValidatorService,
	telemetryService service.TelemetryService,
) *usecase.AuthenticateUseCase {
	return usecase.NewAuthenticateUseCase(
		clienteRepo,
		jwtService,
		validatorService,
		telemetryService,
		cfg.JWT.Expiration,
	)
}
