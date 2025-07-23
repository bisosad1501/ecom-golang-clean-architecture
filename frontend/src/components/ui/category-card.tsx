'use client'

import React, { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ChevronRight, FolderTree, Package } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Category } from '@/types'
import { CompactCategoryPath } from '@/components/ui/category-path-indicator'



interface CategoryCardProps {
  category: Category
  variant?: 'default' | 'compact' | 'featured' | 'list'
  showStats?: boolean
  showSubcategories?: boolean
  className?: string
  priority?: boolean
}

export function CategoryCard({
  category,
  variant = 'default',
  showStats = true,
  showSubcategories = true,
  className,
  priority = false
}: CategoryCardProps) {
  const [imageLoaded, setImageLoaded] = useState(false)
  const [imageError, setImageError] = useState(false)

  const cardVariants = {
    default: 'aspect-[4/3]',
    compact: 'aspect-[4/3]',
    featured: 'aspect-[4/3]',
    list: 'aspect-square w-16 h-16'
  }



  const hasImage = category.image && !imageError
  const productCount = category.product_count || 0
  const hasSubcategories = category.children && category.children.length > 0

  // List variant has different layout
  if (variant === 'list') {
    return (
      <Link href={`/categories/${category.slug || category.id}`} className="group block h-full">
        <Card className={cn(
          'overflow-hidden transition-all duration-300 ease-out h-full',
          'bg-white/[0.08] backdrop-blur-sm border border-white/15',
          'hover:bg-white/[0.12] hover:border-[#ff9000]/50 hover:shadow-lg hover:shadow-[#ff9000]/10',
          'hover:scale-[1.01]',
          'focus-within:ring-2 focus-within:ring-[#ff9000]/50 focus-within:ring-offset-2 focus-within:ring-offset-black',
          className
        )}>
          <CardContent className="p-0 h-full">
            <div className="flex gap-4 p-4 h-full">
              {/* Image Section */}
              <div className={cn(
                'relative overflow-hidden bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/10 rounded-lg flex-shrink-0',
                cardVariants[variant]
              )}>
                {hasImage ? (
                  <>
                    <Image
                      src={category.image}
                      alt={category.name}
                      fill
                      className={cn(
                        'object-cover transition-all duration-500 ease-out',
                        'group-hover:scale-110',
                        imageLoaded ? 'opacity-100' : 'opacity-0'
                      )}
                      onLoad={() => setImageLoaded(true)}
                      onError={() => setImageError(true)}
                      priority={priority}
                      sizes="80px"
                    />
                    {!imageLoaded && (
                      <div className="absolute inset-0 bg-gradient-to-br from-gray-700/50 to-gray-800/50 animate-pulse" />
                    )}
                  </>
                ) : (
                  <div className="absolute inset-0 flex items-center justify-center">
                    <FolderTree className="w-6 h-6 text-[#ff9000] transition-transform duration-300 group-hover:scale-110" />
                  </div>
                )}
              </div>

              {/* Content Section */}
              <div className="flex-1 min-w-0 flex flex-col justify-between">
                <div>
                  <h3 className={cn(
                    'text-white line-clamp-1 transition-colors duration-300 mb-2',
                    'group-hover:text-[#ff9000]',
                    'text-lg font-semibold'
                  )}>
                    {category.name}
                  </h3>

                  {category.description && (
                    <p className={cn(
                      'text-gray-400 line-clamp-2 transition-colors duration-300 mb-3',
                      'group-hover:text-gray-300',
                      'text-sm'
                    )}>
                      {category.description}
                    </p>
                  )}
                </div>

                <div className="flex items-center justify-between mt-auto">
                  {showStats && (
                    <div className="flex items-center text-gray-400">
                      <Package className="w-4 h-4 mr-1" />
                      <span className="text-xs">
                        {productCount} products
                      </span>
                    </div>
                  )}

                  <ChevronRight className={cn(
                    'w-5 h-5 text-[#ff9000] transition-transform duration-300',
                    'group-hover:translate-x-1'
                  )} />
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </Link>
    )
  }

  return (
    <div className="relative group h-full">
      {/* Refined outer glow - matching product cards */}
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
      >
        <Link href={`/categories/${category.slug || category.id}`} className="block h-full">
            {/* Image container with exact same proportions as product card */}
            <div className="relative aspect-[4/3] overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 rounded-t-3xl">
              {/* Subtle shimmer effect */}
              <div className={cn(
                'absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/6 to-transparent',
                'transition-transform duration-1500 ease-out',
                'group-hover:translate-x-full'
              )} />

              {/* Product count badge */}
              {showStats && productCount > 0 && (
                <div className="absolute top-3 right-3 z-20">
                  <Badge
                    className="shadow-md bg-[#ff9000]/90 text-white border border-[#ff9000]/20 text-xs px-2 py-1 rounded-md backdrop-blur-sm"
                    style={{
                      boxShadow: '0 2px 12px rgba(255, 144, 0, 0.2)'
                    }}
                  >
                    {productCount} products
                  </Badge>
                </div>
              )}

              {/* Featured badge */}
              {variant === 'featured' && (
                <div className="absolute top-3 left-3 z-20">
                  <Badge className="shadow-md bg-gradient-to-r from-purple-500 to-pink-500 text-white border-0 text-xs px-2 py-1 rounded-md backdrop-blur-sm">
                    Featured
                  </Badge>
                </div>
              )}

              {/* Category content - always show image with fallback */}
              <div className="relative w-full h-full">
                <Image
                  src={category.image || '/placeholder-product.svg'}
                  alt={category.name}
                  fill
                  className={cn(
                    'object-cover transition-all duration-700 ease-out',
                    imageLoaded ? 'opacity-100' : 'scale-105 blur-sm opacity-0',
                    'group-hover:scale-108 group-hover:brightness-105'
                  )}
                  onLoad={() => setImageLoaded(true)}
                  onError={() => setImageError(true)}
                  priority={priority}
                  sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
                />
                {/* Loading state */}
                {!imageLoaded && (
                  <div className="absolute inset-0 bg-gradient-to-br from-gray-200 to-gray-300">
                    <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/8 to-transparent animate-shimmer" />
                  </div>
                )}

                {/* Gentle overlay - very subtle */}
                <div className={cn(
                  'absolute inset-0 transition-all duration-500',
                  'bg-gradient-to-t from-black/5 via-transparent to-transparent'
                )} />
              </div>
            </div>

            {/* Category info section with exact same structure as product card */}
            <div className="p-4 space-y-3 relative">
              {/* Category path - clean and simple */}
              <div>
                <CompactCategoryPath category={category} />
              </div>

              {/* Clean category name with optimal height */}
              <div>
                <h3 className="text-base font-semibold leading-tight line-clamp-2 text-white transition-colors duration-300 group-hover:text-[#ff9000]/90">
                  {category.name}
                </h3>
              </div>

              {/* Description section - clean and simple */}
              {category.description && variant !== 'compact' && (
                <div>
                  <p className="text-sm text-gray-400 line-clamp-2">
                    {category.description}
                  </p>
                </div>
              )}

              {/* Subcategories indicator - bottom section like product card */}
              {hasSubcategories && showSubcategories && (
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-1">
                    <div className="w-1.5 h-1.5 rounded-full bg-[#ff9000]" />
                    <span className="text-xs text-[#ff9000] font-medium">
                      {category.children?.length} subcategories
                    </span>
                  </div>
                </div>
              )}
              </div>
          </Link>
        </Card>
    </div>
  )
}

// Skeleton loader for category cards
export function CategoryCardSkeleton({ variant = 'default' }: { variant?: 'default' | 'compact' | 'featured' }) {
  const cardVariants = {
    default: 'aspect-[4/3]',
    compact: 'aspect-square',
    featured: 'aspect-[16/9]'
  }

  return (
    <Card className="overflow-hidden bg-white/[0.08] backdrop-blur-sm border border-white/15">
      <CardContent className="p-0">
        <div className={cn(
          'bg-gradient-to-br from-gray-700/50 to-gray-800/50 animate-pulse',
          cardVariants[variant]
        )} />
        <div className={cn(
          'p-4 space-y-3',
          variant === 'compact' && 'p-3 space-y-2'
        )}>
          <div className="h-5 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded animate-pulse" />
          {variant !== 'compact' && (
            <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-3/4 animate-pulse" />
          )}
          <div className="flex justify-between items-center">
            <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-1/3 animate-pulse" />
            <div className="h-4 w-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded animate-pulse" />
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
