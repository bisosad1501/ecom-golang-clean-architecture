// ===== WISHLIST SERVICE =====

import { apiClient } from '@/lib/api'
import type {
  WishlistApiResponse,
  GetWishlistRequest,
  AddToWishlistRequest,
  WishlistCountResponse,
} from '@/types/wishlist'

export class WishlistService {
  // Get user's wishlist
  async getWishlist(params?: GetWishlistRequest): Promise<WishlistApiResponse> {
    const searchParams = new URLSearchParams()

    if (params?.limit) {
      searchParams.append('limit', params.limit.toString())
    }
    if (params?.offset) {
      searchParams.append('offset', params.offset.toString())
    }

    const url = `/wishlist${searchParams.toString() ? `?${searchParams.toString()}` : ''}`

    const response = await apiClient.get<WishlistApiResponse>(url)



    // Handle response structure - backend returns array directly

    const result = {
      data: Array.isArray(response.data) ? response.data : [],
      pagination: {
        total: Array.isArray(response.data) ? response.data.length : 0,
        page: 1,
        limit: params?.limit ?? 10,
        total_pages: 1,
        has_next: false,
        has_prev: false,
        start_index: 0,
        end_index: Array.isArray(response.data) ? response.data.length : 0,
        next_page: null,
        prev_page: null,
        page_sizes: [],
        use_cursor: false,
        cache_key: '',
        cache_ttl: 0,
        is_cached: false,
        query_time: 0
      }
    }

    return result
  }

  // Add product to wishlist
  async addToWishlist(productId: string): Promise<void> {
    const data: AddToWishlistRequest = { product_id: productId }
    await apiClient.post('/wishlist/items', data)
  }

  // Remove product from wishlist
  async removeFromWishlist(productId: string): Promise<void> {
    await apiClient.delete(`/wishlist/items/${productId}`)
  }

  // Clear entire wishlist
  async clearWishlist(): Promise<void> {
    await apiClient.delete('/wishlist/clear')
  }

  // Get wishlist count
  async getWishlistCount(): Promise<WishlistCountResponse> {
    const response = await apiClient.get<WishlistCountResponse>('/wishlist/count')
    // Handle response structure similar to wishlist data
    if (response.data && typeof response.data === 'object' && 'data' in response.data) {
      return response.data as WishlistCountResponse
    }
    // If response.data is direct count object
    return { data: response.data as any }
  }

  // Check if product is in wishlist
  async isInWishlist(productId: string): Promise<boolean> {
    try {
      const response = await apiClient.get(`/wishlist/items/${productId}/status`)
      // Backend returns {data: {in_wishlist: boolean}} with lowercase
      return response.data?.data?.in_wishlist || false
    } catch (error) {
      return false
    }
  }

  // Move wishlist item to cart (if backend supports it)
  async moveToCart(productId: string): Promise<void> {
    await apiClient.post(`/wishlist/items/${productId}/move-to-cart`)
  }
}

// Create and export singleton instance
export const wishlistService = new WishlistService()
export default wishlistService
