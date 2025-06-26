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
    <div className="space-y-6">
      {/* Enhanced Category Filter */}
      <EnhancedCategoryFilter
        selectedCategoryId={currentParams.category_id}
        onCategoryChange={(categoryId) => updateFilters({ category_id: categoryId })}
        showProductCount={true}
        showSearch={true}
      />

      {/* Clear Filters */}
      {hasActiveFilters && (
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">Active Filters</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={clearAllFilters}
                className="text-primary-600 hover:text-primary-700"
              >
                Clear All
              </Button>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Price Range */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Price Range</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Quick price ranges */}
          <div className="space-y-2">
            {PRICE_RANGES.map((range) => (
              <button
                key={`${range.min}-${range.max}`}
                onClick={() => updateFilters({
                  min_price: range.min ?? undefined,
                  max_price: range.max ?? undefined,
                })}
                className={cn(
                  'w-full text-left px-3 py-2 text-sm rounded-md transition-colors',
                  currentParams.min_price === range.min && 
                  currentParams.max_price === range.max
                    ? 'bg-primary-100 text-primary-800'
                    : 'hover:bg-gray-100'
                )}
              >
                {range.label}
              </button>
            ))}
          </div>

          {/* Custom price range */}
          <div className="pt-4 border-t">
            <p className="text-sm font-medium mb-3">Custom Range</p>
            <div className="flex items-center space-x-2">
              <Input
                type="number"
                placeholder="Min"
                value={priceRange.min}
                onChange={(e) => setPriceRange(prev => ({ ...prev, min: e.target.value }))}
                className="text-sm"
              />
              <span className="text-gray-500">-</span>
              <Input
                type="number"
                placeholder="Max"
                value={priceRange.max}
                onChange={(e) => setPriceRange(prev => ({ ...prev, max: e.target.value }))}
                className="text-sm"
              />
            </div>
            <Button
              onClick={applyPriceRange}
              size="sm"
              className="w-full mt-3"
              variant="outline"
            >
              Apply
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Rating Filter */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Customer Rating</CardTitle>
        </CardHeader>
        <CardContent className="space-y-2">
          {[5, 4, 3, 2, 1].map((rating) => (
            <button
              key={rating}
              onClick={() => updateFilters({ rating })}
              className={cn(
                'w-full flex items-center space-x-2 px-3 py-2 text-sm rounded-md transition-colors',
                currentParams.rating === rating
                  ? 'bg-primary-100 text-primary-800'
                  : 'hover:bg-gray-100'
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
                        : 'text-gray-300'
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
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Availability</CardTitle>
        </CardHeader>
        <CardContent>
          <label className="flex items-center space-x-2 cursor-pointer">
            <input
              type="checkbox"
              checked={currentParams.in_stock === true}
              onChange={(e) => updateFilters({ 
                in_stock: e.target.checked ? true : undefined 
              })}
              className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
            />
            <span className="text-sm">In Stock Only</span>
          </label>
        </CardContent>
      </Card>

      {/* Categories */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Categories</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {/* This would be populated with actual categories from API */}
            {[
              { id: '1', name: 'Electronics', count: 150 },
              { id: '2', name: 'Fashion', count: 89 },
              { id: '3', name: 'Home & Garden', count: 67 },
              { id: '4', name: 'Sports', count: 45 },
              { id: '5', name: 'Books', count: 123 },
            ].map((category) => (
              <button
                key={category.id}
                onClick={() => updateFilters({ category_id: category.id })}
                className={cn(
                  'w-full flex items-center justify-between px-3 py-2 text-sm rounded-md transition-colors',
                  currentParams.category_id === category.id
                    ? 'bg-primary-100 text-primary-800'
                    : 'hover:bg-gray-100'
                )}
              >
                <span>{category.name}</span>
                <span className="text-gray-500">({category.count})</span>
              </button>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Popular Tags */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Popular Tags</CardTitle>
        </CardHeader>
        <CardContent>
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
                  className="cursor-pointer"
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
