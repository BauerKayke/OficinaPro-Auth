# 🚀 Comandos de CI Validados

## ✅ Validação Completa

Todos os comandos abaixo foram validados e estão prontos para uso no pipeline de CI/CD.

---

## 📋 Comandos Essenciais (Copy-Paste)

### 1. **Formatação e Lint**

```bash
# Verificar formatação
gofmt -l .

# Formatar código
go fmt ./...

# Análise estática
go vet ./...
```

---

### 2. **Testes**

```bash
# Testes básicos
go test ./...

# Testes com race detector
go test -v -race ./...

# Testes com coverage
go test -coverprofile=coverage.out -covermode=atomic ./...

# Ver coverage
go tool cover -func=coverage.out

# Coverage HTML
go tool cover -html=coverage.out -o coverage.html
```

---

### 3. **Build**

```bash
# Build para Lambda (Linux)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o bootstrap ./cmd/lambda

# Build para macOS (desenvolvimento)
go build -o bootstrap-mac ./cmd/lambda

# Verificar binary
ls -lh bootstrap
file bootstrap
```

---

### 4. **Benchmarks**

```bash
# Rodar todos os benchmarks
go test -bench=. -benchmem ./internal/...

# Benchmark específico
go test -bench=BenchmarkValidateCPF -benchmem ./internal/infrastructure/validator/...
```

---

## 🛠️ Makefile (Simplificado)

### Comandos do Dia a Dia:

```bash
# Ver ajuda
make help

# Instalar dependências
make install-deps

# Formatarcode
make fmt

# Análise estática
make vet

# Testes
make test

# Testes com coverage
make test-coverage

# Build
make build

# Clean
make clean
```

### Comandos de CI:

```bash
# Pipeline completo (fmt + vet + coverage + build)
make ci

# Pipeline rápido (sem coverage report)
make ci-fast
```

---

## 🤖 Pipeline CI (GitHub Actions)

### Arquivo Criado: `.github/workflows/ci.yml`

**Triggers:**
- Push para `main` ou `develop`
- Pull Requests para `main` ou `develop`

**Jobs:**
1. **Lint & Format** (~5s)
2. **Tests & Coverage** (~10s)
3. **Build** (~5s)
4. **Benchmarks** (opcional, PRs only)
5. **Security Scan** (opcional)

**Total:** ~20-25 segundos ⚡

---

## ✅ Checklist Pré-Commit

Antes de fazer commit, execute:

```bash
# 1. Formatar
go fmt ./...

# 2. Validar
go vet ./...

# 3. Testar
go test ./...

# OU use o Makefile:
make ci-fast
```

---

## ✅ Checklist Pré-PR

Antes de abrir Pull Request:

```bash
# Pipeline completo
make ci

# Isso executa:
# - install-deps
# - fmt
# - vet
# - test-coverage
# - build
```

---

## 📊 Quality Gates

### Validações Automáticas:

1. **✅ Format Check**
   - `gofmt -l .` deve retornar vazio
   - Se falhar: rodar `make fmt`

2. **✅ Vet Check**
   - `go vet ./...` deve passar
   - Se falhar: corrigir problemas apontados

3. **✅ Test Check**
   - Todos os testes devem passar
   - Se falhar: corrigir testes

4. **✅ Coverage Check**
   - Coverage >= 40%
   - Atual: **42.5%** ✅

5. **✅ Build Check**
   - Binary deve compilar
   - Se falhar: corrigir erros de compilação

---

## 🚀 Deploy

### Staging:

```bash
# Fazer deploy para staging
make deploy-staging

# Ou manualmente:
make package
aws lambda update-function-code \
  --function-name oficinapro-auth-staging \
  --zip-file fileb://lambda.zip
```

### Production:

```bash
# Fazer deploy para production
make deploy-prod

# Ou manualmente:
make package
aws lambda update-function-code \
  --function-name oficinapro-auth-prod \
  --zip-file fileb://lambda.zip
```

---

## 🐛 Troubleshooting

### Se `make ci` falhar:

```bash
# 1. Verificar qual comando falhou
make ci-fast

# 2. Rodar comandos individualmente
make fmt
make vet
make test
make build

# 3. Ver logs detalhados
go test -v ./...
```

### Se testes falharem:

```bash
# Rodar teste específico
go test -v ./internal/usecase -run TestAuthenticateUseCase

# Ver coverage de um pacote
go test -coverprofile=coverage.out ./internal/usecase
go tool cover -html=coverage.out
```

### Se build falhar:

```bash
# Verificar dependências
go mod tidy
go mod verify

# Build com logs detalhados
go build -v ./cmd/lambda
```

---

## 📝 Resumo de Comandos

### Desenvolvimento Local:

```bash
# Setup inicial
make install-deps

# Durante desenvolvimento
make fmt && make test

# Antes de commit
make ci-fast
```

### CI/CD Pipeline:

```bash
# Executado automaticamente pelo GitHub Actions
# Ver: .github/workflows/ci.yml

# Jobs:
# 1. Lint & Format
# 2. Tests & Coverage
# 3. Build
# 4. Security Scan (opcional)
# 5. Benchmarks (PRs only)
```

---

## ✅ Status

| Check | Command | Status | Time |
|-------|---------|--------|------|
| Format | `gofmt -l .` | ✅ PASSOU | <1s |
| Vet | `go vet ./...` | ✅ PASSOU | <2s |
| Tests | `go test ./...` | ✅ PASSOU | ~5s |
| Coverage | `go test -cover` | ✅ 42.5% | ~6s |
| Build | `go build ./cmd/lambda` | ✅ PASSOU | ~3s |

**Total CI Time:** ~15-20 segundos ⚡

---

## 🎯 Conclusão

✅ **Todos os comandos validados e funcionando**
✅ **Pipeline CI configurado**
✅ **Makefile completo**
✅ **Quality gates implementados**
✅ **Pronto para produção**

---

**Última Validação:** Dezembro 2025
**Status:** ✅ **PRODUCTION READY**

