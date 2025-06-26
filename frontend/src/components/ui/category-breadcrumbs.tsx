'use client'

import { useEffect, useState } from 'react'
import { ChevronRight, Home } from 'lucide-react'
import Link from 'next/link'
import { Category } from '@/types'
import { categoryService } from '@/lib/services/categories'
import { cn } from '@/lib/utils'

interface CategoryBreadcrumbsProps {
  categoryId?: string
  className?: string
  showHome?: boolean
  onCategoryClick?: (categoryId: string) => void
}

export function CategoryBreadcrumbs({ 
  categoryId, 
  className,
  showHome = true,
  onCategoryClick
}: CategoryBreadcrumbsProps) {
  const [categoryPath, setCategoryPath] = useState<Category[]>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!categoryId) {
      setCategoryPath([])
      return
    }

    const loadCategoryPath = async () => {
      setIsLoading(true)
      try {
        const path = await categoryService.getCategoryPath(categoryId)
        setCategoryPath(path)
      } catch (error) {
        console.error('Failed to load category path:', error)
        setCategoryPath([])
      } finally {
        setIsLoading(false)
      }
    }

    loadCategoryPath()
  }, [categoryId])

  if (isLoading) {
    return (
      <div className={cn('flex items-center space-x-2', className)}>
        <div className="h-4 w-12 bg-gray-200 rounded animate-pulse" />
        <ChevronRight className="h-4 w-4 text-gray-400" />
        <div className="h-4 w-20 bg-gray-200 rounded animate-pulse" />
        <ChevronRight className="h-4 w-4 text-gray-400" />
        <div className="h-4 w-16 bg-gray-200 rounded animate-pulse" />
      </div>
    )
  }

  if (!categoryId || categoryPath.length === 0) {
    return null
  }

  return (
    <nav className={cn('flex items-center space-x-1 text-sm', className)}>
      {showHome && (
        <>
          <Link 
            href="/products" 
            className="flex items-center text-gray-500 hover:text-gray-700 transition-colors"
          >
            <Home className="h-4 w-4" />
            <span className="ml-1">Products</span>
          </Link>
          <ChevronRight className="h-4 w-4 text-gray-400" />
        </>
      )}
      
      {categoryPath.map((category, index) => {
        const isLast = index === categoryPath.length - 1
        
        return (
          <div key={category.id} className="flex items-center space-x-1">
            {onCategoryClick ? (
              <button
                onClick={() => onCategoryClick(category.id)}
                className={cn(
                  'transition-colors',
                  isLast 
                    ? 'text-gray-900 font-medium cursor-default' 
                    : 'text-primary-600 hover:text-primary-800 hover:underline'
                )}
                disabled={isLast}
              >
                {category.name}
              </button>
            ) : (
              <Link
                href={`/products?category_id=${category.id}`}
                className={cn(
                  'transition-colors',
                  isLast 
                    ? 'text-gray-900 font-medium' 
                    : 'text-primary-600 hover:text-primary-800 hover:underline'
                )}
              >
                {category.name}
              </Link>
            )}
            
            {!isLast && (
              <ChevronRight className="h-4 w-4 text-gray-400" />
            )}
          </div>
        )
      })}
    </nav>
  )
}
