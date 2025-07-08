'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Minus, Plus, Trash2, ShoppingBag, ArrowLeft, Heart, Star, Shield, Truck, CreditCard, Gift, Zap, CheckCircle, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { useCartStore, getCartTotal, getCartItemCount } from '@/store/cart'
import { formatPrice } from '@/lib/utils'
import { toast } from 'sonner'

export function CartPage() {
  const { 
    cart, 
    isLoading, 
    updateItem, 
    removeItem, 
    clearCart,
    fetchCart
  } = useCartStore()
  
  const [couponCode, setCouponCode] = useState('')
  const [isApplyingCoupon, setIsApplyingCoupon] = useState(false)

  // Fetch cart data when component mounts
  useEffect(() => {
    fetchCart()
  }, [fetchCart])

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

  // Loading state
  if (isLoading && !cart) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black flex items-center justify-center">
        <div className="text-center">
          <div className="relative mb-8">
            <div className="animate-spin rounded-full h-20 w-20 border-4 border-gray-700 border-t-orange-500 mx-auto"></div>
            <div className="absolute inset-0 rounded-full h-20 w-20 border-4 border-transparent border-r-orange-400 animate-pulse mx-auto"></div>
          </div>
          <div className="space-y-3">
            <h3 className="text-xl font-semibold text-white">Loading your cart...</h3>
            <p className="text-gray-400">Preparing your shopping experience</p>
            <div className="flex justify-center space-x-1 mt-4">
              <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce"></div>
              <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce" style={{animationDelay: '0.1s'}}></div>
              <div className="w-2 h-2 bg-orange-500 rounded-full animate-bounce" style={{animationDelay: '0.2s'}}></div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black py-16">
        <div className="container mx-auto px-4">
          <div className="max-w-4xl mx-auto text-center">
            {/* Enhanced Empty Cart Hero */}
            <div className="relative mb-12 animate-fade-in">
              <div className="w-48 h-48 mx-auto rounded-full bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm flex items-center justify-center shadow-2xl border border-gray-700/50 relative overflow-hidden">
                <div className="absolute inset-0 bg-gradient-to-r from-orange-500/10 to-transparent animate-pulse"></div>
                <ShoppingBag className="h-24 w-24 text-gray-400 relative z-10" />
              </div>
              <div className="absolute -top-4 -right-4 w-16 h-16 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-xl animate-bounce">
                <span className="text-white text-xl font-bold">0</span>
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
        {/* Simple Cart Header */}
        <div className="mb-12 animate-fade-in">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-8">
            <div>
              <h1 className="text-3xl lg:text-4xl font-bold text-white mb-3">
                Your Cart
              </h1>
              <p className="text-lg text-gray-300">
                {cartItemCount} {cartItemCount === 1 ? 'item' : 'items'} • Total: <span className="text-orange-400 font-bold">{formatPrice(cartTotal)}</span>
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

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 lg:gap-8 animate-slide-up">
          {/* BiHub Cart Items */}
          <div className="xl:col-span-2 space-y-4 lg:space-y-6">
            {cart.items.map((item, index) => (
              <Card
                key={item.id}
                className="bg-gray-900/50 backdrop-blur-sm border-gray-700 hover:border-orange-500/50 transition-all duration-500 hover:shadow-2xl hover:shadow-orange-500/20 hover:scale-[1.02] group animate-fade-in"
                style={{ animationDelay: `${index * 100}ms` }}
              >
                <CardContent className="p-4 sm:p-6">
                  <div className="flex flex-col sm:flex-row gap-4 sm:gap-6">
                    {/* Enhanced Product Image */}
                    <div className="relative h-28 w-28 sm:h-36 sm:w-36 flex-shrink-0 overflow-hidden rounded-2xl border border-gray-600 group-hover:border-orange-500/50 transition-all duration-500">
                      <div className="absolute inset-0 bg-gradient-to-br from-orange-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500 z-10"></div>
                      <Image
                        src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                        alt={item.product.name}
                        fill
                        className="object-cover transition-all duration-700 group-hover:scale-110 group-hover:rotate-1"
                      />
                      <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500" />

                      {/* Floating Badge */}
                      <div className="absolute -top-2 -right-2 w-8 h-8 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-lg transform group-hover:scale-110 transition-transform duration-300">
                        <span className="text-white text-xs font-bold">#{index + 1}</span>
                      </div>

                      {/* Quick Actions Overlay */}
                      <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-300 z-20">
                        <Button
                          size="sm"
                          variant="secondary"
                          className="bg-white/90 text-gray-900 hover:bg-white transform scale-90 group-hover:scale-100 transition-transform duration-300"
                          asChild
                        >
                          <Link href={`/products/${item.product.id}`}>
                            View Details
                          </Link>
                        </Button>
                      </div>
                    </div>

                    {/* Enhanced Product Details */}
                    <div className="flex-1 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-6">
                        <div className="flex-1 space-y-4">
                          {/* Product Header */}
                          <div className="space-y-3">
                            <div className="flex items-center gap-3">
                              <Badge className="bg-gradient-to-r from-orange-500/20 to-orange-600/20 text-orange-400 border-orange-500/30 text-xs font-medium px-3 py-1">
                                Premium
                              </Badge>
                              <div className="flex items-center gap-1">
                                {[...Array(5)].map((_, i) => (
                                  <Star key={i} className="h-4 w-4 fill-orange-400 text-orange-400 drop-shadow-sm" />
                                ))}
                                <span className="text-sm text-gray-400 ml-2">(4.9)</span>
                              </div>
                            </div>

                            <h3 className="text-2xl font-bold text-white leading-tight group-hover:text-orange-400 transition-colors duration-300">
                              <Link
                                href={`/products/${item.product.id}`}
                                className="hover:text-orange-400 transition-colors duration-300 cursor-pointer"
                              >
                              {item.product.name}
                            </Link>
                          </h3>

                            <div className="flex items-center gap-3 text-sm text-gray-400">
                              <span>SKU: <span className="font-mono text-gray-300">{item.product.sku}</span></span>
                              <Separator orientation="vertical" className="h-4 bg-gray-600" />
                              <span className="flex items-center gap-1">
                                <CheckCircle className="h-4 w-4 text-green-400" />
                                In Stock
                              </span>
                            </div>

                            {/* Enhanced Pricing */}
                            <div className="space-y-2">
                              <div className="flex items-baseline gap-3">
                                <span className="text-3xl font-bold text-orange-400">
                                  {formatPrice(item.price)}
                                </span>
                                <span className="text-sm text-gray-400 line-through">
                                  {formatPrice(item.price * 1.2)}
                                </span>
                                <Badge className="bg-green-500/20 text-green-400 border-green-500/30 text-xs">
                                  20% OFF
                                </Badge>
                              </div>
                              <div className="text-lg font-semibold text-white">
                                Subtotal: <span className="text-orange-300 font-bold">{formatPrice(item.price * item.quantity)}</span>
                              </div>
                            </div>
                          </div>
                        </div>

                        {/* Enhanced Quantity & Actions */}
                        <div className="flex flex-col sm:flex-row sm:items-start gap-4 sm:gap-6 w-full sm:w-auto">
                          {/* Advanced Quantity Controls */}
                          <div className="flex flex-col gap-4 sm:min-w-[140px]">
                            <div className="flex items-center justify-between">
                              <span className="text-sm font-medium text-gray-300">Quantity:</span>
                              <span className="text-xs text-gray-400">Max: 10</span>
                            </div>

                            <div className="relative">
                              <div className="flex items-center bg-gradient-to-r from-gray-800 to-gray-800/80 rounded-2xl p-1 border border-gray-600 hover:border-orange-500/50 transition-all duration-300 shadow-lg">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                                  disabled={isLoading || item.quantity <= 1}
                                  className="h-12 w-12 rounded-xl hover:bg-orange-500/20 text-gray-300 hover:text-orange-400 transition-all duration-300 transform hover:scale-110 disabled:opacity-50 disabled:hover:scale-100"
                                >
                                  <Minus className="h-5 w-5" />
                                </Button>

                                <div className="flex-1 text-center relative">
                                  <span className="text-2xl font-bold text-white block transition-all duration-300">
                                    {item.quantity}
                                  </span>
                                  <div className="absolute inset-0 bg-gradient-to-r from-orange-500/10 to-transparent rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                                </div>

                                <Button
                                  variant="ghost"
                                  size="sm"
                                  onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                                  disabled={isLoading || item.quantity >= 10}
                                  className="h-12 w-12 rounded-xl hover:bg-orange-500/20 text-gray-300 hover:text-orange-400 transition-all duration-300 transform hover:scale-110 disabled:opacity-50 disabled:hover:scale-100"
                                >
                                  <Plus className="h-5 w-5" />
                                </Button>
                              </div>

                              {/* Quantity Progress Bar */}
                              <div className="w-full bg-gray-700 rounded-full h-1 mt-2">
                                <div
                                  className="bg-gradient-to-r from-orange-500 to-orange-600 h-1 rounded-full transition-all duration-500"
                                  style={{ width: `${(item.quantity / 10) * 100}%` }}
                                ></div>
                              </div>
                            </div>
                          </div>

                          {/* Enhanced Action Buttons */}
                          <div className="flex flex-col gap-3">
                            <Button
                              variant="outline"
                              size="sm"
                              className="border-gray-600 text-gray-300 hover:bg-gradient-to-r hover:from-pink-500/20 hover:to-red-500/20 hover:border-pink-500/50 hover:text-pink-400 transition-all duration-300 transform hover:scale-105"
                            >
                              <Heart className="h-4 w-4 mr-2" />
                              Save for Later
                            </Button>

                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleRemoveItem(item.id)}
                              disabled={isLoading}
                              className="border-red-500/50 text-red-400 hover:bg-gradient-to-r hover:from-red-500 hover:to-red-600 hover:text-white hover:border-red-500 transition-all duration-300 transform hover:scale-105 hover:shadow-lg hover:shadow-red-500/25"
                            >
                              <Trash2 className="h-4 w-4 mr-2" />
                              {isLoading ? 'Removing...' : 'Remove'}
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
              {/* Enhanced Promo Code */}
              <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700 hover:border-purple-500/50 transition-all duration-500 hover:shadow-2xl hover:shadow-purple-500/10 group">
                <CardHeader className="pb-4">
                  <CardTitle className="flex items-center gap-3 text-white">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                      <Gift className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <span className="text-lg font-bold">Promo Code</span>
                      <p className="text-sm text-gray-400 font-normal">Save up to 30% off</p>
                    </div>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex gap-3">
                      <div className="relative flex-1">
                        <Input
                          placeholder="Enter promo code"
                          value={couponCode}
                          onChange={(e) => setCouponCode(e.target.value)}
                          className="bg-gray-800/50 border-gray-600 text-white placeholder-gray-400 focus:border-purple-500 focus:ring-purple-500/20 pr-12"
                        />
                        <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                          <Gift className="h-4 w-4 text-gray-400" />
                        </div>
                      </div>
                      <Button
                        onClick={handleApplyCoupon}
                        disabled={!couponCode.trim() || isApplyingCoupon}
                        className="bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white font-semibold px-6 transition-all duration-300 transform hover:scale-105 hover:shadow-lg hover:shadow-purple-500/25"
                      >
                        {isApplyingCoupon ? 'Applying...' : 'Apply'}
                      </Button>
                    </div>

                    {/* Popular Codes */}
                    <div className="flex gap-2 flex-wrap">
                      <Badge className="bg-purple-500/20 text-purple-400 border-purple-500/30 cursor-pointer hover:bg-purple-500/30 transition-colors duration-300">
                        SAVE20
                      </Badge>
                      <Badge className="bg-purple-500/20 text-purple-400 border-purple-500/30 cursor-pointer hover:bg-purple-500/30 transition-colors duration-300">
                        FIRST10
                      </Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Enhanced Order Summary */}
              <Card className="bg-gradient-to-br from-gray-900/60 to-gray-800/60 backdrop-blur-sm border-gray-700 hover:border-orange-500/50 transition-all duration-500 hover:shadow-2xl hover:shadow-orange-500/10 group">
                <CardHeader className="pb-6">
                  <CardTitle className="text-xl text-white flex items-center gap-3">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-300">
                      <ShoppingBag className="h-6 w-6 text-white" />
                    </div>
                    <div>
                      <span className="text-xl font-bold">Order Summary</span>
                      <p className="text-sm text-gray-400 font-normal">{cartItemCount} items in cart</p>
                    </div>
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  {/* Order Details */}
                  <div className="space-y-4">
                    <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                      <span className="text-gray-300 flex items-center gap-2">
                        <div className="w-2 h-2 bg-blue-400 rounded-full"></div>
                        Subtotal
                      </span>
                      <span className="font-semibold text-white text-lg">{formatPrice(cartTotal)}</span>
                    </div>

                    <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                      <span className="text-gray-300 flex items-center gap-2">
                        <Truck className="h-4 w-4 text-green-400" />
                        Shipping
                      </span>
                      <span className="font-semibold">
                        {shippingCost === 0 ? (
                          <div className="flex items-center gap-2">
                            <span className="text-green-400 font-bold">FREE</span>
                            <Badge className="bg-green-500/20 text-green-400 border-green-500/30 text-xs">
                              Over $50
                            </Badge>
                          </div>
                        ) : (
                          <span className="text-white text-lg">{formatPrice(shippingCost)}</span>
                        )}
                      </span>
                    </div>

                    <div className="flex justify-between items-center p-3 bg-gray-800/30 rounded-lg">
                      <span className="text-gray-300 flex items-center gap-2">
                        <div className="w-2 h-2 bg-yellow-400 rounded-full"></div>
                        Tax (8%)
                      </span>
                      <span className="font-semibold text-white text-lg">{formatPrice(tax)}</span>
                    </div>
                  </div>

                  {/* Total Section */}
                  <div className="relative">
                    <div className="absolute inset-0 bg-gradient-to-r from-orange-500/10 to-transparent rounded-xl"></div>
                    <div className="border-2 border-orange-500/30 rounded-xl p-4 bg-gradient-to-r from-orange-500/5 to-transparent">
                      <div className="flex justify-between items-center">
                        <span className="text-white text-xl font-bold">Total</span>
                        <div className="text-right">
                          <span className="text-3xl font-bold text-orange-400 block">{formatPrice(finalTotal)}</span>
                          <span className="text-sm text-gray-400">Including all taxes</span>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Savings Badge */}
                  {cartTotal > 100 && (
                    <div className="bg-gradient-to-r from-green-500/20 to-emerald-500/20 border border-green-500/30 rounded-lg p-3 text-center">
                      <div className="flex items-center justify-center gap-2 text-green-400">
                        <CheckCircle className="h-5 w-5" />
                        <span className="font-semibold">You saved {formatPrice(cartTotal * 0.1)} today!</span>
                      </div>
                    </div>
                  )}

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

                  {/* Enhanced Checkout Section */}
                  <div className="space-y-6">
                    {/* Main Checkout Button */}
                    <div className="relative group">
                      <div className="absolute -inset-1 bg-gradient-to-r from-orange-500 to-orange-600 rounded-2xl blur opacity-25 group-hover:opacity-75 transition duration-1000 group-hover:duration-200"></div>
                      <Button
                        size="lg"
                        className="relative w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-bold text-lg py-6 rounded-xl transition-all duration-300 transform hover:scale-[1.02] hover:shadow-2xl hover:shadow-orange-500/25"
                        asChild
                      >
                        <Link href="/checkout" className="flex items-center justify-center gap-3">
                          <CreditCard className="h-6 w-6" />
                          <span>Proceed to Secure Checkout</span>
                          <div className="ml-2 flex items-center gap-1">
                            <div className="w-2 h-2 bg-white rounded-full animate-pulse"></div>
                            <div className="w-2 h-2 bg-white rounded-full animate-pulse" style={{animationDelay: '0.2s'}}></div>
                            <div className="w-2 h-2 bg-white rounded-full animate-pulse" style={{animationDelay: '0.4s'}}></div>
                          </div>
                        </Link>
                      </Button>
                    </div>

                    {/* Security Features */}
                    <div className="grid grid-cols-3 gap-3 text-center">
                      <div className="flex flex-col items-center gap-2 p-3 bg-gray-800/30 rounded-lg">
                        <Shield className="h-5 w-5 text-green-400" />
                        <span className="text-xs text-gray-400">SSL Secured</span>
                      </div>
                      <div className="flex flex-col items-center gap-2 p-3 bg-gray-800/30 rounded-lg">
                        <CreditCard className="h-5 w-5 text-blue-400" />
                        <span className="text-xs text-gray-400">Safe Payment</span>
                      </div>
                      <div className="flex flex-col items-center gap-2 p-3 bg-gray-800/30 rounded-lg">
                        <Truck className="h-5 w-5 text-orange-400" />
                        <span className="text-xs text-gray-400">Fast Delivery</span>
                      </div>
                    </div>

                    {/* Continue Shopping */}
                    <div className="text-center">
                      <Link
                        href="/products"
                        className="inline-flex items-center gap-2 text-sm text-orange-400 hover:text-orange-300 transition-all duration-300 font-medium hover:gap-3"
                      >
                        <ArrowLeft className="h-4 w-4" />
                        Continue Shopping
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
