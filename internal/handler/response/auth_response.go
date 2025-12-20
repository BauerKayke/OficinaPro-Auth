package response

// AuthResponse representa a resposta de autenticação bem-sucedida.
// Contém o token JWT e informações básicas do cliente autenticado.
type AuthResponse struct {
	Token     string `json:"token"`     // Token JWT para autenticação em outras APIs
	ExpiresIn int    `json:"expiresIn"` // Tempo de expiração em segundos
	ClienteID int64  `json:"clienteId"` // ID único do cliente
	Nome      string `json:"nome"`      // Nome do cliente autenticado
}

// NewAuthResponse cria uma nova resposta de autenticação.
// Factory method para garantir criação consistente da resposta.
func NewAuthResponse(token string, expiresIn int, clienteID int64, nome string) *AuthResponse {
	return &AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		ClienteID: clienteID,
		Nome:      nome,
	}
}
