'use client'

import { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
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
} from 'lucide-react'
import { useOrders } from '@/hooks/use-orders'
import { Order } from '@/types'
import { RequirePermission } from '@/components/auth/require-permission'
import { PERMISSIONS } from '@/constants/permissions'
import { formatPrice, formatDate } from '@/lib/utils'

export default function AdminOrdersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null)
  const [showOrderModal, setShowOrderModal] = useState(false)

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

  const getStatusVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'secondary'
      case 'confirmed':
        return 'default'
      case 'processing':
        return 'default'
      case 'shipped':
        return 'default'
      case 'delivered':
        return 'default'
      case 'cancelled':
        return 'destructive'
      case 'refunded':
        return 'secondary'
      default:
        return 'outline'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return <Clock className="h-4 w-4" />
      case 'confirmed':
        return <CheckCircle className="h-4 w-4" />
      case 'processing':
        return <Package className="h-4 w-4" />
      case 'shipped':
        return <Truck className="h-4 w-4" />
      case 'delivered':
        return <CheckCircle className="h-4 w-4" />
      case 'cancelled':
        return <XCircle className="h-4 w-4" />
      case 'refunded':
        return <RefreshCw className="h-4 w-4" />
      default:
        return <Clock className="h-4 w-4" />
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
    <div className="space-y-8">
      {/* Enhanced Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
        <div>
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-large">
              <ShoppingCart className="h-6 w-6 text-white" />
            </div>
            <span className="text-primary font-semibold">ORDER MANAGEMENT</span>
          </div>
          
          <h1 className="text-4xl lg:text-5xl font-bold text-foreground mb-4">
            Order <span className="text-gradient">Management</span>
          </h1>
          <p className="text-xl text-muted-foreground">
            Track and manage customer orders, payments, and fulfillment
          </p>
          
          {/* Quick stats */}
          <div className="flex items-center gap-6 mt-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-primary">{pagination?.total || 0}</div>
              <div className="text-sm text-muted-foreground">Total Orders</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-orange-600">{stats.pendingOrders}</div>
              <div className="text-sm text-muted-foreground">Pending</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-emerald-600">{stats.completedOrders}</div>
              <div className="text-sm text-muted-foreground">Completed</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600">{formatPrice(stats.totalRevenue)}</div>
              <div className="text-sm text-muted-foreground">Revenue</div>
            </div>
          </div>
        </div>
        
        <Button 
          onClick={handleExportOrders}
          size="xl"
          variant="outline"
          className="shadow-large hover:shadow-xl transition-all duration-200"
        >
          <Download className="mr-2 h-5 w-5" />
          Export Orders
        </Button>
      </div>

      {/* Enhanced Filters */}
      <Card variant="elevated" className="border-0 shadow-large">
        <CardContent className="p-8">
          <div className="flex flex-col lg:flex-row items-center gap-6">
            <div className="flex-1 w-full">
              <div className="relative">
                <Input
                  placeholder="Search orders by ID, customer name, or email..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  leftIcon={<Search className="h-5 w-5" />}
                  size="lg"
                  className="w-full pr-12"
                />
                {searchQuery && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSearchQuery('')}
                    className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0"
                  >
                    ×
                  </Button>
                )}
              </div>
            </div>
            
            <div className="flex items-center gap-4">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="lg" className="border-2 hover:border-primary transition-colors">
                    <Filter className="mr-2 h-5 w-5" />
                    Status: {statusFilter || 'All'}
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-48">
                  <DropdownMenuItem onClick={() => setStatusFilter('')}>
                    All Orders
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={() => setStatusFilter('pending')}>
                    Pending
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setStatusFilter('confirmed')}>
                    Confirmed
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setStatusFilter('processing')}>
                    Processing
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setStatusFilter('shipped')}>
                    Shipped
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setStatusFilter('delivered')}>
                    Delivered
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setStatusFilter('cancelled')}>
                    Cancelled
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              
              <Button variant="outline" size="lg" className="border-2 hover:border-primary transition-colors">
                Export
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Enhanced Orders List */}
      {isLoading ? (
        <div className="space-y-6">
          {[...Array(5)].map((_, i) => (
            <Card key={i} variant="elevated" className="border-0 shadow-large">
              <CardContent className="p-6">
                <div className="animate-pulse">
                  <div className="flex items-center justify-between mb-4">
                    <div className="h-6 bg-muted rounded w-32"></div>
                    <div className="h-6 bg-muted rounded w-20"></div>
                  </div>
                  <div className="space-y-3">
                    <div className="h-4 bg-muted rounded w-3/4"></div>
                    <div className="h-4 bg-muted rounded w-1/2"></div>
                    <div className="h-4 bg-muted rounded w-1/4"></div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : orders.length > 0 ? (
        <div className="space-y-6">
          {orders.map((order) => (
            <Card 
              key={order.id} 
              variant="elevated" 
              className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group"
            >
              <CardContent className="p-8">
                <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
                  {/* Order Info */}
                  <div className="flex-1">
                    <div className="flex items-center gap-4 mb-4">
                      <h3 className="text-xl font-bold text-foreground">
                        Order #{order.order_number || order.id.slice(0, 8)}
                      </h3>
                      <Badge 
                        variant={getStatusVariant(order.status)} 
                        className="font-semibold"
                      >
                        {getStatusIcon(order.status)}
                        <span className="ml-2">{order.status.charAt(0).toUpperCase() + order.status.slice(1)}</span>
                      </Badge>
                    </div>
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 text-sm">
                      <div className="flex items-center gap-2">
                        <User className="h-4 w-4 text-muted-foreground" />
                        <span className="text-muted-foreground">Customer:</span>
                        <span className="font-medium">{order.user?.first_name} {order.user?.last_name}</span>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <DollarSign className="h-4 w-4 text-muted-foreground" />
                        <span className="text-muted-foreground">Total:</span>
                        <span className="font-bold text-primary">{formatPrice(order.total_amount)}</span>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <Calendar className="h-4 w-4 text-muted-foreground" />
                        <span className="text-muted-foreground">Date:</span>
                        <span className="font-medium">{formatDate(order.created_at)}</span>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        <Package className="h-4 w-4 text-muted-foreground" />
                        <span className="text-muted-foreground">Items:</span>
                        <span className="font-medium">{order.items?.length || 0}</span>
                      </div>
                    </div>
                    
                    {order.shipping_address && (
                      <div className="flex items-center gap-2 mt-3 text-sm">
                        <MapPin className="h-4 w-4 text-muted-foreground" />
                        <span className="text-muted-foreground">Ship to:</span>
                        <span className="font-medium">
                          {order.shipping_address.city}, {order.shipping_address.country}
                        </span>
                      </div>
                    )}
                  </div>
                  
                  {/* Actions */}
                  <div className="flex items-center gap-3">
                    <Button 
                      variant="outline" 
                      size="lg"
                      onClick={() => handleViewOrder(order)}
                      className="border-2 hover:border-primary transition-colors"
                    >
                      <Eye className="h-4 w-4 mr-2" />
                      View Details
                    </Button>
                    
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" size="lg" className="border-2 hover:border-primary transition-colors">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end" className="w-48">
                        <DropdownMenuItem onClick={() => handleViewOrder(order)}>
                          <Eye className="mr-2 h-4 w-4" />
                          View Details
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem onClick={() => handleUpdateOrderStatus(order, 'confirmed')}>
                          <CheckCircle className="mr-2 h-4 w-4" />
                          Mark as Confirmed
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleUpdateOrderStatus(order, 'shipped')}>
                          <Truck className="mr-2 h-4 w-4" />
                          Mark as Shipped
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleUpdateOrderStatus(order, 'delivered')}>
                          <CheckCircle className="mr-2 h-4 w-4" />
                          Mark as Delivered
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem 
                          onClick={() => handleUpdateOrderStatus(order, 'cancelled')}
                          className="text-destructive focus:text-destructive"
                        >
                          <XCircle className="mr-2 h-4 w-4" />
                          Cancel Order
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="text-center py-24">
          <div className="relative mb-8">
            <div className="w-32 h-32 mx-auto rounded-full bg-gradient-to-br from-muted to-muted/50 flex items-center justify-center shadow-large">
              <ShoppingCart className="h-16 w-16 text-muted-foreground" />
            </div>
            <div className="absolute -top-2 -right-2 w-8 h-8 bg-gradient-to-br from-primary to-violet-600 rounded-full flex items-center justify-center shadow-medium">
              <span className="text-white text-sm font-bold">0</span>
            </div>
          </div>
          
          <h3 className="text-3xl font-bold text-foreground mb-4">
            {searchQuery ? 'No orders found' : 'No orders yet'}
          </h3>
          <p className="text-xl text-muted-foreground mb-12 max-w-md mx-auto leading-relaxed">
            {searchQuery 
              ? `No orders found matching "${searchQuery}". Try adjusting your search terms.` 
              : 'Orders will appear here once customers start placing them.'
            }
          </p>
          
          {searchQuery && (
            <Button 
              onClick={() => setSearchQuery('')}
              size="xl"
              variant="outline"
            >
              Clear Search
            </Button>
          )}
        </div>
      )}

      {/* Enhanced Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <Card variant="elevated" className="border-0 shadow-large">
          <CardContent className="p-6">
            <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
              <p className="text-muted-foreground">
                Showing <span className="font-semibold">{((pagination.page - 1) * pagination.limit) + 1}</span> to{' '}
                <span className="font-semibold">{Math.min(pagination.page * pagination.limit, pagination.total)}</span> of{' '}
                <span className="font-semibold">{pagination.total}</span> orders
              </p>
              
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="lg"
                  onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                  disabled={!pagination.has_prev}
                  className="border-2 hover:border-primary transition-colors"
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
                        size="lg"
                        onClick={() => setCurrentPage(pageNum)}
                        className="w-12 h-12"
                      >
                        {pageNum}
                      </Button>
                    );
                  })}
                </div>
                
                <Button
                  variant="outline"
                  size="lg"
                  onClick={() => setCurrentPage(prev => prev + 1)}
                  disabled={!pagination.has_next}
                  className="border-2 hover:border-primary transition-colors"
                >
                  Next
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
