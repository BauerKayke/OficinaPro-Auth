package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/oficinapro/auth-service/internal/domain/service"
)

// OTelService implementa TelemetryService com OpenTelemetry completo
type OTelService struct {
	tracer         trace.Tracer
	meter          metric.Meter
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	serviceName    string
}

// NewOTelService cria uma instância completa de OpenTelemetry
func NewOTelService(ctx context.Context, cfg Config) (*OTelService, error) {
	// 1. Criar Resource com atributos do serviço
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
			attribute.String("telemetry.sdk.language", "go"),
			attribute.String("telemetry.sdk.name", "opentelemetry"),
			attribute.String("cloud.provider", "aws"),
			attribute.String("cloud.platform", "aws_lambda"),
			attribute.String("faas.name", cfg.ServiceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 2. Configurar Trace Exporter (New Relic OTLP HTTP)
	// NewRelicEndpoint deve ser "otlp.nr-data.net" (porta 443 HTTPS por padrão)
	traceExporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.NewRelicEndpoint),
		otlptracehttp.WithHeaders(map[string]string{
			"api-key": cfg.NewRelicKey,
		}),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// 3. Configurar Trace Provider
	// IMPORTANTE: Para Lambda, usar SimpleSpanProcessor para envio SÍNCRONO
	// BatchSpanProcessor não funciona bem com Lambda porque o timeout de 5s
	// é maior que a duração típica de uma invocação Lambda (~100-500ms)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(traceExporter)),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.SampleRate)),
	)
	otel.SetTracerProvider(tracerProvider)

	// 4. Configurar Propagators (W3C Trace Context + Baggage)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 5. Configurar Metric Exporter (New Relic OTLP HTTP)
	metricExporter, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(cfg.NewRelicEndpoint),
		otlpmetrichttp.WithHeaders(map[string]string{
			"api-key": cfg.NewRelicKey,
		}),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	// 6. Configurar Metric Provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
			metricExporter,
			sdkmetric.WithInterval(60*time.Second), // Export a cada 60s
		)),
		sdkmetric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	// 7. Criar Tracer e Meter
	tracer := tracerProvider.Tracer(cfg.ServiceName)
	meter := meterProvider.Meter(cfg.ServiceName)

	return &OTelService{
		tracer:         tracer,
		meter:          meter,
		tracerProvider: tracerProvider,
		meterProvider:  meterProvider,
		serviceName:    cfg.ServiceName,
	}, nil
}

// StartSpan cria um novo span
func (s *OTelService) StartSpan(ctx context.Context, name string) (context.Context, service.Span) {
	ctx, span := s.tracer.Start(ctx, name)
	return ctx, &otelSpan{span: span}
}

// RecordMetric registra uma métrica custom
func (s *OTelService) RecordMetric(name string, value float64, attributes map[string]interface{}) {
	counter, err := s.meter.Float64Counter(name)
	if err != nil {
		return // Falha silenciosa para não quebrar app
	}
	attrs := convertAttributes(attributes)
	counter.Add(context.Background(), value, metric.WithAttributes(attrs...))
}

// RecordDuration registra duração como histogram
func (s *OTelService) RecordDuration(name string, duration time.Duration, attributes map[string]interface{}) {
	histogram, err := s.meter.Float64Histogram(name)
	if err != nil {
		return
	}
	attrs := convertAttributes(attributes)
	histogram.Record(context.Background(), duration.Seconds(), metric.WithAttributes(attrs...))
}

// IncrementCounter incrementa contador
func (s *OTelService) IncrementCounter(name string, attributes map[string]interface{}) {
	counter, err := s.meter.Int64Counter(name)
	if err != nil {
		return
	}
	attrs := convertAttributes(attributes)
	counter.Add(context.Background(), 1, metric.WithAttributes(attrs...))
}

// ForceFlush força envio de todos os spans/métricas pendentes
// Importante para Lambda: garante que dados sejam enviados antes do retorno
func (s *OTelService) ForceFlush(ctx context.Context) error {
	var errs []error

	// Flush tracer (enviar spans pendentes)
	if err := s.tracerProvider.ForceFlush(ctx); err != nil {
		errs = append(errs, fmt.Errorf("tracer flush: %w", err))
	}

	// Flush meter (enviar métricas pendentes)
	if err := s.meterProvider.ForceFlush(ctx); err != nil {
		errs = append(errs, fmt.Errorf("meter flush: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("flush errors: %v", errs)
	}
	return nil
}

// Shutdown gracefully fecha providers
func (s *OTelService) Shutdown(ctx context.Context) error {
	var errs []error

	// Shutdown tracer (flush pending spans)
	if err := s.tracerProvider.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("tracer shutdown: %w", err))
	}

	// Shutdown meter (flush pending metrics)
	if err := s.meterProvider.Shutdown(ctx); err != nil {
		errs = append(errs, fmt.Errorf("meter shutdown: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}
	return nil
}

// otelSpan implementa service.Span
type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) SetAttribute(key string, value interface{}) {
	switch v := value.(type) {
	case string:
		s.span.SetAttributes(attribute.String(key, v))
	case int:
		s.span.SetAttributes(attribute.Int(key, v))
	case int64:
		s.span.SetAttributes(attribute.Int64(key, v))
	case float64:
		s.span.SetAttributes(attribute.Float64(key, v))
	case bool:
		s.span.SetAttributes(attribute.Bool(key, v))
	default:
		s.span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}

func (s *otelSpan) SetStatus(err error) {
	if err != nil {
		s.span.RecordError(err)
	}
}

func (s *otelSpan) AddEvent(name string, attributes map[string]interface{}) {
	attrs := convertAttributes(attributes)
	s.span.AddEvent(name, trace.WithAttributes(attrs...))
}

func (s *otelSpan) End() {
	s.span.End()
}

// convertAttributes converte map[string]interface{} para []attribute.KeyValue
func convertAttributes(attrs map[string]interface{}) []attribute.KeyValue {
	if attrs == nil {
		return nil
	}

	result := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		switch val := v.(type) {
		case string:
			result = append(result, attribute.String(k, val))
		case int:
			result = append(result, attribute.Int(k, val))
		case int64:
			result = append(result, attribute.Int64(k, val))
		case float64:
			result = append(result, attribute.Float64(k, val))
		case bool:
			result = append(result, attribute.Bool(k, val))
		default:
			result = append(result, attribute.String(k, fmt.Sprintf("%v", val)))
		}
	}
	return result
}
