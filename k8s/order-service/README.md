# Kubernetes Manifests - Order Service

Manifests Kubernetes para deploy do Order Service no EKS.

## 📁 Estrutura

```
k8s/order-service/
├── namespace.yaml          # Namespace
├── configmap.yaml          # Configurações não-sensíveis
├── secret.yaml             # Secrets (JWT, DB password)
├── deployment.yaml         # Deployment
├── service.yaml            # Service (ClusterIP)
├── ingress.yaml            # Ingress (ALB)
├── hpa.yaml                # Horizontal Pod Autoscaler
└── kustomization.yaml      # Kustomize
```

## 🚀 Deploy

### Opção 1: kubectl apply
```bash
kubectl apply -f k8s/order-service/
```

### Opção 2: Kustomize
```bash
kubectl apply -k k8s/order-service/
```

### Opção 3: Helm (futuro)
```bash
helm install order-service ./helm/order-service
```

## 🔧 Configuração

### Secrets

**⚠️ IMPORTANTE**: Criar secrets ANTES de fazer deploy:

```bash
# JWT Secret (MESMO DO AUTH GATEWAY)
kubectl create secret generic order-service-jwt \
  --from-literal=JWT_SECRET="your-secret-key-min-32-chars-long" \
  -n oficinapro

# DB Password
kubectl create secret generic order-service-db \
  --from-literal=DB_PASSWORD="your-db-password" \
  -n oficinapro
```

Ou usando AWS Secrets Manager (External Secrets):
```yaml
apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: order-service-secrets
spec:
  secretStoreRef:
    name: aws-secrets-manager
  target:
    name: order-service-secret
  data:
  - secretKey: JWT_SECRET
    remoteRef:
      key: /oficinapro/prod/jwt/secret
  - secretKey: DB_PASSWORD
    remoteRef:
      key: /oficinapro/prod/db/password
```

## 📊 Recursos

### Requests/Limits
- **Requests**: 128Mi RAM, 100m CPU
- **Limits**: 256Mi RAM, 200m CPU
- **Replicas**: 3 (para HA)

### Autoscaling
- **Min**: 3 replicas
- **Max**: 10 replicas
- **Target CPU**: 70%
- **Target Memory**: 80%

## 🔍 Monitoring

### Health Checks
- **Liveness**: `GET /health` (10s interval)
- **Readiness**: `GET /health` (5s interval)

### Logs
```bash
# Ver logs
kubectl logs -f deployment/order-service -n oficinapro

# Ver logs de um pod específico
kubectl logs -f pod/order-service-xxx -n oficinapro

# Seguir logs de todos os pods
kubectl logs -f -l app=order-service -n oficinapro
```

### Metrics
```bash
# CPU/Memory usage
kubectl top pods -n oficinapro

# Describe deployment
kubectl describe deployment order-service -n oficinapro

# Get events
kubectl get events -n oficinapro --sort-by='.lastTimestamp'
```

## 🌐 Ingress

### ALB Ingress Controller
```yaml
annotations:
  kubernetes.io/ingress.class: alb
  alb.ingress.kubernetes.io/scheme: internet-facing
  alb.ingress.kubernetes.io/target-type: ip
  alb.ingress.kubernetes.io/listen-ports: '[{"HTTPS":443}]'
  alb.ingress.kubernetes.io/certificate-arn: arn:aws:acm:...
```

### DNS
Após deploy, pegar o ALB endpoint:
```bash
kubectl get ingress order-service -n oficinapro
```

Criar CNAME em Route53:
```
orders.oficinapro.com → <alb-url>.elb.amazonaws.com
```

## 🧪 Testes

### 1. Port-forward (local)
```bash
kubectl port-forward svc/order-service 8080:80 -n oficinapro

# Testar
curl http://localhost:8080/health
```

### 2. Testar com JWT
```bash
# Obter JWT do Auth Gateway
TOKEN=$(curl -s -X POST https://api.oficinapro.com/auth \
  -H "Content-Type: application/json" \
  -d '{"email":"user@test.com","senha":"senha123"}' | jq -r '.token')

# Testar endpoint protegido
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/orders
```

## 🔄 Rolling Update

```bash
# Atualizar imagem
kubectl set image deployment/order-service \
  order-service=123456789.dkr.ecr.us-east-1.amazonaws.com/order-service:v1.1.0 \
  -n oficinapro

# Ver status do rollout
kubectl rollout status deployment/order-service -n oficinapro

# Rollback se necessário
kubectl rollout undo deployment/order-service -n oficinapro
```

## 🛡️ RBAC

ServiceAccount com permissões mínimas:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: order-service
  namespace: oficinapro
```

## 📝 Troubleshooting

### Pod não inicia
```bash
kubectl describe pod order-service-xxx -n oficinapro
kubectl logs order-service-xxx -n oficinapro
```

### Erro de conexão com RDS
- Verificar Security Group do RDS
- Verificar se pods estão nas subnets corretas
- Verificar DNS do RDS

### JWT inválido
- Verificar se JWT_SECRET é o mesmo do Auth Gateway
- Ver logs: `kubectl logs -f deployment/order-service -n oficinapro`

### Load Balancer não funciona
```bash
kubectl describe ingress order-service -n oficinapro
kubectl get events -n oficinapro
```

## 🚀 Production Checklist

- [ ] Secrets criados via Secrets Manager
- [ ] Resource limits configurados
- [ ] HPA configurado
- [ ] Liveness/Readiness probes testados
- [ ] Ingress com TLS (ACM certificate)
- [ ] DNS configurado (Route53)
- [ ] Logs centralizados (CloudWatch/ELK)
- [ ] Métricas (Prometheus/CloudWatch)
- [ ] Alertas configurados
- [ ] Backup strategy definida

---

**Ver também**:
- `examples/order-service/` - Código do serviço
- `INFRASTRUCTURE_STRATEGY.md` - Estratégia completa

