package handler_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/handler"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// Simple mocks
type SimpleUsuarioRepo struct {
	Users map[string]*entity.Usuario
}

func (r *SimpleUsuarioRepo) FindByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	if user, ok := r.Users[email]; ok {
		return user, nil
	}
	return nil, entity.ErrUsuarioNotFound
}

func (r *SimpleUsuarioRepo) FindByID(ctx context.Context, id int64) (*entity.Usuario, error) {
	return nil, entity.ErrUsuarioNotFound
}

func (r *SimpleUsuarioRepo) Create(ctx context.Context, usuario *entity.Usuario) error {
	return nil
}

func (r *SimpleUsuarioRepo) Update(ctx context.Context, usuario *entity.Usuario) error {
	return nil
}

func (r *SimpleUsuarioRepo) Deactivate(ctx context.Context, id int64) error {
	return nil
}

type SimpleJWT struct{}

func (s *SimpleJWT) GenerateToken(payload map[string]interface{}, expiration time.Duration) (string, error) {
	return "test-token", nil
}

func (s *SimpleJWT) ValidateToken(token string) (map[string]interface{}, error) {
	return map[string]interface{}{"user_id": 1}, nil
}

type NoOpTelemetry struct{}

func (n *NoOpTelemetry) StartSpan(ctx context.Context, name string) (context.Context, service.Span) {
	return ctx, &NoOpSpan{}
}

func (n *NoOpTelemetry) RecordMetric(name string, value float64, attributes map[string]interface{}) {}

func (n *NoOpTelemetry) RecordDuration(name string, duration time.Duration, attributes map[string]interface{}) {
}

func (n *NoOpTelemetry) IncrementCounter(name string, attributes map[string]interface{}) {}

func (n *NoOpTelemetry) Shutdown(ctx context.Context) error {
	return nil
}

func (n *NoOpTelemetry) ForceFlush(ctx context.Context) error {
	return nil
}

type NoOpSpan struct{}

func (n *NoOpSpan) SetAttribute(key string, value interface{})              {}
func (n *NoOpSpan) SetStatus(err error)                                     {}
func (n *NoOpSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (n *NoOpSpan) End()                                                    {}

type SimpleLogger struct{}

func (l *SimpleLogger) Info(msg string, fields ...interface{})  {}
func (l *SimpleLogger) Error(msg string, fields ...interface{}) {}
func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {}

// Ensure interfaces are satisfied
var _ repository.UsuarioRepository = (*SimpleUsuarioRepo)(nil)
var _ service.JWTService = (*SimpleJWT)(nil)
var _ service.TelemetryService = (*NoOpTelemetry)(nil)
var _ handler.Logger = (*SimpleLogger)(nil)

func setupHandler() *handler.AuthHandler {
	repo := &SimpleUsuarioRepo{
		Users: make(map[string]*entity.Usuario),
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test@123"), bcrypt.DefaultCost)
	repo.Users["test@example.com"] = &entity.Usuario{
		ID:              1,
		Email:           "test@example.com",
		Senha:           string(hashedPassword),
		Role:            entity.RoleUser,
		IsAtivo:         true,
		DataCriacao:     time.Now(),
		DataAtualizacao: time.Now(),
	}

	jwtService := &SimpleJWT{}
	telemetry := &NoOpTelemetry{}
	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)
	errorMapper := handler.NewDefaultErrorMapper()
	logger := &SimpleLogger{}

	return handler.NewAuthHandler(uc, jwtService, errorMapper, logger)
}

func TestAuthHandler_Handle_PostAuth_Success(t *testing.T) {
	h := setupHandler()

	reqBody := map[string]string{
		"email": "test@example.com",
		"senha": "Test@123",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := handler.HTTPRequest{
		Method: "POST",
		Path:   "/auth",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: bodyBytes,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, resp.Body)

	var authResp map[string]interface{}
	err := json.Unmarshal(resp.Body, &authResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, authResp["token"])
}

func TestAuthHandler_Handle_PostAuth_InvalidCredentials(t *testing.T) {
	h := setupHandler()

	reqBody := map[string]string{
		"email": "test@example.com",
		"senha": "WrongPassword",
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req := handler.HTTPRequest{
		Method: "POST",
		Path:   "/auth",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: bodyBytes,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 401, resp.StatusCode)
}

func TestAuthHandler_Handle_PostAuth_InvalidJSON(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method: "POST",
		Path:   "/auth",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: []byte("invalid json"),
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 400, resp.StatusCode)
}

func TestAuthHandler_Handle_GetHealth(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method:  "GET",
		Path:    "/health",
		Headers: map[string]string{},
		Body:    nil,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthHandler_Handle_Options(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method:  "OPTIONS",
		Path:    "/auth",
		Headers: map[string]string{},
		Body:    nil,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Headers, "Access-Control-Allow-Origin")
}

func TestAuthHandler_Handle_NotFound(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method:  "GET",
		Path:    "/non-existent",
		Headers: map[string]string{},
		Body:    nil,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestAuthHandler_Handle_PostValidate_Success(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method: "POST",
		Path:   "/auth/validate",
		Headers: map[string]string{
			"Authorization": "Bearer test-token",
		},
		Body: nil,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 200, resp.StatusCode)
}

func TestAuthHandler_Handle_PostValidate_NoToken(t *testing.T) {
	h := setupHandler()

	req := handler.HTTPRequest{
		Method:  "POST",
		Path:    "/auth/validate",
		Headers: map[string]string{},
		Body:    nil,
	}

	resp := h.Handle(context.Background(), req)

	assert.Equal(t, 401, resp.StatusCode)
}
