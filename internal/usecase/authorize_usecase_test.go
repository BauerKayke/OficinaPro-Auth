package usecase_test

import (
	"testing"
	"time"

	"github.com/oficinapro/auth-service/internal/infrastructure/jwt"
	"github.com/oficinapro/auth-service/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthorizeUseCase(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")

	// Act
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Assert
	assert.NotNil(t, uc, "UseCase should not be nil")
}

func TestAuthorizeUseCase_Execute_ValidToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Gerar token válido
	payload := map[string]interface{}{
		"email":  "test@test.com",
		"userId": 123,
		"role":   "USER",
	}
	validToken, err := jwtService.GenerateToken(payload, 1*time.Hour)
	require.NoError(t, err)

	// Act
	isAuthorized, err := uc.Execute(validToken)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isAuthorized, "Should be authorized with valid token")
}

func TestAuthorizeUseCase_Execute_ValidTokenWithBearerPrefix(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	payload := map[string]interface{}{
		"email": "admin@test.com",
	}
	validToken, err := jwtService.GenerateToken(payload, 1*time.Hour)
	require.NoError(t, err)

	// Act - com prefixo "Bearer "
	isAuthorized, err := uc.Execute("Bearer " + validToken)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isAuthorized, "Should handle Bearer prefix correctly")
}

func TestAuthorizeUseCase_Execute_EmptyToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Act
	isAuthorized, err := uc.Execute("")

	// Assert
	assert.NoError(t, err, "Empty token should not return error")
	assert.False(t, isAuthorized, "Should not be authorized with empty token")
}

func TestAuthorizeUseCase_Execute_InvalidToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Act
	isAuthorized, err := uc.Execute("invalid-token-string")

	// Assert
	assert.Error(t, err, "Should return error for invalid token")
	assert.False(t, isAuthorized, "Should not be authorized with invalid token")
	assert.Contains(t, err.Error(), "token validation failed")
}

func TestAuthorizeUseCase_Execute_ExpiredToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Gerar token expirado (1 nanosegundo)
	payload := map[string]interface{}{
		"email": "expired@test.com",
	}
	expiredToken, err := jwtService.GenerateToken(payload, 1*time.Nanosecond)
	require.NoError(t, err)

	// Aguardar para garantir expiração
	time.Sleep(10 * time.Millisecond)

	// Act
	isAuthorized, err := uc.Execute(expiredToken)

	// Assert
	assert.Error(t, err, "Should return error for expired token")
	assert.False(t, isAuthorized, "Should not be authorized with expired token")
}

func TestAuthorizeUseCase_Execute_TokenWithoutEmailClaim(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Gerar token sem claim "email"
	payload := map[string]interface{}{
		"userId": 123,
		"role":   "USER",
		// Sem "email"
	}
	tokenWithoutEmail, err := jwtService.GenerateToken(payload, 1*time.Hour)
	require.NoError(t, err)

	// Act
	isAuthorized, err := uc.Execute(tokenWithoutEmail)

	// Assert
	assert.Error(t, err, "Should return error when email claim is missing")
	assert.False(t, isAuthorized, "Should not be authorized without email claim")
	assert.Contains(t, err.Error(), "missing email claim")
}

func TestAuthorizeUseCase_Execute_TokenWithSpaces(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	payload := map[string]interface{}{
		"email": "test@test.com",
	}
	validToken, err := jwtService.GenerateToken(payload, 1*time.Hour)
	require.NoError(t, err)

	// Act - com espaços ao redor (Bearer + token)
	isAuthorized, err := uc.Execute("Bearer " + validToken + "  ")

	// Assert
	assert.NoError(t, err)
	assert.True(t, isAuthorized, "Should handle trailing spaces correctly")
}

func TestAuthorizeUseCase_Execute_TokenFromDifferentSecret(t *testing.T) {
	// Arrange
	jwtService1 := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	jwtService2 := jwt.NewJWTService("different-secret-key-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService1)

	// Gerar token com secret diferente
	payload := map[string]interface{}{
		"email": "test@test.com",
	}
	tokenFromDifferentSecret, err := jwtService2.GenerateToken(payload, 1*time.Hour)
	require.NoError(t, err)

	// Act
	isAuthorized, err := uc.Execute(tokenFromDifferentSecret)

	// Assert
	assert.Error(t, err, "Should return error for token signed with different secret")
	assert.False(t, isAuthorized, "Should not be authorized with wrong secret")
}

func TestAuthorizeUseCase_Execute_MalformedToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	uc := usecase.NewAuthorizeUseCase(jwtService)

	// Act - token malformado (sem partes suficientes)
	isAuthorized, err := uc.Execute("malformed.token")

	// Assert
	assert.Error(t, err, "Should return error for malformed token")
	assert.False(t, isAuthorized, "Should not be authorized with malformed token")
}
