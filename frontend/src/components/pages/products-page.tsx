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
    <div className="min-h-screen bg-black">
      {/* Simplified Header */}
      <div className="bg-gradient-to-r from-gray-900 to-black border-b border-gray-800">
        <div className="container mx-auto px-4 py-6">
          {/* Title and Results Count */}
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
            <div>
              <h1 className="text-2xl lg:text-3xl font-bold text-white mb-2">
                {params.search ? (
                  <>Search Results for <span className="text-orange-400">"{params.search}"</span></>
                ) : (
                  <>Products</>
                )}
              </h1>
              {pagination && (
                <p className="text-sm text-gray-400">
                  Showing {((pagination.page - 1) * pagination.limit) + 1} - {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} products
                </p>
              )}
            </div>

            {/* View Controls */}
            <div className="flex items-center gap-4">
              {/* View mode toggle */}
              <div className="flex items-center bg-gray-800 border border-gray-700 rounded-lg p-1">
                <Button
                  variant={viewMode === 'grid' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('grid')}
                  className="rounded-md h-8 px-3 text-sm"
                  style={viewMode === 'grid' ? {backgroundColor: '#FF9000', color: 'white'} : {color: 'white'}}
                >
                  <Grid className="h-4 w-4 mr-1" />
                  Grid
                </Button>
                <Button
                  variant={viewMode === 'list' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('list')}
                  className="rounded-md h-8 px-3 text-sm"
                  style={viewMode === 'list' ? {backgroundColor: '#FF9000', color: 'white'} : {color: 'white'}}
                >
                  <List className="h-4 w-4 mr-1" />
                  List
                </Button>
              </div>

              {/* Filters toggle (mobile) */}
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowFilters(!showFilters)}
                className="lg:hidden rounded-lg border border-gray-600 text-white hover:border-orange-500 transition-all duration-200"
              >
                <SlidersHorizontal className="h-4 w-4 mr-2" />
                Filters
              </Button>
            </div>
          </div>

          {/* Category Breadcrumbs */}
          {params.category_id && (
            <div className="mb-4">
              <CategoryBreadcrumbs
                categoryId={params.category_id}
                className="bg-gray-800/80 backdrop-blur-sm rounded-lg px-3 py-2 border border-gray-700 shadow-sm"
                showHome={true}
              />
            </div>
          )}
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        <div className="flex flex-col lg:flex-row gap-8">
          {/* Enhanced Sidebar Filters */}
          <aside className={`lg:w-80 flex-shrink-0 ${showFilters ? 'block' : 'hidden lg:block'}`}>
            <div className="sticky top-8">
              <ProductFilters currentParams={params} />
            </div>
          </aside>

          {/* Main Content */}
          <main className="flex-1">
            {/* Sort and Results */}
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                {/* Active filters */}
                <div className="flex items-center space-x-2">
                  {params.search && (
                    <Badge className="bg-gray-700 text-gray-200 text-xs">
                      Search: {params.search}
                    </Badge>
                  )}
                  {params.category_id && (
                    <Badge className="bg-gray-700 text-gray-200 text-xs">
                      Category
                    </Badge>
                  )}
                  {params.min_price && (
                    <Badge className="bg-gray-700 text-gray-200 text-xs">
                      Min: ${params.min_price}
                    </Badge>
                  )}
                  {params.max_price && (
                    <Badge className="bg-gray-700 text-gray-200 text-xs">
                      Max: ${params.max_price}
                    </Badge>
                  )}
                </div>
              </div>

              <ProductSort currentSort={`${params.sort_by}:${params.sort_order}`} />
            </div>

            {/* Products Grid/List */}
            {isLoading ? (
              viewMode === 'grid' ? (
                <div className="grid gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {[...Array(12)].map((_, i) => (
                    <Card key={i} className="animate-pulse bg-gray-800 border-gray-700">
                      <div className="aspect-square bg-gray-700 rounded-t-lg"></div>
                      <CardContent className="p-4">
                        <div className="h-4 bg-gray-600 rounded mb-2"></div>
                        <div className="h-3 bg-gray-600 rounded w-2/3 mb-2"></div>
                        <div className="h-5 bg-gray-600 rounded w-1/2"></div>
                      </CardContent>
                    </Card>
                  ))}
                </div>
              ) : (
                <div className="space-y-4">
                  {[...Array(8)].map((_, i) => (
                    <Card key={i} className="animate-pulse bg-gray-800 border-gray-700">
                      <div className="flex gap-4 p-4">
                        <div className="w-32 h-32 bg-gray-700 rounded-lg flex-shrink-0"></div>
                        <div className="flex-1 space-y-3">
                          <div className="h-3 bg-gray-600 rounded w-1/4"></div>
                          <div className="h-5 bg-gray-600 rounded w-3/4"></div>
                          <div className="h-3 bg-gray-600 rounded w-full"></div>
                          <div className="h-3 bg-gray-600 rounded w-2/3"></div>
                          <div className="flex justify-between items-center mt-4">
                            <div className="h-6 bg-gray-600 rounded w-1/3"></div>
                            <div className="h-8 bg-gray-600 rounded w-24"></div>
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
                    {products.map((product) => (
                      <ProductCard
                        key={product.id}
                        product={product}
                      />
                    ))}
                  </div>
                ) : (
                  <div className="space-y-4">
                    {products.map((product) => (
                      <ProductListCard
                        key={product.id}
                        product={product}
                      />
                    ))}
                  </div>
                )}

                {/* Pagination */}
                {pagination && pagination.total_pages > 1 && (
                  <div className="mt-8 flex justify-center">
                    <Pagination
                      currentPage={pagination.page}
                      totalPages={pagination.total_pages}
                      hasNext={pagination.has_next}
                      hasPrev={pagination.has_prev}
                    />
                  </div>
                )}
              </>
            ) : (
              <Card className="p-8 text-center bg-gray-900 border-gray-700">
                <div className="max-w-md mx-auto">
                  <div className="w-12 h-12 bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-3">
                    <Filter className="h-6 w-6 text-gray-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-white mb-2">
                    No products found
                  </h3>
                  <p className="text-gray-300 mb-4">
                    Try adjusting your search or filter criteria to find what you're looking for.
                  </p>
                  <Button onClick={() => window.location.href = '/products'} variant="gradient">
                    Clear Filters
                  </Button>
                </div>
              </Card>
            )}
          </main>
        </div>
      </div>
    </div>
  )
}
