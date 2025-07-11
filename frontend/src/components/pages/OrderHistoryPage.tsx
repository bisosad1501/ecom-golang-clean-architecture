'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Package, Truck, CheckCircle, Clock, XCircle, Search, Eye, Download, Star, ShoppingBag, SlidersHorizontal, ShieldCheck, Zap } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useOrders } from '@/hooks/use-orders'
import { useAuthStore } from '@/store/auth'
import { formatPrice, formatDate, cn } from '@/lib/utils'
import { Order } from '@/types'

const getStatusIcon = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return <Clock className="h-4 w-4 text-yellow-300" />
    case 'confirmed':
      return <ShieldCheck className="h-4 w-4 text-cyan-300" />
    case 'processing':
      return <Zap className="h-4 w-4 text-blue-300" />
    case 'ready_to_ship':
      return <Package className="h-4 w-4 text-orange-300" />
    case 'shipped':
      return <Truck className="h-4 w-4 text-purple-300" />
    case 'out_for_delivery':
      return <Truck className="h-4 w-4 text-indigo-300" />
    case 'delivered':
      return <CheckCircle className="h-4 w-4 text-green-300" />
    case 'cancelled':
      return <XCircle className="h-4 w-4 text-red-300" />
    case 'refunded':
      return <XCircle className="h-4 w-4 text-gray-300" />
    case 'returned':
      return <XCircle className="h-4 w-4 text-orange-300" />
    case 'exchanged':
      return <XCircle className="h-4 w-4 text-blue-300" />
    default:
      return <Package className="h-4 w-4 text-gray-300" />
  }
}

const getStatusColor = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return 'bg-yellow-500/20 text-yellow-300 border-yellow-500/40 shadow-sm'
    case 'confirmed':
      return 'bg-cyan-500/20 text-cyan-300 border-cyan-500/40 shadow-sm'
    case 'processing':
      return 'bg-blue-500/20 text-blue-300 border-blue-500/40 shadow-sm'
    case 'ready_to_ship':
      return 'bg-orange-500/20 text-orange-300 border-orange-500/40 shadow-sm'
    case 'shipped':
      return 'bg-purple-500/20 text-purple-300 border-purple-500/40 shadow-sm'
    case 'out_for_delivery':
      return 'bg-indigo-500/20 text-indigo-300 border-indigo-500/40 shadow-sm'
    case 'delivered':
      return 'bg-green-500/20 text-green-300 border-green-500/40 shadow-sm'
    case 'cancelled':
      return 'bg-red-500/20 text-red-300 border-red-500/40 shadow-sm'
    case 'refunded':
      return 'bg-gray-500/20 text-gray-300 border-gray-500/40 shadow-sm'
    case 'returned':
      return 'bg-orange-500/20 text-orange-300 border-orange-500/40 shadow-sm'
    case 'exchanged':
      return 'bg-blue-500/20 text-blue-300 border-blue-500/40 shadow-sm'
    default:
      return 'bg-gray-500/20 text-gray-300 border-gray-500/40 shadow-sm'
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
    search: searchQuery.trim(),
    status: statusFilter === 'all' ? undefined : statusFilter
  })

  const orders = ordersData?.data || []
  const totalPages = ordersData?.pagination?.total_pages || 1
  const totalOrders = ordersData?.pagination?.total || orders.length

  // Backend now handles filtering, so we use the filtered results directly
  const filteredOrders = orders

  // Debug logging
  console.log('OrderHistoryPage - ordersData:', ordersData)
  console.log('OrderHistoryPage - totalPages:', totalPages)
  console.log('OrderHistoryPage - totalOrders:', totalOrders)
  console.log('OrderHistoryPage - currentPage:', currentPage)

  // Reset page when filters change
  const handleSearchChange = (value: string) => {
    setSearchQuery(value)
    setCurrentPage(1)
  }

  const handleStatusChange = (value: string) => {
    setStatusFilter(value)
    setCurrentPage(1)
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 py-12 relative z-10">
          <div className="space-y-6">
            {[...Array(3)].map((_, i) => (
              <Card key={i} className={cn(
                'animate-pulse backdrop-blur-sm border transition-all duration-300',
                'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
                'border-gray-700/50 rounded-2xl shadow-xl shadow-black/40'
              )}>
                <CardContent className="p-6">
                  <div className="flex items-center gap-4 mb-4">
                    <div className="w-12 h-12 bg-gradient-to-br from-gray-600/50 to-gray-700/50 rounded-xl"></div>
                    <div className="flex-1">
                      <div className="h-5 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded mb-2"></div>
                      <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-2/3"></div>
                    </div>
                  </div>
                  <div className="grid grid-cols-3 gap-4">
                    <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded"></div>
                    <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded"></div>
                    <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded"></div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 py-12 relative z-10">
          <Card className={cn(
            'p-8 text-center backdrop-blur-sm border transition-all duration-300',
            'bg-gradient-to-br from-slate-900/90 via-gray-900/95 to-slate-800/90',
            'border-gray-700/50 hover:border-red-500/30 rounded-2xl',
            'shadow-xl shadow-black/40 hover:shadow-red-500/10'
          )}>
            <div className="w-16 h-16 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
              <XCircle className="h-8 w-8 text-white" />
            </div>
            <h2 className="text-2xl font-bold text-white mb-4">Error loading orders</h2>
            <p className="text-gray-300 mb-6">We couldn't load your order history. Please try again.</p>
            <Button 
              onClick={() => window.location.reload()} 
              className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg"
            >
              Try Again
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      {/* Enhanced Background Pattern */}
      <AnimatedBackground className="opacity-30" />
      
      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-12 relative z-10">
        {/* Header Section trong body */}
        <div className="mb-12">
          <h1 className="text-3xl lg:text-4xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight mb-3">
            Order <span className="text-[#ff9000]">History</span>
          </h1>
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2">
            <p className="text-gray-300 text-lg">
              Track your Bi<span className="text-[#ff9000]">Hub</span> orders and view purchase history
            </p>
            <div className="flex items-center gap-2 text-sm text-gray-400">
              <span>Total Orders:</span>
              <span className="text-[#ff9000] font-semibold">{totalOrders}</span>
              {statusFilter !== 'all' && (
                <>
                  <span>•</span>
                  <span>Status: {statusFilter}</span>
                </>
              )}
            </div>
          </div>
        </div>
        {/* Compact Search & Filters */}
        <div className="mb-8 flex flex-col sm:flex-row gap-4">
          {/* Search Bar */}
          <div className="flex-1">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search by order number or product..."
                value={searchQuery}
                onChange={(e) => handleSearchChange(e.target.value)}
                className="pl-10 h-10 bg-white/[0.08] border border-white/15 text-white placeholder:text-gray-400 focus:border-[#ff9000]/50 focus:ring-[#ff9000]/20 rounded-lg transition-all duration-300"
              />
            </div>
          </div>
          
          {/* Filters */}
          <div className="flex gap-3">
            <Select value={statusFilter} onValueChange={handleStatusChange}>
              <SelectTrigger className="w-[140px] h-10 bg-white/[0.08] border border-white/15 text-white rounded-lg hover:border-[#ff9000]/30 transition-colors">
                <SelectValue />
              </SelectTrigger>
              <SelectContent className="bg-gray-800 border-gray-600 shadow-xl backdrop-blur-sm">
                <SelectItem value="all" className="text-white hover:bg-[#ff9000]/20 hover:text-[#ff9000] focus:bg-[#ff9000]/20 focus:text-[#ff9000] cursor-pointer">
                  All Status
                </SelectItem>
                <SelectItem value="pending" className="text-white hover:bg-yellow-500/20 hover:text-yellow-300 focus:bg-yellow-500/20 focus:text-yellow-300 cursor-pointer">
                  Pending
                </SelectItem>
                <SelectItem value="confirmed" className="text-white hover:bg-cyan-500/20 hover:text-cyan-300 focus:bg-cyan-500/20 focus:text-cyan-300 cursor-pointer">
                  Confirmed
                </SelectItem>
                <SelectItem value="processing" className="text-white hover:bg-blue-500/20 hover:text-blue-300 focus:bg-blue-500/20 focus:text-blue-300 cursor-pointer">
                  Processing
                </SelectItem>
                <SelectItem value="shipped" className="text-white hover:bg-purple-500/20 hover:text-purple-300 focus:bg-purple-500/20 focus:text-purple-300 cursor-pointer">
                  Shipped
                </SelectItem>
                <SelectItem value="delivered" className="text-white hover:bg-green-500/20 hover:text-green-300 focus:bg-green-500/20 focus:text-green-300 cursor-pointer">
                  Delivered
                </SelectItem>
                <SelectItem value="cancelled" className="text-white hover:bg-red-500/20 hover:text-red-300 focus:bg-red-500/20 focus:text-red-300 cursor-pointer">
                  Cancelled
                </SelectItem>
              </SelectContent>
            </Select>

            {/* Clear Filters Button */}
            {(searchQuery || statusFilter !== 'all') && (
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => {
                  setSearchQuery('')
                  setStatusFilter('all')
                  setCurrentPage(1)
                }}
                className="h-10 px-3 border-gray-600/50 text-gray-300 hover:bg-gray-800/50 hover:border-red-500/30 hover:text-red-400 transition-all duration-300 rounded-lg"
              >
                <XCircle className="h-4 w-4 mr-1" />
                Clear
              </Button>
            )}
          </div>
        </div>

        {/* Filter Summary */}
        {(searchQuery || statusFilter !== 'all') && (
          <div className="mb-6 flex items-center gap-2 text-sm text-gray-400">
            <span>Showing results</span>
            {statusFilter !== 'all' && (
              <Badge variant="outline" className="text-xs border-gray-600 text-gray-300">
                Status: {statusFilter}
              </Badge>
            )}
            {searchQuery && (
              <Badge variant="outline" className="text-xs border-gray-600 text-gray-300">
                Search: "{searchQuery}"
              </Badge>
            )}
            <span>• {filteredOrders.length} order{filteredOrders.length !== 1 ? 's' : ''} found</span>
          </div>
        )}

        {/* Orders List */}
        {filteredOrders.length === 0 ? (
          <Card className={cn(
            'backdrop-blur-sm border transition-all duration-300',
            'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
            'border-gray-700/50 hover:border-[#ff9000]/30 rounded-2xl',
            'shadow-xl shadow-black/40 hover:shadow-[#ff9000]/10'
          )}>
            <CardContent className="p-12 text-center">
              <div className="w-24 h-24 bg-gradient-to-br from-[#ff9000] to-orange-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
                <Package className="h-12 w-12 text-white" />
              </div>
              <h3 className="text-2xl font-bold text-white mb-4">
                {(searchQuery || statusFilter !== 'all') ? 'No Matching Orders Found' : 'No Orders Found'}
              </h3>
              <p className="text-gray-300 mb-8">
                {(searchQuery || statusFilter !== 'all') 
                  ? 'Try adjusting your search criteria or filters to find what you\'re looking for.'
                  : 'You haven\'t placed any orders yet. Start shopping to see your order history here.'
                }
              </p>
              {(searchQuery || statusFilter !== 'all') ? (
                <Button 
                  onClick={() => {
                    setSearchQuery('')
                    setStatusFilter('all')
                    setCurrentPage(1)
                  }}
                  className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg"
                >
                  <XCircle className="h-4 w-4 mr-2" />
                  Clear All Filters
                </Button>
              ) : (
                <Button className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg" asChild>
                  <Link href="/products">
                    <ShoppingBag className="h-4 w-4 mr-2" />
                    Start Shopping
                  </Link>
                </Button>
              )}
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-6">
            {filteredOrders.map((order: any, index: number) => (
              <Card 
                key={order.id} 
                className={cn(
                  'group relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-300 ease-out',
                  'bg-gradient-to-br from-slate-900/95 via-gray-900/90 to-slate-800/95',
                  'hover:shadow-xl hover:shadow-[#ff9000]/20 hover:scale-[1.02]',
                  'rounded-xl border-gray-700/40 hover:border-[#ff9000]/50',
                  'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/[0.03] before:via-transparent before:to-white/[0.01] before:pointer-events-none before:rounded-xl'
                )}
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <CardContent className="p-6">
                  {/* Order Header with Actions */}
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-4 flex-1 min-w-0">
                      <div className={cn(
                        "w-14 h-14 rounded-xl flex items-center justify-center shadow-lg group-hover:shadow-lg transition-all duration-300 flex-shrink-0",
                        order.status.toLowerCase() === 'pending' && "bg-gradient-to-br from-yellow-500/20 to-yellow-600/30 shadow-yellow-500/20",
                        order.status.toLowerCase() === 'confirmed' && "bg-gradient-to-br from-cyan-500/20 to-cyan-600/30 shadow-cyan-500/20",
                        order.status.toLowerCase() === 'processing' && "bg-gradient-to-br from-blue-500/20 to-blue-600/30 shadow-blue-500/20",
                        order.status.toLowerCase() === 'shipped' && "bg-gradient-to-br from-purple-500/20 to-purple-600/30 shadow-purple-500/20",
                        order.status.toLowerCase() === 'delivered' && "bg-gradient-to-br from-green-500/20 to-green-600/30 shadow-green-500/20",
                        order.status.toLowerCase() === 'cancelled' && "bg-gradient-to-br from-red-500/20 to-red-600/30 shadow-red-500/20",
                        (!['pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled'].includes(order.status.toLowerCase())) && "bg-gradient-to-br from-gray-500/20 to-gray-600/30 shadow-gray-500/20"
                      )}>
                        {getStatusIcon(order.status)}
                      </div>
                      
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center justify-between mb-2">
                          <div className="flex items-center gap-4 min-w-0 flex-1">
                            <h3 className="text-xl font-bold text-white truncate">
                              Order #{order.order_number}
                            </h3>
                            
                            {/* Action Buttons next to title */}
                            <div className="flex items-center gap-2 flex-shrink-0">
                              <Button 
                                variant="outline" 
                                size="sm"
                                className="h-7 px-2 bg-white/[0.05] border-gray-600/40 text-gray-300 hover:bg-[#ff9000]/15 hover:border-[#ff9000]/60 hover:text-[#ff9000] transition-all duration-300 rounded-md text-xs"
                                asChild
                              >
                                <Link href={`/orders/${order.id}`}>
                                  <Eye className="h-3 w-3 mr-1" />
                                  Details
                                </Link>
                              </Button>
                              
                              {order.status === 'delivered' && (
                                <Button 
                                  size="sm"
                                  className="h-7 px-2 bg-gradient-to-r from-[#ff9000]/20 to-orange-600/20 border border-[#ff9000]/40 text-[#ff9000] hover:from-[#ff9000]/30 hover:to-orange-600/30 hover:border-[#ff9000]/70 transition-all duration-300 rounded-md text-xs"
                                >
                                  <Star className="h-3 w-3 mr-1" />
                                  Review
                                </Button>
                              )}
                              
                              <Button 
                                variant="outline" 
                                size="sm"
                                className="h-7 px-2 bg-white/[0.05] border-gray-600/40 text-gray-300 hover:bg-gray-800/50 hover:border-gray-500/50 hover:text-white transition-all duration-300 rounded-md text-xs"
                              >
                                <Download className="h-3 w-3 mr-1" />
                                Invoice
                              </Button>
                            </div>
                          </div>
                          
                          <Badge className={cn(
                            getStatusColor(order.status),
                            'px-3 py-1 text-xs font-medium capitalize flex-shrink-0 ml-3'
                          )}>
                            {order.status}
                          </Badge>
                        </div>
                        
                        <p className="text-gray-400 text-sm mb-3">
                          {formatDate(new Date(order.created_at))}
                        </p>
                        
                        {/* Order Details */}
                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3 text-sm">
                          <div className="flex items-center gap-2">
                            <Package className="h-4 w-4 text-[#ff9000] flex-shrink-0" />
                            <span className="text-gray-400">Items:</span>
                            <span className="text-white font-medium">{order.items?.length || order.item_count || 0}</span>
                          </div>
                          
                          <div className="flex items-center gap-2">
                            <div className="w-4 h-4 rounded-full bg-green-500 flex items-center justify-center flex-shrink-0">
                              <span className="text-white text-xs font-bold">$</span>
                            </div>
                            <span className="text-gray-400">Total:</span>
                            <span className="text-white font-bold">{formatPrice(order.total)}</span>
                          </div>
                          
                          <div className="flex items-center gap-2">
                            <CheckCircle className="h-4 w-4 text-blue-400 flex-shrink-0" />
                            <span className="text-gray-400">Payment:</span>
                            <span className={cn(
                              "font-medium capitalize",
                              order.payment_status === 'paid' ? 'text-green-400' : 
                              order.payment_status === 'pending' ? 'text-yellow-400' : 'text-gray-400'
                            )}>
                              {order.payment_status}
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}

        {/* Enhanced Pagination - đồng bộ với ProductsPage */}
        {totalPages > 1 && (
          <div className="flex justify-center mt-12">
            <div className="flex items-center gap-2 bg-white/[0.06] backdrop-blur-md border border-white/10 rounded-lg p-1 shadow-lg">
              <Button
                variant="ghost"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(currentPage - 1)}
                className="border-0 text-gray-300 hover:text-[#ff9000] hover:bg-[#ff9000]/10 disabled:text-gray-500 disabled:hover:bg-transparent rounded-md h-8 px-3 text-sm transition-all duration-200"
              >
                Previous
              </Button>
              
              {[...Array(totalPages)].map((_, i) => (
                <Button
                  key={i}
                  variant={currentPage === i + 1 ? "default" : "ghost"}
                  onClick={() => setCurrentPage(i + 1)}
                  className={currentPage === i + 1 
                    ? "bg-gradient-to-r from-[#ff9000] to-orange-600 text-white border-0 h-8 w-8 p-0 text-sm font-medium shadow-md" 
                    : "border-0 text-gray-300 hover:text-[#ff9000] hover:bg-[#ff9000]/10 h-8 w-8 p-0 text-sm transition-all duration-200 rounded-md"
                  }
                >
                  {i + 1}
                </Button>
              ))}
              
              <Button
                variant="ghost"
                disabled={currentPage === totalPages}
                onClick={() => setCurrentPage(currentPage + 1)}
                className="border-0 text-gray-300 hover:text-[#ff9000] hover:bg-[#ff9000]/10 disabled:text-gray-500 disabled:hover:bg-transparent rounded-md h-8 px-3 text-sm transition-all duration-200"
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
