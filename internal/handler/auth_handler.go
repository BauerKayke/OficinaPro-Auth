package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oficinapro/auth-service/internal/handler/request"
	"github.com/oficinapro/auth-service/internal/handler/response"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// AuthUseCase interface (programar contra abstrações - DIP)
type AuthUseCase interface {
	Execute(ctx context.Context, input usecase.AuthenticateInput) (*usecase.AuthenticateOutput, error)
}

// AuthHandler gerencia autenticação HTTP (framework agnostic)
type AuthHandler struct {
	authUC      AuthUseCase
	errorMapper ErrorMapper
	logger      Logger
}

// NewAuthHandler cria novo handler
func NewAuthHandler(authUC AuthUseCase, errorMapper ErrorMapper, logger Logger) *AuthHandler {
	return &AuthHandler{
		authUC:      authUC,
		errorMapper: errorMapper,
		logger:      logger,
	}
}

// Handle processa request de autenticação
func (h *AuthHandler) Handle(ctx context.Context, req HTTPRequest) HTTPResponse {
	start := time.Now()

	// 1. Validar método
	if req.Method != "POST" {
		return h.methodNotAllowedResponse()
	}

	// 2. Parse request
	authReq, err := request.ParseAuthRequest(string(req.Body))
	if err != nil {
		h.logger.Error("Failed to parse request", "error", err)
		return h.badRequestResponse(err.Error())
	}

	// 3. Validar request básico
	if err := authReq.Validate(); err != nil {
		h.logger.Error("Request validation failed", "error", err)
		return h.badRequestResponse(err.Error())
	}

	// 4. Executar use case
	ucInput := usecase.AuthenticateInput{CPF: authReq.CPF}
	ucOutput, err := h.authUC.Execute(ctx, ucInput)
	if err != nil {
		h.logger.Error("Authentication failed", "cpf", authReq.CPF, "error", err)
		return h.errorMapper.Map(err)
	}

	// 5. Criar response
	authResp := response.NewAuthResponse(
		ucOutput.Token,
		ucOutput.ExpiresIn,
		ucOutput.ClienteID,
		ucOutput.Nome,
	)

	h.logger.Info("Authentication successful", "cpf", authReq.CPF, "duration", time.Since(start))
	return h.successResponse(authResp)
}

// successResponse cria response de sucesso
func (h *AuthHandler) successResponse(data interface{}) HTTPResponse {
	body, err := json.Marshal(data)
	if err != nil {
		h.logger.Error("Failed to marshal response", "error", err)
		return h.internalErrorResponse()
	}

	return NewHTTPResponse(http.StatusOK, body)
}

// badRequestResponse cria response 400
func (h *AuthHandler) badRequestResponse(message string) HTTPResponse {
	errorResp := response.NewErrorResponse(message)
	body, _ := json.Marshal(errorResp)
	return NewHTTPResponse(http.StatusBadRequest, body)
}

// methodNotAllowedResponse cria response 405
func (h *AuthHandler) methodNotAllowedResponse() HTTPResponse {
	errorResp := response.NewErrorResponse("Method not allowed")
	body, _ := json.Marshal(errorResp)
	return NewHTTPResponse(http.StatusMethodNotAllowed, body)
}

// internalErrorResponse cria response 500
func (h *AuthHandler) internalErrorResponse() HTTPResponse {
	errorResp := response.NewErrorResponse("Internal server error")
	body, _ := json.Marshal(errorResp)
	return NewHTTPResponse(http.StatusInternalServerError, body)
}
