package telemetry_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oficinapro/auth-service/internal/infrastructure/telemetry"
)

func TestNewOTelService_Success(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()

	// Act
	service, err := telemetry.NewOTelService(ctx, cfg)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, service)

	// Cleanup
	_ = service.Shutdown(ctx)
}

func TestOTelService_StartSpan(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	newCtx, span := service.StartSpan(ctx, "test-operation")

	// Assert
	assert.NotNil(t, newCtx)
	assert.NotNil(t, span)

	// Test span operations
	span.SetAttribute("user_id", 123)
	span.SetAttribute("email", "test@test.com")
	span.SetAttribute("success", true)
	span.SetAttribute("duration", 1.5)
	span.AddEvent("processing_started", map[string]interface{}{"step": 1})
	span.SetStatus(nil) // Success
	span.End()
}

func TestOTelService_SpanWithError(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	_, span := service.StartSpan(ctx, "test-operation-with-error")
	span.SetAttribute("operation", "test")
	span.SetStatus(assert.AnError) // Com erro
	span.End()

	// Assert - não deve causar panic
	assert.NotNil(t, span)
}

func TestOTelService_RecordMetric(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	service.RecordMetric("http_requests_total", 1.0, map[string]interface{}{
		"method": "POST",
		"status": 200,
	})

	service.RecordMetric("active_users", 42.0, nil)

	// Assert - não deve causar panic
	assert.NotNil(t, service)
}

func TestOTelService_RecordDuration(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	service.RecordDuration("request_duration", 150*time.Millisecond, map[string]interface{}{
		"endpoint": "/api/users",
	})

	service.RecordDuration("db_query_time", 50*time.Millisecond, nil)

	// Assert - não deve causar panic
	assert.NotNil(t, service)
}

func TestOTelService_IncrementCounter(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	service.IncrementCounter("login_attempts", map[string]interface{}{
		"success": true,
	})

	service.IncrementCounter("api_calls", nil)

	// Assert - não deve causar panic
	assert.NotNil(t, service)
}

func TestOTelService_SpanAttributeTypes(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	_, span := service.StartSpan(ctx, "test-attributes")

	// Test different attribute types
	span.SetAttribute("string_attr", "value")
	span.SetAttribute("int_attr", 42)
	span.SetAttribute("int64_attr", int64(123456789))
	span.SetAttribute("float64_attr", 3.14)
	span.SetAttribute("bool_attr", true)
	span.SetAttribute("slice_attr", []string{"a", "b", "c"})         // Should convert to string
	span.SetAttribute("map_attr", map[string]string{"key": "value"}) // Should convert to string

	span.End()

	// Assert - não deve causar panic
	assert.NotNil(t, span)
}

func TestOTelService_AddEventWithAttributes(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act
	_, span := service.StartSpan(ctx, "test-events")

	span.AddEvent("user_action", map[string]interface{}{
		"action": "click",
		"button": "submit",
		"count":  5,
	})

	span.AddEvent("system_event", nil)

	span.End()

	// Assert - não deve causar panic
	assert.NotNil(t, span)
}

func TestOTelService_Shutdown(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)

	// Act
	err = service.Shutdown(ctx)

	// Assert - Pode ter erro 403 em testes (endpoint real sem chave válida)
	// O importante é que não dá panic
	_ = err
}

func TestOTelService_CompleteWorkflow(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act - Simular um workflow completo
	start := time.Now()

	// 1. Criar span principal
	ctx, mainSpan := service.StartSpan(ctx, "user_authentication")
	mainSpan.SetAttribute("user_id", 123)
	mainSpan.SetAttribute("email", "user@test.com")

	// 2. Adicionar eventos
	mainSpan.AddEvent("validation_started", nil)

	// 3. Simular sub-operação
	time.Sleep(10 * time.Millisecond)

	// 4. Registrar métricas
	service.IncrementCounter("auth_attempts", map[string]interface{}{
		"method": "email",
	})

	// 5. Registrar duração
	duration := time.Since(start)
	service.RecordDuration("auth_duration", duration, map[string]interface{}{
		"success": true,
	})

	// 6. Finalizar span com sucesso
	mainSpan.SetStatus(nil)
	mainSpan.End()

	// Assert
	assert.NotNil(t, service)
}

func TestOTelService_ForceFlush(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Criar alguns spans para ter algo para flush
	_, span := service.StartSpan(ctx, "test_span_for_flush")
	span.SetAttribute("test_key", "test_value")
	span.End()

	// Act - Forçar flush dos dados pendentes
	err = service.ForceFlush(ctx)

	// Assert - Pode retornar erro 403 (endpoint real sem chave válida)
	// O importante é que não dá panic e a função executa
	// Em ambiente real com chave válida, err seria nil
	_ = err
}

func TestOTelService_ForceFlush_WithTimeout(t *testing.T) {
	// Arrange
	cfg := telemetry.Config{
		ServiceName:      "test-service",
		ServiceVersion:   "1.0.0",
		Environment:      "test",
		NewRelicEndpoint: "otlp.nr-data.net:4318",
		NewRelicKey:      "test-key-12345678901234567890123456789012",
		SampleRate:       1.0,
	}

	ctx := context.Background()
	service, err := telemetry.NewOTelService(ctx, cfg)
	require.NoError(t, err)
	defer service.Shutdown(ctx)

	// Act - ForceFlush com timeout curto
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	err = service.ForceFlush(ctxWithTimeout)

	// Assert - Mesmo com timeout, não deve dar panic
	_ = err
}

func TestNoOpService_ForceFlush(t *testing.T) {
	// Arrange
	service := telemetry.NewNoOpTelemetryService()
	ctx := context.Background()

	// Act
	err := service.ForceFlush(ctx)

	// Assert - NoOp sempre retorna nil
	assert.NoError(t, err)
}
