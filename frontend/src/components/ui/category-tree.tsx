'use client'

import React, { useState } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { 
  ChevronRight, 
  ChevronDown, 
  FolderTree, 
  Package,
  Eye,
  Plus,
  Minus,
  ArrowRight
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Category } from '@/types'

interface CategoryTreeProps {
  categories: Category[]
  className?: string
  showProductCount?: boolean
  expandAll?: boolean
  onCategorySelect?: (category: Category) => void
  selectedCategoryId?: string
  maxDepth?: number
}

interface CategoryNodeProps {
  category: Category
  level: number
  isExpanded: boolean
  onToggle: (categoryId: string) => void
  showProductCount: boolean
  onCategorySelect?: (category: Category) => void
  selectedCategoryId?: string
  maxDepth: number
}

function CategoryNode({
  category,
  level,
  isExpanded,
  onToggle,
  showProductCount,
  onCategorySelect,
  selectedCategoryId,
  maxDepth
}: CategoryNodeProps) {
  const hasChildren = category.children && category.children.length > 0
  const isSelected = selectedCategoryId === category.id
  const productCount = category.product_count || 0
  
  // Indentation based on level
  const indentClass = level > 0 ? `ml-${Math.min(level * 4, 16)}` : ''
  
  // Colors based on level
  const levelColors = [
    'border-l-[#ff9000]', // Level 0 - Orange
    'border-l-blue-400',  // Level 1 - Blue
    'border-l-green-400', // Level 2 - Green
    'border-l-purple-400', // Level 3 - Purple
    'border-l-pink-400',  // Level 4 - Pink
  ]
  
  const borderColor = levelColors[Math.min(level, levelColors.length - 1)]

  return (
    <div className={cn('relative', level > maxDepth && 'hidden')}>
      {/* Connection lines for tree structure */}
      {level > 0 && (
        <>
          {/* Vertical line */}
          <div 
            className="absolute left-2 top-0 bottom-0 w-px bg-white/20"
            style={{ left: `${(level - 1) * 16 + 8}px` }}
          />
          {/* Horizontal line */}
          <div 
            className="absolute top-6 w-4 h-px bg-white/20"
            style={{ left: `${(level - 1) * 16 + 8}px` }}
          />
        </>
      )}

      <Card className={cn(
        'mb-2 transition-all duration-200 border-l-4',
        borderColor,
        isSelected 
          ? 'bg-[#ff9000]/20 border-[#ff9000]/50 shadow-lg shadow-[#ff9000]/10' 
          : 'bg-white/[0.06] border-white/10 hover:bg-white/[0.08] hover:border-white/20',
        indentClass
      )}>
        <CardContent className="p-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3 flex-1 min-w-0">
              {/* Expand/Collapse Button */}
              {hasChildren ? (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => onToggle(category.id)}
                  className="h-6 w-6 p-0 text-gray-400 hover:text-white flex-shrink-0"
                >
                  {isExpanded ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </Button>
              ) : (
                <div className="w-6 h-6 flex items-center justify-center flex-shrink-0">
                  <div className="w-2 h-2 rounded-full bg-white/30" />
                </div>
              )}

              {/* Category Icon */}
              <div className={cn(
                'w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0',
                level === 0 ? 'bg-[#ff9000]/20' : 'bg-white/10'
              )}>
                <FolderTree className={cn(
                  'w-4 h-4',
                  level === 0 ? 'text-[#ff9000]' : 'text-white/70'
                )} />
              </div>

              {/* Category Info */}
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2">
                  <h4 className={cn(
                    'font-medium truncate transition-colors',
                    isSelected ? 'text-[#ff9000]' : 'text-white hover:text-[#ff9000]',
                    level === 0 ? 'text-base' : 'text-sm'
                  )}>
                    {category.name}
                  </h4>
                  
                  {/* Level indicator */}
                  <Badge 
                    variant="outline" 
                    className={cn(
                      'text-xs px-1.5 py-0.5 border-current',
                      level === 0 ? 'text-[#ff9000] border-[#ff9000]/30' : 'text-white/60 border-white/20'
                    )}
                  >
                    L{level}
                  </Badge>
                </div>
                
                {category.description && (
                  <p className="text-xs text-gray-400 truncate mt-1">
                    {category.description}
                  </p>
                )}
              </div>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-2 flex-shrink-0">
              {/* Product Count */}
              {showProductCount && (
                <Badge variant="secondary" className="bg-white/10 text-white/80 text-xs">
                  <Package className="w-3 h-3 mr-1" />
                  {productCount}
                </Badge>
              )}

              {/* View Button */}
              <Link href={`/categories/${category.slug || category.id}`}>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 px-2 text-xs text-gray-400 hover:text-white hover:bg-white/10"
                  onClick={() => onCategorySelect?.(category)}
                >
                  <Eye className="w-3 h-3 mr-1" />
                  View
                </Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Children */}
      {hasChildren && isExpanded && (
        <div className="ml-4 space-y-1">
          {category.children?.map((child) => (
            <CategoryNode
              key={child.id}
              category={child}
              level={level + 1}
              isExpanded={isExpanded}
              onToggle={onToggle}
              showProductCount={showProductCount}
              onCategorySelect={onCategorySelect}
              selectedCategoryId={selectedCategoryId}
              maxDepth={maxDepth}
            />
          ))}
        </div>
      )}
    </div>
  )
}

export function CategoryTree({
  categories,
  className,
  showProductCount = true,
  expandAll = false,
  onCategorySelect,
  selectedCategoryId,
  maxDepth = 5
}: CategoryTreeProps) {
  const [expandedCategories, setExpandedCategories] = useState<Set<string>>(
    expandAll ? new Set(getAllCategoryIds(categories)) : new Set()
  )

  const toggleCategory = (categoryId: string) => {
    const newExpanded = new Set(expandedCategories)
    if (newExpanded.has(categoryId)) {
      newExpanded.delete(categoryId)
    } else {
      newExpanded.add(categoryId)
    }
    setExpandedCategories(newExpanded)
  }

  const expandAllCategories = () => {
    setExpandedCategories(new Set(getAllCategoryIds(categories)))
  }

  const collapseAllCategories = () => {
    setExpandedCategories(new Set())
  }

  return (
    <div className={cn('space-y-2', className)}>
      {/* Tree Controls */}
      <div className="flex items-center justify-end mb-4">
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={expandAllCategories}
            className="text-xs text-gray-400 hover:text-white"
          >
            <Plus className="w-3 h-3 mr-1" />
            Expand All
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={collapseAllCategories}
            className="text-xs text-gray-400 hover:text-white"
          >
            <Minus className="w-3 h-3 mr-1" />
            Collapse All
          </Button>
        </div>
      </div>

      {/* Tree Nodes */}
      <div className="space-y-1">
        {categories.map((category) => (
          <CategoryNode
            key={category.id}
            category={category}
            level={0}
            isExpanded={expandedCategories.has(category.id)}
            onToggle={toggleCategory}
            showProductCount={showProductCount}
            onCategorySelect={onCategorySelect}
            selectedCategoryId={selectedCategoryId}
            maxDepth={maxDepth}
          />
        ))}
      </div>

      {/* Empty State */}
      {categories.length === 0 && (
        <div className="text-center py-8">
          <FolderTree className="w-12 h-12 text-gray-400 mx-auto mb-3" />
          <p className="text-gray-400">No categories found</p>
        </div>
      )}
    </div>
  )
}

// Helper function to get all category IDs recursively
function getAllCategoryIds(categories: Category[]): string[] {
  const ids: string[] = []
  
  function collectIds(cats: Category[]) {
    cats.forEach(cat => {
      ids.push(cat.id)
      if (cat.children && cat.children.length > 0) {
        collectIds(cat.children)
      }
    })
  }
  
  collectIds(categories)
  return ids
}

// Breadcrumb component for category navigation
interface CategoryBreadcrumbProps {
  category: Category
  className?: string
}

export function CategoryBreadcrumb({ category, className }: CategoryBreadcrumbProps) {
  const breadcrumbs: Category[] = []
  
  // Build breadcrumb path (this would need to be enhanced with actual parent data)
  let current = category
  while (current) {
    breadcrumbs.unshift(current)
    current = current.parent as Category // Type assertion needed
  }

  return (
    <nav className={cn('flex items-center space-x-2 text-sm', className)}>
      {breadcrumbs.map((crumb, index) => (
        <React.Fragment key={crumb.id}>
          {index > 0 && (
            <ChevronRight className="w-4 h-4 text-gray-400" />
          )}
          <Link
            href={`/categories/${crumb.slug || crumb.id}`}
            className={cn(
              'transition-colors',
              index === breadcrumbs.length - 1
                ? 'text-[#ff9000] font-medium'
                : 'text-gray-400 hover:text-white'
            )}
          >
            {crumb.name}
          </Link>
        </React.Fragment>
      ))}
    </nav>
  )
}
