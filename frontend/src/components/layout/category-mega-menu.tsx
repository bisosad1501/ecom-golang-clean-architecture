'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { ChevronDown, ChevronRight, Grid3X3, Package, TrendingUp, Star, ArrowRight, Tag } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'
import { cn } from '@/lib/utils'

interface CategoryMegaMenuProps {
  className?: string
}

export function CategoryMegaMenu({ className }: CategoryMegaMenuProps) {
  const { data: categories, isLoading } = useCategories()
  const [isOpen, setIsOpen] = useState(false)
  const [hoveredCategory, setHoveredCategory] = useState<string | null>(null)
  const [activeTimeout, setActiveTimeout] = useState<NodeJS.Timeout | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)

  const rootCategories = categories?.filter(cat => !cat.parent_id) || []

  // Close menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  const handleCategoryHover = (categoryId: string) => {
    if (activeTimeout) {
      clearTimeout(activeTimeout)
    }
    
    const timeout = setTimeout(() => {
      setHoveredCategory(categoryId)
    }, 100)
    
    setActiveTimeout(timeout)
  }

  const handleCategoryLeave = () => {
    if (activeTimeout) {
      clearTimeout(activeTimeout)
    }
    
    const timeout = setTimeout(() => {
      setHoveredCategory(null)
    }, 150)
    
    setActiveTimeout(timeout)
  }

  const getSubcategories = (parentId: string): Category[] => {
    return categories?.filter(cat => cat.parent_id === parentId) || []
  }

  const getFeaturedCategories = (): Category[] => {
    // Return categories with highest product counts or featured ones
    return rootCategories
      .filter(cat => cat.product_count && cat.product_count > 0)
      .sort((a, b) => (b.product_count || 0) - (a.product_count || 0))
      .slice(0, 4)
  }

  if (isLoading) {
    return (
      <div className={className}>
        <Button variant="ghost" className="h-12 px-4" disabled>
          <Grid3X3 className="h-5 w-5 mr-2" />
          Categories
          <ChevronDown className="h-4 w-4 ml-2" />
        </Button>
      </div>
    )
  }

  return (
    <div className={cn("relative", className)} ref={menuRef}>
      <Button
        variant="ghost"
        className={cn(
          "h-12 px-4 hover:bg-primary-50 transition-all duration-200",
          isOpen && "bg-primary-50 text-primary-700"
        )}
        onMouseEnter={() => setIsOpen(true)}
        onClick={() => setIsOpen(!isOpen)}
      >
        <Grid3X3 className="h-5 w-5 mr-2" />
        <span className="font-semibold">All Categories</span>
        <ChevronDown className={cn(
          "h-4 w-4 ml-2 transition-transform duration-200",
          isOpen && "rotate-180"
        )} />
      </Button>

      {/* Mega Menu Dropdown */}
      {isOpen && (
        <div 
          className="absolute top-full left-0 w-screen max-w-6xl bg-background border border-border rounded-lg shadow-2xl z-50 overflow-hidden"
          onMouseLeave={() => {
            setIsOpen(false)
            setHoveredCategory(null)
          }}
        >
          <div className="flex h-[500px]">
            {/* Left Column - Main Categories */}
            <div className="w-64 bg-gray-50 border-r border-border overflow-y-auto">
              <div className="p-4">
                <h3 className="font-semibold text-sm text-gray-900 mb-3 flex items-center">
                  <Package className="h-4 w-4 mr-2" />
                  Browse Categories
                </h3>
                <div className="space-y-1">
                  {rootCategories.map((category) => {
                    const subcategories = getSubcategories(category.id)
                    const isHovered = hoveredCategory === category.id
                    
                    return (
                      <div key={category.id}>
                        <Link
                          href={`/categories/${category.id}`}
                          className={cn(
                            "flex items-center justify-between p-3 rounded-lg transition-all duration-200 group",
                            "hover:bg-white hover:shadow-md",
                            isHovered && "bg-white shadow-md text-primary-700"
                          )}
                          onMouseEnter={() => handleCategoryHover(category.id)}
                          onMouseLeave={handleCategoryLeave}
                        >
                          <div className="flex items-center">
                            <Tag className="h-4 w-4 mr-3 text-gray-500 group-hover:text-primary-600" />
                            <div>
                              <div className="font-medium text-sm">{category.name}</div>
                              {category.product_count !== undefined && (
                                <div className="text-xs text-gray-500">
                                  {category.product_count} products
                                </div>
                              )}
                            </div>
                          </div>
                          {subcategories.length > 0 && (
                            <ChevronRight className="h-4 w-4 text-gray-400 group-hover:text-primary-600" />
                          )}
                        </Link>
                      </div>
                    )
                  })}
                </div>
              </div>
            </div>

            {/* Right Column - Subcategories & Featured */}
            <div className="flex-1 p-6">
              {hoveredCategory ? (
                <div>
                  {/* Subcategories */}
                  <div className="mb-6">
                    <h4 className="font-semibold text-gray-900 mb-4 flex items-center">
                      <ChevronRight className="h-4 w-4 mr-2" />
                      {rootCategories.find(c => c.id === hoveredCategory)?.name} Subcategories
                    </h4>
                    <div className="grid grid-cols-2 gap-4">
                      {getSubcategories(hoveredCategory).map((subcategory) => (
                        <Link
                          key={subcategory.id}
                          href={`/categories/${subcategory.id}`}
                          className="flex items-center p-3 rounded-lg hover:bg-gray-50 transition-colors group"
                        >
                          <div className="flex-1">
                            <div className="font-medium text-sm group-hover:text-primary-600">
                              {subcategory.name}
                            </div>
                            {subcategory.description && (
                              <div className="text-xs text-gray-500 mt-1">
                                {subcategory.description}
                              </div>
                            )}
                            {subcategory.product_count !== undefined && (
                              <Badge variant="secondary" className="mt-2 text-xs">
                                {subcategory.product_count} items
                              </Badge>
                            )}
                          </div>
                          <ArrowRight className="h-4 w-4 text-gray-400 group-hover:text-primary-600 opacity-0 group-hover:opacity-100 transition-all duration-200" />
                        </Link>
                      ))}
                    </div>
                  </div>

                  {/* View All Link */}
                  <div className="pt-4 border-t border-border">
                    <Link
                      href={`/categories/${hoveredCategory}`}
                      className="inline-flex items-center text-primary-600 hover:text-primary-700 font-medium transition-colors group"
                    >
                      View All {rootCategories.find(c => c.id === hoveredCategory)?.name}
                      <ArrowRight className="h-4 w-4 ml-2 group-hover:translate-x-1 transition-transform duration-200" />
                    </Link>
                  </div>
                </div>
              ) : (
                <div>
                  {/* Featured Categories */}
                  <div className="mb-6">
                    <h4 className="font-semibold text-gray-900 mb-4 flex items-center">
                      <Star className="h-4 w-4 mr-2 text-yellow-500" />
                      Featured Categories
                    </h4>
                    <div className="grid grid-cols-2 gap-4">
                      {getFeaturedCategories().map((category) => (
                        <Link
                          key={category.id}
                          href={`/categories/${category.id}`}
                          className="group"
                        >
                          <div className="bg-gradient-to-br from-primary-50 to-purple-50 p-4 rounded-xl border border-primary-100 hover:border-primary-200 transition-all duration-200 hover:shadow-md">
                            <div className="flex items-center mb-2">
                              <div className="h-8 w-8 rounded-lg bg-primary-100 flex items-center justify-center mr-3">
                                <Tag className="h-4 w-4 text-primary-600" />
                              </div>
                              <div>
                                <div className="font-semibold text-sm group-hover:text-primary-700">
                                  {category.name}
                                </div>
                                <div className="text-xs text-gray-600">
                                  {category.product_count} products
                                </div>
                              </div>
                            </div>
                            {category.description && (
                              <p className="text-xs text-gray-600 line-clamp-2">
                                {category.description}
                              </p>
                            )}
                          </div>
                        </Link>
                      ))}
                    </div>
                  </div>

                  {/* Trending */}
                  <div className="pt-4 border-t border-border">
                    <h4 className="font-semibold text-gray-900 mb-3 flex items-center">
                      <TrendingUp className="h-4 w-4 mr-2 text-emerald-500" />
                      Trending Now
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {['Electronics', 'Fashion', 'Home & Garden', 'Sports', 'Books', 'Toys'].map((trend) => (
                        <Badge key={trend} variant="outline" className="hover:bg-primary-50 cursor-pointer">
                          {trend}
                        </Badge>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
