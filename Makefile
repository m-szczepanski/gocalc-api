# Makefile for gocalc-api
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary name and output directory
BINARY_NAME=gocalc-api
BINARY_DIR=bin
BINARY_PATH=$(BINARY_DIR)/$(BINARY_NAME)

# Main package path
MAIN_PATH=./cmd/api

# golangci-lint version
GOLANGCI_LINT_VERSION=v1.62.2

.PHONY: help
help: ## Display this help message
	@echo "gocalc-api Makefile commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

.PHONY: build
build: ## Build the application binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built: $(BINARY_PATH)"

.PHONY: run
run: ## Run the application
	@echo "Running gocalc-api..."
	$(GORUN) $(MAIN_PATH)

.PHONY: clean
clean: ## Remove build artifacts and clean cache
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	@echo "Clean complete"

.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-coverage-view
test-coverage-view: test-coverage ## Run tests with coverage and open report in browser
	@echo "Opening coverage report..."
	@if command -v xdg-open >/dev/null 2>&1; then xdg-open coverage.html; \
	elif command -v open >/dev/null 2>&1; then open coverage.html; \
	elif command -v start >/dev/null 2>&1; then start coverage.html; \
	else echo "Please open coverage.html manually."; fi

.PHONY: fmt
fmt: ## Format code with gofmt
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Format complete"

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo "Checking code format..."
	@test -z "$$($(GOFMT) -s -l . | tee /dev/stderr)" || (echo "Files need formatting" && exit 1)
	@echo "Format check passed"

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...
	@echo "Vet complete"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run 'make install-tools'" && exit 1)
	golangci-lint run ./...
	@echo "Lint complete"

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run 'make install-tools'" && exit 1)
	golangci-lint run --fix ./...
	@echo "Lint fix complete"

.PHONY: check
check: fmt-check vet lint test ## Run all checks (format, vet, lint, test)
	@echo "All checks passed!"

.PHONY: tidy
tidy: ## Tidy and verify go.mod
	@echo "Tidying go.mod..."
	$(GOMOD) tidy
	$(GOMOD) verify
	@echo "Tidy complete"

.PHONY: install-tools
install-tools: ## Install development tools (golangci-lint)
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@which golangci-lint > /dev/null && echo "golangci-lint already installed" || \
		(curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION))
	@echo "Tools installed"

.PHONY: setup-hooks
setup-hooks: ## Set up git hooks
	@echo "Setting up git hooks..."
	git config core.hooksPath .githooks
	chmod +x .githooks/pre-commit
	@echo "Git hooks configured"

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies downloaded"

.PHONY: all
all: clean build test ## Clean, build, and test

# Default target
.DEFAULT_GOAL := help
