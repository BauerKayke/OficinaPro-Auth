package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// Setup environment variables
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PASSWORD", "dummy_password_for_tests")
	os.Setenv("JWT_SECRET", "dummy_jwt_secret_min_32_chars_long")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)
	require.NotNil(t, config)

	assert.Equal(t, "testhost", config.Database.Host)
	assert.Equal(t, "dummy_password_for_tests", config.Database.Password)
	assert.Equal(t, "dummy_jwt_secret_min_32_chars_long", config.JWT.Secret)
}

func TestLoad_Defaults(t *testing.T) {
	// Setup minimal required env vars
	os.Setenv("DB_PASSWORD", "required")
	os.Setenv("JWT_SECRET", "test-jwt-secret-min-32-characters-long")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	// Verificar defaults
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, "5432", config.Database.Port)
	assert.Equal(t, "oficinapro", config.Database.Database)
	assert.Equal(t, "postgres", config.Database.User)
	assert.Equal(t, "disable", config.Database.SSLMode)
	assert.Equal(t, "us-east-1", config.AWS.Region)
	assert.Equal(t, "development", config.App.Environment)
	assert.Equal(t, "info", config.App.LogLevel)
}

func TestLoad_ProductionDefaults(t *testing.T) {
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("DB_PASSWORD", "required")
	os.Setenv("JWT_SECRET", "test-jwt-secret-min-32-characters-long")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	// Produção tem pool menor
	assert.Equal(t, 2, config.Database.MaxConnections)
	assert.Equal(t, 1, config.Database.MaxIdleConnections)
	assert.Equal(t, 10*time.Second, config.Database.ConnectionTimeout)
}

func TestLoad_DevelopmentDefaults(t *testing.T) {
	os.Setenv("ENVIRONMENT", "development")
	os.Setenv("DB_PASSWORD", "required")
	os.Setenv("JWT_SECRET", "test-jwt-secret-min-32-characters-long")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	// Desenvolvimento tem pool maior
	assert.Equal(t, 10, config.Database.MaxConnections)
	assert.Equal(t, 5, config.Database.MaxIdleConnections)
	assert.Equal(t, 30*time.Second, config.Database.ConnectionTimeout)
}

func TestValidate_MissingDBHost(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "", // Missing
			Password: "dummy_password",
		},
		JWT: JWTConfig{
			Secret: "dummy_jwt_secret_min_32_chars_long",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_HOST")
}

func TestValidate_MissingDBPassword(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Password: "", // Missing
		},
		JWT: JWTConfig{
			Secret: "test-jwt-secret-min-32-characters-long",
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestValidate_MissingJWTSecret(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Password: "dummy_password",
		},
		JWT: JWTConfig{
			Secret: "", // Missing
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}

func TestValidate_ShortJWTSecret(t *testing.T) {
	config := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Password: "dummy_password",
		},
		JWT: JWTConfig{
			Secret: "short", // Too short
		},
	}

	err := config.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 32 characters")
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue int
		expected     int
	}{
		{
			name:         "Valid integer",
			envValue:     "42",
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "Empty value uses default",
			envValue:     "",
			defaultValue: 10,
			expected:     10,
		},
		{
			name:         "Invalid value uses default",
			envValue:     "abc",
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_INT", tt.envValue)
			} else {
				os.Unsetenv("TEST_INT")
			}

			result := getEnvAsInt("TEST_INT", tt.defaultValue)
			assert.Equal(t, tt.expected, result)

			os.Unsetenv("TEST_INT")
		})
	}
}

func TestGetEnvAsDuration(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "Hours as integer",
			envValue:     "24",
			defaultValue: 1 * time.Hour,
			expected:     24 * time.Hour,
		},
		{
			name:         "Duration string",
			envValue:     "30m",
			defaultValue: 1 * time.Hour,
			expected:     30 * time.Minute,
		},
		{
			name:         "Empty uses default",
			envValue:     "",
			defaultValue: 1 * time.Hour,
			expected:     1 * time.Hour,
		},
		{
			name:         "Invalid uses default",
			envValue:     "invalid",
			defaultValue: 1 * time.Hour,
			expected:     1 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_DURATION", tt.envValue)
			} else {
				os.Unsetenv("TEST_DURATION")
			}

			result := getEnvAsDuration("TEST_DURATION", tt.defaultValue)
			assert.Equal(t, tt.expected, result)

			os.Unsetenv("TEST_DURATION")
		})
	}
}

func TestGetEnvAsBool(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "True",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "False",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "1 is true",
			envValue:     "1",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "0 is false",
			envValue:     "0",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "Empty uses default",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "Invalid uses default",
			envValue:     "invalid",
			defaultValue: false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_BOOL", tt.envValue)
			} else {
				os.Unsetenv("TEST_BOOL")
			}

			result := getEnvAsBool("TEST_BOOL", tt.defaultValue)
			assert.Equal(t, tt.expected, result)

			os.Unsetenv("TEST_BOOL")
		})
	}
}

func TestGetEnvAsFloat(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue float64
		expected     float64
	}{
		{
			name:         "Valid float",
			envValue:     "3.14",
			defaultValue: 1.0,
			expected:     3.14,
		},
		{
			name:         "Integer as float",
			envValue:     "42",
			defaultValue: 1.0,
			expected:     42.0,
		},
		{
			name:         "Empty uses default",
			envValue:     "",
			defaultValue: 1.5,
			expected:     1.5,
		},
		{
			name:         "Invalid uses default",
			envValue:     "abc",
			defaultValue: 2.0,
			expected:     2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("TEST_FLOAT", tt.envValue)
			} else {
				os.Unsetenv("TEST_FLOAT")
			}

			result := getEnvAsFloat("TEST_FLOAT", tt.defaultValue)
			assert.Equal(t, tt.expected, result)

			os.Unsetenv("TEST_FLOAT")
		})
	}
}

func TestTelemetryConfig(t *testing.T) {
	os.Setenv("DB_PASSWORD", "dummy_required")
	os.Setenv("JWT_SECRET", "dummy_jwt_secret_min_32_chars_long")
	os.Setenv("TELEMETRY_ENABLED", "true")
	os.Setenv("TELEMETRY_SERVICE_NAME", "test-service")
	os.Setenv("TELEMETRY_SERVICE_VERSION", "2.0.0")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "dummy_new_relic_key")
	os.Setenv("TELEMETRY_SAMPLE_RATE", "0.5")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	assert.True(t, config.Telemetry.Enabled)
	assert.Equal(t, "test-service", config.Telemetry.ServiceName)
	assert.Equal(t, "2.0.0", config.Telemetry.ServiceVersion)
	assert.Equal(t, "dummy_new_relic_key", config.Telemetry.NewRelicKey)
	assert.Equal(t, 0.5, config.Telemetry.SampleRate)
}

func TestCustomDatabaseConfig(t *testing.T) {
	os.Setenv("DB_HOST", "custom-host")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_NAME", "custom-db")
	os.Setenv("DB_USER", "custom-user")
	os.Setenv("DB_PASSWORD", "dummy_custom_pass")
	os.Setenv("DB_SSL_MODE", "require")
	os.Setenv("DB_MAX_CONNECTIONS", "20")
	os.Setenv("DB_MAX_IDLE_CONNECTIONS", "10")
	os.Setenv("DB_CONNECTION_TIMEOUT", "45s")
	os.Setenv("JWT_SECRET", "dummy_jwt_secret_min_32_chars_long")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "custom-host", config.Database.Host)
	assert.Equal(t, "5433", config.Database.Port)
	assert.Equal(t, "custom-db", config.Database.Database)
	assert.Equal(t, "custom-user", config.Database.User)
	assert.Equal(t, "dummy_custom_pass", config.Database.Password)
	assert.Equal(t, "require", config.Database.SSLMode)
	assert.Equal(t, 20, config.Database.MaxConnections)
	assert.Equal(t, 10, config.Database.MaxIdleConnections)
	assert.Equal(t, 45*time.Second, config.Database.ConnectionTimeout)
}

func TestJWTConfig(t *testing.T) {
	os.Setenv("DB_PASSWORD", "dummy_required")
	os.Setenv("JWT_SECRET", "dummy_very_long_jwt_secret_key_for_testing")
	os.Setenv("JWT_EXPIRATION_HOURS", "48")
	os.Setenv("JWT_ISSUER", "custom-issuer")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "dummy_very_long_jwt_secret_key_for_testing", config.JWT.Secret)
	assert.Equal(t, 48*time.Hour, config.JWT.Expiration)
	assert.Equal(t, "custom-issuer", config.JWT.Issuer)
}

func TestAWSConfig(t *testing.T) {
	os.Setenv("DB_PASSWORD", "dummy_required")
	os.Setenv("JWT_SECRET", "dummy_jwt_secret_min_32_chars_long")
	os.Setenv("AWS_REGION", "sa-east-1")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "sa-east-1", config.AWS.Region)
}

func TestAppConfig(t *testing.T) {
	os.Setenv("DB_PASSWORD", "dummy_required")
	os.Setenv("JWT_SECRET", "dummy_jwt_secret_min_32_chars_long")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("LOG_LEVEL", "debug")
	defer cleanupEnv()

	config, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "staging", config.App.Environment)
	assert.Equal(t, "debug", config.App.LogLevel)
}

// cleanupEnv limpa todas as variáveis de ambiente usadas nos testes
func cleanupEnv() {
	envVars := []string{
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD",
		"DB_SSL_MODE", "DB_MAX_CONNECTIONS", "DB_MAX_IDLE_CONNECTIONS",
		"DB_CONNECTION_TIMEOUT", "JWT_SECRET", "JWT_EXPIRATION_HOURS",
		"JWT_ISSUER", "AWS_REGION", "ENVIRONMENT", "LOG_LEVEL",
		"TELEMETRY_ENABLED", "TELEMETRY_SERVICE_NAME",
		"TELEMETRY_SERVICE_VERSION", "NEW_RELIC_LICENSE_KEY",
		"NEW_RELIC_OTLP_ENDPOINT", "TELEMETRY_SAMPLE_RATE",
	}

	for _, v := range envVars {
		os.Unsetenv(v)
	}
}
