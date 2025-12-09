# 📚 Índice Completo - Auth Service

## 🎯 Start Here

1. **[README.md](./README.md)** - Overview geral do projeto
2. **[REFACTORING_COMPLETE.md](./REFACTORING_COMPLETE.md)** - Resumo da refatoração
3. **[CHANGELOG.md](./CHANGELOG.md)** - Histórico de versões

---

## 📖 Documentação Técnica

### Arquitetura
- **[CLEAN_ARCHITECTURE.md](./docs/CLEAN_ARCHITECTURE.md)** - Guia completo de Clean Architecture
- **[SOLID_ANALYSIS.md](./docs/SOLID_ANALYSIS.md)** - Análise SOLID 97.3%
- **[ARCHITECTURE_IMPROVEMENTS.md](./docs/ARCHITECTURE_IMPROVEMENTS.md)** - Melhorias implementadas
- **[ARCHITECTURE_SUMMARY.md](./ARCHITECTURE_SUMMARY.md)** - Resumo da arquitetura simplificada

### API
- **[API.md](./docs/API.md)** - Documentação completa da API REST
- **[Postman Collection](./assets/)** - Collections para testes (se disponível)

### Implementation
- **[IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)** - Resumo completo da implementação

---

## 🏗️ Estrutura de Código

### Entry Point
```
cmd/lambda/main.go              # 25 linhas - Entry point
```

### Dependency Injection (Raiz)
```
di/
├── container.go                # Container principal
├── infrastructure.go           # Database, cache
├── repositories.go             # Repositórios
├── services.go                 # JWT, validators
├── usecases.go                 # Use cases
└── config/
    └── config.go               # Configuration loader
```

### Domain Layer (Puro)
```
internal/domain/
├── entity/
│   ├── cliente.go              # Entidade Cliente
│   └── errors.go               # Domain errors
├── repository/
│   └── cliente_repository.go  # Repository interface
└── service/
    ├── jwt_service.go          # JWT interface
    └── validator_service.go    # Validator interface
```

### Use Cases
```
internal/usecase/
└── authenticate_usecase.go     # Autenticação
```

### Handler Layer
```
internal/handler/
├── lambda_handler.go           # Lambda handler
├── errors.go                   # Error mapping
├── request/
│   └── auth_request.go         # Request DTOs
└── response/
    ├── auth_response.go        # Success response
    └── error_response.go       # Error response
```

### Infrastructure
```
internal/infrastructure/
├── database/
│   ├── models.go               # GORM models
│   ├── gorm_connection.go      # DB connection
│   └── gorm_cliente_repository.go  # Repository impl
├── jwt/
│   └── jwt_service_impl.go     # JWT implementation
└── validator/
    └── cpf_validator.go        # CPF validation
```

---

## 🧪 Testes

```
internal/
├── infrastructure/validator/
│   └── cpf_validator_test.go
└── usecase/
    └── authenticate_usecase_test.go
```

---

## 🛠️ DevOps

### Docker
- **[Dockerfile](./Dockerfile)** - Multi-stage build otimizado
- **[docker-compose.yml](./docker-compose.yml)** - Desenvolvimento local
- **[scripts/init.sql](./scripts/init.sql)** - Inicialização do DB

### CI/CD
- **[.github/workflows/ci.yml](./.github/workflows/ci.yml)** - Pipeline CI
- **[.github/workflows/cd.yml](./.github/workflows/cd.yml)** - Pipeline CD

### Build
- **[Makefile](./Makefile)** - Comandos de build e deploy
- **[go.mod](./go.mod)** - Dependências Go
- **[.gitignore](./.gitignore)** - Arquivos ignorados

---

## 📊 Métricas e Qualidade

| Documento | Foco |
|-----------|------|
| **SOLID_ANALYSIS.md** | Análise SOLID completa (97.3%) |
| **CLEAN_ARCHITECTURE.md** | Implementação Clean Architecture |
| **ARCHITECTURE_IMPROVEMENTS.md** | Comparação antes/depois |
| **REFACTORING_COMPLETE.md** | Resultados da refatoração |

---

## 🎓 Guias de Uso

### Quick Start
```bash
# 1. Setup
cd auth-oficinapro
make install-deps

# 2. Config
cp .env.example .env

# 3. Run
docker-compose up -d

# 4. Test
make test
```

### Deploy
```bash
# Build
make build

# Deploy Staging
make deploy-staging

# Deploy Production
make deploy-prod
```

### Development
```bash
# Lint
make lint

# Coverage
make test-coverage

# Clean
make clean
```

---

## 📝 Convenções

### Logs
- **Produção**: Apenas logs essenciais (2 por request)
- **Desenvolvimento**: Logs WARN do GORM
- **Erros**: Sempre com contexto (`fmt.Errorf`)

### Naming
- **Files**: snake_case (Go convention)
- **Packages**: lowercase, singular
- **Interfaces**: Suffix "able" ou noun
- **Implementations**: Suffix "Impl" ou descriptive

### Architecture
- **Domain**: Zero dependências externas
- **Use Cases**: Dependem de interfaces domain
- **Infrastructure**: Implementa interfaces domain
- **Handler**: Coordena request/response

---

## 🔗 Links Úteis

### Repositórios Relacionados
- **Core Domain Service**: `/core-domain-service`
- **Orders API**: (pendente)
- **Infra Kubernetes**: (pendente)
- **Infra Database**: (pendente)

### Documentação Externa
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [GORM Guide](https://gorm.io/docs/)
- [AWS Lambda Go](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)

---

## 📞 Suporte

### Issues
- GitHub Issues: `/issues`

### Documentation
- Todos os docs em `/docs`
- README principal: `/README.md`

### Contact
- Tech Challenge FIAP - Fase 4

---

## ✅ Checklist de Documentação

- [x] README.md completo
- [x] CHANGELOG.md criado
- [x] Clean Architecture documentado
- [x] SOLID Analysis completo
- [x] API documentada
- [x] Refatoração documentada
- [x] Índice criado
- [x] Quick start guide
- [x] Deploy guide
- [x] Architecture diagrams

---

**Status da Documentação**: ✅ **100% Completa**

*Última atualização: 20/10/2025*

