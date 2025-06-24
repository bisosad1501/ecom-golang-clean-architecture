'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Minus, Plus, Trash2, ShoppingBag, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useCartStore, getCartTotal, getCartItemCount } from '@/store/cart'
import { formatPrice } from '@/lib/utils'
import { toast } from 'sonner'

export function CartPage() {
  const { 
    cart, 
    isLoading, 
    updateItem, 
    removeItem, 
    clearCart 
  } = useCartStore()
  
  const [couponCode, setCouponCode] = useState('')
  const [isApplyingCoupon, setIsApplyingCoupon] = useState(false)

  const cartTotal = getCartTotal(cart)
  const cartItemCount = getCartItemCount(cart)
  const shippingCost = cartTotal > 50 ? 0 : 9.99
  const tax = cartTotal * 0.08 // 8% tax
  const finalTotal = cartTotal + shippingCost + tax

  const handleUpdateQuantity = async (itemId: string, newQuantity: number) => {
    if (newQuantity < 1) return
    
    try {
      await updateItem(itemId, newQuantity)
    } catch (error) {
      toast.error('Failed to update item quantity')
    }
  }

  const handleRemoveItem = async (itemId: string) => {
    try {
      await removeItem(itemId)
      toast.success('Item removed from cart')
    } catch (error) {
      toast.error('Failed to remove item')
    }
  }

  const handleClearCart = async () => {
    if (window.confirm('Are you sure you want to clear your cart?')) {
      try {
        await clearCart()
        toast.success('Cart cleared')
      } catch (error) {
        toast.error('Failed to clear cart')
      }
    }
  }

  const handleApplyCoupon = async () => {
    if (!couponCode.trim()) return
    
    setIsApplyingCoupon(true)
    try {
      // TODO: Implement coupon application
      await new Promise(resolve => setTimeout(resolve, 1000))
      toast.success('Coupon applied successfully!')
    } catch (error) {
      toast.error('Invalid coupon code')
    } finally {
      setIsApplyingCoupon(false)
    }
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gray-50 py-12">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto text-center">
            <ShoppingBag className="h-24 w-24 text-gray-300 mx-auto mb-6" />
            <h1 className="text-3xl font-bold text-gray-900 mb-4">
              Your cart is empty
            </h1>
            <p className="text-gray-600 mb-8">
              Looks like you haven't added any items to your cart yet.
            </p>
            <Button size="lg" asChild>
              <Link href="/products">
                <ArrowLeft className="mr-2 h-5 w-5" />
                Continue Shopping
              </Link>
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <div className="flex items-center justify-between">
            <h1 className="text-3xl font-bold text-gray-900">
              Shopping Cart ({cartItemCount} items)
            </h1>
            <Button
              variant="outline"
              onClick={handleClearCart}
              disabled={isLoading}
            >
              Clear Cart
            </Button>
          </div>
          <Link 
            href="/products" 
            className="inline-flex items-center text-primary-600 hover:text-primary-700 mt-2"
          >
            <ArrowLeft className="mr-1 h-4 w-4" />
            Continue Shopping
          </Link>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle>Cart Items</CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y">
                  {cart.items.map((item) => (
                    <div key={item.id} className="p-6">
                      <div className="flex items-center space-x-4">
                        {/* Product Image */}
                        <div className="relative h-20 w-20 flex-shrink-0 overflow-hidden rounded-md border">
                          <Image
                            src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                            alt={item.product.name}
                            fill
                            className="object-cover"
                          />
                        </div>

                        {/* Product Details */}
                        <div className="flex-1 min-w-0">
                          <h3 className="text-lg font-medium text-gray-900">
                            <Link 
                              href={`/products/${item.product.id}`}
                              className="hover:text-primary-600"
                            >
                              {item.product.name}
                            </Link>
                          </h3>
                          <p className="text-sm text-gray-500 mt-1">
                            SKU: {item.product.sku}
                          </p>
                          <p className="text-lg font-semibold text-gray-900 mt-2">
                            {formatPrice(item.unit_price)}
                          </p>
                        </div>

                        {/* Quantity Controls */}
                        <div className="flex items-center space-x-3">
                          <Button
                            variant="outline"
                            size="icon"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                            disabled={isLoading || item.quantity <= 1}
                          >
                            <Minus className="h-4 w-4" />
                          </Button>
                          
                          <span className="text-lg font-medium w-12 text-center">
                            {item.quantity}
                          </span>
                          
                          <Button
                            variant="outline"
                            size="icon"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                            disabled={isLoading}
                          >
                            <Plus className="h-4 w-4" />
                          </Button>
                        </div>

                        {/* Total Price */}
                        <div className="text-right">
                          <p className="text-lg font-semibold text-gray-900">
                            {formatPrice(item.total_price)}
                          </p>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => handleRemoveItem(item.id)}
                            disabled={isLoading}
                            className="text-red-600 hover:text-red-700 mt-2"
                          >
                            <Trash2 className="h-4 w-4 mr-1" />
                            Remove
                          </Button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Order Summary */}
          <div className="lg:col-span-1">
            <div className="sticky top-8 space-y-6">
              {/* Coupon Code */}
              <Card>
                <CardHeader>
                  <CardTitle>Coupon Code</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex space-x-2">
                    <Input
                      placeholder="Enter coupon code"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                    />
                    <Button
                      onClick={handleApplyCoupon}
                      disabled={!couponCode.trim() || isApplyingCoupon}
                      isLoading={isApplyingCoupon}
                    >
                      Apply
                    </Button>
                  </div>
                </CardContent>
              </Card>

              {/* Order Summary */}
              <Card>
                <CardHeader>
                  <CardTitle>Order Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex justify-between">
                    <span>Subtotal</span>
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

                  <Button size="lg" className="w-full" asChild>
                    <Link href="/checkout">
                      Proceed to Checkout
                    </Link>
                  </Button>

                  <div className="text-center">
                    <Link 
                      href="/products" 
                      className="text-sm text-primary-600 hover:text-primary-700"
                    >
                      Continue Shopping
                    </Link>
                  </div>
                </CardContent>
              </Card>

              {/* Security Notice */}
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center space-x-2 text-sm text-gray-600">
                    <div className="w-4 h-4 bg-green-500 rounded-full flex items-center justify-center">
                      <span className="text-white text-xs">✓</span>
                    </div>
                    <span>Secure checkout with SSL encryption</span>
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
