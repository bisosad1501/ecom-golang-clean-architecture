#!/bin/bash

# Setup test users for the ecommerce application
# This script creates admin, moderator, and customer test accounts

API_BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Setting up test users..."

# Function to create user via API
create_user() {
    local email=$1
    local password=$2
    local first_name=$3
    local last_name=$4
    local role=$5
    
    echo "Creating user: $email"
    
    # Register user
    response=$(curl -s -X POST "$API_BASE_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{
            \"email\": \"$email\",
            \"password\": \"$password\",
            \"first_name\": \"$first_name\",
            \"last_name\": \"$last_name\"
        }")
    
    echo "Register response: $response"
    
    # Check if registration was successful
    if echo "$response" | grep -q "error"; then
        echo "❌ Failed to register $email"
        echo "Response: $response"
        return 1
    else
        echo "✅ Successfully registered $email"
        return 0
    fi
}

# Function to update user role in database
update_user_role() {
    local email=$1
    local role=$2
    
    echo "Updating role for $email to $role"
    
    # Connect to database and update role
    docker exec ecom_postgres psql -U postgres -d ecommerce_db -c "
        UPDATE users SET role = '$role' WHERE email = '$email';
    "
    
    if [ $? -eq 0 ]; then
        echo "✅ Successfully updated role for $email to $role"
    else
        echo "❌ Failed to update role for $email"
    fi
}

# Wait for API to be ready
echo "⏳ Waiting for API to be ready..."
for i in {1..30}; do
    if curl -s "$API_BASE_URL/../health" > /dev/null; then
        echo "✅ API is ready"
        break
    fi
    echo "Waiting... ($i/30)"
    sleep 2
done

# Create test users
echo ""
echo "📝 Creating test users..."

# Admin user
create_user "admin@ecom.com" "admin123" "Admin" "User" "admin"
if [ $? -eq 0 ]; then
    update_user_role "admin@ecom.com" "admin"
fi

echo ""

# Moderator user  
create_user "moderator@ecom.com" "moderator123" "Moderator" "User" "moderator"
if [ $? -eq 0 ]; then
    update_user_role "moderator@ecom.com" "moderator"
fi

echo ""

# Customer user
create_user "customer@ecom.com" "customer123" "Customer" "User" "customer"
if [ $? -eq 0 ]; then
    update_user_role "customer@ecom.com" "customer"
fi

echo ""
echo "🎉 User setup complete!"
echo ""
echo "Test accounts created:"
echo "👑 Admin:     admin@ecom.com     / admin123"
echo "🛠️  Moderator: moderator@ecom.com / moderator123" 
echo "🛒 Customer:  customer@ecom.com  / customer123"
echo ""
echo "You can now login at: http://localhost:3000/auth/login"
