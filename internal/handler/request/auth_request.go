package request

import (
	"encoding/json"
	"errors"
	"strings"
)

// AuthRequest representa a requisição de autenticação via email e senha.
// Alinhado com AuthenticationRequest do core-domain-service (Java).
type AuthRequest struct {
	Email string `json:"email"` // Email do usuário
	Senha string `json:"senha"` // Senha em texto plano (será validada com bcrypt)
}

// Validate valida os campos obrigatórios da requisição.
func (r *AuthRequest) Validate() error {
	// Trim espaços para evitar valores apenas com espaços
	r.Email = strings.TrimSpace(r.Email)
	r.Senha = strings.TrimSpace(r.Senha)

	if r.Email == "" {
		return errors.New("email is required")
	}

	// Validação básica de formato de email
	if !strings.Contains(r.Email, "@") || !strings.Contains(r.Email, ".") {
		return errors.New("invalid email format")
	}

	if r.Senha == "" {
		return errors.New("senha is required")
	}

	// Validação mínima de tamanho (para login, não cadastro)
	if len(r.Senha) < 6 {
		return errors.New("senha must be at least 6 characters")
	}

	return nil
}

// NormalizeEmail normaliza email para lowercase para comparação case-insensitive.
func (r *AuthRequest) NormalizeEmail() {
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
}

// ParseAuthRequest faz parse do body JSON para AuthRequest.
// Retorna erro se o JSON for inválido ou malformado.
func ParseAuthRequest(body string) (*AuthRequest, error) {
	var request AuthRequest

	if err := json.Unmarshal([]byte(body), &request); err != nil {
		return nil, errors.New("invalid request body")
	}

	return &request, nil
}
