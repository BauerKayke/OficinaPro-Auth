# 🏗️ Proposta de Refatoração: DI Container

## 📋 Análise de Problemas Atuais

### ❌ Violações de SOLID Identificadas

| Princípio | Violação | Impacto | Severidade |
|-----------|----------|---------|------------|
| **SRP** | Container tem 6+ responsabilidades | God Object anti-pattern | 🔴 CRÍTICO |
| **OCP** | Adicionar componente = modificar Container | Não extensível | 🔴 CRÍTICO |
| **DIP** | Container depende de implementações concretas | Baixa testabilidade | 🟡 MÉDIO |
| **ISP** | Container expõe tudo para todos | Acoplamento excessivo | 🟡 MÉDIO |

### ⚠️ Outros Problemas

1. **Error handling inconsistente**: métodos retornam `error` mas nunca falham
2. **Acoplamento temporal**: ordem de inicialização importa
3. **Baixa testabilidade**: difícil mockar e testar
4. **Sem lazy loading**: ruim para serverless (cold start)
5. **God Object**: Container conhece tudo sobre tudo

---

## ✅ Solução Proposta: Builder Pattern + Functional Options

### Arquitetura Nova

```
┌─────────────────────────────────────────────────────┐
│                   Container                         │
│  (Apenas armazena dependências, não cria)          │
└─────────────────────────────────────────────────────┘
                      ▲
                      │
            ┌─────────┴─────────┐
            │   ContainerBuilder │
            │  (Orquestração)    │
            └─────────┬──────────┘
                      │
      ┌───────────────┼───────────────┐
      │               │               │
┌─────▼─────┐  ┌─────▼─────┐  ┌─────▼─────┐
│ Infra      │  │ Repository│  │  Service  │
│ Builder    │  │ Builder   │  │  Builder  │
└────────────┘  └───────────┘  └───────────┘
```

### Benefícios

1. ✅ **SRP**: Cada builder tem UMA responsabilidade
2. ✅ **OCP**: Adicionar componente = criar novo builder (não modifica existente)
3. ✅ **DIP**: Builders dependem de interfaces, não implementações
4. ✅ **Testabilidade**: Cada builder pode ser testado isoladamente
5. ✅ **Lazy Loading**: Inicializa apenas o necessário
6. ✅ **Flexibilidade**: Functional options para configurar

---

## 🔧 Implementação

### 1. Container Simplificado

```go
// di/container.go
package di

import (
    "context"
    "github.com/oficinapro/auth-service/internal/domain/repository"
    "github.com/oficinapro/auth-service/internal/domain/service"
    "github.com/oficinapro/auth-service/internal/usecase"
)

// Container holds application dependencies (APENAS armazena, não cria)
type Container struct {
    // Domain interfaces (não implementações)
    clienteRepo      repository.ClienteRepository
    jwtService       service.JWTService
    validatorService service.ValidatorService
    authenticateUC   usecase.AuthenticateUseCase

    // Lifecycle management
    closeFunc func() error
}

// GetClienteRepository returns the cliente repository
func (c *Container) GetClienteRepository() repository.ClienteRepository {
    return c.clienteRepo
}

// GetJWTService returns the JWT service
func (c *Container) GetJWTService() service.JWTService {
    return c.jwtService
}

// GetValidatorService returns the validator service
func (c *Container) GetValidatorService() service.ValidatorService {
    return c.validatorService
}

// GetAuthenticateUseCase returns the authenticate use case
func (c *Container) GetAuthenticateUseCase() usecase.AuthenticateUseCase {
    return c.authenticateUC
}

// Close closes all resources
func (c *Container) Close(ctx context.Context) error {
    if c.closeFunc != nil {
        return c.closeFunc()
    }
    return nil
}
```

**Diferenças**:
- ✅ Container **NÃO cria** nada, apenas **armazena**
- ✅ Expõe apenas **interfaces** (não implementações concretas)
- ✅ Getters explícitos (melhor encapsulamento)
- ✅ Não expõe `*gorm.DB` ou outros detalhes de infra

---

### 2. ContainerBuilder (Orquestração)

```go
// di/builder.go
package di

import (
    "context"
    "fmt"

    "gorm.io/gorm"
    "github.com/oficinapro/auth-service/di/config"
    "github.com/oficinapro/auth-service/internal/domain/repository"
    "github.com/oficinapro/auth-service/internal/domain/service"
    "github.com/oficinapro/auth-service/internal/usecase"
)

// ContainerBuilder builds the DI container
type ContainerBuilder struct {
    config           *config.Config
    db               *gorm.DB
    clienteRepo      repository.ClienteRepository
    jwtService       service.JWTService
    validatorService service.ValidatorService
    authenticateUC   usecase.AuthenticateUseCase
    closeFunc        func() error
}

// NewContainerBuilder creates a new builder
func NewContainerBuilder() *ContainerBuilder {
    return &ContainerBuilder{}
}

// WithConfig sets the configuration
func (b *ContainerBuilder) WithConfig(cfg *config.Config) *ContainerBuilder {
    b.config = cfg
    return b
}

// WithDatabase sets the database connection
func (b *ContainerBuilder) WithDatabase(db *gorm.DB, closeFunc func() error) *ContainerBuilder {
    b.db = db
    b.closeFunc = closeFunc
    return b
}

// WithClienteRepository sets the cliente repository
func (b *ContainerBuilder) WithClienteRepository(repo repository.ClienteRepository) *ContainerBuilder {
    b.clienteRepo = repo
    return b
}

// WithJWTService sets the JWT service
func (b *ContainerBuilder) WithJWTService(svc service.JWTService) *ContainerBuilder {
    b.jwtService = svc
    return b
}

// WithValidatorService sets the validator service
func (b *ContainerBuilder) WithValidatorService(svc service.ValidatorService) *ContainerBuilder {
    b.validatorService = svc
    return b
}

// WithAuthenticateUseCase sets the authenticate use case
func (b *ContainerBuilder) WithAuthenticateUseCase(uc usecase.AuthenticateUseCase) *ContainerBuilder {
    b.authenticateUC = uc
    return b
}

// Build creates the container
func (b *ContainerBuilder) Build() (*Container, error) {
    if b.clienteRepo == nil {
        return nil, fmt.Errorf("clienteRepo is required")
    }
    if b.jwtService == nil {
        return nil, fmt.Errorf("jwtService is required")
    }
    if b.validatorService == nil {
        return nil, fmt.Errorf("validatorService is required")
    }
    if b.authenticateUC == nil {
        return nil, fmt.Errorf("authenticateUC is required")
    }

    return &Container{
        clienteRepo:      b.clienteRepo,
        jwtService:       b.jwtService,
        validatorService: b.validatorService,
        authenticateUC:   b.authenticateUC,
        closeFunc:        b.closeFunc,
    }, nil
}
```

**Benefícios**:
- ✅ **Fluent Interface** (fácil de usar)
- ✅ **Validação** antes de construir
- ✅ **Flexível**: pode injetar mocks para testes
- ✅ **Separação de responsabilidades**: builder vs container

---

### 3. Factories (Criação de Dependências)

```go
// di/factories/infrastructure_factory.go
package factories

import (
    "fmt"
    "gorm.io/gorm"
    "github.com/oficinapro/auth-service/di/config"
    "github.com/oficinapro/auth-service/internal/infrastructure/database"
)

// InfrastructureFactory creates infrastructure dependencies
type InfrastructureFactory struct {
    config *config.Config
}

// NewInfrastructureFactory creates a new infrastructure factory
func NewInfrastructureFactory(cfg *config.Config) *InfrastructureFactory {
    return &InfrastructureFactory{config: cfg}
}

// CreateDatabase creates database connection
func (f *InfrastructureFactory) CreateDatabase() (*gorm.DB, func() error, error) {
    dbConfig := database.Config{
        Host:            f.config.Database.Host,
        Port:            f.config.Database.Port,
        Database:        f.config.Database.Database,
        User:            f.config.Database.User,
        Password:        f.config.Database.Password,
        SSLMode:         f.config.Database.SSLMode,
        MaxConnections:  f.config.Database.MaxConnections,
        MaxIdleConns:    f.config.Database.MaxIdleConnections,
        ConnMaxLifetime: f.config.Database.ConnectionTimeout,
    }

    db, err := database.NewGormConnection(dbConfig)
    if err != nil {
        return nil, nil, fmt.Errorf("gorm connection: %w", err)
    }

    closeFunc := func() error {
        return database.CloseGorm(db)
    }

    return db, closeFunc, nil
}
```

```go
// di/factories/repository_factory.go
package factories

import (
    "gorm.io/gorm"
    "github.com/oficinapro/auth-service/internal/domain/repository"
    "github.com/oficinapro/auth-service/internal/infrastructure/database"
)

// RepositoryFactory creates repository dependencies
type RepositoryFactory struct {
    db *gorm.DB
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
    return &RepositoryFactory{db: db}
}

// CreateClienteRepository creates cliente repository
func (f *RepositoryFactory) CreateClienteRepository() repository.ClienteRepository {
    return database.NewGormClienteRepository(f.db)
}
```

```go
// di/factories/service_factory.go
package factories

import (
    "github.com/oficinapro/auth-service/di/config"
    "github.com/oficinapro/auth-service/internal/domain/service"
    "github.com/oficinapro/auth-service/internal/infrastructure/jwt"
    "github.com/oficinapro/auth-service/internal/infrastructure/validator"
)

// ServiceFactory creates service dependencies
type ServiceFactory struct {
    config *config.Config
}

// NewServiceFactory creates a new service factory
func NewServiceFactory(cfg *config.Config) *ServiceFactory {
    return &ServiceFactory{config: cfg}
}

// CreateJWTService creates JWT service
func (f *ServiceFactory) CreateJWTService() service.JWTService {
    return jwt.NewJWTService(f.config.JWT.Secret, f.config.JWT.Issuer)
}

// CreateValidatorService creates validator service
func (f *ServiceFactory) CreateValidatorService() service.ValidatorService {
    return validator.NewCPFValidator()
}
```

```go
// di/factories/usecase_factory.go
package factories

import (
    "github.com/oficinapro/auth-service/di/config"
    "github.com/oficinapro/auth-service/internal/domain/repository"
    "github.com/oficinapro/auth-service/internal/domain/service"
    "github.com/oficinapro/auth-service/internal/usecase"
)

// UseCaseFactory creates use case dependencies
type UseCaseFactory struct {
    config           *config.Config
    clienteRepo      repository.ClienteRepository
    jwtService       service.JWTService
    validatorService service.ValidatorService
}

// NewUseCaseFactory creates a new use case factory
func NewUseCaseFactory(
    cfg *config.Config,
    clienteRepo repository.ClienteRepository,
    jwtService service.JWTService,
    validatorService service.ValidatorService,
) *UseCaseFactory {
    return &UseCaseFactory{
        config:           cfg,
        clienteRepo:      clienteRepo,
        jwtService:       jwtService,
        validatorService: validatorService,
    }
}

// CreateAuthenticateUseCase creates authenticate use case
func (f *UseCaseFactory) CreateAuthenticateUseCase() usecase.AuthenticateUseCase {
    return usecase.NewAuthenticateUseCase(
        f.clienteRepo,
        f.jwtService,
        f.validatorService,
        f.config.JWT.Expiration,
    )
}
```

**Benefícios**:
- ✅ **SRP**: Cada factory cuida de UMA camada
- ✅ **OCP**: Adicionar novo componente = adicionar método na factory correspondente
- ✅ **Testável**: Factories podem ser testadas isoladamente
- ✅ **Sem error desnecessário**: métodos sem error se não podem falhar

---

### 4. Facade (API Simplificada)

```go
// di/container_setup.go
package di

import (
    "context"
    "fmt"

    "github.com/oficinapro/auth-service/di/config"
    "github.com/oficinapro/auth-service/di/factories"
)

// NewContainer creates a new fully configured container
func NewContainer(ctx context.Context) (*Container, error) {
    // 1. Load config
    cfg, err := config.Load()
    if err != nil {
        return nil, fmt.Errorf("config: %w", err)
    }

    // 2. Create factories
    infraFactory := factories.NewInfrastructureFactory(cfg)

    // 3. Create infrastructure
    db, closeDB, err := infraFactory.CreateDatabase()
    if err != nil {
        return nil, fmt.Errorf("database: %w", err)
    }

    // 4. Create repositories
    repoFactory := factories.NewRepositoryFactory(db)
    clienteRepo := repoFactory.CreateClienteRepository()

    // 5. Create services
    serviceFactory := factories.NewServiceFactory(cfg)
    jwtService := serviceFactory.CreateJWTService()
    validatorService := serviceFactory.CreateValidatorService()

    // 6. Create use cases
    useCaseFactory := factories.NewUseCaseFactory(cfg, clienteRepo, jwtService, validatorService)
    authenticateUC := useCaseFactory.CreateAuthenticateUseCase()

    // 7. Build container
    return NewContainerBuilder().
        WithConfig(cfg).
        WithDatabase(db, closeDB).
        WithClienteRepository(clienteRepo).
        WithJWTService(jwtService).
        WithValidatorService(validatorService).
        WithAuthenticateUseCase(authenticateUC).
        Build()
}
```

**Benefícios**:
- ✅ **API simples**: `NewContainer(ctx)` funciona igual
- ✅ **Backward compatible**: não quebra código existente
- ✅ **Flexível**: pode usar builder diretamente para testes

---

## 🧪 Testabilidade

### Antes (Difícil de testar)

```go
// ❌ Como testar? Container cria implementações concretas
func TestSomething(t *testing.T) {
    container, _ := di.NewContainer(context.Background())
    // Como mockar o jwtService? Não dá!
}
```

### Depois (Fácil de testar)

```go
// ✅ FÁCIL: Injeta mocks via builder
func TestAuthenticateUseCase(t *testing.T) {
    mockRepo := mocks.NewMockClienteRepository()
    mockJWT := mocks.NewMockJWTService()
    mockValidator := mocks.NewMockValidatorService()

    cfg := &config.Config{...}

    container, _ := di.NewContainerBuilder().
        WithClienteRepository(mockRepo).
        WithJWTService(mockJWT).
        WithValidatorService(mockValidator).
        WithAuthenticateUseCase(
            usecase.NewAuthenticateUseCase(
                mockRepo, mockJWT, mockValidator, 24*time.Hour,
            ),
        ).
        Build()

    // Testa com mocks!
}
```

---

## 📊 Comparação: Antes vs Depois

| Aspecto | Antes ❌ | Depois ✅ |
|---------|---------|-----------|
| **SRP** | Container tem 6+ responsabilidades | Cada factory tem 1 responsabilidade |
| **OCP** | Adicionar = modificar Container | Adicionar = novo método em factory |
| **DIP** | Depende de implementações | Depende de interfaces |
| **Testabilidade** | Difícil mockar | Fácil injetar mocks |
| **Error Handling** | Errors desnecessários | Errors apenas onde faz sentido |
| **God Object** | Container sabe tudo | Responsabilidades distribuídas |
| **Lazy Loading** | Não suporta | Suporta (via builder manual) |
| **Acoplamento** | Alto | Baixo |

---

## 🚀 Migração (Passo a Passo)

### Fase 1: Adicionar Estrutura Nova (SEM quebrar existente)

1. Criar `di/factories/` com factories
2. Criar `di/builder.go` com builder
3. Criar `di/container_setup.go` com facade
4. **NÃO mexer** nos arquivos existentes ainda

### Fase 2: Ajustar Container

1. Refatorar `container.go` para ter apenas getters
2. Mover exposição de `*gorm.DB` para interno

### Fase 3: Remover Código Antigo

1. Remover `infrastructure.go`
2. Remover `repositories.go`
3. Remover `services.go`
4. Remover `usecases.go`
5. Manter apenas: `container.go`, `builder.go`, `container_setup.go`, `factories/`

---

## 🎯 Resultado Final

### Estrutura de Pastas:

```
di/
├── config/
│   └── config.go               # Config (mantém)
├── factories/
│   ├── infrastructure_factory.go # NOVO
│   ├── repository_factory.go     # NOVO
│   ├── service_factory.go        # NOVO
│   └── usecase_factory.go        # NOVO
├── builder.go                  # NOVO (ContainerBuilder)
├── container.go                # REFATORADO (simplificado)
└── container_setup.go          # NOVO (facade/setup)
```

### Uso (Mantém compatibilidade):

```go
// main.go
func main() {
    ctx := context.Background()

    // ✅ API simples (backward compatible)
    container, err := di.NewContainer(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer container.Close(ctx)

    // ✅ Usa via getters
    authUC := container.GetAuthenticateUseCase()
    // ...
}
```

---

## ✅ Checklist de Validação SOLID

Após refatoração:

- [x] **SRP**: Cada factory tem UMA responsabilidade
- [x] **OCP**: Extensível sem modificar código existente
- [x] **LSP**: Interfaces podem ser substituídas
- [x] **ISP**: Interfaces segregadas (cada uma com seu propósito)
- [x] **DIP**: Container depende de abstrações, não implementações

---

## 📝 Recomendação Final

**Para um projeto acadêmico (Tech Challenge FIAP)**, essa refatoração demonstra:

1. ✅ Conhecimento profundo de **SOLID**
2. ✅ Aplicação prática de **Design Patterns** (Builder, Factory)
3. ✅ **Clean Architecture** bem implementada
4. ✅ **Go idiomático** e boas práticas
5. ✅ **Testabilidade** e qualidade de código

**Vale MUITO a pena fazer essa refatoração** antes de apresentar o projeto! 🚀

---

## 🔗 Referências

- [Clean Architecture (Uncle Bob)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Go Dependency Injection](https://github.com/google/wire)
- [Builder Pattern](https://refactoring.guru/design-patterns/builder)
- [Factory Pattern](https://refactoring.guru/design-patterns/factory-method)

