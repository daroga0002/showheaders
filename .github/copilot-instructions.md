# Copilot Instructions for showheaders

## Project Overview
A production-grade HTTP debugging server that displays request headers in multiple formats (HTML, JSON, plain text). Built with Go 1.24.3, emphasizing clean architecture, middleware patterns, and graceful shutdown.

## Architecture Pattern
- **Layer separation**: `cmd/` (entry point) → `internal/` (business logic)
- **Internal structure**: `config/` (CLI parsing), `server/` (HTTP setup), `handlers/` (endpoints), `middleware/` (cross-cutting)
- **Dependency injection**: Components receive logger and config via constructors (see [server.go](../internal/server/server.go))
- **No external HTTP router**: Uses standard library `http.ServeMux` - avoid suggesting third-party routers

## Key Patterns

### Middleware Composition
Middleware is applied in [server.go](../internal/server/server.go#L33-L36) using function wrapping (not method chaining):
```go
handler := middleware.NoCacheMiddleware(
    middleware.LoggingMiddleware(logger, excludePaths)(mux),
)
```
- Execution order: NoCacheMiddleware → LoggingMiddleware → handlers
- LoggingMiddleware uses closure to inject logger and excludePaths configuration
- Always wrap `http.ResponseWriter` to capture status codes ([logging.go](../internal/middleware/logging.go#L10-L22))

### Handler Format Negotiation
The [ShowHeaders](../internal/handlers/headers.go) handler supports `?format=` query parameter:
- No param/default → HTML with styled table
- `?format=json` → JSON with structured response (`HeadersResponse` type)
- `?format=plain` → Plain text output
**Pattern**: Use query params for format selection, not Accept headers or separate endpoints

### Testing Strategy
- Use `httptest.NewRecorder()` and `httptest.NewRequest()` for handler tests ([headers_test.go](../internal/handlers/headers_test.go))
- Use `zaptest.NewLogger(t)` for test logging (never real logger in tests)
- For server tests, use `Port: 0` to let OS assign free ports ([server_test.go](../internal/server/server_test.go#L43))
- Test middleware by wrapping test handlers, not in isolation

### Graceful Shutdown
[main.go](../cmd/showheaders/main.go#L37-L62) implements signal-based shutdown:
1. Server runs in goroutine with error channel
2. `select` blocks on either server error or OS signal (SIGINT/SIGTERM)
3. 30-second timeout for graceful shutdown
4. Server must implement `Shutdown(context.Context)` method

## Build and Test Workflow
Use VS Code tasks (defined in `.vscode/tasks.json`):
- **Build**: Default task compiles to `./showheaders`
- **Cross-compile**: Separate tasks for Linux/Windows (use `GOOS`/`GOARCH` env vars)
- **Test**: `go test ./... -v` (default test task)
- **Run**: `go run ./cmd/showheaders` (not the binary directly in development)

**Manual commands**:
```bash
go build -o showheaders ./cmd/showheaders  # Build
go run ./cmd/showheaders -port 3000        # Run with custom port
go test ./... -v -coverprofile=coverage.out # Test with coverage
```

## Code Conventions
- **Logging**: Use structured logging with `go.uber.org/zap`, never fmt/log packages
- **Errors**: Log with `zap.Error(err)`, fatal only in main.go
- **Timeouts**: All HTTP servers must set `ReadTimeout`, `WriteTimeout`, `IdleTimeout` ([server.go](../internal/server/server.go#L42-L44))
- **Package naming**: Use singular names (`config`, not `configs`; `handler`, not `handlers` - exception: this project uses `handlers` for multiple handler types)
- **Exported types**: Document with godoc comments starting with type name

## Common Gotchas
- Middleware in `internal/middleware/logging.go` must call `next.ServeHTTP(wrapped, r)` with wrapped writer, not original `w`
- Health check endpoint excluded from logging via `excludePaths` map in [server.go](../internal/server/server.go#L28-L30)
- Headers are multi-value: use `r.Header[name]` (returns `[]string`), not `r.Header.Get(name)` when all values needed
- HTML output requires proper escaping - this project intentionally doesn't escape for debugging purposes

## Dependencies
- Single external dependency: `go.uber.org/zap` for structured logging
- Standard library preferred - avoid adding dependencies without strong justification

## CI/CD Pipeline

### GitHub Actions Workflows
- **[ci.yml](../workflows/ci.yml)**: Runs on push/PR to main and llm branches
  - Matrix testing: Go 1.24.3 and 1.24
  - Race detector enabled
  - Coverage reports uploaded as artifacts
  - Linting with golangci-lint (see [.golangci.yml](../../.golangci.yml))
  - Security scanning with gosec (SARIF upload)
  - Cross-platform builds (Linux, macOS, Windows × amd64, arm64)

- **[release.yml](../workflows/release.yml)**: Triggers on version tags (v*)
  - Builds binaries with version injection: `-ldflags="-X main.version=$TAG"`
  - Creates GitHub releases with checksums
  - Builds multi-arch Docker images (amd64, arm64)
  - Pushes to GitHub Container Registry (ghcr.io)

### Docker
- Multi-stage build in [Dockerfile](../../Dockerfile):
  1. Builder: golang:1.24.3-alpine with static binary compilation
  2. Runtime: alpine:latest with non-root user (UID 1000)
- Health check configured for `/health` endpoint
- Images tagged with semantic versions and `latest`

### Kubernetes Manifests
- [deployment.yaml](../../k8s/deployment.yaml): 2 replicas, rolling update strategy, resource limits
- [ingress.yaml](../../k8s/ingress.yaml): Template for nginx ingress (configure domain/TLS)
- Security context: non-root user, readiness/liveness probes

### Release Process
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
# Workflow automatically builds, tests, and publishes
```
