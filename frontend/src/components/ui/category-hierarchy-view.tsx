'use client'

import React, { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { CategoryTree } from '@/components/ui/category-tree'
import { CategoryCard } from '@/components/ui/category-card'
import { ResponsiveGrid, GridItem } from '@/components/ui/responsive-grid'
import { 
  TreePine, 
  Grid3X3, 
  List,
  Filter,
  Search,
  ChevronDown,
  Eye,
  Package
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Category } from '@/types'

interface CategoryHierarchyViewProps {
  categories: Category[]
  className?: string
  viewMode?: 'tree' | 'grid' | 'list'
  showFilters?: boolean
  filteredCategories?: Category[]
  searchQuery?: string
}

type ViewMode = 'tree' | 'grid' | 'list'
type FilterMode = 'all' | 'root' | 'children' | 'leaf'

export function CategoryHierarchyView({
  categories,
  className,
  viewMode = 'tree',
  showFilters = true,
  filteredCategories: externalFilteredCategories,
  searchQuery: externalSearchQuery = ''
}: CategoryHierarchyViewProps) {
  const [filterMode, setFilterMode] = useState<FilterMode>('all')
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>()

  // Use external filtered categories if provided, otherwise filter internally
  const filteredCategories = React.useMemo(() => {
    if (externalFilteredCategories) {
      return externalFilteredCategories
    }

    let filtered = categories

    // Apply filter mode
    switch (filterMode) {
      case 'root':
        filtered = categories.filter(cat => !cat.parent_id)
        break
      case 'children':
        filtered = categories.filter(cat => cat.parent_id && cat.children && cat.children.length > 0)
        break
      case 'leaf':
        filtered = categories.filter(cat => !cat.children || cat.children.length === 0)
        break
      default:
        filtered = categories
    }

    return filtered
  }, [categories, filterMode, externalFilteredCategories])

  // Flatten categories for grid/list view
  const flatCategories = React.useMemo(() => {
    const flat: Category[] = []
    
    function flatten(cats: Category[], level = 0) {
      cats.forEach(cat => {
        flat.push({ ...cat, level })
        if (cat.children && cat.children.length > 0) {
          flatten(cat.children, level + 1)
        }
      })
    }
    
    flatten(filteredCategories)
    return flat
  }, [filteredCategories])

  const handleCategorySelect = (category: Category) => {
    setSelectedCategoryId(category.id)
  }

  // Get statistics
  const stats = React.useMemo(() => {
    const totalCategories = flatCategories.length
    const rootCategories = flatCategories.filter(cat => !cat.parent_id).length
    const leafCategories = flatCategories.filter(cat => !cat.children || cat.children.length === 0).length
    const totalProducts = flatCategories.reduce((sum, cat) => sum + (cat.product_count || 0), 0)
    
    return {
      totalCategories,
      rootCategories,
      leafCategories,
      totalProducts
    }
  }, [flatCategories])

  return (
    <div className={cn('space-y-4', className)}>
      {/* Simple Header */}
      <div className="flex items-center gap-3 mb-4">
        <TreePine className="w-5 h-5 text-[#ff9000]" />
        <h3 className="text-lg font-semibold text-white">Category Hierarchy</h3>
      </div>



      {/* Content based on view mode */}
      {viewMode === 'tree' && (
        <CategoryTree
          categories={filteredCategories}
          showProductCount={true}
          onCategorySelect={handleCategorySelect}
          selectedCategoryId={selectedCategoryId}
          maxDepth={5}
        />
      )}

      {viewMode === 'grid' && (
        <ResponsiveGrid variant="categories" gap="lg">
          {flatCategories.map((category, index) => (
            <GridItem key={category.id} index={index} animationDelay={50}>
              <CategoryCard
                category={category}
                variant="default"
                showStats={true}
                showSubcategories={true}
                priority={index < 4}
                className={cn(
                  'relative',
                  category.level && category.level > 0 && 'ml-4 border-l-4 border-l-blue-400/50'
                )}
              />
              {/* Level indicator for grid view */}
              {category.level && category.level > 0 && (
                <Badge 
                  className="absolute top-2 left-2 bg-blue-500/20 text-blue-400 text-xs"
                >
                  L{category.level}
                </Badge>
              )}
            </GridItem>
          ))}
        </ResponsiveGrid>
      )}

      {viewMode === 'list' && (
        <div className="space-y-3">
          {flatCategories.map((category, index) => (
            <div
              key={category.id}
              className={cn(
                'animate-in fade-in slide-in-from-left-4 duration-500',
                category.level && category.level > 0 && `ml-${Math.min(category.level * 8, 32)}`
              )}
              style={{ animationDelay: `${index * 100}ms` }}
            >
              <CategoryCard
                category={category}
                variant="list"
                showStats={true}
                showSubcategories={false}
                className={cn(
                  'relative',
                  category.level && category.level > 0 && 'border-l-4 border-l-blue-400/50'
                )}
              />
              {/* Level indicator for list view */}
              {category.level && category.level > 0 && (
                <div className="absolute left-2 top-1/2 transform -translate-y-1/2">
                  <Badge className="bg-blue-500/20 text-blue-400 text-xs">
                    L{category.level}
                  </Badge>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Empty State */}
      {filteredCategories.length === 0 && (
        <div className="text-center py-12 bg-white/[0.04] rounded-lg border border-white/10">
          <TreePine className="w-12 h-12 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-white mb-2">No Categories Found</h3>
          <p className="text-gray-400 mb-4">
            {externalSearchQuery
              ? `No categories match "${externalSearchQuery}"`
              : 'No categories match the current filter'
            }
          </p>
          {!externalFilteredCategories && filterMode !== 'all' && (
            <Button
              variant="outline"
              onClick={() => setFilterMode('all')}
              className="border-white/20 text-white hover:bg-white/10"
            >
              Clear Filters
            </Button>
          )}
        </div>
      )}
    </div>
  )
}
