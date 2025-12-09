package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/oficinapro/auth-service/di"
	"github.com/oficinapro/auth-service/internal/handler"
)

var lambdaAdapter *LambdaAdapter

func init() {
	ctx := context.Background()

	// 1. Criar container DI
	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// 2. Criar dependências do handler
	errorMapper := handler.NewDefaultErrorMapper()
	logger := &handler.StdLogger{}

	// 3. Criar HTTP handler (framework agnostic)
	authHandler := handler.NewAuthHandler(
		container.AuthenticateUseCase(),
		errorMapper,
		logger,
	)

	// 4. Criar adapter Lambda (camada externa - Adapter Pattern)
	lambdaAdapter = NewLambdaAdapter(authHandler)
}

func main() {
	lambda.Start(lambdaAdapter.Handle)
}
