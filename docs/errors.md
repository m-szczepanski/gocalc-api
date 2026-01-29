# Error Response Documentation

This document describes all error codes, their meanings, HTTP status codes, and troubleshooting guidance for the gocalc-api.

## Error Response Format

All errors follow a consistent JSON structure:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable error message",
  "details": "Optional detailed error information",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `code` | string | Yes | Machine-readable error code |
| `message` | string | Yes | Human-readable error description |
| `details` | string | No | Additional context about the error |
| `request_id` | string | Yes | Unique request identifier for tracing |
| `timestamp` | string | Yes | ISO 8601 timestamp when error occurred |

---

## Error Codes

### INVALID_INPUT

**HTTP Status:** `400 Bad Request`

**Description:** The request body is malformed or cannot be parsed as valid JSON.

**Common Causes:**

- Malformed JSON syntax
- Missing required fields
- Incorrect data types
- Invalid characters in JSON

**Example:**

```json
{
  "code": "INVALID_INPUT",
  "message": "Invalid request body",
  "details": "invalid character '}' looking for beginning of value",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Troubleshooting:**

1. Validate JSON syntax using a JSON validator
2. Ensure all required fields are present
3. Check that field values match expected data types
4. Verify Content-Type header is set to `application/json`

**Example Trigger:**

```bash
curl -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10.5, "b": }'  # Missing value for 'b'
```

### VALIDATION_ERROR

**HTTP Status:** `400 Bad Request`

**Description:** The request body is valid JSON but contains invalid field values that fail validation rules.

**Common Causes:**

- Negative values where only positive numbers are allowed
- Zero or negative values for fields requiring positive numbers
- Invalid enum values (e.g., unsupported unit types)
- Missing required string fields
- NaN or Infinity values

**Example:**

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Validation failed",
  "details": "field 'rate' must be >= 0",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Validation Rules by Endpoint:**

#### Math Operations (`/api/math/*`)

- `a` and `b` must be valid numbers (not NaN, not Inf)

#### VAT (`/api/finance/vat`)

- `amount` must be ≥ 0
- `rate` must be ≥ 0
- `inclusive` must be boolean

#### Compound Interest (`/api/finance/compound-interest`)

- `principal` must be ≥ 0
- `rate` must be ≥ 0
- `time` must be ≥ 0
- `compound_frequency` must be > 0

#### Loan Payment (`/api/finance/loan-payment`)

- `principal` must be ≥ 0
- `annual_rate` must be ≥ 0
- `years` must be > 0
- `payments_per_year` must be > 0

#### BMI (`/api/utils/bmi`)

- `weight` must be > 0
- `height` must be > 0
- `weight_unit` must be non-empty and valid (kg, g, lb, oz)
- `height_unit` must be non-empty and valid (m, cm, ft, in)

#### Unit Conversion (`/api/utils/unit-conversion`)

- `value` can be any valid number (including negative)
- `from_unit` must be non-empty
- `to_unit` must be non-empty
- `unit_type` must be one of: weight, height, temperature, distance, volume

**Troubleshooting:**

1. Check the `details` field for specific field that failed validation
2. Review validation rules for the endpoint you're calling
3. Ensure numeric values are within valid ranges
4. Verify enum values match exactly (case-sensitive except for temperature units)

**Example Triggers:**

Negative rate:

```bash
curl -X POST http://localhost:8080/api/finance/vat \
  -H "Content-Type: application/json" \
  -d '{"amount": 100.0, "rate": -5.0, "inclusive": false}'
```

Invalid unit type:

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 10, "from_unit": "kg", "to_unit": "lb", "unit_type": "mass"}'
```

### DIVISION_BY_ZERO

**HTTP Status:** `400 Bad Request`

**Description:** Attempted to divide by zero, which is mathematically undefined.

**Common Causes:**

- Setting `b` to 0 in division endpoint
- Edge cases in financial calculations with zero values

**Example:**

```json
{
  "code": "DIVISION_BY_ZERO",
  "message": "Division by zero is not allowed",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Troubleshooting:**

1. Check that divisor (`b` field) is not zero
2. Ensure calculations don't result in division by zero
3. Use a non-zero value for the divisor

**Example Trigger:**

```bash
curl -X POST http://localhost:8080/api/math/divide \
  -H "Content-Type: application/json" \
  -d '{"a": 10.0, "b": 0.0}'
```

### METHOD_NOT_ALLOWED

**HTTP Status:** `405 Method Not Allowed`

**Description:** The HTTP method used is not supported for this endpoint.

**Common Causes:**

- Using GET instead of POST for calculation endpoints
- Using PUT, PATCH, or DELETE on any endpoint
- Wrong HTTP verb for the endpoint

**Example:**

```json
{
  "code": "METHOD_NOT_ALLOWED",
  "message": "Method not allowed",
  "details": "Only POST method is allowed",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Allowed Methods:**

- `/health` - Any method (GET, POST, PUT, etc.)
- All other endpoints - POST only

**Troubleshooting:**

1. Use POST method for all calculation endpoints
2. Check API documentation for correct HTTP method
3. Verify curl command uses `-X POST`

**Example Trigger:**

```bash
curl -X GET http://localhost:8080/api/math/add
```

### RATE_LIMIT_EXCEEDED

**HTTP Status:** `429 Too Many Requests`

**Description:** The client has exceeded the rate limit for API requests.

**Rate Limit Configuration:**

- **Limit:** 100 requests per 60 seconds
- **Burst:** 20 requests
- **Scope:** Per IP address

**Example:**

```json
{
  "code": "RATE_LIMIT_EXCEEDED",
  "message": "Rate limit exceeded",
  "details": "Too many requests. Please try again in 60 seconds",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Response Headers:**

```txt
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
Retry-After: 60
```

**Troubleshooting:**

1. Implement exponential backoff in your client

2. Wait for the time specified in `Retry-After` header (60 seconds)
3. Reduce request frequency to stay within limits
4. Consider batching operations if possible
5. Cache results to minimize redundant requests

**Rate Limit Best Practices:**

- Monitor `X-RateLimit-Remaining` header
- Implement client-side rate limiting
- Add delays between requests
- Handle 429 responses gracefully with retry logic

**Example Trigger:**

```bash
# Send more than 100 requests within 60 seconds
for i in {1..101}; do
  curl -X POST http://localhost:8080/api/math/add \
    -H "Content-Type: application/json" \
    -d '{"a": 1, "b": 2}'
done
```

### INTERNAL_ERROR

**HTTP Status:** `500 Internal Server Error`

**Description:** An unexpected error occurred on the server side.

**Common Causes:**

- Unexpected runtime panic (recovered by middleware)
- Request timeout (>30 seconds)
- Calculation errors not caught by validation
- Server resource exhaustion
- Unexpected nil pointer dereference

**Example:**

```json
{
  "code": "INTERNAL_ERROR",
  "message": "Internal server error",
  "details": "An unexpected error occurred",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Timeout Scenario:**

```json
{
  "code": "INTERNAL_ERROR",
  "message": "Internal server error",
  "details": "request timeout",
  "request_id": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "timestamp": "2026-01-29T10:30:00Z"
}
```

**Troubleshooting:**

1. Note the `request_id` for server log correlation
2. Check if request is timing out (>30 seconds)
3. Verify server is running and healthy (`/health` endpoint)
4. Reduce complexity of calculations if possible
5. Retry the request after a brief delay
6. Report persistent errors with `request_id` to maintainers

**Prevention:**

- Keep calculation inputs within reasonable ranges
- Avoid extremely large numbers that may cause overflow
- Monitor server logs for panic messages
- Ensure server has adequate resources

## HTTP Status Code Summary

| Status Code | Error Code(s) | Description |
|-------------|---------------|-------------|
| 400 | INVALID_INPUT, VALIDATION_ERROR, DIVISION_BY_ZERO | Bad Request - Client error |
| 405 | METHOD_NOT_ALLOWED | Method Not Allowed |
| 429 | RATE_LIMIT_EXCEEDED | Too Many Requests |
| 500 | INTERNAL_ERROR | Internal Server Error |

## Error Handling Best Practices

### 1. Check HTTP Status Code First

```bash
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
  -X POST http://localhost:8080/api/math/add \
  -H "Content-Type: application/json" \
  -d '{"a": 1, "b": 2}')

if [ "$HTTP_STATUS" -ne 200 ]; then
  echo "Request failed with status $HTTP_STATUS"
fi
```

### 2. Parse Error Response

```javascript
try {
  const response = await fetch('http://localhost:8080/api/math/add', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ a: 10, b: 0 })
  });
  
  if (!response.ok) {
    const error = await response.json();
    console.error(`Error ${error.code}: ${error.message}`);
    console.error(`Details: ${error.details}`);
    console.error(`Request ID: ${error.request_id}`);
  }
} catch (err) {
  console.error('Network error:', err);
}
```

### 3. Implement Retry Logic for Rate Limits

```python
import requests
import time

def call_api_with_retry(url, data, max_retries=3):
    for attempt in range(max_retries):
        response = requests.post(url, json=data)
        
        if response.status_code == 429:
            retry_after = int(response.headers.get('Retry-After', 60))
            print(f"Rate limited. Waiting {retry_after} seconds...")
            time.sleep(retry_after)
            continue
            
        return response
    
    raise Exception("Max retries exceeded")
```

### 4. Log Request IDs

Always log the `request_id` from error responses for troubleshooting:

```go
if resp.StatusCode != http.StatusOK {
    var errResp ErrorResponse
    json.NewDecoder(resp.Body).Decode(&errResp)
    log.Printf("Error: %s (Request ID: %s)", errResp.Code, errResp.RequestID)
}
```

## Common Error Scenarios

### Scenario 1: Forgot Content-Type Header

```bash
curl -X POST http://localhost:8080/api/math/add \
  -d '{"a": 1, "b": 2}'
```

**Result:** May receive `INVALID_INPUT` error

**Solution:** Always include `-H "Content-Type: application/json"`

### Scenario 2: Using GET for Calculations

```bash
curl http://localhost:8080/api/math/add?a=1&b=2
```

**Result:** `METHOD_NOT_ALLOWED` error

**Solution:** Use POST with JSON body

### Scenario 3: Invalid Unit in Conversion

```bash
curl -X POST http://localhost:8080/api/utils/unit-conversion \
  -H "Content-Type: application/json" \
  -d '{"value": 10, "from_unit": "xyz", "to_unit": "kg", "unit_type": "weight"}'
```

**Result:** `VALIDATION_ERROR` or `INTERNAL_ERROR`

**Solution:** Use only supported units from the documentation

### Scenario 4: Negative Values Where Not Allowed

```bash
curl -X POST http://localhost:8080/api/finance/compound-interest \
  -H "Content-Type: application/json" \
  -d '{"principal": -1000, "rate": 5, "time": 2, "compound_frequency": 12}'
```

**Result:** `VALIDATION_ERROR` - principal must be >= 0

**Solution:** Use non-negative values for financial calculations

## Getting Help

If you encounter an error that doesn't match these descriptions:

1. **Check the `request_id`** - This helps trace the request in server logs
2. **Review endpoint documentation** - Verify request format and validation rules
3. **Check API examples** - See [examples.md](./examples.md) for working requests
4. **Verify OpenAPI spec** - See [openapi.yaml](./openapi.yaml) for complete API contract
5. **Test with `/health` endpoint** - Ensure the server is running
6. **Check server logs** - Look for errors corresponding to your `request_id`

For persistent issues, include:

- Request ID
- Full error response
- Request body (sanitized)
- Expected vs actual behavior
