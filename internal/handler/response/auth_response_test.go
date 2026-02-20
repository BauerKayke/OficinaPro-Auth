package response

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAuthResponse(t *testing.T) {
	token := "fake-jwt-token"
	expiresIn := 3600
	userID := int64(123)
	email := "user@example.com"
	nome := "João Silva"
	role := "USER"

	response := NewAuthResponse(token, expiresIn, userID, email, nome, role)

	assert.Equal(t, token, response.Token)
	assert.Equal(t, expiresIn, response.ExpiresIn)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, email, response.Email)
	assert.Equal(t, nome, response.Nome)
	assert.Equal(t, role, response.Role)
}

func TestNewAuthResponse_WithoutNome(t *testing.T) {
	token := "fake-jwt-token"
	expiresIn := 3600
	userID := int64(456)
	email := "admin@example.com"
	nome := ""
	role := "ADMIN"

	response := NewAuthResponse(token, expiresIn, userID, email, nome, role)

	assert.Equal(t, token, response.Token)
	assert.Equal(t, expiresIn, response.ExpiresIn)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, email, response.Email)
	assert.Empty(t, response.Nome)
	assert.Equal(t, role, response.Role)
}

func TestAuthResponse_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		response *AuthResponse
	}{
		{
			name: "Response completo",
			response: NewAuthResponse(
				"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				3600,
				1,
				"user@test.com",
				"Test User",
				"USER",
			),
		},
		{
			name: "Response sem nome",
			response: NewAuthResponse(
				"token123",
				7200,
				2,
				"admin@test.com",
				"",
				"ADMIN",
			),
		},
		{
			name: "Response role USER",
			response: NewAuthResponse(
				"token456",
				1800,
				3,
				"customer@test.com",
				"Customer Name",
				"USER",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.Marshal(tt.response)
			require.NoError(t, err)
			require.NotEmpty(t, jsonData)

			// Unmarshal back
			var unmarshaledResponse AuthResponse
			err = json.Unmarshal(jsonData, &unmarshaledResponse)
			require.NoError(t, err)

			// Verify fields
			assert.Equal(t, tt.response.Token, unmarshaledResponse.Token)
			assert.Equal(t, tt.response.ExpiresIn, unmarshaledResponse.ExpiresIn)
			assert.Equal(t, tt.response.UserID, unmarshaledResponse.UserID)
			assert.Equal(t, tt.response.Email, unmarshaledResponse.Email)
			assert.Equal(t, tt.response.Nome, unmarshaledResponse.Nome)
			assert.Equal(t, tt.response.Role, unmarshaledResponse.Role)
		})
	}
}

func TestAuthResponse_JSONFields(t *testing.T) {
	response := NewAuthResponse(
		"test-token",
		3600,
		123,
		"user@test.com",
		"Test User",
		"USER",
	)

	jsonData, err := json.Marshal(response)
	require.NoError(t, err)

	// Verify JSON contains expected fields
	jsonStr := string(jsonData)
	assert.Contains(t, jsonStr, `"token":"test-token"`)
	assert.Contains(t, jsonStr, `"expiresIn":3600`)
	assert.Contains(t, jsonStr, `"userId":123`)
	assert.Contains(t, jsonStr, `"email":"user@test.com"`)
	assert.Contains(t, jsonStr, `"nome":"Test User"`)
	assert.Contains(t, jsonStr, `"role":"USER"`)
}

func TestAuthResponse_DifferentRoles(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		expectedRole string
	}{
		{
			name:         "Role USER",
			role:         "USER",
			expectedRole: "USER",
		},
		{
			name:         "Role ADMIN",
			role:         "ADMIN",
			expectedRole: "ADMIN",
		},
		{
			name:         "Role custom",
			role:         "MANAGER",
			expectedRole: "MANAGER",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewAuthResponse(
				"token",
				3600,
				1,
				"test@test.com",
				"Test",
				tt.role,
			)
			assert.Equal(t, tt.expectedRole, response.Role)
		})
	}
}

func TestAuthResponse_RealWorldExample(t *testing.T) {
	// Simular resposta real de autenticação
	// Note: This is a fake JWT token for testing purposes only
	response := NewAuthResponse(
		"fake.test.token",
		3600,
		1,
		"admin@example.com",
		"Admin User",
		"ADMIN",
	)

	// Serialize to JSON
	jsonData, err := json.Marshal(response)
	require.NoError(t, err)

	// Verify structure
	var decoded map[string]interface{}
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)

	assert.Contains(t, decoded, "token")
	assert.Contains(t, decoded, "expiresIn")
	assert.Contains(t, decoded, "userId")
	assert.Contains(t, decoded, "email")
	assert.Contains(t, decoded, "nome")
	assert.Contains(t, decoded, "role")

	// Verify types
	assert.IsType(t, "", decoded["token"])
	assert.IsType(t, float64(0), decoded["expiresIn"]) // JSON numbers are float64
	assert.IsType(t, float64(0), decoded["userId"])
	assert.IsType(t, "", decoded["email"])
	assert.IsType(t, "", decoded["nome"])
	assert.IsType(t, "", decoded["role"])
}
