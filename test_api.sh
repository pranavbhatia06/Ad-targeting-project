#!/bin/bash

# Targeting Engine API Test Script
echo "🎯 Testing Targeting Engine API"
echo "================================"

BASE_URL="http://localhost:7000"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${YELLOW}Testing: $description${NC}"
    echo "Endpoint: $method $endpoint"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
                   -H "Content-Type: application/json" \
                   -d "$data")
    fi
    
    # Extract HTTP status code (last line)
    http_code=$(echo "$response" | tail -n1)
    # Extract response body (all but last line)
    body=$(echo "$response" | head -n -1)
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✅ SUCCESS (HTTP $http_code)${NC}"
        echo "Response: $body" | jq . 2>/dev/null || echo "Response: $body"
    else
        echo -e "${RED}❌ FAILED (HTTP $http_code)${NC}"
        echo "Response: $body"
    fi
}

# Test 1: Health Check
test_endpoint "GET" "/health" "" "Health Check"

# Test 2: Get Campaigns with GET (India + Android)
test_endpoint "GET" "/api/v1/campaigns?app=myapp&country=india&os=android" "" "Get Campaigns for India + Android (GET)"

# Test 3: Get Campaigns with GET (USA + iOS with limit)
test_endpoint "GET" "/api/v1/campaigns?app=myapp&country=usa&os=ios&limit=2" "" "Get Campaigns for USA + iOS with Limit (GET)"

# Test 4: Get Campaigns with POST (India + Android)
test_endpoint "POST" "/api/v1/campaigns" '{
    "app": "myapp",
    "country": "india",
    "os": "android",
    "limit": 3
}' "Get Campaigns for India + Android (POST)"

# Test 5: Get Campaigns with POST (USA + iOS)
test_endpoint "POST" "/api/v1/campaigns" '{
    "app": "otherapp",
    "country": "usa",
    "os": "ios"
}' "Get Campaigns for USA + iOS (POST)"

# Test 6: Get Campaigns with invalid parameters
test_endpoint "GET" "/api/v1/campaigns?app=unknownapp&country=unknown&os=unknown" "" "Get Campaigns for Unknown Parameters"

# Test 7: Missing required parameters
test_endpoint "GET" "/api/v1/campaigns?app=myapp&country=india" "" "Missing OS Parameter"

# Test 8: Refresh Targeting Data
test_endpoint "POST" "/api/v1/refresh-targeting" "" "Refresh Targeting Data"

# Test 9: Get Campaigns after refresh
test_endpoint "GET" "/api/v1/campaigns?app=myapp&country=india&os=android&limit=1" "" "Get Campaigns After Refresh"

echo -e "\n${YELLOW}🏁 Testing Complete!${NC}"
echo "================================" 