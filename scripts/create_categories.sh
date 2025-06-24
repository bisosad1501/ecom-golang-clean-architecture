#!/bin/bash

# Script để tạo categories cơ bản cho test
API_URL="http://localhost:8080/api/v1"

# Lấy admin token (cần thay thế bằng token thực tế)
# Hoặc login để lấy token
echo "Creating basic categories for testing..."

# Categories cơ bản
categories=(
    '{"name":"Electronics","slug":"electronics","description":"Electronic devices and gadgets"}'
    '{"name":"Fashion","slug":"fashion","description":"Clothing and accessories"}'
    '{"name":"Home & Garden","slug":"home-garden","description":"Home and garden products"}'
    '{"name":"Sports","slug":"sports","description":"Sports and outdoor equipment"}'
    '{"name":"Books","slug":"books","description":"Books and educational materials"}'
)

for category in "${categories[@]}"; do
    echo "Creating category: $category"
    curl -X POST "$API_URL/admin/categories" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
        -d "$category"
    echo ""
done

echo "Categories created successfully!"
