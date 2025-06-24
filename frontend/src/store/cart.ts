import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient } from '@/lib/api'
import { Cart, CartItem, Product } from '@/types'

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
          
          const response = await apiClient.get<Cart>('/cart')
          const cart = response.data

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
          
          const response = await apiClient.post<Cart>('/cart/items', {
            product_id: productId,
            quantity,
          })
          const cart = response.data

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

          const response = await apiClient.put<Cart>(`/cart/items/${itemId}`, {
            quantity,
          })
          const cart = response.data

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
          
          const response = await apiClient.delete<Cart>(`/cart/items/${itemId}`)
          const cart = response.data

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
          
          await apiClient.delete('/cart')

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

      openCart: () => set({ isOpen: true }),
      closeCart: () => set({ isOpen: false }),
      toggleCart: () => set((state) => ({ isOpen: !state.isOpen })),
      clearError: () => set({ error: null }),
    }),
    {
      name: 'cart-storage',
      partialize: (state) => ({
        cart: state.cart,
      }),
    }
  )
)

// Helper functions
export const getCartItemCount = (cart: Cart | null): number => {
  return cart?.item_count || 0
}

export const getCartTotal = (cart: Cart | null): number => {
  return cart?.total_amount || 0
}

export const isProductInCart = (cart: Cart | null, productId: string): boolean => {
  return cart?.items.some(item => item.product_id === productId) || false
}

export const getCartItemQuantity = (cart: Cart | null, productId: string): number => {
  const item = cart?.items.find(item => item.product_id === productId)
  return item?.quantity || 0
}
