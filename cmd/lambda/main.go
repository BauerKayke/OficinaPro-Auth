package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/oficinapro/auth-service/di"
	"github.com/oficinapro/auth-service/internal/adapter"
	"github.com/oficinapro/auth-service/internal/handler"
	infraLogger "github.com/oficinapro/auth-service/internal/infrastructure/logger"
)

var (
	lambdaAdapter *adapter.LambdaAdapter
	container     *di.Container
)

func init() {
	ctx := context.Background()

	// Inicializa dependências
	var err error
	container, err = di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Cria AuthHandler
	errorMapper := handler.NewDefaultErrorMapper()
	logger := infraLogger.NewStdLogger()

	authHandler := handler.NewAuthHandler(
		container.AuthenticateUseCase(),
		container.JWTService(),
		errorMapper,
		logger,
	)

	// Cria adapter Lambda
	lambdaAdapter = adapter.NewLambdaAdapter(authHandler)

	// Configura graceful shutdown
	setupGracefulShutdown()
}

func main() {
	lambda.Start(lambdaAdapter.Handle)
}

// setupGracefulShutdown configura limpeza de recursos quando Lambda container for desligado.
// Lambda envia SIGTERM antes de destruir o container.
func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, cleaning up resources...")
		if container != nil {
			ctx := context.Background()
			if err := container.Close(ctx); err != nil {
				log.Printf("Error closing container: %v", err)
			}
		}
		os.Exit(0)
	}()
}
