package validator

import (
	"testing"
)

// Benchmark CPF validation
func BenchmarkValidateCPF(b *testing.B) {
	validator := NewCPFValidator()
	cpfs := []string{
		"12345678909",
		"98765432100",
		"11144477735",
		"55566688899",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateCPF(cpfs[i%len(cpfs)])
	}
}

// Benchmark CPF normalization
func BenchmarkNormalizeCPF(b *testing.B) {
	validator := NewCPFValidator()
	cpfs := []string{
		"123.456.789-09",
		"987.654.321-00",
		"111.444.777-35",
		"555.666.888-99",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.NormalizeCPF(cpfs[i%len(cpfs)])
	}
}

// Benchmark validação completa (normalize + validate)
func BenchmarkValidateCPFComplete(b *testing.B) {
	validator := NewCPFValidator()
	cpf := "123.456.789-09"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		normalized := validator.NormalizeCPF(cpf)
		validator.ValidateCPF(normalized)
	}
}
