'use client'

import { useState, useEffect, useRef } from 'react'
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
  const [isOpen, setIsOpen] = useState(false)
  const dropdownRef = useRef<HTMLDivElement>(null)

  const updateSort = (sortValue: string) => {
    const params = new URLSearchParams(searchParams.toString())
    const [sortBy, sortOrder] = sortValue.split(':')
    
    params.set('sort_by', sortBy)
    params.set('sort_order', sortOrder)
    params.set('page', '1') // Reset to page 1 when sorting changes

    router.push(`/products?${params.toString()}`)
    setIsOpen(false)
  }

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false)
      }
    }

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isOpen])

  const currentOption = PRODUCT_SORT_OPTIONS.find(option => option.value === currentSort)

  return (
    <div className="relative" ref={dropdownRef}>
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setIsOpen(!isOpen)}
        className="bg-white/[0.06] backdrop-blur-md border border-white/10 rounded-lg px-3 py-1.5 pr-7 text-xs text-gray-400 hover:text-white focus:text-white font-medium hover:bg-white/[0.08] hover:border-white/15 transition-all duration-200 min-w-[110px] shadow-sm justify-start"
      >
        <span className="truncate">{currentOption?.label || 'Newest'}</span>
        <ChevronDown className={cn(
          "h-3 w-3 text-gray-500 absolute right-2 top-1/2 transform -translate-y-1/2 transition-transform duration-200",
          isOpen && "rotate-180"
        )} />
      </Button>

      {isOpen && (
        <div className="absolute top-full left-0 mt-1 w-full bg-gray-900/95 backdrop-blur-xl border border-white/15 rounded-lg shadow-lg shadow-black/20 z-20 py-1 animate-in fade-in slide-in-from-top-1 duration-200">
          {PRODUCT_SORT_OPTIONS.map((option) => (
            <button
              key={option.value}
              onClick={() => updateSort(option.value)}
              className={cn(
                "w-full text-left px-3 py-2 text-xs font-medium transition-all duration-200 hover:bg-white/[0.08]",
                option.value === currentSort
                  ? "text-[#ff9000] bg-[#ff9000]/10"
                  : "text-gray-300 hover:text-white"
              )}
            >
              {option.label}
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
