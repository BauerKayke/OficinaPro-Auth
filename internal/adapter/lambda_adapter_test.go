package adapter_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/oficinapro/auth-service/internal/adapter"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/handler"
	"github.com/oficinapro/auth-service/internal/infrastructure/jwt"
	"github.com/oficinapro/auth-service/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func TestNewLambdaAdapter(t *testing.T) {
	// Arrange - criar dependências mínimas
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)

	// AuthHandler precisa de muitas deps, então vamos passar nil e testar apenas a criação
	var authHandler *handler.AuthHandler = nil

	// Act
	lambdaAdapter := adapter.NewLambdaAdapter(authHandler, authorizeUC)

	// Assert
	assert.NotNil(t, lambdaAdapter, "Adapter should not be nil")
}

func TestLambdaAdapter_Handle_AuthorizerRequest_Authorized(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")

	// Gerar um token válido
	payload := map[string]interface{}{
		"email":  "test@test.com",
		"userId": 123,
		"role":   "USER",
	}
	validToken, err := jwtService.GenerateToken(payload, 3600*1000000000) // 1 hora
	assert.NoError(t, err)

	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)
	lambdaAdapter := adapter.NewLambdaAdapter(nil, authorizeUC)

	// Criar evento de authorizer com token válido
	authEvent := events.APIGatewayV2CustomAuthorizerV2Request{
		RouteArn: "arn:aws:execute-api:us-east-1:123456789012:abcdef123/prod/GET/path",
		Headers: map[string]string{
			"authorization": "Bearer " + validToken,
		},
	}

	rawEvent, err := json.Marshal(authEvent)
	assert.NoError(t, err)

	// Act
	result, err := lambdaAdapter.Handle(context.Background(), rawEvent)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)

	authResp, ok := result.(events.APIGatewayV2CustomAuthorizerSimpleResponse)
	assert.True(t, ok, "Expected APIGatewayV2CustomAuthorizerSimpleResponse")
	assert.True(t, authResp.IsAuthorized, "Should be authorized with valid token")
}

func TestLambdaAdapter_Handle_AuthorizerRequest_Unauthorized_InvalidToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)
	lambdaAdapter := adapter.NewLambdaAdapter(nil, authorizeUC)

	// Criar evento com token inválido
	authEvent := events.APIGatewayV2CustomAuthorizerV2Request{
		RouteArn: "arn:aws:execute-api:us-east-1:123456789012:abcdef123/prod/GET/path",
		Headers: map[string]string{
			"authorization": "Bearer invalid-token-here",
		},
	}

	rawEvent, err := json.Marshal(authEvent)
	assert.NoError(t, err)

	// Act
	result, err := lambdaAdapter.Handle(context.Background(), rawEvent)

	// Assert
	assert.NoError(t, err)
	authResp, ok := result.(events.APIGatewayV2CustomAuthorizerSimpleResponse)
	assert.True(t, ok)
	assert.False(t, authResp.IsAuthorized, "Should not be authorized with invalid token")
}

func TestLambdaAdapter_Handle_HTTPRequest_InvalidJSON(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)
	lambdaAdapter := adapter.NewLambdaAdapter(nil, authorizeUC)

	// Evento com JSON inválido
	rawEvent := json.RawMessage(`{invalid json`)

	// Act
	result, err := lambdaAdapter.Handle(context.Background(), rawEvent)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to unmarshal event")
}

func TestLambdaAdapter_Handle_AuthorizerRequest_Unauthorized_NoToken(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")
	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)
	lambdaAdapter := adapter.NewLambdaAdapter(nil, authorizeUC)

	// Criar evento sem authorization token mas com outro header para ser detectado como authorizer
	authEvent := events.APIGatewayV2CustomAuthorizerV2Request{
		RouteArn: "arn:aws:execute-api:us-east-1:123456789012:abcdef123/prod/GET/path",
		Headers: map[string]string{
			"content-type": "application/json", // Header qualquer para len(Headers) > 0
		},
	}

	rawEvent, err := json.Marshal(authEvent)
	assert.NoError(t, err)

	// Act
	result, err := lambdaAdapter.Handle(context.Background(), rawEvent)

	// Assert
	assert.NoError(t, err)
	authResp, ok := result.(events.APIGatewayV2CustomAuthorizerSimpleResponse)
	assert.True(t, ok)
	assert.False(t, authResp.IsAuthorized, "Should not be authorized without token")
}

func TestLambdaAdapter_Handle_AuthorizerRequest_CaseInsensitiveHeader(t *testing.T) {
	// Arrange
	jwtService := jwt.NewJWTService("test-secret-key-minimum-32-chars!!", "test-issuer")

	payload := map[string]interface{}{
		"email": "test@test.com",
	}
	validToken, err := jwtService.GenerateToken(payload, 3600*1000000000)
	assert.NoError(t, err)

	authorizeUC := usecase.NewAuthorizeUseCase(jwtService)
	lambdaAdapter := adapter.NewLambdaAdapter(nil, authorizeUC)

	// Testar com "Authorization" (maiúscula)
	authEvent := events.APIGatewayV2CustomAuthorizerV2Request{
		RouteArn: "arn:aws:execute-api:us-east-1:123456789012:abcdef123/prod/GET/path",
		Headers: map[string]string{
			"Authorization": "Bearer " + validToken,
		},
	}

	rawEvent, err := json.Marshal(authEvent)
	assert.NoError(t, err)

	// Act
	result, err := lambdaAdapter.Handle(context.Background(), rawEvent)

	// Assert
	assert.NoError(t, err)
	authResp, ok := result.(events.APIGatewayV2CustomAuthorizerSimpleResponse)
	assert.True(t, ok)
	assert.True(t, authResp.IsAuthorized, "Should work with capitalized Authorization header")
}

func TestLambdaAdapter_Handle_HTTPRequest(t *testing.T) {
	// Skip: HTTP requests require a real AuthHandler with all dependencies
	// Testing HTTP flow would require mocking database, JWT service, etc.
	// This is better tested through integration tests
	t.Skip("HTTP request handling requires full AuthHandler with dependencies")
}

func TestLambdaAdapter_Handle_InvalidJSON(t *testing.T) {
	// Skip: Invalid JSON may trigger HTTP path which requires AuthHandler
	// This edge case is better tested through integration tests
	t.Skip("Invalid JSON handling may require AuthHandler dependencies")
}

// MockJWTService para testes que precisam controlar o comportamento do JWT
type MockJWTService struct {
	validateFunc func(string) (map[string]interface{}, error)
}

func (m *MockJWTService) GenerateToken(payload map[string]interface{}, expiresIn time.Duration) (string, error) {
	return "mock-token", nil
}

func (m *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	if m.validateFunc != nil {
		return m.validateFunc(token)
	}
	return map[string]interface{}{"email": "test@test.com"}, nil
}

var _ service.JWTService = (*MockJWTService)(nil)
