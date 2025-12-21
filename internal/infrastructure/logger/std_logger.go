package logger

import (
	"log"
	"os"

	"github.com/oficinapro/auth-service/internal/handler"
)

// StdLogger implementa handler.Logger usando log padrão do Go.
// Útil para ambientes de desenvolvimento e testes locais.
type StdLogger struct {
	logger *log.Logger
}

// NewStdLogger cria um logger que escreve para stdout.
func NewStdLogger() handler.Logger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

// Debug registra mensagens de debug.
func (l *StdLogger) Debug(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("[DEBUG] %s | %v", msg, fields)
		return
	}
	l.logger.Printf("[DEBUG] %s", msg)
}

// Info registra mensagens informativas.
func (l *StdLogger) Info(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("[INFO] %s | %v", msg, fields)
		return
	}
	l.logger.Printf("[INFO] %s", msg)
}

// Error registra mensagens de erro.
func (l *StdLogger) Error(msg string, fields ...interface{}) {
	if len(fields) > 0 {
		l.logger.Printf("[ERROR] %s | %v", msg, fields)
		return
	}
	l.logger.Printf("[ERROR] %s", msg)
}
