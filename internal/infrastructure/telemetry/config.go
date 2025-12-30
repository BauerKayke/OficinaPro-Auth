package telemetry

// Config representa configuração do serviço de telemetry.
type Config struct {
	ServiceName      string  // Nome do serviço
	ServiceVersion   string  // Versão do serviço
	Environment      string  // Ambiente (dev, staging, prod)
	NewRelicKey      string  // New Relic License Key
	NewRelicEndpoint string  // Endpoint OTLP do New Relic
	SampleRate       float64 // Taxa de amostragem (0.0 a 1.0)
}
