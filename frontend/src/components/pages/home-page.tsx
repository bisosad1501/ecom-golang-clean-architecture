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
import { DESIGN_TOKENS } from '@/constants/design-tokens'

export function HomePage() {
  const { data: featuredProducts, isLoading, error } = useFeaturedProducts(8)
  const { data: categories, isLoading: categoriesLoading } = useCategories()

  return (
    <div className="min-h-screen bg-black text-white">

      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-black via-gray-900 to-black text-white overflow-hidden min-h-[80vh] flex items-center">
        <AnimatedBackground variant="hero" />

        <div className={`container mx-auto px-4 ${DESIGN_TOKENS.SPACING.SECTION_DEFAULT} relative z-10`}>
          <div className="max-w-5xl mx-auto">
            <div className={`grid lg:grid-cols-2 ${DESIGN_TOKENS.SPACING.GAP_DEFAULT} items-center`}>
              {/* Left Content */}
              <div className="animate-fade-in">
                <div className={`flex items-center ${DESIGN_TOKENS.SPACING.GAP_SMALL} ${DESIGN_TOKENS.SPACING.MARGIN_DEFAULT}`}>
                  <div className={`flex items-center gap-2 bg-orange-500/20 backdrop-blur-md ${DESIGN_TOKENS.RADIUS.FULL} px-4 py-1.5 border border-orange-500/30`}>
                    <Sparkles className={`${DESIGN_TOKENS.ICONS.SMALL} text-orange-400 animate-pulse`} />
                    <span className={`${DESIGN_TOKENS.TYPOGRAPHY.LABEL} text-orange-300`}>New Collection Available</span>
                  </div>
                </div>

                <h1 className={`${DESIGN_TOKENS.TYPOGRAPHY.DISPLAY_LARGE} ${DESIGN_TOKENS.SPACING.MARGIN_DEFAULT} leading-[0.9]`}>
                  <span className="block">Discover</span>
                  <span className="block text-gradient bg-gradient-to-r from-orange-500 via-orange-400 to-orange-500 bg-clip-text text-transparent animate-gradient">
                    Amazing Products
                  </span>
                  <span className="block">at Great Prices</span>
                </h1>

                <p className={`${DESIGN_TOKENS.TYPOGRAPHY.BODY_LARGE} ${DESIGN_TOKENS.SPACING.MARGIN_DEFAULT} text-white/90 leading-relaxed max-w-xl`}>
                  Experience premium shopping with curated collections, lightning-fast delivery, and exceptional customer care.
                </p>

                <div className={`flex flex-col sm:flex-row ${DESIGN_TOKENS.SPACING.GAP_SMALL} ${DESIGN_TOKENS.SPACING.MARGIN_DEFAULT}`}>
                  <Button size="lg" variant="secondary" className="shadow-2xl hover:shadow-3xl transition-all duration-300 group bg-orange-500 text-white hover:bg-orange-600" asChild>
                    <Link href="/products">
                      <Sparkles className={`mr-2 ${DESIGN_TOKENS.ICONS.DEFAULT} group-hover:rotate-12 transition-transform`} />
                      Shop Collection
                      <ArrowRight className={`ml-2 ${DESIGN_TOKENS.ICONS.DEFAULT} group-hover:translate-x-1 transition-transform`} />
                    </Link>
                  </Button>
                  <Button size="lg" variant="outline" className="border-2 border-orange-500/60 text-orange-400 hover:bg-orange-500/10 shadow-xl backdrop-blur-sm" asChild>
                    <Link href="/categories">
                      <TrendingUp className={`mr-2 ${DESIGN_TOKENS.ICONS.DEFAULT}`} />
                      Explore Categories
                    </Link>
                  </Button>
                </div>

                {/* Trust Indicators */}
                <div className={`flex items-center ${DESIGN_TOKENS.SPACING.GAP_DEFAULT} text-white/80`}>
                  <div className={`flex items-center ${DESIGN_TOKENS.SPACING.GAP_TINY}`}>
                    <Shield className={`${DESIGN_TOKENS.ICONS.DEFAULT} text-orange-400`} />
                    <span className={DESIGN_TOKENS.TYPOGRAPHY.LABEL}>Secure Shopping</span>
                  </div>
                  <div className={`flex items-center ${DESIGN_TOKENS.SPACING.GAP_TINY}`}>
                    <Award className={`${DESIGN_TOKENS.ICONS.DEFAULT} text-orange-400`} />
                    <span className={DESIGN_TOKENS.TYPOGRAPHY.LABEL}>Premium Quality</span>
                  </div>
                  <div className={`flex items-center ${DESIGN_TOKENS.SPACING.GAP_TINY}`}>
                    <Truck className={`${DESIGN_TOKENS.ICONS.DEFAULT} text-orange-400`} />
                    <span className={DESIGN_TOKENS.TYPOGRAPHY.LABEL}>Fast Delivery</span>
                  </div>
                </div>
              </div>

              {/* Right Content - Stats & Visual */}
              <div className="animate-scale-in" style={{ animationDelay: '0.3s' }}>
                <div className="relative">
                  {/* Floating Stats Cards */}
                  <div className={`grid grid-cols-2 ${DESIGN_TOKENS.SPACING.GAP_SMALL}`}>
                    <div className={`glass-effect ${DESIGN_TOKENS.CONTAINERS.CARD_PADDING} ${DESIGN_TOKENS.RADIUS.DEFAULT} border border-orange-500/20 animate-float`}>
                      <div className={`${DESIGN_TOKENS.TYPOGRAPHY.HEADING_2} mb-1 text-gradient bg-gradient-to-r from-orange-400 to-orange-500 bg-clip-text text-transparent`}>10K+</div>
                      <div className={`text-white/80 ${DESIGN_TOKENS.TYPOGRAPHY.LABEL}`}>Happy Customers</div>
                    </div>
                    <div className={`glass-effect ${DESIGN_TOKENS.CONTAINERS.CARD_PADDING} ${DESIGN_TOKENS.RADIUS.DEFAULT} border border-white/20 animate-float`} style={{ animationDelay: '0.5s' }}>
                      <div className={`${DESIGN_TOKENS.TYPOGRAPHY.HEADING_2} mb-1 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent`}>5K+</div>
                      <div className={`text-white/80 ${DESIGN_TOKENS.TYPOGRAPHY.LABEL}`}>Premium Products</div>
                    </div>
                    <div className={`glass-effect ${DESIGN_TOKENS.CONTAINERS.CARD_PADDING} ${DESIGN_TOKENS.RADIUS.DEFAULT} border border-white/20 animate-float`} style={{ animationDelay: '1s' }}>
                      <div className={`${DESIGN_TOKENS.TYPOGRAPHY.HEADING_2} mb-1 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent`}>99%</div>
                      <div className={`text-white/80 ${DESIGN_TOKENS.TYPOGRAPHY.LABEL}`}>Satisfaction Rate</div>
                    </div>
                    <div className={`glass-effect ${DESIGN_TOKENS.CONTAINERS.CARD_PADDING} ${DESIGN_TOKENS.RADIUS.DEFAULT} border border-white/20 animate-float`} style={{ animationDelay: '1.5s' }}>
                      <div className={`${DESIGN_TOKENS.TYPOGRAPHY.HEADING_2} mb-1 text-gradient bg-gradient-to-r from-white to-violet-200 bg-clip-text text-transparent`}>24/7</div>
                      <div className={`text-white/80 ${DESIGN_TOKENS.TYPOGRAPHY.LABEL}`}>Customer Support</div>
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
      <section className="py-16 bg-gradient-to-br from-gray-900 via-black to-gray-900 relative overflow-hidden">{/* Dark background */}
        <div className="absolute inset-0 bg-grid-pattern opacity-5"></div>

        <div className="container mx-auto px-4 relative z-10">
          <div className="text-center mb-12 animate-slide-up">{/* Reduced margin */}
            <div className="flex items-center justify-center gap-2 mb-4">{/* Smaller gaps and margin */}
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-orange-500 to-orange-600 flex items-center justify-center">{/* Orange background */}
                <Sparkles className="h-5 w-5 text-white" />
              </div>
              <span className="text-orange-500 font-semibold text-sm">WHY CHOOSE US</span>{/* Orange text */}
            </div>

            <h2 className="text-2xl lg:text-4xl font-bold text-white mb-4">{/* White text */}
              Experience the
              <span className="text-gradient bg-gradient-to-r from-orange-400 to-orange-500 bg-clip-text text-transparent"> Difference</span>
            </h2>
            <p className="text-base text-gray-300 max-w-3xl mx-auto leading-relaxed">{/* Light gray text */}
              We're not just another online store. We're your trusted partner in discovering exceptional products
              that enhance your lifestyle and exceed your expectations.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">{/* Reduced gap */}
            {/* Left side - Benefits */}
            <div className="space-y-6">{/* Reduced spacing */}
              <div className="flex gap-4 group animate-fade-in">{/* Smaller gap */}
                <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-emerald-500 to-emerald-600 rounded-xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">{/* Smaller icon container */}
                  <Award className="h-6 w-6 text-white" />{/* Smaller icon */}
                </div>
                <div>
                  <h3 className="text-lg font-bold text-foreground mb-2">Premium Quality Guarantee</h3>{/* Smaller heading and margin */}
                  <p className="text-muted-foreground text-sm leading-relaxed">Every product is carefully curated and tested to meet our high standards. Your satisfaction is our commitment.</p>{/* Smaller text */}
                </div>
              </div>

              <div className="flex gap-4 group animate-fade-in" style={{ animationDelay: '0.2s' }}>
                <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <TrendingUp className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-foreground mb-2">Trending & Innovative</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed">Stay ahead with the latest trends and innovative products that define tomorrow's lifestyle.</p>
                </div>
              </div>

              <div className="flex gap-4 group animate-fade-in" style={{ animationDelay: '0.4s' }}>
                <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-110 transition-all duration-300">
                  <Star className="h-6 w-6 text-white" />
                </div>
                <div>
                  <h3 className="text-lg font-bold text-foreground mb-2">Exceptional Experience</h3>
                  <p className="text-muted-foreground text-sm leading-relaxed">From browsing to delivery, every touchpoint is designed to delight and exceed expectations.</p>
                </div>
              </div>
            </div>

            {/* Right side - Visual */}
            <div className="relative animate-scale-in" style={{ animationDelay: '0.6s' }}>
              <div className="relative bg-gradient-to-br from-primary-50 to-violet-50 rounded-2xl p-6 shadow-2xl">{/* Smaller border radius and padding */}
                <div className="grid grid-cols-2 gap-4">{/* Smaller gap */}
                  <div className="bg-white rounded-xl p-4 shadow-medium">{/* Smaller padding and border radius */}
                    <div className="text-2xl font-bold text-primary mb-1">4.9★</div>{/* Smaller text and margin */}
                    <div className="text-xs text-muted-foreground">Customer Rating</div>{/* Smaller text */}
                  </div>
                  <div className="bg-white rounded-xl p-4 shadow-medium">
                    <div className="text-2xl font-bold text-emerald-600 mb-1">24h</div>
                    <div className="text-xs text-muted-foreground">Fast Delivery</div>
                  </div>
                  <div className="bg-white rounded-xl p-4 shadow-medium">
                    <div className="text-2xl font-bold text-blue-600 mb-1">100%</div>
                    <div className="text-xs text-muted-foreground">Secure Payment</div>
                  </div>
                  <div className="bg-white rounded-xl p-4 shadow-medium">
                    <div className="text-2xl font-bold text-purple-600 mb-1">30d</div>
                    <div className="text-xs text-muted-foreground">Return Policy</div>
                  </div>
                </div>

                {/* Decorative elements */}
                <div className="absolute -top-2 -right-2 w-16 h-16 bg-gradient-to-br from-yellow-400/30 to-orange-500/30 rounded-full blur-xl animate-pulse"></div>{/* Smaller elements */}
                <div className="absolute -bottom-2 -left-2 w-12 h-12 bg-gradient-to-br from-blue-400/30 to-purple-500/30 rounded-full blur-lg animate-pulse" style={{ animationDelay: '1s' }}></div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Products Section */}
      <section className="py-16 bg-background relative">{/* Reduced padding */}
        <div className="container mx-auto px-4">
          <div className="text-center mb-12 animate-slide-up">{/* Reduced margin */}
            <div className="flex items-center justify-center gap-2 mb-4">{/* Smaller gaps and margin */}
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center">{/* Smaller icon container */}
                <Award className="h-5 w-5 text-white" />{/* Smaller icon */}
              </div>
              <span className="text-primary font-semibold text-sm">FEATURED COLLECTION</span>{/* Smaller text */}
            </div>

            <h2 className="text-2xl lg:text-4xl font-bold text-foreground mb-4">{/* Much smaller heading and margin */}
              Handpicked for
              <span className="text-gradient"> Excellence</span>
            </h2>
            <p className="text-base text-muted-foreground max-w-3xl mx-auto leading-relaxed">{/* Smaller text and max-width */}
              Discover our carefully curated selection of premium products, chosen for their exceptional quality,
              outstanding value, and proven customer satisfaction.
            </p>
          </div>

          {isLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">{/* Reduced gap */}
              {[...Array(8)].map((_, i) => (
                <Card key={i} variant="elevated" className="animate-pulse">
                  <div className="aspect-square bg-gradient-to-br from-muted to-muted/50 rounded-t-xl"></div>
                  <CardContent className="p-4">{/* Smaller padding */}
                    <div className="h-3 bg-muted rounded-lg mb-2"></div>{/* Smaller heights and margins */}
                    <div className="h-3 bg-muted rounded-lg w-2/3 mb-2"></div>
                    <div className="h-4 bg-muted rounded-lg w-1/2"></div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : error ? (
            <div className="text-center py-12 animate-fade-in">{/* Reduced padding */}
              <div className="w-12 h-12 rounded-full bg-destructive/10 flex items-center justify-center mx-auto mb-4">{/* Smaller icon container and margin */}
                <TrendingUp className="h-6 w-6 text-destructive" />{/* Smaller icon */}
              </div>
              <h3 className="text-xl font-bold text-foreground mb-3">Unable to Load Products</h3>{/* Smaller heading and margin */}
              <p className="text-muted-foreground mb-4 max-w-md mx-auto text-sm">{/* Smaller text and margin */}
                We're having trouble connecting to our servers. Please check your connection or try again later.
              </p>
              <div className="flex justify-center gap-3">{/* Smaller gap */}
                <Button variant="gradient" size="default" asChild>{/* Smaller button */}
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" size="default" asChild>
                  <Link href="/test">View Test Page</Link>
                </Button>
              </div>
            </div>
          ) : featuredProducts && featuredProducts.length > 0 ? (
            <div className="animate-fade-in">
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">{/* Reduced gap */}
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

              <div className="text-center mt-12">{/* Reduced margin */}
                <Button size="lg" variant="gradient" className="group" asChild>{/* Smaller button */}
                  <Link href="/products">
                    View All Products
                    <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />{/* Smaller icon */}
                  </Link>
                </Button>
              </div>
            </div>
          ) : (
            <div className="text-center py-12 animate-fade-in">{/* Reduced padding */}
              <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mx-auto mb-4">{/* Smaller icon container and margin */}
                <Sparkles className="h-6 w-6 text-muted-foreground" />{/* Smaller icon */}
              </div>
              <h3 className="text-xl font-bold text-foreground mb-3">Coming Soon</h3>{/* Smaller heading and margin */}
              <p className="text-muted-foreground mb-4 max-w-md mx-auto text-sm">{/* Smaller text and margin */}
                We're preparing amazing products for you. Check back soon or browse our existing collection.
              </p>
              <div className="flex justify-center gap-3">{/* Smaller gap */}
                <Button variant="gradient" size="default" asChild>{/* Smaller button */}
                  <Link href="/products">Browse All Products</Link>
                </Button>
                <Button variant="outline" size="default" asChild>
                  <Link href="/categories">View Categories</Link>
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-16 bg-gradient-to-br from-muted/30 to-primary-50/20">{/* Reduced padding */}
        <div className="container mx-auto px-4">
          <div className="text-center mb-12 animate-slide-up">{/* Reduced margin */}
            <div className="flex items-center justify-center gap-2 mb-4">{/* Smaller gaps and margin */}
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-violet-500 to-primary-600 flex items-center justify-center">{/* Smaller icon container */}
                <TrendingUp className="h-5 w-5 text-white" />{/* Smaller icon */}
              </div>
              <span className="text-primary font-semibold text-sm">EXPLORE CATEGORIES</span>{/* Smaller text */}
            </div>

            <h2 className="text-2xl lg:text-4xl font-bold text-foreground mb-4">{/* Much smaller heading and margin */}
              Shop by
              <span className="text-gradient"> Category</span>
            </h2>
            <p className="text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">{/* Smaller text and max-width */}
              Discover our carefully organized product categories, designed to help you find exactly what you're looking for.
            </p>
          </div>

          {categoriesLoading ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">{/* Reduced gap */}
              {[...Array(4)].map((_, i) => (
                <Card key={i} variant="elevated" className="animate-pulse">
                  <div className="aspect-square bg-gradient-to-br from-muted to-muted/50 rounded-xl"></div>
                  <CardContent className="p-3">{/* Smaller padding */}
                    <div className="h-4 bg-muted rounded-lg"></div>{/* Smaller height */}
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : categories && categories.length > 0 ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-6">{/* Reduced gap */}
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
                        <div className="w-12 h-12 rounded-xl bg-white/20 backdrop-blur-sm flex items-center justify-center">{/* Smaller icon container */}
                          <TrendingUp className="h-6 w-6 text-white" />{/* Smaller icon */}
                        </div>
                      </div>
                      <div className="absolute bottom-3 left-3 right-3 text-white">{/* Smaller positioning */}
                        <h3 className="text-base font-bold group-hover:text-yellow-300 transition-colors">{/* Smaller text */}
                          {category.name}
                        </h3>
                        {category.description && (
                          <p className="text-xs text-white/80 mt-1 line-clamp-2">{/* Smaller text */}
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
            <div className="text-center py-10">{/* Reduced padding */}
              <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mx-auto mb-4">{/* Smaller icon container and margin */}
                <TrendingUp className="h-6 w-6 text-muted-foreground" />{/* Smaller icon */}
              </div>
              <h3 className="text-xl font-bold text-foreground mb-3">Categories Coming Soon</h3>{/* Smaller heading and margin */}
              <p className="text-muted-foreground text-sm">We're organizing our products into categories for better browsing.</p>{/* Smaller text */}
            </div>
          )}

          {categories && categories.length > 8 && (
            <div className="text-center mt-10">{/* Reduced margin */}
              <Button size="lg" variant="outline" asChild>{/* Smaller button */}
                <Link href="/categories">
                  View All Categories
                  <ArrowRight className="ml-2 h-4 w-4" />{/* Smaller icon */}
                </Link>
              </Button>
            </div>
          )}
        </div>
      </section>

      {/* Newsletter Section */}
      <section className="py-16 hero-gradient text-white relative overflow-hidden">{/* Reduced padding */}
        <AnimatedBackground className="opacity-30" />

        <div className="container mx-auto px-4 text-center relative z-10">
          <div className="max-w-3xl mx-auto animate-fade-in">{/* Smaller max-width */}
            <div className="flex items-center justify-center gap-2 mb-6">{/* Smaller gap and margin */}
              <div className="w-10 h-10 rounded-xl bg-white/20 backdrop-blur-sm flex items-center justify-center">{/* Smaller icon container */}
                <Sparkles className="h-5 w-5 text-white" />{/* Smaller icon */}
              </div>
              <span className="text-white/90 font-semibold text-sm">EXCLUSIVE OFFERS</span>{/* Smaller text */}
            </div>

            <h2 className="text-2xl lg:text-4xl font-bold mb-4">{/* Much smaller heading and margin */}
              Stay Updated with Our
              <span className="block text-gradient bg-gradient-to-r from-white via-violet-200 to-white bg-clip-text text-transparent">
                Latest Offers
              </span>
            </h2>

            <p className="text-base text-white/90 mb-8 max-w-2xl mx-auto leading-relaxed">{/* Smaller text, margin and max-width */}
              Join our exclusive community and be the first to discover new products, special promotions,
              and insider deals delivered straight to your inbox.
            </p>

            <div className="max-w-md mx-auto">{/* Smaller max-width */}
              <div className="flex flex-col sm:flex-row gap-3 p-2 bg-white/10 backdrop-blur-sm rounded-xl border border-white/20">{/* Smaller gap and border radius */}
                <input
                  type="email"
                  placeholder="Enter your email address"
                  className="flex-1 px-4 py-3 rounded-lg text-foreground placeholder-muted-foreground bg-white border-0 focus:ring-2 focus:ring-violet-400 transition-all text-sm"
                />
                <Button variant="secondary" size="default" className="px-6 py-3 rounded-lg font-semibold shadow-lg hover:shadow-xl transition-all text-sm">{/* Smaller button and text */}
                  Subscribe
                </Button>
              </div>

              <p className="text-xs text-white/70 mt-4">{/* Smaller text and margin */}
                🔒 We respect your privacy. Unsubscribe at any time.
              </p>

              {/* Trust indicators */}
              <div className="flex items-center justify-center gap-6 mt-6 text-white/60">{/* Smaller gap and margin */}
                <div className="flex items-center gap-1.5">{/* Smaller gap */}
                  <Shield className="h-3.5 w-3.5" />{/* Smaller icon */}
                  <span className="text-xs">Secure</span>{/* Smaller text */}
                </div>
                <div className="flex items-center gap-1.5">
                  <Sparkles className="h-3.5 w-3.5" />
                  <span className="text-xs">Exclusive</span>
                </div>
                <div className="flex items-center gap-1.5">
                  <Award className="h-3.5 w-3.5" />
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
