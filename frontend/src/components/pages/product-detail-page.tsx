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
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="container mx-auto px-4">
          <div className="animate-pulse">
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
              <div className="aspect-square bg-gray-200 rounded-lg"></div>
              <div className="space-y-4">
                <div className="h-8 bg-gray-200 rounded w-3/4"></div>
                <div className="h-6 bg-gray-200 rounded w-1/2"></div>
                <div className="h-12 bg-gray-200 rounded w-full"></div>
                <div className="h-10 bg-gray-200 rounded w-1/3"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (error || !product) {
    return (
      <div className="min-h-screen bg-gray-50 py-8">
        <div className="container mx-auto px-4">
          <Card className="p-8 text-center max-w-md mx-auto">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              Product Not Found
            </h2>
            <p className="text-gray-600 mb-6">
              The product you're looking for doesn't exist or has been removed.
            </p>
            <Button asChild>
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
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4">
        {/* Breadcrumb */}
        <nav className="flex items-center space-x-2 text-sm text-gray-600 mb-8">
          <Link href="/" className="hover:text-gray-900">Home</Link>
          <span>/</span>
          <Link href="/products" className="hover:text-gray-900">Products</Link>
          {product.category && (
            <>
              <span>/</span>
              <Link href={`/products?category=${product.category.id}`} className="hover:text-gray-900">
                {product.category.name}
              </Link>
            </>
          )}
          <span>/</span>
          <span className="text-gray-900">{product.name}</span>
        </nav>

        {/* Product Details */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-12">
          {/* Product Images */}
          <div className="space-y-4">
            {/* Main Image */}
            <div className="relative aspect-square overflow-hidden rounded-lg bg-gray-100">
              {hasDiscount && (
                <Badge variant="destructive" className="absolute top-4 left-4 z-10">
                  -{discountPercentage}%
                </Badge>
              )}
              {isOutOfStock && (
                <Badge variant="secondary" className="absolute top-4 right-4 z-10">
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
                    className="absolute left-4 top-1/2 transform -translate-y-1/2"
                    onClick={() => setSelectedImageIndex(prev => 
                      prev === 0 ? images.length - 1 : prev - 1
                    )}
                  >
                    <ChevronLeft className="h-4 w-4" />
                  </Button>
                  <Button
                    variant="secondary"
                    size="icon"
                    className="absolute right-4 top-1/2 transform -translate-y-1/2"
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
                        ? 'border-primary-500'
                        : 'border-gray-200 hover:border-gray-300'
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
                className="text-sm text-primary-600 hover:text-primary-700 font-medium"
              >
                {product.category.name}
              </Link>
            )}

            {/* Title */}
            <h1 className="text-3xl font-bold text-gray-900">{product.name}</h1>

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
                          : 'text-gray-300'
                      )}
                    />
                  ))}
                </div>
                <span className="text-sm text-gray-600">
                  {product.rating.average.toFixed(1)} ({product.rating.count} reviews)
                </span>
              </div>
            )}

            {/* Price */}
            <div className="space-y-2">
              <div className="flex items-center space-x-3">
                <span className="text-3xl font-bold text-gray-900">
                  {formatPrice(displayPrice)}
                </span>
                {hasDiscount && (
                  <span className="text-xl text-gray-500 line-through">
                    {formatPrice(product.price)}
                  </span>
                )}
              </div>
              {hasDiscount && (
                <p className="text-green-600 font-medium">
                  You save {formatPrice(product.price - product.sale_price!)} ({discountPercentage}% off)
                </p>
              )}
            </div>

            {/* Short Description */}
            {product.short_description && (
              <p className="text-gray-600 text-lg leading-relaxed">
                {product.short_description}
              </p>
            )}

            {/* Stock Status */}
            <div className="flex items-center space-x-2">
              {isOutOfStock ? (
                <Badge variant="destructive">Out of Stock</Badge>
              ) : product.stock <= 5 ? (
                <Badge variant="warning">Only {product.stock} left in stock</Badge>
              ) : (
                <Badge variant="success">In Stock</Badge>
              )}
            </div>

            {/* Quantity and Add to Cart */}
            {!isOutOfStock && (
              <div className="space-y-4">
                <div className="flex items-center space-x-4">
                  <span className="text-sm font-medium">Quantity:</span>
                  <div className="flex items-center border rounded-md">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setQuantity(Math.max(1, quantity - 1))}
                      disabled={quantity <= 1}
                    >
                      <Minus className="h-4 w-4" />
                    </Button>
                    <span className="px-4 py-2 text-center min-w-[60px]">{quantity}</span>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setQuantity(Math.min(product.stock, quantity + 1))}
                      disabled={quantity >= product.stock}
                    >
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                <div className="flex space-x-4">
                  <Button
                    size="lg"
                    className="flex-1"
                    onClick={handleAddToCart}
                    disabled={cartLoading}
                    isLoading={cartLoading}
                  >
                    <ShoppingCart className="mr-2 h-5 w-5" />
                    Add to Cart
                  </Button>
                  
                  {isAuthenticated && (
                    <Button
                      variant="outline"
                      size="lg"
                      onClick={handleAddToWishlist}
                      disabled={addToWishlistMutation.isPending}
                    >
                      <Heart className="h-5 w-5" />
                    </Button>
                  )}
                  
                  <Button
                    variant="outline"
                    size="lg"
                    onClick={handleShare}
                  >
                    <Share2 className="h-5 w-5" />
                  </Button>
                </div>
              </div>
            )}

            {/* Product Features */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-6 border-t">
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <Truck className="h-5 w-5 text-primary-600" />
                <span>Free shipping over $50</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <Shield className="h-5 w-5 text-primary-600" />
                <span>1 year warranty</span>
              </div>
              <div className="flex items-center space-x-2 text-sm text-gray-600">
                <RotateCcw className="h-5 w-5 text-primary-600" />
                <span>30-day returns</span>
              </div>
            </div>
          </div>
        </div>

        {/* Product Tabs */}
        <Card className="mb-12">
          <CardHeader>
            <div className="flex space-x-8 border-b">
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
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700'
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
                <p className="text-gray-600 leading-relaxed">
                  {product.description}
                </p>
                {/* Additional product details would go here */}
              </div>
            )}
            
            {activeTab === 'reviews' && (
              <div className="text-center py-8">
                <p className="text-gray-500">Reviews feature coming soon!</p>
              </div>
            )}
            
            {activeTab === 'shipping' && (
              <div className="space-y-4">
                <div>
                  <h4 className="font-semibold mb-2">Shipping Information</h4>
                  <ul className="text-gray-600 space-y-1">
                    <li>• Free shipping on orders over $50</li>
                    <li>• Standard delivery: 3-5 business days</li>
                    <li>• Express delivery: 1-2 business days</li>
                  </ul>
                </div>
                <div>
                  <h4 className="font-semibold mb-2">Returns & Exchanges</h4>
                  <ul className="text-gray-600 space-y-1">
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
            <h2 className="text-2xl font-bold text-gray-900 mb-6">Related Products</h2>
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
