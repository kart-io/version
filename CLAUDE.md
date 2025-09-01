# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go version management package (`github.com/kart-io/version`) that provides comprehensive version information management for Go applications. The package supports:

- Build-time version injection via Go's `-ldflags` mechanism
- Multi-dimensional version information (Git version, commit, branch, build time, runtime environment)
- Multiple output formats (simple string, JSON, detailed table)
- Dynamic version management at runtime
- Command-line flag integration for `--version` support

## Development Commands

### Testing
```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -cover ./...

# Run specific test
go test -v -run TestGet
```

### Building
```bash
# Basic build (without version injection)
go build ./...

# Build with version injection (example)
go build -ldflags "
  -X 'github.com/kart-io/version.serviceName=myservice'
  -X 'github.com/kart-io/version.gitVersion=v1.0.0'
  -X 'github.com/kart-io/version.gitCommit=$(git rev-parse HEAD)'
  -X 'github.com/kart-io/version.gitBranch=$(git branch --show-current)'
  -X 'github.com/kart-io/version.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'
" ./...
```

### Linting and Formatting
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run golangci-lint (if available)
golangci-lint run
```

## Architecture and Design

### Core Components

1. **Static Variables** (`version.go`): Package-level variables that serve as version information carriers, injected via ldflags during build time
2. **Info Struct** (`version.go`): Central structure that aggregates all version-related information with JSON serialization support
3. **Dynamic Version Management** (`dynamic.go`): Thread-safe runtime version overrides using `atomic.Value` with semantic version validation
4. **Command-Line Integration** (`flag.go`): pflag-based version flag support with multiple output modes

### Key Files

- `version.go` - Core version information structure and output formatters
- `dynamic.go` - Runtime dynamic version management with validation
- `flag.go` - Command-line flag integration for version queries
- `semver.go` - Internal semantic version parsing and comparison
- `doc.go` - Package documentation and usage examples
- `version_test.go` - Unit tests for core functionality
- `semver_test.go` - Unit tests for semantic version parsing

### Data Flow

```
Build Time: Git Repository → Build System (ldflags) → Binary File (embedded variables)
Runtime:    Binary File → Version Package (Get() function) → Application (various output formats)
```

### Version Variable Injection

The package uses these ldflags-injectable variables:
- `github.com/kart-io/version.serviceName` - Service name
- `github.com/kart-io/version.gitVersion` - Git version tag  
- `github.com/kart-io/version.gitCommit` - Git commit SHA
- `github.com/kart-io/version.gitBranch` - Git branch name
- `github.com/kart-io/version.gitTreeState` - Git repository state (clean/dirty)
- `github.com/kart-io/version.buildDate` - ISO8601 build timestamp

### Development Commands Available

Use the provided Makefile for common development tasks:
```bash
make help        # Show all available commands
make fmt         # Format code
make test        # Run tests
make build       # Build with version injection
make check       # Run all checks
make ci          # Run full CI pipeline
```

### Dependencies

- `github.com/gosuri/uitable` - For tabular text output formatting
- `github.com/spf13/pflag` - For command-line flag handling
- `github.com/stretchr/testify/assert` - For unit testing assertions

### Internal Components

- `semver.go` - Internal semantic version parsing and validation (no external dependencies)

## Design Patterns

### Thread-Safe Dynamic Versioning
Uses `atomic.Value` for concurrent-safe runtime version modifications, ensuring no race conditions during version queries or updates.

### Multiple Output Formats
- `String()` - Simple version string for quick display
- `ToJSON()` - Structured JSON for APIs and logging
- `Text()` - Human-readable table format for detailed inspection

### Graceful Degradation  
System gracefully handles missing version information with sensible defaults and fallback mechanisms.

## Common Development Patterns

### Version Information Usage
```go
info := version.Get()
fmt.Printf("Version: %s\n", info.String())     // Simple output
fmt.Printf("Details:\n%s\n", info.Text())      // Detailed table
fmt.Printf("JSON: %s\n", info.ToJSON())        // JSON format
```

### Dynamic Version Setting
```go
if err := version.SetDynamicVersion("v1.2.3-hotfix.1"); err != nil {
    log.Fatal("Invalid version: ", err)
}
```

### Command-Line Integration
```go
version.AddFlags(pflag.CommandLine)
pflag.Parse()
version.PrintAndExitIfRequested()  // Handles --version flag
```

## Testing Considerations

- Tests cover all output formats and edge cases
- Dynamic version validation is thoroughly tested
- Concurrent access patterns are validated
- Build-time injection scenarios are simulated

The package is designed for production use in microservice architectures where comprehensive version tracking and reporting is essential.