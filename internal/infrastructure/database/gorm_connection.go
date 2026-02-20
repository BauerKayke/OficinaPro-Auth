package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewGormConnection cria uma nova conexão GORM com PostgreSQL.
// Configura prepared statements, pool de conexões e logging apropriado.
func NewGormConnection(config Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=auth_schema",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	// Log (sem senha) para debug
	fmt.Printf("DB: Connecting to host=%s port=%s user=%s dbname=%s sslmode=%s\n",
		config.Host, config.Port, config.User, config.Database, config.SSLMode)

	gormLogger := logger.Default.LogMode(logger.Silent)
	if config.SSLMode == "disable" {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	fmt.Printf("DB: Opening connection...\n")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	fmt.Printf("DB: Getting sql.DB instance...\n")
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	fmt.Printf("DB: Configuring connection pool (max=%d, idle=%d, lifetime=%v)...\n",
		config.MaxConnections, config.MaxIdleConns, config.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(config.MaxConnections)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	fmt.Printf("DB: Pinging database...\n")
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	fmt.Printf("DB: Connection successful!\n")
	return db, nil
}

// CloseGorm fecha a conexão GORM e libera recursos.
// Deve ser chamado ao finalizar a aplicação para evitar vazamento de conexões.
func CloseGorm(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
