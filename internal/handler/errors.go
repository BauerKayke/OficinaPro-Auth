package handler

import "errors"

// Handler errors - erros específicos do handler
var (
	// ErrMissingAuthHeader quando header Authorization não está presente
	ErrMissingAuthHeader = errors.New("missing authorization header")
)
