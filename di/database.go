package di

import (
	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/infrastructure/database"
)

// NewDatabaseConnection cria e configura conexão com banco de dados PostgreSQL.
// Retorna conexão GORM, função de cleanup e erro caso falhe.
func NewDatabaseConnection(cfg *config.Config) (*gorm.DB, func() error, error) {
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
