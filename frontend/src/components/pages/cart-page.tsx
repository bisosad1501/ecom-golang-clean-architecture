'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Minus, Plus, Trash2, ShoppingBag, ArrowLeft, Heart, Star, Shield, Truck, CreditCard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
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
      <div className="min-h-screen hero-gradient py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-3xl mx-auto text-center">
            {/* BiHub Empty Cart Hero */}
            <div className="relative mb-12 animate-fade-in">
              <div className="w-40 h-40 mx-auto rounded-full bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center shadow-2xl border border-gray-700">
                <ShoppingBag className="h-20 w-20 text-gray-400" />
              </div>
              <div className="absolute -top-3 -right-3 w-12 h-12 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-xl">
                <span className="text-white text-lg font-bold">0</span>
              </div>
            </div>

            <h1 className="text-5xl lg:text-6xl font-bold text-white mb-6 animate-slide-up">
              Your <span className="text-gradient">BiHub</span> cart is empty
            </h1>
            <p className="text-xl text-gray-300 mb-12 leading-relaxed animate-slide-up">
              Discover amazing products and start building your perfect BiHub collection.
            </p>

            {/* Action buttons */}
            <div className="flex flex-col sm:flex-row gap-6 justify-center mb-16 animate-scale-in">
              <Button size="lg" className="btn-gradient px-8 py-4 text-lg" asChild>
                <Link href="/products">
                  <ShoppingBag className="mr-3 h-6 w-6" />
                  Start Shopping
                </Link>
              </Button>
              <Button size="lg" variant="outline" className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white px-8 py-4 text-lg" asChild>
                <Link href="/">
                  Browse Categories
                </Link>
              </Button>
            </div>

            {/* BiHub Features */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 animate-fade-in">
              <div className="glass-effect p-6 rounded-2xl text-center">
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center mx-auto mb-4">
                  <Truck className="h-8 w-8 text-white" />
                </div>
                <h3 className="font-bold text-white mb-2 text-lg">Free Shipping</h3>
                <p className="text-gray-300">On orders over $50</p>
              </div>
              <div className="glass-effect p-6 rounded-2xl text-center">
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center mx-auto mb-4">
                  <Shield className="h-8 w-8 text-white" />
                </div>
                <h3 className="font-bold text-white mb-2 text-lg">Secure Payment</h3>
                <p className="text-gray-300">100% secure checkout</p>
              </div>
              <div className="glass-effect p-6 rounded-2xl text-center">
                <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center mx-auto mb-4">
                  <ArrowLeft className="h-8 w-8 text-white" />
                </div>
                <h3 className="font-bold text-white mb-2 text-lg">Easy Returns</h3>
                <p className="text-gray-300">30-day return policy</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen hero-gradient py-12">
      <div className="container mx-auto px-4">
        {/* BiHub Cart Header */}
        <div className="mb-12 animate-fade-in">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-8">
            <div>
              <div className="flex items-center gap-3 mb-4">
                <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center shadow-xl">
                  <ShoppingBag className="h-6 w-6 text-white" />
                </div>
                <div>
                  <span className="text-orange-400 font-bold text-sm tracking-wider">BIHUB SHOPPING CART</span>
                  <div className="flex items-center gap-2 mt-1">
                    <Badge className="bg-orange-500/20 text-orange-400 border-orange-500/30">
                      {cartItemCount} {cartItemCount === 1 ? 'item' : 'items'}
                    </Badge>
                  </div>
                </div>
              </div>

              <h1 className="text-3xl lg:text-4xl font-bold text-white mb-3">
                Your <span className="text-gradient">BiHub</span> Shopping Cart
              </h1>
              <p className="text-lg text-gray-300">
                Ready for checkout • Total: <span className="text-orange-400 font-bold">{formatPrice(cartTotal)}</span>
              </p>
            </div>

            <div className="flex flex-col sm:flex-row gap-4">
              <Button
                variant="outline"
                size="lg"
                onClick={handleClearCart}
                disabled={isLoading}
                className="border-red-500/50 text-red-400 hover:bg-red-500 hover:text-white hover:border-red-500 transition-all duration-300"
              >
                <Trash2 className="mr-2 h-5 w-5" />
                Clear Cart
              </Button>

              <Button variant="outline" size="lg" className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white" asChild>
                <Link href="/products" className="flex items-center">
                  <ArrowLeft className="mr-2 h-5 w-5" />
                  Continue Shopping
                </Link>
              </Button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 animate-slide-up">
          {/* BiHub Cart Items */}
          <div className="lg:col-span-2 space-y-6">
            {cart.items.map((item, index) => (
              <Card key={item.id} className="bg-gray-900/50 border-gray-700 hover:border-orange-500/50 transition-all duration-300 hover:shadow-2xl hover:shadow-orange-500/10 card-hover">
                <CardContent className="p-6">
                  <div className="flex flex-col sm:flex-row gap-6">
                    {/* Product Image */}
                    <div className="relative h-32 w-32 flex-shrink-0 overflow-hidden rounded-xl border border-gray-600 group">
                      <Image
                        src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                        alt={item.product.name}
                        fill
                        className="object-cover transition-transform duration-500 group-hover:scale-110"
                      />
                      <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
                    </div>

                    {/* Product Details */}
                    <div className="flex-1 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <Badge className="bg-orange-500/20 text-orange-400 border-orange-500/30 text-xs">
                              #{index + 1}
                            </Badge>
                            <div className="flex items-center gap-1">
                              {[...Array(5)].map((_, i) => (
                                <Star key={i} className="h-3 w-3 fill-orange-400 text-orange-400" />
                              ))}
                            </div>
                          </div>

                          <h3 className="text-xl font-bold text-white mb-2 group">
                            <Link
                              href={`/products/${item.product.id}`}
                              className="hover:text-orange-400 transition-colors duration-300"
                            >
                              {item.product.name}
                            </Link>
                          </h3>

                          <p className="text-sm text-gray-400 mb-3">
                            SKU: <span className="font-mono">{item.product.sku}</span>
                          </p>

                          <div className="flex items-center gap-4 mb-4">
                            <div className="text-2xl font-bold text-orange-400">
                              {formatPrice(item.unit_price)}
                            </div>
                            <div className="text-lg font-semibold text-white">
                              Total: <span className="text-orange-300">{formatPrice(item.unit_price * item.quantity)}</span>
                            </div>
                          </div>
                        </div>

                        {/* Quantity & Actions */}
                        <div className="flex flex-col gap-6">
                          {/* Quantity Controls */}
                          <div className="flex flex-col gap-3">
                            <span className="text-sm font-medium text-gray-400">Quantity:</span>
                            <div className="flex items-center bg-gray-800 rounded-xl p-2 border border-gray-600">
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                                disabled={isLoading || item.quantity <= 1}
                                className="h-10 w-10 rounded-lg hover:bg-gray-700 text-gray-300 hover:text-white"
                              >
                                <Minus className="h-4 w-4" />
                              </Button>

                              <span className="text-xl font-bold w-16 text-center text-white">
                                {item.quantity}
                              </span>

                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                                disabled={isLoading}
                                className="h-10 w-10 rounded-lg hover:bg-gray-700 text-gray-300 hover:text-white"
                              >
                                <Plus className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>

                          {/* Action Buttons */}
                          <div className="flex flex-col gap-3">
                            <Button
                              variant="outline"
                              size="sm"
                              className="border-gray-600 text-gray-300 hover:bg-gray-800 hover:text-white"
                            >
                              <Heart className="h-4 w-4 mr-2" />
                              Save for Later
                            </Button>

                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleRemoveItem(item.id)}
                              disabled={isLoading}
                              className="border-red-500/50 text-red-400 hover:bg-red-500 hover:text-white hover:border-red-500 transition-all duration-300"
                            >
                              <Trash2 className="h-4 w-4 mr-2" />
                              Remove
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* BiHub Order Summary */}
          <div className="lg:col-span-1">
            <div className="sticky top-8 space-y-6">
              {/* Promo Code */}
              <Card className="bg-gray-900/50 border-gray-700 hover:border-orange-500/50 transition-all duration-300">
                <CardHeader className="pb-4">
                  <CardTitle className="flex items-center gap-3 text-white">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center">
                      <span className="text-white text-sm">🎫</span>
                    </div>
                    Promo Code
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="flex gap-3">
                    <Input
                      placeholder="Enter promo code"
                      value={couponCode}
                      onChange={(e) => setCouponCode(e.target.value)}
                      className="flex-1 bg-gray-800 border-gray-600 text-white placeholder-gray-400 focus:border-orange-500"
                    />
                    <Button
                      onClick={handleApplyCoupon}
                      disabled={!couponCode.trim() || isApplyingCoupon}
                      className="btn-gradient"
                    >
                      {isApplyingCoupon ? 'Applying...' : 'Apply'}
                    </Button>
                  </div>
                </CardContent>
              </Card>

              {/* BiHub Order Summary */}
              <Card className="bg-gray-900/50 border-gray-700 hover:border-orange-500/50 transition-all duration-300">
                <CardHeader className="pb-6">
                  <CardTitle className="text-xl text-white flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center">
                      <ShoppingBag className="h-4 w-4 text-white" />
                    </div>
                    Order Summary
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-4">
                    <div className="flex justify-between text-lg">
                      <span className="text-gray-400">Subtotal</span>
                      <span className="font-semibold text-white">{formatPrice(cartTotal)}</span>
                    </div>

                    <div className="flex justify-between text-lg">
                      <span className="text-gray-400">Shipping</span>
                      <span className="font-semibold">
                        {shippingCost === 0 ? (
                          <span className="text-emerald-400 font-bold">Free</span>
                        ) : (
                          <span className="text-white">{formatPrice(shippingCost)}</span>
                        )}
                      </span>
                    </div>

                    <div className="flex justify-between text-lg">
                      <span className="text-gray-400">Tax</span>
                      <span className="font-semibold text-white">{formatPrice(tax)}</span>
                    </div>
                  </div>

                  <div className="border-t-2 border-gray-700 pt-6">
                    <div className="flex justify-between text-2xl font-bold">
                      <span className="text-white">Total</span>
                      <span className="text-orange-400">{formatPrice(finalTotal)}</span>
                    </div>
                  </div>

                  {cartTotal < 50 && (
                    <div className="bg-gradient-to-r from-blue-500/20 to-indigo-500/20 border border-blue-500/30 rounded-2xl p-4">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                          <Truck className="h-5 w-5 text-white" />
                        </div>
                        <div>
                          <p className="font-semibold text-blue-300">Almost there!</p>
                          <p className="text-sm text-blue-400">
                            Add {formatPrice(50 - cartTotal)} more for free shipping
                          </p>
                        </div>
                      </div>
                    </div>
                  )}

                  <div className="space-y-4">
                    <Button size="lg" className="w-full btn-gradient text-lg py-4" asChild>
                      <Link href="/checkout">
                        <CreditCard className="mr-3 h-5 w-5" />
                        Proceed to Checkout
                      </Link>
                    </Button>

                    <div className="text-center">
                      <Link
                        href="/products"
                        className="text-sm text-orange-400 hover:text-orange-300 transition-colors font-medium"
                      >
                        ← Continue Shopping
                      </Link>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* BiHub Trust Signals */}
              <Card className="bg-gray-900/50 border-gray-700">
                <CardContent className="p-6">
                  <h3 className="text-lg font-bold text-white mb-4">BiHub Guarantee</h3>
                  <div className="space-y-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-full flex items-center justify-center">
                        <Shield className="h-5 w-5 text-white" />
                      </div>
                      <span className="text-sm font-medium text-gray-300">Secure SSL encryption</span>
                    </div>

                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-full flex items-center justify-center">
                        <CreditCard className="h-5 w-5 text-white" />
                      </div>
                      <span className="text-sm font-medium text-gray-300">Multiple payment options</span>
                    </div>

                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-gradient-to-br from-purple-500 to-purple-600 rounded-full flex items-center justify-center">
                        <ArrowLeft className="h-5 w-5 text-white" />
                      </div>
                      <span className="text-sm font-medium text-gray-300">30-day return policy</span>
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
