package logger_test

import (
	"testing"

	"github.com/oficinapro/auth-service/internal/handler"
	infraLogger "github.com/oficinapro/auth-service/internal/infrastructure/logger"
	"github.com/stretchr/testify/assert"
)

func TestStdLogger_Interface(t *testing.T) {
	logger := infraLogger.NewStdLogger()

	// Verifica que implementa a interface handler.Logger
	var _ handler.Logger = logger

	assert.NotNil(t, logger)
}

func TestStdLogger_Methods(t *testing.T) {
	logger := infraLogger.NewStdLogger()

	// Testa que os métodos não causam panic
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Debug("debug with fields", "key", "value")

		logger.Info("info message")
		logger.Info("info with fields", "key", "value")

		logger.Error("error message")
		logger.Error("error with fields", "key", "value")
	})
}

func TestStdLogger_Output(t *testing.T) {
	// Captura output para testar (em ambiente real, seria mais complexo)
	// Aqui apenas testamos que não há panic
	logger := infraLogger.NewStdLogger()

	testCases := []struct {
		name   string
		logFn  func()
		expect string
	}{
		{
			name:   "Debug without fields",
			logFn:  func() { logger.Debug("test message") },
			expect: "DEBUG",
		},
		{
			name:   "Info with fields",
			logFn:  func() { logger.Info("test info", "user", "john") },
			expect: "INFO",
		},
		{
			name:   "Error message",
			logFn:  func() { logger.Error("test error") },
			expect: "ERROR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Apenas verifica que não causa panic
			assert.NotPanics(t, tc.logFn)
		})
	}
}
