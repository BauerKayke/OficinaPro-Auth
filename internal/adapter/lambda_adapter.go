package adapter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"

	"github.com/oficinapro/auth-service/internal/handler"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// LambdaAdapter adapta eventos AWS para lógica interna
type LambdaAdapter struct {
	authHandler *handler.AuthHandler
	authorizeUC *usecase.AuthorizeUseCase
}

// NewLambdaAdapter cria novo adapter
func NewLambdaAdapter(authHandler *handler.AuthHandler, authorizeUC *usecase.AuthorizeUseCase) *LambdaAdapter {
	return &LambdaAdapter{
		authHandler: authHandler,
		authorizeUC: authorizeUC,
	}
}

// Handle é o ponto de entrada genérico
func (a *LambdaAdapter) Handle(ctx context.Context, rawEvent json.RawMessage) (interface{}, error) {
	// 1. Tentar detectar se é um Authorizer Request
	var authEvent events.APIGatewayV2CustomAuthorizerV2Request
	if err := json.Unmarshal(rawEvent, &authEvent); err == nil {
		// Verificamos campos chaves para confirmar
		if authEvent.RouteArn != "" && len(authEvent.Headers) > 0 {
			return a.handleAuthorizer(ctx, authEvent)
		}
	}

	// 2. Assumir que é HTTP Request (Login) - Formato V2
	var httpEvent events.APIGatewayV2HTTPRequest
	if err := json.Unmarshal(rawEvent, &httpEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return a.handleHttp(ctx, httpEvent)
}

// handleAuthorizer trata validação de token
func (a *LambdaAdapter) handleAuthorizer(ctx context.Context, req events.APIGatewayV2CustomAuthorizerV2Request) (events.APIGatewayV2CustomAuthorizerSimpleResponse, error) {
	token := req.Headers["authorization"]
	if token == "" {
		token = req.Headers["Authorization"]
	}

	isAuthorized, err := a.authorizeUC.Execute(token)
	if err != nil {
		// Log erro, mas não retorna erro para o Gateway (apenas nega acesso)
		fmt.Printf("Authorization failed: %v\n", err)
		return events.APIGatewayV2CustomAuthorizerSimpleResponse{IsAuthorized: false}, nil
	}

	return events.APIGatewayV2CustomAuthorizerSimpleResponse{
		IsAuthorized: isAuthorized,
		Context: map[string]interface{}{
			"user": "authenticated",
		},
	}, nil
}

// handleHttp trata login
func (a *LambdaAdapter) handleHttp(ctx context.Context, apiReq events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Health check independente (não depende do banco)
	if apiReq.RequestContext.HTTP.Method == "GET" && apiReq.RequestContext.HTTP.Path == "/health" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 200,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"status":"healthy","service":"auth-lambda"}`,
		}, nil
	}

	// Verificar se handler está disponível (container inicializado)
	if a.authHandler == nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 503,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"error":"Service temporarily unavailable","message":"Database connection failed"}`,
		}, nil
	}

	// Converter APIGateway V2 → HTTP abstrato
	httpReq := handler.HTTPRequest{
		Method:  apiReq.RequestContext.HTTP.Method,
		Path:    apiReq.RequestContext.HTTP.Path,
		Body:    []byte(apiReq.Body),
		Headers: apiReq.Headers,
	}

	// Processar com handler HTTP (framework agnostic)
	httpResp := a.authHandler.Handle(ctx, httpReq)

	// Converter HTTP abstrato → APIGateway V2
	return events.APIGatewayV2HTTPResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Headers,
		Body:       string(httpResp.Body),
	}, nil
}
