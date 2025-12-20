package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAuthResponse(t *testing.T) {
	token := "fake-jwt-token"
	expiresIn := 3600
	clienteID := int64(123)
	nome := "João Silva"

	response := NewAuthResponse(token, expiresIn, clienteID, nome)

	assert.Equal(t, token, response.Token)
	assert.Equal(t, expiresIn, response.ExpiresIn)
	assert.Equal(t, clienteID, response.ClienteID)
	assert.Equal(t, nome, response.Nome)
}
