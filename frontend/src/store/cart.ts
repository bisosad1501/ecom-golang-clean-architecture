import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient, cartApi } from '@/lib/api'
import { Cart, CartItem, Product, ApiResponse, CartConflictInfo, MergeStrategy } from '@/types'
import { toast } from 'sonner'

// Generate or get session ID for guest users
const getSessionId = (): string => {
  if (typeof window === 'undefined') return ''

  let sessionId = localStorage.getItem('guest_session_id')
  if (!sessionId) {
    sessionId = `guest-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    localStorage.setItem('guest_session_id', sessionId)
  }
  return sessionId
}

// Clear guest session when user logs in
const clearGuestSession = () => {
  if (typeof window !== 'undefined') {
    localStorage.removeItem('guest_session_id')
  }
}

// Types are now imported from @/types

interface CartState {
  cart: Cart | null
  isLoading: boolean
  error: string | null
}

interface CartActions {
  fetchCart: () => Promise<void>
  addItem: (productId: string, quantity?: number) => Promise<void>
  updateItem: (itemId: string, quantity: number) => Promise<void>
  removeItem: (itemId: string) => Promise<void>
  clearCart: () => Promise<void>
  clearCartLocal: () => void  // Clear cart from local state only (for logout)
  mergeGuestCart: (strategy?: MergeStrategy) => Promise<void>  // Merge guest cart when user logs in
  checkMergeConflict: () => Promise<CartConflictInfo | null>  // Check merge conflicts
  clearError: () => void
}

type CartStore = CartState & CartActions

export const useCartStore = create<CartStore>()(
  persist(
    (set, get) => ({
      // Initial state
      cart: null,
      isLoading: false,
      error: null,

      // Actions
      fetchCart: async () => {
        try {
          set({ isLoading: true, error: null })

          // Check if user is authenticated
          const { useAuthStore } = await import('@/store/auth')
          const { isAuthenticated } = useAuthStore.getState()

          let response: any

          if (isAuthenticated) {
            // Authenticated user - use protected endpoint
            response = await apiClient.get<ApiResponse<Cart>>('/cart')
          } else {
            // Guest user - use public endpoint with session ID
            const sessionId = getSessionId()
            response = await apiClient.get<ApiResponse<Cart>>('/public/cart', {
              headers: {
                'X-Session-ID': sessionId
              }
            })
          }

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to fetch cart',
          })
        }
      },

      addItem: async (productId: string, quantity = 1) => {
        try {
          set({ isLoading: true, error: null })

          // Validate input
          if (!productId || quantity <= 0) {
            throw new Error('Invalid product ID or quantity')
          }

          // Check if user is authenticated
          const { useAuthStore } = await import('@/store/auth')
          const { isAuthenticated } = useAuthStore.getState()

          let response: any

          if (isAuthenticated) {
            // Authenticated user - use protected endpoint
            response = await apiClient.post<ApiResponse<Cart>>('/cart/items', {
              product_id: productId,
              quantity,
            })
          } else {
            // Guest user - use public endpoint with session ID
            const sessionId = getSessionId()
            response = await apiClient.post<ApiResponse<Cart>>('/public/cart/items', {
              product_id: productId,
              quantity,
            }, {
              headers: {
                'X-Session-ID': sessionId
              }
            })
          }

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
          })

          toast.success('Item added to cart')
        } catch (error: any) {
          console.error('Failed to add item to cart:', error)

          // Handle specific error types
          let errorMessage = 'Failed to add item to cart'
          if (error.code === 'INSUFFICIENT_STOCK') {
            errorMessage = 'Not enough stock available'
          } else if (error.code === 'PRODUCT_NOT_FOUND') {
            errorMessage = 'Product not found'
          } else if (error.code === 'VALIDATION_ERROR') {
            errorMessage = error.details?.join(', ') || 'Invalid input'
          } else if (error.code === 'RATE_LIMITED') {
            errorMessage = 'Too many requests. Please wait a moment.'
          } else if (error.message) {
            errorMessage = error.message
          }

          set({
            isLoading: false,
            error: errorMessage,
          })
          toast.error(errorMessage)
          throw error
        }
      },

      updateItem: async (itemId: string, quantity: number) => {
        try {
          set({ isLoading: true, error: null })

          if (quantity <= 0) {
            await get().removeItem(itemId)
            return
          }

          // Find the item to get the product_id
          const currentCart = get().cart
          const item = currentCart?.items.find(item => item.id === itemId)
          if (!item) {
            throw new Error('Item not found in cart')
          }

          const productId = item.product_id || item.product?.id
          if (!productId) {
            throw new Error('Product ID not found')
          }

          const response = await apiClient.put<ApiResponse<Cart>>('/cart/items', {
            product_id: productId,
            quantity,
          })

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to update cart item',
          })
          throw error
        }
      },

      removeItem: async (itemId: string) => {
        try {
          set({ isLoading: true, error: null })

          // Find the item to get the product_id
          const currentCart = get().cart
          const item = currentCart?.items.find(item => item.id === itemId)
          if (!item) {
            throw new Error('Item not found in cart')
          }

          const productId = item.product_id || item.product?.id
          if (!productId) {
            throw new Error('Product ID not found')
          }

          const response = await apiClient.delete<ApiResponse<Cart>>(`/cart/items/${productId}`)

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to remove cart item',
          })
          throw error
        }
      },

      clearCart: async () => {
        try {
          set({ isLoading: true, error: null })
          
          // Only call API if user is authenticated
          const { useAuthStore } = await import('@/store/auth')
          const { isAuthenticated } = useAuthStore.getState()
          
          if (isAuthenticated) {
            await apiClient.delete('/cart')
          }

          set({
            cart: null,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to clear cart',
          })
          throw error
        }
      },

      clearCartLocal: () => {
        // Clear cart from local state only (used for logout)
        set({
          cart: null,
          error: null,
        })
      },

      mergeGuestCart: async (strategy: MergeStrategy = 'auto') => {
        try {
          const sessionId = localStorage.getItem('guest_session_id')
          if (!sessionId) {
            // No guest cart to merge
            return
          }

          set({ isLoading: true, error: null })

          // Call merge API with strategy
          const response = await apiClient.post<ApiResponse<Cart>>('/cart/merge', {
            session_id: sessionId,
            strategy: strategy
          })

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
          })

          // Clear guest session after successful merge
          clearGuestSession()

          // Don't show toast here - let the calling component handle it
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to merge guest cart',
          })
          // Don't throw error for merge failure - just log it
          console.error('Failed to merge guest cart:', error)
        }
      },

      checkMergeConflict: async (): Promise<CartConflictInfo | null> => {
        try {
          const sessionId = localStorage.getItem('guest_session_id')
          if (!sessionId) {
            return null
          }

          const response = await apiClient.post<ApiResponse<CartConflictInfo>>('/cart/check-conflict', {
            session_id: sessionId
          })

          return response.data?.data || response.data
        } catch (error: any) {
          console.error('Failed to check merge conflict:', error)
          return null
        }
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({
        // Only persist cart if we have items
        // Don't persist when cart is null (i.e., user logged out)
        cart: state.cart,
      }),
      onRehydrateStorage: () => (state) => {
        // After rehydration, validate cart belongs to current user
        if (state?.cart) {
          import('@/store/auth').then(({ useAuthStore }) => {
            const { isAuthenticated } = useAuthStore.getState()
            if (!isAuthenticated) {
              // Clear cart if user is not authenticated
              state.clearCartLocal()
            } else {
              // Fetch fresh cart from server to ensure consistency
              state.fetchCart().catch(console.error)
            }
          })
        }
      },
    }
  )
)

// Helper functions
export const getCartItemCount = (cart: Cart | null): number => {
  // Use calculated field from backend if available, fallback to items length
  return cart?.item_count ?? cart?.items?.length ?? 0
}

export const getCartTotal = (cart: Cart | null): number => {
  // Use calculated field from backend if available
  return cart?.total ?? 0
}

export const getCartSubtotal = (cart: Cart | null): number => {
  // Use calculated field from backend if available
  return cart?.subtotal ?? 0
}

export const getCartTaxAmount = (cart: Cart | null): number => {
  // Use calculated field from backend
  return cart?.tax_amount ?? 0
}

export const getCartShippingAmount = (cart: Cart | null): number => {
  // Use calculated field from backend
  return cart?.shipping_amount ?? 0
}

export const isProductInCart = (cart: Cart | null, productId: string): boolean => {
  return cart?.items.some(item => item.product_id === productId || item.product?.id === productId) || false
}

export const getCartItemQuantity = (cart: Cart | null, productId: string): number => {
  const item = cart?.items.find(item => item.product_id === productId || item.product?.id === productId)
  return item?.quantity || 0
}

export const isGuestCart = (cart: Cart | null): boolean => {
  return cart?.session_id != null && cart?.user_id === '00000000-0000-0000-0000-000000000000'
}
