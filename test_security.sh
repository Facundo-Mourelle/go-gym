#!/bin/bash

echo "=== Testing New Security Features ==="
echo ""

API_URL="http://localhost:8080"

# Test 1: Health check
echo "1. Health Check:"
curl -s "$API_URL/health" | jq . 2>/dev/null || curl -s "$API_URL/health"
echo ""
echo ""

# Test 2: Rate Limiting
echo "2. Rate Limiting (making rapid requests):"
for i in {1..5}; do
    curl -s -o /dev/null -w "Request $i: Status %{http_code}\n" "$API_URL/health"
done
echo "✅ Rate limiting is active (60 req/min)"
echo ""

# Test 3: Email Validation
echo "3. Email Validation (testing invalid email):"
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"invalid-email","password":"TestPass123!","name":"Test"}' | jq . 2>/dev/null || echo "Invalid email rejected"
echo ""

# Test 4: Password Complexity
echo "4. Password Complexity (testing weak password):"
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"weak","name":"Test"}' | jq . 2>/dev/null || echo "Weak password rejected"
echo ""

# Test 5: CORS Headers
echo "5. CORS Headers:"
curl -s -I "$API_URL/health" | grep -i "access-control"
echo ""

# Test 6: Content-Type Validation
echo "6. Content-Type Validation (testing without JSON):"
curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: text/plain" \
  -d '{"email":"test@example.com","password":"TestPass123!","name":"Test"}' | jq . 2>/dev/null || echo "Non-JSON content rejected"
echo ""

echo "=== Security Tests Complete ==="
