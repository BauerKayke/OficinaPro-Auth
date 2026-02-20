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

func TestJWTServiceImpl_ValidateToken_WrongSecret(t *testing.T) {
	// Criar token com um secret
	secret1 := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service1 := NewJWTService(secret1, issuer)

	payload := map[string]interface{}{
		"id": int64(123),
	}
	token, err := service1.GenerateToken(payload, 1*time.Hour)
	assert.NoError(t, err)

	// Tentar validar com secret diferente
	secret2 := "another-secret-key-min-32-characters"
	service2 := NewJWTService(secret2, issuer)

	claims, err := service2.ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTServiceImpl_ValidateToken_EmptyToken(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	claims, err := service.ValidateToken("")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestJWTServiceImpl_GenerateToken_WithComplexPayload(t *testing.T) {
	secret := "test-secret-key-min-32-characters-long"
	issuer := "test-issuer"
	service := NewJWTService(secret, issuer)

	payload := map[string]interface{}{
		"id":          int64(123),
		"nome":        "João Silva",
		"email":       "joao@example.com",
		"role":        "ADMIN",
		"permissions": []string{"read", "write", "delete"},
		"metadata": map[string]interface{}{
			"departamento": "TI",
			"nivel":        5,
		},
	}

	token, err := service.GenerateToken(payload, 1*time.Hour)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validar e verificar payload
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), claims["id"]) // JSON numbers são float64
	assert.Equal(t, "João Silva", claims["nome"])
	assert.Equal(t, "joao@example.com", claims["email"])
	assert.Equal(t, "ADMIN", claims["role"])
}
