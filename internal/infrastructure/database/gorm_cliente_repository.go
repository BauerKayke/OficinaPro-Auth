package database

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

type GormClienteRepository struct {
	db *gorm.DB
}

func NewGormClienteRepository(db *gorm.DB) *GormClienteRepository {
	return &GormClienteRepository{db: db}
}

func (r *GormClienteRepository) FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error) {
	var cliente Cliente

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

func (r *GormClienteRepository) HealthCheck(ctx context.Context) error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
