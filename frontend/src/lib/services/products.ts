import { apiClient, buildQueryString, getPaginated } from '@/lib/api'
import { ProductFilters, Product, ProductResponse } from '@/types'

export interface ProductsParams extends ProductFilters {
  page?: number
  limit?: number
  search?: string
  category_id?: string
  min_price?: number
  max_price?: number
  rating?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'

  // Enhanced filters
  featured?: boolean
  on_sale?: boolean
  stock_status?: 'in_stock' | 'out_of_stock' | 'on_backorder'
  product_type?: string
  brand_id?: string
}

// Backend-compatible request interfaces
export interface CreateProductRequest {
  name: string
  description: string
  short_description: string
  sku: string

  // SEO and Metadata
  slug: string // Required in backend
  meta_title?: string
  meta_description?: string
  keywords?: string
  featured?: boolean
  visibility?: 'visible' | 'hidden' | 'private'

  // Pricing
  price: number
  compare_price?: number
  cost_price?: number

  // Sale Pricing
  sale_price?: number
  sale_start_date?: string
  sale_end_date?: string

  // Inventory
  stock: number
  low_stock_threshold?: number
  track_quantity?: boolean
  allow_backorder?: boolean

  // Physical Properties
  weight?: number
  dimensions?: {
    length: number
    width: number
    height: number
  }

  // Shipping and Tax
  requires_shipping?: boolean
  shipping_class?: string
  tax_class?: string
  country_of_origin?: string

  // Categorization
  category_id: string
  brand_id?: string

  // Content
  images?: Array<{
    url: string
    alt_text?: string
    position?: number
  }>
  tags?: string[]
  attributes?: Array<{
    attribute_id: string
    term_id?: string
    value: string
    position?: number
  }>
  variants?: Array<{
    sku: string
    price: number
    compare_price?: number
    cost_price?: number
    stock: number
    weight?: number
    dimensions?: {
      length: number
      width: number
      height: number
    }
    image?: string
    position?: number
    is_active?: boolean
    attributes?: Array<{
      attribute_id: string
      term_id: string
      value: string
    }>
  }>

  // Status and Type
  status?: 'active' | 'draft' | 'archived' | 'inactive'
  product_type?: 'simple' | 'variable' | 'grouped' | 'external' | 'digital'
  is_digital?: boolean
}

// Update request - all fields optional (matching backend UpdateProductRequest)
export interface UpdateProductRequest {
  name?: string
  description?: string
  short_description?: string

  // SEO and Metadata
  slug?: string
  meta_title?: string
  meta_description?: string
  keywords?: string
  featured?: boolean
  visibility?: 'visible' | 'hidden' | 'private'

  // Pricing
  price?: number
  compare_price?: number
  cost_price?: number

  // Sale Pricing
  sale_price?: number
  sale_start_date?: string
  sale_end_date?: string

  // Inventory
  stock?: number
  low_stock_threshold?: number
  track_quantity?: boolean
  allow_backorder?: boolean

  // Physical Properties
  weight?: number
  dimensions?: {
    length: number
    width: number
    height: number
  }

  // Shipping and Tax
  requires_shipping?: boolean
  shipping_class?: string
  tax_class?: string
  country_of_origin?: string

  // Categorization
  category_id?: string
  brand_id?: string

  // Content - arrays replace all existing items
  images?: Array<{
    url: string
    alt_text?: string
    position?: number
  }>
  tags?: string[]
  attributes?: Array<{
    attribute_id: string
    term_id?: string
    value: string
    position?: number
  }>
  variants?: Array<{
    sku: string
    price: number
    compare_price?: number
    cost_price?: number
    stock: number
    weight?: number
    dimensions?: {
      length: number
      width: number
      height: number
    }
    image?: string
    position?: number
    is_active?: boolean
    attributes?: Array<{
      attribute_id: string
      term_id: string
      value: string
    }>
  }>

  // Status and Type
  status?: 'active' | 'draft' | 'archived' | 'inactive'
  product_type?: 'simple' | 'variable' | 'grouped' | 'external' | 'digital'
  is_digital?: boolean
}

class ProductService {
  // Get all products with filters and pagination
  async getProducts(params: ProductsParams = {}): Promise<any> {
    try {
      // Build query parameters
      const queryParams: Record<string, any> = {}

      // Add pagination parameters
      if (params.page) queryParams.page = params.page
      if (params.limit) queryParams.limit = params.limit

      // Add search and filter parameters
      if (params.search) queryParams.search = params.search
      if (params.category_id) queryParams.category_id = params.category_id
      if (params.min_price) queryParams.min_price = params.min_price
      if (params.max_price) queryParams.max_price = params.max_price
      if (params.rating) queryParams.rating = params.rating
      if (params.sort_by) queryParams.sort_by = params.sort_by
      if (params.sort_order) queryParams.sort_order = params.sort_order
      if (params.featured !== undefined) queryParams.featured = params.featured
      if (params.on_sale !== undefined) queryParams.on_sale = params.on_sale
      if (params.stock_status) queryParams.stock_status = params.stock_status
      if (params.product_type) queryParams.product_type = params.product_type
      if (params.brand_id) queryParams.brand_id = params.brand_id

      // Build query string
      const queryString = Object.keys(queryParams).length > 0
        ? '?' + new URLSearchParams(queryParams).toString()
        : ''

      const response = await apiClient.get(`/products${queryString}`)

      // API client already returns the correct data structure
      return response
    } catch (error) {
      console.error('ProductService.getProducts error:', error)
      throw error
    }
  }

  // Get single product by ID
  async getProduct(id: string): Promise<Product> {
    const response = await apiClient.get(`/products/${id}`)
    return response as Product
  }

  // Get product by slug
  async getProductBySlug(slug: string): Promise<Product> {
    const response = await apiClient.get(`/products/slug/${slug}`)
    return response as Product
  }

  // Get featured products
  async getFeaturedProducts(limit = 8): Promise<{ data: Product[], pagination?: any }> {
    const response = await apiClient.get(`/products/featured?limit=${limit}`)
    return response || { data: [], pagination: null }
  }

  // Get related products
  async getRelatedProducts(productId: string, limit = 4): Promise<{ data: Product[], pagination?: any }> {
    const response = await apiClient.get(`/products/${productId}/related?limit=${limit}`)
    return response || { data: [], pagination: null }
  }

  // Get products by category
  async getProductsByCategory(categoryId: string, params: ProductsParams = {}): Promise<{ data: Product[], pagination?: any }> {
    // Build query parameters
    const queryParams: Record<string, any> = { category_id: categoryId }

    // Add pagination parameters
    if (params.page) queryParams.page = params.page
    if (params.limit) queryParams.limit = params.limit

    // Add other filter parameters
    if (params.search) queryParams.search = params.search
    if (params.min_price) queryParams.min_price = params.min_price
    if (params.max_price) queryParams.max_price = params.max_price
    if (params.rating) queryParams.rating = params.rating
    if (params.sort_by) queryParams.sort_by = params.sort_by
    if (params.sort_order) queryParams.sort_order = params.sort_order
    if (params.featured !== undefined) queryParams.featured = params.featured
    if (params.on_sale !== undefined) queryParams.on_sale = params.on_sale
    if (params.stock_status) queryParams.stock_status = params.stock_status
    if (params.product_type) queryParams.product_type = params.product_type
    if (params.brand_id) queryParams.brand_id = params.brand_id

    // Build query string
    const queryString = Object.keys(queryParams).length > 0
      ? '?' + new URLSearchParams(queryParams).toString()
      : ''

    const response = await apiClient.get(`/products/category/${categoryId}${queryString}`)
    return response as { data: Product[], pagination?: any }
  }

  // Search products
  async searchProducts(query: string, params: ProductsParams = {}): Promise<{ data: Product[], pagination?: any }> {
    // Build query parameters
    const queryParams: Record<string, any> = { q: query }

    // Add pagination parameters
    if (params.page) queryParams.page = params.page
    if (params.limit) queryParams.limit = params.limit

    // Add filter parameters
    if (params.category_id) queryParams.category_id = params.category_id
    if (params.min_price) queryParams.min_price = params.min_price
    if (params.max_price) queryParams.max_price = params.max_price
    if (params.rating) queryParams.rating = params.rating
    if (params.sort_by) queryParams.sort_by = params.sort_by
    if (params.sort_order) queryParams.sort_order = params.sort_order
    if (params.featured !== undefined) queryParams.featured = params.featured
    if (params.on_sale !== undefined) queryParams.on_sale = params.on_sale
    if (params.stock_status) queryParams.stock_status = params.stock_status
    if (params.product_type) queryParams.product_type = params.product_type
    if (params.brand_id) queryParams.brand_id = params.brand_id

    // Build query string
    const queryString = Object.keys(queryParams).length > 0
      ? '?' + new URLSearchParams(queryParams).toString()
      : ''

    const response = await apiClient.get(`/products/search${queryString}`)
    return response as { data: Product[], pagination?: any }
  }

  // Get product suggestions for autocomplete
  async getProductSuggestions(query: string, limit = 5): Promise<Product[]> {
    const response = await apiClient.get(`/products/suggestions?q=${query}&limit=${limit}`)
    return (response?.data || []) as Product[]
  }

  // Admin methods
  async getAdminProducts(params: ProductsParams = {}): Promise<{ data: Product[], pagination?: any }> {
    console.log('ProductService.getAdminProducts called with params:', params)
    try {
      // Build query parameters
      const queryParams: Record<string, any> = {}

      // Add pagination parameters
      if (params.page) queryParams.page = params.page
      if (params.limit) queryParams.limit = params.limit

      // Add search and filter parameters
      if (params.search) queryParams.search = params.search
      if (params.category_id) queryParams.category_id = params.category_id
      if (params.min_price) queryParams.min_price = params.min_price
      if (params.max_price) queryParams.max_price = params.max_price
      if (params.rating) queryParams.rating = params.rating
      if (params.sort_by) queryParams.sort_by = params.sort_by
      if (params.sort_order) queryParams.sort_order = params.sort_order
      if (params.featured !== undefined) queryParams.featured = params.featured
      if (params.on_sale !== undefined) queryParams.on_sale = params.on_sale
      if (params.stock_status) queryParams.stock_status = params.stock_status
      if (params.product_type) queryParams.product_type = params.product_type
      if (params.brand_id) queryParams.brand_id = params.brand_id

      // Build query string
      const queryString = Object.keys(queryParams).length > 0
        ? '?' + new URLSearchParams(queryParams).toString()
        : ''

      const response = await apiClient.get(`/admin/products${queryString}`)
      console.log('ProductService.getAdminProducts response:', response)

      // API client already returns the correct data structure
      return response as { data: Product[], pagination?: any }
    } catch (error) {
      console.error('ProductService.getAdminProducts error:', error)
      throw error
    }
  }

  async createProduct(data: CreateProductRequest): Promise<Product> {
    // FE: Tự động generate slug nếu chưa có
    if (!data.slug || validateSlug(data.slug)) {
      data.slug = generateSlug(data.name)
    }
    // FE: Validate lại slug
    const slugError = validateSlug(data.slug)
    if (slugError) throw new Error(slugError)
    const response = await apiClient.post('/admin/products', data)
    return response as Product
  }

  async updateProduct(id: string, data: UpdateProductRequest): Promise<Product> {
    // FE: Tự động generate slug nếu chưa có
    if (!data.slug && data.name) {
      data.slug = generateSlug(data.name as string)
    }
    // FE: Validate lại slug nếu có
    if (data.slug) {
      const slugError = validateSlug(data.slug)
      if (slugError) throw new Error(slugError)
    }
    const response = await apiClient.put(`/admin/products/${id}`, data)
    return response as Product
  }

  async deleteProduct(id: string): Promise<void> {
    try {
      const response = await apiClient.delete(`/admin/products/${id}`)
      console.log('Delete product response:', response)
    } catch (error: any) {
      console.error('Delete product API error:', error)
      // Re-throw with better error message
      const errorMessage = error?.message || error?.error || 'Failed to delete product. Please try again.'
      throw new Error(errorMessage)
    }
  }

  async uploadProductImage(productId: string, file: File): Promise<any> {
    const response = await apiClient.upload(`/admin/products/${productId}/images`, file)
    return response
  }

  async deleteProductImage(productId: string, imageId: string): Promise<void> {
    await apiClient.delete(`/admin/products/${productId}/images/${imageId}`)
  }

  // Inventory management
  async updateStock(productId: string, quantity: number): Promise<Product> {
    const response = await apiClient.patch(`/admin/products/${productId}/stock`, { quantity })
    return response as Product
  }

  async getInventory(productId: string): Promise<any> {
    const response = await apiClient.get(`/admin/products/${productId}/inventory`)
    return response
  }

  // Product analytics
  async getProductAnalytics(productId: string, period = '30d'): Promise<any> {
    const response = await apiClient.get(`/admin/products/${productId}/analytics?period=${period}`)
    return response
  }

  // Bulk operations
  async bulkUpdateProducts(productIds: string[], data: UpdateProductRequest): Promise<void> {
    await apiClient.patch('/admin/products/bulk', { product_ids: productIds, ...data })
  }

  async bulkDeleteProducts(productIds: string[]): Promise<void> {
    await apiClient.delete('/admin/products/bulk', { data: { product_ids: productIds } })
  }

  // Export/Import
  async exportProducts(filters: ProductFilters = {}): Promise<Blob> {
    const queryString = buildQueryString(filters)
    const url = `/admin/products/export${queryString ? `?${queryString}` : ''}`
    
    const response = await apiClient.getClient().get(url, {
      responseType: 'blob',
    })
    
    return response.data
  }

  async importProducts(file: File): Promise<any> {
    const response = await apiClient.upload('/admin/products/import', file)
    return response
  }

  // Product reviews
  async getProductReviews(productId: string, params: { page?: number; limit?: number } = {}): Promise<any> {
    const response = await getPaginated(`/products/${productId}/reviews`, params)
    return response
  }

  async createReview(productId: string, data: { rating: number; title?: string; comment?: string }): Promise<any> {
    const response = await apiClient.post(`/products/${productId}/reviews`, data)
    return response
  }

  // Wishlist
  async addToWishlist(productId: string): Promise<void> {
    await apiClient.post('/wishlist/items', { product_id: productId })
  }

  async removeFromWishlist(productId: string): Promise<void> {
    await apiClient.delete(`/wishlist/items/${productId}`)
  }

  // Product comparison
  async compareProducts(productIds: string[]): Promise<any> {
    const response = await apiClient.post('/products/compare', { product_ids: productIds })
    return response
  }

  // Recently viewed products
  async getRecentlyViewed(limit = 10): Promise<any> {
    const response = await apiClient.get(`/products/recently-viewed?limit=${limit}`)
    return response
  }

  async addToRecentlyViewed(productId: string): Promise<void> {
    await apiClient.post('/products/recently-viewed', { product_id: productId })
  }
}

export const productService = new ProductService()

// Utility: Generate slug from name
export function generateSlug(name: string): string {
  if (!name) return ''
  let slug = name.toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .replace(/--+/g, '-')
  if (!slug) slug = `product-${Date.now()}`
  return slug
}

export function validateSlug(slug: string): string | null {
  if (!slug) return 'Slug cannot be empty'
  if (slug.length > 255) return 'Slug too long'
  if (!/^[a-z0-9-]+$/.test(slug)) return 'Slug can only contain lowercase letters, numbers, and hyphens'
  if (/^-|-$/.test(slug)) return 'Slug cannot start or end with hyphen'
  if (/--/.test(slug)) return 'Slug cannot contain consecutive hyphens'
  return null
}
