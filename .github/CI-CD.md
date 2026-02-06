# CI/CD Pipeline Documentation

This document describes the continuous integration and deployment pipelines for the showheaders project.

## Overview

The project uses GitHub Actions for automated CI/CD with the following workflows:

- **CI Pipeline**: Runs on every push and pull request
- **Release Pipeline**: Triggers on version tags

## CI Pipeline (`.github/workflows/ci.yml`)

### Triggers
- Push to `main` branch
- Pull requests to `main` branch

### Jobs

#### 1. Lint
- Runs `golangci-lint` with comprehensive linter suite
- Configuration: [.golangci.yml](../.golangci.yml)
- Enabled linters: gofmt, goimports, govet, errcheck, staticcheck, gosec, and more

#### 2. Test (Matrix)
- **Platforms**: Ubuntu, macOS, Windows
- **Go version**: 1.24.3
- **Steps**:
  - Download and verify dependencies
  - Run tests with race detector
  - Generate coverage reports
  - Upload coverage to Codecov (Ubuntu only)

#### 3. Build (Matrix)
- **Targets**: 
  - linux/amd64, linux/arm64
  - darwin/amd64, darwin/arm64
  - windows/amd64
- Produces binaries for all platforms
- Uploads artifacts (retained for 7 days)

#### 4. Security
- **Gosec**: Static security scanner for Go code
- **Trivy**: Vulnerability scanner for dependencies
- Results uploaded to GitHub Security tab (SARIF format)

## Release Pipeline (`.github/workflows/release.yml`)

### Triggers
- Git tags matching pattern `v*` (e.g., `v1.0.0`)

### Process
1. Runs full test suite
2. Builds optimized binaries for all platforms:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64)
3. Generates SHA256 checksums
4. Creates GitHub Release with:
   - Automated release notes from git commits
   - All platform binaries
   - Checksums file

### Creating a Release

```bash
# Tag the release
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# - Build binaries
# - Create release
# - Upload artifacts
```

## Dependency Management

**Dependabot** (`.github/dependabot.yml`) automatically:
- Updates Go modules weekly (Mondays)
- Updates GitHub Actions weekly
- Creates pull requests with dependency updates
- Labels: `dependencies`, `go`, `github-actions`

## Local Development

### Run Linting
```bash
# Install golangci-lint
brew install golangci-lint  # macOS
# or
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run linters
golangci-lint run ./...
```

### Run Tests
```bash
# Basic tests
go test ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Security Scans
```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Run scan
gosec ./...
```

### Build locally
```bash
# Current platform
go build -o showheaders ./cmd/showheaders

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o showheaders-linux-amd64 ./cmd/showheaders
```

## CI/CD Metrics

### Success Criteria
- ✅ All linters pass
- ✅ 100% test pass rate on all platforms
- ✅ No security vulnerabilities (critical/high)
- ✅ Successful builds for all target platforms

### Pipeline Performance
- **Lint job**: ~2-3 minutes
- **Test matrix**: ~5-8 minutes (parallel)
- **Build matrix**: ~3-5 minutes (parallel)
- **Security scan**: ~3-4 minutes
- **Total CI time**: ~10-15 minutes

## Troubleshooting

### Lint Failures
Check [.golangci.yml](../.golangci.yml) configuration and run locally:
```bash
golangci-lint run ./...
```

### Test Failures
Run with verbose output:
```bash
go test -v ./...
```

### Security Issues
Review GitHub Security tab for detailed vulnerability reports.

### Build Failures
Verify Go version matches CI:
```bash
go version  # Should be 1.24.3
```

## Best Practices

1. **Before pushing**: Run tests and linters locally
2. **Pull requests**: Wait for CI to pass before merging
3. **Dependencies**: Review Dependabot PRs weekly
4. **Releases**: Use semantic versioning (vX.Y.Z)
5. **Security**: Address vulnerabilities promptly

## Future Enhancements

- [ ] Docker image builds
- [ ] Kubernetes manifests
- [ ] Performance benchmarking
- [ ] Integration tests
- [ ] Deployment to staging/production
- [ ] Chaos engineering tests

## Related Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Gosec Security Scanner](https://github.com/securego/gosec)
- [Trivy Vulnerability Scanner](https://github.com/aquasecurity/trivy)
