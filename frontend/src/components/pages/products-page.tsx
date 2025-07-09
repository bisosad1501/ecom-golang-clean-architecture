'use client'

import { useState, useMemo } from 'react'
import { useSearchParams } from 'next/navigation'
import { Filter, Grid, List, SlidersHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProductCard } from '@/components/products/product-card'
import { ProductListCard } from '@/components/products/product-list-card'
import { ProductFilters } from '@/components/products/product-filters'
import { ProductSort } from '@/components/products/product-sort'
import { Pagination } from '@/components/ui/pagination'
import { CategoryBreadcrumbs } from '@/components/ui/category-breadcrumbs'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { useProducts } from '@/hooks/use-products'
import { ProductsParams } from '@/lib/services/products'
import { DEFAULT_PAGE_SIZE, PRODUCT_SORT_OPTIONS } from '@/constants'

export function ProductsPage() {
  const searchParams = useSearchParams()
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [showFilters, setShowFilters] = useState(false)

  // Parse URL parameters
  const params = useMemo((): ProductsParams => {
    const urlParams = new URLSearchParams(searchParams.toString())
    return {
      page: parseInt(urlParams.get('page') || '1'),
      limit: parseInt(urlParams.get('limit') || DEFAULT_PAGE_SIZE.toString()),
      search: urlParams.get('search') || undefined,
      category_id: urlParams.get('category') || undefined,
      min_price: urlParams.get('min_price') ? parseFloat(urlParams.get('min_price')!) : undefined,
      max_price: urlParams.get('max_price') ? parseFloat(urlParams.get('max_price')!) : undefined,
      in_stock: urlParams.get('in_stock') === 'true' ? true : undefined,
      rating: urlParams.get('rating') ? parseInt(urlParams.get('rating')!) : undefined,
      tags: urlParams.get('tags')?.split(',') || undefined,
      sort_by: urlParams.get('sort_by') || 'created_at',
      sort_order: (urlParams.get('sort_order') as 'asc' | 'desc') || 'desc',
    }
  }, [searchParams])

  const { data, isLoading, error } = useProducts(params)

  const products = data?.data || []
  const pagination = data?.pagination

  if (error) {
    return (
      <div className="min-h-screen bg-black">
        <div className="container mx-auto px-4 py-8">
          <Card className="p-8 text-center bg-gray-900 border-gray-700">
            <h2 className="text-2xl font-bold text-white mb-4">
              Oops! Something went wrong
            </h2>
            <p className="text-gray-300 mb-6">
              We couldn't load the products. Please try again later.
            </p>
            <Button onClick={() => window.location.reload()} variant="gradient">
              Try Again
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      {/* Enhanced Background Pattern - Matching Newsletter Section */}
      <AnimatedBackground className="opacity-30" />
      
      {/* Main Content Area */}
      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-8 relative z-10">
        <div className="flex flex-col lg:flex-row gap-8">
          {/* Mobile Filters Overlay */}
          {showFilters && (
            <div className="fixed inset-0 z-50 lg:hidden">
              <div 
                className="absolute inset-0 bg-black/70 backdrop-blur-sm"
                onClick={() => setShowFilters(false)}
              />
              <div className="absolute top-0 left-0 w-80 h-full bg-gray-900/95 backdrop-blur-2xl border-r border-white/20 rounded-r-xl shadow-2xl p-6 overflow-y-auto">
                <div className="flex items-center justify-between mb-6">
                  <h3 className="text-lg font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">Filters</h3>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setShowFilters(false)}
                    className="text-gray-400 hover:text-white hover:bg-white/10 rounded-lg h-8 w-8 p-0 transition-all duration-300 border border-white/10"
                  >
                    <div className="text-sm">✕</div>
                  </Button>
                </div>
                <ProductFilters currentParams={params} />
              </div>
            </div>
          )}
          
          {/* Compact Sidebar Filters */}
          <aside className="hidden lg:block lg:w-60 xl:w-64 flex-shrink-0">
            <div className="sticky top-8">
              <ProductFilters currentParams={params} />
            </div>
          </aside>

          {/* Main Content */}
          <main className="flex-1 min-w-0">
            {/* Compact Page Header */}
            <div className="mb-8">
              {/* Title and Results */}
              <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6 mb-6">
                <div className="space-y-2">
                  <h1 className="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
                    {params.search ? (
                      <>Search Results for <span className="text-[#ff9000]">"{params.search}"</span></>
                    ) : (
                      <>Discover Amazing <span className="text-[#ff9000]">Products</span></>
                    )}
                  </h1>
                  {pagination && (
                    <div className="flex items-center gap-3">
                      <p className="text-gray-400 text-sm">
                        Showing <span className="text-[#ff9000] font-medium">{((pagination.page - 1) * pagination.limit) + 1} - {Math.min(pagination.page * pagination.limit, pagination.total)}</span> of <span className="text-[#ff9000] font-medium">{pagination.total}</span> products
                      </p>
                      <Badge className="bg-white/8 text-gray-300 border-white/15 px-2 py-1 text-xs backdrop-blur-sm font-medium">
                        {pagination.total} Total
                      </Badge>
                    </div>
                  )}
                </div>

                  {/* Compact View Controls */}
                <div className="flex items-center gap-3">
                  {/* View mode toggle */}
                  <div className="flex items-center bg-white/[0.06] backdrop-blur-md border border-white/10 rounded-lg p-0.5 shadow-sm">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setViewMode('grid')}
                      className={`rounded-md h-7 px-2 text-xs font-medium transition-all duration-200 border-0 hover:text-[#ff9000] hover:bg-[#ff9000]/10 ${
                        viewMode === 'grid'
                          ? 'bg-white/10 text-white'
                          : 'text-gray-400'
                      }`}
                    >
                      <Grid className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => setViewMode('list')}
                      className={`rounded-md h-7 px-2 text-xs font-medium transition-all duration-200 border-0 hover:text-[#ff9000] hover:bg-[#ff9000]/10 ${
                        viewMode === 'list'
                          ? 'bg-white/10 text-white'
                          : 'text-gray-400'
                      }`}
                    >
                      <List className="h-3.5 w-3.5" />
                    </Button>
                  </div>

                  {/* Sort Dropdown */}
                  <div className="hidden lg:block">
                    <ProductSort currentSort={`${params.sort_by}:${params.sort_order}`} />
                  </div>

                  {/* Filters toggle (mobile) */}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowFilters(!showFilters)}
                    className="lg:hidden rounded-lg border border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-400 hover:bg-white/[0.08] hover:border-white/15 hover:text-white transition-all duration-200 h-7 px-2.5 text-xs font-medium shadow-sm"
                  >
                    <SlidersHorizontal className="h-3 w-3 mr-1" />
                    Filters
                  </Button>
                </div>
              </div>

              {/* Category Breadcrumbs */}
              {params.category_id && (
                <div className="mb-8">
                  <CategoryBreadcrumbs
                    categoryId={params.category_id}
                    className="bg-white/[0.08] backdrop-blur-md rounded-lg px-4 py-3 border border-white/15 shadow-md"
                    showHome={true}
                  />
                </div>
              )}
            </div>

            {/* Active Filters Bar with Better Spacing */}
            {(params.search || params.category_id || params.min_price || params.max_price) && (
              <div className="flex flex-wrap items-center justify-between gap-3 mb-10 mt-8">
                <div className="flex flex-wrap items-center gap-2">
                  {/* Active filters with compact styling */}
                  {params.search && (
                    <Badge className="bg-white/10 text-gray-300 border border-white/20 px-2.5 py-1 text-xs backdrop-blur-sm font-medium hover:bg-white/15 transition-colors">
                      <span className="text-[#ff9000] mr-1">🔍</span>
                      Search: {params.search}
                    </Badge>
                  )}
                  {params.category_id && (
                    <Badge className="bg-white/10 text-gray-300 border border-white/20 px-2.5 py-1 text-xs backdrop-blur-sm font-medium hover:bg-white/15 transition-colors">
                      <span className="text-[#ff9000] mr-1">📁</span>
                      Category
                    </Badge>
                  )}
                  {params.min_price && (
                    <Badge className="bg-white/10 text-gray-300 border border-white/20 px-2.5 py-1 text-xs backdrop-blur-sm font-medium hover:bg-white/15 transition-colors">
                      <span className="text-[#ff9000] mr-1">💰</span>
                      Min: ${params.min_price}
                    </Badge>
                  )}
                  {params.max_price && (
                    <Badge className="bg-white/10 text-gray-300 border border-white/20 px-2.5 py-1 text-xs backdrop-blur-sm font-medium hover:bg-white/15 transition-colors">
                      <span className="text-[#ff9000] mr-1">💰</span>
                      Max: ${params.max_price}
                    </Badge>
                  )}
                </div>

                {/* Mobile Sort */}
                <div className="lg:hidden">
                  <ProductSort currentSort={`${params.sort_by}:${params.sort_order}`} />
                </div>
              </div>
            )}

            {/* Products Section with Conditional Spacing */}
            <div className={params.search || params.category_id || params.min_price || params.max_price ? 'mt-0' : 'mt-8'}>
            {/* Optimized Products Grid/List */}
            {isLoading ? (
              viewMode === 'grid' ? (
                <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {[...Array(12)].map((_, i) => (
                    <Card key={i} className="animate-pulse bg-white/[0.08] backdrop-blur-sm border border-white/15 rounded-xl overflow-hidden shadow-md">
                      <div className="aspect-square bg-gradient-to-br from-gray-700/50 to-gray-800/50 rounded-t-xl"></div>
                      <CardContent className="p-4">
                        <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded mb-2"></div>
                        <div className="h-3 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-2/3 mb-2"></div>
                        <div className="h-5 bg-gradient-to-r from-[#ff9000]/40 to-[#ff9000]/40 rounded w-1/2"></div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              ) : (
                <div className="space-y-4">
                  {[...Array(8)].map((_, i) => (
                    <Card key={i} className="animate-pulse bg-white/[0.08] backdrop-blur-sm border border-white/15 rounded-xl overflow-hidden shadow-md">
                      <div className="flex gap-4 p-4">
                        <div className="w-24 h-24 bg-gradient-to-br from-gray-700/50 to-gray-800/50 rounded-lg flex-shrink-0"></div>
                        <div className="flex-1 space-y-3">
                          <div className="h-3 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-1/4"></div>
                          <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-3/4"></div>
                          <div className="h-3 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-full"></div>
                          <div className="h-3 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-2/3"></div>
                          <div className="flex justify-between items-center mt-4">
                            <div className="h-5 bg-gradient-to-r from-[#ff9000]/40 to-[#ff9000]/40 rounded w-1/3"></div>
                            <div className="h-8 bg-gradient-to-r from-[#ff9000]/40 to-[#ff9000]/40 rounded w-20"></div>
                          </div>
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              )
            ) : products.length > 0 ? (
              <>
                {viewMode === 'grid' ? (
                  <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                    {products.map((product, index) => (
                      <div
                        key={product.id}
                        className="animate-in fade-in slide-in-from-bottom-4 duration-500"
                        style={{ animationDelay: `${index * 50}ms` }}
                      >
                        <ProductCard
                          product={product}
                          className="transform hover:scale-[1.02] transition-all duration-300 hover:shadow-lg hover:shadow-[#ff9000]/10"
                        />
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="space-y-4">
                    {products.map((product, index) => (
                      <div
                        key={product.id}
                        className="animate-in fade-in slide-in-from-left-4 duration-500"
                        style={{ animationDelay: `${index * 100}ms` }}
                      >
                        <ProductListCard
                          product={product}
                          className="bg-white/[0.08] backdrop-blur-xl border border-white/15 rounded-xl hover:bg-white/[0.12] transition-all duration-300 hover:shadow-lg hover:shadow-[#ff9000]/10"
                        />
                      </div>
                    ))}
                  </div>
                )}

                {/* Compact Pagination */}
                {pagination && pagination.total_pages > 1 && (
                  <div className="mt-8 flex justify-center">
                    <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 p-2 shadow-md">
                      <Pagination
                        currentPage={pagination.page}
                        totalPages={pagination.total_pages}
                        hasNext={pagination.has_next}
                        hasPrev={pagination.has_prev}
                      />
                    </div>
                  </div>
                )}
              </>
            ) : (
              <Card className="p-8 text-center bg-white/[0.08] backdrop-blur-xl border border-white/15 rounded-xl shadow-md">
                <div className="max-w-sm mx-auto">
                  <div className="w-16 h-16 bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/20 rounded-full flex items-center justify-center mx-auto mb-4 border border-[#ff9000]/30">
                    <Filter className="h-8 w-8 text-[#ff9000]" />
                  </div>
                  <h3 className="text-xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-3">
                    No products found
                  </h3>
                  <p className="text-gray-400 mb-6 text-sm leading-relaxed">
                    Try adjusting your search or filter criteria to find what you're looking for.
                  </p>
                  <Button
                    onClick={() => window.location.href = '/products'}
                    className="bg-gradient-to-r from-[#ff9000] to-[#ff9000] hover:from-[#ff9000]/90 hover:to-[#ff9000]/90 text-white px-6 py-2.5 rounded-lg font-medium shadow-md shadow-[#ff9000]/20 transition-all duration-300 hover:scale-105"
                  >
                    Clear All Filters
                  </Button>
                </div>
              </Card>
            )}
            </div>
          </main>
        </div>
      </div>
    </div>
  )
}
