'use client'

import { useState, useMemo } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useRouter, useSearchParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ChevronRight, Home } from 'lucide-react'
import { ArrowLeft, Grid, List, SlidersHorizontal, Package, FolderTree, Layers, SortAsc, SortDesc } from 'lucide-react'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { categoryService } from '@/lib/services/categories'
import { productService } from '@/lib/services/products'
import { ProductCard } from '@/components/products/product-card'
import { CategoryCard } from '@/components/ui/category-card'
import { CategoryFilters } from '@/components/categories/category-filters'
import { Category, Product } from '@/types'
import { cn } from '@/lib/utils'

interface Props {
  slug: string
  searchParams: {
    page?: string
    limit?: string
    sort_by?: string
    sort_order?: string
  }
}

export function CategoryDetailPage({ slug, searchParams }: Props) {
  const router = useRouter()
  const urlSearchParams = useSearchParams()
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [showFilters, setShowFilters] = useState(false)

  // Parse search params
  const page = parseInt(searchParams.page || '1')
  const limit = parseInt(searchParams.limit || '20')
  const sortBy = searchParams.sort_by || 'created_at'
  const sortOrder = searchParams.sort_order || 'desc'

  // Build params object for filters and API calls
  const params = {
    page,
    limit,
    sort_by: sortBy,
    sort_order: sortOrder as 'asc' | 'desc',
    search: searchParams.search,
    min_price: searchParams.min_price ? parseInt(searchParams.min_price) : undefined,
    max_price: searchParams.max_price ? parseInt(searchParams.max_price) : undefined,
    rating: searchParams.rating ? parseInt(searchParams.rating) : undefined,
    featured: searchParams.featured === 'true' ? true : undefined,
    on_sale: searchParams.on_sale === 'true' ? true : undefined,
    stock_status: searchParams.stock_status as 'in_stock' | 'out_of_stock' | 'on_backorder' | undefined,
    requires_shipping: searchParams.requires_shipping === 'false' ? false : undefined,
  }

  // Update URL with new search params
  const updateSearchParams = (newParams: Record<string, string>) => {
    const params = new URLSearchParams(urlSearchParams.toString())

    Object.entries(newParams).forEach(([key, value]) => {
      if (value) {
        params.set(key, value)
      } else {
        params.delete(key)
      }
    })

    router.push(`/categories/${slug}?${params.toString()}`)
  }

  const handleSortChange = (newSortBy: string, newSortOrder: 'asc' | 'desc') => {
    updateSearchParams({
      sort_by: newSortBy,
      sort_order: newSortOrder,
      page: '1' // Reset to first page
    })
  }

  // Fetch category by slug
  const { 
    data: category, 
    isLoading: isLoadingCategory,
    error: categoryError 
  } = useQuery({
    queryKey: ['category', slug],
    queryFn: () => categoryService.getCategoryBySlug(slug),
  })

  // Fetch category products
  const {
    data: productsData,
    isLoading: isLoadingProducts
  } = useQuery({
    queryKey: ['category-products', category?.id, params],
    queryFn: () => {
      if (!category?.id) return null
      return productService.getProductsByCategory(category.id, params)
    },
    enabled: !!category?.id,
  })

  // Fetch subcategories
  const {
    data: subcategories
  } = useQuery({
    queryKey: ['subcategories', category?.id],
    queryFn: () => {
      if (!category?.id) return []
      return categoryService.getSubcategories(category.id)
    },
    enabled: !!category?.id,
  })

  if (categoryError) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white">
        <div className="container mx-auto px-4 py-8">
          <Card className="bg-gray-900/50 border-gray-700 max-w-md mx-auto">
            <CardContent className="p-8 text-center">
              <h2 className="text-2xl font-bold text-white mb-4">Category Not Found</h2>
              <p className="text-gray-400 mb-6">
                The category you're looking for doesn't exist or has been removed.
              </p>
              <Button asChild variant="outline" className="bg-gray-800 border-gray-600 text-white hover:bg-gray-700">
                <Link href="/categories">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Back to Categories
                </Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      {/* Enhanced Background Pattern - Matching Products Page */}
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
              <div className="absolute top-0 left-0 w-80 h-full bg-gray-900/95 backdrop-blur-2xl border-r border-white/20 rounded-r-xl shadow-2xl p-6 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-gray-800">
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
                <CategoryFilters currentParams={params} categoryId={category?.id || ''} categorySlug={slug} />
              </div>
            </div>
          )}

          {/* Compact Sidebar Filters */}
          <aside className="hidden lg:block lg:w-60 xl:w-64 flex-shrink-0">
            <div className="sticky top-4">
              <CategoryFilters currentParams={params} categoryId={category?.id || ''} categorySlug={slug} />
            </div>
          </aside>

          {/* Main Content */}
          <main className="flex-1 min-w-0">
            {/* Breadcrumb Navigation */}
            <nav className="flex items-center space-x-2 text-sm mb-6">
              <Link href="/" className="text-gray-400 hover:text-white transition-colors flex items-center">
                <Home className="h-4 w-4" />
                <span className="ml-1">Home</span>
              </Link>
              <ChevronRight className="h-4 w-4 text-gray-500" />
              <Link href="/categories" className="text-gray-400 hover:text-white transition-colors">
                Categories
              </Link>
              <ChevronRight className="h-4 w-4 text-gray-500" />
              <span className="text-[#ff9000] font-medium">
                {category?.name || slug}
              </span>
            </nav>

            {/* Page Header */}
            <div className="mb-8">
              {/* Title and Results */}
              <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6 mb-6">
                <div className="space-y-2">
                  {isLoadingCategory ? (
                    <div className="h-8 w-64 bg-gray-800 animate-pulse rounded" />
                  ) : category ? (
                    <h1 className="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
                      {category.name}
                    </h1>
                  ) : (
                    <h1 className="text-2xl lg:text-3xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent leading-tight">
                      Category
                    </h1>
                  )}

                  {category?.description && (
                    <p className="text-gray-400 text-sm lg:text-base leading-relaxed max-w-2xl">
                      {category.description}
                    </p>
                  )}

                  {productsData?.data && (
                    <div className="flex items-center gap-3 pt-2">
                      <div className="text-gray-400 text-sm">
                        {productsData.data.length} {productsData.data.length === 1 ? 'product' : 'products'} found
                      </div>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Subcategories */}
            {subcategories && subcategories.length > 0 && (
              <div className="mb-8">
                <h2 className="text-lg font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-4 flex items-center">
                  <Layers className="w-4 h-4 mr-2 text-[#ff9000]" />
                  Subcategories
                </h2>
                <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
                  {subcategories.map((subcat) => (
                    <CategoryCard
                      key={subcat.id}
                      category={subcat}
                      variant="compact"
                      showStats={false}
                    />
                  ))}
                </div>
              </div>
            )}

            {/* Products Section */}
            <div>
              {/* Products Header with Controls */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
                <div className="flex items-center gap-3">
                  <h2 className="text-lg font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent flex items-center">
                    <Package className="w-4 h-4 mr-2 text-[#ff9000]" />
                    Products
                  </h2>
                </div>

                {/* Controls */}
                <div className="flex items-center gap-3">
                  {/* Mobile Filter Button */}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowFilters(true)}
                    className="lg:hidden bg-white/[0.08] backdrop-blur-sm border border-white/15 text-gray-300 hover:bg-white/[0.12] hover:text-white transition-all duration-300 text-xs"
                  >
                    <SlidersHorizontal className="h-3 w-3 mr-1" />
                    Filters
                  </Button>

                  {/* Sort Button */}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleSortChange(sortBy, sortOrder === 'asc' ? 'desc' : 'asc')}
                    className="bg-white/[0.08] backdrop-blur-sm border border-white/15 text-gray-300 hover:bg-white/[0.12] hover:text-white transition-all duration-300 text-xs"
                  >
                    {sortOrder === 'asc' ? <SortAsc className="h-3 w-3 mr-1" /> : <SortDesc className="h-3 w-3 mr-1" />}
                    {sortBy === 'created_at' ? 'Newest' : sortBy === 'name' ? 'Name' : 'Price'}
                  </Button>

                  {/* View Toggle */}
                  <div className="flex items-center gap-1 bg-white/[0.08] backdrop-blur-sm border border-white/15 rounded-lg p-1">
                    <Button
                      variant={viewMode === 'grid' ? 'default' : 'ghost'}
                      size="sm"
                      onClick={() => setViewMode('grid')}
                      className={cn(
                        'h-7 w-7 p-0 rounded-md transition-all duration-300',
                        viewMode === 'grid'
                          ? 'bg-[#ff9000] text-white shadow-sm'
                          : 'text-gray-400 hover:text-white hover:bg-white/10'
                      )}
                    >
                      <Grid className="h-3 w-3" />
                    </Button>
                    <Button
                      variant={viewMode === 'list' ? 'default' : 'ghost'}
                      size="sm"
                      onClick={() => setViewMode('list')}
                      className={cn(
                        'h-7 w-7 p-0 rounded-md transition-all duration-300',
                        viewMode === 'list'
                          ? 'bg-[#ff9000] text-white shadow-sm'
                          : 'text-gray-400 hover:text-white hover:bg-white/10'
                      )}
                    >
                      <List className="h-3 w-3" />
                    </Button>
                  </div>
                </div>
              </div>

              {/* Products Grid/List */}
              {isLoadingProducts ? (
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
              ) : productsData?.data && productsData.data.length > 0 ? (
                <>
                  {viewMode === 'grid' ? (
                    <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                      {productsData.data.map((product) => (
                        <ProductCard
                          key={product.id}
                          product={product}
                          showQuickView={true}
                          showWishlist={true}
                        />
                      ))}
                    </div>
                  ) : (
                    <div className="space-y-4">
                      {productsData.data.map((product) => (
                        <ProductCard
                          key={product.id}
                          product={product}
                          variant="list"
                          showQuickView={true}
                          showWishlist={true}
                        />
                      ))}
                    </div>
                  )}
                </>
              ) : (
                <Card className="p-8 text-center bg-white/[0.08] backdrop-blur-xl border border-white/15 rounded-xl shadow-md">
                  <div className="max-w-sm mx-auto">
                    <div className="w-16 h-16 bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/20 rounded-full flex items-center justify-center mx-auto mb-4 border border-[#ff9000]/30">
                      <Package className="h-8 w-8 text-[#ff9000]" />
                    </div>
                    <h3 className="text-xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-3">
                      No products found
                    </h3>
                    <p className="text-gray-400 mb-6 text-sm leading-relaxed">
                      This category doesn't have any products yet.
                    </p>
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
