'use client'

import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { XCircle, ArrowLeft, ShoppingCart, Home, RefreshCw, CreditCard, Shield, Headphones, ShoppingBag } from 'lucide-react'

export default function CheckoutCancelPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      {/* Cancel Header */}
      <div className="bg-gradient-to-r from-red-500/20 to-orange-500/20 border-b border-red-500/30">
        <div className="container mx-auto px-4 py-12">
          <div className="text-center">
            <div className="relative mb-8">
              <div className="w-32 h-32 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center mx-auto shadow-2xl animate-bounce-in">
                <XCircle className="h-16 w-16 text-white" />
              </div>
              <div className="absolute -top-4 -right-4 w-12 h-12 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center animate-pulse">
                <span className="text-white text-2xl">⚠️</span>
              </div>

              {/* Floating particles */}
              <div className="absolute -top-8 -left-8 w-4 h-4 bg-red-400 rounded-full animate-bounce opacity-60"></div>
              <div className="absolute -bottom-4 -right-8 w-3 h-3 bg-orange-400 rounded-full animate-bounce opacity-60" style={{animationDelay: '0.5s'}}></div>
              <div className="absolute -top-4 left-8 w-2 h-2 bg-red-300 rounded-full animate-bounce opacity-60" style={{animationDelay: '1s'}}></div>
            </div>

            <h1 className="text-5xl font-bold text-white mb-6 animate-fade-in">
              Payment Cancelled
            </h1>
            <p className="text-2xl text-gray-300 mb-4">
              No worries! Your order from Bi<span className="text-orange-400">Hub</span> is still waiting
            </p>
            <p className="text-gray-400">
              Your cart items are saved and ready when you're ready to complete your purchase
            </p>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        <div className="max-w-6xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Why This Happened */}
          <div className="lg:col-span-2 space-y-8">
            <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700">
              <CardHeader>
                <CardTitle className="flex items-center gap-3 text-white text-2xl">
                  <div className="p-3 bg-blue-500/20 rounded-xl">
                    <RefreshCw className="h-8 w-8 text-blue-400" />
                  </div>
                  What Happened?
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-6">
                  <p className="text-gray-300 text-lg">
                    Your payment was cancelled, but don't worry - this happens for various reasons and your cart is still saved.
                  </p>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="p-6 bg-gray-800/30 rounded-xl">
                      <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center mb-4">
                        <CreditCard className="h-6 w-6 text-blue-400" />
                      </div>
                      <h4 className="font-bold text-white mb-2">Payment Issues</h4>
                      <p className="text-gray-400 text-sm">
                        Sometimes payment providers have temporary issues or cards need verification
                      </p>
                    </div>

                    <div className="p-6 bg-gray-800/30 rounded-xl">
                      <div className="w-12 h-12 bg-orange-500/20 rounded-lg flex items-center justify-center mb-4">
                        <Shield className="h-6 w-6 text-orange-400" />
                      </div>
                      <h4 className="font-bold text-white mb-2">Security Check</h4>
                      <p className="text-gray-400 text-sm">
                        Your bank might have flagged the transaction for additional security verification
                      </p>
                    </div>
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
                    <Link href="/cart">
                      <ShoppingBag className="mr-3 h-5 w-5" />
                      Return to Cart
                    </Link>
                  </Button>
                </div>

                <Button
                  variant="outline"
                  className="w-full border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                  asChild
                >
                  <Link href="/checkout">
                    <CreditCard className="mr-3 h-5 w-5" />
                    Try Payment Again
                  </Link>
                </Button>

                <Button
                  variant="outline"
                  className="w-full border-gray-600 text-gray-300 hover:bg-gray-800 hover:border-gray-500 transition-all duration-300"
                  asChild
                >
                  <Link href="/products">
                    <ArrowLeft className="mr-3 h-5 w-5" />
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

            {/* Support */}
            <Card className="bg-gradient-to-br from-blue-900/30 to-purple-900/30 border-blue-500/30">
              <CardHeader>
                <CardTitle className="text-white flex items-center gap-2">
                  <Headphones className="h-5 w-5 text-blue-400" />
                  Need Help?
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <p className="text-gray-300 text-sm">
                    Our support team is available 24/7 to help resolve any payment issues.
                  </p>
                  <Button variant="outline" size="sm" className="w-full border-blue-500/50 text-blue-400 hover:bg-blue-500/10">
                    Contact Support
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Security Notice */}
            <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700">
              <CardContent className="pt-6">
                <div className="flex items-center gap-3 mb-3">
                  <Shield className="h-5 w-5 text-green-400" />
                  <span className="text-white font-semibold">Your Data is Safe</span>
                </div>
                <p className="text-gray-400 text-sm">
                  No payment information was stored or charged. Your personal and financial data remains secure.
                </p>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Footer Message */}
        <div className="text-center mt-12 pt-8 border-t border-gray-700">
          <p className="text-gray-400 text-sm">
            We're here to help make your shopping experience smooth. Don't hesitate to reach out if you need assistance.
          </p>
        </div>
      </div>
    </div>
  )
}
