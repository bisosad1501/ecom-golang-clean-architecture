'use client'

import { useState, useMemo } from 'react'
import { useSearchParams } from 'next/navigation'
import { Filter, Grid, List, SlidersHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProductCard } from '@/components/products/product-card'
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
      <div className="container mx-auto px-4 py-8">
        <Card className="p-8 text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">
            Oops! Something went wrong
          </h2>
          <p className="text-gray-600 mb-6">
            We couldn't load the products. Please try again later.
          </p>
          <Button onClick={() => window.location.reload()}>
            Try Again
          </Button>
        </Card>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background">
      {/* Enhanced Header */}
      <div className="bg-gradient-to-r from-primary-50/50 to-violet-50/50 border-b border-border/30">
        <div className="container mx-auto px-4 py-12">
          <div className="text-center mb-8">
            <div className="flex items-center justify-center gap-3 mb-6">
              <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large">
                <Grid className="h-6 w-6 text-white" />
              </div>
              <span className="text-primary font-semibold">PRODUCT COLLECTION</span>
            </div>

            <h1 className="text-4xl lg:text-6xl font-bold text-foreground mb-4">
              {params.search ? (
                <>Search Results for <span className="text-gradient">"{params.search}"</span></>
              ) : (
                <>Discover Our <span className="text-gradient">Premium Products</span></>
              )}
            </h1>

            {pagination && (
              <p className="text-lg text-muted-foreground mb-6">
                Showing {((pagination.page - 1) * pagination.limit) + 1} - {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} products
              </p>
            )}

            {/* Quick stats */}
            <div className="flex items-center justify-center gap-8">
              <div className="text-center">
                <div className="text-2xl font-bold text-primary">{pagination?.total || 0}</div>
                <div className="text-sm text-muted-foreground">Products</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-emerald-600">4.9★</div>
                <div className="text-sm text-muted-foreground">Avg Rating</div>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">24h</div>
                <div className="text-sm text-muted-foreground">Fast Delivery</div>
              </div>
            </div>
          </div>

          {/* Category Breadcrumbs */}
          {params.category_id && (
            <div className="mt-6 flex justify-center">
              <CategoryBreadcrumbs 
                categoryId={params.category_id}
                className="bg-white/80 backdrop-blur-sm rounded-xl px-4 py-2 border border-border/50 shadow-sm"
                showHome={true}
              />
            </div>
          )}

          <div className="flex items-center justify-center gap-6">
            {/* View mode toggle */}
            <div className="flex items-center bg-background border-2 border-border rounded-2xl p-1 shadow-medium">
              <Button
                variant={viewMode === 'grid' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setViewMode('grid')}
                className="rounded-xl h-10 px-4"
              >
                <Grid className="h-4 w-4 mr-2" />
                Grid
              </Button>
              <Button
                variant={viewMode === 'list' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setViewMode('list')}
                className="rounded-xl h-10 px-4"
              >
                <List className="h-4 w-4 mr-2" />
                List
              </Button>
            </div>

            {/* Filters toggle (mobile) */}
            <Button
              variant="outline"
              size="lg"
              onClick={() => setShowFilters(!showFilters)}
              className="lg:hidden rounded-2xl border-2 hover:border-primary transition-all duration-200"
            >
              <SlidersHorizontal className="h-5 w-5 mr-2" />
              Filters
            </Button>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-12">
        <div className="flex flex-col lg:flex-row gap-12">
          {/* Enhanced Sidebar Filters */}
          <aside className={`lg:w-96 flex-shrink-0 ${showFilters ? 'block' : 'hidden lg:block'}`}>
            <div className="sticky top-8">
              <ProductFilters currentParams={params} />
            </div>
          </aside>

          {/* Main Content */}
          <main className="flex-1">
            {/* Sort and Results */}
            <div className="flex items-center justify-between mb-6">
              <div className="flex items-center space-x-4">
                {/* Active filters */}
                <div className="flex items-center space-x-2">
                  {params.search && (
                    <Badge variant="secondary">
                      Search: {params.search}
                    </Badge>
                  )}
                  {params.category_id && (
                    <Badge variant="secondary">
                      Category
                    </Badge>
                  )}
                  {params.min_price && (
                    <Badge variant="secondary">
                      Min: ${params.min_price}
                    </Badge>
                  )}
                  {params.max_price && (
                    <Badge variant="secondary">
                      Max: ${params.max_price}
                    </Badge>
                  )}
                </div>
              </div>

              <ProductSort currentSort={`${params.sort_by}:${params.sort_order}`} />
            </div>

            {/* Products Grid/List */}
            {isLoading ? (
              <div className={`grid gap-6 ${
                viewMode === 'grid' 
                  ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4' 
                  : 'grid-cols-1'
              }`}>
                {[...Array(12)].map((_, i) => (
                  <Card key={i} className="animate-pulse">
                    <div className="aspect-square bg-gray-200 rounded-t-lg"></div>
                    <CardContent className="p-4">
                      <div className="h-4 bg-gray-200 rounded mb-2"></div>
                      <div className="h-4 bg-gray-200 rounded w-2/3 mb-2"></div>
                      <div className="h-6 bg-gray-200 rounded w-1/2"></div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : products.length > 0 ? (
              <>
                <div className={`grid gap-6 ${
                  viewMode === 'grid' 
                    ? 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4' 
                    : 'grid-cols-1'
                }`}>
                  {products.map((product) => (
                    <ProductCard 
                      key={product.id} 
                      product={product}
                      className={viewMode === 'list' ? 'flex-row' : ''}
                    />
                  ))}
                </div>

                {/* Pagination */}
                {pagination && pagination.total_pages > 1 && (
                  <div className="mt-12 flex justify-center">
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
              <Card className="p-12 text-center">
                <div className="max-w-md mx-auto">
                  <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                    <Filter className="h-8 w-8 text-gray-400" />
                  </div>
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    No products found
                  </h3>
                  <p className="text-gray-600 mb-6">
                    Try adjusting your search or filter criteria to find what you're looking for.
                  </p>
                  <Button onClick={() => window.location.href = '/products'}>
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
