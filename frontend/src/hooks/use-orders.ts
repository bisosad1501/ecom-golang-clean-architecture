import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import apiClient from '@/lib/api'
import { Order, PaginatedResponse, CreateOrderRequest, OrderEvent } from '@/types'
import { useEffect, useState } from 'react'

// Query keys
export const orderKeys = {
  all: ['orders'] as const,
  lists: () => [...orderKeys.all, 'list'] as const,
  list: (filters: Record<string, any>) => [...orderKeys.lists(), filters] as const,
  details: () => [...orderKeys.all, 'detail'] as const,
  detail: (id: string) => [...orderKeys.details(), id] as const,
  user: (userId: string) => [...orderKeys.all, 'user', userId] as const,
  admin: () => [...orderKeys.all, 'admin'] as const,
  events: (orderId: string) => [...orderKeys.detail(orderId), 'events'] as const,
}

// Get orders (user's own orders)
export function useOrders(params: {
  page?: number
  limit?: number
  search?: string
  status?: string
  user_id?: string
} = {}) {
  console.log('🔥 useOrders called with params:', params)

  const query = useQuery({
    queryKey: ['orders-v2', Date.now(), params.page, params.limit, params.search, params.status],
    queryFn: async (): Promise<PaginatedResponse<Order>> => {
      console.log('🔥 useOrders queryFn executing...')
      try {
        const queryParams = new URLSearchParams()
        
        // Use page-based pagination (backend expects page parameter)
        const limit = params.limit || 10
        const page = params.page || 1

        queryParams.append('limit', limit.toString())
        queryParams.append('page', page.toString())
        if (params.search) queryParams.append('search', params.search)
        if (params.status) queryParams.append('status', params.status)
        if (params.user_id) queryParams.append('user_id', params.user_id)
        
        const url = `/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`

        const response = await apiClient.get<any>(url)


        // The backend returns: { data: [...], pagination: {...} }
        // response.data contains the full backend response
        const responseData = response.data
        // Check if the response has the pagination structure
        if (responseData && responseData.data && responseData.pagination) {
          const pag = responseData.pagination
          const result = {
            data: responseData.data,
            pagination: {
            page: pag.page,
            limit: pag.limit,
            total: pag.total,
            total_pages: pag.total_pages,
            has_next: pag.has_next,
            has_prev: pag.has_prev,
            start_index: pag.start_index ?? ((pag.page - 1) * pag.limit + 1),
            end_index: pag.end_index ?? (Math.min(pag.page * pag.limit, pag.total)),
            next_page: pag.next_page ?? (pag.has_next ? pag.page + 1 : undefined),
            prev_page: pag.prev_page ?? (pag.has_prev ? pag.page - 1 : undefined),
            canonical_url: pag.canonical_url ?? undefined,
            page_sizes: pag.page_sizes ?? undefined,
            use_cursor: pag.use_cursor ?? false
            }
          }

          return result
        }

        // Fallback for old structure (array of orders) or unknown structure

        let ordersArray: Order[] = []

        // Check if responseData is directly an array
        if (Array.isArray(responseData)) {
          ordersArray = responseData
        }
        // Check if it's a SuccessResponse wrapper with data array
        else if (responseData && responseData.data && Array.isArray(responseData.data)) {
          ordersArray = responseData.data
        }
        // Check if it's just an object with orders property
        else if (responseData && Array.isArray(responseData.orders)) {
          ordersArray = responseData.orders
        }

        // Create fallback paginated response structure
        return {
          data: ordersArray,
          pagination: {
            page: page,
            limit: limit,
            total: ordersArray.length,
            total_pages: Math.ceil(ordersArray.length / limit) || 1,
            has_next: ordersArray.length === limit,
            has_prev: page > 1,
            start_index: (page - 1) * limit + 1,
            end_index: Math.min(page * limit, ordersArray.length),
            next_page: (ordersArray.length === limit) ? page + 1 : undefined,
            prev_page: (page > 1) ? page - 1 : undefined,
            canonical_url: undefined,
            page_sizes: undefined,
            use_cursor: false
          }
        }
      } catch (error) {
        console.error('useOrders - API call failed:', error)
        throw error
      }
    },
    staleTime: 0, // Force fresh data
    gcTime: 0, // No cache (React Query v5)
    refetchOnMount: true,
    refetchOnWindowFocus: true,
    enabled: true
  })

  console.log('🔥 useOrders query state:', {
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isError: query.isError,
    error: query.error,
    data: query.data,
    status: query.status,
    fetchStatus: query.fetchStatus
  })

  return query
}

// Simple implementation without React Query for debugging
export function useOrdersSimple(params: {
  page?: number
  limit?: number
  search?: string
  status?: string
  user_id?: string
} = {}) {
  console.log('🔥 useOrdersSimple called with params:', params)

  const [data, setData] = useState<PaginatedResponse<Order> | undefined>(undefined)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<any>(null)

  useEffect(() => {
    const fetchOrders = async () => {
      console.log('🔥 Simple fetch starting...')
      try {
        setIsLoading(true)
        setError(null)

        const queryParams = new URLSearchParams()
        if (params.page) queryParams.append('page', params.page.toString())
        if (params.limit) queryParams.append('limit', params.limit.toString())
        if (params.search) queryParams.append('search', params.search)
        if (params.status) queryParams.append('status', params.status)
        if (params.user_id) queryParams.append('user_id', params.user_id)

        const url = `/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
        console.log('🔥 Simple fetch URL:', url)

        const response = await apiClient.get<any>(url)
        console.log('🔥 Simple fetch response:', response)


        // response is already unwrapped by apiClient.get()
        // Backend returns: {data: [...], pagination: {...}}
        // apiClient.get() returns: response.data = {data: [...], pagination: {...}}
        const responseData = response as any
        console.log('🔥 Simple fetch responseData:', responseData)
        console.log('🔥 Simple fetch responseData type:', typeof responseData)
        console.log('🔥 Simple fetch responseData keys:', Object.keys(responseData || {}))
        console.log('🔥 Simple fetch responseData.data:', responseData?.data)
        console.log('🔥 Simple fetch responseData.pagination:', responseData?.pagination)

        // Check if the response has the pagination structure
        if (responseData && responseData.data && responseData.pagination) {
          console.log('🔥 Simple fetch - using paginated structure')
          const result = {
            data: responseData.data,
            pagination: responseData.pagination
          }
          setData(result)
        } else {
          console.log('🔥 Simple fetch - using fallback structure')
          // Fallback
          const ordersArray = Array.isArray(responseData) ? responseData : []
          const page = params.page || 1
          const limit = params.limit || 10
          setData({
            data: ordersArray,
            pagination: {
              page: page,
              limit: limit,
              total: ordersArray.length,
              total_pages: Math.ceil(ordersArray.length / limit) || 1,
              has_next: ordersArray.length === limit,
              has_prev: page > 1
            }
          })
        }
      } catch (err) {
        console.error('🔥 Simple fetch error:', err)
        setError(err)
      } finally {
        setIsLoading(false)
      }
    }

    fetchOrders()
  }, [params.page, params.limit, params.search, params.status, params.user_id])

  return { data, isLoading, error }
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

      // Use page-based pagination for admin orders
      if (params.page) queryParams.append('page', params.page.toString())
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

// Get order events
export function useOrderEvents(orderId: string, publicOnly: boolean = false) {
  return useQuery({
    queryKey: orderKeys.events(orderId),
    queryFn: async (): Promise<OrderEvent[]> => {
      const params = new URLSearchParams()
      if (publicOnly) params.append('public', 'true')

      const url = `/orders/${orderId}/events${params.toString() ? `?${params.toString()}` : ''}`
      const response = await apiClient.get<{ data: OrderEvent[] }>(url)
      return response.data.data
    },
    enabled: !!orderId,
    staleTime: 30 * 1000,
  })
}

// Add order note
export function useAddOrderNote() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ orderId, note, isPublic }: {
      orderId: string
      note: string
      isPublic: boolean
    }) => {
      const response = await apiClient.post(`/orders/${orderId}/notes`, {
        note,
        is_public: isPublic
      })
      return response.data
    },
    onSuccess: (_, variables) => {
      toast.success('Note added successfully')
      queryClient.invalidateQueries({ queryKey: orderKeys.events(variables.orderId) })
      queryClient.invalidateQueries({ queryKey: orderKeys.detail(variables.orderId) })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to add note')
    },
  })
}

// Update shipping info (Admin)
export function useUpdateShippingInfo() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ orderId, shippingData }: {
      orderId: string
      shippingData: {
        tracking_number: string
        carrier: string
        shipping_method: string
        tracking_url?: string
      }
    }) => {
      const response = await apiClient.put(`/admin/orders/${orderId}/shipping`, shippingData)
      return response.data
    },
    onSuccess: (_, variables) => {
      toast.success('Shipping info updated successfully')
      queryClient.invalidateQueries({ queryKey: orderKeys.events(variables.orderId) })
      queryClient.invalidateQueries({ queryKey: orderKeys.detail(variables.orderId) })
      queryClient.invalidateQueries({ queryKey: orderKeys.admin() })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update shipping info')
    },
  })
}

// Update delivery status (Admin)
export function useUpdateDeliveryStatus() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ orderId, status }: {
      orderId: string
      status: string
    }) => {
      const response = await apiClient.put(`/admin/orders/${orderId}/delivery`, { status })
      return response.data
    },
    onSuccess: (_, variables) => {
      toast.success('Delivery status updated successfully')
      queryClient.invalidateQueries({ queryKey: orderKeys.events(variables.orderId) })
      queryClient.invalidateQueries({ queryKey: orderKeys.detail(variables.orderId) })
      queryClient.invalidateQueries({ queryKey: orderKeys.admin() })
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update delivery status')
    },
  })
}
