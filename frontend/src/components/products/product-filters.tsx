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
    (currentParams as any).featured ||
    (currentParams as any).on_sale ||
    (currentParams as any).stock_status ||
    currentParams.tags?.length
  )

  return (
    <div className="space-y-5">
      {/* Enhanced Category Filter */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg p-5 transition-all duration-300 hover:bg-white/[0.10]">
        <EnhancedCategoryFilter
          selectedCategoryId={currentParams.category_id}
          onCategoryChange={(categoryId) => updateFilters({ category_id: categoryId })}
          showProductCount={true}
          showSearch={true}
        />
      </div>

      {/* Clear Filters */}
      {hasActiveFilters && (
        <div className="bg-gradient-to-r from-[#ff9000]/10 to-[#ff9000]/10 backdrop-blur-xl rounded-xl border border-[#ff9000]/20 shadow-lg p-4 transition-all duration-300">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <div className="w-2 h-2 bg-[#ff9000] rounded-full animate-pulse"></div>
              <span className="text-sm font-medium text-white">Active Filters</span>
            </div>
            <Button
              variant="ghost"
              size="sm"
              onClick={clearAllFilters}
              className="text-[#ff9000] hover:text-white hover:bg-[#ff9000]/20 border border-[#ff9000]/30 hover:border-[#ff9000]/50 rounded-lg px-3 py-1.5 text-xs font-medium transition-all duration-300 hover:scale-105"
            >
              Clear All
            </Button>
          </div>
        </div>
      )}

      {/* Price Range */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/20 rounded-lg border border-[#ff9000]/30">
              <div className="w-4 h-4 bg-gradient-to-br from-[#ff9000] to-[#ff9000] rounded"></div>
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Price Range
            </h3>
          </div>
        </div>
        
        <div className="px-5 pb-5 space-y-4">
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
                  'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
                  currentParams.min_price === range.min &&
                  currentParams.max_price === range.max
                    ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                    : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
                )}
              >
                <div className="flex items-center justify-between">
                  <span>{range.label}</span>
                  <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                    →
                  </div>
                </div>
              </button>
            ))}
          </div>

          {/* Custom price range */}
          <div className="pt-4 border-t border-white/10">
            <p className="text-sm font-medium mb-3 text-gray-300">Custom Range</p>
            <div className="flex items-center space-x-2">
              <Input
                type="number"
                placeholder="Min"
                value={priceRange.min}
                onChange={(e) => setPriceRange(prev => ({ ...prev, min: e.target.value }))}
                className="text-sm bg-white/[0.05] border-white/15 text-white placeholder-gray-400 backdrop-blur-sm focus:bg-white/[0.08] focus:border-[#ff9000]/50 transition-all duration-300"
              />
              <div className="px-2 text-gray-400 font-medium">—</div>
              <Input
                type="number"
                placeholder="Max"
                value={priceRange.max}
                onChange={(e) => setPriceRange(prev => ({ ...prev, max: e.target.value }))}
                className="text-sm bg-white/[0.05] border-white/15 text-white placeholder-gray-400 backdrop-blur-sm focus:bg-white/[0.08] focus:border-[#ff9000]/50 transition-all duration-300"
              />
            </div>
            <Button
              onClick={applyPriceRange}
              size="sm"
              className="w-full mt-3 bg-gradient-to-r from-[#ff9000] to-[#ff9000] hover:from-[#ff9000]/90 hover:to-[#ff9000]/90 text-white font-medium shadow-md shadow-[#ff9000]/20 transition-all duration-300 hover:scale-[1.02] rounded-lg"
            >
              Apply Range
            </Button>
          </div>
        </div>
      </div>

      {/* Rating Filter */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-yellow-500/20 to-yellow-600/20 rounded-lg border border-yellow-500/30">
              <Star className="w-4 h-4 text-yellow-500 fill-current" />
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Customer Rating
            </h3>
          </div>
        </div>
        
        <div className="px-5 pb-5 space-y-2">
          {[5, 4, 3, 2, 1].map((rating) => (
            <button
              key={rating}
              onClick={() => updateFilters({ rating })}
              className={cn(
                'w-full flex items-center justify-between px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
                currentParams.rating === rating
                  ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                  : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
              )}
            >
              <div className="flex items-center gap-3">
                <div className="flex">
                  {[...Array(5)].map((_, i) => (
                    <Star
                      key={i}
                      className={cn(
                        'h-4 w-4 transition-colors duration-300',
                        i < rating
                          ? 'text-yellow-500 fill-current'
                          : 'text-gray-500'
                      )}
                    />
                  ))}
                </div>
                <span>& Up</span>
              </div>
              <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                →
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Availability & Stock Status */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-green-500/20 to-green-600/20 rounded-lg border border-green-500/30">
              <div className="w-4 h-4 bg-gradient-to-br from-green-400 to-green-500 rounded-full"></div>
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Availability
            </h3>
          </div>
        </div>

        <div className="px-5 pb-5 space-y-3">
          <label className="flex items-center justify-between p-4 rounded-lg bg-white/[0.05] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all duration-300 cursor-pointer group">
            <div className="flex items-center gap-3">
              <div className="relative">
                <input
                  type="checkbox"
                  checked={currentParams.in_stock === true}
                  onChange={(e) => updateFilters({
                    in_stock: e.target.checked ? true : undefined
                  })}
                  className="w-4 h-4 rounded border-white/20 bg-white/10 text-[#ff9000] focus:ring-[#ff9000]/50 focus:ring-2 transition-all duration-300"
                />
              </div>
              <span className="text-sm font-medium text-white group-hover:text-[#ff9000]/80 transition-colors duration-300">
                In Stock Only
              </span>
            </div>
            <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-[#ff9000]">
              →
            </div>
          </label>

          {/* Stock Status Options */}
          <div className="space-y-2">
            {[
              { value: 'in_stock', label: 'In Stock', color: 'emerald' },
              { value: 'out_of_stock', label: 'Out of Stock', color: 'red' },
              { value: 'on_backorder', label: 'On Backorder', color: 'blue' },
            ].map((status) => (
              <label key={status.value} className="flex items-center justify-between p-3 rounded-lg bg-white/[0.03] border border-white/5 hover:bg-white/[0.06] hover:border-white/15 transition-all duration-300 cursor-pointer group">
                <div className="flex items-center gap-3">
                  <div className="relative">
                    <input
                      type="radio"
                      name="stock_status"
                      checked={(currentParams as any).stock_status === status.value}
                      onChange={(e) => updateFilters({
                        stock_status: e.target.checked ? status.value : undefined
                      } as any)}
                      className="w-4 h-4 rounded-full border-white/20 bg-white/10 text-[#ff9000] focus:ring-[#ff9000]/50 focus:ring-2 transition-all duration-300"
                    />
                  </div>
                  <span className="text-sm font-medium text-white group-hover:text-[#ff9000]/80 transition-colors duration-300">
                    {status.label}
                  </span>
                </div>
              </label>
            ))}
          </div>
        </div>
      </div>

      {/* Special Filters */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-purple-500/20 to-pink-500/20 rounded-lg border border-purple-500/30">
              <div className="w-4 h-4 bg-gradient-to-br from-purple-400 to-pink-500 rounded-full"></div>
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Special Offers
            </h3>
          </div>
        </div>

        <div className="px-5 pb-5 space-y-3">
          <label className="flex items-center justify-between p-4 rounded-lg bg-white/[0.05] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all duration-300 cursor-pointer group">
            <div className="flex items-center gap-3">
              <div className="relative">
                <input
                  type="checkbox"
                  checked={(currentParams as any).featured === true}
                  onChange={(e) => updateFilters({
                    featured: e.target.checked ? true : undefined
                  } as any)}
                  className="w-4 h-4 rounded border-white/20 bg-white/10 text-purple-500 focus:ring-purple-500/50 focus:ring-2 transition-all duration-300"
                />
              </div>
              <span className="text-sm font-medium text-white group-hover:text-purple-400 transition-colors duration-300">
                Featured Products
              </span>
            </div>
            <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-purple-400">
              ⭐
            </div>
          </label>

          <label className="flex items-center justify-between p-4 rounded-lg bg-white/[0.05] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition-all duration-300 cursor-pointer group">
            <div className="flex items-center gap-3">
              <div className="relative">
                <input
                  type="checkbox"
                  checked={(currentParams as any).on_sale === true}
                  onChange={(e) => updateFilters({
                    on_sale: e.target.checked ? true : undefined
                  } as any)}
                  className="w-4 h-4 rounded border-white/20 bg-white/10 text-[#ff9000] focus:ring-[#ff9000]/50 focus:ring-2 transition-all duration-300"
                />
              </div>
              <span className="text-sm font-medium text-white group-hover:text-[#ff9000]/80 transition-colors duration-300">
                On Sale
              </span>
            </div>
            <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300 text-[#ff9000]">
              🔥
            </div>
          </label>
        </div>
      </div>

      {/* Popular Tags */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-purple-500/20 to-purple-600/20 rounded-lg border border-purple-500/30">
              <div className="w-4 h-4 bg-gradient-to-br from-purple-400 to-purple-500 rounded" style={{clipPath: 'polygon(0% 0%, 75% 0%, 100% 50%, 75% 100%, 0% 100%)'}}></div>
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Popular Tags
            </h3>
          </div>
        </div>
        
        <div className="px-5 pb-5">
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
                    "cursor-pointer transition-all duration-300 font-medium text-xs px-3 py-1.5 rounded-lg",
                    isSelected
                      ? "bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white border-[#ff9000]/50 shadow-md shadow-[#ff9000]/20 hover:scale-105"
                      : "border-white/20 text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] hover:border-white/30 backdrop-blur-sm hover:scale-105"
                  )}
                  onClick={() => {
                    const currentTags = currentParams.tags || []
                    const newTags = isSelected
                      ? currentTags.filter(t => t !== tag)
                      : [...currentTags, tag]
                    updateFilters({ tags: newTags.length > 0 ? newTags : undefined })
                  }}
                >
                  <div className="flex items-center gap-1.5">
                    <span>{tag}</span>
                    {isSelected && (
                      <X className="h-3 w-3 opacity-80" />
                    )}
                  </div>
                </Badge>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
