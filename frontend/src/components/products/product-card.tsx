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
  
  const hasDiscount = product.sale_price && product.sale_price < product.price
  const discountPercentage = hasDiscount 
    ? Math.round(((product.price - product.sale_price!) / product.price) * 100)
    : 0

  const displayPrice = product.sale_price || product.price
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
      className={cn(
        'group relative overflow-hidden transition-all duration-300 hover:shadow-lg',
        className
      )}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <Link href={`/products/${product.id}`}>
        {/* Image container */}
        <div className="relative aspect-square overflow-hidden bg-gray-100">
          {/* Discount badge */}
          {hasDiscount && (
            <Badge 
              variant="destructive" 
              className="absolute top-2 left-2 z-10"
            >
              -{discountPercentage}%
            </Badge>
          )}

          {/* Out of stock badge */}
          {isOutOfStock && (
            <Badge 
              variant="secondary" 
              className="absolute top-2 right-2 z-10"
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
            'absolute inset-0 bg-black bg-opacity-0 transition-all duration-300 flex items-center justify-center',
            isHovered && 'bg-opacity-20'
          )}>
            <div className={cn(
              'flex space-x-2 transform transition-all duration-300',
              isHovered ? 'translate-y-0 opacity-100' : 'translate-y-4 opacity-0'
            )}>
              {/* Quick view */}
              {showQuickView && (
                <Button
                  variant="secondary"
                  size="icon"
                  className="h-10 w-10 rounded-full bg-white hover:bg-gray-100"
                  onClick={handleQuickView}
                >
                  <Eye className="h-4 w-4" />
                </Button>
              )}

              {/* Wishlist */}
              {showWishlist && isAuthenticated && (
                <Button
                  variant="secondary"
                  size="icon"
                  className="h-10 w-10 rounded-full bg-white hover:bg-gray-100"
                  onClick={handleWishlistToggle}
                  disabled={addToWishlistMutation.isPending || removeFromWishlistMutation.isPending}
                >
                  <Heart className="h-4 w-4" />
                </Button>
              )}

              {/* Add to cart */}
              {!isOutOfStock && (
                <Button
                  variant="default"
                  size="icon"
                  className="h-10 w-10 rounded-full"
                  onClick={handleAddToCart}
                  disabled={cartLoading}
                >
                  <ShoppingCart className="h-4 w-4" />
                </Button>
              )}
            </div>
          </div>
        </div>

        {/* Product info */}
        <div className="p-4">
          {/* Category */}
          {product.category && (
            <p className="text-xs text-gray-500 uppercase tracking-wide mb-1">
              {product.category.name}
            </p>
          )}

          {/* Product name */}
          <h3 className="font-medium text-gray-900 mb-2 line-clamp-2 group-hover:text-primary-600 transition-colors">
            {product.name}
          </h3>

          {/* Rating */}
          {product.rating && (
            <div className="flex items-center space-x-1 mb-2">
              <div className="flex">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    className={cn(
                      'h-3 w-3',
                      i < Math.floor(product.rating!.average)
                        ? 'text-yellow-400 fill-current'
                        : 'text-gray-300'
                    )}
                  />
                ))}
              </div>
              <span className="text-xs text-gray-500">
                ({product.rating.count})
              </span>
            </div>
          )}

          {/* Price */}
          <div className="flex items-center space-x-2">
            <span className="text-lg font-semibold text-gray-900">
              {formatPrice(displayPrice)}
            </span>
            {hasDiscount && (
              <span className="text-sm text-gray-500 line-through">
                {formatPrice(product.price)}
              </span>
            )}
          </div>

          {/* Stock status */}
          {product.stock <= 5 && product.stock > 0 && (
            <p className="text-xs text-orange-600 mt-1">
              Only {product.stock} left in stock
            </p>
          )}
        </div>
      </Link>

      {/* Quick add to cart button (bottom) */}
      {!isOutOfStock && (
        <div className={cn(
          'absolute bottom-0 left-0 right-0 p-4 transform transition-all duration-300',
          isHovered ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0'
        )}>
          <Button
            className="w-full"
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
