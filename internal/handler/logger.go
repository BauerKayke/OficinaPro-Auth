package handler

import "log"

// Logger interface para logging estruturado
type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NoOpLogger logger que não faz nada (útil para testes)
type NoOpLogger struct{}

func (l *NoOpLogger) Info(msg string, fields ...interface{})  {}
func (l *NoOpLogger) Error(msg string, fields ...interface{}) {}
func (l *NoOpLogger) Debug(msg string, fields ...interface{}) {}

// StdLogger wrapper para log.Printf padrão
type StdLogger struct{}

func (l *StdLogger) Info(msg string, fields ...interface{}) {
	log.Printf("[INFO] %s %v", msg, fields)
}

func (l *StdLogger) Error(msg string, fields ...interface{}) {
	log.Printf("[ERROR] %s %v", msg, fields)
}

func (l *StdLogger) Debug(msg string, fields ...interface{}) {
	log.Printf("[DEBUG] %s %v", msg, fields)
}
