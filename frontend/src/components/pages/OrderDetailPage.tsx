'use client'

import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { ArrowLeft, Package, Truck, CheckCircle, Clock, XCircle, MapPin, CreditCard, Download, Star, MessageCircle, RefreshCw, Copy, Check, ShieldCheck, Zap, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { InvoicePreviewModal } from '@/components/modals/InvoicePreviewModal'
import { OrderTimeline } from '@/components/ui/order-timeline'
import { useOrder, useOrderEvents } from '@/hooks/use-orders'
import { formatPrice, formatDate, cn } from '@/lib/utils'
import { toast } from 'sonner'

interface Props {
  orderId: string
}

const getStatusIcon = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return <Clock className="h-5 w-5 text-yellow-300" />
    case 'confirmed':
      return <ShieldCheck className="h-5 w-5 text-cyan-300" />
    case 'processing':
      return <Zap className="h-5 w-5 text-blue-300" />
    case 'ready_to_ship':
      return <Package className="h-5 w-5 text-orange-300" />
    case 'shipped':
      return <Truck className="h-5 w-5 text-purple-300" />
    case 'out_for_delivery':
      return <MapPin className="h-5 w-5 text-indigo-300" />
    case 'delivered':
      return <CheckCircle className="h-5 w-5 text-green-300" />
    case 'cancelled':
      return <XCircle className="h-5 w-5 text-red-300" />
    case 'refunded':
      return <RefreshCw className="h-5 w-5 text-gray-300" />
    case 'returned':
      return <RefreshCw className="h-5 w-5 text-orange-300" />
    case 'exchanged':
      return <RefreshCw className="h-5 w-5 text-blue-300" />
    default:
      return <Package className="h-5 w-5 text-gray-300" />
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

export function OrderDetailPage({ orderId }: Props) {
  const [copiedTrackingId, setCopiedTrackingId] = useState(false)
  const [showInvoiceModal, setShowInvoiceModal] = useState(false)
  const { data: order, isLoading, error } = useOrder(orderId)
  const { data: events = [] } = useOrderEvents(orderId)

  const copyTrackingNumber = () => {
    if (order?.tracking_number) {
      navigator.clipboard.writeText(order.tracking_number)
      setCopiedTrackingId(true)
      toast.success('Tracking number copied!')
      setTimeout(() => setCopiedTrackingId(false), 2000)
    }
  }

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 py-12 relative z-10">
          <div className="space-y-6">
            {[...Array(4)].map((_, i) => (
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
        <div className="container mx-auto px-4 py-12 relative z-10 flex items-center justify-center">
          <Card className={cn(
            'p-8 text-center backdrop-blur-sm border transition-all duration-300 max-w-md w-full',
            'bg-gradient-to-br from-slate-900/90 via-gray-900/95 to-slate-800/90',
            'border-gray-700/50 hover:border-red-500/30 rounded-2xl',
            'shadow-xl shadow-black/40 hover:shadow-red-500/10'
          )}>
            <div className="w-16 h-16 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
              <XCircle className="h-8 w-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-white mb-4">Order Not Found</h3>
            <p className="text-gray-300 mb-6">The order you're looking for doesn't exist or you don't have permission to view it.</p>
            <Button 
              asChild 
              className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg"
            >
              <Link href="/orders">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Orders
              </Link>
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 py-12 relative z-10 flex items-center justify-center">
          <Card className={cn(
            'p-8 text-center backdrop-blur-sm border transition-all duration-300 max-w-md w-full',
            'bg-gradient-to-br from-slate-900/90 via-gray-900/95 to-slate-800/90',
            'border-gray-700/50 hover:border-yellow-500/30 rounded-2xl',
            'shadow-xl shadow-black/40 hover:shadow-yellow-500/10'
          )}>
            <div className="w-16 h-16 bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-lg">
              <XCircle className="h-8 w-8 text-white" />
            </div>
            <h3 className="text-2xl font-bold text-white mb-4">Order Not Found</h3>
            <p className="text-gray-300 mb-6">The order data could not be loaded. Please try refreshing the page.</p>
            <Button 
              asChild 
              className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg"
            >
              <Link href="/orders">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Orders
              </Link>
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
        {/* Compact Header */}
        <div className="mb-6">
          <div className="flex items-center justify-between mb-3">
            <Button 
              variant="ghost" 
              asChild 
              size="sm"
              className="text-gray-300 hover:text-[#ff9000] hover:bg-[#ff9000]/10 transition-all duration-300"
            >
              <Link href="/orders">
                <ArrowLeft className="h-4 w-4 mr-1" />
                Back
              </Link>
            </Button>

            <div className="flex items-center gap-2">
              {order?.status && (
                <Badge className={cn(
                  getStatusColor(order.status),
                  'text-sm px-3 py-1 font-medium capitalize flex items-center gap-1.5'
                )}>
                  {getStatusIcon(order.status)}
                  <span>{order.status.charAt(0).toUpperCase() + order.status.slice(1)}</span>
                </Badge>
              )}
              
              <Button 
                variant="outline" 
                size="sm"
                onClick={() => setShowInvoiceModal(true)}
                className="h-8 px-2 border-gray-600/50 text-gray-300 hover:bg-gray-800/50 hover:border-[#ff9000]/30 hover:text-[#ff9000] transition-all duration-300"
              >
                <Download className="h-3 w-3 mr-1" />
                Invoice
              </Button>
            </div>
          </div>
          
          <div>
            <h1 className="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
              Order <span className="text-[#ff9000]">#{order?.order_number || 'Loading...'}</span>
            </h1>
            <p className="text-gray-400 text-sm mt-1">
              {order?.created_at ? `${formatDate(new Date(order.created_at))}` : 'Loading...'}
            </p>
          </div>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Content - Order Items */}
          <div className="lg:col-span-2">
            {/* Order Items */}
            <Card className={cn(
              'backdrop-blur-sm border transition-all duration-300',
              'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
              'border-gray-700/50 hover:border-[#ff9000]/30 rounded-xl',
              'shadow-lg shadow-black/40 hover:shadow-[#ff9000]/10'
            )}>
              <CardHeader className="pb-3">
                <CardTitle className="flex items-center gap-3 text-white text-xl font-semibold">
                  <div className="w-10 h-10 bg-gradient-to-br from-[#ff9000]/20 to-orange-600/30 rounded-lg flex items-center justify-center">
                    <Package className="h-5 w-5 text-[#ff9000]" />
                  </div>
                  Order Items ({order?.items?.length || order?.item_count || 0})
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <div className="space-y-4">
                  {order?.items?.map((item, index) => (
                    <div key={item.id} className="relative group">
                      {/* Refined outer glow effect */}
                      <div className={cn(
                        'absolute -inset-0.5 rounded-2xl opacity-0 group-hover:opacity-40 transition-all duration-700 ease-out',
                        'bg-gradient-to-br from-[#ff9000]/15 via-orange-500/8 to-amber-400/10 blur-lg'
                      )} />
                      
                      <div className={cn(
                        'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-300 ease-out',
                        'bg-gradient-to-br from-white/[0.02] via-white/[0.03] to-white/[0.02]',
                        'hover:shadow-lg hover:shadow-[#ff9000]/8 hover:-translate-y-0.5',
                        'rounded-2xl backdrop-saturate-150 border-gray-700/40 hover:border-[#ff9000]/30',
                        'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/1 before:via-transparent before:to-white/0.5 before:pointer-events-none before:rounded-2xl',
                        'p-4'
                      )}>
                        <div className="flex gap-4 items-start">
                          {/* Enhanced Product Image */}
                          <div className="relative flex-shrink-0">
                            <div className="relative w-36 h-36 rounded-xl overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 shadow-lg">
                              {/* Product discount badge if available */}
                              {item.product?.compare_price && item.price < item.product.compare_price && (
                                <div className="absolute top-2 left-2 z-20">
                                  <span className="text-xs font-bold text-white bg-[#ff9000] px-2 py-1 rounded-md shadow-md">
                                    -{Math.round(((item.product.compare_price - item.price) / item.product.compare_price) * 100)}%
                                  </span>
                                </div>
                              )}

                              <Image
                                src={item.product?.images?.[0]?.url || '/placeholder-product.svg'}
                                alt={item.product?.name || item.product_name || 'Product'}
                                fill
                                className="object-cover transition-all duration-300 ease-out group-hover:scale-105"
                                sizes="(max-width: 144px) 100vw, 144px"
                              />
                              
                              {/* Quick View overlay */}
                              <div className={cn(
                                'absolute inset-0 flex items-center justify-center transition-all duration-300',
                                'bg-black/60 backdrop-blur-sm opacity-0 group-hover:opacity-100'
                              )}>
                                <Button
                                  variant="ghost"
                                  size="icon"
                                  className="h-9 w-9 bg-white/90 hover:bg-white text-slate-700 hover:text-blue-600 rounded-lg"
                                  asChild
                                >
                                  <Link href={`/products/${item.product?.id}`}>
                                    <Eye className="h-4 w-4" />
                                  </Link>
                                </Button>
                              </div>
                            </div>
                          </div>

                          {/* Product Info Section - Exactly like cart */}
                          <div className="flex-1 min-w-0 space-y-2">
                            {/* Category Row */}
                            <div className="flex items-center gap-2">
                              {item.product?.category && (
                                <span className="text-sm font-semibold text-[#ff9000] bg-[#ff9000]/15 px-3 py-1.5 rounded-lg border border-[#ff9000]/40 shadow-sm">
                                  {item.product.category.name}
                                </span>
                              )}
                            </div>

                            {/* Product Name - Much Larger */}
                            <h3 className="text-xl font-bold leading-tight line-clamp-2 text-white group-hover:text-[#ff9000]/90 transition-colors">
                              <Link href={`/products/${item.product?.id}`} className="hover:text-[#ff9000] transition-colors">
                                {item.product?.name || item.product_name}
                              </Link>
                            </h3>

                            {/* Price Section - Larger and more prominent */}
                            <div className="flex items-center gap-3">
                              <span className="text-2xl font-bold text-white">
                                {formatPrice(item.price)}
                              </span>
                              {item.product?.has_discount && item.product?.original_price && (
                                <span className="text-lg line-through text-gray-500">
                                  {formatPrice(item.product.original_price)}
                                </span>
                              )}
                            </div>
                          </div>

                          {/* Controls Section - Exactly like cart */}
                          <div className="flex flex-col items-end gap-3 flex-shrink-0 min-w-[200px]">
                            {/* Quantity Controls - Larger and cleaner */}
                            <div className="flex items-center gap-2">
                              <span className="text-sm text-gray-400 font-medium whitespace-nowrap">Qty:</span>
                              <div className="flex items-center bg-white/[0.08] rounded-lg border border-white/20 shadow-sm">
                                <span className="px-4 py-2 text-base font-bold text-white min-w-[3rem] text-center">
                                  {item.quantity}
                                </span>
                              </div>
                            </div>

                            {/* Item Subtotal - Clear but less prominent than final total */}
                            <div className="text-right">
                              <div className="text-sm text-gray-400 mb-1">Subtotal</div>
                              <div className="text-xl font-bold text-white">
                                {formatPrice(item.price * item.quantity)}
                              </div>
                            </div>

                            {/* Action Button - More prominent */}
                            {order?.status === 'delivered' && (
                              <Button 
                                variant="outline" 
                                size="sm" 
                                className="h-10 px-4 text-[#ff9000] hover:text-red-400 hover:bg-[#ff9000]/15 border border-[#ff9000]/50 hover:border-[#ff9000] rounded-lg transition-colors font-medium"
                              >
                                <Star className="h-4 w-4 mr-2" />
                                Review
                              </Button>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                {/* Order Notes & Delivery Info */}
                {(order?.estimated_delivery || order?.customer_notes) && (
                  <div className="mt-6 space-y-3">
                    {order?.estimated_delivery && (
                      <div className="flex items-center gap-3 p-3 bg-green-500/10 border border-green-500/30 rounded-lg">
                        <CheckCircle className="h-5 w-5 text-green-400 flex-shrink-0" />
                        <div>
                          <p className="text-green-400 font-medium">Estimated Delivery</p>
                          <p className="text-gray-300 text-sm">{formatDate(new Date(order.estimated_delivery))}</p>
                        </div>
                      </div>
                    )}

                    {order?.customer_notes && (
                      <div className="p-4 bg-blue-500/10 border border-blue-500/30 rounded-lg">
                        <p className="text-blue-300 font-medium mb-2">Order Notes</p>
                        <p className="text-gray-300 text-sm">{order.customer_notes}</p>
                      </div>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Sidebar - Order Details & Summary */}
          <div className="space-y-4">
            {/* Order Summary */}
            <Card className={cn(
              'backdrop-blur-sm border transition-all duration-300',
              'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
              'border-gray-700/50 hover:border-[#ff9000]/30 rounded-xl',
              'shadow-lg shadow-black/40 hover:shadow-[#ff9000]/10'
            )}>
              <CardHeader className="pb-3">
                <CardTitle className="text-white flex items-center gap-2 text-xl font-semibold">
                  <Package className="h-4 w-4 text-[#ff9000]" />
                  Order Summary
                </CardTitle>
              </CardHeader>
              <CardContent className="pt-0">
                <div className="space-y-3">
                  <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-400">Subtotal</span>
                    <span className="text-white font-semibold">{formatPrice(order?.subtotal || 0)}</span>
                  </div>
                  
                  <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-400">Shipping</span>
                    <span className="text-white font-semibold">
                      {(order?.shipping_amount || 0) === 0 ? (
                        <span className="text-green-400 font-bold">FREE</span>
                      ) : (
                        formatPrice(order?.shipping_amount || 0)
                      )}
                    </span>
                  </div>
                  
                  <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-400">Tax</span>
                    <span className="text-white font-semibold">{formatPrice(order?.tax_amount || 0)}</span>
                  </div>
                  
                  {/* Debug: Let's see what's causing the 0 */}
                  {/* {JSON.stringify({
                    discount_amount: order?.discount_amount,
                    hasDiscount: order?.discount_amount && order.discount_amount > 0
                  })} */}
                  
                  {/* Fixed: Use proper conditional rendering to avoid 0 being displayed */}
                  {order?.discount_amount && order.discount_amount > 0 ? (
                    <div className="flex justify-between items-center text-sm">
                      <span className="text-gray-400">Discount</span>
                      <span className="text-green-400 font-semibold">-{formatPrice(order.discount_amount)}</span>
                    </div>
                  ) : null}
                  
                  <div className="h-px bg-gray-700/50 my-3"></div>
                  
                  <div className="flex justify-between items-center py-3 bg-[#ff9000]/10 border border-[#ff9000]/30 rounded-lg px-3 backdrop-blur-sm">
                    <span className="text-white font-bold text-lg">Total</span>
                    <span className="text-[#ff9000] font-bold text-xl">{formatPrice(order?.total || 0)}</span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Shipping & Payment Info */}
            <Card className={cn(
              'backdrop-blur-sm border transition-all duration-300',
              'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
              'border-gray-700/50 hover:border-[#ff9000]/30 rounded-xl',
              'shadow-lg shadow-black/40 hover:shadow-[#ff9000]/10'
            )}>
              <CardContent className="p-4">
                <div className="space-y-4">
                  {/* Shipping Address */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <MapPin className="h-4 w-4 text-blue-300" />
                      <h3 className="text-white font-semibold text-xl">Shipping Address</h3>
                    </div>
                    {order?.shipping_address ? (
                      <div className="text-gray-300 text-sm space-y-1 pl-6">
                        <p className="font-medium text-white">
                          {order.shipping_address.first_name} {order.shipping_address.last_name}
                        </p>
                        {order.shipping_address.company && <p>{order.shipping_address.company}</p>}
                        <p>{order.shipping_address.address1}</p>
                        {order.shipping_address.address2 && <p>{order.shipping_address.address2}</p>}
                        <p>{order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.zip_code}</p>
                        <p>{order.shipping_address.country}</p>
                        {order.shipping_address.phone && <p>📞 {order.shipping_address.phone}</p>}
                      </div>
                    ) : (
                      <p className="text-gray-400 text-sm pl-6">No shipping address available</p>
                    )}
                  </div>

                  <div className="h-px bg-gray-700/50"></div>

                  {/* Payment Method */}
                  <div>
                    <div className="flex items-center gap-2 mb-3">
                      <CreditCard className="h-4 w-4 text-green-300" />
                      <h3 className="text-white font-semibold text-xl">Payment Details</h3>
                    </div>
                    <div className="text-gray-300 text-sm space-y-2 pl-6">
                      <div className="flex justify-between items-center">
                        <span>Status:</span>
                        <span className={cn(
                          "font-semibold",
                          order?.payment_status === 'paid' ? 'text-green-400' : 'text-yellow-400'
                        )}>
                          {order?.payment_status?.charAt(0).toUpperCase() + (order?.payment_status?.slice(1) || '')}
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span>Method:</span>
                        <span>{order?.payment?.method || 'Stripe'}</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span>Currency:</span>
                        <span>{order?.currency || 'USD'}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Tracking Information */}
            {(order?.tracking_number || order?.carrier || order?.shipping_method) && (
              <Card className={cn(
                'backdrop-blur-sm border transition-all duration-300',
                'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
                'border-gray-700/50 hover:border-[#ff9000]/30 rounded-xl',
                'shadow-lg shadow-black/40 hover:shadow-[#ff9000]/10'
              )}>
                <CardHeader className="pb-3">
                  <CardTitle className="text-white flex items-center gap-2 text-xl font-semibold">
                    <Truck className="h-4 w-4 text-[#ff9000]" />
                    Tracking Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="pt-0 space-y-3">
                  {order?.tracking_number && (
                    <div className="space-y-2">
                      <div className="flex justify-between items-center">
                        <span className="text-gray-400 text-sm">Tracking Number</span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={copyTrackingNumber}
                          className="h-auto p-1 text-blue-400 hover:text-blue-300"
                        >
                          {copiedTrackingId ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                        </Button>
                      </div>
                      <div className="text-white font-mono text-sm bg-gray-800/50 p-2 rounded border border-gray-700">
                        {order.tracking_number}
                      </div>
                    </div>
                  )}

                  {order?.carrier && (
                    <div className="space-y-1">
                      <span className="text-gray-400 text-sm">Carrier</span>
                      <div className="text-white font-medium">{order.carrier}</div>
                    </div>
                  )}

                  {order?.shipping_method && (
                    <div className="space-y-1">
                      <span className="text-gray-400 text-sm">Shipping Method</span>
                      <div className="text-white font-medium">{order.shipping_method}</div>
                    </div>
                  )}

                  {order?.estimated_delivery && (
                    <div className="space-y-1">
                      <span className="text-gray-400 text-sm">Estimated Delivery</span>
                      <div className="text-white font-medium">{formatDate(order.estimated_delivery)}</div>
                    </div>
                  )}

                  {order?.actual_delivery && (
                    <div className="space-y-1">
                      <span className="text-gray-400 text-sm">Delivered On</span>
                      <div className="text-green-400 font-medium">{formatDate(order.actual_delivery)}</div>
                    </div>
                  )}

                  {order?.tracking_url && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => window.open(order.tracking_url, '_blank')}
                      className="w-full border-purple-500/50 text-purple-400 hover:bg-purple-500/10 hover:border-purple-500/70 transition-all duration-300"
                    >
                      <MapPin className="h-4 w-4 mr-2" />
                      Track Package
                    </Button>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Quick Actions */}
            <Card className={cn(
              'backdrop-blur-sm border transition-all duration-300',
              'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
              'border-gray-700/50 hover:border-[#ff9000]/30 rounded-xl',
              'shadow-lg shadow-black/40 hover:shadow-[#ff9000]/10'
            )}>
              <CardHeader className="pb-2">
                <CardTitle className="text-white text-xl font-semibold">Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="pt-0 space-y-2">
                <Button 
                  variant="outline" 
                  size="sm"
                  className="w-full border-gray-600/50 text-gray-300 hover:bg-gray-800/50 hover:border-[#ff9000]/30 hover:text-[#ff9000] transition-all duration-300"
                >
                  <MessageCircle className="h-4 w-4 mr-2" />
                  Contact Support
                </Button>
                
                <Button 
                  variant="outline" 
                  size="sm"
                  className="w-full border-gray-600/50 text-gray-300 hover:bg-gray-800/50 hover:border-[#ff9000]/30 hover:text-[#ff9000] transition-all duration-300"
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Reorder Items
                </Button>

                {order?.can_be_cancelled && order?.status !== 'delivered' && order?.status !== 'cancelled' && (
                  <Button 
                    variant="outline" 
                    size="sm"
                    className="w-full border-red-500/50 text-red-400 hover:bg-red-500/10 hover:border-red-500/70 transition-all duration-300"
                  >
                    <XCircle className="h-4 w-4 mr-2" />
                    Cancel Order
                  </Button>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Order Timeline */}
          <div className="lg:col-span-2">
            <OrderTimeline
              orderId={orderId}
              events={events}
              showPrivateEvents={false}
            />
          </div>
        </div>
      </div>

      {/* Invoice Preview Modal */}
      {order && (
        <InvoicePreviewModal
          isOpen={showInvoiceModal}
          onClose={() => setShowInvoiceModal(false)}
          order={order}
        />
      )}
    </div>
  )
}
