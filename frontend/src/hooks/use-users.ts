import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import apiClient from '@/lib/api'
import { User, PaginatedResponse } from '@/types'

// Query keys
export const userKeys = {
  all: ['users'] as const,
  lists: () => [...userKeys.all, 'list'] as const,
  list: (filters: Record<string, any>) => [...userKeys.lists(), filters] as const,
  details: () => [...userKeys.all, 'detail'] as const,
  detail: (id: string) => [...userKeys.details(), id] as const,
  profile: () => [...userKeys.all, 'profile'] as const,
  admin: () => [...userKeys.all, 'admin'] as const,
}

// Get current user profile
export function useProfile() {
  return useQuery({
    queryKey: userKeys.profile(),
    queryFn: async (): Promise<User> => {
      const response = await apiClient.get<User>('/users/profile')
      return response.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Get users (admin only)
export function useUsers(params: {
  page?: number
  limit?: number
  search?: string
  role?: string
  is_active?: boolean
} = {}) {
  return useQuery({
    queryKey: [...userKeys.admin(), params],
    queryFn: async (): Promise<PaginatedResponse<User>> => {
      const queryParams = new URLSearchParams()
      
      // Convert page to offset for backend
      const offset = params.page ? (params.page - 1) * (params.limit || 10) : 0
      if (offset > 0) queryParams.append('offset', offset.toString())
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.role) queryParams.append('role', params.role)
      if (params.is_active !== undefined) queryParams.append('is_active', params.is_active.toString())
      
      const url = `/admin/users${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useUsers - Making API call to:', url)
      console.log('useUsers - Query params:', params)
      
      try {
        const response = await apiClient.get<any>(url)
        console.log('useUsers - Raw API response:', response.data)
        console.log('useUsers - Response data field:', response.data.data)
        console.log('useUsers - Response pagination field:', response.data.pagination)
        
        // Fix: API response structure - backend returns {message, data: {users, total, pagination}}
        const responseData = response.data.data || response.data
        
        // Transform backend user data to match frontend User interface
        const transformUser = (backendUser: any): User => ({
          ...backendUser,
          is_active: backendUser.status === 'active', // Convert status to is_active
          status: backendUser.status === 'active' ? 'active' : 'inactive' as any, // Keep status for compatibility
        })
        
        const users = Array.isArray(responseData.users) 
          ? responseData.users.map(transformUser)
          : Array.isArray(responseData) 
            ? responseData.map(transformUser) 
            : []
        
        const finalData: PaginatedResponse<User> = {
          data: users,
          pagination: responseData.pagination || {
            page: 1,
            limit: 20,
            total: responseData.total || users.length,
            total_pages: Math.ceil((responseData.total || users.length) / 20),
            has_next: false,
            has_prev: false,
          }
        }
        
        console.log('useUsers - Final transformed data:', finalData)
        return finalData
      } catch (error) {
        console.error('useUsers - API error:', error)
        throw error
      }
    },
    staleTime: 30 * 1000, // 30 seconds
  })
}

// Get single user (admin only)
export function useUser(id: string) {
  return useQuery({
    queryKey: userKeys.detail(id),
    queryFn: async (): Promise<User> => {
      const response = await apiClient.get<User>(`/admin/users/${id}`)
      return response.data
    },
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
  })
}

// Update profile
export function useUpdateProfile() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: {
      first_name?: string
      last_name?: string
      phone?: string
      profile?: {
        date_of_birth?: string
        gender?: string
        address?: string
        city?: string
        country?: string
        postal_code?: string
      }
    }): Promise<User> => {
      const response = await apiClient.put<User>('/users/profile', data)
      return response.data
    },
    onSuccess: (data) => {
      // Update profile in cache
      queryClient.setQueryData(userKeys.profile(), data)
      toast.success('Profile updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update profile')
    },
  })
}

// Change password
export function useChangePassword() {
  return useMutation({
    mutationFn: async (data: {
      current_password: string
      new_password: string
      confirm_password: string
    }): Promise<void> => {
      await apiClient.post('/users/change-password', data)
    },
    onSuccess: () => {
      toast.success('Password changed successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to change password')
    },
  })
}

// Activate user (admin only)
export function useActivateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string): Promise<User> => {
      const response = await apiClient.post<User>(`/admin/users/${id}/activate`)
      return response.data
    },
    onSuccess: (data, id) => {
      // Update user in cache
      queryClient.setQueryData(userKeys.detail(id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('User activated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to activate user')
    },
  })
}

// Deactivate user (admin only)
export function useDeactivateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string): Promise<User> => {
      const response = await apiClient.post<User>(`/admin/users/${id}/deactivate`)
      return response.data
    },
    onSuccess: (data, id) => {
      // Update user in cache
      queryClient.setQueryData(userKeys.detail(id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('User deactivated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to deactivate user')
    },
  })
}

// Update user role (admin only)
export function useUpdateUserRole() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, role }: {
      id: string
      role: 'admin' | 'moderator' | 'customer'
    }): Promise<User> => {
      const response = await apiClient.patch<User>(`/admin/users/${id}/role`, { role })
      return response.data
    },
    onSuccess: (data, variables) => {
      // Update user in cache
      queryClient.setQueryData(userKeys.detail(variables.id), data)
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('User role updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update user role')
    },
  })
}

// Create user (admin only)
export function useCreateUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: {
      email: string
      password: string
      first_name: string
      last_name: string
      phone?: string
      role?: 'admin' | 'moderator' | 'customer'
    }): Promise<User> => {
      const response = await apiClient.post<User>('/admin/users', data)
      return response.data
    },
    onSuccess: () => {
      // Invalidate and refetch users
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('User created successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create user')
    },
  })
}

// Delete user (admin only)
export function useDeleteUser() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await apiClient.delete(`/admin/users/${id}`)
    },
    onSuccess: (_, id) => {
      // Remove from cache
      queryClient.removeQueries({ queryKey: userKeys.detail(id) })
      // Invalidate lists to refetch
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('User deleted successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to delete user')
    },
  })
}

// Get user statistics (admin only)
export function useUserStatistics(params: {
  date_from?: string
  date_to?: string
  group_by?: 'day' | 'week' | 'month'
} = {}) {
  return useQuery({
    queryKey: [...userKeys.admin(), 'statistics', params],
    queryFn: async (): Promise<{
      total_users: number
      active_users: number
      new_users_today: number
      new_users_this_week: number
      new_users_this_month: number
      users_by_role: Record<string, number>
      registrations_by_period: Array<{
        period: string
        count: number
      }>
    }> => {
      const queryParams = new URLSearchParams()
      
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      if (params.group_by) queryParams.append('group_by', params.group_by)
      
      const url = `/admin/users/statistics${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      const response = await apiClient.get(url)
      return response.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Export users (admin only)
export function useExportUsers() {
  return useMutation({
    mutationFn: async (params: {
      format: 'csv' | 'xlsx'
      role?: string
      is_active?: boolean
      date_from?: string
      date_to?: string
    }): Promise<Blob> => {
      const queryParams = new URLSearchParams()
      
      queryParams.append('format', params.format)
      if (params.role) queryParams.append('role', params.role)
      if (params.is_active !== undefined) queryParams.append('is_active', params.is_active.toString())
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      
      const response = await apiClient.get(`/admin/users/export?${queryParams.toString()}`, {
        responseType: 'blob'
      })
      return response.data
    },
    onSuccess: (blob, variables) => {
      // Create download link
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `users.${variables.format}`
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      toast.success('Users exported successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to export users')
    },
  })
}

// Search users
export function useSearchUsers(query: string, options?: {
  enabled?: boolean
  limit?: number
}) {
  return useQuery({
    queryKey: [...userKeys.admin(), 'search', query],
    queryFn: async (): Promise<User[]> => {
      if (!query.trim()) return []

      const params = new URLSearchParams({
        search: query,
        limit: (options?.limit || 20).toString()
      })

      const response = await apiClient.get<PaginatedResponse<User>>(`/admin/users?${params}`)
      return response.data.data
    },
    enabled: options?.enabled !== false && !!query.trim(),
    staleTime: 30 * 1000,
  })
}

// Customer search and segmentation hooks
export interface CustomerSearchFilters {
  query?: string
  role?: string
  status?: string
  is_active?: boolean
  email_verified?: boolean
  phone_verified?: boolean
  two_factor_enabled?: boolean
  membership_tier?: string
  customer_segment?: string
  min_total_spent?: number
  max_total_spent?: number
  min_total_orders?: number
  max_total_orders?: number
  min_loyalty_points?: number
  max_loyalty_points?: number
  created_from?: string
  created_to?: string
  last_login_from?: string
  last_login_to?: string
  last_activity_from?: string
  last_activity_to?: string
  include_inactive?: boolean
  include_unverified?: boolean
  sort_by?: string
  sort_order?: string
  limit?: number
  offset?: number
}

export interface CustomerSearchResult {
  id: string
  email: string
  first_name: string
  last_name: string
  phone?: string
  role: string
  status: string
  is_active: boolean
  email_verified: boolean
  phone_verified: boolean
  two_factor_enabled: boolean
  last_login?: string
  last_activity?: string
  order_count: number
  total_spent: number
  loyalty_points: number
  membership_tier: string
  customer_segment: string
  security_level: string
  is_high_value: boolean
  is_vip: boolean
  created_at: string
  updated_at: string
}

export interface CustomerSearchResponse {
  customers: CustomerSearchResult[]
  total: number
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
  facets?: {
    roles: Array<{ value: string; count: number }>
    statuses: Array<{ value: string; count: number }>
    membership_tiers: Array<{ value: string; count: number }>
    customer_segments: Array<{ value: string; count: number }>
    security_levels: Array<{ value: string; count: number }>
    verification_status: {
      email_verified: number
      phone_verified: number
      two_factor_enabled: number
    }
  }
}

export function useCustomerSearch(filters: CustomerSearchFilters) {
  return useQuery({
    queryKey: ['admin', 'customers', 'search', filters],
    queryFn: async (): Promise<CustomerSearchResponse> => {
      const params = new URLSearchParams()

      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, value.toString())
        }
      })

      const response = await apiClient.get<{
        data: CustomerSearchResponse
      }>(`/admin/customers/search?${params}`)

      return response.data.data
    },
    staleTime: 30 * 1000,
  })
}

export interface CustomerSegmentInfo {
  segment: string
  count: number
  percentage: number
  avg_spent: number
  avg_orders: number
  description: string
}

export interface CustomerSegmentsResponse {
  segments: CustomerSegmentInfo[]
  total: number
}

export function useCustomerSegments() {
  return useQuery({
    queryKey: ['admin', 'customers', 'segments'],
    queryFn: async (): Promise<CustomerSegmentsResponse> => {
      const response = await apiClient.get<{
        data: CustomerSegmentsResponse
      }>('/admin/customers/segments')

      return response.data.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export interface CustomerAnalyticsResponse {
  overview: {
    total_customers: number
    active_customers: number
    new_customers: number
    returning_customers: number
    churn_rate: number
    avg_lifetime_value: number
    avg_order_value: number
  }
  segment_breakdown: CustomerSegmentInfo[]
  tier_distribution: Array<{
    tier: string
    count: number
    percentage: number
    revenue: number
  }>
  geographic_distribution: Array<{
    country: string
    count: number
    percentage: number
  }>
  acquisition_trends: Array<{
    date: string
    count: number
  }>
  retention_metrics: {
    day_30_retention: number
    day_90_retention: number
    day_365_retention: number
    repeat_purchase_rate: number
  }
}

export function useCustomerAnalytics(filters?: {
  date_from?: string
  date_to?: string
  segment?: string
}) {
  return useQuery({
    queryKey: ['admin', 'customers', 'analytics', filters],
    queryFn: async (): Promise<CustomerAnalyticsResponse> => {
      const params = new URLSearchParams()

      if (filters) {
        Object.entries(filters).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== '') {
            params.append(key, value.toString())
          }
        })
      }

      const response = await apiClient.get<{
        data: CustomerAnalyticsResponse
      }>(`/admin/customers/analytics?${params}`)

      return response.data.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export function useHighValueCustomers(limit: number = 50) {
  return useQuery({
    queryKey: ['admin', 'customers', 'high-value', limit],
    queryFn: async (): Promise<{
      customers: CustomerSearchResult[]
      total: number
      criteria: {
        min_total_spent: number
        min_total_orders: number
      }
    }> => {
      const response = await apiClient.get<{
        data: {
          customers: CustomerSearchResult[]
          total: number
          criteria: {
            min_total_spent: number
            min_total_orders: number
          }
        }
      }>(`/admin/customers/high-value?limit=${limit}`)

      return response.data.data
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Bulk operations
export function useBulkUpdateUsers() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (updates: Array<{
      id: string
      data: Partial<User>
    }>): Promise<User[]> => {
      const promises = updates.map(({ id, data }) =>
        apiClient.patch<User>(`/admin/users/${id}`, data)
      )
      const responses = await Promise.all(promises)
      return responses.map(response => response.data)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: userKeys.admin() })
      toast.success('Users updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update users')
    },
  })
}
