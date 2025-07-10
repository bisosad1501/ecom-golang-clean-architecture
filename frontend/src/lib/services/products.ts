import { apiClient, buildQueryString, getPaginated } from '@/lib/api'
import { Product, PaginatedResponse, ProductFilters, ApiResponse } from '@/types'

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

export interface CreateProductRequest {
  name: string
  description: string
  short_description: string
  sku: string

  // SEO and Metadata
  slug?: string
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

  // Status and Type
  status?: string
  product_type?: 'simple' | 'variable' | 'grouped' | 'external' | 'digital'
  is_digital?: boolean
}

export interface UpdateProductRequest extends Partial<CreateProductRequest> {}

class ProductService {
  // Get all products with filters and pagination
  async getProducts(params: ProductsParams = {}): Promise<PaginatedResponse<Product>> {
    console.log('ProductService.getProducts called with params:', params)
    try {
      const response = await apiClient.get<Product[]>('/products')
      console.log('ProductService.getProducts response:', response)
      console.log('ProductService.getProducts response.data:', response.data)
      
      // Backend response structure: { data: Product[] }
      const products = response.data || []
      console.log('ProductService.getProducts products array:', products)
      console.log('ProductService.getProducts products length:', products.length)
      
      // Transform backend response to match frontend expectations
      const paginatedResponse: PaginatedResponse<Product> = {
        data: products,
        pagination: {
          page: params.page || 1,
          limit: params.limit || 10,
          total: products.length,
          total_pages: 1,
          has_next: false,
          has_prev: false,
        }
      }
      
      console.log('ProductService.getProducts transformed response:', paginatedResponse)
      return paginatedResponse
    } catch (error) {
      console.error('ProductService.getProducts error:', error)
      throw error
    }
  }

  // Get single product by ID
  async getProduct(id: string): Promise<Product> {
    const response = await apiClient.get<Product>(`/products/${id}`)
    return response.data
  }

  // Get product by slug
  async getProductBySlug(slug: string): Promise<Product> {
    const response = await apiClient.get<Product>(`/products/slug/${slug}`)
    return response.data
  }

  // Get featured products
  async getFeaturedProducts(limit = 8): Promise<Product[]> {
    const response = await apiClient.get<Product[]>(`/products/featured?limit=${limit}`)
    return response.data || []
  }

  // Get related products
  async getRelatedProducts(productId: string, limit = 4): Promise<Product[]> {
    const response = await apiClient.get<Product[]>(`/products/${productId}/related?limit=${limit}`)
    return response.data || []
  }

  // Get products by category
  async getProductsByCategory(categoryId: string, params: ProductsParams = {}): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<Product[]>(`/products/category/${categoryId}`)
    return {
      data: response.data || [],
      pagination: {
        page: params.page || 1,
        limit: params.limit || 10,
        total: response.data?.length || 0,
        total_pages: 1,
        has_next: false,
        has_prev: false,
      }
    }
  }

  // Search products
  async searchProducts(query: string, params: ProductsParams = {}): Promise<PaginatedResponse<Product>> {
    const response = await apiClient.get<Product[]>('/products/search', {
      params: { ...params, search: query }
    })
    return {
      data: response.data || [],
      pagination: {
        page: params.page || 1,
        limit: params.limit || 10,
        total: response.data?.length || 0,
        total_pages: 1,
        has_next: false,
        has_prev: false,
      }
    }
  }

  // Get product suggestions for autocomplete
  async getProductSuggestions(query: string, limit = 5): Promise<Product[]> {
    const response = await apiClient.get<Product[]>(`/products/suggestions?q=${query}&limit=${limit}`)
    return response.data || []
  }

  // Admin methods
  async getAdminProducts(params: ProductsParams = {}): Promise<PaginatedResponse<Product>> {
    console.log('ProductService.getAdminProducts called with params:', params)
    try {
      const queryString = buildQueryString(params)
      const endpoint = `/admin/products${queryString ? `?${queryString}` : ''}`
      console.log('ProductService.getAdminProducts endpoint:', endpoint)
      
      const response = await apiClient.get<PaginatedResponse<Product>>(endpoint)
      console.log('ProductService.getAdminProducts response:', response)
      console.log('ProductService.getAdminProducts response.data:', response.data)
      
      return response.data
    } catch (error) {
      console.error('ProductService.getAdminProducts error:', error)
      throw error
    }
  }

  async createProduct(data: CreateProductRequest): Promise<Product> {
    const response = await apiClient.post<Product>('/admin/products', data)
    return response.data
  }

  async updateProduct(id: string, data: UpdateProductRequest): Promise<Product> {
    console.log('=== ProductService.updateProduct START ===')
    console.log('Product ID:', id)
    console.log('Update data:', data)
    console.log('Update data category_id:', data.category_id)
    console.log('Update data category_id type:', typeof data.category_id)
    
    const response = await apiClient.put<Product>(`/admin/products/${id}`, data)
    
    console.log('=== ProductService.updateProduct RESPONSE ===')
    console.log('Response:', response)
    console.log('Response data:', response.data)
    console.log('Response data category:', response.data?.category)
    console.log('Response data category_id:', response.data?.category_id)
    
    return response.data
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

  async uploadProductImage(productId: string, file: File): Promise<{ url: string }> {
    const response = await apiClient.upload<{ url: string }>(`/admin/products/${productId}/images`, file)
    return response.data
  }

  async deleteProductImage(productId: string, imageId: string): Promise<void> {
    await apiClient.delete(`/admin/products/${productId}/images/${imageId}`)
  }

  // Inventory management
  async updateStock(productId: string, quantity: number): Promise<Product> {
    const response = await apiClient.patch<Product>(`/admin/products/${productId}/stock`, { quantity })
    return response.data
  }

  async getInventory(productId: string): Promise<any> {
    const response = await apiClient.get(`/admin/products/${productId}/inventory`)
    return response.data
  }

  // Product analytics
  async getProductAnalytics(productId: string, period = '30d'): Promise<any> {
    const response = await apiClient.get(`/admin/products/${productId}/analytics?period=${period}`)
    return response.data
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

  async importProducts(file: File): Promise<{ success: number; errors: any[] }> {
    const response = await apiClient.upload<{ success: number; errors: any[] }>('/admin/products/import', file)
    return response.data
  }

  // Product reviews
  async getProductReviews(productId: string, params: { page?: number; limit?: number } = {}): Promise<any> {
    const response = await getPaginated(`/products/${productId}/reviews`, params)
    return response.data
  }

  async createReview(productId: string, data: { rating: number; title?: string; comment?: string }): Promise<any> {
    const response = await apiClient.post(`/products/${productId}/reviews`, data)
    return response.data
  }

  // Wishlist
  async addToWishlist(productId: string): Promise<void> {
    await apiClient.post('/wishlist/items', { product_id: productId })
  }

  async removeFromWishlist(productId: string): Promise<void> {
    await apiClient.delete(`/wishlist/items/${productId}`)
  }

  // Product comparison
  async compareProducts(productIds: string[]): Promise<Product[]> {
    const response = await apiClient.post<Product[]>('/products/compare', { product_ids: productIds })
    return response.data
  }

  // Recently viewed products
  async getRecentlyViewed(limit = 10): Promise<Product[]> {
    const response = await apiClient.get<Product[]>(`/products/recently-viewed?limit=${limit}`)
    return response.data
  }

  async addToRecentlyViewed(productId: string): Promise<void> {
    await apiClient.post('/products/recently-viewed', { product_id: productId })
  }
}

export const productService = new ProductService()
