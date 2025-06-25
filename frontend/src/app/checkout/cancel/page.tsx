'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { XCircle, ArrowLeft, ShoppingCart, Home } from 'lucide-react'

export default function CheckoutCancelPage() {
  return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        <Card className="text-center">
          <CardHeader className="pb-4">
            <div className="mx-auto w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
              <XCircle className="w-8 h-8 text-red-600" />
            </div>
            <CardTitle className="text-2xl text-red-600">
              Payment Cancelled
            </CardTitle>
          </CardHeader>
          
          <CardContent className="space-y-6">
            <div className="text-gray-600">
              <p className="mb-2">
                Your payment was cancelled. No charges have been made to your account.
              </p>
              <p className="text-sm text-gray-500">
                Your cart items are still saved and ready for checkout.
              </p>
            </div>

            <div className="bg-yellow-50 p-4 rounded-lg">
              <div className="flex items-center justify-center mb-2">
                <ShoppingCart className="w-5 h-5 text-yellow-600 mr-2" />
                <span className="font-medium text-yellow-800">Your Cart is Safe</span>
              </div>
              <p className="text-sm text-yellow-700">
                All items in your cart have been preserved. 
                You can continue shopping or try checking out again.
              </p>
            </div>

            <div className="space-y-3">
              <Button 
                asChild 
                className="w-full"
                size="lg"
              >
                <Link href="/checkout">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Return to Checkout
                </Link>
              </Button>
              
              <Button 
                asChild 
                variant="outline" 
                className="w-full"
                size="lg"
              >
                <Link href="/cart">
                  <ShoppingCart className="mr-2 h-4 w-4" />
                  View Cart
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
                Having trouble? Contact our{' '}
                <Link href="/support" className="text-blue-600 hover:underline">
                  customer support
                </Link>{' '}
                for assistance.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
