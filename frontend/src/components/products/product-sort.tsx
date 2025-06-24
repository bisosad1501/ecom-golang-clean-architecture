'use client'

import { useRouter, useSearchParams } from 'next/navigation'
import { ChevronDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { PRODUCT_SORT_OPTIONS } from '@/constants'
import { cn } from '@/lib/utils'

interface ProductSortProps {
  currentSort: string
}

export function ProductSort({ currentSort }: ProductSortProps) {
  const router = useRouter()
  const searchParams = useSearchParams()

  const updateSort = (sortValue: string) => {
    const params = new URLSearchParams(searchParams.toString())
    const [sortBy, sortOrder] = sortValue.split(':')
    
    params.set('sort_by', sortBy)
    params.set('sort_order', sortOrder)
    params.set('page', '1') // Reset to page 1 when sorting changes

    router.push(`/products?${params.toString()}`)
  }

  const currentOption = PRODUCT_SORT_OPTIONS.find(option => option.value === currentSort)

  return (
    <div className="relative">
      <div className="flex items-center space-x-2">
        <span className="text-sm text-gray-600">Sort by:</span>
        <div className="relative">
          <select
            value={currentSort}
            onChange={(e) => updateSort(e.target.value)}
            className="appearance-none bg-white border border-gray-300 rounded-md px-4 py-2 pr-8 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
          >
            {PRODUCT_SORT_OPTIONS.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          <ChevronDown className="absolute right-2 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400 pointer-events-none" />
        </div>
      </div>
    </div>
  )
}
