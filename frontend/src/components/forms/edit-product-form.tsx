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
import { Badge } from '@/components/ui/badge'
import { X, Image as ImageIcon, AlertCircle, ChevronDown } from 'lucide-react'
import { useUpdateProduct } from '@/hooks/use-products'
import { transformUpdateProductData } from '@/lib/utils/product-transform'
import { Product, Category } from '@/types'
import { categoryService } from '@/lib/services/categories'
import { toast } from 'sonner'
import Image from 'next/image'
import { cn } from '@/lib/utils'
import { uploadMultipleImageFiles } from '@/lib/utils/image-upload'
import { ImageUploadGrid, type ImageUploadItem } from '@/components/ui/image-upload-grid'

const editProductSchema = z.object({
  name: z.string().min(1, 'Product name is required'),
  description: z.string().min(1, 'Description is required'),
  short_description: z.string().optional(),
  sku: z.string().min(1, 'SKU is required'),
  price: z.number().min(0, 'Price must be positive'),
  compare_price: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  cost_price: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  stock: z.number().int().min(0, 'Stock must be non-negative'),
  category_id: z.string().min(1, 'Category is required'),
  weight: z.number().min(0).optional().transform(val => val === 0 ? undefined : val),
  status: z.enum(['active', 'draft', 'archived']),
  is_digital: z.boolean(),
})

type EditProductFormData = z.infer<typeof editProductSchema>

interface ProductImage extends ImageUploadItem {}

interface EditProductFormProps {
  product: Product
  onSuccess?: () => void
  onCancel?: () => void
}

export function EditProductForm({ product, onSuccess, onCancel }: EditProductFormProps) {
  // ... existing state ...
  const [hoveredParentId, setHoveredParentId] = useState<string | null>(null)
  console.log('=== EditProductForm RENDER ===', { productId: product.id, timestamp: new Date().toISOString() })
  
  const [images, setImages] = useState<ProductImage[]>([])
  const [originalImages, setOriginalImages] = useState<ProductImage[]>([]) // Track original images
  const [imagesChanged, setImagesChanged] = useState(false) // Track if images were modified
  const [tags, setTags] = useState<string[]>([])
  const [newTag, setNewTag] = useState('')
  const [isUploadingImages, setIsUploadingImages] = useState(false)
  const [isSubmittingRef, setIsSubmittingRef] = useState(false) // Prevent double submission
  const [categories, setCategories] = useState<Category[]>([])
  const [loadingCategories, setLoadingCategories] = useState(true)
  const [categoryDropdownOpen, setCategoryDropdownOpen] = useState(false)

  const updateProductMutation = useUpdateProduct()

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setValue,
    watch,
    reset,
  } = useForm<EditProductFormData>({
    resolver: zodResolver(editProductSchema),
    defaultValues: {
      name: product.name,
      description: product.description,
      short_description: (product as any).short_description || '',
      sku: product.sku,
      price: (product as any).price,
      compare_price: (product as any).compare_price || undefined,
      cost_price: (product as any).cost_price || undefined,
      stock: (product as any).stock,
      // Backend uses ProductCategory many-to-many, but accepts category_id for primary category
      category_id: product.category?.id || '',
      weight: (product as any).weight || undefined,
      status: (product.status as any) || 'active',
      is_digital: Boolean((product as any).is_digital),
    },
  })

  // Register category_id field for custom dropdown
  useEffect(() => {
    register('category_id')
  }, [register])

  // Initialize images and tags from product
  useEffect(() => {
    console.log('=== DEBUG: useEffect initializing ===')
    console.log('Product:', product)
    console.log('Product.images:', product.images)
    console.log('Current isSubmittingRef:', isSubmittingRef)
    
    // Don't reset images if we're currently submitting or just finished submitting
    // This prevents overwriting the local state with potentially stale server data
    if (isSubmittingRef) {
      console.log('Skipping images reset - currently submitting')
      return
    }
    
    if (product.images) {
      const productImages = product.images.map((img, index) => ({
        url: img.url,
        alt_text: (img as any).alt_text || '',
        position: index,
      }))
      console.log('Mapped product images:', productImages)
      
      // Only update if the images are actually different from current state
      // This prevents unnecessary resets when the same product data comes from cache
      const currentImagesStr = JSON.stringify(images)
      const productImagesStr = JSON.stringify(productImages)
      
      if (currentImagesStr !== productImagesStr) {
        console.log('Product images differ from current state, updating...')
        setImages([...productImages]) // Create a copy
        setOriginalImages([...productImages]) // Create a separate copy for comparison
        setImagesChanged(false) // Reset flag
        console.log('Set images and originalImages to:', productImages)
      } else {
        console.log('Product images same as current state, no update needed')
      }
    } else {
      // Only clear images if we currently have images but product doesn't
      if (images.length > 0) {
        console.log('Product has no images, clearing current images')
        setImages([])
        setOriginalImages([])
        setImagesChanged(false)
      }
    }
    
    if (product.tags) {
      const productTags = (product as any).tags?.map((tag: any) => typeof tag === 'string' ? tag : tag.name) || []
      const currentTagsStr = JSON.stringify(tags)
      const productTagsStr = JSON.stringify(productTags)
      
      if (currentTagsStr !== productTagsStr) {
        setTags(productTags)
      }
    }
  }, [product, isSubmittingRef]) // Add isSubmittingRef as dependency

  // Reset form when product changes
  useEffect(() => {
    console.log('Product changed, resetting form with category_id:', product.category?.id)
    reset({
      name: product.name,
      description: product.description,
      short_description: product.short_description || '',
      sku: product.sku,
      price: product.price,
      compare_price: product.compare_price || undefined,
      cost_price: product.cost_price || undefined,
      stock: product.stock,
      category_id: product.category?.id || '',
      weight: product.weight || undefined,
      status: product.status || 'active',
      is_digital: Boolean(product.is_digital),
    })
  }, [product, reset])

  // Fetch categories on component mount
  useEffect(() => {
    const fetchCategories = async () => {
      try {
        setLoadingCategories(true)
        const categoryData = await categoryService.getCategories()
        setCategories(categoryData)
      } catch (error) {
        console.error('Failed to fetch categories:', error)
        toast.error('Failed to load categories')
      } finally {
        setLoadingCategories(false)
      }
    }

    fetchCategories()
  }, [])

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (categoryDropdownOpen) {
        const target = event.target as Element
        if (!target.closest('[data-category-dropdown]')) {
          setCategoryDropdownOpen(false)
        }
      }
    }

    document.addEventListener('mousedown', handleClickOutside)
    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [categoryDropdownOpen])

  const watchedFields = watch()

  // Helper function to compare images arrays intelligently
  const compareImages = (current: ProductImage[], original: ProductImage[]) => {
    console.log('=== compareImages ===')
    console.log('Current:', current)
    console.log('Original:', original)
    
    // Quick check: different lengths means definitely changed
    if (current.length !== original.length) {
      console.log('Length differs:', current.length, 'vs', original.length)
      return true
    }
    
    // Check each image for changes (URL, alt_text, position)
    for (let i = 0; i < current.length; i++) {
      const currentImg = current[i]
      const originalImg = original[i]
      
      if (
        currentImg.url !== originalImg.url ||
        currentImg.alt_text !== originalImg.alt_text ||
        currentImg.position !== originalImg.position
      ) {
        console.log('Image differs at index', i)
        console.log('Current img:', currentImg)
        console.log('Original img:', originalImg)
        return true
      }
    }
    
    console.log('No changes detected')
    return false
  }

  const addTag = () => {
    if (!newTag.trim()) return
    if (tags.includes(newTag.trim())) {
      toast.error('Tag already exists')
      return
    }
    if (tags.length >= 10) {
      toast.error('Maximum 10 tags allowed')
      return
    }

    setTags(prev => [...prev, newTag.trim()])
    setNewTag('')
  }

  const removeTag = (tagToRemove: string) => {
    setTags(prev => prev.filter(tag => tag !== tagToRemove))
  }

  // Function to analyze image changes in detail
  const analyzeImageChanges = (current: ProductImage[], original: ProductImage[]) => {
    console.log('=== analyzeImageChanges ===')
    
    const changes = {
      added: [] as ProductImage[],
      removed: [] as ProductImage[],
      updated: [] as ProductImage[],
      reordered: false,
      hasChanges: false
    }
    
    // Find added images (in current but not in original)
    changes.added = current.filter(currentImg => 
      !original.some(originalImg => originalImg.url === currentImg.url)
    )
    
    // Find removed images (in original but not in current)
    changes.removed = original.filter(originalImg => 
      !current.some(currentImg => currentImg.url === originalImg.url)
    )
    
    // Find updated images (same URL but different alt_text)
    changes.updated = current.filter(currentImg => {
      const originalImg = original.find(orig => orig.url === currentImg.url)
      return originalImg && (
        originalImg.alt_text !== currentImg.alt_text
      )
    })
    
    // Check if order changed (same images but different positions)
    if (current.length === original.length && changes.added.length === 0 && changes.removed.length === 0) {
      changes.reordered = current.some((currentImg, index) => {
        const originalImg = original[index]
        return originalImg && originalImg.url !== currentImg.url
      })
    }
    
    changes.hasChanges = changes.added.length > 0 || 
                        changes.removed.length > 0 || 
                        changes.updated.length > 0 || 
                        changes.reordered
    
    console.log('Image changes analysis:', changes)
    return changes
  }

  const onSubmit = async (data: EditProductFormData) => {
    // Prevent double submission
    if (isSubmittingRef) {
      console.log('=== PREVENTING DOUBLE SUBMISSION ===')
      return
    }
    
    setIsSubmittingRef(true)
    
    try {
      // Debug: Log current state
      console.log('=== DEBUG: onSubmit START ===')
      console.log('Form data:', data)
      console.log('Category ID from form:', data.category_id)
      console.log('Original product category:', product.category)
      // Note: product.category_id removed - Backend uses ProductCategory many-to-many
      console.log('Original product category ID:', product.category?.id)
      console.log('Current watchedFields:', watchedFields)
      console.log('Current categories available:', categories)
      console.log('Current images:', images)
      console.log('Original images:', originalImages)
      console.log('Images length:', images.length)
      console.log('Original images length:', originalImages.length)
      console.log('Images changed flag:', imagesChanged)
      
      // Analyze detailed changes
      const imageChanges = analyzeImageChanges(images, originalImages)
      console.log('Detailed image changes:', imageChanges)
      
      // Transform form data với images và tags
      const formDataWithExtras: any = {
        ...data,
        tags,
      }

      // Only include images if they have been modified
      if (imageChanges.hasChanges) {
        // Option 1: Send full images array (current approach)
        formDataWithExtras.images = images.map(img => ({
          url: img.url,
          alt_text: (img as any).alt_text || '',
          position: img.position,
        }))

        // Option 2: Send detailed changes (for future optimization)
        formDataWithExtras.image_changes = {
          added: imageChanges.added,
          removed: imageChanges.removed,
          updated: imageChanges.updated,
          reordered: imageChanges.reordered,
          full_list: formDataWithExtras.images // Fallback for current backend
        }

        console.log('Including images in request:', formDataWithExtras.images)
        console.log('Image changes details:', formDataWithExtras.image_changes)
      } else {
        console.log('No image changes detected, not including images in request')
      }

      console.log('EditProductForm payload:', formDataWithExtras)
      console.log('Product ID:', product.id)

      // Transform data using the utility function
      const transformedData = transformUpdateProductData(formDataWithExtras)
      console.log('Transformed data:', transformedData)

      const result = await updateProductMutation.mutateAsync({
        id: product.id,
        data: transformedData,
      })
      
      console.log('Update result:', result)

      // Reset the images changed flag after successful update
      setImagesChanged(false)
      
      // Update originalImages to current images to prevent future false positives
      setOriginalImages([...images])

      toast.success('Product updated successfully!')
      onSuccess?.()
    } catch (error: any) {
      console.error('EditProductForm error:', error)
      toast.error(error?.message || 'Failed to update product')
    } finally {
      setIsSubmittingRef(false)
    }
  }

  const handleReset = () => {
    reset()
    const productImages = product.images?.map((img, index) => ({
      url: img.url,
      alt_text: (img as any).alt_text || '',
      position: index,
    })) || []
    
    setImages(productImages)
    setOriginalImages(productImages)
    setImagesChanged(false) // Reset flag
    setTags((product as any).tags?.map((tag: any) => typeof tag === 'string' ? tag : tag.name) || [])
    setNewTag('')
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
      {/* Basic Information */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
        <h3 className="text-lg font-semibold text-white mb-6 flex items-center gap-2">
          <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
          Basic Information
        </h3>
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name" className="text-sm font-medium text-gray-300">Product Name</Label>
              <Input
                id="name"
                {...register('name')}
                className={`bg-gray-900 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500 ${errors.name ? 'border-red-500' : ''}`}
              />
              {errors.name && (
                <p className="text-sm text-red-400 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="sku" className="text-sm font-medium text-gray-300">SKU</Label>
              <Input
                id="sku"
                {...register('sku')}
                className={`bg-gray-900 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500 font-mono ${errors.sku ? 'border-red-500' : ''}`}
              />
              {errors.sku && (
                <p className="text-sm text-red-400 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.sku.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description" className="text-sm font-medium text-gray-300">Description</Label>
            <Textarea
              id="description"
              rows={4}
              {...register('description')}
              className={`bg-gray-900 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500 ${errors.description ? 'border-red-500' : ''}`}
            />
            {errors.description && (
              <p className="text-sm text-red-400 flex items-center gap-1">
                <AlertCircle className="h-4 w-4" />
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="short_description" className="text-sm font-medium text-gray-300">Short Description</Label>
            <Textarea
              id="short_description"
              rows={2}
              {...register('short_description')}
              className="bg-gray-900 border-gray-600 text-white focus:border-blue-500 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      {/* Images */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
        <div className="flex items-center gap-2 mb-6">
          <ImageIcon className="h-5 w-5 text-white" />
          <h3 className="text-lg font-semibold text-white">Product Images</h3>
          {imagesChanged && (
            <Badge variant="outline" className="text-orange-400 border-orange-400 bg-orange-500/10">
              Modified
            </Badge>
          )}
        </div>
        <p className="text-sm text-gray-400 mb-6">
          Upload images to showcase your product. The first image will be used as the primary product image.
        </p>
        <div>
          <ImageUploadGrid
            images={images}
            onImagesChange={(newImages) => {
              setImages(newImages)
              // Check if images actually changed using intelligent comparison
              const actuallyChanged = compareImages(newImages, originalImages)
              setImagesChanged(actuallyChanged)
            }}
            onUploadFiles={async (files) => {
              return await uploadMultipleImageFiles(files, '/admin/upload/image')
            }}
            maxImages={10}
            isUploading={isUploadingImages}
          />
        </div>
      </div>

      {/* Pricing & Inventory */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
        <h3 className="text-lg font-semibold text-white mb-6 flex items-center gap-2">
          <span className="w-2 h-2 bg-[#FF9000] rounded-full"></span>
          Pricing & Inventory
        </h3>
        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-gray-900 rounded-lg p-4 border border-gray-600 space-y-2">
              <Label htmlFor="price" className="text-sm font-medium text-gray-300">Price *</Label>
              <Input
                id="price"
                type="number"
                step="0.01"
                {...register('price', { valueAsNumber: true })}
                className={`bg-gray-800 border-gray-600 text-white focus:border-[#FF9000] focus:ring-[#FF9000] ${errors.price ? 'border-red-500' : ''}`}
              />
              {errors.price && (
                <p className="text-sm text-red-400 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.price.message}
                </p>
              )}
            </div>

            <div className="bg-gray-900 rounded-lg p-4 border border-gray-600 space-y-2">
              <Label htmlFor="compare_price" className="text-sm font-medium text-gray-300">Compare Price</Label>
              <Input
                id="compare_price"
                type="number"
                step="0.01"
                {...register('compare_price', { valueAsNumber: true })}
                className="bg-gray-800 border-gray-600 text-white focus:border-[#FF9000] focus:ring-[#FF9000]"
              />
            </div>

            <div className="bg-gray-900 rounded-lg p-4 border border-gray-600 space-y-2">
              <Label htmlFor="cost_price" className="text-sm font-medium text-gray-300">Cost Price</Label>
              <Input
                id="cost_price"
                type="number"
                step="0.01"
                {...register('cost_price', { valueAsNumber: true })}
                className="bg-gray-800 border-gray-600 text-white focus:border-[#FF9000] focus:ring-[#FF9000]"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="stock" className="text-sm font-medium text-gray-300">Stock Quantity</Label>
              <Input
                id="stock"
                type="number"
                {...register('stock', { valueAsNumber: true })}
                className={`bg-gray-900 border-gray-600 text-white focus:border-[#FF9000] focus:ring-[#FF9000] ${errors.stock ? 'border-red-500' : ''}`}
              />
              {errors.stock && (
                <p className="text-sm text-red-400 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.stock?.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="category_id">Category</Label>
              {loadingCategories ? (
                <div className="animate-pulse">
                  <div className="h-10 bg-gray-200 rounded-md"></div>
                </div>
              ) : (
                <div className="relative" data-category-dropdown>
                  <Button
                    type="button"
                    variant="outline"
                    role="combobox"
                    aria-expanded={categoryDropdownOpen}
                    className={cn(
                      "w-full justify-between",
                      !watchedFields.category_id && "text-muted-foreground",
                      errors.category_id && "border-red-500"
                    )}
                    onClick={() => setCategoryDropdownOpen(!categoryDropdownOpen)}
                  >
                    {watchedFields.category_id ? (
                      (() => {
                        const selectedCategory = categories.find(cat => cat.id === watchedFields.category_id)
                        if (!selectedCategory) return "Select a category"
                        const level = selectedCategory.level || 0
                        
                        if (level > 0) {
                          // Show parent > child for sub-categories
                          const parent = categories.find(cat => cat.id === selectedCategory.parent_id)
                          return `${parent?.name || ''} > ${selectedCategory.name}`
                        } else {
                          // Show just category name for parent categories
                          return selectedCategory.name
                        }
                      })()
                    ) : (
                      "Select a category"
                    )}
                    <ChevronDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                  
                  {categoryDropdownOpen && (
                    <div className="absolute z-50 w-full mt-1 bg-gray-900 border border-gray-600 rounded-md shadow-lg max-h-60">
                      <div 
                        className="flex"
                        onMouseLeave={() => setHoveredParentId(null)}
                      >
                        {/* Left Column - Parent Categories */}
                        <div className="flex-1 p-1 border-r border-gray-700">
                          {(() => {
                            // Group categories by parent
                            const parentCategories = categories.filter(cat => (cat.level || 0) === 0)
                            
                            return parentCategories.map((parentCategory) => {
                              const children = categories.filter(cat => cat.parent_id === parentCategory.id)
                              const isParentSelected = watchedFields.category_id === parentCategory.id
                              
                              return (
                                <div 
                                  key={parentCategory.id} 
                                  className="relative group/parent"
                                  onMouseEnter={() => {
                                    // Set which parent is being hovered for sub-categories display
                                    if (children.length > 0) {
                                      setHoveredParentId(parentCategory.id)
                                    }
                                  }}
                                >
                                  {/* Parent Category */}
                                  <div
                                    className={cn(
                                      "cursor-pointer select-none py-2.5 px-3 text-sm rounded transition-colors flex items-center justify-between",
                                      isParentSelected 
                                        ? "bg-blue-600 text-white" 
                                        : "hover:bg-gray-800 text-gray-300"
                                    )}
                                    onClick={() => {
                                      console.log('=== PARENT CATEGORY SELECTED ===')
                                      console.log('Selected parent category:', parentCategory.id)
                                      setValue('category_id', parentCategory.id)
                                      setCategoryDropdownOpen(false)
                                    }}
                                  >
                                    <div className="flex items-center">
                                      <span className={cn(
                                        "font-medium",
                                        isParentSelected ? "text-white" : "text-gray-200"
                                      )}>
                                        {parentCategory.name}
                                      </span>
                                      {isParentSelected && (
                                        <span className="text-blue-300 ml-2 font-bold">✓</span>
                                      )}
                                    </div>
                                    <div className="flex items-center">
                                      {children.length > 0 && (
                                        <>
                                          <span className="text-xs text-gray-400 mr-2">
                                            {children.length}
                                          </span>
                                          <svg 
                                            className="h-3 w-3 text-gray-400" 
                                            fill="none" 
                                            stroke="currentColor" 
                                            viewBox="0 0 24 24"
                                          >
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                                          </svg>
                                        </>
                                      )}
                                    </div>
                                  </div>
                                </div>
                              )
                            })
                          })()}
                        </div>
                        
                        {/* Right Column - Sub Categories */}
                        <div className="flex-1 p-1">
                          {hoveredParentId ? (
                            (() => {
                              const hoveredParent = categories.find(cat => cat.id === hoveredParentId)
                              const children = categories.filter(cat => cat.parent_id === hoveredParentId)
                              
                              if (!hoveredParent || children.length === 0) {
                                return (
                                  <div className="text-center py-4 text-gray-400 text-sm">
                                    Hover over a category to see subcategories
                                  </div>
                                )
                              }
                              
                              return (
                                <div>
                                  <div className="px-3 py-2 text-xs font-medium text-gray-400 border-b border-gray-700 mb-1">
                                    {hoveredParent.name}
                                  </div>
                                  {children.map((childCategory) => {
                                    const isChildSelected = watchedFields.category_id === childCategory.id
                                    
                                    return (
                                      <div
                                        key={childCategory.id}
                                        className={cn(
                                          "cursor-pointer select-none py-2 px-3 text-sm rounded transition-colors flex items-center justify-between",
                                          isChildSelected 
                                            ? "bg-blue-600 text-white font-medium" 
                                            : "hover:bg-gray-800 text-gray-300"
                                        )}
                                        onClick={() => {
                                          console.log('=== CHILD CATEGORY SELECTED ===')
                                          console.log('Selected child category:', childCategory.id)
                                          setValue('category_id', childCategory.id)
                                          setCategoryDropdownOpen(false)
                                        }}
                                      >
                                        <span>
                                          {childCategory.name}
                                        </span>
                                        {isChildSelected && (
                                          <span className="text-blue-300 font-bold">✓</span>
                                        )}
                                      </div>
                                    )
                                  })}
                                </div>
                              )
                            })()
                          ) : (
                            <div className="text-center py-4 text-gray-400 text-sm">
                              Hover over a category to see subcategories
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  )}
                </div>
              )}
              {errors.category_id && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.category_id.message}
                </p>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Product Details */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
        <h3 className="text-lg font-semibold text-white mb-6 flex items-center gap-2">
          <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
          Product Details
        </h3>
        <div className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="weight" className="text-sm font-medium text-gray-300">Weight (kg)</Label>
              <Input
                id="weight"
                type="number"
                step="0.01"
                {...register('weight', { valueAsNumber: true })}
                className="bg-gray-900 border-gray-600 text-white focus:border-purple-500 focus:ring-purple-500"
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="status" className="text-sm font-medium text-gray-300">Status</Label>
              <select
                {...register('status')}
                className="flex h-10 w-full rounded-md border border-gray-600 bg-gray-900 px-3 py-2 text-sm text-white ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              >
                <option value="active">Active</option>
                <option value="draft">Draft</option>
                <option value="archived">Archived</option>
              </select>
            </div>
          </div>

          <div className="flex items-center space-x-2">
            <input
              id="is_digital"
              type="checkbox"
              {...register('is_digital')}
              className="h-4 w-4 rounded border-gray-600 bg-gray-900 text-purple-500 focus:ring-purple-500 focus:ring-offset-gray-800"
            />
            <Label htmlFor="is_digital" className="text-gray-300">Digital Product</Label>
          </div>
        </div>
      </div>

      {/* Tags */}
      <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
        <h3 className="text-lg font-semibold text-white mb-6 flex items-center gap-2">
          <span className="w-2 h-2 bg-green-500 rounded-full"></span>
          Tags ({tags.length}/10)
        </h3>
        <div className="space-y-4">
          <div className="flex gap-2">
            <Input
              placeholder="Enter tag name"
              value={newTag}
              onChange={(e) => setNewTag(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
              className="bg-gray-900 border-gray-600 text-white focus:border-green-500 focus:ring-green-500 placeholder:text-gray-400"
            />
            <Button
              type="button"
              variant="outline"
              onClick={addTag}
              disabled={!newTag.trim() || tags.length >= 10}
            >
              Add Tag
            </Button>
          </div>

          {tags.length > 0 && (
            <div className="flex flex-wrap gap-2">
              {tags.map((tag) => (
                <Badge
                  key={tag}
                  variant="secondary"
                  className="flex items-center gap-1"
                >
                  {tag}
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-4 w-4 p-0 hover:bg-transparent"
                    onClick={() => removeTag(tag)}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </Badge>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Actions */}
      <div className="flex justify-end space-x-4">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting || updateProductMutation.isPending}
        >
          Cancel
        </Button>
        <Button
          type="button"
          variant="outline"
          onClick={handleReset}
          disabled={isSubmitting || updateProductMutation.isPending}
        >
          Reset
        </Button>
        <Button
          type="submit"
          disabled={isSubmitting || updateProductMutation.isPending}
          className="min-w-[120px]"
        >
          {isSubmitting || updateProductMutation.isPending ? 'Updating...' : 'Update Product'}
        </Button>
      </div>
    </form>
  )
}
