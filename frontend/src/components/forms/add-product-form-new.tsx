'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { AdminFormLayout } from '@/components/ui/admin-form-layout'
import { FormSection } from '@/components/ui/form-section'
import { FormField } from '@/components/ui/form-field'
import { FormActions } from '@/components/ui/form-actions'
import { CategorySelect } from '@/components/ui/category-select'
import { TagsInput } from '@/components/ui/tags-input'
import { MultiImageUpload } from '@/components/ui/multi-image-upload'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productService } from '@/lib/services/products'
import { categoryService } from '@/lib/services/categories'
import { CreateProductRequest } from '@/lib/services/products'
import { toast } from 'sonner'

const addProductSchema = z.object({
  name: z.string()
    .min(1, 'Product name is required')
    .max(200, 'Name too long')
    .regex(/^[a-zA-Z0-9\s\-_.,()&]+$/, 'Name contains invalid characters'),
  description: z.string()
    .min(10, 'Description must be at least 10 characters')
    .max(5000, 'Description too long'),
  short_description: z.string()
    .max(500, 'Short description too long')
    .optional(),
  sku: z.string()
    .min(1, 'SKU is required')
    .max(100, 'SKU too long')
    .regex(/^[A-Z0-9\-_]+$/, 'SKU must contain only uppercase letters, numbers, hyphens, and underscores'),
  price: z.number()
    .min(0.01, 'Price must be greater than 0')
    .max(999999.99, 'Price too high'),
  compare_price: z.number()
    .min(0)
    .max(999999.99, 'Compare price too high')
    .optional()
    .transform(val => val === 0 ? undefined : val),
  cost_price: z.number()
    .min(0)
    .max(999999.99, 'Cost price too high')
    .optional()
    .transform(val => val === 0 ? undefined : val),
  stock: z.number()
    .int('Stock must be a whole number')
    .min(0, 'Stock must be non-negative')
    .max(999999, 'Stock too high'),
  category_id: z.string().min(1, 'Category is required'),
  weight: z.number()
    .min(0)
    .max(999999, 'Weight too high')
    .optional()
    .transform(val => val === 0 ? undefined : val),
  length: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  width: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  height: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  status: z.enum(['active', 'draft', 'archived']),
  is_digital: z.boolean(),
}).refine((data) => {
  // Compare price should be higher than regular price
  if (data.compare_price && data.compare_price <= data.price) {
    return false
  }
  return true
}, {
  message: "Compare price must be higher than regular price",
  path: ["compare_price"]
}).refine((data) => {
  // Cost price should be lower than regular price
  if (data.cost_price && data.cost_price >= data.price) {
    return false
  }
  return true
}, {
  message: "Cost price should be lower than selling price",
  path: ["cost_price"]
})

type AddProductFormData = z.infer<typeof addProductSchema>

interface ImageItem {
  url: string
  alt_text?: string
  position: number
}

interface AddProductFormProps {
  onSuccess?: () => void
  onCancel?: () => void
}

export function AddProductForm({ onSuccess, onCancel }: AddProductFormProps) {
  const [tags, setTags] = useState<string[]>([])
  const [images, setImages] = useState<ImageItem[]>([])
  const queryClient = useQueryClient()

  const { data: categories = [], isLoading: categoriesLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryService.getCategories(),
  })

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
    reset,
  } = useForm<AddProductFormData>({
    resolver: zodResolver(addProductSchema),
    defaultValues: {
      name: '',
      description: '',
      short_description: '',
      sku: '',
      price: 0,
      compare_price: undefined,
      cost_price: undefined,
      stock: 0,
      category_id: '',
      weight: undefined,
      length: undefined,
      width: undefined,
      height: undefined,
      status: 'draft',
      is_digital: false,
    },
  })

  const createProductMutation = useMutation({
    mutationFn: async (data: CreateProductRequest) => {
      return await productService.createProduct(data)
    },
    onSuccess: () => {
      toast.success('Product created successfully!')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      
      // Reset form
      reset()
      setTags([])
      setImages([])
      onSuccess?.()
    },
    onError: (error: any) => {
      console.error('Product creation failed:', error)
      let errorMessage = 'Failed to create product'
      
      if (error.response?.data?.message) {
        errorMessage = error.response.data.message
      } else if (error.message) {
        errorMessage = error.message
      }
      
      toast.error(errorMessage)
    },
  })

  const onSubmit = async (data: AddProductFormData) => {
    try {
      const cleanData: CreateProductRequest = {
        ...data,
        tags: tags.length > 0 ? tags : undefined,
        images: images.length > 0 ? images.map(img => ({
          url: img.url,
          alt_text: img.alt_text || '',
          position: img.position,
        })) : undefined,
        dimensions: (data.length || data.width || data.height) ? {
          length: data.length || 0,
          width: data.width || 0,
          height: data.height || 0,
        } : undefined,
      }

      await createProductMutation.mutateAsync(cleanData)
    } catch (error) {
      // Error handling is done in the mutation
    }
  }

  const handleCancel = () => {
    reset()
    setTags([])
    setImages([])
    onCancel?.()
  }

  // Auto-generate SKU from name
  const nameValue = watch('name')
  const skuValue = watch('sku')

  useEffect(() => {
    if (nameValue && !skuValue) {
      const autoSku = nameValue
        .toUpperCase()
        .trim()
        .replace(/[^A-Z0-9\s]/g, '')
        .replace(/\s+/g, '-')
        .slice(0, 20)
      setValue('sku', autoSku)
    }
  }, [nameValue, skuValue, setValue])

  return (
    <AdminFormLayout
      title="Add New Product"
      description="Create a new product for your store"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Basic Information */}
        <FormSection title="Basic Information">
          <FormField
            label="Product Name"
            required
            error={errors.name?.message}
          >
            <Input
              {...register('name')}
              placeholder="Enter product name"
            />
          </FormField>

          <FormField
            label="SKU"
            required
            error={errors.sku?.message}
            hint="Stock Keeping Unit - auto-generated from name"
          >
            <Input
              {...register('sku')}
              placeholder="PRODUCT-SKU"
            />
          </FormField>

          <FormField
            label="Description"
            required
            error={errors.description?.message}
          >
            <Textarea
              {...register('description')}
              placeholder="Enter detailed product description"
              rows={4}
            />
          </FormField>

          <FormField
            label="Short Description"
            error={errors.short_description?.message}
            hint="Brief summary for product listings"
          >
            <Textarea
              {...register('short_description')}
              placeholder="Enter brief product summary"
              rows={2}
            />
          </FormField>
        </FormSection>

        {/* Pricing & Inventory */}
        <FormSection title="Pricing & Inventory">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <FormField
              label="Price"
              required
              error={errors.price?.message}
            >
              <Input
                type="number"
                step="0.01"
                min="0"
                {...register('price', { valueAsNumber: true })}
                placeholder="0.00"
              />
            </FormField>

            <FormField
              label="Compare Price"
              error={errors.compare_price?.message}
              hint="Original price for discounts"
            >
              <Input
                type="number"
                step="0.01"
                min="0"
                {...register('compare_price', { valueAsNumber: true })}
                placeholder="0.00"
              />
            </FormField>

            <FormField
              label="Cost Price"
              error={errors.cost_price?.message}
              hint="Your cost (not shown to customers)"
            >
              <Input
                type="number"
                step="0.01"
                min="0"
                {...register('cost_price', { valueAsNumber: true })}
                placeholder="0.00"
              />
            </FormField>
          </div>

          <FormField
            label="Stock Quantity"
            required
            error={errors.stock?.message}
          >
            <Input
              type="number"
              min="0"
              {...register('stock', { valueAsNumber: true })}
              placeholder="0"
            />
          </FormField>
        </FormSection>

        {/* Organization */}
        <FormSection title="Organization">
          <CategorySelect
            categories={categories}
            value={watch('category_id') || ''}
            onChange={(value) => setValue('category_id', value)}
            error={errors.category_id?.message}
            disabled={categoriesLoading}
          />

          <TagsInput
            tags={tags}
            onTagsChange={setTags}
            placeholder="Add product tags..."
          />
        </FormSection>

        {/* Images */}
        <FormSection title="Product Images">
          <MultiImageUpload
            images={images}
            onImagesChange={setImages}
            maxImages={10}
            endpoint="admin"
          />
        </FormSection>

        {/* Physical Properties */}
        <FormSection title="Physical Properties">
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <label className="text-sm font-medium">Digital Product</label>
                <p className="text-xs text-gray-500">
                  No shipping required for digital products
                </p>
              </div>
              <Switch
                checked={watch('is_digital')}
                onCheckedChange={(checked) => setValue('is_digital', checked)}
              />
            </div>

            {!watch('is_digital') && (
              <>
                <FormField
                  label="Weight (kg)"
                  error={errors.weight?.message}
                >
                  <Input
                    type="number"
                    step="0.01"
                    min="0"
                    {...register('weight', { valueAsNumber: true })}
                    placeholder="0.00"
                  />
                </FormField>

                <div className="grid grid-cols-3 gap-4">
                  <FormField
                    label="Length (cm)"
                    error={errors.length?.message}
                  >
                    <Input
                      type="number"
                      step="0.01"
                      min="0"
                      {...register('length', { valueAsNumber: true })}
                      placeholder="0.00"
                    />
                  </FormField>

                  <FormField
                    label="Width (cm)"
                    error={errors.width?.message}
                  >
                    <Input
                      type="number"
                      step="0.01"
                      min="0"
                      {...register('width', { valueAsNumber: true })}
                      placeholder="0.00"
                    />
                  </FormField>

                  <FormField
                    label="Height (cm)"
                    error={errors.height?.message}
                  >
                    <Input
                      type="number"
                      step="0.01"
                      min="0"
                      {...register('height', { valueAsNumber: true })}
                      placeholder="0.00"
                    />
                  </FormField>
                </div>
              </>
            )}
          </div>
        </FormSection>

        {/* Status */}
        <FormSection title="Status">
          <FormField
            label="Product Status"
            required
            error={errors.status?.message}
          >
            <select
              {...register('status')}
              className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="archived">Archived</option>
            </select>
          </FormField>
        </FormSection>

        {/* Actions */}
        <FormActions
          onCancel={handleCancel}
          submitLabel="Create Product"
          isSubmitting={isSubmitting || createProductMutation.isPending}
        />
      </form>
    </AdminFormLayout>
  )
}
