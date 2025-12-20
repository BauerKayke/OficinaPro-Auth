package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		request     AuthRequest
		expectedErr bool
	}{
		{
			name:        "CPF válido",
			request:     AuthRequest{CPF: "123.456.789-09"},
			expectedErr: false,
		},
		{
			name:        "CPF vazio - inválido",
			request:     AuthRequest{CPF: ""},
			expectedErr: true,
		},
		{
			name:        "CPF com espaços - inválido",
			request:     AuthRequest{CPF: "   "},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseAuthRequest(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		expectedCPF string
		expectedErr bool
	}{
		{
			name:        "JSON válido",
			body:        `{"cpf":"123.456.789-09"}`,
			expectedCPF: "123.456.789-09",
			expectedErr: false,
		},
		{
			name:        "JSON inválido",
			body:        `{invalid json}`,
			expectedCPF: "",
			expectedErr: true,
		},
		{
			name:        "Corpo vazio",
			body:        ``,
			expectedCPF: "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAuthRequest(tt.body)
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedCPF, result.CPF)
			}
		})
	}
}
