package database

import (
	"time"

	"gorm.io/gorm"
)

// Pessoa representa a tabela pessoa no banco
type Pessoa struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	Nome           string         `gorm:"type:varchar(255);not null"`
	Email          string         `gorm:"type:varchar(255)"`
	Documento      string         `gorm:"type:varchar(14);not null;uniqueIndex"`
	TipoPessoa     string         `gorm:"type:varchar(20);not null"`
	DataNascimento *time.Time     `gorm:"type:date"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	// Relacionamentos
	Cliente *Cliente `gorm:"foreignKey:PessoaID"`
}

// TableName especifica o nome da tabela
func (Pessoa) TableName() string {
	return "pessoa"
}

// Cliente representa a tabela cliente no banco
type Cliente struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	PessoaID        int64          `gorm:"not null;uniqueIndex"`
	Ativo           bool           `gorm:"default:true;index"`
	DataCriacao     time.Time      `gorm:"column:data_criacao;autoCreateTime"`
	DataAtualizacao time.Time      `gorm:"column:data_atualizacao;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relacionamentos
	Pessoa Pessoa `gorm:"foreignKey:PessoaID;references:ID"`
}

// TableName especifica o nome da tabela
func (Cliente) TableName() string {
	return "cliente"
}

// BeforeCreate hook executado antes de criar
func (c *Cliente) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	c.DataCriacao = now
	c.DataAtualizacao = now
	return nil
}

// BeforeUpdate hook executado antes de atualizar
func (c *Cliente) BeforeUpdate(tx *gorm.DB) error {
	c.DataAtualizacao = time.Now()
	return nil
}
