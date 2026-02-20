package telemetry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oficinapro/auth-service/internal/infrastructure/telemetry"
)

func TestNoOpService(t *testing.T) {
	service := telemetry.NewNoOpTelemetryService()

	ctx := context.Background()

	// Test StartSpan
	newCtx, span := service.StartSpan(ctx, "test-span")
	assert.NotNil(t, newCtx)
	assert.NotNil(t, span)

	// Test span methods
	span.SetAttribute("key", "value")
	span.SetStatus(nil)
	span.AddEvent("event", nil)
	span.End()

	// Test metrics methods
	service.RecordMetric("test-metric", 1.0, nil)
	service.RecordDuration("test-duration", 0, nil)
	service.IncrementCounter("test-counter", nil)

	// Test shutdown
	err := service.Shutdown(ctx)
	assert.NoError(t, err)
}
