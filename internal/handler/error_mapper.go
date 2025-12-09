package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/handler/response"
)

// ErrorMapper mapeia domain errors para HTTP responses
type ErrorMapper interface {
	Map(err error) HTTPResponse
	Register(err error, statusCode int, message string)
}

// DefaultErrorMapper implementação padrão
type DefaultErrorMapper struct {
	mappings map[error]errorMapping
}

type errorMapping struct {
	statusCode int
	message    string
}

// NewDefaultErrorMapper cria mapper com mapeamentos padrão
func NewDefaultErrorMapper() *DefaultErrorMapper {
	return &DefaultErrorMapper{
		mappings: map[error]errorMapping{
			entity.ErrInvalidCPF:      {http.StatusUnauthorized, "CPF inválido"},
			entity.ErrClienteNotFound: {http.StatusUnauthorized, "Cliente não encontrado"},
			entity.ErrClienteInativo:  {http.StatusUnauthorized, "Cliente inativo"},
			entity.ErrInvalidDocument: {http.StatusBadRequest, "Documento inválido"},
		},
	}
}

// Map mapeia erro para HTTP response
func (m *DefaultErrorMapper) Map(err error) HTTPResponse {
	// Buscar mapeamento específico
	if mapping, exists := m.mappings[err]; exists {
		return m.errorResponse(mapping.statusCode, mapping.message)
	}

	// Erro desconhecido = 500
	return m.errorResponse(http.StatusInternalServerError, "Internal server error")
}

// Register permite adicionar novos mapeamentos (OCP)
func (m *DefaultErrorMapper) Register(err error, statusCode int, message string) {
	m.mappings[err] = errorMapping{statusCode, message}
}

// errorResponse cria response de erro
func (m *DefaultErrorMapper) errorResponse(statusCode int, message string) HTTPResponse {
	errorResp := response.NewErrorResponse(message)
	body, err := json.Marshal(errorResp)
	if err != nil {
		// Fallback se marshalling falhar
		body = []byte(`{"error":"Internal server error","timestamp":""}`)
	}

	return NewHTTPResponse(statusCode, body)
}
