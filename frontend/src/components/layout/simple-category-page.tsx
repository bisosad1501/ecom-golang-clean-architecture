'use client'

import { useState, useEffect } from 'react'
import { useSearchParams } from 'next/navigation'
import { Package, Grid, List } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { CategoryBreadcrumbs } from '@/components/layout/category-breadcrumbs'
import { EnhancedCategorySidebar } from '@/components/layout/enhanced-category-sidebar'
import { useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'
import { cn } from '@/lib/utils'
import { apiClient } from '@/lib/api'

interface SimpleCategoryPageProps {
  categoryId?: string
  className?: string
}

interface CategoryInfo {
  category: Category
  productCount: number
  subcategories: Category[]
}

export function SimpleCategoryPage({ categoryId, className }: SimpleCategoryPageProps) {
  const searchParams = useSearchParams()
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [categoryInfo, setCategoryInfo] = useState<CategoryInfo | null>(null)
  const [isLoading, setIsLoading] = useState(false)

  const { data: categories } = useCategories()

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
        const productCount = (countResponse.data as any)?.data?.product_count || 0
        
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
                      <span>{categoryInfo.subcategories.length} subcategories</span>
                    </div>
                  )}
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
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-6">
        <div className="flex gap-6">
          {/* Sidebar - Desktop */}
          <div className="hidden lg:block w-80 flex-shrink-0">
            <div className="sticky top-6">
              <EnhancedCategorySidebar
                selectedCategoryId={categoryId}
                onCategoryChange={handleCategoryChange}
                showProductCount={true}
                showSearch={true}
                showFilters={true}
              />
            </div>
          </div>

          {/* Main Content Area */}
          <div className="flex-1 min-w-0">
            {/* Toolbar */}
            <div className="bg-background rounded-lg border p-4 mb-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-4">
                  <div className="text-sm text-gray-600">
                    {isLoading ? (
                      <div className="h-5 w-32 bg-gray-200 rounded animate-pulse" />
                    ) : (
                      <span>
                        Category: {categoryInfo?.category.name || 'All Categories'}
                      </span>
                    )}
                  </div>
                </div>

                <div className="flex items-center space-x-2">
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

            {/* Demo Content */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center">
                  <Package className="h-5 w-5 mr-2" />
                  Category System Demo
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <p className="text-gray-600">
                    ✅ <strong>Category Mega Menu</strong> - Hover over "All Categories" in header
                  </p>
                  <p className="text-gray-600">
                    ✅ <strong>Category Breadcrumbs</strong> - Shows path from root to current category
                  </p>
                  <p className="text-gray-600">
                    ✅ <strong>Enhanced Category Sidebar</strong> - Smart filtering, search, and hierarchy
                  </p>
                  <p className="text-gray-600">
                    ✅ <strong>Clean Header Design</strong> - Single unified header instead of multiple bars
                  </p>
                  <p className="text-gray-600">
                    ✅ <strong>Responsive Design</strong> - Works on all devices
                  </p>
                  
                  {categoryInfo && (
                    <div className="mt-6 p-4 bg-blue-50 rounded-lg">
                      <h3 className="font-semibold text-blue-900 mb-2">Current Category Info:</h3>
                      <ul className="space-y-1 text-blue-800">
                        <li><strong>Name:</strong> {categoryInfo.category.name}</li>
                        <li><strong>Products:</strong> {categoryInfo.productCount}</li>
                        <li><strong>Subcategories:</strong> {categoryInfo.subcategories.length}</li>
                        {categoryInfo.category.description && (
                          <li><strong>Description:</strong> {categoryInfo.category.description}</li>
                        )}
                      </ul>
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  )
}
