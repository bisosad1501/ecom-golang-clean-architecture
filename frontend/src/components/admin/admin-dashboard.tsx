'use client'

import { 
  TrendingUp, 
  TrendingDown, 
  DollarSign, 
  ShoppingCart, 
  Users, 
  Package,
  Eye,
  MoreHorizontal
} from 'lucide-react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { hasPermission, PERMISSIONS } from '@/lib/permissions'
import { RequirePermission } from '@/components/auth/permission-guard'
import { formatPrice } from '@/lib/utils'

export function AdminDashboard() {
  const { user } = useAuthStore()

  // Mock data - would come from API
  const stats = {
    revenue: {
      current: 45231.89,
      previous: 42150.32,
      change: 7.3,
    },
    orders: {
      current: 1234,
      previous: 1156,
      change: 6.7,
    },
    customers: {
      current: 8945,
      previous: 8721,
      change: 2.6,
    },
    products: {
      current: 567,
      previous: 543,
      change: 4.4,
    },
  }

  const recentOrders = [
    {
      id: 'ORD-001',
      customer: 'John Doe',
      amount: 129.99,
      status: 'completed',
      date: '2024-01-15',
    },
    {
      id: 'ORD-002',
      customer: 'Jane Smith',
      amount: 89.50,
      status: 'processing',
      date: '2024-01-15',
    },
    {
      id: 'ORD-003',
      customer: 'Bob Johnson',
      amount: 199.99,
      status: 'shipped',
      date: '2024-01-14',
    },
    {
      id: 'ORD-004',
      customer: 'Alice Brown',
      amount: 59.99,
      status: 'pending',
      date: '2024-01-14',
    },
  ]

  const topProducts = [
    {
      id: '1',
      name: 'Wireless Headphones',
      sales: 234,
      revenue: 18720,
    },
    {
      id: '2',
      name: 'Smartphone Case',
      sales: 189,
      revenue: 3780,
    },
    {
      id: '3',
      name: 'USB Cable',
      sales: 156,
      revenue: 3120,
    },
    {
      id: '4',
      name: 'Bluetooth Speaker',
      sales: 98,
      revenue: 9800,
    },
  ]

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

  return (
    <div className="space-y-6">
      {/* Welcome Message */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900">
          Welcome back, {user?.first_name}!
        </h1>
        <p className="text-gray-600 mt-2">
          Here's what's happening with your store today.
        </p>
      </div>

      {/* Stats Cards */}
      <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
              <DollarSign className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{formatPrice(stats.revenue.current)}</div>
              <div className="flex items-center text-xs text-muted-foreground">
                <TrendingUp className="mr-1 h-3 w-3 text-green-500" />
                +{stats.revenue.change}% from last month
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Orders</CardTitle>
              <ShoppingCart className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.orders.current.toLocaleString()}</div>
              <div className="flex items-center text-xs text-muted-foreground">
                <TrendingUp className="mr-1 h-3 w-3 text-green-500" />
                +{stats.orders.change}% from last month
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Customers</CardTitle>
              <Users className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.customers.current.toLocaleString()}</div>
              <div className="flex items-center text-xs text-muted-foreground">
                <TrendingUp className="mr-1 h-3 w-3 text-green-500" />
                +{stats.customers.change}% from last month
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Products</CardTitle>
              <Package className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats.products.current}</div>
              <div className="flex items-center text-xs text-muted-foreground">
                <TrendingUp className="mr-1 h-3 w-3 text-green-500" />
                +{stats.products.change}% from last month
              </div>
            </CardContent>
          </Card>
        </div>
      </RequirePermission>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Orders */}
        <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Recent Orders</CardTitle>
                <Button variant="outline" size="sm">
                  <Eye className="mr-2 h-4 w-4" />
                  View All
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {recentOrders.map((order) => (
                  <div key={order.id} className="flex items-center justify-between">
                    <div className="space-y-1">
                      <p className="text-sm font-medium">{order.id}</p>
                      <p className="text-sm text-gray-600">{order.customer}</p>
                    </div>
                    <div className="text-right space-y-1">
                      <p className="text-sm font-medium">{formatPrice(order.amount)}</p>
                      <Badge variant={getStatusColor(order.status) as any} className="text-xs">
                        {order.status}
                      </Badge>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </RequirePermission>

        {/* Top Products */}
        <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Top Products</CardTitle>
                <Button variant="outline" size="sm">
                  <Eye className="mr-2 h-4 w-4" />
                  View All
                </Button>
              </div>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {topProducts.map((product, index) => (
                  <div key={product.id} className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
                        <span className="text-primary-600 text-sm font-medium">
                          {index + 1}
                        </span>
                      </div>
                      <div>
                        <p className="text-sm font-medium">{product.name}</p>
                        <p className="text-sm text-gray-600">{product.sales} sales</p>
                      </div>
                    </div>
                    <p className="text-sm font-medium">{formatPrice(product.revenue)}</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </RequirePermission>
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
              <Button variant="outline" className="h-20 flex-col">
                <Package className="h-6 w-6 mb-2" />
                Add Product
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
              <Button variant="outline" className="h-20 flex-col">
                <ShoppingCart className="h-6 w-6 mb-2" />
                View Orders
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.USERS_VIEW_ALL}>
              <Button variant="outline" className="h-20 flex-col">
                <Users className="h-6 w-6 mb-2" />
                Manage Users
              </Button>
            </RequirePermission>

            <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
              <Button variant="outline" className="h-20 flex-col">
                <TrendingUp className="h-6 w-6 mb-2" />
                View Analytics
              </Button>
            </RequirePermission>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
