'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useCartStore } from '@/store/cart'
import { useOrderStore } from '@/store/order'
import { usePaymentStore } from '@/store/payment'
import { redirectToCheckout } from '@/lib/stripe'
import { Loader2, CreditCard, MapPin } from 'lucide-react'
import toast from 'react-hot-toast'

const addressSchema = z.object({
  street: z.string().min(1, 'Street address is required'),
  city: z.string().min(1, 'City is required'),
  state: z.string().min(1, 'State is required'),
  zip_code: z.string().min(1, 'ZIP code is required'),
  country: z.string().min(1, 'Country is required'),
})

const checkoutSchema = z.object({
  shipping_address: addressSchema,
  billing_address: addressSchema,
  same_as_shipping: z.boolean().default(false),
  coupon_code: z.string().optional(),
  notes: z.string().optional(),
})

type CheckoutFormData = z.infer<typeof checkoutSchema>

export default function CheckoutForm() {
  const router = useRouter()
  const [isProcessing, setIsProcessing] = useState(false)
  const [sameAsShipping, setSameAsShipping] = useState(false)
  
  const { cart } = useCartStore()
  const { createOrder } = useOrderStore()
  const { createCheckoutSession } = usePaymentStore()

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<CheckoutFormData>({
    resolver: zodResolver(checkoutSchema),
    defaultValues: {
      same_as_shipping: false,
      shipping_address: {
        country: 'US',
      },
      billing_address: {
        country: 'US',
      },
    },
  })

  const shippingAddress = watch('shipping_address')

  // Update billing address when "same as shipping" is checked
  const handleSameAsShipping = (checked: boolean) => {
    setSameAsShipping(checked)
    setValue('same_as_shipping', checked)
    
    if (checked) {
      setValue('billing_address', shippingAddress)
    }
  }

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
        billing_address: data.same_as_shipping ? data.shipping_address : data.billing_address,
        coupon_code: data.coupon_code,
        notes: data.notes,
      }

      const orderResponse = await createOrder(orderData)
      const orderId = orderResponse.data.id

      // Create Stripe checkout session
      const checkoutData = {
        order_id: orderId,
        amount: cart.total,
        currency: 'usd',
        description: `Order ${orderResponse.data.order_number}`,
        success_url: `${window.location.origin}/checkout/success?session_id={CHECKOUT_SESSION_ID}`,
        cancel_url: `${window.location.origin}/checkout/cancel`,
        metadata: {
          order_id: orderId,
          order_number: orderResponse.data.order_number,
        },
      }

      const session = await createCheckoutSession(checkoutData)
      
      // Redirect to Stripe Checkout
      await redirectToCheckout(session.session_id)
      
    } catch (error: any) {
      console.error('Checkout error:', error)
      toast.error(error.message || 'Failed to process checkout')
    } finally {
      setIsProcessing(false)
    }
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500 mb-4">Your cart is empty</p>
        <Button onClick={() => router.push('/products')}>
          Continue Shopping
        </Button>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-3xl font-bold mb-8">Checkout</h1>
      
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* Checkout Form */}
        <div className="space-y-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Shipping Address */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Shipping Address
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="shipping_street">Street Address</Label>
                  <Input
                    id="shipping_street"
                    {...register('shipping_address.street')}
                    placeholder="123 Main St"
                  />
                  {errors.shipping_address?.street && (
                    <p className="text-sm text-red-500 mt-1">
                      {errors.shipping_address.street.message}
                    </p>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="shipping_city">City</Label>
                    <Input
                      id="shipping_city"
                      {...register('shipping_address.city')}
                      placeholder="New York"
                    />
                    {errors.shipping_address?.city && (
                      <p className="text-sm text-red-500 mt-1">
                        {errors.shipping_address.city.message}
                      </p>
                    )}
                  </div>
                  <div>
                    <Label htmlFor="shipping_state">State</Label>
                    <Input
                      id="shipping_state"
                      {...register('shipping_address.state')}
                      placeholder="NY"
                    />
                    {errors.shipping_address?.state && (
                      <p className="text-sm text-red-500 mt-1">
                        {errors.shipping_address.state.message}
                      </p>
                    )}
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="shipping_zip">ZIP Code</Label>
                    <Input
                      id="shipping_zip"
                      {...register('shipping_address.zip_code')}
                      placeholder="10001"
                    />
                    {errors.shipping_address?.zip_code && (
                      <p className="text-sm text-red-500 mt-1">
                        {errors.shipping_address.zip_code.message}
                      </p>
                    )}
                  </div>
                  <div>
                    <Label htmlFor="shipping_country">Country</Label>
                    <Input
                      id="shipping_country"
                      {...register('shipping_address.country')}
                      placeholder="US"
                    />
                    {errors.shipping_address?.country && (
                      <p className="text-sm text-red-500 mt-1">
                        {errors.shipping_address.country.message}
                      </p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Billing Address */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <CreditCard className="h-5 w-5" />
                  Billing Address
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="same_as_shipping"
                    checked={sameAsShipping}
                    onChange={(e) => handleSameAsShipping(e.target.checked)}
                    className="rounded border-gray-300"
                  />
                  <Label htmlFor="same_as_shipping">
                    Same as shipping address
                  </Label>
                </div>

                {!sameAsShipping && (
                  <>
                    <div>
                      <Label htmlFor="billing_street">Street Address</Label>
                      <Input
                        id="billing_street"
                        {...register('billing_address.street')}
                        placeholder="123 Main St"
                      />
                      {errors.billing_address?.street && (
                        <p className="text-sm text-red-500 mt-1">
                          {errors.billing_address.street.message}
                        </p>
                      )}
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="billing_city">City</Label>
                        <Input
                          id="billing_city"
                          {...register('billing_address.city')}
                          placeholder="New York"
                        />
                        {errors.billing_address?.city && (
                          <p className="text-sm text-red-500 mt-1">
                            {errors.billing_address.city.message}
                          </p>
                        )}
                      </div>
                      <div>
                        <Label htmlFor="billing_state">State</Label>
                        <Input
                          id="billing_state"
                          {...register('billing_address.state')}
                          placeholder="NY"
                        />
                        {errors.billing_address?.state && (
                          <p className="text-sm text-red-500 mt-1">
                            {errors.billing_address.state.message}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="grid grid-cols-2 gap-4">
                      <div>
                        <Label htmlFor="billing_zip">ZIP Code</Label>
                        <Input
                          id="billing_zip"
                          {...register('billing_address.zip_code')}
                          placeholder="10001"
                        />
                        {errors.billing_address?.zip_code && (
                          <p className="text-sm text-red-500 mt-1">
                            {errors.billing_address.zip_code.message}
                          </p>
                        )}
                      </div>
                      <div>
                        <Label htmlFor="billing_country">Country</Label>
                        <Input
                          id="billing_country"
                          {...register('billing_address.country')}
                          placeholder="US"
                        />
                        {errors.billing_address?.country && (
                          <p className="text-sm text-red-500 mt-1">
                            {errors.billing_address.country.message}
                          </p>
                        )}
                      </div>
                    </div>
                  </>
                )}
              </CardContent>
            </Card>

            {/* Additional Options */}
            <Card>
              <CardContent className="pt-6 space-y-4">
                <div>
                  <Label htmlFor="coupon_code">Coupon Code (Optional)</Label>
                  <Input
                    id="coupon_code"
                    {...register('coupon_code')}
                    placeholder="Enter coupon code"
                  />
                </div>

                <div>
                  <Label htmlFor="notes">Order Notes (Optional)</Label>
                  <textarea
                    id="notes"
                    {...register('notes')}
                    placeholder="Special instructions for your order"
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    rows={3}
                  />
                </div>
              </CardContent>
            </Card>

            <Button
              type="submit"
              className="w-full"
              size="lg"
              disabled={isProcessing}
            >
              {isProcessing ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Processing...
                </>
              ) : (
                <>
                  <CreditCard className="mr-2 h-4 w-4" />
                  Proceed to Payment
                </>
              )}
            </Button>
          </form>
        </div>

        {/* Order Summary */}
        <div>
          <Card className="sticky top-6">
            <CardHeader>
              <CardTitle>Order Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {cart.items.map((item) => (
                <div key={item.id} className="flex justify-between items-center">
                  <div className="flex-1">
                    <p className="font-medium">{item.product.name}</p>
                    <p className="text-sm text-gray-500">Qty: {item.quantity}</p>
                  </div>
                  <p className="font-medium">${item.subtotal.toFixed(2)}</p>
                </div>
              ))}
              
              <Separator />
              
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Subtotal</span>
                  <span>${cart.total.toFixed(2)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Shipping</span>
                  <span>Free</span>
                </div>
                <div className="flex justify-between">
                  <span>Tax</span>
                  <span>$0.00</span>
                </div>
                <Separator />
                <div className="flex justify-between text-lg font-bold">
                  <span>Total</span>
                  <span>${cart.total.toFixed(2)}</span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
