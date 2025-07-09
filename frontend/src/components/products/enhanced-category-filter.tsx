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
            'flex items-center space-x-2 p-3 rounded-lg cursor-pointer transition-all duration-300 font-medium group',
            'hover:bg-white/[0.08] hover:border-white/20 border border-transparent backdrop-blur-sm',
            isSelected && 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white border-[#ff9000]/50 shadow-md shadow-[#ff9000]/20',
            isInPath && !isSelected && 'bg-white/[0.05] border-white/15 text-gray-200',
          )}
          style={{ paddingLeft: `${level * 16 + 8}px` }}
          onClick={() => onCategoryChange(isSelected ? undefined : category.id)}
        >
          {hasChildren && (
            <Button
              variant="ghost"
              size="sm"
              className="h-6 w-6 p-0 hover:bg-white/10 rounded-md transition-all duration-300"
              onClick={(e) => {
                e.stopPropagation()
                toggleExpanded(category.id)
              }}
            >
              <ChevronRight
                className={cn(
                  'h-3 w-3 transition-transform duration-300',
                  isExpanded && 'rotate-90',
                  isSelected ? 'text-white' : 'text-gray-400'
                )}
              />
            </Button>
          )}

          <div className="p-1 rounded-md transition-all duration-300">
            <Tag className={cn(
              'h-3.5 w-3.5 transition-colors duration-300',
              isSelected ? 'text-white' : 'text-[#ff9000]'
            )} />
          </div>

          <span className={cn(
            "flex-1 text-sm font-medium transition-colors duration-300",
            isSelected ? 'text-white' : 'text-gray-300 group-hover:text-white'
          )}>{category.name}</span>

          {showProductCount && category.product_count !== undefined && (
            <Badge
              className={cn(
                "text-xs font-medium transition-all duration-300",
                isSelected
                  ? "bg-white/20 text-white border border-white/30"
                  : "bg-white/10 text-gray-300 border border-white/15"
              )}
            >
              {category.product_count}
            </Badge>
          )}

          <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-orange-400">
            →
          </div>
        </div>

        {hasChildren && isExpanded && (
          <div className="space-y-1 pl-2 border-l border-white/10 ml-3">
            {category.children!.map(child => renderCategory(child, level + 1))}
          </div>
        )}
      </div>
    )
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="flex items-center gap-2 mb-4">
          <div className="p-1.5 bg-gradient-to-br from-orange-500/20 to-orange-600/20 rounded-lg border border-orange-500/30">
            <Package className="w-4 h-4 text-orange-400" />
          </div>
          <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
            Categories
          </h3>
        </div>
        <div className="space-y-2">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-12 bg-white/[0.05] border border-white/10 rounded-lg animate-pulse backdrop-blur-sm" />
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="p-1.5 bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/20 rounded-lg border border-[#ff9000]/30">
            <Package className="w-4 h-4 text-[#ff9000]" />
          </div>
          <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
            Categories
          </h3>
        </div>
        {selectedCategoryId && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onCategoryChange(undefined)}
            className="text-xs text-[#ff9000] hover:text-white hover:bg-[#ff9000]/20 border border-[#ff9000]/30 hover:border-[#ff9000]/50 rounded-lg px-2 py-1 font-medium transition-all duration-300 hover:scale-105"
          >
            Clear
          </Button>
        )}
      </div>

      {/* Search Categories */}
      {showSearch && (
        <div className="relative">
          <div className="absolute left-3 top-1/2 transform -translate-y-1/2 z-10">
            <div className="p-1 bg-[#ff9000]/10 rounded border border-[#ff9000]/20">
              <Search className="h-3 w-3 text-[#ff9000]" />
            </div>
          </div>
          <Input
            placeholder="Search categories..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-11 text-sm bg-white/[0.05] border-white/15 text-white placeholder-gray-400 backdrop-blur-sm focus:bg-white/[0.08] focus:border-[#ff9000]/50 transition-all duration-300 rounded-lg"
          />
        </div>
      )}

      {/* Category Breadcrumbs */}
      {categoryPath.length > 0 && (
        <div className="p-4 bg-gradient-to-r from-[#ff9000]/10 to-[#ff9000]/10 rounded-lg border border-[#ff9000]/20 backdrop-blur-sm">
          <p className="text-xs font-medium text-[#ff9000] mb-2">Current Selection:</p>
          <div className="flex items-center space-x-1 text-sm">
            {categoryPath.map((category, index) => (
              <div key={category.id} className="flex items-center space-x-1">
                <button
                  onClick={() => onCategoryChange(category.id)}
                  className="text-[#ff9000] hover:text-white hover:underline font-medium transition-colors duration-300"
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
      <div className="space-y-2 max-h-80 overflow-y-auto scrollbar-thin scrollbar-thumb-white/10 scrollbar-track-transparent">
        {filteredCategories?.filter(cat => !cat.parent_id).map(category => renderCategory(category))}
      </div>

      {/* Search Behavior Info */}
      {selectedCategoryId && (
        <div className="p-4 bg-gradient-to-r from-blue-500/10 to-blue-600/10 rounded-lg border border-blue-500/20 backdrop-blur-sm">
          <div className="flex items-start space-x-2">
            <div className="p-1 bg-blue-500/20 rounded border border-blue-500/30">
              <Search className="h-3 w-3 text-blue-400" />
            </div>
            <div>
              <p className="font-medium mb-1 text-blue-300">🔍 Smart Search Active</p>
              <p className="text-xs text-blue-200 leading-relaxed">
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
          <div className="p-3 bg-white/[0.05] rounded-lg border border-white/10 max-w-sm mx-auto">
            <Package className="h-8 w-8 mx-auto mb-2 opacity-50" />
            <p className="text-sm">No categories found for "{searchTerm}"</p>
          </div>
        </div>
      )}
    </div>
  )
}
