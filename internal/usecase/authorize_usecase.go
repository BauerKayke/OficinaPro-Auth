package usecase

import (
	"fmt"
	"strings"

	"github.com/oficinapro/auth-service/internal/domain/service"
)

type AuthorizeUseCase struct {
	jwtService service.JWTService
}

func NewAuthorizeUseCase(jwtService service.JWTService) *AuthorizeUseCase {
	return &AuthorizeUseCase{
		jwtService: jwtService,
	}
}

// Execute valida o token e retorna true se válido, false caso contrário.
// O token pode vir com prefixo "Bearer ".
func (uc *AuthorizeUseCase) Execute(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, nil
	}

	// Remover prefixo Bearer se existir
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	claims, err := uc.jwtService.ValidateToken(tokenString)
	if err != nil {
		return false, fmt.Errorf("token validation failed: %w", err)
	}

	// Validação extra opcional: verificar se o CPF existe nos claims
	if _, ok := claims["cpf"]; !ok {
		return false, fmt.Errorf("token missing cpf claim")
	}

	return true, nil
}

