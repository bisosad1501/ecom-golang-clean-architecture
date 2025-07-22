import type { Product } from './product'

// Wishlist item interface (matches backend WishlistItemResponse)
export interface WishlistItem {
  id: string
  product: Product
  added_at: string
}

// Wishlist API response interface (matches backend PaginatedResponse)
export interface WishlistApiResponse {
  data: WishlistItem[]
  pagination: {
    total: number
    page: number
    limit: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
    start_index: number
    end_index: number
    next_page: number | null
    prev_page: number | null
    page_sizes: number[]
    use_cursor: boolean
    cache_key: string
    cache_ttl: number
    is_cached: boolean
    query_time: number
  }
}

// Wishlist response for hooks (just the data array)
export type WishlistResponse = WishlistItem[]

// Get wishlist request interface (matches backend GetWishlistRequest)
export interface GetWishlistRequest {
  limit?: number
  offset?: number
}

// Add to wishlist request interface
export interface AddToWishlistRequest {
  product_id: string
}

// Wishlist count response
export interface WishlistCountResponse {
  data: {
    count: number
  }
}

// Wishlist filters for frontend
export interface WishlistFilters {
  search?: string
  category?: string
  price_min?: number
  price_max?: number
  sort_by?: 'added_at' | 'product_name' | 'price'
  sort_order?: 'asc' | 'desc'
}
