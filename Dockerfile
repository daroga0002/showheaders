# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24.3

FROM golang:${GO_VERSION} AS build
WORKDIR /src

ARG TARGETOS
ARG TARGETARCH

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
	go build -trimpath -ldflags="-s -w" -o /out/showheaders ./cmd/showheaders

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/showheaders /showheaders

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/showheaders"]
