// Package usecase contém os casos de uso da aplicação (camada de aplicação).
// Orquestra a lógica de negócio usando entidades de domínio e serviços.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
)

// AuthenticateInput representa os dados de entrada para autenticação.
type AuthenticateInput struct {
	CPF string `json:"cpf"`
}

// AuthenticateOutput representa os dados de saída após autenticação bem-sucedida.
type AuthenticateOutput struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
	ClienteID int64  `json:"clienteId"`
	Nome      string `json:"nome"`
}

// AuthenticateUseCase orquestra o processo de autenticação de clientes via CPF.
// Responsável por validar CPF, buscar cliente, verificar permissões e gerar JWT.
type AuthenticateUseCase struct {
	clienteRepo      repository.ClienteRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	telemetryService service.TelemetryService
	jwtExpiration    time.Duration
}

// NewAuthenticateUseCase cria uma nova instância do caso de uso de autenticação.
// Todas as dependências são injetadas via parâmetros (Dependency Inversion Principle).
func NewAuthenticateUseCase(
	clienteRepo repository.ClienteRepository,
	jwtService service.JWTService,
	validatorService service.ValidatorService,
	telemetryService service.TelemetryService,
	jwtExpiration time.Duration,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		clienteRepo:      clienteRepo,
		jwtService:       jwtService,
		validatorService: validatorService,
		telemetryService: telemetryService,
		jwtExpiration:    jwtExpiration,
	}
}

// Execute executa o caso de uso de autenticação.
// Valida o CPF, busca o cliente, verifica se pode autenticar e gera o token JWT.
// Retorna AuthenticateOutput com o token e dados do cliente ou erro de domínio.
func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
	// Span principal do caso de uso
	ctx, span := uc.telemetryService.StartSpan(ctx, "authenticate.execute")
	defer span.End()

	startTime := time.Now()

	// Normaliza e valida CPF
	normalizedCPF := uc.validatorService.NormalizeCPF(input.CPF)
	if !uc.validatorService.ValidateCPF(normalizedCPF) {
		span.SetStatus(entity.ErrInvalidCPF)
		uc.recordMetric("error", "invalid_cpf")
		return nil, entity.ErrInvalidCPF
	}

	// Busca cliente no repositório
	cliente, err := uc.clienteRepo.FindByCPF(ctx, normalizedCPF)
	if err != nil {
		span.SetStatus(err)
		uc.recordMetric("error", "repository_error")
		return nil, fmt.Errorf("failed to find cliente: %w", err)
	}

	if cliente == nil {
		span.SetStatus(entity.ErrClienteNotFound)
		uc.recordMetric("error", "not_found")
		return nil, entity.ErrClienteNotFound
	}

	// Verifica se cliente pode autenticar
	if err := cliente.CanAuthenticate(); err != nil {
		span.SetStatus(err)
		uc.recordMetric("error", "not_allowed")
		return nil, err
	}

	// Gera token JWT
	payload := cliente.ToAuthPayload()
	token, err := uc.jwtService.GenerateToken(payload.ToMap(), uc.jwtExpiration)
	if err != nil {
		span.SetStatus(err)
		uc.recordMetric("error", "token_failed")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Registrar sucesso
	uc.recordMetric("success", "")
	uc.telemetryService.RecordDuration("authenticate.duration", time.Since(startTime), nil)

	return &AuthenticateOutput{
		Token:     token,
		ExpiresIn: int(uc.jwtExpiration.Seconds()),
		ClienteID: cliente.ID,
		Nome:      cliente.Nome,
	}, nil
}

// recordMetric registra métrica simples de resultado (success/error).
func (uc *AuthenticateUseCase) recordMetric(status, errorType string) {
	attrs := map[string]interface{}{"status": status}
	if errorType != "" {
		attrs["error_type"] = errorType
	}
	uc.telemetryService.IncrementCounter("authenticate.requests", attrs)
}
