package request

import (
	"encoding/json"
	"errors"
	"strings"
)

// AuthRequest representa a requisição de autenticação via CPF.
// Contém apenas os dados necessários para identificar o usuário.
type AuthRequest struct {
	CPF string `json:"cpf"` // CPF do cliente (pode conter pontos e traços)
}

// Validate valida os campos obrigatórios da requisição.
// Não valida formato do CPF (isso é responsabilidade do ValidatorService).
func (r *AuthRequest) Validate() error {
	// Trim espaços para evitar valores apenas com espaços
	r.CPF = strings.TrimSpace(r.CPF)

	if r.CPF == "" {
		return errors.New("CPF is required")
	}

	return nil
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
