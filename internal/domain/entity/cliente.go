package entity

import (
	"time"
)

// Cliente representa um cliente no domínio
type Cliente struct {
	ID          int64
	Nome        string
	Email       string
	Documento   string // CPF
	TipoPessoa  string
	Ativo       bool
	DataCriacao time.Time
}

// IsValid verifica se o cliente é válido para autenticação
func (c *Cliente) IsValid() bool {
	return c.ID > 0 && c.Ativo && c.Documento != ""
}

// CanAuthenticate verifica se o cliente pode se autenticar
func (c *Cliente) CanAuthenticate() error {
	if c.ID == 0 {
		return ErrClienteNotFound
	}

	if !c.Ativo {
		return ErrClienteInativo
	}

	if c.Documento == "" {
		return ErrInvalidDocument
	}

	return nil
}

// ToAuthPayload converte cliente para payload de autenticação
func (c *Cliente) ToAuthPayload() map[string]interface{} {
	return map[string]interface{}{
		"cliente_id": c.ID,
		"nome":       c.Nome,
		"email":      c.Email,
		"documento":  c.Documento,
	}
}
