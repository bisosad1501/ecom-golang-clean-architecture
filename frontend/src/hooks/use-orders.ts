import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import apiClient from '@/lib/api'
import { Order, PaginatedResponse, CreateOrderRequest } from '@/types'

// Query keys
export const orderKeys = {
  all: ['orders'] as const,
  lists: () => [...orderKeys.all, 'list'] as const,
  list: (filters: Record<string, any>) => [...orderKeys.lists(), filters] as const,
  details: () => [...orderKeys.all, 'detail'] as const,
  detail: (id: string) => [...orderKeys.details(), id] as const,
  user: (userId: string) => [...orderKeys.all, 'user', userId] as const,
  admin: () => [...orderKeys.all, 'admin'] as const,
}

// Get orders (user's own orders)
export function useOrders(params: {
  page?: number
  limit?: number
  search?: string
  status?: string
  user_id?: string
} = {}) {
  return useQuery({
    queryKey: orderKeys.list(params),
    queryFn: async (): Promise<PaginatedResponse<Order>> => {
      try {
        const queryParams = new URLSearchParams()
        
        // Convert page to offset for backend compatibility
        const limit = params.limit || 10
        const page = params.page || 1
        const offset = (page - 1) * limit
        
        queryParams.append('limit', limit.toString())
        queryParams.append('offset', offset.toString())
        if (params.search) queryParams.append('search', params.search)
        if (params.status) queryParams.append('status', params.status)
        if (params.user_id) queryParams.append('user_id', params.user_id)
        
        const url = `/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
        console.log('useOrders - Making API call to:', url)
        console.log('useOrders - Query params (original):', params)
        console.log('useOrders - Query params (converted):', { limit, offset, page })

        const response = await apiClient.get<any>(url)
        console.log('useOrders - Raw API response:', response)
        console.log('useOrders - Response data:', response.data)
        
        // Handle the SuccessResponse structure from backend
        let ordersArray = response.data
        
        // If the response is wrapped in SuccessResponse structure
        if (response.data && response.data.data) {
          console.log('useOrders - Detected SuccessResponse wrapper, unwrapping...')
          ordersArray = response.data.data
        }
        
        // Ensure we have an array
        if (!Array.isArray(ordersArray)) {
          console.log('useOrders - Response is not an array, using empty array')
          ordersArray = []
        }
        
        console.log('useOrders - Extracted orders array:', ordersArray)
        console.log('useOrders - Orders count:', ordersArray.length)
        
        // Create paginated response structure
        const responseData = {
          data: ordersArray,
          pagination: {
            page: page,
            limit: limit,
            total: ordersArray.length, // This is not accurate but we don't have total from backend
            total_pages: Math.ceil(ordersArray.length / limit) || 1,
            has_next: ordersArray.length === limit, // Assume there might be more if we got a full page
            has_prev: page > 1
          }
        }
        
        console.log('useOrders - Final processed data:', responseData)
        return responseData
      } catch (error) {
        console.error('useOrders - API call failed:', error)
        throw error
      }
    },
    staleTime: 30 * 1000, // 30 seconds
  })
}

// Get admin orders (all orders)
export function useAdminOrders(params: {
  page?: number
  limit?: number
  search?: string
  status?: string
  user_id?: string
  date_from?: string
  date_to?: string
} = {}) {
  console.log('useAdminOrders - Hook called with params:', params)

  return useQuery({
    queryKey: ['admin-orders', params],
    queryFn: async (): Promise<PaginatedResponse<Order>> => {
      console.log('useAdminOrders - QueryFn executing...')

      const queryParams = new URLSearchParams()

      // Convert page to offset for backend
      const offset = params.page ? (params.page - 1) * (params.limit || 20) : 0
      queryParams.append('offset', offset.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.status) queryParams.append('status', params.status)
      if (params.user_id) queryParams.append('user_id', params.user_id)
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)

      const url = `/admin/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useAdminOrders - Making API call to:', url)

      const response = await apiClient.get<any>(url)
      console.log('useAdminOrders - Raw response:', response)

      // Handle SuccessResponse wrapper
      let ordersData = response.data
      if (response.data && response.data.data) {
        console.log('useAdminOrders - Unwrapping SuccessResponse')
        ordersData = response.data.data
      }

      // Transform backend response to frontend format
      const result = {
        data: ordersData.orders || [],
        pagination: ordersData.pagination || {
          current_page: 1,
          per_page: 20,
          total_pages: 1,
          total_items: 0,
          has_next: false,
          has_prev: false
        }
      }

      console.log('useAdminOrders - Final result:', result)
      return result
    },
    staleTime: 30 * 1000,
    retry: 1,
  })
}

// Get single order
export function useOrder(id: string) {
  return useQuery({
    queryKey: orderKeys.detail(id),
    queryFn: async (): Promise<Order> => {
      console.log('useOrder - Fetching order:', id)

      // Try public endpoint first (no auth required)
      try {
        const response = await apiClient.get<any>(`/orders/${id}/public`)
        console.log('useOrder - Public endpoint response:', response)

        // Handle SuccessResponse wrapper
        let orderData = response.data
        if (response.data && response.data.data) {
          console.log('useOrder - Unwrapping SuccessResponse')
          orderData = response.data.data
        }

        console.log('useOrder - Final order data:', orderData)
        return orderData
      } catch (publicError) {
        console.log('useOrder - Public endpoint failed, trying authenticated endpoint:', publicError)

        // Fallback to authenticated endpoint
        const response = await apiClient.get<any>(`/orders/${id}`)
        console.log('useOrder - Authenticated endpoint response:', response)

        // Handle SuccessResponse wrapper
        let orderData = response.data
        if (response.data && response.data.data) {
          console.log('useOrder - Unwrapping SuccessResponse')
          orderData = response.data.data
        }

        console.log('useOrder - Final order data:', orderData)
        return orderData
      }
    },
    enabled: !!id,
    staleTime: 30 * 1000,
  })
}

// Get user orders
export function useUserOrders(userId?: string) {
  return useQuery({
    queryKey: orderKeys.user(userId || ''),
    queryFn: async (): Promise<Order[]> => {
      const response = await apiClient.get<Order[]>('/orders')
      return response.data
    },
    enabled: !!userId,
    staleTime: 30 * 1000,
  })
}

// Create order
export function useCreateOrder() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: CreateOrderRequest): Promise<Order> => {
      const response = await apiClient.post<Order>('/orders', data)
      return response.data
    },
    onSuccess: (data) => {
      // Invalidate and refetch orders
      queryClient.invalidateQueries({ queryKey: orderKeys.all })
      toast.success('Order created successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create order')
    },
  })
}

// Update order status (admin only)
export function useUpdateOrderStatus() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, status, notes }: {
      id: string
      status: string
      notes?: string
    }): Promise<any> => {
      const response = await apiClient.patch(`/admin/orders/${id}/status`, {
        status,
        notes
      })
      return response.data
    },
    onSuccess: (data, variables) => {
      console.log('useUpdateOrderStatus - onSuccess called, invalidating queries...')
      // Invalidate all order queries to refetch with updated data
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() })
      queryClient.invalidateQueries({ queryKey: orderKeys.admin() })
      queryClient.invalidateQueries({ queryKey: ['admin-orders'] })  // Add specific admin orders key
      queryClient.invalidateQueries({ queryKey: orderKeys.detail(variables.id) })
      // Also invalidate dashboard to update revenue if needed
      queryClient.invalidateQueries({ queryKey: ['admin-dashboard'] })
      console.log('useUpdateOrderStatus - All queries invalidated')
      toast.success('Order status updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update order status')
    },
  })
}

// Cancel order
export function useCancelOrder() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, reason }: {
      id: string
      reason?: string
    }): Promise<Order> => {
      const response = await apiClient.post<Order>(`/orders/${id}/cancel`, {
        reason
      })
      return response.data
    },
    onSuccess: (data, variables) => {
      // Update specific order in cache
      queryClient.setQueryData(orderKeys.detail(variables.id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() })
      toast.success('Order cancelled successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to cancel order')
    },
  })
}

// Refund order (admin only)
export function useRefundOrder() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, amount, reason }: {
      id: string
      amount?: number
      reason?: string
    }): Promise<Order> => {
      const response = await apiClient.post<Order>(`/admin/orders/${id}/refund`, {
        amount,
        reason
      })
      return response.data
    },
    onSuccess: (data, variables) => {
      // Update specific order in cache
      queryClient.setQueryData(orderKeys.detail(variables.id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() })
      queryClient.invalidateQueries({ queryKey: orderKeys.admin() })
      toast.success('Order refunded successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to refund order')
    },
  })
}

// Update order shipping (admin only)
export function useUpdateOrderShipping() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, tracking_number, carrier, notes }: {
      id: string
      tracking_number?: string
      carrier?: string
      notes?: string
    }): Promise<Order> => {
      const response = await apiClient.patch<Order>(`/admin/orders/${id}/shipping`, {
        tracking_number,
        carrier,
        notes
      })
      return response.data
    },
    onSuccess: (data, variables) => {
      // Update specific order in cache
      queryClient.setQueryData(orderKeys.detail(variables.id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() })
      queryClient.invalidateQueries({ queryKey: orderKeys.admin() })
      toast.success('Shipping information updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update shipping information')
    },
  })
}

// Get order statistics (admin only)
export function useOrderStatistics(params: {
  date_from?: string
  date_to?: string
  group_by?: 'day' | 'week' | 'month'
} = {}) {
  return useQuery({
    queryKey: [...orderKeys.admin(), 'statistics', params],
    queryFn: async (): Promise<{
      total_orders: number
      total_revenue: number
      average_order_value: number
      orders_by_status: Record<string, number>
      revenue_by_period: Array<{
        period: string
        revenue: number
        orders: number
      }>
    }> => {
      const queryParams = new URLSearchParams()
      
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      if (params.group_by) queryParams.append('group_by', params.group_by)
      
      const url = `/admin/orders/statistics${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      const response = await apiClient.get(url)
      return response.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Get admin order details (detailed information)
export function useAdminOrderDetails(orderId: string) {
  return useQuery({
    queryKey: ['admin-order-details', orderId],
    queryFn: async (): Promise<Order> => {
      console.log('useAdminOrderDetails - Fetching order details:', orderId)
      
      const response = await apiClient.get<any>(`/admin/orders/${orderId}`)
      console.log('useAdminOrderDetails - Raw response:', response)

      // Handle SuccessResponse wrapper
      let orderData = response.data
      if (response.data && response.data.data) {
        console.log('useAdminOrderDetails - Unwrapping SuccessResponse')
        orderData = response.data.data
      }

      console.log('useAdminOrderDetails - Backend order data:', orderData)

      // Transform backend response to frontend Order format
      const transformedOrder: Partial<Order> = {
        // Basic order info
        ...(orderData.order?.id && { id: orderData.order.id }),
        ...(orderData.id && { id: orderData.id }),
        order_number: orderData.order?.order_number || orderData.order_number || '',
        status: orderData.order?.status || orderData.status || 'pending',
        payment_status: orderData.order?.payment_status || orderData.payment_status || 'pending',
        
        // Financial data
        subtotal: orderData.order?.subtotal || orderData.subtotal || 0,
        tax_amount: orderData.order?.tax_amount || orderData.tax_amount || 0,
        shipping_amount: orderData.order?.shipping_amount || orderData.shipping_amount || 0,
        discount_amount: orderData.order?.discount_amount || orderData.discount_amount || 0,
        total: orderData.order?.total || orderData.total || 0,
        
        // Meta data
        item_count: orderData.items?.length || orderData.item_count || 0,
        created_at: orderData.order?.created_at || orderData.created_at || new Date().toISOString(),
        updated_at: orderData.order?.updated_at || orderData.updated_at || new Date().toISOString(),
        currency: orderData.order?.currency || orderData.currency || 'USD',
        can_be_cancelled: orderData.order?.can_be_cancelled || orderData.can_be_cancelled || false,
        can_be_refunded: orderData.order?.can_be_refunded || orderData.can_be_refunded || false,
        
        // Transform customer to user format
        user: orderData.customer ? {
          id: orderData.customer.id,
          email: orderData.customer.email,
          first_name: orderData.customer.first_name,
          last_name: orderData.customer.last_name,
        } : orderData.user,

        // Transform items
        items: orderData.items || orderData.order_items || [],

        // Handle shipping address with field name transformation
        shipping_address: orderData.shipping_address ? {
          first_name: orderData.shipping_address.first_name,
          last_name: orderData.shipping_address.last_name,
          address_line_1: orderData.shipping_address.address_line_1 || orderData.shipping_address.AddressLine1,
          address_line_2: orderData.shipping_address.address_line_2 || orderData.shipping_address.AddressLine2,
          city: orderData.shipping_address.city,
          state: orderData.shipping_address.state,
          postal_code: orderData.shipping_address.postal_code || orderData.shipping_address.PostalCode,
          country: orderData.shipping_address.country,
          phone: orderData.shipping_address.phone,
        } : orderData.order?.shipping_address,

        // Handle billing address
        billing_address: orderData.billing_address ? {
          first_name: orderData.billing_address.first_name,
          last_name: orderData.billing_address.last_name,
          address_line_1: orderData.billing_address.address_line_1 || orderData.billing_address.AddressLine1,
          address_line_2: orderData.billing_address.address_line_2 || orderData.billing_address.AddressLine2,
          city: orderData.billing_address.city,
          state: orderData.billing_address.state,
          postal_code: orderData.billing_address.postal_code || orderData.billing_address.PostalCode,
          country: orderData.billing_address.country,
          phone: orderData.billing_address.phone,
        } : orderData.order?.billing_address,

        // Additional fields that might exist
        notes: orderData.order?.notes || orderData.notes,
      }

      console.log('useAdminOrderDetails - Transformed order:', transformedOrder)
      return transformedOrder as Order
    },
    enabled: !!orderId,
    staleTime: 30 * 1000,
    retry: 1,
  })
}

// Export orders (admin only)
export function useExportOrders() {
  return useMutation({
    mutationFn: async (params: {
      format: 'csv' | 'xlsx'
      date_from?: string
      date_to?: string
      status?: string
    }): Promise<Blob> => {
      const queryParams = new URLSearchParams()
      
      queryParams.append('format', params.format)
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      if (params.status) queryParams.append('status', params.status)
      
      const response = await apiClient.get(`/admin/orders/export?${queryParams.toString()}`, {
        responseType: 'blob'
      })
      return response.data
    },
    onSuccess: (blob, variables) => {
      // Create download link
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `orders.${variables.format}`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      toast.success('Orders exported successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to export orders')
    },
  })
}

// Search orders
export function useSearchOrders(query: string, options?: {
  enabled?: boolean
  limit?: number
}) {
  return useQuery({
    queryKey: [...orderKeys.lists(), 'search', query],
    queryFn: async (): Promise<Order[]> => {
      if (!query.trim()) return []
      
      const params = new URLSearchParams({
        search: query,
        limit: (options?.limit || 20).toString()
      })
      
      const response = await apiClient.get<PaginatedResponse<Order>>(`/orders?${params}`)
      return response.data.data
    },
    enabled: options?.enabled !== false && !!query.trim(),
    staleTime: 30 * 1000,
  })
}
