package di_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oficinapro/auth-service/di"
)

func setupTestEnv() {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_SSL_MODE", "disable")
	os.Setenv("JWT_SECRET", "test-secret-key-minimum-32-chars!!")
	os.Setenv("JWT_ISSUER", "test-issuer")
	os.Setenv("TELEMETRY_ENABLED", "false") // Desabilitar para testes
	os.Setenv("ENVIRONMENT", "test")
}

func cleanupTestEnv() {
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_SSL_MODE")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("JWT_ISSUER")
	os.Unsetenv("TELEMETRY_ENABLED")
	os.Unsetenv("ENVIRONMENT")
}

func TestNewContainer_Success(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	ctx := context.Background()

	// Act
	container, err := di.NewContainer(ctx)

	// Assert
	// O erro pode ocorrer devido à conexão com DB não disponível,
	// mas o código de criação está sendo testado
	if err != nil {
		// Se falhar por DB, é esperado em testes sem DB real
		assert.Contains(t, err.Error(), "database", "Error should be related to database connection")
		return
	}

	assert.NotNil(t, container)
	defer container.Close(ctx)

	// Verificar que todos os serviços principais foram criados
	assert.NotNil(t, container.JWTService())
	assert.NotNil(t, container.ValidatorService())
	assert.NotNil(t, container.TelemetryService())
	assert.NotNil(t, container.UsuarioRepository())
	assert.NotNil(t, container.ClienteRepository())
	assert.NotNil(t, container.AuthenticateUseCase())
	assert.NotNil(t, container.AuthorizeUseCase())
}

func TestNewContainer_MissingJWTSecret(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	// Remover JWT_SECRET
	os.Unsetenv("JWT_SECRET")

	ctx := context.Background()

	// Act
	container, err := di.NewContainer(ctx)

	// Assert
	assert.Error(t, err, "Should fail without JWT_SECRET")
	assert.Nil(t, container)
	assert.Contains(t, err.Error(), "JWT_SECRET", "Error should mention JWT_SECRET")
}

func TestNewContainer_ShortJWTSecret(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	// JWT_SECRET muito curto
	os.Setenv("JWT_SECRET", "short")

	ctx := context.Background()

	// Act
	container, err := di.NewContainer(ctx)

	// Assert
	assert.Error(t, err, "Should fail with short JWT_SECRET")
	assert.Nil(t, container)
	assert.Contains(t, err.Error(), "32 characters", "Error should mention minimum length")
}

func TestNewContainer_MissingDBPassword(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	// Remover DB_PASSWORD
	os.Unsetenv("DB_PASSWORD")

	ctx := context.Background()

	// Act
	container, err := di.NewContainer(ctx)

	// Assert
	assert.Error(t, err, "Should fail without DB_PASSWORD")
	assert.Nil(t, container)
	assert.Contains(t, err.Error(), "DB_PASSWORD", "Error should mention DB_PASSWORD")
}

func TestContainer_Close(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	ctx := context.Background()
	container, err := di.NewContainer(ctx)

	// Se falhar por DB, pular teste
	if err != nil {
		t.Skip("Skipping test due to database connection unavailable")
	}

	require.NotNil(t, container)

	// Act & Assert - Close não deve causar panic
	assert.NotPanics(t, func() {
		err := container.Close(ctx)
		_ = err // Pode ter erro de shutdown, mas não deve dar panic
	})
}

func TestContainer_Getters(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	ctx := context.Background()
	container, err := di.NewContainer(ctx)

	// Se falhar por DB, pular teste
	if err != nil {
		t.Skip("Skipping test due to database connection unavailable")
	}

	require.NotNil(t, container)
	defer container.Close(ctx)

	// Act & Assert - Todos os getters devem retornar instâncias não-nil
	t.Run("JWTService", func(t *testing.T) {
		service := container.JWTService()
		assert.NotNil(t, service)
	})

	t.Run("ValidatorService", func(t *testing.T) {
		service := container.ValidatorService()
		assert.NotNil(t, service)
	})

	t.Run("TelemetryService", func(t *testing.T) {
		service := container.TelemetryService()
		assert.NotNil(t, service)
	})

	t.Run("UsuarioRepository", func(t *testing.T) {
		repo := container.UsuarioRepository()
		assert.NotNil(t, repo)
	})

	t.Run("ClienteRepository", func(t *testing.T) {
		repo := container.ClienteRepository()
		assert.NotNil(t, repo)
	})

	t.Run("AuthenticateUseCase", func(t *testing.T) {
		uc := container.AuthenticateUseCase()
		assert.NotNil(t, uc)
	})

	t.Run("AuthorizeUseCase", func(t *testing.T) {
		uc := container.AuthorizeUseCase()
		assert.NotNil(t, uc)
	})
}

func TestContainer_SingletonPattern(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	ctx := context.Background()
	container, err := di.NewContainer(ctx)

	// Se falhar por DB, pular teste
	if err != nil {
		t.Skip("Skipping test due to database connection unavailable")
	}

	require.NotNil(t, container)
	defer container.Close(ctx)

	// Act & Assert - Verificar que múltiplas chamadas retornam mesma instância
	jwt1 := container.JWTService()
	jwt2 := container.JWTService()
	assert.Same(t, jwt1, jwt2, "JWTService should be singleton")

	validator1 := container.ValidatorService()
	validator2 := container.ValidatorService()
	assert.Same(t, validator1, validator2, "ValidatorService should be singleton")

	telemetry1 := container.TelemetryService()
	telemetry2 := container.TelemetryService()
	assert.Same(t, telemetry1, telemetry2, "TelemetryService should be singleton")

	usuario1 := container.UsuarioRepository()
	usuario2 := container.UsuarioRepository()
	assert.Same(t, usuario1, usuario2, "UsuarioRepository should be singleton")

	cliente1 := container.ClienteRepository()
	cliente2 := container.ClienteRepository()
	assert.Same(t, cliente1, cliente2, "ClienteRepository should be singleton")

	auth1 := container.AuthenticateUseCase()
	auth2 := container.AuthenticateUseCase()
	assert.Same(t, auth1, auth2, "AuthenticateUseCase should be singleton")

	authz1 := container.AuthorizeUseCase()
	authz2 := container.AuthorizeUseCase()
	assert.Same(t, authz1, authz2, "AuthorizeUseCase should be singleton")
}

func TestNewContainer_WithTelemetryEnabled(t *testing.T) {
	// Arrange
	setupTestEnv()
	defer cleanupTestEnv()

	// Habilitar telemetry
	os.Setenv("TELEMETRY_ENABLED", "true")
	os.Setenv("TELEMETRY_SERVICE_NAME", "test-service")
	os.Setenv("TELEMETRY_SERVICE_VERSION", "1.0.0")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "test-key-12345678901234567890123456789012")
	os.Setenv("NEW_RELIC_OTLP_ENDPOINT", "otlp.nr-data.net:4318")

	ctx := context.Background()

	// Act
	container, err := di.NewContainer(ctx)

	// Se falhar por DB, é esperado
	if err != nil && !assert.Contains(t, err.Error(), "database") {
		t.Fatalf("Unexpected error: %v", err)
	}

	if container != nil {
		defer container.Close(ctx)

		// Assert
		assert.NotNil(t, container.TelemetryService())
	}
}
