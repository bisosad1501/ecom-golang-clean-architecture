'use client'

import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { X, Star } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ProductsParams } from '@/lib/services/products'
import { PRICE_RANGES } from '@/constants'
import { cn } from '@/lib/utils'
import { EnhancedCategoryFilter } from './enhanced-category-filter'

interface ProductFiltersProps {
  currentParams: ProductsParams
}

export function ProductFilters({ currentParams }: ProductFiltersProps) {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [priceRange, setPriceRange] = useState({
    min: currentParams.min_price?.toString() || '',
    max: currentParams.max_price?.toString() || '',
  })

  const updateFilters = (newParams: Partial<ProductsParams>) => {
    const params = new URLSearchParams(searchParams.toString())
    
    // Update or remove parameters
    Object.entries(newParams).forEach(([key, value]) => {
      if (value === undefined || value === null || value === '') {
        params.delete(key)
      } else {
        params.set(key, value.toString())
      }
    })

    // Reset to page 1 when filters change
    params.set('page', '1')

    router.push(`/products?${params.toString()}`)
  }

  const clearAllFilters = () => {
    router.push('/products')
  }

  const applyPriceRange = () => {
    updateFilters({
      min_price: priceRange.min ? parseFloat(priceRange.min) : undefined,
      max_price: priceRange.max ? parseFloat(priceRange.max) : undefined,
    })
  }

  const hasActiveFilters = !!(
    currentParams.search ||
    currentParams.category_id ||
    currentParams.min_price ||
    currentParams.max_price ||
    currentParams.rating ||
    currentParams.in_stock ||
    currentParams.tags?.length
  )

  return (
    <div className="space-y-4">
      {/* Enhanced Category Filter */}
      <EnhancedCategoryFilter
        selectedCategoryId={currentParams.category_id}
        onCategoryChange={(categoryId) => updateFilters({ category_id: categoryId })}
        showProductCount={true}
        showSearch={true}
      />

      {/* Clear Filters */}
      {hasActiveFilters && (
        <Card className="bg-gray-800 border-gray-700">
          <CardContent className="p-3">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium text-white">Active Filters</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={clearAllFilters}
                className="text-orange-400 hover:text-orange-300 hover:bg-orange-500/10"
              >
                Clear All
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Price Range */}
      <Card className="bg-gray-800 border-gray-700">
        <CardHeader className="pb-3">
          <CardTitle className="text-base text-white">Price Range</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 pt-0">
          {/* Quick price ranges */}
          <div className="space-y-1">
            {PRICE_RANGES.map((range) => (
              <button
                key={`${range.min}-${range.max}`}
                onClick={() => updateFilters({
                  min_price: range.min ?? undefined,
                  max_price: range.max ?? undefined,
                })}
                className={cn(
                  'w-full text-left px-3 py-2 text-sm rounded-lg transition-all duration-200',
                  currentParams.min_price === range.min &&
                  currentParams.max_price === range.max
                    ? 'bg-orange-500 text-white shadow-md'
                    : 'text-gray-300 hover:bg-gray-700 hover:text-white'
                )}
              >
                {range.label}
              </button>
            ))}
          </div>

          {/* Custom price range */}
          <div className="pt-3 border-t border-gray-700">
            <p className="text-sm font-medium mb-2 text-white">Custom Range</p>
            <div className="flex items-center space-x-2">
              <Input
                type="number"
                placeholder="Min"
                value={priceRange.min}
                onChange={(e) => setPriceRange(prev => ({ ...prev, min: e.target.value }))}
                className="text-sm bg-gray-700 border-gray-600 text-white placeholder-gray-400"
              />
              <span className="text-gray-400">-</span>
              <Input
                type="number"
                placeholder="Max"
                value={priceRange.max}
                onChange={(e) => setPriceRange(prev => ({ ...prev, max: e.target.value }))}
                className="text-sm bg-gray-700 border-gray-600 text-white placeholder-gray-400"
              />
            </div>
            <Button
              onClick={applyPriceRange}
              size="sm"
              className="w-full mt-2 bg-orange-500 hover:bg-orange-600 text-white"
            >
              Apply
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Rating Filter */}
      <Card className="bg-gray-800 border-gray-700">
        <CardHeader className="pb-3">
          <CardTitle className="text-base text-white">Customer Rating</CardTitle>
        </CardHeader>
        <CardContent className="space-y-1 pt-0">
          {[5, 4, 3, 2, 1].map((rating) => (
            <button
              key={rating}
              onClick={() => updateFilters({ rating })}
              className={cn(
                'w-full flex items-center space-x-2 px-3 py-2 text-sm rounded-lg transition-all duration-200',
                currentParams.rating === rating
                  ? 'bg-orange-500 text-white shadow-md'
                  : 'text-gray-300 hover:bg-gray-700 hover:text-white'
              )}
            >
              <div className="flex">
                {[...Array(5)].map((_, i) => (
                  <Star
                    key={i}
                    className={cn(
                      'h-4 w-4',
                      i < rating
                        ? 'text-yellow-400 fill-current'
                        : 'text-gray-500'
                    )}
                  />
                ))}
              </div>
              <span>& Up</span>
            </button>
          ))}
        </CardContent>
      </Card>

      {/* Availability */}
      <Card className="bg-gray-800 border-gray-700">
        <CardHeader className="pb-3">
          <CardTitle className="text-base text-white">Availability</CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <label className="flex items-center space-x-2 cursor-pointer">
            <input
              type="checkbox"
              checked={currentParams.in_stock === true}
              onChange={(e) => updateFilters({
                in_stock: e.target.checked ? true : undefined
              })}
              className="rounded border-gray-600 text-orange-500 focus:ring-orange-500 bg-gray-700"
            />
            <span className="text-sm text-white">In Stock Only</span>
          </label>
        </CardContent>
      </Card>

      {/* Popular Tags */}
      <Card className="bg-gray-800 border-gray-700">
        <CardHeader className="pb-3">
          <CardTitle className="text-base text-white">Popular Tags</CardTitle>
        </CardHeader>
        <CardContent className="pt-0">
          <div className="flex flex-wrap gap-2">
            {[
              'New Arrival',
              'Best Seller',
              'On Sale',
              'Premium',
              'Eco-Friendly',
              'Limited Edition',
            ].map((tag) => {
              const isSelected = currentParams.tags?.includes(tag)
              return (
                <Badge
                  key={tag}
                  variant={isSelected ? 'default' : 'outline'}
                  className={cn(
                    "cursor-pointer transition-all duration-200",
                    isSelected
                      ? "bg-orange-500 text-white hover:bg-orange-600 border-orange-500"
                      : "border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white hover:border-gray-500"
                  )}
                  onClick={() => {
                    const currentTags = currentParams.tags || []
                    const newTags = isSelected
                      ? currentTags.filter(t => t !== tag)
                      : [...currentTags, tag]
                    updateFilters({ tags: newTags.length > 0 ? newTags : undefined })
                  }}
                >
                  {tag}
                  {isSelected && (
                    <X className="ml-1 h-3 w-3" />
                  )}
                </Badge>
              )
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
