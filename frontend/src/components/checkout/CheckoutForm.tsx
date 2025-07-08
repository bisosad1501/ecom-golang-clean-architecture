'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Badge } from '@/components/ui/badge'
import { useCartStore, getCartTotal, getCartItemCount } from '@/store/cart'
import { useOrderStore } from '@/store/order'
import { usePaymentStore } from '@/store/payment'
import { useAuthStore } from '@/store/auth'
import { redirectToCheckout } from '@/lib/stripe'
import { formatPrice } from '@/lib/utils'
import { Loader2, CreditCard, MapPin, ShoppingBag, Lock, ArrowLeft, Truck, Shield, CheckCircle, AlertCircle } from 'lucide-react'
import { toast } from 'sonner'

const addressSchema = z.object({
  street: z.string().min(1, 'Street address is required'),
  city: z.string().min(1, 'City is required'),
  state: z.string().min(1, 'State is required'),
  zip_code: z.string().min(1, 'ZIP code is required'),
  country: z.string().min(1, 'Country is required'),
})

const checkoutSchema = z.object({
  shipping_address: addressSchema,
  coupon_code: z.string().optional(),
  notes: z.string().optional(),
})

type CheckoutFormData = z.infer<typeof checkoutSchema>

export default function CheckoutForm() {
  const router = useRouter()
  const [isProcessing, setIsProcessing] = useState(false)

  const { cart, fetchCart } = useCartStore()
  const { createOrder } = useOrderStore()
  const paymentStore = usePaymentStore()
  const { user, isAuthenticated } = useAuthStore()

  // Fetch cart data when component mounts
  useEffect(() => {
    if (isAuthenticated) {
      fetchCart()
    }
  }, [fetchCart, isAuthenticated])

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      shipping_address: {
        country: 'US',
      },
    },
  })



  const onSubmit = async (data: CheckoutFormData) => {
    if (!cart || cart.items.length === 0) {
      toast.error('Your cart is empty')
      return
    }

    setIsProcessing(true)
    
    try {
      // Create order first
      const orderData = {
        shipping_address: data.shipping_address,
        billing_address: data.shipping_address, // Use shipping address as billing
        payment_method: 'stripe',
        tax_rate: 0.08,
        shipping_cost: shippingCost,
        discount_amount: 0,
        notes: data.notes,
      }

      const orderResponse = await createOrder(orderData)
      const orderId = orderResponse.data.id

      // Create Stripe checkout session
      const checkoutData = {
        order_id: orderId,
        amount: finalTotal,
        currency: 'usd',
        description: `Order ${orderResponse.data.order_number}`,
        success_url: `${window.location.origin}/checkout/success?session_id={CHECKOUT_SESSION_ID}&order_id=${orderId}`,
        cancel_url: `${window.location.origin}/checkout/cancel?order_id=${orderId}`,
        metadata: {
          order_id: orderId,
          order_number: orderResponse.data.order_number,
        },
      }

      const session = await paymentStore.createCheckoutSession(checkoutData)

      // Redirect to Stripe Checkout
      if (session.session_id) {
        await redirectToCheckout(session.session_id)
      } else if (session.session_url) {
        window.location.href = session.session_url
      } else {
        throw new Error('No session ID or URL received from payment service')
      }
      
    } catch (error: any) {
      console.error('Checkout error:', error)
      toast.error(error.message || 'Failed to process checkout')
    } finally {
      setIsProcessing(false)
    }
  }

  const cartTotal = getCartTotal(cart)
  const cartItemCount = getCartItemCount(cart)
  const shippingCost = cartTotal > 50 ? 0 : 9.99
  const tax = cartTotal * 0.08 // 8% tax
  const finalTotal = cartTotal + shippingCost + tax

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white">
        <div className="container mx-auto px-4 py-16">
          <div className="max-w-md mx-auto text-center">
            <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-2xl p-8">
              <Lock className="h-16 w-16 text-orange-500 mx-auto mb-6" />
              <h2 className="text-2xl font-bold mb-4">Login Required</h2>
              <p className="text-gray-400 mb-6">Please log in to continue with checkout</p>
              <Button
                onClick={() => router.push('/auth/login')}
                className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700"
              >
                Login to Continue
              </Button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white">
        <div className="container mx-auto px-4 py-16">
          <div className="max-w-md mx-auto text-center">
            <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-2xl p-8">
              <ShoppingBag className="h-16 w-16 text-orange-500 mx-auto mb-6" />
              <h2 className="text-2xl font-bold mb-4">Your Cart is Empty</h2>
              <p className="text-gray-400 mb-6">Add some items to your cart before checkout</p>
              <Button
                onClick={() => router.push('/products')}
                className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700"
              >
                Continue Shopping
              </Button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white">
      {/* Enhanced BiHub Header with Progress */}
      <div className="bg-black/50 backdrop-blur-sm border-b border-gray-700">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center space-x-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => router.back()}
                className="text-gray-400 hover:text-white transition-all duration-300 hover:scale-105"
              >
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Cart
              </Button>
              <div>
                <h1 className="text-3xl font-bold">
                  Bi<span className="bg-gradient-to-r from-orange-400 to-orange-600 bg-clip-text text-transparent">Hub</span> Checkout
                </h1>
                <p className="text-gray-400 text-sm">Secure payment processing • Step 1 of 2</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <Shield className="h-5 w-5 text-green-500" />
                <span className="text-sm text-gray-400">256-bit SSL</span>
              </div>
              <div className="flex items-center space-x-2">
                <Lock className="h-5 w-5 text-blue-500" />
                <span className="text-sm text-gray-400">Encrypted</span>
              </div>
            </div>
          </div>

          {/* Progress Steps */}
          <div className="flex items-center justify-center space-x-8">
            <div className="flex items-center space-x-3">
              <div className="w-10 h-10 bg-gradient-to-r from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-lg">
                <span className="text-white font-bold">1</span>
              </div>
              <div>
                <p className="text-white font-semibold">Shipping Details</p>
                <p className="text-gray-400 text-xs">Current step</p>
              </div>
            </div>

            <div className="flex-1 h-1 bg-gray-700 rounded-full mx-4">
              <div className="h-1 bg-gradient-to-r from-orange-500 to-orange-600 rounded-full w-0 transition-all duration-1000"></div>
            </div>

            <div className="flex items-center space-x-3 opacity-50">
              <div className="w-10 h-10 bg-gray-700 rounded-full flex items-center justify-center">
                <span className="text-gray-400 font-bold">2</span>
              </div>
              <div>
                <p className="text-gray-400 font-semibold">Payment</p>
                <p className="text-gray-500 text-xs">Next step</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Checkout Form */}
            <div className="lg:col-span-2 space-y-6">
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
                {/* Shipping Address */}
                <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 hover:border-orange-500/50 transition-all duration-300">
                  <CardHeader className="pb-4">
                    <CardTitle className="flex items-center gap-3 text-white">
                      <div className="p-2 bg-orange-500/20 rounded-lg">
                        <MapPin className="h-5 w-5 text-orange-500" />
                      </div>
                      Shipping Address
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="shipping_street" className="text-gray-300 font-medium flex items-center gap-2">
                        Street Address
                        <span className="text-red-400">*</span>
                      </Label>
                      <div className="relative">
                        <Input
                          id="shipping_street"
                          {...register('shipping_address.street')}
                          placeholder="123 Main St, Apartment, suite, etc."
                          className={`bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20 transition-all duration-300 ${
                            errors.shipping_address?.street ? 'border-red-500 focus:border-red-500' : ''
                          }`}
                          defaultValue={user?.address || ''}
                        />
                        {!errors.shipping_address?.street && (
                          <CheckCircle className="absolute right-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-green-500 opacity-0 transition-opacity duration-300" />
                        )}
                        {errors.shipping_address?.street && (
                          <AlertCircle className="absolute right-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-red-500" />
                        )}
                      </div>
                      {errors.shipping_address?.street && (
                        <p className="text-sm text-red-400 mt-1 flex items-center gap-1 animate-slide-in-right">
                          <AlertCircle className="h-3 w-3" />
                          {errors.shipping_address.street.message}
                        </p>
                      )}
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="shipping_city" className="text-gray-300 font-medium">City</Label>
                        <Input
                          id="shipping_city"
                          {...register('shipping_address.city')}
                          placeholder="New York"
                          className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20"
                        />
                        {errors.shipping_address?.city && (
                          <p className="text-sm text-red-400 mt-1">
                            {errors.shipping_address.city.message}
                          </p>
                        )}
                      </div>
                      <div>
                        <Label htmlFor="shipping_state" className="text-gray-300 font-medium">State</Label>
                        <Input
                          id="shipping_state"
                          {...register('shipping_address.state')}
                          placeholder="NY"
                          className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20"
                        />
                        {errors.shipping_address?.state && (
                          <p className="text-sm text-red-400 mt-1">
                            {errors.shipping_address.state.message}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="shipping_zip" className="text-gray-300 font-medium">ZIP Code</Label>
                        <Input
                          id="shipping_zip"
                          {...register('shipping_address.zip_code')}
                          placeholder="10001"
                          className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20"
                        />
                        {errors.shipping_address?.zip_code && (
                          <p className="text-sm text-red-400 mt-1">
                            {errors.shipping_address.zip_code.message}
                          </p>
                        )}
                      </div>
                      <div>
                        <Label htmlFor="shipping_country" className="text-gray-300 font-medium">Country</Label>
                        <Input
                          id="shipping_country"
                          {...register('shipping_address.country')}
                          placeholder="US"
                          className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20"
                        />
                        {errors.shipping_address?.country && (
                          <p className="text-sm text-red-400 mt-1">
                            {errors.shipping_address.country.message}
                          </p>
                        )}
                      </div>
                    </div>
                  </CardContent>
                </Card>





                {/* Additional Options */}
                <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 hover:border-orange-500/50 transition-all duration-300">
                  <CardContent className="pt-6 space-y-4">
                    <div>
                      <Label htmlFor="coupon_code" className="text-gray-300 font-medium">Coupon Code (Optional)</Label>
                      <Input
                        id="coupon_code"
                        {...register('coupon_code')}
                        placeholder="Enter coupon code"
                        className="bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500 focus:ring-orange-500/20"
                      />
                    </div>

                    <div>
                      <Label htmlFor="notes" className="text-gray-300 font-medium">Order Notes (Optional)</Label>
                      <textarea
                        id="notes"
                        {...register('notes')}
                        placeholder="Special instructions for your order"
                        className="w-full px-3 py-2 bg-gray-800/50 border border-gray-600 text-white placeholder:text-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500"
                        rows={3}
                      />
                    </div>
                  </CardContent>
                </Card>

                {/* Payment Method */}
                <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700">
                  <CardHeader className="pb-4">
                    <CardTitle className="flex items-center gap-3 text-white">
                      <div className="p-2 bg-blue-500/20 rounded-lg">
                        <CreditCard className="h-5 w-5 text-blue-400" />
                      </div>
                      Payment Method
                    </CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="bg-gradient-to-r from-blue-500/10 to-purple-500/10 border border-blue-500/20 rounded-xl p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="w-12 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-md flex items-center justify-center">
                            <span className="text-white font-bold text-xs">STRIPE</span>
                          </div>
                          <div>
                            <p className="text-white font-medium">Secure Card Payment</p>
                            <p className="text-gray-400 text-sm">Visa, Mastercard, American Express</p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Shield className="h-4 w-4 text-green-400" />
                          <span className="text-green-400 text-sm font-medium">Secured</span>
                        </div>
                      </div>
                      <div className="mt-3 pt-3 border-t border-gray-700">
                        <p className="text-xs text-gray-400">
                          🔒 Your payment information is encrypted and secure. We never store your card details.
                        </p>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Submit Button */}
                <div className="pt-4">
                  <Button
                    type="submit"
                    className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold py-4 px-6 rounded-xl transition-all duration-300 transform hover:scale-[1.02] hover:shadow-2xl hover:shadow-orange-500/25"
                    size="lg"
                    disabled={isProcessing}
                  >
                    {isProcessing ? (
                      <>
                        <Loader2 className="mr-3 h-5 w-5 animate-spin" />
                        Creating Order...
                      </>
                    ) : (
                      <>
                        <CreditCard className="mr-3 h-5 w-5" />
                        Pay ${finalTotal.toFixed(2)} with Stripe
                      </>
                    )}
                  </Button>

                  <div className="text-center mt-4">
                    <p className="text-xs text-gray-400">
                      By clicking "Pay", you agree to our{' '}
                      <a href="/terms" className="text-orange-400 hover:text-orange-300 underline">Terms of Service</a>
                      {' '}and{' '}
                      <a href="/privacy" className="text-orange-400 hover:text-orange-300 underline">Privacy Policy</a>
                    </p>
                  </div>
                </div>
              </form>
            </div>

            {/* Order Summary */}
            <div className="space-y-6">
              <Card className="bg-gray-900/50 backdrop-blur-sm border-gray-700 sticky top-6">
                <CardHeader className="pb-4">
                  <CardTitle className="flex items-center gap-3 text-white">
                    <div className="p-2 bg-orange-500/20 rounded-lg">
                      <ShoppingBag className="h-5 w-5 text-orange-500" />
                    </div>
                    Order Summary
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  {/* Cart Items */}
                  <div className="space-y-3 max-h-64 overflow-y-auto">
                    {cart.items.map((item) => (
                      <div key={item.id} className="flex items-center space-x-3 p-3 bg-gray-800/30 rounded-lg">
                        <div className="relative h-12 w-12 flex-shrink-0 overflow-hidden rounded-lg border border-gray-600">
                          <Image
                            src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                            alt={item.product.name}
                            fill
                            className="object-cover"
                          />
                        </div>
                        <div className="flex-1 min-w-0">
                          <p className="font-medium text-white truncate">{item.product.name}</p>
                          <p className="text-sm text-gray-400">Qty: {item.quantity}</p>
                        </div>
                        <p className="font-semibold text-orange-400">{formatPrice(item.subtotal || item.price * item.quantity)}</p>
                      </div>
                    ))}
                  </div>

                  <Separator className="bg-gray-700" />

                  {/* Order Totals */}
                  <div className="space-y-3">
                    <div className="flex justify-between text-gray-300">
                      <span>Subtotal ({cartItemCount} items)</span>
                      <span>{formatPrice(cartTotal)}</span>
                    </div>
                    <div className="flex justify-between text-gray-300">
                      <span className="flex items-center gap-2">
                        <Truck className="h-4 w-4" />
                        Shipping
                      </span>
                      <span className={shippingCost === 0 ? "text-green-400 font-medium" : ""}>
                        {shippingCost === 0 ? "FREE" : formatPrice(shippingCost)}
                      </span>
                    </div>
                    <div className="flex justify-between text-gray-300">
                      <span>Tax (8%)</span>
                      <span>{formatPrice(tax)}</span>
                    </div>

                    <Separator className="bg-gray-700" />

                    <div className="flex justify-between text-xl font-bold text-white">
                      <span>Total</span>
                      <span className="text-orange-400">{formatPrice(finalTotal)}</span>
                    </div>
                  </div>

                  {/* Security Notice */}
                  <div className="bg-green-900/20 border border-green-700/50 rounded-lg p-3 mt-4">
                    <div className="flex items-center space-x-2">
                      <Lock className="h-4 w-4 text-green-400" />
                      <span className="text-sm text-green-300">
                        Your payment information is secure and encrypted
                      </span>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
