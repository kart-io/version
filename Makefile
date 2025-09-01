# Makefile for kart-io/version package

# Project information
PROJECT_NAME := version
MODULE_NAME := github.com/kart-io/version
VERSION_PKG := $(MODULE_NAME)

# Git information
GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git branch --show-current 2>/dev/null || echo "unknown")
GIT_TREE_STATE := $(shell if [ -n "$$(git status --porcelain 2>/dev/null)" ]; then echo "dirty"; else echo "clean"; fi)
BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Build flags
LDFLAGS := -w -s \
	-X '$(VERSION_PKG).serviceName=version-service' \
	-X '$(VERSION_PKG).gitVersion=$(GIT_VERSION)' \
	-X '$(VERSION_PKG).gitCommit=$(GIT_COMMIT)' \
	-X '$(VERSION_PKG).gitBranch=$(GIT_BRANCH)' \
	-X '$(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE)' \
	-X '$(VERSION_PKG).buildDate=$(BUILD_DATE)'

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}'

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(GREEN)[INFO]$(NC) Formatting Go code..."
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	@echo "$(GREEN)[INFO]$(NC) Running go vet..."
	go vet ./...

.PHONY: lint
lint: vet ## Run linters (go vet)
	@echo "$(GREEN)[INFO]$(NC) Linting completed"

.PHONY: test
test: ## Run tests
	@echo "$(GREEN)[INFO]$(NC) Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(GREEN)[INFO]$(NC) Running tests with coverage..."
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)[INFO]$(NC) Coverage report generated: coverage.html"

.PHONY: test-race
test-race: ## Run tests with race detection
	@echo "$(GREEN)[INFO]$(NC) Running tests with race detection..."
	go test -v -race ./...

.PHONY: build
build: ## Build the library (validation build)
	@echo "$(GREEN)[INFO]$(NC) Building library..."
	@echo "Version: $(GIT_VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Branch: $(GIT_BRANCH)"
	@echo "Date: $(BUILD_DATE)"
	go build -ldflags "$(LDFLAGS)" ./...

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(GREEN)[INFO]$(NC) Cleaning build artifacts..."
	rm -f coverage.out coverage.html

.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Module: $(MODULE_NAME)"
	@echo "Version: $(GIT_VERSION)"
	@echo "Commit: $(GIT_COMMIT)"
	@echo "Branch: $(GIT_BRANCH)"
	@echo "Tree State: $(GIT_TREE_STATE)"
	@echo "Build Date: $(BUILD_DATE)"

.PHONY: deps
deps: ## Download dependencies
	@echo "$(GREEN)[INFO]$(NC) Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "$(GREEN)[INFO]$(NC) Updating dependencies..."
	go get -u ./...
	go mod tidy

.PHONY: check
check: fmt vet test ## Run all checks (format, vet, test)
	@echo "$(GREEN)[INFO]$(NC) All checks completed successfully"

.PHONY: ci
ci: deps check test-race test-coverage ## Run CI pipeline locally
	@echo "$(GREEN)[INFO]$(NC) CI pipeline completed successfully"

.PHONY: example
example: ## Run example usage
	@echo "$(GREEN)[INFO]$(NC) Running example..."
	@go run -ldflags "$(LDFLAGS)" -c 'package main; import ("fmt"; "$(MODULE_NAME)"); func main() { info := version.Get(); fmt.Printf("Version: %s\nDetails:\n%s\n", info.String(), info.Text()) }' || echo "Run this as a demo in your own application"

# Development helpers
.PHONY: install-tools
install-tools: ## Install development tools
	@echo "$(GREEN)[INFO]$(NC) Installing development tools..."
	@command -v golangci-lint >/dev/null 2>&1 || { \
		echo "$(YELLOW)[WARN]$(NC) golangci-lint not found. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2"; \
	}

.PHONY: doc
doc: ## Generate and serve documentation
	@echo "$(GREEN)[INFO]$(NC) Serving documentation at http://localhost:6060"
	@echo "$(YELLOW)[INFO]$(NC) Press Ctrl+C to stop"
	godoc -http=:6060