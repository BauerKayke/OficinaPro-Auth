package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/oficinapro/auth-service/internal/handler"
)

// JSONLogger implementa handler.Logger com logs estruturados em JSON.
// Ideal para ambientes de produção e integração com sistemas de observabilidade.
type JSONLogger struct {
	serviceName    string
	serviceVersion string
	environment    string
}

// LogEntry representa uma entrada de log estruturada.
type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service.name"`
	Version     string                 `json:"service.version"`
	Environment string                 `json:"environment"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

// NewJSONLogger cria um logger JSON estruturado para produção.
func NewJSONLogger(serviceName, serviceVersion, environment string) handler.Logger {
	return &JSONLogger{
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
		environment:    environment,
	}
}

// parseFields converte campos variádicos em um mapa.
func parseFields(fields []interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	result := make(map[string]interface{})
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if ok {
			result[key] = fields[i+1]
		}
	}
	return result
}

// log escreve uma entrada de log em JSON.
func (l *JSONLogger) log(level, msg string, fields ...interface{}) {
	entry := LogEntry{
		Timestamp:   time.Now().UTC().Format(time.RFC3339Nano),
		Level:       level,
		Message:     msg,
		Service:     l.serviceName,
		Version:     l.serviceVersion,
		Environment: l.environment,
		Fields:      parseFields(fields),
	}

	// Serializar para JSON e escrever em stdout
	data, err := json.Marshal(entry)
	if err != nil {
		// Fallback para formato simples se JSON falhar
		os.Stdout.WriteString(`{"level":"ERROR","message":"failed to marshal log entry"}` + "\n")
		return
	}

	os.Stdout.Write(data)
	os.Stdout.WriteString("\n")
}

// Debug registra mensagens de debug em JSON.
func (l *JSONLogger) Debug(msg string, fields ...interface{}) {
	l.log("DEBUG", msg, fields...)
}

// Info registra mensagens informativas em JSON.
func (l *JSONLogger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

// Error registra mensagens de erro em JSON.
func (l *JSONLogger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}
