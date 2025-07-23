'use client'

import React, { useState, useMemo, useEffect } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Grid, List, ChevronRight, Home, FolderTree, TreePine, Filter } from 'lucide-react'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { CategoryCard } from '@/components/ui/category-card'
import { Pagination } from '@/components/ui/pagination'
import { ResponsiveGrid, GridItem, GridSkeleton } from '@/components/ui/responsive-grid'
import { CategoryFilters, CategoryFilterState } from '@/components/ui/category-filters'
import { CategoryHierarchyView } from '@/components/ui/category-hierarchy-view'
import { Category } from '@/types'
import { categoryService } from '@/lib/services/categories'
import { useQuery } from '@tanstack/react-query'
import { cn } from '@/lib/utils'

export function CategoriesPage() {
  const searchParams = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get('search') || '')
  const [viewMode, setViewMode] = useState<'grid' | 'list' | 'hierarchy'>('hierarchy')
  const [showFilters, setShowFilters] = useState(false)

  // Handle view mode change with localStorage persistence
  const handleViewModeChange = (mode: 'grid' | 'list' | 'hierarchy') => {
    setViewMode(mode)
    localStorage.setItem('categories-view-mode', mode)
  }

  // Load view mode from localStorage on mount
  useEffect(() => {
    const savedViewMode = localStorage.getItem('categories-view-mode') as 'grid' | 'list' | 'hierarchy'
    if (savedViewMode && ['grid', 'list', 'hierarchy'].includes(savedViewMode)) {
      setViewMode(savedViewMode)
    }
  }, [])
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage, setItemsPerPage] = useState(12)
  const [filters, setFilters] = useState<CategoryFilterState>({
    searchQuery: searchParams.get('search') || '',
    sortBy: 'name',
    sortOrder: 'asc',
    showActiveOnly: true,
    productCountRange: [0, 1000]
  })

  // Fetch category tree
  const {
    data: categoryTree,
    isLoading,
    error
  } = useQuery({
    queryKey: ['categories', 'tree'],
    queryFn: () => categoryService.getCategoryTree(),
  })

  // Flatten tree for easier filtering
  const flattenCategories = (categories: Category[]): Category[] => {
    const result: Category[] = []
    const flatten = (cats: Category[], level = 0) => {
      cats.forEach(cat => {
        result.push({ ...cat, level })
        if (cat.children && cat.children.length > 0) {
          flatten(cat.children, level + 1)
        }
      })
    }
    flatten(categories)
    return result
  }

  const flatCategories = useMemo(() => {
    return categoryTree ? flattenCategories(categoryTree) : []
  }, [categoryTree])

  // Enhanced filtering with better search capabilities
  const filteredCategories = useMemo(() => {
    let catsToFilter = flatCategories

    // Active filter
    if (filters.showActiveOnly) {
      catsToFilter = catsToFilter.filter(category => category.is_active !== false)
    }

    // Enhanced search filter - search in name, description, and slug
    if (filters.searchQuery.trim()) {
      const query = filters.searchQuery.toLowerCase().trim()
      catsToFilter = catsToFilter.filter(category => {
        const searchableText = [
          category.name,
          category.description || '',
          category.slug || ''
        ].join(' ').toLowerCase()

        return searchableText.includes(query)
      })
    }

    // Parent category filter
    if (filters.parentCategory) {
      if (filters.parentCategory === 'root') {
        catsToFilter = catsToFilter.filter(category => !category.parent_id)
      } else {
        catsToFilter = catsToFilter.filter(category =>
          category.parent_id === filters.parentCategory
        )
      }
    }

    // Subcategories filter
    if (filters.hasSubcategories !== undefined) {
      catsToFilter = catsToFilter.filter(category =>
        filters.hasSubcategories
          ? category.children && category.children.length > 0
          : !category.children || category.children.length === 0
      )
    }

    // Recently updated filter
    if (filters.recentlyUpdated) {
      const sevenDaysAgo = new Date()
      sevenDaysAgo.setDate(sevenDaysAgo.getDate() - 7)
      catsToFilter = catsToFilter.filter(category =>
        new Date(category.updated_at) > sevenDaysAgo
      )
    }

    // Product count range filter
    if (filters.productCountRange) {
      const [min, max] = filters.productCountRange
      catsToFilter = catsToFilter.filter(category => {
        const count = category.product_count || 0
        return count >= min && count <= max
      })
    }

    // Enhanced sorting
    catsToFilter.sort((a, b) => {
      let result = 0
      switch (filters.sortBy) {
        case 'name':
          result = a.name.localeCompare(b.name)
          break
        case 'product_count':
          result = (a.product_count || 0) - (b.product_count || 0)
          break
        case 'created_at':
          result = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
          break
        case 'updated_at':
          result = new Date(a.updated_at).getTime() - new Date(b.updated_at).getTime()
          break
        default:
          result = 0
      }
      return filters.sortOrder === 'desc' ? -result : result
    })

    return catsToFilter
  }, [flatCategories, filters])

  // Pagination logic
  const paginatedCategories = useMemo(() => {
    const startIndex = (currentPage - 1) * itemsPerPage
    const endIndex = startIndex + itemsPerPage
    return filteredCategories.slice(startIndex, endIndex)
  }, [filteredCategories, currentPage, itemsPerPage])

  const totalPages = Math.ceil(filteredCategories.length / itemsPerPage)

  // Reset to first page when filters change
  const resetPagination = () => {
    setCurrentPage(1)
  }

  // Update pagination when search or filters change
  useEffect(() => {
    resetPagination()
  }, [filters])

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    setFilters((prev: CategoryFilterState) => ({ ...prev, searchQuery: query }))
    const url = new URL(window.location.href)
    if (query) {
      url.searchParams.set('search', query)
    } else {
      url.searchParams.delete('search')
    }
    window.history.replaceState({}, '', url.toString())
  }

  const handleFiltersChange = (newFilters: CategoryFilterState) => {
    setFilters(newFilters)
    setSearchQuery(newFilters.searchQuery)
    resetPagination()
  }

  if (error) {
    return (
      <div className="min-h-screen bg-black">
        <div className="container mx-auto px-4 py-8">
          <Card className="p-8 text-center bg-gray-900 border-gray-700">
            <h2 className="text-2xl font-bold text-white mb-4">
              Oops! Something went wrong
            </h2>
            <p className="text-gray-300 mb-6">
              We couldn't load the categories. Please try again later.
            </p>
            <Button onClick={() => window.location.reload()} variant="outline" className="bg-[#ff9000] text-white border-[#ff9000] hover:bg-[#ff9000]/90">
              Try Again
            </Button>
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
                <CategoryFilters filters={filters} onFiltersChange={handleFiltersChange} />
              </div>
            </div>
          )}

          {/* Compact Sidebar Filters */}
          <aside className="hidden lg:block lg:w-60 xl:w-64 flex-shrink-0">
            <div className="sticky top-4 max-h-[calc(100vh-2rem)] overflow-y-auto">
              <CategoryFilters
                filters={filters}
                onFiltersChange={handleFiltersChange}
              />
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
                    Browse <span className="text-[#ff9000]">Categories</span>
                  </h1>
                  {filteredCategories.length > 0 && (
                    <div className="flex items-center gap-3">
                      <p className="text-gray-400 text-sm">
                        Showing <span className="text-[#ff9000] font-medium">{paginatedCategories.length}</span> of <span className="text-[#ff9000] font-medium">{filteredCategories.length}</span> categories
                      </p>
                      <div className="bg-white/8 text-gray-300 border border-white/15 px-2 py-1 text-xs backdrop-blur-sm font-medium rounded-md">
                        {filteredCategories.length} Total
                      </div>
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
                      onClick={() => handleViewModeChange('hierarchy')}
                      className={`rounded-md h-7 px-2 text-xs font-medium transition-all duration-200 border-0 hover:text-[#ff9000] hover:bg-[#ff9000]/10 ${
                        viewMode === 'hierarchy'
                          ? 'bg-white/10 text-white'
                          : 'text-gray-400'
                      }`}
                    >
                      <TreePine className="h-3.5 w-3.5" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleViewModeChange('grid')}
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
                      onClick={() => handleViewModeChange('list')}
                      className={`rounded-md h-7 px-2 text-xs font-medium transition-all duration-200 border-0 hover:text-[#ff9000] hover:bg-[#ff9000]/10 ${
                        viewMode === 'list'
                          ? 'bg-white/10 text-white'
                          : 'text-gray-400'
                      }`}
                    >
                      <List className="h-3.5 w-3.5" />
                    </Button>
                  </div>

                  {/* Items per page selector */}
                  <Select value={itemsPerPage.toString()} onValueChange={(value) => setItemsPerPage(Number(value))}>
                    <SelectTrigger className="w-16 h-7 bg-white/[0.06] backdrop-blur-md border-white/10 text-white text-xs rounded-lg">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent className="bg-gray-900 border-gray-700">
                      <SelectItem value="12">12</SelectItem>
                      <SelectItem value="24">24</SelectItem>
                      <SelectItem value="48">48</SelectItem>
                    </SelectContent>
                  </Select>

                  {/* Filters toggle (mobile) */}
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowFilters(!showFilters)}
                    className="lg:hidden rounded-lg border border-white/10 bg-white/[0.06] backdrop-blur-md text-gray-400 hover:bg-white/[0.08] hover:border-white/15 hover:text-white transition-all duration-200 h-7 px-2.5 text-xs font-medium shadow-sm"
                  >
                    <Filter className="h-3 w-3 mr-1" />
                    Filters
                  </Button>
                </div>
              </div>

              {/* Breadcrumb */}
              <nav className="flex items-center space-x-2 text-sm text-gray-400 mb-6">
                <Link href="/" className="hover:text-white transition-colors">
                  <Home className="h-4 w-4" />
                </Link>
                <ChevronRight className="h-4 w-4" />
                <span className="text-white">Categories</span>
              </nav>
            </div>

            {/* Content */}
            <div className="space-y-8">
              {isLoading ? (
                <GridSkeleton count={8} variant="categories" gap="lg" />
              ) : filteredCategories.length > 0 ? (
                <>
                  {viewMode === 'hierarchy' ? (
                    <div className="bg-white/[0.02] backdrop-blur-sm border border-white/10 rounded-xl p-6 shadow-lg">
                      <CategoryHierarchyView
                        categories={categoryTree || []}
                        filteredCategories={filteredCategories}
                        searchQuery={filters.searchQuery}
                        viewMode="tree"
                        showFilters={false}
                      />
                    </div>
                  ) : viewMode === 'grid' ? (
                    <ResponsiveGrid variant="categories" gap="lg">
                      {paginatedCategories.map((category, index) => (
                        <GridItem key={category.id} index={index} animationDelay={50}>
                          <CategoryCard
                            category={category}
                            variant="default"
                            showStats={true}
                            showSubcategories={true}
                            priority={index < 4} // Prioritize first 4 images
                          />
                        </GridItem>
                      ))}
                    </ResponsiveGrid>
                  ) : (
                    <div className="space-y-4">
                      {paginatedCategories.map((category, index) => (
                        <GridItem
                          key={category.id}
                          index={index}
                          animationDelay={100}
                          className="animate-in fade-in slide-in-from-left-4 duration-500"
                        >
                          <CategoryCard
                            category={category}
                            variant="list"
                            showStats={true}
                            showSubcategories={false}
                          />
                        </GridItem>
                      ))}
                    </div>
                  )}
                </>
              ) : (
                <div className="text-center py-20">
                  <div className="bg-white/[0.03] backdrop-blur-sm border border-white/10 rounded-xl p-12 shadow-lg max-w-md mx-auto">
                    <FolderTree className="h-16 w-16 text-gray-500 mx-auto mb-6" />
                    <h3 className="text-xl font-semibold text-white mb-3">
                      {searchQuery ? 'No matching categories' : 'No categories found'}
                    </h3>
                    <p className="text-gray-400 mb-6">
                      {searchQuery
                        ? `We couldn't find any categories matching "${searchQuery}".`
                        : 'No categories are available at the moment.'
                      }
                    </p>
                    {searchQuery && (
                      <Button
                        onClick={() => handleSearch('')}
                        variant="outline"
                        className="bg-[#ff9000] text-white border-[#ff9000] hover:bg-[#ff9000]/90"
                      >
                        Clear Search
                      </Button>
                    )}
                  </div>
                </div>
              )}
            </div>

            {/* Pagination - only show for grid/list views */}
            {viewMode !== 'hierarchy' && filteredCategories.length > itemsPerPage && (
              <div className="mt-8">
                <div className="flex justify-center">
                  <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    hasNext={currentPage < totalPages}
                    hasPrev={currentPage > 1}
                  />
                </div>
              </div>
            )}
          </main>
        </div>
      </div>
    </div>
  )
}