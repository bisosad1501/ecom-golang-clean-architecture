'use client'

import { useState, useEffect, useMemo } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { 
  ChevronRight, 
  ChevronDown, 
  Tag, 
  Package, 
  Search, 
  Filter,
  X,
  Star,
  BarChart3,
  TrendingUp,
  Clock
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from '@/components/ui/collapsible'
import { useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'
import { cn } from '@/lib/utils'
import { apiClient } from '@/lib/api'

interface EnhancedCategorySidebarProps {
  selectedCategoryId?: string
  onCategoryChange: (categoryId: string | undefined) => void
  showProductCount?: boolean
  showSearch?: boolean
  showFilters?: boolean
  className?: string
}

interface CategoryStats {
  [categoryId: string]: {
    productCount: number
    subcategoryCount: number
    avgRating?: number
    isPopular?: boolean
    isTrending?: boolean
  }
}

export function EnhancedCategorySidebar({ 
  selectedCategoryId, 
  onCategoryChange,
  showProductCount = true,
  showSearch = true,
  showFilters = true,
  className
}: EnhancedCategorySidebarProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { data: categories, isLoading } = useCategories()
  
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set())
  const [categoryPath, setCategoryPath] = useState<Category[]>([])
  const [searchTerm, setSearchTerm] = useState('')
  const [filteredCategories, setFilteredCategories] = useState<Category[]>([])
  const [categoryStats, setCategoryStats] = useState<CategoryStats>({})
  const [showOnlyWithProducts, setShowOnlyWithProducts] = useState(false)
  const [sortBy, setSortBy] = useState<'name' | 'products' | 'popularity'>('name')

  // Build category tree structure
  const categoryTree = useMemo(() => {
    if (!categories) return []
    
    const categoryMap = new Map<string, Category & { children: Category[] }>()
    
    // Initialize all categories with empty children array
    categories.forEach(cat => {
      categoryMap.set(cat.id, { ...cat, children: [] })
    })
    
    // Build tree structure
    const rootCategories: (Category & { children: Category[] })[] = []
    
    categories.forEach(cat => {
      const categoryWithChildren = categoryMap.get(cat.id)!
      
      if (cat.parent_id) {
        const parent = categoryMap.get(cat.parent_id)
        if (parent) {
          parent.children.push(categoryWithChildren)
        }
      } else {
        rootCategories.push(categoryWithChildren)
      }
    })
    
    return rootCategories
  }, [categories])

  // Fetch category statistics
  useEffect(() => {
    const fetchCategoryStats = async () => {
      if (!categories) return
      
      const stats: CategoryStats = {}
      
      try {
        // Fetch product counts for all categories in parallel
        const promises = categories.map(async (category) => {
          try {
            const response = await apiClient.get(`/categories/${category.id}/count`)
            const data = (response.data as any)?.data
            return {
              categoryId: category.id,
              productCount: data?.product_count || 0
            }
          } catch (error) {
            return {
              categoryId: category.id,
              productCount: 0
            }
          }
        })
        
        const results = await Promise.all(promises)
        
        results.forEach(({ categoryId, productCount }) => {
          const category = categories.find(c => c.id === categoryId)
          const subcategoryCount = categories.filter(c => c.parent_id === categoryId).length
          
          stats[categoryId] = {
            productCount,
            subcategoryCount,
            isPopular: productCount > 10, // Arbitrary threshold
            isTrending: productCount > 5 && subcategoryCount > 0
          }
        })
        
        setCategoryStats(stats)
      } catch (error) {
        console.error('Error fetching category stats:', error)
      }
    }
    
    fetchCategoryStats()
  }, [categories])

  // Auto-expand path to selected category
  useEffect(() => {
    if (selectedCategoryId && categories) {
      const findPathAndExpand = (categoryId: string): Category[] | null => {
        const category = categories.find(c => c.id === categoryId)
        if (!category) return null
        
        if (category.parent_id) {
          const parentPath = findPathAndExpand(category.parent_id)
          if (parentPath) {
            return [...parentPath, category]
          }
        }
        
        return [category]
      }
      
      const path = findPathAndExpand(selectedCategoryId)
      if (path) {
        setCategoryPath(path)
        // Auto-expand all categories in path
        const pathIds = path.map(c => c.id)
        const parentIds = path.slice(0, -1).map(c => c.id)
        setExpandedCategories(prev => new Set([...prev, ...parentIds]))
      }
    } else {
      setCategoryPath([])
    }
  }, [selectedCategoryId, categories])

  // Filter and sort categories
  useEffect(() => {
    if (!categoryTree) return
    
    const filterAndSort = (cats: (Category & { children: Category[] })[]): (Category & { children: Category[] })[] => {
      let filtered = cats.filter(category => {
        // Text search filter
        const matchesSearch = !searchTerm.trim() || 
          category.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
          category.description?.toLowerCase().includes(searchTerm.toLowerCase())
        
        // Products filter
        const hasProducts = !showOnlyWithProducts || 
          (categoryStats[category.id]?.productCount || 0) > 0
        
        return matchesSearch && hasProducts
      })
      
      // Recursively filter children
      filtered = filtered.map(category => ({
        ...category,
        children: filterAndSort(category.children as (Category & { children: Category[] })[])
      }))
      
      // Sort categories
      filtered.sort((a, b) => {
        switch (sortBy) {
          case 'products':
            return (categoryStats[b.id]?.productCount || 0) - (categoryStats[a.id]?.productCount || 0)
          case 'popularity':
            const aScore = (categoryStats[a.id]?.productCount || 0) + (categoryStats[a.id]?.subcategoryCount || 0) * 2
            const bScore = (categoryStats[b.id]?.productCount || 0) + (categoryStats[b.id]?.subcategoryCount || 0) * 2
            return bScore - aScore
          default:
            return a.name.localeCompare(b.name)
        }
      })
      
      return filtered
    }
    
    setFilteredCategories(filterAndSort(categoryTree))
  }, [categoryTree, searchTerm, showOnlyWithProducts, sortBy, categoryStats])

  const toggleExpanded = (categoryId: string) => {
    setExpandedCategories(prev => {
      const newSet = new Set(prev)
      if (newSet.has(categoryId)) {
        newSet.delete(categoryId)
      } else {
        newSet.add(categoryId)
      }
      return newSet
    })
  }

  const handleCategorySelect = (categoryId: string) => {
    const newCategoryId = selectedCategoryId === categoryId ? undefined : categoryId
    onCategoryChange(newCategoryId)
    
    // Update URL
    const params = new URLSearchParams(searchParams.toString())
    if (newCategoryId) {
      params.set('category', newCategoryId)
    } else {
      params.delete('category')
    }
    router.push(`?${params.toString()}`)
  }

  const renderCategory = (category: Category & { children?: Category[] }, level = 0) => {
    const hasChildren = category.children && category.children.length > 0
    const isExpanded = expandedCategories.has(category.id)
    const isSelected = selectedCategoryId === category.id
    const isInPath = categoryPath.some(c => c.id === category.id)
    const stats = categoryStats[category.id]

    return (
      <div key={category.id} className="space-y-1">
        <div 
          className={cn(
            'flex items-center space-x-2 p-2 rounded-lg cursor-pointer transition-all duration-200 group',
            'hover:bg-gray-50',
            isSelected && 'bg-primary-100 text-primary-800 ring-1 ring-primary-200 shadow-sm',
            isInPath && !isSelected && 'bg-primary-25 text-primary-700',
          )}
          style={{ paddingLeft: `${level * 12 + 8}px` }}
        >
          {/* Expand/Collapse Button */}
          {hasChildren && (
            <Button
              variant="ghost"
              size="sm"
              className="h-5 w-5 p-0 hover:bg-gray-200 rounded"
              onClick={(e) => {
                e.stopPropagation()
                toggleExpanded(category.id)
              }}
            >
              <ChevronRight 
                className={cn(
                  'h-3 w-3 transition-transform duration-200',
                  isExpanded && 'rotate-90'
                )} 
              />
            </Button>
          )}
          
          {/* Category Content */}
          <div 
            className="flex-1 flex items-center justify-between"
            onClick={() => handleCategorySelect(category.id)}
          >
            <div className="flex items-center space-x-2">
              <Tag className={cn(
                'h-4 w-4',
                isSelected ? 'text-primary-600' : 'text-gray-500 group-hover:text-gray-700'
              )} />
              
              <div className="flex-1">
                <div className="flex items-center space-x-2">
                  <span className="text-sm font-medium">{category.name}</span>
                  
                  {/* Category badges */}
                  <div className="flex items-center space-x-1">
                    {stats?.isPopular && (
                      <Badge variant="outline" className="text-xs px-1 py-0 h-4">
                        <Star className="h-2 w-2 mr-1" />
                        Popular
                      </Badge>
                    )}
                    
                    {stats?.isTrending && (
                      <Badge variant="outline" className="text-xs px-1 py-0 h-4 border-emerald-200 text-emerald-700">
                        <TrendingUp className="h-2 w-2 mr-1" />
                        Trending
                      </Badge>
                    )}
                  </div>
                </div>
                
                {category.description && (
                  <div className="text-xs text-gray-500 truncate max-w-40">
                    {category.description}
                  </div>
                )}
              </div>
            </div>
            
            {/* Product Count */}
            {showProductCount && stats && (
              <div className="flex items-center space-x-1">
                {stats.productCount > 0 && (
                  <Badge 
                    variant={isSelected ? "default" : "secondary"} 
                    className="text-xs"
                  >
                    {stats.productCount}
                  </Badge>
                )}
                
                {stats.subcategoryCount > 0 && (
                  <Badge variant="outline" className="text-xs">
                    +{stats.subcategoryCount}
                  </Badge>
                )}
              </div>
            )}
          </div>
        </div>

        {/* Children */}
        {hasChildren && isExpanded && (
          <div className="space-y-1">
            {category.children!.map(child => renderCategory(child, level + 1))}
          </div>
        )}
      </div>
    )
  }

  if (isLoading) {
    return (
      <Card className={className}>
        <CardHeader>
          <CardTitle>Categories</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="h-10 bg-gray-200 rounded animate-pulse" />
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className={className}>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <Package className="h-5 w-5" />
            <span>Categories</span>
          </div>
          {selectedCategoryId && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onCategoryChange(undefined)}
              className="text-xs text-primary-600 hover:text-primary-800"
            >
              <X className="h-3 w-3 mr-1" />
              Clear
            </Button>
          )}
        </CardTitle>
      </CardHeader>
      
      <CardContent className="space-y-4">
        {/* Search */}
        {showSearch && (
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input
              placeholder="Search categories..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 text-sm"
            />
          </div>
        )}

        {/* Filters */}
        {showFilters && (
          <Collapsible>
            <CollapsibleTrigger asChild>
              <Button variant="ghost" className="w-full justify-between p-2 h-auto">
                <div className="flex items-center space-x-2">
                  <Filter className="h-4 w-4" />
                  <span className="text-sm">Filters & Sort</span>
                </div>
                <ChevronDown className="h-4 w-4" />
              </Button>
            </CollapsibleTrigger>
            <CollapsibleContent className="space-y-3 pt-2">
              <div className="space-y-2">
                <label className="flex items-center space-x-2 text-sm">
                  <input
                    type="checkbox"
                    checked={showOnlyWithProducts}
                    onChange={(e) => setShowOnlyWithProducts(e.target.checked)}
                    className="rounded"
                  />
                  <span>Only categories with products</span>
                </label>
              </div>
              
              <Separator />
              
              <div className="space-y-2">
                <div className="text-sm font-medium">Sort by:</div>
                <div className="space-y-1">
                  {[
                    { value: 'name', label: 'Name', icon: Tag },
                    { value: 'products', label: 'Product Count', icon: BarChart3 },
                    { value: 'popularity', label: 'Popularity', icon: TrendingUp }
                  ].map(({ value, label, icon: Icon }) => (
                    <label key={value} className="flex items-center space-x-2 text-sm cursor-pointer">
                      <input
                        type="radio"
                        name="sortBy"
                        value={value}
                        checked={sortBy === value}
                        onChange={(e) => setSortBy(e.target.value as any)}
                        className="rounded"
                      />
                      <Icon className="h-3 w-3" />
                      <span>{label}</span>
                    </label>
                  ))}
                </div>
              </div>
            </CollapsibleContent>
          </Collapsible>
        )}

        {/* Category Breadcrumbs */}
        {categoryPath.length > 0 && (
          <div className="p-3 bg-gradient-to-r from-primary-50 to-blue-50 rounded-lg border border-primary-100">
            <p className="text-xs font-medium text-primary-700 mb-2 flex items-center">
              <Clock className="h-3 w-3 mr-1" />
              Current Selection:
            </p>
            <div className="flex items-center space-x-1 text-sm">
              {categoryPath.map((category, index) => (
                <div key={category.id} className="flex items-center space-x-1">
                  <button
                    onClick={() => handleCategorySelect(category.id)}
                    className="text-primary-600 hover:text-primary-800 hover:underline font-medium transition-colors"
                  >
                    {category.name}
                  </button>
                  {index < categoryPath.length - 1 && (
                    <ChevronRight className="h-3 w-3 text-gray-400" />
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Category Tree */}
        <div className="space-y-1 max-h-96 overflow-y-auto">
          {filteredCategories.map(category => renderCategory(category))}
        </div>
        
        {/* Empty State */}
        {searchTerm && filteredCategories.length === 0 && (
          <div className="text-center py-8 text-gray-500">
            <Package className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No categories found for "{searchTerm}"</p>
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={() => setSearchTerm('')}
              className="mt-2 text-xs"
            >
              Clear search
            </Button>
          </div>
        )}

        {/* Smart Search Info */}
        {selectedCategoryId && (
          <div className="p-3 bg-blue-50 rounded-lg text-xs text-blue-700 border border-blue-200">
            <div className="flex items-start space-x-2">
              <Search className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <div>
                <p className="font-medium mb-1">🔍 Smart Search Active</p>
                <p className="text-blue-600">
                  Showing products from this category and all subcategories
                </p>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
