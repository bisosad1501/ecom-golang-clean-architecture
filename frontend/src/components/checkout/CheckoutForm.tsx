'use client'

import { useState, useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useCartStore } from '@/store/cart'
import { useOrderStore } from '@/store/order'
import { usePaymentStore } from '@/store/payment'
import { useAuthStore } from '@/store/auth'
import { redirectToCheckout } from '@/lib/stripe'
import { formatPrice } from '@/lib/utils'
import { Loader2, CreditCard, Lock, Shield, ChevronLeft, ChevronRight } from 'lucide-react'
import { toast } from 'sonner'
import { CheckoutProgress, CheckoutProgressMobile } from './CheckoutProgress'
import {
  ContactInfoStep,
  ShippingAddressStep,
  BillingAddressStep,
  PaymentMethodStep,
  ShippingOptionsStep,
  ReviewStep
} from './CheckoutSteps'

// Schema matching Backend exactly
const addressSchema = z.object({
  first_name: z.string().min(1, 'First name is required'),
  last_name: z.string().min(1, 'Last name is required'),
  company: z.string().optional(),
  address1: z.string().min(1, 'Street address is required'),
  address2: z.string().optional(),
  city: z.string().min(1, 'City is required'),
  state: z.string().min(1, 'State is required'),
  zip_code: z.string().min(1, 'ZIP code is required'),
  country: z.string().min(1, 'Country is required'),
  phone: z.string().optional(),
})

const checkoutSchema = z.object({
  // Contact Information
  email: z.string().email('Valid email is required'),
  phone: z.string().min(1, 'Phone number is required'),

  // Shipping Address
  shipping_address: addressSchema,

  // Billing Address
  use_shipping_for_billing: z.boolean().default(true),
  billing_address: z.object({
    first_name: z.string().optional(),
    last_name: z.string().optional(),
    company: z.string().optional(),
    address1: z.string().optional(),
    address2: z.string().optional(),
    city: z.string().optional(),
    state: z.string().optional(),
    zip_code: z.string().optional(),
    country: z.string().optional(),
    phone: z.string().optional(),
  }).optional(),
  
  // Payment Method
  payment_method: z.enum(['stripe', 'cash', 'bank_transfer']).default('stripe'),

  // Shipping Options
  shipping_method: z.enum(['standard', 'express', 'overnight']).default('standard'),

  // Gift Options
  is_gift: z.boolean().default(false),
  gift_message: z.string().optional(),
  gift_wrap: z.boolean().default(false),
  
  // Order Preferences
  delivery_instructions: z.string().optional(),
  priority: z.enum(['normal', 'high', 'urgent']).default('normal'),
  tip_amount: z.number().min(0).default(0),
  
  // Promo & Notes
  coupon_code: z.string().optional(),
  notes: z.string().optional(),
}).refine((data) => {
  // If not using shipping for billing, billing address fields are required
  if (!data.use_shipping_for_billing) {
    const billing = data.billing_address
    if (!billing) return false
    return !!(billing.first_name && billing.last_name && billing.address1 &&
             billing.city && billing.state && billing.zip_code)
  }
  return true
}, {
  message: "Complete billing address is required when not using shipping address",
  path: ["billing_address"]
})

type CheckoutFormData = z.infer<typeof checkoutSchema>

export default function CheckoutForm() {
  const router = useRouter()
  const [isProcessing, setIsProcessing] = useState(false)
  const [processingMessage, setProcessingMessage] = useState('')
  const [currentStep, setCurrentStep] = useState(1)
  
  const steps = [
    { id: 1, title: 'Contact Info', description: 'Your contact information' },
    { id: 2, title: 'Shipping', description: 'Where to send your order' },
    { id: 3, title: 'Billing', description: 'Payment address' },
    { id: 4, title: 'Payment', description: 'Choose payment method' },
    { id: 5, title: 'Options', description: 'Shipping & gift options' },
    { id: 6, title: 'Review', description: 'Review and pay' },
  ]

  const { cart, fetchCart } = useCartStore()
  const { createOrder } = useOrderStore()
  const paymentStore = usePaymentStore()
  const { user, isAuthenticated, token } = useAuthStore()

  // Check authentication
  useEffect(() => {
    if (!isAuthenticated || !token) {
      toast.error('Please log in to continue with checkout')
      router.push('/auth/login?redirect=/checkout')
      return
    }
  }, [isAuthenticated, token, router])

  // Fetch cart data when component mounts
  useEffect(() => {
    // Always fetch cart - it will handle guest vs authenticated users
    fetchCart()
  }, [fetchCart])

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    getValues,
    formState: { errors },
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      email: user?.email || '',
      phone: '',
      shipping_address: {
        first_name: user?.first_name || '',
        last_name: user?.last_name || '',
        company: '',
        address1: '',
        address2: '',
        city: '',
        state: '',
        zip_code: '',
        country: 'US',
        phone: '',
      },
      billing_address: {
        first_name: '',
        last_name: '',
        company: '',
        address1: '',
        address2: '',
        city: '',
        state: '',
        zip_code: '',
        country: 'US',
        phone: '',
      },
      use_shipping_for_billing: true,
      payment_method: 'stripe',
      shipping_method: 'standard',
      is_gift: false,
      gift_message: '',
      gift_wrap: false,
      delivery_instructions: '',
      priority: 'normal',
      tip_amount: 0,
      coupon_code: '',
      notes: '',
    },
  })

  // Mock cart for testing if empty
  const mockCart = {
    id: 'mock-cart',
    user_id: user?.id || 'test-user',
    items: [
      {
        id: 'mock-item-1',
        product: {
          id: 'mock-product-1',
          name: 'Test Product',
          price: 99.99,
          sku: 'TEST-001'
        },
        quantity: 2,
        price: 99.99,
        total: 199.98
      }
    ],
    total: 199.98,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  }

  const effectiveCart = cart && cart.items.length > 0 ? cart : mockCart

  // Watch form values for dynamic updates
  const watchedValues = watch()
  const useShippingForBilling = watch('use_shipping_for_billing')
  const shippingAddress = watch('shipping_address')
  const isGift = watch('is_gift')
  const shippingMethod = watch('shipping_method')

  // Copy shipping address to billing address when checkbox is checked
  useEffect(() => {
    console.log('🔄 useEffect triggered:', { useShippingForBilling, shippingAddress })
    if (useShippingForBilling && shippingAddress) {
      console.log('✅ Copying shipping to billing:', shippingAddress)
      setValue('billing_address', shippingAddress)
    }
  }, [useShippingForBilling, shippingAddress, setValue])

  // Navigation functions
  const nextStep = () => {
    if (currentStep < steps.length) {
      setCurrentStep(currentStep + 1)
    }
  }

  const prevStep = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1)
    }
  }

  const canProceedToNext = () => {
    switch (currentStep) {
      case 1: // Contact Info
        return watchedValues.email && watchedValues.phone
      case 2: // Shipping Address
        return watchedValues.shipping_address?.first_name && 
               watchedValues.shipping_address?.last_name &&
               watchedValues.shipping_address?.address1 &&
               watchedValues.shipping_address?.city &&
               watchedValues.shipping_address?.state &&
               watchedValues.shipping_address?.zip_code
      case 3: // Billing Address
        return true // Always can proceed (using shipping address)
      case 4: // Payment Method
        return watchedValues.payment_method
      case 5: // Shipping Options
        return watchedValues.shipping_method
      case 6: // Review
        return true
      default:
        return false
    }
  }

  const onSubmit = async (data: CheckoutFormData) => {
    console.log('🚀 onSubmit called!')
    console.log('Cart data:', cart)
    console.log('Form data:', data)
    console.log('Form errors:', errors)

    if (!cart || cart.items.length === 0) {
      toast.error('Your cart is empty')
      return
    }

    setIsProcessing(true)

    try {
      // Calculate shipping cost based on method
      const shippingCosts = {
        standard: cartTotal > 50 ? 0 : 9.99,
        express: 9.99,
        overnight: 19.99,
      }
      const selectedShippingCost = shippingCosts[data.shipping_method]

      // Handle different payment methods
      if (data.payment_method === 'stripe') {
        setProcessingMessage('Creating secure checkout session...')

        // Use new checkout session API for Stripe
        const checkoutSessionData = {
          shipping_address: data.shipping_address,
          billing_address: data.use_shipping_for_billing ? data.shipping_address : data.billing_address,
          payment_method: 'stripe' as const,
          notes: data.notes || '',
          tax_rate: 0.08,
          shipping_cost: selectedShippingCost,
          discount_amount: 0,
        }

        console.log('Checkout session data being sent:', JSON.stringify(checkoutSessionData, null, 2))

        // Create checkout session (includes Stripe URL)
        setProcessingMessage('Connecting to Stripe...')

        // Get token from auth store
        const token = useAuthStore.getState().token
        console.log('Token available:', !!token)
        console.log('Token length:', token?.length)
        console.log('Is authenticated:', isAuthenticated)

        if (!token) {
          throw new Error('Please log in to continue with checkout')
        }

        const response = await fetch('http://localhost:8080/api/v1/checkout/session', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify(checkoutSessionData),
        })

        const sessionResult = await response.json()

        if (!response.ok) {
          throw new Error(sessionResult.error || 'Failed to create checkout session')
        }

        // Redirect to Stripe hosted page
        if (sessionResult.data?.stripe_url) {
          setProcessingMessage('Redirecting to secure payment page...')
          toast.success('Redirecting to Stripe for secure payment...')
          setTimeout(() => {
            window.location.href = sessionResult.data.stripe_url
          }, 1000)
        } else {
          throw new Error('No Stripe URL received from checkout session')
        }
      } else {
        setProcessingMessage(`Creating ${data.payment_method === 'cash' ? 'COD' : 'bank transfer'} order...`)

        // Use direct order creation for COD and Bank Transfer
        const orderData = {
          shipping_address: data.shipping_address,
          billing_address: data.use_shipping_for_billing ? data.shipping_address : data.billing_address,
          payment_method: data.payment_method,
          notes: data.notes || '',
          tax_rate: 0.08,
          shipping_cost: selectedShippingCost,
          discount_amount: 0,
        }

        console.log('Order data being sent:', JSON.stringify(orderData, null, 2))

        // Get token from auth store
        const token = useAuthStore.getState().token
        if (!token) {
          throw new Error('Please log in to continue with checkout')
        }

        let orderResponse
        if (data.payment_method === 'cash') {
          orderResponse = await fetch('http://localhost:8080/api/v1/checkout/cod', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(orderData),
          })
        } else {
          orderResponse = await fetch('http://localhost:8080/api/v1/orders', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(orderData),
          })
        }

        const orderResult = await orderResponse.json()

        if (!orderResponse.ok) {
          throw new Error(orderResult.error || 'Failed to create order')
        }

        // Redirect to order confirmation
        const order = orderResult.data || orderResult
        setProcessingMessage('Order created successfully! Redirecting...')
        toast.success(`${data.payment_method === 'cash' ? 'COD' : 'Bank transfer'} order created successfully!`)
        setTimeout(() => {
          router.push(`/order-confirmation?order_id=${order.id}`)
        }, 1000)
      }
      
    } catch (error: any) {
      console.error('Checkout error:', error)
      toast.error(error.message || 'Failed to process checkout')
    } finally {
      setIsProcessing(false)
    }
  }

  // Calculate costs
  const getCartTotal = (cart: any) => {
    if (!cart?.items) return 0
    return cart.items.reduce((total: number, item: any) => {
      // Use current_price (backend computed) for accurate pricing
      const itemPrice = item.product?.current_price || item.product?.price || 0
      return total + itemPrice * item.quantity
    }, 0)
  }

  const getCartItemCount = (cart: any) => {
    if (!cart?.items) return 0
    return cart.items.reduce((total: number, item: any) => total + item.quantity, 0)
  }

  const cartTotal = getCartTotal(effectiveCart)
  const cartItemCount = getCartItemCount(effectiveCart)
  
  // Calculate shipping cost based on selected method
  const getShippingCost = (method: string) => {
    switch (method) {
      case 'standard':
        return cartTotal > 50 ? 0 : 9.99
      case 'express':
        return 9.99
      case 'overnight':
        return 19.99
      default:
        return 9.99
    }
  }
  
  const shippingCost = getShippingCost(shippingMethod)
  const tipAmount = watchedValues.tip_amount || 0
  const tax = cartTotal * 0.08 // 8% tax
  const finalTotal = cartTotal + shippingCost + tax + tipAmount

  // Allow guest checkout - remove authentication requirement
  // if (!isAuthenticated) {
  //   return (
  //     <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white flex items-center justify-center">
  //       <Card className="bg-gray-900/50 border-gray-800 p-8">
  //         <CardContent className="text-center">
  //           <h2 className="text-2xl font-bold mb-4">Please Sign In</h2>
  //           <p className="text-gray-400 mb-6">You need to be signed in to checkout</p>
  //           <Button onClick={() => router.push('/auth/login')} className="bg-[#ff9000] hover:bg-[#e68100]">
  //             Sign In
  //           </Button>
  //         </CardContent>
  //       </Card>
  //     </div>
  //   )
  // }

  if (!effectiveCart || effectiveCart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white flex items-center justify-center">
        <Card className="bg-gray-900/50 border-gray-800 p-8">
          <CardContent className="text-center">
            <h2 className="text-2xl font-bold mb-4">Your Cart is Empty</h2>
            <p className="text-gray-400 mb-6">Add some items to your cart before checkout</p>
            <Button onClick={() => router.push('/products')} className="bg-[#ff9000] hover:bg-[#e68100]">
              Continue Shopping
            </Button>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      {/* Enhanced Background */}
      <AnimatedBackground className="opacity-30" />
      
      <div className="container mx-auto px-4 lg:px-6 py-8 relative z-10">
        <div className="w-full space-y-8">
          {/* Header */}
          <div className="max-w-6xl mx-auto">
            <div className="mb-6">
              <div className="flex items-center justify-between gap-4">
                <div className="flex items-center space-x-4">
                  <Button
                    variant="ghost"
                    onClick={() => currentStep > 1 ? prevStep() : router.push('/cart')}
                    className="text-gray-400 hover:text-white"
                  >
                    <ChevronLeft className="h-4 w-4 mr-2" />
                    {currentStep > 1 ? 'Previous' : 'Back to Cart'}
                  </Button>
                  <h1 className="text-4xl font-bold text-white">
                    Checkout
                  </h1>
                </div>
                
                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-2">
                    <Shield className="h-4 w-4 text-green-400" />
                    <span className="text-sm text-gray-400">256-bit SSL</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Lock className="h-4 w-4 text-blue-400" />
                    <span className="text-sm text-gray-400">Encrypted</span>
                  </div>
                </div>
              </div>
            </div>

            {/* Progress Indicator */}
            <div className="hidden md:block">
              <CheckoutProgress steps={steps} currentStep={currentStep} />
            </div>
            <div className="md:hidden">
              <CheckoutProgressMobile steps={steps} currentStep={currentStep} />
            </div>
          </div>

          {/* Main Content */}
          <div className="max-w-6xl mx-auto">
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
              {/* Checkout Form */}
              <div className="lg:col-span-2 space-y-6">
                <form onSubmit={handleSubmit((data) => {
                  console.log('🚀 onSubmit called!')
                  console.log('Cart data:', cart)
                  console.log('Form data:', data)
                  console.log('Form errors:', errors)
                  onSubmit(data)
                }, (errors) => {
                  console.log('❌ Form validation failed:', errors)
                })} className="space-y-6">
                  {/* Step Content */}
                  {currentStep === 1 && (
                    <ContactInfoStep 
                      register={register} 
                      errors={errors} 
                      watch={watch} 
                      setValue={setValue} 
                    />
                  )}
                  
                  {currentStep === 2 && (
                    <ShippingAddressStep 
                      register={register} 
                      errors={errors} 
                      watch={watch} 
                      setValue={setValue} 
                    />
                  )}
                  
                  {currentStep === 3 && (
                    <BillingAddressStep 
                      register={register} 
                      errors={errors} 
                      watch={watch} 
                      setValue={setValue} 
                    />
                  )}
                  
                  {currentStep === 4 && (
                    <PaymentMethodStep
                      register={register}
                      errors={errors}
                      watch={watch}
                      setValue={setValue}
                    />
                  )}

                  {currentStep === 5 && (
                    <ShippingOptionsStep
                      register={register}
                      errors={errors}
                      watch={watch}
                      setValue={setValue}
                    />
                  )}

                  {currentStep === 6 && (
                    <ReviewStep 
                      register={register} 
                      errors={errors} 
                      watch={watch} 
                      setValue={setValue}
                      cartTotal={cartTotal}
                      shippingCost={shippingCost}
                      tax={tax}
                      finalTotal={finalTotal}
                      cartItems={effectiveCart?.items || []}
                    />
                  )}

                  {/* Debug Info */}
                  <div className="text-sm text-gray-500 mb-4">
                    Current Step: {currentStep} / 6
                  </div>

                  {/* Navigation Buttons */}
                  <div className="flex justify-between pt-6">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={prevStep}
                      disabled={currentStep === 1}
                      className="border-gray-700 text-gray-300 hover:bg-gray-800"
                    >
                      <ChevronLeft className="h-4 w-4 mr-2" />
                      Previous
                    </Button>

                    {currentStep < 6 ? (
                      <Button
                        type="button"
                        onClick={() => {
                          console.log('🔄 Next button clicked, currentStep:', currentStep)
                          nextStep()
                        }}
                        disabled={!canProceedToNext()}
                        className="bg-[#ff9000] hover:bg-[#e68100] text-white"
                      >
                        Next
                        <ChevronRight className="h-4 w-4 ml-2" />
                      </Button>
                    ) : (
                      <Button
                        type="submit"
                        disabled={isProcessing}
                        className="bg-[#ff9000] hover:bg-[#e68100] text-white"
                        onClick={() => {
                          console.log('🔥 Complete Order button clicked!')
                          console.log('Current form values:', getValues())
                          console.log('Form is valid:', formState.isValid)
                          console.log('Form errors:', errors)
                        }}
                      >
                        {isProcessing ? (
                          <>
                            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                            {processingMessage || 'Processing...'}
                          </>
                        ) : (
                          <>
                            Complete Order
                            <CreditCard className="h-4 w-4 ml-2" />
                          </>
                        )}
                      </Button>
                    )}
                  </div>
                </form>
              </div>

              {/* Order Summary */}
              <div className="lg:col-span-1">
                <Card className="bg-gray-900/50 border-gray-800 sticky top-8">
                  <CardHeader>
                    <CardTitle className="text-white">Order Summary</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    {/* Cart Items */}
                    <div className="space-y-3">
                      {effectiveCart.items.map((item: any) => (
                        <div key={item.id} className="flex items-center justify-between py-2 border-b border-gray-700 last:border-b-0">
                          <div className="flex items-center space-x-3">
                            <div className="w-12 h-12 bg-gray-800 rounded-lg flex items-center justify-center">
                              <span className="text-gray-400 text-xs">{item.quantity}</span>
                            </div>
                            <div>
                              <div className="text-white font-medium text-sm">{item.product?.name}</div>
                              <div className="text-gray-400 text-xs">Qty: {item.quantity}</div>
                            </div>
                          </div>
                          <div className="text-white font-semibold text-sm">
                            {formatPrice(item.product?.price * item.quantity)}
                          </div>
                        </div>
                      ))}
                    </div>
                    
                    {/* Cost Breakdown */}
                    <div className="space-y-2 pt-4 border-t border-gray-700">
                      <div className="flex justify-between text-gray-300">
                        <span>Subtotal</span>
                        <span>{formatPrice(cartTotal)}</span>
                      </div>
                      <div className="flex justify-between text-gray-300">
                        <span>Shipping</span>
                        <span>{shippingCost === 0 ? 'Free' : formatPrice(shippingCost)}</span>
                      </div>
                      <div className="flex justify-between text-gray-300">
                        <span>Tax</span>
                        <span>{formatPrice(tax)}</span>
                      </div>
                      {tipAmount > 0 && (
                        <div className="flex justify-between text-gray-300">
                          <span>Tip</span>
                          <span>{formatPrice(tipAmount)}</span>
                        </div>
                      )}
                      <div className="flex justify-between text-white font-semibold text-lg pt-2 border-t border-gray-700">
                        <span>Total</span>
                        <span>{formatPrice(finalTotal)}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
