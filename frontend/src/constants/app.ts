// Application Configuration Constants
export const APP_CONFIG = {
  // App metadata
  NAME: 'BiHub',
  DESCRIPTION: 'Your Ultimate Shopping Destination',
  VERSION: '1.0.0',
  AUTHOR: 'BiHub Team',
  
  // URLs
  DOMAIN: process.env.NEXT_PUBLIC_DOMAIN || 'localhost:3000',
  BASE_URL: process.env.NEXT_PUBLIC_BASE_URL || 'http://localhost:3000',
  
  // Social links
  SOCIAL: {
    FACEBOOK: 'https://facebook.com/bihub',
    TWITTER: 'https://twitter.com/bihub',
    INSTAGRAM: 'https://instagram.com/bihub',
    LINKEDIN: 'https://linkedin.com/company/bihub',
    YOUTUBE: 'https://youtube.com/bihub',
  },

  // Contact information
  CONTACT: {
    EMAIL: 'support@bihub.com',
    PHONE: '+1 (555) 123-4567',
    ADDRESS: '123 Commerce Street, Business District, City 12345',
    SUPPORT_HOURS: 'Monday - Friday: 9AM - 6PM EST',
  },
  
  // Feature flags
  FEATURES: {
    DARK_MODE: true,
    WISHLIST: true,
    REVIEWS: true,
    RECOMMENDATIONS: true,
    LIVE_CHAT: false,
    NOTIFICATIONS: true,
    ANALYTICS: true,
    A_B_TESTING: false,
  },
  
  // Limits and constraints
  LIMITS: {
    CART_MAX_ITEMS: 100,
    WISHLIST_MAX_ITEMS: 500,
    SEARCH_HISTORY_MAX: 10,
    RECENT_PRODUCTS_MAX: 20,
    PRODUCT_IMAGES_MAX: 10,
    REVIEW_LENGTH_MAX: 1000,
    USERNAME_LENGTH_MAX: 50,
    PASSWORD_LENGTH_MIN: 8,
  },
} as const

// Navigation routes
export const ROUTES = {
  // Public routes
  HOME: '/',
  PRODUCTS: '/products',
  PRODUCT_DETAIL: '/products/[id]',
  CATEGORIES: '/categories',
  CATEGORY_DETAIL: '/categories/[slug]',
  SEARCH: '/search',
  ABOUT: '/about',
  CONTACT: '/contact',
  FAQ: '/faq',
  TERMS: '/terms',
  PRIVACY: '/privacy',
  
  // Auth routes
  LOGIN: '/auth/login',
  REGISTER: '/auth/register',
  FORGOT_PASSWORD: '/auth/forgot-password',
  RESET_PASSWORD: '/auth/reset-password',
  VERIFY_EMAIL: '/auth/verify-email',
  
  // User routes
  PROFILE: '/profile',
  ORDERS: '/profile/orders',
  ORDER_DETAIL: '/profile/orders/[id]',
  ADDRESSES: '/profile/addresses',
  WISHLIST: '/profile/wishlist',
  SETTINGS: '/profile/settings',
  
  // Shopping routes
  CART: '/cart',
  CHECKOUT: '/checkout',
  CHECKOUT_SUCCESS: '/checkout/success',
  CHECKOUT_CANCEL: '/checkout/cancel',
  
  // Admin routes
  ADMIN: '/admin',
  ADMIN_DASHBOARD: '/admin/dashboard',
  ADMIN_PRODUCTS: '/admin/products',
  ADMIN_PRODUCT_NEW: '/admin/products/new',
  ADMIN_PRODUCT_EDIT: '/admin/products/[id]/edit',
  ADMIN_ORDERS: '/admin/orders',
  ADMIN_ORDER_DETAIL: '/admin/orders/[id]',
  ADMIN_USERS: '/admin/users',
  ADMIN_CATEGORIES: '/admin/categories',
  ADMIN_ANALYTICS: '/admin/analytics',
  ADMIN_SETTINGS: '/admin/settings',
} as const

// Local storage keys
export const STORAGE_KEYS = {
  THEME: 'bihub-theme',
  CART: 'bihub-cart',
  WISHLIST: 'bihub-wishlist',
  SEARCH_HISTORY: 'bihub-search-history',
  RECENT_PRODUCTS: 'bihub-recent-products',
  USER_PREFERENCES: 'bihub-user-preferences',
  AUTH_TOKEN: 'bihub-auth-token',
  REFRESH_TOKEN: 'bihub-refresh-token',
  LANGUAGE: 'bihub-language',
  CURRENCY: 'bihub-currency',
} as const

// Event names for analytics
export const ANALYTICS_EVENTS = {
  // Page views
  PAGE_VIEW: 'page_view',
  
  // Product events
  PRODUCT_VIEW: 'product_view',
  PRODUCT_SEARCH: 'product_search',
  PRODUCT_FILTER: 'product_filter',
  PRODUCT_SORT: 'product_sort',
  
  // Cart events
  ADD_TO_CART: 'add_to_cart',
  REMOVE_FROM_CART: 'remove_from_cart',
  UPDATE_CART: 'update_cart',
  VIEW_CART: 'view_cart',
  CLEAR_CART: 'clear_cart',
  
  // Checkout events
  BEGIN_CHECKOUT: 'begin_checkout',
  ADD_PAYMENT_INFO: 'add_payment_info',
  PURCHASE: 'purchase',
  
  // User events
  SIGN_UP: 'sign_up',
  LOGIN: 'login',
  LOGOUT: 'logout',
  
  // Wishlist events
  ADD_TO_WISHLIST: 'add_to_wishlist',
  REMOVE_FROM_WISHLIST: 'remove_from_wishlist',
  
  // Engagement events
  SHARE: 'share',
  REVIEW_SUBMIT: 'review_submit',
  NEWSLETTER_SIGNUP: 'newsletter_signup',
} as const

// Toast notification types
export const TOAST_TYPES = {
  SUCCESS: 'success',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
} as const

// Modal types
export const MODAL_TYPES = {
  CONFIRM: 'confirm',
  ALERT: 'alert',
  PRODUCT_QUICK_VIEW: 'product_quick_view',
  AUTH: 'auth',
  CART: 'cart',
  SEARCH: 'search',
} as const

// Product sort options
export const PRODUCT_SORT_OPTIONS = [
  { value: 'created_at', label: 'Newest First', order: 'desc' },
  { value: 'created_at', label: 'Oldest First', order: 'asc' },
  { value: 'name', label: 'Name A-Z', order: 'asc' },
  { value: 'name', label: 'Name Z-A', order: 'desc' },
  { value: 'price', label: 'Price Low to High', order: 'asc' },
  { value: 'price', label: 'Price High to Low', order: 'desc' },
  { value: 'rating', label: 'Highest Rated', order: 'desc' },
  { value: 'popularity', label: 'Most Popular', order: 'desc' },
] as const

// Order status options
export const ORDER_STATUS = {
  PENDING: 'pending',
  CONFIRMED: 'confirmed',
  PROCESSING: 'processing',
  SHIPPED: 'shipped',
  DELIVERED: 'delivered',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

// Payment status options (synchronized with backend)
export const PAYMENT_STATUS = {
  PENDING: 'pending',
  PROCESSING: 'processing',
  PAID: 'paid',
  COMPLETED: 'completed',  // Alias for paid
  FAILED: 'failed',
  CANCELLED: 'cancelled',
  REFUNDED: 'refunded',
} as const

// User roles
export const USER_ROLES = {
  CUSTOMER: 'customer',
  ADMIN: 'admin',
  SUPER_ADMIN: 'super_admin',
  MODERATOR: 'moderator',
} as const

// Product availability status
export const PRODUCT_STATUS = {
  ACTIVE: 'active',
  INACTIVE: 'inactive',
  OUT_OF_STOCK: 'out_of_stock',
  DISCONTINUED: 'discontinued',
  DRAFT: 'draft',
} as const

// Currency options
export const CURRENCIES = [
  { code: 'USD', symbol: '$', name: 'US Dollar' },
  { code: 'EUR', symbol: '€', name: 'Euro' },
  { code: 'GBP', symbol: '£', name: 'British Pound' },
  { code: 'JPY', symbol: '¥', name: 'Japanese Yen' },
  { code: 'CAD', symbol: 'C$', name: 'Canadian Dollar' },
  { code: 'AUD', symbol: 'A$', name: 'Australian Dollar' },
] as const

// Language options
export const LANGUAGES = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'es', name: 'Español', flag: '🇪🇸' },
  { code: 'fr', name: 'Français', flag: '🇫🇷' },
  { code: 'de', name: 'Deutsch', flag: '🇩🇪' },
  { code: 'it', name: 'Italiano', flag: '🇮🇹' },
  { code: 'pt', name: 'Português', flag: '🇵🇹' },
] as const

// Default values
export const DEFAULTS = {
  CURRENCY: 'USD',
  LANGUAGE: 'en',
  THEME: 'dark',
  PRODUCTS_PER_PAGE: 12,
  PAGINATION_RANGE: 5,
  DEBOUNCE_DELAY: 300,
  TOAST_DURATION: 5000,
  MODAL_ANIMATION_DURATION: 200,
} as const

// Type helpers
export type Route = typeof ROUTES[keyof typeof ROUTES]
export type StorageKey = typeof STORAGE_KEYS[keyof typeof STORAGE_KEYS]
export type AnalyticsEvent = typeof ANALYTICS_EVENTS[keyof typeof ANALYTICS_EVENTS]
export type ToastType = typeof TOAST_TYPES[keyof typeof TOAST_TYPES]
export type OrderStatus = typeof ORDER_STATUS[keyof typeof ORDER_STATUS]
export type PaymentStatus = typeof PAYMENT_STATUS[keyof typeof PAYMENT_STATUS]
export type UserRole = typeof USER_ROLES[keyof typeof USER_ROLES]
export type ProductStatus = typeof PRODUCT_STATUS[keyof typeof PRODUCT_STATUS]
