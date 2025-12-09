package telemetry

import (
	"context"
	"time"

	"github.com/oficinapro/auth-service/internal/domain/service"
)

// OTelService implementação stub temporária até instalar dependências corretas.
// TODO: Implementar com OpenTelemetry real após resolver dependências.
type OTelService struct {
	serviceName string
}

// NewOTelService cria stub temporário.
func NewOTelService(ctx context.Context, cfg Config) (*OTelService, error) {
	return &OTelService{
		serviceName: cfg.ServiceName,
	}, nil
}

func (s *OTelService) StartSpan(ctx context.Context, name string) (context.Context, service.Span) {
	return ctx, &noOpSpan{}
}

func (s *OTelService) RecordMetric(name string, value float64, attributes map[string]interface{}) {}

func (s *OTelService) RecordDuration(name string, duration time.Duration, attributes map[string]interface{}) {
}

func (s *OTelService) IncrementCounter(name string, attributes map[string]interface{}) {}

func (s *OTelService) Shutdown(ctx context.Context) error {
	return nil
}
