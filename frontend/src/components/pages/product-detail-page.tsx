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
  ChevronRight
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProductCard } from '@/components/products/product-card'
import { ProductReviewsSection } from '@/components/reviews/product-reviews-section'
import { useProduct, useRelatedProducts, useAddToWishlist } from '@/hooks/use-products'
import { useCartStore } from '@/store/cart'
import { useAuthStore } from '@/store/auth'
import { formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'

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
      <div className="min-h-screen bg-gradient-to-br from-slate-950 via-gray-900 to-black">
        <div className="container mx-auto px-4 py-12">
          <div className="animate-pulse">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 mb-16">
              <div className="space-y-4">
                <div className="aspect-square bg-gradient-to-br from-gray-800/50 to-gray-900/50 rounded-2xl backdrop-blur-sm border border-white/10"></div>
                <div className="grid grid-cols-4 gap-3">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="aspect-square bg-gradient-to-br from-gray-800/30 to-gray-900/30 rounded-xl"></div>
                  ))}
                </div>
              </div>
              <div className="space-y-6">
                <div className="h-10 bg-gradient-to-r from-gray-800/50 to-gray-700/50 rounded-xl w-3/4"></div>
                <div className="h-6 bg-gradient-to-r from-orange-500/30 to-orange-600/30 rounded-lg w-1/2"></div>
                <div className="h-16 bg-gradient-to-r from-gray-800/50 to-gray-700/50 rounded-xl w-full"></div>
                <div className="h-12 bg-gradient-to-r from-orange-500/30 to-orange-600/30 rounded-xl w-1/3"></div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="h-12 bg-gradient-to-r from-gray-800/50 to-gray-700/50 rounded-xl"></div>
                  <div className="h-12 bg-gradient-to-r from-orange-500/30 to-orange-600/30 rounded-xl"></div>
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
      <div className="min-h-screen bg-gradient-to-br from-slate-950 via-gray-900 to-black flex items-center justify-center">
        <div className="container mx-auto px-4">
          <Card className="p-12 text-center max-w-lg mx-auto bg-white/5 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl">
            <div className="w-20 h-20 bg-gradient-to-br from-red-500/20 to-red-600/20 rounded-full flex items-center justify-center mx-auto mb-6 border border-red-500/30">
              <span className="text-4xl">😞</span>
            </div>
            <h2 className="text-3xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-4">
              Product Not Found
            </h2>
            <p className="text-gray-300 mb-8 text-lg leading-relaxed">
              The product you're looking for doesn't exist or has been removed.
            </p>
            <Button asChild className="bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white px-8 py-3 rounded-xl font-semibold shadow-lg shadow-orange-500/25 transition-all duration-300 hover:scale-105">
              <Link href="/products">Browse Products</Link>
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  const images = product.images || []
  const hasDiscount = product.sale_price && product.sale_price < product.price
  const discountPercentage = hasDiscount 
    ? Math.round(((product.price - product.sale_price!) / product.price) * 100)
    : 0
  const displayPrice = product.sale_price || product.price
  const isOutOfStock = product.stock <= 0

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
    <div className="min-h-screen bg-black">
      <div className="container mx-auto px-4 py-8">
        {/* Breadcrumb */}
        <nav className="flex items-center space-x-2 text-sm text-gray-400 mb-8">
          <Link href="/" className="hover:text-orange-400 transition-colors">Home</Link>
          <span>/</span>
          <Link href="/products" className="hover:text-orange-400 transition-colors">Products</Link>
          {product.category && (
            <>
              <span>/</span>
              <Link href={`/products?category=${product.category.id}`} className="hover:text-orange-400 transition-colors">
                {product.category.name}
              </Link>
            </>
          )}
          <span>/</span>
          <span className="text-white">{product.name}</span>
        </nav>

        {/* Product Details - More compact */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Product Images */}
          <div className="space-y-3">
            {/* Main Image */}
            <div className="relative aspect-square overflow-hidden rounded-lg bg-gray-800">
              {hasDiscount && (
                <Badge className="absolute top-6 left-6 z-10 bg-gradient-to-r from-red-500 to-red-600 text-white px-4 py-2 text-lg font-bold shadow-lg">
                  -{discountPercentage}% OFF
                </Badge>
              )}
              {isOutOfStock && (
                <Badge className="absolute top-6 right-6 z-10 bg-gradient-to-r from-gray-700 to-gray-800 text-white px-4 py-2 text-lg font-bold shadow-lg">
                  Out of Stock
                </Badge>
              )}
              <Image
                src={images[selectedImageIndex]?.url || '/placeholder-product.jpg'}
                alt={product.name}
                fill
                className="object-cover transition-all duration-500 hover:scale-105"
                priority
              />

              {/* Enhanced Image Navigation */}
              {images.length > 1 && (
                <>
                  <Button
                    variant="secondary"
                    size="icon"
                    className="absolute left-6 top-1/2 transform -translate-y-1/2 bg-white/10 backdrop-blur-md hover:bg-white/20 text-white border border-white/20 rounded-full w-12 h-12 shadow-xl transition-all duration-300 hover:scale-110"
                    onClick={() => setSelectedImageIndex(prev =>
                      prev === 0 ? images.length - 1 : prev - 1
                    )}
                  >
                    <ChevronLeft className="h-5 w-5" />
                  </Button>
                  <Button
                    variant="secondary"
                    size="icon"
                    className="absolute right-6 top-1/2 transform -translate-y-1/2 bg-white/10 backdrop-blur-md hover:bg-white/20 text-white border border-white/20 rounded-full w-12 h-12 shadow-xl transition-all duration-300 hover:scale-110"
                    onClick={() => setSelectedImageIndex(prev =>
                      prev === images.length - 1 ? 0 : prev + 1
                    )}
                  >
                    <ChevronRight className="h-5 w-5" />
                  </Button>
                </>
              )}

              {/* Image Indicators */}
              {images.length > 1 && (
                <div className="absolute bottom-6 left-1/2 transform -translate-x-1/2 flex space-x-2">
                  {images.map((_, index) => (
                    <button
                      key={index}
                      onClick={() => setSelectedImageIndex(index)}
                      className={cn(
                        'w-3 h-3 rounded-full transition-all duration-300',
                        selectedImageIndex === index
                          ? 'bg-orange-500 scale-125'
                          : 'bg-white/30 hover:bg-white/50'
                      )}
                    />
                  ))}
                </div>
              )}
            </div>

            {/* Enhanced Thumbnail Images */}
            {images.length > 1 && (
              <div className="flex space-x-4 overflow-x-auto pb-2">
                {images.map((image, index) => (
                  <button
                    key={image.id}
                    onClick={() => setSelectedImageIndex(index)}
                    className={cn(
                      'relative w-24 h-24 flex-shrink-0 overflow-hidden rounded-xl border-2 transition-all duration-300 hover:scale-105',
                      selectedImageIndex === index
                        ? 'border-orange-500 shadow-lg shadow-orange-500/25 scale-105'
                        : 'border-white/20 hover:border-orange-400/50 bg-white/5 backdrop-blur-sm'
                    )}
                  >
                    <Image
                      src={image.url}
                      alt={image.alt_text || product.name}
                      fill
                      className="object-cover"
                    />
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Enhanced Product Info */}
          <div className="space-y-8">
            {/* Category Badge */}
            {product.category && (
              <Link
                href={`/products?category=${product.category.id}`}
                className="inline-block"
              >
                <Badge className="bg-gradient-to-r from-orange-500/20 to-orange-600/20 text-orange-300 border border-orange-500/30 px-4 py-2 text-sm font-medium hover:from-orange-500/30 hover:to-orange-600/30 transition-all duration-300 backdrop-blur-sm">
                  📁 {product.category.name}
                </Badge>
              </Link>
            )}

            {/* Enhanced Title */}
            <div className="space-y-4">
              <h1 className="text-3xl lg:text-5xl font-bold bg-gradient-to-r from-white via-gray-100 to-orange-200 bg-clip-text text-transparent leading-tight">
                {product.name}
              </h1>

              {/* Enhanced Rating */}
              {product.rating && (
                <div className="flex items-center space-x-4">
                  <div className="flex items-center space-x-2">
                    <div className="flex">
                      {[...Array(5)].map((_, i) => (
                        <Star
                          key={i}
                          className={cn(
                            'h-6 w-6 transition-all duration-200',
                            i < Math.floor(product.rating!.average)
                              ? 'text-yellow-400 fill-current drop-shadow-sm'
                              : 'text-gray-500'
                          )}
                        />
                      ))}
                    </div>
                    <span className="text-lg font-semibold text-yellow-400">
                      {product.rating.average.toFixed(1)}
                    </span>
                  </div>
                  <Badge className="bg-white/10 text-gray-300 border border-white/20 px-3 py-1">
                    {product.rating.count} reviews
                  </Badge>
                </div>
              )}
            </div>

            {/* Enhanced Price Section */}
            <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6 space-y-4">
              <div className="flex items-center space-x-4">
                <span className="text-3xl lg:text-4xl font-bold bg-gradient-to-r from-orange-400 to-orange-500 bg-clip-text text-transparent">
                  {formatPrice(displayPrice)}
                </span>
                {hasDiscount && (
                  <div className="flex items-center space-x-3">
                    <span className="text-xl lg:text-2xl text-gray-400 line-through">
                      {formatPrice(product.price)}
                    </span>
                    <Badge className="bg-gradient-to-r from-red-500 to-red-600 text-white px-3 py-1 text-sm font-bold">
                      -{discountPercentage}% OFF
                    </Badge>
                  </div>
                )}
              </div>
              {hasDiscount && (
                <div className="bg-green-500/10 border border-green-500/20 rounded-xl p-4">
                  <p className="text-green-400 font-semibold text-lg">
                    💰 You save {formatPrice(product.price - product.sale_price!)} ({discountPercentage}% off)
                  </p>
                </div>
              )}
            </div>

            {/* Enhanced Short Description */}
            {product.short_description && (
              <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6">
                <p className="text-gray-200 text-lg leading-relaxed">
                  {product.short_description}
                </p>
              </div>
            )}

            {/* Enhanced Stock Status */}
            <div className="flex items-center space-x-3">
              {isOutOfStock ? (
                <Badge className="bg-gradient-to-r from-red-500 to-red-600 text-white px-4 py-2 text-lg font-semibold shadow-lg">
                  ❌ Out of Stock
                </Badge>
              ) : product.stock <= 5 ? (
                <Badge className="bg-gradient-to-r from-yellow-500 to-yellow-600 text-black px-4 py-2 text-lg font-semibold shadow-lg animate-pulse">
                  ⚠️ Only {product.stock} left in stock
                </Badge>
              ) : (
                <Badge className="bg-gradient-to-r from-green-500 to-green-600 text-white px-4 py-2 text-lg font-semibold shadow-lg">
                  ✅ In Stock
                </Badge>
              )}
            </div>

            {/* Enhanced Quantity and Add to Cart */}
            {!isOutOfStock && (
              <div className="space-y-6">
                <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6">
                  <div className="flex items-center justify-between mb-4">
                    <span className="text-lg font-semibold text-white">Quantity:</span>
                    <div className="flex items-center bg-white/10 backdrop-blur-md border border-white/20 rounded-xl overflow-hidden">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.max(1, quantity - 1))}
                        disabled={quantity <= 1}
                        className="text-white hover:bg-white/20 px-4 py-3 rounded-none border-r border-white/20"
                      >
                        <Minus className="h-5 w-5" />
                      </Button>
                      <span className="px-6 py-3 text-center min-w-[80px] text-white font-bold text-lg bg-white/5">
                        {quantity}
                      </span>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                        disabled={quantity >= product.stock}
                        className="text-white hover:bg-white/20 px-4 py-3 rounded-none border-l border-white/20"
                      >
                        <Plus className="h-5 w-5" />
                      </Button>
                    </div>
                  </div>

                  {/* Enhanced Action Buttons */}
                  <div className="grid grid-cols-1 gap-4">
                    <Button
                      size="lg"
                      className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-bold py-4 text-lg shadow-xl shadow-orange-500/25 transition-all duration-300 hover:scale-105 hover:shadow-2xl hover:shadow-orange-500/40"
                      onClick={handleAddToCart}
                      disabled={cartLoading}
                    >
                      <ShoppingCart className="mr-3 h-6 w-6" />
                      {cartLoading ? 'Adding to Cart...' : 'Add to Cart'}
                    </Button>

                    <div className="grid grid-cols-2 gap-3">
                      {isAuthenticated && (
                        <Button
                          variant="outline"
                          size="lg"
                          onClick={handleAddToWishlist}
                          disabled={addToWishlistMutation.isPending}
                          className="border-white/20 text-white hover:bg-white/10 backdrop-blur-sm py-3 transition-all duration-300 hover:scale-105"
                        >
                          <Heart className="h-5 w-5 mr-2" />
                          Wishlist
                        </Button>
                      )}

                      <Button
                        variant="outline"
                        size="lg"
                        onClick={handleShare}
                        className="border-white/20 text-white hover:bg-white/10 backdrop-blur-sm py-3 transition-all duration-300 hover:scale-105"
                      >
                        <Share2 className="h-5 w-5 mr-2" />
                        Share
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Enhanced Product Features */}
            <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-6">
              <h3 className="text-lg font-semibold text-white mb-4">Product Benefits</h3>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                <div className="flex items-center space-x-3 text-gray-200">
                  <div className="w-10 h-10 bg-gradient-to-br from-orange-500/20 to-orange-600/20 rounded-full flex items-center justify-center border border-orange-500/30">
                    <Truck className="h-5 w-5 text-orange-400" />
                  </div>
                  <span className="font-medium">Free shipping over $50</span>
                </div>
                <div className="flex items-center space-x-3 text-gray-200">
                  <div className="w-10 h-10 bg-gradient-to-br from-green-500/20 to-green-600/20 rounded-full flex items-center justify-center border border-green-500/30">
                    <Shield className="h-5 w-5 text-green-400" />
                  </div>
                  <span className="font-medium">1 year warranty</span>
                </div>
                <div className="flex items-center space-x-3 text-gray-200">
                  <div className="w-10 h-10 bg-gradient-to-br from-blue-500/20 to-blue-600/20 rounded-full flex items-center justify-center border border-blue-500/30">
                    <RotateCcw className="h-5 w-5 text-blue-400" />
                  </div>
                  <span className="font-medium">30-day returns</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Enhanced Product Tabs */}
        <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 shadow-2xl mb-16">
          <div className="border-b border-white/10">
            <div className="flex space-x-1 p-2">
              {[
                { key: 'description', label: 'Description', icon: '📝' },
                { key: 'reviews', label: 'Reviews', icon: '⭐' },
                { key: 'shipping', label: 'Shipping & Returns', icon: '🚚' },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as any)}
                  className={cn(
                    'flex items-center space-x-2 px-6 py-4 rounded-xl font-semibold text-sm transition-all duration-300',
                    activeTab === tab.key
                      ? 'bg-gradient-to-r from-orange-500 to-orange-600 text-white shadow-lg shadow-orange-500/25 scale-105'
                      : 'text-gray-300 hover:text-white hover:bg-white/10'
                  )}
                >
                  <span className="text-lg">{tab.icon}</span>
                  <span>{tab.label}</span>
                </button>
              ))}
            </div>
          </div>
          <div className="p-8">
            {activeTab === 'description' && (
              <div className="prose max-w-none">
                <div className="bg-white/5 backdrop-blur-md rounded-xl border border-white/10 p-6">
                  <h3 className="text-xl font-bold text-white mb-4">Product Description</h3>
                  <p className="text-gray-200 leading-relaxed text-lg">
                    {product.description}
                  </p>
                </div>
                {/* Additional product details would go here */}
              </div>
            )}

            {activeTab === 'reviews' && (
              <ProductReviewsSection productId={productId} />
            )}

            {activeTab === 'shipping' && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                <div className="bg-white/5 backdrop-blur-md rounded-xl border border-white/10 p-6">
                  <div className="flex items-center space-x-3 mb-4">
                    <div className="w-10 h-10 bg-gradient-to-br from-blue-500/20 to-blue-600/20 rounded-full flex items-center justify-center border border-blue-500/30">
                      <Truck className="h-5 w-5 text-blue-400" />
                    </div>
                    <h4 className="text-xl font-bold text-white">Shipping Information</h4>
                  </div>
                  <ul className="text-gray-200 space-y-3 text-lg">
                    <li className="flex items-center space-x-3">
                      <span className="text-green-400">✓</span>
                      <span>Free shipping on orders over $50</span>
                    </li>
                    <li className="flex items-center space-x-3">
                      <span className="text-blue-400">📦</span>
                      <span>Standard delivery: 3-5 business days</span>
                    </li>
                    <li className="flex items-center space-x-3">
                      <span className="text-orange-400">⚡</span>
                      <span>Express delivery: 1-2 business days</span>
                    </li>
                  </ul>
                </div>
                <div className="bg-white/5 backdrop-blur-md rounded-xl border border-white/10 p-6">
                  <div className="flex items-center space-x-3 mb-4">
                    <div className="w-10 h-10 bg-gradient-to-br from-green-500/20 to-green-600/20 rounded-full flex items-center justify-center border border-green-500/30">
                      <RotateCcw className="h-5 w-5 text-green-400" />
                    </div>
                    <h4 className="text-xl font-bold text-white">Returns & Exchanges</h4>
                  </div>
                  <ul className="text-gray-200 space-y-3 text-lg">
                    <li className="flex items-center space-x-3">
                      <span className="text-green-400">✓</span>
                      <span>30-day return policy</span>
                    </li>
                    <li className="flex items-center space-x-3">
                      <span className="text-blue-400">📋</span>
                      <span>Items must be in original condition</span>
                    </li>
                    <li className="flex items-center space-x-3">
                      <span className="text-orange-400">🚚</span>
                      <span>Free return shipping</span>
                    </li>
                  </ul>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Enhanced Related Products */}
        {relatedProducts && relatedProducts.length > 0 && (
          <div className="bg-white/5 backdrop-blur-xl rounded-2xl border border-white/10 p-8 shadow-2xl">
            <div className="text-center mb-8">
              <h2 className="text-3xl font-bold bg-gradient-to-r from-white via-gray-100 to-orange-200 bg-clip-text text-transparent mb-2">
                You Might Also Like
              </h2>
              <p className="text-gray-300 text-lg">Discover more amazing products</p>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
              {relatedProducts.map((relatedProduct, index) => (
                <div
                  key={relatedProduct.id}
                  className="animate-in fade-in slide-in-from-bottom-4 duration-700"
                  style={{ animationDelay: `${index * 150}ms` }}
                >
                  <ProductCard
                    product={relatedProduct}
                    className="transform hover:scale-105 transition-all duration-300 hover:shadow-2xl hover:shadow-orange-500/10"
                  />
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
