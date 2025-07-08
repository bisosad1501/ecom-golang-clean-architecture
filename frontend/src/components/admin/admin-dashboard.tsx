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
  Star
} from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { hasPermission, PERMISSIONS } from '@/lib/permissions'
import { RequirePermission } from '@/components/auth/permission-guard'
import { formatPrice } from '@/lib/utils'
import { useAdminOrders } from '@/hooks/use-orders'
import { useUsers } from '@/hooks/use-users'
import { useProducts } from '@/hooks/use-products'

export function AdminDashboard() {
  console.log('=== AdminDashboard COMPONENT RENDERING ===')
  
  const { user, isAuthenticated } = useAuthStore()

  console.log('AdminDashboard - Auth state:', { user: !!user, isAuthenticated, userRole: user?.role })

  // Fetch real data from APIs
  const { data: ordersData, isLoading: ordersLoading } = useAdminOrders({ limit: 4 })
  const { data: usersData, isLoading: usersLoading, error: usersError } = useUsers({ limit: 100 })
  const { data: productsData, isLoading: productsLoading } = useProducts({ limit: 100 })

  // Debug logging
  console.log('Debug - usersData:', usersData)
  console.log('Debug - usersLoading:', usersLoading)
  console.log('Debug - usersError:', usersError)
  console.log('Debug - ordersData:', ordersData)
  console.log('Debug - productsData:', productsData)

  const recentOrders = ordersData?.data || []
  const totalUsers = usersData?.pagination?.total || 0
  const totalProducts = productsData?.pagination?.total || 0
  const totalOrders = ordersData?.pagination?.total || 0

  // Calculate revenue from orders
  const totalRevenue = recentOrders.reduce((sum, order) => sum + order.total_amount, 0)

  // Mock previous data for comparison (would come from a different API call)
  const stats = {
    revenue: {
      current: totalRevenue,
      previous: totalRevenue * 0.85, // Mock 15% growth
      change: 15.0,
    },
    orders: {
      current: totalOrders,
      previous: Math.floor(totalOrders * 0.9), // Mock 10% growth
      change: 10.0,
    },
    customers: {
      current: totalUsers,
      previous: Math.floor(totalUsers * 0.95), // Mock 5% growth
      change: 5.0,
    },
    products: {
      current: totalProducts,
      previous: Math.floor(totalProducts * 0.92), // Mock 8% growth
      change: 8.0,
    },
  }

  // Get top products from products data (mock popularity)
  const topProducts = (productsData?.data || [])
    .slice(0, 4)
    .map((product, index) => ({
      id: product.id,
      name: product.name,
      sales: Math.floor(Math.random() * 200) + 50, // Mock sales data
      revenue: product.price * (Math.floor(Math.random() * 100) + 20), // Mock revenue
    }))
    .sort((a, b) => b.revenue - a.revenue)

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'completed':
        return 'success'
      case 'processing':
        return 'warning'
      case 'shipped':
        return 'default'
      case 'pending':
        return 'secondary'
      default:
        return 'secondary'
    }
  }

  // Loading state
  if (ordersLoading || usersLoading || productsLoading) {
    return (
      <div className="space-y-8">
        {/* Loading skeleton for stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {[...Array(4)].map((_, i) => (
            <Card key={i} variant="elevated" className="border-0 shadow-large animate-pulse bg-gray-800/50">
              <CardContent className="p-8">
                <div className="flex items-center justify-between mb-6">
                  <div className="w-16 h-16 rounded-3xl bg-gray-700"></div>
                  <div className="w-20 h-6 bg-gray-700 rounded-full"></div>
                </div>
                <div className="space-y-2">
                  <div className="h-4 bg-gray-700 rounded w-24"></div>
                  <div className="h-8 bg-gray-700 rounded w-32"></div>
                  <div className="h-3 bg-gray-700 rounded w-28"></div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Loading skeleton for content */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {[...Array(2)].map((_, i) => (
            <Card key={i} variant="elevated" className="border-0 shadow-large animate-pulse bg-gray-800/50">
              <CardHeader className="pb-6">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-2xl bg-gray-700"></div>
                    <div className="h-6 bg-gray-700 rounded w-32"></div>
                  </div>
                  <div className="h-8 bg-gray-700 rounded w-20"></div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {[...Array(3)].map((_, j) => (
                    <div key={j} className="h-16 bg-gray-700 rounded-2xl"></div>
                  ))}
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-8">
      {/* Enhanced Welcome Header */}
      <div className="relative overflow-hidden bg-gradient-to-br from-[#FF9000] via-[#e67e00] to-[#cc6600] rounded-3xl shadow-2xl">
        <div className="absolute inset-0 bg-black/10"></div>
        <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -translate-y-32 translate-x-32"></div>
        <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24"></div>

        <div className="relative p-8 lg:p-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div className="text-white">
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center">
                  <TrendingUp className="h-6 w-6 text-white" />
                </div>
                <span className="text-white/90 font-semibold">ADMIN DASHBOARD</span>
              </div>

              <h1 className="text-4xl lg:text-5xl font-bold mb-4">
                Welcome back, <span className="text-white/90">{user?.first_name}!</span>
              </h1>
              <p className="text-xl text-white/80 leading-relaxed">
                Here's your store performance overview and key metrics for today.
              </p>
            </div>

            <div className="flex flex-col sm:flex-row gap-4">
              <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6 text-center">
                <p className="text-white/70 text-sm font-medium">Last Updated</p>
                <p className="text-white font-semibold text-lg">
                  {new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </p>
              </div>

              <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6 text-center">
                <p className="text-white/70 text-sm font-medium">Store Status</p>
                <div className="flex items-center justify-center gap-2 mt-1">
                  <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse"></div>
                  <p className="text-white font-semibold">Online</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Enhanced Stats Cards */}
      <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <Card variant="elevated" className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group bg-gray-800/50 border-gray-700/50">
            <CardContent className="p-8">
              <div className="flex items-center justify-between mb-6">
                <div className="w-16 h-16 rounded-3xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center shadow-large group-hover:scale-110 transition-transform duration-300">
                  <DollarSign className="h-8 w-8 text-white" />
                </div>
                <div className="flex items-center gap-2 px-3 py-1 bg-emerald-900/30 rounded-full">
                  <TrendingUp className="h-4 w-4 text-emerald-400" />
                  <span className="text-sm font-semibold text-emerald-400">+{stats.revenue.change}%</span>
                </div>
              </div>

              <div>
                <p className="text-sm font-medium text-gray-400 mb-2">Total Revenue</p>
                <p className="text-3xl font-bold text-white">{formatPrice(stats.revenue.current)}</p>
                <p className="text-sm text-gray-400 mt-2">
                  vs {formatPrice(stats.revenue.previous)} last month
                </p>
              </div>
            </CardContent>
          </Card>

          <Card variant="elevated" className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group bg-gray-800/50 border-gray-700/50">
            <CardContent className="p-8">
              <div className="flex items-center justify-between mb-6">
                <div className="w-16 h-16 rounded-3xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-large group-hover:scale-110 transition-transform duration-300">
                  <ShoppingCart className="h-8 w-8 text-white" />
                </div>
                <div className="flex items-center gap-2 px-3 py-1 bg-blue-900/30 rounded-full">
                  <TrendingUp className="h-4 w-4 text-blue-400" />
                  <span className="text-sm font-semibold text-blue-400">+{stats.orders.change}%</span>
                </div>
              </div>

              <div>
                <p className="text-sm font-medium text-gray-400 mb-2">Total Orders</p>
                <p className="text-3xl font-bold text-white">{stats.orders.current.toLocaleString()}</p>
                <p className="text-sm text-gray-400 mt-2">
                  vs {stats.orders.previous.toLocaleString()} last month
                </p>
              </div>
            </CardContent>
          </Card>

          <Card variant="elevated" className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group bg-gray-800/50 border-gray-700/50">
            <CardContent className="p-8">
              <div className="flex items-center justify-between mb-6">
                <div className="w-16 h-16 rounded-3xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center shadow-large group-hover:scale-110 transition-transform duration-300">
                  <Users className="h-8 w-8 text-white" />
                </div>
                <div className="flex items-center gap-2 px-3 py-1 bg-purple-900/30 rounded-full">
                  <TrendingUp className="h-4 w-4 text-purple-400" />
                  <span className="text-sm font-semibold text-purple-400">+{stats.customers.change}%</span>
                </div>
              </div>

              <div>
                <p className="text-sm font-medium text-gray-400 mb-2">Total Customers</p>
                <p className="text-3xl font-bold text-white">{stats.customers.current.toLocaleString()}</p>
                <p className="text-sm text-gray-400 mt-2">
                  vs {stats.customers.previous.toLocaleString()} last month
                </p>
              </div>
            </CardContent>
          </Card>

          <Card variant="elevated" className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group bg-gray-800/50 border-gray-700/50">
            <CardContent className="p-8">
              <div className="flex items-center justify-between mb-6">
                <div className="w-16 h-16 rounded-3xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center shadow-large group-hover:scale-110 transition-transform duration-300">
                  <Package className="h-8 w-8 text-white" />
                </div>
                <div className="flex items-center gap-2 px-3 py-1 bg-orange-900/30 rounded-full">
                  <TrendingUp className="h-4 w-4 text-[#FF9000]" />
                  <span className="text-sm font-semibold text-[#FF9000]">+{stats.products.change}%</span>
                </div>
              </div>

              <div>
                <p className="text-sm font-medium text-gray-400 mb-2">Total Products</p>
                <p className="text-3xl font-bold text-white">{stats.products.current}</p>
                <p className="text-sm text-gray-400 mt-2">
                  vs {stats.products.previous} last month
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </RequirePermission>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Enhanced Recent Orders */}
        <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
          <Card variant="elevated" className="border-0 shadow-large bg-gray-800/50 border-gray-700/50">
            <CardHeader className="pb-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                    <ShoppingCart className="h-5 w-5 text-white" />
                  </div>
                  <CardTitle className="text-xl text-white">Recent Orders</CardTitle>
                </div>
                <Button variant="outline" size="sm" className="border-2 border-gray-600 hover:border-[#FF9000] transition-colors text-gray-300 hover:text-white">
                  <Eye className="mr-2 h-4 w-4" />
                  View All
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {recentOrders.length > 0 ? (
                  recentOrders.map((order) => (
                    <div key={order.id} className="flex items-center justify-between p-4 bg-gray-700/30 rounded-2xl hover:bg-gray-700/50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000]/20 to-[#e67e00]/20 flex items-center justify-center border border-[#FF9000]/30">
                          <span className="text-[#FF9000] font-bold text-sm">{order.order_number}</span>
                        </div>
                        <div>
                          <p className="font-semibold text-white">{order.order_number}</p>
                          <p className="text-sm text-gray-400">
                            {order.user ? `${order.user.first_name} ${order.user.last_name}` : 'Guest'}
                          </p>
                          <p className="text-xs text-gray-500">
                            {new Date(order.created_at).toLocaleDateString()}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-lg font-bold text-white">{formatPrice(order.total_amount)}</p>
                        <Badge variant={getStatusColor(order.status) as any} className="text-xs font-semibold">
                          {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                        </Badge>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-center py-8">
                    <div className="w-16 h-16 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-4">
                      <ShoppingCart className="h-8 w-8 text-gray-400" />
                    </div>
                    <h3 className="text-lg font-semibold text-white mb-2">No orders yet</h3>
                    <p className="text-gray-400">When customers place orders, they'll appear here.</p>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </RequirePermission>

        {/* Enhanced Top Products */}
        <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
          <Card variant="elevated" className="border-0 shadow-large bg-gray-800/50 border-gray-700/50">
            <CardHeader className="pb-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                    <Package className="h-5 w-5 text-white" />
                  </div>
                  <CardTitle className="text-xl text-white">Top Products</CardTitle>
                </div>
                <Button variant="outline" size="sm" className="border-2 border-gray-600 hover:border-[#FF9000] transition-colors text-gray-300 hover:text-white">
                  <Eye className="mr-2 h-4 w-4" />
                  View All
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                {topProducts.length > 0 ? (
                  topProducts.map((product, index) => (
                    <div key={product.id} className="flex items-center justify-between p-4 bg-gray-700/30 rounded-2xl hover:bg-gray-700/50 transition-colors">
                      <div className="flex items-center gap-4">
                        <div className="relative">
                          <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center shadow-medium">
                            <span className="text-white font-bold text-lg">#{index + 1}</span>
                          </div>
                          {index === 0 && (
                            <div className="absolute -top-1 -right-1 w-6 h-6 bg-gradient-to-br from-yellow-400 to-yellow-500 rounded-full flex items-center justify-center">
                              <Star className="h-3 w-3 text-white" />
                            </div>
                          )}
                        </div>
                        <div>
                          <p className="font-semibold text-white">{product.name}</p>
                          <p className="text-sm text-gray-400">{product.sales} sales this month</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-lg font-bold text-white">{formatPrice(product.revenue)}</p>
                        <p className="text-sm text-gray-400">Revenue</p>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-center py-8">
                    <div className="w-16 h-16 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-4">
                      <Package className="h-8 w-8 text-gray-400" />
                    </div>
                    <h3 className="text-lg font-semibold text-white mb-2">No products yet</h3>
                    <p className="text-gray-400">Add products to see top performers here.</p>
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </RequirePermission>
      </div>

      {/* Enhanced Quick Actions */}
      <Card variant="elevated" className="border-0 shadow-large bg-gray-800/50 border-gray-700/50">
        <CardHeader className="pb-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
              <TrendingUp className="h-5 w-5 text-white" />
            </div>
            <CardTitle className="text-xl text-white">Quick Actions</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
              <Button
                variant="outline"
                className="h-24 flex-col gap-3 border-2 border-gray-600 hover:border-[#FF9000] hover:bg-[#FF9000]/5 transition-all duration-200 group text-gray-300 hover:text-white"
              >
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <Package className="h-6 w-6 text-white" />
                </div>
                <span className="font-semibold">Add Product</span>
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
              <Button
                variant="outline"
                className="h-24 flex-col gap-3 border-2 border-gray-600 hover:border-[#FF9000] hover:bg-[#FF9000]/5 transition-all duration-200 group text-gray-300 hover:text-white"
              >
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <ShoppingCart className="h-6 w-6 text-white" />
                </div>
                <span className="font-semibold">View Orders</span>
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.USERS_VIEW_ALL}>
              <Button
                variant="outline"
                className="h-24 flex-col gap-3 border-2 border-gray-600 hover:border-[#FF9000] hover:bg-[#FF9000]/5 transition-all duration-200 group text-gray-300 hover:text-white"
              >
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <Users className="h-6 w-6 text-white" />
                </div>
                <span className="font-semibold">Manage Users</span>
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
              <Button
                variant="outline"
                className="h-24 flex-col gap-3 border-2 border-gray-600 hover:border-[#FF9000] hover:bg-[#FF9000]/5 transition-all duration-200 group text-gray-300 hover:text-white"
              >
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
                  <TrendingUp className="h-6 w-6 text-white" />
                </div>
                <span className="font-semibold">View Analytics</span>
              </Button>
            </RequirePermission>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
