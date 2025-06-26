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

function CategoryTree({
  categories,
  parentId = null,
  expandedId,
  setExpandedId,
  activeId,
  setActiveId,
  level = 0
}: {
  categories: Category[];
  parentId?: string | null;
  expandedId: string | null;
  setExpandedId: (id: string | null) => void;
  activeId: string | null;
  setActiveId: (id: string | null) => void;
  level?: number;
}) {
  const children = categories.filter(cat => cat.parent_id === parentId)
  if (!children.length) return null
  return (
    <div className={level === 0 ? "space-y-1" : "ml-4 border-l border-gray-800 pl-3"}>
      {children.map(category => {
        const hasChildren = categories.some(cat => cat.parent_id === category.id)
        const isExpanded = expandedId === category.id
        const isActive = activeId === category.id
        return (
          <div key={category.id}>
            <button
              type="button"
              className={cn(
                "flex items-center justify-between w-full p-2 rounded-lg transition-all duration-200 group text-left",
                isActive ? "bg-gray-800 border border-orange-500/30 text-orange-500" : "border border-transparent text-gray-300 hover:bg-gray-800 hover:text-orange-500 hover:border-orange-500/30",
                hasChildren && "pr-2"
              )}
              onClick={() => {
                setActiveId(category.id)
                if (hasChildren) setExpandedId(isExpanded ? null : category.id)
              }}
              onMouseEnter={() => setActiveId(category.id)}
            >
              <div className="flex items-center">
                <Tag className={cn("h-4 w-4 mr-2", isActive ? "text-orange-500" : "text-gray-500 group-hover:text-orange-500")} />
                <span className={cn("font-medium text-sm", isActive ? "text-orange-500" : "")}>{category.name}</span>
                {category.product_count !== undefined && (
                  <span className="ml-2 text-xs text-gray-500">{category.product_count}</span>
                )}
              </div>
              {hasChildren && (
                <ChevronRight className={cn(
                  "h-4 w-4 transition-transform",
                  isExpanded ? "rotate-90 text-orange-500" : "text-gray-400 group-hover:text-orange-500"
                )} />
              )}
            </button>
            {hasChildren && isExpanded && (
              <CategoryTree
                categories={categories}
                parentId={category.id}
                expandedId={expandedId}
                setExpandedId={setExpandedId}
                activeId={activeId}
                setActiveId={setActiveId}
                level={level + 1}
              />
            )}
          </div>
        )
      })}
    </div>
  )
}

export function CategoryMegaMenu({ className }: CategoryMegaMenuProps) {
  const { data: categories, isLoading } = useCategories()
  const [isOpen, setIsOpen] = useState(false)
  const [expandedId, setExpandedId] = useState<string | null>(null)
  const [activeId, setActiveId] = useState<string | null>(null)
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
    }, 50) // Giảm delay từ 100ms xuống 50ms
    
    setActiveTimeout(timeout)
  }

  const handleCategoryLeave = () => {
    if (activeTimeout) {
      clearTimeout(activeTimeout)
    }
    
    const timeout = setTimeout(() => {
      setHoveredCategory(null)
    }, 300) // Tăng delay từ 150ms lên 300ms để có thời gian di chuột
    
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
        <Button variant="ghost" className="h-14 px-4 hover:bg-orange-500/10 text-gray-300 hover:text-orange-500" disabled>
          <Grid3X3 className="h-4 w-4 mr-2 text-gray-400" />
          <span className="font-medium">All Categories</span>
          <ChevronDown className="h-4 w-4 ml-2 text-gray-400" />
        </Button>
      </div>
    )
  }

  return (
    <div className={cn("relative", className)} ref={menuRef}>
      <Button
        variant="ghost"
        className={cn(
          "h-14 px-4 hover:bg-orange-500/10 transition-all duration-200 text-gray-300 hover:text-orange-500",
          isOpen && "bg-orange-500/10 text-orange-500"
        )}
        onMouseEnter={() => setIsOpen(true)}
        onClick={() => setIsOpen(!isOpen)}
      >
        <Grid3X3 className="h-4 w-4 mr-2" />
        <span className="font-medium">All Categories</span>
        <ChevronDown className={cn(
          "h-4 w-4 ml-2 transition-transform duration-200",
          isOpen && "rotate-180"
        )} />
      </Button>

      {/* Mega Menu Dropdown */}
      {isOpen && (
        <div 
          className="absolute top-full left-0 min-w-[16rem] w-auto bg-gray-900 border border-gray-700 rounded-xl shadow-xl z-50 overflow-visible mt-1"
          onMouseLeave={() => {
            setIsOpen(false)
            setActiveId(null)
          }}
          onMouseEnter={() => {
            if (activeTimeout) {
              clearTimeout(activeTimeout)
              setActiveTimeout(null)
            }
          }}
        >
          <div className="py-2 px-2">
            <CategoryTree
              categories={categories || []}
              expandedId={expandedId}
              setExpandedId={setExpandedId}
              activeId={activeId}
              setActiveId={setActiveId}
            />
          </div>
        </div>
      )}
    </div>
  )
}
