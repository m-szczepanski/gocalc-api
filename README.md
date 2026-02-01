# gocalc-api

It's a **stateless HTTP API** written in **Go** that exposes endpoints for performing various calculations such as mathematical, financial, and utility-based computations.

This project was created **purely as a learning exercise** to explore Go as a backend language, understand its ecosystem, and practice building clean, testable APIs using Go’s standard library.

## Project Goals

- Learn Go fundamentals through a practical backend project
- Understand idiomatic Go project structure
- Build a stateless HTTP API without heavy frameworks
- Practice clean architecture and separation of concerns
- Write unit and HTTP tests using Go tooling

## What “Stateless” Means Here

- No database
- No sessions
- No stored user data
- Every request contains all the information required to perform a calculation

This makes the API:

- Easy to scale
- Simple to reason about
- Perfect for learning backend fundamentals

## Tech Stack

- **Language:** Go
- **HTTP:** `net/http`
- **Routing:** Standard library (with possible future extensions)
- **JSON:** `encoding/json`
- **Testing:** `testing`, `httptest`
- **Linting:** `golangci-lint`
- **Build Automation:** `Makefile`
- **Code Quality:** Pre-commit hooks

## Planned Features

### Mathematical Calculations

- Addition, subtraction, multiplication, division
- Factorials, percentages

### Financial Calculations

- VAT calculation
- Compound interest
- Loan payment estimation

### Utility Calculations

- BMI
- Unit conversions

## Project Structure

```text
gocalc-api/
├── cmd/
│   └── api/
│       └── main.go        # Application entry point
├── internal/
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Business logic
│   ├── models/            # Request/response models
│   └── validation/        # Input validation
├── pkg/
│   └── calculations/      # Pure calculation functions
├── tests/
├── go.mod
└── README.md
```

## Quick Start

### First-time setup

```bash
# Install development tools (golangci-lint)
make install-tools

# Set up git hooks for code quality checks
make setup-hooks
```

### Running the Project

#### Local Development

```bash
# Run with default configuration
make run

# Run with custom port
PORT=3000 make run

# Run with custom configuration
PORT=3000 RATE_LIMIT_RPM=200 REQUEST_TIMEOUT=45s make run
```

#### Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run

# Run in background
make docker-run-detached

# View logs
make docker-logs

# Stop container
make docker-stop
```

The API will start on `http://localhost:8080` (or your configured port).

### Configuration

The application is configured via environment variables with sensible defaults:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `8080` |
| `READ_TIMEOUT` | HTTP read timeout | `10s` |
| `WRITE_TIMEOUT` | HTTP write timeout | `10s` |
| `IDLE_TIMEOUT` | HTTP idle timeout | `120s` |
| `SHUTDOWN_TIMEOUT` | Graceful shutdown timeout | `15s` |
| `REQUEST_TIMEOUT` | Request processing timeout | `30s` |
| `RATE_LIMIT_RPM` | Rate limit (requests/min) | `100.0` |
| `RATE_LIMIT_BURST` | Rate limit burst size | `20` |

See **[docs/deployment.md](docs/deployment.md)** for complete deployment guide.

### Health check

```bash
# Health endpoint (detailed status)
curl http://localhost:8080/health

# Readiness endpoint (for load balancers)
curl http://localhost:8080/ready
```

### Testing

```bash
make test
```

## Development Workflow

The project includes a Makefile for common development tasks:

```bash
make help          # Show all available commands
make build         # Build the binary to bin/gocalc-api
make run           # Run the application
make test          # Run all tests
make test-coverage # Generate test coverage report
make fmt           # Format code with gofmt
make lint          # Run golangci-lint
make vet           # Run go vet
make check         # Run all checks (fmt, vet, lint, test)
make clean         # Remove build artifacts
make tidy          # Tidy go.mod and go.sum

# Docker commands
make docker-build          # Build Docker image
make docker-run            # Run Docker container
make docker-run-detached   # Run in background
make docker-stop           # Stop container
make docker-logs           # View logs
make docker-clean          # Remove image
```

### Pre-commit Hooks

The project uses git hooks to ensure code quality before commits:

1. **Setup** (first time only):

   ```bash
   make setup-hooks
   ```

2. **What runs automatically** on each commit:
   - Code formatting check (`gofmt`)
   - Go vet analysis
   - golangci-lint checks
   - go.mod tidiness verification

3. **To bypass** hooks (not recommended):

   ```bash
   git commit --no-verify
   ```

### Code Quality

Before pushing code, run all quality checks:

```bash
make check
```

This runs formatting checks, linting, vetting, and tests in one command.

## API Documentation

See **[docs/API.md](docs/API.md)** for complete API documentation including endpoints, examples, and features.

Quick resources:

- **[API Documentation](docs/API.md)** - Complete endpoint documentation and quick start
- **[Deployment Guide](docs/deployment.md)** - Docker, Kubernetes, and cloud deployment
- **[OpenAPI Spec](docs/openapi.yaml)** - Full OpenAPI 3.0 specification
- **[Examples](docs/examples.md)** - curl examples for every endpoint
- **[Error Reference](docs/errors.md)** - Error codes and troubleshooting

## Why This Project Exists

gocalc-api is not meant to be a production-ready system.
It exists to:

- Learn Go by doing
- Experiment with API design
- Make mistakes and fix them
- Serve as a reference for future Go projects
