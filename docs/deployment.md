# Deployment Guide

This guide covers deploying gocalc-api in various environments including Docker, Kubernetes, and cloud platforms.

## Table of Contents

- [Environment Variables](#environment-variables)
- [Docker Deployment](#docker-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Health Checks](#health-checks)
- [Configuration Best Practices](#configuration-best-practices)

## Environment Variables

The application is configured via environment variables with sensible defaults:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | HTTP server port | `8080` | `3000` |
| `READ_TIMEOUT` | HTTP read timeout | `10s` | `15s` |
| `WRITE_TIMEOUT` | HTTP write timeout | `10s` | `20s` |
| `IDLE_TIMEOUT` | HTTP idle timeout | `120s` | `180s` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `15s` | `30s` |
| `REQUEST_TIMEOUT` | Request processing timeout | `30s` | `45s` |
| `RATE_LIMIT_RPM` | Rate limit (requests per minute) | `100.0` | `200.0` |
| `RATE_LIMIT_BURST` | Rate limit burst size | `20` | `50` |

Duration values accept standard Go time formats: `10s`, `2m`, `1h`, etc.

## Docker Deployment

### Building the Image

```bash
# Build using Makefile
make docker-build

# Build with custom tag
make docker-build DOCKER_TAG=v1.0.0

# Or build directly with docker
docker build -t gocalc-api:latest .
```

The Dockerfile uses a multi-stage build:

- **Builder stage**: Compiles the Go application
- **Runtime stage**: Minimal Alpine-based image (~15MB)
- Non-root user for security
- Built-in health check

### Running the Container

```bash
# Run on default port (8080)
make docker-run

# Run in background
make docker-run-detached

# Run on custom port
make docker-run DOCKER_PORT=3000

# Run with custom environment variables
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e RATE_LIMIT_RPM=200 \
  -e REQUEST_TIMEOUT=45s \
  gocalc-api:latest
```

### Managing Containers

```bash
# View logs
make docker-logs

# Stop container
make docker-stop

# Open shell in container
make docker-shell

# Clean up image
make docker-clean
```

### Docker Compose Example

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    image: gocalc-api:latest
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - RATE_LIMIT_RPM=200
      - REQUEST_TIMEOUT=45s
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 3s
      retries: 3
      start_period: 5s
```

Run with:

```bash
docker-compose up -d
```

## Kubernetes Deployment

### Deployment Manifest

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gocalc-api
  labels:
    app: gocalc-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: gocalc-api
  template:
    metadata:
      labels:
        app: gocalc-api
    spec:
      containers:
      - name: gocalc-api
        image: gocalc-api:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: PORT
          value: "8080"
        - name: RATE_LIMIT_RPM
          value: "200"
        - name: REQUEST_TIMEOUT
          value: "45s"
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 5
          timeoutSeconds: 2
          failureThreshold: 2
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          readOnlyRootFilesystem: true
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
```

### Service Manifest

Create `k8s/service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: gocalc-api
  labels:
    app: gocalc-api
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: gocalc-api
```

### ConfigMap for Environment Variables

Create `k8s/configmap.yaml`:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: gocalc-api-config
data:
  PORT: "8080"
  RATE_LIMIT_RPM: "200"
  REQUEST_TIMEOUT: "45s"
  READ_TIMEOUT: "10s"
  WRITE_TIMEOUT: "10s"
  IDLE_TIMEOUT: "120s"
  SHUTDOWN_TIMEOUT: "15s"
  RATE_LIMIT_BURST: "20"
```

Update deployment to use ConfigMap:

```yaml
        envFrom:
        - configMapRef:
            name: gocalc-api-config
```

### Deploy to Kubernetes

```bash
# Apply all manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get deployments
kubectl get pods
kubectl get services

# View logs
kubectl logs -l app=gocalc-api -f

# Port forward for local testing
kubectl port-forward service/gocalc-api 8080:80
```

## Health Checks

The application exposes two health endpoints:

### Health Endpoint: `/health`

Provides detailed health information including:

- Service status
- Uptime
- Version
- Memory usage check
- Goroutine count check
- Timestamp

**Example Response:**

```json
{
  "data": {
    "status": "healthy",
    "timestamp": "2026-02-01T12:34:56Z",
    "uptime": "2h15m30s",
    "version": "1.0.0",
    "checks": {
      "memory": "ok",
      "goroutines": "ok"
    }
  },
  "request_id": "abc123",
  "timestamp": "2026-02-01T12:34:56Z"
}
```

**Use for:**

- Kubernetes liveness probes
- Docker health checks
- Monitoring systems

### Readiness Endpoint: `/ready`

Indicates if the service is ready to accept traffic.

**Example Response:**

```json
{
  "data": {
    "ready": true,
    "status": "ready"
  },
  "request_id": "abc123",
  "timestamp": "2026-02-01T12:34:56Z"
}
```

**Use for:**

- Kubernetes readiness probes
- Load balancer health checks
- Service mesh configurations

## Configuration Best Practices

### Development

Use defaults or minimal configuration:

```bash
# Run locally with defaults
make run

# Or with custom port
PORT=3000 make run
```

### Staging/Production

Use environment-specific configuration:

```bash
# Staging
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e RATE_LIMIT_RPM=150 \
  -e REQUEST_TIMEOUT=45s \
  gocalc-api:latest

# Production
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e RATE_LIMIT_RPM=300 \
  -e REQUEST_TIMEOUT=30s \
  -e READ_TIMEOUT=15s \
  -e WRITE_TIMEOUT=15s \
  gocalc-api:latest
```

### Security Considerations

1. **Non-root user**: The Docker image runs as a non-root user (UID 1000)
2. **Read-only filesystem**: Kubernetes deployments should use read-only root filesystem
3. **Resource limits**: Always set memory and CPU limits in production
4. **Rate limiting**: Adjust rate limits based on expected load
5. **Timeouts**: Configure timeouts based on expected request processing time

### Monitoring

Monitor these metrics:

- Response times (via logging middleware)
- Error rates (via error responses)
- Health check status (`/health`)
- Memory usage (from health endpoint)
- Goroutine count (from health endpoint)
- Rate limit hits (429 responses)

### Scaling

The application is stateless and can be scaled horizontally:

```bash
# Kubernetes
kubectl scale deployment gocalc-api --replicas=5

# Docker Compose
docker-compose up --scale api=5
```

Adjust rate limits when scaling:

- Total capacity = `RATE_LIMIT_RPM * number_of_instances`
- Use a reverse proxy/load balancer for distributed rate limiting

### Troubleshooting

**Container won't start:**

```bash
# Check logs
docker logs gocalc-api

# Common issues:
# - Invalid PORT (must be a number)
# - Invalid timeout format (use Go duration: 10s, 2m, etc.)
# - Port already in use
```

**Health check failing:**

```bash
# Test health endpoint
curl http://localhost:8080/health

# Check if service is listening
docker exec gocalc-api netstat -ln | grep 8080
```

**High memory usage:**

```bash
# Check health endpoint for memory stats
curl http://localhost:8080/health | jq '.data.checks.memory'

# Adjust resource limits
# Kubernetes: Update deployment.yaml resources.limits.memory
# Docker: Use --memory flag
```

## Cloud Platform Examples

### AWS ECS Task Definition

```json
{
  "family": "gocalc-api",
  "containerDefinitions": [
    {
      "name": "gocalc-api",
      "image": "gocalc-api:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "PORT",
          "value": "8080"
        },
        {
          "name": "RATE_LIMIT_RPM",
          "value": "200"
        }
      ],
      "healthCheck": {
        "command": ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1"],
        "interval": 30,
        "timeout": 3,
        "retries": 3,
        "startPeriod": 5
      }
    }
  ],
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512"
}
```

### Google Cloud Run

```bash
# Deploy to Cloud Run
gcloud run deploy gocalc-api \
  --image gocalc-api:latest \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --port 8080 \
  --set-env-vars "RATE_LIMIT_RPM=200,REQUEST_TIMEOUT=45s" \
  --max-instances 10 \
  --memory 256Mi \
  --cpu 1
```

### Azure Container Instances

```bash
# Create container instance
az container create \
  --resource-group myResourceGroup \
  --name gocalc-api \
  --image gocalc-api:latest \
  --ports 8080 \
  --environment-variables \
    PORT=8080 \
    RATE_LIMIT_RPM=200 \
  --cpu 1 \
  --memory 0.5
```
