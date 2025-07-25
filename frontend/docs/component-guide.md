# 🧩 Component Organization Guide

This guide helps you understand how components are organized and where to find what you need.

## 📁 Directory Structure

```
src/components/
├── admin/              # Admin panel components
├── auth/               # Authentication components
├── cart/               # Shopping cart components
├── categories/         # Category-related components
├── checkout/           # Checkout process components
├── forms/              # Reusable form components
├── layout/             # Layout components (header, footer, etc.)
├── modals/             # Modal dialogs
├── notifications/      # Notification components
├── pages/              # 🎯 PAGE COMPONENTS (main content)
├── products/           # Product-related components
├── profile/            # User profile components
├── reviews/            # Review and rating components
├── search/             # Search functionality
└── ui/                 # Reusable UI components
```

## 🎯 Page Components (`src/components/pages/`)

These are the main content components for each page:

### 🏠 General Pages
- `home-page.tsx` - Landing page with hero, features, products
- `about-page.tsx` - Company information and team
- `contact-page.tsx` - Contact form and information

### 🛍️ Shopping Pages
- `products-page.tsx` - Product listing with filters
- `product-detail-page.tsx` - Single product view
- `categories-page.tsx` - Category listing
- `category-detail-page.tsx` - Products in a category
- `search-page.tsx` - Search results

### 🛒 Cart & Orders
- `cart-page.tsx` - Shopping cart management
- `wishlist-page.tsx` - User's saved products
- `order-history-page.tsx` - List of user's orders
- `order-detail-page.tsx` - Detailed order view
- `order-confirmation-page.tsx` - Order success page

### 👤 User Pages
- `profile-page.tsx` - User profile and stats
- `account-page.tsx` - Account management
- `settings-page.tsx` - App preferences

### 🎨 Special Pages
- `categories-page-simple.tsx` - Simplified category view

## 🔧 UI Components (`src/components/ui/`)

Reusable building blocks:

### Basic UI
- `button.tsx` - Button variants
- `input.tsx` - Form inputs
- `card.tsx` - Content cards
- `badge.tsx` - Status badges

### Complex UI
- `pagination.tsx` - Page navigation
- `category-card.tsx` - Category display
- `animated-background.tsx` - Animated backgrounds
- `order-timeline.tsx` - Order status timeline

## 🏗️ Feature Components

### Products (`src/components/products/`)
- `product-card.tsx` - Product display card
- `product-list-card.tsx` - List view product
- `product-filters.tsx` - Filter controls
- `product-sort.tsx` - Sort controls

### Cart (`src/components/cart/`)
- `cart-item.tsx` - Individual cart item
- `cart-summary.tsx` - Cart totals
- `global-cart-conflict-modal.tsx` - Cart conflict handling

### Auth (`src/components/auth/`)
- `login-form.tsx` - Login form
- `auth-layout.tsx` - Auth page layout
- `permission-guard.tsx` - Access control

## 🎯 How to Find Components

### 1. **By Feature**
```bash
# Product-related
ls src/components/products/

# User-related
ls src/components/profile/
ls src/components/auth/

# Shopping-related
ls src/components/cart/
ls src/components/checkout/
```

### 2. **By Type**
```bash
# Page components (main content)
ls src/components/pages/

# Reusable UI
ls src/components/ui/

# Layout components
ls src/components/layout/
```

### 3. **By Search**
```bash
# Find component by name
find src/components -name "*product*"

# Find component by content
grep -r "ProductCard" src/components/
```

## 📝 Component Patterns

### Page Component
```tsx
'use client'

export function PageName() {
  return (
    <div className="container mx-auto">
      {/* Page content */}
    </div>
  )
}
```

### UI Component
```tsx
interface Props {
  // Props definition
}

export function ComponentName({ ...props }: Props) {
  return (
    <div className="...">
      {/* Component content */}
    </div>
  )
}
```

### Feature Component
```tsx
'use client'

import { useHook } from '@/hooks/use-hook'

export function FeatureComponent() {
  const { data } = useHook()
  
  return (
    <div>
      {/* Feature implementation */}
    </div>
  )
}
```

## 🚀 Quick Tips

1. **Page components** = Main content for routes
2. **UI components** = Reusable building blocks  
3. **Feature components** = Business logic + UI
4. **Layout components** = Page structure
5. **Use search** when you can't find something
6. **Check imports** to understand dependencies
