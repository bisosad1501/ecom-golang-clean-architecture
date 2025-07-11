import { create } from 'zustand'
import { devtools } from 'zustand/middleware'
import { apiClient } from '@/lib/api'
import { 
  PaymentStore, 
  PaymentMethod, 
  Payment, 
  CheckoutSession,
  CreateCheckoutSessionRequest,
  ProcessPaymentRequest,
  RefundRequest 
} from '@/types'
import toast from 'react-hot-toast'

export const usePaymentStore = create<PaymentStore>()(
  devtools(
    (set, get) => ({
      paymentMethods: [],
      currentPayment: null,
      isLoading: false,
      error: null,

      // Fetch user's payment methods
      fetchPaymentMethods: async () => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.getPaymentMethods()
          set({ paymentMethods: response.data || [], isLoading: false })
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to fetch payment methods'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
        }
      },

      // Save a new payment method
      savePaymentMethod: async (data: any) => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.savePaymentMethod(data)
          
          // Refresh payment methods list
          await get().fetchPaymentMethods()
          
          toast.success('Payment method saved successfully')
          set({ isLoading: false })
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to save payment method'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },

      // Delete a payment method
      deletePaymentMethod: async (methodId: string) => {
        set({ isLoading: true, error: null })
        try {
          await apiClient.deletePaymentMethod(methodId)
          
          // Remove from local state
          const currentMethods = get().paymentMethods
          const updatedMethods = currentMethods.filter(method => method.id !== methodId)
          set({ paymentMethods: updatedMethods, isLoading: false })
          
          toast.success('Payment method deleted successfully')
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to delete payment method'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },

      // Set default payment method
      setDefaultPaymentMethod: async (methodId: string) => {
        set({ isLoading: true, error: null })
        try {
          await apiClient.setDefaultPaymentMethod(methodId)
          
          // Update local state
          const currentMethods = get().paymentMethods
          const updatedMethods = currentMethods.map(method => ({
            ...method,
            is_default: method.id === methodId
          }))
          set({ paymentMethods: updatedMethods, isLoading: false })
          
          toast.success('Default payment method updated')
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to update default payment method'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },

      // Create Stripe checkout session
      createCheckoutSession: async (data: CreateCheckoutSessionRequest): Promise<CheckoutSession> => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.post<CheckoutSession>('/payments/checkout-session', data)
          set({ isLoading: false })

          const checkoutData = response.data?.data || response.data

          console.log('Checkout session response:', checkoutData)

          // Check if we have session data
          if (checkoutData.session_id || checkoutData.session_url) {
            toast.success('Redirecting to payment...')
            return {
              id: checkoutData.session_id,
              url: checkoutData.session_url,
              expires_at: '', // Not provided by backend
            }
          } else {
            throw new Error(checkoutData.message || 'Failed to create checkout session')
          }
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || error.message || 'Failed to create checkout session'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },

      // Process direct payment
      processPayment: async (data: ProcessPaymentRequest) => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.processPayment(data)
          set({ currentPayment: response.data, isLoading: false })
          toast.success('Payment processed successfully')
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Payment failed'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },

      // Get payment details
      getPayment: async (paymentId: string) => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.getPayment(paymentId)
          set({ currentPayment: response.data, isLoading: false })
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to fetch payment details'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
        }
      },

      // Process refund
      processRefund: async (paymentId: string, data: RefundRequest) => {
        set({ isLoading: true, error: null })
        try {
          const response = await apiClient.processRefund(paymentId, data)
          
          // Update current payment if it's the same one
          const currentPayment = get().currentPayment
          if (currentPayment && currentPayment.id === paymentId) {
            set({ currentPayment: { ...currentPayment, status: 'refunded' } })
          }
          
          set({ isLoading: false })
          toast.success('Refund processed successfully')
        } catch (error: any) {
          const errorMessage = error.response?.data?.message || 'Failed to process refund'
          set({ error: errorMessage, isLoading: false })
          toast.error(errorMessage)
          throw error
        }
      },
    }),
    {
      name: 'payment-store',
    }
  )
)
