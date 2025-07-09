'use client'

import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { ArrowLeft, Package, Truck, CheckCircle, Clock, XCircle, MapPin, CreditCard, Download, Star, MessageCircle, RefreshCw, Copy, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { useOrder } from '@/hooks/use-orders'
import { formatPrice, formatDate } from '@/lib/utils'
import { toast } from 'sonner'

interface Props {
  orderId: string
}

const getStatusIcon = (status: string) => {
  switch (status.toLowerCase()) {
    case 'pending':
      return <Clock className="h-5 w-5" />
    case 'processing':
      return <Package className="h-5 w-5" />
    case 'shipped':
      return <Truck className="h-5 w-5" />
    case 'delivered':
      return <CheckCircle className="h-5 w-5" />
    case 'cancelled':
      return <XCircle className="h-5 w-5" />
    default:
      return <Package className="h-5 w-5" />
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

export function OrderDetailPage({ orderId }: Props) {
  const [copiedTrackingId, setCopiedTrackingId] = useState(false)
  const { data: order, isLoading, error } = useOrder(orderId)

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
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
        <div className="container mx-auto px-4 py-12">
          <div className="space-y-6">
            {[...Array(4)].map((_, i) => (
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
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black flex items-center justify-center">
        <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 max-w-md w-full mx-4">
          <CardContent className="p-8 text-center">
            <XCircle className="h-16 w-16 text-red-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Order Not Found</h3>
            <p className="text-gray-400 mb-6">The order you're looking for doesn't exist or you don't have permission to view it.</p>
            <Button asChild className="bg-gradient-to-r from-orange-500 to-orange-600">
              <Link href="/orders">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Orders
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (!order) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black flex items-center justify-center">
        <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 max-w-md w-full mx-4">
          <CardContent className="p-8 text-center">
            <XCircle className="h-16 w-16 text-yellow-400 mx-auto mb-4" />
            <h3 className="text-xl font-bold text-white mb-2">Order Not Found</h3>
            <p className="text-gray-400 mb-6">The order data could not be loaded. Please try refreshing the page.</p>
            <Button asChild className="bg-gradient-to-r from-orange-500 to-orange-600">
              <Link href="/orders">
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Orders
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      {/* Header */}
      <div className="bg-black/50 backdrop-blur-sm border-b border-gray-700">
        <div className="container mx-auto px-4 py-8">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div className="flex items-center gap-4">
              <Button variant="ghost" asChild className="text-gray-400 hover:text-white">
                <Link href="/orders">
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  Back to Orders
                </Link>
              </Button>
              <div>
                <h1 className="text-3xl font-bold text-white">
                  Order #{order?.order_number || 'Loading...'}
                </h1>
                <p className="text-gray-400">
                  {order?.created_at ? `Placed on ${formatDate(new Date(order.created_at))}` : 'Loading...'}
                </p>
              </div>
            </div>
            
            <div className="flex items-center gap-4">
              {order?.status && (
                <Badge className={`${getStatusColor(order.status)} text-lg px-4 py-2`}>
                  {getStatusIcon(order.status)}
                  <span className="ml-2">{order.status.charAt(0).toUpperCase() + order.status.slice(1)}</span>
                </Badge>
              )}
              
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800">
                <Download className="h-4 w-4 mr-2" />
                Download Invoice
              </Button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-8">
            {/* Order Status Timeline */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-3 text-white">
                  <div className="p-3 bg-blue-500/20 rounded-xl">
                    <Truck className="h-6 w-6 text-blue-400" />
                  </div>
                  Order Status & Tracking
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-6">
                  {/* Order Status */}
                  <div className="p-4 bg-gray-800/30 rounded-lg">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="text-gray-400 text-sm">Order Status</p>
                        <p className="text-white font-semibold text-lg">{order?.status?.charAt(0).toUpperCase() + (order?.status?.slice(1) || '')}</p>
                      </div>
                      <div>
                        <p className="text-gray-400 text-sm">Payment Status</p>
                        <p className={`font-semibold text-lg ${order?.payment_status === 'paid' ? 'text-green-400' : 'text-yellow-400'}`}>
                          {order?.payment_status?.charAt(0).toUpperCase() + (order?.payment_status?.slice(1) || '')}
                        </p>
                      </div>
                    </div>
                  </div>

                  {order?.tracking_number && (
                    <div className="p-4 bg-gray-800/30 rounded-lg">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-gray-400 text-sm">Tracking Number</p>
                          <p className="text-white font-mono text-lg">{order.tracking_number}</p>
                        </div>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={copyTrackingNumber}
                          className="border-gray-600 text-gray-300 hover:bg-gray-800"
                        >
                          {copiedTrackingId ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                        </Button>
                      </div>
                    </div>
                  )}

                  {order?.estimated_delivery && (
                    <div className="p-4 bg-green-500/10 border border-green-500/30 rounded-lg">
                      <p className="text-green-400 font-semibold">
                        Estimated Delivery: {formatDate(new Date(order.estimated_delivery))}
                      </p>
                    </div>
                  )}

                  {order?.notes && (
                    <div className="p-4 bg-blue-500/10 border border-blue-500/30 rounded-lg">
                      <p className="text-blue-400 font-semibold">Order Notes:</p>
                      <p className="text-gray-300 mt-1">{order.notes}</p>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Order Items */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-3 text-white">
                  <div className="p-3 bg-orange-500/20 rounded-xl">
                    <Package className="h-6 w-6 text-orange-400" />
                  </div>
                  Order Items ({order?.items?.length || order?.item_count || 0})
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {order?.items?.map((item, index) => (
                    <div key={item.id} className="flex items-center gap-4 p-4 bg-gray-800/30 rounded-lg hover:bg-gray-800/50 transition-colors">
                      {/* Product Image */}
                      <div className="relative w-20 h-20 bg-gray-700 rounded-lg overflow-hidden flex-shrink-0">
                        <Image
                          src={item.product?.images?.[0]?.url || '/placeholder-product.svg'}
                          alt={item.product?.name || item.product_name || 'Product'}
                          fill
                          className="object-cover"
                        />
                      </div>

                      {/* Product Details */}
                      <div className="flex-1 min-w-0">
                        <h4 className="font-semibold text-white text-lg mb-1">
                          {item.product?.name || item.product_name}
                        </h4>
                        <div className="space-y-1">
                          <p className="text-gray-400 text-sm">
                            <span className="font-medium">SKU:</span> {item.product_sku}
                          </p>
                          <p className="text-gray-400 text-sm">
                            <span className="font-medium">Quantity:</span> {item.quantity}
                          </p>
                          {item.product?.description && (
                            <p className="text-gray-400 text-sm line-clamp-2">
                              {item.product.description}
                            </p>
                          )}
                        </div>
                      </div>

                      {/* Pricing */}
                      <div className="text-right flex-shrink-0">
                        <div className="space-y-1">
                          <p className="text-gray-400 text-sm">Unit Price</p>
                          <p className="font-semibold text-white text-lg">{formatPrice(item.price)}</p>
                        </div>
                      </div>

                      <div className="text-right flex-shrink-0 ml-4">
                        <div className="space-y-1">
                          <p className="text-gray-400 text-sm">Total</p>
                          <p className="font-bold text-orange-400 text-xl">{formatPrice(item.total)}</p>
                        </div>
                      </div>

                      {/* Review Button for delivered orders */}
                      {order?.status === 'delivered' && (
                        <div className="flex-shrink-0 ml-4">
                          <Button 
                            variant="outline" 
                            size="sm" 
                            className="border-orange-500/50 text-orange-400 hover:bg-orange-500/10 hover:border-orange-500"
                          >
                            <Star className="h-4 w-4 mr-1" />
                            Review
                          </Button>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Order Summary */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="text-white flex items-center gap-2">
                  <div className="w-2 h-2 rounded-full bg-orange-400"></div>
                  Order Summary
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                {/* Subtotal */}
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-400 font-medium">Subtotal</span>
                  <span className="text-white font-semibold text-lg">{formatPrice(order?.subtotal || 0)}</span>
                </div>
                
                {/* Shipping */}
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-400 font-medium">Shipping</span>
                  <span className="text-white font-semibold text-lg">
                    {(order?.shipping_amount || 0) === 0 ? (
                      <span className="text-green-400 font-bold">FREE</span>
                    ) : (
                      formatPrice(order?.shipping_amount || 0)
                    )}
                  </span>
                </div>
                
                {/* Tax */}
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-400 font-medium">Tax</span>
                  <span className="text-white font-semibold text-lg">{formatPrice(order?.tax_amount || 0)}</span>
                </div>
                
                {/* Discount */}
                {order?.discount_amount && order.discount_amount > 0 ? (
                  <div className="flex justify-between items-center py-2">
                    <span className="text-gray-400 font-medium">Discount</span>
                    <span className="text-green-400 font-semibold text-lg">-{formatPrice(order.discount_amount)}</span>
                  </div>
                ) : null}
                
                <div className="h-px bg-gray-700 my-4"></div>
                
                {/* Total */}
                <div className="flex justify-between items-center py-3 bg-gray-800/50 rounded-lg px-4">
                  <span className="text-white font-bold text-xl">Total</span>
                  <span className="text-orange-400 font-bold text-2xl">{formatPrice(order?.total || 0)}</span>
                </div>
                
                <div className="h-px bg-gray-700 my-4"></div>
                
                {/* Payment Status */}
                <div className="flex justify-between items-center py-2">
                  <span className="text-gray-400 font-medium">Payment Status</span>
                  <div className="flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${order?.payment_status === 'paid' ? 'bg-green-400' : 'bg-yellow-400'}`}></div>
                    <span className={`font-semibold text-lg ${order?.payment_status === 'paid' ? 'text-green-400' : 'text-yellow-400'}`}>
                      {order?.payment_status?.charAt(0).toUpperCase() + (order?.payment_status?.slice(1) || '')}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Shipping Address */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-white">
                  <MapPin className="h-5 w-5 text-blue-400" />
                  Shipping Address
                </CardTitle>
              </CardHeader>
              <CardContent>
                {order?.shipping_address ? (
                  <div className="text-gray-300 space-y-1">
                    <p className="font-semibold text-white">
                      {order.shipping_address.first_name} {order.shipping_address.last_name}
                    </p>
                    {order.shipping_address.company && <p>{order.shipping_address.company}</p>}
                    <p>{order.shipping_address.address1}</p>
                    {order.shipping_address.address2 && <p>{order.shipping_address.address2}</p>}
                    <p>{order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.zip_code}</p>
                    <p>{order.shipping_address.country}</p>
                    {order.shipping_address.phone && <p>Phone: {order.shipping_address.phone}</p>}
                  </div>
                ) : (
                  <p className="text-gray-400">No shipping address available</p>
                )}
              </CardContent>
            </Card>

            {/* Payment Method */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-white">
                  <CreditCard className="h-5 w-5 text-green-400" />
                  Payment Method
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="text-gray-300 space-y-2">
                  <p>Payment Status: <span className={`font-semibold ${order?.payment_status === 'paid' ? 'text-green-400' : 'text-yellow-400'}`}>
                    {order?.payment_status?.charAt(0).toUpperCase() + (order?.payment_status?.slice(1) || '')}
                  </span></p>
                  {order?.payment && (
                    <p>Method: {order.payment.method || 'Stripe'}</p>
                  )}
                  <p>Currency: {order?.currency || 'USD'}</p>
                </div>
              </CardContent>
            </Card>

            {/* Actions */}
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="text-white">Need Help?</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button variant="outline" className="w-full border-gray-600 text-gray-300 hover:bg-gray-800">
                  <MessageCircle className="h-4 w-4 mr-2" />
                  Contact Support
                </Button>
                
                {order?.can_be_cancelled && order?.status !== 'delivered' && order?.status !== 'cancelled' && (
                  <Button variant="outline" className="w-full border-red-500/50 text-red-400 hover:bg-red-500/10">
                    <XCircle className="h-4 w-4 mr-2" />
                    Cancel Order
                  </Button>
                )}
                
                <Button variant="outline" className="w-full border-gray-600 text-gray-300 hover:bg-gray-800">
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Reorder Items
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
