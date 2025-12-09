package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// Mock ClienteRepository
type MockClienteRepository struct {
	mock.Mock
}

func (m *MockClienteRepository) FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error) {
	args := m.Called(ctx, cpf)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Cliente), args.Error(1)
}

func (m *MockClienteRepository) HealthCheck(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Mock JWTService
type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) GenerateToken(payload map[string]interface{}, expiresIn time.Duration) (string, error) {
	args := m.Called(payload, expiresIn)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock ValidatorService
type MockValidatorService struct {
	mock.Mock
}

func (m *MockValidatorService) ValidateCPF(cpf string) bool {
	args := m.Called(cpf)
	return args.Bool(0)
}

func (m *MockValidatorService) NormalizeCPF(cpf string) string {
	args := m.Called(cpf)
	return args.String(0)
}

func TestAuthenticateUseCase_Execute_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		24*time.Hour,
	)

	ctx := context.Background()
	input := AuthenticateInput{CPF: "123.456.789-09"}

	cliente := &entity.Cliente{
		ID:         1,
		Nome:       "João Silva",
		Email:      "joao@example.com",
		Documento:  "12345678909",
		TipoPessoa: "FISICA",
		Ativo:      true,
	}

	// Setup mocks
	mockValidator.On("NormalizeCPF", input.CPF).Return("12345678909")
	mockValidator.On("ValidateCPF", "12345678909").Return(true)
	mockRepo.On("FindByCPF", ctx, "12345678909").Return(cliente, nil)
	mockJWT.On("GenerateToken", mock.Anything, mock.Anything).Return("fake-token", nil)

	// Act
	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "fake-token", output.Token)
	assert.Equal(t, int64(1), output.ClienteID)
	assert.Equal(t, "João Silva", output.Nome)
	assert.Equal(t, 86400, output.ExpiresIn) // 24 hours in seconds

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestAuthenticateUseCase_Execute_InvalidCPF(t *testing.T) {
	// Arrange
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		24*time.Hour,
	)

	ctx := context.Background()
	input := AuthenticateInput{CPF: "111.111.111-11"}

	// Setup mocks
	mockValidator.On("NormalizeCPF", input.CPF).Return("11111111111")
	mockValidator.On("ValidateCPF", "11111111111").Return(false)

	// Act
	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrInvalidCPF, err)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "FindByCPF")
	mockJWT.AssertNotCalled(t, "GenerateToken")
}

func TestAuthenticateUseCase_Execute_ClienteNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		24*time.Hour,
	)

	ctx := context.Background()
	input := AuthenticateInput{CPF: "123.456.789-09"}

	// Setup mocks
	mockValidator.On("NormalizeCPF", input.CPF).Return("12345678909")
	mockValidator.On("ValidateCPF", "12345678909").Return(true)
	mockRepo.On("FindByCPF", ctx, "12345678909").Return(nil, entity.ErrClienteNotFound)

	// Act
	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "erro ao buscar cliente")

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertNotCalled(t, "GenerateToken")
}

func TestAuthenticateUseCase_Execute_ClienteInativo(t *testing.T) {
	// Arrange
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		24*time.Hour,
	)

	ctx := context.Background()
	input := AuthenticateInput{CPF: "123.456.789-09"}

	cliente := &entity.Cliente{
		ID:         1,
		Nome:       "João Silva",
		Email:      "joao@example.com",
		Documento:  "12345678909",
		TipoPessoa: "FISICA",
		Ativo:      false, // Cliente inativo
	}

	// Setup mocks
	mockValidator.On("NormalizeCPF", input.CPF).Return("12345678909")
	mockValidator.On("ValidateCPF", "12345678909").Return(true)
	mockRepo.On("FindByCPF", ctx, "12345678909").Return(cliente, nil)

	// Act
	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Equal(t, entity.ErrClienteInativo, err)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertNotCalled(t, "GenerateToken")
}

func TestAuthenticateUseCase_Execute_JWTGenerationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		24*time.Hour,
	)

	ctx := context.Background()
	input := AuthenticateInput{CPF: "123.456.789-09"}

	cliente := &entity.Cliente{
		ID:         1,
		Nome:       "João Silva",
		Email:      "joao@example.com",
		Documento:  "12345678909",
		TipoPessoa: "FISICA",
		Ativo:      true,
	}

	// Setup mocks
	mockValidator.On("NormalizeCPF", input.CPF).Return("12345678909")
	mockValidator.On("ValidateCPF", "12345678909").Return(true)
	mockRepo.On("FindByCPF", ctx, "12345678909").Return(cliente, nil)
	mockJWT.On("GenerateToken", mock.Anything, mock.Anything).
		Return("", errors.New("JWT generation failed"))

	// Act
	output, err := useCase.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, output)
	assert.Contains(t, err.Error(), "erro ao gerar token")

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}
