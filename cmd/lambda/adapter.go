package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/oficinapro/auth-service/internal/handler"
)

// LambdaAdapter adapta APIGateway para HTTPRequest (Adapter Pattern)
type LambdaAdapter struct {
	authHandler *handler.AuthHandler
}

// NewLambdaAdapter cria novo adapter
func NewLambdaAdapter(authHandler *handler.AuthHandler) *LambdaAdapter {
	return &LambdaAdapter{
		authHandler: authHandler,
	}
}

// Handle adapta APIGatewayProxyRequest → HTTPRequest → HTTPResponse → APIGatewayProxyResponse
func (a *LambdaAdapter) Handle(ctx context.Context, apiReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Converter APIGateway → HTTP abstrato
	httpReq := handler.HTTPRequest{
		Method:  apiReq.HTTPMethod,
		Body:    []byte(apiReq.Body),
		Headers: apiReq.Headers,
	}

	// Processar com handler HTTP (framework agnostic)
	httpResp := a.authHandler.Handle(ctx, httpReq)

	// Converter HTTP abstrato → APIGateway
	return events.APIGatewayProxyResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Headers,
		Body:       string(httpResp.Body),
	}, nil
}
