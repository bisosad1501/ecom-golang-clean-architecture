// ===== ORDER STORE =====

import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { api } from '@/lib/api'
import type { Order, Address } from '@/types'

interface OrderState {
  // State
  currentOrder: Order | null
  orders: Order[]
  isLoading: boolean
  error: string | null
  
  // Checkout state
  shippingAddress: Address | null
  billingAddress: Address | null
  paymentMethod: string | null
  
  // Actions
  createOrder: (orderData: any) => Promise<Order>
  getOrder: (orderId: string) => Promise<Order>
  getOrders: () => Promise<Order[]>
  updateOrder: (orderId: string, updates: Partial<Order>) => Promise<Order>
  cancelOrder: (orderId: string) => Promise<void>
  
  // Checkout actions
  setShippingAddress: (address: Address) => void
  setBillingAddress: (address: Address) => void
  setPaymentMethod: (method: string) => void
  clearCheckoutData: () => void
  
  // Utility actions
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  clearError: () => void
}

export const useOrderStore = create<OrderState>()(
  persist(
    (set, get) => ({
      // Initial state
      currentOrder: null,
      orders: [],
      isLoading: false,
      error: null,
      shippingAddress: null,
      billingAddress: null,
      paymentMethod: null,

      // Order actions
      createOrder: async (orderData) => {
        set({ isLoading: true, error: null })
        try {
          // This would be replaced with actual API call
          const order: Order = {
            id: Math.random().toString(36).substr(2, 9),
            order_number: `ORD-${Date.now()}`,
            user_id: 'user-1',
            status: 'pending',
            payment_status: 'pending',
            items: orderData.items || [],
            shipping_address: orderData.shippingAddress,
            billing_address: orderData.billingAddress,
            subtotal: orderData.subtotal || 0,
            tax_amount: orderData.taxAmount || 0,
            shipping_amount: orderData.shippingAmount || 0,
            discount_amount: orderData.discountAmount || 0,
            total_amount: orderData.totalAmount || 0,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          }
          
          set({ 
            currentOrder: order,
            orders: [...get().orders, order],
            isLoading: false 
          })
          
          return order
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Failed to create order'
          set({ error: errorMessage, isLoading: false })
          throw error
        }
      },

      getOrder: async (orderId) => {
        set({ isLoading: true, error: null })
        try {
          // This would be replaced with actual API call
          const existingOrder = get().orders.find(order => order.id === orderId)
          if (existingOrder) {
            set({ currentOrder: existingOrder, isLoading: false })
            return existingOrder
          }
          
          // Mock API call
          const order: Order = {
            id: orderId,
            order_number: `ORD-${Date.now()}`,
            user_id: 'user-1',
            status: 'pending',
            payment_status: 'pending',
            items: [],
            shipping_address: {
              id: '1',
              type: 'shipping',
              first_name: 'John',
              last_name: 'Doe',
              address_line_1: '123 Main St',
              city: 'New York',
              state: 'NY',
              postal_code: '10001',
              country: 'US',
              is_default: true,
            },
            billing_address: {
              id: '2',
              type: 'billing',
              first_name: 'John',
              last_name: 'Doe',
              address_line_1: '123 Main St',
              city: 'New York',
              state: 'NY',
              postal_code: '10001',
              country: 'US',
              is_default: true,
            },
            subtotal: 0,
            tax_amount: 0,
            shipping_amount: 0,
            discount_amount: 0,
            total_amount: 0,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString(),
          }
          
          set({ 
            currentOrder: order,
            orders: [...get().orders, order],
            isLoading: false 
          })
          
          return order
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Failed to get order'
          set({ error: errorMessage, isLoading: false })
          throw error
        }
      },

      getOrders: async () => {
        set({ isLoading: true, error: null })
        try {
          // This would be replaced with actual API call
          const orders: Order[] = []
          
          set({ orders, isLoading: false })
          return orders
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Failed to get orders'
          set({ error: errorMessage, isLoading: false })
          throw error
        }
      },

      updateOrder: async (orderId, updates) => {
        set({ isLoading: true, error: null })
        try {
          // This would be replaced with actual API call
          const orders = get().orders
          const orderIndex = orders.findIndex(order => order.id === orderId)
          
          if (orderIndex === -1) {
            throw new Error('Order not found')
          }
          
          const updatedOrder = {
            ...orders[orderIndex],
            ...updates,
            updated_at: new Date().toISOString(),
          }
          
          const newOrders = [...orders]
          newOrders[orderIndex] = updatedOrder
          
          set({ 
            orders: newOrders,
            currentOrder: get().currentOrder?.id === orderId ? updatedOrder : get().currentOrder,
            isLoading: false 
          })
          
          return updatedOrder
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Failed to update order'
          set({ error: errorMessage, isLoading: false })
          throw error
        }
      },

      cancelOrder: async (orderId) => {
        set({ isLoading: true, error: null })
        try {
          await get().updateOrder(orderId, { status: 'cancelled' })
          set({ isLoading: false })
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Failed to cancel order'
          set({ error: errorMessage, isLoading: false })
          throw error
        }
      },

      // Checkout actions
      setShippingAddress: (address) => {
        set({ shippingAddress: address })
      },

      setBillingAddress: (address) => {
        set({ billingAddress: address })
      },

      setPaymentMethod: (method) => {
        set({ paymentMethod: method })
      },

      clearCheckoutData: () => {
        set({
          shippingAddress: null,
          billingAddress: null,
          paymentMethod: null,
        })
      },

      // Utility actions
      setLoading: (loading) => {
        set({ isLoading: loading })
      },

      setError: (error) => {
        set({ error })
      },

      clearError: () => {
        set({ error: null })
      },
    }),
    {
      name: 'order-store',
      partialize: (state) => ({
        orders: state.orders,
        shippingAddress: state.shippingAddress,
        billingAddress: state.billingAddress,
        paymentMethod: state.paymentMethod,
      }),
    }
  )
)
