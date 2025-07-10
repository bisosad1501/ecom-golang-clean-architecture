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
import { ProductReviewsSection } from '@/components/reviews/product-reviews-section'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useProduct, useRelatedProducts, useAddToWishlist } from '@/hooks/use-products'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'

interface ProductDetailPageProps {
  productId: string
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

  const images = product.images || []
  
  // Enhanced pricing logic using new backend fields
  const currentPrice = (product as any).current_price || (product as any).pricing?.price || (product as any).price || 0
  const originalPrice = (product as any).price || (product as any).pricing?.price || 0
  const salePrice = (product as any).sale_price
  const isOnSale = (product as any).is_on_sale || false
  const hasDiscount = (product as any).has_discount || isOnSale || false
  const discountPercentage = (product as any).sale_discount_percentage || 0
  const stockQuantity = (product as any).stock || (product as any).inventory?.stock_quantity || 0
  const stockStatus = (product as any).stock_status || 'in_stock'
  const isLowStock = (product as any).is_low_stock || false
  const featured = (product as any).featured || false

  const displayPrice = currentPrice
  const comparePrice = isOnSale ? originalPrice : ((product as any).compare_price || (product as any).pricing?.compare_price)
  const isOutOfStock = stockStatus === 'out_of_stock' || stockQuantity <= 0

  const handleAddToCart = async () => {
    if (isOutOfStock) return
    
    try {
      await addItem(product.id, quantity)
      toast.success(`Added ${quantity} item(s) to cart!`)
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
      await addToWishlistMutation.mutateAsync(product.id)
    } catch (error) {
      // Error handled by mutation
    }
  }

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: product.name,
          text: product.short_description || product.description,
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

      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-6 relative z-10">
        {/* Compact Header */}
        <div className="mb-6">
          <nav className="flex items-center space-x-2 text-sm text-gray-400 mb-3">
            <Link href="/" className="hover:text-[#ff9000] transition-colors">Home</Link>
            <span>/</span>
            <Link href="/products" className="hover:text-[#ff9000] transition-colors">Products</Link>
            {product.category && (
              <>
                <span>/</span>
                <Link href={`/products?category=${product.category.id}`} className="hover:text-[#ff9000] transition-colors">
                  {product.category.name}
                </Link>
              </>
            )}
          </nav>
          
          <Button variant="outline" size="sm" asChild className="border-white/20 hover:bg-white/10 transition-all duration-300">
            <Link href="/products">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to Products
            </Link>
          </Button>
        </div>

        {/* Main Product Section - Unified Card */}
        <div className="bg-white/[0.03] backdrop-blur-sm border border-white/10 rounded-xl overflow-hidden mb-6">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-0">
            {/* Product Images */}
            <div className="p-6 border-b lg:border-b-0 lg:border-r border-white/10">
              {/* Main Image */}
              <div className="relative group mb-4">
                <div className="aspect-square overflow-hidden rounded-lg bg-white/[0.02] border border-white/10">
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
                    alt={product.name}
                    fill
                    className="object-cover transition-transform duration-300 group-hover:scale-[1.02]"
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
                      <Image src={image.url} alt="" width={64} height={64} className="object-cover w-full h-full" />
                    </button>
                  ))}
                </div>
              )}
            </div>

            {/* Product Info */}
            <div className="p-6 space-y-4">
              {/* Category */}
              {product.category && (
                <Link href={`/products?category=${product.category.id}`}>
                  <Badge className="bg-[#ff9000]/20 text-[#ff9000] border border-[#ff9000]/30 hover:bg-[#ff9000]/30 transition-colors">
                    {product.category.name}
                  </Badge>
                </Link>
              )}

              {/* Title */}
              <h1 className="text-2xl lg:text-3xl font-bold text-white leading-tight">
                {product.name}
              </h1>

              {/* Rating */}
              {product.rating_average > 0 && (
                <div className="flex items-center gap-3">
                  <div className="flex items-center gap-1">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={cn(
                          'h-4 w-4',
                          i < Math.floor(product.rating_average)
                            ? 'text-yellow-400 fill-current'
                            : 'text-gray-600'
                        )}
                      />
                    ))}
                    <span className="text-sm font-medium text-yellow-400 ml-1">
                      {product.rating_average.toFixed(1)}
                    </span>
                  </div>
                  <span className="text-sm text-gray-400">({product.rating_count} reviews)</span>
                </div>
              )}

              {/* Enhanced Price Display */}
              <div className="border border-white/10 rounded-lg p-4 bg-white/[0.02]">
                <div className="flex items-center gap-3 mb-2">
                  <span className="text-3xl font-bold text-[#ff9000]">
                    {formatPrice(displayPrice)}
                  </span>
                  {hasDiscount && (
                    <div className="flex items-center gap-2">
                      {isOnSale ? (
                        <span className="text-lg text-gray-500 line-through">
                          {formatPrice(originalPrice)}
                        </span>
                      ) : comparePrice && (
                        <span className="text-lg text-gray-500 line-through">
                          {formatPrice(comparePrice)}
                        </span>
                      )}
                      {discountPercentage > 0 && (
                        <Badge className="bg-[#ff9000] text-white text-xs px-2 py-1">
                          -{Math.round(discountPercentage)}%
                        </Badge>
                      )}
                    </div>
                  )}
                </div>
                {hasDiscount && discountPercentage > 0 && (
                  <p className="text-green-400 text-sm font-medium">
                    💰 You save {formatPrice((isOnSale ? originalPrice : comparePrice || 0) - displayPrice)}
                  </p>
                )}
                {isOnSale && (
                  <p className="text-[#ff9000] text-sm font-medium mt-1">
                    🔥 Limited time sale!
                  </p>
                )}
              </div>

              {/* Enhanced Stock Status */}
              <div className="flex items-center gap-2">
                {isOutOfStock ? (
                  <Badge className="bg-red-500/20 text-red-400 border border-red-500/30">
                    ❌ Sold Out
                  </Badge>
                ) : stockStatus === 'on_backorder' ? (
                  <Badge className="bg-blue-500/20 text-blue-400 border border-blue-500/30">
                    📦 Available on Backorder
                  </Badge>
                ) : isLowStock ? (
                  <Badge className="bg-amber-500/20 text-amber-400 border border-amber-500/30">
                    ⚠️ Low Stock - Only {stockQuantity} left
                  </Badge>
                ) : stockQuantity <= 10 ? (
                  <Badge className="bg-yellow-500/20 text-yellow-400 border border-yellow-500/30">
                    ⏰ Limited Stock - {stockQuantity} available
                  </Badge>
                ) : (
                  <Badge className="bg-green-500/20 text-green-400 border border-green-500/30">
                    ✅ In Stock ({stockQuantity} available)
                  </Badge>
                )}
              </div>

              {/* Short Description */}
              {product.short_description && (
                <p className="text-gray-300 text-sm leading-relaxed border-t border-white/10 pt-4">
                  {product.short_description}
                </p>
              )}

              {/* Purchase Actions */}
              {!isOutOfStock && (
                <div className="border-t border-white/10 pt-4 space-y-4">
                  {/* Quantity */}
                  <div className="flex items-center justify-between">
                    <span className="text-sm font-medium text-white">Quantity:</span>
                    <div className="flex items-center border border-white/20 rounded-lg overflow-hidden bg-white/[0.02]">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.max(1, quantity - 1))}
                        disabled={quantity <= 1}
                        className="px-3 py-2 h-9 border-r border-white/20"
                      >
                        <Minus className="h-3 w-3" />
                      </Button>
                      <span className="px-4 py-2 min-w-[50px] text-center font-medium">
                        {quantity}
                      </span>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.min(stockQuantity, quantity + 1))}
                        disabled={quantity >= stockQuantity}
                        className="px-3 py-2 h-9 border-l border-white/20"
                      >
                        <Plus className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>

                  {/* Action Buttons */}
                  <div className="space-y-3">
                    <Button
                      className="w-full bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold py-3 transition-all duration-300 hover:scale-[1.02]"
                      onClick={handleAddToCart}
                      disabled={cartLoading}
                    >
                      <ShoppingCart className="mr-2 h-4 w-4" />
                      {cartLoading ? 'Adding...' : 'Add to Cart'}
                    </Button>

                    <div className="grid grid-cols-2 gap-2">
                      {isAuthenticated && (
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={handleAddToWishlist}
                          disabled={addToWishlistMutation.isPending}
                          className="border-white/20 hover:bg-white/10 transition-colors"
                        >
                          <Heart className="h-4 w-4 mr-1" />
                          Wishlist
                        </Button>
                      )}
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={handleShare}
                        className="border-white/20 hover:bg-white/10 transition-colors"
                      >
                        <Share2 className="h-4 w-4 mr-1" />
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
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div className="bg-white/[0.03] border border-white/10 rounded-lg p-4 text-center">
            <Truck className="h-8 w-8 text-[#ff9000] mx-auto mb-2" />
            <h3 className="font-semibold text-white text-sm mb-1">Free Shipping</h3>
            <p className="text-gray-400 text-xs">On orders over $50</p>
          </div>
          <div className="bg-white/[0.03] border border-white/10 rounded-lg p-4 text-center">
            <Shield className="h-8 w-8 text-green-400 mx-auto mb-2" />
            <h3 className="font-semibold text-white text-sm mb-1">1 Year Warranty</h3>
            <p className="text-gray-400 text-xs">Quality guaranteed</p>
          </div>
          <div className="bg-white/[0.03] border border-white/10 rounded-lg p-4 text-center">
            <RotateCcw className="h-8 w-8 text-blue-400 mx-auto mb-2" />
            <h3 className="font-semibold text-white text-sm mb-1">Easy Returns</h3>
            <p className="text-gray-400 text-xs">30-day policy</p>
          </div>
        </div>

        {/* Compact Tabs */}
        <div className="bg-white/[0.03] backdrop-blur-sm border border-white/10 rounded-xl mb-6">
          <div className="border-b border-white/10 p-1">
            <div className="flex space-x-1">
              {[
                { key: 'description', label: 'Description', icon: '📝' },
                { key: 'reviews', label: 'Reviews', icon: '⭐' },
                { key: 'shipping', label: 'Shipping', icon: '🚚' },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as any)}
                  className={cn(
                    'flex items-center space-x-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300',
                    activeTab === tab.key
                      ? 'bg-[#ff9000] text-white'
                      : 'text-gray-400 hover:text-white hover:bg-white/8'
                  )}
                >
                  <span>{tab.icon}</span>
                  <span>{tab.label}</span>
                  {tab.key === 'reviews' && product.rating_count > 0 && (
                    <Badge className="bg-white/20 text-white text-xs px-1.5 py-0.5">
                      {product.rating_count}
                    </Badge>
                  )}
                </button>
              ))}
            </div>
          </div>
          <div className="p-6">
            {activeTab === 'description' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold text-white">Product Description</h3>
                <div className="text-gray-300 leading-relaxed">
                  {product.description ? (
                    <p>{product.description}</p>
                  ) : (
                    <p>No detailed description available for this product.</p>
                  )}
                </div>
              </div>
            )}

            {activeTab === 'reviews' && (
              <ProductReviewsSection productId={productId} />
            )}

            {activeTab === 'shipping' && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                    <Truck className="h-5 w-5 text-blue-400" />
                    Shipping Options
                  </h4>
                  <div className="space-y-3">
                    <div className="flex items-center gap-3 text-sm">
                      <CheckCircle className="h-4 w-4 text-green-400" />
                      <span className="text-gray-300">Free shipping on orders over $50</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm">
                      <Truck className="h-4 w-4 text-blue-400" />
                      <span className="text-gray-300">Standard delivery: 3-5 business days</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm">
                      <Gift className="h-4 w-4 text-[#ff9000]" />
                      <span className="text-gray-300">Express delivery: 1-2 business days</span>
                    </div>
                  </div>
                </div>
                <div>
                  <h4 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                    <RotateCcw className="h-5 w-5 text-green-400" />
                    Return Policy
                  </h4>
                  <div className="space-y-3">
                    <div className="flex items-center gap-3 text-sm">
                      <CheckCircle className="h-4 w-4 text-green-400" />
                      <span className="text-gray-300">30-day return policy</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm">
                      <Shield className="h-4 w-4 text-blue-400" />
                      <span className="text-gray-300">Items must be in original condition</span>
                    </div>
                    <div className="flex items-center gap-3 text-sm">
                      <Truck className="h-4 w-4 text-[#ff9000]" />
                      <span className="text-gray-300">Free return shipping</span>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Related Products */}
        {relatedProducts && relatedProducts.length > 0 && (
          <div className="bg-white/[0.03] backdrop-blur-sm border border-white/10 rounded-xl">
            <div className="p-4 border-b border-white/10">
              <h2 className="text-xl font-bold text-white mb-2">You Might Also Like</h2>
              <p className="text-gray-400 text-sm">Discover similar products</p>
            </div>
            <div className="p-4">
              <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
                {relatedProducts.map((relatedProduct, index) => (
                  <div
                    key={relatedProduct.id}
                    className="animate-in fade-in duration-500"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <ProductCard
                      product={relatedProduct}
                      className="transition-all duration-300 hover:scale-[1.02] hover:shadow-lg hover:shadow-[#ff9000]/10"
                    />
                  </div>
                ))}
              </div>
              <div className="text-center mt-6">
                <Button
                  asChild
                  variant="outline"
                  className="border-white/20 hover:bg-white/10 transition-colors"
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
