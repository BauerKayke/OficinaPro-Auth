package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	message := "Erro de validação"
	response := NewErrorResponse(message)

	assert.Equal(t, message, response.Error)
	assert.Empty(t, response.Details)
	assert.NotEmpty(t, response.Timestamp)
}

func TestNewErrorResponseWithDetails(t *testing.T) {
	message := "Erro com detalhes"
	details := "Campo CPF é inválido"

	response := NewErrorResponseWithDetails(message, details)

	assert.Equal(t, message, response.Error)
	assert.Equal(t, details, response.Details)
	assert.NotEmpty(t, response.Timestamp)
}
