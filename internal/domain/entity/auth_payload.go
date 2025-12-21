package entity

// AuthPayload representa o payload JWT de autenticação de forma tipada.
// Usar struct ao invés de map[string]interface{} reduz alocações e melhora performance.
// Alinhado com JWT claims do core-domain-service (Java).
type AuthPayload struct {
	UserID int64  `json:"userId"`         // ID do usuário (mudou de ClienteID para UserID)
	Email  string `json:"email"`          // Email do usuário
	Nome   string `json:"nome,omitempty"` // Nome do usuário (opcional)
	Role   string `json:"role"`           // Role do usuário (USER, ADMIN, etc)
}

// ToMap converte AuthPayload para map (apenas quando necessário para JWT library).
// Minimiza uso de maps mantendo type-safety no domínio.
func (p *AuthPayload) ToMap() map[string]interface{} {
	claims := map[string]interface{}{
		"userId": p.UserID,
		"email":  p.Email,
		"role":   p.Role,
	}

	// Adiciona nome apenas se não vazio
	if p.Nome != "" {
		claims["nome"] = p.Nome
	}

	return claims
}
