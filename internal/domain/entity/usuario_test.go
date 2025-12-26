package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestUsuario_NewUsuario(t *testing.T) {
	// This function doesn't exist, so we test the struct creation manually
	usuario := &Usuario{
		Email:           "newuser@test.com",
		Senha:           "hashed-password",
		Role:            RoleUser,
		IsAtivo:         true,
		DataCriacao:     time.Now(),
		DataAtualizacao: time.Now(),
	}

	assert.Equal(t, "newuser@test.com", usuario.Email)
	assert.True(t, usuario.IsAtivo)
	assert.Equal(t, RoleUser, usuario.Role)
}

func TestUsuario_ValidatePassword(t *testing.T) {
	// Hash da senha "senha123"
	hashedPassword, err := HashPassword("senha123")
	require.NoError(t, err)

	tests := []struct {
		name          string
		storedHash    string
		plainPassword string
		expectedValid bool
	}{
		{
			name:          "Senha correta",
			storedHash:    hashedPassword,
			plainPassword: "senha123",
			expectedValid: true,
		},
		{
			name:          "Senha incorreta",
			storedHash:    hashedPassword,
			plainPassword: "senhaerrada",
			expectedValid: false,
		},
		{
			name:          "Senha vazia",
			storedHash:    hashedPassword,
			plainPassword: "",
			expectedValid: false,
		},
		{
			name:          "Case sensitive",
			storedHash:    hashedPassword,
			plainPassword: "SENHA123",
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usuario := &Usuario{
				ID:    1,
				Email: "test@test.com",
				Senha: tt.storedHash,
			}

			isValid := usuario.ValidatePassword(tt.plainPassword)
			assert.Equal(t, tt.expectedValid, isValid)
		})
	}
}

func TestUsuario_IsEnabled(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name            string
		isAtivo         bool
		dataDesativacao *time.Time
		expectedEnabled bool
	}{
		{
			name:            "Usuário ativo sem data de desativação",
			isAtivo:         true,
			dataDesativacao: nil,
			expectedEnabled: true,
		},
		{
			name:            "Usuário inativo",
			isAtivo:         false,
			dataDesativacao: nil,
			expectedEnabled: false,
		},
		{
			name:            "Usuário com data de desativação",
			isAtivo:         true,
			dataDesativacao: &now,
			expectedEnabled: false,
		},
		{
			name:            "Usuário inativo com data de desativação",
			isAtivo:         false,
			dataDesativacao: &now,
			expectedEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usuario := &Usuario{
				ID:              1,
				Email:           "test@test.com",
				IsAtivo:         tt.isAtivo,
				DataDesativacao: tt.dataDesativacao,
			}

			isEnabled := usuario.IsEnabled()
			assert.Equal(t, tt.expectedEnabled, isEnabled)
		})
	}
}

func TestUsuario_GetDisplayName(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		pessoaID     *int64
		expectedName string
	}{
		{
			name:         "Email como display name",
			email:        "user@test.com",
			pessoaID:     nil,
			expectedName: "user@test.com",
		},
		{
			name:         "Com pessoaID (futuro: buscar nome da pessoa)",
			email:        "admin@test.com",
			pessoaID:     ptrInt64(123),
			expectedName: "admin@test.com", // Por enquanto retorna email
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usuario := &Usuario{
				ID:       1,
				Email:    tt.email,
				PessoaID: tt.pessoaID,
			}

			displayName := usuario.GetDisplayName()
			assert.Equal(t, tt.expectedName, displayName)
		})
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name          string
		plainPassword string
		expectError   bool
	}{
		{
			name:          "Senha válida",
			plainPassword: "senha123",
			expectError:   false,
		},
		{
			name:          "Senha longa",
			plainPassword: "senhamuitolongacomcaracteresvariadosabc123!@#",
			expectError:   false,
		},
		{
			name:          "Senha curta",
			plainPassword: "abc",
			expectError:   false, // bcrypt aceita qualquer tamanho
		},
		{
			name:          "Senha vazia",
			plainPassword: "",
			expectError:   false, // bcrypt aceita string vazia
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.plainPassword)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)

				// Verificar que hash começa com $2a$ (bcrypt)
				assert.Contains(t, hash, "$2a$")

				// Verificar que consegue validar a senha
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.plainPassword))
				assert.NoError(t, err)
			}
		})
	}
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	// Mesmo password deve gerar hashes diferentes (salt automático)
	hash1, err1 := HashPassword("senha123")
	hash2, err2 := HashPassword("senha123")

	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "Hashes should be different due to random salt")

	// Mas ambos devem validar a mesma senha
	usuario1 := &Usuario{Senha: hash1}
	usuario2 := &Usuario{Senha: hash2}
	assert.True(t, usuario1.ValidatePassword("senha123"))
	assert.True(t, usuario2.ValidatePassword("senha123"))
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		expectedErr error
	}{
		{
			name:        "Email válido",
			email:       "user@example.com",
			expectedErr: nil,
		},
		{
			name:        "Email válido com subdomínio",
			email:       "admin@mail.example.com",
			expectedErr: nil,
		},
		{
			name:        "Email válido com números",
			email:       "user123@test123.com",
			expectedErr: nil,
		},
		{
			name:        "Email vazio",
			email:       "",
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Email sem @",
			email:       "invalidemail.com",
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Email sem domínio",
			email:       "user@",
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Email sem local part",
			email:       "@example.com",
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Email sem ponto no domínio",
			email:       "user@example",
			expectedErr: ErrInvalidEmail,
		},
		{
			name:        "Email muito curto",
			email:       "a@b",
			expectedErr: ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "Senha válida (6 chars)",
			password:    "senha1",
			expectedErr: nil,
		},
		{
			name:        "Senha válida longa",
			password:    "senhaSegura123!@#",
			expectedErr: nil,
		},
		{
			name:        "Senha muito curta (5 chars)",
			password:    "12345",
			expectedErr: ErrWeakPassword,
		},
		{
			name:        "Senha vazia",
			password:    "",
			expectedErr: ErrWeakPassword,
		},
		{
			name:        "Senha com 1 char",
			password:    "a",
			expectedErr: ErrWeakPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUsuario_RoleConstants(t *testing.T) {
	assert.Equal(t, "USER", RoleUser)
	assert.Equal(t, "ADMIN", RoleAdmin)
}

func TestUsuario_DomainErrors(t *testing.T) {
	// Verificar que todos os erros de domínio estão definidos
	assert.NotNil(t, ErrUsuarioNotFound)
	assert.NotNil(t, ErrInvalidCredentials)
	assert.NotNil(t, ErrUsuarioInactive)
	assert.NotNil(t, ErrInvalidEmail)
	assert.NotNil(t, ErrWeakPassword)
	assert.NotNil(t, ErrEmailAlreadyExists)

	// Verificar mensagens
	assert.Contains(t, ErrUsuarioNotFound.Error(), "usuário")
	assert.Contains(t, ErrInvalidCredentials.Error(), "credenciais")
	assert.Contains(t, ErrUsuarioInactive.Error(), "inativo")
	assert.Contains(t, ErrInvalidEmail.Error(), "email")
	assert.Contains(t, ErrWeakPassword.Error(), "senha")
	assert.Contains(t, ErrEmailAlreadyExists.Error(), "email")
}

func TestUsuario_ToAuthPayload(t *testing.T) {
	tests := []struct {
		name     string
		usuario  *Usuario
		expected *AuthPayload
	}{
		{
			name: "Usuário completo",
			usuario: &Usuario{
				ID:    123,
				Email: "admin@oficinapro.com",
				Role:  RoleAdmin,
			},
			expected: &AuthPayload{
				UserID: 123,
				Email:  "admin@oficinapro.com",
				Nome:   "admin@oficinapro.com",
				Role:   RoleAdmin,
			},
		},
		{
			name: "Usuário regular",
			usuario: &Usuario{
				ID:    456,
				Email: "user@example.com",
				Role:  RoleUser,
			},
			expected: &AuthPayload{
				UserID: 456,
				Email:  "user@example.com",
				Nome:   "user@example.com",
				Role:   RoleUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := tt.usuario.ToAuthPayload()
			assert.Equal(t, tt.expected.UserID, payload.UserID)
			assert.Equal(t, tt.expected.Email, payload.Email)
			assert.Equal(t, tt.expected.Nome, payload.Nome)
			assert.Equal(t, tt.expected.Role, payload.Role)
		})
	}
}

func TestUsuario_CompleteFlow(t *testing.T) {
	// Simular fluxo completo de criação e autenticação de usuário

	// 1. Criar senha hash
	plainPassword := "senhaSegura123"
	hashedPassword, err := HashPassword(plainPassword)
	require.NoError(t, err)

	// 2. Criar usuário
	usuario := &Usuario{
		ID:              1,
		Email:           "admin@oficinapro.com",
		Senha:           hashedPassword,
		Role:            RoleAdmin,
		PessoaID:        nil,
		DataCriacao:     time.Now(),
		DataAtualizacao: time.Now(),
		DataDesativacao: nil,
		IsAtivo:         true,
	}

	// 3. Verificar que está habilitado
	assert.True(t, usuario.IsEnabled())

	// 4. Validar senha correta
	assert.True(t, usuario.ValidatePassword(plainPassword))

	// 5. Validar senha incorreta
	assert.False(t, usuario.ValidatePassword("senhaErrada"))

	// 6. Display name
	assert.Equal(t, "admin@oficinapro.com", usuario.GetDisplayName())

	// 7. ToAuthPayload
	payload := usuario.ToAuthPayload()
	assert.Equal(t, usuario.ID, payload.UserID)
	assert.Equal(t, usuario.Email, payload.Email)
	assert.Equal(t, usuario.Role, payload.Role)

	// 8. Desativar usuário
	now := time.Now()
	usuario.IsAtivo = false
	usuario.DataDesativacao = &now

	// 9. Verificar que não está mais habilitado
	assert.False(t, usuario.IsEnabled())
}

// Helper functions
func ptrInt64(v int64) *int64 {
	return &v
}
