#!/bin/bash

# Create sample products for testing
API_BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Creating sample products..."

# First, login as admin to get token
echo "🔐 Logging in as admin..."
login_response=$(curl -s -X POST "$API_BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@ecom.com",
        "password": "admin123"
    }')

# Extract token
token=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)

if [ -z "$token" ]; then
    echo "❌ Failed to get admin token"
    echo "Response: $login_response"
    exit 1
fi

echo "✅ Got admin token"

# Get categories
echo "📂 Getting categories..."
categories_response=$(curl -s -X GET "$API_BASE_URL/categories" \
    -H "Content-Type: application/json")

echo "Categories: $categories_response"

# Extract category IDs (assuming we have Electronics, Clothing, Books, Home & Garden)
electronics_id=$(echo "$categories_response" | grep -o '"id":"[^"]*","name":"Electronics"' | cut -d'"' -f4)
clothing_id=$(echo "$categories_response" | grep -o '"id":"[^"]*","name":"Clothing"' | cut -d'"' -f4)
books_id=$(echo "$categories_response" | grep -o '"id":"[^"]*","name":"Books"' | cut -d'"' -f4)
home_id=$(echo "$categories_response" | grep -o '"id":"[^"]*","name":"Home & Garden"' | cut -d'"' -f4)

echo "Electronics ID: $electronics_id"
echo "Clothing ID: $clothing_id"
echo "Books ID: $books_id"
echo "Home ID: $home_id"

# Function to create product
create_product() {
    local name="$1"
    local description="$2"
    local sku="$3"
    local price="$4"
    local stock="$5"
    local category_id="$6"
    
    echo "Creating product: $name"
    
    response=$(curl -s -X POST "$API_BASE_URL/admin/products" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $token" \
        -d "{
            \"name\": \"$name\",
            \"description\": \"$description\",
            \"sku\": \"$sku\",
            \"price\": $price,
            \"stock\": $stock,
            \"category_id\": \"$category_id\",
            \"status\": \"active\"
        }")
    
    if echo "$response" | grep -q "error"; then
        echo "❌ Failed to create $name"
        echo "Response: $response"
    else
        echo "✅ Successfully created $name"
    fi
}

# Create sample products
echo ""
echo "📱 Creating Electronics products..."

if [ ! -z "$electronics_id" ]; then
    create_product "iPhone 15 Pro" "Latest iPhone with advanced camera system and A17 Pro chip" "IPHONE15PRO-001" 999.99 50 "$electronics_id"
    create_product "Samsung Galaxy S24" "Flagship Android phone with AI features" "GALAXY-S24-001" 899.99 30 "$electronics_id"
    create_product "MacBook Air M3" "Lightweight laptop with M3 chip" "MACBOOK-AIR-M3" 1299.99 20 "$electronics_id"
    create_product "Sony WH-1000XM5" "Premium noise-canceling headphones" "SONY-WH1000XM5" 399.99 100 "$electronics_id"
fi

echo ""
echo "👕 Creating Clothing products..."

if [ ! -z "$clothing_id" ]; then
    create_product "Nike Air Max 270" "Comfortable running shoes with Air Max technology" "NIKE-AIRMAX270" 150.00 75 "$clothing_id"
    create_product "Adidas Ultraboost 22" "High-performance running shoes" "ADIDAS-UB22" 180.00 60 "$clothing_id"
    create_product "Levi's 501 Jeans" "Classic straight-leg jeans" "LEVIS-501-BLUE" 89.99 120 "$clothing_id"
    create_product "Nike Dri-FIT T-Shirt" "Moisture-wicking athletic t-shirt" "NIKE-DRIFIT-TEE" 29.99 200 "$clothing_id"
fi

echo ""
echo "📚 Creating Books products..."

if [ ! -z "$books_id" ]; then
    create_product "Clean Code" "A Handbook of Agile Software Craftsmanship by Robert C. Martin" "BOOK-CLEANCODE" 42.99 50 "$books_id"
    create_product "The Pragmatic Programmer" "Your Journey to Mastery by David Thomas and Andrew Hunt" "BOOK-PRAGPROG" 39.99 40 "$books_id"
    create_product "Design Patterns" "Elements of Reusable Object-Oriented Software" "BOOK-DESIGNPAT" 54.99 30 "$books_id"
fi

echo ""
echo "🏠 Creating Home & Garden products..."

if [ ! -z "$home_id" ]; then
    create_product "Dyson V15 Detect" "Cordless vacuum cleaner with laser detection" "DYSON-V15-001" 749.99 25 "$home_id"
    create_product "Instant Pot Duo 7-in-1" "Multi-use pressure cooker" "INSTANTPOT-DUO7" 99.99 80 "$home_id"
    create_product "Philips Hue Smart Bulbs" "Color-changing LED smart bulbs (4-pack)" "PHILIPS-HUE-4PK" 199.99 60 "$home_id"
fi

echo ""
echo "🎉 Sample products creation complete!"
echo "You can now browse products at: http://localhost:3000"
