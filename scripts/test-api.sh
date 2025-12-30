#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo "========================================="
echo "🧪 Testando Auth Service"
echo "========================================="
echo ""

# 1. Health Check
echo "1️⃣  Testando health check..."
HEALTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$BASE_URL/health")
HTTP_STATUS=$(echo "$HEALTH_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$HEALTH_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Health check OK"
    echo "   Response: $BODY"
else
    echo "❌ Health check FALHOU (HTTP $HTTP_STATUS)"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 2. Autenticação com credenciais corretas (Admin)
echo "2️⃣  Testando autenticação com admin@oficinapro.com..."
AUTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X POST "$BASE_URL/auth" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@oficinapro.com",
        "senha": "Admin@123"
    }')

HTTP_STATUS=$(echo "$AUTH_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$AUTH_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Autenticação admin OK"
    echo "   Response: $BODY"

    # Extrai o token
    ADMIN_TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "   Token extraído: ${ADMIN_TOKEN:0:50}..."
else
    echo "❌ Autenticação admin FALHOU (HTTP $HTTP_STATUS)"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 3. Autenticação com credenciais corretas (Client)
echo "3️⃣  Testando autenticação com client@oficinapro.com..."
AUTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X POST "$BASE_URL/auth" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "client@oficinapro.com",
        "senha": "Client@123"
    }')

HTTP_STATUS=$(echo "$AUTH_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$AUTH_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Autenticação client OK"
    echo "   Response: $BODY"

    # Extrai o token
    CLIENT_TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    echo "   Token extraído: ${CLIENT_TOKEN:0:50}..."
else
    echo "❌ Autenticação client FALHOU (HTTP $HTTP_STATUS)"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 4. Autenticação com senha incorreta
echo "4️⃣  Testando autenticação com senha incorreta..."
AUTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X POST "$BASE_URL/auth" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@oficinapro.com",
        "senha": "SenhaErrada123"
    }')

HTTP_STATUS=$(echo "$AUTH_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$AUTH_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 401 ]; then
    echo "✅ Rejeição de senha incorreta OK"
    echo "   Response: $BODY"
else
    echo "❌ Deveria retornar 401, mas retornou HTTP $HTTP_STATUS"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 5. Autenticação com usuário inexistente
echo "5️⃣  Testando autenticação com usuário inexistente..."
AUTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X POST "$BASE_URL/auth" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "usuario_inexistente@oficinapro.com",
        "senha": "Senha123"
    }')

HTTP_STATUS=$(echo "$AUTH_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$AUTH_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 401 ]; then
    echo "✅ Rejeição de usuário inexistente OK"
    echo "   Response: $BODY"
else
    echo "❌ Deveria retornar 401, mas retornou HTTP $HTTP_STATUS"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 6. Validação de token (Admin)
echo "6️⃣  Testando validação de token do admin..."
VALIDATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X POST "$BASE_URL/auth/validate" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_TOKEN")

HTTP_STATUS=$(echo "$VALIDATE_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$VALIDATE_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Validação de token admin OK"
    echo "   Response: $BODY"
else
    echo "❌ Validação de token admin FALHOU (HTTP $HTTP_STATUS)"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

# 7. CORS Preflight
echo "7️⃣  Testando CORS preflight..."
CORS_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" \
    -X OPTIONS "$BASE_URL/auth" \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST")

HTTP_STATUS=$(echo "$CORS_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ CORS preflight OK"
else
    echo "❌ CORS preflight FALHOU (HTTP $HTTP_STATUS)"
    exit 1
fi
echo ""

# 8. Endpoint não encontrado
echo "8️⃣  Testando endpoint não encontrado..."
NOT_FOUND_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "$BASE_URL/nao-existe")

HTTP_STATUS=$(echo "$NOT_FOUND_RESPONSE" | grep "HTTP_STATUS" | cut -d':' -f2)
BODY=$(echo "$NOT_FOUND_RESPONSE" | sed -e 's/HTTP_STATUS.*$//')

if [ "$HTTP_STATUS" -eq 404 ]; then
    echo "✅ Endpoint não encontrado OK"
    echo "   Response: $BODY"
else
    echo "❌ Deveria retornar 404, mas retornou HTTP $HTTP_STATUS"
    echo "   Response: $BODY"
    exit 1
fi
echo ""

echo "========================================="
echo "✅ Todos os testes passaram!"
echo "========================================="

