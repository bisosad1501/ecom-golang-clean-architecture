'use client'

import { useState, useEffect } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Minus, Plus, Trash2, ShoppingBag, ArrowLeft, Heart, Star, Shield, Truck, CreditCard, Gift, Zap, CheckCircle, AlertCircle, Eye } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useCartStore, getCartTotal, getCartItemCount, getCartSubtotal, isGuestCart } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'
import { useProductRatingSummary } from '@/hooks/use-reviews'
import { CompactReviewSummary } from '@/components/reviews'

// Cart Item Component đồng bộ với ProductCard design
function CartItemCard({ item, isLoading, onUpdateQuantity, onRemove, onAddToWishlist }: {
  item: any,
  isLoading: boolean,
  onUpdateQuantity: (id: string, quantity: number) => void,
  onRemove: (id: string) => void,
  onAddToWishlist: (productId: string) => void
}) {
  const { data: ratingSummary } = useProductRatingSummary(item.product.id)
  
  // Enhanced price logic giống ProductCard
  const currentPrice = item.price
  const originalPrice = item.product.compare_price || item.price * 1.2 // fallback for demo
  const hasDiscount = originalPrice && currentPrice < originalPrice
  const discountPercentage = hasDiscount
    ? Math.round(((originalPrice - currentPrice) / originalPrice) * 100)
    : 0

  return (
    <div className="relative group">
      {/* Refined outer glow giống ProductCard */}
      <div className={cn(
        'absolute -inset-0.5 rounded-3xl opacity-0 group-hover:opacity-40 transition-all duration-700 ease-out',
        'bg-gradient-to-br from-[#ff9000]/15 via-orange-500/8 to-amber-400/10 blur-lg'
      )} />
      
      <Card className={cn(
        'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-300 ease-out',
        'bg-gradient-to-br from-slate-900/80 via-gray-900/85 to-slate-800/80',
        'hover:shadow-lg hover:shadow-[#ff9000]/8 hover:-translate-y-0.5',
        'rounded-2xl backdrop-saturate-150 border-gray-700/50 hover:border-[#ff9000]/30',
        'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/1 before:via-transparent before:to-white/0.5 before:pointer-events-none before:rounded-2xl'
      )}>
        <CardContent className="p-4">
          <div className="flex gap-4 items-start">
            {/* Larger Product Image */}
            <div className="relative flex-shrink-0">
              <div className="relative w-36 h-36 rounded-xl overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 shadow-lg">
                {/* Discount badge */}
                {hasDiscount && (
                  <div className="absolute top-2 left-2 z-20">
                    <span className="text-sm font-bold text-white bg-[#ff9000] px-2.5 py-1 rounded-md shadow-md">
                      -{discountPercentage}%
                    </span>
                  </div>
                )}

                <Image
                  src={item.product.images?.[0]?.url || '/placeholder-product.svg'}
                  alt={item.product.name}
                  fill
                  className="object-cover transition-all duration-300 ease-out group-hover:scale-105"
                  sizes="(max-width: 144px) 100vw, 144px"
                />
                
                {/* Quick View on hover */}
                <div className={cn(
                  'absolute inset-0 flex items-center justify-center transition-all duration-300',
                  'bg-black/60 backdrop-blur-sm opacity-0 group-hover:opacity-100'
                )}>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-9 w-9 bg-white/90 hover:bg-white text-slate-700 hover:text-blue-600 rounded-lg"
                    asChild
                  >
                    <Link href={`/products/${item.product.id}`}>
                      <Eye className="h-4 w-4" />
                    </Link>
                  </Button>
                </div>
              </div>
            </div>

            {/* Expanded Product Info Section */}
            <div className="flex-1 min-w-0 space-y-2">
              {/* Category and Wishlist Row */}
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-center gap-2">
                  {item.product.category && (
                    <span className="text-sm font-semibold text-[#ff9000] bg-[#ff9000]/15 px-3 py-1.5 rounded-lg border border-[#ff9000]/40 shadow-sm">
                      {item.product.category.name}
                    </span>
                  )}
                  
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onAddToWishlist(item.product.id)}
                    className="h-8 w-8 p-0 text-gray-400 hover:text-red-400 hover:bg-red-500/15 rounded-lg border border-transparent hover:border-red-500/30"
                  >
                    <Heart className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {/* Product Name - Much Larger */}
              <h3 className="text-xl font-bold leading-tight line-clamp-2 text-white group-hover:text-[#ff9000]/90 transition-colors">
                <Link href={`/products/${item.product.id}`} className="hover:text-[#ff9000] transition-colors">
                  {item.product.name}
                </Link>
              </h3>

              {/* Price Section - Larger and more prominent */}
              <div className="flex items-center gap-3">
                <span className="text-2xl font-bold text-white">
                  {formatPrice(currentPrice)}
                </span>
                {hasDiscount && (
                  <span className="text-lg line-through text-gray-500">
                    {formatPrice(originalPrice)}
                  </span>
                )}
              </div>

              {/* Review Section - Well positioned */}
              {ratingSummary && (
                <div>
                  <CompactReviewSummary
                    summary={ratingSummary}
                    className="text-sm"
                  />
                </div>
              )}
            </div>

            {/* Controls Section - Better organized */}
            <div className="flex flex-col items-end gap-3 flex-shrink-0 min-w-[200px]">
              {/* Quantity Controls - Larger and cleaner */}
              <div className="flex items-center gap-2">
                <span className="text-sm text-gray-400 font-medium whitespace-nowrap">Qty:</span>
                <div className="flex items-center bg-white/[0.08] rounded-lg border border-white/20 shadow-sm">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onUpdateQuantity(item.id, item.quantity - 1)}
                    disabled={isLoading || item.quantity <= 1}
                    className="h-9 w-9 p-0 rounded-l-lg hover:bg-[#ff9000]/20 text-gray-300 hover:text-[#ff9000] transition-colors disabled:opacity-50"
                  >
                    <Minus className="h-4 w-4" />
                  </Button>
                  
                  <span className="px-4 py-2 text-base font-bold text-white min-w-[3rem] text-center">
                    {item.quantity}
                  </span>
                  
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => onUpdateQuantity(item.id, item.quantity + 1)}
                    disabled={isLoading || item.quantity >= 10}
                    className="h-9 w-9 p-0 rounded-r-lg hover:bg-[#ff9000]/20 text-gray-300 hover:text-[#ff9000] transition-colors disabled:opacity-50"
                  >
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              {/* Item Subtotal - Clear but less prominent than final total */}
              <div className="text-right">
                <div className="text-sm text-gray-400 mb-1">Subtotal</div>
                <div className="text-xl font-bold text-white">
                  {formatPrice(currentPrice * item.quantity)}
                </div>
              </div>

              {/* Remove Button - More prominent */}
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onRemove(item.id)}
                className="h-10 px-4 text-gray-400 hover:text-red-400 hover:bg-red-500/15 border border-transparent hover:border-red-500/30 rounded-lg transition-colors font-medium"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Remove
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

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
  const [showClearConfirm, setShowClearConfirm] = useState(false)

  // Fetch cart data when component mounts
  useEffect(() => {
    fetchCart()
  }, [fetchCart])

  const cartTotal = getCartTotal(cart)
  const cartSubtotal = getCartSubtotal(cart)
  const cartItemCount = getCartItemCount(cart)
  const isGuest = isGuestCart(cart)
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
    setShowClearConfirm(true)
  }

  const confirmClearCart = async () => {
    try {
      await clearCart()
      toast.success('Cart cleared successfully!')
      setShowClearConfirm(false)
    } catch (error) {
      toast.error('Failed to clear cart')
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

  const handleAddToWishlist = async (productId: string) => {
    try {
      // Add to wishlist logic here
      toast.success('Added to wishlist!')
    } catch (error) {
      toast.error('Failed to add to wishlist')
    }
  }

  // Loading state
  if (isLoading && !cart) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-8 relative z-10 flex items-center justify-center min-h-screen">
          <div className="text-center">
            <div className="relative mb-8">
              <div className="animate-spin rounded-full h-20 w-20 border-4 border-gray-700/40 border-t-[#ff9000] mx-auto"></div>
              <div className="absolute inset-0 rounded-full h-20 w-20 border-4 border-transparent border-r-[#ff9000]/60 animate-pulse mx-auto"></div>
            </div>
            <div className="space-y-3">
              <h3 className="text-xl font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">Loading your cart...</h3>
              <p className="text-gray-400">Preparing your shopping experience</p>
              <div className="flex justify-center space-x-1 mt-4">
                <div className="w-2 h-2 bg-[#ff9000] rounded-full animate-bounce"></div>
                <div className="w-2 h-2 bg-[#ff9000] rounded-full animate-bounce" style={{animationDelay: '0.1s'}}></div>
                <div className="w-2 h-2 bg-[#ff9000] rounded-full animate-bounce" style={{animationDelay: '0.2s'}}></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-16 relative z-10">
          <div className="max-w-3xl mx-auto text-center">
            <div className="space-y-6">
              {/* Compact Empty Cart Hero */}
              <div className="relative animate-in fade-in slide-in-from-bottom-4 duration-500">
                <div className="w-32 h-32 mx-auto rounded-full bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80 backdrop-blur-sm flex items-center justify-center shadow-xl border border-gray-700/40 relative overflow-hidden">
                  <div className="absolute inset-0 bg-gradient-to-r from-[#ff9000]/10 to-transparent animate-pulse"></div>
                  <ShoppingBag className="h-16 w-16 text-gray-400 relative z-10" />
                </div>
                <div className="absolute -top-2 -right-2 w-10 h-10 bg-gradient-to-br from-[#ff9000] to-[#ff9000] rounded-full flex items-center justify-center shadow-lg animate-bounce">
                  <span className="text-white text-sm font-bold">0</span>
                </div>
              </div>

              <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500" style={{animationDelay: '0.2s'}}>
                <h1 className="text-3xl lg:text-4xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
                  Your <span className="text-[#ff9000]">Shopping Cart</span> is empty
                </h1>
                <p className="text-lg text-gray-300 leading-relaxed">
                  Discover amazing products and start building your perfect collection.
                </p>
              </div>

              {/* Compact Action buttons */}
              <div className="flex flex-col sm:flex-row gap-3 justify-center animate-in fade-in slide-in-from-bottom-4 duration-500" style={{animationDelay: '0.4s'}}>
                <Button 
                  size="lg" 
                  className="bg-gradient-to-r from-[#ff9000] to-[#ff9000] hover:from-[#ff9000]/90 hover:to-[#ff9000]/90 text-white font-bold py-3 px-6 text-base rounded-xl transition-all duration-300 hover:scale-105 hover:shadow-xl hover:shadow-[#ff9000]/25" 
                  asChild
                >
                  <Link href="/products">
                    <ShoppingBag className="mr-2 h-5 w-5" />
                    Start Shopping
                  </Link>
                </Button>
                <Button 
                  size="lg" 
                  variant="outline" 
                  className="border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-300 hover:bg-white/[0.08] hover:border-white/15 hover:text-white py-3 px-6 text-base rounded-xl transition-all duration-300" 
                  asChild
                >
                  <Link href="/">
                    Browse Categories
                  </Link>
                </Button>
              </div>

              {/* Compact Features */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-8 animate-in fade-in slide-in-from-bottom-4 duration-500" style={{animationDelay: '0.6s'}}>
                <Card className={cn(
                  'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-500 ease-out',
                  'bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80',
                  'rounded-xl backdrop-saturate-150 border-gray-700/40',
                  'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/2 before:via-transparent before:to-white/1 before:pointer-events-none before:rounded-xl'
                )}>
                  <CardContent className="p-4 text-center">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center mx-auto mb-3">
                      <Truck className="h-6 w-6 text-white" />
                    </div>
                    <h3 className="font-bold text-white mb-1 text-base">Free Shipping</h3>
                    <p className="text-gray-300 text-sm">On orders over $50</p>
                  </CardContent>
                </Card>
                <Card className={cn(
                  'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-500 ease-out',
                  'bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80',
                  'rounded-xl backdrop-saturate-150 border-gray-700/40',
                  'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/2 before:via-transparent before:to-white/1 before:pointer-events-none before:rounded-xl'
                )}>
                  <CardContent className="p-4 text-center">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center mx-auto mb-3">
                      <Shield className="h-6 w-6 text-white" />
                    </div>
                    <h3 className="font-bold text-white mb-1 text-base">Secure Payment</h3>
                    <p className="text-gray-300 text-sm">100% secure checkout</p>
                  </CardContent>
                </Card>
                <Card className={cn(
                  'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-500 ease-out',
                  'bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80',
                  'rounded-xl backdrop-saturate-150 border-gray-700/40',
                  'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/2 before:via-transparent before:to-white/1 before:pointer-events-none before:rounded-xl'
                )}>
                  <CardContent className="p-4 text-center">
                    <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center mx-auto mb-3">
                      <ArrowLeft className="h-6 w-6 text-white" />
                    </div>
                    <h3 className="font-bold text-white mb-1 text-base">Easy Returns</h3>
                    <p className="text-gray-300 text-sm">30-day return policy</p>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <>
      {/* Modern Clear Cart Confirmation Dialog */}
      {showClearConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black/70 backdrop-blur-sm"
            onClick={() => setShowClearConfirm(false)}
          />
          
          {/* Dialog */}
          <div className="relative z-10 w-full max-w-md mx-4 animate-in fade-in slide-in-from-bottom-4 duration-300">
            <Card className={cn(
              'relative overflow-hidden backdrop-blur-sm border text-white',
              'bg-gradient-to-br from-slate-900/90 via-gray-900/90 to-slate-800/90',
              'rounded-2xl backdrop-saturate-150 border-gray-700/40',
              'shadow-2xl shadow-red-500/10',
              'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/5 before:via-transparent before:to-white/2 before:pointer-events-none before:rounded-2xl'
            )}>
              <CardContent className="p-6">
                <div className="text-center space-y-6">
                  {/* Icon */}
                  <div className="w-16 h-16 rounded-full bg-gradient-to-br from-red-500/20 to-red-600/20 border border-red-500/30 flex items-center justify-center mx-auto">
                    <Trash2 className="h-8 w-8 text-red-400" />
                  </div>
                  
                  {/* Content */}
                  <div className="space-y-3">
                    <h3 className="text-xl font-bold text-white">Clear Shopping Cart</h3>
                    <p className="text-gray-300 leading-relaxed">
                      Are you sure you want to remove all items from your cart? This action cannot be undone.
                    </p>
                    <div className="bg-red-500/10 border border-red-500/20 rounded-lg p-3">
                      <p className="text-red-300 text-sm flex items-center gap-2">
                        <AlertCircle className="h-4 w-4" />
                        {cartItemCount} {cartItemCount === 1 ? 'item' : 'items'} will be removed
                      </p>
                    </div>
                  </div>
                  
                  {/* Actions */}
                  <div className="flex flex-col sm:flex-row gap-3">
                    <Button
                      variant="outline"
                      onClick={() => setShowClearConfirm(false)}
                      className="flex-1 border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-300 hover:bg-white/[0.08] hover:border-white/15 hover:text-white transition-all duration-300 rounded-lg"
                    >
                      Cancel
                    </Button>
                    <Button
                      onClick={confirmClearCart}
                      disabled={isLoading}
                      className="flex-1 bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white font-medium transition-all duration-300 hover:scale-105 hover:shadow-lg hover:shadow-red-500/25 rounded-lg"
                    >
                      {isLoading ? (
                        <>
                          <div className="animate-spin rounded-full h-4 w-4 border-2 border-white/30 border-t-white mr-2"></div>
                          Clearing...
                        </>
                      ) : (
                        <>
                          <Trash2 className="h-4 w-4 mr-2" />
                          Clear Cart
                        </>
                      )}
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      )}

      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      {/* Enhanced Background Pattern - Matching Products Page */}
      <AnimatedBackground className="opacity-30" />
      
      {/* Main Content Area */}
      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-8 relative z-10">
        {/* Cart Header - Matching Products Page Style */}
        <div className="mb-8">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6 mb-6">
            <div className="space-y-2">
              <h1 className="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
                Your <span className="text-[#ff9000]">Shopping Cart</span>
              </h1>
              <div className="flex items-center gap-3">
                <p className="text-gray-400 text-sm">
                  <span className="text-[#ff9000] font-medium">{cartItemCount}</span> {cartItemCount === 1 ? 'item' : 'items'} • Cart value: <span className="text-[#ff9000] font-medium">{formatPrice(cartTotal)}</span>
                </p>
                <Badge className="bg-white/8 text-gray-300 border-white/15 px-2 py-1 text-xs backdrop-blur-sm font-medium">
                  {cartItemCount} Items
                </Badge>
                {isGuest && (
                  <Badge className="bg-blue-500/20 text-blue-300 border-blue-500/30 px-2 py-1 text-xs backdrop-blur-sm font-medium">
                    Guest Cart
                  </Badge>
                )}
              </div>

              {/* Guest Cart Notice */}
              {isGuest && (
                <div className="mt-4 bg-gradient-to-r from-blue-500/20 to-indigo-500/20 border border-blue-500/30 rounded-lg p-3">
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                      <AlertCircle className="h-4 w-4 text-white" />
                    </div>
                    <div>
                      <p className="font-semibold text-blue-300 text-sm">Shopping as Guest</p>
                      <p className="text-xs text-blue-400">
                        <Link href="/auth/login" className="text-blue-300 hover:text-blue-200 underline">
                          Sign in
                        </Link> to save your cart and access order history
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Action Controls */}
            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="sm"
                onClick={handleClearCart}
                disabled={isLoading}
                className="rounded-lg border border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-400 hover:bg-red-500/10 hover:border-red-500/30 hover:text-red-400 transition-all duration-200 h-7 px-2.5 text-xs font-medium shadow-sm"
              >
                <Trash2 className="h-3 w-3 mr-1" />
                Clear Cart
              </Button>

              <Button 
                variant="outline" 
                size="sm" 
                className="rounded-lg border border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-400 hover:bg-white/[0.08] hover:border-white/15 hover:text-white transition-all duration-200 h-7 px-2.5 text-xs font-medium shadow-sm" 
                asChild
              >
                <Link href="/products" className="flex items-center">
                  <ArrowLeft className="h-3 w-3 mr-1" />
                  Continue Shopping
                </Link>
              </Button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 xl:grid-cols-3 gap-6 lg:gap-8">
          {/* Cart Items - Using New Component */}
          <div className="xl:col-span-2 space-y-6">
            {cart.items.map((item, index) => (
              <CartItemCard
                key={item.id}
                item={item}
                isLoading={isLoading}
                onUpdateQuantity={handleUpdateQuantity}
                onRemove={handleRemoveItem}
                onAddToWishlist={handleAddToWishlist}
              />
            ))}
          </div>

          {/* Order Summary - Matching Products Page Style */}
          <div className="lg:col-span-1">
            <div className="sticky top-8 space-y-4">
              {/* Promo Code - Balanced prominence */}
              <Card className={cn(
                'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-300 ease-out',
                'bg-gradient-to-br from-slate-900/75 via-gray-900/80 to-slate-800/75',
                'hover:shadow-md hover:shadow-purple-500/8 hover:-translate-y-0.5',
                'rounded-xl backdrop-saturate-150 border-gray-700/50 hover:border-purple-500/30',
                'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/1 before:via-transparent before:to-white/0.5 before:pointer-events-none before:rounded-xl'
              )}>
                <CardContent className="p-4">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center shadow-lg">
                      <Gift className="h-4 w-4 text-white" />
                    </div>
                    <div>
                      <h3 className="text-base font-semibold text-white">Promo Code</h3>
                      <p className="text-sm text-gray-400">Save up to 30% off</p>
                    </div>
                  </div>
                  
                  <div className="space-y-3">
                    <div className="flex gap-3">
                      <Input
                        placeholder="Enter promo code"
                        value={couponCode}
                        onChange={(e) => setCouponCode(e.target.value)}
                        className="flex-1 bg-white/[0.08] border-white/20 text-white placeholder-gray-400 focus:border-purple-500/50 focus:ring-purple-500/20 rounded-lg"
                      />
                      <Button
                        onClick={handleApplyCoupon}
                        disabled={!couponCode.trim() || isApplyingCoupon}
                        className="bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white font-medium px-4 transition-all duration-300 rounded-lg"
                      >
                        {isApplyingCoupon ? 'Applying...' : 'Apply'}
                      </Button>
                    </div>

                    <div className="flex gap-2 flex-wrap">
                      <Badge 
                        className="bg-purple-500/20 text-purple-400 border-purple-500/30 cursor-pointer hover:bg-purple-500/30 transition-colors duration-300 text-xs"
                        onClick={() => setCouponCode('SAVE20')}
                      >
                        SAVE20
                      </Badge>
                      <Badge 
                        className="bg-purple-500/20 text-purple-400 border-purple-500/30 cursor-pointer hover:bg-purple-500/30 transition-colors duration-300 text-xs"
                        onClick={() => setCouponCode('FIRST10')}
                      >
                        FIRST10
                      </Badge>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Order Summary - Main focus */}
              <Card className={cn(
                'relative overflow-hidden backdrop-blur-sm border text-white transition-all duration-300 ease-out',
                'bg-gradient-to-br from-slate-900/85 via-gray-900/90 to-slate-800/85',
                'hover:shadow-lg hover:shadow-[#ff9000]/8 hover:-translate-y-0.5',
                'rounded-2xl backdrop-saturate-150 border-gray-700/60 hover:border-[#ff9000]/30',
                'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/1 before:via-transparent before:to-white/0.5 before:pointer-events-none before:rounded-2xl'
              )}>
                <CardContent className="p-5">
                  <div className="flex items-center gap-3 mb-4">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#ff9000] to-[#ff9000] flex items-center justify-center shadow-lg">
                      <ShoppingBag className="h-5 w-5 text-white" />
                    </div>
                    <div>
                      <h3 className="text-lg font-bold text-white">Order Summary</h3>
                      <p className="text-sm text-gray-400">{cartItemCount} items in cart</p>
                    </div>
                  </div>

                  <div className="space-y-4">
                    {/* Breakdown Items - Simplified and less prominent */}
                    <div className="space-y-2 py-3 border-b border-white/10">
                      {/* Subtotal - Less prominent */}
                      <div className="flex justify-between items-center text-sm">
                        <span className="text-gray-400">Subtotal ({cartItemCount} items)</span>
                        <span className="text-gray-300">{formatPrice(cartSubtotal)}</span>
                      </div>

                      {/* Shipping - Less prominent */}
                      <div className="flex justify-between items-center text-sm">
                        <span className="text-gray-400">Shipping</span>
                        <span className="text-gray-300">
                          {shippingCost === 0 ? (
                            <div className="flex items-center gap-2">
                              <span className="text-green-400 font-medium">FREE</span>
                              <Badge className="bg-green-500/20 text-green-400 border-green-500/30 text-xs">
                                Over $50
                              </Badge>
                            </div>
                          ) : (
                            <span>{formatPrice(shippingCost)}</span>
                          )}
                        </span>
                      </div>

                      {/* Tax - Less prominent */}
                      <div className="flex justify-between items-center text-sm">
                        <span className="text-gray-400">Tax (8%)</span>
                        <span className="text-gray-300">{formatPrice(tax)}</span>
                      </div>
                    </div>

                    {/* FINAL TOTAL - The only prominent total */}
                    <div className="relative">
                      <div className="absolute inset-0 bg-gradient-to-r from-[#ff9000]/20 to-[#ff9000]/5 rounded-xl blur-lg"></div>
                      <div className="relative border-2 border-[#ff9000]/40 rounded-xl p-4 bg-gradient-to-r from-[#ff9000]/10 to-[#ff9000]/5 backdrop-blur-sm">
                        <div className="flex justify-between items-center">
                          <span className="text-white font-bold text-xl">Total</span>
                          <div className="text-right">
                            <span className="text-3xl font-bold text-[#ff9000] drop-shadow-lg">{formatPrice(finalTotal)}</span>
                            <p className="text-xs text-gray-400 mt-1">Including all taxes & fees</p>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Savings Notice */}
                    {cartTotal > 100 && (
                      <div className="bg-gradient-to-r from-green-500/20 to-emerald-500/20 border border-green-500/30 rounded-lg p-3 text-center">
                        <div className="flex items-center justify-center gap-2 text-green-400 text-sm">
                          <CheckCircle className="h-4 w-4" />
                          <span className="font-semibold">You saved {formatPrice(cartTotal * 0.1)} today!</span>
                        </div>
                      </div>
                    )}

                    {cartTotal < 50 && (
                      <div className="bg-gradient-to-r from-blue-500/20 to-indigo-500/20 border border-blue-500/30 rounded-lg p-3">
                        <div className="flex items-center gap-3">
                          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                            <Truck className="h-4 w-4 text-white" />
                          </div>
                          <div>
                            <p className="font-semibold text-blue-300 text-sm">Almost there!</p>
                            <p className="text-xs text-blue-400">
                              Add {formatPrice(50 - cartTotal)} more for free shipping
                            </p>
                          </div>
                        </div>
                      </div>
                    )}

                    {/* Checkout Button - Most prominent action */}
                    <Button
                      size="lg"
                      className="w-full bg-gradient-to-r from-[#ff9000] to-[#ff9000] hover:from-[#ff9000]/90 hover:to-[#ff9000]/90 text-white font-bold py-4 text-lg rounded-xl transition-all duration-300 hover:scale-[1.02] hover:shadow-2xl hover:shadow-[#ff9000]/30 ring-2 ring-[#ff9000]/20 hover:ring-[#ff9000]/40"
                      asChild
                    >
                      <Link href="/checkout" className="flex items-center justify-center gap-2">
                        <CreditCard className="h-5 w-5" />
                        <span>Proceed to Checkout</span>
                      </Link>
                    </Button>

                    {/* Security Features - Smaller and less prominent */}
                    <div className="grid grid-cols-3 gap-2">
                      <div className="flex flex-col items-center gap-1 p-2 bg-white/[0.05] rounded-lg border border-white/10">
                        <Shield className="h-3 w-3 text-green-400" />
                        <span className="text-xs text-gray-500">SSL Secured</span>
                      </div>
                      <div className="flex flex-col items-center gap-1 p-2 bg-white/[0.05] rounded-lg border border-white/10">
                        <CreditCard className="h-3 w-3 text-blue-400" />
                        <span className="text-xs text-gray-500">Safe Payment</span>
                      </div>
                      <div className="flex flex-col items-center gap-1 p-2 bg-white/[0.05] rounded-lg border border-white/10">
                        <Truck className="h-3 w-3 text-[#ff9000]" />
                        <span className="text-xs text-gray-500">Fast Delivery</span>
                      </div>
                    </div>

                    {/* Continue Shopping - Subtle */}
                    <div className="text-center">
                      <Link
                        href="/products"
                        className="inline-flex items-center gap-2 text-xs text-gray-500 hover:text-[#ff9000]/70 transition-all duration-300 font-medium"
                      >
                        <ArrowLeft className="h-3 w-3" />
                        Continue Shopping
                      </Link>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </div>
    </>
  )
}
