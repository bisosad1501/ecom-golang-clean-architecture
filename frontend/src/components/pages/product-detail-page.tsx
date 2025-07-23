'use client'

import { useState } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { 
  Star, 
  Heart, 
  Share2, 
  ShoppingCart, 
  Minus, 
  Plus,
  Truck,
  Shield,
  RotateCcw,
  ChevronLeft,
  ChevronRight,
  ArrowLeft,
  Eye,
  CheckCircle,
  Gift
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProductCard } from '@/components/products/product-card'

import { AnimatedBackground } from '@/components/ui/animated-background'
import { useProduct, useRelatedProducts, useAddToWishlist } from '@/hooks/use-products'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { cn } from '@/lib/utils'
import { useProductReviews, useProductRatingSummary } from '@/hooks/use-reviews'
import { toast } from 'sonner'

interface ProductDetailPageProps {
  productId: string
}

// Simple price formatter
const formatPrice = (price: number): string => {
  return `$${price.toFixed(2)}`
}

export function ProductDetailPage({ productId }: ProductDetailPageProps) {
  const [selectedImageIndex, setSelectedImageIndex] = useState(0)
  const [quantity, setQuantity] = useState(1)
  const [activeTab, setActiveTab] = useState<'description' | 'reviews' | 'shipping'>('description')

  const { data: product, isLoading, error } = useProduct(productId)
  const { data: relatedProducts } = useRelatedProducts(productId, 4)
  const { addItem, isLoading: cartLoading } = useCartStore()
  const { isAuthenticated } = useAuthStore()
  const addToWishlistMutation = useAddToWishlist()

  // Import reviews hooks
  const { data: reviewsData } = useProductReviews(productId, { limit: 5, sort_by: 'created_at', sort_order: 'desc' })
  const { data: ratingSummary } = useProductRatingSummary(productId)

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-6 relative z-10">
          <div className="animate-pulse">
            {/* Header skeleton */}
            <div className="flex items-center justify-between mb-6">
              <div className="h-8 bg-white/8 rounded-lg w-32"></div>
              <div className="h-4 bg-white/8 rounded w-48"></div>
            </div>
            
            <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 mb-8">
              {/* Image skeleton - 6 columns */}
              <div className="lg:col-span-6 space-y-3">
                <div className="aspect-square bg-white/[0.03] border border-white/10 rounded-xl"></div>
                <div className="grid grid-cols-6 gap-2">
                  {[...Array(6)].map((_, i) => (
                    <div key={i} className="aspect-square bg-white/[0.03] border border-white/10 rounded-lg"></div>
                  ))}
                </div>
              </div>
              {/* Content skeleton - 6 columns */}
              <div className="lg:col-span-6 space-y-4">
                <div className="h-4 bg-white/8 rounded w-20"></div>
                <div className="h-6 bg-white/8 rounded w-3/4"></div>
                <div className="h-4 bg-white/8 rounded w-1/2"></div>
                <div className="bg-white/[0.03] border border-white/10 rounded-xl p-4">
                  <div className="h-6 bg-white/8 rounded w-1/3"></div>
                </div>
                <div className="bg-white/[0.03] border border-white/10 rounded-xl p-4">
                  <div className="h-16 bg-white/8 rounded w-full"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden flex items-center justify-center">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 relative z-10">
          <Card className="p-12 text-center max-w-lg mx-auto bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl">
            <div className="w-20 h-20 bg-gradient-to-br from-red-500/20 to-red-600/20 rounded-full flex items-center justify-center mx-auto mb-6 border border-red-500/30">
              <span className="text-4xl">😞</span>
            </div>
            <h2 className="text-3xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-4">
              Product Not Found
            </h2>
            <p className="text-gray-300 mb-8 leading-relaxed">
              The product you're looking for doesn't exist or has been removed.
            </p>
            <Button asChild className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white px-8 py-3 rounded-xl font-semibold shadow-lg shadow-[#ff9000]/25 transition-all duration-300 hover:scale-105">
              <Link href="/products">Browse Products</Link>
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  // Extract product data from API response
  const productData = (product as any).data || product;

  // Handle images - use main_image if images is null or empty
  let images: Array<{ id: string; url: string; alt_text?: string; position?: number }> = [];
  if (Array.isArray(productData.images) && productData.images.length > 0) {
    images = productData.images;
  } else if (productData.main_image && typeof productData.main_image === 'string' && productData.main_image.trim() !== '') {
    images = [{ id: 'main', url: productData.main_image, alt_text: productData.name, position: 0 }];
  } else {
    images = [{ id: 'placeholder', url: '/placeholder-product.svg', alt_text: productData.name, position: 0 }];
  }



  // Use backend computed fields directly
  const currentPrice = typeof productData.current_price === 'number' ? productData.current_price : 0;
  const originalPrice = typeof productData.original_price === 'number' ? productData.original_price : null;
  const hasDiscount = productData.has_discount
  const isOnSale = productData.is_on_sale
  const discountPercentage = productData.discount_percentage
  const stockQuantity = productData.stock
  const stockStatus = productData.stock_status
  const isLowStock = productData.is_low_stock
  const featured = productData.featured

  const isOutOfStock = stockStatus === 'out_of_stock' || stockQuantity <= 0

  // Handle rating data - use reviews data if available, fallback to product data
  const ratingAverage = ratingSummary?.average_rating || productData.rating_average || 0
  const ratingCount = ratingSummary?.total_reviews || productData.rating_count || 0
  const actualReviews = reviewsData?.reviews || []

  const handleAddToCart = async () => {
    if (isOutOfStock) return
    
    try {
      await addItem(productData.id, quantity)
    } catch (error) {
      toast.error('Failed to add to cart')
    }
  }

  const handleAddToWishlist = async () => {
    if (!isAuthenticated) {
      toast.error('Please sign in to add to wishlist')
      return
    }

    try {
      await addToWishlistMutation.mutateAsync(productData.id)
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: productData.name,
          text: productData.short_description || productData.description,
          url: window.location.href,
        })
      } catch (error) {
        // User cancelled sharing
      }
    } else {
      // Fallback: copy to clipboard
      navigator.clipboard.writeText(window.location.href)
      toast.success('Product link copied to clipboard!')
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      <AnimatedBackground className="opacity-30" />

      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-8 relative z-10">





        {/* Main Product Section */}
        <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl overflow-hidden mb-8 shadow-2xl shadow-black/20">
          <div className="grid grid-cols-1 lg:grid-cols-2">
            {/* Product Images */}
            <div className="p-8 border-b lg:border-b-0 lg:border-r border-white/15">
              {/* Main Image */}
              <div className="relative group mb-4">
                <div className="aspect-square overflow-hidden rounded-lg bg-gradient-to-br from-white/[0.08] to-white/[0.02] border border-white/20 shadow-lg cursor-zoom-in relative">
                  {/* Zoom overlay */}
                  <div className="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition-all duration-300 z-20 flex items-center justify-center">
                    <div className="w-10 h-10 bg-white/20 backdrop-blur-sm rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-300 scale-75 group-hover:scale-100">
                      <Eye className="w-5 h-5 text-white" />
                    </div>
                  </div>
                  {/* Enhanced Badges */}
                  {featured && (
                    <Badge className="absolute top-3 left-3 z-10 bg-gradient-to-r from-purple-500 to-pink-500 text-white px-2 py-1 text-xs font-bold">
                      Featured
                    </Badge>
                  )}
                  {hasDiscount && discountPercentage > 0 && (
                    <Badge className={cn(
                      "absolute top-3 z-10 bg-[#ff9000] text-white px-2 py-1 text-xs font-bold",
                      featured ? "left-20" : "left-3"
                    )}>
                      -{Math.round(discountPercentage)}%
                    </Badge>
                  )}
                  {isOutOfStock && (
                    <Badge className="absolute top-3 right-3 z-10 bg-red-500 text-white px-2 py-1 text-xs">
                      Sold Out
                    </Badge>
                  )}
                  {!isOutOfStock && isLowStock && (
                    <Badge className="absolute top-3 right-3 z-10 bg-amber-500 text-white px-2 py-1 text-xs">
                      Low Stock
                    </Badge>
                  )}

                  <Image
                    src={images[selectedImageIndex]?.url || '/placeholder-product.svg'}
                    alt={
                      typeof images[selectedImageIndex]?.alt_text === 'string' && images[selectedImageIndex]?.alt_text.trim() !== ''
                        ? images[selectedImageIndex]?.alt_text
                        : (productData.name || 'Product image')
                    }
                    fill
                    className="object-cover transition-transform duration-700 ease-out group-hover:scale-110"
                    priority
                  />

                  {/* Navigation Arrows */}
                  {images.length > 1 && (
                    <>
                      <Button
                        variant="secondary"
                        size="icon"
                        className="absolute left-2 top-1/2 transform -translate-y-1/2 bg-black/40 hover:bg-black/60 text-white border-0 w-8 h-8 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={() => setSelectedImageIndex(prev => prev === 0 ? images.length - 1 : prev - 1)}
                      >
                        <ChevronLeft className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="secondary"
                        size="icon"
                        className="absolute right-2 top-1/2 transform -translate-y-1/2 bg-black/40 hover:bg-black/60 text-white border-0 w-8 h-8 opacity-0 group-hover:opacity-100 transition-opacity"
                        onClick={() => setSelectedImageIndex(prev => prev === images.length - 1 ? 0 : prev + 1)}
                      >
                        <ChevronRight className="h-4 w-4" />
                      </Button>
                    </>
                  )}
                </div>
              </div>

              {/* Thumbnail Strip */}
              {images.length > 1 && (
                <div className="flex space-x-2 overflow-x-auto pb-2">
                  {images.slice(0, 5).map((image, index) => (
                    <button
                      key={image.id}
                      onClick={() => setSelectedImageIndex(index)}
                      className={cn(
                        'flex-shrink-0 w-16 h-16 rounded border-2 overflow-hidden transition-all',
                        selectedImageIndex === index 
                          ? 'border-[#ff9000]' 
                          : 'border-white/20 hover:border-[#ff9000]/50'
                      )}
                    >
                      <Image src={image.url} alt={image.alt_text || productData.name || `Product image ${index + 1}`} width={64} height={64} className="object-cover w-full h-full" />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Product Info */}
            <div className="p-8 space-y-6">
              {/* Category and Brand */}
              <div className="flex items-center gap-3 flex-wrap">
                {productData.category && (
                  <Link href={`/products?category=${productData.category.id}`}>
                    <Badge className="bg-gradient-to-r from-[#ff9000]/20 to-orange-600/20 text-[#ff9000] border border-[#ff9000]/40 hover:from-[#ff9000]/30 hover:to-orange-600/30 transition-all duration-300 font-semibold px-3 py-1.5 text-sm rounded-full">
                      {productData.category.name}
                    </Badge>
                  </Link>
                )}
                {productData.brand && (
                  <Badge className="bg-gradient-to-r from-blue-500/20 to-cyan-500/20 text-blue-400 border border-blue-500/40 hover:from-blue-500/30 hover:to-cyan-500/30 transition-all duration-300 font-semibold px-3 py-1.5 text-sm rounded-full">
                    {productData.brand.name}
                  </Badge>
                )}
              </div>

              {/* Title */}
              <h1 className="text-3xl lg:text-4xl font-bold text-white leading-tight tracking-tight bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                {productData.name}
              </h1>

              {/* Rating */}
              {ratingAverage > 0 && (
                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-1">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={cn(
                          'h-4 w-4',
                          i < Math.floor(ratingAverage)
                            ? 'text-yellow-400 fill-current'
                            : 'text-gray-600'
                        )}
                      />
                    ))}
                    <span className="text-sm font-medium text-yellow-400 ml-1">
                      {ratingAverage.toFixed(1)}
                    </span>
                  </div>
                  <span className="text-sm text-gray-400">({ratingCount} reviews)</span>
                </div>
              )}

              {/* Enhanced Price Display */}
              <div className="border border-white/30 rounded-2xl p-6 bg-gradient-to-br from-white/[0.12] to-white/[0.04] shadow-xl backdrop-blur-sm">
                <div className="flex items-center gap-4 mb-3">
                  <span className="text-4xl font-bold bg-gradient-to-r from-[#ff9000] to-orange-600 bg-clip-text text-transparent tracking-tight">
                    {formatPrice(currentPrice)}
                  </span>
                  {hasDiscount && originalPrice && (
                    <div className="flex items-center gap-2">
                      <span className="text-lg text-gray-400 line-through">
                        {formatPrice(originalPrice)}
                      </span>
                      {discountPercentage > 0 && (
                        <Badge className="bg-gradient-to-r from-[#ff9000] to-orange-600 text-white text-xs px-2 py-1 font-bold shadow-lg">
                          -{Math.round(discountPercentage)}% OFF
                        </Badge>
                      )}
                    </div>
                  )}
                </div>
                {hasDiscount && discountPercentage > 0 && originalPrice && (
                  <p className="text-green-400 text-sm font-semibold flex items-center gap-2">
                    <span className="text-green-300">💰</span>
                    You save {formatPrice(originalPrice - currentPrice)}
                  </p>
                )}
                {isOnSale && (
                  <p className="text-[#ff9000] text-sm font-semibold mt-2 flex items-center gap-2">
                    <span className="text-orange-400">🔥</span>
                    Limited time sale!
                  </p>
                )}
              </div>

              {/* Enhanced Stock Status */}
              <div className="flex items-center gap-3">
                {isOutOfStock ? (
                  <Badge className="bg-red-500/20 text-red-300 border border-red-500/40 px-4 py-2 text-sm font-semibold">
                    <span className="mr-2">❌</span>
                    Sold Out
                  </Badge>
                ) : stockStatus === 'on_backorder' ? (
                  <Badge className="bg-blue-500/20 text-blue-300 border border-blue-500/40 px-4 py-2 text-sm font-semibold">
                    <span className="mr-2">📦</span>
                    Available on Backorder
                  </Badge>
                ) : isLowStock ? (
                  <Badge className="bg-amber-500/20 text-amber-300 border border-amber-500/40 px-4 py-2 text-sm font-semibold">
                    <span className="mr-2">⚠️</span>
                    Low Stock - Only {stockQuantity} left
                  </Badge>
                ) : stockQuantity <= 10 ? (
                  <Badge className="bg-yellow-500/20 text-yellow-300 border border-yellow-500/40 px-4 py-2 text-sm font-semibold">
                    <span className="mr-2">⏰</span>
                    Limited Stock - {stockQuantity} available
                  </Badge>
                ) : (
                  <Badge className="bg-green-500/20 text-green-300 border border-green-500/40 px-4 py-2 text-sm font-semibold">
                    <span className="mr-2">✅</span>
                    In Stock ({stockQuantity} available)
                  </Badge>
                )}
              </div>

              {/* Short Description */}
              {productData.short_description && (
                <div className="border-t border-white/15 pt-6">
                  <p className="text-gray-200 text-base leading-relaxed">
                    {productData.short_description}
                  </p>
                </div>
              )}



              {/* Product Tags */}
              {productData.tags && productData.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {productData.tags.slice(0, 5).map((tag: any) => (
                    <Badge
                      key={tag.id}
                      variant="outline"
                      className="border-gray-500/50 text-gray-300 hover:border-[#ff9000]/50 hover:text-[#ff9000] transition-colors text-xs px-2 py-1"
                    >
                      #{tag.name}
                    </Badge>
                  ))}
                </div>
              )}

              {/* Purchase Actions */}
              {!isOutOfStock && (
                <div className="border-t border-white/15 pt-4 space-y-4">
                  {/* Quantity */}
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-semibold text-white">Quantity:</span>
                    <div className="flex items-center border border-white/30 rounded-lg overflow-hidden bg-white/[0.05] shadow-lg">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.max(1, quantity - 1))}
                        disabled={quantity <= 1}
                        className="px-3 py-2 h-10 border-r border-white/30 hover:bg-white/10 text-white"
                      >
                        <Minus className="h-3 w-3" />
                      </Button>
                      <span className="px-4 py-2 min-w-[50px] text-center font-bold text-base text-white">
                        {quantity}
                      </span>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.min(stockQuantity, quantity + 1))}
                        disabled={quantity >= stockQuantity}
                        className="px-3 py-2 h-10 border-l border-white/30 hover:bg-white/10 text-white"
                      >
                        <Plus className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="space-y-4">
                    <Button
                      className="w-full bg-gradient-to-r from-[#ff9000] via-orange-500 to-orange-600 hover:from-orange-600 hover:via-orange-600 hover:to-orange-700 text-white font-bold py-4 text-lg transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:shadow-[#ff9000]/30 rounded-xl shadow-lg"
                      onClick={handleAddToCart}
                      disabled={cartLoading}
                    >
                      <ShoppingCart className="mr-3 h-5 w-5" />
                      {cartLoading ? 'Adding to Cart...' : 'Add to Cart'}
                    </Button>

                    <div className="grid grid-cols-2 gap-3">
                      {isAuthenticated && (
                        <Button
                          variant="outline"
                          onClick={handleAddToWishlist}
                          disabled={addToWishlistMutation.isPending}
                          className="border-white/40 hover:bg-gradient-to-r hover:from-pink-500/10 hover:to-red-500/10 transition-all duration-300 py-3 text-white font-semibold rounded-xl hover:border-pink-500/60 hover:text-pink-400 backdrop-blur-sm"
                        >
                          <Heart className="h-4 w-4 mr-2" />
                          Wishlist
                        </Button>
                      )}
                      <Button
                        variant="outline"
                        onClick={handleShare}
                        className="border-white/40 hover:bg-gradient-to-r hover:from-blue-500/10 hover:to-cyan-500/10 transition-all duration-300 py-3 text-white font-semibold rounded-xl hover:border-blue-500/60 hover:text-blue-400 backdrop-blur-sm"
                      >
                        <Share2 className="h-4 w-4 mr-2" />
                        Share
                      </Button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Features & Benefits */}
        <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 mb-8 shadow-xl shadow-black/10">
          <h2 className="text-2xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-8 text-center">Why Choose This Product</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center group cursor-pointer">
              <div className="w-20 h-20 bg-gradient-to-br from-[#ff9000] via-orange-500 to-orange-600 rounded-2xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 group-hover:rotate-3 transition-all duration-300 shadow-xl shadow-orange-500/25">
                <Truck className="h-10 w-10 text-white" />
              </div>
              <h3 className="font-bold text-white text-lg mb-3">Free Shipping</h3>
              <p className="text-gray-300 text-sm leading-relaxed">Fast and free delivery on orders over $50. Get your products delivered to your doorstep.</p>
            </div>
            <div className="text-center group cursor-pointer">
              <div className="w-20 h-20 bg-gradient-to-br from-green-500 via-emerald-500 to-emerald-600 rounded-2xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 group-hover:rotate-3 transition-all duration-300 shadow-xl shadow-green-500/25">
                <Shield className="h-10 w-10 text-white" />
              </div>
              <h3 className="font-bold text-white text-lg mb-3">Quality Guarantee</h3>
              <p className="text-gray-300 text-sm leading-relaxed">100% authentic products with quality assurance. Your satisfaction is our priority.</p>
            </div>
            <div className="text-center group cursor-pointer">
              <div className="w-20 h-20 bg-gradient-to-br from-blue-500 via-cyan-500 to-cyan-600 rounded-2xl flex items-center justify-center mx-auto mb-4 group-hover:scale-110 group-hover:rotate-3 transition-all duration-300 shadow-xl shadow-blue-500/25">
                <RotateCcw className="h-10 w-10 text-white" />
              </div>
              <h3 className="font-bold text-white text-lg mb-3">Easy Returns</h3>
              <p className="text-gray-300 text-sm leading-relaxed">Hassle-free 30-day return policy. Not satisfied? Return it with no questions asked.</p>
            </div>
          </div>
        </div>

        {/* Product Information Tabs */}
        <div className="bg-white/[0.03] backdrop-blur-sm border border-white/10 rounded-xl mb-8 overflow-hidden">
          {/* Tab Navigation */}
          <div className="border-b border-white/10">
            <div className="flex">
              {[
                { key: 'description', label: 'Description', icon: Eye },
                { key: 'reviews', label: 'Reviews', icon: Star },
                { key: 'shipping', label: 'Shipping & Returns', icon: Truck },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as any)}
                  className={cn(
                    'flex items-center space-x-3 px-8 py-4 text-sm font-medium transition-all duration-200 border-b-2 relative',
                    activeTab === tab.key
                      ? 'text-[#ff9000] border-[#ff9000] bg-white/[0.05]'
                      : 'text-gray-400 border-transparent hover:text-white hover:bg-white/[0.02]'
                  )}
                >
                  <tab.icon className="h-5 w-5" />
                  <span>{tab.label}</span>
                  {tab.key === 'reviews' && ratingCount > 0 && (
                    <Badge className="bg-[#ff9000] text-white text-xs px-2 py-1 font-semibold">
                      {ratingCount}
                    </Badge>
                  )}
                </button>
              ))}
            </div>
          </div>

          {/* Tab Content */}
          <div className="p-8">
            {activeTab === 'description' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-xl font-bold text-white mb-4">Product Description</h3>
                  <div className="text-gray-300 leading-relaxed text-base">
                    {productData.description ? (
                      <p className="whitespace-pre-line">{productData.description}</p>
                    ) : (
                      <p className="text-gray-400 italic">No detailed description available for this product.</p>
                    )}
                  </div>
                </div>

                {/* Product Specifications */}
                <div className="border-t border-white/10 pt-6">
                  <h4 className="text-lg font-semibold text-white mb-4">Specifications</h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-3">
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">SKU</span>
                        <span className="text-white font-medium">{productData.sku}</span>
                      </div>
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">Category</span>
                        <span className="text-white font-medium">{productData.category?.name || 'N/A'}</span>
                      </div>
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">Brand</span>
                        <span className="text-white font-medium">{productData.brand?.name || 'N/A'}</span>
                      </div>
                    </div>
                    <div className="space-y-3">
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">Weight</span>
                        <span className="text-white font-medium">{productData.weight ? `${productData.weight}kg` : 'N/A'}</span>
                      </div>
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">Shipping Required</span>
                        <span className="text-white font-medium">{productData.requires_shipping ? 'Yes' : 'No'}</span>
                      </div>
                      <div className="flex justify-between py-2 border-b border-white/5">
                        <span className="text-gray-400">Digital Product</span>
                        <span className="text-white font-medium">{productData.is_digital ? 'Yes' : 'No'}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            )}



            {activeTab === 'reviews' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-xl font-bold text-white mb-6">Customer Reviews</h3>

                  {/* Review Summary */}
                  <div className="bg-white/[0.02] border border-white/10 rounded-xl p-6 mb-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="text-center">
                        <div className="text-3xl font-bold text-[#ff9000] mb-2">
                          {ratingAverage > 0 ? ratingAverage.toFixed(1) : '0.0'}
                        </div>
                        <div className="flex items-center justify-center gap-1 mb-2">
                          {[...Array(5)].map((_, i) => (
                            <Star
                              key={i}
                              className={cn(
                                'h-4 w-4',
                                i < Math.floor(ratingAverage)
                                  ? 'text-yellow-400 fill-current'
                                  : 'text-gray-600'
                              )}
                            />
                          ))}
                        </div>
                        <p className="text-gray-400 text-sm">
                          Based on {ratingCount} {ratingCount === 1 ? 'review' : 'reviews'}
                        </p>
                      </div>

                      <div className="space-y-2">
                        {[5, 4, 3, 2, 1].map((rating) => {
                          const count = ratingSummary?.rating_counts?.[rating.toString()] || 0
                          const percentage = ratingCount > 0 ? (count / ratingCount) * 100 : 0
                          return (
                            <div key={rating} className="flex items-center gap-3">
                              <span className="text-sm text-gray-400 w-8">{rating}★</span>
                              <div className="flex-1 bg-gray-700 rounded-full h-2">
                                <div
                                  className="bg-[#ff9000] h-2 rounded-full transition-all duration-300"
                                  style={{ width: `${percentage}%` }}
                                ></div>
                              </div>
                              <span className="text-sm text-gray-400 w-8">{count}</span>
                            </div>
                          )
                        })}
                      </div>
                    </div>
                  </div>

                  {/* Actual Reviews */}
                  {actualReviews.length > 0 ? (
                    <div className="space-y-4">
                      {actualReviews.slice(0, 3).map((review) => (
                        <div key={review.id} className="bg-white/[0.02] border border-white/10 rounded-xl p-4">
                          <div className="flex items-start justify-between mb-3">
                            <div>
                              <div className="flex items-center gap-2 mb-1">
                                <span className="font-semibold text-white">
                                  {review.user.first_name} {review.user.last_name.charAt(0)}.
                                </span>
                                <div className="flex items-center gap-1">
                                  {[...Array(5)].map((_, i) => (
                                    <Star
                                      key={i}
                                      className={cn(
                                        'h-3 w-3',
                                        i < review.rating
                                          ? 'text-yellow-400 fill-current'
                                          : 'text-gray-600'
                                      )}
                                    />
                                  ))}
                                </div>
                                {review.is_verified && (
                                  <span className="text-xs bg-green-500/20 text-green-400 px-2 py-0.5 rounded">
                                    Verified
                                  </span>
                                )}
                              </div>
                              <p className="text-xs text-gray-400">
                                {new Date(review.created_at).toLocaleDateString()}
                              </p>
                            </div>
                          </div>
                          {review.title && (
                            <h4 className="font-semibold text-white text-sm mb-2">{review.title}</h4>
                          )}
                          <p className="text-gray-300 text-sm">{review.comment}</p>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-8">
                      <div className="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-3">
                        <Star className="h-6 w-6 text-gray-500" />
                      </div>
                      <h4 className="text-lg font-semibold text-white mb-2">No Reviews Yet</h4>
                      <p className="text-gray-400 mb-4">Be the first to share your experience!</p>
                      <Button className="bg-[#ff9000] hover:bg-orange-600 text-white">
                        Write a Review
                      </Button>
                    </div>
                  )}
                </div>
              </div>
            )}

            {activeTab === 'shipping' && (
              <div className="space-y-8">
                <h3 className="text-xl font-bold text-white mb-6">Shipping & Returns</h3>

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                  {/* Shipping Information */}
                  <div className="bg-white/[0.02] border border-white/10 rounded-xl p-6">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center">
                        <Truck className="h-6 w-6 text-blue-400" />
                      </div>
                      <h4 className="text-lg font-semibold text-white">Shipping Options</h4>
                    </div>

                    <div className="space-y-4">
                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <CheckCircle className="h-5 w-5 text-green-400 mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">Free Standard Shipping</h5>
                          <p className="text-gray-400 text-sm">3-5 business days • Orders over $50</p>
                        </div>
                      </div>

                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <Truck className="h-5 w-5 text-blue-400 mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">Express Delivery</h5>
                          <p className="text-gray-400 text-sm">1-2 business days • $9.99</p>
                        </div>
                      </div>

                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <Gift className="h-5 w-5 text-[#ff9000] mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">Same Day Delivery</h5>
                          <p className="text-gray-400 text-sm">Available in select cities • $19.99</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Return Policy */}
                  <div className="bg-white/[0.02] border border-white/10 rounded-xl p-6">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center">
                        <RotateCcw className="h-6 w-6 text-green-400" />
                      </div>
                      <h4 className="text-lg font-semibold text-white">Return Policy</h4>
                    </div>

                    <div className="space-y-4">
                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <CheckCircle className="h-5 w-5 text-green-400 mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">30-Day Return Window</h5>
                          <p className="text-gray-400 text-sm">Return items within 30 days of delivery</p>
                        </div>
                      </div>

                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <Shield className="h-5 w-5 text-blue-400 mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">Original Condition Required</h5>
                          <p className="text-gray-400 text-sm">Items must be unused and in original packaging</p>
                        </div>
                      </div>

                      <div className="flex items-start gap-4 p-4 bg-white/[0.02] rounded-lg border border-white/5">
                        <Truck className="h-5 w-5 text-[#ff9000] mt-0.5" />
                        <div>
                          <h5 className="font-medium text-white mb-1">Free Return Shipping</h5>
                          <p className="text-gray-400 text-sm">We'll cover the return shipping costs</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Additional Information */}
                <div className="bg-gradient-to-r from-[#ff9000]/10 to-orange-600/10 border border-[#ff9000]/20 rounded-xl p-6">
                  <h4 className="text-lg font-semibold text-white mb-3 flex items-center gap-2">
                    <Shield className="h-5 w-5 text-[#ff9000]" />
                    Purchase Protection
                  </h4>
                  <p className="text-gray-300 text-sm leading-relaxed">
                    Your purchase is protected by our buyer protection program. If you're not completely satisfied with your order,
                    we'll make it right with a full refund or replacement.
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Related Products */}
        {relatedProducts && Array.isArray(relatedProducts.data) && relatedProducts.data.length > 0 && (
          <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.03] backdrop-blur-sm border border-white/20 rounded-2xl shadow-lg">
            <div className="p-6 border-b border-white/15">
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-xl font-bold text-white mb-1">You Might Also Like</h2>
                  <p className="text-gray-300 text-sm">Discover similar products</p>
                </div>
                <Badge className="bg-[#ff9000]/20 text-[#ff9000] border border-[#ff9000]/30 px-3 py-1 font-semibold text-sm">
                  {relatedProducts.data.length} Products
                </Badge>
              </div>
            </div>
            <div className="p-6">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {relatedProducts.data.map((relatedProduct: any, index: number) => (
                  <div
                    key={relatedProduct.id}
                    className="animate-in fade-in duration-500"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <ProductCard
                      product={relatedProduct}
                      className="transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:shadow-[#ff9000]/15"
                    />
                  </div>
                ))}
              </div>
              <div className="text-center mt-8">
                <Button
                  asChild
                  variant="outline"
                  className="border-white/30 hover:bg-white/10 transition-all duration-200 py-3 px-8 text-white font-semibold rounded-xl hover:border-[#ff9000]/50 hover:text-[#ff9000]"
                >
                  <Link href="/products">
                    View All Products
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
