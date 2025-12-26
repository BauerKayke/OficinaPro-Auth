package adapter_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oficinapro/auth-service/internal/adapter"
	"github.com/oficinapro/auth-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHTTPHandler mock do HTTPHandler para testes
type MockHTTPHandler struct {
	mock.Mock
}

func (m *MockHTTPHandler) Handle(ctx context.Context, req handler.HTTPRequest) handler.HTTPResponse {
	args := m.Called(ctx, req)
	return args.Get(0).(handler.HTTPResponse)
}

func TestHTTPAdapter_ServeHTTP_Success(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	adapter := adapter.NewHTTPAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       []byte(`{"message":"success"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return req.Method == "POST" &&
			req.Path == "/auth/login" &&
			string(req.Body) == `{"email":"test@test.com"}`
	})).Return(expectedResp)

	// Act
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(`{"email":"test@test.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"success"}`, w.Body.String())
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	mockHandler.AssertExpectations(t)
}

func TestHTTPAdapter_ServeHTTP_WithHeaders(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	adapter := adapter.NewHTTPAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       []byte(`{}`),
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return req.Headers["Authorization"] == "Bearer token123"
	})).Return(expectedResp)

	// Act
	req := httptest.NewRequest("GET", "/auth/validate", nil)
	req.Header.Set("Authorization", "Bearer token123")
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "custom-value", w.Header().Get("X-Custom-Header"))
	mockHandler.AssertExpectations(t)
}

func TestHTTPAdapter_ServeHTTP_EmptyBody(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	adapter := adapter.NewHTTPAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: http.StatusNoContent,
		Body:       nil,
		Headers:    map[string]string{},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return len(req.Body) == 0
	})).Return(expectedResp)

	// Act
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
	mockHandler.AssertExpectations(t)
}

func TestHTTPAdapter_ServeHTTP_ErrorResponse(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	adapter := adapter.NewHTTPAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: http.StatusBadRequest,
		Body:       []byte(`{"error":"invalid request"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.Anything).Return(expectedResp)

	// Act
	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBufferString(`invalid json`))
	w := httptest.NewRecorder()

	adapter.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error":"invalid request"}`, w.Body.String())
	mockHandler.AssertExpectations(t)
}

