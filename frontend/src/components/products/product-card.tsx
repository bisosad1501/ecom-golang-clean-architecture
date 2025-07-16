'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Heart, ShoppingCart, Eye, Star } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { useAddToWishlist, useRemoveFromWishlist } from '@/hooks/use-products'
import { useProductRatingSummary } from '@/hooks/use-reviews'
import { CompactReviewSummary } from '@/components/reviews'
import { Product } from '@/types'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'
import { DESIGN_TOKENS } from '@/constants/design-tokens'
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'

interface ProductCardProps {
  product: Product
  className?: string
  showQuickView?: boolean
  showWishlist?: boolean
}

export function ProductCard({ 
  product, 
  className,
  showQuickView = true,
  showWishlist = true 
}: ProductCardProps) {
  const [isHovered, setIsHovered] = useState(false)
  const [imageLoading, setImageLoading] = useState(true)
  
  const { addItem, isLoading: cartLoading } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const addToWishlistMutation = useAddToWishlist()
  const removeFromWishlistMutation = useRemoveFromWishlist()
  const { data: ratingSummary } = useProductRatingSummary(product.id)

  const primaryImage = product.images?.[0]?.url || product.main_image || '/placeholder-product.jpg'
  const secondaryImage = product.images?.[1]?.url

  // Use backend computed fields directly
  const currentPrice = product.current_price
  const originalPrice = product.original_price
  const hasDiscount = product.has_discount
  const isOnSale = product.is_on_sale
  const discountPercentage = product.discount_percentage

  const stockQuantity = product.stock
  const stockStatus = product.stock_status
  const isLowStock = product.is_low_stock
  const isOutOfStock = stockStatus === 'out_of_stock' || stockQuantity <= 0

  const handleAddToCart = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    
    if (isOutOfStock) return
    
    try {
      await addItem(product.id, 1)
      toast.success('Added to cart!')
    } catch (error) {
      toast.error('Failed to add to cart')
    }
  }

  const handleWishlistToggle = async (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    
    if (!isAuthenticated) {
      toast.error('Please sign in to add to wishlist')
      return
    }

    try {
      // TODO: Check if product is already in wishlist
      const isInWishlist = false // This should come from wishlist state
      
      if (isInWishlist) {
        await removeFromWishlistMutation.mutateAsync(product.id)
      } else {
        await addToWishlistMutation.mutateAsync(product.id)
      }
    } catch (error) {
      // Error is handled by the mutation
    }
  }

  const handleQuickView = (e: React.MouseEvent) => {
    e.preventDefault()
    e.stopPropagation()
    // TODO: Open quick view modal
    toast.info('Quick view coming soon!')
  }

  return (
    <div className="relative group h-full">
      {/* Refined outer glow - subtle and elegant */}
      <div className={cn(
        'absolute -inset-0.5 rounded-3xl opacity-0 group-hover:opacity-40 transition-all duration-700 ease-out',
        'bg-gradient-to-br from-[#ff9000]/15 via-orange-500/8 to-amber-400/10 blur-lg'
      )} />
      
      <Card
        variant="elevated"
        padding="none"
        className={cn(
          'relative h-full overflow-hidden backdrop-blur-sm border text-white transition-all duration-500 ease-out',
          'bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80',
          'hover:shadow-xl hover:shadow-[#ff9000]/12 hover:-translate-y-2 hover:scale-[1.01]',
          'rounded-3xl backdrop-saturate-150 border-gray-700/40 hover:border-[#ff9000]/25',
          'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/2 before:via-transparent before:to-white/1 before:pointer-events-none before:rounded-3xl',
          className
        )}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <Link href={`/products/${product.id}`} className="block h-full">
        {/* Optimized image container with perfect proportions */}
        <div className="relative aspect-[4/3] overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 rounded-t-3xl">
          {/* Subtle shimmer effect */}
          <div className={cn(
            'absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/6 to-transparent',
            'transition-transform duration-1500 ease-out',
            isHovered && 'translate-x-full'
          )} />
          
          {/* Featured badge */}
          {(product as any).featured && (
            <div className="absolute top-3 left-3 z-20">
              <Badge className="shadow-md bg-gradient-to-r from-purple-500 to-pink-500 text-white border-0 text-xs px-2 py-1 rounded-md backdrop-blur-sm">
                Featured
              </Badge>
            </div>
          )}

          {/* Discount badge */}
          {hasDiscount && discountPercentage > 0 && (
            <div className={cn(
              "absolute top-3 z-20",
              (product as any).featured ? "left-20" : "left-3"
            )}>
              <span className="text-xs font-medium text-white bg-[#ff9000] px-2 py-1 rounded-md shadow-md">
                -{Math.round(discountPercentage)}%
              </span>
            </div>
          )}

          {/* Stock status badges */}
          {isOutOfStock && (
            <div className="absolute top-3 right-3 z-20">
              <Badge
                className="shadow-md bg-red-500/90 text-white border border-red-400/20 text-xs px-2 py-1 rounded-md backdrop-blur-sm"
                style={{
                  boxShadow: '0 2px 12px rgba(239, 68, 68, 0.2)'
                }}
              >
                Sold Out
              </Badge>
            </div>
          )}

          {/* Low stock warning */}
          {!isOutOfStock && isLowStock && (
            <div className="absolute top-3 right-3 z-20">
              <Badge className="shadow-md bg-amber-500/90 text-white border border-amber-400/20 text-xs px-2 py-1 rounded-md backdrop-blur-sm">
                Low Stock
              </Badge>
            </div>
          )}

          {/* Enhanced product image */}
          <div className="relative w-full h-full">
            {/* Loading state */}
            {imageLoading && (
              <div className="absolute inset-0 bg-gradient-to-br from-gray-200 to-gray-300">
                <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/8 to-transparent animate-shimmer" />
              </div>
            )}
            
            <Image
              src={isHovered && secondaryImage ? secondaryImage : primaryImage}
              alt={product.name}
              fill
              className={cn(
                'object-cover transition-all duration-700 ease-out',
                imageLoading && 'scale-105 blur-sm opacity-0',
                isHovered ? 'scale-108 brightness-105' : 'scale-100 brightness-100',
                !imageLoading && 'opacity-100'
              )}
              onLoad={() => setImageLoading(false)}
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            />
            
            {/* Gentle overlay - very subtle */}
            <div className={cn(
              'absolute inset-0 transition-all duration-500',
              'bg-gradient-to-t from-black/5 via-transparent to-transparent',
              isHovered && 'from-black/10 via-transparent to-transparent'
            )} />
          </div>

          {/* Clean and balanced action buttons */}
          <div className={cn(
            'absolute inset-0 flex items-center justify-center transition-all duration-400 ease-out',
            isHovered ? 'bg-black/8 backdrop-blur-[1px]' : 'bg-transparent'
          )}>
            <div className={cn(
              'flex gap-2 transform transition-all duration-400 ease-out',
              isHovered ? 'translate-y-0 opacity-100 scale-100' : 'translate-y-3 opacity-0 scale-95'
            )}>
              {/* Minimalist quick view button */}
              {showQuickView && (
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    'h-9 w-9 rounded-lg backdrop-blur-sm border border-white/15 shadow-md transition-all duration-250',
                    'bg-white/85 hover:bg-white text-slate-700 hover:text-blue-600',
                    'hover:scale-105',
                    'transform-gpu'
                  )}
                  onClick={handleQuickView}
                  style={{
                    boxShadow: '0 2px 12px rgba(0, 0, 0, 0.12)'
                  }}
                >
                  <Eye className="h-3.5 w-3.5" />
                </Button>
              )}

              {/* Minimalist wishlist button */}
              {showWishlist && isAuthenticated && (
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    'h-9 w-9 rounded-lg backdrop-blur-sm border border-white/15 shadow-md transition-all duration-250',
                    'bg-white/85 hover:bg-white text-slate-700 hover:text-red-500',
                    'hover:scale-105',
                    'transform-gpu'
                  )}
                  onClick={handleWishlistToggle}
                  disabled={addToWishlistMutation.isPending || removeFromWishlistMutation.isPending}
                  style={{
                    boxShadow: '0 2px 12px rgba(0, 0, 0, 0.12)'
                  }}
                >
                  <Heart className="h-3.5 w-3.5" />
                </Button>
              )}

              {/* Clean add to cart button */}
              {!isOutOfStock && (
                <Button
                  size="icon"
                  className={cn(
                    'h-9 w-9 rounded-lg border border-[#ff9000]/25 backdrop-blur-sm shadow-md transition-all duration-250',
                    'text-white hover:scale-105',
                    'transform-gpu'
                  )}
                  style={{
                    background: 'linear-gradient(135deg, rgba(255, 144, 0, 0.9) 0%, rgba(230, 126, 0, 0.85) 100%)',
                    boxShadow: '0 2px 12px rgba(255, 144, 0, 0.25)'
                  }}
                  onClick={handleAddToCart}
                  disabled={cartLoading}
                >
                  <ShoppingCart className="h-3.5 w-3.5" />
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Optimized product info section with perfect balance */}
        <div className="p-4 space-y-3 relative">
          {/* Category and Rating in one line for better space utilization */}
          <div className="flex items-center justify-between">
            {/* Minimal category badge */}
            {product.category && (
              <span className="text-xs font-medium text-[#ff9000] bg-[#ff9000]/10 px-2 py-1 rounded-md border border-[#ff9000]/20">
                {product.category.name}
              </span>
            )}
            
            {/* Compact rating display */}
            {product.rating_average > 0 && (
              <div className="flex items-center gap-1">
                <div className="flex items-center gap-0.5">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={cn(
                        'h-3 w-3 transition-all duration-200',
                        i < Math.floor(product.rating_average)
                          ? 'text-amber-400 fill-amber-400'
                          : 'text-gray-600'
                      )}
                    />
                  ))}
                </div>
                <span className="text-xs text-gray-400 ml-1">
                  {product.rating_average.toFixed(1)}
                </span>
              </div>
            )}
          </div>

          {/* Clean product name with optimal height */}
          <div>
            <h3 className="text-base font-semibold leading-tight line-clamp-2 text-white transition-colors duration-300 group-hover:text-[#ff9000]/90">
              {product.name}
            </h3>
          </div>

          {/* Streamlined price section */}
          <div className="flex items-center justify-between">
            <div className="flex items-baseline gap-2">
              <span className="text-lg font-bold text-white">
                {formatPrice(currentPrice)}
              </span>
              {hasDiscount && originalPrice && (
                <span className="text-sm line-through text-gray-500">
                  {formatPrice(originalPrice)}
                </span>
              )}
            </div>
            
            {/* Enhanced stock status indicator */}
            {!isOutOfStock && (
              <div className="flex items-center gap-1">
                {isLowStock ? (
                  <>
                    <div className="w-1.5 h-1.5 rounded-full bg-amber-400" />
                    <span className="text-xs text-amber-400 font-medium">
                      {stockQuantity} left
                    </span>
                  </>
                ) : stockQuantity <= 10 && stockQuantity > 0 ? (
                  <>
                    <div className="w-1.5 h-1.5 rounded-full bg-green-400" />
                    <span className="text-xs text-green-400 font-medium">
                      In Stock
                    </span>
                  </>
                ) : null}
              </div>
            )}
          </div>

          {/* Reviews summary - minimal and clean */}
          {ratingSummary && (
            <div className="pt-2 border-t border-gray-700/30">
              <CompactReviewSummary
                summary={ratingSummary}
                className="text-xs opacity-70"
              />
            </div>
          )}
        </div>
      </Link>
    </Card>
    </div>
  )
}
