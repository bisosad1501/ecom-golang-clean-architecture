'use client'

import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import { 
  Grid, 
  List, 
  SlidersHorizontal, 
  ArrowUpDown,
  Filter,
  X,
  ChevronDown,
  Package,
  Star,
  TrendingUp
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'
import { CategoryBreadcrumbs } from '@/components/layout/category-breadcrumbs'
import { EnhancedCategorySidebar } from '@/components/layout/enhanced-category-sidebar'
import { ProductCard } from '@/components/products/product-card'
import { ProductFilters } from '@/components/products/product-filters'
import { useProducts } from '@/hooks/use-products'
import { useCategories } from '@/hooks/use-categories'
import { Category, Product } from '@/types'
import { cn } from '@/lib/utils'
import { apiClient } from '@/lib/api'

interface CategoryPageLayoutProps {
  categoryId?: string
  className?: string
}

interface CategoryInfo {
  category: Category
  productCount: number
  subcategories: Category[]
}

export function CategoryPageLayout({ categoryId, className }: CategoryPageLayoutProps) {
  const searchParams = useSearchParams()
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [sortBy, setSortBy] = useState('name')
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [activeFilters, setActiveFilters] = useState<Record<string, any>>({})
  const [categoryInfo, setCategoryInfo] = useState<CategoryInfo | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const { data: categories } = useCategories()
  
  // Get products with current filters
  const {
    data: productsResponse,
    isLoading: productsLoading,
    error: productsError
  } = useProducts({
    category_id: categoryId,
    sort_by: sortBy,
    sort_order: sortOrder,
    ...activeFilters
  })

  const products = productsResponse?.data || []

  // Fetch category information
  useEffect(() => {
    if (!categoryId || !categories) return

    const fetchCategoryInfo = async () => {
      try {
        setIsLoading(true)
        
        const category = categories.find(c => c.id === categoryId)
        if (!category) return

        // Get product count
        const countResponse = await apiClient.get(`/categories/${categoryId}/count`)
        const productCount = countResponse.data?.product_count || 0
        
        // Get subcategories
        const subcategories = categories.filter(c => c.parent_id === categoryId)
        
        setCategoryInfo({
          category,
          productCount,
          subcategories
        })
      } catch (error) {
        console.error('Error fetching category info:', error)
      } finally {
        setIsLoading(false)
      }
    }

    fetchCategoryInfo()
  }, [categoryId, categories])

  const handleCategoryChange = (newCategoryId: string | undefined) => {
    const params = new URLSearchParams(searchParams.toString())
    if (newCategoryId) {
      params.set('category', newCategoryId)
    } else {
      params.delete('category')
    }
    window.location.href = `?${params.toString()}`
  }

  const handleFilterChange = (filters: Record<string, any>) => {
    setActiveFilters(filters)
  }

  const sortOptions = [
    { value: 'name', label: 'Name' },
    { value: 'price', label: 'Price' },
    { value: 'created_at', label: 'Newest' },
    { value: 'rating', label: 'Rating' },
    { value: 'popularity', label: 'Popularity' }
  ]

  const activeFilterCount = Object.keys(activeFilters).filter(key => 
    activeFilters[key] && activeFilters[key] !== '' && activeFilters[key] !== 0
  ).length

  return (
    <div className={cn("min-h-screen bg-gray-50", className)}>
      {/* Header */}
      <div className="bg-background border-b">
        <div className="container mx-auto px-4 py-6">
          {/* Breadcrumbs */}
          <CategoryBreadcrumbs 
            categoryId={categoryId}
            showProductCount={true}
            className="mb-4"
          />
          
          {/* Category Header */}
          {categoryInfo && (
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <h1 className="text-3xl font-bold text-gray-900 mb-2">
                  {categoryInfo.category.name}
                </h1>
                
                {categoryInfo.category.description && (
                  <p className="text-gray-600 text-lg mb-4 max-w-2xl">
                    {categoryInfo.category.description}
                  </p>
                )}
                
                <div className="flex items-center space-x-4 text-sm text-gray-500">
                  <div className="flex items-center">
                    <Package className="h-4 w-4 mr-1" />
                    <span>{categoryInfo.productCount} products</span>
                  </div>
                  
                  {categoryInfo.subcategories.length > 0 && (
                    <div className="flex items-center">
                      <TrendingUp className="h-4 w-4 mr-1" />
                      <span>{categoryInfo.subcategories.length} subcategories</span>
                    </div>
                  )}
                  
                  <div className="flex items-center">
                    <Star className="h-4 w-4 mr-1" />
                    <span>Featured category</span>
                  </div>
                </div>
              </div>
              
              {categoryInfo.category.image && (
                <div className="ml-6">
                  <img
                    src={categoryInfo.category.image}
                    alt={categoryInfo.category.name}
                    className="w-24 h-24 object-cover rounded-lg border"
                  />
                </div>
              )}
            </div>
          )}
          
          {/* Subcategories */}
          {categoryInfo?.subcategories && categoryInfo.subcategories.length > 0 && (
            <div className="mt-6">
              <h3 className="text-lg font-semibold mb-3">Browse Subcategories</h3>
              <div className="flex flex-wrap gap-2">
                {categoryInfo.subcategories.map((subcat) => (
                  <Button
                    key={subcat.id}
                    variant="outline"
                    size="sm"
                    onClick={() => handleCategoryChange(subcat.id)}
                    className="hover:bg-primary-50 hover:border-primary-200"
                  >
                    {subcat.name}
                    {subcat.product_count !== undefined && (
                      <Badge variant="secondary" className="ml-2">
                        {subcat.product_count}
                      </Badge>
                    )}
                  </Button>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-6">
        <div className="flex gap-6">
          {/* Sidebar - Desktop */}
          <div className="hidden lg:block w-80 flex-shrink-0">
            <div className="sticky top-6 space-y-6">
              <EnhancedCategorySidebar
                selectedCategoryId={categoryId}
                onCategoryChange={handleCategoryChange}
                showProductCount={true}
                showSearch={true}
                showFilters={true}
              />
              
              <ProductFilters currentParams={{}} />
            </div>
          </div>

          {/* Main Content Area */}
          <div className="flex-1 min-w-0">
            {/* Toolbar */}
            <div className="bg-background rounded-lg border p-4 mb-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  {/* Mobile Filter Toggle */}
                  <Sheet>
                    <SheetTrigger asChild>
                      <Button variant="outline" size="sm" className="lg:hidden">
                        <Filter className="h-4 w-4 mr-2" />
                        Filters
                        {activeFilterCount > 0 && (
                          <Badge variant="default" className="ml-2">
                            {activeFilterCount}
                          </Badge>
                        )}
                      </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-80">
                      <SheetHeader>
                        <SheetTitle>Filters & Categories</SheetTitle>
                      </SheetHeader>
                      <div className="mt-6 space-y-6">
                        <EnhancedCategorySidebar
                          selectedCategoryId={categoryId}
                          onCategoryChange={handleCategoryChange}
                          showProductCount={true}
                          showSearch={true}
                          showFilters={true}
                        />
                        
                        <ProductFilters currentParams={{}} />
                      </div>
                    </SheetContent>
                  </Sheet>

                  {/* Results Count */}
                  <div className="text-sm text-gray-600">
                    {productsLoading ? (
                      <div className="h-5 w-32 bg-gray-200 rounded animate-pulse" />
                    ) : (
                      <span>
                        {products?.length || 0} products found
                        {categoryInfo && ` in ${categoryInfo.category.name}`}
                      </span>
                    )}
                  </div>

                  {/* Active Filters */}
                  {activeFilterCount > 0 && (
                    <div className="hidden md:flex items-center space-x-2">
                      {Object.entries(activeFilters).map(([key, value]) => {
                        if (!value || value === '' || value === 0) return null
                        return (
                          <Badge key={key} variant="secondary" className="flex items-center gap-1">
                            {key}: {String(value)}
                            <X
                              className="h-3 w-3 cursor-pointer hover:text-red-500"
                              onClick={() => setActiveFilters(prev => ({ ...prev, [key]: undefined }))}
                            />
                          </Badge>
                        )
                      })}
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => setActiveFilters({})}
                        className="text-xs"
                      >
                        Clear all
                      </Button>
                    </div>
                  )}
                </div>

                <div className="flex items-center space-x-2">
                  {/* Sort */}
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" size="sm">
                        <ArrowUpDown className="h-4 w-4 mr-2" />
                        Sort
                        <ChevronDown className="h-4 w-4 ml-2" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      {sortOptions.map((option) => (
                        <DropdownMenuItem
                          key={option.value}
                          onClick={() => setSortBy(option.value)}
                          className={cn(sortBy === option.value && "bg-primary-50")}
                        >
                          {option.label}
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>

                  {/* View Mode */}
                  <div className="flex border rounded-md overflow-hidden">
                    <Button
                      variant={viewMode === 'grid' ? 'default' : 'ghost'}
                      size="sm"
                      onClick={() => setViewMode('grid')}
                      className="rounded-none"
                    >
                      <Grid className="h-4 w-4" />
                    </Button>
                    <Button
                      variant={viewMode === 'list' ? 'default' : 'ghost'}
                      size="sm"
                      onClick={() => setViewMode('list')}
                      className="rounded-none"
                    >
                      <List className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>

            {/* Products Grid/List */}
            {productsLoading ? (
              <div className={cn(
                "grid gap-6",
                viewMode === 'grid' 
                  ? "grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
                  : "grid-cols-1"
              )}>
                {[...Array(12)].map((_, i) => (
                  <div key={i} className="bg-background rounded-lg border p-4">
                    <div className="aspect-square bg-gray-200 rounded-lg mb-4 animate-pulse" />
                    <div className="space-y-2">
                      <div className="h-4 bg-gray-200 rounded animate-pulse" />
                      <div className="h-4 bg-gray-200 rounded w-2/3 animate-pulse" />
                    </div>
                  </div>
                ))}
              </div>
            ) : productsError ? (
              <Card>
                <CardContent className="text-center py-12">
                  <Package className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                  <h3 className="text-lg font-semibold mb-2">Error loading products</h3>
                  <p className="text-gray-600">There was an error loading the products. Please try again.</p>
                </CardContent>
              </Card>
            ) : !products || products.length === 0 ? (
              <Card>
                <CardContent className="text-center py-12">
                  <Package className="h-12 w-12 mx-auto mb-4 text-gray-400" />
                  <h3 className="text-lg font-semibold mb-2">No products found</h3>
                  <p className="text-gray-600 mb-4">
                    {categoryInfo 
                      ? `No products found in ${categoryInfo.category.name} category.`
                      : "No products match your current filters."
                    }
                  </p>
                  {activeFilterCount > 0 && (
                    <Button onClick={() => setActiveFilters({})}>
                      Clear all filters
                    </Button>
                  )}
                </CardContent>
              </Card>
            ) : (
              <div className={cn(
                "grid gap-6",
                viewMode === 'grid' 
                  ? "grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4"
                  : "grid-cols-1"
              )}>
                {products.map((product: Product) => (
                  <ProductCard
                    key={product.id}
                    product={product}
                  />
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
