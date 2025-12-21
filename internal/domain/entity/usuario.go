// Package entity contém as entidades de domínio da aplicação.
// Entidades representam conceitos fundamentais do negócio com identidade única.
package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Usuario representa um usuário do sistema com credenciais de autenticação.
// Alinhado com a entidade Usuario do core-domain-service (Java).
type Usuario struct {
	ID              int64      `json:"id"`
	Email           string     `json:"email"`
	Senha           string     `json:"-"` // Never expose in JSON
	Role            string     `json:"role"`
	PessoaID        *int64     `json:"pessoaId,omitempty"`
	DataCriacao     time.Time  `json:"dataCriacao"`
	DataAtualizacao time.Time  `json:"dataAtualizacao"`
	DataDesativacao *time.Time `json:"dataDesativacao,omitempty"`
	IsAtivo         bool       `json:"isAtivo"`
}

// Role constants (alinhado com enum Role do Java)
const (
	RoleUser  = "USER"
	RoleAdmin = "ADMIN"
)

// ValidatePassword verifica se a senha fornecida corresponde ao hash armazenado.
// Usa bcrypt para comparação segura.
func (u *Usuario) ValidatePassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(plainPassword))
	return err == nil
}

// IsEnabled verifica se o usuário está ativo e pode autenticar.
// Implementa mesma lógica do UserDetails.isEnabled() do Spring Security.
func (u *Usuario) IsEnabled() bool {
	return u.IsAtivo && u.DataDesativacao == nil
}

// GetDisplayName retorna nome para exibição (usa email como fallback).
func (u *Usuario) GetDisplayName() string {
	return u.Email
}

// ToAuthPayload converte Usuario para payload de autenticação tipado.
func (u *Usuario) ToAuthPayload() *AuthPayload {
	return &AuthPayload{
		UserID: u.ID,
		Email:  u.Email,
		Nome:   u.GetDisplayName(),
		Role:   u.Role,
	}
}

// HashPassword gera hash bcrypt de uma senha em texto plano.
// Usado para criar ou atualizar senhas de usuários.
func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
