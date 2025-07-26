'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { FormSection } from '@/components/ui/form-section'
import { FormField } from '@/components/ui/form-field'
import { FormActions } from '@/components/ui/form-actions'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { productService } from '@/lib/services/products'
import { categoryService } from '@/lib/services/categories'
import { CreateProductRequest } from '@/lib/services/products'
import { useAuthStore } from '@/store/auth'
import { validatePriceInputs, cleanPriceInput } from '@/lib/utils/price'
import { X, Upload, Link2, Image as ImageIcon, GripVertical } from 'lucide-react'
import { DragDropContext, Droppable, Draggable, type DropResult } from '@hello-pangea/dnd'

interface AddProductFormProps {
  onSuccess?: () => void
  onCancel?: () => void
}

export function AddProductForm({ onSuccess, onCancel }: AddProductFormProps) {
  const [formData, setFormData] = useState<CreateProductRequest>({
    name: '',
    description: '',
    short_description: '',
    sku: '',

    // SEO and Metadata
    slug: '', // Will be auto-generated from name if empty
    meta_title: '',
    meta_description: '',
    keywords: '',
    featured: false,
    visibility: 'visible',

    // Pricing
    price: 0,
    compare_price: undefined,
    cost_price: undefined,

    // Sale Pricing
    sale_price: undefined,
    sale_start_date: undefined,
    sale_end_date: undefined,

    // Inventory
    stock: 0,
    low_stock_threshold: 5,
    track_quantity: true,
    allow_backorder: false,

    // Shipping and Tax
    requires_shipping: true,
    shipping_class: 'standard',
    tax_class: 'standard',
    country_of_origin: '',

    // Categorization - Backend uses ProductCategory many-to-many internally
    // but accepts category_id for primary category assignment
    category_id: '',
    brand_id: undefined,

    // Product Type
    product_type: 'simple',
    is_digital: false,

    // Physical Properties
    weight: undefined,
    dimensions: {
      length: 0,
      width: 0,
      height: 0,
    },

    tags: [],
    status: 'draft',
  } as CreateProductRequest)
  
  const [tagInput, setTagInput] = useState('')
  const [imageFiles, setImageFiles] = useState<File[]>([])
  const [imagePreviewUrls, setImagePreviewUrls] = useState<string[]>([])
  const [imageUrlInput, setImageUrlInput] = useState('')
  const [isDragOver, setIsDragOver] = useState(false)
  const queryClient = useQueryClient()
  const { user, token } = useAuthStore()
  
  // Debug auth state
  console.log('Auth state:', { 
    user: user?.role, 
    hasToken: !!token,
    token: token?.substring(0, 20) + '...',
    localStorage: typeof window !== 'undefined' ? localStorage.getItem('auth_token')?.substring(0, 20) + '...' : 'N/A'
  })
  
  // Fetch categories
  const { data: categories = [], isLoading: categoriesLoading, error: categoriesError } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryService.getCategories(),
  })
  
  console.log('Categories:', { 
    count: categories.length, 
    categories: categories.map(c => ({ id: c.id, name: c.name })),
    loading: categoriesLoading,
    error: categoriesError
  })

  // Create product mutation
  const createProductMutation = useMutation({
    mutationFn: async (data: CreateProductRequest) => {
      console.log('API call starting with data:', data)
      console.log('Current auth token:', token)
      console.log('User role:', user?.role)
      
      try {
        const result = await productService.createProduct(data)
        console.log('API call successful:', result)
        return result
      } catch (err) {
        console.error('API call error caught:', err)
        throw err
      }
    },
    onSuccess: (response) => {
      console.log('Product created successfully:', response)
      toast.success('Product created successfully!')
      queryClient.invalidateQueries({ queryKey: ['products'] })
      // Reset form
      setFormData({
        name: '',
        description: '',
        short_description: '',
        sku: '',

        // SEO and Metadata
        slug: '',
        meta_title: '',
        meta_description: '',
        keywords: '',
        featured: false,
        visibility: 'visible',

        // Pricing
        price: 0,
        compare_price: undefined,
        cost_price: undefined,

        // Sale Pricing
        sale_price: undefined,
        sale_start_date: undefined,
        sale_end_date: undefined,

        // Inventory
        stock: 0,
        low_stock_threshold: 5,
        track_quantity: true,
        allow_backorder: false,

        // Shipping and Tax
        requires_shipping: true,
        shipping_class: 'standard',
        tax_class: 'standard',
        country_of_origin: '',

        // Categorization
        category_id: '',
        brand_id: undefined,

        // Product Type
        product_type: 'simple',
        is_digital: false,

        // Physical Properties
        weight: undefined,
        dimensions: {
          length: 0,
          width: 0,
          height: 0,
        },

        tags: [],
        status: 'draft',
      } as CreateProductRequest)
      setImageFiles([])
      setImagePreviewUrls([])
      setImageUrlInput('')
      onSuccess?.()
    },
    onError: (error: any) => {
      console.error('Product creation failed:', error)
      console.error('Error type:', typeof error)
      console.error('Error constructor:', error.constructor.name)
      console.error('Error keys:', Object.keys(error))
      console.error('Full error object:', JSON.stringify(error, null, 2))
      
      // More detailed response logging
      if (error.response) {
        console.error('Response status:', error.response.status)
        console.error('Response data:', error.response.data)
        console.error('Response headers:', error.response.headers)
      }
      
      console.error('Error details:', {
        message: error.message,
        response: error.response,
        status: error.status,
        data: error.data,
        code: error.code,
        config: error.config,
      })
      
      let errorMessage = 'Failed to create product'
      
      if (error.response?.data?.message) {
        errorMessage = error.response.data.message
      } else if (error.response?.data?.error) {
        errorMessage = error.response.data.error
      } else if (error.message) {
        errorMessage = error.message
      }
      
      toast.error(errorMessage)
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    
    console.log('Form data being submitted:', formData)
    
    // Basic validation
    if (!formData.name.trim()) {
      toast.error('Product name is required')
      return
    }
    
    if (!formData.sku.trim()) {
      toast.error('SKU is required')
      return
    }
    
    if (!formData.description.trim()) {
      toast.error('Description is required')
      return
    }
    
    if (!formData.category_id) {
      toast.error('Please select a category')
      return
    }
    
    // Price validation using utility
    const priceErrors = validatePriceInputs(
      formData.price, 
      formData.compare_price, 
      formData.cost_price
    )
    
    if (priceErrors.length > 0) {
      toast.error(priceErrors[0])
      return
    }
    
    if (formData.stock < 0) {
      toast.error('Stock cannot be negative')
      return
    }
    
    // Clean up the data before sending
    const cleanData = {
      name: formData.name.trim(),
      description: formData.description.trim(),
      short_description: formData.short_description?.trim() || undefined,
      sku: formData.sku.trim(),
      price: formData.price,
      compare_price: cleanPriceInput(formData.compare_price),
      cost_price: cleanPriceInput(formData.cost_price),
      stock: formData.stock,
      category_id: formData.category_id,
      is_digital: formData.is_digital,
      weight: formData.weight || undefined,
      dimensions: formData.dimensions && (formData.dimensions.length > 0 || formData.dimensions.width > 0 || formData.dimensions.height > 0) 
        ? formData.dimensions 
        : undefined,
      tags: formData.tags || [],
      status: 'active', // Set as active instead of draft
      images: imagePreviewUrls.length > 0 ? imagePreviewUrls.map((url, index) => ({
        url: url.startsWith('data:') 
          ? `https://via.placeholder.com/400x400?text=${encodeURIComponent(formData.name.trim())}-${index + 1}`
          : url,
        alt_text: `${formData.name.trim()} - Image ${index + 1}`,
        position: index
      })) : undefined,
    }
    
    console.log('Clean data being sent to API:', cleanData)
    createProductMutation.mutate(cleanData)
  }

  const handleInputChange = (field: keyof CreateProductRequest, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const addTag = () => {
    if (tagInput.trim() && !formData.tags?.includes(tagInput.trim())) {
      handleInputChange('tags', [...(formData.tags || []), tagInput.trim()])
      setTagInput('')
    }
  }

  const removeTag = (tagToRemove: string) => {
    handleInputChange('tags', formData.tags?.filter(tag => tag !== tagToRemove) || [])
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      addTag()
    }
  }

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || [])
    addImageFiles(files)
  }

  const addImageFiles = (files: File[]) => {
    if (files.length === 0) return
    
    // Check total limit (max 10 images)
    const totalImages = imageFiles.length + imagePreviewUrls.length + files.length
    if (totalImages > 10) {
      toast.error('Maximum 10 images allowed')
      return
    }

    setImageFiles(prev => [...prev, ...files])
    
    // Create preview URLs
    files.forEach(file => {
      const reader = new FileReader()
      reader.onload = (e) => {
        const result = e.target?.result
        if (result && typeof result === 'string') {
          setImagePreviewUrls(prev => [...prev, result])
        }
      }
      reader.readAsDataURL(file)
    })
  }

  const addImageUrl = () => {
    const url = imageUrlInput.trim()
    if (!url) return
    
    // Simple URL validation
    try {
      new URL(url)
    } catch {
      toast.error('Please enter a valid URL')
      return
    }
    
    // Check total limit
    const totalImages = imageFiles.length + imagePreviewUrls.length + 1
    if (totalImages > 10) {
      toast.error('Maximum 10 images allowed')
      return
    }
    
    setImagePreviewUrls(prev => [...prev, url])
    setImageUrlInput('')
  }

  const removeImage = (index: number) => {
    setImageFiles(prev => prev.filter((_, i) => i !== index))
    setImagePreviewUrls(prev => prev.filter((_, i) => i !== index))
  }

  const moveImage = (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex) return
    
    setImagePreviewUrls(prev => {
      const newUrls = [...prev]
      const [movedUrl] = newUrls.splice(fromIndex, 1)
      newUrls.splice(toIndex, 0, movedUrl)
      return newUrls
    })
    
    setImageFiles(prev => {
      const newFiles = [...prev]
      if (fromIndex < newFiles.length && toIndex < newFiles.length) {
        const [movedFile] = newFiles.splice(fromIndex, 1)
        newFiles.splice(toIndex, 0, movedFile)
      }
      return newFiles
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

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
    
    const files = Array.from(e.dataTransfer.files).filter(file => 
      file.type.startsWith('image/')
    )
    
    if (files.length === 0) {
      toast.error('Please drop image files only')
      return
    }
    
    addImageFiles(files)
  }

  return (
    <div className="bihub-admin-form bg-transparent">
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Product Information */}
        <div className="lg:col-span-2 space-y-6">
          <FormSection title="Product Information">
            <FormField
              label="Product Name"
              required
            >
              <Input
                value={formData.name}
                onChange={(e) => handleInputChange('name', e.target.value)}
                placeholder="Enter product name"
                required
              />
            </FormField>

            <FormField
              label="SKU"
              required
            >
              <Input
                value={formData.sku}
                onChange={(e) => handleInputChange('sku', e.target.value)}
                placeholder="Enter product SKU"
                required
              />
            </FormField>

            <FormField
              label="Description"
              required
            >
              <Textarea
                value={formData.description}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleInputChange('description', e.target.value)}
                placeholder="Enter product description"
                rows={4}
                required
              />
            </FormField>

            <FormField
              label="Short Description"
            >
              <Textarea
                value={formData.short_description || ''}
                onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => handleInputChange('short_description', e.target.value)}
                placeholder="Enter short description"
                rows={2}
              />
            </FormField>
          </FormSection>

          {/* SEO & Metadata */}
          <FormSection title="SEO & Metadata">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField label="URL Slug">
                <Input
                  value={(formData as any).slug || ''}
                  onChange={(e) => handleInputChange('slug', e.target.value)}
                  placeholder="product-url-slug"
                />
              </FormField>

              <FormField label="Meta Title">
                <Input
                  value={(formData as any).meta_title || ''}
                  onChange={(e) => handleInputChange('meta_title', e.target.value)}
                  placeholder="SEO title for search engines"
                />
              </FormField>
            </div>

            <FormField label="Meta Description">
              <Textarea
                value={(formData as any).meta_description || ''}
                onChange={(e) => handleInputChange('meta_description', e.target.value)}
                placeholder="SEO description for search engines"
                rows={2}
              />
            </FormField>

            <FormField label="Keywords">
              <Input
                value={(formData as any).keywords || ''}
                onChange={(e) => handleInputChange('keywords', e.target.value)}
                placeholder="keyword1, keyword2, keyword3"
              />
            </FormField>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField label="Featured Product">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="featured"
                    checked={(formData as any).featured || false}
                    onChange={(e) => handleInputChange('featured', e.target.checked)}
                    className="h-4 w-4 shrink-0 rounded-sm border border-white/20 ring-offset-black focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#ff9000] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-[#ff9000] checked:text-white checked:border-[#ff9000]"
                  />
                  <span className="text-sm text-gray-300">Mark as featured</span>
                </div>
              </FormField>

              <FormField label="Visibility">
                <select
                  value={(formData as any).visibility || 'visible'}
                  onChange={(e) => handleInputChange('visibility', e.target.value)}
                  className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
                >
                  <option value="visible" className="bg-gray-800 text-white">Visible</option>
                  <option value="hidden" className="bg-gray-800 text-white">Hidden</option>
                  <option value="private" className="bg-gray-800 text-white">Private</option>
                </select>
              </FormField>
            </div>
          </FormSection>

          {/* Pricing */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Pricing</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {/* Price */}
                <div>
                  <Label htmlFor="price" className="text-sm font-medium text-white">Regular Price *</Label>
                  <Input
                    id="price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.price}
                    onChange={(e) => handleInputChange('price', parseFloat(e.target.value) || 0)}
                    placeholder="0.00"
                    required
                  />
                </div>

                {/* Sale Price */}
                <div>
                  <Label htmlFor="sale_price" className="text-sm font-medium text-white">Sale Price</Label>
                  <Input
                    id="sale_price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={(formData as any).sale_price || ''}
                    onChange={(e) => handleInputChange('sale_price', parseFloat(e.target.value) || undefined)}
                    placeholder="0.00"
                  />
                </div>

                {/* Compare Price */}
                <div>
                  <Label htmlFor="compare_price" className="text-sm font-medium text-white">Compare Price</Label>
                  <Input
                    id="compare_price"
                    type="number"
                    step="0.01"
                    min="0"
                    value={formData.compare_price || ''}
                    onChange={(e) => handleInputChange('compare_price', parseFloat(e.target.value) || undefined)}
                    placeholder="0.00"
                  />
                </div>
              </div>

              {/* Sale Date Range */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="sale_start_date" className="text-sm font-medium text-white">Sale Start Date</Label>
                  <Input
                    id="sale_start_date"
                    type="datetime-local"
                    value={(formData as any).sale_start_date || ''}
                    onChange={(e) => handleInputChange('sale_start_date', e.target.value || undefined)}
                  />
                </div>

                <div>
                  <Label htmlFor="sale_end_date" className="text-sm font-medium text-white">Sale End Date</Label>
                  <Input
                    id="sale_end_date"
                    type="datetime-local"
                    value={(formData as any).sale_end_date || ''}
                    onChange={(e) => handleInputChange('sale_end_date', e.target.value || undefined)}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Enhanced Inventory */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Inventory Management</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {/* Stock */}
                <div>
                  <Label htmlFor="stock" className="text-sm font-medium text-white">Stock Quantity *</Label>
                  <Input
                    id="stock"
                    type="number"
                    min="0"
                    value={formData.stock}
                    onChange={(e) => handleInputChange('stock', parseInt(e.target.value) || 0)}
                    placeholder="0"
                    required
                  />
                </div>

                {/* Low Stock Threshold */}
                <div>
                  <Label htmlFor="low_stock_threshold" className="text-sm font-medium text-white">Low Stock Threshold</Label>
                  <Input
                    id="low_stock_threshold"
                    type="number"
                    min="0"
                    value={(formData as any).low_stock_threshold || 5}
                    onChange={(e) => handleInputChange('low_stock_threshold', parseInt(e.target.value) || 5)}
                    placeholder="5"
                  />
                </div>

                {/* Weight (if not digital) */}
                {!formData.is_digital && (
                  <div>
                    <Label htmlFor="weight" className="text-sm font-medium text-white">Weight (kg)</Label>
                    <Input
                      id="weight"
                      type="number"
                      step="0.01"
                      min="0"
                      value={formData.weight || ''}
                      onChange={(e) => handleInputChange('weight', parseFloat(e.target.value) || undefined)}
                      placeholder="0.00"
                    />
                  </div>
                )}
              </div>

              {/* Inventory Options */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="track_quantity"
                    checked={(formData as any).track_quantity !== false}
                    onChange={(e) => handleInputChange('track_quantity', e.target.checked)}
                    className="h-4 w-4 shrink-0 rounded-sm border border-white/20 ring-offset-black focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#ff9000] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-[#ff9000] checked:text-white checked:border-[#ff9000]"
                  />
                  <Label htmlFor="track_quantity" className="text-sm font-medium text-white">Track Quantity</Label>
                </div>

                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="allow_backorder"
                    checked={(formData as any).allow_backorder || false}
                    onChange={(e) => handleInputChange('allow_backorder', e.target.checked)}
                    className="h-4 w-4 shrink-0 rounded-sm border border-white/20 ring-offset-black focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#ff9000] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-[#ff9000] checked:text-white checked:border-[#ff9000]"
                  />
                  <Label htmlFor="allow_backorder" className="text-sm font-medium text-white">Allow Backorders</Label>
                </div>
              </div>

              {/* Dimensions (if not digital) */}
              {!formData.is_digital && (
                <div>
                  <Label className="text-sm font-medium text-white">Dimensions (cm)</Label>
                  <div className="grid grid-cols-3 gap-2">
                    <div>
                      <Label htmlFor="length" className="text-xs text-gray-300">Length</Label>
                      <Input
                        id="length"
                        type="number"
                        step="0.01"
                        min="0"
                        value={formData.dimensions?.length || ''}
                        onChange={(e) => {
                          const value = parseFloat(e.target.value) || 0
                          handleInputChange('dimensions', {
                            ...formData.dimensions,
                            length: value
                          })
                        }}
                        placeholder="0"
                      />
                    </div>
                    <div>
                      <Label htmlFor="width" className="text-xs text-gray-300">Width</Label>
                      <Input
                        id="width"
                        type="number"
                        step="0.01"
                        min="0"
                        value={formData.dimensions?.width || ''}
                        onChange={(e) => {
                          const value = parseFloat(e.target.value) || 0
                          handleInputChange('dimensions', {
                            ...formData.dimensions,
                            length: formData.dimensions?.length || 0,
                            width: value,
                            height: formData.dimensions?.height || 0,
                          })
                        }}
                        placeholder="0"
                      />
                    </div>
                    <div>
                      <Label htmlFor="height" className="text-xs text-gray-300">Height</Label>
                      <Input
                        id="height"
                        type="number"
                        step="0.01"
                        min="0"
                        value={formData.dimensions?.height || ''}
                        onChange={(e) => {
                          const value = parseFloat(e.target.value) || 0
                          handleInputChange('dimensions', {
                            ...formData.dimensions,
                            length: formData.dimensions?.length || 0,
                            width: formData.dimensions?.width || 0,
                            height: value,
                          })
                        }}
                        placeholder="0"
                      />
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Shipping & Tax */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Shipping & Tax</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Shipping Class */}
                <div>
                  <Label htmlFor="shipping_class" className="text-sm font-medium text-white">Shipping Class</Label>
                  <select
                    id="shipping_class"
                    value={(formData as any).shipping_class || 'standard'}
                    onChange={(e) => handleInputChange('shipping_class', e.target.value)}
                    className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
                  >
                    <option value="standard" className="bg-gray-800 text-white">Standard</option>
                    <option value="express" className="bg-gray-800 text-white">Express</option>
                    <option value="overnight" className="bg-gray-800 text-white">Overnight</option>
                    <option value="free" className="bg-gray-800 text-white">Free Shipping</option>
                  </select>
                </div>

                {/* Tax Class */}
                <div>
                  <Label htmlFor="tax_class" className="text-sm font-medium text-white">Tax Class</Label>
                  <select
                    id="tax_class"
                    value={(formData as any).tax_class || 'standard'}
                    onChange={(e) => handleInputChange('tax_class', e.target.value)}
                    className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
                  >
                    <option value="standard" className="bg-gray-800 text-white">Standard</option>
                    <option value="reduced" className="bg-gray-800 text-white">Reduced Rate</option>
                    <option value="zero" className="bg-gray-800 text-white">Zero Rate</option>
                    <option value="exempt" className="bg-gray-800 text-white">Tax Exempt</option>
                  </select>
                </div>
              </div>

              {/* Country of Origin */}
              <div>
                <Label htmlFor="country_of_origin" className="text-sm font-medium text-white">Country of Origin</Label>
                <Input
                  id="country_of_origin"
                  value={(formData as any).country_of_origin || ''}
                  onChange={(e) => handleInputChange('country_of_origin', e.target.value)}
                  placeholder="e.g., United States"
                />
              </div>

              {/* Requires Shipping */}
              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="requires_shipping"
                  checked={(formData as any).requires_shipping !== false}
                  onChange={(e) => handleInputChange('requires_shipping', e.target.checked)}
                  className="h-4 w-4 shrink-0 rounded-sm border border-white/20 ring-offset-black focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#ff9000] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-[#ff9000] checked:text-white checked:border-[#ff9000]"
                />
                <Label htmlFor="requires_shipping" className="text-sm font-medium text-white">Requires Shipping</Label>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Product Settings */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Product Settings</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Category */}
              <div>
                <Label htmlFor="category" className="text-sm font-medium text-white">Category *</Label>
                <select
                  id="category"
                  value={formData.category_id}
                  onChange={(e) => handleInputChange('category_id', e.target.value)}
                  className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
                  required
                >
                  <option value="" className="bg-gray-800 text-white">Select category</option>
                  {categories.map((category: any) => (
                    <option key={category.id} value={category.id} className="bg-gray-800 text-white">
                      {category.name}
                    </option>
                  ))}
                </select>
                {categoriesLoading && (
                  <p className="text-sm text-gray-400 mt-1">Loading categories...</p>
                )}
              </div>

              {/* Product Type */}
              <div>
                <Label htmlFor="product_type" className="text-sm font-medium text-white">Product Type</Label>
                <select
                  id="product_type"
                  value={(formData as any).product_type || 'simple'}
                  onChange={(e) => handleInputChange('product_type', e.target.value)}
                  className="flex h-12 w-full rounded-lg border border-gray-600/50 bg-gray-700/50 px-4 py-3 text-sm text-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#FF9000]/20 focus-visible:border-[#FF9000] disabled:cursor-not-allowed disabled:opacity-50 transition-colors duration-200"
                >
                  <option value="simple" className="bg-gray-800 text-white">Simple Product</option>
                  <option value="variable" className="bg-gray-800 text-white">Variable Product</option>
                  <option value="grouped" className="bg-gray-800 text-white">Grouped Product</option>
                  <option value="external" className="bg-gray-800 text-white">External Product</option>
                  <option value="digital" className="bg-gray-800 text-white">Digital Product</option>
                </select>
              </div>

              {/* Digital Product */}
              <div className="flex items-center space-x-2">
                <input
                  type="checkbox"
                  id="is_digital"
                  checked={formData.is_digital}
                  onChange={(e) => handleInputChange('is_digital', e.target.checked)}
                  className="h-4 w-4 shrink-0 rounded-sm border border-white/20 ring-offset-black focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#ff9000] focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 checked:bg-[#ff9000] checked:text-white checked:border-[#ff9000]"
                />
                <Label htmlFor="is_digital" className="text-sm font-medium text-white">Digital Product</Label>
              </div>
            </CardContent>
          </Card>

          {/* Images */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Product Images</CardTitle>
              <p className="text-sm text-gray-400">Add 1-10 images. First image will be the main product image.</p>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Drag & Drop Area */}
              <div
                className={`border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200 ${
                  isDragOver
                    ? 'border-[#FF9000] bg-[#FF9000]/10 scale-105'
                    : 'border-gray-600 hover:border-gray-500 hover:bg-gray-800/30'
                }`}
                onDragOver={handleDragOver}
                onDragLeave={handleDragLeave}
                onDrop={handleDrop}
              >
                <div className="space-y-4">
                  <div className={`transition-colors ${isDragOver ? 'text-[#FF9000]' : 'text-gray-400'}`}>
                    <Upload className="mx-auto h-16 w-16" />
                  </div>
                  <div>
                    <p className={`text-xl font-medium transition-colors ${isDragOver ? 'text-[#FF9000]' : 'text-white'}`}>
                      {isDragOver ? 'Drop images here!' : 'Upload Product Images'}
                    </p>
                    <p className="text-sm text-gray-400 mt-1">
                      Drag and drop images here, or click to browse
                    </p>
                    <p className="text-xs text-gray-500 mt-1">
                      Supports: JPG, PNG, WEBP • Max 10 images • Up to 5MB each
                    </p>
                  </div>
                  <Input
                    type="file"
                    multiple
                    accept="image/*"
                    onChange={handleImageUpload}
                    className="hidden"
                    id="file-upload"
                  />
                  <label
                    htmlFor="file-upload"
                    className="inline-flex items-center gap-2 px-6 py-3 border border-gray-600 rounded-lg shadow-sm text-sm font-medium text-gray-300 bg-gray-800/50 hover:bg-gray-700/50 hover:border-[#FF9000] hover:text-white cursor-pointer transition-colors"
                  >
                    <Upload className="h-4 w-4" />
                    Choose Files
                  </label>
                </div>
              </div>

              {/* URL Input */}
              <div className="flex gap-2">
                <div className="flex-1 relative">
                  <Link2 className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                  <Input
                    value={imageUrlInput}
                    onChange={(e) => setImageUrlInput(e.target.value)}
                    placeholder="Or paste image URL here..."
                    className="pl-10"
                    onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), addImageUrl())}
                  />
                </div>
                <Button type="button" onClick={addImageUrl} variant="outline" className="flex items-center gap-2">
                  <Link2 className="h-4 w-4" />
                  Add URL
                </Button>
              </div>
              
              {/* Image Previews */}
              {imagePreviewUrls.length > 0 && (
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                      <ImageIcon className="h-5 w-5 text-gray-400" />
                      <p className="text-sm font-medium text-white">Images ({imagePreviewUrls.length}/10)</p>
                    </div>
                    <div className="text-right">
                      <p className="text-xs text-gray-400">Use grip handle to reorder • First image is main</p>
                    </div>
                  </div>
                  
                  <DragDropContext onDragEnd={handleDragEnd}>
                    <Droppable droppableId="images-list" direction="horizontal">
                      {(provided, snapshot) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.droppableProps}
                          className={`grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 ${
                            snapshot.isDraggingOver ? 'bg-[#FF9000]/10 rounded-lg p-2' : ''
                          }`}
                        >
                          {imagePreviewUrls.map((url, index) => (
                            <Draggable
                              key={`image-${index}`}
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
                                  <div className="aspect-square relative overflow-hidden rounded-lg border-2 border-gray-600 bg-gray-800">
                                    {/* Drag Handle */}
                                    <div
                                      {...provided.dragHandleProps}
                                      className="absolute top-1 left-1 z-10 opacity-0 group-hover:opacity-100 transition-opacity cursor-grab active:cursor-grabbing bg-black/50 rounded p-1"
                                      title="Drag to reorder"
                                    >
                                      <GripVertical className="h-3 w-3 text-white" />
                                    </div>

                                    <img
                                      src={url}
                                      alt={`Preview ${index + 1}`}
                                      className="w-full h-full object-cover transition-transform group-hover:scale-105"
                                      onError={(e) => {
                                        const target = e.target as HTMLImageElement;
                                        target.src = '/placeholder-product.svg';
                                      }}
                                    />
                                    
                                    {/* Main Image Badge */}
                                    {index === 0 && (
                                      <div className="absolute top-2 left-2 bg-[#FF9000] text-white text-xs px-2 py-1 rounded-full font-medium shadow-lg">
                                        Main
                                      </div>
                                    )}
                                    
                                    {/* Position Number */}
                                    <div className="absolute bottom-2 left-2 bg-black bg-opacity-60 text-white text-xs px-2 py-1 rounded-full">
                                      {index + 1}
                                    </div>
                                    
                                    {/* Remove Button */}
                                    <button
                                      type="button"
                                      onClick={() => removeImage(index)}
                                      className="absolute top-2 right-2 bg-red-500 text-white rounded-full w-7 h-7 flex items-center justify-center text-sm hover:bg-red-600 opacity-0 group-hover:opacity-100 transition-all duration-200 shadow-lg hover:scale-110"
                                      title="Remove image"
                                    >
                                      <X className="h-4 w-4" />
                                    </button>
                                  </div>
                                  
                                  {/* Image Info */}
                                  <div className="mt-2 text-center">
                                    <p className="text-xs text-gray-400">
                                      {index === 0 ? 'Main Image' : `Image ${index + 1}`}
                                    </p>
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
                  
                  {/* Instructions */}
                  <div className="bg-gray-800/50 border border-gray-600/50 rounded-lg p-3">
                    <div className="flex items-start gap-2">
                      <ImageIcon className="h-4 w-4 text-[#FF9000] mt-0.5 flex-shrink-0" />
                      <div className="text-sm text-gray-300">
                        <p className="font-medium text-white">Image Tips:</p>
                        <ul className="mt-1 space-y-1 text-gray-400">
                          <li>• Drag the grip icon to reorder images</li>
                          <li>• First image will be the main product image</li>
                          <li>• Recommended size: 800x800px or larger</li>
                          <li>• Use high-quality images for better sales</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Tags */}
          <Card className="bg-gray-800/50 border-gray-700/50">
            <CardHeader>
              <CardTitle className="text-white">Tags</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex gap-2">
                <Input
                  value={tagInput}
                  onChange={(e) => setTagInput(e.target.value)}
                  onKeyDown={handleKeyDown}
                  placeholder="Add a tag"
                />
                <Button type="button" onClick={addTag} variant="outline">
                  Add
                </Button>
              </div>
              
              {formData.tags && formData.tags.length > 0 && (
                <div className="flex flex-wrap gap-2">
                  {formData.tags.map((tag) => (
                    <span
                      key={tag}
                      className="inline-flex items-center gap-1 px-2 py-1 bg-[#FF9000]/20 text-[#FF9000] text-sm rounded-md border border-[#FF9000]/30"
                    >
                      {tag}
                      <button
                        type="button"
                        onClick={() => removeTag(tag)}
                        className="hover:text-[#FF9000]/80"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </span>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-4 pt-6 border-t border-gray-700/50">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={createProductMutation.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          disabled={createProductMutation.isPending}
        >
          {createProductMutation.isPending ? 'Creating...' : 'Create Product'}
        </Button>
      </div>
    </form>
    </div>
  )
}
