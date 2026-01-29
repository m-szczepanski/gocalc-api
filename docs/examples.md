# API Examples

This document provides practical curl examples for all gocalc-api endpoints.

## Base URL

```txt
http://localhost:8080
```

## Table of Contents

- [Health Check](#health-check)
- [Math Operations](#math-operations)
  - [Addition](#addition)
  - [Subtraction](#subtraction)
  - [Multiplication](#multiplication)
  - [Division](#division)
- [Finance Calculations](#finance-calculations)
  - [VAT Calculation](#vat-calculation)
  - [Compound Interest](#compound-interest)
  - [Loan Payment](#loan-payment)
- [Utility Calculations](#utility-calculations)
  - [BMI Calculator](#bmi-calculator)
  - [Unit Conversion](#unit-conversion)

## Health Check

Check if the API is running and healthy.

```bash
curl http://localhost:8080/health
```

**Response:**

```bash
OK
```

## Math Operations

### Addition

Add two numbers together.

**Success case:**

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{
    "a": 10.5,
    "b": 5.3
  }'
```

**Response:**

```json
{
  "data": {
    "result": 15.8
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**With negative numbers:**

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{
    "a": -15.2,
    "b": 8.7
  }'
```

**Validation error (missing field):**

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{
    "a": 10.5
  }'
```

**Response:**

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": "field 'b' is required",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Subtraction

Subtract one number from another.

```bash
curl -X POST http://localhost:8080/api/math/subtract \
  -H "Content-Type: application/json" \
  -d '{
    "a": 20.0,
    "b": 8.5
  }'
```

**Response:**

```json
{
  "data": {
    "result": 11.5
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Multiplication

Multiply two numbers.

```bash
curl -X POST http://localhost:8080/api/math/multiply \
  -H "Content-Type: application/json" \
  -d '{
    "a": 4.5,
    "b": 3.0
  }'
```

**Response:**

```json
{
  "data": {
    "result": 13.5
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Division

Divide one number by another.

**Success case:**

```bash
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{
    "a": 15.0,
    "b": 3.0
  }'
```

**Response:**

```json
{
  "data": {
    "result": 5
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Division by zero (error):**

```bash
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{
    "a": 10.0,
    "b": 0.0
  }'
```

**Response:**

```json
{
  "code": "DIVISION_BY_ZERO",
  "message": "Division by zero is not allowed",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

## Finance Calculations

### VAT Calculation

Calculate Value Added Tax for a given amount.

**Add VAT to net amount:**

```bash
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.0,
    "rate": 23.0,
    "inclusive": false
  }'
```

**Response:**

```json
{
  "data": {
    "vat_amount": 23,
    "net_amount": 100,
    "gross_amount": 123
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Extract VAT from gross amount:**

```bash
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 123.0,
    "rate": 23.0,
    "inclusive": true
  }'
```

**Response:**

```json
{
  "data": {
    "vat_amount": 23,
    "net_amount": 100,
    "gross_amount": 123
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Validation error (negative rate):**

```bash
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 100.0,
    "rate": -5.0,
    "inclusive": false
  }'
```

**Response:**

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": "field 'rate' must be >= 0",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Compound Interest

Calculate compound interest over time.

**Monthly compounding:**

```bash
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{
    "principal": 1000.0,
    "rate": 5.0,
    "time": 2.0,
    "compound_frequency": 12
  }'
```

**Response:**

```json
{
  "data": {
    "final_amount": 1104.94,
    "interest_earned": 104.94
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Annual compounding:**

```bash
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{
    "principal": 5000.0,
    "rate": 3.5,
    "time": 10.0,
    "compound_frequency": 1
  }'
```

**Response:**

```json
{
  "data": {
    "final_amount": 7052.99,
    "interest_earned": 2052.99
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Loan Payment

Calculate periodic loan payment amounts.

**30-year mortgage:**

```bash
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{
    "principal": 200000.0,
    "annual_rate": 3.5,
    "years": 30.0,
    "payments_per_year": 12
  }'
```

**Response:**

```json
{
  "data": {
    "payment_amount": 898.09,
    "total_payment": 323312.18,
    "total_interest": 123312.18
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Car loan:**

```bash
curl -X POST http://localhost:8080/api/finance/loan-payment \
  -H "Content-Type: application/json" \
  -d '{
    "principal": 25000.0,
    "annual_rate": 4.2,
    "years": 5.0,
    "payments_per_year": 12
  }'
```

**Response:**

```json
{
  "data": {
    "payment_amount": 461.45,
    "total_payment": 27687,
    "total_interest": 2687
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

## Utility Calculations

### BMI Calculator

Calculate Body Mass Index with automatic unit conversion.

**Using metric units:**

```bash
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 70.0,
    "weight_unit": "kg",
    "height": 1.75,
    "height_unit": "m"
  }'
```

**Response:**

```json
{
  "data": {
    "bmi": 22.86,
    "category": "normal"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Using imperial units:**

```bash
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{
    "weight": 154.0,
    "weight_unit": "lb",
    "height": 68.0,
    "height_unit": "in"
  }'
```

**Response:**

```json
{
  "data": {
    "bmi": 23.4,
    "category": "normal"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**BMI Categories:**

- `underweight`: BMI < 18.5
- `normal`: 18.5 ≤ BMI < 25
- `overweight`: 25 ≤ BMI < 30
- `obesity_class_1`: 30 ≤ BMI < 35
- `obesity_class_2`: 35 ≤ BMI < 40
- `obesity_class_3`: BMI ≥ 40

### Unit Conversion

Convert values between different units.

**Temperature: Celsius to Fahrenheit**

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{
    "value": 25.0,
    "from_unit": "C",
    "to_unit": "F",
    "unit_type": "temperature"
  }'
```

**Response:**

```json
{
  "data": {
    "result": 77,
    "from_unit": "C",
    "to_unit": "F",
    "unit_type": "temperature"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Weight: Kilograms to Pounds**

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{
    "value": 70.0,
    "from_unit": "kg",
    "to_unit": "lb",
    "unit_type": "weight"
  }'
```

**Response:**

```json
{
  "data": {
    "result": 154.323584,
    "from_unit": "kg",
    "to_unit": "lb",
    "unit_type": "weight"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Distance: Miles to Kilometers**

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{
    "value": 10.0,
    "from_unit": "mi",
    "to_unit": "km",
    "unit_type": "distance"
  }'
```

**Response:**

```json
{
  "data": {
    "result": 16.09344,
    "from_unit": "mi",
    "to_unit": "km",
    "unit_type": "distance"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Volume: Liters to Gallons**

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{
    "value": 10.0,
    "from_unit": "L",
    "to_unit": "gal",
    "unit_type": "volume"
  }'
```

**Response:**

```json
{
  "data": {
    "result": 2.641721,
    "from_unit": "L",
    "to_unit": "gal",
    "unit_type": "volume"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Height: Feet to Meters**

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{
    "value": 6.0,
    "from_unit": "ft",
    "to_unit": "m",
    "unit_type": "height"
  }'
```

**Response:**

```json
{
  "data": {
    "result": 1.8288,
    "from_unit": "ft",
    "to_unit": "m",
    "unit_type": "height"
  },
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

#### Supported Units

| Unit Type | Supported Units |
|-----------|----------------|
| **weight** | kg, g, lb, oz |
| **height** | m, cm, ft, in |
| **temperature** | C, F, K (case-insensitive) |
| **distance** | m, km, mi, ft, yd |
| **volume** | L, ml, l, gal, fl_oz |

## Common Error Scenarios

### Invalid JSON

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10.5, "b": }'
```

**Response (400):**

```json
{
  "code": "INVALID_INPUT",
  "message": "Invalid request body",
  "details": "invalid character '}' looking for beginning of value",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Wrong HTTP Method

```bash
curl -X GET http://localhost:8080/api/math/add
```

**Response (405):**

```json
{
  "code": "METHOD_NOT_ALLOWED",
  "message": "Method not allowed",
  "details": "Only POST method is allowed",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Rate Limit Exceeded

After exceeding 100 requests per minute:

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 1, "b": 2}'
```

**Response (429):**

```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Rate limit exceeded",
  "details": "Too many requests. Please try again in 60 seconds",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Headers:**

- `X-RateLimit-Limit: 100`
- `X-RateLimit-Remaining: 0`
- `Retry-After: 60`

## Tips

1. **All calculation endpoints use POST method** (except `/health` which accepts any method)
2. **Always set Content-Type header** to `application/json` for POST requests
3. **Check the X-Request-ID** in responses for troubleshooting
4. **Rate limit:** 100 requests per minute per IP address
5. **Request timeout:** 30 seconds
6. **All numeric results** are returned as float64 values
7. **Temperature units** are case-insensitive (C, c, F, f, K, k all work)
