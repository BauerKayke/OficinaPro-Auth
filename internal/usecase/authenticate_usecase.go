package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/oficinapro/auth-service/internal/domain/entity"
	"github.com/oficinapro/auth-service/internal/domain/repository"
	"github.com/oficinapro/auth-service/internal/domain/service"
)

type AuthenticateInput struct {
	CPF string `json:"cpf"`
}

type AuthenticateOutput struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
	ClienteID int64  `json:"clienteId"`
	Nome      string `json:"nome"`
}

type AuthenticateUseCase struct {
	clienteRepo      repository.ClienteRepository
	jwtService       service.JWTService
	validatorService service.ValidatorService
	jwtExpiration    time.Duration
}

func NewAuthenticateUseCase(
	clienteRepo repository.ClienteRepository,
	jwtService service.JWTService,
	validatorService service.ValidatorService,
	jwtExpiration time.Duration,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		clienteRepo:      clienteRepo,
		jwtService:       jwtService,
		validatorService: validatorService,
		jwtExpiration:    jwtExpiration,
	}
}

func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
	normalizedCPF := uc.validatorService.NormalizeCPF(input.CPF)

	if !uc.validatorService.ValidateCPF(normalizedCPF) {
		return nil, entity.ErrInvalidCPF
	}

	cliente, err := uc.clienteRepo.FindByCPF(ctx, normalizedCPF)
	if err != nil {
		return nil, fmt.Errorf("failed to find cliente: %w", err)
	}

	if cliente == nil {
		return nil, entity.ErrClienteNotFound
	}

	if err := cliente.CanAuthenticate(); err != nil {
		return nil, err
	}

	payload := cliente.ToAuthPayload()
	token, err := uc.jwtService.GenerateToken(payload, uc.jwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	log.Printf("[Auth] Success: cliente_id=%d", cliente.ID)

	return &AuthenticateOutput{
		Token:     token,
		ExpiresIn: int(uc.jwtExpiration.Seconds()),
		ClienteID: cliente.ID,
		Nome:      cliente.Nome,
	}, nil
}
