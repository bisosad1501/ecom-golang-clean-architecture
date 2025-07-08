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
import { useOrders } from '@/hooks/use-orders'
import { Order } from '@/types'
import { RequirePermission } from '@/components/auth/require-permission'
import { PERMISSIONS } from '@/constants/permissions'
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
  const [showOrderModal, setShowOrderModal] = useState(false)
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('list')

  const {
    data: ordersData,
    isLoading,
    error
  } = useOrders({
    page: currentPage,
    limit: 20,
    search: searchQuery,
    status: statusFilter
  })

  const orders = ordersData?.data || []
  const pagination = ordersData?.pagination

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return <Clock className="h-4 w-4 text-white" />
      case 'confirmed':
        return <CheckCircle className="h-4 w-4 text-white" />
      case 'processing':
        return <Package className="h-4 w-4 text-white" />
      case 'shipped':
        return <Truck className="h-4 w-4 text-white" />
      case 'delivered':
        return <CheckCircle className="h-4 w-4 text-white" />
      case 'cancelled':
        return <XCircle className="h-4 w-4 text-white" />
      case 'refunded':
        return <RefreshCw className="h-4 w-4 text-white" />
      default:
        return <ShoppingCart className="h-4 w-4 text-white" />
    }
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'from-yellow-500 to-yellow-600'
      case 'confirmed':
        return 'from-blue-500 to-blue-600'
      case 'processing':
        return 'from-purple-500 to-purple-600'
      case 'shipped':
        return 'from-indigo-500 to-indigo-600'
      case 'delivered':
        return 'from-green-500 to-green-600'
      case 'cancelled':
        return 'from-red-500 to-red-600'
      case 'refunded':
        return 'from-gray-500 to-gray-600'
      default:
        return 'from-gray-500 to-gray-600'
    }
  }

  const handleViewOrder = (order: Order) => {
    setSelectedOrder(order)
    setShowOrderModal(true)
  }

  const handleUpdateOrderStatus = (order: Order, newStatus: string) => {
    // TODO: Implement status update
    console.log('Update order status:', order.id, newStatus)
  }

  const handleExportOrders = () => {
    // TODO: Implement export functionality
    console.log('Export orders')
  }

  const getTotalStats = () => {
    const totalRevenue = orders.reduce((sum, order) => sum + order.total_amount, 0)
    const pendingOrders = orders.filter(order => order.status === 'pending').length
    const completedOrders = orders.filter(order => order.status === 'delivered').length
    
    return { totalRevenue, pendingOrders, completedOrders }
  }

  const stats = getTotalStats()

  return (
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

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <BiHubStatCard
          title="Total Orders"
          value={pagination?.total || 0}
          icon={<ShoppingCart className="h-8 w-8 text-white" />}
          color="primary"
        />
        <BiHubStatCard
          title="Pending Orders"
          value={stats.pendingOrders}
          icon={<Clock className="h-8 w-8 text-white" />}
          color="warning"
        />
        <BiHubStatCard
          title="Completed Orders"
          value={stats.completedOrders}
          icon={<CheckCircle className="h-8 w-8 text-white" />}
          color="success"
        />
        <BiHubStatCard
          title="Total Revenue"
          value={formatPrice(stats.totalRevenue)}
          icon={<DollarSign className="h-8 w-8 text-white" />}
          color="info"
        />
      </div>

      {/* Search & Filters */}
      <BiHubAdminCard
        title="Search & Filter Orders"
        subtitle="Find and filter BiHub orders by status and customer"
        icon={<Search className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className={cn(
                BIHUB_ADMIN_THEME.components.button.ghost,
                'h-10 w-10 p-0'
              )}
            >
              {viewMode === 'grid' ? (
                <List className="h-4 w-4" />
              ) : (
                <Grid className="h-4 w-4" />
              )}
            </Button>
          </div>
        }
      >
        <div className="flex flex-col lg:flex-row items-center gap-4">
          <div className="flex-1 w-full">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                placeholder="Search orders by ID, customer name, or email..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={cn(
                  BIHUB_ADMIN_THEME.components.input.base,
                  'pl-10 pr-12 h-12'
                )}
              />
              {searchQuery && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSearchQuery('')}
                  className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 text-gray-400 hover:text-white"
                >
                  ×
                </Button>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.secondary}>
                  <Filter className="mr-2 h-4 w-4" />
                  Status: {statusFilter || 'All'}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48 bg-gray-900 border-gray-700">
                <DropdownMenuItem
                  onClick={() => setStatusFilter('')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  All Orders
                </DropdownMenuItem>
                <DropdownMenuSeparator className="bg-gray-700" />
                {['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled'].map((status) => (
                  <DropdownMenuItem
                    key={status}
                    onClick={() => setStatusFilter(status)}
                    className="text-gray-300 hover:text-white hover:bg-gray-800"
                  >
                    {status.charAt(0).toUpperCase() + status.slice(1)}
                  </DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </BiHubAdminCard>

      {/* Orders List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className={cn(
              BIHUB_ADMIN_THEME.components.card.base,
              'p-6 animate-pulse'
            )}>
              <div className="flex items-center justify-between mb-4">
                <div className="h-6 bg-gray-700 rounded w-32"></div>
                <div className="h-6 bg-gray-700 rounded w-20"></div>
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-700 rounded w-3/4"></div>
                <div className="h-4 bg-gray-700 rounded w-1/2"></div>
                <div className="h-4 bg-gray-700 rounded w-1/4"></div>
              </div>
            </div>
          ))}
        </div>
      ) : orders.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {orders.map((order) => (
            <div
              key={order.id}
              className={cn(
                BIHUB_ADMIN_THEME.components.card.base,
                BIHUB_ADMIN_THEME.components.card.hover,
                'group',
                viewMode === 'list' && 'flex items-center gap-6 p-6'
              )}
            >
              {/* Order Header */}
              <div className={cn(
                'flex items-center justify-between mb-4',
                viewMode === 'list' && 'flex-1'
              )}>
                <div className="flex items-center gap-3">
                  {/* Status Icon */}
                  <div className={cn(
                    'w-10 h-10 rounded-xl bg-gradient-to-br flex items-center justify-center',
                    getStatusColor(order.status)
                  )}>
                    {getStatusIcon(order.status)}
                  </div>

                  <div>
                    <h3 className={cn(
                      BIHUB_ADMIN_THEME.typography.heading.h4,
                      'group-hover:text-[#FF9000] transition-colors'
                    )}>
                      #{order.order_number || order.id.slice(0, 8)}
                    </h3>
                    <BiHubStatusBadge status={getBadgeVariant(order.status)}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <p className="text-2xl font-bold text-[#FF9000]">
                      {formatPrice(order.total_amount)}
                    </p>
                  </div>
                )}
              </div>

              {/* Order Details */}
              <div className={cn(
                'space-y-3 text-sm',
                viewMode === 'list' && 'flex-1 grid grid-cols-2 md:grid-cols-4 gap-4'
              )}>
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4 text-gray-400" />
                  <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                    {order.user?.first_name} {order.user?.last_name}
                  </span>
                </div>

                {viewMode === 'list' && (
                  <div className="flex items-center gap-2">
                    <DollarSign className="h-4 w-4 text-gray-400" />
                    <span className="font-bold text-[#FF9000]">
                      {formatPrice(order.total_amount)}
                    </span>
                  </div>
                )}

                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-gray-400" />
                  <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                    {formatDate(order.created_at)}
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  <Package className="h-4 w-4 text-gray-400" />
                  <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                    {order.items?.length || 0} items
                  </span>
                </div>

                {order.shipping_address && (
                  <div className="flex items-center gap-2">
                    <MapPin className="h-4 w-4 text-gray-400" />
                    <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                      {order.shipping_address.city}, {order.shipping_address.country}
                    </span>
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className={cn(
                'flex items-center gap-2 mt-4',
                viewMode === 'list' && 'flex-shrink-0 mt-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleViewOrder(order)}
                  className={BIHUB_ADMIN_THEME.components.button.ghost}
                >
                  <Eye className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'View' : ''}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className={BIHUB_ADMIN_THEME.components.button.ghost}
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-48 bg-gray-900 border-gray-700">
                    <DropdownMenuItem
                      onClick={() => handleViewOrder(order)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <Eye className="mr-2 h-4 w-4" />
                      View Details
                    </DropdownMenuItem>
                    <DropdownMenuSeparator className="bg-gray-700" />
                    <DropdownMenuItem
                      onClick={() => handleUpdateOrderStatus(order, 'confirmed')}
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <CheckCircle className="mr-2 h-4 w-4" />
                      Mark as Confirmed
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={() => handleUpdateOrderStatus(order, 'shipped')}
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <Truck className="mr-2 h-4 w-4" />
                      Mark as Shipped
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      onClick={() => handleUpdateOrderStatus(order, 'delivered')}
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <CheckCircle className="mr-2 h-4 w-4" />
                      Mark as Delivered
                    </DropdownMenuItem>
                    <DropdownMenuSeparator className="bg-gray-700" />
                    <DropdownMenuItem
                      onClick={() => handleUpdateOrderStatus(order, 'cancelled')}
                      className="text-red-400 hover:text-red-300 hover:bg-red-900/20"
                    >
                      <XCircle className="mr-2 h-4 w-4" />
                      Cancel Order
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <BiHubEmptyState
          icon={<ShoppingCart className="h-8 w-8 text-gray-400" />}
          title={searchQuery ? 'No orders found' : 'No orders yet'}
          description={
            searchQuery
              ? `No orders found matching "${searchQuery}". Try adjusting your search terms.`
              : 'Orders will appear here once BiHub customers start placing them.'
          }
          action={
            searchQuery && (
              <Button
                onClick={() => setSearchQuery('')}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Clear Search
              </Button>
            )
          }
        />
      )}

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <BiHubAdminCard
          title="Pagination"
          subtitle={`Showing ${((pagination.page - 1) * pagination.limit) + 1} to ${Math.min(pagination.page * pagination.limit, pagination.total)} of ${pagination.total} orders`}
        >
          <div className="flex items-center justify-between">
            <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
              Page {pagination.page} of {pagination.total_pages}
            </p>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={!pagination.has_prev}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Previous
              </Button>

              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                  const pageNum = i + 1;
                  return (
                    <Button
                      key={pageNum}
                      variant={pageNum === pagination.page ? "default" : "ghost"}
                      onClick={() => setCurrentPage(pageNum)}
                      className={cn(
                        'w-10 h-10',
                        pageNum === pagination.page
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
                disabled={!pagination.has_next}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Next
              </Button>
            </div>
          </div>
        </BiHubAdminCard>
      )}
    </div>
  )
}
