package telemetry

import (
	"context"
	"time"

	"github.com/oficinapro/auth-service/internal/domain/service"
)

// NoOpTelemetryService implementação que não faz nada (Null Object Pattern).
// Usado quando telemetry está desabilitado.
type NoOpTelemetryService struct{}

// NewNoOpTelemetryService cria serviço de telemetry no-op.
func NewNoOpTelemetryService() *NoOpTelemetryService {
	return &NoOpTelemetryService{}
}

// StartSpan retorna contexto inalterado e span no-op.
func (s *NoOpTelemetryService) StartSpan(ctx context.Context, name string) (context.Context, service.Span) {
	return ctx, &noOpSpan{}
}

// RecordMetric não faz nada.
func (s *NoOpTelemetryService) RecordMetric(name string, value float64, attributes map[string]interface{}) {
}

// RecordDuration não faz nada.
func (s *NoOpTelemetryService) RecordDuration(name string, duration time.Duration, attributes map[string]interface{}) {
}

// IncrementCounter não faz nada.
func (s *NoOpTelemetryService) IncrementCounter(name string, attributes map[string]interface{}) {}

// Shutdown não faz nada.
func (s *NoOpTelemetryService) Shutdown(ctx context.Context) error {
	return nil
}

// noOpSpan implementação de Span que não faz nada.
type noOpSpan struct{}

// SetAttribute não faz nada.
func (s *noOpSpan) SetAttribute(key string, value interface{}) {}

// SetStatus não faz nada.
func (s *noOpSpan) SetStatus(err error) {}

// AddEvent não faz nada.
func (s *noOpSpan) AddEvent(name string, attributes map[string]interface{}) {}

// End não faz nada.
func (s *noOpSpan) End() {}
