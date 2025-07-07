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

  const primaryImage = product.images?.[0]?.url || '/placeholder-product.jpg'
  const secondaryImage = product.images?.[1]?.url
  
  const hasDiscount = product.compare_price && product.compare_price < product.price
  const discountPercentage = hasDiscount 
    ? Math.round(((product.price - product.compare_price!) / product.price) * 100)
    : 0

  const displayPrice = product.compare_price || product.price
  const isOutOfStock = product.stock <= 0

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
    <Card
      variant="elevated"
      padding="none"
      className={cn(
        'group relative overflow-hidden card-hover border border-gray-600 bg-gray-900 hover:border-gray-500 text-white transition-all duration-300',
        className
      )}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Link href={`/products/${product.id}`}>
        {/* Image container */}
        <div className={`relative aspect-square overflow-hidden bg-gradient-to-br from-gray-800 to-gray-900 ${DESIGN_TOKENS.RADIUS.LARGE}`}>
          {/* Discount badge */}
          {hasDiscount && (
            <Badge
              className="absolute top-2 left-2 z-10 shadow-lg font-bold text-white text-xs px-2 py-1"
              style={{backgroundColor: '#FF9000'}}
            >
              -{discountPercentage}%
            </Badge>
          )}

          {/* Out of stock badge */}
          {isOutOfStock && (
            <Badge
              className="absolute top-2 right-2 z-10 shadow-lg bg-gray-800 text-gray-200 border border-gray-600 text-xs px-2 py-1"
            >
              Out of Stock
            </Badge>
          )}

          {/* Product image */}
          <div className="relative w-full h-full">
            <Image
              src={isHovered && secondaryImage ? secondaryImage : primaryImage}
              alt={product.name}
              fill
              className={cn(
                'object-cover transition-all duration-500',
                imageLoading && 'scale-110 blur-sm',
                isHovered && 'scale-105'
              )}
              onLoad={() => setImageLoading(false)}
              sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
            />
          </div>

          {/* Overlay actions */}
          <div className={cn(
            'absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent transition-all duration-300 flex items-center justify-center',
            isHovered && 'from-black/40'
          )}>
            <div className={cn(
              `flex ${DESIGN_TOKENS.SPACING.GAP_SMALL} transform transition-all duration-300`,
              isHovered ? 'translate-y-0 opacity-100 scale-100' : 'translate-y-4 opacity-0 scale-95'
            )}>
              {/* Quick view */}
              {showQuickView && (
                <Button
                  variant="ghost"
                  size="icon"
                  className={`${DESIGN_TOKENS.BUTTONS.ICON_DEFAULT} ${DESIGN_TOKENS.RADIUS.FULL} glass-effect hover:bg-white/90 shadow-medium`}
                  onClick={handleQuickView}
                >
                  <Eye className={DESIGN_TOKENS.ICONS.DEFAULT} />
                </Button>
              )}

              {/* Wishlist */}
              {showWishlist && isAuthenticated && (
                <Button
                  variant="ghost"
                  size="icon"
                  className={`${DESIGN_TOKENS.BUTTONS.ICON_DEFAULT} ${DESIGN_TOKENS.RADIUS.FULL} glass-effect hover:bg-white/90 shadow-medium`}
                  onClick={handleWishlistToggle}
                  disabled={addToWishlistMutation.isPending || removeFromWishlistMutation.isPending}
                >
                  <Heart className={DESIGN_TOKENS.ICONS.DEFAULT} />
                </Button>
              )}

              {/* Add to cart */}
              {!isOutOfStock && (
                <Button
                  variant="gradient"
                  size="icon"
                  className={`${DESIGN_TOKENS.BUTTONS.ICON_DEFAULT} ${DESIGN_TOKENS.RADIUS.FULL} shadow-large`}
                  onClick={handleAddToCart}
                  disabled={cartLoading}
                >
                  <ShoppingCart className={DESIGN_TOKENS.ICONS.DEFAULT} />
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Product info */}
        <div className="p-3">
          {/* Category */}
          {product.category && (
            <p className="text-xs uppercase tracking-wider text-gray-400 mb-1">
              {product.category.name}
            </p>
          )}

          {/* Product name */}
          <h3 className="text-sm font-semibold text-white mb-2 line-clamp-2 transition-colors duration-200 leading-tight"
            style={{
              color: 'white'
            }}
            onMouseEnter={(e) => e.currentTarget.style.color = '#FF9000'}
            onMouseLeave={(e) => e.currentTarget.style.color = 'white'}
          >
            {product.name}
          </h3>

          {/* Rating */}
          {product.rating && (
            <div className="flex items-center gap-1 mb-2">
              <div className="flex">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    className={cn(
                      'h-3 w-3',
                      i < Math.floor(product.rating!.average)
                        ? 'text-amber-400 fill-current'
                        : 'text-gray-500'
                    )}
                  />
                ))}
              </div>
              <span className="text-xs font-medium text-gray-400">
                ({product.rating.count})
              </span>
            </div>
          )}

          {/* Price */}
          <div className="flex items-baseline gap-2 mb-1">
            <span className="text-base font-bold text-white">
              {formatPrice(displayPrice)}
            </span>
            {hasDiscount && (
              <span className="text-xs line-through text-gray-500">
                {formatPrice(product.price)}
              </span>
            )}
          </div>

          {/* Stock status */}
          {product.stock <= 5 && product.stock > 0 && (
            <p className="text-xs font-medium" style={{color: '#FF9000'}}>
              Only {product.stock} left in stock
            </p>
          )}
        </div>
      </Link>

      {/* Quick add to cart button (bottom) */}
      {!isOutOfStock && (
        <div className={cn(
          'absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-gray-900 via-gray-900 to-transparent transform transition-all duration-300',
          isHovered ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0'
        )}>
          <Button
            variant="gradient"
            className="w-full shadow-lg text-white"
            onClick={handleAddToCart}
            disabled={cartLoading}
            size="sm"
          >
            {cartLoading ? 'Adding...' : 'Add to Cart'}
          </Button>
        </div>
      )}
    </Card>
  )
}
