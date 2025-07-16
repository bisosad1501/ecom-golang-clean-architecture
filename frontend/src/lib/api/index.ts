// ===== API SERVICES EXPORT =====

// Export API client
export { apiClient, ApiClient, ApiError } from './client'
export type { RequestConfig } from './client'

// Export services
export { authService, AuthService } from './auth'
// Note: ProductsService removed - using services/products.ts instead

// Re-export commonly used services for convenience
export { authService as auth } from './auth'
// Note: products service removed - using services/products.ts instead

// Service factory for creating custom instances
// Note: Commented out to avoid TypeScript errors since this is not actively used
// export class ServiceFactory {
//   private baseURL: string
//   private timeout: number
//   private retries: number

//   constructor(options: {
//     baseURL?: string
//     timeout?: number
//     retries?: number
//   } = {}) {
//     this.baseURL = options.baseURL || process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'
//     this.timeout = options.timeout || 10000
//     this.retries = options.retries || 3
//   }

//   createAuthService(): AuthService {
//     return new AuthService()
//   }

//   createProductsService(): ProductsService {
//     return new ProductsService()
//   }

//   createApiClient(): ApiClient {
//     return new ApiClient(this.baseURL, {
//       timeout: this.timeout,
//       retries: this.retries,
//     })
//   }
// }

// Default service factory
// export const serviceFactory = new ServiceFactory()

// Utility functions for common API operations
// Note: These are commented out to avoid import errors since the services are not actively used
// export const api = {
//   // Authentication shortcuts
//   login: authService.login.bind(authService),
//   logout: authService.logout.bind(authService),
//   register: authService.register.bind(authService),
//   getProfile: authService.getProfile.bind(authService),
//
//   // Products shortcuts
//   getProducts: productsService.getProducts.bind(productsService),
//   getProduct: productsService.getProduct.bind(productsService),
//   searchProducts: productsService.searchProducts.bind(productsService),
//   getFeaturedProducts: productsService.getFeaturedProducts.bind(productsService),
//
//   // Wishlist shortcuts
//   getWishlist: productsService.getWishlist.bind(productsService),
//   addToWishlist: productsService.addToWishlist.bind(productsService),
//   removeFromWishlist: productsService.removeFromWishlist.bind(productsService),
//
//   // Categories shortcuts
//   getCategories: productsService.getCategories.bind(productsService),
//   getCategory: productsService.getCategory.bind(productsService),
// }

// Error handling utilities
// Note: Commented out to avoid TypeScript errors since this is not actively used
// export const handleApiError = (error: any) => {
//   if (error instanceof ApiError) {
//     switch (error.status) {
//       case 401:
//         // Unauthorized - redirect to login
//         if (typeof window !== 'undefined') {
//           authService.logout()
//           window.location.href = '/auth/login'
//         }
//         break
//       case 403:
//         // Forbidden - show access denied message
//         console.error('Access denied:', error.message)
//         break
//       case 404:
//         // Not found
//         console.error('Resource not found:', error.message)
//         break
//       case 422:
//         // Validation error
//         console.error('Validation error:', error.details)
//         break
//       case 500:
//         // Server error
//         console.error('Server error:', error.message)
//         break
//       default:
//         console.error('API error:', error.message)
//     }
//     return error
//   }
//
//   // Network or other errors
//   console.error('Network error:', error)
//   return new ApiError('Network error occurred', 0, 'NETWORK_ERROR')
// }

// Request interceptor for global error handling
export const setupGlobalErrorHandling = () => {
  // This would be implemented with the specific HTTP client being used
  // For now, it's a placeholder for future implementation
  console.log('Global error handling setup')
}

// Response interceptor for token refresh
export const setupTokenRefresh = () => {
  // This would automatically refresh tokens when they expire
  // Implementation depends on the specific requirements
  console.log('Token refresh setup')
}

// Health check utility
// Note: Commented out to avoid TypeScript errors since this is not actively used
// export const checkApiHealth = async (): Promise<boolean> => {
//   try {
//     await apiClient.healthCheck()
//     return true
//   } catch {
//     return false
//   }
// }

// API status utility
// export const getApiStatus = async (): Promise<{
//   healthy: boolean
//   latency: number
//   timestamp: string
// }> => {
//   const start = Date.now()
//
//   try {
//     const response = await apiClient.healthCheck()
//     const latency = Date.now() - start
//
//     return {
//       healthy: true,
//       latency,
//       timestamp: response.timestamp,
//     }
//   } catch {
//     return {
//       healthy: false,
//       latency: Date.now() - start,
//       timestamp: new Date().toISOString(),
//     }
//   }
// }

// Cache utilities (for future implementation with React Query or SWR)
export const cacheKeys = {
  // Auth
  profile: ['auth', 'profile'],
  sessions: ['auth', 'sessions'],
  
  // Products
  products: (params?: any) => ['products', 'list', params],
  product: (id: string) => ['products', 'detail', id],
  featuredProducts: ['products', 'featured'],
  relatedProducts: (id: string) => ['products', 'related', id],
  productReviews: (id: string, params?: any) => ['products', id, 'reviews', params],
  productQuestions: (id: string, params?: any) => ['products', id, 'questions', params],
  
  // Categories
  categories: ['categories', 'list'],
  category: (id: string) => ['categories', 'detail', id],
  
  // Wishlist
  wishlist: ['user', 'wishlist'],
  
  // Search
  searchSuggestions: (query: string) => ['search', 'suggestions', query],
  productFilters: (categoryId?: string) => ['products', 'filters', categoryId],
}

// Export cache keys for use with React Query or SWR
export { cacheKeys as queryKeys }

// Type exports for better TypeScript support
export type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
  Product,
  ProductListItem,
  ProductSearchParams,
  Category,
  Brand,
  WishlistItem,
} from '@/types'
