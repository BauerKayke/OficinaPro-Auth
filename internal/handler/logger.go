package handler

import (
	"fmt"
	"log"
	"strings"
)

// Logger interface para logging estruturado com key-value pairs.
// Permite trocar implementação sem afetar o código (Dependency Inversion).
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NoOpLogger é uma implementação que não realiza logging.
// Útil para testes onde logs não são necessários.
type NoOpLogger struct{}

// Info não faz nada (no-op).
func (l *NoOpLogger) Info(msg string, fields ...interface{}) {}

// Error não faz nada (no-op).
func (l *NoOpLogger) Error(msg string, fields ...interface{}) {}

// Debug não faz nada (no-op).
func (l *NoOpLogger) Debug(msg string, fields ...interface{}) {}

// StdLogger implementa logging estruturado usando log padrão do Go.
// Formata logs com key-value pairs para melhor rastreabilidade.
type StdLogger struct{}

// Info registra mensagem informativa com campos estruturados.
func (l *StdLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %s", msg, formatFields(fields))
}

// Error registra mensagem de erro com campos estruturados.
func (l *StdLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %s", msg, formatFields(fields))
}

// Debug registra mensagem de debug com campos estruturados.
func (l *StdLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %s", msg, formatFields(fields))
}

// formatFields formata campos variádicos em formato key=value.
// Espera que fields venha em pares: key1, value1, key2, value2, ...
func formatFields(fields []interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var pairs []string
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprintf("%v", fields[i])
			value := fmt.Sprintf("%v", fields[i+1])
			pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(pairs, " ")
}
