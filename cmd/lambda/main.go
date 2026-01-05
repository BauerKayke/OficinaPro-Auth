package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"

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
	// Inicializar container DI (incluindo OpenTelemetry)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	container, err = di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Criar AuthHandler
	errorMapper := handler.NewDefaultErrorMapper()
	logger := infraLogger.NewStdLogger()

	authHandler := handler.NewAuthHandler(
		container.AuthenticateUseCase(),
		container.JWTService(),
		errorMapper,
		logger,
	)

	// Criar adapter Lambda
	lambdaAdapter = adapter.NewLambdaAdapter(authHandler, container.AuthorizeUseCase())

	// Configurar graceful shutdown
	setupGracefulShutdown()

	log.Println("Lambda initialized with OpenTelemetry instrumentation")
}

func main() {
	// Wrapper Lambda com OpenTelemetry
	// Usa otellambda para instrumentação automática
	wrappedHandler := otellambda.InstrumentHandler(lambdaAdapter.Handle)

	// Iniciar Lambda
	lambda.Start(wrappedHandler)
}

// setupGracefulShutdown configura limpeza de recursos quando Lambda container for desligado.
// Lambda envia SIGTERM antes de destruir o container.
func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, cleaning up resources...")

		// Cleanup container (incluindo telemetry)
		if container != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := container.Close(ctx); err != nil {
				log.Printf("Error closing container: %v", err)
			}
		}

		os.Exit(0)
	}()
}
