# 📝 Naming Conventions Guide

Consistent naming makes the codebase easier to navigate and maintain.

## 📁 File Naming

### Page Components
```
✅ CORRECT:
- home-page.tsx
- product-detail-page.tsx
- order-history-page.tsx

❌ INCORRECT:
- HomePage.tsx
- ProductDetailPage.tsx
- OrderHistoryPage.tsx
```

### Regular Components
```
✅ CORRECT:
- product-card.tsx
- user-profile.tsx
- cart-item.tsx

❌ INCORRECT:
- ProductCard.tsx
- UserProfile.tsx
- CartItem.tsx
```

### Hooks
```
✅ CORRECT:
- use-products.ts
- use-auth.ts
- use-cart.ts

❌ INCORRECT:
- useProducts.ts
- useAuth.ts
- useCart.ts
```

### Types
```
✅ CORRECT:
- product.types.ts
- auth.types.ts
- api.types.ts

❌ INCORRECT:
- ProductTypes.ts
- AuthTypes.ts
- ApiTypes.ts
```

## 🏷️ Component Naming

### Component Functions
```tsx
✅ CORRECT:
export function HomePage() { ... }
export function ProductCard() { ... }
export function UserProfile() { ... }

❌ INCORRECT:
export function homePage() { ... }
export function productCard() { ... }
export function userProfile() { ... }
```

### Component Props
```tsx
✅ CORRECT:
interface HomePageProps { ... }
interface ProductCardProps { ... }

❌ INCORRECT:
interface homePageProps { ... }
interface productCardProps { ... }
```

## 📂 Directory Naming

### Directories
```
✅ CORRECT:
- components/
- pages/
- auth/
- products/
- ui/

❌ INCORRECT:
- Components/
- Pages/
- Auth/
- Products/
- UI/
```

### Feature Directories
```
✅ CORRECT:
- src/components/products/
- src/components/auth/
- src/hooks/use-products/

❌ INCORRECT:
- src/components/Products/
- src/components/Auth/
- src/hooks/useProducts/
```

## 🔗 Import/Export Patterns

### Named Exports (Preferred)
```tsx
✅ CORRECT:
// In component file
export function HomePage() { ... }

// In importing file
import { HomePage } from '@/components/pages/home-page'
```

### Default Exports (App Router only)
```tsx
✅ CORRECT (App Router):
// In app/page.tsx
export default function Home() {
  return <HomePage />
}
```

### Index Files
```tsx
✅ CORRECT:
// In index.ts
export { HomePage } from './home-page'
export { ProductCard } from './product-card'

// Usage
import { HomePage, ProductCard } from '@/components/pages'
```

## 🎯 Variable Naming

### Constants
```tsx
✅ CORRECT:
const API_BASE_URL = 'https://api.example.com'
const MAX_ITEMS = 100
const DEFAULT_PAGE_SIZE = 20

❌ INCORRECT:
const apiBaseUrl = 'https://api.example.com'
const maxItems = 100
const defaultPageSize = 20
```

### Functions
```tsx
✅ CORRECT:
const handleSubmit = () => { ... }
const fetchProducts = async () => { ... }
const validateForm = () => { ... }

❌ INCORRECT:
const HandleSubmit = () => { ... }
const FetchProducts = async () => { ... }
const ValidateForm = () => { ... }
```

### State Variables
```tsx
✅ CORRECT:
const [isLoading, setIsLoading] = useState(false)
const [products, setProducts] = useState([])
const [currentPage, setCurrentPage] = useState(1)

❌ INCORRECT:
const [IsLoading, setIsLoading] = useState(false)
const [Products, setProducts] = useState([])
const [CurrentPage, setCurrentPage] = useState(1)
```

## 🗂️ URL and Route Naming

### App Router Files
```
✅ CORRECT:
- app/page.tsx (for /)
- app/products/page.tsx (for /products)
- app/products/[id]/page.tsx (for /products/123)

❌ INCORRECT:
- app/home.tsx
- app/products.tsx
- app/product-detail.tsx
```

### Route Segments
```
✅ CORRECT:
- /products
- /categories
- /order-confirmation
- /auth/login

❌ INCORRECT:
- /Products
- /Categories
- /orderConfirmation
- /auth/Login
```

## 🎨 CSS Class Naming

### Tailwind Classes (Preferred)
```tsx
✅ CORRECT:
<div className="flex items-center justify-between p-4 bg-white rounded-lg shadow-md">

❌ INCORRECT:
<div className="custom-header-style">
```

### Custom CSS (When needed)
```css
✅ CORRECT:
.product-card { ... }
.user-profile { ... }
.cart-item { ... }

❌ INCORRECT:
.ProductCard { ... }
.UserProfile { ... }
.CartItem { ... }
```

## 🚀 Quick Reference

| Type | Convention | Example |
|------|------------|---------|
| **Files** | kebab-case | `product-card.tsx` |
| **Components** | PascalCase | `ProductCard` |
| **Variables** | camelCase | `isLoading` |
| **Constants** | UPPER_SNAKE_CASE | `API_BASE_URL` |
| **Directories** | lowercase | `components/` |
| **URLs** | kebab-case | `/order-confirmation` |
| **CSS Classes** | kebab-case | `.product-card` |

## ✅ Checklist

Before creating new files:

- [ ] File name uses kebab-case
- [ ] Component name uses PascalCase
- [ ] Directory is lowercase
- [ ] Exports are named (not default)
- [ ] Imports use consistent paths
- [ ] Variables use camelCase
- [ ] Constants use UPPER_SNAKE_CASE
