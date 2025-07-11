'use client'

import { UseFormRegister, FieldErrors, UseFormWatch, UseFormSetValue } from 'react-hook-form'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  User, 
  Mail, 
  Phone, 
  MapPin, 
  Building, 
  Truck, 
  Gift, 
  Clock,
  DollarSign,
  Star,
  Package
} from 'lucide-react'
import { cn, formatPrice } from '@/lib/utils'

interface CheckoutFormData {
  email: string
  phone: string
  shipping_address: {
    first_name: string
    last_name: string
    company?: string
    address1: string
    address2?: string
    city: string
    state: string
    zip_code: string
    country: string
    phone?: string
  }
  use_shipping_for_billing: boolean
  billing_address?: any
  shipping_method: 'standard' | 'express' | 'overnight'
  is_gift: boolean
  gift_message?: string
  gift_wrap: boolean
  delivery_instructions?: string
  priority: 'normal' | 'high' | 'urgent'
  tip_amount: number
  coupon_code?: string
  notes?: string
}

interface StepProps {
  register: UseFormRegister<CheckoutFormData>
  errors: FieldErrors<CheckoutFormData>
  watch: UseFormWatch<CheckoutFormData>
  setValue: UseFormSetValue<CheckoutFormData>
}

// Step 1: Contact Information
export function ContactInfoStep({ register, errors }: StepProps) {
  return (
    <Card className="bg-gray-900/50 border-gray-800">
      <CardHeader>
        <CardTitle className="text-white flex items-center gap-2">
          <User className="h-5 w-5 text-[#ff9000]" />
          Contact Information
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="email" className="text-gray-300">
              Email Address *
            </Label>
            <div className="relative">
              <Mail className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
              <Input
                id="email"
                type="email"
                {...register('email')}
                className="pl-10 bg-gray-800 border-gray-700 text-white"
                placeholder="your@email.com"
              />
            </div>
            {errors.email && (
              <p className="text-red-400 text-sm">{errors.email.message}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="phone" className="text-gray-300">
              Phone Number *
            </Label>
            <div className="relative">
              <Phone className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
              <Input
                id="phone"
                type="tel"
                {...register('phone')}
                className="pl-10 bg-gray-800 border-gray-700 text-white"
                placeholder="+1 (555) 123-4567"
              />
            </div>
            {errors.phone && (
              <p className="text-red-400 text-sm">{errors.phone.message}</p>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

// Step 2: Shipping Address
export function ShippingAddressStep({ register, errors }: StepProps) {
  return (
    <Card className="bg-gray-900/50 border-gray-800">
      <CardHeader>
        <CardTitle className="text-white flex items-center gap-2">
          <MapPin className="h-5 w-5 text-[#ff9000]" />
          Shipping Address
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <Label htmlFor="shipping_first_name" className="text-gray-300">
              First Name *
            </Label>
            <Input
              id="shipping_first_name"
              {...register('shipping_address.first_name')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="John"
            />
            {errors.shipping_address?.first_name && (
              <p className="text-red-400 text-sm">{errors.shipping_address.first_name.message}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="shipping_last_name" className="text-gray-300">
              Last Name *
            </Label>
            <Input
              id="shipping_last_name"
              {...register('shipping_address.last_name')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="Doe"
            />
            {errors.shipping_address?.last_name && (
              <p className="text-red-400 text-sm">{errors.shipping_address.last_name.message}</p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="shipping_company" className="text-gray-300">
            Company (Optional)
          </Label>
          <div className="relative">
            <Building className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
            <Input
              id="shipping_company"
              {...register('shipping_address.company')}
              className="pl-10 bg-gray-800 border-gray-700 text-white"
              placeholder="Your Company"
            />
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="shipping_address_1" className="text-gray-300">
            Street Address *
          </Label>
          <Input
            id="shipping_address_1"
            {...register('shipping_address.address1')}
            className="bg-gray-800 border-gray-700 text-white"
            placeholder="123 Main Street"
          />
          {errors.shipping_address?.address1 && (
            <p className="text-red-400 text-sm">{errors.shipping_address.address1.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label htmlFor="shipping_address_2" className="text-gray-300">
            Apartment, Suite, etc. (Optional)
          </Label>
          <Input
            id="shipping_address_2"
            {...register('shipping_address.address2')}
            className="bg-gray-800 border-gray-700 text-white"
            placeholder="Apt 4B"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <Label htmlFor="shipping_city" className="text-gray-300">
              City *
            </Label>
            <Input
              id="shipping_city"
              {...register('shipping_address.city')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="New York"
            />
            {errors.shipping_address?.city && (
              <p className="text-red-400 text-sm">{errors.shipping_address.city.message}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="shipping_state" className="text-gray-300">
              State *
            </Label>
            <Input
              id="shipping_state"
              {...register('shipping_address.state')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="NY"
            />
            {errors.shipping_address?.state && (
              <p className="text-red-400 text-sm">{errors.shipping_address.state.message}</p>
            )}
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="shipping_zip_code" className="text-gray-300">
              ZIP Code *
            </Label>
            <Input
              id="shipping_zip_code"
              {...register('shipping_address.zip_code')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="10001"
            />
            {errors.shipping_address?.zip_code && (
              <p className="text-red-400 text-sm">{errors.shipping_address.zip_code.message}</p>
            )}
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="shipping_country" className="text-gray-300">
            Country *
          </Label>
          <Input
            id="shipping_country"
            {...register('shipping_address.country')}
            className="bg-gray-800 border-gray-700 text-white"
            placeholder="United States"
            defaultValue="US"
          />
          {errors.shipping_address?.country && (
            <p className="text-red-400 text-sm">{errors.shipping_address.country.message}</p>
          )}
        </div>
      </CardContent>
    </Card>
  )
}

// Step 3: Billing Address
export function BillingAddressStep({ register, errors, watch, setValue }: StepProps) {
  const useShippingForBilling = watch('use_shipping_for_billing')

  return (
    <Card className="bg-gray-900/50 border-gray-800">
      <CardHeader>
        <CardTitle className="text-white flex items-center gap-2">
          <DollarSign className="h-5 w-5 text-[#ff9000]" />
          Billing Address
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex items-center space-x-2">
          <input
            type="checkbox"
            id="use_shipping_for_billing"
            {...register('use_shipping_for_billing')}
            className="rounded border-gray-700 bg-gray-800 text-[#ff9000] focus:ring-[#ff9000]"
          />
          <Label htmlFor="use_shipping_for_billing" className="text-gray-300">
            Use shipping address for billing
          </Label>
        </div>

        {!useShippingForBilling && (
          <div className="space-y-4 pt-4 border-t border-gray-700">
            <p className="text-gray-400 text-sm">
              Enter a different billing address below:
            </p>
            {/* Billing address fields would go here - similar to shipping */}
            <div className="text-center py-8 text-gray-400">
              <p>Billing address form would be implemented here</p>
              <p className="text-sm">For now, shipping address will be used</p>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}

// Step 4: Shipping & Gift Options
export function ShippingOptionsStep({ register, errors, watch, setValue }: StepProps) {
  const shippingMethod = watch('shipping_method')
  const isGift = watch('is_gift')
  const giftWrap = watch('gift_wrap')

  const shippingOptions = [
    {
      id: 'standard',
      name: 'Standard Shipping',
      description: 'Free on orders over $50',
      price: 0, // Will be calculated based on cart total
      estimatedDays: '5-7 business days',
      icon: Package
    },
    {
      id: 'express',
      name: 'Express Shipping',
      description: 'Faster delivery',
      price: 9.99,
      estimatedDays: '2-3 business days',
      icon: Truck
    },
    {
      id: 'overnight',
      name: 'Overnight Shipping',
      description: 'Next business day',
      price: 19.99,
      estimatedDays: '1 business day',
      icon: Clock
    }
  ]

  return (
    <div className="space-y-6">
      {/* Shipping Options */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <Truck className="h-5 w-5 text-[#ff9000]" />
            Shipping Options
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {shippingOptions.map((option) => {
            const IconComponent = option.icon
            return (
              <div
                key={option.id}
                className={cn(
                  "border rounded-lg p-4 cursor-pointer transition-all duration-200",
                  shippingMethod === option.id
                    ? "border-[#ff9000] bg-[#ff9000]/10"
                    : "border-gray-700 hover:border-gray-600"
                )}
                onClick={() => setValue('shipping_method', option.id as any)}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-3">
                    <input
                      type="radio"
                      {...register('shipping_method')}
                      value={option.id}
                      className="text-[#ff9000] focus:ring-[#ff9000]"
                    />
                    <IconComponent className="h-5 w-5 text-gray-400" />
                    <div>
                      <div className="text-white font-medium">{option.name}</div>
                      <div className="text-gray-400 text-sm">{option.description}</div>
                      <div className="text-gray-300 text-sm">{option.estimatedDays}</div>
                    </div>
                  </div>
                  <div className="text-white font-semibold">
                    {option.price === 0 ? 'Free' : formatPrice(option.price)}
                  </div>
                </div>
              </div>
            )
          })}
        </CardContent>
      </Card>

      {/* Gift Options */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <Gift className="h-5 w-5 text-[#ff9000]" />
            Gift Options
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="is_gift"
              {...register('is_gift')}
              className="rounded border-gray-700 bg-gray-800 text-[#ff9000] focus:ring-[#ff9000]"
            />
            <Label htmlFor="is_gift" className="text-gray-300">
              This is a gift
            </Label>
          </div>

          {isGift && (
            <div className="space-y-4 pt-4 border-t border-gray-700">
              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="gift_wrap"
                  {...register('gift_wrap')}
                  className="rounded border-gray-700 bg-gray-800 text-[#ff9000] focus:ring-[#ff9000]"
                />
                <Label htmlFor="gift_wrap" className="text-gray-300">
                  Add gift wrapping (+$4.99)
                </Label>
              </div>

              <div className="space-y-2">
                <Label htmlFor="gift_message" className="text-gray-300">
                  Gift Message (Optional)
                </Label>
                <Textarea
                  id="gift_message"
                  {...register('gift_message')}
                  className="bg-gray-800 border-gray-700 text-white"
                  placeholder="Write a personal message..."
                  rows={3}
                />
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delivery Instructions */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <MapPin className="h-5 w-5 text-[#ff9000]" />
            Delivery Instructions
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="delivery_instructions" className="text-gray-300">
              Special Instructions (Optional)
            </Label>
            <Textarea
              id="delivery_instructions"
              {...register('delivery_instructions')}
              className="bg-gray-800 border-gray-700 text-white"
              placeholder="Leave at front door, ring doorbell, etc..."
              rows={3}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="priority" className="text-gray-300">
              Order Priority
            </Label>
            <Select defaultValue="normal">
              <SelectTrigger className="bg-gray-800 border-gray-700 text-white">
                <SelectValue placeholder="Select priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="normal">Normal</SelectItem>
                <SelectItem value="high">High (+$5.00)</SelectItem>
                <SelectItem value="urgent">Urgent (+$15.00)</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="tip_amount" className="text-gray-300">
              Tip for Delivery (Optional)
            </Label>
            <div className="relative">
              <DollarSign className="absolute left-3 top-3 h-4 w-4 text-gray-400" />
              <Input
                id="tip_amount"
                type="number"
                step="0.01"
                min="0"
                {...register('tip_amount', { valueAsNumber: true })}
                className="pl-10 bg-gray-800 border-gray-700 text-white"
                placeholder="0.00"
              />
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

// Step 5: Review & Payment
export function ReviewStep({
  watch,
  cartTotal,
  shippingCost,
  tax,
  finalTotal,
  cartItems
}: StepProps & {
  cartTotal: number
  shippingCost: number
  tax: number
  finalTotal: number
  cartItems: any[]
}) {
  const formData = watch()

  return (
    <div className="space-y-6">
      {/* Order Summary */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <Package className="h-5 w-5 text-[#ff9000]" />
            Order Summary
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Cart Items */}
          <div className="space-y-3">
            {cartItems.map((item: any) => (
              <div key={item.id} className="flex items-center justify-between py-2 border-b border-gray-700 last:border-b-0">
                <div className="flex items-center space-x-3">
                  <div className="w-12 h-12 bg-gray-800 rounded-lg flex items-center justify-center">
                    <Package className="h-6 w-6 text-gray-400" />
                  </div>
                  <div>
                    <div className="text-white font-medium">{item.product?.name}</div>
                    <div className="text-gray-400 text-sm">Qty: {item.quantity}</div>
                  </div>
                </div>
                <div className="text-white font-semibold">
                  {formatPrice(item.product?.price * item.quantity)}
                </div>
              </div>
            ))}
          </div>

          <Separator className="bg-gray-700" />

          {/* Cost Breakdown */}
          <div className="space-y-2">
            <div className="flex justify-between text-gray-300">
              <span>Subtotal</span>
              <span>{formatPrice(cartTotal)}</span>
            </div>
            <div className="flex justify-between text-gray-300">
              <span>Shipping ({formData.shipping_method})</span>
              <span>{shippingCost === 0 ? 'Free' : formatPrice(shippingCost)}</span>
            </div>
            <div className="flex justify-between text-gray-300">
              <span>Tax</span>
              <span>{formatPrice(tax)}</span>
            </div>
            {formData.tip_amount > 0 && (
              <div className="flex justify-between text-gray-300">
                <span>Tip</span>
                <span>{formatPrice(formData.tip_amount)}</span>
              </div>
            )}
            {formData.gift_wrap && (
              <div className="flex justify-between text-gray-300">
                <span>Gift Wrapping</span>
                <span>{formatPrice(4.99)}</span>
              </div>
            )}
            <Separator className="bg-gray-700" />
            <div className="flex justify-between text-white font-semibold text-lg">
              <span>Total</span>
              <span>{formatPrice(finalTotal)}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Shipping Information */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <MapPin className="h-5 w-5 text-[#ff9000]" />
            Shipping Information
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div>
            <div className="text-gray-400 text-sm">Contact</div>
            <div className="text-white">{formData.email}</div>
            <div className="text-white">{formData.phone}</div>
          </div>

          <div>
            <div className="text-gray-400 text-sm">Ship to</div>
            <div className="text-white">
              {formData.shipping_address.first_name} {formData.shipping_address.last_name}
            </div>
            <div className="text-white">{formData.shipping_address.address1}</div>
            {formData.shipping_address.address2 && (
              <div className="text-white">{formData.shipping_address.address2}</div>
            )}
            <div className="text-white">
              {formData.shipping_address.city}, {formData.shipping_address.state} {formData.shipping_address.zip_code}
            </div>
            <div className="text-white">{formData.shipping_address.country}</div>
          </div>

          <div>
            <div className="text-gray-400 text-sm">Shipping Method</div>
            <div className="text-white capitalize">{formData.shipping_method} shipping</div>
          </div>

          {formData.delivery_instructions && (
            <div>
              <div className="text-gray-400 text-sm">Delivery Instructions</div>
              <div className="text-white">{formData.delivery_instructions}</div>
            </div>
          )}

          {formData.is_gift && (
            <div>
              <div className="text-gray-400 text-sm">Gift Options</div>
              <div className="text-white">
                This is a gift
                {formData.gift_wrap && ' with gift wrapping'}
              </div>
              {formData.gift_message && (
                <div className="text-white italic">"{formData.gift_message}"</div>
              )}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
