package service

import (
	"time"
)

// JWTService define o contrato para geração e validação de JWT
// Interface (Port) - Dependency Inversion Principle
type JWTService interface {
	// GenerateToken gera um token JWT para o payload fornecido
	GenerateToken(payload map[string]interface{}, expiresIn time.Duration) (string, error)

	// ValidateToken valida um token JWT
	ValidateToken(token string) (map[string]interface{}, error)
}
