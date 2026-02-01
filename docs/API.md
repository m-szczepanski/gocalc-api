# API Documentation

Complete documentation for the gocalc-api HTTP endpoints.

## Base URL

```bash
http://localhost:8080
```

## Quick Reference

- **[OpenAPI Specification](openapi.yaml)** - Full OpenAPI 3.0 spec
- **[API Examples](examples.md)** - Working curl examples for every endpoint
- **[Error Reference](errors.md)** - Error codes and troubleshooting guide

## Available Endpoints

### Math Operations

- `POST /api/math/add` - Add two numbers
- `POST /api/math/subtract` - Subtract two numbers
- `POST /api/math/multiply` - Multiply two numbers
- `POST /api/math/divide` - Divide two numbers

### Finance Calculations

- `POST /api/finance/vat` - Calculate VAT (inclusive or exclusive)
- `POST /api/finance/compound-interest` - Calculate compound interest
- `POST /api/finance/loan-payment` - Calculate loan payment amounts

### Utility Calculations

- `POST /api/utils/bmi` - Calculate Body Mass Index
- `POST /api/utils/unit-conversion` - Convert between units

### Health

- `GET /health` - Health check endpoint (accepts any HTTP method)

## Quick Start Examples

### Health Check

```bash
curl http://localhost:8080/health
```

**Response:**

```bash
OK
```

### Add Two Numbers

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10.5, "b": 5.3}'
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

### Calculate VAT

```bash
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.0, "rate": 23.0, "inclusive": false}'
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

### Calculate BMI

```bash
curl -X POST http://localhost:8080/api/utils/bmi \
  -H "Content-Type: application/json" \
  -d '{"weight": 70.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
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

### Convert Units (Temperature)

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 25.0, "from_unit": "C", "to_unit": "F", "unit_type": "temperature"}'
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

## Response Format

All successful responses (except `/health`) follow this structure:

```json
{
  "data": {
    // Endpoint-specific response data
  },
  "request_id": "unique-request-identifier",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

## Error Response Format

All errors follow this structure:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": "Optional detailed information",
  "request_id": "unique-request-identifier",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

See [errors.md](errors.md) for detailed error documentation.

## API Features

### Rate Limiting

- **Limit:** 100 requests per minute per IP address
- **Burst:** 20 requests
- **Headers:** `X-RateLimit-Limit`, `X-RateLimit-Remaining`, `Retry-After`

When rate limit is exceeded, you'll receive a `429 Too Many Requests` response.

### Request Tracing

Every request receives a unique `X-Request-ID` header and includes the `request_id` in the response body. Use this for debugging and log correlation.

### Timeouts

- **Request timeout:** 30 seconds
- **Server read timeout:** 10 seconds
- **Server write timeout:** 10 seconds
- **Idle timeout:** 120 seconds

### Content Type

All POST endpoints require `Content-Type: application/json` header.

## Unit Conversion Support

The `/api/utils/unit-conversion` endpoint supports the following unit types:

| Unit Type | Supported Units |
| ----------- | ---------------- |
| **weight** | kg, g, lb, oz |
| **height** | m, cm, ft, in |
| **temperature** | C, F, K (case-insensitive) |
| **distance** | m, km, mi, ft, yd |
| **volume** | L, ml, l, gal, fl_oz |

## BMI Categories

The BMI endpoint categorizes results as:

- `underweight`: BMI < 18.5
- `normal`: 18.5 ≤ BMI < 25
- `overweight`: 25 ≤ BMI < 30
- `obesity_class_1`: 30 ≤ BMI < 35
- `obesity_class_2`: 35 ≤ BMI < 40
- `obesity_class_3`: BMI ≥ 40

## Common Error Codes

| Code | HTTP Status | Description |
| ------ | ------------- | ------------- |
| `INVALID_INPUT` | 400 | Malformed JSON or invalid request body |
| `VALIDATION_ERROR` | 400 | Field validation failed |
| `DIVISION_BY_ZERO` | 400 | Attempted division by zero |
| `METHOD_NOT_ALLOWED` | 405 | Wrong HTTP method |
| `RATE_LIMIT_EXCEEDED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

For detailed error documentation with examples, see [errors.md](errors.md).

## More Examples

For comprehensive examples including edge cases, validation errors, and all endpoint variations, see [examples.md](examples.md).

## OpenAPI Specification

For the complete API specification including all schemas, validation rules, and detailed descriptions, see [openapi.yaml](openapi.yaml).
