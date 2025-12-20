package response

import (
	"time"
)

// ErrorResponse representa uma resposta de erro padronizada da API.
// Formato consistente facilita tratamento de erros no cliente.
type ErrorResponse struct {
	Error     string `json:"error"`             // Mensagem de erro legível
	Timestamp string `json:"timestamp"`         // Timestamp UTC do erro
	Details   string `json:"details,omitempty"` // Detalhes adicionais opcionais
}

// NewErrorResponse cria uma nova resposta de erro com timestamp atual.
// Usado para erros simples sem detalhes adicionais.
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewErrorResponseWithDetails cria uma resposta de erro com detalhes adicionais.
// Útil para fornecer contexto extra sobre o erro ao cliente.
func NewErrorResponseWithDetails(message, details string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   details,
	}
}
