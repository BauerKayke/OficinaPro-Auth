package handler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

func TestNewDefaultErrorMapper(t *testing.T) {
	mapper := NewDefaultErrorMapper()
	assert.NotNil(t, mapper)
}

func TestErrorMapper_Map_KnownErrors(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "ErrUsuarioNotFound",
			err:            entity.ErrUsuarioNotFound,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Credenciais inválidas",
		},
		{
			name:           "ErrInvalidCredentials",
			err:            entity.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedMsg:    "Credenciais inválidas",
		},
		{
			name:           "ErrUsuarioInactive",
			err:            entity.ErrUsuarioInactive,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Usuário inativo",
		},
		{
			name:           "ErrInvalidEmail",
			err:            entity.ErrInvalidEmail,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Email inválido",
		},
		{
			name:           "ErrWeakPassword",
			err:            entity.ErrWeakPassword,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Senha muito fraca",
		},
		{
			name:           "ErrEmailAlreadyExists",
			err:            entity.ErrEmailAlreadyExists,
			expectedStatus: http.StatusConflict,
			expectedMsg:    "Email já cadastrado",
		},
		{
			name:           "ErrInvalidCPF",
			err:            entity.ErrInvalidCPF,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "CPF inválido",
		},
		{
			name:           "ErrClienteNotFound",
			err:            entity.ErrClienteNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Cliente não encontrado",
		},
		{
			name:           "ErrClienteInativo",
			err:            entity.ErrClienteInativo,
			expectedStatus: http.StatusForbidden,
			expectedMsg:    "Cliente inativo",
		},
		{
			name:           "ErrInvalidDocument",
			err:            entity.ErrInvalidDocument,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Documento inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := mapper.Map(tt.err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Contains(t, string(resp.Body), tt.expectedMsg)
		})
	}
}

func TestErrorMapper_Map_UnknownError(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	unknownErr := errors.New("some unknown error")
	resp := mapper.Map(unknownErr)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "Internal server error")
}

func TestErrorMapper_Register(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	// Registrar novo erro
	customErr := errors.New("custom error")
	mapper.Register(customErr, http.StatusTeapot, "I'm a teapot")

	// Verificar que o erro foi registrado
	resp := mapper.Map(customErr)
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "I'm a teapot")
}

func TestErrorMapper_Register_OverrideExisting(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	// Sobrescrever erro existente
	mapper.Register(entity.ErrInvalidEmail, http.StatusNotAcceptable, "Custom message")

	// Verificar que foi sobrescrito
	resp := mapper.Map(entity.ErrInvalidEmail)
	assert.Equal(t, http.StatusNotAcceptable, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "Custom message")
}

func TestErrorMapper_ErrorResponse(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	resp := mapper.errorResponse(http.StatusBadRequest, "Test message")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)
	assert.Contains(t, string(resp.Body), "Test message")
	assert.Contains(t, string(resp.Body), "timestamp")
}

func TestErrorMapper_MultipleMappings(t *testing.T) {
	mapper := NewDefaultErrorMapper()

	// Registrar múltiplos erros customizados
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	mapper.Register(err1, http.StatusBadRequest, "Message 1")
	mapper.Register(err2, http.StatusForbidden, "Message 2")
	mapper.Register(err3, http.StatusNotFound, "Message 3")

	// Verificar todos
	resp1 := mapper.Map(err1)
	assert.Equal(t, http.StatusBadRequest, resp1.StatusCode)
	assert.Contains(t, string(resp1.Body), "Message 1")

	resp2 := mapper.Map(err2)
	assert.Equal(t, http.StatusForbidden, resp2.StatusCode)
	assert.Contains(t, string(resp2.Body), "Message 2")

	resp3 := mapper.Map(err3)
	assert.Equal(t, http.StatusNotFound, resp3.StatusCode)
	assert.Contains(t, string(resp3.Body), "Message 3")
}
