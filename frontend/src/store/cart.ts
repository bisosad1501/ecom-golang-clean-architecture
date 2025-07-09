import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient } from '@/lib/api'
import { Cart, CartItem, Product, ApiResponse } from '@/types'

interface CartState {
  cart: Cart | null
  isLoading: boolean
  error: string | null
  isOpen: boolean
}

interface CartActions {
  fetchCart: () => Promise<void>
  addItem: (productId: string, quantity?: number) => Promise<void>
  updateItem: (itemId: string, quantity: number) => Promise<void>
  removeItem: (itemId: string) => Promise<void>
  clearCart: () => Promise<void>
  clearCartLocal: () => void  // Clear cart from local state only (for logout)
  openCart: () => void
  closeCart: () => void
  toggleCart: () => void
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
      isOpen: false,

      // Actions
      fetchCart: async () => {
        try {
          set({ isLoading: true, error: null })
          
          // Check if user is authenticated
          const { useAuthStore } = await import('@/store/auth')
          const { isAuthenticated } = useAuthStore.getState()
          
          if (!isAuthenticated) {
            // Clear cart if not authenticated
            set({
              cart: null,
              isLoading: false,
              error: null,
            })
            return
          }

          const response = await apiClient.get<ApiResponse<Cart>>('/cart')

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

          const response = await apiClient.post<ApiResponse<Cart>>('/cart/items', {
            product_id: productId,
            quantity,
          })

          // Handle response format - check if data is nested
          const cart = response.data?.data || response.data

          set({
            cart,
            isLoading: false,
            error: null,
            isOpen: true, // Open cart when item is added
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Failed to add item to cart',
          })
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
          isOpen: false,
          error: null,
        })
      },

      openCart: () => set({ isOpen: true }),
      closeCart: () => set({ isOpen: false }),
      toggleCart: () => set((state) => ({ isOpen: !state.isOpen })),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({
        // Only persist cart if we have items and isOpen state
        // Don't persist when cart is null (i.e., user logged out)
        cart: state.cart,
        isOpen: state.isOpen,
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
  return cart?.items?.length || 0
}

export const getCartTotal = (cart: Cart | null): number => {
  return cart?.total || 0
}

export const isProductInCart = (cart: Cart | null, productId: string): boolean => {
  return cart?.items.some(item => item.product_id === productId) || false
}

export const getCartItemQuantity = (cart: Cart | null, productId: string): number => {
  const item = cart?.items.find(item => item.product_id === productId)
  return item?.quantity || 0
}
