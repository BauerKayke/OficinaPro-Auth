package di

import (
	"context"

	"github.com/oficinapro/auth-service/di/config"
	"github.com/oficinapro/auth-service/internal/domain/service"
	"github.com/oficinapro/auth-service/internal/infrastructure/telemetry"
)

// NewTelemetryService configura serviço de observabilidade (OpenTelemetry + New Relic).
// Retorna telemetry service, função de cleanup e erro caso falhe.
func NewTelemetryService(ctx context.Context, cfg *config.Config) (service.TelemetryService, func(), error) {
	// Se telemetry estiver desabilitado, retorna NoOp
	if !cfg.Telemetry.Enabled {
		return telemetry.NewNoOpTelemetryService(), func() {}, nil
	}

	// Criar configuração de telemetry
	telemetryCfg := telemetry.Config{
		ServiceName:      cfg.Telemetry.ServiceName,
		ServiceVersion:   cfg.Telemetry.ServiceVersion,
		Environment:      cfg.App.Environment,
		NewRelicKey:      cfg.Telemetry.NewRelicKey,
		NewRelicEndpoint: cfg.Telemetry.NewRelicEndpoint,
		SampleRate:       cfg.Telemetry.SampleRate,
	}

	// Criar serviço OpenTelemetry
	otelService, err := telemetry.NewOTelService(ctx, telemetryCfg)
	if err != nil {
		return nil, nil, err
	}

	// Função de cleanup
	cleanup := func() {
		_ = otelService.Shutdown(context.Background())
	}

	return otelService, cleanup, nil
}
