# User Setup Guide

## Quick Setup

To create admin, moderator, and customer test accounts:

```bash
make setup-users
```

This will create 3 test accounts with appropriate roles.

## Test Accounts

| Role | Email | Password | Permissions |
|------|-------|----------|-------------|
| **Admin** | admin@ecom.com | admin123 | Full system access, manage users, products, orders |
| **Moderator** | moderator@ecom.com | moderator123 | Manage products and inventory only |
| **Customer** | customer@ecom.com | customer123 | Browse products, manage cart, place orders |

## Login

### Frontend
- Open: http://localhost:3000
- Login with any of the above credentials

### API
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@ecom.com",
    "password": "admin123"
  }'
```

## User Roles & Permissions

### Admin (`admin`)
- Full access to all features
- Can manage users (activate/deactivate)
- Can manage products, categories
- Can view and manage all orders
- Access to admin panel

### Moderator (`moderator`) 
- Can create, update, delete products
- Can manage product inventory
- Cannot manage users
- Cannot access user management

### Customer (`customer`)
- Can browse products and categories  
- Can manage shopping cart
- Can place and view own orders
- Can update own profile
- No admin access

## Notes

- The setup script can be run multiple times safely
- Users are created via API then roles are updated in database
- All accounts are active by default
- Passwords are simple for testing purposes only
