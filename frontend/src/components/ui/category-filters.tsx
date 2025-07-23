'use client'

import { useState, useEffect } from 'react'
import { Search, X, Filter, Package, TreePine, Layers, Leaf, SlidersHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Slider } from '@/components/ui/slider'
import { Category } from '@/types'
import { cn } from '@/lib/utils'

interface CategoryFiltersProps {
  filters: CategoryFilterState
  onFiltersChange: (filters: CategoryFilterState) => void
  className?: string
}

export interface CategoryFilterState {
  searchQuery: string
  sortBy: 'name' | 'product_count' | 'created_at' | 'updated_at'
  sortOrder: 'asc' | 'desc'
  showActiveOnly: boolean
  parentCategory?: string
  hasSubcategories?: boolean
  recentlyUpdated?: boolean
  productCountRange: [number, number]
}

const QUICK_FILTERS = [
  {
    key: 'root',
    label: 'Root Only',
    icon: TreePine,
    color: 'bg-blue-500/20 text-blue-400 border-blue-400/40'
  },
  {
    key: 'parent',
    label: 'With Children',
    icon: Layers,
    color: 'bg-green-500/20 text-green-400 border-green-400/40'
  },
  {
    key: 'leaf',
    label: 'Leaf Only',
    icon: Leaf,
    color: 'bg-purple-500/20 text-purple-400 border-purple-400/40'
  },
  {
    key: 'recent',
    label: 'Recently Updated',
    icon: Package,
    color: 'bg-[#ff9000]/20 text-[#ff9000] border-[#ff9000]/40'
  }
]

export function CategoryFilters({
  filters,
  onFiltersChange,
  className
}: CategoryFiltersProps) {
  const [activeQuickFilters, setActiveQuickFilters] = useState<string[]>([])

  // Handle filter changes
  const updateFilters = (newFilters: Partial<CategoryFilterState>) => {
    const updated = { ...filters, ...newFilters }
    onFiltersChange(updated)
  }

  // Handle search input
  const handleSearchChange = (value: string) => {
    updateFilters({ searchQuery: value })
  }

  // Handle quick filter toggle
  const toggleQuickFilter = (filterKey: string) => {
    const newActiveFilters = activeQuickFilters.includes(filterKey)
      ? activeQuickFilters.filter(f => f !== filterKey)
      : [...activeQuickFilters, filterKey]
    
    setActiveQuickFilters(newActiveFilters)
    
    // Apply quick filter logic
    const filterUpdates: Partial<CategoryFilterState> = {}
    
    if (newActiveFilters.includes('root')) {
      filterUpdates.parentCategory = 'root'
    }
    if (newActiveFilters.includes('parent')) {
      filterUpdates.hasSubcategories = true
    }
    if (newActiveFilters.includes('leaf')) {
      filterUpdates.hasSubcategories = false
    }
    if (newActiveFilters.includes('recent')) {
      filterUpdates.recentlyUpdated = true
      filterUpdates.sortBy = 'updated_at'
      filterUpdates.sortOrder = 'desc'
    }
    
    updateFilters(filterUpdates)
  }

  // Clear all filters
  const clearFilters = () => {
    const defaultFilters: CategoryFilterState = {
      searchQuery: '',
      sortBy: 'name',
      sortOrder: 'asc',
      showActiveOnly: true,
      productCountRange: [0, 1000]
    }
    setFilters(defaultFilters)
    setActiveQuickFilters([])
    onSearchChange('')
    onFiltersChange(defaultFilters)
  }

  // Get active filter count
  const activeFilterCount = [
    filters.parentCategory,
    filters.hasSubcategories,
    filters.recentlyUpdated,
    ...activeQuickFilters
  ].filter(Boolean).length

  return (
    <div className={cn('space-y-4', className)}>
      {/* Search */}
      <div className="space-y-3">
        <div className="flex items-center gap-2">
          <Search className="h-3.5 w-3.5 text-[#ff9000]" />
          <h3 className="text-white font-medium text-sm">Search</h3>
        </div>

        <div className="relative">
          <div className="absolute left-3 top-1/2 transform -translate-y-1/2 z-10">
            <Search className="h-3 w-3 text-gray-400" />
          </div>
          <Input
            placeholder="Search categories..."
            value={filters.searchQuery}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="pl-9 text-sm bg-white/[0.05] border-white/15 text-white placeholder-gray-400 focus:bg-white/[0.08] focus:border-[#ff9000]/50 transition-all duration-300 rounded-lg h-8"
          />
          {filters.searchQuery && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleSearchChange('')}
              className="absolute right-1 top-1/2 transform -translate-y-1/2 h-6 w-6 p-0 text-gray-400 hover:text-white"
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </div>
      </div>

      {/* Quick Filters */}
      <div className="space-y-3">
        <div className="flex items-center gap-2">
          <Filter className="h-3.5 w-3.5 text-[#ff9000]" />
          <h3 className="text-white font-medium text-sm">Quick Filters</h3>
        </div>
        
        <div className="flex flex-wrap gap-2">
          {QUICK_FILTERS.map((filter) => {
            const Icon = filter.icon
            const isActive = activeQuickFilters.includes(filter.key)
            
            return (
              <Button
                key={filter.key}
                variant="ghost"
                size="sm"
                onClick={() => toggleQuickFilter(filter.key)}
                className={cn(
                  'h-8 px-3 text-xs font-medium transition-all duration-300 rounded-lg border',
                  isActive
                    ? `${filter.color} shadow-lg scale-105`
                    : 'bg-white/[0.05] text-gray-300 border-white/20 hover:bg-white/[0.08] hover:text-white hover:border-white/30 hover:scale-105 hover:shadow-md backdrop-blur-sm'
                )}
              >
                <Icon className="w-3 h-3 mr-1.5" />
                {filter.label}
              </Button>
            )
          })}
        </div>
      </div>

      {/* Sort & Filters */}
      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            <SlidersHorizontal className="h-3.5 w-3.5 text-[#ff9000]" />
            <h3 className="text-white font-medium text-sm">Sort & Filters</h3>
          </div>
          {activeFilterCount > 0 && (
            <div className="flex items-center gap-2">
              <Badge className="bg-[#ff9000] text-white text-xs px-1.5 py-0.5">
                {activeFilterCount}
              </Badge>
              <Button
                variant="ghost"
                size="sm"
                onClick={clearFilters}
                className="text-gray-400 hover:text-white hover:bg-red-500/20 text-xs h-6 px-2"
              >
                Clear
              </Button>
            </div>
          )}
        </div>

        <div className="space-y-3">
          <div className="grid grid-cols-2 gap-2">
            {/* Sort Options */}
            <div className="space-y-1">
              <label className="text-xs font-medium text-gray-400">Sort By</label>
              <Select
                value={filters.sortBy}
                onValueChange={(value: any) => updateFilters({ sortBy: value })}
              >
                <SelectTrigger className="bg-white/[0.05] border-white/15 text-white h-8 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent className="bg-gray-900 border-gray-700">
                  <SelectItem value="name">Name</SelectItem>
                  <SelectItem value="product_count">Products</SelectItem>
                  <SelectItem value="created_at">Created</SelectItem>
                  <SelectItem value="updated_at">Updated</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {/* Sort Order */}
            <div className="space-y-1">
              <label className="text-xs font-medium text-gray-400">Order</label>
              <Select
                value={filters.sortOrder}
                onValueChange={(value: any) => updateFilters({ sortOrder: value })}
              >
                <SelectTrigger className="bg-white/[0.05] border-white/15 text-white h-8 text-xs">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent className="bg-gray-900 border-gray-700">
                  <SelectItem value="asc">A-Z</SelectItem>
                  <SelectItem value="desc">Z-A</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Product Count Range */}
          <div className="space-y-2">
            <label className="text-xs font-medium text-gray-400">Product Count</label>
            <div className="px-1">
              <Slider
                value={filters.productCountRange}
                onValueChange={(value) => updateFilters({ productCountRange: value as [number, number] })}
                max={1000}
                min={0}
                step={10}
                className="w-full"
              />
              <div className="flex justify-between text-xs text-gray-500 mt-1">
                <span>{filters.productCountRange[0]}</span>
                <span>{filters.productCountRange[1]}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
