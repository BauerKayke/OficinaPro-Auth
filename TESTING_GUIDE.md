# 🧪 Guia de Testes - Auth Gateway Email+Senha

## Teste Manual da API

### 1. **Credenciais de Teste Disponíveis**

Após executar `scripts/init.sql`, você terá 3 usuários criados:

| Email | Senha | Role | Descrição |
|-------|-------|------|-----------|
| `admin@oficinapro.com` | `senha123` | `ADMIN` | Administrador do sistema |
| `user@oficinapro.com` | `senha123` | `USER` | Usuário padrão |
| `joao.silva@example.com` | `senha123` | `USER` | Usuário com Pessoa vinculada |

---

### 2. **Testes com curl**

#### ✅ **Autenticação bem-sucedida (Admin)**
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@oficinapro.com",
    "senha": "senha123"
  }' | jq
```

**Resposta esperada** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImVtYWlsIjoiYWRtaW5Ab2ZpY2luYXByby5jb20iLCJyb2xlIjoiQURNSU4iLCJleHAiOjE3MzQ1Njc4OTAsImlhdCI6MTczNDU2NDI5MCwiaXNzIjoiYXV0aC1vZmljaW5hcHJvIn0.xxx",
  "expiresIn": 3600,
  "userId": 1,
  "email": "admin@oficinapro.com",
  "nome": "admin@oficinapro.com",
  "role": "ADMIN"
}
```

#### ❌ **Email não existe**
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "naoexiste@example.com",
    "senha": "senha123"
  }' | jq
```

**Resposta esperada** (401 Unauthorized):
```json
{
  "error": "Credenciais inválidas",
  "timestamp": "2024-12-20T12:00:00Z"
}
```

#### ❌ **Senha incorreta**
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@oficinapro.com",
    "senha": "senhaerrada"
  }' | jq
```

**Resposta esperada** (401 Unauthorized):
```json
{
  "error": "Credenciais inválidas",
  "timestamp": "2024-12-20T12:00:00Z"
}
```

#### ❌ **Email inválido (formato)**
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "emailinvalido",
    "senha": "senha123"
  }' | jq
```

**Resposta esperada** (400 Bad Request):
```json
{
  "error": "invalid email format",
  "timestamp": "2024-12-20T12:00:00Z"
}
```

#### ❌ **Campo obrigatório faltando**
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@oficinapro.com"
  }' | jq
```

**Resposta esperada** (400 Bad Request):
```json
{
  "error": "senha is required",
  "timestamp": "2024-12-20T12:00:00Z"
}
```

---

### 3. **Testes com Postman/Insomnia**

#### **Request**
```
POST http://localhost:8080/auth
Content-Type: application/json

{
  "email": "admin@oficinapro.com",
  "senha": "senha123"
}
```

#### **Validar JWT Gerado**

Copie o token JWT da resposta e cole em [jwt.io](https://jwt.io) para decodificar e validar:

**Payload esperado**:
```json
{
  "userId": 1,
  "email": "admin@oficinapro.com",
  "nome": "admin@oficinapro.com",
  "role": "ADMIN",
  "exp": 1734567890,
  "iat": 1734564290,
  "iss": "auth-oficinapro"
}
```

---

### 4. **Testes de Usuário Inativo**

Para testar um usuário desativado, execute no banco:

```sql
-- Desativar usuário
UPDATE usuario
SET is_ativo = false, data_desativacao = NOW()
WHERE email = 'user@oficinapro.com';

-- Tentar autenticar (deve retornar 403 Forbidden)
```

**curl**:
```bash
curl -X POST http://localhost:8080/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@oficinapro.com",
    "senha": "senha123"
  }' | jq
```

**Resposta esperada** (403 Forbidden):
```json
{
  "error": "Usuário inativo",
  "timestamp": "2024-12-20T12:00:00Z"
}
```

**Reativar usuário**:
```sql
UPDATE usuario
SET is_ativo = true, data_desativacao = NULL
WHERE email = 'user@oficinapro.com';
```

---

### 5. **Health Check**

```bash
curl http://localhost:8080/health | jq
```

**Resposta esperada** (200 OK):
```json
{
  "status": "ok",
  "service": "auth-oficinapro"
}
```

---

## Teste com Script Bash

Crie `test-auth.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== Teste 1: Autenticação Admin (Sucesso) ==="
curl -X POST $BASE_URL/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@oficinapro.com","senha":"senha123"}' \
  -w "\nHTTP Status: %{http_code}\n\n" | jq

echo "=== Teste 2: Senha Incorreta ==="
curl -X POST $BASE_URL/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@oficinapro.com","senha":"senhaerrada"}' \
  -w "\nHTTP Status: %{http_code}\n\n" | jq

echo "=== Teste 3: Email Inexistente ==="
curl -X POST $BASE_URL/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"naoexiste@example.com","senha":"senha123"}' \
  -w "\nHTTP Status: %{http_code}\n\n" | jq

echo "=== Teste 4: Email Inválido ==="
curl -X POST $BASE_URL/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"emailinvalido","senha":"senha123"}' \
  -w "\nHTTP Status: %{http_code}\n\n" | jq

echo "=== Teste 5: Health Check ==="
curl $BASE_URL/health -w "\nHTTP Status: %{http_code}\n\n" | jq
```

**Executar**:
```bash
chmod +x test-auth.sh
./test-auth.sh
```

---

## Teste de Carga (Opcional)

Use `hey` ou `ab` para testes de performance:

```bash
# Instalar hey (Go)
go install github.com/rakyll/hey@latest

# 1000 requests, 50 concorrentes
hey -n 1000 -c 50 -m POST \
  -H "Content-Type: application/json" \
  -d '{"email":"user@oficinapro.com","senha":"senha123"}' \
  http://localhost:8080/auth
```

---

## Logs Esperados

### Sucesso
```
INFO Authentication successful email=us***@oficinapro.com userId=2 duration=45ms
```

### Falha (senha errada)
```
ERROR Authentication failed email=us***@oficinapro.com error="credenciais inválidas"
```

### Falha (email não existe)
```
ERROR Authentication failed email=no***@example.com error="usuário não encontrado"
```

**Nota**: Emails são mascarados nos logs para segurança.

---

## Troubleshooting

### Erro: "connection refused"
- Verificar se servidor está rodando: `ps aux | grep lambda`
- Verificar porta: `lsof -i :8080`

### Erro: "database connection failed"
- Verificar se Postgres está rodando: `docker ps | grep postgres`
- Testar conexão: `psql -h localhost -U postgres -d oficinapro -c "\dt"`

### Token JWT inválido
- Verificar `JWT_SECRET` configurado corretamente
- Verificar `JWT_ISSUER` consistente entre geração e validação

### Senha não valida (mesmo estando correta)
- Verificar hash bcrypt no banco: `SELECT senha FROM usuario WHERE email = 'user@oficinapro.com';`
- Re-gerar hash:
  ```go
  hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
  fmt.Println(string(hash))
  ```

---

## Próximos Testes

1. **Testes unitários**: Atualizar todos os `*_test.go`
2. **Testes de integração**: E2E com banco real
3. **Testes de segurança**: Brute force, SQL injection, XSS
4. **Testes de performance**: Load testing com > 10k req/s

