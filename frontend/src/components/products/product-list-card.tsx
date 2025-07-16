'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { Heart, ShoppingCart, Star } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { useAddToWishlist, useRemoveFromWishlist } from '@/hooks/use-products'
import { Product } from '@/types'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'
import { DESIGN_SYSTEM, getCardClasses } from '@/constants/design-system'

interface ProductListCardProps {
  product: Product
  className?: string
}

export function ProductListCard({ product, className }: ProductListCardProps) {
  const [imageLoading, setImageLoading] = useState(true)
  
  const { addItem, isLoading: cartLoading } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const addToWishlistMutation = useAddToWishlist()
  const removeFromWishlistMutation = useRemoveFromWishlist()

  const primaryImage = product.images?.[0]?.url || product.main_image || '/placeholder-product.jpg'

  // Use backend computed fields directly
  const currentPrice = product.current_price
  const originalPrice = product.original_price
  const isOnSale = product.is_on_sale
  const hasDiscount = product.has_discount
  const discountPercentage = product.discount_percentage
  const featured = product.featured
  const stockStatus = product.stock_status
  const isLowStock = product.is_low_stock
  const stockQuantity = product.stock
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

  return (
    <Card
      className={cn(
        getCardClasses('sm'),
        'group overflow-hidden text-white transition-all duration-300',
        className
      )}
    >
      <Link href={`/products/${product.id}`}>
        <div className="flex gap-4 p-4">
          {/* Product Image */}
          <div className="relative w-32 h-32 flex-shrink-0 overflow-hidden rounded-lg bg-gray-800">
            {/* Enhanced badges */}
            {featured && (
              <Badge className="absolute top-2 left-2 z-10 bg-gradient-to-r from-purple-500 to-pink-500 text-white text-xs px-2 py-1">
                Featured
              </Badge>
            )}

            {hasDiscount && discountPercentage > 0 && (
              <Badge
                className={cn(
                  "absolute top-2 z-10 bg-[#ff9000] text-white text-xs px-2 py-1",
                  featured ? "left-20" : "left-2"
                )}
              >
                -{Math.round(discountPercentage)}%
              </Badge>
            )}

            {/* Stock status badges */}
            {isOutOfStock && (
              <Badge
                className="absolute top-2 right-2 z-10 bg-gray-700 text-gray-200 text-xs px-2 py-1"
              >
                Out of Stock
              </Badge>
            )}

            <Image
              src={primaryImage}
              alt={product.name}
              fill
              className={cn(
                'object-cover transition-all duration-300',
                imageLoading && 'scale-110 blur-sm'
              )}
              onLoad={() => setImageLoading(false)}
              sizes="128px"
            />
          </div>

          {/* Product Info */}
          <div className="flex-1 flex flex-col justify-between">
            <div>
              {/* Category */}
              {product.category && (
                <p className="text-xs uppercase tracking-wider text-gray-400 mb-1">
                  {product.category.name}
                </p>
              )}

              {/* Product name */}
              <h3 className="text-lg font-semibold text-white mb-2 line-clamp-2 hover:text-orange-400 transition-colors">
                {product.name}
              </h3>

              {/* Rating */}
              {product.rating && (
                <div className="flex items-center gap-2 mb-2">
                  <div className="flex">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={cn(
                          'h-4 w-4',
                          i < Math.floor(product.rating!.average)
                            ? 'text-amber-400 fill-current'
                            : 'text-gray-500'
                        )}
                      />
                    ))}
                  </div>
                  <span className="text-sm text-gray-400">
                    ({product.rating.count})
                  </span>
                </div>
              )}

              {/* Short description */}
              {product.short_description && (
                <p className="text-sm text-gray-400 line-clamp-2 mb-3">
                  {product.short_description}
                </p>
              )}
            </div>

            {/* Enhanced Price and Actions */}
            <div className="flex items-center justify-between">
              <div>
                {/* Enhanced Price Display */}
                <div className="flex items-baseline gap-2 mb-1">
                  <span className="text-xl font-bold text-white">
                    {formatPrice(currentPrice)}
                  </span>
                  {hasDiscount && originalPrice && (
                    <div className="flex items-center gap-2">
                      <span className="text-sm line-through text-gray-500">
                        {formatPrice(originalPrice)}
                      </span>
                      {discountPercentage > 0 && (
                        <Badge className="bg-[#ff9000]/20 text-[#ff9000] text-xs px-1.5 py-0.5">
                          -{Math.round(discountPercentage)}%
                        </Badge>
                      )}
                    </div>
                  )}
                </div>

                {/* Sale indicator */}
                {isOnSale && (
                  <div className="flex items-center gap-1 mb-1">
                    <span className="text-xs text-[#ff9000] font-medium">🔥 On Sale</span>
                  </div>
                )}

                {/* Enhanced Stock status */}
                {!isOutOfStock && (
                  <div className="flex items-center gap-2">
                    {isLowStock ? (
                      <p className="text-xs font-medium text-amber-400">
                        ⚠️ Only {stockQuantity} left
                      </p>
                    ) : stockQuantity <= 10 ? (
                      <p className="text-xs font-medium text-green-400">
                        ✅ {stockQuantity} in stock
                      </p>
                    ) : (
                      <p className="text-xs font-medium text-green-400">
                        ✅ In Stock
                      </p>
                    )}

                    {featured && (
                      <Badge className="bg-purple-500/20 text-purple-400 text-xs px-1.5 py-0.5">
                        Featured
                      </Badge>
                    )}
                  </div>
                )}

                {isOutOfStock && (
                  <p className="text-xs font-medium text-red-400">
                    ❌ Out of Stock
                  </p>
                )}
              </div>

              {/* Action buttons */}
              <div className="flex items-center gap-2">
                {/* Wishlist */}
                {isAuthenticated && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-9 w-9 text-gray-400 hover:text-white hover:bg-gray-800"
                    onClick={handleWishlistToggle}
                    disabled={addToWishlistMutation.isPending || removeFromWishlistMutation.isPending}
                  >
                    <Heart className="h-4 w-4" />
                  </Button>
                )}

                {/* Add to cart */}
                {!isOutOfStock && (
                  <Button
                    size="sm"
                    className="bg-orange-500 hover:bg-orange-600 text-white"
                    onClick={handleAddToCart}
                    disabled={cartLoading}
                  >
                    <ShoppingCart className="h-4 w-4 mr-2" />
                    {cartLoading ? 'Adding...' : 'Add to Cart'}
                  </Button>
                )}
              </div>
            </div>
          </div>
        </div>
      </Link>
    </Card>
  )
}
