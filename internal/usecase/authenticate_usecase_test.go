package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/usecase"
)

// Simple mocks without telemetry complexity
type SimpleUsuarioRepository struct {
	Users map[string]*entity.Usuario
}

func NewSimpleUsuarioRepository() *SimpleUsuarioRepository {
	return &SimpleUsuarioRepository{
		Users: make(map[string]*entity.Usuario),
	}
}

func (r *SimpleUsuarioRepository) FindByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	if user, ok := r.Users[email]; ok {
		return user, nil
	}
	return nil, entity.ErrUsuarioNotFound
}

func (r *SimpleUsuarioRepository) FindByID(ctx context.Context, id int64) (*entity.Usuario, error) {
	for _, user := range r.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, entity.ErrUsuarioNotFound
}

func (r *SimpleUsuarioRepository) Create(ctx context.Context, usuario *entity.Usuario) error {
	if _, exists := r.Users[usuario.Email]; exists {
		return entity.ErrEmailAlreadyExists
	}
	r.Users[usuario.Email] = usuario
	return nil
}

func (r *SimpleUsuarioRepository) Update(ctx context.Context, usuario *entity.Usuario) error {
	if _, exists := r.Users[usuario.Email]; !exists {
		return entity.ErrUsuarioNotFound
	}
	r.Users[usuario.Email] = usuario
	return nil
}

func (r *SimpleUsuarioRepository) Deactivate(ctx context.Context, id int64) error {
	for _, user := range r.Users {
		if user.ID == id {
			user.IsAtivo = false
			return nil
		}
	}
	return entity.ErrUsuarioNotFound
}

type SimpleJWTService struct{}

func (s *SimpleJWTService) GenerateToken(payload map[string]interface{}, expiration time.Duration) (string, error) {
	return "fake-jwt-token", nil
}

func (s *SimpleJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	return map[string]interface{}{"user_id": 1}, nil
}

type NoOpTelemetryService struct{}

func (n *NoOpTelemetryService) StartSpan(ctx context.Context, name string) (context.Context, service.Span) {
	return ctx, &NoOpSpan{}
}

func (n *NoOpTelemetryService) RecordMetric(name string, value float64, attributes map[string]interface{}) {
}

func (n *NoOpTelemetryService) RecordDuration(name string, duration time.Duration, attributes map[string]interface{}) {
}

func (n *NoOpTelemetryService) IncrementCounter(name string, attributes map[string]interface{}) {}

func (n *NoOpTelemetryService) Shutdown(ctx context.Context) error {
	return nil
}

type NoOpSpan struct{}

func (n *NoOpSpan) SetAttribute(key string, value interface{})              {}
func (n *NoOpSpan) SetStatus(err error)                                     {}
func (n *NoOpSpan) AddEvent(name string, attributes map[string]interface{}) {}
func (n *NoOpSpan) End()                                                    {}

// Helper to create test user
func createUser(email, password, role string) *entity.Usuario {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return &entity.Usuario{
		ID:              1,
		Email:           email,
		Senha:           string(hashedPassword),
		Role:            role,
		IsAtivo:         true,
		DataCriacao:     time.Now(),
		DataAtualizacao: time.Now(),
	}
}

// Make sure these interfaces are satisfied
var _ repository.UsuarioRepository = (*SimpleUsuarioRepository)(nil)
var _ service.JWTService = (*SimpleJWTService)(nil)
var _ service.TelemetryService = (*NoOpTelemetryService)(nil)

func TestAuthenticateUseCase_Success(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	user := createUser("user@test.com", "Password@123", entity.RoleUser)
	repo.Users["user@test.com"] = user

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "user@test.com",
		Senha: "Password@123",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "fake-jwt-token", output.Token)
	assert.Greater(t, output.ExpiresIn, 0)
}

func TestAuthenticateUseCase_UserNotFound(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "notfound@test.com",
		Senha: "Password@123",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrUsuarioNotFound, err)
}

func TestAuthenticateUseCase_InvalidPassword(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	user := createUser("user@test.com", "CorrectPassword", entity.RoleUser)
	repo.Users["user@test.com"] = user

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "user@test.com",
		Senha: "WrongPassword",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidCredentials, err)
}

func TestAuthenticateUseCase_InactiveUser(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	user := createUser("user@test.com", "Password@123", entity.RoleUser)
	user.IsAtivo = false
	repo.Users["user@test.com"] = user

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "user@test.com",
		Senha: "Password@123",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrUsuarioInactive, err)
}

func TestAuthenticateUseCase_EmptyEmail(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "",
		Senha: "Password@123",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
}

func TestAuthenticateUseCase_EmptyPassword(t *testing.T) {
	repo := NewSimpleUsuarioRepository()
	jwtService := &SimpleJWTService{}
	telemetry := &NoOpTelemetryService{}

	uc := usecase.NewAuthenticateUseCase(repo, jwtService, telemetry, time.Hour)

	input := usecase.AuthenticateInput{
		Email: "user@test.com",
		Senha: "",
	}

	output, err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	assert.Nil(t, output)
}
