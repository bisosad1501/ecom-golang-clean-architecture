'use client'

import { useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Star, DollarSign, Package, Truck, Shield } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ProductsParams } from '@/lib/services/products'
import { PRICE_RANGES } from '@/constants'
import { cn } from '@/lib/utils'

interface CategoryFiltersProps {
  currentParams: ProductsParams
  categoryId: string
  categorySlug: string
}

export function CategoryFilters({ currentParams, categoryId, categorySlug }: CategoryFiltersProps) {
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

    router.push(`/categories/${categorySlug}?${params.toString()}`)
  }

  const clearAllFilters = () => {
    router.push(`/categories/${categorySlug}`)
  }

  const applyPriceRange = () => {
    updateFilters({
      min_price: priceRange.min ? parseInt(priceRange.min) : undefined,
      max_price: priceRange.max ? parseInt(priceRange.max) : undefined,
    })
  }

  const clearPriceRange = () => {
    setPriceRange({ min: '', max: '' })
    updateFilters({
      min_price: undefined,
      max_price: undefined,
    })
  }

  const hasActiveFilters = Boolean(
    currentParams.min_price ||
    currentParams.max_price ||
    currentParams.rating ||
    currentParams.featured ||
    currentParams.on_sale ||
    currentParams.stock_status ||
    currentParams.requires_shipping === false
  )

  return (
    <div className="space-y-5">
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
                'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
                currentParams.rating === rating
                  ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                  : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
              )}
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <div className="flex">
                    {[...Array(5)].map((_, i) => (
                      <Star
                        key={i}
                        className={cn(
                          "w-4 h-4",
                          i < rating ? "text-yellow-400 fill-current" : "text-gray-500"
                        )}
                      />
                    ))}
                  </div>
                  <span className="ml-2">& Up</span>
                </div>
                <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                  →
                </div>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Stock Status Filter */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-green-500/20 to-green-600/20 rounded-lg border border-green-500/30">
              <Package className="w-4 h-4 text-green-500" />
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Availability
            </h3>
          </div>
        </div>

        <div className="px-5 pb-5 space-y-2">
          {[
            { value: 'in_stock', label: 'In Stock' },
            { value: 'out_of_stock', label: 'Out of Stock' },
            { value: 'on_backorder', label: 'On Backorder' },
          ].map((status) => (
            <button
              key={status.value}
              onClick={() => updateFilters({
                stock_status: currentParams.stock_status === status.value ? undefined : status.value as any
              })}
              className={cn(
                'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
                currentParams.stock_status === status.value
                  ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                  : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
              )}
            >
              <div className="flex items-center justify-between">
                <span>{status.label}</span>
                <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                  →
                </div>
              </div>
            </button>
          ))}
        </div>
      </div>

      {/* Special Offers Filter */}
      <div className="bg-white/[0.08] backdrop-blur-xl rounded-xl border border-white/15 shadow-lg overflow-hidden transition-all duration-300 hover:bg-white/[0.10]">
        <div className="p-5 pb-4">
          <div className="flex items-center gap-2 mb-4">
            <div className="p-1.5 bg-gradient-to-br from-purple-500/20 to-purple-600/20 rounded-lg border border-purple-500/30">
              <Shield className="w-4 h-4 text-purple-500" />
            </div>
            <h3 className="text-base font-semibold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
              Special Offers
            </h3>
          </div>
        </div>

        <div className="px-5 pb-5 space-y-2">
          <button
            onClick={() => updateFilters({
              featured: currentParams.featured ? undefined : true
            })}
            className={cn(
              'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
              currentParams.featured
                ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
            )}
          >
            <div className="flex items-center justify-between">
              <span>Featured Products</span>
              <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                →
              </div>
            </div>
          </button>

          <button
            onClick={() => updateFilters({
              on_sale: currentParams.on_sale ? undefined : true
            })}
            className={cn(
              'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
              currentParams.on_sale
                ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
            )}
          >
            <div className="flex items-center justify-between">
              <span>On Sale</span>
              <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                →
              </div>
            </div>
          </button>

          <button
            onClick={() => updateFilters({
              requires_shipping: currentParams.requires_shipping ? undefined : false
            })}
            className={cn(
              'w-full text-left px-4 py-3 text-sm rounded-lg transition-all duration-300 font-medium group',
              currentParams.requires_shipping === false
                ? 'bg-gradient-to-r from-[#ff9000] to-[#ff9000] text-white shadow-md shadow-[#ff9000]/20 border border-[#ff9000]/50'
                : 'text-gray-300 hover:text-white bg-white/[0.05] hover:bg-white/[0.08] border border-white/10 hover:border-white/20 backdrop-blur-sm'
            )}
          >
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Truck className="w-4 h-4" />
                <span>Free Shipping</span>
              </div>
              <div className="opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                →
              </div>
            </div>
          </button>
        </div>
      </div>
    </div>
  )
}
