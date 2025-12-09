# 📊 Diagramas Visuais para Apresentação - Tech Challenge

## 🎯 Diagramas Prontos para Slides

---

## 1. Visão Geral da Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                    USUÁRIO / FRONTEND                           │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         │ HTTPS
                         ▼
            ┌────────────────────────────┐
            │      API GATEWAY           │
            │   (AWS API Gateway)        │
            │                            │
            │  • Roteamento              │
            │  • Rate Limiting           │
            │  • CORS                    │
            └────────┬──────────┬────────┘
                     │          │
            /auth    │          │ /api/*
                     │          │
        ┌────────────▼──┐   ┌──▼─────────────┐
        │  AUTH LAMBDA  │   │  CORE SERVICE  │
        │  ═══════════  │   │  ═════════════  │
        │  Go 1.21+     │   │  Java 21       │
        │  Serverless   │   │  Spring Boot   │
        │               │   │  Kubernetes    │
        │  • CPF        │   │                │
        │  • JWT        │   │  • Clientes    │
        │               │   │  • Veículos    │
        │               │   │  • Ordens      │
        └───────┬───────┘   └───────┬────────┘
                │                   │
                │ Read              │ Read/Write
                │                   │
                └──────┬────────────┘
                       │
                       ▼
        ┌──────────────────────────────┐
        │   RDS POSTGRESQL             │
        │   db.t3.micro                │
        │   20GB Storage               │
        │                              │
        │  Tabelas:                    │
        │  • pessoa, cliente           │
        │  • veiculo                   │
        │  • ordem_servico             │
        │  • peca, servico             │
        └──────────────────────────────┘
```

---

## 2. Separação de Repositórios

```
┌──────────────────────────────────────────────────────────────────┐
│                      4 REPOSITÓRIOS                              │
└──────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  1️⃣  auth-oficinapro (Go + Lambda)                              │
│  ───────────────────────────────────                            │
│  Responsabilidade: Autenticação                                 │
│  Deploy: AWS Lambda                                             │
│  Database: Read-only (pessoa, cliente)                          │
│  Status: ✅ 91% completo                                        │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  2️⃣  core-domain-service (Java + Spring Boot)                   │
│  ────────────────────────────────────────────                   │
│  Responsabilidade: Lógica de negócio                            │
│  Deploy: Kubernetes (K3s) em EC2                                │
│  Database: Read/Write (todas tabelas)                           │
│  Status: ✅ Clean Architecture + DDD                            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  3️⃣  OficinaPro-Database (Terraform)                            │
│  ───────────────────────────────────                            │
│  Responsabilidade: Provisionamento RDS                          │
│  Recursos: db.t3.micro + Security Groups                        │
│  Status: ✅ Terraform IaC pronto                                │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│  4️⃣  OficinaPro-DevOps (Terraform)                              │
│  ─────────────────────────────────                              │
│  Responsabilidade: Infraestrutura AWS                           │
│  Recursos: VPC + EC2 + K3s + Networking                         │
│  Status: ✅ Terraform IaC pronto                                │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. Fluxo de Autenticação (Sequence Diagram)

```
┌────────┐   ┌─────────────┐   ┌──────────────┐   ┌──────────┐
│ USER   │   │ API Gateway │   │ Auth Lambda  │   │ Database │
└───┬────┘   └──────┬──────┘   └──────┬───────┘   └────┬─────┘
    │               │                  │                │
    │ POST /auth    │                  │                │
    │ {cpf: "..."}  │                  │                │
    ├──────────────>│                  │                │
    │               │                  │                │
    │               │ Invoke Lambda    │                │
    │               ├─────────────────>│                │
    │               │                  │                │
    │               │                  │ 1. Valida CPF  │
    │               │                  │                │
    │               │                  │ 2. Query       │
    │               │                  │ SELECT cliente │
    │               │                  ├───────────────>│
    │               │                  │                │
    │               │                  │    Result      │
    │               │                  │<───────────────┤
    │               │                  │                │
    │               │                  │ 3. Gera JWT    │
    │               │                  │    (24h)       │
    │               │                  │                │
    │               │   Response       │                │
    │               │   {token, ...}   │                │
    │               │<─────────────────┤                │
    │               │                  │                │
    │  200 OK       │                  │                │
    │  {token, ...} │                  │                │
    │<──────────────┤                  │                │
    │               │                  │                │
    │  ✅ JWT       │                  │                │
    │  Armazenado   │                  │                │
    │               │                  │                │

TEMPO: ~10-50ms (warm) / ~100-200ms (cold start)
```

---

## 4. Fluxo de Ordem de Serviço (Sequence Diagram)

```
┌────────┐  ┌───────────┐  ┌─────────────┐  ┌──────────┐
│ USER   │  │ Gateway   │  │ Core Service│  │ Database │
│ (JWT)  │  │           │  │ (K8s Pod)   │  │          │
└───┬────┘  └─────┬─────┘  └──────┬──────┘  └────┬─────┘
    │             │                │              │
    │ POST /api/ordens-servico    │              │
    │ Authorization: Bearer JWT   │              │
    ├────────────>│                │              │
    │             │                │              │
    │             │ Forward +JWT   │              │
    │             ├───────────────>│              │
    │             │                │              │
    │             │                │1. Valida JWT │
    │             │                │              │
    │             │                │2. Extrai     │
    │             │                │   clienteId  │
    │             │                │              │
    │             │                │3. Business   │
    │             │                │   Logic      │
    │             │                │              │
    │             │                │4. BEGIN TX   │
    │             │                │ INSERT ordem │
    │             │                ├─────────────>│
    │             │                │              │
    │             │                │   Result     │
    │             │                │<─────────────┤
    │             │                │              │
    │             │                │5. COMMIT TX  │
    │             │                │              │
    │             │                │6. Event      │
    │             │                │   ORDEM_CRIADA│
    │             │                │              │
    │             │   201 Created  │              │
    │             │   {id, ...}    │              │
    │             │<───────────────┤              │
    │             │                │              │
    │ 201 Created │                │              │
    │ {id, ...}   │                │              │
    │<────────────┤                │              │
    │             │                │              │

TEMPO: ~100-500ms (depende da complexidade)
```

---

## 5. Database Access Pattern

```
┌─────────────────────────────────────────────────────────┐
│            RDS POSTGRESQL (db.t3.micro)                 │
│            ═══════════════════════════                  │
│                                                         │
│  Database: oficinapro                                   │
│  Engine: PostgreSQL 15.10                               │
│  Storage: 20GB GP2                                      │
│                                                         │
│  Tabelas:                                               │
│  ┌─────────────────────────────────────────┐           │
│  │ pessoa (id, nome, cpf, email)           │           │
│  │ cliente (id, pessoa_id, ativo)          │           │
│  │ veiculo (id, placa, cliente_id)         │           │
│  │ ordem_servico (id, veiculo_id, status)  │           │
│  │ peca, servico, ordem_servico_item, ...  │           │
│  └─────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────┘
            ▲                           ▲
            │                           │
    ┌───────┴────────┐         ┌────────┴────────┐
    │                │         │                 │
    │  READ ONLY     │         │  READ/WRITE     │
    │                │         │                 │
┌───┴────────────────┴──┐   ┌──┴──────────────────┴───┐
│  Auth Lambda (Go)     │   │  Core Service (Java)    │
│  ═════════════════     │   │  ═══════════════════     │
│                       │   │                         │
│  Acesso:              │   │  Acesso:                │
│  • pessoa (read)      │   │  • All tables (R/W)     │
│  • cliente (read)     │   │  • Liquibase migrations │
│                       │   │  • JPA/Hibernate        │
│  Connection Pool:     │   │                         │
│  • Max: 10            │   │  Connection Pool:       │
│  • Idle: 5            │   │  • Max: 50              │
│  • Timeout: 30s       │   │  • Idle: 10             │
│                       │   │  • Min Idle: 5          │
└───────────────────────┘   └─────────────────────────┘
```

---

## 6. Security & Authentication

```
┌──────────────────────────────────────────────────────────┐
│                   JWT FLOW                               │
└──────────────────────────────────────────────────────────┘

1️⃣  GERAÇÃO JWT (Auth Lambda)
   ┌─────────────────────────────────────┐
   │  JWT_SECRET (AWS Secrets Manager)   │
   │  <32+ caracteres aleatórios>        │
   └──────────────┬──────────────────────┘
                  │
                  ▼
   ┌─────────────────────────────────────┐
   │  jwt.NewWithClaims(                 │
   │    jwt.SigningMethodHS256,          │ ✅ HMAC-SHA256
   │    claims: {                        │
   │      cliente_id: 1,                 │
   │      nome: "João",                  │
   │      email: "joao@...",             │
   │      documento: "12345678909",      │
   │      iat: 1701234567,               │
   │      exp: 1701320967,  // +24h      │
   │      iss: "oficinapro-auth"         │
   │    }                                │
   │  )                                  │
   └──────────────┬──────────────────────┘
                  │
                  ▼
   ┌─────────────────────────────────────┐
   │  token.SignedString(secretKey)      │
   │  → eyJhbGciOiJIUzI1NiIsInR5cCI...  │
   └─────────────────────────────────────┘


2️⃣  VALIDAÇÃO JWT (Core Service)
   ┌─────────────────────────────────────┐
   │  JWT_SECRET (AWS Secrets Manager)   │
   │  <MESMO secret do Lambda>           │
   └──────────────┬──────────────────────┘
                  │
                  ▼
   ┌─────────────────────────────────────┐
   │  JwtAuthenticationFilter:           │
   │                                     │
   │  1. Extrai token do header          │
   │     Authorization: Bearer eyJ...    │
   │                                     │
   │  2. Valida assinatura               │
   │     if (!token.verify(secret))      │ ✅ HMAC-SHA256
   │       throw UnauthorizedException   │
   │                                     │
   │  3. Verifica expiração              │
   │     if (now > token.exp)            │
   │       throw TokenExpiredException   │
   │                                     │
   │  4. Extrai claims                   │
   │     clienteId = token.cliente_id    │
   │                                     │
   │  5. Popula SecurityContext          │
   │     SecurityContext.setAuthentication()│
   │                                     │
   └─────────────────────────────────────┘
```

---

## 7. Infrastructure Overview (AWS)

```
┌────────────────────────────────────────────────────────────────┐
│                         AWS CLOUD                              │
│  ══════════════════════════════════════════════════════════    │
│                                                                │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │               VPC (10.0.0.0/16)                          │ │
│  │               ═══════════════════                         │ │
│  │                                                           │ │
│  │  ┌─────────────────┐       ┌──────────────────┐         │ │
│  │  │ Public Subnet 1 │       │ Public Subnet 2  │         │ │
│  │  │ (AZ-1a)         │       │ (AZ-1b)          │         │ │
│  │  │                 │       │                  │         │ │
│  │  │ ┌─────────────┐ │       │ ┌──────────────┐ │         │ │
│  │  │ │ EC2 t3.micro│ │       │ │RDS db.t3.micro│ │         │ │
│  │  │ │             │ │       │ │              │ │         │ │
│  │  │ │ K3s Cluster │◄┼───────┼─┤ PostgreSQL   │ │         │ │
│  │  │ │  • Core App │ │       │ │ 20GB Storage │ │         │ │
│  │  │ │  • Redis    │ │       │ └──────────────┘ │         │ │
│  │  │ │  • Nginx    │ │       │                  │         │ │
│  │  │ │             │ │       │  Private Subnet  │         │ │
│  │  │ │ Elastic IP  │ │       │  (no public IP)  │         │ │
│  │  │ └─────────────┘ │       │                  │         │ │
│  │  └─────────────────┘       └──────────────────┘         │ │
│  │                                                           │ │
│  │  Internet Gateway                                         │ │
│  └───────────┬───────────────────────────────────────────────┘ │
│              │                                                 │
│  ┌───────────▼──────────┐      ┌──────────────────────┐      │
│  │  API Gateway         │      │  Lambda              │      │
│  │  • /auth → Lambda    │      │  auth-oficinapro     │      │
│  │  • /api/* → EIP      │      │  (Go 1.21)           │      │
│  └──────────────────────┘      │  128MB Memory        │      │
│                                │  30s Timeout         │      │
│                                └──────────────────────┘      │
│                                                                │
│  MANAGED BY:                                                   │
│  • VPC/Subnets/EC2 → OficinaPro-DevOps (Terraform)           │
│  • RDS → OficinaPro-Database (Terraform)                      │
│  • Lambda → auth-oficinapro (SAM/Manual)                      │
│  • K3s Apps → core-domain-service (kubectl)                   │
└────────────────────────────────────────────────────────────────┘
```

---

## 8. Deployment Flow

```
┌────────────────────────────────────────────────────────────────┐
│                   DEPLOYMENT ORDER                             │
└────────────────────────────────────────────────────────────────┘

STEP 1: BOOTSTRAP
┌─────────────────────────────────────────┐
│  cd OficinaPro-DevOps/bootstrap         │
│  terraform init                         │
│  terraform apply                        │
│                                         │
│  Cria:                                  │
│  ✅ S3 Bucket (tfstate storage)         │
│  ✅ DynamoDB (state locking)            │
└─────────────────────────────────────────┘
                  │
                  ▼
STEP 2: DATABASE
┌─────────────────────────────────────────┐
│  cd OficinaPro-Database                 │
│  terraform init -backend-config=...     │
│  terraform apply                        │
│                                         │
│  Cria:                                  │
│  ✅ RDS PostgreSQL (db.t3.micro)        │
│  ✅ Security Group (RDS)                │
│  ✅ Subnet Group                        │
│                                         │
│  Outputs:                               │
│  • db_instance_address                  │
│  • db_instance_name                     │
│  • db_credentials                       │
└─────────────────────────────────────────┘
                  │
                  ▼
STEP 3: APP INFRASTRUCTURE
┌─────────────────────────────────────────┐
│  cd OficinaPro-DevOps/app-infra         │
│  terraform init -backend-config=...     │
│  terraform apply                        │
│                                         │
│  Cria:                                  │
│  ✅ VPC + Subnets + IGW                 │
│  ✅ Security Groups                     │
│  ✅ EC2 t3.micro (K3s)                  │
│  ✅ Elastic IP                          │
│  ✅ User Data (install K3s)             │
│                                         │
│  Outputs:                               │
│  • k3s_node_public_ip                   │
│  • vpc_id                               │
└─────────────────────────────────────────┘
                  │
                  ▼
STEP 4: AUTH LAMBDA
┌─────────────────────────────────────────┐
│  cd auth-oficinapro                     │
│  make build                             │
│  make deploy                            │
│                                         │
│  Cria:                                  │
│  ✅ Lambda Function (Go)                │
│  ✅ Lambda Execution Role               │
│  ✅ API Gateway Integration             │
│                                         │
│  Env Vars (from Terraform outputs):     │
│  • DB_HOST, DB_NAME, DB_USER            │
│  • JWT_SECRET (Secrets Manager)         │
└─────────────────────────────────────────┘
                  │
                  ▼
STEP 5: CORE SERVICE
┌─────────────────────────────────────────┐
│  cd core-domain-service                 │
│  kubectl apply -f deployment/kubernetes/│
│                                         │
│  Cria:                                  │
│  ✅ Deployment (Java Spring Boot)       │
│  ✅ Service (ClusterIP)                 │
│  ✅ Ingress (Traefik)                   │
│  ✅ ConfigMap (DB connection)           │
│  ✅ Secret (credentials)                │
└─────────────────────────────────────────┘
```

---

## 9. Cost Breakdown (Free Tier)

```
┌────────────────────────────────────────────────────────────────┐
│                     CUSTOS - FREE TIER                         │
└────────────────────────────────────────────────────────────────┘

┌──────────────────┬──────────────┬─────────────┬──────────────┐
│   SERVIÇO        │  FREE TIER   │  USO ATUAL  │  CUSTO/MÊS   │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ Auth Lambda      │ 1M requests  │ ~100K       │    $0 ✅     │
│                  │ 400K GB-sec  │ ~10K GB-sec │              │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ EC2 t3.micro     │ 750h/mês     │ 750h        │    $0 ✅     │
│ (K3s)            │              │             │              │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ RDS db.t3.micro  │ 750h/mês     │ 750h        │    $0 ✅     │
│                  │ 20GB storage │ 20GB        │              │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ EBS Storage      │ 30GB/mês     │ 30GB        │    $0 ✅     │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ Data Transfer    │ 15GB/mês     │ <5GB        │    $0 ✅     │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ CloudWatch Logs  │ 5GB/mês      │ <1GB        │    $0 ✅     │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ S3 (tfstate)     │ 5GB/mês      │ <1MB        │    $0 ✅     │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ Elastic IP       │ Always free  │ 1 IP        │    $0 ✅     │
│ (associated)     │              │             │              │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ API Gateway      │ ❌ NOT FREE  │ ~100K req   │  ~$3.50 ⚠️   │
├──────────────────┼──────────────┼─────────────┼──────────────┤
│ TOTAL            │              │             │  ~$3.50/MÊS  │
└──────────────────┴──────────────┴─────────────┴──────────────┘

PÓS-FREE TIER (após 12 meses):
• EC2 t3.micro: $7.50/mês (Spot: $2.50/mês)
• RDS db.t3.micro: $15/mês
• Lambda: ~$0.20/mês
• Total: ~$26-30/mês
```

---

## 10. Benefits Summary

```
┌────────────────────────────────────────────────────────────────┐
│                      BENEFÍCIOS                                │
└────────────────────────────────────────────────────────────────┘

SEPARAÇÃO DE RESPONSABILIDADES
┌─────────────────────────────────────┐
│ ✅ Auth: Microserviço independente  │ → Lambda (Go)
│ ✅ Core: Lógica de negócio isolada  │ → K8s (Java)
│ ✅ DB: Infraestrutura gerenciada    │ → RDS (Terraform)
│ ✅ Infra: IaC versionado            │ → Terraform
└─────────────────────────────────────┘

ESCALABILIDADE
┌─────────────────────────────────────┐
│ ✅ Lambda: Automática (AWS)         │ → 0 a 1000s instâncias
│ ✅ K8s: HPA (Horizontal Scaling)    │ → 1 a 10 pods
│ ✅ Database: Connection pooling     │ → 60 conexões
└─────────────────────────────────────┘

CUSTO-BENEFÍCIO
┌─────────────────────────────────────┐
│ ✅ Free Tier: $0/mês (12 meses)     │
│ ✅ Pós-Free: ~$26-30/mês            │
│ ✅ ROI: 85% economia vs. tradicional│
└─────────────────────────────────────┘

CLEAN ARCHITECTURE
┌─────────────────────────────────────┐
│ ✅ Auth: SOLID 9.1/10               │
│ ✅ Core: DDD + CQRS + Clean Arch    │
│ ✅ Testabilidade: 85%+ coverage     │
│ ✅ Documentação: 40+ docs           │
└─────────────────────────────────────┘

SEGURANÇA
┌─────────────────────────────────────┐
│ ✅ JWT: HMAC-SHA256 shared secret   │
│ ✅ Database: Private subnet         │
│ ✅ SQL Injection: Protected (ORM)   │
│ ✅ Secrets: AWS Secrets Manager     │
└─────────────────────────────────────┘

OBSERVABILIDADE
┌─────────────────────────────────────┐
│ ✅ Logs: CloudWatch (central)       │
│ ✅ Metrics: Spring Actuator         │
│ ✅ Health: /actuator/health         │
│ ✅ Tracing: Ready for OpenTelemetry │
└─────────────────────────────────────┘
```

---

## 📝 Como Usar os Diagramas

### Para Slides (PowerPoint/Google Slides)

1. **Copie os diagramas ASCII** diretamente
2. Use fonte **monospace** (Courier New, Consolas, Monaco)
3. Tamanho da fonte: **8-10pt** para legibilidade
4. Background: Preto com texto verde (terminal style) ou branco

### Para Demo ao Vivo

1. **Abra este arquivo** no terminal ou IDE
2. Use **zoom** para mostrar cada diagrama
3. Navegue pelos diagramas durante a apresentação
4. Explique cada componente apontando para o diagrama

### Para Documentação

Estes diagramas já estão prontos para:
- README.md dos repositórios
- Wiki do projeto
- Confluence
- Documentação técnica

---

## 🎯 Ordem Sugerida para Apresentação

### Slide 1: Visão Geral
- Use **Diagrama 1** (Visão Geral da Arquitetura)
- Mostre os 4 componentes principais
- Destaque a separação de responsabilidades

### Slide 2: Repositórios
- Use **Diagrama 2** (Separação de Repositórios)
- Explique cada repositório
- Mostre tecnologias usadas

### Slide 3: Fluxo de Autenticação
- Use **Diagrama 3** (Fluxo de Autenticação)
- Mostre o sequence diagram
- Destaque tempo de resposta (~10-50ms)

### Slide 4: Fluxo de Negócio
- Use **Diagrama 4** (Fluxo de Ordem de Serviço)
- Mostre integração completa
- Destaque validação JWT

### Slide 5: Database
- Use **Diagrama 5** (Database Access Pattern)
- Explique compartilhamento temporário
- Mostre connection pooling

### Slide 6: Segurança
- Use **Diagrama 6** (Security & Authentication)
- Explique JWT flow completo
- Destaque secret compartilhado

### Slide 7: Infraestrutura
- Use **Diagrama 7** (Infrastructure Overview)
- Mostre AWS resources
- Destaque Terraform IaC

### Slide 8: Deploy
- Use **Diagrama 8** (Deployment Flow)
- Mostre ordem de deploy
- Explique dependências

### Slide 9: Custos
- Use **Diagrama 9** (Cost Breakdown)
- Mostre Free Tier
- Destaque economia

### Slide 10: Benefícios
- Use **Diagrama 10** (Benefits Summary)
- Recapitule vantagens
- Conclua com métricas

---

**"Arquitetura cloud-native com diagramas profissionais!"** 📊🚀

---

*Diagramas prontos para apresentação*
*Status: ✅ COMPLETO*
*Use com confiança no Tech Challenge!*



