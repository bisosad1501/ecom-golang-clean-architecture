// API Configuration Constants
export const API_CONFIG = {
  // Base URLs
  BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  API_VERSION: '/api/v1',
  
  // Endpoints
  ENDPOINTS: {
    // Auth
    AUTH: {
      LOGIN: '/auth/login',
      REGISTER: '/auth/register',
      LOGOUT: '/auth/logout',
      REFRESH: '/auth/refresh',
      PROFILE: '/auth/profile',
      FORGOT_PASSWORD: '/auth/forgot-password',
      RESET_PASSWORD: '/auth/reset-password',
      VERIFY_EMAIL: '/auth/verify-email',
    },
    
    // Products
    PRODUCTS: {
      LIST: '/products',
      DETAIL: '/products/:id',
      SEARCH: '/products/search',
      CATEGORIES: '/products/categories',
      FEATURED: '/products/featured',
      RELATED: '/products/:id/related',
    },
    
    // Categories
    CATEGORIES: {
      LIST: '/categories',
      DETAIL: '/categories/:id',
      PRODUCTS: '/categories/:id/products',
    },
    
    // Cart
    CART: {
      GET: '/cart',
      ADD: '/cart/add',
      UPDATE: '/cart/update',
      REMOVE: '/cart/remove',
      CLEAR: '/cart/clear',
    },
    
    // Orders
    ORDERS: {
      LIST: '/orders',
      DETAIL: '/orders/:id',
      CREATE: '/orders',
      UPDATE: '/orders/:id',
      CANCEL: '/orders/:id/cancel',
    },
    
    // Users
    USERS: {
      PROFILE: '/users/profile',
      UPDATE_PROFILE: '/users/profile',
      CHANGE_PASSWORD: '/users/change-password',
      ADDRESSES: '/users/addresses',
      WISHLIST: '/users/wishlist',
    },
    
    // Admin
    ADMIN: {
      DASHBOARD: '/admin/dashboard',
      PRODUCTS: '/admin/products',
      ORDERS: '/admin/orders',
      USERS: '/admin/users',
      CATEGORIES: '/admin/categories',
      ANALYTICS: '/admin/analytics',
    },
    
    // Payment
    PAYMENT: {
      STRIPE_INTENT: '/payment/stripe/create-intent',
      STRIPE_CONFIRM: '/payment/stripe/confirm',
      WEBHOOK: '/payment/webhook',
    },
  },
  
  // Request timeouts
  TIMEOUT: {
    DEFAULT: 10000,     // 10 seconds
    UPLOAD: 30000,      // 30 seconds
    DOWNLOAD: 60000,    // 60 seconds
  },
  
  // Retry configuration
  RETRY: {
    ATTEMPTS: 3,
    DELAY: 1000,        // 1 second
    BACKOFF: 2,         // Exponential backoff multiplier
  },
} as const

// HTTP Status Codes
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  CONFLICT: 409,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
  BAD_GATEWAY: 502,
  SERVICE_UNAVAILABLE: 503,
} as const

// Request headers
export const REQUEST_HEADERS = {
  CONTENT_TYPE: 'Content-Type',
  AUTHORIZATION: 'Authorization',
  ACCEPT: 'Accept',
  USER_AGENT: 'User-Agent',
} as const

// Content types
export const CONTENT_TYPES = {
  JSON: 'application/json',
  FORM_DATA: 'multipart/form-data',
  URL_ENCODED: 'application/x-www-form-urlencoded',
  TEXT: 'text/plain',
} as const

// Cache keys
export const CACHE_KEYS = {
  PRODUCTS: 'products',
  CATEGORIES: 'categories',
  USER_PROFILE: 'user-profile',
  CART: 'cart',
  WISHLIST: 'wishlist',
  ORDERS: 'orders',
} as const

// Cache durations (in milliseconds)
export const CACHE_DURATION = {
  SHORT: 5 * 60 * 1000,      // 5 minutes
  MEDIUM: 30 * 60 * 1000,    // 30 minutes
  LONG: 2 * 60 * 60 * 1000,  // 2 hours
  VERY_LONG: 24 * 60 * 60 * 1000, // 24 hours
} as const

// Query keys for React Query
export const QUERY_KEYS = {
  PRODUCTS: ['products'],
  PRODUCT_DETAIL: (id: string) => ['products', id],
  PRODUCT_SEARCH: (query: string) => ['products', 'search', query],
  CATEGORIES: ['categories'],
  CATEGORY_DETAIL: (id: string) => ['categories', id],
  CART: ['cart'],
  ORDERS: ['orders'],
  ORDER_DETAIL: (id: string) => ['orders', id],
  USER_PROFILE: ['user', 'profile'],
  WISHLIST: ['user', 'wishlist'],
  ADMIN_DASHBOARD: ['admin', 'dashboard'],
  ADMIN_PRODUCTS: ['admin', 'products'],
  ADMIN_ORDERS: ['admin', 'orders'],
  ADMIN_USERS: ['admin', 'users'],
} as const

// Error messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection.',
  UNAUTHORIZED: 'You are not authorized to perform this action.',
  FORBIDDEN: 'Access denied.',
  NOT_FOUND: 'The requested resource was not found.',
  VALIDATION_ERROR: 'Please check your input and try again.',
  SERVER_ERROR: 'Something went wrong. Please try again later.',
  TIMEOUT: 'Request timed out. Please try again.',
  RATE_LIMIT: 'Too many requests. Please try again later.',
} as const

// Success messages
export const SUCCESS_MESSAGES = {
  LOGIN: 'Successfully logged in!',
  LOGOUT: 'Successfully logged out!',
  REGISTER: 'Account created successfully!',
  PROFILE_UPDATED: 'Profile updated successfully!',
  PASSWORD_CHANGED: 'Password changed successfully!',
  PRODUCT_ADDED_TO_CART: 'Product added to cart!',
  PRODUCT_REMOVED_FROM_CART: 'Product removed from cart!',
  ORDER_PLACED: 'Order placed successfully!',
  ORDER_CANCELLED: 'Order cancelled successfully!',
} as const

// Pagination defaults (synchronized with backend)
export const PAGINATION = {
  DEFAULT_PAGE: 1,
  DEFAULT_LIMIT: 20,
  MAX_LIMIT: 100,
  MIN_LIMIT: 1,

  // Entity-specific limits (matching backend)
  PRODUCTS_PER_PAGE: 12,
  ORDERS_PER_PAGE: 10,
  REVIEWS_PER_PAGE: 5,
  NOTIFICATIONS_PER_PAGE: 15,
  SEARCH_RESULTS_PER_PAGE: 20,
  WISHLIST_PER_PAGE: 12,
  ADMIN_USERS_PER_PAGE: 25,
  ADMIN_ORDERS_PER_PAGE: 20,

  // Page size options for UI
  PRODUCT_PAGE_SIZES: [12, 24, 48, 96],
  ORDER_PAGE_SIZES: [10, 20, 50],
  REVIEW_PAGE_SIZES: [5, 10, 20],
  SEARCH_PAGE_SIZES: [20, 40, 60],
  DEFAULT_PAGE_SIZES: [10, 20, 50, 100],

  // Performance limits
  MAX_SEARCH_RESULTS: 1000,
  MAX_ORDER_HISTORY: 500,
} as const

// File upload limits
export const UPLOAD_LIMITS = {
  MAX_FILE_SIZE: 5 * 1024 * 1024, // 5MB
  ALLOWED_IMAGE_TYPES: ['image/jpeg', 'image/png', 'image/webp'],
  ALLOWED_DOCUMENT_TYPES: ['application/pdf', 'text/plain'],
  MAX_FILES: 10,
} as const

// Type helpers
export type ApiEndpoint = typeof API_CONFIG.ENDPOINTS
export type HttpStatus = typeof HTTP_STATUS[keyof typeof HTTP_STATUS]
export type CacheKey = typeof CACHE_KEYS[keyof typeof CACHE_KEYS]
