package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Host               string
	Port               string
	Database           string
	User               string
	Password           string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	ConnectionTimeout  time.Duration
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
	Issuer     string
}

type AWSConfig struct {
	Region string
}

type AppConfig struct {
	Environment string
	LogLevel    string
}

func Load() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "production" {
		_ = godotenv.Load()
	}

	config := &Config{
		Database: DatabaseConfig{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnv("DB_PORT", "5432"),
			Database:           getEnv("DB_NAME", "oficinapro"),
			User:               getEnv("DB_USER", "postgres"),
			Password:           getEnv("DB_PASSWORD", ""),
			SSLMode:            getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:     getEnvAsInt("DB_MAX_CONNECTIONS", 10),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 5),
			ConnectionTimeout:  getEnvAsDuration("DB_CONNECTION_TIMEOUT", 30*time.Second),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: getEnvAsDuration("JWT_EXPIRATION_HOURS", 24*time.Hour),
			Issuer:     getEnv("JWT_ISSUER", "oficinapro-auth"),
		},
		AWS: AWSConfig{
			Region: getEnv("AWS_REGION", "us-east-1"),
		},
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	if hours, err := strconv.Atoi(valueStr); err == nil {
		return time.Duration(hours) * time.Hour
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
