package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// Benchmark autenticação completa (cenário de sucesso)
func BenchmarkAuthenticateUseCase_Success(b *testing.B) {
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)
	mockTelemetry := new(MockTelemetryService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		mockTelemetry,
		24*time.Hour,
	)

	cliente := &entity.Cliente{
		ID:         1,
		Nome:       "João Silva",
		Email:      "joao@example.com",
		Documento:  "12345678909",
		TipoPessoa: "FISICA",
		Ativo:      true,
	}

	mockTelemetry.On("StartSpan", mock.Anything, mock.Anything).Return(context.Background(), &MockSpan{})
	mockTelemetry.On("IncrementCounter", mock.Anything, mock.Anything).Return()
	mockTelemetry.On("RecordDuration", mock.Anything, mock.Anything, mock.Anything).Return()
	mockValidator.On("NormalizeCPF", mock.Anything).Return("12345678909")
	mockValidator.On("ValidateCPF", mock.Anything).Return(true)
	mockRepo.On("FindByCPF", mock.Anything, mock.Anything).Return(cliente, nil)
	mockJWT.On("GenerateToken", mock.Anything, mock.Anything).Return("fake-token", nil)

	input := AuthenticateInput{CPF: "123.456.789-09"}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.Execute(ctx, input)
	}
}

// Benchmark cenário de CPF inválido (falha rápida)
func BenchmarkAuthenticateUseCase_InvalidCPF(b *testing.B) {
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)
	mockTelemetry := new(MockTelemetryService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		mockTelemetry,
		24*time.Hour,
	)

	mockTelemetry.On("StartSpan", mock.Anything, mock.Anything).Return(context.Background(), &MockSpan{})
	mockTelemetry.On("IncrementCounter", mock.Anything, mock.Anything).Return()
	mockValidator.On("NormalizeCPF", mock.Anything).Return("11111111111")
	mockValidator.On("ValidateCPF", mock.Anything).Return(false)

	input := AuthenticateInput{CPF: "111.111.111-11"}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.Execute(ctx, input)
	}
}

// Benchmark com diferentes tamanhos de documento
func BenchmarkAuthenticateUseCase_ParallelExecution(b *testing.B) {
	mockRepo := new(MockClienteRepository)
	mockJWT := new(MockJWTService)
	mockValidator := new(MockValidatorService)
	mockTelemetry := new(MockTelemetryService)

	useCase := NewAuthenticateUseCase(
		mockRepo,
		mockJWT,
		mockValidator,
		mockTelemetry,
		24*time.Hour,
	)

	cliente := &entity.Cliente{
		ID:         1,
		Nome:       "João Silva",
		Email:      "joao@example.com",
		Documento:  "12345678909",
		TipoPessoa: "FISICA",
		Ativo:      true,
	}

	mockTelemetry.On("StartSpan", mock.Anything, mock.Anything).Return(context.Background(), &MockSpan{})
	mockTelemetry.On("IncrementCounter", mock.Anything, mock.Anything).Return()
	mockTelemetry.On("RecordDuration", mock.Anything, mock.Anything, mock.Anything).Return()
	mockValidator.On("NormalizeCPF", mock.Anything).Return("12345678909")
	mockValidator.On("ValidateCPF", mock.Anything).Return(true)
	mockRepo.On("FindByCPF", mock.Anything, mock.Anything).Return(cliente, nil)
	mockJWT.On("GenerateToken", mock.Anything, mock.Anything).Return("fake-token", nil)

	input := AuthenticateInput{CPF: "123.456.789-09"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		for pb.Next() {
			_, _ = useCase.Execute(ctx, input)
		}
	})
}
