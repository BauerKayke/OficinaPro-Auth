package service

import (
	"context"
	"time"
)

// TelemetryService define operações de observabilidade (traces, métricas, logs).
// Interface de domínio que permite trocar implementação (New Relic, Datadog, etc).
type TelemetryService interface {
	// StartSpan inicia um novo span de tracing
	StartSpan(ctx context.Context, name string) (context.Context, Span)

	// RecordMetric registra uma métrica customizada
	RecordMetric(name string, value float64, attributes map[string]interface{})

	// RecordDuration registra duração de operação
	RecordDuration(name string, duration time.Duration, attributes map[string]interface{})

	// IncrementCounter incrementa contador
	IncrementCounter(name string, attributes map[string]interface{})

	// Shutdown finaliza o telemetry service
	Shutdown(ctx context.Context) error
}

// Span representa um span de tracing distribuído.
type Span interface {
	// SetAttribute adiciona atributo ao span
	SetAttribute(key string, value interface{})

	// SetStatus define status do span (sucesso/erro)
	SetStatus(err error)

	// AddEvent adiciona evento ao span
	AddEvent(name string, attributes map[string]interface{})

	// End finaliza o span
	End()
}
