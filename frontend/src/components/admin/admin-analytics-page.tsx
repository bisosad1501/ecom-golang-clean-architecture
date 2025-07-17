'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  BarChart3,
  PieChart,
  TrendingUp,
  TrendingDown,
  DollarSign,
  ShoppingCart,
  Users,
  Package,
  Calendar,
  Download,
  Filter,
  ArrowUpRight,
  ArrowDownRight,
  Target,
  Activity,
  Clock,
  Eye,
  Zap,
} from 'lucide-react'
import { PermissionGuard } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/constants/permissions'
import { formatPrice, formatDate } from '@/lib/utils'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export default function AdminAnalyticsPage() {
  const [timeRange, setTimeRange] = useState('7d')
  const [isLoading, setIsLoading] = useState(false)

  // Mock analytics data - replace with real API calls
  const analyticsData = {
    revenue: {
      current: 125000,
      previous: 98000,
      change: 27.6,
      isPositive: true,
      data: [
        { date: '2024-01-01', value: 15000 },
        { date: '2024-01-02', value: 18000 },
        { date: '2024-01-03', value: 22000 },
        { date: '2024-01-04', value: 19000 },
        { date: '2024-01-05', value: 25000 },
        { date: '2024-01-06', value: 21000 },
        { date: '2024-01-07', value: 28000 },
      ]
    },
    orders: {
      current: 1250,
      previous: 980,
      change: 27.6,
      isPositive: true,
    },
    customers: {
      current: 850,
      previous: 720,
      change: 18.1,
      isPositive: true,
    },
    conversion: {
      current: 3.2,
      previous: 2.8,
      change: 14.3,
      isPositive: true,
    },
    topProducts: [
      { name: 'Gaming Laptop Pro', sales: 145, revenue: 289000, growth: 23.5 },
      { name: 'Wireless Headphones', sales: 320, revenue: 96000, growth: 18.2 },
      { name: 'Smart Watch Ultra', sales: 280, revenue: 84000, growth: 15.7 },
      { name: 'Mechanical Keyboard', sales: 190, revenue: 38000, growth: 12.3 },
      { name: 'Gaming Mouse RGB', sales: 240, revenue: 24000, growth: 8.9 },
    ],
    salesByCategory: [
      { category: 'Electronics', value: 45, color: 'from-blue-500 to-blue-600' },
      { category: 'Gaming', value: 30, color: 'from-purple-500 to-purple-600' },
      { category: 'Accessories', value: 15, color: 'from-green-500 to-green-600' },
      { category: 'Software', value: 10, color: 'from-yellow-500 to-yellow-600' },
    ],
    recentActivity: [
      { type: 'sale', message: 'New sale: Gaming Laptop Pro', amount: 1999, time: '2 min ago' },
      { type: 'customer', message: 'New customer registration', time: '5 min ago' },
      { type: 'order', message: 'Order #1234 completed', amount: 299, time: '8 min ago' },
      { type: 'review', message: 'New 5-star review received', time: '12 min ago' },
      { type: 'sale', message: 'New sale: Wireless Headphones', amount: 299, time: '15 min ago' },
    ]
  }

  const getTimeRangeLabel = (range: string) => {
    switch (range) {
      case '24h': return 'Last 24 Hours'
      case '7d': return 'Last 7 Days'
      case '30d': return 'Last 30 Days'
      case '90d': return 'Last 90 Days'
      case '1y': return 'Last Year'
      default: return 'Last 7 Days'
    }
  }

  const getActivityIcon = (type: string) => {
    switch (type) {
      case 'sale': return <DollarSign className="h-4 w-4" />
      case 'customer': return <Users className="h-4 w-4" />
      case 'order': return <ShoppingCart className="h-4 w-4" />
      case 'review': return <Target className="h-4 w-4" />
      default: return <Activity className="h-4 w-4" />
    }
  }

  const getActivityColor = (type: string) => {
    switch (type) {
      case 'sale': return 'bg-green-600'
      case 'customer': return 'bg-blue-600'
      case 'order': return 'bg-purple-600'
      case 'review': return 'bg-yellow-600'
      default: return 'bg-gray-600'
    }
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Analytics & Reports"
        subtitle="Comprehensive insights into your BiHub store performance"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Analytics' }
        ]}
        action={
          <div className="flex items-center gap-4">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.secondary}>
                  <Calendar className="mr-2 h-4 w-4" />
                  {getTimeRangeLabel(timeRange)}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="bg-gray-900 border-gray-700">
                {['24h', '7d', '30d', '90d', '1y'].map((range) => (
                  <DropdownMenuItem
                    key={range}
                    onClick={() => setTimeRange(range)}
                    className="text-gray-300 hover:text-white hover:bg-gray-800"
                  >
                    {getTimeRangeLabel(range)}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
            
            <Button className={BIHUB_ADMIN_THEME.components.button.primary}>
              <Download className="mr-2 h-4 w-4" />
              Export Report
            </Button>
          </div>
        }
      />

      {/* Key Metrics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <BiHubStatCard
          title="Total Revenue"
          value={formatPrice(analyticsData.revenue.current)}
          change={analyticsData.revenue.change}
          isPositive={analyticsData.revenue.isPositive}
          icon={<DollarSign className="h-8 w-8 text-white" />}
          color="primary"
        />
        <BiHubStatCard
          title="Total Orders"
          value={analyticsData.orders.current}
          change={analyticsData.orders.change}
          isPositive={analyticsData.orders.isPositive}
          icon={<ShoppingCart className="h-8 w-8 text-white" />}
          color="success"
        />
        <BiHubStatCard
          title="New Customers"
          value={analyticsData.customers.current}
          change={analyticsData.customers.change}
          isPositive={analyticsData.customers.isPositive}
          icon={<Users className="h-8 w-8 text-white" />}
          color="info"
        />
        <BiHubStatCard
          title="Conversion Rate"
          value={`${analyticsData.conversion.current}%`}
          change={analyticsData.conversion.change}
          isPositive={analyticsData.conversion.isPositive}
          icon={<Target className="h-8 w-8 text-white" />}
          color="warning"
        />
      </div>

      {/* Charts and Analytics */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Revenue Chart */}
        <BiHubAdminCard
          title="Revenue Trend"
          subtitle="Revenue performance over time"
          icon={<BarChart3 className="h-5 w-5 text-white" />}
          headerAction={
            <Button
              variant="outline"
              size="sm"
              className={BIHUB_ADMIN_THEME.components.button.ghost}
            >
              <Eye className="h-4 w-4 mr-2" />
              View Details
            </Button>
          }
        >
          <div className="h-64 flex items-center justify-center bg-gray-800 rounded-lg">
            <div className="text-center">
              <BarChart3 className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
                Revenue Chart
              </p>
              <p className={BIHUB_ADMIN_THEME.typography.body.small}>
                Chart component would be integrated here
              </p>
            </div>
          </div>
        </BiHubAdminCard>

        {/* Sales by Category */}
        <BiHubAdminCard
          title="Sales by Category"
          subtitle="Product category performance breakdown"
          icon={<PieChart className="h-5 w-5 text-white" />}
        >
          <div className="space-y-4">
            {analyticsData.salesByCategory.map((category, index) => (
              <div key={index} className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={cn(
                    'w-4 h-4 rounded-full bg-gradient-to-r',
                    category.color
                  )}></div>
                  <span className={BIHUB_ADMIN_THEME.typography.body.medium}>
                    {category.category}
                  </span>
                </div>
                <div className="flex items-center gap-2">
                  <span className="text-lg font-bold text-[#FF9000]">
                    {category.value}%
                  </span>
                  <div className="w-20 h-2 bg-gray-700 rounded-full overflow-hidden">
                    <div 
                      className={cn('h-full bg-gradient-to-r', category.color)}
                      style={{ width: `${category.value}%` }}
                    ></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      </div>

      {/* Bottom Section */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Products */}
        <BiHubAdminCard
          title="Top Performing Products"
          subtitle="Best selling products this period"
          icon={<Package className="h-5 w-5 text-white" />}
        >
          <div className="space-y-4">
            {analyticsData.topProducts.map((product, index) => (
              <div key={index} className="flex items-center justify-between p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
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
            ))}
          </div>
        </BiHubAdminCard>

        {/* Recent Activity */}
        <BiHubAdminCard
          title="Recent Activity"
          subtitle="Latest activities in your store"
          icon={<Activity className="h-5 w-5 text-white" />}
        >
          <div className="space-y-4">
            {analyticsData.recentActivity.map((activity, index) => (
              <div key={index} className="flex items-center gap-3 p-3 bg-gray-800 rounded-lg">
                <div className={cn(
                  'w-8 h-8 rounded-lg flex items-center justify-center text-white',
                  getActivityColor(activity.type)
                )}>
                  {getActivityIcon(activity.type)}
                </div>
                <div className="flex-1">
                  <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
                    {activity.message}
                  </p>
                  <div className="flex items-center gap-2">
                    <p className={cn(BIHUB_ADMIN_THEME.typography.body.small, 'text-gray-500')}>
                      {activity.time}
                    </p>
                    {activity.amount && (
                      <>
                        <span className="text-gray-500">•</span>
                        <span className="text-sm font-bold text-[#FF9000]">
                          {formatPrice(activity.amount)}
                        </span>
                      </>
                    )}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      </div>
    </div>
  )
}
