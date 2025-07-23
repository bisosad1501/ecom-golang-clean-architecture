'use client'

import React from 'react'
import Link from 'next/link'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  ChevronRight, 
  Home, 
  FolderTree,
  ArrowUp,
  ArrowDown,
  Package,
  Users
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Category } from '@/types'

interface CategoryPathIndicatorProps {
  category: Category
  className?: string
  showStats?: boolean
  showActions?: boolean
}

export function CategoryPathIndicator({
  category,
  className,
  showStats = true,
  showActions = true
}: CategoryPathIndicatorProps) {
  // Build breadcrumb path
  const buildPath = (cat: Category): Category[] => {
    const path: Category[] = []
    let current = cat
    
    // Add current category
    path.unshift(current)
    
    // Add parent if exists (simplified - in real app you'd fetch full parent data)
    if (current.parent) {
      path.unshift(current.parent)
    }
    
    return path
  }

  const breadcrumbPath = buildPath(category)
  const hasChildren = category.children && category.children.length > 0
  const hasParent = category.parent_id !== null
  const productCount = category.product_count || 0
  const childrenCount = category.children?.length || 0

  return (
    <div className={cn('space-y-4', className)}>
      {/* Breadcrumb Navigation */}
      <nav className="flex items-center space-x-2 text-sm">
        <Link 
          href="/" 
          className="text-gray-400 hover:text-white transition-colors flex items-center"
        >
          <Home className="w-4 h-4" />
        </Link>
        
        <ChevronRight className="w-4 h-4 text-gray-400" />
        
        <Link 
          href="/categories" 
          className="text-gray-400 hover:text-white transition-colors"
        >
          Categories
        </Link>

        {breadcrumbPath.map((crumb, index) => (
          <React.Fragment key={crumb.id}>
            <ChevronRight className="w-4 h-4 text-gray-400" />
            <Link
              href={`/categories/${crumb.slug || crumb.id}`}
              className={cn(
                'transition-colors flex items-center gap-1',
                index === breadcrumbPath.length - 1
                  ? 'text-[#ff9000] font-medium'
                  : 'text-gray-400 hover:text-white'
              )}
            >
              <FolderTree className="w-3 h-3" />
              {crumb.name}
            </Link>
          </React.Fragment>
        ))}
      </nav>

      {/* Category Info Card */}
      <div className="bg-white/[0.06] backdrop-blur-sm border border-white/10 rounded-xl p-4">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            {/* Category Header */}
            <div className="flex items-center gap-3 mb-3">
              <div className="w-12 h-12 bg-[#ff9000]/20 rounded-xl flex items-center justify-center">
                <FolderTree className="w-6 h-6 text-[#ff9000]" />
              </div>
              
              <div>
                <h1 className="text-2xl font-bold text-white">{category.name}</h1>
                {category.description && (
                  <p className="text-gray-400 mt-1">{category.description}</p>
                )}
              </div>
            </div>

            {/* Stats */}
            {showStats && (
              <div className="flex items-center gap-4 mb-4">
                <Badge variant="secondary" className="bg-[#ff9000]/20 text-[#ff9000] border-[#ff9000]/30">
                  <Package className="w-3 h-3 mr-1" />
                  {productCount} products
                </Badge>
                
                {hasChildren && (
                  <Badge variant="secondary" className="bg-blue-500/20 text-blue-400 border-blue-500/30">
                    <FolderTree className="w-3 h-3 mr-1" />
                    {childrenCount} subcategories
                  </Badge>
                )}
                
                <Badge variant="outline" className="border-white/20 text-white/70">
                  Level {category.level || 0}
                </Badge>
              </div>
            )}
          </div>

          {/* Actions */}
          {showActions && (
            <div className="flex flex-col gap-2">
              {hasParent && (
                <Link href={`/categories/${category.parent?.slug || category.parent_id}`}>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-gray-400 hover:text-white hover:bg-white/10"
                  >
                    <ArrowUp className="w-4 h-4 mr-1" />
                    Parent
                  </Button>
                </Link>
              )}
              
              {hasChildren && (
                <Button
                  variant="ghost"
                  size="sm"
                  className="text-gray-400 hover:text-white hover:bg-white/10"
                  onClick={() => {
                    // Scroll to children section
                    document.getElementById('subcategories')?.scrollIntoView({ 
                      behavior: 'smooth' 
                    })
                  }}
                >
                  <ArrowDown className="w-4 h-4 mr-1" />
                  Children
                </Button>
              )}
            </div>
          )}
        </div>

        {/* Hierarchy Visualization */}
        <div className="mt-4 pt-4 border-t border-white/10">
          <div className="flex items-center justify-center space-x-4 text-sm">
            {/* Parent */}
            {hasParent && category.parent && (
              <Link 
                href={`/categories/${category.parent.slug || category.parent_id}`}
                className="flex items-center gap-2 text-gray-400 hover:text-white transition-colors"
              >
                <div className="w-8 h-8 bg-white/10 rounded-lg flex items-center justify-center">
                  <FolderTree className="w-4 h-4" />
                </div>
                <span className="truncate max-w-24">{category.parent.name}</span>
              </Link>
            )}

            {hasParent && (
              <ChevronRight className="w-4 h-4 text-gray-400" />
            )}

            {/* Current Category */}
            <div className="flex items-center gap-2 text-[#ff9000]">
              <div className="w-10 h-10 bg-[#ff9000]/20 rounded-lg flex items-center justify-center border-2 border-[#ff9000]/50">
                <FolderTree className="w-5 h-5" />
              </div>
              <span className="font-medium truncate max-w-32">{category.name}</span>
            </div>

            {hasChildren && (
              <>
                <ChevronRight className="w-4 h-4 text-gray-400" />
                
                {/* Children Preview */}
                <div className="flex items-center gap-1">
                  {category.children?.slice(0, 3).map((child, index) => (
                    <Link
                      key={child.id}
                      href={`/categories/${child.slug || child.id}`}
                      className="flex items-center gap-1 text-gray-400 hover:text-white transition-colors"
                    >
                      <div className="w-6 h-6 bg-white/10 rounded flex items-center justify-center">
                        <FolderTree className="w-3 h-3" />
                      </div>
                      {index === 0 && (
                        <span className="text-xs truncate max-w-16">{child.name}</span>
                      )}
                    </Link>
                  ))}
                  
                  {childrenCount > 3 && (
                    <Badge variant="outline" className="border-white/20 text-white/60 text-xs">
                      +{childrenCount - 3}
                    </Badge>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Quick Navigation */}
      {(hasParent || hasChildren) && (
        <div className="flex items-center justify-between bg-white/[0.04] backdrop-blur-sm border border-white/10 rounded-lg p-3">
          <div className="text-sm text-gray-400">Quick Navigation</div>
          
          <div className="flex items-center gap-2">
            {hasParent && (
              <Link href={`/categories/${category.parent?.slug || category.parent_id}`}>
                <Button variant="ghost" size="sm" className="text-xs">
                  ← Parent Category
                </Button>
              </Link>
            )}
            
            {hasChildren && (
              <Button 
                variant="ghost" 
                size="sm" 
                className="text-xs"
                onClick={() => {
                  document.getElementById('subcategories')?.scrollIntoView({ 
                    behavior: 'smooth' 
                  })
                }}
              >
                View Subcategories →
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  )
}

// Compact version for use in cards
interface CompactCategoryPathProps {
  category: Category
  className?: string
}

export function CompactCategoryPath({ category, className }: CompactCategoryPathProps) {
  const hasParent = category.parent_id !== null
  const hasChildren = category.children && category.children.length > 0
  const level = category.level || 0

  return (
    <div className={cn('flex items-center gap-2 text-xs', className)}>
      {/* Level indicator */}
      <Badge 
        variant="outline" 
        className={cn(
          'px-1.5 py-0.5 text-xs border-current',
          level === 0 ? 'text-[#ff9000] border-[#ff9000]/30' : 'text-white/60 border-white/20'
        )}
      >
        L{level}
      </Badge>

      {/* Hierarchy indicators */}
      <div className="flex items-center gap-1">
        {hasParent && (
          <div className="w-2 h-2 rounded-full bg-blue-400/60" title="Has parent" />
        )}
        
        <div className="w-2 h-2 rounded-full bg-[#ff9000]" title="Current category" />
        
        {hasChildren && (
          <div className="w-2 h-2 rounded-full bg-green-400/60" title="Has children" />
        )}
      </div>

      {/* Path text */}
      <span className="text-gray-400 truncate">
        {hasParent && '... / '}
        <span className="text-white">{category.name}</span>
        {hasChildren && ` / ${category.children?.length} items`}
      </span>
    </div>
  )
}
