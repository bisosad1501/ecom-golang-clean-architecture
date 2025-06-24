// API Configuration
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

// App Configuration
export const APP_NAME = 'EcomStore'
export const APP_DESCRIPTION = 'Modern E-commerce Platform'
export const APP_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'

// Authentication
export const AUTH_TOKEN_KEY = 'auth_token'
export const REFRESH_TOKEN_KEY = 'refresh_token'
export const USER_KEY = 'user'

// Pagination
export const DEFAULT_PAGE_SIZE = 12
export const DEFAULT_ADMIN_PAGE_SIZE = 20

// File Upload
export const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB
export const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp']

// Product Status
export const PRODUCT_STATUS = {
  DRAFT: 'draft',
  PUBLISHED: 'published',
  ARCHIVED: 'archived',
} as const

// Order Status
export const ORDER_STATUS = {
  PENDING: 'pending',
  CONFIRMED: 'confirmed',
  PROCESSING: 'processing',
  SHIPPED: 'shipped',
  DELIVERED: 'delivered',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

// Payment Status
export const PAYMENT_STATUS = {
  PENDING: 'pending',
  PROCESSING: 'processing',
  COMPLETED: 'completed',
  FAILED: 'failed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

// User Roles
export const USER_ROLES = {
  CUSTOMER: 'customer',
  ADMIN: 'admin',
  MODERATOR: 'moderator',
  SUPER_ADMIN: 'super_admin',
} as const

// Navigation
export const MAIN_NAV = [
  {
    title: 'Home',
    href: '/',
  },
  {
    title: 'Products',
    href: '/products',
  },
  {
    title: 'Categories',
    href: '/categories',
  },
  {
    title: 'About',
    href: '/about',
  },
  {
    title: 'Contact',
    href: '/contact',
  },
]

export const ADMIN_NAV = [
  {
    title: 'Dashboard',
    href: '/admin',
    icon: 'LayoutDashboard',
  },
  {
    title: 'Products',
    href: '/admin/products',
    icon: 'Package',
  },
  {
    title: 'Categories',
    href: '/admin/categories',
    icon: 'FolderTree',
  },
  {
    title: 'Orders',
    href: '/admin/orders',
    icon: 'ShoppingCart',
  },
  {
    title: 'Users',
    href: '/admin/users',
    icon: 'Users',
  },
  {
    title: 'Analytics',
    href: '/admin/analytics',
    icon: 'BarChart3',
  },
  {
    title: 'Settings',
    href: '/admin/settings',
    icon: 'Settings',
  },
]

export const USER_NAV = [
  {
    title: 'Profile',
    href: '/profile',
    icon: 'User',
  },
  {
    title: 'Orders',
    href: '/orders',
    icon: 'Package',
  },
  {
    title: 'Wishlist',
    href: '/wishlist',
    icon: 'Heart',
  },
  {
    title: 'Addresses',
    href: '/addresses',
    icon: 'MapPin',
  },
  {
    title: 'Settings',
    href: '/settings',
    icon: 'Settings',
  },
]

// Sort Options
export const PRODUCT_SORT_OPTIONS = [
  { label: 'Newest', value: 'created_at:desc' },
  { label: 'Oldest', value: 'created_at:asc' },
  { label: 'Price: Low to High', value: 'price:asc' },
  { label: 'Price: High to Low', value: 'price:desc' },
  { label: 'Name: A to Z', value: 'name:asc' },
  { label: 'Name: Z to A', value: 'name:desc' },
  { label: 'Best Selling', value: 'sales:desc' },
  { label: 'Top Rated', value: 'rating:desc' },
]

// Filter Options
export const PRICE_RANGES = [
  { label: 'Under $25', min: 0, max: 25 },
  { label: '$25 - $50', min: 25, max: 50 },
  { label: '$50 - $100', min: 50, max: 100 },
  { label: '$100 - $200', min: 100, max: 200 },
  { label: 'Over $200', min: 200, max: null },
]

// Social Links
export const SOCIAL_LINKS = {
  twitter: 'https://twitter.com/ecomstore',
  facebook: 'https://facebook.com/ecomstore',
  instagram: 'https://instagram.com/ecomstore',
  linkedin: 'https://linkedin.com/company/ecomstore',
  youtube: 'https://youtube.com/ecomstore',
}

// Contact Information
export const CONTACT_INFO = {
  email: 'support@ecomstore.com',
  phone: '+1 (555) 123-4567',
  address: '123 Commerce St, Business City, BC 12345',
  hours: 'Mon-Fri 9AM-6PM EST',
}

// SEO
export const DEFAULT_SEO = {
  title: 'EcomStore - Modern E-commerce Platform',
  description: 'Discover amazing products at great prices. Fast shipping, easy returns, and excellent customer service.',
  keywords: 'ecommerce, online shopping, products, deals, fashion, electronics',
  ogImage: '/images/og-image.jpg',
}

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  FORBIDDEN: 'Access denied.',
  NOT_FOUND: 'The requested resource was not found.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  SERVER_ERROR: 'Something went wrong. Please try again later.',
  RATE_LIMIT: 'Too many requests. Please try again later.',
}

// Success Messages
export const SUCCESS_MESSAGES = {
  PRODUCT_ADDED: 'Product added successfully!',
  PRODUCT_UPDATED: 'Product updated successfully!',
  PRODUCT_DELETED: 'Product deleted successfully!',
  ORDER_PLACED: 'Order placed successfully!',
  PROFILE_UPDATED: 'Profile updated successfully!',
  PASSWORD_CHANGED: 'Password changed successfully!',
  EMAIL_SENT: 'Email sent successfully!',
  ITEM_ADDED_TO_CART: 'Item added to cart!',
  ITEM_REMOVED_FROM_CART: 'Item removed from cart!',
  WISHLIST_UPDATED: 'Wishlist updated!',
}

// Loading States
export const LOADING_STATES = {
  IDLE: 'idle',
  LOADING: 'loading',
  SUCCESS: 'success',
  ERROR: 'error',
} as const

// Theme
export const THEME_OPTIONS = {
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system',
} as const

// Currency
export const CURRENCIES = [
  { code: 'USD', symbol: '$', name: 'US Dollar' },
  { code: 'EUR', symbol: '€', name: 'Euro' },
  { code: 'GBP', symbol: '£', name: 'British Pound' },
  { code: 'VND', symbol: '₫', name: 'Vietnamese Dong' },
]

// Countries
export const COUNTRIES = [
  { code: 'US', name: 'United States' },
  { code: 'CA', name: 'Canada' },
  { code: 'GB', name: 'United Kingdom' },
  { code: 'VN', name: 'Vietnam' },
  { code: 'AU', name: 'Australia' },
  { code: 'DE', name: 'Germany' },
  { code: 'FR', name: 'France' },
  { code: 'JP', name: 'Japan' },
]

// Regex Patterns
export const REGEX_PATTERNS = {
  EMAIL: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
  PHONE: /^\+?[\d\s\-\(\)]+$/,
  PASSWORD: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$!%*?&]{8,}$/,
  SLUG: /^[a-z0-9]+(?:-[a-z0-9]+)*$/,
  HEX_COLOR: /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/,
}

// Date Formats
export const DATE_FORMATS = {
  SHORT: 'MMM d, yyyy',
  LONG: 'MMMM d, yyyy',
  FULL: 'EEEE, MMMM d, yyyy',
  TIME: 'h:mm a',
  DATETIME: 'MMM d, yyyy h:mm a',
  ISO: "yyyy-MM-dd'T'HH:mm:ss.SSSxxx",
}

// Animation Durations
export const ANIMATION_DURATION = {
  FAST: 150,
  NORMAL: 300,
  SLOW: 500,
} as const

// Breakpoints (matching Tailwind CSS)
export const BREAKPOINTS = {
  SM: 640,
  MD: 768,
  LG: 1024,
  XL: 1280,
  '2XL': 1536,
} as const
