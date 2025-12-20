#!/bin/bash
# Script para instalar dependências OpenTelemetry
# Uso: ./INSTALL_TELEMETRY_DEPS.sh

set -e

echo "🚀 Instalando dependências OpenTelemetry..."
echo ""

# Forçar uso do proxy público do Go
export GOPROXY=https://proxy.golang.org,direct

# Core OpenTelemetry
echo "📦 Instalando core OpenTelemetry..."
go get go.opentelemetry.io/otel@v1.21.0
go get go.opentelemetry.io/otel/trace@v1.21.0
go get go.opentelemetry.io/otel/metric@v1.21.0
go get go.opentelemetry.io/otel/sdk@v1.21.0

# OTLP Exporters
echo "📦 Instalando OTLP exporters..."
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc@v1.21.0

# Semantic conventions
echo "📦 Instalando semantic conventions..."
go get go.opentelemetry.io/otel/semconv/v1.17.0

# gRPC (para exporter)
echo "📦 Instalando gRPC..."
go get google.golang.org/grpc@latest

# Tidying
echo "🧹 Limpando go.mod..."
go mod tidy

echo ""
echo "✅ Dependências instaladas com sucesso!"
echo ""
echo "Para verificar:"
echo "  go mod verify"
echo ""
echo "Para buildar:"
echo "  go build ./..."

