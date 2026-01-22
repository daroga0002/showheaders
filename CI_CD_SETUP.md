# CI/CD Setup Summary

## ✅ Completed Infrastructure

### GitHub Actions Workflows

#### 1. CI Pipeline (`.github/workflows/ci.yml`)
**Triggers**: Push/PR to `main` and `llm` branches

**Jobs**:
- **Test**: Matrix testing with Go 1.24.3 and 1.24
  - Race detection enabled
  - Coverage reports generated and uploaded as artifacts
  
- **Lint**: Code quality checks
  - golangci-lint with comprehensive ruleset (see `.golangci.yml`)
  
- **Security**: Security scanning
  - gosec vulnerability scanner
  - SARIF results uploaded to GitHub Security tab
  
- **Build**: Cross-platform binary builds
  - Platforms: Linux, macOS, Windows
  - Architectures: amd64, arm64
  - Binaries uploaded as artifacts (7-day retention)

#### 2. Release Pipeline (`.github/workflows/release.yml`)
**Triggers**: Git tags matching `v*` (e.g., v1.0.0)

**Jobs**:
- **Build and Release**:
  - Runs full test suite
  - Builds binaries for all platforms with version injection
  - Generates SHA256 checksums
  - Creates GitHub release with artifacts
  
- **Docker**:
  - Multi-stage build (builder + minimal runtime)
  - Multi-architecture images (linux/amd64, linux/arm64)
  - Publishes to GitHub Container Registry (ghcr.io)
  - Tags: semantic version, major.minor, major, latest

### Docker Setup

**Dockerfile**:
- Multi-stage build optimized for size
- Builder: golang:1.24.3-alpine
- Runtime: alpine:latest with non-root user (UID 1000)
- Static binary with CGO disabled
- Health check on `/health` endpoint
- Default port: 8080

**Build locally**:
```bash
docker build -t showheaders:latest .
docker run -p 8080:8080 showheaders:latest
```

### Kubernetes Manifests

**deployment.yaml**:
- 2 replicas with rolling updates (maxUnavailable: 0, maxSurge: 1)
- Security context: non-root user (UID 1000)
- Resource limits: 50m-200m CPU, 64Mi-128Mi memory
- Probes: readiness (5s interval) and liveness (10s interval)
- Health checks via `/health` endpoint

**ingress.yaml**:
- Template for nginx ingress controller
- Configurable domain and TLS settings
- Ready for cert-manager integration

**Deploy**:
```bash
kubectl apply -f k8s/
kubectl rollout status deployment/showheaders
```

### Linting Configuration

**.golangci.yml**:
- Enabled linters: errcheck, gosimple, govet, staticcheck, gofmt, goimports, misspell, gocritic, revive, gosec, bodyclose, unconvert
- Type assertions and blank checks enabled
- Variable naming and exported symbols checked
- 5-minute timeout for large codebases

### Docker Ignore

**.dockerignore**:
- Excludes build artifacts, IDE files, git history
- Reduces context size and build time
- Improves layer caching

## 📋 Usage Guide

### Development Workflow

1. **Local development**:
   ```bash
   go run ./cmd/showheaders -port 8080
   go test ./... -v
   golangci-lint run ./...
   ```

2. **Push changes**:
   - Create PR → CI runs automatically
   - Merge to main → CI runs on main branch

3. **Create release**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
   - Release workflow builds binaries
   - Docker images pushed to ghcr.io/daroga0002/showheaders
   - GitHub release created with artifacts

### Kubernetes Deployment

1. **Update image tag** in `k8s/deployment.yaml`:
   ```yaml
   image: ghcr.io/daroga0002/showheaders:v1.0.0
   ```

2. **Apply manifests**:
   ```bash
   kubectl apply -f k8s/deployment.yaml
   kubectl rollout status deployment/showheaders
   ```

3. **Configure ingress** (optional):
   - Edit domain in `k8s/ingress.yaml`
   - Enable TLS if using cert-manager
   - Apply: `kubectl apply -f k8s/ingress.yaml`

### Docker Usage

**Pull from registry**:
```bash
docker pull ghcr.io/daroga0002/showheaders:latest
docker run -p 8080:8080 ghcr.io/daroga0002/showheaders:latest
```

**Build locally**:
```bash
docker build -t showheaders:dev .
docker run -p 8080:8080 showheaders:dev
```

**Multi-arch build**:
```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t showheaders:multi .
```

## 🔐 Security Features

1. **Non-root container**: Runs as UID 1000
2. **Gosec scanning**: Automated vulnerability detection
3. **Dependency verification**: `go mod verify` in CI
4. **SARIF uploads**: Security results in GitHub Security tab
5. **Read-only filesystem**: Compatible (no writes needed)
6. **Minimal base image**: Alpine Linux reduces attack surface

## 📊 Observability

**Logs**:
- Structured JSON logging with ZAP
- Request/response logging (excluding /health)
- Fields: method, path, status, duration, remote_addr, user_agent

**Health checks**:
- Kubernetes: readiness/liveness probes on `/health`
- Docker: HEALTHCHECK with 30s interval
- Response: `{"status":"healthy"}`

**Metrics** (future enhancement):
- Consider adding Prometheus endpoint
- Export request count, duration, status codes

## 🚀 Next Steps

Optional improvements:
1. Add Prometheus metrics endpoint
2. Implement distributed tracing (OpenTelemetry)
3. Add Helm chart for easier K8s deployment
4. Setup branch protection rules requiring CI success
5. Add performance/load testing to CI
6. Implement automated rollback on failed health checks
7. Add canary deployment strategy

## 📝 Notes

- **GHCR Authentication**: Workflows use `GITHUB_TOKEN` (automatic)
- **Manual releases**: Use GitHub UI or push tags
- **Branch strategy**: Main receives merges, llm for development
- **Artifact retention**: 7 days for CI, indefinite for releases
- **Image retention**: Configure via GHCR settings (default: keep all)
