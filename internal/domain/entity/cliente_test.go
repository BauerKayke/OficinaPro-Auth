package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCliente_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		cliente  Cliente
		expected bool
	}{
		{
			name: "Cliente válido",
			cliente: Cliente{
				ID:         1,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "12345678909",
				TipoPessoa: "FISICA",
				Ativo:      true,
			},
			expected: true,
		},
		{
			name: "ID zerado - inválido",
			cliente: Cliente{
				ID:         0,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "12345678909",
				TipoPessoa: "FISICA",
				Ativo:      true,
			},
			expected: false,
		},
		{
			name: "Documento vazio - inválido",
			cliente: Cliente{
				ID:         1,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "",
				TipoPessoa: "FISICA",
				Ativo:      true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cliente.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCliente_CanAuthenticate(t *testing.T) {
	tests := []struct {
		name        string
		cliente     Cliente
		expectedErr error
	}{
		{
			name: "Cliente pode autenticar",
			cliente: Cliente{
				ID:         1,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "12345678909",
				TipoPessoa: "FISICA",
				Ativo:      true,
			},
			expectedErr: nil,
		},
		{
			name: "Cliente inativo",
			cliente: Cliente{
				ID:         1,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "12345678909",
				TipoPessoa: "FISICA",
				Ativo:      false,
			},
			expectedErr: ErrClienteInativo,
		},
		{
			name: "Documento vazio",
			cliente: Cliente{
				ID:         1,
				Nome:       "João Silva",
				Email:      "joao@example.com",
				Documento:  "",
				TipoPessoa: "FISICA",
				Ativo:      true,
			},
			expectedErr: ErrInvalidDocument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cliente.CanAuthenticate()
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCliente_ToAuthPayload(t *testing.T) {
	cliente := Cliente{
		ID:         123,
		Nome:       "Maria Oliveira",
		Email:      "maria@example.com",
		Documento:  "98765432100",
		TipoPessoa: "FISICA",
		Ativo:      true,
	}

	result := cliente.ToAuthPayload()

	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.UserID)
	assert.Equal(t, "Maria Oliveira", result.Nome)
	assert.Equal(t, "maria@example.com", result.Email)
	assert.Equal(t, "USER", result.Role) // Clientes legados sempre USER

	// Testar conversão para map
	resultMap := result.ToMap()
	
	// Verificar campos principais
	assert.Equal(t, int64(123), resultMap["userId"])
	assert.Equal(t, "Maria Oliveira", resultMap["nome"])
	assert.Equal(t, "maria@example.com", resultMap["email"])
	assert.Equal(t, "USER", resultMap["role"])
}
