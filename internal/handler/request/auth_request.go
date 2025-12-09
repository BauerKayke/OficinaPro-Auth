package request

import (
	"encoding/json"
	"errors"
)

// AuthRequest representa o request de autenticação
type AuthRequest struct {
	CPF string `json:"cpf"`
}

// Validate valida o request de autenticação
func (r *AuthRequest) Validate() error {
	if r.CPF == "" {
		return errors.New("CPF is required")
	}

	return nil
}

// ParseAuthRequest faz parse do body JSON para AuthRequest
func ParseAuthRequest(body string) (*AuthRequest, error) {
	var request AuthRequest

	if err := json.Unmarshal([]byte(body), &request); err != nil {
		return nil, errors.New("invalid request body")
	}

	return &request, nil
}
