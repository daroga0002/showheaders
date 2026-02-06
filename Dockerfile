# Build stage
FROM golang:1.24.13-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build arguments for versioning
ARG VERSION=dev
ARG BUILD_TIME=unknown

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
    -trimpath \
    -o showheaders \
    ./cmd/showheaders

# Final stage
FROM scratch

# Copy CA certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/showheaders /showheaders

# Use non-root user
USER 65534:65534

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/showheaders"]
