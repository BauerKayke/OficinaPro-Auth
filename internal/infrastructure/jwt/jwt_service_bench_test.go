package jwt

import (
	"testing"
	"time"
)

var (
	testSecret = "test-secret-key-min-32-characters-long-for-security"
	testIssuer = "oficinapro-bench"
)

// Benchmark JWT token generation
func BenchmarkGenerateToken(b *testing.B) {
	service := NewJWTService(testSecret, testIssuer)
	payload := map[string]interface{}{
		"cliente_id": int64(123),
		"nome":       "João Silva",
		"email":      "joao@example.com",
		"documento":  "12345678909",
	}
	expiration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(payload, expiration)
	}
}

// Benchmark JWT token validation
func BenchmarkValidateToken(b *testing.B) {
	service := NewJWTService(testSecret, testIssuer)
	payload := map[string]interface{}{
		"cliente_id": int64(123),
		"nome":       "João Silva",
	}

	token, _ := service.GenerateToken(payload, 24*time.Hour)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(token)
	}
}

// Benchmark ciclo completo (generate + validate)
func BenchmarkJWTFullCycle(b *testing.B) {
	service := NewJWTService(testSecret, testIssuer)
	payload := map[string]interface{}{
		"cliente_id": int64(123),
		"nome":       "João Silva",
	}
	expiration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _ := service.GenerateToken(payload, expiration)
		_, _ = service.ValidateToken(token)
	}
}

// Benchmark diferentes tamanhos de payload
func BenchmarkGenerateTokenSmallPayload(b *testing.B) {
	service := NewJWTService(testSecret, testIssuer)
	payload := map[string]interface{}{
		"id": int64(123),
	}
	expiration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(payload, expiration)
	}
}

func BenchmarkGenerateTokenLargePayload(b *testing.B) {
	service := NewJWTService(testSecret, testIssuer)
	payload := map[string]interface{}{
		"cliente_id":  int64(123),
		"nome":        "João Silva",
		"email":       "joao@example.com",
		"documento":   "12345678909",
		"tipo_pessoa": "FISICA",
		"telefone":    "11999999999",
		"endereco":    "Rua Teste, 123",
		"cidade":      "São Paulo",
		"estado":      "SP",
		"cep":         "01234-567",
	}
	expiration := 24 * time.Hour

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GenerateToken(payload, expiration)
	}
}
