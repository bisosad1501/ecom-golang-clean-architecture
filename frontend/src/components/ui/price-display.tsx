'use client'

import { cn } from '@/lib/utils'
import { formatPrice, formatDiscountPercentage, getPriceInfo, getCurrentPrice, getOriginalPrice } from '@/lib/utils/price'
import { Product } from '@/types'

interface PriceDisplayProps {
  product: Product
  className?: string
  showDiscount?: boolean
  showOriginalPrice?: boolean
  size?: 'sm' | 'md' | 'lg'
}

export function PriceDisplay({
  product,
  className,
  showDiscount = true,
  showOriginalPrice = true,
  size = 'md'
}: PriceDisplayProps) {
  const priceInfo = getPriceInfo(product)
  const currentPrice = getCurrentPrice(product)
  const originalPrice = getOriginalPrice(product)
  
  const sizeClasses = {
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg'
  }
  
  const discountSizeClasses = {
    sm: 'text-xs',
    md: 'text-sm', 
    lg: 'text-base'
  }

  return (
    <div className={cn('flex items-center gap-2', className)}>
      {/* Current Price */}
      <span className={cn(
        'font-semibold text-green-600',
        sizeClasses[size]
      )}>
        {formatPrice(currentPrice)}
      </span>

      {/* Original Price (if discount) */}
      {priceInfo.hasDiscount && showOriginalPrice && originalPrice && (
        <span className={cn(
          'line-through text-gray-500',
          discountSizeClasses[size]
        )}>
          {formatPrice(originalPrice)}
        </span>
      )}

      {/* Discount Badge */}
      {priceInfo.hasDiscount && showDiscount && priceInfo.discountPercentage > 0 && (
        <span className={cn(
          'bg-red-500 text-white px-2 py-1 rounded text-xs font-medium',
          discountSizeClasses[size]
        )}>
          {formatDiscountPercentage(priceInfo.discountPercentage)}
        </span>
      )}

      {/* Availability indicator */}
      {!priceInfo.isAvailable && (
        <span className="text-red-500 text-xs font-medium">
          Out of Stock
        </span>
      )}
    </div>
  )
}

interface AdminPriceDisplayProps {
  product: Product
  className?: string
}

export function AdminPriceDisplay({ product, className }: AdminPriceDisplayProps) {
  const priceInfo = getPriceInfo(product)
  
  return (
    <div className={cn('space-y-1', className)}>
      <div className="flex items-center gap-2">
        <span className="text-sm text-gray-600">Price:</span>
        <span className="font-medium">{formatPrice(priceInfo.price)}</span>
      </div>
      
      {priceInfo.comparePrice && (
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-600">Compare:</span>
          <span className="text-sm">{formatPrice(priceInfo.comparePrice)}</span>
          {priceInfo.hasDiscount && (
            <span className="text-green-600 text-xs">
              ({formatDiscountPercentage(priceInfo.discountPercentage)})
            </span>
          )}
        </div>
      )}
      
      {priceInfo.costPrice && (
        <div className="flex items-center gap-2">
          <span className="text-sm text-gray-600">Cost:</span>
          <span className="text-sm">{formatPrice(priceInfo.costPrice)}</span>
          <span className="text-green-600 text-xs">
            (Profit: {formatPrice(priceInfo.price - priceInfo.costPrice)})
          </span>
        </div>
      )}
      
      <div className="flex items-center gap-2">
        <span className="text-sm text-gray-600">Status:</span>
        <span className={cn(
          'text-xs px-2 py-1 rounded',
          priceInfo.isAvailable 
            ? 'bg-green-100 text-green-800' 
            : 'bg-red-100 text-red-800'
        )}>
          {priceInfo.isAvailable ? 'Available' : 'Unavailable'}
        </span>
      </div>
    </div>
  )
}
