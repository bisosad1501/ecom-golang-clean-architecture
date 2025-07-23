'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { 
  Heart, 
  ShoppingCart, 
  Trash2, 
  ArrowLeft, 
  Star, 
  Eye, 
  Filter,
  Grid3X3,
  List,
  Search,
  SortAsc,
  SortDesc,
  Package,
  Sparkles,
  Gift
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useWishlist, useRemoveFromWishlist, useClearWishlist } from '@/hooks/use-wishlist'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'
// import { useProductRatingSummary } from '@/hooks/use-reviews'
// import { CompactReviewSummary } from '@/components/reviews'

// Wishlist Item Component - Redesigned to match ProductCard
function WishlistItemCard({ item, onRemove, onAddToCart }: {
  item: any,
  onRemove: (productId: string) => void,
  onAddToCart: (product: any) => void
}) {
  const [isHovered, setIsHovered] = useState(false)
  const [imageLoading, setImageLoading] = useState(true)

  // Use backend computed fields directly
  const currentPrice = item.product.current_price
  const originalPrice = item.product.original_price
  const hasDiscount = item.product.has_discount
  const discountPercentage = item.product.discount_percentage
  const stockQuantity = item.product.stock
  const stockStatus = item.product.stock_status
  const isLowStock = item.product.is_low_stock
  const isOutOfStock = stockStatus === 'out_of_stock' || stockQuantity <= 0

  const primaryImage = item.product.images?.[0]?.url || item.product.main_image || '/placeholder-product.svg'
  const secondaryImage = item.product.images?.[1]?.url

  return (
    <div className="relative group h-full">
      {/* Refined outer glow - subtle and elegant */}
      <div className={cn(
        'absolute -inset-0.5 rounded-3xl opacity-0 group-hover:opacity-40 transition-all duration-700 ease-out',
        'bg-gradient-to-br from-[#ff9000]/15 via-orange-500/8 to-amber-400/10 blur-lg'
      )} />

      <Card
        className={cn(
          'relative h-full overflow-hidden backdrop-blur-sm border text-white transition-all duration-500 ease-out',
          'bg-gradient-to-br from-slate-900/70 via-gray-900/75 to-slate-800/80',
          'hover:shadow-xl hover:shadow-[#ff9000]/12 hover:-translate-y-2 hover:scale-[1.01]',
          'rounded-3xl backdrop-saturate-150 border-gray-700/40 hover:border-[#ff9000]/25',
          'before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/2 before:via-transparent before:to-white/1 before:pointer-events-none before:rounded-3xl'
        )}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <Link href={`/products/${item.product.id}`} className="block h-full">
          {/* Optimized image container with perfect proportions */}
          <div className="relative aspect-[4/3] overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 rounded-t-3xl">
            {/* Subtle shimmer effect */}
            <div className={cn(
              'absolute inset-0 -translate-x-full bg-gradient-to-r from-transparent via-white/6 to-transparent',
              'transition-transform duration-1500 ease-out',
              isHovered && 'translate-x-full'
            )} />

            {/* Featured badge */}
            {(item.product as any).featured && (
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
                (item.product as any).featured ? "left-20" : "left-3"
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
                alt={item.product.name}
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
                {/* Quick view button */}
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    'h-9 w-9 rounded-lg backdrop-blur-sm border border-white/15 shadow-md transition-all duration-250',
                    'bg-white/85 hover:bg-white text-slate-700 hover:text-blue-600',
                    'hover:scale-105',
                    'transform-gpu'
                  )}
                  onClick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    window.open(`/products/${item.product.id}`, '_blank')
                  }}
                  style={{
                    boxShadow: '0 2px 12px rgba(0, 0, 0, 0.12)'
                  }}
                >
                  <Eye className="h-3.5 w-3.5" />
                </Button>

                {/* Remove from wishlist button */}
                <Button
                  variant="ghost"
                  size="icon"
                  className={cn(
                    'h-9 w-9 rounded-lg backdrop-blur-sm border border-white/15 shadow-md transition-all duration-250',
                    'bg-white/85 hover:bg-white text-slate-700 hover:text-red-500',
                    'hover:scale-105',
                    'transform-gpu'
                  )}
                  onClick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    onRemove(item.product.id)
                  }}
                  style={{
                    boxShadow: '0 2px 12px rgba(0, 0, 0, 0.12)'
                  }}
                >
                  <Trash2 className="h-3.5 w-3.5" />
                </Button>

                {/* Add to cart button */}
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
                    onClick={(e) => {
                      e.preventDefault()
                      e.stopPropagation()
                      onAddToCart(item.product)
                    }}
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
              {item.product.category && (
                <span className="text-xs font-medium text-[#ff9000] bg-[#ff9000]/10 px-2 py-1 rounded-md border border-[#ff9000]/20">
                  {item.product.category.name}
                </span>
              )}

              {/* Compact rating display */}
              {item.product.rating_average > 0 && (
                <div className="flex items-center gap-1">
                  <div className="flex items-center gap-0.5">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={cn(
                          'h-3 w-3 transition-all duration-200',
                          i < Math.floor(item.product.rating_average)
                            ? 'text-amber-400 fill-amber-400'
                            : 'text-gray-600'
                        )}
                      />
                    ))}
                  </div>
                  <span className="text-xs text-gray-400 ml-1">
                    {item.product.rating_average.toFixed(1)}
                  </span>
                </div>
              )}
            </div>

            {/* Clean product name with optimal height */}
            <div>
              <h3 className="text-base font-semibold leading-tight line-clamp-2 text-white transition-colors duration-300 group-hover:text-[#ff9000]/90">
                {item.product.name}
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

            {/* Added date - subtle and minimal */}
            <div className="pt-2 border-t border-gray-700/30">
              <span className="text-xs text-gray-400">
                Added {new Date(item.added_at).toLocaleDateString()}
              </span>
            </div>
          </div>
        </Link>
      </Card>
    </div>
  )
}

export function WishlistPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [sortBy, setSortBy] = useState<'added_at' | 'product_name' | 'price'>('added_at')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('desc')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')

  const { data: wishlistData, isLoading, error } = useWishlist({ limit: 50 })
  const removeFromWishlistMutation = useRemoveFromWishlist()
  const clearWishlistMutation = useClearWishlist()
  const { addItem } = useCartStore()
  const { isAuthenticated } = useAuthStore()



  // Filter and sort items
  const items = Array.isArray(wishlistData?.data) ? wishlistData.data : [];
  console.log('wishlistData:', wishlistData);
  console.log('items:', items);
  const filteredItems = items.filter(item => {
    const nameMatch = item.product.name.toLowerCase().includes(searchQuery.toLowerCase());
    const categoryMatch = item.product.category && item.product.category.name
      ? item.product.category.name.toLowerCase().includes(searchQuery.toLowerCase())
      : false;
    return nameMatch || categoryMatch;
  })

  const sortedItems = [...filteredItems].sort((a, b) => {
    let comparison = 0
    switch (sortBy) {
      case 'added_at':
        comparison = new Date(a.added_at).getTime() - new Date(b.added_at).getTime()
        break
      case 'product_name':
        comparison = a.product.name.localeCompare(b.product.name)
        break
      case 'price':
        comparison = a.product.current_price - b.product.current_price
        break
    }
    return sortOrder === 'asc' ? comparison : -comparison
  })

  const handleRemoveFromWishlist = (productId: string) => {
    removeFromWishlistMutation.mutate(productId)
  }

  const handleClearWishlist = () => {
    if (window.confirm('Are you sure you want to clear your entire wishlist?')) {
      clearWishlistMutation.mutate()
    }
  }

  const handleAddToCart = (product: any) => {
    addItem(product.id)
    toast.success('Added to cart!')
  }

  if (!isAuthenticated) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-12 relative z-10">
          <div className="max-w-md mx-auto text-center">
            <div className="w-20 h-20 bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 rounded-full flex items-center justify-center mx-auto mb-6 border border-[#ff9000]/30">
              <Heart className="h-10 w-10 text-[#ff9000]" />
            </div>
            <h1 className="text-3xl font-bold mb-4">Sign In Required</h1>
            <p className="text-gray-400 mb-6">
              Please sign in to view your wishlist and save your favorite products.
            </p>
            <div className="space-y-3">
              <Button asChild className="w-full bg-[#ff9000] hover:bg-[#e67e00] text-white">
                <Link href="/auth/login">Sign In</Link>
              </Button>
              <Button asChild variant="outline" className="w-full border-gray-600 text-gray-300 hover:bg-gray-700">
                <Link href="/auth/register">Create Account</Link>
              </Button>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      <AnimatedBackground className="opacity-30" />
      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-6 relative z-10">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Button
              variant="ghost"
              size="icon"
              asChild
              className="text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg"
            >
              <Link href="/">
                <ArrowLeft className="h-5 w-5" />
              </Link>
            </Button>
            <div>
              <h1 className="text-3xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                My Wishlist
              </h1>
              <p className="text-gray-400 mt-1">
                {wishlistData?.data?.length || 0} saved items
              </p>
            </div>
          </div>

          {wishlistData && wishlistData.data && wishlistData.data.length > 0 && (
            <Button
              variant="outline"
              onClick={handleClearWishlist}
              disabled={clearWishlistMutation.isPending}
              className="border-red-600 text-red-400 hover:bg-red-600 hover:text-white"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Clear All
            </Button>
          )}
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="animate-pulse">
                <div className="bg-white/[0.03] border border-white/10 rounded-3xl overflow-hidden">
                  <div className="aspect-[4/3] bg-white/[0.03]"></div>
                  <div className="p-4 space-y-3">
                    <div className="h-4 bg-white/[0.03] rounded w-3/4"></div>
                    <div className="h-6 bg-white/[0.03] rounded w-1/2"></div>
                    <div className="h-4 bg-white/[0.03] rounded w-1/3"></div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="text-center py-12">
            <div className="w-20 h-20 bg-gradient-to-br from-red-500/20 to-red-600/20 rounded-full flex items-center justify-center mx-auto mb-6 border border-red-500/30">
              <span className="text-4xl">😞</span>
            </div>
            <h2 className="text-2xl font-bold mb-4">Failed to Load Wishlist</h2>
            <p className="text-gray-400 mb-6">
              There was an error loading your wishlist. Please try again.
            </p>
            <Button onClick={() => window.location.reload()} className="bg-[#ff9000] hover:bg-[#e67e00]">
              Try Again
            </Button>
          </div>
        ) : !wishlistData || !wishlistData.data || wishlistData.data.length === 0 ? (
          <div className="text-center py-12">
            <div className="w-20 h-20 bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 rounded-full flex items-center justify-center mx-auto mb-6 border border-[#ff9000]/30">
              <Heart className="h-10 w-10 text-[#ff9000]" />
            </div>
            <h2 className="text-2xl font-bold mb-4">Your Wishlist is Empty</h2>
            <p className="text-gray-400 mb-6">
              Start adding products you love to your wishlist and never lose track of them.
            </p>
            <Button asChild className="bg-[#ff9000] hover:bg-[#e67e00] text-white">
              <Link href="/products">
                <Sparkles className="h-4 w-4 mr-2" />
                Discover Products
              </Link>
            </Button>
          </div>
        ) : (
          <>
            {/* Filters and Controls */}
            <div className="flex flex-col sm:flex-row gap-4 mb-6">
              {/* Search */}
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search your wishlist..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 bg-gray-800/50 border-gray-600 text-white placeholder:text-gray-400 focus:border-[#ff9000]"
                />
              </div>

              {/* Sort */}
              <div className="flex gap-2">
                <select
                  value={`${sortBy}-${sortOrder}`}
                  onChange={(e) => {
                    const [newSortBy, newSortOrder] = e.target.value.split('-') as [typeof sortBy, typeof sortOrder]
                    setSortBy(newSortBy)
                    setSortOrder(newSortOrder)
                  }}
                  className="px-3 py-2 bg-gray-800/50 border border-gray-600 rounded-lg text-white focus:border-[#ff9000] focus:outline-none"
                >
                  <option value="added_at-desc">Newest First</option>
                  <option value="added_at-asc">Oldest First</option>
                  <option value="product_name-asc">Name A-Z</option>
                  <option value="product_name-desc">Name Z-A</option>
                  <option value="price-asc">Price Low-High</option>
                  <option value="price-desc">Price High-Low</option>
                </select>

                {/* View Mode Toggle */}
                <div className="flex border border-gray-600 rounded-lg overflow-hidden">
                  <Button
                    variant={viewMode === 'grid' ? 'default' : 'ghost'}
                    size="icon"
                    onClick={() => setViewMode('grid')}
                    className="rounded-none"
                  >
                    <Grid3X3 className="h-4 w-4" />
                  </Button>
                  <Button
                    variant={viewMode === 'list' ? 'default' : 'ghost'}
                    size="icon"
                    onClick={() => setViewMode('list')}
                    className="rounded-none"
                  >
                    <List className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>

            {/* Wishlist Items */}
            <div className={cn(
              'grid gap-6',
              viewMode === 'grid'
                ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4'
                : 'grid-cols-1 max-w-4xl mx-auto'
            )}>
              {sortedItems.map((item) => (
                <WishlistItemCard
                  key={item.id}
                  item={item}
                  onRemove={handleRemoveFromWishlist}
                  onAddToCart={handleAddToCart}
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}
