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
    const pendingOrders = orders.filter(order => order.status === 'pending').length
    const completedOrders = orders.filter(order => order.status === 'delivered').length
    
    return { totalRevenue, pendingOrders, completedOrders }
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

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <BiHubStatCard
          title="Total Orders"
          value={(pagination as any)?.total_items || 0}
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
              key={(order as any).id}
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
                      #{order.order_number || (order as any).id?.slice(0, 8)}
                    </h3>
                    <BiHubStatusBadge status={getBadgeVariant(order.status)}>
                      {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <p className="text-2xl font-bold text-[#FF9000]">
                      {formatPrice(order.total)}
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
                      {formatPrice(order.total)}
                    </span>
                  </div>
                )}

                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4 text-gray-400" />
                  <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                    {formatDate((order as any).created_at)}
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  <Package className="h-4 w-4 text-gray-400" />
                  <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                    {order.items?.length || order.item_count || 0} items
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
