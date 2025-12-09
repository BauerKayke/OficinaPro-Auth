package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPFValidator_ValidateCPF(t *testing.T) {
	validator := NewCPFValidator()

	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{
			name:     "CPF válido com pontos e traço",
			cpf:      "123.456.789-09",
			expected: true,
		},
		{
			name:     "CPF válido sem formatação",
			cpf:      "12345678909",
			expected: true,
		},
		{
			name:     "CPF inválido - todos dígitos iguais",
			cpf:      "111.111.111-11",
			expected: false,
		},
		{
			name:     "CPF inválido - dígito verificador errado",
			cpf:      "123.456.789-00",
			expected: false,
		},
		{
			name:     "CPF inválido - tamanho incorreto",
			cpf:      "123.456.789",
			expected: false,
		},
		{
			name:     "CPF vazio",
			cpf:      "",
			expected: false,
		},
		{
			name:     "CPF com letras",
			cpf:      "ABC.DEF.GHI-JK",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateCPF(tt.cpf)
			assert.Equal(t, tt.expected, result, "Validation result should match expected")
		})
	}
}

func TestCPFValidator_NormalizeCPF(t *testing.T) {
	validator := NewCPFValidator()

	tests := []struct {
		name     string
		cpf      string
		expected string
	}{
		{
			name:     "CPF com pontos e traço",
			cpf:      "123.456.789-09",
			expected: "12345678909",
		},
		{
			name:     "CPF já normalizado",
			cpf:      "12345678909",
			expected: "12345678909",
		},
		{
			name:     "CPF com espaços",
			cpf:      "123 456 789 09",
			expected: "12345678909",
		},
		{
			name:     "CPF vazio",
			cpf:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.NormalizeCPF(tt.cpf)
			assert.Equal(t, tt.expected, result, "Normalized CPF should match expected")
		})
	}
}

func BenchmarkCPFValidator_ValidateCPF(b *testing.B) {
	validator := NewCPFValidator()
	cpf := "123.456.789-09"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateCPF(cpf)
	}
}
