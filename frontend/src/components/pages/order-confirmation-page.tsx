'use client'

import { useEffect } from 'react'
import Link from 'next/link'
import { CheckCircle, Package, Truck, Mail, ArrowRight } from 'lucide-react'
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
    <div className="min-h-screen bg-gradient-to-br from-emerald-50 via-background to-green-50 py-16">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          {/* Success Header */}
          <div className="text-center mb-16">
            {/* Success Animation */}
            <div className="relative mb-8">
              <div className="w-32 h-32 bg-gradient-to-br from-emerald-500 to-green-600 rounded-full flex items-center justify-center mx-auto shadow-2xl animate-bounce-slow">
                <CheckCircle className="h-20 w-20 text-white" />
              </div>
              <div className="absolute inset-0 w-32 h-32 bg-emerald-400/30 rounded-full mx-auto animate-ping"></div>
            </div>

            <h1 className="text-5xl lg:text-6xl font-bold text-foreground mb-6">
              Order <span className="text-gradient bg-gradient-to-r from-emerald-600 to-green-600 bg-clip-text text-transparent">Confirmed!</span>
            </h1>
            <p className="text-2xl text-muted-foreground mb-8 leading-relaxed">
              Thank you for your purchase! Your order has been received and is being processed with care.
            </p>

            {/* Order Info */}
            <div className="flex flex-col sm:flex-row items-center justify-center gap-6 text-lg">
              <div className="flex items-center gap-2">
                <Package className="h-5 w-5 text-primary" />
                <span className="font-semibold">Order #{order.number}</span>
              </div>
              <div className="hidden sm:block w-2 h-2 bg-muted-foreground rounded-full"></div>
              <div className="flex items-center gap-2">
                <span className="text-muted-foreground">{order.date}</span>
              </div>
            </div>
          </div>

          {/* Order Status */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <Package className="h-5 w-5" />
                <span>Order Status</span>
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                    <CheckCircle className="h-5 w-5 text-green-600" />
                  </div>
                  <div>
                    <p className="font-medium">Order Confirmed</p>
                    <p className="text-sm text-gray-600">We've received your order</p>
                  </div>
                </div>
                <Badge variant="success">Confirmed</Badge>
              </div>
              
              <div className="mt-6 flex items-center space-x-4 text-sm text-gray-600">
                <div className="flex items-center space-x-2">
                  <Truck className="h-4 w-4" />
                  <span>Estimated delivery: {order.estimated_delivery}</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Mail className="h-4 w-4" />
                  <span>Confirmation email sent</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Order Summary */}
          <Card className="mb-8">
            <CardHeader>
              <CardTitle>Order Summary</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {order.items.map((item) => (
                  <div key={item.id} className="flex items-center space-x-4">
                    <div className="w-16 h-16 bg-gray-100 rounded-md flex items-center justify-center">
                      <Package className="h-8 w-8 text-gray-400" />
                    </div>
                    <div className="flex-1">
                      <h4 className="font-medium">{item.name}</h4>
                      <p className="text-sm text-gray-600">Qty: {item.quantity}</p>
                    </div>
                    <p className="font-medium">{formatPrice(item.price)}</p>
                  </div>
                ))}
              </div>
              
              <div className="border-t mt-6 pt-6 space-y-2">
                <div className="flex justify-between">
                  <span>Subtotal</span>
                  <span>{formatPrice(order.subtotal)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Shipping</span>
                  <span>{formatPrice(order.shipping)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Tax</span>
                  <span>{formatPrice(order.tax)}</span>
                </div>
                <div className="flex justify-between text-lg font-semibold border-t pt-2">
                  <span>Total</span>
                  <span>{formatPrice(order.total)}</span>
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

          {/* Action Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            {isAuthenticated && (
              <Button variant="outline" asChild>
                <Link href="/account">
                  View Order History
                </Link>
              </Button>
            )}
            
            <Button asChild>
              <Link href="/products">
                Continue Shopping
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>

          {/* Support */}
          <div className="text-center mt-12 p-6 bg-gray-100 rounded-lg">
            <h3 className="font-semibold mb-2">Need Help?</h3>
            <p className="text-sm text-gray-600 mb-4">
              If you have any questions about your order, our customer support team is here to help.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button variant="outline" size="sm" asChild>
                <Link href="/contact">Contact Support</Link>
              </Button>
              <Button variant="outline" size="sm" asChild>
                <Link href="/help">View FAQ</Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
