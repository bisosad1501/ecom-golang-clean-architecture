'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { SingleImageUpload } from '@/components/ui/single-image-upload'
import { AdminFormLayout } from '@/components/ui/admin-form-layout'
import { FormSection } from '@/components/ui/form-section'
import { FormField } from '@/components/ui/form-field'
import { FormActions } from '@/components/ui/form-actions'
import { CategorySelect } from '@/components/ui/category-select'
import { useUpdateCategory, useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'

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

  // Filter out current category and its descendants to prevent circular references
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
      
      onSuccess?.()
    } catch (error: any) {
      console.error('Failed to update category:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    reset()
    onCancel?.()
  }

  return (
    <AdminFormLayout
      title="Edit Category"
      description={`Update details for "${category.name}"`}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Basic Information */}
        <FormSection title="Basic Information">
          <FormField
            label="Category Name"
            required
            error={errors.name?.message}
          >
            <Input
              {...register('name')}
              placeholder="Enter category name"
            />
          </FormField>

          <FormField
            label="Slug"
            error={errors.slug?.message}
            hint="Leave empty to auto-generate from name"
          >
            <Input
              {...register('slug')}
              placeholder="category-slug"
            />
          </FormField>

          <FormField
            label="Description"
            error={errors.description?.message}
          >
            <Textarea
              {...register('description')}
              placeholder="Enter category description"
              rows={3}
            />
          </FormField>
        </FormSection>

        {/* Organization */}
        <FormSection title="Organization">
          <CategorySelect
            categories={availableParentCategories}
            value={watch('parent_id') || ''}
            onChange={(value) => setValue('parent_id', value)}
            error={errors.parent_id?.message}
            placeholder="No parent (Top-level category)"
          />

          <FormField
            label="Sort Order"
            error={errors.sort_order?.message}
            hint="Lower numbers appear first (0 = highest priority)"
          >
            <Input
              type="number"
              min="0"
              {...register('sort_order', { valueAsNumber: true })}
              placeholder="0"
            />
          </FormField>
        </FormSection>

        {/* Image */}
        <FormSection title="Category Image">
          <SingleImageUpload
            label=""
            value={watch('image') || ''}
            onChange={(url: string) => setValue('image', url)}
            onRemove={() => setValue('image', '')}
            placeholder="Enter category image URL"
          />
          {errors.image && (
            <p className="text-sm text-red-600 mt-2">{errors.image.message}</p>
          )}
        </FormSection>

        {/* Status */}
        <FormSection title="Status">
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <label className="text-sm font-medium">Active Status</label>
              <p className="text-xs text-gray-500">
                Enable to make this category visible to customers
              </p>
            </div>
            <Switch
              checked={watch('is_active')}
              onCheckedChange={(checked) => setValue('is_active', checked)}
            />
          </div>
        </FormSection>

        {/* Actions */}
        <FormActions
          onCancel={handleCancel}
          submitLabel="Update Category"
          isSubmitting={isSubmitting}
        />
      </form>
    </AdminFormLayout>
  )
}
