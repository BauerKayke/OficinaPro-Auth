package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oficinapro/auth-service/internal/handler/response"
)

func TestResponseBuilder_Success(t *testing.T) {
	rb := NewResponseBuilder()

	data := map[string]interface{}{
		"message": "success",
		"code":    200,
	}

	resp := rb.Success(data)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var result map[string]interface{}
	err := json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	assert.Equal(t, "success", result["message"])
}

func TestResponseBuilder_BadRequest(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.BadRequest("Invalid input")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Invalid input", errorResp.Error)
}

func TestResponseBuilder_Unauthorized(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.Unauthorized("Invalid credentials")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Invalid credentials", errorResp.Error)
}

func TestResponseBuilder_Forbidden(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.Forbidden("Access denied")

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Access denied", errorResp.Error)
}

func TestResponseBuilder_NotFound(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.NotFound("Resource not found")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Resource not found", errorResp.Error)
}

func TestResponseBuilder_MethodNotAllowed(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.MethodNotAllowed("Method not allowed")

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Method not allowed", errorResp.Error)
}

func TestResponseBuilder_InternalError(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.InternalError()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Internal server error", errorResp.Error)
}

func TestResponseBuilder_CORSPreflight(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.CORSPreflight()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, resp.Body)
	assert.NotEmpty(t, resp.Headers)
}

// Teste com struct que não pode ser serializada (canal, por exemplo)
type UnmarshallableStruct struct {
	Channel chan int `json:"channel"` // canais não podem ser serializados
}

func TestResponseBuilder_BuildResponse_MarshalError(t *testing.T) {
	rb := NewResponseBuilder()

	// Tentar serializar algo que causa erro no Marshal
	data := UnmarshallableStruct{
		Channel: make(chan int),
	}

	resp := rb.buildResponse(http.StatusOK, data)

	// Deve retornar InternalError quando marshaling falha
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestResponseBuilder_BuildErrorResponse(t *testing.T) {
	rb := NewResponseBuilder()

	resp := rb.buildErrorResponse(http.StatusBadRequest, "Test error")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.NotEmpty(t, resp.Body)

	var errorResp response.ErrorResponse
	err := json.Unmarshal(resp.Body, &errorResp)
	require.NoError(t, err)
	assert.Equal(t, "Test error", errorResp.Error)
	assert.NotEmpty(t, errorResp.Timestamp)
}

func TestNewResponseBuilder(t *testing.T) {
	rb := NewResponseBuilder()
	assert.NotNil(t, rb)
}

func TestResponseBuilder_ComplexData(t *testing.T) {
	rb := NewResponseBuilder()

	type ComplexData struct {
		ID       int64             `json:"id"`
		Name     string            `json:"name"`
		Metadata map[string]string `json:"metadata"`
		Tags     []string          `json:"tags"`
		Nested   map[string]int    `json:"nested"`
	}

	data := ComplexData{
		ID:   123,
		Name: "Test",
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Tags: []string{"tag1", "tag2", "tag3"},
		Nested: map[string]int{
			"count": 42,
		},
	}

	resp := rb.Success(data)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result ComplexData
	err := json.Unmarshal(resp.Body, &result)
	require.NoError(t, err)
	assert.Equal(t, int64(123), result.ID)
	assert.Equal(t, "Test", result.Name)
	assert.Len(t, result.Tags, 3)
	assert.Equal(t, 42, result.Nested["count"])
}
