package service

// ValidatorService define o contrato para validações
// Interface (Port) - Dependency Inversion Principle
type ValidatorService interface {
	// ValidateCPF valida um CPF brasileiro
	ValidateCPF(cpf string) bool

	// NormalizeCPF normaliza um CPF removendo caracteres não numéricos
	NormalizeCPF(cpf string) string
}
