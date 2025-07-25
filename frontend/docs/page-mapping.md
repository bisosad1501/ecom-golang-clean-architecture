# 🗺️ Page Mapping Guide

This document helps you quickly find which component corresponds to which page route.

## 📍 Quick Reference

| URL Route | Component File | Component Name | Description |
|-----------|----------------|----------------|-------------|
| `/` | `home-page.tsx` | `HomePage` | Landing page |
| `/about` | `about-page.tsx` | `AboutPage` | About us page |
| `/contact` | `contact-page.tsx` | `ContactPage` | Contact form |
| `/products` | `products-page.tsx` | `ProductsPage` | Product listing |
| `/products/[id]` | `product-detail-page.tsx` | `ProductDetailPage` | Single product |
| `/categories` | `categories-page.tsx` | `CategoriesPage` | Category listing |
| `/categories/[slug]` | `category-detail-page.tsx` | `CategoryDetailPage` | Category products |
| `/search` | `search-page.tsx` | `SearchPage` | Search results |
| `/cart` | `cart-page.tsx` | `CartPage` | Shopping cart |
| `/wishlist` | `wishlist-page.tsx` | `WishlistPage` | User wishlist |
| `/orders` | `order-history-page.tsx` | `OrderHistoryPage` | Order history |
| `/orders/[id]` | `order-detail-page.tsx` | `OrderDetailPage` | Order details |
| `/order-confirmation` | `order-confirmation-page.tsx` | `OrderConfirmationPage` | Order success |
| `/profile` | `profile-page.tsx` | `ProfilePage` | User profile |
| `/account` | `account-page.tsx` | `AccountPage` | Account settings |
| `/settings` | `settings-page.tsx` | `SettingsPage` | App settings |

## 🏗️ File Structure

```
src/
├── app/                          # Next.js App Router (URL routes)
│   ├── page.tsx                  # → HomePage
│   ├── about/page.tsx            # → AboutPage
│   ├── products/
│   │   ├── page.tsx              # → ProductsPage
│   │   └── [id]/page.tsx         # → ProductDetailPage
│   ├── orders/
│   │   ├── page.tsx              # → OrderHistoryPage
│   │   └── [id]/page.tsx         # → OrderDetailPage
│   └── ...
│
├── components/pages/             # Page Components (UI implementation)
│   ├── home-page.tsx             # HomePage component
│   ├── about-page.tsx            # AboutPage component
│   ├── products-page.tsx         # ProductsPage component
│   ├── product-detail-page.tsx   # ProductDetailPage component
│   ├── order-history-page.tsx    # OrderHistoryPage component
│   ├── order-detail-page.tsx     # OrderDetailPage component
│   └── ...
```

## 🔍 How to Find What You Need

### 1. **Looking for a specific page?**
- Check the URL in your browser
- Find the corresponding row in the table above
- The component file is in `src/components/pages/`

### 2. **Want to modify a page?**
- Edit the component file in `src/components/pages/`
- The app route file just imports and renders the component

### 3. **Adding a new page?**
1. Create component in `src/components/pages/new-page.tsx`
2. Create route in `src/app/new-route/page.tsx`
3. Import and render your component in the route file

## 📝 Naming Conventions

- **App Routes**: Follow Next.js conventions (`page.tsx`, `[id]/page.tsx`)
- **Page Components**: Use kebab-case (`home-page.tsx`, `product-detail-page.tsx`)
- **Component Names**: Use PascalCase (`HomePage`, `ProductDetailPage`)
- **Export**: Always use named exports (`export function HomePage()`)

## 🎯 Common Patterns

### App Route File Pattern:
```tsx
import { ComponentName } from '@/components/pages/component-file'

export default function RouteName() {
  return <ComponentName />
}
```

### Page Component Pattern:
```tsx
'use client'

import { ... } from '...'

export function ComponentName() {
  return (
    <div>
      {/* Page content */}
    </div>
  )
}
```

## 🚀 Quick Commands

```bash
# Find a component by name
find src/components/pages -name "*product*"

# Find which route uses a component
grep -r "ProductDetailPage" src/app/

# List all page components
ls src/components/pages/
```
