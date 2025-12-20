package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewJWTService(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"

	service := NewJWTService(secret, issuer)

	assert.NotNil(t, service)
	// Campos privados, apenas verificamos que a criação não retorna nil
}

func TestJWTServiceImpl_GenerateToken(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	payload := map[string]interface{}{
		"id":    int64(123),
		"nome":  "João Silva",
		"email": "joao@example.com",
	}
	expiration := 1 * time.Hour

	token, err := service.GenerateToken(payload, expiration)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTServiceImpl_ValidateToken_Success(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	payload := map[string]interface{}{
		"id":    int64(123),
		"nome":  "João Silva",
		"email": "joao@example.com",
	}
	expiration := 1 * time.Hour

	token, err := service.GenerateToken(payload, expiration)
	assert.NoError(t, err)

	// Validar token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	// Verificar que o issuer está nos claims
	assert.Equal(t, issuer, claims["iss"])
}

func TestJWTServiceImpl_ValidateToken_InvalidToken(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	// Token inválido
	invalidToken := "invalid.jwt.token"

	claims, err := service.ValidateToken(invalidToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTServiceImpl_ValidateToken_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	payload := map[string]interface{}{
		"id": int64(123),
	}
	// Token expirado (duração negativa simula token já expirado)
	expiration := -1 * time.Hour

	token, err := service.GenerateToken(payload, expiration)
	assert.NoError(t, err)

	// Validar token expirado
	claims, err := service.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}
