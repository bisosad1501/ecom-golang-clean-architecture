// ===== PRODUCTS API SERVICE =====

import { apiClient } from './client'
import { API_CONFIG, PAGINATION } from '@/constants'
import type {
  Product,
  ProductListItem,
  ProductSearchParams,
  ProductFilters,
  PaginatedResponse,
  Category,
  Brand,
  ProductReview,
  ProductQuestion,
  WishlistItem,
} from '@/types'

export class ProductsService {
  // Get paginated list of products
  async getProducts(params: ProductSearchParams = {}): Promise<PaginatedResponse<ProductListItem>> {
    const searchParams = new URLSearchParams()
    
    // Add pagination
    searchParams.set('page', String(params.page || PAGINATION.DEFAULT_PAGE))
    searchParams.set('limit', String(params.limit || PAGINATION.PRODUCTS_PER_PAGE))
    
    // Add search query
    if (params.query) {
      searchParams.set('q', params.query)
    }
    
    // Add sorting
    if (params.sort_by) {
      searchParams.set('sort_by', params.sort_by)
      searchParams.set('sort_order', params.sort_order || 'asc')
    }
    
    // Add filters
    if (params.filters) {
      Object.entries(params.filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          if (Array.isArray(value)) {
            value.forEach(v => searchParams.append(key, String(v)))
          } else {
            searchParams.set(key, String(value))
          }
        }
      })
    }
    
    return apiClient.get<PaginatedResponse<ProductListItem>>(
      `${API_CONFIG.ENDPOINTS.PRODUCTS.LIST}?${searchParams.toString()}`,
      { skipAuth: true }
    )
  }

  // Get single product by ID
  async getProduct(id: string): Promise<Product> {
    return apiClient.get<Product>(
      API_CONFIG.ENDPOINTS.PRODUCTS.DETAIL.replace(':id', id),
      { skipAuth: true }
    )
  }

  // Search products
  async searchProducts(query: string, params: Omit<ProductSearchParams, 'query'> = {}): Promise<PaginatedResponse<ProductListItem>> {
    return this.getProducts({ ...params, query })
  }

  // Get featured products
  async getFeaturedProducts(limit: number = 12): Promise<ProductListItem[]> {
    return apiClient.get<ProductListItem[]>(
      `${API_CONFIG.ENDPOINTS.PRODUCTS.FEATURED}?limit=${limit}`,
      { skipAuth: true }
    )
  }

  // Get related products
  async getRelatedProducts(productId: string, limit: number = 8): Promise<ProductListItem[]> {
    return apiClient.get<ProductListItem[]>(
      API_CONFIG.ENDPOINTS.PRODUCTS.RELATED.replace(':id', productId) + `?limit=${limit}`,
      { skipAuth: true }
    )
  }

  // Get products by category
  async getProductsByCategory(categoryId: string, params: ProductSearchParams = {}): Promise<PaginatedResponse<ProductListItem>> {
    return this.getProducts({
      ...params,
      filters: {
        ...params.filters,
        category_ids: [categoryId]
      }
    })
  }

  // Get products by brand
  async getProductsByBrand(brandId: string, params: ProductSearchParams = {}): Promise<PaginatedResponse<ProductListItem>> {
    return this.getProducts({
      ...params,
      filters: {
        ...params.filters,
        brand_ids: [brandId]
      }
    })
  }

  // Get product categories
  async getCategories(): Promise<Category[]> {
    return apiClient.get<Category[]>(
      API_CONFIG.ENDPOINTS.CATEGORIES.LIST,
      { skipAuth: true }
    )
  }

  // Get single category
  async getCategory(id: string): Promise<Category> {
    return apiClient.get<Category>(
      API_CONFIG.ENDPOINTS.CATEGORIES.DETAIL.replace(':id', id),
      { skipAuth: true }
    )
  }

  // Get product brands
  async getBrands(): Promise<Brand[]> {
    return apiClient.get<Brand[]>('/brands', { skipAuth: true })
  }

  // Get product reviews
  async getProductReviews(
    productId: string,
    params: { page?: number; limit?: number } = {}
  ): Promise<PaginatedResponse<ProductReview>> {
    const searchParams = new URLSearchParams()
    searchParams.set('page', String(params.page || 1))
    searchParams.set('limit', String(params.limit || 10))
    
    return apiClient.get<PaginatedResponse<ProductReview>>(
      `/products/${productId}/reviews?${searchParams.toString()}`,
      { skipAuth: true }
    )
  }

  // Add product review
  async addProductReview(
    productId: string,
    review: {
      rating: number
      title?: string
      comment?: string
      images?: File[]
    }
  ): Promise<ProductReview> {
    if (review.images && review.images.length > 0) {
      // Handle image upload
      const formData = new FormData()
      formData.append('rating', String(review.rating))
      if (review.title) formData.append('title', review.title)
      if (review.comment) formData.append('comment', review.comment)
      
      review.images.forEach((image, index) => {
        formData.append(`images[${index}]`, image)
      })
      
      return apiClient.upload<ProductReview>(`/products/${productId}/reviews`, formData)
    } else {
      return apiClient.post<ProductReview>(`/products/${productId}/reviews`, {
        rating: review.rating,
        title: review.title,
        comment: review.comment,
      })
    }
  }

  // Update product review
  async updateProductReview(
    productId: string,
    reviewId: string,
    review: {
      rating?: number
      title?: string
      comment?: string
    }
  ): Promise<ProductReview> {
    return apiClient.put<ProductReview>(`/products/${productId}/reviews/${reviewId}`, review)
  }

  // Delete product review
  async deleteProductReview(productId: string, reviewId: string): Promise<void> {
    return apiClient.delete(`/products/${productId}/reviews/${reviewId}`)
  }

  // Get product questions
  async getProductQuestions(
    productId: string,
    params: { page?: number; limit?: number } = {}
  ): Promise<PaginatedResponse<ProductQuestion>> {
    const searchParams = new URLSearchParams()
    searchParams.set('page', String(params.page || 1))
    searchParams.set('limit', String(params.limit || 10))
    
    return apiClient.get<PaginatedResponse<ProductQuestion>>(
      `/products/${productId}/questions?${searchParams.toString()}`,
      { skipAuth: true }
    )
  }

  // Ask product question
  async askProductQuestion(
    productId: string,
    question: string
  ): Promise<ProductQuestion> {
    return apiClient.post<ProductQuestion>(`/products/${productId}/questions`, {
      question
    })
  }

  // Get wishlist
  async getWishlist(): Promise<WishlistItem[]> {
    return apiClient.get<WishlistItem[]>(API_CONFIG.ENDPOINTS.USERS.WISHLIST)
  }

  // Add to wishlist
  async addToWishlist(productId: string): Promise<WishlistItem> {
    return apiClient.post<WishlistItem>(API_CONFIG.ENDPOINTS.USERS.WISHLIST, {
      product_id: productId
    })
  }

  // Remove from wishlist
  async removeFromWishlist(productId: string): Promise<void> {
    return apiClient.delete(`${API_CONFIG.ENDPOINTS.USERS.WISHLIST}/${productId}`)
  }

  // Check if product is in wishlist
  async isInWishlist(productId: string): Promise<boolean> {
    try {
      const wishlist = await this.getWishlist()
      return wishlist.some(item => item.product_id === productId)
    } catch {
      return false
    }
  }

  // Get product recommendations
  async getRecommendations(
    type: 'related' | 'similar' | 'trending' | 'personalized' = 'personalized',
    productId?: string,
    limit: number = 8
  ): Promise<ProductListItem[]> {
    const params = new URLSearchParams()
    params.set('type', type)
    params.set('limit', String(limit))
    if (productId) params.set('product_id', productId)
    
    return apiClient.get<ProductListItem[]>(`/products/recommendations?${params.toString()}`)
  }

  // Get recently viewed products
  async getRecentlyViewed(limit: number = 10): Promise<ProductListItem[]> {
    return apiClient.get<ProductListItem[]>(`/products/recently-viewed?limit=${limit}`)
  }

  // Track product view
  async trackProductView(productId: string): Promise<void> {
    return apiClient.post('/products/track-view', { product_id: productId }, { skipAuth: true })
  }

  // Get product filters/facets
  async getProductFilters(categoryId?: string): Promise<{
    price_range: { min: number; max: number }
    brands: Array<{ id: string; name: string; count: number }>
    attributes: Array<{
      id: string
      name: string
      type: string
      options: Array<{ value: string; label: string; count: number }>
    }>
  }> {
    const params = categoryId ? `?category_id=${categoryId}` : ''
    return apiClient.get(`/products/filters${params}`, { skipAuth: true })
  }

  // Get search suggestions
  async getSearchSuggestions(query: string, limit: number = 5): Promise<string[]> {
    return apiClient.get<string[]>(
      `/products/search-suggestions?q=${encodeURIComponent(query)}&limit=${limit}`,
      { skipAuth: true }
    )
  }
}

// Create and export service instance
export const productsService = new ProductsService()

// Export for testing and custom instances
export default ProductsService
