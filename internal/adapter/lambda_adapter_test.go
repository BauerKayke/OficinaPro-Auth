package adapter_test

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oficinapro/auth-service/internal/adapter"
	"github.com/oficinapro/auth-service/internal/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLambdaAdapter_Handle_Success(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	lambdaAdapter := adapter.NewLambdaAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: 200,
		Body:       []byte(`{"token":"abc123"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return req.Method == "POST" &&
			req.Path == "/auth/login" &&
			string(req.Body) == `{"email":"test@test.com"}`
	})).Return(expectedResp)

	apiReq := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/auth/login",
		Body:       `{"email":"test@test.com"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// Act
	apiResp, err := lambdaAdapter.Handle(context.Background(), apiReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, apiResp.StatusCode)
	assert.Equal(t, `{"token":"abc123"}`, apiResp.Body)
	assert.Equal(t, "application/json", apiResp.Headers["Content-Type"])
	mockHandler.AssertExpectations(t)
}

func TestLambdaAdapter_Handle_WithHeaders(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	lambdaAdapter := adapter.NewLambdaAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: 200,
		Body:       []byte(`{"valid":true}`),
		Headers: map[string]string{
			"X-User-Id": "123",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return req.Headers["Authorization"] == "Bearer token123"
	})).Return(expectedResp)

	apiReq := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/auth/validate",
		Headers: map[string]string{
			"Authorization": "Bearer token123",
		},
	}

	// Act
	apiResp, err := lambdaAdapter.Handle(context.Background(), apiReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, apiResp.StatusCode)
	assert.Equal(t, "123", apiResp.Headers["X-User-Id"])
	mockHandler.AssertExpectations(t)
}

func TestLambdaAdapter_Handle_ErrorResponse(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	lambdaAdapter := adapter.NewLambdaAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: 401,
		Body:       []byte(`{"error":"unauthorized"}`),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	mockHandler.On("Handle", mock.Anything, mock.Anything).Return(expectedResp)

	apiReq := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/auth/login",
		Body:       `{"email":"invalid"}`,
	}

	// Act
	apiResp, err := lambdaAdapter.Handle(context.Background(), apiReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 401, apiResp.StatusCode)
	assert.Equal(t, `{"error":"unauthorized"}`, apiResp.Body)
	mockHandler.AssertExpectations(t)
}

func TestLambdaAdapter_Handle_EmptyBody(t *testing.T) {
	// Arrange
	mockHandler := new(MockHTTPHandler)
	lambdaAdapter := adapter.NewLambdaAdapter(mockHandler)

	expectedResp := handler.HTTPResponse{
		StatusCode: 200,
		Body:       []byte(`{"status":"healthy"}`),
		Headers:    map[string]string{},
	}

	mockHandler.On("Handle", mock.Anything, mock.MatchedBy(func(req handler.HTTPRequest) bool {
		return len(req.Body) == 0
	})).Return(expectedResp)

	apiReq := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/health",
	}

	// Act
	apiResp, err := lambdaAdapter.Handle(context.Background(), apiReq)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 200, apiResp.StatusCode)
	assert.Equal(t, `{"status":"healthy"}`, apiResp.Body)
	mockHandler.AssertExpectations(t)
}
