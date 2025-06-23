# API Documentation

Complete API reference for the E-commerce system.

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Response Format

All API responses follow this format:

### Success Response
```json
{
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "error": "Error message",
  "details": "Additional error details"
}
```

## Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `422` - Unprocessable Entity
- `500` - Internal Server Error

## Endpoints

### Authentication

#### Register User
```http
POST /auth/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "phone": "+1234567890"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "customer",
    "is_active": true,
    "created_at": "2025-01-23T00:00:00Z"
  }
}
```

#### Login User
```http
POST /auth/login
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "token": "jwt_token_here",
    "expires_at": "2025-01-24T00:00:00Z",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "customer"
    }
  }
}
```

### User Management

#### Get User Profile
```http
GET /users/profile
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone": "+1234567890",
    "role": "customer",
    "is_active": true,
    "profile": {
      "date_of_birth": "1990-01-01",
      "gender": "male",
      "address": "123 Main St",
      "city": "New York",
      "country": "USA"
    },
    "created_at": "2025-01-23T00:00:00Z",
    "updated_at": "2025-01-23T00:00:00Z"
  }
}
```

#### Update User Profile
```http
PUT /users/profile
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "first_name": "John Updated",
  "last_name": "Doe Updated",
  "phone": "+1234567891",
  "date_of_birth": "1990-01-01",
  "gender": "male",
  "address": "456 New St",
  "city": "Los Angeles",
  "country": "USA"
}
```

### Categories

#### Get Categories
```http
GET /categories?limit=10&offset=0
```

**Query Parameters:**
- `limit` (optional): Number of items to return (default: 10, max: 100)
- `offset` (optional): Number of items to skip (default: 0)

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Electronics",
      "description": "Electronic devices and accessories",
      "slug": "electronics",
      "image": "https://example.com/image.jpg",
      "parent_id": null,
      "is_active": true,
      "sort_order": 1,
      "level": 0,
      "path": "Electronics",
      "created_at": "2025-01-23T00:00:00Z",
      "updated_at": "2025-01-23T00:00:00Z"
    }
  ]
}
```

#### Get Category Tree
```http
GET /categories/tree
```

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Electronics",
      "slug": "electronics",
      "children": [
        {
          "id": "uuid",
          "name": "Smartphones",
          "slug": "smartphones",
          "parent_id": "parent_uuid",
          "children": []
        }
      ]
    }
  ]
}
```

### Products

#### Get Products
```http
GET /products?limit=10&offset=0
```

**Query Parameters:**
- `limit` (optional): Number of items to return
- `offset` (optional): Number of items to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "iPhone 15",
      "description": "Latest iPhone model",
      "sku": "IPHONE15-001",
      "price": 999.99,
      "stock": 50,
      "category": {
        "id": "uuid",
        "name": "Smartphones",
        "slug": "smartphones"
      },
      "images": [
        {
          "id": "uuid",
          "url": "https://example.com/image1.jpg",
          "alt_text": "iPhone 15 front view",
          "sort_order": 1
        }
      ],
      "tags": [
        {
          "id": "uuid",
          "name": "Featured",
          "slug": "featured"
        }
      ],
      "status": "active",
      "is_digital": false,
      "weight": 0.2,
      "dimensions": "150x75x8",
      "created_at": "2025-01-23T00:00:00Z",
      "updated_at": "2025-01-23T00:00:00Z"
    }
  ]
}
```

#### Search Products
```http
GET /products/search?q=iphone&category_id=uuid&min_price=100&max_price=1000&limit=10
```

**Query Parameters:**
- `q` (optional): Search query
- `category_id` (optional): Filter by category
- `min_price` (optional): Minimum price filter
- `max_price` (optional): Maximum price filter
- `status` (optional): Product status filter
- `sort_by` (optional): Sort field (default: created_at)
- `sort_order` (optional): Sort order (asc/desc, default: desc)
- `limit` (optional): Number of items to return
- `offset` (optional): Number of items to skip

### Shopping Cart

#### Get Cart
```http
GET /cart
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "items": [
      {
        "id": "uuid",
        "product": {
          "id": "uuid",
          "name": "iPhone 15",
          "price": 999.99,
          "image": "https://example.com/image.jpg"
        },
        "quantity": 2,
        "price": 999.99,
        "subtotal": 1999.98
      }
    ],
    "item_count": 2,
    "total": 1999.98,
    "created_at": "2025-01-23T00:00:00Z",
    "updated_at": "2025-01-23T00:00:00Z"
  }
}
```

#### Add to Cart
```http
POST /cart/items
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "product_id": "uuid",
  "quantity": 2
}
```

### Orders

#### Create Order
```http
POST /orders
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "shipping_address": {
    "first_name": "John",
    "last_name": "Doe",
    "address1": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zip_code": "10001",
    "country": "USA",
    "phone": "+1234567890"
  },
  "billing_address": {
    // Same structure as shipping_address (optional)
  },
  "payment_method": "credit_card",
  "tax_rate": 0.08,
  "shipping_cost": 10.00,
  "discount_amount": 0.00,
  "notes": "Please deliver after 5 PM"
}
```

**Response:**
```json
{
  "message": "Order created successfully",
  "data": {
    "id": "uuid",
    "order_number": "ORD-2025-001",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe"
    },
    "items": [
      {
        "id": "uuid",
        "product": {
          "id": "uuid",
          "name": "iPhone 15",
          "sku": "IPHONE15-001"
        },
        "quantity": 2,
        "price": 999.99,
        "total": 1999.98
      }
    ],
    "status": "pending",
    "payment_status": "pending",
    "subtotal": 1999.98,
    "tax_amount": 159.998,
    "shipping_amount": 10.00,
    "discount_amount": 0.00,
    "total": 2169.978,
    "currency": "USD",
    "shipping_address": {
      "first_name": "John",
      "last_name": "Doe",
      "address1": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip_code": "10001",
      "country": "USA",
      "phone": "+1234567890"
    },
    "payment": {
      "id": "uuid",
      "amount": 2169.978,
      "currency": "USD",
      "method": "credit_card",
      "status": "pending"
    },
    "item_count": 2,
    "can_be_cancelled": true,
    "can_be_refunded": false,
    "created_at": "2025-01-23T00:00:00Z",
    "updated_at": "2025-01-23T00:00:00Z"
  }
}
```

## Error Handling

### Validation Errors
```json
{
  "error": "Invalid request format",
  "details": "email: must be a valid email address"
}
```

### Authentication Errors
```json
{
  "error": "Invalid token"
}
```

### Business Logic Errors
```json
{
  "error": "Insufficient stock",
  "details": "Only 5 items available in stock"
}
```

## Rate Limiting

API endpoints are rate limited:
- **Authentication endpoints**: 5 requests per minute
- **General endpoints**: 100 requests per minute
- **Admin endpoints**: 200 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1642723200
```
