'use client'

import { useState, useEffect } from 'react'
import { ChevronRight, Tag, Package, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'
import { cn } from '@/lib/utils'

interface EnhancedCategoryFilterProps {
  selectedCategoryId?: string
  onCategoryChange: (categoryId: string | undefined) => void
  showProductCount?: boolean
  showSearch?: boolean
  className?: string
}

export function EnhancedCategoryFilter({ 
  selectedCategoryId, 
  onCategoryChange,
  showProductCount = true,
  showSearch = true,
  className
}: EnhancedCategoryFilterProps) {
  const { data: categories, isLoading } = useCategories()
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(new Set())
  const [categoryPath, setCategoryPath] = useState<Category[]>([])
  const [searchTerm, setSearchTerm] = useState('')
  const [filteredCategories, setFilteredCategories] = useState<Category[]>([])

  // Auto-expand path to selected category
  useEffect(() => {
    if (selectedCategoryId && categories) {
      const findPathAndExpand = (categoryId: string, currentPath: Category[] = []): Category[] | null => {
        for (const category of categories) {
          if (category.id === categoryId) {
            const fullPath = [...currentPath, category]
            // Auto-expand all categories in path
            const pathIds = fullPath.map(c => c.id)
            setExpandedCategories(prev => new Set([...prev, ...pathIds]))
            return fullPath
          }
          
          if (category.children && category.children.length > 0) {
            const childPath = findPathAndExpand(categoryId, [...currentPath, category])
            if (childPath) return childPath
          }
        }
        return null
      }
      
      const path = findPathAndExpand(selectedCategoryId)
      setCategoryPath(path || [])
    } else {
      setCategoryPath([])
    }
  }, [selectedCategoryId, categories])

  // Filter categories by search term
  useEffect(() => {
    if (!categories) return
    
    if (!searchTerm.trim()) {
      setFilteredCategories(categories)
      return
    }

    const filterRecursive = (cats: Category[]): Category[] => {
      return cats.reduce((acc: Category[], category) => {
        const matchesSearch = category.name.toLowerCase().includes(searchTerm.toLowerCase())
        const filteredChildren = category.children ? filterRecursive(category.children) : []
        
        if (matchesSearch || filteredChildren.length > 0) {
          acc.push({
            ...category,
            children: filteredChildren.length > 0 ? filteredChildren : category.children
          })
          // Auto-expand categories that match search
          if (matchesSearch || filteredChildren.length > 0) {
            setExpandedCategories(prev => new Set([...prev, category.id]))
          }
        }
        
        return acc
      }, [])
    }

    setFilteredCategories(filterRecursive(categories))
  }, [categories, searchTerm])

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

  const renderCategory = (category: Category, level = 0) => {
    const hasChildren = category.children && category.children.length > 0
    const isExpanded = expandedCategories.has(category.id)
    const isSelected = selectedCategoryId === category.id
    const isInPath = categoryPath.some(c => c.id === category.id)

    return (
      <div key={category.id} className="space-y-1">
        <div
          className={cn(
            'flex items-center space-x-2 p-2 rounded-lg cursor-pointer transition-all duration-200',
            'hover:bg-gray-700',
            isSelected && 'bg-orange-500 text-white ring-1 ring-orange-400',
            isInPath && !isSelected && 'bg-gray-700 text-gray-200',
          )}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          onClick={() => onCategoryChange(isSelected ? undefined : category.id)}
        >
          {hasChildren && (
            <Button
              variant="ghost"
              size="sm"
              className="h-4 w-4 p-0 hover:bg-gray-600 text-gray-300"
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

          <Tag className={cn(
            'h-4 w-4',
            isSelected ? 'text-white' : 'text-gray-400'
          )} />

          <span className={cn(
            "flex-1 text-sm font-medium",
            isSelected ? 'text-white' : 'text-gray-300'
          )}>{category.name}</span>

          {showProductCount && category.product_count !== undefined && (
            <Badge
              className={cn(
                "text-xs",
                isSelected
                  ? "bg-white/20 text-white"
                  : "bg-gray-700 text-gray-300"
              )}
            >
              {category.product_count}
            </Badge>
          )}
        </div>

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
      <Card className={cn(className, "bg-gray-800 border-gray-700")}>
        <CardHeader>
          <CardTitle className="text-white">Categories</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="h-8 bg-gray-700 rounded animate-pulse" />
            ))}
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card className={cn(className, "bg-gray-800 border-gray-700")}>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center justify-between text-base">
          <div className="flex items-center space-x-2">
            <Package className="h-5 w-5 text-orange-400" />
            <span className="text-white">Categories</span>
          </div>
          {selectedCategoryId && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onCategoryChange(undefined)}
              className="text-xs text-orange-400 hover:text-orange-300 hover:bg-orange-500/10"
            >
              Clear
            </Button>
          )}
        </CardTitle>
      </CardHeader>

      <CardContent className="pt-0">
        {/* Search Categories */}
        {showSearch && (
          <div className="mb-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search categories..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 text-sm bg-gray-700 border-gray-600 text-white placeholder-gray-400"
              />
            </div>
          </div>
        )}

        {/* Category Breadcrumbs */}
        {categoryPath.length > 0 && (
          <div className="mb-3 p-3 bg-gradient-to-r from-orange-500/10 to-orange-400/10 rounded-lg border border-orange-500/20">
            <p className="text-xs font-medium text-orange-400 mb-2">Current Selection:</p>
            <div className="flex items-center space-x-1 text-sm">
              {categoryPath.map((category, index) => (
                <div key={category.id} className="flex items-center space-x-1">
                  <button
                    onClick={() => onCategoryChange(category.id)}
                    className="text-orange-300 hover:text-orange-200 hover:underline font-medium transition-colors"
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
        <div className="space-y-1 max-h-80 overflow-y-auto">
          {filteredCategories?.filter(cat => !cat.parent_id).map(category => renderCategory(category))}
        </div>

        {/* Search Behavior Info */}
        {selectedCategoryId && (
          <div className="mt-3 p-3 bg-blue-500/10 rounded-lg text-xs text-blue-300 border border-blue-500/20">
            <div className="flex items-start space-x-2">
              <Search className="h-4 w-4 mt-0.5 flex-shrink-0" />
              <div>
                <p className="font-medium mb-1">🔍 Smart Search Active</p>
                <p className="text-blue-200">
                  • Products in this category<br/>
                  • Products in all subcategories<br/>
                  • Includes nested hierarchy
                </p>
              </div>
            </div>
          </div>
        )}

        {searchTerm && filteredCategories.length === 0 && (
          <div className="text-center py-8 text-gray-400">
            <Package className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No categories found for "{searchTerm}"</p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
