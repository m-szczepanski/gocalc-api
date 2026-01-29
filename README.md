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

## Running the Project

```bash
go run ./cmd/api
```

## The API will start on

```bash
http://localhost:8080
```

## Health check

```bash
curl http://localhost:8080/health
```

## Testing

```bash
go test ./...
```

## API Documentation

See **[docs/API.md](docs/API.md)** for complete API documentation including endpoints, examples, and features.

Quick resources:

- **[API Documentation](docs/API.md)** - Complete endpoint documentation and quick start
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
