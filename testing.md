# gocalc-api Testing Guide

Comprehensive testing guide and manual checklist for all API endpoints and functionalities.

## Prerequisites

1. **Start the API:**

   ```bash
   make run
   # OR with Docker:
   make docker-run
   ```

2. **Verify API is running:**

   ```bash
   curl http://localhost:8080/health
   # Expected: OK or JSON with health status
   ```

## Automated Testing

### Run All Tests

```bash
# Run unit and integration tests
make test

# Run full quality checks (format, lint, vet, test)
make check

# Run automated API tests (requires API to be running)
./test-api.sh
```

## Manual API Testing

Base URL: `http://localhost:8080`

### 1. Health & Monitoring Endpoints

#### Health Check

```bash
curl http://localhost:8080/health
```

**Expected:** JSON with status, uptime, version, memory and goroutine checks

#### Readiness Check

```bash
curl http://localhost:8080/ready
```

**Expected:** JSON with ready status

### 2. Math Operations

#### Addition

```bash
# Basic addition
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10.5, "b": 5.3}'
# Expected: {"data": {"result": 15.8}, ...}

# Negative numbers
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": -15.2, "b": 8.7}'
# Expected: {"data": {"result": -6.5}, ...}

# Large numbers
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 999999.99, "b": 0.01}'
# Expected: {"data": {"result": 1000000}, ...}
```

#### Subtraction

```bash
# Basic subtraction
curl -X POST http://localhost:8080/api/math/subtract \
  -H "Content-Type: application/json" \
  -d '{"a": 20.0, "b": 8.5}'
# Expected: {"data": {"result": 11.5}, ...}

# Negative result
curl -X POST http://localhost:8080/api/math/subtract \
  -H "Content-Type: application/json" \
  -d '{"a": 5.0, "b": 10.0}'
# Expected: {"data": {"result": -5}, ...}
```

#### Multiplication

```bash
# Basic multiplication
curl -X POST http://localhost:8080/api/math/multiply \
  -H "Content-Type: application/json" \
  -d '{"a": 4.5, "b": 3.0}'
# Expected: {"data": {"result": 13.5}, ...}

# Multiply by zero
curl -X POST http://localhost:8080/api/math/multiply \
  -H "Content-Type: application/json" \
  -d '{"a": 100.0, "b": 0.0}'
# Expected: {"data": {"result": 0}, ...}

# Negative numbers
curl -X POST http://localhost:8080/api/math/multiply \
  -H "Content-Type: application/json" \
  -d '{"a": -5.0, "b": -3.0}'
# Expected: {"data": {"result": 15}, ...}
```

#### Division

```bash
# Basic division
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 15.0, "b": 3.0}'
# Expected: {"data": {"result": 5}, ...}

# Division with decimals
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10.0, "b": 3.0}'
# Expected: {"data": {"result": 3.333...}, ...}

# Division by zero (ERROR)
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10.0, "b": 0.0}'
# Expected: {"code": "DIVISION_BY_ZERO", ...}
```

### 3. Finance Calculations

#### VAT Calculation

```bash
# VAT exclusive (add VAT to net amount)
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.0, "rate": 23.0, "inclusive": false}'
# Expected: vat_amount: 23, net_amount: 100, gross_amount: 123

# VAT inclusive (extract VAT from gross amount)
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 123.0, "rate": 23.0, "inclusive": true}'
# Expected: vat_amount: 23, net_amount: 100, gross_amount: 123

# Zero VAT rate
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.0, "rate": 0.0, "inclusive": false}'
# Expected: vat_amount: 0, net_amount: 100, gross_amount: 100

# High VAT rate (25%)
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 1000.0, "rate": 25.0, "inclusive": false}'
# Expected: vat_amount: 250, net_amount: 1000, gross_amount: 1250
```

#### Compound Interest

```bash
# Annual compounding
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{"principal": 1000.0, "rate": 5.0, "time": 10.0, "compound_frequency": 1}'
# Expected: final_amount, interest_earned

# Monthly compounding
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{"principal": 5000.0, "rate": 3.5, "time": 5.0, "compound_frequency": 12}'
# Expected: final_amount, interest_earned

# Daily compounding
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{"principal": 10000.0, "rate": 2.0, "time": 2.0, "compound_frequency": 365}'
# Expected: final_amount, interest_earned
```

#### Loan Payment

```bash
# Mortgage (30 years)
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{"principal": 200000.0, "annual_rate": 4.5, "years": 30, "payments_per_year": 12}'
# Expected: monthly_payment, total_payment, total_interest

# Car loan (5 years)
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{"principal": 25000.0, "annual_rate": 6.0, "years": 5, "payments_per_year": 12}'
# Expected: monthly_payment, total_payment, total_interest

# Short-term loan (2 years)
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{"principal": 5000.0, "annual_rate": 8.0, "years": 2, "payments_per_year": 12}'
# Expected: monthly_payment, total_payment, total_interest
```

### 4. Utility Calculations

#### BMI Calculator

```bash
# Normal BMI (metric)
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 70.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
# Expected: bmi: ~22.86, category: "normal"

# Imperial units
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 150.0, "weight_unit": "lb", "height": 5.5, "height_unit": "ft"}'
# Expected: bmi and category

# Underweight
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 50.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
# Expected: category: "underweight"

# Overweight
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 90.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
# Expected: category: "overweight"

# Obesity
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 100.0, "weight_unit": "kg", "height": 1.65, "height_unit": "m"}'
# Expected: category: obesity class
```

#### Unit Conversion

**Temperature:**

```bash
# Celsius to Fahrenheit
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 25.0, "from_unit": "C", "to_unit": "F", "unit_type": "temperature"}'
# Expected: result: 77

# Fahrenheit to Celsius
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 77.0, "from_unit": "F", "to_unit": "C", "unit_type": "temperature"}'
# Expected: result: 25

# Celsius to Kelvin
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 0.0, "from_unit": "C", "to_unit": "K", "unit_type": "temperature"}'
# Expected: result: 273.15
```

**Weight:**

```bash
# Kilograms to pounds
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 70.0, "from_unit": "kg", "to_unit": "lb", "unit_type": "weight"}'
# Expected: result: ~154.32

# Pounds to kilograms
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 154.0, "from_unit": "lb", "to_unit": "kg", "unit_type": "weight"}'
# Expected: result: ~69.85

# Kilograms to grams
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 1.5, "from_unit": "kg", "to_unit": "g", "unit_type": "weight"}'
# Expected: result: 1500
```

**Distance:**

```bash
# Meters to kilometers
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 1000.0, "from_unit": "m", "to_unit": "km", "unit_type": "distance"}'
# Expected: result: 1

# Miles to kilometers
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 10.0, "from_unit": "mi", "to_unit": "km", "unit_type": "distance"}'
# Expected: result: ~16.09

# Feet to meters
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 100.0, "from_unit": "ft", "to_unit": "m", "unit_type": "distance"}'
# Expected: result: ~30.48
```

**Volume:**

```bash
# Liters to milliliters
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 2.0, "from_unit": "L", "to_unit": "ml", "unit_type": "volume"}'
# Expected: result: 2000

# Gallons to liters
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 5.0, "from_unit": "gal", "to_unit": "L", "unit_type": "volume"}'
# Expected: result: ~18.93
```

### 5. Error Handling Tests

#### Invalid JSON

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": invalid}'
# Expected: HTTP 400, INVALID_INPUT error
```

#### Missing Required Fields

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10.0}'
# Expected: HTTP 400, VALIDATION_ERROR
```

#### Wrong HTTP Method

```bash
curl -X GET http://localhost:8080/api/math/add
# Expected: HTTP 405, METHOD_NOT_ALLOWED
```

#### Invalid Values

```bash
# Negative VAT rate
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.0, "rate": -5.0, "inclusive": false}'
# Expected: HTTP 400, VALIDATION_ERROR

# Zero principal for loan
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{"principal": 0.0, "annual_rate": 5.0, "years": 10}'
# Expected: HTTP 400, VALIDATION_ERROR

# Negative BMI weight
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": -70.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
# Expected: HTTP 400, VALIDATION_ERROR
```

#### Unsupported Unit Types

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 100.0, "from_unit": "invalid", "to_unit": "kg", "unit_type": "weight"}'
# Expected: HTTP 400, VALIDATION_ERROR
```

### 6. Configuration & Environment Testing

Test environment variable configuration (requires restart):

```bash
# Custom port
PORT=3000 make run

# Custom rate limiting
RATE_LIMIT_RPM=200 RATE_LIMIT_BURST=50 make run

# Custom timeouts
REQUEST_TIMEOUT=45s READ_TIMEOUT=15s make run
```

## Test Coverage

Check test coverage:

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
make test-coverage-view

# Or manually:
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Docker Testing

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run-detached

# Test health endpoint
curl http://localhost:8080/health

# View logs
make docker-logs

# Stop container
make docker-stop

# Test with custom env vars
docker run --rm -d -p 8080:8080 \
  -e PORT=8080 \
  -e RATE_LIMIT_RPM=200 \
  --name gocalc-api \
  gocalc-api:latest
```

## Summary Checklist

- [ ] All unit tests pass (`make test`)
- [ ] Code quality checks pass (`make check`)
- [ ] Automated API tests pass (`./test-api.sh`)
- [ ] Health endpoints respond correctly
- [ ] All math operations work (add, subtract, multiply, divide)
- [ ] Finance calculations work (VAT, compound interest, loan payment)
- [ ] BMI calculator works with different units and categories
- [ ] Unit conversions work for all types (temp, weight, distance, volume)
- [ ] Error handling works correctly (invalid JSON, missing fields, wrong methods)
- [ ] Validation catches invalid values
- [ ] Configuration via environment variables works
- [ ] Docker build and run work correctly
- [ ] Rate limiting works (test with repeated requests)
- [ ] Request IDs are included in responses
- [ ] Timestamps are in RFC3339 format

## Troubleshooting

**API not responding:**

```bash
# Check if API is running
curl http://localhost:8080/health

# Check process
ps aux | grep gocalc

# Check logs if running in Docker
make docker-logs
```

**Tests failing:**

```bash
# Clean and rebuild
make clean
make build

# Check for errors
make vet
make lint
```

**Port already in use:**

```bash
# Find process using port 8080
lsof -i :8080

# Or use different port
PORT=3000 make run
```
