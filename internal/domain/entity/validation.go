package entity

import (
	"strings"
)

// ValidateEmail valida formato básico de email.
// Retorna ErrInvalidEmail se formato for inválido.
func ValidateEmail(email string) error {
	if email == "" {
		return ErrInvalidEmail
	}

	// Email muito curto (ex: a@b.c tem 5 chars min)
	if len(email) < 5 {
		return ErrInvalidEmail
	}

	// Verificar se tem @ e dividir email em partes
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ErrInvalidEmail // Sem @ ou mais de um @
	}

	localPart := parts[0]
	domain := parts[1]

	// Local part não pode ser vazio
	if localPart == "" {
		return ErrInvalidEmail
	}

	// Domínio deve ter pelo menos um ponto e não pode estar vazio
	if domain == "" || !strings.Contains(domain, ".") {
		return ErrInvalidEmail
	}

	// Domínio não pode começar ou terminar com ponto
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return ErrInvalidEmail
	}

	return nil
}

// ValidatePasswordStrength valida requisitos mínimos de senha.
// Retorna ErrWeakPassword se senha não atender aos requisitos.
func ValidatePasswordStrength(password string) error {
	if len(password) < 6 {
		return ErrWeakPassword
	}
	return nil
}
