'use client'

import {
  TrendingUp,
  TrendingDown,
  DollarSign,
  ShoppingCart,
  Users,
  Package,
  Eye,
  MoreHorizontal,
  Star,
  Activity,
  Calendar,
  Clock,
  ArrowUpRight,
  ArrowDownRight,
  BarChart3,
  PieChart,
  Target,
  Zap
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuthStore } from '@/store/auth'
import { hasPermission, PERMISSIONS } from '@/lib/permissions'
import { RequirePermission } from '@/components/auth/permission-guard'
import { formatPrice, formatDate } from '@/lib/utils'
import { useAdminOrders } from '@/hooks/use-orders'
import { useUsers } from '@/hooks/use-users'
import { useProducts } from '@/hooks/use-products'
import { useAdminDashboard, useTopProducts, useRecentActivity } from '@/hooks/use-admin-dashboard'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export function AdminDashboard() {
  const { user, isAuthenticated } = useAuthStore()

  // Debug log
  console.log('AdminDashboard - Auth status:', { user, isAuthenticated })
  console.log('AdminDashboard - User role:', user?.role)

  // Fetch real data from APIs
  const { data: dashboardData, isLoading: dashboardLoading, error: dashboardError } = useAdminDashboard()
  const { data: ordersData, isLoading: ordersLoading } = useAdminOrders({ limit: 5 })
  const { data: usersData, isLoading: usersLoading } = useUsers({ limit: 100 })
  const { data: productsData, isLoading: productsLoading } = useProducts({ limit: 100 })
  const { data: topProductsData, isLoading: topProductsLoading } = useTopProducts({ limit: 5 })
  const { data: recentActivityData, isLoading: activityLoading } = useRecentActivity({ limit: 5 })

  // Debug log
  console.log('AdminDashboard - Dashboard data:', dashboardData)
  console.log('AdminDashboard - Dashboard loading:', dashboardLoading)
  console.log('AdminDashboard - Dashboard error:', dashboardError)

  // Use dashboard data if available, otherwise fallback to individual API calls
  const recentOrders = dashboardData?.recent_orders || ordersData?.data || []
  const totalUsers = dashboardData?.overview?.total_customers || usersData?.pagination?.total || 0
  const totalProducts = dashboardData?.overview?.total_products || productsData?.pagination?.total || 0
  const totalOrders = dashboardData?.overview?.total_orders || ordersData?.pagination?.total || 0
  const totalRevenue = dashboardData?.overview?.total_revenue || 0
  
  // Debug log for revenue
  console.log('AdminDashboard - Total revenue:', totalRevenue)
  console.log('AdminDashboard - Dashboard overview:', dashboardData?.overview)

  // Calculate stats with growth indicators (use real data if available)
  const stats = {
    revenue: {
      current: typeof totalRevenue === 'number' ? totalRevenue : 0,
      change: dashboardData?.revenue_growth || 15.2,
      isPositive: (dashboardData?.revenue_growth || 15.2) > 0,
    },
    orders: {
      current: typeof totalOrders === 'number' ? totalOrders : 0,
      change: dashboardData?.orders_growth || 8.5,
      isPositive: (dashboardData?.orders_growth || 8.5) > 0,
    },
    customers: {
      current: typeof totalUsers === 'number' ? totalUsers : 0,
      change: dashboardData?.customers_growth || 12.3,
      isPositive: (dashboardData?.customers_growth || 12.3) > 0,
    },
    products: {
      current: typeof totalProducts === 'number' ? totalProducts : 0,
      change: dashboardData?.products_growth || 5.7,
      isPositive: (dashboardData?.products_growth || 5.7) > 0,
    },
  }

  // Get recent activity data (use real data if available)
  const recentActivity = recentActivityData || dashboardData?.recent_activity || [
    { id: 'activity-1', type: 'order', message: 'New order received', time: '2 minutes ago', status: 'success' },
    { id: 'activity-2', type: 'user', message: 'New customer registered', time: '5 minutes ago', status: 'info' },
    { id: 'activity-3', type: 'product', message: 'Product updated', time: '10 minutes ago', status: 'warning' },
    { id: 'activity-4', type: 'order', message: 'Order shipped', time: '15 minutes ago', status: 'success' },
    { id: 'activity-5', type: 'system', message: 'System backup completed', time: '1 hour ago', status: 'info' },
  ]

  // Get top products (use real data if available)
  const topProducts = topProductsData || dashboardData?.charts?.top_products || (productsData?.data || [])
    .slice(0, 5)
    .map((product, index) => ({
      id: product.id,
      name: product.name,
      sales: Math.floor(Math.random() * 200) + 50,
      revenue: Math.floor((product.pricing?.price || 0) * (Math.floor(Math.random() * 100) + 20)),
      growth: Math.floor(Math.random() * 30) + 5,
    }))

  // Additional stats from dashboard
  const pendingOrders = dashboardData?.overview?.pending_orders || 0
  const lowStockItems = dashboardData?.overview?.low_stock_items || 0
  const pendingReviews = dashboardData?.overview?.pending_reviews || 0
  const activeUsers = dashboardData?.overview?.active_users || totalUsers

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'order': return <ShoppingCart className="h-4 w-4" />
      case 'user': return <Users className="h-4 w-4" />
      case 'product': return <Package className="h-4 w-4" />
      case 'system': return <Activity className="h-4 w-4" />
      default: return <Clock className="h-4 w-4" />
    }
  }

  // Loading state
  if (dashboardLoading || ordersLoading || usersLoading || productsLoading) {
    return (
      <div className={BIHUB_ADMIN_THEME.spacing.section}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {[...Array(4)].map((_, i) => (
            <div key={i} className={cn(
              BIHUB_ADMIN_THEME.components.card.base,
              'p-6 animate-pulse'
            )}>
              <div className="flex items-center justify-between mb-4">
                <div className="w-12 h-12 rounded-xl bg-gray-700"></div>
                <div className="w-16 h-6 bg-gray-700 rounded-full"></div>
              </div>
              <div className="space-y-2">
                <div className="h-4 bg-gray-700 rounded w-20"></div>
                <div className="h-8 bg-gray-700 rounded w-24"></div>
                <div className="h-3 bg-gray-700 rounded w-16"></div>
              </div>
            </div>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Welcome Header */}
      <BiHubPageHeader
        title={`Welcome back, ${user?.first_name}!`}
        subtitle="Here's your BiHub store performance overview and key metrics"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Dashboard' }
        ]}
        action={
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 px-4 py-2 bg-gray-800 rounded-lg">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span className="text-sm text-gray-300">Store Online</span>
            </div>
            <div className="text-sm text-gray-400">
              Last updated: {new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </div>
          </div>
        }
      />

      {/* Enhanced Quick Stats with Modern Design - Similar to Orders Page */}
      <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
        <div className="space-y-6">
          {/* Primary Stats Row - Modern Glass Design */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Total Revenue */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-green-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-green-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-green-500/5 to-emerald-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex items-center justify-between">
                <div>
                  <p className="text-green-400/80 text-sm font-medium uppercase tracking-wide">Total Revenue</p>
                  <p className="text-2xl font-bold text-green-100 mt-1">{formatPrice(stats.revenue.current)}</p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.revenue.isPositive ? (
                      <TrendingUp className="h-3 w-3 text-green-400" />
                    ) : (
                      <TrendingDown className="h-3 w-3 text-red-400" />
                    )}
                    <span className={cn(
                      "text-xs font-medium",
                      stats.revenue.isPositive ? "text-green-400" : "text-red-400"
                    )}>
                      {stats.revenue.change}%
                    </span>
                  </div>
                </div>
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-green-500/20 to-emerald-600/20 border border-green-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <DollarSign className="h-5 w-5 text-green-400" />
                </div>
              </div>
            </div>

            {/* Total Orders */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-blue-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-blue-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-indigo-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex items-center justify-between">
                <div>
                  <p className="text-blue-400/80 text-sm font-medium uppercase tracking-wide">Total Orders</p>
                  <p className="text-2xl font-bold text-blue-100 mt-1">{stats.orders.current}</p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.orders.isPositive ? (
                      <TrendingUp className="h-3 w-3 text-green-400" />
                    ) : (
                      <TrendingDown className="h-3 w-3 text-red-400" />
                    )}
                    <span className={cn(
                      "text-xs font-medium",
                      stats.orders.isPositive ? "text-green-400" : "text-red-400"
                    )}>
                      {stats.orders.change}%
                    </span>
                  </div>
                </div>
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <ShoppingCart className="h-5 w-5 text-blue-400" />
                </div>
              </div>
            </div>

            {/* Total Customers */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-purple-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-purple-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-violet-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex items-center justify-between">
                <div>
                  <p className="text-purple-400/80 text-sm font-medium uppercase tracking-wide">Total Customers</p>
                  <p className="text-2xl font-bold text-purple-100 mt-1">{stats.customers.current}</p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.customers.isPositive ? (
                      <TrendingUp className="h-3 w-3 text-green-400" />
                    ) : (
                      <TrendingDown className="h-3 w-3 text-red-400" />
                    )}
                    <span className={cn(
                      "text-xs font-medium",
                      stats.customers.isPositive ? "text-green-400" : "text-red-400"
                    )}>
                      {stats.customers.change}%
                    </span>
                  </div>
                </div>
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-purple-500/20 to-violet-600/20 border border-purple-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <Users className="h-5 w-5 text-purple-400" />
                </div>
              </div>
            </div>

            {/* Total Products */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-amber-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-amber-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-amber-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 to-orange-500/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex items-center justify-between">
                <div>
                  <p className="text-amber-400/80 text-sm font-medium uppercase tracking-wide">Total Products</p>
                  <p className="text-2xl font-bold text-amber-100 mt-1">{stats.products.current}</p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.products.isPositive ? (
                      <TrendingUp className="h-3 w-3 text-green-400" />
                    ) : (
                      <TrendingDown className="h-3 w-3 text-red-400" />
                    )}
                    <span className={cn(
                      "text-xs font-medium",
                      stats.products.isPositive ? "text-green-400" : "text-red-400"
                    )}>
                      {stats.products.change}%
                    </span>
                  </div>
                </div>
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                  <Package className="h-5 w-5 text-amber-400" />
                </div>
              </div>
            </div>
          </div>

          {/* Secondary Stats - Ultra Compact Cards */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-2">
            {/* Pending Orders */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-amber-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-amber-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-amber-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 to-orange-500/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex flex-col items-center text-center">
                <div className="w-6 h-6 rounded-md bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                  <Clock className="h-3 w-3 text-amber-400" />
                </div>
                <p className="text-xs text-amber-400/90 font-medium uppercase tracking-wider">Pending</p>
                <p className="text-base font-bold text-amber-300 mt-0.5">{pendingOrders}</p>
              </div>
            </div>

            {/* Low Stock */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-red-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-red-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-red-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-red-500/5 to-rose-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex flex-col items-center text-center">
                <div className="w-6 h-6 rounded-md bg-gradient-to-br from-red-500/20 to-rose-600/20 border border-red-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                  <Package className="h-3 w-3 text-red-400" />
                </div>
                <p className="text-xs text-red-400/90 font-medium uppercase tracking-wider">Low Stock</p>
                <p className="text-base font-bold text-red-300 mt-0.5">{lowStockItems}</p>
              </div>
            </div>

            {/* Pending Reviews */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-sky-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-sky-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-sky-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-sky-500/5 to-blue-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex flex-col items-center text-center">
                <div className="w-6 h-6 rounded-md bg-gradient-to-br from-sky-500/20 to-blue-600/20 border border-sky-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                  <Star className="h-3 w-3 text-sky-400" />
                </div>
                <p className="text-xs text-sky-400/90 font-medium uppercase tracking-wider">Reviews</p>
                <p className="text-base font-bold text-sky-300 mt-0.5">{pendingReviews}</p>
              </div>
            </div>

            {/* Active Users */}
            <div className="group relative bg-white/5 backdrop-blur-sm border border-emerald-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-emerald-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-emerald-500/10">
              <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-green-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
              <div className="relative flex flex-col items-center text-center">
                <div className="w-6 h-6 rounded-md bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                  <Users className="h-3 w-3 text-emerald-400" />
                </div>
                <p className="text-xs text-emerald-400/90 font-medium uppercase tracking-wider">Active</p>
                <p className="text-base font-bold text-emerald-300 mt-0.5">{activeUsers}</p>
              </div>
            </div>
          </div>
        </div>
      </RequirePermission>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Orders */}
        <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
          <BiHubAdminCard
            title="Recent Orders"
            subtitle="Latest orders from BiHub customers"
            icon={<ShoppingCart className="h-5 w-5 text-white" />}
            headerAction={
              <Button
                variant="outline"
                size="sm"
                className={BIHUB_ADMIN_THEME.components.button.ghost}
              >
                <Eye className="h-4 w-4 mr-2" />
                View All
              </Button>
            }
          >
            <div className="space-y-3">
              {recentOrders.length > 0 ? (
                recentOrders.slice(0, 5).map((order) => (
                  <div key={order.id} className="group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 hover:scale-[1.01] transition-all duration-200 shadow-sm hover:shadow-lg">
                    {/* Gradient Background */}
                    <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-purple-500/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>

                    <div className="relative flex items-center justify-between">
                      <div className="flex items-center gap-4">
                        {/* Modern Status Icon */}
                        <div className="relative w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center shadow-sm">
                          <ShoppingCart className="h-4 w-4 text-blue-400" />
                          <div className="absolute inset-0 bg-white/10 rounded-lg blur-sm"></div>
                        </div>

                        <div>
                          <p className="text-sm font-semibold text-white group-hover:text-[#FF9000] transition-colors">
                            Order #{order.order_number || order.id.slice(0, 8)}
                          </p>
                          <p className="text-xs text-gray-400 mt-0.5">
                            {order.user?.first_name} {order.user?.last_name}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p className="text-lg font-bold text-[#FF9000] group-hover:text-[#FF9000]/80 transition-colors">
                          {formatPrice((order?.total || order?.total_amount || 0))}
                        </p>
                        <div className="mt-1">
                          <BiHubStatusBadge status={getBadgeVariant(order.status)}>
                            {order.status}
                          </BiHubStatusBadge>
                        </div>
                      </div>
                    </div>
                  </div>
                ))
              ) : (
                <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-8 text-center">
                  <div className="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-slate-500/5 rounded-xl"></div>
                  <div className="relative">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-gray-500/20 to-slate-500/20 border border-gray-400/30 flex items-center justify-center mx-auto mb-4">
                      <ShoppingCart className="h-6 w-6 text-gray-400" />
                    </div>
                    <p className="text-lg font-semibold text-white mb-1">No recent orders</p>
                    <p className="text-sm text-gray-400">Orders will appear here once customers start placing them</p>
                  </div>
                </div>
              )}
            </div>
          </BiHubAdminCard>
        </RequirePermission>

        {/* Recent Activity */}
        <BiHubAdminCard
          title="Recent Activity"
          subtitle="Latest activities in your BiHub store"
          icon={<Activity className="h-5 w-5 text-white" />}
        >
          <div className="space-y-4">
            {recentActivity.map((activity, index) => (
              <div key={activity.id || `activity-${index}`} className="flex items-center gap-3 p-3 bg-gray-800 rounded-lg">
                <div className={cn(
                  'w-8 h-8 rounded-lg flex items-center justify-center',
                  activity.status === 'success' && 'bg-green-600',
                  activity.status === 'info' && 'bg-blue-600',
                  activity.status === 'warning' && 'bg-yellow-600',
                  activity.status === 'error' && 'bg-red-600'
                )}>
                  {getActivityIcon(activity.type)}
                </div>
                <div className="flex-1">
                  <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
                    {activity.message}
                  </p>
                  <p className={cn(BIHUB_ADMIN_THEME.typography.body.small, 'text-gray-500')}>
                    {activity.time}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      </div>

      {/* Top Products */}
      <RequirePermission permission={PERMISSIONS.PRODUCTS_VIEW}>
        <BiHubAdminCard
          title="Top Performing Products"
          subtitle="Best selling products in your BiHub store"
          icon={<Target className="h-5 w-5 text-white" />}
          headerAction={
            <Button
              variant="outline"
              size="sm"
              className={BIHUB_ADMIN_THEME.components.button.ghost}
            >
              <BarChart3 className="h-4 w-4 mr-2" />
              View Analytics
            </Button>
          }
        >
          <div className="space-y-4">
            {topProducts.length > 0 ? (
              topProducts.map((product: any, index: number) => (
                <div key={product.id || `product-${index}`} className="flex items-center justify-between p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm">
                      #{index + 1}
                    </div>
                    <div>
                      <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
                        {product.name}
                      </p>
                      <p className={BIHUB_ADMIN_THEME.typography.body.small}>
                        {product.sales} sales
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-lg font-bold text-[#FF9000]">
                      {formatPrice(typeof product.revenue === 'number' ? product.revenue : 0)}
                    </p>
                    <div className="flex items-center gap-1">
                      <ArrowUpRight className="h-3 w-3 text-green-500" />
                      <span className="text-xs text-green-500">+{product.growth || 0}%</span>
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <div className="text-center py-8">
                <Package className="h-12 w-12 text-gray-600 mx-auto mb-4" />
                <p className={BIHUB_ADMIN_THEME.typography.body.medium}>No products data</p>
              </div>
            )}
          </div>
        </BiHubAdminCard>
      </RequirePermission>
    </div>
  )
}
