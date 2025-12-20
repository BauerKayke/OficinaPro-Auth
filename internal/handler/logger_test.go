package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoOpLogger(t *testing.T) {
	logger := &NoOpLogger{}

	// Não deve paniquear
	assert.NotPanics(t, func() {
		logger.Info("test", "key", "value")
		logger.Error("test", "key", "value")
		logger.Debug("test", "key", "value")
	})
}

func TestStdLogger(t *testing.T) {
	logger := &StdLogger{}

	// Não deve paniquear
	assert.NotPanics(t, func() {
		logger.Info("test message", "key1", "value1", "key2", "value2")
		logger.Error("error message", "error", "details")
		logger.Debug("debug message")
	})
}

func TestFormatFields(t *testing.T) {
	tests := []struct {
		name     string
		fields   []interface{}
		expected string
	}{
		{
			name:     "campos vazios",
			fields:   []interface{}{},
			expected: "",
		},
		{
			name:     "um par key-value",
			fields:   []interface{}{"key", "value"},
			expected: "key=value",
		},
		{
			name:     "múltiplos pares",
			fields:   []interface{}{"key1", "value1", "key2", "value2"},
			expected: "key1=value1 key2=value2",
		},
		{
			name:     "campo ímpar (sem par)",
			fields:   []interface{}{"key1", "value1", "key2"},
			expected: "key1=value1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFields(tt.fields)
			assert.Equal(t, tt.expected, result)
		})
	}
}
