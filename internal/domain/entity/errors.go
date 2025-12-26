package entity

import "errors"

// Domain errors - Erros de negócio consolidados

// Erros de Usuario
var (
	ErrUsuarioNotFound    = errors.New("usuário não encontrado")
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrUsuarioInactive    = errors.New("usuário inativo")
	ErrInvalidEmail       = errors.New("email inválido")
	ErrWeakPassword       = errors.New("senha muito fraca")
	ErrEmailAlreadyExists = errors.New("email já cadastrado")
)

// Erros de Cliente (legado)
var (
	ErrClienteNotFound = errors.New("cliente não encontrado")
	ErrClienteInativo  = errors.New("cliente inativo")
	ErrInvalidDocument = errors.New("documento inválido")
	ErrInvalidCPF      = errors.New("CPF inválido")
)
