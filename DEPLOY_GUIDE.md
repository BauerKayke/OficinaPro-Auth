# 🚀 Guia Rápido de Deploy - Auth Gateway Lambda

## ✅ Pré-requisitos

- [x] AWS CLI configurado (`aws configure`)
- [x] Terraform instalado (>= 1.5)
- [x] Go instalado (>= 1.23)
- [x] Make instalado
- [x] jq instalado (para testes)

## 📋 Checklist de Deploy

### 1. **Preparação da Infraestrutura de Database** (Primeiro Deploy)

⚠️ **IMPORTANTE**: Este repositório depende do `oficinapro-infra-db` para VPC, RDS e Secrets.

Se ainda não existe, você precisa:

**Opção A: Criar `oficinapro-infra-db` separadamente** (Recomendado para produção)
- Seguir `INFRASTRUCTURE_STRATEGY.md` para criar repositório de DB
- Deploy VPC + RDS primeiro
- Capturar outputs (VPC ID, Subnets, RDS endpoint, etc)

**Opção B: Deploy standalone** (Apenas para dev/testes)
- Configurar variáveis manualmente em `terraform.tfvars`
- Passar credenciais direto (não usar Secrets Manager)

---

## 🎯 Deploy Rápido (Opção B - Standalone Dev)

### Passo 1: Configurar Terraform

```bash
cd terraform/lambda

# Copiar exemplo de configuração
cp terraform.tfvars.example terraform.tfvars

# Editar terraform.tfvars com seus valores
vim terraform.tfvars
```

**Configuração mínima para dev** (`terraform.tfvars`):
```hcl
environment  = "dev"
aws_region   = "us-east-1"
project_name = "auth-oficinapro"

# Lambda
lambda_memory  = 512
lambda_timeout = 30

# VPC (Desabilitar se não tiver VPC/RDS ainda)
vpc_enabled = false

# Database (passar direto - apenas para DEV!)
db_host     = "localhost" # ou RDS endpoint se tiver
db_port     = 5432
db_name     = "oficinapro"
db_user     = "postgres"
db_password = "postgres"
db_ssl_mode = "disable"

# JWT (GERAR UM SECRET SEGURO!)
jwt_secret     = "your-secret-key-must-be-at-least-32-characters-long-change-this"
jwt_issuer     = "auth-oficinapro"
jwt_expiration = 3600

# Telemetry (Desabilitar se não tiver New Relic)
telemetry_enabled = false

# Monitoramento (Desabilitar alarmes se não tiver SNS)
enable_cloudwatch_alarms = false
```

### Passo 2: Build + Deploy

```bash
# Na raiz do projeto

# Opção 1: Deploy completo (recomendado)
make deploy-infra

# Opção 2: Passo a passo
make build                  # Build do binário Go
make terraform-init         # Inicializar Terraform
make terraform-plan         # Ver o que será criado
make terraform-apply        # Aplicar mudanças

# Opção 3: Quick deploy (auto-approve - cuidado!)
make redeploy
```

### Passo 3: Testar

```bash
# Testar automaticamente
make test-lambda

# Ou manualmente
export API_URL=$(cd terraform/lambda && terraform output -raw api_gateway_url)

# Health check
curl $API_URL/health

# Autenticação
curl -X POST $API_URL/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@oficinapro.com",
    "senha": "senha123"
  }' | jq
```

### Passo 4: Ver Logs

```bash
# Ver logs da Lambda em tempo real
make logs-lambda

# Ou diretamente
aws logs tail /aws/lambda/auth-oficinapro-dev --follow
```

---

## 🔐 Deploy com VPC + RDS (Opção A - Produção)

### Pré-requisitos

1. **Deploy do `oficinapro-infra-db` completo**
2. **Capturar outputs**:

```bash
cd oficinapro-infra-db/terraform/environments/dev

# Capturar outputs
terraform output -json > /tmp/infra-outputs.json

export VPC_ID=$(terraform output -raw vpc_id)
export SUBNET_IDS=$(terraform output -json private_subnet_ids)
export RDS_SG=$(terraform output -raw rds_security_group_id)
export DB_SECRET=$(terraform output -raw db_connection_secret_arn)
```

### Configuração `terraform.tfvars`

```hcl
environment  = "dev"
aws_region   = "us-east-1"
project_name = "auth-oficinapro"

# VPC (do oficinapro-infra-db)
vpc_enabled          = true
vpc_id               = "vpc-0abc123..."
private_subnet_ids   = ["subnet-0aaa...", "subnet-0bbb..."]
rds_security_group_id = "sg-0123456..."

# Database (via Secrets Manager - RECOMENDADO)
db_secret_arn = "arn:aws:secretsmanager:us-east-1:123456789:secret:db-creds-xxx"

# NÃO passar db_host, db_user, db_password quando usar Secrets Manager

# JWT
jwt_secret     = "GENERATE-A-STRONG-SECRET-32-CHARS-MIN"
jwt_issuer     = "auth-oficinapro"
jwt_expiration = 3600

# Telemetry
telemetry_enabled = true
newrelic_key      = "YOUR-NEW-RELIC-LICENSE-KEY"

# Monitoramento
enable_cloudwatch_alarms = true
alarm_sns_topic_arn       = "arn:aws:sns:us-east-1:123456789:alerts"
```

### Deploy

```bash
# Build + Deploy
make deploy-infra

# Testar
make test-lambda
```

---

## 🔄 Redeploy (Atualizar Código)

```bash
# Quick redeploy (apenas código Lambda)
make redeploy

# Full redeploy (com plan)
make build
make terraform-plan
make terraform-apply
```

---

## 🧹 Limpar (Destroy)

```bash
# Destruir infraestrutura
make terraform-destroy

# Limpar arquivos Terraform locais
make terraform-clean
```

---

## 📊 Monitoramento

### CloudWatch Logs

```bash
# Lambda logs
make logs-lambda

# API Gateway logs (se habilitado)
aws logs tail /aws/apigateway/auth-oficinapro-api-dev --follow
```

### CloudWatch Metrics

Acesse o console AWS:
- Lambda: `https://console.aws.amazon.com/lambda/home?region=us-east-1#/functions/auth-oficinapro-dev`
- API Gateway: Dashboard de métricas
- Alarmes: Se configurado, verifique SNS

---

## 🐛 Troubleshooting

### Erro: "Lambda binary not found"
```bash
# Build do binário primeiro
make build
ls -l bootstrap
```

### Erro: "too many connections" no RDS
```hcl
# Ajustar em terraform.tfvars
db_max_connections      = 2
db_max_idle_connections = 1
```

### Erro: "unable to connect to database"
1. Lambda está na mesma VPC do RDS?
2. Security Group do RDS permite conexões do SG da Lambda?
3. Credenciais corretas no Secrets Manager?

```bash
# Testar conexão manual
aws secretsmanager get-secret-value \
  --secret-id db-credentials-xxx \
  --query SecretString \
  --output text | jq
```

### Erro: "JWT secret too short"
```bash
# Gerar secret seguro (min 32 chars)
openssl rand -base64 32
```

---

## 📋 Workflow Completo (First Time)

```bash
# 1. Clone do repositório
git clone https://github.com/your-org/auth-oficinapro.git
cd auth-oficinapro

# 2. Instalar dependências
make install-deps

# 3. Rodar testes
make test

# 4. Build
make build

# 5. Configurar Terraform
cd terraform/lambda
cp terraform.tfvars.example terraform.tfvars
vim terraform.tfvars # Editar com seus valores

# 6. Deploy
cd ../..
make deploy-infra

# 7. Testar
make test-lambda

# 8. Ver logs
make logs-lambda
```

---

## 🔑 Variáveis Obrigatórias

| Variável | Obrigatório | Onde Obter |
|----------|-------------|-----------|
| `vpc_id` | Se `vpc_enabled=true` | oficinapro-infra-db |
| `private_subnet_ids` | Se `vpc_enabled=true` | oficinapro-infra-db |
| `rds_security_group_id` | Se `vpc_enabled=true` | oficinapro-infra-db |
| `db_secret_arn` | Se usar Secrets Manager | oficinapro-infra-db |
| `jwt_secret` | ✅ Sempre | Gerar com `openssl rand -base64 32` |

---

## 🎯 Próximos Passos

Após deploy bem-sucedido:

1. **Configurar DNS**: Apontar domínio para API Gateway
2. **Configurar TLS**: API Gateway custom domain + ACM certificate
3. **CI/CD**: GitHub Actions para deploy automatizado
4. **Testes E2E**: Automatizar testes pós-deploy
5. **Monitoramento**: Configurar dashboards e alertas

---

## 📚 Referências

- **Terraform Docs**: `terraform/lambda/README.md`
- **Estratégia de Infra**: `INFRASTRUCTURE_STRATEGY.md`
- **Arquitetura**: `HYBRID_ARCHITECTURE_LAMBDA_K8S.md`
- **Testes**: `TESTING_GUIDE.md`

