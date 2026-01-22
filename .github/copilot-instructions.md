# Copilot instructions (showheaders)

## Big picture
- Single-binary Go HTTP server. Entry point: [cmd/showheaders/main.go](cmd/showheaders/main.go).
- HTTP wiring lives in [internal/server/server.go](internal/server/server.go): `http.ServeMux` routes + middleware chain.
- Request handlers are simple `http.HandlerFunc`s in [internal/handlers](internal/handlers).
- Cross-cutting concerns live in [internal/middleware](internal/middleware) (Zap logging + no-cache headers).

## Key flows & conventions
- Server lifecycle: `main` creates a Zap production logger, parses CLI config, starts server, then handles SIGINT/SIGTERM with a 30s graceful shutdown timeout.
- Routing:
  - `/` → `handlers.ShowHeaders` (supports `?format=json|plain`, default HTML)
  - `/health` → `handlers.Health` (JSON `{status:"healthy"}`)
- Middleware order (important): in `server.New`, the mux is wrapped by `LoggingMiddleware` (with excluded paths), then by `NoCacheMiddleware`.
- Logging convention: structured Zap `Info` message name is `"request"`; `/health` is excluded from logging by default via `excludePaths` in `server.New`.

## Developer workflows
- Build: `go build -o showheaders ./cmd/showheaders`
- Run: `go run ./cmd/showheaders` (flags: `-port`, `-hostname`; see [internal/config/config.go](internal/config/config.go))
- Tests: `go test ./... -v`
- VS Code tasks (optional): use the built-in tasks for `Build`, `Test`, `Run`, and `Tidy Modules` instead of typing commands.

## Repo-specific testing patterns
- Handlers/middleware are tested with `net/http/httptest` (see [internal/handlers/headers_test.go](internal/handlers/headers_test.go)).
- For Zap in tests, prefer `zaptest.NewLogger(t)` (see [internal/server/server_test.go](internal/server/server_test.go)).
- When testing CLI parsing, reset `flag.CommandLine` and `os.Args` before calling `config.Parse()` (see [internal/config/config_test.go](internal/config/config_test.go)).

## When making changes
- Prefer keeping behavior in small, package-local helpers (e.g., `showHeadersJSON/plain/HTML` in [internal/handlers/headers.go](internal/handlers/headers.go)).
- If you add a new endpoint, wire it in [internal/server/server.go](internal/server/server.go) and decide whether it should be excluded from request logging.
  - Add “noisy” endpoints (e.g., `/metrics`, `/favicon.ico`) to `excludePaths` to avoid log spam.
