'use client'

import { useEffect } from 'react'
import Link from 'next/link'
import { CheckCircle, Package, Truck, Mail, ArrowRight, ShoppingBag, Star, Clock, CreditCard, Shield } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { formatPrice } from '@/lib/utils'

export function OrderConfirmationPage() {
  const { user, isAuthenticated } = useAuthStore()

  // Mock order data - would come from URL params or API
  const order = {
    id: 'ORD-2024-001',
    number: 'ORD-2024-001',
    date: new Date().toLocaleDateString(),
    status: 'confirmed',
    total: 129.99,
    shipping: 9.99,
    tax: 10.40,
    subtotal: 109.60,
    items: [
      {
        id: '1',
        name: 'Wireless Bluetooth Headphones',
        price: 79.99,
        quantity: 1,
        image: '/placeholder-product.jpg',
      },
      {
        id: '2',
        name: 'USB-C Cable',
        price: 19.99,
        quantity: 1,
        image: '/placeholder-product.jpg',
      },
      {
        id: '3',
        name: 'Phone Case',
        price: 9.62,
        quantity: 1,
        image: '/placeholder-product.jpg',
      },
    ],
    shipping_address: {
      name: user ? `${user.first_name} ${user.last_name}` : 'John Doe',
      address: '123 Main Street',
      city: 'New York',
      state: 'NY',
      zip: '10001',
      country: 'United States',
    },
    estimated_delivery: '3-5 business days',
  }

  useEffect(() => {
    // Track order confirmation event
    if (typeof window !== 'undefined') {
      // Analytics tracking would go here
      console.log('Order confirmed:', order.id)
    }
  }, [order.id])

  return (
    <div className="min-h-screen hero-gradient py-16">
      <div className="container mx-auto px-4">
        <div className="max-w-5xl mx-auto">
          {/* BiHub Success Header */}
          <div className="text-center mb-16 animate-fade-in">
            {/* Success Animation */}
            <div className="relative mb-12">
              <div className="w-40 h-40 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-full flex items-center justify-center mx-auto shadow-2xl animate-scale-in">
                <CheckCircle className="h-24 w-24 text-white" />
              </div>
              <div className="absolute -top-4 -right-4 w-16 h-16 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-xl animate-bounce-slow">
                <span className="text-white text-2xl font-bold">✓</span>
              </div>
              <div className="absolute inset-0 w-40 h-40 bg-emerald-400/20 rounded-full mx-auto animate-ping"></div>
            </div>

            <h1 className="text-5xl lg:text-6xl font-bold text-white mb-6 animate-slide-up">
              Order <span className="text-gradient">Confirmed!</span>
            </h1>
            <p className="text-xl text-gray-300 mb-8 leading-relaxed animate-slide-up max-w-3xl mx-auto">
              Thank you for choosing <span className="text-orange-400 font-bold">BiHub</span>! Your order has been received and is being processed with care.
            </p>

            {/* BiHub Order Info */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-6 text-lg animate-scale-in">
              <Badge className="bg-emerald-500/20 text-emerald-400 border-emerald-500/30 px-6 py-3 text-lg">
                <Package className="h-5 w-5 mr-2" />
                Order #{order.number}
              </Badge>
              <Badge className="bg-orange-500/20 text-orange-400 border-orange-500/30 px-6 py-3 text-lg">
                <Clock className="h-5 w-5 mr-2" />
                {order.date}
              </Badge>
              <Badge className="bg-blue-500/20 text-blue-400 border-blue-500/30 px-6 py-3 text-lg">
                <CreditCard className="h-5 w-5 mr-2" />
                {formatPrice(order.total)}
              </Badge>
            </div>
          </div>

          {/* BiHub Order Status */}
          <Card className="mb-8 bg-gray-900/50 border-gray-700 hover:border-orange-500/50 transition-all duration-300 animate-slide-up">
            <CardHeader>
              <CardTitle className="flex items-center space-x-3 text-white">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center">
                  <Package className="h-5 w-5 text-white" />
                </div>
                <span>Order Status</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center space-x-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-full flex items-center justify-center">
                    <CheckCircle className="h-6 w-6 text-white" />
                  </div>
                  <div>
                    <p className="font-bold text-white text-lg">Order Confirmed</p>
                    <p className="text-gray-400">We've received your BiHub order and it's being processed</p>
                  </div>
                </div>
                <Badge className="bg-emerald-500/20 text-emerald-400 border-emerald-500/30 px-4 py-2">
                  Confirmed
                </Badge>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-center space-x-3 p-3 bg-gray-800 rounded-lg">
                  <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                    <Truck className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-white">Estimated delivery</p>
                    <p className="text-xs text-gray-400">{order.estimated_delivery}</p>
                  </div>
                </div>
                <div className="flex items-center space-x-3 p-3 bg-gray-800 rounded-lg">
                  <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center">
                    <Mail className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-white">Confirmation email</p>
                    <p className="text-xs text-gray-400">Sent to your email</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* BiHub Order Summary */}
          <Card className="mb-8 bg-gray-900/50 border-gray-700 hover:border-orange-500/50 transition-all duration-300 animate-slide-up">
            <CardHeader>
              <CardTitle className="text-white flex items-center gap-3">
                <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center">
                  <ShoppingBag className="h-5 w-5 text-white" />
                </div>
                Order Summary
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order.items.map((item) => (
                  <div key={item.id} className="flex items-center space-x-4 p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
                    <div className="w-16 h-16 bg-gradient-to-br from-gray-700 to-gray-800 rounded-lg flex items-center justify-center border border-gray-600">
                      <Package className="h-8 w-8 text-orange-400" />
                    </div>
                    <div className="flex-1">
                      <h4 className="font-bold text-white">{item.name}</h4>
                      <div className="flex items-center gap-2 mt-1">
                        <Badge className="bg-orange-500/20 text-orange-400 border-orange-500/30 text-xs">
                          Qty: {item.quantity}
                        </Badge>
                        <div className="flex items-center gap-1">
                          {[...Array(5)].map((_, i) => (
                            <Star key={i} className="h-3 w-3 fill-orange-400 text-orange-400" />
                          ))}
                        </div>
                      </div>
                    </div>
                    <p className="font-bold text-orange-400 text-lg">{formatPrice(item.price)}</p>
                  </div>
                ))}
              </div>

              <div className="border-t border-gray-700 mt-6 pt-6 space-y-3">
                <div className="flex justify-between text-gray-300">
                  <span>Subtotal</span>
                  <span className="font-semibold">{formatPrice(order.subtotal)}</span>
                </div>
                <div className="flex justify-between text-gray-300">
                  <span>Shipping</span>
                  <span className="font-semibold">{formatPrice(order.shipping)}</span>
                </div>
                <div className="flex justify-between text-gray-300">
                  <span>Tax</span>
                  <span className="font-semibold">{formatPrice(order.tax)}</span>
                </div>
                <div className="flex justify-between text-xl font-bold border-t border-gray-700 pt-3">
                  <span className="text-white">Total</span>
                  <span className="text-orange-400">{formatPrice(order.total)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Shipping Information */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle>Shipping Information</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-2">
                <p className="font-medium">{order.shipping_address.name}</p>
                <p className="text-gray-600">{order.shipping_address.address}</p>
                <p className="text-gray-600">
                  {order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.zip}
                </p>
                <p className="text-gray-600">{order.shipping_address.country}</p>
              </div>
            </CardContent>
          </Card>

          {/* Next Steps */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle>What's Next?</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-primary-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <span className="text-primary-600 text-sm font-medium">1</span>
                  </div>
                  <div>
                    <p className="font-medium">Order Processing</p>
                    <p className="text-sm text-gray-600">
                      We're preparing your items for shipment. This usually takes 1-2 business days.
                    </p>
                  </div>
                </div>
                
                <div className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <span className="text-gray-600 text-sm font-medium">2</span>
                  </div>
                  <div>
                    <p className="font-medium">Shipping</p>
                    <p className="text-sm text-gray-600">
                      Once shipped, you'll receive a tracking number to monitor your package.
                    </p>
                  </div>
                </div>
                
                <div className="flex items-start space-x-3">
                  <div className="w-6 h-6 bg-gray-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                    <span className="text-gray-600 text-sm font-medium">3</span>
                  </div>
                  <div>
                    <p className="font-medium">Delivery</p>
                    <p className="text-sm text-gray-600">
                      Your order will arrive within {order.estimated_delivery}.
                    </p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* BiHub Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-6 justify-center animate-scale-in">
            {isAuthenticated && (
              <Button variant="outline" size="lg" className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white px-8 py-4" asChild>
                <Link href="/profile">
                  <Package className="mr-3 h-5 w-5" />
                  View Order History
                </Link>
              </Button>
            )}

            <Button size="lg" className="btn-gradient px-8 py-4" asChild>
              <Link href="/products">
                <ShoppingBag className="mr-3 h-5 w-5" />
                Continue Shopping
                <ArrowRight className="ml-3 h-5 w-5" />
              </Link>
            </Button>
          </div>

          {/* BiHub Support */}
          <div className="text-center mt-16 p-8 bg-gray-900/50 border border-gray-700 rounded-2xl animate-fade-in">
            <div className="w-16 h-16 rounded-full bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center mx-auto mb-4">
              <Shield className="h-8 w-8 text-white" />
            </div>
            <h3 className="font-bold text-white mb-3 text-xl">Need Help with Your BiHub Order?</h3>
            <p className="text-gray-300 mb-6 max-w-2xl mx-auto">
              Our BiHub customer support team is available 24/7 to assist you with any questions about your order, shipping, or returns.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white" asChild>
                <Link href="/contact">
                  <Mail className="mr-2 h-4 w-4" />
                  Contact Support
                </Link>
              </Button>
              <Button variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white" asChild>
                <Link href="/help">
                  <Package className="mr-2 h-4 w-4" />
                  View FAQ
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
