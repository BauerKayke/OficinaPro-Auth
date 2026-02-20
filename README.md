# 🔐 Auth Gateway - OficinaPro

[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://golang.org/)
[![AWS Lambda](https://img.shields.io/badge/AWS-Lambda-FF9900?logo=aws-lambda)](https://aws.amazon.com/lambda/)
[![JWT](https://img.shields.io/badge/JWT-HMAC--SHA256-000000?logo=json-web-tokens)](https://jwt.io/)
[![Coverage](https://img.shields.io/badge/Coverage-91.8%25-brightgreen)](TEST_COVERAGE_REPORT.md)
[![License](https://img.shields.io/badge/License-Academic-blue)](LICENSE)

## 📋 Descrição

Serviço de autenticação **serverless** do projeto OficinaPro, responsável por validar credenciais e gerar tokens JWT. Desenvolvido em **Go** e executado como **AWS Lambda** integrado com **API Gateway**.

## 🎯 Propósito

- Autenticar usuários via email e senha
- Gerar tokens JWT para autorização em outros serviços
- Validar tokens JWT existentes
- Prover health checks para monitoramento
- Escalar automaticamente com demanda (serverless)

## 🛠️ Tecnologias Utilizadas

| Tecnologia | Versão | Descrição |
|------------|--------|-----------|
| **Go** | 1.23+ | Linguagem de programação |
| **Chi Router** | 5.x | HTTP router |
| **GORM** | 2.x | ORM para PostgreSQL |
| **golang-jwt/jwt** | 5.x | Geração/validação JWT |
| **bcrypt** | - | Hash de senhas |
| **AWS Lambda** | - | Runtime serverless |
| **API Gateway** | HTTP API | Gateway de entrada |
| **Terraform** | 1.5+ | Infrastructure as Code |

## 📁 Estrutura do Projeto

```
auth-oficinapro/
├── cmd/
│   ├── http/                   # Entrypoint HTTP server (dev)
│   │   └── main.go
│   └── lambda/                 # Entrypoint AWS Lambda
│       └── main.go
├── internal/
│   ├── domain/                 # Camada de domínio
│   │   ├── entity/            # Entidades (User)
│   │   ├── repository/        # Interfaces de repositório
│   │   └── service/           # Serviços de domínio
│   ├── usecase/               # Casos de uso (Login, Validate)
│   ├── handler/               # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── health_handler.go
│   │   └── response/          # DTOs de resposta
│   ├── adapter/               # Adaptadores (Repository impl)
│   └── infrastructure/        # Infraestrutura
│       ├── database/          # Conexão GORM
│       ├── jwt/               # Serviço JWT
│       └── telemetry/         # OpenTelemetry
├── di/                        # Dependency Injection
├── terraform/lambda/          # IaC para Lambda + API GW
├── k8s/order-service/         # Manifests K8s (exemplo)
├── examples/order-service/    # Serviço exemplo em Go
├── scripts/                   # Scripts utilitários
├── Dockerfile                 # Build Lambda
├── Makefile                   # Comandos úteis
└── README.md
```

## 🏗️ Arquitetura do Serviço

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                       AUTH GATEWAY - ARQUITETURA                             │
│                                                                              │
│   ┌────────────────┐                                                        │
│   │    Cliente     │                                                        │
│   │  (App/Browser) │                                                        │
│   └───────┬────────┘                                                        │
│           │                                                                  │
│           │ HTTPS POST /auth                                                 │
│           │ {email, senha}                                                   │
│           ▼                                                                  │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                       AWS API GATEWAY                                │   │
│   │                   (HTTP API, CORS, Throttling)                       │   │
│   └───────────────────────────────┬─────────────────────────────────────┘   │
│                                   │                                          │
│                                   │ Lambda Invoke                            │
│                                   ▼                                          │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                      AWS LAMBDA FUNCTION                             │   │
│   │                       (Auth Gateway Go)                              │   │
│   │                                                                      │   │
│   │   ┌─────────────────────────────────────────────────────────────┐   │   │
│   │   │                     CLEAN ARCHITECTURE                       │   │   │
│   │   │                                                              │   │   │
│   │   │  ┌─────────────┐   ┌─────────────┐   ┌─────────────────┐   │   │   │
│   │   │  │   Handler   │──>│   UseCase   │──>│   Repository    │   │   │   │
│   │   │  │ (HTTP/Chi)  │   │  (Login)    │   │   (GORM)        │   │   │   │
│   │   │  └─────────────┘   └──────┬──────┘   └────────┬────────┘   │   │   │
│   │   │                          │                    │            │   │   │
│   │   │                          │                    │            │   │   │
│   │   │                          ▼                    ▼            │   │   │
│   │   │                   ┌─────────────┐      ┌───────────────┐   │   │   │
│   │   │                   │ JWT Service │      │ Domain Entity │   │   │   │
│   │   │                   │(HMAC-SHA256)│      │    (User)     │   │   │   │
│   │   │                   └─────────────┘      └───────────────┘   │   │   │
│   │   │                                                              │   │   │
│   │   └──────────────────────────────────────────────────────────────┘   │   │
│   │                                                                      │   │
│   └──────────────────────────────────────────────────────────────────────┘   │
│                                   │                                          │
│                                   │ PostgreSQL (5432)                        │
│                                   ▼                                          │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                        RDS PostgreSQL                                │   │
│   │                       (VPC Private Subnet)                           │   │
│   │                                                                      │   │
│   │   ┌─────────────────────────────────────────────────────────────┐   │   │
│   │   │                    Tabela: usuarios                          │   │   │
│   │   │  • id, email, senha_hash, nome, role, created_at, updated_at │   │   │
│   │   └─────────────────────────────────────────────────────────────┘   │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │                         JWT TOKEN                                    │   │
│   │                                                                      │   │
│   │  Header: {"alg": "HS256", "typ": "JWT"}                             │   │
│   │  Payload: {sub, email, name, role, iat, exp, iss}                   │   │
│   │  Signature: HMACSHA256(base64(header).base64(payload), secret)      │   │
│   │                                                                      │   │
│   └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
```

## 🚀 Passos para Execução e Deploy

### Pré-requisitos

- [Go 1.23+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Terraform](https://www.terraform.io/) >= 1.5
- [AWS CLI](https://aws.amazon.com/cli/) configurado
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (opcional, para K8s)

### Variáveis de Ambiente

```bash
# Database
export DB_HOST=postgres.rds.amazonaws.com
export DB_PORT=5432
export DB_NAME=oficinapro
export DB_USER=oficinapro_app
export DB_PASSWORD=secret

# JWT
export JWT_SECRET=your-secret-key-min-32-chars
export JWT_ISSUER=auth-oficinapro
export JWT_EXPIRATION=3600

# Application
export ENVIRONMENT=production
```

### 1. Build Local

```bash
# Build binário
make build

# Ou manualmente
go build -o bin/auth-service ./cmd/http/
```

### 2. Executar Localmente

```bash
# Via Makefile
make run-local

# Ou diretamente
./bin/auth-service
```

### 3. Testes

```bash
# Todos os testes
make test

# Com cobertura
make test-coverage

# Testes específicos
go test -v ./internal/handler/...
go test -v ./internal/usecase/...
```

### 4. Deploy na AWS

```bash
# 1. Build para Lambda
make build-lambda

# 2. Deploy infraestrutura
cd terraform/lambda
terraform init
terraform plan
terraform apply
```

## 📊 Endpoints da API

| Método | Path | Descrição | Auth | Body |
|--------|------|-----------|------|------|
| `POST` | `/auth` | Login com email+senha | ❌ | `{email, senha}` |
| `POST` | `/auth/validate` | Validar token JWT | ✅ | - |
| `GET` | `/health` | Health check | ❌ | - |
| `OPTIONS` | `/*` | CORS preflight | ❌ | - |

### Autenticação

**Request:**
```bash
curl -X POST https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/auth \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "senha": "password123"}'
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 3600,
  "userId": 123,
  "email": "user@example.com",
  "nome": "User Name",
  "role": "USER"
}
```

**Response (401 Unauthorized):**
```json
{
  "error": "Credenciais inválidas",
  "timestamp": "2026-01-06T10:30:00Z"
}
```

### Validação de Token

```bash
curl -X POST https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/auth/validate \
  -H "Authorization: Bearer <jwt_token>"
```

### Uso do Token em Outros Serviços

```bash
# 1. Obter token
TOKEN=$(curl -s -X POST https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","senha":"senha123"}' | jq -r '.token')

# 2. Usar token no Core Service
curl -H "Authorization: Bearer $TOKEN" \
  http://3.89.176.39/api/ordem-servico
```

## 🧪 Testes e Cobertura

| Pacote | Cobertura |
|--------|-----------|
| `domain/entity` | 95% |
| `domain/service` | 92% |
| `usecase` | 90% |
| `handler` | 88% |
| `adapter` | 85% |
| **Total** | **91.8%** |

```bash
# Gerar relatório de cobertura
make test-coverage

# Visualizar no browser
open coverage.html
```

Veja o relatório completo em [`TEST_COVERAGE_REPORT.md`](TEST_COVERAGE_REPORT.md).

## 🔐 Segurança

| Aspecto | Implementação |
|---------|---------------|
| Password Hash | bcrypt (cost 10) |
| Token Signing | HMAC-SHA256 |
| Token Expiration | 1 hora (configurável) |
| Transport | HTTPS only (API Gateway) |
| Secrets | AWS Secrets Manager |
| Network | VPC, Security Groups |
| IAM | Least privilege |
| Validation | Email RFC 5322, senha min 6 chars |

## 📈 Performance

| Métrica | Valor |
|---------|-------|
| Cold start | ~500ms |
| Warm execution | ~50ms |
| Memory | 256MB |
| Timeout | 30s |
| Concurrency | Auto-scaling |

## 🛠️ Comandos Úteis (Makefile)

```bash
# Build
make build              # Build local
make build-lambda       # Build para Lambda
make docker-build       # Build Docker image

# Testes
make test              # Rodar testes
make test-coverage     # Coverage report
make test-integration  # Testes de integração

# Deploy
make deploy-infra      # Deploy Terraform
make deploy-k8s        # Deploy Kubernetes

# Limpeza
make clean            # Limpar binários
make clean-test       # Limpar coverage files

# Desenvolvimento
make fmt              # Format código
make lint             # Lint
make run-local        # Rodar localmente
```

## 🏆 Features

- ✅ Clean Architecture
- ✅ SOLID Principles
- ✅ Dependency Injection
- ✅ Bcrypt password hashing
- ✅ JWT HMAC-SHA256
- ✅ OpenTelemetry + New Relic ready
- ✅ Graceful shutdown
- ✅ Error handling centralizado
- ✅ CORS support
- ✅ Health checks
- ✅ **91.8% test coverage**
- ✅ Multi-stage Dockerfile
- ✅ Terraform IaC
- ✅ Kubernetes manifests
- ✅ CI/CD com GitHub Actions
- ✅ Dual mode: Lambda + HTTP server

## 📚 Documentação

### Documentação Local

| Documento | Descrição |
|-----------|-----------|
| [`CHANGELOG.md`](CHANGELOG.md) | Histórico de mudanças |
| [`DEPLOY_GUIDE.md`](DEPLOY_GUIDE.md) | Guia de deploy completo |
| [`TESTING_GUIDE.md`](TESTING_GUIDE.md) | Guia de testes e qualidade |
| [`TEST_COVERAGE_REPORT.md`](TEST_COVERAGE_REPORT.md) | Relatório de cobertura (91.8%) |
| [`terraform/lambda/README.md`](terraform/lambda/README.md) | Infraestrutura Lambda |
| [`k8s/order-service/README.md`](k8s/order-service/README.md) | Exemplo Kubernetes |

### Documentação Arquitetural (Centralizada)

| Documento | Descrição |
|-----------|-----------|
| [📊 Diagrama de Componentes](../core-domain-service/docs/architecture/COMPONENT_DIAGRAM.md) | Visão geral da arquitetura em nuvem |
| [📈 Diagramas de Sequência](../core-domain-service/docs/architecture/SEQUENCE_DIAGRAMS.md) | Fluxo de autenticação detalhado |
| [📝 RFC-003: Estratégia de Autenticação](../core-domain-service/docs/rfcs/RFC-003-AUTHENTICATION_STRATEGY.md) | Justificativa técnica do JWT + Lambda |
| [📋 ADR-003: Arquitetura Híbrida](../core-domain-service/docs/adrs/ADR-003-SERVERLESS_AUTH.md) | Decisão Lambda + K8s |
| [📚 Índice da Documentação](../core-domain-service/docs/INDEX.md) | Índice completo |

### Repositórios Relacionados

| Repositório | Descrição |
|-------------|-----------|
| [core-domain-service](../core-domain-service/) | Aplicação principal (Java/Spring Boot) |
| [OficinaPro-Database](../OficinaPro-Database/) | Infraestrutura do banco de dados |
| [OficinaPro-DevOps](../OficinaPro-DevOps/) | Infraestrutura e CI/CD |
| [OficinaPro-Payments](../OficinaPro-Payments/) | Serviço de pagamentos (Python) |

## 🔗 Links Úteis (Produção)

| Recurso | URL |
|---------|-----|
| **Auth Endpoint** | `https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/auth` |
| **Validate Endpoint** | `https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/auth/validate` |
| **Health Check** | `https://tzqb0s3zd0.execute-api.us-east-1.amazonaws.com/health` |
| **API Gateway Console** | [AWS Console](https://console.aws.amazon.com/apigateway) |
| **Lambda Console** | [AWS Console](https://console.aws.amazon.com/lambda) |
| **CloudWatch Logs** | [AWS Console](https://console.aws.amazon.com/cloudwatch) |

## 💰 Custos

| Recurso | Free Tier | Após Free Tier |
|---------|-----------|----------------|
| Lambda (1M requests) | Incluído | ~$0.20/M |
| API Gateway (1M requests) | Incluído | ~$1.00/M |
| CloudWatch Logs | 5GB | $0.50/GB |
| **Total estimado** | **$0.00** | **~$2/mês** |

## 👥 Equipe

Desenvolvido para o **Tech Challenge da FIAP**.

---

**Status**: ✅ Produção Ready
**Test Coverage**: 91.8%
**Runtime**: AWS Lambda (Go 1.23)
**Última atualização**: Janeiro de 2026
