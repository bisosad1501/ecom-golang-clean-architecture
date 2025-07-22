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

// Wishlist Item Component
function WishlistItemCard({ item, onRemove, onAddToCart }: {
  item: any,
  onRemove: (productId: string) => void,
  onAddToCart: (product: any) => void
}) {
  // const { data: ratingSummary } = useProductRatingSummary(item.product.id)
  
  // Use backend computed fields directly
  const currentPrice = item.product.current_price
  const originalPrice = item.product.original_price
  const hasDiscount = item.product.has_discount
  const discountPercentage = item.product.discount_percentage

  return (
    <div className="relative group">
      {/* Refined outer glow */}
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
            {/* Product Image */}
            <div className="relative flex-shrink-0">
              <div className="relative w-32 h-32 rounded-xl overflow-hidden bg-gradient-to-br from-gray-100 to-gray-200 shadow-lg">
                {/* Discount badge */}
                {hasDiscount && discountPercentage > 0 && (
                  <div className="absolute top-2 left-2 z-20">
                    <span className="text-sm font-bold text-white bg-[#ff9000] px-2.5 py-1 rounded-md shadow-md">
                      -{Math.round(discountPercentage)}%
                    </span>
                  </div>
                )}

                <Image
                  src={item.product.images?.[0]?.url || '/placeholder-product.svg'}
                  alt={item.product.name}
                  fill
                  className="object-cover transition-all duration-300 ease-out group-hover:scale-105"
                  sizes="(max-width: 128px) 100vw, 128px"
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

            {/* Product Info */}
            <div className="flex-1 min-w-0 space-y-2">
              {/* Category and Remove Button */}
              <div className="flex items-start justify-between gap-2">
                <div className="flex items-center gap-2">
                  {item.product.category && (
                    <span className="text-sm font-semibold text-[#ff9000] bg-[#ff9000]/15 px-3 py-1.5 rounded-lg border border-[#ff9000]/40 shadow-sm">
                      {item.product.category.name}
                    </span>
                  )}
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => onRemove(item.product.id)}
                  className="h-8 w-8 text-gray-400 hover:text-red-400 hover:bg-red-500/10 rounded-lg transition-colors"
                >
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>

              {/* Product Name */}
              <Link 
                href={`/products/${item.product.id}`}
                className="block group/link"
              >
                <h3 className="text-lg font-bold text-white group-hover/link:text-[#ff9000] transition-colors line-clamp-2 leading-tight">
                  {item.product.name}
                </h3>
              </Link>

              {/* Rating - Temporarily disabled */}
              {/* {ratingSummary && ratingSummary.total_reviews > 0 && (
                <CompactReviewSummary
                  rating={ratingSummary.average_rating}
                  reviewCount={ratingSummary.total_reviews}
                  className="text-sm"
                />
              )} */}

              {/* Price and Stock */}
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="text-2xl font-bold text-[#ff9000]">
                    {formatPrice(currentPrice)}
                  </span>
                  {hasDiscount && originalPrice > currentPrice && (
                    <span className="text-lg text-gray-400 line-through">
                      {formatPrice(originalPrice)}
                    </span>
                  )}
                </div>
                
                <div className="flex items-center gap-2">
                  {item.product.stock > 0 ? (
                    <Badge variant="outline" className="text-green-400 border-green-400/50 bg-green-400/10">
                      In Stock
                    </Badge>
                  ) : (
                    <Badge variant="outline" className="text-red-400 border-red-400/50 bg-red-400/10">
                      Out of Stock
                    </Badge>
                  )}
                </div>
              </div>

              {/* Added Date */}
              <div className="text-sm text-gray-400">
                Added {new Date(item.added_at).toLocaleDateString()}
              </div>

              {/* Action Buttons */}
              <div className="flex gap-2 pt-2">
                <Button
                  onClick={() => onAddToCart(item.product)}
                  disabled={item.product.stock === 0}
                  className="flex-1 bg-[#ff9000] hover:bg-[#e67e00] text-white font-semibold py-2 px-4 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <ShoppingCart className="h-4 w-4 mr-2" />
                  Add to Cart
                </Button>
                <Button
                  variant="outline"
                  asChild
                  className="px-4 border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white rounded-lg"
                >
                  <Link href={`/products/${item.product.id}`}>
                    <Eye className="h-4 w-4 mr-2" />
                    View
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </CardContent>
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
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="animate-pulse">
                <div className="bg-white/[0.03] border border-white/10 rounded-2xl p-4">
                  <div className="flex gap-4">
                    <div className="w-32 h-32 bg-white/[0.03] rounded-xl"></div>
                    <div className="flex-1 space-y-3">
                      <div className="h-4 bg-white/[0.03] rounded w-3/4"></div>
                      <div className="h-6 bg-white/[0.03] rounded w-1/2"></div>
                      <div className="h-4 bg-white/[0.03] rounded w-1/3"></div>
                      <div className="h-10 bg-white/[0.03] rounded"></div>
                    </div>
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
                ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3' 
                : 'grid-cols-1'
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
