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
      <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-2xl mx-auto text-center">
            {/* Empty cart illustration */}
            <div className="relative mb-8">
              <div className="w-32 h-32 mx-auto rounded-full bg-gradient-to-br from-muted to-muted/50 flex items-center justify-center shadow-large">
                <ShoppingBag className="h-16 w-16 text-muted-foreground" />
              </div>
              <div className="absolute -top-2 -right-2 w-8 h-8 bg-gradient-to-br from-primary to-violet-600 rounded-full flex items-center justify-center shadow-medium">
                <span className="text-white text-sm font-bold">0</span>
              </div>
            </div>

            <h1 className="text-4xl lg:text-5xl font-bold text-foreground mb-6">
              Your cart is <span className="text-gradient">empty</span>
            </h1>
            <p className="text-xl text-muted-foreground mb-12 leading-relaxed">
              Discover amazing products and start building your perfect collection.
            </p>

            {/* Action buttons */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <Button size="xl" variant="gradient" asChild>
                <Link href="/products">
                  <ArrowLeft className="mr-2 h-5 w-5" />
                  Start Shopping
                </Link>
              </Button>
              <Button size="xl" variant="outline" asChild>
                <Link href="/">
                  Browse Categories
                </Link>
              </Button>
            </div>

            {/* Features */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16">
              <div className="text-center">
                <div className="w-12 h-12 rounded-2xl bg-emerald-100 flex items-center justify-center mx-auto mb-3">
                  <span className="text-emerald-600 text-lg">🚚</span>
                </div>
                <h3 className="font-semibold text-foreground mb-1">Free Shipping</h3>
                <p className="text-sm text-muted-foreground">On orders over $50</p>
              </div>
              <div className="text-center">
                <div className="w-12 h-12 rounded-2xl bg-blue-100 flex items-center justify-center mx-auto mb-3">
                  <span className="text-blue-600 text-lg">🔒</span>
                </div>
                <h3 className="font-semibold text-foreground mb-1">Secure Payment</h3>
                <p className="text-sm text-muted-foreground">100% secure checkout</p>
              </div>
              <div className="text-center">
                <div className="w-12 h-12 rounded-2xl bg-purple-100 flex items-center justify-center mx-auto mb-3">
                  <span className="text-purple-600 text-lg">↩️</span>
                </div>
                <h3 className="font-semibold text-foreground mb-1">Easy Returns</h3>
                <p className="text-sm text-muted-foreground">30-day return policy</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-12">
      <div className="container mx-auto px-4">
        {/* Enhanced Header */}
        <div className="mb-12">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 mb-3">
                <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center shadow-lg">
                  <ShoppingBag className="h-4 w-4 text-white" />
                </div>
                <span className="text-orange-500 font-semibold text-sm">SHOPPING CART</span>
              </div>

              <h1 className="text-2xl lg:text-3xl font-bold text-white mb-2">
                Your <span className="text-orange-400">Shopping Cart</span>
              </h1>
              <p className="text-base text-gray-300">
                {cartItemCount} {cartItemCount === 1 ? 'item' : 'items'} ready for checkout
              </p>
            </div>

            <div className="flex flex-col sm:flex-row gap-3">
              <Button
                variant="outline"
                size="default"
                onClick={handleClearCart}
                disabled={isLoading}
                className="border border-red-600 text-red-400 hover:bg-red-600 hover:text-white transition-all duration-200"
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Clear Cart
              </Button>

              <Button variant="ghost" size="default" asChild>
                <Link href="/products" className="flex items-center text-gray-300 hover:text-white">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Continue Shopping
                </Link>
              </Button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Cart Items */}
          <div className="lg:col-span-2 space-y-4">
            {cart.items.map((item) => (
              <Card key={item.id} className="bg-gray-800 border-gray-700 hover:border-gray-600 transition-all duration-300">
                <CardContent className="p-4">
                  <div className="flex flex-col sm:flex-row gap-4">
                    {/* Product Image - More compact */}
                    <div className="relative h-24 w-24 flex-shrink-0 overflow-hidden rounded-lg border border-gray-600">
                      <Image
                        src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                        alt={item.product.name}
                        fill
                        className="object-cover transition-transform duration-300 hover:scale-105"
                      />
                    </div>

                    {/* Product Details - More compact */}
                    <div className="flex-1 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-3">
                        <div className="flex-1">
                          <h3 className="text-lg font-bold text-white mb-1">
                            <Link
                              href={`/products/${item.product.id}`}
                              className="hover:text-orange-400 transition-colors"
                            >
                              {item.product.name}
                            </Link>
                          </h3>
                          <p className="text-xs text-gray-400 mb-2">
                            SKU: {item.product.sku}
                          </p>
                          <div className="flex items-center gap-3">
                            <p className="text-lg font-bold text-orange-400">
                              {formatPrice(item.unit_price)}
                            </p>
                            <p className="text-sm font-semibold text-white">
                              Total: {formatPrice(item.unit_price * item.quantity)}
                            </p>
                          </div>
                        </div>

                        {/* Quantity & Actions */}
                        <div className="flex flex-col gap-4">
                          {/* Quantity Controls */}
                          <div className="flex items-center gap-3">
                            <span className="text-sm font-medium text-muted-foreground">Qty:</span>
                            <div className="flex items-center bg-muted rounded-xl p-1">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                                disabled={isLoading || item.quantity <= 1}
                                className="h-8 w-8 rounded-lg hover:bg-background"
                              >
                                <Minus className="h-4 w-4" />
                              </Button>

                              <span className="text-lg font-bold w-12 text-center">
                                {item.quantity}
                              </span>

                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                                disabled={isLoading}
                                className="h-8 w-8 rounded-lg hover:bg-background"
                              >
                                <Plus className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>

                          {/* Remove Button */}
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleRemoveItem(item.id)}
                            disabled={isLoading}
                            className="border-destructive/30 text-destructive hover:bg-destructive hover:text-destructive-foreground transition-all duration-200"
                          >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Remove
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Enhanced Order Summary */}
          <div className="lg:col-span-1">
            <div className="sticky top-8 space-y-8">
              {/* Coupon Code */}
              <Card variant="elevated" className="border-0 shadow-large">
                <CardHeader className="pb-4">
                  <CardTitle className="flex items-center gap-2">
                    <span className="text-lg">🎫</span>
                    Promo Code
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-3">
                    <Input
                      placeholder="Enter promo code"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                      size="lg"
                      className="flex-1"
                    />
                    <Button
                      onClick={handleApplyCoupon}
                      disabled={!couponCode.trim() || isApplyingCoupon}
                      isLoading={isApplyingCoupon}
                      size="lg"
                      variant="outline"
                    >
                      Apply
                    </Button>
                  </div>
                </CardContent>
              </Card>

              {/* Order Summary */}
              <Card variant="elevated" className="border-0 shadow-large">
                <CardHeader className="pb-6">
                  <CardTitle className="text-xl">Order Summary</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-4">
                    <div className="flex justify-between text-lg">
                      <span className="text-muted-foreground">Subtotal</span>
                      <span className="font-semibold">{formatPrice(cartTotal)}</span>
                    </div>

                    <div className="flex justify-between text-lg">
                      <span className="text-muted-foreground">Shipping</span>
                      <span className="font-semibold">
                        {shippingCost === 0 ? (
                          <span className="text-emerald-600 font-bold">Free</span>
                        ) : (
                          formatPrice(shippingCost)
                        )}
                      </span>
                    </div>

                    <div className="flex justify-between text-lg">
                      <span className="text-muted-foreground">Tax</span>
                      <span className="font-semibold">{formatPrice(tax)}</span>
                    </div>
                  </div>

                  <div className="border-t-2 border-border pt-6">
                    <div className="flex justify-between text-2xl font-bold">
                      <span>Total</span>
                      <span className="text-primary">{formatPrice(finalTotal)}</span>
                    </div>
                  </div>

                  {cartTotal < 50 && (
                    <div className="bg-gradient-to-r from-blue-50 to-indigo-50 border-2 border-blue-200 rounded-2xl p-4">
                      <div className="flex items-center gap-3">
                        <span className="text-2xl">🚚</span>
                        <div>
                          <p className="font-semibold text-blue-900">Almost there!</p>
                          <p className="text-sm text-blue-700">
                            Add {formatPrice(50 - cartTotal)} more for free shipping
                          </p>
                        </div>
                      </div>
                    </div>
                  )}

                  <div className="space-y-4">
                    <Button size="xl" variant="gradient" className="w-full" asChild>
                      <Link href="/checkout">
                        <span className="text-lg font-semibold">Proceed to Checkout</span>
                      </Link>
                    </Button>

                    <div className="text-center">
                      <Link
                        href="/products"
                        className="text-sm text-primary hover:text-primary-600 transition-colors font-medium"
                      >
                        ← Continue Shopping
                      </Link>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Trust Signals */}
              <Card variant="elevated" className="border-0 shadow-medium">
                <CardContent className="p-6">
                  <div className="space-y-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-emerald-100 rounded-full flex items-center justify-center">
                        <span className="text-emerald-600 text-sm">🔒</span>
                      </div>
                      <span className="text-sm font-medium text-foreground">Secure SSL encryption</span>
                    </div>

                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                        <span className="text-blue-600 text-sm">💳</span>
                      </div>
                      <span className="text-sm font-medium text-foreground">Multiple payment options</span>
                    </div>

                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                        <span className="text-purple-600 text-sm">↩️</span>
                      </div>
                      <span className="text-sm font-medium text-foreground">30-day return policy</span>
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
