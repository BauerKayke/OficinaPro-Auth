package adapter

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"github.com/oficinapro/auth-service/internal/handler"
)

// LambdaAdapter adapta APIGatewayProxyRequest para handler.HTTPRequest/Response.
// Implementa o padrão Adapter para desacoplar AWS Lambda do handler.
type LambdaAdapter struct {
	handler HTTPHandler
}

// NewLambdaAdapter cria um novo adapter Lambda.
func NewLambdaAdapter(h HTTPHandler) *LambdaAdapter {
	return &LambdaAdapter{
		handler: h,
	}
}

// Handle processa evento do API Gateway via handler HTTP abstrato.
func (a *LambdaAdapter) Handle(ctx context.Context, apiReq events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Converte APIGatewayProxyRequest → handler.HTTPRequest
	httpReq := handler.HTTPRequest{
		Method:  apiReq.HTTPMethod,
		Path:    apiReq.Path,
		Body:    []byte(apiReq.Body),
		Headers: apiReq.Headers,
	}

	// Processa com handler (framework agnostic)
	httpResp := a.handler.Handle(ctx, httpReq)

	// Converte handler.HTTPResponse → APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		StatusCode: httpResp.StatusCode,
		Headers:    httpResp.Headers,
		Body:       string(httpResp.Body),
	}, nil
}

