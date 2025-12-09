# 📊 Telemetry Guide - OpenTelemetry + New Relic

## 🎯 Visão Geral

Sistema de observabilidade **clean e minimalista** usando OpenTelemetry para enviar traces e métricas para New Relic.

**Filosofia:** Instrumentação enxuta e extensível - apenas o essencial, fácil de expandir quando necessário.

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────┐
│                   Application                        │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │  UseCase     │  │  Handler     │  │  Handler  │ │
│  │  (spans)     │  │  (opcional)  │  │ (opcional)│ │
│  └──────┬───────┘  └──────────────┘  └───────────┘ │
│         │                                            │
│  ┌──────▼──────────────────────────────────────┐   │
│  │  TelemetryService (domain interface)        │   │
│  └──────┬──────────────────────────────────────┘   │
└─────────┼───────────────────────────────────────────┘
          │
┌─────────▼───────────────────────────────────────────┐
│         Infrastructure Layer                         │
│  ┌──────────────────┐  ┌──────────────────┐         │
│  │  OTelService     │  │  NoOpService     │         │
│  │  (production)    │  │  (disabled)      │         │
│  └────────┬─────────┘  └──────────────────┘         │
└───────────┼──────────────────────────────────────────┘
            │
┌───────────▼──────────────────────────────────────────┐
│              OpenTelemetry SDK                        │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐  │
│  │   Tracer     │  │    Meter     │  │  Exporter │  │
│  └──────────────┘  └──────────────┘  └─────┬─────┘  │
└──────────────────────────────────────────────┼────────┘
                                               │
                                    ┌──────────▼─────────┐
                                    │   New Relic OTLP   │
                                    │  (gRPC endpoint)   │
                                    └────────────────────┘
```

---

## 📦 O Que Foi Implementado

### ✅ 1. Domain Layer (Interfaces)

**`internal/domain/service/telemetry_service.go`**

```go
type TelemetryService interface {
    StartSpan(ctx context.Context, name string) (context.Context, Span)
    RecordMetric(name string, value float64, attributes map[string]interface{})
    RecordDuration(name string, duration time.Duration, attributes map[string]interface{})
    IncrementCounter(name string, attributes map[string]interface{})
    Shutdown(ctx context.Context) error
}
```

- ✅ Interface limpa de domínio
- ✅ Framework-agnostic (pode trocar OpenTelemetry por Datadog, etc)
- ✅ Programação contra abstrações (DIP)

---

### ✅ 2. Infrastructure Layer

#### **OTelService** (Produção)
`internal/infrastructure/telemetry/otel_service.go`

- ✅ Implementação completa OpenTelemetry
- ✅ Exporta para New Relic via OTLP/gRPC
- ✅ Lazy creation de métricas (cria sob demanda)
- ✅ Propagação de contexto automática
- ✅ Sampling configurável

#### **NoOpService** (Telemetry Desabilitado)
`internal/infrastructure/telemetry/noop_service.go`

- ✅ Null Object Pattern
- ✅ Zero overhead quando desabilitado
- ✅ Útil para testes e ambientes sem telemetry

---

### ✅ 3. Configuração

**`di/config/config.go`**

```go
type TelemetryConfig struct {
    Enabled         bool    // Liga/desliga telemetry
    ServiceName     string  // Nome do serviço
    ServiceVersion  string  // Versão
    NewRelicKey     string  // License Key
    NewRelicEndpoint string // Endpoint OTLP
    SampleRate      float64 // Taxa de amostragem (0.0-1.0)
}
```

**Variáveis de Ambiente:**
```bash
TELEMETRY_ENABLED=true
TELEMETRY_SERVICE_NAME=oficinapro-auth
TELEMETRY_SERVICE_VERSION=1.0.0
NEW_RELIC_LICENSE_KEY=your-license-key
NEW_RELIC_OTLP_ENDPOINT=https://otlp.nr-data.net:4317
TELEMETRY_SAMPLE_RATE=1.0  # 100% das requisições
```

---

### ✅ 4. Instrumentação no UseCase

**O que foi implementado:**

```go
func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
    // 1. Span principal
    ctx, span := uc.telemetryService.StartSpan(ctx, "authenticate.execute")
    defer span.End()

    startTime := time.Now()

    // ... lógica de negócio ...

    // 2. Em caso de erro
    if err != nil {
        span.SetStatus(err)
        uc.recordMetric("error", "invalid_cpf")
        return nil, err
    }

    // 3. Em caso de sucesso
    uc.recordMetric("success", "")
    uc.telemetryService.RecordDuration("authenticate.duration", time.Since(startTime), nil)

    return output, nil
}
```

**Métricas Capturadas:**
- ✅ `authenticate.requests` (counter) - Total de requests com status
- ✅ `authenticate.duration` (histogram) - Duração em segundos

**Traces Capturados:**
- ✅ Span `authenticate.execute` - Execução completa do usecase
- ✅ Status do span (success/error)
- ✅ Propagação de contexto automática

---

## 📈 Métricas Disponíveis

### Métricas Essenciais (Implementadas)

| Métrica | Tipo | Descrição | Attributes |
|---------|------|-----------|------------|
| `authenticate.requests` | Counter | Total de requisições | `status` (success/error), `error_type` (opcional) |
| `authenticate.duration` | Histogram | Duração da operação | Nenhum (pode adicionar depois) |

### Como Expandir (Fácil)

Para adicionar novas métricas, basta chamar os métodos:

```go
// Contador simples
telemetry.IncrementCounter("minha.metrica", map[string]interface{}{
    "tipo": "exemplo",
})

// Duração
telemetry.RecordDuration("operacao.duration", duration, map[string]interface{}{
    "complexidade": "alta",
})

// Métrica customizada
telemetry.RecordMetric("cache.hit.rate", 0.95, map[string]interface{}{
    "cache_type": "redis",
})
```

✅ Métricas são criadas automaticamente (lazy creation)
✅ Sem necessidade de pré-registrar

---

## 🔍 Traces e Spans

### Estrutura de Spans

```
authenticate.execute (UseCase)
  └─ duração total
  └─ status (success/error)
  └─ error details (se houver)
```

### Como Adicionar Novos Spans

```go
// Span filho dentro de uma operação
ctx, childSpan := telemetry.StartSpan(ctx, "operacao.secundaria")
defer childSpan.End()

// Adicionar atributos
childSpan.SetAttribute("param_importante", valor)

// Marcar erro
if err != nil {
    childSpan.SetStatus(err)
}
```

---

## 🚀 Como Usar

### 1. Configurar Variáveis de Ambiente

```bash
# .env ou exportar no ambiente
export TELEMETRY_ENABLED=true
export NEW_RELIC_LICENSE_KEY=your-key-here
export TELEMETRY_SERVICE_NAME=oficinapro-auth
export TELEMETRY_SERVICE_VERSION=1.0.0
```

### 2. A Telemetry já está integrada no DI

O container DI cuida de tudo automaticamente:

```go
container, err := di.NewContainer(ctx)
// TelemetryService já está configurado e injetado
```

### 3. Verificar no New Relic

Após deploy, acesse New Relic:

**Traces:**
- APM → Your Service → Distributed Tracing
- Verá spans `authenticate.execute`

**Métricas:**
- Metrics Explorer → Search: `authenticate.requests`
- Metrics Explorer → Search: `authenticate.duration`

---

## 🎛️ Configuração Avançada

### Ajustar Sample Rate

Para reduzir volume em produção:

```bash
# Captura apenas 10% das requisições
TELEMETRY_SAMPLE_RATE=0.1
```

### Desabilitar Telemetry

Para ambientes de teste/desenvolvimento:

```bash
TELEMETRY_ENABLED=false
# Usa NoOpService automaticamente (zero overhead)
```

---

## 📊 Dashboards Recomendados no New Relic

### 1. Dashboard de Performance

```nrql
# Latência P50, P95, P99
SELECT percentile(authenticate.duration, 50, 95, 99)
FROM Metric
WHERE service.name = 'oficinapro-auth'

# Taxa de sucesso
SELECT count(*)
FROM Metric
WHERE metricName = 'authenticate.requests'
FACET status
```

### 2. Dashboard de Erros

```nrql
# Erros por tipo
SELECT count(*)
FROM Metric
WHERE metricName = 'authenticate.requests'
  AND status = 'error'
FACET error_type

# Taxa de erro
SELECT (count(*) WHERE status = 'error') / count(*) * 100 as 'Error Rate %'
FROM Metric
WHERE metricName = 'authenticate.requests'
```

---

## 🔧 Troubleshooting

### Telemetry não aparece no New Relic

1. ✅ Verificar license key está correta
2. ✅ Verificar endpoint OTLP está acessível
3. ✅ Verificar `TELEMETRY_ENABLED=true`
4. ✅ Logs da aplicação devem mostrar inicialização do tracer

### Performance impactada

1. ✅ Reduzir `TELEMETRY_SAMPLE_RATE` (ex: 0.1 = 10%)
2. ✅ Remover atributos desnecessários dos spans
3. ✅ Verificar se há spans filhos excessivos

---

## 📚 Referências

- [OpenTelemetry Go Docs](https://opentelemetry.io/docs/instrumentation/go/)
- [New Relic OTLP](https://docs.newrelic.com/docs/more-integrations/open-source-telemetry-integrations/opentelemetry/opentelemetry-introduction/)
- [OTLP Specification](https://opentelemetry.io/docs/reference/specification/protocol/)

---

## ✅ Checklist de Implementação

- [x] Interface de domínio `TelemetryService`
- [x] Implementação OpenTelemetry (OTelService)
- [x] NoOp implementation (zero overhead quando desabilitado)
- [x] Configuração via env vars
- [x] Integração no DI Container
- [x] Instrumentação no UseCase
- [x] Testes com mocks
- [x] Documentação completa

---

**Status:** ✅ Implementação Completa - Clean & Extensível

**Filosofia:** Menos é mais. Instrumentação mínima com máximo valor.

