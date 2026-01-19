---
name: github-actions-pipeline
description: Guide for building GitHub Actions workflows for Go projects. Use this skill when asked to design, optimize, or troubleshoot CI/CD pipelines in GitHub Actions for Go.
---

# GitHub Actions Pipeline for Go Projects

## Overview

This skill provides comprehensive guidance for building **production-ready CI/CD pipelines** for Go applications using GitHub Actions. Optimized for the DevOps infinity loop: **Plan → Code → Build → Test → Release → Deploy → Operate → Monitor**.

## Core Principles

- **Fast Feedback**: Fail fast with parallel jobs and caching
- **Cross-Platform**: Test on Linux adm64 and arm64
- **Security First**: Scan dependencies and code for vulnerabilities
--**Linting & Quality**: Enforce coding standards automatically using golangci-lint
- **Artifact Management**: Preserve coverage reports and binaries
- **Progressive Enhancement**: Start simple, add complexity as needed

## Use This Skill When

- Designing a new GitHub Actions pipeline for a Go project
- Optimizing build times and reducing CI costs
- Adding security scanning or compliance checks
- Troubleshooting failing workflows or flaky tests

---

## Pipeline Stages

### Stage 1: Build & Test (Foundation)

**Purpose**: Validate code compiles and tests pass across platforms

**Implementation** (`.github/workflows/ci.yml`):

```yaml
name: CI

on:
	push:
		branches: [main]
	pull_request:
		branches: [main]
	workflow_dispatch:

permissions:
	contents: read

jobs:
	build-test:
		name: Build & Test (${{ matrix.os }} • Go ${{ matrix.go }})
		runs-on: ${{ matrix.os }}
		strategy:
			fail-fast: false
			matrix:
				os: [ubuntu-latest, macos-latest, windows-latest]
		timeout-minutes: 20
		concurrency:
			group: ci-${{ github.ref }}-${{ matrix.os }}-${{ matrix.go }}
			cancel-in-progress: true

		steps:
			- name: Checkout
				uses: actions/checkout@v4

			- name: Setup Go
				uses: actions/setup-go@v5
				with:
					go-version: ${{ matrix.go }}
					cache: true

			- name: Validate module
				run: |
					go mod tidy
					git diff --exit-code || (echo "Run 'go mod tidy' locally" && exit 1)

			- name: Build
				run: go build -v ./cmd/showheaders

			- name: Test (race detector)
				run: go test ./... -v -race -timeout=5m

			- name: Coverage (Linux only)
				if: matrix.os == 'ubuntu-latest' && matrix.go == '1.24.x'
				run: go test ./... -v -coverprofile=coverage.out -covermode=atomic

			- name: Upload coverage
				if: matrix.os == 'ubuntu-latest' && matrix.go == '1.24.x'
				uses: actions/upload-artifact@v4
				with:
					name: coverage
					path: coverage.out
					retention-days: 7
```

**Key Features**:
- Matrix builds catch platform-specific issues early
- Race detector finds concurrency bugs (critical for servers)
- Module validation prevents go.mod/go.sum drift
- Coverage only on one platform (faster, less duplication)

### Stage 2: Code Quality (Enhancement)

**Purpose**: Enforce code standards and catch common bugs

**Implementation**:

```yaml
	lint:
		name: Lint
		runs-on: ubuntu-latest
		steps:
			- uses: actions/checkout@v4
      
			- uses: actions/setup-go@v5
				with:
					go-version: "1.24.x"
					cache: true
      
			- name: golangci-lint
				uses: golangci/golangci-lint-action@v6
				with:
					version: v1.61.0
					args: --timeout=5m --out-format=github-actions

			- name: Go fmt check
				run: |
					if [ -n "$(gofmt -l .)" ]; then
						echo "Run 'gofmt -w .' to fix formatting"
						gofmt -l .
						exit 1
					fi

			- name: Go vet
				run: go vet ./...
```

**Recommended .golangci.yml**:

```yaml
run:
	timeout: 5m
	tests: true

linters:
	enable:
		- errcheck
		- gosimple
		- govet
		- ineffassign
		- staticcheck
		- unused
		- gofmt
		- revive
		- gosec

issues:
	exclude-dirs:
		- vendor
	max-same-issues: 0
```

### Stage 3: Security Scanning (Critical)

**Purpose**: Detect vulnerabilities early in development

**Implementation**:

```yaml
	security:
		name: Security Scan
		runs-on: ubuntu-latest
		permissions:
			contents: read
			security-events: write
		steps:
			- uses: actions/checkout@v4
      
			- uses: actions/setup-go@v5
				with:
					go-version: "1.24.x"
					cache: true

			- name: Go vulnerability check
				run: |
					go install golang.org/x/vuln/cmd/govulncheck@latest
					govulncheck ./...

			- name: Dependency scanning
				uses: aquasecurity/trivy-action@master
				with:
					scan-type: 'fs'
					scan-ref: '.'
					format: 'sarif'
					output: 'trivy-results.sarif'

			- name: Upload to Security tab
				uses: github/codeql-action/upload-sarif@v3
				with:
					sarif_file: 'trivy-results.sarif'
```

### Stage 4: Build Artifacts (Release Prep)

**Purpose**: Create distributable binaries for multiple platforms

**Implementation**:

```yaml
	build-artifacts:
		name: Build Artifacts
		runs-on: ubuntu-latest
		if: github.event_name == 'push' && startsWith(github.ref, 'refs/tags/')
		needs: [build-test, lint, security]
		strategy:
			matrix:
				include:
					- os: linux
						arch: amd64
					- os: linux
						arch: arm64
					- os: darwin
						arch: amd64
					- os: darwin
						arch: arm64
					- os: windows
						arch: amd64
		steps:
			- uses: actions/checkout@v4
      
			- uses: actions/setup-go@v5
				with:
					go-version: "1.24.x"
					cache: true

			- name: Build binary
				env:
					GOOS: ${{ matrix.os }}
					GOARCH: ${{ matrix.arch }}
				run: |
					EXT=""
					if [ "${{ matrix.os }}" = "windows" ]; then EXT=".exe"; fi
					go build -v -trimpath -ldflags="-s -w" \
						-o showheaders-${{ matrix.os }}-${{ matrix.arch }}${EXT} \
						./cmd/showheaders

			- name: Upload artifact
				uses: actions/upload-artifact@v4
				with:
					name: showheaders-${{ matrix.os }}-${{ matrix.arch }}
					path: showheaders-*
```

**Build Flags Explained**:
- -trimpath: Remove absolute paths (reproducible builds)
- -ldflags="-s -w": Strip debug info (smaller binaries)
- -v: Verbose output for debugging

### Stage 5: Container Image (Kubernetes Ready)

**Purpose**: Build and push Docker images for deployment

**Implementation**:

```yaml
	docker:
		name: Build & Push Docker Image
		runs-on: ubuntu-latest
		if: github.event_name == 'push' && (github.ref == 'refs/heads/main' || startsWith(github.ref, 'refs/tags/'))
		needs: [build-test, lint, security]
		permissions:
			contents: read
			packages: write
		steps:
			- uses: actions/checkout@v4

			- name: Set up Docker Buildx
				uses: docker/setup-buildx-action@v3

			- name: Login to GitHub Container Registry
				uses: docker/login-action@v3
				with:
					registry: ghcr.io
					username: ${{ github.actor }}
					password: ${{ secrets.GITHUB_TOKEN }}

			- name: Extract metadata
				id: meta
				uses: docker/metadata-action@v5
				with:
					images: ghcr.io/${{ github.repository }}
					tags: |
						type=ref,event=branch
						type=ref,event=pr
						type=semver,pattern={{version}}
						type=semver,pattern={{major}}.{{minor}}
						type=sha,prefix={{branch}}-

			- name: Build and push
				uses: docker/build-push-action@v5
				with:
					context: .
					platforms: linux/amd64,linux/arm64
					push: true
					tags: ${{ steps.meta.outputs.tags }}
					labels: ${{ steps.meta.outputs.labels }}
					cache-from: type=gha
					cache-to: type=gha,mode=max
```

**Dockerfile Best Practices**:

```dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
		-o showheaders ./cmd/showheaders

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /build/showheaders /showheaders

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/showheaders"]
```

### Stage 6: Deployment (Kubernetes)

**Purpose**: Automated deployment to Kubernetes

**Implementation**:

```yaml
	deploy:
		name: Deploy to Kubernetes
		runs-on: ubuntu-latest
		if: github.event_name == 'push' && github.ref == 'refs/heads/main'
		needs: [docker]
		environment:
			name: production
			url: https://showheaders.example.com
		steps:
			- uses: actions/checkout@v4

			- name: Install kubectl
				uses: azure/setup-kubectl@v4
				with:
					version: 'v1.30.0'

			- name: Configure kubeconfig
				run: |
					mkdir -p $HOME/.kube
					echo "${{ secrets.KUBE_CONFIG }}" | base64 -d > $HOME/.kube/config

			- name: Update deployment image
				run: |
					kubectl set image deployment/showheaders \
						showheaders=ghcr.io/${{ github.repository }}:main-${{ github.sha }} \
						-n production

			- name: Wait for rollout
				run: kubectl rollout status deployment/showheaders -n production --timeout=5m

			- name: Verify deployment
				run: |
					kubectl get pods -n production -l app=showheaders
					kubectl logs -n production -l app=showheaders --tail=50
```

**Required Secrets**:
- KUBE_CONFIG: Base64-encoded kubeconfig file

---

## Performance Optimization

### Caching Strategy

**Go Module Cache** (automatic):
```yaml
- uses: actions/setup-go@v5
	with:
		go-version: "1.24.x"
		cache: true
```

**Docker Layer Cache**:
```yaml
- uses: docker/build-push-action@v5
	with:
		cache-from: type=gha
		cache-to: type=gha,mode=max
```

### Parallelization

Run independent jobs concurrently:

```yaml
jobs:
	build-test:
		# Runs immediately
  
	lint:
		# Runs in parallel with build-test
  
	security:
		# Runs in parallel with build-test
  
	deploy:
		needs: [build-test, lint, security]
```

### Path Filtering

Skip CI on documentation changes:

```yaml
on:
	push:
		paths:
			- '**/*.go'
			- 'go.mod'
			- 'go.sum'
			- '.github/workflows/**'
```

---

## Security Best Practices

### Minimal Permissions

```yaml
permissions:
	contents: read
	packages: write
	security-events: write
```

### Secret Management

**Use GitHub Secrets**:
```yaml
- run: echo "TOKEN=${{ secrets.MY_TOKEN }}" >> $GITHUB_ENV
```

### Dependency Pinning

```yaml
# Pin to commit SHA
- uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11
```

---

## Troubleshooting

### Common Issues

**Tests pass locally but fail in CI**
→ Check for race conditions, time-dependent tests, environment assumptions

**Slow builds**
→ Enable caching, reduce matrix size, parallelize jobs

**Flaky tests**
→ Add retries:
```yaml
- uses: nick-fields/retry@v3
	with:
		timeout_minutes: 10
		max_attempts: 3
		command: go test ./... -v
```

**Out of disk space**
→ Clean up:
```yaml
- run: |
		sudo rm -rf /usr/share/dotnet
		docker system prune -af
```

---

## Complete Pipeline Architecture

For the showheaders project:

1. On every push/PR: Build, test, lint, security scan
2. On tag push: Build artifacts, create GitHub release
3. On main push: Build Docker image, deploy to Kubernetes
4. Coverage reports: Uploaded as artifacts
5. Security findings: Uploaded to GitHub Security tab

---

## Quick Start Checklist

- [ ] Create .github/workflows/ci.yml
- [ ] Add matrix for OS and Go versions
- [ ] Enable caching with actions/setup-go@v5
- [ ] Add linting with golangci-lint-action
- [ ] Add security scanning (govulncheck, Trivy)
- [ ] Create Dockerfile with multi-stage build
- [ ] Add Docker build job
- [ ] Configure GHCR authentication
- [ ] Create Kubernetes manifests in k8s/
- [ ] Add deployment job with environment protection
- [ ] Set up required secrets
- [ ] Add status badge to README

---

## Local Testing

Test workflows locally using act:

```bash
# Install act
brew install act

# Run CI workflow
act -W .github/workflows/ci.yml

# Run specific job
act -j build-test
```

---

## Metrics to Track

- Build duration (target: <5 minutes)
- Test coverage (target: >80%)
- Time to deployment (target: <10 minutes)
- Failure rate (target: <5% flakiness)

---

## Resources

- https://docs.github.com/en/actions
- https://golang.org/pkg/testing/
- https://golangci-lint.run/usage/linters/
- https://docs.docker.com/develop/dev-best-practices/
- https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

---

Remember: Start simple, add complexity incrementally. Every addition should solve a real problem.
name: github-actions-pipeline
description: Guide for building GitHub Actions workflows for Go projects. Use this skill when asked to design, optimize, or troubleshoot CI/CD pipelines in GitHub Actions for Go.
---

# GitHub Actions Pipeline for Go Projects

## Overview

This skill provides comprehensive guidance for building **production-ready CI/CD pipelines** for Go applications using GitHub Actions. Optimized for the DevOps infinity loop: **Plan → Code → Build → Test → Release → Deploy → Operate → Monitor**.

## Core Principles

- **Fast Feedback**: Fail fast with parallel jobs and caching
- **Cross-Platform**: Test on Linux adm64 and arm64
- **Security First**: Scan dependencies and code for vulnerabilities
--**Linting & Quality**: Enforce coding standards automatically using golangci-lint
- **Artifact Management**: Preserve coverage reports and binaries
- **Progressive Enhancement**: Start simple, add complexity as needed

## Use This Skill When

- Designing a new GitHub Actions pipeline for a Go project
- Optimizing build times and reducing CI costs
- Adding security scanning or compliance checks
- Troubleshooting failing workflows or flaky tests

---
