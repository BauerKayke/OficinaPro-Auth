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
	lambdaAdapter    *adapter.LambdaAdapter
	container        *di.Container
	containerInitErr error
)

func init() {
	// Tentar inicializar container DI (incluindo OpenTelemetry)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	container, err = di.NewContainer(ctx)
	if err != nil {
		// NÃO falhar o init, apenas logar o erro
		// Isso permite que health check funcione mesmo sem DB
		containerInitErr = err
		log.Printf("WARNING: Failed to initialize container (will work in degraded mode): %v", err)

		// Criar adapter com handler nulo para health check básico
		lambdaAdapter = adapter.NewLambdaAdapter(nil, nil)
		return
	}

	// Criar AuthHandler
	errorMapper := handler.NewDefaultErrorMapper()

	// Usar logger JSON estruturado em produção para observabilidade
	var logger handler.Logger
	if os.Getenv("ENVIRONMENT") == "production" {
		logger = infraLogger.NewJSONLogger(
			os.Getenv("TELEMETRY_SERVICE_NAME"),
			os.Getenv("TELEMETRY_SERVICE_VERSION"),
			"production",
		)
	} else {
		logger = infraLogger.NewStdLogger()
	}

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

	log.Println("Lambda initialized successfully with OpenTelemetry instrumentation")
}

func main() {
	// Verificar se houve erro na inicialização
	if containerInitErr != nil {
		log.Printf("WARNING: Running in degraded mode (health check only): %v", containerInitErr)
	}

	// Wrapper Lambda com OpenTelemetry (se disponível)
	var wrappedHandler interface{}
	if container != nil && container.TelemetryService() != nil {
		wrappedHandler = otellambda.InstrumentHandler(
			lambdaAdapter.Handle,
			otellambda.WithFlusher(container.TelemetryService()),
		)
	} else {
		// Sem telemetry se container falhou
		wrappedHandler = lambdaAdapter.Handle
	}

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
			// Lambda tem até 2 segundos após SIGTERM antes de hard kill
			// Usar 1.5s para ter margem de segurança
			ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
			defer cancel()

			// Forçar flush de telemetria antes do shutdown
			log.Println("Flushing telemetry data...")
			if err := container.TelemetryService().ForceFlush(ctx); err != nil {
				// Não logar "context canceled" como erro crítico
				if err != context.Canceled && err != context.DeadlineExceeded {
					log.Printf("Warning: failed to flush telemetry: %v", err)
				}
			}

			// Shutdown completo
			if err := container.Close(ctx); err != nil {
				// Não logar "context canceled" como erro crítico
				if err != context.Canceled && err != context.DeadlineExceeded {
					log.Printf("Warning: telemetry shutdown error: %v", err)
				}
			} else {
				log.Println("Telemetry shutdown completed successfully")
			}
		}

		os.Exit(0)
	}()
}
