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
	"github.com/oficinapro/auth-service/internal/infrastructure/validator"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// NewContainer cria container com todas as dependências
// Approach: criar tudo inline de forma clara (KISS principle)
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

	// 3. Repositories
	clienteRepo := setupRepositories(db)

	// 4. Services
	jwtService, validatorService := setupServices(cfg)

	// 5. Use Cases
	authenticateUC, authorizeUC := setupUseCases(cfg, clienteRepo, jwtService, validatorService)

	// 6. Montar container
	return &Container{
		clienteRepo:      clienteRepo,
		jwtService:       jwtService,
		validatorService: validatorService,
		authenticateUC:   authenticateUC,
		authorizeUC:      authorizeUC,
		closeFunc:        closeDB,
	}, nil
}

// --- Helper functions (funções puras, fácil de testar) ---

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

func setupRepositories(db *gorm.DB) repository.ClienteRepository {
	return database.NewGormClienteRepository(db)
}

func setupServices(cfg *config.Config) (service.JWTService, service.ValidatorService) {
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Issuer)
	validatorService := validator.NewCPFValidator()
	return jwtService, validatorService
}

func setupUseCases(
	cfg *config.Config,
	clienteRepo repository.ClienteRepository,
	jwtService service.JWTService,
	validatorService service.ValidatorService,
) (*usecase.AuthenticateUseCase, *usecase.AuthorizeUseCase) {
	authUC := usecase.NewAuthenticateUseCase(
		clienteRepo,
		jwtService,
		validatorService,
		cfg.JWT.Expiration,
	)
	
	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)

	return authUC, authorizeUC
}
