# 📡 API Documentation

## Endpoint de Autenticação

### POST /auth

Autentica um cliente via CPF e retorna um token JWT.

#### Request

**Method:** `POST`
**Content-Type:** `application/json`

```json
{
  "cpf": "123.456.789-09"
}
```

**Validações:**
- CPF é obrigatório
- CPF deve ser válido (dígitos verificadores)
- Cliente deve existir no banco
- Cliente deve estar ativo

#### Response - Sucesso (200)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 86400,
  "clienteId": 1,
  "nome": "João Silva"
}
```

**Campos:**
- `token`: Token JWT para autenticação
- `expiresIn`: Tempo de expiração em segundos (24h = 86400)
- `clienteId`: ID do cliente autenticado
- `nome`: Nome do cliente

#### Response - Erro (4xx/5xx)

```json
{
  "error": "CPF inválido",
  "timestamp": "2025-10-19T10:30:00Z"
}
```

**Possíveis Erros:**

| Status | Error | Descrição |
|--------|-------|-----------|
| 400 | Invalid request body | JSON malformado |
| 400 | CPF is required | CPF não informado |
| 401 | CPF inválido | CPF com formato/dígitos inválidos |
| 401 | Cliente não encontrado | Cliente não existe |
| 401 | Cliente inativo | Cliente desativado |
| 405 | Method not allowed | Método diferente de POST |
| 500 | Internal server error | Erro interno |

#### JWT Token

O token JWT contém os seguintes claims:

```json
{
  "cliente_id": 1,
  "nome": "João Silva",
  "email": "joao@example.com",
  "documento": "12345678909",
  "iat": 1697712000,
  "exp": 1697798400,
  "iss": "oficinapro-auth"
}
```

**Claims:**
- `cliente_id`: ID do cliente
- `nome`: Nome do cliente
- `email`: Email do cliente
- `documento`: CPF do cliente
- `iat`: Issued at (timestamp)
- `exp`: Expiration (timestamp)
- `iss`: Issuer (oficinapro-auth)

#### Exemplos de Uso

**cURL:**
```bash
curl -X POST https://api.oficinapro.com/auth \
  -H "Content-Type: application/json" \
  -d '{"cpf":"123.456.789-09"}'
```

**JavaScript:**
```javascript
const response = await fetch('https://api.oficinapro.com/auth', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ cpf: '123.456.789-09' })
});

const data = await response.json();
console.log(data.token);
```

**Go:**
```go
type AuthRequest struct {
    CPF string `json:"cpf"`
}

request := AuthRequest{CPF: "123.456.789-09"}
body, _ := json.Marshal(request)

resp, err := http.Post(
    "https://api.oficinapro.com/auth",
    "application/json",
    bytes.NewBuffer(body),
)
```

#### CORS Headers

A API retorna os seguintes headers CORS:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Headers: Content-Type,Authorization
Access-Control-Allow-Methods: POST,OPTIONS
```

#### Rate Limiting

- **Limite:** 100 requests por minuto por IP
- **Header:** `X-RateLimit-Remaining`
- **Resposta ao atingir limite:** 429 Too Many Requests

#### Security

- ✅ SQL Injection: Protegido (GORM)
- ✅ HTTPS Only: Sim (via API Gateway)
- ✅ Input Validation: Sim (CPF validation)
- ✅ JWT Signature: HMAC-SHA256
- ✅ Logs: CPF mascarado (ex: 123***09)

#### Performance

- **Cold Start:** ~100ms
- **Warm Execution:** ~10ms
- **Timeout:** 30s
- **Memory:** 128MB

