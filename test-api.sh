#!/usr/bin/env bash

# gocalc-api Comprehensive API Testing Script
# This script tests all endpoints with various scenarios
# Make sure the API is running on http://localhost:8080 before executing

BASE_URL="${API_URL:-http://localhost:8080}"
FAILED=0
PASSED=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     gocalc-api Comprehensive API Test Suite          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "Testing API at: ${YELLOW}${BASE_URL}${NC}"
echo ""

# Pre-flight checks
if ! command -v curl &> /dev/null; then
    echo -e "${RED}Error: curl is not installed. Please install curl to run these tests.${NC}"
    exit 1
fi

echo -n "Checking API accessibility... "
if ! curl -s -f "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}FAIL${NC}"
    echo -e "${RED}Error: API is not accessible at $BASE_URL${NC}"
    echo -e "${YELLOW}Make sure the API is running with: make run${NC}"
    exit 1
fi
echo -e "${GREEN}OK${NC}"
echo ""

# Function to test an endpoint
test_endpoint() {
    local name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local expected_status="${5:-200}"
    
    echo -n "Testing: $name ... "
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    status_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASS${NC} (HTTP $status_code)"
        ((PASSED++))
        # Delay to respect rate limit (100 req/min = 0.6s/req minimum)
        sleep 0.65
        return 0
    else
        echo -e "${RED}✗ FAIL${NC} (Expected HTTP $expected_status, got $status_code)"
        echo -e "${RED}Response: $body${NC}"
        ((FAILED++))
        # Delay to respect rate limit (100 req/min = 0.6s/req minimum)
        sleep 0.65
        return 1
    fi
}

echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}1. HEALTH & READINESS CHECKS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

test_endpoint "Health Check" "GET" "/health"
test_endpoint "Readiness Check" "GET" "/ready"

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}2. MATH OPERATIONS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

# Addition
test_endpoint "Addition (positive)" "POST" "/api/math/add" '{"a": 10.5, "b": 5.3}'
test_endpoint "Addition (negative)" "POST" "/api/math/add" '{"a": -15.2, "b": 8.7}'
test_endpoint "Addition (zeros)" "POST" "/api/math/add" '{"a": 0, "b": 0}'
test_endpoint "Addition (large numbers)" "POST" "/api/math/add" '{"a": 999999.99, "b": 0.01}'

# Subtraction
test_endpoint "Subtraction (positive)" "POST" "/api/math/subtract" '{"a": 20.0, "b": 8.5}'
test_endpoint "Subtraction (negative result)" "POST" "/api/math/subtract" '{"a": 5.0, "b": 10.0}'
test_endpoint "Subtraction (negatives)" "POST" "/api/math/subtract" '{"a": -10.0, "b": -5.0}'

# Multiplication
test_endpoint "Multiplication (positive)" "POST" "/api/math/multiply" '{"a": 4.5, "b": 3.0}'
test_endpoint "Multiplication (by zero)" "POST" "/api/math/multiply" '{"a": 100.0, "b": 0.0}'
test_endpoint "Multiplication (negatives)" "POST" "/api/math/multiply" '{"a": -5.0, "b": -3.0}'

# Division
test_endpoint "Division (even)" "POST" "/api/math/divide" '{"a": 15.0, "b": 3.0}'
test_endpoint "Division (decimal)" "POST" "/api/math/divide" '{"a": 10.0, "b": 3.0}'
test_endpoint "Division (negative)" "POST" "/api/math/divide" '{"a": -20.0, "b": 4.0}'
test_endpoint "Division by zero (error)" "POST" "/api/math/divide" '{"a": 10.0, "b": 0.0}' 400

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}3. FINANCE CALCULATIONS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

# VAT
test_endpoint "VAT (exclusive)" "POST" "/api/finance/vat" '{"amount": 100.0, "rate": 23.0, "inclusive": false}'
test_endpoint "VAT (inclusive)" "POST" "/api/finance/vat" '{"amount": 123.0, "rate": 23.0, "inclusive": true}'
test_endpoint "VAT (zero rate)" "POST" "/api/finance/vat" '{"amount": 100.0, "rate": 0.0, "inclusive": false}'
test_endpoint "VAT (high rate)" "POST" "/api/finance/vat" '{"amount": 1000.0, "rate": 25.0, "inclusive": false}'

# Compound Interest
test_endpoint "Compound Interest (annual)" "POST" "/api/finance/compound-interest" '{"principal": 1000.0, "rate": 5.0, "time": 10.0, "compound_frequency": 1}'
test_endpoint "Compound Interest (monthly)" "POST" "/api/finance/compound-interest" '{"principal": 5000.0, "rate": 3.5, "time": 5.0, "compound_frequency": 12}'
test_endpoint "Compound Interest (daily)" "POST" "/api/finance/compound-interest" '{"principal": 10000.0, "rate": 2.0, "time": 2.0, "compound_frequency": 365}'

# Loan Payment
test_endpoint "Loan Payment (mortgage)" "POST" "/api/finance/loan-payment" '{"principal": 200000.0, "annual_rate": 4.5, "years": 30, "payments_per_year": 12}'
test_endpoint "Loan Payment (car loan)" "POST" "/api/finance/loan-payment" '{"principal": 25000.0, "annual_rate": 6.0, "years": 5, "payments_per_year": 12}'
test_endpoint "Loan Payment (short term)" "POST" "/api/finance/loan-payment" '{"principal": 5000.0, "annual_rate": 8.0, "years": 2, "payments_per_year": 12}'

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}4. UTILITY CALCULATIONS${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

# BMI
test_endpoint "BMI (metric - normal)" "POST" "/api/utils/bmi" '{"weight": 70.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
test_endpoint "BMI (imperial)" "POST" "/api/utils/bmi" '{"weight": 150.0, "weight_unit": "lb", "height": 5.5, "height_unit": "ft"}'
test_endpoint "BMI (underweight)" "POST" "/api/utils/bmi" '{"weight": 50.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
test_endpoint "BMI (overweight)" "POST" "/api/utils/bmi" '{"weight": 90.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}'
test_endpoint "BMI (obesity)" "POST" "/api/utils/bmi" '{"weight": 100.0, "weight_unit": "kg", "height": 1.65, "height_unit": "m"}'

# Unit Conversion - Temperature
test_endpoint "Convert C to F" "POST" "/api/utils/unit-conversion" '{"value": 25.0, "from_unit": "C", "to_unit": "F", "unit_type": "temperature"}'
test_endpoint "Convert F to C" "POST" "/api/utils/unit-conversion" '{"value": 77.0, "from_unit": "F", "to_unit": "C", "unit_type": "temperature"}'
test_endpoint "Convert C to K" "POST" "/api/utils/unit-conversion" '{"value": 0.0, "from_unit": "C", "to_unit": "K", "unit_type": "temperature"}'

# Unit Conversion - Weight
test_endpoint "Convert kg to lb" "POST" "/api/utils/unit-conversion" '{"value": 70.0, "from_unit": "kg", "to_unit": "lb", "unit_type": "weight"}'
test_endpoint "Convert lb to kg" "POST" "/api/utils/unit-conversion" '{"value": 154.0, "from_unit": "lb", "to_unit": "kg", "unit_type": "weight"}'
test_endpoint "Convert kg to g" "POST" "/api/utils/unit-conversion" '{"value": 1.5, "from_unit": "kg", "to_unit": "g", "unit_type": "weight"}'

# Unit Conversion - Distance
test_endpoint "Convert m to km" "POST" "/api/utils/unit-conversion" '{"value": 1000.0, "from_unit": "m", "to_unit": "km", "unit_type": "distance"}'
test_endpoint "Convert mi to km" "POST" "/api/utils/unit-conversion" '{"value": 10.0, "from_unit": "mi", "to_unit": "km", "unit_type": "distance"}'
test_endpoint "Convert ft to m" "POST" "/api/utils/unit-conversion" '{"value": 100.0, "from_unit": "ft", "to_unit": "m", "unit_type": "distance"}'

# Unit Conversion - Volume
test_endpoint "Convert L to ml" "POST" "/api/utils/unit-conversion" '{"value": 2.0, "from_unit": "L", "to_unit": "ml", "unit_type": "volume"}'
test_endpoint "Convert gal to L" "POST" "/api/utils/unit-conversion" '{"value": 5.0, "from_unit": "gal", "to_unit": "L", "unit_type": "volume"}'

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}5. ERROR HANDLING & VALIDATION${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

# Invalid JSON
test_endpoint "Invalid JSON" "POST" "/api/math/add" '{"a": invalid}' 400

# Invalid method
test_endpoint "Wrong method (GET on POST endpoint)" "GET" "/api/math/add" "" 405

# Invalid values
test_endpoint "Negative VAT rate" "POST" "/api/finance/vat" '{"amount": 100.0, "rate": -5.0, "inclusive": false}' 400
test_endpoint "Negative principal (loan)" "POST" "/api/finance/loan-payment" '{"principal": -1000.0, "annual_rate": 5.0, "years": 10, "payments_per_year": 12}' 400
test_endpoint "Negative annual rate (loan)" "POST" "/api/finance/loan-payment" '{"principal": 10000.0, "annual_rate": -5.0, "years": 10, "payments_per_year": 12}' 400
test_endpoint "Invalid BMI weight" "POST" "/api/utils/bmi" '{"weight": -70.0, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}' 400
test_endpoint "Invalid BMI height" "POST" "/api/utils/bmi" '{"weight": 70.0, "weight_unit": "kg", "height": 0.0, "height_unit": "m"}' 400

# Unsupported units
test_endpoint "Unsupported unit conversion" "POST" "/api/utils/unit-conversion" '{"value": 100.0, "from_unit": "invalid", "to_unit": "kg", "unit_type": "weight"}' 400

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST SUMMARY${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "Total tests: $((PASSED + FAILED))"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║              ALL TESTS PASSED! ✓                      ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════╝${NC}"
    exit 0
else
    echo -e "${RED}╔════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║           SOME TESTS FAILED! ✗                        ║${NC}"
    echo -e "${RED}╚════════════════════════════════════════════════════════╝${NC}"
    exit 1
fi
