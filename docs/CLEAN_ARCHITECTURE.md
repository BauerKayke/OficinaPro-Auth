# 🏛️ Clean Architecture em Go

## Visão Geral

Este projeto segue os princípios da **Clean Architecture** de Robert C. Martin, adaptados para Go de forma idiomática.

## 📦 Estrutura de Camadas

```
┌─────────────────────────────────────────────────────┐
│              cmd/lambda/main.go                     │
│            (Lambda Handler - Entry Point)           │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│           internal/usecase/                         │
│         (Application Business Rules)                │
│  • AuthenticateUseCase                             │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│           internal/domain/                          │
│         (Enterprise Business Rules)                 │
│  • entity/     - Domain entities                   │
│  • repository/ - Port interfaces                   │
│  • service/    - Service interfaces                │
└─────────────────────────────────────────────────────┘
                        ▲
                        │ implements
┌─────────────────────────────────────────────────────┐
│        internal/infrastructure/                     │
│              (Adapters)                             │
│  • database/   - PostgreSQL adapter                │
│  • jwt/        - JWT service implementation        │
│  • validator/  - CPF validator implementation      │
└─────────────────────────────────────────────────────┘
```

## 🎯 Princípios SOLID

### 1. Single Responsibility Principle (SRP)

Cada componente tem uma única responsabilidade:

- **`AuthenticateUseCase`**: Apenas autenticação
- **`PostgresClienteRepository`**: Apenas acesso a dados de clientes
- **`JWTService`**: Apenas geração/validação de JWT
- **`CPFValidator`**: Apenas validação de CPF

### 2. Open/Closed Principle (OCP)

Aberto para extensão, fechado para modificação:

```go
// Fácil adicionar novos validadores sem modificar código existente
type ValidatorService interface {
    ValidateCPF(cpf string) bool
    NormalizeCPF(cpf string) string
}

// Nova implementação sem modificar interface
type EnhancedCPFValidator struct {
    // ...
}
```

### 3. Liskov Substitution Principle (LSP)

Qualquer implementação da interface pode substituir outra:

```go
// Qualquer implementação de ClienteRepository pode ser usada
var repo repository.ClienteRepository

repo = database.NewPostgresClienteRepository(db)  // PostgreSQL
// repo = database.NewMongoClienteRepository(mongo) // Poderia ser MongoDB
```

### 4. Interface Segregation Principle (ISP)

Interfaces específicas e coesas:

```go
// Interface focada apenas em operações de cliente
type ClienteRepository interface {
    FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error)
    HealthCheck(ctx context.Context) error
}

// Não tem métodos desnecessários como Update, Delete, etc.
```

### 5. Dependency Inversion Principle (DIP)

Dependências apontam para abstrações:

```go
// Use case depende de interfaces, não implementações concretas
type AuthenticateUseCase struct {
    clienteRepo      repository.ClienteRepository      // interface
    jwtService       service.JWTService                // interface
    validatorService service.ValidatorService          // interface
}
```

## 📋 Fluxo de Dados

### Request Flow

```
1. Lambda Handler (main.go)
   ↓
2. Parse & Validate Input
   ↓
3. AuthenticateUseCase.Execute()
   ↓
4. ValidatorService.ValidateCPF()
   ↓
5. ClienteRepository.FindByCPF()
   ↓
6. Cliente.CanAuthenticate()
   ↓
7. JWTService.GenerateToken()
   ↓
8. Return AuthenticateOutput
```

### Dependency Injection

Todas as dependências são injetadas via construtor:

```go
// No init() do Lambda
clienteRepo := database.NewPostgresClienteRepository(db)
jwtService := jwt.NewJWTService(secret, issuer)
validatorService := validator.NewCPFValidator()

// Injetar no use case
authenticateUC := usecase.NewAuthenticateUseCase(
    clienteRepo,
    jwtService,
    validatorService,
    24*time.Hour,
)
```

## 🧪 Testabilidade

### Unit Tests com Mocks

```go
// Mocks fáceis devido às interfaces
type MockClienteRepository struct {
    mock.Mock
}

func (m *MockClienteRepository) FindByCPF(ctx context.Context, cpf string) (*entity.Cliente, error) {
    args := m.Called(ctx, cpf)
    return args.Get(0).(*entity.Cliente), args.Error(1)
}

// Teste isolado sem dependências reais
func TestAuthenticateUseCase_Execute(t *testing.T) {
    mockRepo := new(MockClienteRepository)
    mockJWT := new(MockJWTService)
    mockValidator := new(MockValidatorService)

    useCase := NewAuthenticateUseCase(mockRepo, mockJWT, mockValidator, 24*time.Hour)

    // Testar lógica de negócio pura
    // ...
}
```

## 🔧 Boas Práticas Go

### 1. Errors são Valores

```go
// Domain errors como valores
var (
    ErrClienteNotFound = errors.New("cliente não encontrado")
    ErrClienteInativo  = errors.New("cliente inativo")
)

// Wrap errors para contexto
return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
```

### 2. Interfaces Pequenas

```go
// Interface focada
type ValidatorService interface {
    ValidateCPF(cpf string) bool
    NormalizeCPF(cpf string) string
}

// Não: interface gigante com 20 métodos
```

### 3. Accept Interfaces, Return Structs

```go
// Use case aceita interfaces
func NewAuthenticateUseCase(
    clienteRepo repository.ClienteRepository,  // interface
    jwtService service.JWTService,              // interface
) *AuthenticateUseCase {                       // struct concreto
    return &AuthenticateUseCase{...}
}
```

### 4. Context como Primeiro Parâmetro

```go
func (r *PostgresClienteRepository) FindByCPF(
    ctx context.Context,  // sempre primeiro
    cpf string,
) (*entity.Cliente, error) {
    // ...
}
```

## 🎨 Clean Code

### Nomenclatura Clara

```go
// BOM: nomes descritivos
func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput)

// RUIM: nomes genéricos
func (uc *AuthenticateUseCase) Do(ctx context.Context, in Input)
```

### Funções Pequenas

```go
// Cada função faz UMA coisa
func (uc *AuthenticateUseCase) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
    // 1. Validar
    if err := uc.validateInput(input); err != nil {
        return nil, err
    }

    // 2. Buscar
    cliente, err := uc.findCliente(ctx, input.CPF)
    if err != nil {
        return nil, err
    }

    // 3. Autenticar
    return uc.generateAuthOutput(cliente)
}
```

### Logging Estruturado

```go
log.Printf("[UseCase] Iniciando autenticação para CPF: %s", maskCPF(input.CPF))
log.Printf("[Repository] Cliente encontrado: ID=%d, Ativo=%v", cliente.ID, cliente.Ativo)
```

## 📈 Benefícios Alcançados

### ✅ Testabilidade

- Unit tests sem infraestrutura
- Mocks fáceis via interfaces
- Coverage >80%

### ✅ Manutenibilidade

- Código organizado e coeso
- Fácil localizar funcionalidades
- Mudanças isoladas

### ✅ Escalabilidade

- Fácil adicionar novos use cases
- Trocar implementações (ex: PostgreSQL → DynamoDB)
- Independência de frameworks

### ✅ Performance

- Go nativo (compilado)
- Cold start ~100ms
- Warm execution ~10ms

## 🔄 Extensibilidade

### Adicionar Novo Use Case

```go
// 1. Criar novo use case
type ValidateTokenUseCase struct {
    jwtService service.JWTService
}

// 2. Implementar lógica
func (uc *ValidateTokenUseCase) Execute(token string) error {
    _, err := uc.jwtService.ValidateToken(token)
    return err
}

// 3. Registrar no Lambda handler
validateTokenUC := NewValidateTokenUseCase(jwtService)
```

### Trocar Implementação

```go
// Fácil trocar PostgreSQL por DynamoDB
// clienteRepo := database.NewPostgresClienteRepository(db)
clienteRepo := database.NewDynamoDBClienteRepository(dynamoDB)

// Use case não muda!
authenticateUC := usecase.NewAuthenticateUseCase(clienteRepo, ...)
```

## 📚 Referências

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://golang.org/doc/effective_go)
- [SOLID Principles in Go](https://dave.cheney.net/2016/08/20/solid-go-design)

