'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Package, Truck, CheckCircle, Clock, XCircle, Search, Filter, Calendar, ArrowRight, Eye, Download, Star, ShoppingBag, Home } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { useOrders } from '@/hooks/use-orders'
import { useAuthStore } from '@/store/auth'
import { formatPrice, formatDate } from '@/lib/utils'
import { Order } from '@/types'

const getStatusIcon = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return <Clock className="h-4 w-4" />
    case 'processing':
      return <Package className="h-4 w-4" />
    case 'shipped':
      return <Truck className="h-4 w-4" />
    case 'delivered':
      return <CheckCircle className="h-4 w-4" />
    case 'cancelled':
      return <XCircle className="h-4 w-4" />
    default:
      return <Package className="h-4 w-4" />
  }
}

const getStatusColor = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return 'bg-yellow-500/20 text-yellow-400 border-yellow-500/30'
    case 'processing':
      return 'bg-blue-500/20 text-blue-400 border-blue-500/30'
    case 'shipped':
      return 'bg-purple-500/20 text-purple-400 border-purple-500/30'
    case 'delivered':
      return 'bg-green-500/20 text-green-400 border-green-500/30'
    case 'cancelled':
      return 'bg-red-500/20 text-red-400 border-red-500/30'
    default:
      return 'bg-gray-500/20 text-gray-400 border-gray-500/30'
  }
}

export function OrderHistoryPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState('all')
  const [currentPage, setCurrentPage] = useState(1)

  const { user, isAuthenticated } = useAuthStore()

  const { data: ordersData, isLoading, error } = useOrders({
    page: currentPage,
    limit: 10,
    search: searchQuery,
    status: statusFilter === 'all' ? '' : statusFilter
  })

  const orders = ordersData?.data || []
  const totalPages = ordersData?.pagination?.total_pages || 1

  // Debug logging
  console.log('OrderHistoryPage - user:', user)
  console.log('OrderHistoryPage - isAuthenticated:', isAuthenticated)
  console.log('OrderHistoryPage - ordersData:', ordersData)
  console.log('OrderHistoryPage - orders:', orders)
  console.log('OrderHistoryPage - isLoading:', isLoading)
  console.log('OrderHistoryPage - error:', error)
  console.log('OrderHistoryPage - orders.length:', orders.length)
  console.log('OrderHistoryPage - totalPages:', totalPages)
  console.log('OrderHistoryPage - Query params:', {
    page: currentPage,
    limit: 10,
    search: searchQuery,
    status: statusFilter === 'all' ? '' : statusFilter
  })

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
        <div className="container mx-auto px-4 py-12">
          <div className="space-y-6">
            {[...Array(3)].map((_, i) => (
              <div key={i} className="animate-pulse">
                <div className="h-32 bg-gray-800 rounded-xl"></div>
              </div>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
        <div className="container mx-auto px-4 py-12">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-white mb-4">Error loading orders</h2>
            <p className="text-gray-400 mb-6">We couldn't load your order history. Please try again.</p>
            <Button 
              onClick={() => window.location.reload()} 
              className="bg-orange-500 hover:bg-orange-600"
            >
              Try Again
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      {/* Header */}
      <div className="bg-black/50 backdrop-blur-sm border-b border-gray-700">
        <div className="container mx-auto px-4 py-8">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div>
              <h1 className="text-4xl font-bold text-white mb-2">
                Order History
              </h1>
              <p className="text-gray-400">
                Track your Bi<span className="text-orange-400">Hub</span> orders and view purchase history
              </p>
            </div>
            
            <div className="flex items-center gap-4">
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800" asChild>
                <Link href="/products">
                  <ShoppingBag className="h-4 w-4 mr-2" />
                  Continue Shopping
                </Link>
              </Button>
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800" asChild>
                <Link href="/">
                  <Home className="h-4 w-4 mr-2" />
                  Home
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        {/* Filters */}
        <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 mb-8">
          <CardContent className="p-6">
            <div className="flex flex-col lg:flex-row gap-4">
              <div className="flex-1">
                <div className="relative">
                  <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <Input
                    placeholder="Search orders by number or product..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="pl-10 bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500"
                  />
                </div>
              </div>
              
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="w-full lg:w-48 bg-gray-800/50 border-gray-600 text-white">
                  <SelectValue placeholder="Filter by status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Orders</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="processing">Processing</SelectItem>
                  <SelectItem value="shipped">Shipped</SelectItem>
                  <SelectItem value="delivered">Delivered</SelectItem>
                  <SelectItem value="cancelled">Cancelled</SelectItem>
                </SelectContent>
              </Select>
              
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800">
                <Filter className="h-4 w-4 mr-2" />
                More Filters
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* Orders List */}
        {orders.length === 0 ? (
          <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700">
            <CardContent className="p-12 text-center">
              <div className="w-24 h-24 bg-gray-800 rounded-full flex items-center justify-center mx-auto mb-6">
                <Package className="h-12 w-12 text-gray-400" />
              </div>
              <h3 className="text-2xl font-bold text-white mb-4">No Orders Found</h3>
              <p className="text-gray-400 mb-8">
                You haven't placed any orders yet. Start shopping to see your order history here.
              </p>
              <Button className="bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700" asChild>
                <Link href="/products">
                  <ShoppingBag className="h-4 w-4 mr-2" />
                  Start Shopping
                </Link>
              </Button>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            {orders.map((order: any, index: number) => (
              <Card 
                key={order.id} 
                className="bg-gray-900/50 backdrop-blur-sm border-gray-700 hover:border-orange-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-orange-500/10 group"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <CardContent className="p-6">
                  <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
                    {/* Order Info */}
                    <div className="flex-1">
                      <div className="flex items-center gap-4 mb-4">
                        <div className="w-12 h-12 bg-gradient-to-br from-orange-500 to-orange-600 rounded-xl flex items-center justify-center shadow-lg">
                          {getStatusIcon(order.status)}
                        </div>
                        <div>
                          <h3 className="text-xl font-bold text-white">
                            Order #{order.order_number}
                          </h3>
                          <p className="text-gray-400">
                            Placed on {formatDate(new Date(order.created_at))}
                          </p>
                        </div>
                        <Badge className={getStatusColor(order.status)}>
                          {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                        </Badge>
                      </div>
                      
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                        <div>
                          <p className="text-gray-400">Items</p>
                          <p className="text-white font-semibold">{order.items?.length || order.item_count || 0} products</p>
                        </div>
                        <div>
                          <p className="text-gray-400">Total</p>
                          <p className="text-white font-semibold">{formatPrice(order.total)}</p>
                        </div>
                        <div>
                          <p className="text-gray-400">Payment Status</p>
                          <p className="text-white font-semibold">{order.payment_status}</p>
                        </div>
                      </div>
                    </div>

                    {/* Actions */}
                    <div className="flex flex-col sm:flex-row gap-3">
                      <Button 
                        variant="outline" 
                        className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                        asChild
                      >
                        <Link href={`/orders/${order.id}`}>
                          <Eye className="h-4 w-4 mr-2" />
                          View Details
                        </Link>
                      </Button>
                      
                      {order.status === 'delivered' && (
                        <Button 
                          variant="outline" 
                          className="border-orange-500/50 text-orange-400 hover:bg-orange-500/10 hover:border-orange-500 transition-all duration-300"
                        >
                          <Star className="h-4 w-4 mr-2" />
                          Review
                        </Button>
                      )}
                      
                      <Button 
                        variant="outline" 
                        className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                      >
                        <Download className="h-4 w-4 mr-2" />
                        Invoice
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-12">
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
                className="border-gray-600 text-gray-300 hover:bg-gray-800"
              >
                Previous
              </Button>
              
              {[...Array(totalPages)].map((_, i) => (
                <Button
                  key={i}
                  variant={currentPage === i + 1 ? "default" : "outline"}
                  onClick={() => setCurrentPage(i + 1)}
                  className={currentPage === i + 1 
                    ? "bg-gradient-to-r from-orange-500 to-orange-600" 
                    : "border-gray-600 text-gray-300 hover:bg-gray-800"
                  }
                >
                  {i + 1}
                </Button>
              ))}
              
              <Button
                variant="outline"
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage(currentPage + 1)}
                className="border-gray-600 text-gray-300 hover:bg-gray-800"
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
