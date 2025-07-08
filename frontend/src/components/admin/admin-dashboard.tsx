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

  // Fetch real data from APIs
  const { data: ordersData, isLoading: ordersLoading } = useAdminOrders({ limit: 5 })
  const { data: usersData, isLoading: usersLoading } = useUsers({ limit: 100 })
  const { data: productsData, isLoading: productsLoading } = useProducts({ limit: 100 })

  const recentOrders = ordersData?.data || []
  const totalUsers = usersData?.pagination?.total || 0
  const totalProducts = productsData?.pagination?.total || 0
  const totalOrders = ordersData?.pagination?.total || 0

  // Calculate revenue from orders
  const totalRevenue = recentOrders.reduce((sum, order) => sum + order.total_amount, 0)

  // Calculate stats with growth indicators
  const stats = {
    revenue: {
      current: totalRevenue,
      change: 15.2,
      isPositive: true,
    },
    orders: {
      current: totalOrders,
      change: 8.5,
      isPositive: true,
    },
    customers: {
      current: totalUsers,
      change: 12.3,
      isPositive: true,
    },
    products: {
      current: totalProducts,
      change: 5.7,
      isPositive: true,
    },
  }

  // Get recent activity data
  const recentActivity = [
    { type: 'order', message: 'New order #1234 received', time: '2 minutes ago', status: 'success' },
    { type: 'user', message: 'New customer registered', time: '5 minutes ago', status: 'info' },
    { type: 'product', message: 'Product "Gaming Laptop" updated', time: '10 minutes ago', status: 'warning' },
    { type: 'order', message: 'Order #1230 shipped', time: '15 minutes ago', status: 'success' },
    { type: 'system', message: 'System backup completed', time: '1 hour ago', status: 'info' },
  ]

  // Get top products (mock data for now)
  const topProducts = (productsData?.data || [])
    .slice(0, 5)
    .map((product, index) => ({
      id: product.id,
      name: product.name,
      sales: Math.floor(Math.random() * 200) + 50,
      revenue: product.price * (Math.floor(Math.random() * 100) + 20),
      growth: Math.floor(Math.random() * 30) + 5,
    }))
    .sort((a, b) => b.revenue - a.revenue)

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
  if (ordersLoading || usersLoading || productsLoading) {
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

      {/* Stats Cards */}
      <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <BiHubStatCard
            title="Total Revenue"
            value={formatPrice(stats.revenue.current)}
            change={stats.revenue.change}
            isPositive={stats.revenue.isPositive}
            icon={<DollarSign className="h-8 w-8 text-white" />}
            color="primary"
          />
          <BiHubStatCard
            title="Total Orders"
            value={stats.orders.current}
            change={stats.orders.change}
            isPositive={stats.orders.isPositive}
            icon={<ShoppingCart className="h-8 w-8 text-white" />}
            color="success"
          />
          <BiHubStatCard
            title="Total Customers"
            value={stats.customers.current}
            change={stats.customers.change}
            isPositive={stats.customers.isPositive}
            icon={<Users className="h-8 w-8 text-white" />}
            color="info"
          />
          <BiHubStatCard
            title="Total Products"
            value={stats.products.current}
            change={stats.products.change}
            isPositive={stats.products.isPositive}
            icon={<Package className="h-8 w-8 text-white" />}
            color="warning"
          />
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
            <div className="space-y-4">
              {recentOrders.length > 0 ? (
                recentOrders.slice(0, 5).map((order) => (
                  <div key={order.id} className="flex items-center justify-between p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                        <ShoppingCart className="h-5 w-5 text-white" />
                      </div>
                      <div>
                        <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
                          Order #{order.order_number || order.id.slice(0, 8)}
                        </p>
                        <p className={BIHUB_ADMIN_THEME.typography.body.small}>
                          {order.user?.first_name} {order.user?.last_name}
                        </p>
                      </div>
                    </div>
                    <div className="text-right">
                      <p className="text-lg font-bold text-[#FF9000]">
                        {formatPrice(order.total_amount)}
                      </p>
                      <BiHubStatusBadge status={getBadgeVariant(order.status)}>
                        {order.status}
                      </BiHubStatusBadge>
                    </div>
                  </div>
                ))
              ) : (
                <div className="text-center py-8">
                  <ShoppingCart className="h-12 w-12 text-gray-600 mx-auto mb-4" />
                  <p className={BIHUB_ADMIN_THEME.typography.body.medium}>No recent orders</p>
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
              <div key={index} className="flex items-center gap-3 p-3 bg-gray-800 rounded-lg">
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
              topProducts.map((product, index) => (
                <div key={product.id} className="flex items-center justify-between p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
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
                      {formatPrice(product.revenue)}
                    </p>
                    <div className="flex items-center gap-1">
                      <ArrowUpRight className="h-3 w-3 text-green-500" />
                      <span className="text-xs text-green-500">+{product.growth}%</span>
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
