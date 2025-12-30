# Auth Gateway - OficinaPro

Auth Gateway usando AWS Lambda para autenticação via email+senha com JWT.

## 🎯 Visão Geral

Serviço de autenticação serverless que:
- Valida credenciais (email + senha com bcrypt)
- Gera tokens JWT para autenticação
- Integra-se com infraestrutura existente (VPC, RDS, EKS)
- Roda como AWS Lambda + API Gateway

## 🏗️ Arquitetura

```
Cliente
  ↓
API Gateway (HTTPS)
  ↓
Lambda (Auth Gateway)
  ↓
RDS PostgreSQL
  ↓
JWT Token ← Cliente → Kubernetes Services
```

**Arquitetura Híbrida**:
- **Lambda**: Autenticação (stateless, escalável, pay-per-use)
- **Kubernetes**: Aplicações de negócio (stateful, validam JWT localmente)

## 📦 Estrutura do Projeto

```
auth-oficinapro/
├── cmd/lambda/              # Entrypoint Lambda
├── internal/                # Código da aplicação
│   ├── domain/             # Entidades, repositórios, serviços
│   ├── usecase/            # Casos de uso
│   ├── handler/            # HTTP handlers
│   └── infrastructure/     # GORM, JWT, Telemetry
├── di/                      # Dependency Injection
├── terraform/lambda/        # Infraestrutura Lambda
├── k8s/order-service/       # Manifests Kubernetes exemplo
├── examples/order-service/  # Serviço Go exemplo
├── scripts/                 # Scripts SQL
├── Dockerfile              # Build Lambda
├── Makefile               # Comandos úteis
└── README.md
```

## 🚀 Quick Start

### Pré-requisitos
- Go 1.23+
- Docker
- Terraform
- AWS CLI configurado
- kubectl (para Kubernetes)

### 1. Build Local
```bash
make build
```

### 2. Testes
```bash
make test
make test-coverage
```

### 3. Deploy
```bash
# Deploy infraestrutura (Lambda + API Gateway)
cd terraform/lambda
terraform init
terraform apply

# Deploy serviços Kubernetes (opcional)
kubectl apply -k k8s/order-service/
```

## 🔐 Autenticação

### Endpoint: POST /auth

**Request**:
```json
{
  "email": "user@example.com",
  "senha": "senha123"
}
```

**Response (200 OK)**:
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

**Response (401 Unauthorized)**:
```json
{
  "error": "Credenciais inválidas",
  "timestamp": "2024-12-20T10:30:00Z"
}
```

### Usar JWT em outros serviços

```bash
# Obter token
TOKEN=$(curl -s -X POST https://api.oficinapro.com/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","senha":"senha123"}' | jq -r '.token')

# Usar token
curl -H "Authorization: Bearer $TOKEN" \
  https://orders.oficinapro.com/orders
```

## 🧪 Testes

```bash
# Rodar todos os testes
make test

# Com coverage
make test-coverage

# Específico
go test -v ./internal/handler/...
go test -v ./internal/domain/entity/...
```

**Cobertura de Testes**:
- ✅ **91.8%** nos pacotes de negócio (domain, usecase, handler)
- ✅ 150+ casos de teste
- ✅ Testes unitários e de integração
- ✅ Ver [`TEST_COVERAGE_REPORT.md`](TEST_COVERAGE_REPORT.md) para detalhes

## 📊 Endpoints

| Método | Path | Descrição | Auth |
|--------|------|-----------|------|
| POST | `/auth` | Autenticação email+senha | ❌ |
| GET | `/health` | Health check | ❌ |
| POST | `/auth/validate` | Validar JWT | ✅ |
| OPTIONS | `/*` | CORS preflight | ❌ |

## 🔧 Configuração

### Variáveis de Ambiente

```bash
# Database
DB_HOST=postgres.rds.amazonaws.com
DB_PORT=5432
DB_NAME=oficinapro
DB_USER=oficinapro_app
DB_PASSWORD=secret

# JWT
JWT_SECRET=your-secret-key-min-32-chars
JWT_ISSUER=auth-oficinapro
JWT_EXPIRATION=3600  # segundos

# Application
ENVIRONMENT=production
```

## 🔗 Integração com app-infra

Este projeto integra-se com o repositório `OficinaPro-DevOps/app-infra` que provisiona:
- VPC e Subnets
- RDS PostgreSQL
- EKS Cluster
- Security Groups

## 📚 Documentação

| Documento | Descrição |
|-----------|-----------|
| [`CHANGELOG.md`](CHANGELOG.md) | Histórico de mudanças |
| [`DEPLOY_GUIDE.md`](DEPLOY_GUIDE.md) | Guia de deploy completo |
| [`TESTING_GUIDE.md`](TESTING_GUIDE.md) | Guia de testes e qualidade |
| [`TEST_COVERAGE_REPORT.md`](TEST_COVERAGE_REPORT.md) | Relatório de cobertura (91.8%) |
| [`terraform/lambda/README.md`](terraform/lambda/README.md) | Infraestrutura Lambda |
| [`k8s/order-service/README.md`](k8s/order-service/README.md) | Exemplo Kubernetes |

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

## 📈 Performance

- **Cold start**: ~500ms
- **Warm execution**: ~50ms
- **Memory**: 256MB-512MB
- **Timeout**: 30s
- **Concurrency**: Auto-scaling

## 🔒 Segurança

- ✅ Bcrypt (cost 10)
- ✅ JWT HMAC-SHA256
- ✅ AWS Secrets Manager
- ✅ VPC integration
- ✅ Security Groups
- ✅ IAM least privilege
- ✅ HTTPS only (API Gateway)
- ✅ Password validation (min 6 chars)
- ✅ Email validation

## 📝 Licença

Projeto acadêmico - FIAP Tech Challenge

## 👥 Equipe

Desenvolvido para o Tech Challenge da FIAP

---

**Status**: ✅ Produção Ready
**Test Coverage**: 91.8%
**Última atualização**: 21 de dezembro de 2025
