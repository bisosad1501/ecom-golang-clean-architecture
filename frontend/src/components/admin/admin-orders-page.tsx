'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import {
  Search,
  Filter,
  MoreHorizontal,
  Eye,
  ShoppingCart,
  Package,
  Truck,
  CheckCircle,
  XCircle,
  Clock,
  DollarSign,
  Calendar,
  User,
  MapPin,
  Download,
  RefreshCw,
  Grid,
  List,
} from 'lucide-react'
import { useAdminOrders, useAdminOrderDetails, useUpdateOrderStatus } from '@/hooks/use-orders'
import { useAdminDashboard } from '@/hooks/use-admin-dashboard'
import { Order } from '@/types'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { formatPrice, formatDate } from '@/lib/utils'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export default function AdminOrdersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [selectedOrderId, setSelectedOrderId] = useState<string>('')
  const [showOrderModal, setShowOrderModal] = useState(false)
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('list')

  const {
    data: ordersData,
    isLoading,
    error
  } = useAdminOrders({
    page: currentPage,
    limit: 20,
    search: searchQuery || undefined,
    status: statusFilter || undefined
  })

  // Get dashboard data for accurate total revenue
  const { data: dashboardData } = useAdminDashboard()

  // Hook for updating order status
  const { mutateAsync: updateOrderStatus } = useUpdateOrderStatus()

  const orders = ordersData?.data || []
  const pagination = ordersData?.pagination

  // Fetch detailed order information when modal is opened
  const {
    data: orderDetails,
    isLoading: isLoadingDetails,
    error: detailsError
  } = useAdminOrderDetails(selectedOrderId)

  // Debug logs
  console.log('AdminOrdersPage - ordersData:', ordersData)
  console.log('AdminOrdersPage - orders:', orders)
  console.log('AdminOrdersPage - pagination:', pagination)
  console.log('AdminOrdersPage - isLoading:', isLoading)
  console.log('AdminOrdersPage - error:', error)

  // Enhanced status icons with better visual hierarchy
  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return <Clock className="h-5 w-5 text-white drop-shadow-sm" />
      case 'confirmed':
        return <CheckCircle className="h-5 w-5 text-white drop-shadow-sm" />
      case 'processing':
        return <Package className="h-5 w-5 text-white drop-shadow-sm" />
      case 'shipped':
        return <Truck className="h-5 w-5 text-white drop-shadow-sm" />
      case 'delivered':
        return <CheckCircle className="h-5 w-5 text-white drop-shadow-sm" />
      case 'cancelled':
        return <XCircle className="h-5 w-5 text-white drop-shadow-sm" />
      case 'refunded':
        return <RefreshCw className="h-5 w-5 text-white drop-shadow-sm" />
      default:
        return <ShoppingCart className="h-5 w-5 text-white drop-shadow-sm" />
    }
  }

  // Enhanced status color system with better differentiation
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'from-amber-500 to-orange-500' // Warm amber for pending
      case 'confirmed':
        return 'from-sky-500 to-blue-600' // Sky blue for confirmed
      case 'processing':
        return 'from-purple-500 to-violet-600' // Purple for processing
      case 'shipped':
        return 'from-teal-500 to-cyan-600' // Teal/cyan for shipping
      case 'delivered':
        return 'from-emerald-500 to-green-600' // Green for success
      case 'cancelled':
        return 'from-red-500 to-rose-600' // Red for cancelled
      case 'refunded':
        return 'from-slate-500 to-gray-600' // Gray for refunded
      default:
        return 'from-gray-500 to-gray-600'
    }
  }

  // Enhanced status styling with better color differentiation
  const getStatusBadgeClass = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'bg-amber-100 text-amber-800 border-amber-200 dark:bg-amber-950/30 dark:text-amber-300 dark:border-amber-800'
      case 'confirmed':
        return 'bg-sky-100 text-sky-800 border-sky-200 dark:bg-sky-950/30 dark:text-sky-300 dark:border-sky-800'
      case 'processing':
        return 'bg-purple-100 text-purple-800 border-purple-200 dark:bg-purple-950/30 dark:text-purple-300 dark:border-purple-800'
      case 'shipped':
        return 'bg-teal-100 text-teal-800 border-teal-200 dark:bg-teal-950/30 dark:text-teal-300 dark:border-teal-800'
      case 'delivered':
        return 'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-800'
      case 'cancelled':
        return 'bg-red-100 text-red-800 border-red-200 dark:bg-red-950/30 dark:text-red-300 dark:border-red-800'
      case 'refunded':
        return 'bg-slate-100 text-slate-800 border-slate-200 dark:bg-slate-950/30 dark:text-slate-300 dark:border-slate-800'
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200 dark:bg-gray-950/30 dark:text-gray-300 dark:border-gray-800'
    }
  }

  const handleViewOrder = (order: Order) => {
    console.log('handleViewOrder - Opening modal for order:', order)
    setSelectedOrder(order)
    setSelectedOrderId((order as any).id)
    setShowOrderModal(true)
  }

  const handleCloseModal = () => {
    console.log('handleCloseModal - Closing modal and resetting state')
    setShowOrderModal(false)
    setSelectedOrder(null)
    setSelectedOrderId('')
  }

  const handleUpdateOrderStatus = async (order: Order, newStatus: string) => {
    try {
      console.log('handleUpdateOrderStatus - Starting update:', (order as any).id, newStatus)
      await updateOrderStatus({ 
        id: (order as any).id, 
        status: newStatus 
      })
      console.log('handleUpdateOrderStatus - Update completed successfully')
      // Data will be automatically updated via React Query invalidation
    } catch (error) {
      console.error('handleUpdateOrderStatus - Failed to update order status:', error)
    }
  }

  const handleExportOrders = async () => {
    try {
      console.log('Exporting orders with filters:', { searchQuery, statusFilter })
      // TODO: Call actual export API
      // await exportOrders({ 
      //   format: 'csv',
      //   search: searchQuery,
      //   status: statusFilter 
      // })
    } catch (error) {
      console.error('Failed to export orders:', error)
    }
  }

  const getTotalStats = () => {
    // Use dashboard data for accurate total revenue across all orders
    const totalRevenue = dashboardData?.overview?.total_revenue || 0
    
    // Calculate statistics for each status
    const pendingOrders = orders.filter(order => order.status === 'pending').length
    const confirmedOrders = orders.filter(order => order.status === 'confirmed').length
    const processingOrders = orders.filter(order => order.status === 'processing').length
    const shippedOrders = orders.filter(order => order.status === 'shipped').length
    const deliveredOrders = orders.filter(order => order.status === 'delivered').length
    const cancelledOrders = orders.filter(order => order.status === 'cancelled').length
    const refundedOrders = orders.filter(order => order.status === 'refunded').length
    
    return { 
      totalRevenue, 
      pendingOrders, 
      confirmedOrders,
      processingOrders,
      shippedOrders,
      deliveredOrders, 
      cancelledOrders,
      refundedOrders
    }
  }

  const stats = getTotalStats()

  return (
    <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
      <div className={BIHUB_ADMIN_THEME.spacing.section}>
        {/* BiHub Page Header */}
        <BiHubPageHeader
          title="Order Management"
          subtitle="Track and manage BiHub customer orders, payments, and fulfillment"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Orders' }
        ]}
        action={
          <Button
            onClick={handleExportOrders}
            className={BIHUB_ADMIN_THEME.components.button.secondary}
          >
            <Download className="mr-2 h-5 w-5" />
            Export Orders
          </Button>
        }
      />

      {/* Enhanced Quick Stats with Modern Design */}
      <div className="space-y-6">
        {/* Primary Stats Row - Modern Glass Design */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          {/* Total Orders */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-blue-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-blue-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-indigo-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-blue-400/80 text-sm font-medium uppercase tracking-wide">Total Orders</p>
                <p className="text-2xl font-bold text-blue-100 mt-1">{(pagination as any)?.total_items || 0}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <ShoppingCart className="h-5 w-5 text-blue-400" />
              </div>
            </div>
          </div>

          {/* Total Revenue */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-green-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-green-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-green-500/5 to-emerald-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-green-400/80 text-sm font-medium uppercase tracking-wide">Total Revenue</p>
                <p className="text-2xl font-bold text-green-100 mt-1">{formatPrice(stats.totalRevenue)}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-green-500/20 to-emerald-600/20 border border-green-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <DollarSign className="h-5 w-5 text-green-400" />
              </div>
            </div>
          </div>

        </div>

        {/* Status Overview - Ultra Compact Cards */}
        <div className="grid grid-cols-3 sm:grid-cols-4 lg:grid-cols-7 gap-2">
          {/* Pending */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-amber-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-amber-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-amber-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 to-orange-500/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <Clock className="h-3 w-3 text-amber-400" />
              </div>
              <p className="text-xs text-amber-400/90 font-medium uppercase tracking-wider">Pending</p>
              <p className="text-base font-bold text-amber-300 mt-0.5">{stats.pendingOrders}</p>
            </div>
          </div>

          {/* Confirmed */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-sky-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-sky-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-sky-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-sky-500/5 to-blue-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-sky-500/20 to-blue-600/20 border border-sky-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <CheckCircle className="h-3 w-3 text-sky-400" />
              </div>
              <p className="text-xs text-sky-400/90 font-medium uppercase tracking-wider">Confirmed</p>
              <p className="text-base font-bold text-sky-300 mt-0.5">{stats.confirmedOrders}</p>
            </div>
          </div>

          {/* Processing */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-purple-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-purple-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-violet-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-purple-500/20 to-violet-600/20 border border-purple-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <Package className="h-3 w-3 text-purple-400" />
              </div>
              <p className="text-xs text-purple-400/90 font-medium uppercase tracking-wider">Processing</p>
              <p className="text-base font-bold text-purple-300 mt-0.5">{stats.processingOrders}</p>
            </div>
          </div>

          {/* Shipped */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-teal-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-teal-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-teal-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-teal-500/5 to-cyan-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-teal-500/20 to-cyan-600/20 border border-teal-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <Truck className="h-3 w-3 text-teal-400" />
              </div>
              <p className="text-xs text-teal-400/90 font-medium uppercase tracking-wider">Shipped</p>
              <p className="text-base font-bold text-teal-300 mt-0.5">{stats.shippedOrders}</p>
            </div>
          </div>

          {/* Delivered */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-emerald-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-emerald-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-emerald-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-green-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <CheckCircle className="h-3 w-3 text-emerald-400" />
              </div>
              <p className="text-xs text-emerald-400/90 font-medium uppercase tracking-wider">Delivered</p>
              <p className="text-base font-bold text-emerald-300 mt-0.5">{stats.deliveredOrders}</p>
            </div>
          </div>

          {/* Cancelled */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-red-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-red-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-red-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-red-500/5 to-rose-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-red-500/20 to-rose-600/20 border border-red-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <XCircle className="h-3 w-3 text-red-400" />
              </div>
              <p className="text-xs text-red-400/90 font-medium uppercase tracking-wider">Cancelled</p>
              <p className="text-base font-bold text-red-300 mt-0.5">{stats.cancelledOrders}</p>
            </div>
          </div>

          {/* Refunded */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-slate-300/20 rounded-lg p-2.5 hover:bg-white/10 hover:border-slate-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-md hover:shadow-slate-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-slate-500/5 to-gray-600/5 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex flex-col items-center text-center">
              <div className="w-6 h-6 rounded-md bg-gradient-to-br from-slate-500/20 to-gray-600/20 border border-slate-400/30 flex items-center justify-center mb-1.5 group-hover:scale-110 transition-transform">
                <RefreshCw className="h-3 w-3 text-slate-400" />
              </div>
              <p className="text-xs text-slate-400/90 font-medium uppercase tracking-wider">Refunded</p>
              <p className="text-base font-bold text-slate-300 mt-0.5">{stats.refundedOrders}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Search & Filters */}
      <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-6 shadow-lg">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-purple-500/5 rounded-2xl"></div>
        <div className="relative">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 border border-blue-400/30 flex items-center justify-center">
                <Search className="h-5 w-5 text-blue-400" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">Search & Filter Orders</h3>
                <p className="text-sm text-gray-400">Find and filter orders by status and customer</p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className="group relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white transition-all duration-200"
            >
              <div className="flex items-center gap-2">
                {viewMode === 'grid' ? (
                  <List className="h-4 w-4" />
                ) : (
                  <Grid className="h-4 w-4" />
                )}
                <span className="text-sm font-medium">
                  {viewMode === 'grid' ? 'List View' : 'Grid View'}
                </span>
              </div>
            </Button>
          </div>

          {/* Search and Filter Controls */}
          <div className="flex flex-col lg:flex-row items-center gap-4">
            <div className="flex-1 w-full">
              <div className="relative group">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 group-focus-within:text-blue-400 transition-colors" />
                <Input
                  placeholder="Search orders by ID, customer name, or email..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full h-12 pl-12 pr-12 bg-white/5 border-gray-600/50 rounded-xl text-white placeholder:text-gray-400 focus:bg-white/10 focus:border-blue-500/50 focus:ring-2 focus:ring-blue-500/20 transition-all duration-200"
                />
                {searchQuery && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSearchQuery('')}
                    className="absolute right-3 top-1/2 -translate-y-1/2 h-6 w-6 p-0 text-gray-400 hover:text-white hover:bg-white/10 rounded-md transition-all duration-200"
                  >
                    ×
                  </Button>
                )}
              </div>
            </div>

            <div className="flex items-center gap-3">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button className="group relative bg-white/5 border border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white px-4 py-2 h-12 rounded-xl transition-all duration-200">
                    <Filter className="mr-2 h-4 w-4" />
                    <span className="font-medium">Status: {statusFilter || 'All'}</span>
                    <div className="ml-2 w-2 h-2 rounded-full bg-blue-400 animate-pulse opacity-60"></div>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                  <DropdownMenuItem
                    onClick={() => setStatusFilter('')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gray-600/50 flex items-center justify-center">
                        <Filter className="h-3 w-3 text-white" />
                      </div>
                      <span className="font-medium">All Orders</span>
                    </div>
                  </DropdownMenuItem>
                  <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                  {['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled'].map((status) => (
                    <DropdownMenuItem
                      key={status}
                      onClick={() => setStatusFilter(status)}
                      className={cn(
                        'text-gray-300 hover:text-white rounded-lg m-1 p-3 transition-all duration-200',
                        statusFilter === status ? 'bg-gray-800/70' : 'hover:bg-gray-800/50'
                      )}
                    >
                      <div className="flex items-center gap-3 w-full">
                        <div className={cn(
                          'w-6 h-6 rounded-lg bg-gradient-to-br flex items-center justify-center border',
                          getStatusColor(status),
                          'border-white/20'
                        )}>
                          <div className="text-white scale-75">
                            {getStatusIcon(status)}
                          </div>
                        </div>
                        <span className="capitalize font-medium flex-1">
                          {status}
                        </span>
                        {statusFilter === status && (
                          <CheckCircle className="h-4 w-4 text-[#FF9000]" />
                        )}
                      </div>
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Orders List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 animate-pulse">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gray-700/50 rounded-xl"></div>
                  <div>
                    <div className="h-5 bg-gray-700/50 rounded w-24 mb-2"></div>
                    <div className="h-4 bg-gray-700/50 rounded w-16"></div>
                  </div>
                </div>
                {viewMode === 'grid' && (
                  <div className="h-6 bg-gray-700/50 rounded w-20"></div>
                )}
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-700/50 rounded w-3/4"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/2"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/4"></div>
              </div>
              <div className="flex items-center gap-3 mt-6 pt-4 border-t border-gray-700/50">
                <div className="h-8 bg-gray-700/50 rounded w-20"></div>
                <div className="h-8 bg-gray-700/50 rounded w-8"></div>
              </div>
            </div>
          ))}
        </div>
      ) : orders.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {orders.map((order) => (
            <div
              key={(order as any).id}
              className={cn(
                'group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 hover:bg-white/10 hover:border-gray-600/50 hover:scale-[1.02] transition-all duration-200 shadow-lg hover:shadow-xl',
                viewMode === 'list' && 'flex items-center gap-6'
              )}
            >
              {/* Gradient Background */}
              <div className={cn(
                'absolute inset-0 bg-gradient-to-br opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-2xl',
                getStatusColor(order.status).includes('amber') ? 'from-amber-500/5 to-orange-500/5' :
                getStatusColor(order.status).includes('sky') ? 'from-sky-500/5 to-blue-600/5' :
                getStatusColor(order.status).includes('purple') ? 'from-purple-500/5 to-violet-600/5' :
                getStatusColor(order.status).includes('teal') ? 'from-teal-500/5 to-cyan-600/5' :
                getStatusColor(order.status).includes('emerald') ? 'from-emerald-500/5 to-green-600/5' :
                getStatusColor(order.status).includes('red') ? 'from-red-500/5 to-rose-600/5' :
                'from-slate-500/5 to-gray-600/5'
              )} />

              {/* Order Header */}
              <div className={cn(
                'relative flex items-center justify-between mb-4',
                viewMode === 'list' && 'flex-1 mb-0'
              )}>
                <div className="flex items-center gap-4">
                  {/* Modern Status Icon */}
                  <div className={cn(
                    'relative w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center shadow-lg border border-white/10',
                    getStatusColor(order.status),
                    order.status === 'pending' || order.status === 'processing' ? 'animate-pulse' : ''
                  )}>
                    <div className="relative z-10 text-white">
                      {getStatusIcon(order.status)}
                    </div>
                    <div className="absolute inset-0 bg-white/10 rounded-xl blur-sm"></div>
                  </div>

                  <div>
                    <h3 className="text-lg font-bold text-white group-hover:text-[#FF9000] transition-colors">
                      #{order.order_number || (order as any).id?.slice(0, 8)}
                    </h3>
                    {/* Modern Status Badge */}
                    <div className={cn(
                      'inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold mt-2 border border-white/10',
                      getStatusBadgeClass(order.status),
                      'backdrop-blur-sm'
                    )}>
                      <div className={cn(
                        'w-2 h-2 rounded-full mr-2',
                        getStatusColor(order.status),
                        order.status === 'pending' || order.status === 'processing' ? 'animate-pulse' : ''
                      )} />
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </div>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <p className="text-2xl font-bold text-[#FF9000] group-hover:text-[#FF9000]/80 transition-colors">
                      {formatPrice(order.total)}
                    </p>
                    <p className="text-xs text-gray-400 mt-1">Total Amount</p>
                  </div>
                )}
              </div>

              {/* Order Details */}
              <div className={cn(
                'relative space-y-3',
                viewMode === 'list' && 'flex-1 grid grid-cols-2 md:grid-cols-4 gap-4 space-y-0'
              )}>
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                    <User className="h-4 w-4 text-gray-400" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 uppercase tracking-wide">Customer</p>
                    <p className="text-sm text-white font-medium">
                      {order.user?.first_name} {order.user?.last_name}
                    </p>
                  </div>
                </div>

                {viewMode === 'list' && (
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                      <DollarSign className="h-4 w-4 text-[#FF9000]" />
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 uppercase tracking-wide">Amount</p>
                      <p className="text-sm font-bold text-[#FF9000]">
                        {formatPrice(order.total)}
                      </p>
                    </div>
                  </div>
                )}

                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                    <Calendar className="h-4 w-4 text-gray-400" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 uppercase tracking-wide">Date</p>
                    <p className="text-sm text-white font-medium">
                      {formatDate((order as any).created_at)}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                    <Package className="h-4 w-4 text-gray-400" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 uppercase tracking-wide">Items</p>
                    <p className="text-sm text-white font-medium">
                      {order.items?.length || order.item_count || 0} items
                    </p>
                  </div>
                </div>

                {order.shipping_address && (
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                      <MapPin className="h-4 w-4 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 uppercase tracking-wide">Location</p>
                      <p className="text-sm text-white font-medium">
                        {order.shipping_address.city}, {order.shipping_address.country}
                      </p>
                    </div>
                  </div>
                )}
              </div>

              {/* Modern Action Buttons */}
              <div className={cn(
                'relative flex items-center gap-3 mt-6 pt-4 border-t border-white/10',
                viewMode === 'list' && 'flex-shrink-0 mt-0 pt-0 border-t-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleViewOrder(order)}
                  className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-blue-500/10 hover:border-blue-500/50 hover:text-blue-400 transition-all duration-200"
                >
                  <Eye className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'View Details' : 'View'}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-purple-500/10 hover:border-purple-500/50 hover:text-purple-400 transition-all duration-200"
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                    <DropdownMenuItem
                      onClick={() => handleViewOrder(order)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Eye className="mr-3 h-4 w-4" />
                      View Details
                    </DropdownMenuItem>
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    
                    {/* Status-specific actions with modern design */}
                    {order.status === 'pending' && (
                      <DropdownMenuItem
                        onClick={() => handleUpdateOrderStatus(order, 'confirmed')}
                        className="text-blue-400 hover:text-blue-300 hover:bg-blue-900/20 rounded-lg m-1 p-3"
                      >
                        <CheckCircle className="mr-3 h-4 w-4" />
                        Confirm Order
                      </DropdownMenuItem>
                    )}
                    
                    {order.status === 'confirmed' && (
                      <DropdownMenuItem
                        onClick={() => handleUpdateOrderStatus(order, 'processing')}
                        className="text-purple-400 hover:text-purple-300 hover:bg-purple-900/20 rounded-lg m-1 p-3"
                      >
                        <Package className="mr-3 h-4 w-4" />
                        Start Processing
                      </DropdownMenuItem>
                    )}
                    
                    {order.status === 'processing' && (
                      <DropdownMenuItem
                        onClick={() => handleUpdateOrderStatus(order, 'shipped')}
                        className="text-teal-400 hover:text-teal-300 hover:bg-teal-900/20 rounded-lg m-1 p-3"
                      >
                        <Truck className="mr-3 h-4 w-4" />
                        Mark as Shipped
                      </DropdownMenuItem>
                    )}
                    
                    {order.status === 'shipped' && (
                      <DropdownMenuItem
                        onClick={() => handleUpdateOrderStatus(order, 'delivered')}
                        className="text-emerald-400 hover:text-emerald-300 hover:bg-emerald-900/20 rounded-lg m-1 p-3"
                      >
                        <CheckCircle className="mr-3 h-4 w-4" />
                        Mark as Delivered
                      </DropdownMenuItem>
                    )}
                    
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    <DropdownMenuItem
                      onClick={() => handleUpdateOrderStatus(order, 'cancelled')}
                      className="text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded-lg m-1 p-3"
                    >
                      <XCircle className="mr-3 h-4 w-4" />
                      Cancel Order
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-12 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-slate-500/5 rounded-2xl"></div>
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-500/20 to-slate-500/20 border border-gray-400/30 flex items-center justify-center mx-auto mb-6">
              <ShoppingCart className="h-8 w-8 text-gray-400" />
            </div>
            <h3 className="text-xl font-bold text-white mb-2">
              {searchQuery ? 'No orders found' : 'No orders yet'}
            </h3>
            <p className="text-gray-400 mb-6 max-w-md mx-auto">
              {searchQuery
                ? `No orders found matching "${searchQuery}". Try adjusting your search terms.`
                : 'Orders will appear here once BiHub customers start placing them.'
              }
            </p>
            {searchQuery && (
              <Button
                onClick={() => setSearchQuery('')}
                className="bg-blue-500/20 border border-blue-400/30 text-blue-400 hover:bg-blue-500/30 hover:border-blue-400/50 hover:text-blue-300 transition-all duration-200"
              >
                Clear Search
              </Button>
            )}
          </div>
        </div>
      )}

      {/* Pagination */}
      {pagination && (pagination as any).total_pages > 1 && (
        <BiHubAdminCard
          title="Pagination"
          subtitle={`Showing ${(((pagination as any).current_page - 1) * (pagination as any).per_page) + 1} to ${Math.min((pagination as any).current_page * (pagination as any).per_page, (pagination as any).total_items)} of ${(pagination as any).total_items} orders`}
        >
          <div className="flex items-center justify-between">
            <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
              Page {(pagination as any).current_page} of {(pagination as any).total_pages}
            </p>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={!(pagination as any).has_prev}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Previous
              </Button>

              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, (pagination as any).total_pages) }, (_, i) => {
                  const pageNum = i + 1;
                  return (
                    <Button
                      key={pageNum}
                      variant={pageNum === (pagination as any).current_page ? "default" : "ghost"}
                      onClick={() => setCurrentPage(pageNum)}
                      className={cn(
                        'w-10 h-10',
                        pageNum === (pagination as any).current_page
                          ? BIHUB_ADMIN_THEME.components.button.primary
                          : BIHUB_ADMIN_THEME.components.button.ghost
                      )}
                    >
                      {pageNum}
                    </Button>
                  );
                })}
              </div>

              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => prev + 1)}
                disabled={!(pagination as any).has_next}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Next
              </Button>
            </div>
          </div>
        </BiHubAdminCard>
      )}

      {/* Order Details Modal */}
      <Dialog open={showOrderModal} onOpenChange={handleCloseModal}>
        <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto bg-gray-900 border-gray-700">
          <DialogHeader>
            <DialogTitle className="text-white">
              Order Details - #{orderDetails?.order_number || selectedOrder?.order_number || (selectedOrder as any)?.id?.slice(0, 8)}
            </DialogTitle>
            <DialogDescription className="text-gray-400">
              Complete order information and customer details
            </DialogDescription>
          </DialogHeader>

          {isLoadingDetails ? (
            <div className="flex items-center justify-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[#FF9000]"></div>
              <span className="ml-2 text-gray-400">Loading order details...</span>
            </div>
          ) : detailsError ? (
            <div className="text-center py-8">
              <XCircle className="h-8 w-8 text-red-400 mx-auto mb-2" />
              <p className="text-red-400">Failed to load order details</p>
              <p className="text-gray-400 text-sm mt-1">{(detailsError as any)?.message}</p>
              <Button 
                variant="outline" 
                className="mt-4" 
                onClick={() => window.location.reload()}
              >
                Try Again
              </Button>
            </div>
          ) : orderDetails ? (
            <div className="space-y-6">
              {/* Order Overview */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <BiHubStatCard
                  title="Order Total"
                  value={formatPrice(orderDetails.total)}
                  icon={<DollarSign className="h-6 w-6 text-white" />}
                  color="primary"
                />
                <BiHubStatCard
                  title="Items Count"
                  value={orderDetails.items?.length || orderDetails.item_count || 0}
                  icon={<Package className="h-6 w-6 text-white" />}
                  color="info"
                />
                <BiHubStatCard
                  title="Order Date"
                  value={formatDate((orderDetails as any).created_at)}
                  icon={<Calendar className="h-6 w-6 text-white" />}
                  color="info"
                />
              </div>

              {/* Order Status */}
              <BiHubAdminCard title="Order Status" icon={<Package className="h-5 w-5 text-white" />}>
                <div className="flex items-center gap-4">
                  <div className="flex items-center gap-2">
                    <span className="text-gray-400">Status:</span>
                    <BiHubStatusBadge status={getBadgeVariant(orderDetails.status)}>
                      {orderDetails.status.charAt(0).toUpperCase() + orderDetails.status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="text-gray-400">Payment:</span>
                    <BiHubStatusBadge status={getBadgeVariant(orderDetails.payment_status)}>
                      {orderDetails.payment_status.charAt(0).toUpperCase() + orderDetails.payment_status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                </div>
              </BiHubAdminCard>

              {/* Customer Information */}
              <BiHubAdminCard title="Customer Information" icon={<User className="h-5 w-5 text-white" />}>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Name:</span>
                    <span className="text-white">
                      {orderDetails.user?.first_name} {orderDetails.user?.last_name}
                    </span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Email:</span>
                    <span className="text-white">{orderDetails.user?.email}</span>
                  </div>
                  {orderDetails.user?.id && (
                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Customer ID:</span>
                      <span className="text-white font-mono text-sm">
                        {orderDetails.user.id.slice(0, 8)}...
                      </span>
                    </div>
                  )}
                </div>
              </BiHubAdminCard>

              {/* Order Items */}
              <BiHubAdminCard title="Order Items" icon={<Package className="h-5 w-5 text-white" />}>
                {orderDetails.items && orderDetails.items.length > 0 ? (
                  <div className="space-y-3">
                    {orderDetails.items.map((item, index) => (
                      <div key={item.id || index} className="p-4 bg-gray-800 rounded-lg">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-medium text-white">{item.product_name}</h4>
                            <p className="text-sm text-gray-400">SKU: {item.product_sku}</p>
                          </div>
                          <div className="text-right">
                            <p className="text-white">
                              {item.quantity} × {formatPrice(item.price)}
                            </p>
                            <p className="text-[#FF9000] font-bold">
                              {formatPrice(item.total)}
                            </p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-4">
                    <Package className="h-8 w-8 text-gray-400 mx-auto mb-2" />
                    <p className="text-gray-400">No items found for this order</p>
                  </div>
                )}
              </BiHubAdminCard>

              {/* Shipping Address */}
              {orderDetails.shipping_address ? (
                <BiHubAdminCard title="Shipping Address" icon={<MapPin className="h-5 w-5 text-white" />}>
                  <div className="space-y-2">
                    <div className="text-white">
                      {orderDetails.shipping_address.first_name} {orderDetails.shipping_address.last_name}
                    </div>
                    <div className="text-gray-400">
                      {orderDetails.shipping_address.address_line_1}
                      {orderDetails.shipping_address.address_line_2 && (
                        <>
                          <br />
                          {orderDetails.shipping_address.address_line_2}
                        </>
                      )}
                    </div>
                    <div className="text-gray-400">
                      {orderDetails.shipping_address.city}, {orderDetails.shipping_address.state} {orderDetails.shipping_address.postal_code}
                    </div>
                    <div className="text-gray-400">
                      {orderDetails.shipping_address.country}
                    </div>
                    {orderDetails.shipping_address.phone && (
                      <div className="text-gray-400">
                        Phone: {orderDetails.shipping_address.phone}
                      </div>
                    )}
                  </div>
                </BiHubAdminCard>
              ) : (
                <BiHubAdminCard title="Shipping Address" icon={<MapPin className="h-5 w-5 text-white" />}>
                  <div className="text-center py-4">
                    <MapPin className="h-8 w-8 text-gray-400 mx-auto mb-2" />
                    <p className="text-gray-400">No shipping address found</p>
                  </div>
                </BiHubAdminCard>
              )}

              {/* Order Summary */}
              <BiHubAdminCard title="Order Summary" icon={<DollarSign className="h-5 w-5 text-white" />}>
                <div className="space-y-2">
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Subtotal:</span>
                    <span className="text-white">{formatPrice(orderDetails.subtotal)}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Tax:</span>
                    <span className="text-white">{formatPrice(orderDetails.tax_amount)}</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-gray-400">Shipping:</span>
                    <span className="text-white">{formatPrice(orderDetails.shipping_amount)}</span>
                  </div>
                  {orderDetails.discount_amount > 0 && (
                    <div className="flex items-center justify-between">
                      <span className="text-gray-400">Discount:</span>
                      <span className="text-green-400">-{formatPrice(orderDetails.discount_amount)}</span>
                    </div>
                  )}
                  <hr className="border-gray-700" />
                  <div className="flex items-center justify-between text-lg font-bold">
                    <span className="text-white">Total:</span>
                    <span className="text-[#FF9000]">{formatPrice(orderDetails.total)}</span>
                  </div>
                </div>
              </BiHubAdminCard>
            </div>
          ) : selectedOrder ? (
            <div className="text-center py-8">
              <Package className="h-8 w-8 text-gray-400 mx-auto mb-2" />
              <p className="text-gray-400">Unable to load detailed order information</p>
              <p className="text-gray-500 text-sm mt-1">
                Showing basic information from order list
              </p>
              
              {/* Basic order info from list */}
              <div className="mt-6 space-y-4 text-left max-w-md mx-auto">
                <div className="p-4 bg-gray-800 rounded-lg">
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-400">Order Number:</span>
                    <span className="text-white font-mono">
                      #{selectedOrder.order_number || (selectedOrder as any)?.id?.slice(0, 8)}
                    </span>
                  </div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-400">Status:</span>
                    <BiHubStatusBadge status={getBadgeVariant(selectedOrder.status)}>
                      {selectedOrder.status.charAt(0).toUpperCase() + selectedOrder.status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                  <div className="flex justify-between items-center mb-2">
                    <span className="text-gray-400">Total:</span>
                    <span className="text-[#FF9000] font-bold">{formatPrice(selectedOrder.total)}</span>
                  </div>
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Items:</span>
                    <span className="text-white">{selectedOrder.items?.length || selectedOrder.item_count || 0}</span>
                  </div>
                </div>
              </div>
            </div>
          ) : null}
        </DialogContent>
      </Dialog>
      </div>
    </RequirePermission>
  )
}
