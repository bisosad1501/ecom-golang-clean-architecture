'use client'

import { useQuery } from '@tanstack/react-query'
import { apiClient } from '@/lib/api'

// Query keys
const adminDashboardKeys = {
  all: ['admin-dashboard'] as const,
  dashboard: () => [...adminDashboardKeys.all, 'dashboard'] as const,
  stats: () => [...adminDashboardKeys.all, 'stats'] as const,
  analytics: (params: Record<string, any>) => [...adminDashboardKeys.all, 'analytics', params] as const,
}

// Dashboard data types (matching backend response)
export interface DashboardStats {
  overview: {
    total_revenue: number    // Net revenue (current)
    gross_revenue: number    // Before discounts
    product_revenue: number  // Only product sales
    tax_collected: number    // Total tax amount
    shipping_revenue: number // Shipping fees
    discounts_given: number  // Total discounts
    total_orders: number
    total_customers: number
    total_products: number
    pending_orders: number
    low_stock_items: number
    pending_reviews: number
    active_users: number
  }
  // Growth metrics
  revenue_growth?: number
  orders_growth?: number
  customers_growth?: number
  products_growth?: number
  charts: {
    revenue_chart: Array<{
      date: string
      revenue: number
      orders: number
    }>
    orders_chart: any
    top_products: any
    top_categories: any
  }
  recent_activity: Array<{
    type: string
    message: string
    time: string
    status: string
  }>
  recent_orders: any[]
}

export interface SystemStats {
  server_uptime: string
  database_size: string
  total_files: number
  cache_hit_rate: number
  active_sessions: number
  memory_usage: number
  cpu_usage: number
}

// Get admin dashboard data
export function useAdminDashboard(params: {
  period?: string
  date_from?: string
  date_to?: string
} = {}) {
  return useQuery({
    queryKey: adminDashboardKeys.dashboard(),
    queryFn: async (): Promise<DashboardStats> => {
      const queryParams = new URLSearchParams()
      
      if (params.period) queryParams.append('period', params.period)
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      
      const url = `/admin/dashboard${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useAdminDashboard - Making API call to:', url)
      
      try {
        const response = await apiClient.get<any>(url)
        console.log('useAdminDashboard - Raw response:', response)
        console.log('useAdminDashboard - Response data structure:', JSON.stringify(response.data, null, 2))

        // Axios response format: response.data contains the backend data directly
        const dashboardData = response.data

        console.log('useAdminDashboard - Final dashboard data:', dashboardData)
        console.log('useAdminDashboard - Overview object:', dashboardData?.overview)
        console.log('useAdminDashboard - Total revenue value:', dashboardData?.overview?.total_revenue)
        return dashboardData
      } catch (error) {
        console.error('useAdminDashboard - API error:', error)
        throw error
      }
    },
    staleTime: 0, // Disable cache for testing
    refetchInterval: false, // Disable auto refetch
    gcTime: 0, // Disable cache completely (new syntax)
  })
}

// Get system stats
export function useSystemStats() {
  return useQuery({
    queryKey: adminDashboardKeys.stats(),
    queryFn: async (): Promise<SystemStats> => {
      console.log('useSystemStats - Making API call to: /admin/dashboard/stats')
      
      try {
        const response = await apiClient.get<any>('/admin/dashboard/stats')
        console.log('useSystemStats - Raw response:', response)

        // Axios response format: response.data contains the backend data directly
        const statsData = response.data

        console.log('useSystemStats - Final stats data:', statsData)
        return statsData
      } catch (error) {
        console.error('useSystemStats - API error:', error)
        throw error
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchInterval: 10 * 60 * 1000, // Refetch every 10 minutes
  })
}

// Get analytics data
export function useAdminAnalytics(params: {
  period?: string
  date_from?: string
  date_to?: string
  metric?: string
} = {}) {
  return useQuery({
    queryKey: adminDashboardKeys.analytics(params),
    queryFn: async (): Promise<any> => {
      const queryParams = new URLSearchParams()
      
      if (params.period) queryParams.append('period', params.period)
      if (params.date_from) queryParams.append('date_from', params.date_from)
      if (params.date_to) queryParams.append('date_to', params.date_to)
      if (params.metric) queryParams.append('metric', params.metric)
      
      const url = `/admin/analytics/sales${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useAdminAnalytics - Making API call to:', url)
      
      try {
        const response = await apiClient.get<any>(url)
        console.log('useAdminAnalytics - Raw response:', response)

        // Axios response format: response.data contains the backend data directly
        const analyticsData = response.data

        console.log('useAdminAnalytics - Final analytics data:', analyticsData)
        return analyticsData
      } catch (error) {
        console.error('useAdminAnalytics - API error:', error)
        throw error
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
    enabled: !!params.metric || !!params.period,
  })
}

// Get top products analytics
export function useTopProducts(params: {
  limit?: number
  period?: string
} = {}) {
  return useQuery({
    queryKey: [...adminDashboardKeys.all, 'top-products', params],
    queryFn: async (): Promise<any[]> => {
      const queryParams = new URLSearchParams()
      
      if (params.limit) queryParams.append('limit', params.limit.toString())
      if (params.period) queryParams.append('period', params.period)
      
      const url = `/admin/analytics/top-products${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useTopProducts - Making API call to:', url)
      
      try {
        const response = await apiClient.get<any>(url)
        console.log('useTopProducts - Raw response:', response)
        
        // Handle SuccessResponse wrapper
        let productsData = response.data
        if (response.data && response.data.data) {
          console.log('useTopProducts - Unwrapping SuccessResponse')
          productsData = response.data.data
        }
        
        console.log('useTopProducts - Final products data:', productsData)
        return productsData || []
      } catch (error) {
        console.error('useTopProducts - API error:', error)
        return []
      }
    },
    staleTime: 10 * 60 * 1000, // 10 minutes
  })
}

// Get recent activity
export function useRecentActivity(params: {
  limit?: number
} = {}) {
  return useQuery({
    queryKey: [...adminDashboardKeys.all, 'recent-activity', params],
    queryFn: async (): Promise<any[]> => {
      const queryParams = new URLSearchParams()
      
      if (params.limit) queryParams.append('limit', params.limit.toString())
      
      const url = `/admin/dashboard/activity${queryParams.toString() ? `?${queryParams.toString()}` : ''}`
      console.log('useRecentActivity - Making API call to:', url)
      
      try {
        const response = await apiClient.get<any>(url)
        console.log('useRecentActivity - Raw response:', response)
        
        // Handle SuccessResponse wrapper
        let activityData = response.data
        if (response.data && response.data.data) {
          console.log('useRecentActivity - Unwrapping SuccessResponse')
          activityData = response.data.data
        }
        
        console.log('useRecentActivity - Final activity data:', activityData)
        return activityData || []
      } catch (error) {
        console.error('useRecentActivity - API error:', error)
        // Return mock data as fallback
        return [
          { type: 'order', message: 'New order received', time: '2 minutes ago', status: 'success' },
          { type: 'user', message: 'New customer registered', time: '5 minutes ago', status: 'info' },
          { type: 'product', message: 'Product updated', time: '10 minutes ago', status: 'warning' },
        ]
      }
    },
    staleTime: 2 * 60 * 1000, // 2 minutes
  })
}
