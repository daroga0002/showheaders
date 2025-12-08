# showheaders

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
git clone <repository-url>
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
