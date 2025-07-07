'use client'

import { useState, useEffect } from 'react'
import { useSearchParams, useRouter } from 'next/navigation'
import { Search, Filter, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProductCard } from '@/components/products/product-card'
import { ProductFilters } from '@/components/products/product-filters'
import { ProductSort } from '@/components/products/product-sort'
import { Pagination } from '@/components/ui/pagination'
import { useSearchProducts, useProductSuggestions } from '@/hooks/use-products'
import { ProductsParams } from '@/lib/services/products'
import { DEFAULT_PAGE_SIZE } from '@/constants'
import { cn } from '@/lib/utils'
import { PageWrapper, PageHeader, PageSection, PageGrid } from '@/components/layout'

export function SearchPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [searchQuery, setSearchQuery] = useState('')
  const [showFilters, setShowFilters] = useState(false)
  const [showSuggestions, setShowSuggestions] = useState(false)

  const query = searchParams.get('q') || ''
  
  // Parse URL parameters
  const params: ProductsParams = {
    page: parseInt(searchParams.get('page') || '1'),
    limit: parseInt(searchParams.get('limit') || DEFAULT_PAGE_SIZE.toString()),
    category_id: searchParams.get('category') || undefined,
    min_price: searchParams.get('min_price') ? parseFloat(searchParams.get('min_price')!) : undefined,
    max_price: searchParams.get('max_price') ? parseFloat(searchParams.get('max_price')!) : undefined,
    in_stock: searchParams.get('in_stock') === 'true' ? true : undefined,
    rating: searchParams.get('rating') ? parseInt(searchParams.get('rating')!) : undefined,
    tags: searchParams.get('tags')?.split(',') || undefined,
    sort_by: searchParams.get('sort_by') || 'relevance',
    sort_order: (searchParams.get('sort_order') as 'asc' | 'desc') || 'desc',
  }

  const { data, isLoading, error } = useSearchProducts(query, params)
  const { data: suggestions } = useProductSuggestions(searchQuery, 5)

  const products = data?.data || []
  const pagination = data?.pagination

  useEffect(() => {
    setSearchQuery(query)
  }, [query])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      const newParams = new URLSearchParams()
      newParams.set('q', searchQuery.trim())
      router.push(`/search?${newParams.toString()}`)
      setShowSuggestions(false)
    }
  }

  const handleSuggestionClick = (suggestion: string) => {
    setSearchQuery(suggestion)
    const newParams = new URLSearchParams()
    newParams.set('q', suggestion)
    router.push(`/search?${newParams.toString()}`)
    setShowSuggestions(false)
  }

  const clearSearch = () => {
    setSearchQuery('')
    router.push('/products')
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Search Header */}
      <div className="bg-white border-b">
        <div className="container mx-auto px-4 py-6">
          <div className="max-w-2xl mx-auto">
            <form onSubmit={handleSearch} className="relative">
              <div className="relative">
                <Input
                  type="search"
                  placeholder="Search for products..."
                  value={searchQuery}
                  onChange={(e) => {
                    setSearchQuery(e.target.value)
                    setShowSuggestions(e.target.value.length > 1)
                  }}
                  onFocus={() => setShowSuggestions(searchQuery.length > 1)}
                  onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
                  leftIcon={<Search className="h-4 w-4" />}
                  rightIcon={
                    searchQuery && (
                      <button
                        type="button"
                        onClick={() => setSearchQuery('')}
                        className="text-gray-400 hover:text-gray-600"
                      >
                        <X className="h-4 w-4" />
                      </button>
                    )
                  }
                  className="text-lg h-12"
                />
                
                {/* Search Suggestions */}
                {showSuggestions && suggestions && suggestions.length > 0 && (
                  <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-md shadow-lg z-50">
                    {suggestions.map((product) => (
                      <button
                        key={product.id}
                        onClick={() => handleSuggestionClick(product.name)}
                        className="w-full px-4 py-3 text-left hover:bg-gray-50 flex items-center space-x-3"
                      >
                        <Search className="h-4 w-4 text-gray-400" />
                        <span>{product.name}</span>
                      </button>
                    ))}
                  </div>
                )}
              </div>
              
              <Button type="submit" className="mt-4 w-full sm:w-auto">
                Search
              </Button>
            </form>
          </div>
        </div>
      </div>

      {/* Results */}
      <div className="container mx-auto px-4 py-8">
        {query && (
          <div className="mb-6">
            <div className="flex items-center justify-between">
              <div>
                <h1 className="text-2xl font-bold text-gray-900">
                  Search Results for "{query}"
                </h1>
                {pagination && (
                  <p className="text-sm text-gray-500 mt-1">
                    {pagination.total} products found
                  </p>
                )}
              </div>
              
              <Button
                variant="outline"
                onClick={clearSearch}
                className="text-sm"
              >
                Clear Search
              </Button>
            </div>
          </div>
        )}

        {query ? (
          <div className="flex gap-8">
            {/* Sidebar Filters */}
            <aside className={`w-80 flex-shrink-0 ${showFilters ? 'block' : 'hidden lg:block'}`}>
              <div className="sticky top-8">
                <ProductFilters currentParams={params} />
              </div>
            </aside>

            {/* Main Content */}
            <main className="flex-1">
              {/* Sort and Filters Toggle */}
              <div className="flex items-center justify-between mb-6">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowFilters(!showFilters)}
                  className="lg:hidden"
                >
                  <Filter className="h-4 w-4 mr-2" />
                  Filters
                </Button>

                <ProductSort currentSort={`${params.sort_by}:${params.sort_order}`} />
              </div>

              {/* Products Grid */}
              {isLoading ? (
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
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
              ) : error ? (
                <Card className="p-8 text-center">
                  <h3 className="text-xl font-semibold text-gray-900 mb-2">
                    Search Error
                  </h3>
                  <p className="text-gray-600 mb-4">
                    Something went wrong while searching. Please try again.
                  </p>
                  <Button onClick={() => window.location.reload()}>
                    Try Again
                  </Button>
                </Card>
              ) : products.length > 0 ? (
                <>
                  <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                    {products.map((product) => (
                      <ProductCard key={product.id} product={product} />
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
                    <Search className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      No products found
                    </h3>
                    <p className="text-gray-600 mb-6">
                      We couldn't find any products matching "{query}". Try adjusting your search terms or browse our categories.
                    </p>
                    <div className="flex flex-col sm:flex-row gap-4 justify-center">
                      <Button onClick={clearSearch}>
                        Browse All Products
                      </Button>
                      <Button variant="outline" onClick={() => setSearchQuery('')}>
                        Try Different Search
                      </Button>
                    </div>
                  </div>
                </Card>
              )}
            </main>
          </div>
        ) : (
          /* No Search Query - Show Popular Searches or Categories */
          <div className="text-center py-12">
            <Search className="h-16 w-16 text-gray-300 mx-auto mb-4" />
            <h2 className="text-2xl font-bold text-gray-900 mb-4">
              What are you looking for?
            </h2>
            <p className="text-gray-600 mb-8">
              Search for products, brands, or categories to find exactly what you need.
            </p>
            
            {/* Popular Searches */}
            <div className="max-w-2xl mx-auto">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                Popular Searches
              </h3>
              <div className="flex flex-wrap gap-2 justify-center">
                {[
                  'iPhone',
                  'Laptop',
                  'Headphones',
                  'Sneakers',
                  'Watch',
                  'Camera',
                  'Gaming',
                  'Books',
                ].map((term) => (
                  <Badge
                    key={term}
                    variant="outline"
                    className="cursor-pointer hover:bg-primary-50"
                    onClick={() => handleSuggestionClick(term)}
                  >
                    {term}
                  </Badge>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
