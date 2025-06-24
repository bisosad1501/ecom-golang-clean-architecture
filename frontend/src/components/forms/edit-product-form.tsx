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
import { X, Upload, Link, Star, Image as ImageIcon, AlertCircle, ArrowUp, ArrowDown, Edit2, GripVertical } from 'lucide-react'
import { useUpdateProduct } from '@/hooks/use-products'
import { transformUpdateProductData } from '@/lib/utils/product-transform'
import { Product } from '@/types'
import { toast } from 'sonner'
import Image from 'next/image'
import { DragDropContext, Droppable, Draggable, type DropResult } from '@hello-pangea/dnd'

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

interface ProductImage {
  url: string
  alt_text?: string
  position: number
}

interface EditProductFormProps {
  product: Product
  onSuccess?: () => void
  onCancel?: () => void
}

export function EditProductForm({ product, onSuccess, onCancel }: EditProductFormProps) {
  console.log('=== EditProductForm RENDER ===', { productId: product.id, timestamp: new Date().toISOString() })
  
  const [images, setImages] = useState<ProductImage[]>([])
  const [originalImages, setOriginalImages] = useState<ProductImage[]>([]) // Track original images
  const [imagesChanged, setImagesChanged] = useState(false) // Track if images were modified
  const [tags, setTags] = useState<string[]>([])
  const [newTag, setNewTag] = useState('')
  const [imageUrl, setImageUrl] = useState('')
  const [isSubmittingRef, setIsSubmittingRef] = useState(false) // Prevent double submission

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
      short_description: product.short_description || '',
      sku: product.sku,
      price: product.price,
      compare_price: product.compare_price || undefined,
      cost_price: product.cost_price || undefined,
      stock: product.stock,
      category_id: product.category?.id || '',
      weight: product.weight || undefined,
      status: (product.status as any) || 'active',
      is_digital: Boolean(product.is_digital),
    },
  })

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
        alt_text: img.alt_text || '',
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
      const productTags = product.tags.map(tag => tag.name)
      const currentTagsStr = JSON.stringify(tags)
      const productTagsStr = JSON.stringify(productTags)
      
      if (currentTagsStr !== productTagsStr) {
        setTags(productTags)
      }
    }
  }, [product, isSubmittingRef]) // Add isSubmittingRef as dependency

  const watchedFields = watch()

  const handleImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (!files) return

    Array.from(files).forEach((file) => {
      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result && typeof e.target.result === 'string') {
          addImage(e.target.result)
        }
      }
      reader.readAsDataURL(file)
    })
  }

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

  const addImage = (url: string, altText = '') => {
    if (images.length >= 10) {
      toast.error('Maximum 10 images allowed')
      return
    }

    const newImage: ProductImage = {
      url,
      alt_text: altText,
      position: images.length,
    }

    console.log('addImage: Adding new image:', newImage)
    setImages(prev => {
      const newImages = [...prev, newImage]
      console.log('addImage: New images array:', newImages)
      
      // Check if actually changed using intelligent comparison
      const actuallyChanged = compareImages(newImages, originalImages)
      setImagesChanged(actuallyChanged)
      console.log('addImage: Images changed?', actuallyChanged)
      
      return newImages
    })
  }

  const removeImage = (index: number) => {
    console.log('removeImage: Removing image at index:', index)
    setImages(prev => {
      const newImages = prev.filter((_, i) => i !== index).map((img, i) => ({ ...img, position: i }))
      
      // Check if images actually changed using intelligent comparison
      const actuallyChanged = compareImages(newImages, originalImages)
      setImagesChanged(actuallyChanged)
      console.log('removeImage: Images changed?', actuallyChanged)
      console.log('removeImage: New images count:', newImages.length)
      console.log('removeImage: Original images count:', originalImages.length)
      
      return newImages
    })
  }

  // Function to update image alt text
  const updateImageAltText = (index: number, newAltText: string) => {
    console.log('updateImageAltText: Updating alt text at index:', index, 'to:', newAltText)
    setImages(prev => {
      const newImages = [...prev]
      newImages[index] = { ...newImages[index], alt_text: newAltText }
      
      // Check if actually changed using intelligent comparison
      const actuallyChanged = compareImages(newImages, originalImages)
      setImagesChanged(actuallyChanged)
      console.log('updateImageAltText: Images changed?', actuallyChanged)
      
      return newImages
    })
  }

  // Function to reorder images
  const moveImage = (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex) return
    
    console.log('moveImage: Moving image from', fromIndex, 'to', toIndex)
    setImages(prev => {
      const newImages = [...prev]
      const [movedImage] = newImages.splice(fromIndex, 1)
      newImages.splice(toIndex, 0, movedImage)
      
      // Update positions
      const reorderedImages = newImages.map((img, i) => ({ ...img, position: i }))
      
      // Check if actually changed using intelligent comparison
      const actuallyChanged = compareImages(reorderedImages, originalImages)
      setImagesChanged(actuallyChanged)
      console.log('moveImage: Images changed?', actuallyChanged)
      
      return reorderedImages
    })
  }

  // Handle drag and drop reordering
  const handleDragEnd = (result: DropResult) => {
    const { destination, source } = result
    
    // Dropped outside the list
    if (!destination) {
      return
    }
    
    // Dropped in the same position
    if (destination.index === source.index) {
      return
    }
    
    console.log('handleDragEnd: Moving from', source.index, 'to', destination.index)
    moveImage(source.index, destination.index)
  }

  const addImageFromUrl = () => {
    if (!imageUrl.trim()) return
    
    try {
      new URL(imageUrl) // Validate URL
      addImage(imageUrl)
      setImageUrl('')
    } catch {
      toast.error('Please enter a valid image URL')
    }
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
      console.log('=== DEBUG: onSubmit ===')
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
          alt_text: img.alt_text || '',
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
      alt_text: img.alt_text || '',
      position: index,
    })) || []
    
    setImages(productImages)
    setOriginalImages(productImages)
    setImagesChanged(false) // Reset flag
    setTags(product.tags?.map(tag => tag.name) || [])
    setNewTag('')
    setImageUrl('')
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Basic Information */}
      <Card>
        <CardHeader>
          <CardTitle>Basic Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="name">Product Name</Label>
              <Input
                id="name"
                {...register('name')}
                className={errors.name ? 'border-red-500' : ''}
              />
              {errors.name && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.name.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="sku">SKU</Label>
              <Input
                id="sku"
                {...register('sku')}
                className={errors.sku ? 'border-red-500' : ''}
              />
              {errors.sku && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.sku.message}
                </p>
              )}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              rows={4}
              {...register('description')}
              className={errors.description ? 'border-red-500' : ''}
            />
            {errors.description && (
              <p className="text-sm text-red-600 flex items-center gap-1">
                <AlertCircle className="h-4 w-4" />
                {errors.description.message}
              </p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="short_description">Short Description</Label>
            <Textarea
              id="short_description"
              rows={2}
              {...register('short_description')}
            />
          </div>
        </CardContent>
      </Card>

      {/* Images */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ImageIcon className="h-5 w-5" />
            Product Images ({images.length}/10)
            {imagesChanged && (
              <Badge variant="outline" className="text-orange-600 border-orange-600">
                Modified
              </Badge>
            )}
          </CardTitle>
          <div className="text-sm text-gray-600 dark:text-gray-400">
            <p>The first image will be used as the primary product image. You can:</p>
            <ul className="list-disc list-inside mt-1 space-y-1">
              <li>Drag images to reorder them</li>
              <li>Use arrow buttons to move images up/down</li>
              <li>Click the edit icon to add alt text</li>
              <li>Click the X button to remove images</li>
            </ul>
          </div>
          {/* Debug Info */}
          <div className="text-xs text-gray-500 space-y-1">
            <div>Original images: {originalImages.length}</div>
            <div>Current images: {images.length}</div>
            <div>Changes detected: {imagesChanged ? 'Yes' : 'No'}</div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Upload Controls */}
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <Label htmlFor="image-upload" className="cursor-pointer">
                <div className="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center hover:border-gray-400 transition-colors">
                  <Upload className="h-8 w-8 mx-auto text-gray-400 mb-2" />
                  <p className="text-sm text-gray-600">Click to upload images</p>
                  <p className="text-xs text-gray-500">PNG, JPG up to 5MB each</p>
                </div>
                <Input
                  id="image-upload"
                  type="file"
                  multiple
                  accept="image/*"
                  onChange={handleImageUpload}
                  className="sr-only"
                />
              </Label>
            </div>

            <div className="flex-1">
              <div className="flex gap-2">
                <Input
                  placeholder="Enter image URL"
                  value={imageUrl}
                  onChange={(e) => setImageUrl(e.target.value)}
                  onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addImageFromUrl())}
                />
                <Button
                  type="button"
                  variant="outline"
                  onClick={addImageFromUrl}
                  disabled={!imageUrl.trim()}
                >
                  <Link className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          {/* Images Grid */}
          {images.length > 0 && (
            <DragDropContext onDragEnd={handleDragEnd}>
              <Droppable droppableId="images-list" direction="horizontal">
                {(provided, snapshot) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                    className={`grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 ${
                      snapshot.isDraggingOver ? 'bg-blue-50 dark:bg-blue-950/20 rounded-lg p-2' : ''
                    }`}
                  >
                    {images.map((image, index) => (
                      <Draggable
                        key={`${image.url}-${index}`}
                        draggableId={`image-${index}`}
                        index={index}
                      >
                        {(provided, snapshot) => (
                          <div
                            ref={provided.innerRef}
                            {...provided.draggableProps}
                            className={`relative group ${
                              snapshot.isDragging ? 'rotate-6 scale-105 shadow-lg z-50' : ''
                            }`}
                          >
                            <div className="aspect-square rounded-lg border overflow-hidden bg-gray-100">
                              {/* Drag Handle */}
                              <div
                                {...provided.dragHandleProps}
                                className="absolute top-1 left-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity cursor-grab active:cursor-grabbing bg-black/50 rounded p-1"
                                title="Drag to reorder"
                              >
                                <GripVertical className="h-3 w-3 text-white" />
                              </div>

                              <Image
                                src={image.url}
                                alt={image.alt_text || `Product image ${index + 1}`}
                                fill
                                className="object-cover"
                                unoptimized={true}
                                onError={(e) => {
                                  const target = e.target as HTMLImageElement;
                                  // Set a simple placeholder that won't cause additional requests
                                  target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iI0Y3RjhGOSIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBkb21pbmFudC1iYXNlbGluZT0ibWlkZGxlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic3lzdGVtLXVpIiBmb250LXNpemU9IjE2IiBmaWxsPSIjOUNBM0FGIj5JbWFnZTwvdGV4dD48L3N2Zz4=';
                                  // Prevent further error events
                                  target.onerror = null;
                                }}
                              />
                              
                              {/* Primary Badge */}
                              {index === 0 && (
                                <Badge className="absolute top-2 left-2 bg-blue-600">
                                  <Star className="h-3 w-3 mr-1" />
                                  Primary
                                </Badge>
                              )}
                              
                              {/* Image Controls */}
                              <div className="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-30 transition-all duration-200 flex items-center justify-center">
                                <div className="opacity-0 group-hover:opacity-100 flex gap-1 transition-opacity">
                                  {/* Move Up */}
                                  {index > 0 && (
                                    <Button
                                      type="button"
                                      size="sm"
                                      variant="secondary"
                                      className="h-6 w-6 p-0"
                                      onClick={() => moveImage(index, index - 1)}
                                      title="Move up"
                                    >
                                      <ArrowUp className="h-3 w-3" />
                                    </Button>
                                  )}
                                  
                                  {/* Move Down */}
                                  {index < images.length - 1 && (
                                    <Button
                                      type="button"
                                      size="sm"
                                      variant="secondary"
                                      className="h-6 w-6 p-0"
                                      onClick={() => moveImage(index, index + 1)}
                                      title="Move down"
                                    >
                                      <ArrowDown className="h-3 w-3" />
                                    </Button>
                                  )}
                                  
                                  {/* Edit Alt Text */}
                                  <Button
                                    type="button"
                                    size="sm"
                                    variant="secondary"
                                    className="h-6 w-6 p-0"
                                    onClick={() => {
                                      const newAltText = prompt('Enter alt text for this image:', image.alt_text || '')
                                      if (newAltText !== null) {
                                        updateImageAltText(index, newAltText)
                                      }
                                    }}
                                    title="Edit alt text"
                                  >
                                    <Edit2 className="h-3 w-3" />
                                  </Button>
                                </div>
                              </div>
                              
                              {/* Remove Button */}
                              <Button
                                type="button"
                                size="sm"
                                variant="destructive"
                                className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity h-6 w-6 p-0"
                                onClick={() => removeImage(index)}
                                title="Remove image"
                              >
                                <X className="h-3 w-3" />
                              </Button>
                            </div>
                            
                            {/* Image Info */}
                            <div className="mt-1 text-xs text-gray-500 text-center">
                              <div className="truncate">Position: {index + 1}</div>
                              {image.alt_text && (
                                <div className="truncate" title={image.alt_text}>
                                  Alt: {image.alt_text}
                                </div>
                              )}
                            </div>
                          </div>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                  </div>
                )}
              </Droppable>
            </DragDropContext>
          )}

          {images.length === 0 && (
            <div className="text-center py-8 text-gray-500">
              <ImageIcon className="h-12 w-12 mx-auto mb-4 opacity-50" />
              <p>No images uploaded yet</p>
              <p className="text-sm">Upload at least one image for your product</p>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Pricing & Inventory */}
      <Card>
        <CardHeader>
          <CardTitle>Pricing & Inventory</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label htmlFor="price">Price</Label>
              <Input
                id="price"
                type="number"
                step="0.01"
                {...register('price', { valueAsNumber: true })}
                className={errors.price ? 'border-red-500' : ''}
              />
              {errors.price && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.price.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="compare_price">Compare Price</Label>
              <Input
                id="compare_price"
                type="number"
                step="0.01"
                {...register('compare_price', { valueAsNumber: true })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="cost_price">Cost Price</Label>
              <Input
                id="cost_price"
                type="number"
                step="0.01"
                {...register('cost_price', { valueAsNumber: true })}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="stock">Stock Quantity</Label>
              <Input
                id="stock"
                type="number"
                {...register('stock', { valueAsNumber: true })}
                className={errors.stock ? 'border-red-500' : ''}
              />
              {errors.stock && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.stock.message}
                </p>
              )}
            </div>

            <div className="space-y-2">
              <Label htmlFor="category_id">Category ID</Label>
              <Input
                id="category_id"
                {...register('category_id')}
                placeholder="Enter category ID"
                className={errors.category_id ? 'border-red-500' : ''}
              />
              {errors.category_id && (
                <p className="text-sm text-red-600 flex items-center gap-1">
                  <AlertCircle className="h-4 w-4" />
                  {errors.category_id.message}
                </p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Product Details */}
      <Card>
        <CardHeader>
          <CardTitle>Product Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="weight">Weight (kg)</Label>
              <Input
                id="weight"
                type="number"
                step="0.01"
                {...register('weight', { valueAsNumber: true })}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="status">Status</Label>
              <select
                {...register('status')}
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
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
              className="h-4 w-4 rounded border-gray-300"
            />
            <Label htmlFor="is_digital">Digital Product</Label>
          </div>
        </CardContent>
      </Card>

      {/* Tags */}
      <Card>
        <CardHeader>
          <CardTitle>Tags ({tags.length}/10)</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex gap-2">
            <Input
              placeholder="Enter tag name"
              value={newTag}
              onChange={(e) => setNewTag(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addTag())}
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
        </CardContent>
      </Card>

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
