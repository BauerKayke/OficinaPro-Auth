package response

import (
	"time"
)

// ErrorResponse representa uma resposta de erro padronizada
type ErrorResponse struct {
	Error     string `json:"error"`
	Timestamp string `json:"timestamp"`
	Details   string `json:"details,omitempty"`
}

// NewErrorResponse cria um novo ErrorResponse
func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// NewErrorResponseWithDetails cria um ErrorResponse com detalhes adicionais
func NewErrorResponseWithDetails(message, details string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details:   details,
	}
}
