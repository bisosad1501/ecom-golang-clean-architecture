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
      <div className="min-h-screen bg-black py-8">
        <div className="container mx-auto px-4">
          <div className="animate-pulse">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
              <div className="aspect-square bg-gray-800 rounded-lg"></div>
              <div className="space-y-4">
                <div className="h-8 bg-gray-800 rounded w-3/4"></div>
                <div className="h-6 bg-gray-800 rounded w-1/2"></div>
                <div className="h-12 bg-gray-800 rounded w-full"></div>
                <div className="h-10 bg-gray-800 rounded w-1/3"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-black py-8">
        <div className="container mx-auto px-4">
          <Card className="p-8 text-center max-w-md mx-auto bg-gray-900 border-gray-700">
            <h2 className="text-2xl font-bold text-white mb-4">
              Product Not Found
            </h2>
            <p className="text-gray-300 mb-6">
              The product you're looking for doesn't exist or has been removed.
            </p>
            <Button asChild className="bg-orange-500 hover:bg-orange-600 text-white">
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
    <div className="min-h-screen bg-black py-8">
      <div className="container mx-auto px-4">
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
                <Badge className="absolute top-4 left-4 z-10 bg-red-500 text-white">
                  -{discountPercentage}%
                </Badge>
              )}
              {isOutOfStock && (
                <Badge className="absolute top-4 right-4 z-10 bg-gray-700 text-gray-200">
                  Out of Stock
                </Badge>
              )}
              <Image
                src={images[selectedImageIndex]?.url || '/placeholder-product.jpg'}
                alt={product.name}
                fill
                className="object-cover"
                priority
              />

              {/* Image Navigation */}
              {images.length > 1 && (
                <>
                  <Button
                    variant="secondary"
                    size="icon"
                    className="absolute left-4 top-1/2 transform -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white border-gray-600"
                    onClick={() => setSelectedImageIndex(prev =>
                      prev === 0 ? images.length - 1 : prev - 1
                    )}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="secondary"
                    size="icon"
                    className="absolute right-4 top-1/2 transform -translate-y-1/2 bg-black/50 hover:bg-black/70 text-white border-gray-600"
                    onClick={() => setSelectedImageIndex(prev =>
                      prev === images.length - 1 ? 0 : prev + 1
                    )}
                  >
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </>
              )}
            </div>

            {/* Thumbnail Images */}
            {images.length > 1 && (
              <div className="flex space-x-2 overflow-x-auto">
                {images.map((image, index) => (
                  <button
                    key={image.id}
                    onClick={() => setSelectedImageIndex(index)}
                    className={cn(
                      'relative w-20 h-20 flex-shrink-0 overflow-hidden rounded-md border-2 transition-colors',
                      selectedImageIndex === index
                        ? 'border-orange-500'
                        : 'border-gray-600 hover:border-gray-500'
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

          {/* Product Info */}
          <div className="space-y-6">
            {/* Category */}
            {product.category && (
              <Link
                href={`/products?category=${product.category.id}`}
                className="text-sm text-orange-400 hover:text-orange-300 font-medium transition-colors"
              >
                {product.category.name}
              </Link>
            )}

            {/* Title - More compact */}
            <h1 className="text-2xl lg:text-3xl font-bold text-white">{product.name}</h1>

            {/* Rating */}
            {product.rating && (
              <div className="flex items-center space-x-2">
                <div className="flex">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={cn(
                        'h-5 w-5',
                        i < Math.floor(product.rating!.average)
                          ? 'text-yellow-400 fill-current'
                          : 'text-gray-500'
                      )}
                    />
                  ))}
                </div>
                <span className="text-sm text-gray-300">
                  {product.rating.average.toFixed(1)} ({product.rating.count} reviews)
                </span>
              </div>
            )}

            {/* Price - More compact */}
            <div className="space-y-1.5">
              <div className="flex items-center space-x-3">
                <span className="text-2xl lg:text-3xl font-bold text-white">
                  {formatPrice(displayPrice)}
                </span>
                {hasDiscount && (
                  <span className="text-lg lg:text-xl text-gray-500 line-through">
                    {formatPrice(product.price)}
                  </span>
                )}
              </div>
              {hasDiscount && (
                <p className="text-green-400 font-medium">
                  You save {formatPrice(product.price - product.sale_price!)} ({discountPercentage}% off)
                </p>
              )}
            </div>

            {/* Short Description - More compact */}
            {product.short_description && (
              <p className="text-gray-300 text-base leading-relaxed">
                {product.short_description}
              </p>
            )}

            {/* Stock Status */}
            <div className="flex items-center space-x-2">
              {isOutOfStock ? (
                <Badge className="bg-red-500 text-white">Out of Stock</Badge>
              ) : product.stock <= 5 ? (
                <Badge className="bg-yellow-500 text-black">Only {product.stock} left in stock</Badge>
              ) : (
                <Badge className="bg-green-500 text-white">In Stock</Badge>
              )}
            </div>

            {/* Quantity and Add to Cart */}
            {!isOutOfStock && (
              <div className="space-y-4">
                <div className="flex items-center space-x-4">
                  <span className="text-sm font-medium text-white">Quantity:</span>
                  <div className="flex items-center border border-gray-600 rounded-md bg-gray-800">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setQuantity(Math.max(1, quantity - 1))}
                      disabled={quantity <= 1}
                      className="text-white hover:bg-gray-700"
                    >
                      <Minus className="h-4 w-4" />
                    </Button>
                    <span className="px-4 py-2 text-center min-w-[60px] text-white">{quantity}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                      disabled={quantity >= product.stock}
                      className="text-white hover:bg-gray-700"
                    >
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                <div className="flex space-x-3">
                  <Button
                    size="default"
                    className="flex-1 bg-orange-500 hover:bg-orange-600 text-white"
                    onClick={handleAddToCart}
                    disabled={cartLoading}
                  >
                    <ShoppingCart className="mr-2 h-4 w-4" />
                    {cartLoading ? 'Adding...' : 'Add to Cart'}
                  </Button>

                  {isAuthenticated && (
                    <Button
                      variant="outline"
                      size="default"
                      onClick={handleAddToWishlist}
                      disabled={addToWishlistMutation.isPending}
                      className="border-gray-600 text-white hover:bg-gray-800"
                    >
                      <Heart className="h-5 w-5" />
                    </Button>
                  )}

                  <Button
                    variant="outline"
                    size="lg"
                    onClick={handleShare}
                    className="border-gray-600 text-white hover:bg-gray-800"
                  >
                    <Share2 className="h-5 w-5" />
                  </Button>
                </div>
              </div>
            )}

            {/* Product Features */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-6 border-t border-gray-700">
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <Truck className="h-5 w-5 text-orange-400" />
                <span>Free shipping over $50</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <Shield className="h-5 w-5 text-orange-400" />
                <span>1 year warranty</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-300">
                <RotateCcw className="h-5 w-5 text-orange-400" />
                <span>30-day returns</span>
              </div>
            </div>
          </div>
        </div>

        {/* Product Tabs */}
        <Card className="mb-12 bg-gray-900 border-gray-700">
          <CardHeader>
            <div className="flex space-x-8 border-b border-gray-700">
              {[
                { key: 'description', label: 'Description' },
                { key: 'reviews', label: 'Reviews' },
                { key: 'shipping', label: 'Shipping & Returns' },
              ].map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key as any)}
                  className={cn(
                    'pb-4 px-1 border-b-2 font-medium text-sm transition-colors',
                    activeTab === tab.key
                      ? 'border-orange-500 text-orange-400'
                      : 'border-transparent text-gray-400 hover:text-gray-200'
                  )}
                >
                  {tab.label}
                </button>
              ))}
            </div>
          </CardHeader>
          <CardContent className="pt-6">
            {activeTab === 'description' && (
              <div className="prose max-w-none">
                <p className="text-gray-300 leading-relaxed">
                  {product.description}
                </p>
                {/* Additional product details would go here */}
              </div>
            )}

            {activeTab === 'reviews' && (
              <div className="text-center py-8">
                <p className="text-gray-400">Reviews feature coming soon!</p>
              </div>
            )}

            {activeTab === 'shipping' && (
              <div className="space-y-4">
                <div>
                  <h4 className="font-semibold mb-2 text-white">Shipping Information</h4>
                  <ul className="text-gray-300 space-y-1">
                    <li>• Free shipping on orders over $50</li>
                    <li>• Standard delivery: 3-5 business days</li>
                    <li>• Express delivery: 1-2 business days</li>
                  </ul>
                </div>
                <div>
                  <h4 className="font-semibold mb-2 text-white">Returns & Exchanges</h4>
                  <ul className="text-gray-300 space-y-1">
                    <li>• 30-day return policy</li>
                    <li>• Items must be in original condition</li>
                    <li>• Free return shipping</li>
                  </ul>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Related Products */}
        {relatedProducts && relatedProducts.length > 0 && (
          <div>
            <h2 className="text-2xl font-bold text-white mb-6">Related Products</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {relatedProducts.map((relatedProduct) => (
                <ProductCard key={relatedProduct.id} product={relatedProduct} />
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
