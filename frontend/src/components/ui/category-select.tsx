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
          "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background",
          "file:border-0 file:bg-transparent file:text-sm file:font-medium",
          "placeholder:text-muted-foreground",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
          "disabled:cursor-not-allowed disabled:opacity-50"
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
