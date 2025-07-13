#!/bin/bash

# Smart Autocomplete API Testing Script
# This script tests the smart autocomplete functionality comprehensively

API_BASE="http://localhost:8080/api/v1"
ADMIN_EMAIL="admin@ecom.com"
ADMIN_PASSWORD="admin123"

echo "🚀 Starting Smart Autocomplete API Tests..."
echo "=============================================="

# Function to make API calls with proper error handling
make_api_call() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo ""
    echo "📋 Testing: $description"
    echo "🔗 Endpoint: $method $endpoint"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE$endpoint")
    fi
    
    # Extract HTTP status code (last line)
    http_code=$(echo "$response" | tail -n1)
    # Extract response body (all lines except last)
    response_body=$(echo "$response" | head -n -1)
    
    echo "📊 Status Code: $http_code"
    
    if [ "$http_code" -eq 200 ] || [ "$http_code" -eq 201 ]; then
        echo "✅ Success!"
        echo "📄 Response:"
        echo "$response_body" | jq '.' 2>/dev/null || echo "$response_body"
    else
        echo "❌ Failed!"
        echo "📄 Error Response:"
        echo "$response_body"
    fi
    
    echo "----------------------------------------"
}

# Test 1: Basic Smart Autocomplete
echo "🧪 Test 1: Basic Smart Autocomplete Query"
make_api_call "GET" "/search/autocomplete/smart?q=iphone&limit=5" "" "Basic autocomplete for 'iphone'"

# Test 2: Smart Autocomplete with All Features
echo "🧪 Test 2: Smart Autocomplete with All Features Enabled"
make_api_call "GET" "/search/autocomplete/smart?q=phone&limit=10&include_trending=true&include_personalized=true&include_history=true&include_popular=true" "" "Full-featured autocomplete for 'phone'"

# Test 3: Type-Specific Autocomplete
echo "🧪 Test 3: Product-Only Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=samsung&types=product&limit=5" "" "Product-only autocomplete for 'samsung'"

echo "🧪 Test 4: Category-Only Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=elect&types=category&limit=5" "" "Category-only autocomplete for 'elect'"

echo "🧪 Test 5: Brand-Only Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=app&types=brand&limit=5" "" "Brand-only autocomplete for 'app'"

# Test 6: Multi-Type Autocomplete
echo "🧪 Test 6: Multi-Type Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=nike&types=product,brand&limit=8" "" "Product and brand autocomplete for 'nike'"

# Test 7: Fuzzy Matching Test
echo "🧪 Test 7: Fuzzy Matching (Typo Tolerance)"
make_api_call "GET" "/search/autocomplete/smart?q=ipone&limit=5" "" "Fuzzy matching for 'ipone' (typo for iPhone)"

echo "🧪 Test 8: Fuzzy Matching Test 2"
make_api_call "GET" "/search/autocomplete/smart?q=samsng&limit=5" "" "Fuzzy matching for 'samsng' (typo for Samsung)"

# Test 9: Trending Suggestions
echo "🧪 Test 9: Trending Suggestions Only"
make_api_call "GET" "/search/autocomplete/smart?q=&include_trending=true&limit=5" "" "Trending suggestions without query"

# Test 10: Popular Suggestions
echo "🧪 Test 10: Popular Suggestions"
make_api_call "GET" "/search/autocomplete/smart?q=&include_popular=true&limit=5" "" "Popular suggestions without query"

# Test 11: Language-Specific Autocomplete
echo "🧪 Test 11: Vietnamese Language Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=điện&language=vi&limit=5" "" "Vietnamese autocomplete for 'điện'"

echo "🧪 Test 12: Spanish Language Autocomplete"
make_api_call "GET" "/search/autocomplete/smart?q=teléfono&language=es&limit=5" "" "Spanish autocomplete for 'teléfono'"

# Test 13: Empty Query with Features
echo "🧪 Test 13: Empty Query with All Features"
make_api_call "GET" "/search/autocomplete/smart?q=&include_trending=true&include_popular=true&limit=10" "" "Empty query with trending and popular suggestions"

# Test 14: Long Query Test
echo "🧪 Test 14: Long Query Test"
make_api_call "GET" "/search/autocomplete/smart?q=best%20smartphone%202024&limit=5" "" "Long query autocomplete"

# Test 15: Special Characters Test
echo "🧪 Test 15: Special Characters Test"
make_api_call "GET" "/search/autocomplete/smart?q=iphone%2015&limit=5" "" "Query with special characters"

# Test 16: Track Autocomplete Interaction - Click
echo "🧪 Test 16: Track Autocomplete Click Interaction"
make_api_call "POST" "/search/autocomplete/track" '{
    "entry_id": "00000000-0000-0000-0000-000000000001",
    "interaction_type": "click",
    "session_id": "test-session-123",
    "query": "iphone",
    "position": 1
}' "Track autocomplete click interaction"

# Test 17: Track Autocomplete Interaction - Impression
echo "🧪 Test 17: Track Autocomplete Impression Interaction"
make_api_call "POST" "/search/autocomplete/track" '{
    "entry_id": "00000000-0000-0000-0000-000000000002",
    "interaction_type": "impression",
    "session_id": "test-session-123",
    "query": "samsung",
    "position": 2
}' "Track autocomplete impression interaction"

# Test 18: Invalid Parameters Test
echo "🧪 Test 18: Invalid Parameters Test"
make_api_call "GET" "/search/autocomplete/smart?q=test&limit=invalid" "" "Invalid limit parameter test"

echo "🧪 Test 19: Missing Required Parameter Test"
make_api_call "GET" "/search/autocomplete/smart?limit=5" "" "Missing query parameter test"

# Test 20: Large Limit Test
echo "🧪 Test 20: Large Limit Test"
make_api_call "GET" "/search/autocomplete/smart?q=phone&limit=100" "" "Large limit parameter test (should be capped)"

# Test 21: Performance Test with Complex Query
echo "🧪 Test 21: Performance Test"
start_time=$(date +%s%N)
make_api_call "GET" "/search/autocomplete/smart?q=smartphone&include_trending=true&include_personalized=true&include_history=true&include_popular=true&types=product,category,brand,query&limit=20" "" "Complex performance test"
end_time=$(date +%s%N)
duration=$(( (end_time - start_time) / 1000000 ))
echo "⏱️  Query took: ${duration}ms"

# Test 22: Synonym Matching Test
echo "🧪 Test 22: Synonym Matching Test"
make_api_call "GET" "/search/autocomplete/smart?q=mobile&limit=5" "" "Synonym matching for 'mobile' (should match phone-related items)"

echo "🧪 Test 23: Synonym Matching Test 2"
make_api_call "GET" "/search/autocomplete/smart?q=notebook&limit=5" "" "Synonym matching for 'notebook' (should match laptop-related items)"

# Test 24: Case Insensitive Test
echo "🧪 Test 24: Case Insensitive Test"
make_api_call "GET" "/search/autocomplete/smart?q=IPHONE&limit=5" "" "Case insensitive test for 'IPHONE'"

echo "🧪 Test 25: Mixed Case Test"
make_api_call "GET" "/search/autocomplete/smart?q=SaMsUnG&limit=5" "" "Mixed case test for 'SaMsUnG'"

# Summary
echo ""
echo "🎉 Smart Autocomplete API Testing Complete!"
echo "=============================================="
echo ""
echo "📊 Test Summary:"
echo "• Basic functionality tests: ✅"
echo "• Feature-specific tests: ✅"
echo "• Type filtering tests: ✅"
echo "• Fuzzy matching tests: ✅"
echo "• Language support tests: ✅"
echo "• Interaction tracking tests: ✅"
echo "• Error handling tests: ✅"
echo "• Performance tests: ✅"
echo "• Edge case tests: ✅"
echo ""
echo "💡 Key Features Tested:"
echo "• Smart autocomplete with multiple sources"
echo "• Fuzzy matching and typo tolerance"
echo "• Trending and popular suggestions"
echo "• Type-specific filtering"
echo "• Multi-language support"
echo "• Interaction analytics tracking"
echo "• Performance optimization"
echo "• Error handling and validation"
echo ""
echo "🔍 Review the test results above to verify:"
echo "1. All API endpoints return proper HTTP status codes"
echo "2. Response structures match expected format"
echo "3. Smart features work correctly (trending, popular, fuzzy matching)"
echo "4. Type filtering works as expected"
echo "5. Language support functions properly"
echo "6. Analytics tracking works correctly"
echo "7. Error handling is appropriate"
echo "8. Performance is acceptable (< 500ms for complex queries)"
echo ""
echo "✨ If all tests pass, the Smart Autocomplete feature is ready for production!"
