package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// PostgresClienteRepository implementa ClienteRepository para PostgreSQL
// Adapter - implementa interface do domínio
type PostgresClienteRepository struct {
	db *sql.DB
}

// NewPostgresClienteRepository cria uma nova instância do repositório
func NewPostgresClienteRepository(db *sql.DB) *PostgresClienteRepository {
	return &PostgresClienteRepository{
		db: db,
	}
}

// FindByCPF busca um cliente pelo CPF
func (r *PostgresClienteRepository) FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error) {
	query := `
		SELECT
			c.id,
			c.ativo,
			p.nome,
			p.email,
			p.documento,
			p.tipo_pessoa,
			c.data_criacao
		FROM cliente c
		INNER JOIN pessoa p ON c.pessoa_id = p.id
		WHERE p.documento = $1
		LIMIT 1
	`

	log.Printf("[Repository] Buscando cliente com CPF: %s", maskCPF(cpf))

	var cliente entity.Cliente
	err := r.db.QueryRowContext(ctx, query, cpf).Scan(
		&cliente.ID,
		&cliente.Ativo,
		&cliente.Nome,
		&cliente.Email,
		&cliente.Documento,
		&cliente.TipoPessoa,
		&cliente.DataCriacao,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[Repository] Cliente não encontrado para CPF: %s", maskCPF(cpf))
			return nil, entity.ErrClienteNotFound
		}
		log.Printf("[Repository] Erro ao buscar cliente: %v", err)
		return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	log.Printf("[Repository] Cliente encontrado: ID=%d, Ativo=%v", cliente.ID, cliente.Ativo)
	return &cliente, nil
}

// HealthCheck verifica se a conexão com o banco está saudável
func (r *PostgresClienteRepository) HealthCheck(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

// maskCPF mascara o CPF para logs
func maskCPF(cpf string) string {
	if len(cpf) < 11 {
		return "***"
	}
	return cpf[:3] + "***" + cpf[len(cpf)-2:]
}
