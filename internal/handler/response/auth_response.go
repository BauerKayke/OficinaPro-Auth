package response

// AuthResponse representa a resposta de autenticação bem-sucedida.
// Contém o token JWT e informações básicas do usuário autenticado.
// Alinhado com AuthenticationResponse do core-domain-service (Java).
type AuthResponse struct {
	Token     string `json:"token"`          // Token JWT para autenticação em outras APIs
	ExpiresIn int    `json:"expiresIn"`      // Tempo de expiração em segundos
	UserID    int64  `json:"userId"`         // ID único do usuário
	Email     string `json:"email"`          // Email do usuário
	Nome      string `json:"nome,omitempty"` // Nome do usuário (se disponível)
	Role      string `json:"role"`           // Role do usuário (USER, ADMIN, etc)
}

// NewAuthResponse cria uma nova resposta de autenticação.
// Factory method para garantir criação consistente da resposta.
func NewAuthResponse(token string, expiresIn int, userID int64, email, nome, role string) *AuthResponse {
	return &AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		UserID:    userID,
		Email:     email,
		Nome:      nome,
		Role:      role,
	}
}
