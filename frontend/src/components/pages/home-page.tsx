'use client'

import Link from 'next/link'
import { ArrowRight, Star, Truck, Shield, CreditCard, Sparkles, TrendingUp, Award } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ProductCard } from '@/components/products/product-card'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useFeaturedProducts } from '@/hooks/use-products'
import { useCategories } from '@/hooks/use-categories'
import { APP_CONFIG } from '@/constants/app'
import { DESIGN_TOKENS } from '@/constants/design-tokens'

export function HomePage() {
  const { data, isLoading, error } = useFeaturedProducts(8)
  const featuredProducts = data?.data || []
  const { data: categories, isLoading: categoriesLoading } = useCategories()

  return (
    <div className="min-h-screen bg-black text-white">

      {/* Hero Section - More compact and professional */}
      <section className="relative bg-gradient-to-br from-black via-gray-900 to-black text-white overflow-hidden min-h-[50vh] flex items-center py-8">
        <AnimatedBackground variant="hero" />

        <div className="container mx-auto px-4 relative z-10">
          <div className="max-w-5xl mx-auto">
            <div className="grid lg:grid-cols-2 gap-6 items-center">
              {/* Left Content */}
              <div className="animate-fade-in">
                <div className="flex items-center gap-2 mb-3">
                  <div className="flex items-center gap-2 rounded-full px-3 py-1 border text-xs" style={{backgroundColor: 'rgba(255, 144, 0, 0.2)', borderColor: 'rgba(255, 144, 0, 0.3)'}}>
                    <Sparkles className="h-3 w-3 animate-pulse" style={{color: '#FF9000'}} />
                    <span className="font-medium" style={{color: '#FF9000'}}>New Collection</span>
                  </div>
                </div>

                <h1 className="text-3xl lg:text-4xl xl:text-5xl font-bold mb-3 leading-tight">
                  <span className="block">Discover</span>
                  <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 50%, #FF9000 100%)', WebkitBackgroundClip: 'text'}}>
                    Amazing Products
                  </span>
                  <span className="block">at Great Prices</span>
                </h1>

                <p className="text-base text-gray-300 mb-5 leading-relaxed max-w-lg">
                  Experience premium shopping with curated collections, lightning-fast delivery, and exceptional customer care.
                </p>

                <div className="flex flex-col sm:flex-row gap-3 mb-5">
                  <Button size="default" className="shadow-lg hover:shadow-xl transition-all duration-300 group text-white" style={{backgroundColor: '#FF9000'}} asChild
                    onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                    onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
                  >
                    <Link href="/products">
                      <Sparkles className="mr-2 h-4 w-4 group-hover:rotate-12 transition-transform" />
                      Shop Collection
                      <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                    </Link>
                  </Button>
                  <Button size="default" variant="outline" className="border-2 text-white shadow-lg backdrop-blur-sm" style={{borderColor: 'rgba(255, 144, 0, 0.6)', color: '#FF9000'}} asChild
                    onMouseEnter={(e) => {
                      e.currentTarget.style.backgroundColor = 'rgba(255, 144, 0, 0.1)';
                      e.currentTarget.style.color = '#FF9000';
                    }}
                    onMouseLeave={(e) => {
                      e.currentTarget.style.backgroundColor = 'transparent';
                      e.currentTarget.style.color = '#FF9000';
                    }}
                  >
                    <Link href="/categories">
                      <TrendingUp className="mr-2 h-4 w-4" />
                      Explore Categories
                    </Link>
                  </Button>
                </div>

                {/* Trust Indicators - More compact */}
                <div className="flex flex-wrap items-center gap-3 text-gray-300">
                  <div className="flex items-center gap-1.5">
                    <Shield className="h-3.5 w-3.5" style={{color: '#FF9000'}} />
                    <span className="text-xs font-medium">Secure Shopping</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <Award className="h-3.5 w-3.5" style={{color: '#FF9000'}} />
                    <span className="text-xs font-medium">Premium Quality</span>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <Truck className="h-3.5 w-3.5" style={{color: '#FF9000'}} />
                    <span className="text-xs font-medium">Fast Delivery</span>
                  </div>
                </div>
              </div>

              {/* Right Content - Stats & Visual */}
              <div className="animate-scale-in" style={{ animationDelay: '0.3s' }}>
                <div className="relative">
                  {/* Floating Stats Cards - More compact */}
                  <div className="grid grid-cols-2 gap-2">
                    <div className="glass-effect p-3 rounded-lg border animate-float" style={{borderColor: 'rgba(255, 144, 0, 0.2)'}}>
                      <div className="text-base lg:text-lg font-semibold mb-0.5 bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}>10K+</div>
                      <div className="text-gray-300 text-xs font-medium">Happy Customers</div>
                    </div>
                    <div className="glass-effect p-3 rounded-lg border border-gray-600 animate-float" style={{ animationDelay: '0.5s' }}>
                      <div className="text-base lg:text-lg font-semibold mb-0.5 text-white">5K+</div>
                      <div className="text-gray-300 text-xs font-medium">Premium Products</div>
                    </div>
                    <div className="glass-effect p-3 rounded-lg border border-gray-600 animate-float" style={{ animationDelay: '1s' }}>
                      <div className="text-base lg:text-lg font-semibold mb-0.5 text-white">99%</div>
                      <div className="text-gray-300 text-xs font-medium">Satisfaction Rate</div>
                    </div>
                    <div className="glass-effect p-3 rounded-lg border border-gray-600 animate-float" style={{ animationDelay: '1.5s' }}>
                      <div className="text-base lg:text-lg font-semibold mb-0.5 text-white">24/7</div>
                      <div className="text-gray-300 text-xs font-medium">Customer Support</div>
                    </div>
                  </div>

                  {/* Decorative elements */}
                  <div className="absolute -top-4 -right-4 w-20 h-20 rounded-full blur-xl animate-pulse" style={{background: 'linear-gradient(135deg, rgba(255, 144, 0, 0.2) 0%, rgba(255, 119, 0, 0.2) 100%)'}}></div>
                  <div className="absolute -bottom-4 -left-4 w-24 h-24 bg-gradient-to-br from-blue-400/20 to-purple-500/20 rounded-full blur-xl animate-pulse" style={{ animationDelay: '1s' }}></div>
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
      <section className="py-12 bg-gradient-to-br from-gray-900 via-black to-gray-900 relative overflow-hidden">
        <div className="absolute inset-0 bg-grid-pattern opacity-5"></div>

        <div className="container mx-auto px-4 relative z-10">
          <div className="text-center mb-8 animate-slide-up">
            <div className="flex items-center justify-center gap-2 mb-3">
              <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{background: 'linear-gradient(135deg, #FF9000 0%, #e67e00 100%)'}}>
                <Sparkles className="h-4 w-4 text-white" />
              </div>
              <span className="font-semibold text-sm" style={{color: '#FF9000'}}>WHY CHOOSE US</span>
            </div>

            <h2 className="text-2xl lg:text-3xl font-bold text-white mb-3">
              Experience the
              <span className="bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Difference</span>
            </h2>
            <p className="text-base text-gray-300 max-w-2xl mx-auto leading-relaxed">
              We're not just another online store. We're your trusted partner in discovering exceptional products
              that enhance your lifestyle and exceed your expectations.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 items-center">
            {/* Left side - Benefits */}
            <div className="space-y-4">
              <div className="flex gap-3 group animate-fade-in">
                <div className="flex-shrink-0 w-10 h-10 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-lg flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <Award className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white mb-1">Premium Quality Guarantee</h3>
                  <p className="text-gray-300 text-sm leading-relaxed">Every product is carefully curated and tested to meet our high standards. Your satisfaction is our commitment.</p>
                </div>
              </div>

              <div className="flex gap-3 group animate-fade-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex-shrink-0 w-10 h-10 bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <TrendingUp className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white mb-1">Trending & Innovative</h3>
                  <p className="text-gray-300 text-sm leading-relaxed">Stay ahead with the latest trends and innovative products that define tomorrow's lifestyle.</p>
                </div>
              </div>

              <div className="flex gap-3 group animate-fade-in" style={{ animationDelay: '0.4s' }}>
                <div className="flex-shrink-0 w-10 h-10 bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <Star className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-base font-bold text-white mb-1">Exceptional Experience</h3>
                  <p className="text-gray-300 text-sm leading-relaxed">From browsing to delivery, every touchpoint is designed to delight and exceed expectations.</p>
                </div>
              </div>
            </div>

            {/* Right side - Visual */}
            <div className="relative animate-scale-in" style={{ animationDelay: '0.6s' }}>
              <div className="relative bg-gradient-to-br from-gray-800 to-gray-900 rounded-lg p-4 shadow-xl border border-gray-700">
                <div className="grid grid-cols-2 gap-3">
                  <div className="bg-gray-700 rounded-lg p-3 shadow-md">
                    <div className="text-lg font-bold mb-1" style={{color: '#FF9000'}}>4.9★</div>
                    <div className="text-xs text-gray-300">Customer Rating</div>
                  </div>
                  <div className="bg-gray-700 rounded-lg p-3 shadow-md">
                    <div className="text-lg font-bold text-emerald-400 mb-1">24h</div>
                    <div className="text-xs text-gray-300">Fast Delivery</div>
                  </div>
                  <div className="bg-gray-700 rounded-lg p-3 shadow-md">
                    <div className="text-lg font-bold text-blue-400 mb-1">100%</div>
                    <div className="text-xs text-gray-300">Secure Payment</div>
                  </div>
                  <div className="bg-gray-700 rounded-lg p-3 shadow-md">
                    <div className="text-lg font-bold text-purple-400 mb-1">30d</div>
                    <div className="text-xs text-gray-300">Return Policy</div>
                  </div>
                </div>

                {/* Decorative elements */}
                <div className="absolute -top-2 -right-2 w-12 h-12 rounded-full blur-lg animate-pulse" style={{background: 'linear-gradient(135deg, rgba(255, 144, 0, 0.3) 0%, rgba(255, 119, 0, 0.3) 100%)'}}></div>
                <div className="absolute -bottom-2 -left-2 w-10 h-10 bg-gradient-to-br from-blue-400/30 to-purple-500/30 rounded-full blur-lg animate-pulse" style={{ animationDelay: '1s' }}></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Products Section - More compact */}
      <section className="py-8 bg-gray-900 relative">
        <div className="container mx-auto px-4">
          <div className="text-center mb-6 animate-slide-up">
            <div className="flex items-center justify-center gap-2 mb-2">
              <div className="w-6 h-6 rounded-lg flex items-center justify-center" style={{background: 'linear-gradient(135deg, #FF9000 0%, #e67e00 100%)'}}>
                <Award className="h-3 w-3 text-white" />
              </div>
              <span className="font-semibold text-xs" style={{color: '#FF9000'}}>FEATURED COLLECTION</span>
            </div>

            <h2 className="text-xl lg:text-2xl font-bold text-white mb-2">
              Handpicked for
              <span className="bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Excellence</span>
            </h2>
            <p className="text-sm text-gray-300 max-w-xl mx-auto leading-relaxed">
              Discover our carefully curated selection of premium products, chosen for their exceptional quality and value.
            </p>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
              {[...Array(8)].map((_, i) => (
                <Card key={i} className="animate-pulse bg-gray-800 border-gray-700">
                  <div className="aspect-square bg-gradient-to-br from-gray-700 to-gray-600 rounded-t-lg"></div>
                  <CardContent className="p-3">
                    <div className="h-3 bg-gray-600 rounded mb-1.5"></div>
                    <div className="h-3 bg-gray-600 rounded w-2/3 mb-1.5"></div>
                    <div className="h-4 bg-gray-600 rounded w-1/2"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <div className="text-center py-8 animate-fade-in">
              <div className="w-10 h-10 rounded-full bg-red-500/10 flex items-center justify-center mx-auto mb-3">
                <TrendingUp className="h-5 w-5 text-red-400" />
              </div>
              <h3 className="text-lg font-bold text-white mb-2">Unable to Load Products</h3>
              <p className="text-gray-300 mb-4 max-w-md mx-auto text-sm">
                We're having trouble connecting to our servers. Please check your connection or try again later.
              </p>
              <div className="flex justify-center gap-3">
                <Button variant="gradient" size="default" asChild>
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" size="default" className="border-gray-600 text-gray-300" asChild>
                  <Link href="/test">View Test Page</Link>
                </Button>
              </div>
            </div>
          ) : featuredProducts.length > 0 ? (
            <div className="animate-fade-in">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
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

              <div className="text-center mt-8">
                <Button size="lg" variant="gradient" className="group" asChild>
                  <Link href="/products">
                    View All Products
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                  </Link>
                </Button>
              </div>
            </div>
          ) : (
            <div className="text-center py-8 animate-fade-in">
              <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-3">
                <Sparkles className="h-5 w-5 text-gray-400" />
              </div>
              <h3 className="text-lg font-bold text-white mb-2">Coming Soon</h3>
              <p className="text-gray-300 mb-4 max-w-md mx-auto text-sm">
                We're preparing amazing products for you. Check back soon or browse our existing collection.
              </p>
              <div className="flex justify-center gap-3">
                <Button variant="gradient" size="default" asChild>
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" size="default" className="border-gray-600 text-gray-300" asChild>
                  <Link href="/categories">View Categories</Link>
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-12 bg-gradient-to-br from-gray-800 to-black">
        <div className="container mx-auto px-4">
          <div className="text-center mb-8 animate-slide-up">
            <div className="flex items-center justify-center gap-2 mb-3">
              <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{background: 'linear-gradient(135deg, #FF9000 0%, #e67e00 100%)'}}>
                <TrendingUp className="h-4 w-4 text-white" />
              </div>
              <span className="font-semibold text-sm" style={{color: '#FF9000'}}>EXPLORE CATEGORIES</span>
            </div>

            <h2 className="text-2xl lg:text-3xl font-bold text-white mb-3">
              Shop by
              <span className="bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}> Category</span>
            </h2>
            <p className="text-base text-gray-300 max-w-2xl mx-auto leading-relaxed">
              Discover our carefully organized product categories, designed to help you find exactly what you're looking for.
            </p>
          </div>

          {categoriesLoading ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[...Array(4)].map((_, i) => (
                <Card key={i} className="animate-pulse bg-gray-800 border-gray-700">
                  <div className="aspect-square bg-gradient-to-br from-gray-700 to-gray-600 rounded-lg"></div>
                  <CardContent className="p-3">
                    <div className="h-4 bg-gray-600 rounded"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : categories && categories.length > 0 ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {categories.slice(0, 8).map((category, index) => (
                <Link
                  key={category.id}
                  href={`/categories/${category.id}`}
                  className="group animate-scale-in"
                  style={{ animationDelay: `${index * 0.1}s` }}
                >
                  <Card className="overflow-hidden card-hover border-gray-700 bg-gray-800 hover:bg-gray-700 transition-all duration-300">
                    <div className="aspect-square bg-gradient-to-br from-gray-700 to-gray-600 relative">
                      <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
                      <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-10 h-10 rounded-lg bg-white/20 backdrop-blur-sm flex items-center justify-center">
                          <TrendingUp className="h-5 w-5 text-white" />
                        </div>
                      </div>
                      <div className="absolute bottom-3 left-3 right-3 text-white">
                        <h3 className="text-sm font-bold group-hover:text-orange-300 transition-colors">
                          {category.name}
                        </h3>
                        {category.description && (
                          <p className="text-xs text-white/80 mt-1 line-clamp-2">
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
            <div className="text-center py-8">
              <div className="w-10 h-10 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-3">
                <TrendingUp className="h-5 w-5 text-gray-400" />
              </div>
              <h3 className="text-lg font-bold text-white mb-2">Categories Coming Soon</h3>
              <p className="text-gray-300 text-sm">We're organizing our products into categories for better browsing.</p>
            </div>
          )}

          {categories && categories.length > 8 && (
            <div className="text-center mt-6">
              <Button size="lg" variant="outline" className="border-gray-600 text-gray-300" asChild>
                <Link href="/categories">
                  View All Categories
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          )}
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="py-12 bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
        <AnimatedBackground className="opacity-30" />

        <div className="container mx-auto px-4 text-center relative z-10">
          <div className="max-w-2xl mx-auto animate-fade-in">
            <div className="flex items-center justify-center gap-2 mb-4">
              <div className="w-8 h-8 rounded-lg bg-white/20 backdrop-blur-sm flex items-center justify-center">
                <Sparkles className="h-4 w-4 text-white" />
              </div>
              <span className="text-white/90 font-semibold text-sm">EXCLUSIVE OFFERS</span>
            </div>

            <h2 className="text-2xl lg:text-3xl font-bold mb-3">
              Stay Updated with Our
              <span className="block bg-clip-text text-transparent" style={{background: 'linear-gradient(90deg, #FF9000 0%, #ff7700 100%)', WebkitBackgroundClip: 'text'}}>
                Latest Offers
              </span>
            </h2>

            <p className="text-base text-gray-300 mb-6 max-w-xl mx-auto leading-relaxed">
              Join our exclusive community and be the first to discover new products, special promotions,
              and insider deals delivered straight to your inbox.
            </p>

            <div className="max-w-md mx-auto">
              <div className="flex flex-col sm:flex-row gap-3 p-2 bg-white/10 backdrop-blur-sm rounded-lg border border-white/20">
                <input
                  type="email"
                  placeholder="Enter your email address"
                  className="flex-1 px-4 py-2 rounded-lg text-gray-900 placeholder-gray-500 bg-white border-0 focus:ring-2 transition-all text-sm"
                  style={{'--tw-ring-color': '#FF9000'} as React.CSSProperties}
                  onFocus={(e) => e.currentTarget.style.boxShadow = '0 0 0 2px rgba(255, 144, 0, 0.5)'}
                  onBlur={(e) => e.currentTarget.style.boxShadow = ''}
                />
                <Button className="px-6 py-2 rounded-lg font-semibold shadow-lg hover:shadow-xl transition-all text-sm text-white" style={{backgroundColor: '#FF9000'}}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
                >
                  Subscribe
                </Button>
              </div>

              <p className="text-xs text-white/70 mt-3">
                🔒 We respect your privacy. Unsubscribe at any time.
              </p>

              {/* Trust indicators */}
              <div className="flex items-center justify-center gap-4 mt-4 text-white/60">
                <div className="flex items-center gap-1">
                  <Shield className="h-3 w-3" />
                  <span className="text-xs">Secure</span>
                </div>
                <div className="flex items-center gap-1">
                  <Sparkles className="h-3 w-3" />
                  <span className="text-xs">Exclusive</span>
                </div>
                <div className="flex items-center gap-1">
                  <Award className="h-3 w-3" />
                  <span className="text-xs">Premium</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
