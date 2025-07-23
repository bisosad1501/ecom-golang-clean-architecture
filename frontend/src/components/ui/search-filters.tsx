'use client'

import React, { useState, useEffect, useRef } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Checkbox } from '@/components/ui/checkbox'
import { Slider } from '@/components/ui/slider'
import { 
  Search, 
  Filter, 
  X, 
  ChevronDown, 
  ChevronUp,
  SlidersHorizontal,
  Sparkles,
  Clock,
  TrendingUp,
  Package,
  Eye,
  Star
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { Category } from '@/types'

interface SearchFiltersProps {
  searchQuery: string
  onSearchChange: (query: string) => void
  categories: Category[]
  onFiltersChange: (filters: FilterState) => void
  className?: string
}

export interface FilterState {
  searchQuery: string
  sortBy: 'name' | 'product_count' | 'created_at' | 'updated_at'
  sortOrder: 'asc' | 'desc'
  showActiveOnly: boolean
  parentCategory?: string
  productCountRange: [number, number]
  hasSubcategories?: boolean
  recentlyUpdated?: boolean
}

const SORT_OPTIONS = [
  { value: 'name', label: 'Name', icon: '🔤' },
  { value: 'product_count', label: 'Product Count', icon: '📦' },
  { value: 'created_at', label: 'Date Created', icon: '📅' },
  { value: 'updated_at', label: 'Last Updated', icon: '🔄' }
]

const QUICK_FILTERS = [
  { key: 'trending', label: 'Trending', icon: TrendingUp, color: 'bg-red-500/20 text-red-400' },
  { key: 'popular', label: 'Popular', icon: Star, color: 'bg-yellow-500/20 text-yellow-400' },
  { key: 'recent', label: 'Recently Added', icon: Clock, color: 'bg-blue-500/20 text-blue-400' },
  { key: 'featured', label: 'Featured', icon: Sparkles, color: 'bg-purple-500/20 text-purple-400' }
]

export function SearchFilters({
  searchQuery,
  onSearchChange,
  categories,
  onFiltersChange,
  className
}: SearchFiltersProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  const [filters, setFilters] = useState<FilterState>({
    searchQuery,
    sortBy: 'name',
    sortOrder: 'asc',
    showActiveOnly: true,
    productCountRange: [0, 1000]
  })
  const [suggestions, setSuggestions] = useState<Category[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [activeQuickFilters, setActiveQuickFilters] = useState<string[]>([])
  const searchRef = useRef<HTMLInputElement>(null)
  const suggestionsRef = useRef<HTMLDivElement>(null)

  // Update filters when props change
  useEffect(() => {
    setFilters(prev => ({ ...prev, searchQuery }))
  }, [searchQuery])

  // Generate search suggestions
  useEffect(() => {
    if (searchQuery.length > 0) {
      const filtered = categories
        .filter(cat => 
          cat.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          cat.description?.toLowerCase().includes(searchQuery.toLowerCase())
        )
        .slice(0, 5)
      setSuggestions(filtered)
      setShowSuggestions(true)
    } else {
      setSuggestions([])
      setShowSuggestions(false)
    }
  }, [searchQuery, categories])

  // Handle filter changes
  const updateFilters = (newFilters: Partial<FilterState>) => {
    const updated = { ...filters, ...newFilters }
    setFilters(updated)
    onFiltersChange(updated)
  }

  // Handle search input
  const handleSearchChange = (value: string) => {
    onSearchChange(value)
    updateFilters({ searchQuery: value })
  }

  // Handle quick filter toggle
  const toggleQuickFilter = (filterKey: string) => {
    const newActiveFilters = activeQuickFilters.includes(filterKey)
      ? activeQuickFilters.filter(f => f !== filterKey)
      : [...activeQuickFilters, filterKey]
    
    setActiveQuickFilters(newActiveFilters)
    
    // Apply quick filter logic
    const filterUpdates: Partial<FilterState> = {}
    
    if (newActiveFilters.includes('trending')) {
      filterUpdates.sortBy = 'updated_at'
      filterUpdates.sortOrder = 'desc'
    }
    if (newActiveFilters.includes('popular')) {
      filterUpdates.sortBy = 'product_count'
      filterUpdates.sortOrder = 'desc'
    }
    if (newActiveFilters.includes('recent')) {
      filterUpdates.sortBy = 'created_at'
      filterUpdates.sortOrder = 'desc'
      filterUpdates.recentlyUpdated = true
    }
    
    updateFilters(filterUpdates)
  }

  // Clear all filters
  const clearFilters = () => {
    const defaultFilters: FilterState = {
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
      {/* Search Bar */}
      <div className="relative">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
          <Input
            ref={searchRef}
            type="text"
            placeholder="Search categories..."
            value={searchQuery}
            onChange={(e) => handleSearchChange(e.target.value)}
            className="pl-10 pr-4 bg-white/[0.06] border-white/10 text-white placeholder-gray-400 focus:border-[#ff9000]/50 focus:ring-[#ff9000]/20"
          />
          {searchQuery && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => handleSearchChange('')}
              className="absolute right-2 top-1/2 transform -translate-y-1/2 h-6 w-6 p-0 text-gray-400 hover:text-white"
            >
              <X className="h-3 w-3" />
            </Button>
          )}
        </div>

        {/* Search Suggestions */}
        {showSuggestions && suggestions.length > 0 && (
          <Card 
            ref={suggestionsRef}
            className="absolute top-full left-0 right-0 mt-1 z-50 bg-gray-900/95 backdrop-blur-md border-gray-700 shadow-xl"
          >
            <CardContent className="p-2">
              {suggestions.map((category) => (
                <Button
                  key={category.id}
                  variant="ghost"
                  className="w-full justify-start text-left p-2 h-auto hover:bg-white/10"
                  onClick={() => {
                    handleSearchChange(category.name)
                    setShowSuggestions(false)
                  }}
                >
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 bg-[#ff9000]/20 rounded-lg flex items-center justify-center">
                      <Package className="w-4 h-4 text-[#ff9000]" />
                    </div>
                    <div>
                      <div className="font-medium text-white">{category.name}</div>
                      {category.description && (
                        <div className="text-xs text-gray-400 line-clamp-1">
                          {category.description}
                        </div>
                      )}
                    </div>
                  </div>
                </Button>
              ))}
            </CardContent>
          </Card>
        )}
      </div>

      {/* Quick Filters */}
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
                'h-8 px-3 text-xs font-medium transition-all duration-300 rounded-full border',
                isActive
                  ? `${filter.color} border-current shadow-lg scale-105`
                  : 'bg-white/[0.08] text-gray-300 border-white/20 hover:bg-white/[0.15] hover:text-white hover:border-white/40 hover:scale-105 hover:shadow-md backdrop-blur-sm'
              )}
            >
              <Icon className="w-3 h-3 mr-1.5" />
              {filter.label}
            </Button>
          )
        })}
      </div>

      {/* Advanced Filters Toggle */}
      <div className="flex items-center justify-between">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setIsExpanded(!isExpanded)}
          className="text-gray-300 hover:text-white hover:bg-white/[0.12] px-3 py-2 rounded-lg border border-white/20 hover:border-white/40 transition-all duration-300 hover:scale-105 hover:shadow-md backdrop-blur-sm"
        >
          <SlidersHorizontal className="w-4 h-4 mr-2" />
          Advanced Filters
          {activeFilterCount > 0 && (
            <Badge variant="secondary" className="ml-2 bg-[#ff9000] text-white text-xs">
              {activeFilterCount}
            </Badge>
          )}
          {isExpanded ? (
            <ChevronUp className="w-4 h-4 ml-2" />
          ) : (
            <ChevronDown className="w-4 h-4 ml-2" />
          )}
        </Button>

        {activeFilterCount > 0 && (
          <Button
            variant="ghost"
            size="sm"
            onClick={clearFilters}
            className="text-gray-300 hover:text-white hover:bg-red-500/20 hover:border-red-500/40 px-3 py-2 rounded-lg border border-white/20 transition-all duration-300 hover:scale-105 hover:shadow-md backdrop-blur-sm text-xs"
          >
            Clear all
          </Button>
        )}
      </div>

      {/* Advanced Filters Panel */}
      {isExpanded && (
        <Card className="bg-white/[0.06] backdrop-blur-md border-white/10">
          <CardContent className="p-4 space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Sort Options */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-gray-300">Sort By</label>
                <Select
                  value={filters.sortBy}
                  onValueChange={(value: any) => updateFilters({ sortBy: value })}
                >
                  <SelectTrigger className="bg-white/[0.06] border-white/10 text-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700">
                    {SORT_OPTIONS.map((option) => (
                      <SelectItem key={option.value} value={option.value}>
                        <span className="flex items-center gap-2">
                          <span>{option.icon}</span>
                          {option.label}
                        </span>
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>

              {/* Sort Order */}
              <div className="space-y-2">
                <label className="text-sm font-medium text-gray-300">Order</label>
                <Select
                  value={filters.sortOrder}
                  onValueChange={(value: any) => updateFilters({ sortOrder: value })}
                >
                  <SelectTrigger className="bg-white/[0.06] border-white/10 text-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700">
                    <SelectItem value="asc">Ascending</SelectItem>
                    <SelectItem value="desc">Descending</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Checkboxes */}
            <div className="space-y-3">
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="active-only"
                  checked={filters.showActiveOnly}
                  onCheckedChange={(checked) => updateFilters({ showActiveOnly: !!checked })}
                />
                <label htmlFor="active-only" className="text-sm text-gray-300">
                  Show active categories only
                </label>
              </div>

              <div className="flex items-center space-x-2">
                <Checkbox
                  id="has-subcategories"
                  checked={filters.hasSubcategories}
                  onCheckedChange={(checked) => updateFilters({ hasSubcategories: checked ? true : undefined })}
                />
                <label htmlFor="has-subcategories" className="text-sm text-gray-300">
                  Has subcategories
                </label>
              </div>

              <div className="flex items-center space-x-2">
                <Checkbox
                  id="recently-updated"
                  checked={filters.recentlyUpdated}
                  onCheckedChange={(checked) => updateFilters({ recentlyUpdated: checked ? true : undefined })}
                />
                <label htmlFor="recently-updated" className="text-sm text-gray-300">
                  Recently updated (last 7 days)
                </label>
              </div>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
