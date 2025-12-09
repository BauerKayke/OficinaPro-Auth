package validator

import (
	"regexp"
	"strconv"
)

// CPFValidator implementa validação de CPF brasileiro
// Adapter - implementa interface do domínio
type CPFValidator struct{}

// NewCPFValidator cria uma nova instância do validador
func NewCPFValidator() *CPFValidator {
	return &CPFValidator{}
}

// ValidateCPF valida um CPF brasileiro
// Implementa o algoritmo oficial de validação de CPF
func (v *CPFValidator) ValidateCPF(cpf string) bool {
	// Normalizar CPF
	cpf = v.NormalizeCPF(cpf)

	// Verificar tamanho
	if len(cpf) != 11 {
		return false
	}

	// Verificar se todos os dígitos são iguais (CPF inválido)
	if allDigitsEqual(cpf) {
		return false
	}

	// Converter string para slice de inteiros
	digits := make([]int, 11)
	for i, char := range cpf {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		digits[i] = digit
	}

	// Validar primeiro dígito verificador
	if !validateDigit(digits[:9], digits[9]) {
		return false
	}

	// Validar segundo dígito verificador
	if !validateDigit(digits[:10], digits[10]) {
		return false
	}

	return true
}

// NormalizeCPF remove caracteres não numéricos do CPF
func (v *CPFValidator) NormalizeCPF(cpf string) string {
	// Remover tudo que não é número
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(cpf, "")
}

// validateDigit valida um dígito verificador do CPF
func validateDigit(digits []int, expectedDigit int) bool {
	sum := 0
	length := len(digits)

	for i, digit := range digits {
		sum += digit * (length + 1 - i)
	}

	remainder := sum % 11
	calculatedDigit := 0
	if remainder >= 2 {
		calculatedDigit = 11 - remainder
	}

	return calculatedDigit == expectedDigit
}

// allDigitsEqual verifica se todos os dígitos são iguais
func allDigitsEqual(cpf string) bool {
	if len(cpf) == 0 {
		return false
	}

	firstDigit := cpf[0]
	for _, digit := range cpf {
		if digit != firstDigit {
			return false
		}
	}

	return true
}
