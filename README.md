# showheaders

[![CI](https://github.com/daroga0002/showheaders/actions/workflows/ci.yml/badge.svg)](https://github.com/daroga0002/showheaders/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/daroga0002/showheaders)](https://goreportcard.com/report/github.com/daroga0002/showheaders)

A simple HTTP server written in Go that displays all HTTP headers from incoming requests.

## Features

- Displays all HTTP request headers in a styled HTML table
- Health check endpoint (`/health`)
- Structured JSON logging with ZAP
- Configurable port and hostname via CLI arguments
- Graceful shutdown handling

## Requirements

- Go 1.21 or later

## Building

```bash
# Clone the repository
git clone https://github.com/daroga0002/showheaders.git
cd showheaders

# Download dependencies
go mod download

# Build the application
go build -o showheaders ./cmd/showheaders
```

## Running

```bash
# Run with default settings (port 8080, all interfaces)
./showheaders

# Run on a custom port
./showheaders -port 3000

# Bind to a specific hostname
./showheaders -hostname localhost -port 8080

# Run directly with go
go run ./cmd/showheaders -port 8080
```

## CLI Arguments

| Argument    | Default | Description                                      |
|-------------|---------|--------------------------------------------------|
| `-port`     | `8080`  | Port number to listen on                         |
| `-hostname` | `""`    | Hostname to bind to (empty = all interfaces)     |

## Endpoints

| Endpoint  | Method | Description                              |
|-----------|--------|------------------------------------------|
| `/`       | GET    | Displays all HTTP headers from request   |
| `/health` | GET    | Health check endpoint (returns JSON)     |

### Query Parameters

The `/` endpoint supports a `format` query parameter to control the output format:

| Format   | URL Example         | Description                        |
|----------|---------------------|------------------------------------|
| (none)   | `/`                 | HTML table (default)               |
| `json`   | `/?format=json`     | JSON object with all headers       |
| `plain`  | `/?format=plain`    | Plain text output                  |

#### JSON Format Example

```bash
curl "http://localhost:8080/?format=json"
```

```json
{
  "method": "GET",
  "url": "/?format=json",
  "remote_addr": "127.0.0.1:54321",
  "headers": {
    "User-Agent": ["curl/8.0.0"],
    "Accept": ["*/*"]
  }
}
```

#### Plain Text Format Example

```bash
curl "http://localhost:8080/?format=plain"
```

```
Method: GET
URL: /?format=plain
Remote Address: 127.0.0.1:54321

--- Headers ---
Accept: */*
User-Agent: curl/8.0.0
```

### Health Check Response

```json
{
  "status": "healthy"
}
```

## Project Structure

```
.
├── cmd/
│   └── showheaders/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # CLI argument parsing
│   ├── handlers/
│   │   ├── headers.go           # Show headers handler
│   │   └── health.go            # Health check handler
│   ├── middleware/
│   │   └── logging.go           # Request logging middleware
│   └── server/
│       └── server.go            # HTTP server setup
├── go.mod
├── go.sum
└── README.md
```

## Logging

The application uses [ZAP](https://github.com/uber-go/zap) for structured JSON logging. Each request (except health checks) is logged with the following fields:

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:00.000Z",
  "msg": "request",
  "method": "GET",
  "path": "/",
  "query": "",
  "status": 200,
  "duration": "1.234ms",
  "remote_addr": "127.0.0.1:54321",
  "user_agent": "Mozilla/5.0..."
}
```

**Note:** Health check requests (`/health`) are excluded from logging to reduce noise in production environments.

## CI/CD Pipeline

The project includes a comprehensive GitHub Actions CI/CD pipeline that runs on every push and pull request:

### Pipeline Stages

1. **Security Scanning**
   - Secret scanning with TruffleHog and Gitleaks
   - Static Application Security Testing (SAST) with CodeQL
   - Go-specific security checks with Gosec
   - Vulnerability scanning with govulncheck

2. **Code Quality**
   - Linting with golangci-lint
   - Format checking with gofmt
   - Dependency verification

3. **Testing**
   - Unit tests with race detection
   - Coverage report upload to Codecov
   - Cross-platform testing (Linux, macOS, Windows)

4. **Building**
   - Multi-platform builds (Linux, macOS, Windows)
   - Multi-architecture support (amd64, arm64)
   - Build artifacts uploaded for 30 days

5. **Docker Image** (main branch only)
   - Multi-architecture Docker images
   - Published to GitHub Container Registry (ghcr.io)
   - Automatic tagging and versioning

6. **Release** (tags only)
   - Automatic GitHub release creation
   - All platform binaries attached
   - SHA256 checksums included

### Running Pipeline Locally

You can run various checks locally before pushing:

```bash
# Run tests
go test -v -race ./...

# Run linter
golangci-lint run

# Check formatting
gofmt -s -l .

# Build for current platform
go build -o showheaders ./cmd/showheaders

# Build Docker image
docker build -t showheaders:local .
```

### Creating a Release

To create a new release with binaries:

1. Tag your commit: `git tag -a v1.0.0 -m "Release v1.0.0"`
2. Push the tag: `git push origin v1.0.0`
3. The pipeline will automatically create a GitHub release with binaries

## Graceful Shutdown

The server handles graceful shutdown when receiving `SIGINT` (Ctrl+C) or `SIGTERM` signals:

1. Stops accepting new connections
2. Waits up to 30 seconds for in-flight requests to complete
3. Logs shutdown progress

## Development

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o showheaders-linux ./cmd/showheaders

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o showheaders-darwin ./cmd/showheaders

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o showheaders-darwin-arm64 ./cmd/showheaders

# Windows
GOOS=windows GOARCH=amd64 go build -o showheaders.exe ./cmd/showheaders
```

## License

MIT
