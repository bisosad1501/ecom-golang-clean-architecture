'use client'

import Link from 'next/link'
import { ArrowRight, Star, Truck, Shield, CreditCard, Sparkles, TrendingUp, Award } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ProductCard } from '@/components/products/product-card'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useFeaturedProducts } from '@/hooks/use-products'
import { useCategories } from '@/hooks/use-categories'
import { APP_NAME } from '@/constants'

export function HomePage() {
  const { data: featuredProducts, isLoading, error } = useFeaturedProducts(8)
  const { data: categories, isLoading: categoriesLoading } = useCategories()

  return (
    <div className="min-h-screen">

      {/* Hero Section */}
      <section className="relative hero-gradient text-white overflow-hidden min-h-screen flex items-center">
        <AnimatedBackground variant="hero" />

        <div className="container mx-auto px-4 py-24 lg:py-32 relative z-10">
          <div className="max-w-6xl mx-auto">
            <div className="grid lg:grid-cols-2 gap-12 items-center">
              {/* Left Content */}
              <div className="animate-fade-in">
                <div className="flex items-center gap-3 mb-8">
                  <div className="flex items-center gap-2 bg-white/10 backdrop-blur-md rounded-full px-5 py-2 border border-white/20">
                    <Sparkles className="h-4 w-4 text-yellow-300 animate-pulse" />
                    <span className="text-sm font-semibold">New Collection Available</span>
                  </div>
                </div>

                <h1 className="text-5xl lg:text-7xl xl:text-8xl font-bold mb-8 leading-[0.9]">
                  <span className="block">Discover</span>
                  <span className="block text-gradient bg-gradient-to-r from-white via-violet-200 to-white bg-clip-text text-transparent animate-gradient">
                    Amazing Products
                  </span>
                  <span className="block">at Great Prices</span>
                </h1>

                <p className="text-xl lg:text-2xl mb-10 text-white/90 leading-relaxed max-w-2xl">
                  Experience premium shopping with curated collections, lightning-fast delivery, and exceptional customer care.
                </p>

                <div className="flex flex-col sm:flex-row gap-4 mb-12">
                  <Button size="xl" variant="secondary" className="shadow-2xl hover:shadow-3xl transition-all duration-300 group bg-white text-primary-600 hover:bg-white/90" asChild>
                    <Link href="/products">
                      <Sparkles className="mr-2 h-5 w-5 group-hover:rotate-12 transition-transform" />
                      Shop Collection
                      <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                    </Link>
                  </Button>
                  <Button size="xl" variant="outline" className="border-2 border-white/60 text-white hover:bg-white/10 shadow-xl backdrop-blur-sm" asChild>
                    <Link href="/categories">
                      <TrendingUp className="mr-2 h-5 w-5" />
                      Explore Categories
                    </Link>
                  </Button>
                </div>

                {/* Trust Indicators */}
                <div className="flex items-center gap-8 text-white/80">
                  <div className="flex items-center gap-2">
                    <Shield className="h-5 w-5 text-green-400" />
                    <span className="text-sm font-medium">Secure Shopping</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Award className="h-5 w-5 text-yellow-400" />
                    <span className="text-sm font-medium">Premium Quality</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Truck className="h-5 w-5 text-blue-400" />
                    <span className="text-sm font-medium">Fast Delivery</span>
                  </div>
                </div>
              </div>

              {/* Right Content - Stats & Visual */}
              <div className="animate-scale-in" style={{ animationDelay: '0.3s' }}>
                <div className="relative">
                  {/* Floating Stats Cards */}
                  <div className="grid grid-cols-2 gap-6">
                    <div className="glass-effect p-6 rounded-2xl border border-white/20 animate-float">
                      <div className="text-3xl font-bold mb-2 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent">10K+</div>
                      <div className="text-white/80 text-sm font-medium">Happy Customers</div>
                    </div>
                    <div className="glass-effect p-6 rounded-2xl border border-white/20 animate-float" style={{ animationDelay: '0.5s' }}>
                      <div className="text-3xl font-bold mb-2 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent">5K+</div>
                      <div className="text-white/80 text-sm font-medium">Premium Products</div>
                    </div>
                    <div className="glass-effect p-6 rounded-2xl border border-white/20 animate-float" style={{ animationDelay: '1s' }}>
                      <div className="text-3xl font-bold mb-2 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent">99%</div>
                      <div className="text-white/80 text-sm font-medium">Satisfaction Rate</div>
                    </div>
                    <div className="glass-effect p-6 rounded-2xl border border-white/20 animate-float" style={{ animationDelay: '1.5s' }}>
                      <div className="text-3xl font-bold mb-2 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent">24/7</div>
                      <div className="text-white/80 text-sm font-medium">Customer Support</div>
                    </div>
                  </div>

                  {/* Decorative elements */}
                  <div className="absolute -top-4 -right-4 w-24 h-24 bg-gradient-to-br from-yellow-400/20 to-orange-500/20 rounded-full blur-xl animate-pulse"></div>
                  <div className="absolute -bottom-4 -left-4 w-32 h-32 bg-gradient-to-br from-blue-400/20 to-purple-500/20 rounded-full blur-xl animate-pulse" style={{ animationDelay: '1s' }}></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Scroll indicator */}
        <div className="absolute bottom-8 left-1/2 transform -translate-x-1/2 animate-bounce">
          <div className="w-6 h-10 border-2 border-white/30 rounded-full flex justify-center">
            <div className="w-1 h-3 bg-white/60 rounded-full mt-2 animate-pulse"></div>
          </div>
        </div>
      </section>

      {/* Value Proposition Section */}
      <section className="py-24 bg-gradient-to-br from-background via-muted/30 to-background relative overflow-hidden">
        <div className="absolute inset-0 bg-grid-pattern opacity-5"></div>

        <div className="container mx-auto px-4 relative z-10">
          <div className="text-center mb-16 animate-slide-up">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center">
                <Sparkles className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">WHY CHOOSE US</span>
            </div>

            <h2 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
              Experience the
              <span className="text-gradient"> Difference</span>
            </h2>
            <p className="text-xl text-muted-foreground max-w-4xl mx-auto leading-relaxed">
              We're not just another online store. We're your trusted partner in discovering exceptional products
              that enhance your lifestyle and exceed your expectations.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            {/* Left side - Benefits */}
            <div className="space-y-8">
              <div className="flex gap-6 group animate-fade-in">
                <div className="flex-shrink-0 w-16 h-16 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-2xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <Award className="h-8 w-8 text-white" />
                </div>
                <div>
                  <h3 className="text-2xl font-bold text-foreground mb-3">Premium Quality Guarantee</h3>
                  <p className="text-muted-foreground text-lg leading-relaxed">Every product is carefully curated and tested to meet our high standards. Your satisfaction is our commitment.</p>
                </div>
              </div>

              <div className="flex gap-6 group animate-fade-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex-shrink-0 w-16 h-16 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <TrendingUp className="h-8 w-8 text-white" />
                </div>
                <div>
                  <h3 className="text-2xl font-bold text-foreground mb-3">Trending & Innovative</h3>
                  <p className="text-muted-foreground text-lg leading-relaxed">Stay ahead with the latest trends and innovative products that define tomorrow's lifestyle.</p>
                </div>
              </div>

              <div className="flex gap-6 group animate-fade-in" style={{ animationDelay: '0.4s' }}>
                <div className="flex-shrink-0 w-16 h-16 bg-gradient-to-br from-purple-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <Star className="h-8 w-8 text-white" />
                </div>
                <div>
                  <h3 className="text-2xl font-bold text-foreground mb-3">Exceptional Experience</h3>
                  <p className="text-muted-foreground text-lg leading-relaxed">From browsing to delivery, every touchpoint is designed to delight and exceed expectations.</p>
                </div>
              </div>
            </div>

            {/* Right side - Visual */}
            <div className="relative animate-scale-in" style={{ animationDelay: '0.6s' }}>
              <div className="relative bg-gradient-to-br from-primary-50 to-violet-50 rounded-3xl p-8 shadow-2xl">
                <div className="grid grid-cols-2 gap-6">
                  <div className="bg-white rounded-2xl p-6 shadow-medium">
                    <div className="text-3xl font-bold text-primary mb-2">4.9★</div>
                    <div className="text-sm text-muted-foreground">Customer Rating</div>
                  </div>
                  <div className="bg-white rounded-2xl p-6 shadow-medium">
                    <div className="text-3xl font-bold text-emerald-600 mb-2">24h</div>
                    <div className="text-sm text-muted-foreground">Fast Delivery</div>
                  </div>
                  <div className="bg-white rounded-2xl p-6 shadow-medium">
                    <div className="text-3xl font-bold text-blue-600 mb-2">100%</div>
                    <div className="text-sm text-muted-foreground">Secure Payment</div>
                  </div>
                  <div className="bg-white rounded-2xl p-6 shadow-medium">
                    <div className="text-3xl font-bold text-purple-600 mb-2">30d</div>
                    <div className="text-sm text-muted-foreground">Return Policy</div>
                  </div>
                </div>

                {/* Decorative elements */}
                <div className="absolute -top-4 -right-4 w-20 h-20 bg-gradient-to-br from-yellow-400/30 to-orange-500/30 rounded-full blur-xl animate-pulse"></div>
                <div className="absolute -bottom-4 -left-4 w-16 h-16 bg-gradient-to-br from-blue-400/30 to-purple-500/30 rounded-full blur-lg animate-pulse" style={{ animationDelay: '1s' }}></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Products Section */}
      <section className="py-24 bg-background relative">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16 animate-slide-up">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center">
                <Award className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">FEATURED COLLECTION</span>
            </div>

            <h2 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
              Handpicked for
              <span className="text-gradient"> Excellence</span>
            </h2>
            <p className="text-xl text-muted-foreground max-w-4xl mx-auto leading-relaxed">
              Discover our carefully curated selection of premium products, chosen for their exceptional quality,
              outstanding value, and proven customer satisfaction.
            </p>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
              {[...Array(8)].map((_, i) => (
                <Card key={i} variant="elevated" className="animate-pulse">
                  <div className="aspect-square bg-gradient-to-br from-muted to-muted/50 rounded-t-xl"></div>
                  <CardContent className="p-6">
                    <div className="h-4 bg-muted rounded-lg mb-3"></div>
                    <div className="h-4 bg-muted rounded-lg w-2/3 mb-3"></div>
                    <div className="h-6 bg-muted rounded-lg w-1/2"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <div className="text-center py-16 animate-fade-in">
              <div className="w-16 h-16 rounded-full bg-destructive/10 flex items-center justify-center mx-auto mb-6">
                <TrendingUp className="h-8 w-8 text-destructive" />
              </div>
              <h3 className="text-2xl font-bold text-foreground mb-4">Unable to Load Products</h3>
              <p className="text-muted-foreground mb-6 max-w-md mx-auto">
                We're having trouble connecting to our servers. Please check your connection or try again later.
              </p>
              <div className="flex justify-center gap-4">
                <Button variant="gradient" asChild>
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" asChild>
                  <Link href="/test">View Test Page</Link>
                </Button>
              </div>
            </div>
          ) : featuredProducts && featuredProducts.length > 0 ? (
            <div className="animate-fade-in">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
                {featuredProducts.map((product, index) => (
                  <div
                    key={product.id}
                    className="animate-scale-in"
                    style={{ animationDelay: `${index * 0.1}s` }}
                  >
                    <ProductCard product={product} />
                  </div>
                ))}
              </div>

              <div className="text-center mt-16">
                <Button size="xl" variant="gradient" className="group" asChild>
                  <Link href="/products">
                    View All Products
                    <ArrowRight className="ml-2 h-5 w-5 group-hover:translate-x-1 transition-transform" />
                  </Link>
                </Button>
              </div>
            </div>
          ) : (
            <div className="text-center py-16 animate-fade-in">
              <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mx-auto mb-6">
                <Sparkles className="h-8 w-8 text-muted-foreground" />
              </div>
              <h3 className="text-2xl font-bold text-foreground mb-4">Coming Soon</h3>
              <p className="text-muted-foreground mb-6 max-w-md mx-auto">
                We're preparing amazing products for you. Check back soon or browse our existing collection.
              </p>
              <div className="flex justify-center gap-4">
                <Button variant="gradient" asChild>
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" asChild>
                  <Link href="/categories">View Categories</Link>
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-24 bg-gradient-to-br from-muted/30 to-primary-50/20">
        <div className="container mx-auto px-4">
          <div className="text-center mb-16 animate-slide-up">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-violet-500 to-primary-600 flex items-center justify-center">
                <TrendingUp className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">EXPLORE CATEGORIES</span>
            </div>

            <h2 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
              Shop by
              <span className="text-gradient"> Category</span>
            </h2>
            <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
              Discover our carefully organized product categories, designed to help you find exactly what you're looking for.
            </p>
          </div>

          {categoriesLoading ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
              {[...Array(4)].map((_, i) => (
                <Card key={i} variant="elevated" className="animate-pulse">
                  <div className="aspect-square bg-gradient-to-br from-muted to-muted/50 rounded-xl"></div>
                  <CardContent className="p-4">
                    <div className="h-6 bg-muted rounded-lg"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : categories && categories.length > 0 ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
              {categories.slice(0, 8).map((category, index) => (
                <Link
                  key={category.id}
                  href={`/categories/${category.id}`}
                  className="group animate-scale-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <Card variant="elevated" className="overflow-hidden card-hover border-0">
                    <div className="aspect-square bg-gradient-to-br from-primary-100 to-violet-100 relative">
                      <div className="absolute inset-0 bg-gradient-to-t from-primary-600/80 via-primary-500/40 to-transparent"></div>
                      <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-16 h-16 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center">
                          <TrendingUp className="h-8 w-8 text-white" />
                        </div>
                      </div>
                      <div className="absolute bottom-4 left-4 right-4 text-white">
                        <h3 className="text-lg font-bold group-hover:text-yellow-300 transition-colors">
                          {category.name}
                        </h3>
                        {category.description && (
                          <p className="text-sm text-white/80 mt-1 line-clamp-2">
                            {category.description}
                          </p>
                        )}
                      </div>
                    </div>
                  </Card>
                </Link>
              ))}
            </div>
          ) : (
            <div className="text-center py-12">
              <div className="w-16 h-16 rounded-full bg-muted flex items-center justify-center mx-auto mb-6">
                <TrendingUp className="h-8 w-8 text-muted-foreground" />
              </div>
              <h3 className="text-2xl font-bold text-foreground mb-4">Categories Coming Soon</h3>
              <p className="text-muted-foreground">We're organizing our products into categories for better browsing.</p>
            </div>
          )}

          {categories && categories.length > 8 && (
            <div className="text-center mt-12">
              <Button size="xl" variant="outline" asChild>
                <Link href="/categories">
                  View All Categories
                  <ArrowRight className="ml-2 h-5 w-5" />
                </Link>
              </Button>
            </div>
          )}
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="py-24 hero-gradient text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />

        <div className="container mx-auto px-4 text-center relative z-10">
          <div className="max-w-4xl mx-auto animate-fade-in">
            <div className="flex items-center justify-center gap-3 mb-8">
              <div className="w-12 h-12 rounded-xl bg-white/20 backdrop-blur-sm flex items-center justify-center">
                <Sparkles className="h-6 w-6 text-white" />
              </div>
              <span className="text-white/90 font-semibold">EXCLUSIVE OFFERS</span>
            </div>

            <h2 className="text-4xl lg:text-6xl font-bold mb-6">
              Stay Updated with Our
              <span className="block text-gradient bg-gradient-to-r from-white via-violet-200 to-white bg-clip-text text-transparent">
                Latest Offers
              </span>
            </h2>

            <p className="text-xl text-white/90 mb-12 max-w-3xl mx-auto leading-relaxed">
              Join our exclusive community and be the first to discover new products, special promotions,
              and insider deals delivered straight to your inbox.
            </p>

            <div className="max-w-lg mx-auto">
              <div className="flex flex-col sm:flex-row gap-4 p-2 bg-white/10 backdrop-blur-sm rounded-2xl border border-white/20">
                <input
                  type="email"
                  placeholder="Enter your email address"
                  className="flex-1 px-6 py-4 rounded-xl text-foreground placeholder-muted-foreground bg-white border-0 focus:ring-2 focus:ring-violet-400 transition-all"
                />
                <Button variant="secondary" size="lg" className="px-8 py-4 rounded-xl font-semibold shadow-lg hover:shadow-xl transition-all">
                  Subscribe
                </Button>
              </div>

              <p className="text-sm text-white/70 mt-6">
                🔒 We respect your privacy. Unsubscribe at any time.
              </p>

              {/* Trust indicators */}
              <div className="flex items-center justify-center gap-8 mt-8 text-white/60">
                <div className="flex items-center gap-2">
                  <Shield className="h-4 w-4" />
                  <span className="text-sm">Secure</span>
                </div>
                <div className="flex items-center gap-2">
                  <Sparkles className="h-4 w-4" />
                  <span className="text-sm">Exclusive</span>
                </div>
                <div className="flex items-center gap-2">
                  <Award className="h-4 w-4" />
                  <span className="text-sm">Premium</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
