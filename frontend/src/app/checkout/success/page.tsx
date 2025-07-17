'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CheckCircle, Package, Home, Truck, CreditCard, Star, Gift, Shield, ShoppingBag } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { apiClient } from '@/lib/api'
import { useQueryClient } from '@tanstack/react-query'
import { orderKeys } from '@/hooks/use-orders'

export default function CheckoutSuccessPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionId = searchParams.get('session_id')
  const orderId = searchParams.get('order_id')
  const checkoutSessionId = searchParams.get('checkout_session_id')

  const [isLoading, setIsLoading] = useState(true)
  const [orderDetails, setOrderDetails] = useState<any>(null)

  const { clearCart } = useCartStore()
  const { isAuthenticated, token, user } = useAuthStore()
  const queryClient = useQueryClient()

  useEffect(() => {
    const handleSuccess = async () => {
      if (!sessionId && !checkoutSessionId) {
        router.push('/checkout')
        return
      }

      // Check authentication status
      const localToken = localStorage.getItem('token')
      const authToken = token || localToken

      console.log('Success page - Auth check:', {
        isAuthenticated,
        storeToken: !!token,
        localToken: !!localToken,
        finalToken: !!authToken,
        sessionId,
        orderId
      })

      // Debug: Make auth info available globally for debugging
      if (typeof window !== 'undefined') {
        (window as any).debugAuth = {
          isAuthenticated,
          storeToken: token,
          localToken,
          finalToken: authToken,
          user
        }
        console.log('Debug auth info available at window.debugAuth')
      }

      try {
        // Clear the cart since payment was successful
        await clearCart()

        // Confirm payment success with backend (fallback if webhook fails)
        if (sessionId) {
          try {
            console.log('🔄 Attempting payment confirmation for session:', sessionId, 'order:', orderId || 'unknown')

            // Try without auth first (public endpoint)
            const confirmResponse = await fetch(`http://localhost:8080/api/v1/payments/confirm-success`, {
              method: 'POST',
              headers: {
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                session_id: sessionId,
                ...(orderId && { order_id: orderId })
              })
            })

            if (confirmResponse.ok) {
              console.log('✅ Payment confirmation sent to backend successfully')
              const confirmData = await confirmResponse.json()
              console.log('Confirmation response:', confirmData)

              // Invalidate orders cache to refresh data
              queryClient.invalidateQueries({ queryKey: orderKeys.all })
              console.log('🔄 Orders cache invalidated - data will refresh')
            } else {
              console.warn('❌ Failed to confirm payment with backend:', confirmResponse.status, confirmResponse.statusText)
              const errorText = await confirmResponse.text()
              console.warn('Error details:', errorText)

              // Try with auth token if available
              const authToken = token || localStorage.getItem('token')
              if (authToken) {
                console.log('🔄 Retrying with authentication...')
                const authConfirmResponse = await fetch(`http://localhost:8080/api/v1/payments/confirm-success`, {
                  method: 'POST',
                  headers: {
                    'Authorization': `Bearer ${authToken}`,
                    'Content-Type': 'application/json',
                  },
                  body: JSON.stringify({
                    session_id: sessionId,
                    ...(orderId && { order_id: orderId })
                  })
                })

                if (authConfirmResponse.ok) {
                  console.log('✅ Payment confirmation with auth successful')

                  // Invalidate orders cache to refresh data
                  queryClient.invalidateQueries({ queryKey: orderKeys.all })
                  console.log('🔄 Orders cache invalidated - data will refresh')
                } else {
                  console.warn('❌ Auth confirmation also failed:', authConfirmResponse.status)
                }
              }
            }
          } catch (confirmError) {
            console.warn('Error confirming payment with backend:', confirmError)
          }
        }

        // Fetch order details if order ID is available
        if (orderId) {
          try {
            console.log('Fetching order details for order:', orderId)

            // Try to fetch order details using public endpoint first (for success page)
            const response = await fetch(`http://localhost:8080/api/v1/orders/${orderId}/public`, {
              headers: {
                'Content-Type': 'application/json',
              },
            })

            if (response.ok) {
              const data = await response.json()
              setOrderDetails(data.data)
              console.log('Order details fetched:', data.data)
            } else if (response.status === 401) {
              // Try with authentication if available
              const authToken = token || localStorage.getItem('token')
              if (authToken) {
                console.log('Retrying with authentication...')
                const authResponse = await fetch(`http://localhost:8080/api/v1/orders/${orderId}`, {
                  headers: {
                    'Authorization': `Bearer ${authToken}`,
                    'Content-Type': 'application/json',
                  },
                })

                if (authResponse.ok) {
                  const authData = await authResponse.json()
                  setOrderDetails(authData.data)
                  console.log('Order details fetched with auth:', authData.data)
                } else {
                  console.warn('Auth request also failed:', authResponse.status)
                  setOrderDetails({
                    id: orderId,
                    order_number: `ORD-${orderId.slice(-8)}`,
                    status: 'pending',
                    payment_status: 'processing'
                  })
                }
              } else {
                console.warn('No auth token available for retry')
                setOrderDetails({
                  id: orderId,
                  order_number: `ORD-${orderId.slice(-8)}`,
                  status: 'pending',
                  payment_status: 'processing'
                })
              }
            } else {
              console.error('Failed to fetch order details:', response.status, response.statusText)
              setOrderDetails({
                id: orderId,
                order_number: `ORD-${orderId.slice(-8)}`,
                status: 'pending',
                payment_status: 'processing'
              })
            }
          } catch (fetchError) {
            console.error('Error fetching order details:', fetchError)
            setOrderDetails({
              id: orderId,
              order_number: `ORD-${orderId.slice(-8)}`,
              status: 'pending',
              payment_status: 'processing'
            })
          }
        }

        setIsLoading(false)
      } catch (error) {
        console.error('Error handling checkout success:', error)
        setIsLoading(false)
      }
    }

    handleSuccess()
  }, [sessionId, orderId, clearCart, router])

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black flex items-center justify-center">
        <div className="text-center">
          <div className="relative mb-8">
            <div className="animate-spin rounded-full h-20 w-20 border-4 border-gray-700 border-t-green-500 mx-auto"></div>
            <div className="absolute inset-0 rounded-full h-20 w-20 border-4 border-transparent border-r-green-400 animate-pulse mx-auto"></div>
          </div>
          <div className="space-y-3">
            <h3 className="text-xl font-semibold text-white">Processing your order...</h3>
            <p className="text-gray-400">Please wait while we confirm your payment</p>
            <div className="flex justify-center space-x-1 mt-4">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-bounce"></div>
              <div className="w-2 h-2 bg-green-500 rounded-full animate-bounce" style={{animationDelay: '0.1s'}}></div>
              <div className="w-2 h-2 bg-green-500 rounded-full animate-bounce" style={{animationDelay: '0.2s'}}></div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      {/* Success Header */}
      <div className="bg-gradient-to-r from-green-500/20 to-emerald-500/20 border-b border-green-500/30">
        <div className="container mx-auto px-4 py-12">
          <div className="text-center">
            <div className="relative mb-8">
              <div className="w-32 h-32 bg-gradient-to-br from-green-500 to-green-600 rounded-full flex items-center justify-center mx-auto shadow-2xl animate-scale-in">
                <CheckCircle className="h-16 w-16 text-white" />
              </div>
              <div className="absolute -top-4 -right-4 w-12 h-12 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center animate-pulse">
                <Star className="h-6 w-6 text-white" />
              </div>

              {/* Floating particles */}
              <div className="absolute -top-8 -left-8 w-4 h-4 bg-green-400 rounded-full animate-bounce opacity-60"></div>
              <div className="absolute -bottom-4 -right-8 w-3 h-3 bg-emerald-400 rounded-full animate-bounce opacity-60" style={{animationDelay: '0.5s'}}></div>
              <div className="absolute -top-4 left-8 w-2 h-2 bg-green-300 rounded-full animate-bounce opacity-60" style={{animationDelay: '1s'}}></div>
            </div>

            <h1 className="text-5xl font-bold text-white mb-6 animate-fade-in">
              Payment Successful! 🎉
            </h1>
            <p className="text-2xl text-gray-300 mb-4">
              Thank you for shopping with Bi<span className="text-orange-400">Hub</span>
            </p>
            <div className="flex items-center justify-center gap-4 text-gray-400">
              {orderDetails ? (
                <>
                  <span className="text-sm">Order: #{orderDetails.order_number}</span>
                  <Separator orientation="vertical" className="h-4 bg-gray-600" />
                  <span className="text-sm">Status: {orderDetails.status}</span>
                  <Separator orientation="vertical" className="h-4 bg-gray-600" />
                </>
              ) : sessionId && (
                <>
                  <span className="text-sm">Session: #{sessionId.slice(-8).toUpperCase()}</span>
                  <Separator orientation="vertical" className="h-4 bg-gray-600" />
                </>
              )}
              <span className="text-sm">Confirmation sent to your email</span>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        <div className="max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* What's Next Section */}
          <div className="lg:col-span-2 space-y-8">
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-3 text-white text-2xl">
                  <div className="p-3 bg-blue-500/20 rounded-xl">
                    <Package className="h-8 w-8 text-blue-400" />
                  </div>
                  What Happens Next?
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center p-6 bg-gray-800/30 rounded-xl hover:bg-gray-800/50 transition-all duration-300">
                    <div className="w-16 h-16 bg-gradient-to-br from-green-500 to-green-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                      <CreditCard className="h-8 w-8 text-white" />
                    </div>
                    <h4 className="font-bold text-white mb-2">Payment Confirmed</h4>
                    <p className="text-gray-400 text-sm">Your payment has been processed securely</p>
                    <Badge className="mt-3 bg-green-500/20 text-green-400 border-green-500/30">
                      ✓ Complete
                    </Badge>
                  </div>

                  <div className="text-center p-6 bg-gray-800/30 rounded-xl hover:bg-gray-800/50 transition-all duration-300">
                    <div className="w-16 h-16 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                      <Package className="h-8 w-8 text-white" />
                    </div>
                    <h4 className="font-bold text-white mb-2">Order Processing</h4>
                    <p className="text-gray-400 text-sm">We're preparing your items for shipment</p>
                    <Badge className="mt-3 bg-orange-500/20 text-orange-400 border-orange-500/30">
                      ⏳ In Progress
                    </Badge>
                  </div>

                  <div className="text-center p-6 bg-gray-800/30 rounded-xl hover:bg-gray-800/50 transition-all duration-300">
                    <div className="w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-lg">
                      <Truck className="h-8 w-8 text-white" />
                    </div>
                    <h4 className="font-bold text-white mb-2">Fast Delivery</h4>
                    <p className="text-gray-400 text-sm">Estimated delivery: 2-3 business days</p>
                    <Badge className="mt-3 bg-blue-500/20 text-blue-400 border-blue-500/30">
                      🚚 Coming Soon
                    </Badge>
                  </div>
                </div>
              </CardContent>
            </Card>

          </div>

          {/* Actions Sidebar */}
          <div className="space-y-6">
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="text-white">Quick Actions</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="relative group">
                  <div className="absolute -inset-1 bg-gradient-to-r from-orange-500 to-orange-600 rounded-xl blur opacity-25 group-hover:opacity-75 transition duration-1000 group-hover:duration-200"></div>
                  <Button
                    className="relative w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold py-3 rounded-lg transition-all duration-300 transform hover:scale-[1.02]"
                    asChild
                  >
                    <Link href="/orders">
                      <Package className="mr-3 h-5 w-5" />
                      Track Your Order
                    </Link>
                  </Button>
                </div>

                <Button
                  variant="outline"
                  className="w-full border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                  asChild
                >
                  <Link href="/products">
                    <ShoppingBag className="mr-3 h-5 w-5" />
                    Continue Shopping
                  </Link>
                </Button>

                <Button
                  variant="outline"
                  className="w-full border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                  asChild
                >
                  <Link href="/">
                    <Home className="mr-3 h-5 w-5" />
                    Back to Home
                  </Link>
                </Button>
              </CardContent>
            </Card>

            {/* Special Offers */}
            <Card className="bg-gradient-to-br from-purple-900/30 to-pink-900/30 border-purple-500/30">
              <CardHeader>
                <CardTitle className="text-white flex items-center gap-2">
                  <Gift className="h-5 w-5 text-purple-400" />
                  Special Offer
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <p className="text-gray-300 text-sm">
                    Get 20% off your next order with code <span className="font-bold text-purple-400">THANKS20</span>
                  </p>
                  <Badge className="bg-purple-500/20 text-purple-400 border-purple-500/30">
                    Valid for 30 days
                  </Badge>
                </div>
              </CardContent>
            </Card>

            {/* Support */}
            <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="text-white text-lg flex items-center gap-2">
                  <Shield className="h-5 w-5 text-blue-400" />
                  Need Help?
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-400 text-sm mb-4">
                  Our support team is here 24/7 to help with any questions about your order.
                </p>
                <Button variant="outline" size="sm" className="w-full border-gray-600 text-gray-300 hover:bg-gray-800">
                  Contact Support
                </Button>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Footer Message */}
        <div className="text-center mt-12 pt-8 border-t border-gray-700">
          <p className="text-gray-400 mb-4">
            Thank you for choosing Bi<span className="text-orange-400">Hub</span>! We appreciate your business and look forward to serving you again.
          </p>
          <p className="text-gray-500 text-xs mt-2">
            Need help? Contact our{' '}
            <Link href="/support" className="text-orange-400 hover:text-orange-300 transition-colors">
              customer support team
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}
