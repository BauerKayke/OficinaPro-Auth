package di_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oficinapro/auth-service/di"
	"github.com/oficinapro/auth-service/di/config"
)

func TestNewTelemetryService_Disabled(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		Telemetry: config.TelemetryConfig{
			Enabled:        false,
			ServiceName:    "test",
			ServiceVersion: "1.0.0",
		},
	}

	ctx := context.Background()

	// Act
	telemetryService, closeFunc, err := di.NewTelemetryService(ctx, cfg)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, telemetryService)
	assert.NotNil(t, closeFunc)

	// Verificar que é NoOp (não causa erros)
	_, span := telemetryService.StartSpan(ctx, "test")
	span.SetAttribute("key", "value")
	span.End()

	closeFunc()
}

func TestNewTelemetryService_Enabled(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		Telemetry: config.TelemetryConfig{
			Enabled:          true,
			ServiceName:      "test-service",
			ServiceVersion:   "1.0.0",
			NewRelicKey:      "test-key-12345678901234567890123456789012",
			NewRelicEndpoint: "otlp.nr-data.net:4318",
			SampleRate:       1.0,
		},
		App: config.AppConfig{
			Environment: "test",
		},
	}

	ctx := context.Background()

	// Act
	telemetryService, closeFunc, err := di.NewTelemetryService(ctx, cfg)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, telemetryService)
	assert.NotNil(t, closeFunc)

	// Cleanup
	closeFunc()
}
