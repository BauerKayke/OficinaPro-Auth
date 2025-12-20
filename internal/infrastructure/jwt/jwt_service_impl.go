package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTServiceImpl implementa geração e validação de tokens JWT.
// Usa HMAC-SHA256 para assinatura dos tokens.
type JWTServiceImpl struct {
	secretKey []byte // Chave secreta para assinar tokens
	issuer    string // Emissor do token (identificação do serviço)
}

// NewJWTService cria nova instância do serviço JWT.
// secretKey deve ter no mínimo 32 caracteres para segurança adequada.
func NewJWTService(secretKey string, issuer string) *JWTServiceImpl {
	return &JWTServiceImpl{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken gera um novo token JWT com o payload fornecido.
// Adiciona claims padrão (iat, exp, iss) e mescla com o payload customizado.
// Retorna o token assinado ou erro em caso de falha na assinatura.
func (s *JWTServiceImpl) GenerateToken(payload map[string]interface{}, expiresIn time.Duration) (string, error) {
	now := time.Now()

	// Claims padrão do JWT
	claims := jwt.MapClaims{
		"iat": now.Unix(),                // issued at
		"exp": now.Add(expiresIn).Unix(), // expiration
		"iss": s.issuer,                  // issuer
	}

	// Mescla payload customizado com claims padrão
	for key, value := range payload {
		claims[key] = value
	}

	// Cria e assina token usando HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken valida um token JWT e retorna seus claims.
// Verifica assinatura, expiração e formato do token.
// Retorna os claims extraídos ou erro se o token for inválido.
func (s *JWTServiceImpl) ValidateToken(tokenString string) (map[string]interface{}, error) {
	// Parse e valida token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifica se o método de assinatura é HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Verifica se token é válido (não expirado, assinatura correta)
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extrai claims do token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
