package logger_test

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oficinapro/auth-service/internal/handler"
	infraLogger "github.com/oficinapro/auth-service/internal/infrastructure/logger"
)

func TestJSONLogger_Interface(t *testing.T) {
	logger := infraLogger.NewJSONLogger("test-service", "1.0.0", "test")

	// Verifica que implementa a interface handler.Logger
	var _ handler.Logger = logger
}

func TestJSONLogger_Methods(t *testing.T) {
	logger := infraLogger.NewJSONLogger("test-service", "1.0.0", "test")

	// Testa que os métodos não causam panic
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Error("error message")
	})
}

func TestJSONLogger_Output(t *testing.T) {
	// Captura stdout para testar o output JSON
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := infraLogger.NewJSONLogger("test-service", "1.0.0", "test")
	logger.Info("test message", "key1", "value1", "key2", 123)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verificar que é JSON válido
	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err, "Output should be valid JSON")

	// Verificar campos obrigatórios
	assert.Equal(t, "INFO", logEntry["level"])
	assert.Equal(t, "test message", logEntry["message"])
	assert.Equal(t, "test-service", logEntry["service.name"])
	assert.Equal(t, "1.0.0", logEntry["service.version"])
	assert.Equal(t, "test", logEntry["environment"])
	assert.NotEmpty(t, logEntry["timestamp"])

	// Verificar campos customizados
	fields, ok := logEntry["fields"].(map[string]interface{})
	require.True(t, ok, "Fields should be a map")
	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, float64(123), fields["key2"]) // JSON unmarshals numbers as float64
}

func TestJSONLogger_Levels(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(handler.Logger)
		expected string
	}{
		{
			name:     "Debug level",
			logFunc:  func(l handler.Logger) { l.Debug("debug msg") },
			expected: "DEBUG",
		},
		{
			name:     "Info level",
			logFunc:  func(l handler.Logger) { l.Info("info msg") },
			expected: "INFO",
		},
		{
			name:     "Error level",
			logFunc:  func(l handler.Logger) { l.Error("error msg") },
			expected: "ERROR",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger := infraLogger.NewJSONLogger("test-service", "1.0.0", "test")
			tc.logFunc(logger)

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			var logEntry map[string]interface{}
			err := json.Unmarshal([]byte(output), &logEntry)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, logEntry["level"])
		})
	}
}

func TestJSONLogger_WithoutFields(t *testing.T) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger := infraLogger.NewJSONLogger("test-service", "1.0.0", "test")
	logger.Info("message without fields")

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	require.NoError(t, err)

	// Fields should be nil/empty when not provided
	_, hasFields := logEntry["fields"]
	assert.False(t, hasFields, "Fields should not be present when empty")
}
