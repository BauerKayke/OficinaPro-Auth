package handler

import (
	"context"
	"strings"
	"time"

	"github.com/oficinapro/auth-service/internal/handler/request"
	"github.com/oficinapro/auth-service/internal/handler/response"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// AuthUseCase interface (programar contra abstrações - DIP)
type AuthUseCase interface {
	Execute(ctx context.Context, input usecase.AuthenticateInput) (*usecase.AuthenticateOutput, error)
}

// JWTService representa operações de JWT
type JWTService interface {
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

// AuthHandler gerencia autenticação HTTP (framework agnostic)
type AuthHandler struct {
	authUC          AuthUseCase
	jwtService      JWTService
	errorMapper     ErrorMapper
	logger          Logger
	responseBuilder *ResponseBuilder
}

// NewAuthHandler cria novo handler
func NewAuthHandler(authUC AuthUseCase, jwtService JWTService, errorMapper ErrorMapper, logger Logger) *AuthHandler {
	return &AuthHandler{
		authUC:          authUC,
		jwtService:      jwtService,
		errorMapper:     errorMapper,
		logger:          logger,
		responseBuilder: NewResponseBuilder(),
	}
}

// Handle processa request de autenticação com roteamento baseado em path
func (h *AuthHandler) Handle(ctx context.Context, req HTTPRequest) HTTPResponse {
	start := time.Now()

	// Roteamento baseado em método + path
	switch {
	case req.Method == "POST" && req.Path == "/auth":
		return h.handleAuthenticate(ctx, req, start)

	case req.Method == "POST" && req.Path == "/auth/validate":
		return h.handleValidate(ctx, req, start)

	case req.Method == "GET" && req.Path == "/health":
		return h.handleHealth(ctx, req, start)

	case req.Method == "OPTIONS":
		return h.responseBuilder.CORSPreflight()

	default:
		return h.responseBuilder.NotFound("Endpoint not found")
	}
}

// handleAuthenticate processa autenticação via Email+Senha (endpoint principal)
func (h *AuthHandler) handleAuthenticate(ctx context.Context, req HTTPRequest, start time.Time) HTTPResponse {
	// Parse e valida request
	authReq, err := request.ParseAuthRequest(string(req.Body))
	if err != nil {
		h.logger.Error("Failed to parse request", "error", err)
		return h.responseBuilder.BadRequest(err.Error())
	}

	authReq.NormalizeEmail()

	if err := authReq.Validate(); err != nil {
		h.logger.Error("Request validation failed", "error", err)
		return h.responseBuilder.BadRequest(err.Error())
	}

	// Executa use case
	ucInput := usecase.AuthenticateInput{
		Email: authReq.Email,
		Senha: authReq.Senha,
	}
	ucOutput, err := h.authUC.Execute(ctx, ucInput)
	if err != nil {
		h.logger.Error("Authentication failed", "email", maskEmail(authReq.Email), "error", err)
		return h.errorMapper.Map(err)
	}

	// Cria response de sucesso
	authResp := response.NewAuthResponse(
		ucOutput.Token,
		ucOutput.ExpiresIn,
		ucOutput.UserID,
		ucOutput.Email,
		ucOutput.Nome,
		ucOutput.Role,
	)

	h.logger.Info("Authentication successful", "email", maskEmail(authReq.Email), "userId", ucOutput.UserID, "duration", time.Since(start))
	return h.responseBuilder.Success(authResp)
}

// handleValidate valida um JWT token (introspection endpoint)
func (h *AuthHandler) handleValidate(ctx context.Context, req HTTPRequest, start time.Time) HTTPResponse {
	// Extrai token do header
	token, err := extractBearerToken(req.Headers)
	if err != nil {
		h.logger.Error("Missing authorization header")
		return h.responseBuilder.Unauthorized("Missing authorization header")
	}

	// Valida token
	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		h.logger.Error("Token validation failed", "error", err.Error())
		return h.responseBuilder.Unauthorized("Invalid token")
	}

	h.logger.Info("Token validated successfully", "duration", time.Since(start))
	return h.responseBuilder.Success(claims)
}

// handleHealth retorna status de saúde do serviço
func (h *AuthHandler) handleHealth(ctx context.Context, req HTTPRequest, start time.Time) HTTPResponse {
	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"service":   "auth-gateway",
		"version":   "1.0.0",
	}

	h.logger.Info("Health check", "duration", time.Since(start))
	return h.responseBuilder.Success(health)
}

// extractBearerToken extrai token do header Authorization.
// Suporta formato "Bearer <token>" ou apenas "<token>".
func extractBearerToken(headers map[string]string) (string, error) {
	authHeader, exists := headers["Authorization"]
	if !exists {
		authHeader, exists = headers["authorization"]
	}

	if !exists || authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	// Remover prefixo "Bearer " se presente
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:], nil
	}

	return authHeader, nil
}

// maskEmail mascara email para logging seguro.
// Exemplo: "joao.silva@example.com" -> "jo***@example.com"
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 2 {
		return "***@" + domain
	}

	return local[:2] + "***@" + domain
}
