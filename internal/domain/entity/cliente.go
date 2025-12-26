package entity

import (
	"time"
)

// Cliente representa um cliente no domínio da Oficina Pro (LEGADO).
// NOTA: Cliente é mantido apenas para compatibilidade com sistema legado CPF.
// Para novos fluxos de autenticação, use Usuario (email+senha).
type Cliente struct {
	ID          int64     // Identificador único do cliente
	Nome        string    // Nome completo do cliente
	Email       string    // Email do cliente
	Documento   string    // CPF do cliente (normalizado, apenas números)
	TipoPessoa  string    // Tipo de pessoa (física/jurídica)
	Ativo       bool      // Status de ativação do cliente
	DataCriacao time.Time // Data de criação do registro
}

// IsValid verifica se o cliente é válido para autenticação.
func (c *Cliente) IsValid() bool {
	return c.ID > 0 && c.Ativo && c.Documento != ""
}

// CanAuthenticate verifica se o cliente pode se autenticar.
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

// ToAuthPayload converte cliente para payload de autenticação tipado.
// DEPRECATED: Mantido apenas para compatibilidade com sistema legado CPF.
// Para novos fluxos de auth, use Usuario.ToAuthPayload().
func (c *Cliente) ToAuthPayload() *AuthPayload {
	return &AuthPayload{
		UserID: c.ID,
		Email:  c.Email,
		Nome:   c.Nome,
		Role:   RoleUser, // Clientes legados sempre tem role USER
	}
}
