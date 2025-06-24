'use client'

import { useState } from 'react'
import { Switch } from '@/components/ui/switch'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Category } from '@/types'
import { useUpdateCategory } from '@/hooks/use-categories'
import { toast } from 'sonner'

interface CategoryStatusToggleProps {
  category: Category
  onUpdate?: () => void
}

export function CategoryStatusToggle({ category, onUpdate }: CategoryStatusToggleProps) {
  const [isUpdating, setIsUpdating] = useState(false)
  const updateCategory = useUpdateCategory()

  const handleStatusToggle = async (isActive: boolean) => {
    if (isUpdating) return

    setIsUpdating(true)
    try {
      await updateCategory.mutateAsync({
        id: category.id,
        data: { is_active: isActive }
      })
      
      toast.success(
        isActive 
          ? 'Category activated successfully!' 
          : 'Category deactivated successfully!'
      )
      onUpdate?.()
    } catch (error: any) {
      console.error('Failed to update category status:', error)
      toast.error('Failed to update category status')
    } finally {
      setIsUpdating(false)
    }
  }

  return (
    <Card className="border-l-4 border-l-blue-500">
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div>
            <h4 className="font-medium text-gray-900">Category Status</h4>
            <p className="text-sm text-gray-600 mt-1">
              {category.is_active 
                ? 'This category is currently active and visible to customers.'
                : 'This category is currently inactive and hidden from customers.'
              }
            </p>
          </div>
          <div className="flex items-center space-x-3">
            <span className="text-sm text-gray-500">
              {category.is_active ? 'Active' : 'Inactive'}
            </span>
            <Switch
              checked={category.is_active}
              onCheckedChange={handleStatusToggle}
              disabled={isUpdating}
            />
          </div>
        </div>
        
        {!category.is_active && (
          <div className="mt-3 p-3 bg-amber-50 border border-amber-200 rounded-lg">
            <p className="text-sm text-amber-800">
              <strong>Alternative to deletion:</strong> Instead of deleting this category, 
              you can keep it deactivated. This preserves data integrity while hiding 
              it from customers.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
