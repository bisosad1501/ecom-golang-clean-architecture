#!/bin/bash

# Script to create admin, moderator, and customer users
# This script creates users via API and updates their roles directly in database

BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Setting up admin and test accounts..."
echo ""

# Function to create user and update role
create_user_with_role() {
    local email=$1
    local password=$2
    local first_name=$3
    local last_name=$4
    local phone=$5
    local role=$6

    echo "Creating $role user: $email"
    
    # Create user via API
    RESPONSE=$(curl -s -X POST ${BASE_URL}/auth/register \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$email\",
        \"password\": \"$password\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\",
        \"phone\": \"$phone\"
      }")

    # Check if user was created or already exists
    if echo "$RESPONSE" | grep -q "User registered successfully\|already exists"; then
        # Update role in database if not customer
        if [ "$role" != "customer" ]; then
            docker exec -i ecom_postgres psql -U postgres -d ecommerce_db -c \
                "UPDATE users SET role = '$role', updated_at = NOW() WHERE email = '$email';" > /dev/null 2>&1
        fi
        echo "✅ $role user created/updated: $email"
    else
        echo "⚠️  Issue with $email: $RESPONSE"
    fi
}

# Create users
create_user_with_role "admin@ecom.com" "admin123" "Admin" "User" "+1234567890" "admin"
create_user_with_role "moderator@ecom.com" "moderator123" "Moderator" "User" "+1234567891" "moderator"  
create_user_with_role "customer@ecom.com" "customer123" "Customer" "User" "+1234567892" "customer"

echo ""
echo "🔍 Verifying users in database..."
docker exec -i ecom_postgres psql -U postgres -d ecommerce_db <<EOF
SELECT email, first_name, last_name, role, is_active 
FROM users 
WHERE email IN ('admin@ecom.com', 'moderator@ecom.com', 'customer@ecom.com')
ORDER BY 
  CASE role 
    WHEN 'admin' THEN 1 
    WHEN 'moderator' THEN 2 
    WHEN 'customer' THEN 3 
  END;
EOF

echo ""
echo "✅ Setup completed successfully!"
echo ""
echo "📝 Login Credentials:"
echo "=================================="
echo "1. ADMIN"
echo "   Email: admin@ecom.com"
echo "   Password: admin123"
echo "   Role: admin"
echo "   Permissions: Full access to all features"
echo ""
echo "2. MODERATOR"  
echo "   Email: moderator@ecom.com"
echo "   Password: moderator123"
echo "   Role: moderator"
echo "   Permissions: Manage products, inventory (not users)"
echo ""
echo "3. CUSTOMER"
echo "   Email: customer@ecom.com" 
echo "   Password: customer123"
echo "   Role: customer"
echo "   Permissions: Browse, cart, orders, profile"
echo ""
echo "🌐 Access Points:"
echo "   Frontend: http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   Login API: POST ${BASE_URL}/auth/login"
echo ""
echo "🔑 User Role Hierarchy:"
echo "   ADMIN > MODERATOR > CUSTOMER"
