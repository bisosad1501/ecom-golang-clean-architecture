'use client'

import Link from 'next/link'
import { ArrowRight, Star, Truck, Shield, CreditCard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ProductCard } from '@/components/products/product-card'
import { useFeaturedProducts } from '@/hooks/use-products'
import { APP_NAME } from '@/constants'

export function HomePage() {
  const { data: featuredProducts, isLoading, error } = useFeaturedProducts(8)

  return (
    <div className="min-h-screen">

      {/* Hero Section */}
      <section className="relative bg-gradient-to-r from-primary-600 to-primary-800 text-white">
        <div className="container mx-auto px-4 py-20 lg:py-32">
          <div className="max-w-3xl">
            <h1 className="text-4xl lg:text-6xl font-bold mb-6">
              Discover Amazing Products at Great Prices
            </h1>
            <p className="text-xl lg:text-2xl mb-8 text-primary-100">
              Shop the latest trends with fast shipping, easy returns, and excellent customer service.
            </p>
            <div className="flex flex-col sm:flex-row gap-4">
              <Button size="lg" variant="secondary" asChild>
                <Link href="/products">
                  Shop Now
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
              <Button size="lg" variant="outline" className="border-white text-white hover:bg-white hover:text-primary-600" asChild>
                <Link href="/categories">
                  Browse Categories
                </Link>
              </Button>
            </div>
          </div>
        </div>
        
        {/* Hero background pattern */}
        <div className="absolute inset-0 bg-black bg-opacity-10">
          <div className="absolute inset-0 bg-gradient-to-br from-transparent to-black opacity-20"></div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-muted">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-primary-100 text-primary-600 rounded-full mb-4">
                <Truck className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-foreground">Free Shipping</h3>
              <p className="text-muted-foreground">Free shipping on all orders over $50. Fast and reliable delivery.</p>
            </div>
            
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-primary-100 text-primary-600 rounded-full mb-4">
                <Shield className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-foreground">Secure Payment</h3>
              <p className="text-muted-foreground">Your payment information is always safe and secure with us.</p>
            </div>
            
            <div className="text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 bg-primary-100 text-primary-600 rounded-full mb-4">
                <CreditCard className="h-8 w-8" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-foreground">Easy Returns</h3>
              <p className="text-muted-foreground">30-day return policy. No questions asked, hassle-free returns.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Products Section */}
      <section className="py-16 bg-background">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-4">
              Featured Products
            </h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Discover our handpicked selection of the best products, chosen for their quality and value.
            </p>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {[...Array(8)].map((_, i) => (
                <Card key={i} className="animate-pulse">
                  <div className="aspect-square bg-muted rounded-t-lg"></div>
                  <CardContent className="p-4">
                    <div className="h-4 bg-muted rounded mb-2"></div>
                    <div className="h-4 bg-muted rounded w-2/3 mb-2"></div>
                    <div className="h-6 bg-muted rounded w-1/2"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <div className="text-center py-12">
              <p className="text-destructive text-lg mb-4">Unable to load featured products</p>
              <p className="text-muted-foreground">Please check if the backend API is running</p>
              <Button className="mt-4" asChild>
                <Link href="/test">View Test Page</Link>
              </Button>
            </div>
          ) : featuredProducts && featuredProducts.length > 0 ? (
            <>
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
                {featuredProducts.map((product) => (
                  <ProductCard key={product.id} product={product} />
                ))}
              </div>
              
              <div className="text-center mt-12">
                <Button size="lg" variant="outline" asChild>
                  <Link href="/products">
                    View All Products
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Link>
                </Button>
              </div>
            </>
          ) : (
            <div className="text-center py-12">
              <p className="text-muted-foreground text-lg mb-4">No featured products available at the moment.</p>
              <p className="text-muted-foreground/80 mb-4">This might be because the backend API is not running.</p>
              <div className="flex justify-center gap-4">
                <Button asChild>
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" asChild>
                  <Link href="/test">View Test Page</Link>
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-16 bg-muted">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-4">
              Shop by Category
            </h2>
            <p className="text-xl text-muted-foreground">
              Find exactly what you're looking for in our organized categories.
            </p>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-6">
            {[
              { name: 'Electronics', image: '/categories/electronics.jpg', href: '/categories/electronics' },
              { name: 'Fashion', image: '/categories/fashion.jpg', href: '/categories/fashion' },
              { name: 'Home & Garden', image: '/categories/home.jpg', href: '/categories/home-garden' },
              { name: 'Sports', image: '/categories/sports.jpg', href: '/categories/sports' },
            ].map((category) => (
              <Link key={category.name} href={category.href} className="group">
                <Card className="overflow-hidden transition-transform group-hover:scale-105">
                  <div className="aspect-square bg-muted relative">
                    <div className="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent"></div>
                    <div className="absolute bottom-4 left-4 text-white">
                      <h3 className="text-lg font-semibold">{category.name}</h3>
                    </div>
                  </div>
                </Card>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* Testimonials Section */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="text-center mb-12">
            <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 mb-4">
              What Our Customers Say
            </h2>
            <p className="text-xl text-gray-600">
              Don't just take our word for it - hear from our satisfied customers.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {[
              {
                name: 'Sarah Johnson',
                rating: 5,
                comment: 'Amazing quality products and super fast shipping. Will definitely shop here again!',
              },
              {
                name: 'Mike Chen',
                rating: 5,
                comment: 'Great customer service and easy returns. The product quality exceeded my expectations.',
              },
              {
                name: 'Emily Davis',
                rating: 5,
                comment: 'Love the variety of products and competitive prices. Highly recommend this store!',
              },
            ].map((testimonial, index) => (
              <Card key={index} className="p-6">
                <div className="flex items-center mb-4">
                  {[...Array(testimonial.rating)].map((_, i) => (
                    <Star key={i} className="h-5 w-5 text-yellow-400 fill-current" />
                  ))}
                </div>
                <p className="text-muted-foreground mb-4">"{testimonial.comment}"</p>
                <p className="font-semibold text-foreground">- {testimonial.name}</p>
              </Card>
            ))}
          </div>
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="py-16 bg-primary text-primary-foreground">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl lg:text-4xl font-bold mb-4">
            Stay Updated with Our Latest Offers
          </h2>
          <p className="text-xl text-primary-100 mb-8 max-w-2xl mx-auto">
            Subscribe to our newsletter and be the first to know about new products, exclusive deals, and special promotions.
          </p>
          
          <div className="max-w-md mx-auto flex gap-4">
            <input
              type="email"
              placeholder="Enter your email"
              className="flex-1 px-4 py-3 rounded-lg text-foreground placeholder-muted-foreground bg-background border border-input"
            />
            <Button variant="secondary" size="lg">
              Subscribe
            </Button>
          </div>
          
          <p className="text-sm text-primary-200 mt-4">
            We respect your privacy. Unsubscribe at any time.
          </p>
        </div>
      </section>
    </div>
  )
}
