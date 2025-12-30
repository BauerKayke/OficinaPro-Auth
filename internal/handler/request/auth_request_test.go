package request

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRequest_Validate(t *testing.T) {
	tests := []struct {
		name        string
		request     AuthRequest
		expectedErr bool
		errContains string
	}{
		{
			name: "Email e senha válidos",
			request: AuthRequest{
				Email: "user@example.com",
				Senha: "senha123",
			},
			expectedErr: false,
		},
		{
			name: "Email válido com diferentes domínios",
			request: AuthRequest{
				Email: "admin@oficinapro.com.br",
				Senha: "senhaSegura123",
			},
			expectedErr: false,
		},
		{
			name: "Email vazio - inválido",
			request: AuthRequest{
				Email: "",
				Senha: "senha123",
			},
			expectedErr: true,
			errContains: "email is required",
		},
		{
			name: "Email com espaços - inválido",
			request: AuthRequest{
				Email: "   ",
				Senha: "senha123",
			},
			expectedErr: true,
			errContains: "email is required",
		},
		{
			name: "Email sem @ - inválido",
			request: AuthRequest{
				Email: "invalidemail.com",
				Senha: "senha123",
			},
			expectedErr: true,
			errContains: "invalid email format",
		},
		{
			name: "Email sem domínio - inválido",
			request: AuthRequest{
				Email: "user@",
				Senha: "senha123",
			},
			expectedErr: true,
			errContains: "invalid email format",
		},
		{
			name: "Senha vazia - inválido",
			request: AuthRequest{
				Email: "user@example.com",
				Senha: "",
			},
			expectedErr: true,
			errContains: "senha is required",
		},
		{
			name: "Senha com espaços apenas - inválido",
			request: AuthRequest{
				Email: "user@example.com",
				Senha: "     ",
			},
			expectedErr: true,
			errContains: "senha is required",
		},
		{
			name: "Senha muito curta - inválido",
			request: AuthRequest{
				Email: "user@example.com",
				Senha: "12345",
			},
			expectedErr: true,
			errContains: "senha must be at least 6 characters",
		},
		{
			name: "Senha exatamente 6 caracteres - válido",
			request: AuthRequest{
				Email: "user@example.com",
				Senha: "123456",
			},
			expectedErr: false,
		},
		{
			name: "Email e senha vazios - inválido",
			request: AuthRequest{
				Email: "",
				Senha: "",
			},
			expectedErr: true,
			errContains: "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.expectedErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthRequest_NormalizeEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedEmail string
	}{
		{
			name:          "Email em lowercase - mantém",
			email:         "user@example.com",
			expectedEmail: "user@example.com",
		},
		{
			name:          "Email em uppercase - converte para lowercase",
			email:         "USER@EXAMPLE.COM",
			expectedEmail: "user@example.com",
		},
		{
			name:          "Email mixed case - converte para lowercase",
			email:         "User@Example.Com",
			expectedEmail: "user@example.com",
		},
		{
			name:          "Email com espaços - remove espaços",
			email:         "  user@example.com  ",
			expectedEmail: "user@example.com",
		},
		{
			name:          "Email com espaços e uppercase",
			email:         "  USER@EXAMPLE.COM  ",
			expectedEmail: "user@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := AuthRequest{Email: tt.email}
			req.NormalizeEmail()
			assert.Equal(t, tt.expectedEmail, req.Email)
		})
	}
}

func TestParseAuthRequest(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		expectedEmail string
		expectedSenha string
		expectedErr   bool
	}{
		{
			name:          "JSON válido",
			body:          `{"email":"user@example.com","senha":"senha123"}`,
			expectedEmail: "user@example.com",
			expectedSenha: "senha123",
			expectedErr:   false,
		},
		{
			name:          "JSON válido com campos extras (ignorados)",
			body:          `{"email":"admin@test.com","senha":"pass","extra":"field"}`,
			expectedEmail: "admin@test.com",
			expectedSenha: "pass",
			expectedErr:   false,
		},
		{
			name:          "JSON inválido",
			body:          `{invalid json}`,
			expectedEmail: "",
			expectedSenha: "",
			expectedErr:   true,
		},
		{
			name:          "Corpo vazio",
			body:          ``,
			expectedEmail: "",
			expectedSenha: "",
			expectedErr:   true,
		},
		{
			name:          "JSON vazio (objeto)",
			body:          `{}`,
			expectedEmail: "",
			expectedSenha: "",
			expectedErr:   false, // Parse OK, mas Validate vai falhar
		},
		{
			name:          "JSON com apenas email",
			body:          `{"email":"user@test.com"}`,
			expectedEmail: "user@test.com",
			expectedSenha: "",
			expectedErr:   false, // Parse OK, mas Validate vai falhar
		},
		{
			name:          "JSON com apenas senha",
			body:          `{"senha":"password"}`,
			expectedEmail: "",
			expectedSenha: "password",
			expectedErr:   false, // Parse OK, mas Validate vai falhar
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseAuthRequest(tt.body)
			if tt.expectedErr {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expectedEmail, result.Email)
				assert.Equal(t, tt.expectedSenha, result.Senha)
			}
		})
	}
}

func TestAuthRequest_Integration(t *testing.T) {
	t.Run("Parse + Normalize + Validate - fluxo completo válido", func(t *testing.T) {
		body := `{"email":"USER@EXAMPLE.COM","senha":"senha123"}`

		req, err := ParseAuthRequest(body)
		require.NoError(t, err)
		require.NotNil(t, req)

		// Normalizar email
		req.NormalizeEmail()
		assert.Equal(t, "user@example.com", req.Email)

		// Validar
		err = req.Validate()
		assert.NoError(t, err)
	})

	t.Run("Parse + Validate - fluxo com erro", func(t *testing.T) {
		body := `{"email":"invalidemail","senha":"abc"}`

		req, err := ParseAuthRequest(body)
		require.NoError(t, err) // Parse OK

		// Validar deve falhar
		err = req.Validate()
		require.Error(t, err)
	})
}
