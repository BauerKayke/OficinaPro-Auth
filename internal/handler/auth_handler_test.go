package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/handler/response"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// MockAuthUseCase mock do usecase
type MockAuthUseCase struct {
	mock.Mock
}

func (m *MockAuthUseCase) Execute(ctx context.Context, input usecase.AuthenticateInput) (*usecase.AuthenticateOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.AuthenticateOutput), args.Error(1)
}

func TestNewAuthHandler(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &StdLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	assert.NotNil(t, handler)
}

func TestAuthHandler_Handle_Success(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &NoOpLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	reqBody := `{"cpf":"123.456.789-09"}`
	req := HTTPRequest{
		Method: "POST",
		Body:   []byte(reqBody),
	}

	expectedOutput := &usecase.AuthenticateOutput{
		Token:     "fake-token",
		ExpiresIn: 3600,
		ClienteID: 123,
		Nome:      "João Silva",
	}

	mockUC.On("Execute", mock.Anything, mock.MatchedBy(func(input usecase.AuthenticateInput) bool {
		return input.CPF == "123.456.789-09"
	})).Return(expectedOutput, nil)

	resp := handler.Handle(context.Background(), req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var authResp response.AuthResponse
	err := json.Unmarshal(resp.Body, &authResp)
	assert.NoError(t, err)
	assert.Equal(t, "fake-token", authResp.Token)

	mockUC.AssertExpectations(t)
}

func TestAuthHandler_Handle_MethodNotAllowed(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &NoOpLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	req := HTTPRequest{
		Method: "GET",
		Body:   []byte(`{}`),
	}

	resp := handler.Handle(context.Background(), req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestAuthHandler_Handle_InvalidJSON(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &NoOpLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	req := HTTPRequest{
		Method: "POST",
		Body:   []byte(`{invalid json}`),
	}

	resp := handler.Handle(context.Background(), req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Handle_ValidationError(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &NoOpLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	// CPF vazio
	reqBody := `{"cpf":""}`
	req := HTTPRequest{
		Method: "POST",
		Body:   []byte(reqBody),
	}

	resp := handler.Handle(context.Background(), req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthHandler_Handle_UseCaseError(t *testing.T) {
	mockUC := new(MockAuthUseCase)
	mockMapper := NewDefaultErrorMapper()
	mockLogger := &NoOpLogger{}

	handler := NewAuthHandler(mockUC, mockMapper, mockLogger)

	reqBody := `{"cpf":"123.456.789-09"}`
	req := HTTPRequest{
		Method: "POST",
		Body:   []byte(reqBody),
	}

	mockUC.On("Execute", mock.Anything, mock.Anything).
		Return(nil, entity.ErrClienteNotFound)

	resp := handler.Handle(context.Background(), req)

	// DefaultErrorMapper mapeia ClienteNotFound para 401 Unauthorized
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	mockUC.AssertExpectations(t)
}

func TestNewHTTPResponse(t *testing.T) {
	statusCode := http.StatusOK
	body := []byte(`{"message":"test"}`)

	resp := NewHTTPResponse(statusCode, body)

	assert.Equal(t, statusCode, resp.StatusCode)
	assert.Equal(t, body, resp.Body)
	assert.Equal(t, "application/json", resp.Headers["Content-Type"])
}
