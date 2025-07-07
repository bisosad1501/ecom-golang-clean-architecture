// ===== UNIFIED CONSTANTS EXPORT =====
// Export all constants from organized files
export * from './design-tokens'
export * from './api'
export * from './app'

// Re-export commonly used constants for convenience
export { DESIGN_TOKENS } from './design-tokens'
export { API_CONFIG, HTTP_STATUS, QUERY_KEYS, ERROR_MESSAGES, SUCCESS_MESSAGES } from './api'
export { APP_CONFIG, ROUTES, STORAGE_KEYS, ANALYTICS_EVENTS, ORDER_STATUS, PAYMENT_STATUS, USER_ROLES } from './app'

// Legacy constants (for backward compatibility - will be deprecated)
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
export const APP_NAME = 'ShopHub'
export const APP_DESCRIPTION = 'Your Ultimate Shopping Destination'
export const APP_URL = process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'

// Authentication (migrated to STORAGE_KEYS)
export const AUTH_TOKEN_KEY = 'shophub-auth-token'
export const REFRESH_TOKEN_KEY = 'shophub-refresh-token'
export const USER_KEY = 'shophub-user'

// Pagination (migrated to PAGINATION)
export const DEFAULT_PAGE_SIZE = 12
export const DEFAULT_ADMIN_PAGE_SIZE = 20

// File Upload (migrated to UPLOAD_LIMITS)
export const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB
export const ALLOWED_IMAGE_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/webp']

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
  twitter: 'https://twitter.com/bihub',
  facebook: 'https://facebook.com/bihub',
  instagram: 'https://instagram.com/bihub',
  linkedin: 'https://linkedin.com/company/bihub',
  youtube: 'https://youtube.com/bihub',
}

// Contact Information
export const CONTACT_INFO = {
  email: 'support@bihub.com',
  phone: '+1 (555) 123-4567',
  address: '123 Commerce St, Business City, BC 12345',
  hours: 'Mon-Fri 9AM-6PM EST',
}

// SEO
export const DEFAULT_SEO = {
  title: 'BiHub - Your Ultimate Shopping Hub',
  description: 'Discover premium products at BiHub. Fast shipping, easy returns, and exceptional shopping experience.',
  keywords: 'ecommerce, online shopping, products, deals, fashion, electronics, bihub, shopping hub',
  ogImage: '/images/og-image.jpg',
}

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
