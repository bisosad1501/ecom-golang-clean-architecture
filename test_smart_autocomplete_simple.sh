#!/bin/bash

# Simple Smart Autocomplete Test Script
echo "🚀 Testing Smart Autocomplete API..."

# Test 1: Basic query
echo "📋 Test 1: Basic query for 'iphone'"
curl -s "http://localhost:8080/api/v1/search/autocomplete/smart?q=iphone&limit=5" | jq .

echo ""
echo "📋 Test 2: Query for 'phone'"
curl -s "http://localhost:8080/api/v1/search/autocomplete/smart?q=phone&limit=5" | jq .

echo ""
echo "📋 Test 3: Query for 'samsung'"
curl -s "http://localhost:8080/api/v1/search/autocomplete/smart?q=samsung&limit=5" | jq .

echo ""
echo "📋 Test 4: Query for 'apple'"
curl -s "http://localhost:8080/api/v1/search/autocomplete/smart?q=apple&limit=5" | jq .

echo ""
echo "📋 Test 5: Query for 'elect' (should match electronics)"
curl -s "http://localhost:8080/api/v1/search/autocomplete/smart?q=elect&limit=5" | jq .

echo ""
echo "✅ Tests completed!"
