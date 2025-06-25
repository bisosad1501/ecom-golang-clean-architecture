'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CheckCircle, Package, ArrowRight, Home } from 'lucide-react'
import { useOrderStore } from '@/store/order'
import { useCartStore } from '@/store/cart'

export default function CheckoutSuccessPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const sessionId = searchParams.get('session_id')
  
  const [isLoading, setIsLoading] = useState(true)
  const [orderDetails, setOrderDetails] = useState<any>(null)
  
  const { clearCart } = useCartStore()

  useEffect(() => {
    const handleSuccess = async () => {
      if (!sessionId) {
        router.push('/checkout')
        return
      }

      try {
        // Clear the cart since payment was successful
        await clearCart()
        
        // You could also fetch order details here if needed
        // const orderResponse = await getOrderBySessionId(sessionId)
        // setOrderDetails(orderResponse.data)
        
        setIsLoading(false)
      } catch (error) {
        console.error('Error handling checkout success:', error)
        setIsLoading(false)
      }
    }

    handleSuccess()
  }, [sessionId, clearCart, router])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-green-500"></div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        <Card className="text-center">
          <CardHeader className="pb-4">
            <div className="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
              <CheckCircle className="w-8 h-8 text-green-600" />
            </div>
            <CardTitle className="text-2xl text-green-600">
              Payment Successful!
            </CardTitle>
          </CardHeader>
          
          <CardContent className="space-y-6">
            <div className="text-gray-600">
              <p className="mb-2">
                Thank you for your order! Your payment has been processed successfully.
              </p>
              {sessionId && (
                <p className="text-sm text-gray-500">
                  Session ID: {sessionId}
                </p>
              )}
            </div>

            <div className="bg-blue-50 p-4 rounded-lg">
              <div className="flex items-center justify-center mb-2">
                <Package className="w-5 h-5 text-blue-600 mr-2" />
                <span className="font-medium text-blue-800">What's Next?</span>
              </div>
              <p className="text-sm text-blue-700">
                You'll receive an order confirmation email shortly. 
                We'll notify you when your order ships.
              </p>
            </div>

            <div className="space-y-3">
              <Button 
                asChild 
                className="w-full"
                size="lg"
              >
                <Link href="/orders">
                  <Package className="mr-2 h-4 w-4" />
                  View My Orders
                </Link>
              </Button>
              
              <Button 
                asChild 
                variant="outline" 
                className="w-full"
                size="lg"
              >
                <Link href="/products">
                  <ArrowRight className="mr-2 h-4 w-4" />
                  Continue Shopping
                </Link>
              </Button>
              
              <Button 
                asChild 
                variant="ghost" 
                className="w-full"
                size="lg"
              >
                <Link href="/">
                  <Home className="mr-2 h-4 w-4" />
                  Back to Home
                </Link>
              </Button>
            </div>

            <div className="text-xs text-gray-500 pt-4 border-t">
              <p>
                Need help? Contact our{' '}
                <Link href="/support" className="text-blue-600 hover:underline">
                  customer support
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
