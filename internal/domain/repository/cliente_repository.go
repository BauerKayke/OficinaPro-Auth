package repository

import (
	"context"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// ClienteRepository define o contrato para acesso a dados de clientes
// Interface (Port) - Dependency Inversion Principle
type ClienteRepository interface {
	// FindByCPF busca um cliente pelo CPF
	FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error)

	// HealthCheck verifica se a conexão está saudável
	HealthCheck(ctx context.Context) error
}
