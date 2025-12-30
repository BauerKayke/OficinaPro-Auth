package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// GormClienteRepository implementa ClienteRepository usando GORM.
// Responsável por persistência e recuperação de dados de clientes.
type GormClienteRepository struct {
	db *gorm.DB
}

// NewGormClienteRepository cria nova instância do repositório de clientes.
func NewGormClienteRepository(db *gorm.DB) *GormClienteRepository {
	return &GormClienteRepository{db: db}
}

// FindByCPF busca cliente por CPF no banco de dados.
// Usa prepared statement do GORM para prevenir SQL injection.
// Retorna nil se cliente não for encontrado, ou erro em caso de falha na query.
func (r *GormClienteRepository) FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error) {
	var cliente Cliente

	// Query com JOIN e prepared statement (proteção contra SQL injection)
	result := r.db.WithContext(ctx).
		Preload("Pessoa").
		Joins("JOIN pessoa ON pessoa.id = cliente.pessoa_id").
		Where("pessoa.documento = ?", cpf).
		First(&cliente)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, entity.ErrClienteNotFound
		}
		return nil, result.Error
	}

	// Mapeia modelo de banco para entidade de domínio
	return &entity.Cliente{
		ID:          cliente.ID,
		Nome:        cliente.Pessoa.Nome,
		Email:       cliente.Pessoa.Email,
		Documento:   cliente.Pessoa.Documento,
		TipoPessoa:  cliente.Pessoa.TipoPessoa,
		Ativo:       cliente.Ativo,
		DataCriacao: cliente.DataCriacao,
	}, nil
}

// HealthCheck verifica conectividade com o banco de dados.
// Útil para health checks e monitoramento da aplicação.
func (r *GormClienteRepository) HealthCheck(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
