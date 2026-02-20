package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/oficinapro/auth-service/di"
	"github.com/oficinapro/auth-service/internal/adapter"
	"github.com/oficinapro/auth-service/internal/handler"
	infraLogger "github.com/oficinapro/auth-service/internal/infrastructure/logger"
)

func main() {
	ctx := context.Background()

	// Inicializa dependências
	container, authHandler := initializeDependencies(ctx)
	defer container.Close(ctx)

	// Cria adapter HTTP
	httpAdapter := adapter.NewHTTPAdapter(authHandler)

	// Configura e inicia servidor
	server := createServer(httpAdapter)
	startServerWithGracefulShutdown(server)
}

// initializeDependencies inicializa DI container e cria AuthHandler.
func initializeDependencies(ctx context.Context) (*di.Container, *handler.AuthHandler) {
	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	errorMapper := handler.NewDefaultErrorMapper()
	logger := infraLogger.NewStdLogger()

	authHandler := handler.NewAuthHandler(
		container.AuthenticateUseCase(),
		container.JWTService(),
		errorMapper,
		logger,
	)

	return container, authHandler
}

// createServer cria servidor HTTP configurado.
func createServer(handler http.Handler) *http.Server {
	port := getPort()

	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// getPort retorna porta configurada ou padrão 8080.
func getPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

// startServerWithGracefulShutdown inicia servidor com graceful shutdown.
func startServerWithGracefulShutdown(server *http.Server) {
	// Configura graceful shutdown
	go setupGracefulShutdown(server)

	// Inicia servidor
	log.Printf("Auth Service HTTP server listening on port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}

// setupGracefulShutdown configura shutdown gracioso ao receber sinais.
func setupGracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
}
