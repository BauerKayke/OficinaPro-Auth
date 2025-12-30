package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/oficinapro/auth-service/internal/domain/entity"
)

// GormUsuarioRepository implementa UsuarioRepository usando GORM ORM.
// Adaptador entre domínio (entity.Usuario) e infraestrutura (GORM).
type GormUsuarioRepository struct {
	db *gorm.DB
}

// NewGormUsuarioRepository cria nova instância do repository.
func NewGormUsuarioRepository(db *gorm.DB) *GormUsuarioRepository {
	return &GormUsuarioRepository{
		db: db,
	}
}

// UsuarioModel representa a tabela usuario no banco de dados.
// Alinhado com schema do core-domain-service (Java/Spring).
type UsuarioModel struct {
	ID              int64      `gorm:"column:id;primaryKey;autoIncrement"`
	Email           string     `gorm:"column:email;unique;not null"`
	Senha           string     `gorm:"column:senha;not null"`
	Role            string     `gorm:"column:role;not null;default:USER"`
	PessoaID        *int64     `gorm:"column:pessoa_id"`
	DataCriacao     time.Time  `gorm:"column:data_criacao;autoCreateTime"`
	DataAtualizacao time.Time  `gorm:"column:data_atualizacao;autoUpdateTime"`
	DataDesativacao *time.Time `gorm:"column:data_desativacao"`
	IsAtivo         bool       `gorm:"column:is_ativo;default:true"`
}

// TableName especifica nome da tabela no banco (convenção GORM).
func (UsuarioModel) TableName() string {
	return "usuario"
}

// FindByEmail busca usuário por email.
func (r *GormUsuarioRepository) FindByEmail(ctx context.Context, email string) (*entity.Usuario, error) {
	var model UsuarioModel

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&model).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUsuarioNotFound
		}
		return nil, err
	}

	return toUsuarioEntity(&model), nil
}

// FindByID busca usuário por ID.
func (r *GormUsuarioRepository) FindByID(ctx context.Context, id int64) (*entity.Usuario, error) {
	var model UsuarioModel

	err := r.db.WithContext(ctx).
		First(&model, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, entity.ErrUsuarioNotFound
		}
		return nil, err
	}

	return toUsuarioEntity(&model), nil
}

// Create cria novo usuário.
func (r *GormUsuarioRepository) Create(ctx context.Context, usuario *entity.Usuario) error {
	model := toUsuarioModel(usuario)

	// Verificar se email já existe
	var count int64
	r.db.WithContext(ctx).Model(&UsuarioModel{}).Where("email = ?", usuario.Email).Count(&count)
	if count > 0 {
		return entity.ErrEmailAlreadyExists
	}

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return err
	}

	// Atualizar ID gerado
	usuario.ID = model.ID
	return nil
}

// Update atualiza usuário existente.
func (r *GormUsuarioRepository) Update(ctx context.Context, usuario *entity.Usuario) error {
	model := toUsuarioModel(usuario)

	result := r.db.WithContext(ctx).
		Model(&UsuarioModel{}).
		Where("id = ?", usuario.ID).
		Updates(model)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrUsuarioNotFound
	}

	return nil
}

// Deactivate marca usuário como inativo (soft delete).
func (r *GormUsuarioRepository) Deactivate(ctx context.Context, id int64) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&UsuarioModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_ativo":         false,
			"data_desativacao": now,
			"data_atualizacao": now,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return entity.ErrUsuarioNotFound
	}

	return nil
}

// toUsuarioEntity converte model GORM para entidade de domínio.
func toUsuarioEntity(model *UsuarioModel) *entity.Usuario {
	return &entity.Usuario{
		ID:              model.ID,
		Email:           model.Email,
		Senha:           model.Senha,
		Role:            model.Role,
		PessoaID:        model.PessoaID,
		DataCriacao:     model.DataCriacao,
		DataAtualizacao: model.DataAtualizacao,
		DataDesativacao: model.DataDesativacao,
		IsAtivo:         model.IsAtivo,
	}
}

// toUsuarioModel converte entidade de domínio para model GORM.
func toUsuarioModel(usuario *entity.Usuario) *UsuarioModel {
	return &UsuarioModel{
		ID:              usuario.ID,
		Email:           usuario.Email,
		Senha:           usuario.Senha,
		Role:            usuario.Role,
		PessoaID:        usuario.PessoaID,
		DataCriacao:     usuario.DataCriacao,
		DataAtualizacao: usuario.DataAtualizacao,
		DataDesativacao: usuario.DataDesativacao,
		IsAtivo:         usuario.IsAtivo,
	}
}
