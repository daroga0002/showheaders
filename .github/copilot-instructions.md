# Copilot Instructions for showheaders

## Project Overview

`showheaders` is a diagnostic HTTP server that displays request headers in multiple formats (HTML, JSON, plain text). It's built with Go 1.24.3, using minimal dependencies: only `go.uber.org/zap` for structured logging.

## Architecture

### Module Structure (Go internal convention)
- `cmd/showheaders/` - Application entry point with main.go
- `internal/` - Private packages not importable outside this module
  - `config/` - CLI parsing using stdlib `flag` package
  - `handlers/` - HTTP handlers with format negotiation via `?format=` query param
  - `middleware/` - HTTP middleware for logging (excludes `/health`) and no-cache headers
  - `server/` - Server initialization and graceful shutdown logic

### Request Flow
1. Request → `NoCacheMiddleware` → `LoggingMiddleware` → Handler
2. `LoggingMiddleware` wraps `http.ResponseWriter` to capture status codes
3. `/health` endpoint explicitly excluded from logging (see `excludePaths` map in [server.go](internal/server/server.go#L31-L33))

## Development Workflows

### Building & Running
```bash
# Default build (uses VS Code task: "Build")
go build -o showheaders ./cmd/showheaders

# Run with custom port
go run ./cmd/showheaders -port 3000

# Cross-platform builds (see tasks.json for Linux/Windows variants)
GOOS=linux GOARCH=amd64 go build -o showheaders-linux ./cmd/showheaders
```

### Testing
```bash
# Run all tests with verbose output (VS Code task: "Test")
go test ./... -v

# Coverage report (VS Code task: "Test with Coverage")
go test ./... -v -coverprofile=coverage.out
```

### Debugging
Three debug configurations in `.vscode/launch.json`:
- Default (port 8080)
- Custom port 3000
- Localhost-only binding

## Code Conventions

### Handler Pattern
All handlers in `internal/handlers/` support format negotiation:
```go
func SomeHandler(w http.ResponseWriter, r *http.Request) {
    format := r.URL.Query().Get("format")
    switch format {
    case "json": // JSON response
    case "plain": // Plain text
    default: // HTML (styled)
    }
}
```

### Testing Pattern
- Use `httptest.NewRecorder()` for handler testing
- Table-driven tests for format variations (see [headers_test.go](internal/handlers/headers_test.go))
- Test structure: setup → execute → verify status/content-type/body

### Middleware Pattern
- Middleware functions return `func(http.Handler) http.Handler`
- Exclude paths via map: `excludePaths map[string]bool` (e.g., health checks)
- Chain middleware in [server.go](internal/server/server.go#L35-L37): `NoCacheMiddleware(LoggingMiddleware(...)(mux))`

### Logging
- Use ZAP structured logging: `logger.Info("msg", zap.String("key", val))`
- Health checks (`/health`) are intentionally NOT logged to reduce noise
- Server lifecycle events (start/shutdown/errors) always logged

### Graceful Shutdown
Pattern in [main.go](cmd/showheaders/main.go):
1. Server starts in goroutine with error channel
2. Signal handling (`os.Interrupt`, `SIGTERM`) via `signal.Notify`
3. `select{}` blocks until error or signal
4. 30-second timeout for shutdown (`context.WithTimeout`)

## Key Files
- [cmd/showheaders/main.go](cmd/showheaders/main.go) - Application lifecycle and graceful shutdown
- [internal/server/server.go](internal/server/server.go) - Server setup, middleware chain, timeouts (15s read/write, 60s idle)
- [internal/handlers/headers.go](internal/handlers/headers.go) - Format negotiation pattern, sorted header display
- [internal/middleware/logging.go](internal/middleware/logging.go) - responseWriter wrapper to capture status codes
