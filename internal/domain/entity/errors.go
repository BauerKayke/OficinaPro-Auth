package entity

import "errors"

// Domain errors - Erros de negócio
var (
	// ErrClienteNotFound quando cliente não existe
	ErrClienteNotFound = errors.New("cliente não encontrado")

	// ErrClienteInativo quando cliente está inativo
	ErrClienteInativo = errors.New("cliente inativo")

	// ErrInvalidDocument quando documento é inválido
	ErrInvalidDocument = errors.New("documento inválido")

	// ErrInvalidCPF quando CPF é inválido
	ErrInvalidCPF = errors.New("CPF inválido")
)
