package response

// AuthResponse representa o response de autenticação bem-sucedida
type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
	ClienteID int64  `json:"clienteId"`
	Nome      string `json:"nome"`
}

// NewAuthResponse cria um novo AuthResponse
func NewAuthResponse(token string, expiresIn int, clienteID int64, nome string) *AuthResponse {
	return &AuthResponse{
		Token:     token,
		ExpiresIn: expiresIn,
		ClienteID: clienteID,
		Nome:      nome,
	}
}
