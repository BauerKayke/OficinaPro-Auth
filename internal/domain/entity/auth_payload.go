package entity

// AuthPayload representa o payload JWT de autenticação de forma tipada.
// Usar struct ao invés de map[string]interface{} reduz alocações e melhora performance.
type AuthPayload struct {
	ClienteID int64  `json:"cliente_id"` // ID do cliente
	Nome      string `json:"nome"`       // Nome do cliente
	Email     string `json:"email"`      // Email do cliente
	Documento string `json:"documento"`  // CPF do cliente
}

// ToMap converte AuthPayload para map (apenas quando necessário para JWT library).
// Minimiza uso de maps mantendo type-safety no domínio.
func (p *AuthPayload) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"cliente_id": p.ClienteID,
		"nome":       p.Nome,
		"email":      p.Email,
		"documento":  p.Documento,
	}
}
