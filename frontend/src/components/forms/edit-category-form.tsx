'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { SingleImageUpload } from '@/components/ui/single-image-upload'
import { AlertCircle, FolderTree, Loader2 } from 'lucide-react'
import { useUpdateCategory, useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'

const editCategorySchema = z.object({
  name: z.string().min(1, 'Category name is required').max(100, 'Name too long'),
  slug: z.string().optional(),
  description: z.string().optional(),
  parent_id: z.string().optional(),
  image: z.string().url().optional().or(z.literal('')),
  is_active: z.boolean(),
  sort_order: z.number().int().min(0).optional(),
})

type EditCategoryFormData = z.infer<typeof editCategorySchema>

interface EditCategoryFormProps {
  category: Category
  onSuccess?: () => void
  onCancel?: () => void
}

export function EditCategoryForm({ category, onSuccess, onCancel }: EditCategoryFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const updateCategory = useUpdateCategory()
  const { data: categories } = useCategories()

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<EditCategoryFormData>({
    resolver: zodResolver(editCategorySchema),
    defaultValues: {
      name: category.name,
      slug: category.slug || '',
      description: category.description || '',
      parent_id: category.parent_id || '',
      image: category.image || '',
      is_active: category.is_active ?? true,
      sort_order: category.sort_order || 0,
    },
  })

  // Watch name to auto-generate slug
  const nameValue = watch('name')
  const slugValue = watch('slug')

  useEffect(() => {
    if (nameValue && !slugValue) {
      const autoSlug = nameValue
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')
      setValue('slug', autoSlug)
    }
  }, [nameValue, slugValue, setValue])

  // Filter out the current category and its descendants to prevent circular references
  const getAvailableParentCategories = () => {
    if (!categories) return []
    
    const findDescendants = (categoryId: string): string[] => {
      const directChildren = categories.filter(cat => cat.parent_id === categoryId)
      const allDescendants = [categoryId]
      
      directChildren.forEach(child => {
        allDescendants.push(...findDescendants(child.id))
      })
      
      return allDescendants
    }

    const excludeIds = findDescendants(category.id)
    return categories.filter(cat => !excludeIds.includes(cat.id))
  }

  const availableParentCategories = getAvailableParentCategories()

  const onSubmit = async (data: EditCategoryFormData) => {
    if (isSubmitting) return

    setIsSubmitting(true)
    try {
      // Clean data before sending
      const cleanData = {
        name: data.name.trim(),
        slug: data.slug?.trim() || undefined,
        description: data.description?.trim() || undefined,
        parent_id: data.parent_id || undefined,
        image: data.image?.trim() === '' ? '' : (data.image?.trim() || undefined),
        is_active: data.is_active,
        sort_order: data.sort_order || 0,
      }

      await updateCategory.mutateAsync({
        id: category.id,
        data: cleanData,
      })

      toast.success('Category updated successfully!')
      onSuccess?.()
    } catch (error: any) {
      console.error('Failed to update category:', error)
      const errorMessage = error?.message || error?.error || 'Failed to update category. Please try again.'
      toast.error(errorMessage)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    reset()
    onCancel?.()
  }

  return (
    <div className="h-full max-h-[70vh] overflow-y-auto">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Basic Information Section */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Basic Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="name">Category Name *</Label>
                  <Input
                    id="name"
                    {...register('name')}
                    placeholder="Enter category name"
                    className={cn(errors.name && 'border-red-500')}
                  />
                  {errors.name && (
                    <div className="flex items-center text-sm text-red-600">
                      <AlertCircle className="h-4 w-4 mr-1" />
                      {errors.name.message}
                    </div>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="slug">Slug</Label>
                  <Input
                    id="slug"
                    {...register('slug')}
                    placeholder="category-slug (auto-generated)"
                    className={cn(errors.slug && 'border-red-500')}
                  />
                  {errors.slug && (
                    <div className="flex items-center text-sm text-red-600">
                      <AlertCircle className="h-4 w-4 mr-1" />
                      {errors.slug.message}
                    </div>
                  )}
                  <p className="text-xs text-gray-500">
                    Leave empty to auto-generate from name
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="description">Description</Label>
                  <Textarea
                    id="description"
                    {...register('description')}
                    placeholder="Enter category description"
                    rows={3}
                    className={cn(errors.description && 'border-red-500')}
                  />
                  {errors.description && (
                    <div className="flex items-center text-sm text-red-600">
                      <AlertCircle className="h-4 w-4 mr-1" />
                      {errors.description.message}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Organization Section */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Organization</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2">
                  <Label htmlFor="parent_id">Parent Category</Label>
                  <select
                    id="parent_id"
                    {...register('parent_id')}
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    <option value="">No parent (Top-level category)</option>
                    {availableParentCategories
                      .filter(cat => !cat.parent_id) // Only show parent categories
                      .map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.name}
                        </option>
                      ))}
                  </select>
                  {errors.parent_id && (
                    <div className="flex items-center text-sm text-red-600">
                      <AlertCircle className="h-4 w-4 mr-1" />
                      {errors.parent_id.message}
                    </div>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="sort_order">Sort Order</Label>
                  <Input
                    id="sort_order"
                    type="number"
                    min="0"
                    {...register('sort_order', { valueAsNumber: true })}
                    placeholder="0"
                    className={cn(errors.sort_order && 'border-red-500')}
                  />
                  {errors.sort_order && (
                    <div className="flex items-center text-sm text-red-600">
                      <AlertCircle className="h-4 w-4 mr-1" />
                      {errors.sort_order.message}
                    </div>
                  )}
                  <p className="text-xs text-gray-500">
                    Lower numbers appear first (0 = highest priority)
                  </p>
                </div>
              </CardContent>
            </Card>

            {/* Image Section */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Category Image</CardTitle>
              </CardHeader>
              <CardContent>
                <SingleImageUpload
                  label=""
                  value={watch('image') || ''}
                  onChange={(url: string) => setValue('image', url)}
                  onRemove={() => setValue('image', '')}
                  placeholder="Enter category image URL"
                />
                {errors.image && (
                  <div className="flex items-center text-sm text-red-600 mt-2">
                    <AlertCircle className="h-4 w-4 mr-1" />
                    {errors.image.message}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Status Section */}
            <Card>
              <CardHeader>
                <CardTitle className="text-base">Status</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <div className="space-y-1">
                    <Label htmlFor="is_active">Active Status</Label>
                    <p className="text-xs text-gray-500">
                      Enable to make this category visible to customers
                    </p>
                  </div>
                  <Switch
                    id="is_active"
                    checked={watch('is_active')}
                    onCheckedChange={(checked) => setValue('is_active', checked)}
                  />
                </div>
              </CardContent>
            </Card>

            {/* Form Actions */}
            <div className="flex justify-end space-x-4 pt-6 border-t">
              <Button
                type="button"
                variant="outline"
                onClick={handleCancel}
                disabled={isSubmitting}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                disabled={isSubmitting}
                className="min-w-[120px]"
              >
                {isSubmitting ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Updating...
                  </>
                ) : (
                  'Update Category'
                )}
              </Button>
            </div>
          </form>
    </div>
  )
}
