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
// Mudou de CPF para Email+Senha (alinhado com core-domain-service Java).
type AuthenticateInput struct {
	Email string
	Senha string
}

// AuthenticateOutput representa os dados de saída após autenticação bem-sucedida.
type AuthenticateOutput struct {
	Token     string
	ExpiresIn int
	UserID    int64
	Email     string
	Nome      string
	Role      string
}

// AuthenticateUseCase orquestra o processo de autenticação de usuários via Email+Senha.
// Responsável por validar credenciais, verificar permissões e gerar JWT.
type AuthenticateUseCase struct {
	usuarioRepo      repository.UsuarioRepository
	jwtService       service.JWTService
	telemetryService service.TelemetryService
	jwtExpiration    time.Duration
}

// NewAuthenticateUseCase cria uma nova instância do caso de uso de autenticação.
// Todas as dependências são injetadas via parâmetros (Dependency Inversion Principle).
func NewAuthenticateUseCase(
	usuarioRepo repository.UsuarioRepository,
	jwtService service.JWTService,
	telemetryService service.TelemetryService,
	jwtExpiration time.Duration,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		usuarioRepo:      usuarioRepo,
		jwtService:       jwtService,
		telemetryService: telemetryService,
		jwtExpiration:    jwtExpiration,
	}
}

// Execute executa o caso de uso de autenticação.
// Valida email+senha, verifica se usuário está ativo e gera o token JWT.
// Retorna AuthenticateOutput com o token e dados do usuário ou erro de domínio.
func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
	// Span principal do caso de uso
	ctx, span := uc.telemetryService.StartSpan(ctx, "authenticate.execute")
	defer span.End()

	// 1. Buscar usuário por email
	usuario, err := uc.usuarioRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		span.SetStatus(err)
		uc.recordMetric("error", "usuario_not_found")
		return nil, entity.ErrUsuarioNotFound // Não expor se é email ou senha
	}

	// 2. Validar senha usando bcrypt
	if !usuario.ValidatePassword(input.Senha) {
		span.SetStatus(entity.ErrInvalidCredentials)
		uc.recordMetric("error", "invalid_password")
		return nil, entity.ErrInvalidCredentials
	}

	// 3. Verificar se usuário está ativo
	if !usuario.IsEnabled() {
		span.SetStatus(entity.ErrUsuarioInactive)
		uc.recordMetric("error", "usuario_inactive")
		return nil, entity.ErrUsuarioInactive
	}

	// 4. Preparar payload JWT
	payload := entity.AuthPayload{
		UserID: usuario.ID,
		Email:  usuario.Email,
		Nome:   usuario.GetDisplayName(),
		Role:   usuario.Role,
	}

	// 5. Gerar token JWT
	token, err := uc.jwtService.GenerateToken(payload.ToMap(), uc.jwtExpiration)
	if err != nil {
		span.SetStatus(err)
		uc.recordMetric("error", "token_generation_failed")
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// 6. Registrar sucesso
	uc.recordMetric("success", "authenticated")
	span.SetAttribute("user_id", usuario.ID)
	span.SetAttribute("role", usuario.Role)

	// 7. Retornar output
	return &AuthenticateOutput{
		Token:     token,
		ExpiresIn: int(uc.jwtExpiration.Seconds()),
		UserID:    usuario.ID,
		Email:     usuario.Email,
		Nome:      usuario.GetDisplayName(),
		Role:      usuario.Role,
	}, nil
}

// recordMetric registra métrica de autenticação (sucesso/erro).
func (uc *AuthenticateUseCase) recordMetric(status, reason string) {
	// Converter tags para format esperado pelo telemetry service
	tags := map[string]interface{}{
		"status": status,
		"reason": reason,
	}
	uc.telemetryService.IncrementCounter("authenticate.attempts", tags)
}
