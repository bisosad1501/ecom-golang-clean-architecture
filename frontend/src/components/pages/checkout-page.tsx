'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { CreditCard, Lock, ArrowLeft, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useCartStore, getCartTotal, getCartItemCount } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { formatPrice } from '@/lib/utils'
import { toast } from 'sonner'
import Image from 'next/image'

const checkoutSchema = z.object({
  // Shipping Information
  shipping_first_name: z.string().min(2, 'First name is required'),
  shipping_last_name: z.string().min(2, 'Last name is required'),
  shipping_email: z.string().email('Valid email is required'),
  shipping_phone: z.string().min(10, 'Valid phone number is required'),
  shipping_address: z.string().min(5, 'Address is required'),
  shipping_city: z.string().min(2, 'City is required'),
  shipping_state: z.string().min(2, 'State is required'),
  shipping_zip: z.string().min(5, 'ZIP code is required'),
  shipping_country: z.string().min(2, 'Country is required'),
  
  // Billing Information
  billing_same_as_shipping: z.boolean(),
  billing_first_name: z.string().optional(),
  billing_last_name: z.string().optional(),
  billing_address: z.string().optional(),
  billing_city: z.string().optional(),
  billing_state: z.string().optional(),
  billing_zip: z.string().optional(),
  billing_country: z.string().optional(),
  
  // Payment Information
  payment_method: z.enum(['credit_card', 'paypal', 'apple_pay']),
  card_number: z.string().optional(),
  card_expiry: z.string().optional(),
  card_cvc: z.string().optional(),
  card_name: z.string().optional(),
})

type CheckoutFormData = z.infer<typeof checkoutSchema>

export function CheckoutPage() {
  const router = useRouter()
  const [currentStep, setCurrentStep] = useState(1)
  const [isProcessing, setIsProcessing] = useState(false)
  const { cart, clearCart } = useCartStore()
  const { user, isAuthenticated } = useAuthStore()

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
    setValue,
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      shipping_first_name: user?.first_name || '',
      shipping_last_name: user?.last_name || '',
      shipping_email: user?.email || '',
      billing_same_as_shipping: true,
      payment_method: 'credit_card',
    },
  })

  const billingSameAsShipping = watch('billing_same_as_shipping')
  const paymentMethod = watch('payment_method')

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto text-center">
            <h1 className="text-3xl font-bold text-gray-900 mb-4">
              Your cart is empty
            </h1>
            <p className="text-gray-600 mb-8">
              Add some items to your cart before proceeding to checkout.
            </p>
            <Button asChild>
              <a href="/products">Continue Shopping</a>
            </Button>
          </div>
        </div>
      </div>
    )
  }

  const cartTotal = getCartTotal(cart)
  const cartItemCount = getCartItemCount(cart)
  const shippingCost = cartTotal > 50 ? 0 : 9.99
  const tax = cartTotal * 0.08
  const finalTotal = cartTotal + shippingCost + tax

  const onSubmit = async (data: CheckoutFormData) => {
    setIsProcessing(true)
    try {
      // Simulate payment processing
      await new Promise(resolve => setTimeout(resolve, 2000))
      
      // Clear cart and redirect to success page
      await clearCart()
      toast.success('Order placed successfully!')
      router.push('/order-confirmation')
    } catch (error) {
      toast.error('Payment failed. Please try again.')
    } finally {
      setIsProcessing(false)
    }
  }

  const steps = [
    { id: 1, name: 'Shipping', completed: currentStep > 1 },
    { id: 2, name: 'Payment', completed: currentStep > 2 },
    { id: 3, name: 'Review', completed: false },
  ]

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <button
            onClick={() => router.back()}
            className="flex items-center text-primary-600 hover:text-primary-700 mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Cart
          </button>
          <h1 className="text-3xl font-bold text-gray-900">Checkout</h1>
        </div>

        {/* Progress Steps */}
        <div className="mb-8">
          <div className="flex items-center justify-center">
            {steps.map((step, index) => (
              <div key={step.id} className="flex items-center">
                <div
                  className={`flex items-center justify-center w-10 h-10 rounded-full border-2 ${
                    step.completed
                      ? 'bg-primary-600 border-primary-600 text-white'
                      : currentStep === step.id
                      ? 'border-primary-600 text-primary-600'
                      : 'border-gray-300 text-gray-400'
                  }`}
                >
                  {step.completed ? (
                    <Check className="h-5 w-5" />
                  ) : (
                    <span>{step.id}</span>
                  )}
                </div>
                <span
                  className={`ml-2 text-sm font-medium ${
                    step.completed || currentStep === step.id
                      ? 'text-gray-900'
                      : 'text-gray-400'
                  }`}
                >
                  {step.name}
                </span>
                {index < steps.length - 1 && (
                  <div className="w-16 h-0.5 bg-gray-300 mx-4"></div>
                )}
              </div>
            ))}
          </div>
        </div>

        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-8">
              {/* Step 1: Shipping Information */}
              {currentStep === 1 && (
                <Card>
                  <CardHeader>
                    <CardTitle>Shipping Information</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      <Input
                        {...register('shipping_first_name')}
                        label="First Name"
                        error={errors.shipping_first_name?.message}
                        required
                      />
                      <Input
                        {...register('shipping_last_name')}
                        label="Last Name"
                        error={errors.shipping_last_name?.message}
                        required
                      />
                    </div>
                    
                    <Input
                      {...register('shipping_email')}
                      type="email"
                      label="Email"
                      error={errors.shipping_email?.message}
                      required
                    />
                    
                    <Input
                      {...register('shipping_phone')}
                      type="tel"
                      label="Phone Number"
                      error={errors.shipping_phone?.message}
                      required
                    />
                    
                    <Input
                      {...register('shipping_address')}
                      label="Address"
                      error={errors.shipping_address?.message}
                      required
                    />
                    
                    <div className="grid grid-cols-2 gap-4">
                      <Input
                        {...register('shipping_city')}
                        label="City"
                        error={errors.shipping_city?.message}
                        required
                      />
                      <Input
                        {...register('shipping_state')}
                        label="State"
                        error={errors.shipping_state?.message}
                        required
                      />
                    </div>
                    
                    <div className="grid grid-cols-2 gap-4">
                      <Input
                        {...register('shipping_zip')}
                        label="ZIP Code"
                        error={errors.shipping_zip?.message}
                        required
                      />
                      <Input
                        {...register('shipping_country')}
                        label="Country"
                        error={errors.shipping_country?.message}
                        required
                      />
                    </div>

                    <Button
                      type="button"
                      onClick={() => setCurrentStep(2)}
                      className="w-full"
                    >
                      Continue to Payment
                    </Button>
                  </CardContent>
                </Card>
              )}

              {/* Step 2: Payment Information */}
              {currentStep === 2 && (
                <div className="space-y-6">
                  {/* Billing Address */}
                  <Card>
                    <CardHeader>
                      <CardTitle>Billing Address</CardTitle>
                    </CardHeader>
                    <CardContent>
                      <label className="flex items-center space-x-2 mb-4">
                        <input
                          {...register('billing_same_as_shipping')}
                          type="checkbox"
                          className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                        />
                        <span>Same as shipping address</span>
                      </label>

                      {!billingSameAsShipping && (
                        <div className="space-y-4">
                          <div className="grid grid-cols-2 gap-4">
                            <Input
                              {...register('billing_first_name')}
                              label="First Name"
                              error={errors.billing_first_name?.message}
                            />
                            <Input
                              {...register('billing_last_name')}
                              label="Last Name"
                              error={errors.billing_last_name?.message}
                            />
                          </div>
                          <Input
                            {...register('billing_address')}
                            label="Address"
                            error={errors.billing_address?.message}
                          />
                          <div className="grid grid-cols-2 gap-4">
                            <Input
                              {...register('billing_city')}
                              label="City"
                              error={errors.billing_city?.message}
                            />
                            <Input
                              {...register('billing_state')}
                              label="State"
                              error={errors.billing_state?.message}
                            />
                          </div>
                          <div className="grid grid-cols-2 gap-4">
                            <Input
                              {...register('billing_zip')}
                              label="ZIP Code"
                              error={errors.billing_zip?.message}
                            />
                            <Input
                              {...register('billing_country')}
                              label="Country"
                              error={errors.billing_country?.message}
                            />
                          </div>
                        </div>
                      )}
                    </CardContent>
                  </Card>

                  {/* Payment Method */}
                  <Card>
                    <CardHeader>
                      <CardTitle>Payment Method</CardTitle>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      {/* Payment Options */}
                      <div className="space-y-3">
                        <label className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-gray-50">
                          <input
                            {...register('payment_method')}
                            type="radio"
                            value="credit_card"
                            className="text-primary-600 focus:ring-primary-500"
                          />
                          <CreditCard className="h-5 w-5" />
                          <span>Credit/Debit Card</span>
                        </label>
                        
                        <label className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-gray-50">
                          <input
                            {...register('payment_method')}
                            type="radio"
                            value="paypal"
                            className="text-primary-600 focus:ring-primary-500"
                          />
                          <div className="w-5 h-5 bg-blue-600 rounded flex items-center justify-center">
                            <span className="text-white text-xs font-bold">P</span>
                          </div>
                          <span>PayPal</span>
                        </label>
                        
                        <label className="flex items-center space-x-3 p-3 border rounded-lg cursor-pointer hover:bg-gray-50">
                          <input
                            {...register('payment_method')}
                            type="radio"
                            value="apple_pay"
                            className="text-primary-600 focus:ring-primary-500"
                          />
                          <div className="w-5 h-5 bg-black rounded flex items-center justify-center">
                            <span className="text-white text-xs">🍎</span>
                          </div>
                          <span>Apple Pay</span>
                        </label>
                      </div>

                      {/* Credit Card Form */}
                      {paymentMethod === 'credit_card' && (
                        <div className="space-y-4 pt-4 border-t">
                          <Input
                            {...register('card_name')}
                            label="Cardholder Name"
                            error={errors.card_name?.message}
                            placeholder="John Doe"
                          />
                          <Input
                            {...register('card_number')}
                            label="Card Number"
                            error={errors.card_number?.message}
                            placeholder="1234 5678 9012 3456"
                            leftIcon={<CreditCard className="h-4 w-4" />}
                          />
                          <div className="grid grid-cols-2 gap-4">
                            <Input
                              {...register('card_expiry')}
                              label="Expiry Date"
                              error={errors.card_expiry?.message}
                              placeholder="MM/YY"
                            />
                            <Input
                              {...register('card_cvc')}
                              label="CVC"
                              error={errors.card_cvc?.message}
                              placeholder="123"
                            />
                          </div>
                        </div>
                      )}

                      <div className="flex space-x-4">
                        <Button
                          type="button"
                          variant="outline"
                          onClick={() => setCurrentStep(1)}
                          className="flex-1"
                        >
                          Back
                        </Button>
                        <Button
                          type="button"
                          onClick={() => setCurrentStep(3)}
                          className="flex-1"
                        >
                          Review Order
                        </Button>
                      </div>
                    </CardContent>
                  </Card>
                </div>
              )}

              {/* Step 3: Review Order */}
              {currentStep === 3 && (
                <Card>
                  <CardHeader>
                    <CardTitle>Review Your Order</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-6">
                    {/* Order Items */}
                    <div>
                      <h4 className="font-semibold mb-4">Order Items</h4>
                      <div className="space-y-4">
                        {cart.items.map((item) => (
                          <div key={item.id} className="flex items-center space-x-4">
                            <div className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-md border">
                              <Image
                                src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                                alt={item.product.name}
                                fill
                                className="object-cover"
                              />
                            </div>
                            <div className="flex-1">
                              <h5 className="font-medium">{item.product.name}</h5>
                              <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                            </div>
                            <p className="font-semibold">{formatPrice(item.total_price)}</p>
                          </div>
                        ))}
                      </div>
                    </div>

                    {/* Security Notice */}
                    <div className="bg-green-50 border border-green-200 rounded-md p-4">
                      <div className="flex items-center space-x-2">
                        <Lock className="h-5 w-5 text-green-600" />
                        <span className="text-sm text-green-800">
                          Your payment information is secure and encrypted
                        </span>
                      </div>
                    </div>

                    <div className="flex space-x-4">
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setCurrentStep(2)}
                        className="flex-1"
                      >
                        Back
                      </Button>
                      <Button
                        type="submit"
                        disabled={isProcessing}
                        isLoading={isProcessing}
                        loadingText="Processing..."
                        className="flex-1"
                      >
                        Place Order
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Order Summary Sidebar */}
            <div className="lg:col-span-1">
              <div className="sticky top-8">
                <Card>
                  <CardHeader>
                    <CardTitle>Order Summary</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="flex justify-between">
                      <span>Subtotal ({cartItemCount} items)</span>
                      <span>{formatPrice(cartTotal)}</span>
                    </div>
                    
                    <div className="flex justify-between">
                      <span>Shipping</span>
                      <span>
                        {shippingCost === 0 ? (
                          <span className="text-green-600">Free</span>
                        ) : (
                          formatPrice(shippingCost)
                        )}
                      </span>
                    </div>
                    
                    <div className="flex justify-between">
                      <span>Tax</span>
                      <span>{formatPrice(tax)}</span>
                    </div>
                    
                    <div className="border-t pt-4">
                      <div className="flex justify-between text-lg font-semibold">
                        <span>Total</span>
                        <span>{formatPrice(finalTotal)}</span>
                      </div>
                    </div>

                    {cartTotal < 50 && (
                      <div className="bg-blue-50 border border-blue-200 rounded-md p-3">
                        <p className="text-sm text-blue-800">
                          Add {formatPrice(50 - cartTotal)} more to get free shipping!
                        </p>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  )
}
