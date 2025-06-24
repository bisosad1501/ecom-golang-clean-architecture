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
      const queryParams = new URLSearchParams()
      
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.status) queryParams.append('status', params.status)
      if (params.user_id) queryParams.append('user_id', params.user_id)
      
      const url = `/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      const response = await apiClient.get<PaginatedResponse<Order>>(url)
      return response.data
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
  return useQuery({
    queryKey: [...orderKeys.admin(), params],
    queryFn: async (): Promise<PaginatedResponse<Order>> => {
      const queryParams = new URLSearchParams()
      
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.status) queryParams.append('status', params.status)
      if (params.user_id) queryParams.append('user_id', params.user_id)
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      
      const url = `/admin/orders${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      const response = await apiClient.get<PaginatedResponse<Order>>(url)
      return response.data
    },
    staleTime: 30 * 1000,
  })
}

// Get single order
export function useOrder(id: string) {
  return useQuery({
    queryKey: orderKeys.detail(id),
    queryFn: async (): Promise<Order> => {
      const response = await apiClient.get<Order>(`/orders/${id}`)
      return response.data
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
    }): Promise<Order> => {
      const response = await apiClient.patch<Order>(`/admin/orders/${id}/status`, {
        status,
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
