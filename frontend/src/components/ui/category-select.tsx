'use client'

import React from 'react'
import { FormField } from '@/components/ui/form-field'
import { cn } from '@/lib/utils'

interface CategorySelectProps {
  categories: Array<{ id: string; name: string; parent_id?: string }>
  value: string
  onChange: (value: string) => void
  error?: string
  placeholder?: string
  excludeIds?: string[]
  className?: string
  disabled?: boolean
}

export function CategorySelect({
  categories,
  value,
  onChange,
  error,
  placeholder = "Select a category",
  excludeIds = [],
  className,
  disabled = false,
}: CategorySelectProps) {
  // Filter out excluded categories
  const availableCategories = categories.filter(cat => !excludeIds.includes(cat.id))
  
  // Group categories by parent
  const rootCategories = availableCategories.filter(cat => !cat.parent_id)
  const childCategories = availableCategories.filter(cat => cat.parent_id)
  
  return (
    <FormField
      label="Category"
      required
      error={error}
      className={className}
    >
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        disabled={disabled}
        className={cn(
          "flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000]",
          "disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
        )}
      >
        <option value="">{placeholder}</option>
        
        {/* Root categories */}
        {rootCategories.map((category) => (
          <option key={category.id} value={category.id}>
            {category.name}
          </option>
        ))}
        
        {/* Child categories grouped by parent */}
        {rootCategories.map((parent) => {
          const children = childCategories.filter(child => child.parent_id === parent.id)
          if (children.length === 0) return null
          
          return (
            <optgroup key={`group-${parent.id}`} label={parent.name}>
              {children.map((child) => (
                <option key={child.id} value={child.id}>
                  {child.name}
                </option>
              ))}
            </optgroup>
          )
        })}
      </select>
    </FormField>
  )
}
