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

## Graceful Shutdown

The server handles graceful shutdown when receiving `SIGINT` (Ctrl+C) or `SIGTERM` signals:

1. Stops accepting new connections
2. Waits up to 30 seconds for in-flight requests to complete
3. Logs shutdown progress

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run tests with race detector
go test -race ./...
```

### Linting

```bash
# Install golangci-lint (if not installed)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run ./...
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

## Docker

### Building the Image

```bash
docker build -t showheaders:latest .
```

### Running with Docker

```bash
# Run on port 8080
docker run -p 8080:8080 showheaders:latest

# Run on custom port
docker run -p 3000:8080 showheaders:latest -port 8080
```

### Multi-architecture Build

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t showheaders:latest .
```

## Kubernetes Deployment

### Prerequisites
- Kubernetes cluster (local or remote)
- kubectl configured
- Container image pushed to registry (e.g., GHCR)

### Deploy to Kubernetes

```bash
# Apply deployment and service
kubectl apply -f k8s/deployment.yaml

# Check deployment status
kubectl rollout status deployment/showheaders

# Get pods
kubectl get pods -l app=showheaders

# Test the service
kubectl port-forward svc/showheaders 8080:80
```

### Configure Ingress (optional)

```bash
# Edit k8s/ingress.yaml with your domain
# Then apply:
kubectl apply -f k8s/ingress.yaml
```

## CI/CD

The project includes GitHub Actions workflows:

- **CI Pipeline** (`.github/workflows/ci.yml`): Runs on every push and PR
  - Tests with multiple Go versions
  - Linting with golangci-lint
  - Security scanning with gosec
  - Cross-platform builds (Linux, macOS, Windows)

- **Release Pipeline** (`.github/workflows/release.yml`): Triggers on version tags
  - Builds binaries for all platforms
  - Creates GitHub releases with checksums
  - Builds and pushes Docker images to GHCR
  - Supports multi-architecture images (amd64, arm64)

### Creating a Release

```bash
# Tag a new version
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# The release workflow will automatically:
# - Build binaries for all platforms
# - Create a GitHub release
# - Build and push Docker images
```

## License

MIT
