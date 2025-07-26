'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { SingleImageUpload } from '@/components/ui/single-image-upload'
import { FormSection } from '@/components/ui/form-section'
import { FormField } from '@/components/ui/form-field'
import { FormActions } from '@/components/ui/form-actions'
import { useCreateCategory, useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'

const addCategorySchema = z.object({
  name: z.string().min(1, 'Category name is required').max(100, 'Name too long'),
  slug: z.string().optional(),
  description: z.string().optional(),
  parent_id: z.string().optional(),
  image: z.string().url().optional().or(z.literal('')),
  is_active: z.boolean(),
  sort_order: z.number().int().min(0).optional(),
})

type AddCategoryFormData = z.infer<typeof addCategorySchema>

interface AddCategoryFormProps {
  parentCategory?: Category // Optional parent category to auto-select
  onSuccess?: () => void
  onCancel?: () => void
}

export function AddCategoryForm({ parentCategory, onSuccess, onCancel }: AddCategoryFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const createCategory = useCreateCategory()
  const { data: categories } = useCategories()

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    watch,
    reset,
  } = useForm<AddCategoryFormData>({
    resolver: zodResolver(addCategorySchema),
    defaultValues: {
      name: '',
      slug: '',
      description: '',
      parent_id: parentCategory?.id || '',
      image: '',
      is_active: true,
      sort_order: 0,
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

  const onSubmit = async (data: AddCategoryFormData) => {
    if (isSubmitting) return

    setIsSubmitting(true)
    try {
      // Generate slug if not provided
      const slug = data.slug?.trim() || data.name
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '')

      const cleanData = {
        name: data.name.trim(),
        slug,
        description: data.description?.trim() || undefined,
        parent_id: data.parent_id || undefined,
        image: data.image?.trim() || undefined,
        is_active: data.is_active,
        sort_order: data.sort_order || 0,
      }

      await createCategory.mutateAsync(cleanData)
      onSuccess?.()
    } catch (error: any) {
      console.error('Failed to create category:', error)
      // Error handling is done in the hook
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleCancel = () => {
    reset()
    onCancel?.()
  }

  // Filter root categories for parent selection
  const rootCategories = categories?.filter(cat => !cat.parent_id) || []

  return (
    <div className="bihub-admin-form h-full max-h-[70vh] overflow-y-auto">
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Basic Information Section */}
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
              placeholder="category-slug (auto-generated)"
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

        {/* Organization Section */}
        <FormSection title="Organization">
          <FormField
            label="Parent Category"
            error={errors.parent_id?.message}
          >
            <select
              {...register('parent_id')}
              className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
            >
              <option value="">No parent (Top-level category)</option>
              {rootCategories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </FormField>

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

        {/* Image Section */}
        <FormSection title="Category Image">
          <FormField
            label="Image URL"
            error={errors.image?.message}
          >
            <SingleImageUpload
              label=""
              value={watch('image') || ''}
              onChange={(url: string) => setValue('image', url)}
              onRemove={() => setValue('image', '')}
              placeholder="Enter category image URL"
            />
          </FormField>
        </FormSection>

        {/* Status Section */}
        <FormSection title="Status">
          <FormField
            label="Active Status"
            hint="Enable to make this category visible to customers"
          >
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-300">Category is active</span>
              <Switch
                checked={watch('is_active')}
                onCheckedChange={(checked) => setValue('is_active', checked)}
              />
            </div>
          </FormField>
        </FormSection>

        {/* Form Actions */}
        <FormActions
          onCancel={handleCancel}
          submitLabel="Create Category"
          isSubmitting={isSubmitting}
        />
      </form>
    </div>
  )
}
