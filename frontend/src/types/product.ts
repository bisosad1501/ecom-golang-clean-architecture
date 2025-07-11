// ===== PRODUCT & CATALOG TYPES =====

import { BaseEntity, Image, SeoMetadata, PriceRange, Dimensions, Weight } from './common'

// Product status types
export type ProductStatus = 'draft' | 'active' | 'inactive' | 'out_of_stock' | 'discontinued'

// Product visibility
export type ProductVisibility = 'public' | 'private' | 'hidden'

// Product type
export type ProductType = 'simple' | 'variable' | 'grouped' | 'external' | 'digital'

// Stock status
export type StockStatus = 'in_stock' | 'out_of_stock' | 'on_backorder' | 'discontinued'

// Category interface
export interface Category extends BaseEntity {
  name: string
  slug: string
  description?: string
  image?: Image
  parent_id?: string
  parent?: Category
  children?: Category[]
  level: number
  sort_order: number
  is_active: boolean
  product_count: number
  seo: SeoMetadata
}

// Brand interface
export interface Brand extends BaseEntity {
  name: string
  slug: string
  description?: string
  logo?: Image
  website?: string
  is_active: boolean
  product_count: number
  seo: SeoMetadata
}

// Product attribute
export interface ProductAttribute {
  id: string
  name: string
  slug: string
  type: 'text' | 'number' | 'boolean' | 'select' | 'multiselect' | 'color' | 'image'
  is_required: boolean
  is_filterable: boolean
  is_searchable: boolean
  sort_order: number
  options?: ProductAttributeOption[]
}

// Product attribute option
export interface ProductAttributeOption {
  id: string
  attribute_id: string
  value: string
  label: string
  color?: string
  image?: Image
  sort_order: number
}

// Product attribute value
export interface ProductAttributeValue {
  attribute_id: string
  attribute_name: string
  value: any
  display_value: string
}

// Product variant
export interface ProductVariant extends BaseEntity {
  product_id: string
  sku: string
  name?: string
  price: number
  compare_price?: number
  cost_price?: number
  stock_quantity: number
  stock_status: StockStatus
  weight?: Weight
  dimensions?: Dimensions
  images: Image[]
  attributes: ProductAttributeValue[]
  is_active: boolean
  sort_order: number
}

// Product pricing
export interface ProductPricing {
  price: number
  compare_price?: number
  cost_price?: number
  currency: string
  tax_inclusive: boolean
  tax_rate?: number
  discount_percentage?: number
  bulk_pricing?: Array<{
    min_quantity: number
    max_quantity?: number
    price: number
    discount_percentage?: number
  }>
}

// Product inventory
export interface ProductInventory {
  sku: string
  stock_quantity: number
  reserved_quantity: number
  available_quantity: number
  stock_status: StockStatus
  low_stock_threshold?: number
  track_inventory: boolean
  allow_backorders: boolean
  backorder_limit?: number
}

// Product shipping
export interface ProductShipping {
  weight?: Weight
  dimensions?: Dimensions
  shipping_class?: string
  free_shipping: boolean
  shipping_cost?: number
  handling_time?: number // days
  origin_country?: string
}

// Product SEO
export interface ProductSeo extends SeoMetadata {
  meta_title?: string
  meta_description?: string
  focus_keyword?: string
  schema_markup?: Record<string, any>
}

// Main product interface - Updated to match backend ProductResponse structure
export interface Product extends BaseEntity {
  name: string
  slug: string
  description?: string
  short_description?: string
  sku: string

  // SEO and Metadata
  meta_title?: string
  meta_description?: string
  keywords?: string
  featured: boolean
  visibility: ProductVisibility

  // Pricing - Direct fields matching backend
  price: number
  compare_price?: number
  cost_price?: number

  // Sale Pricing - Backend computed fields
  sale_price?: number
  sale_start_date?: string
  sale_end_date?: string
  current_price: number
  is_on_sale: boolean
  sale_discount_percentage: number

  // Inventory - Direct fields matching backend
  stock: number
  low_stock_threshold: number
  track_quantity: boolean
  allow_backorder: boolean
  stock_status: StockStatus
  is_low_stock: boolean

  // Physical Properties
  weight?: number
  dimensions?: Dimensions

  // Shipping and Tax
  requires_shipping: boolean
  shipping_class?: string
  tax_class?: string
  country_of_origin?: string

  // Categorization
  category?: Category
  brand?: Brand

  // Content
  images: Image[]
  tags?: Array<{ id: string; name: string; slug: string }>
  attributes: ProductAttributeValue[]
  variants?: ProductVariant[]

  // Status and Type
  status: ProductStatus
  product_type: ProductType
  is_digital: boolean
  is_available: boolean
  has_discount: boolean
  has_variants: boolean
  main_image?: string

  // Legacy nested structure for backward compatibility (optional)
  pricing?: ProductPricing
  inventory?: ProductInventory
  shipping?: ProductShipping
  seo?: ProductSeo

  // Stats (optional - may not be in all responses)
  view_count?: number
  sales_count?: number
  rating_average?: number
  rating_count?: number
  review_count?: number

  // Timestamps
  published_at?: string
  featured_until?: string

  // Relations
  related_products?: Product[]
  cross_sell_products?: Product[]
  up_sell_products?: Product[]
}

// Product list item (simplified for lists) - Updated to match backend
export interface ProductListItem {
  id: string
  name: string
  slug: string
  sku: string

  // Pricing - Direct backend fields
  price: number
  compare_price?: number
  current_price: number
  is_on_sale: boolean
  sale_discount_percentage: number
  has_discount: boolean

  // Inventory - Direct backend fields
  stock: number
  stock_status: StockStatus
  is_low_stock: boolean

  // Media
  images?: Image[]
  main_image?: string

  // Categorization
  category?: Pick<Category, 'id' | 'name' | 'slug'>
  brand?: Pick<Brand, 'id' | 'name' | 'slug'>

  // Status
  status: ProductStatus
  featured: boolean
  is_available: boolean

  // Stats (optional)
  rating_average?: number
  rating_count?: number

  // Timestamps
  created_at: string
  updated_at: string
}

// Product search filters
export interface ProductFilters {
  category_ids?: string[]
  brand_ids?: string[]
  price_range?: PriceRange
  attributes?: Record<string, any>
  tags?: string[]
  rating_min?: number
  in_stock?: boolean
  on_sale?: boolean
  featured?: boolean
  status?: ProductStatus[]
}

// Product search parameters
export interface ProductSearchParams {
  query?: string
  category?: string
  brand?: string
  min_price?: number
  max_price?: number
  sort_by?: 'name' | 'price' | 'created_at' | 'rating' | 'popularity' | 'sales'
  sort_order?: 'asc' | 'desc'
  page?: number
  limit?: number
  filters?: ProductFilters
}

// Product review
export interface ProductReview extends BaseEntity {
  product_id: string
  user_id: string
  user_name: string
  user_avatar?: string
  rating: number
  title?: string
  comment?: string
  images?: Image[]
  verified_purchase: boolean
  helpful_count: number
  status: 'pending' | 'approved' | 'rejected'
  admin_reply?: {
    comment: string
    created_at: string
    admin_name: string
  }
}

// Product question
export interface ProductQuestion extends BaseEntity {
  product_id: string
  user_id: string
  user_name: string
  question: string
  answer?: string
  answered_by?: string
  answered_at?: string
  status: 'pending' | 'answered' | 'closed'
  helpful_count: number
}

// Product comparison
export interface ProductComparison {
  products: Product[]
  attributes: ProductAttribute[]
  comparison_matrix: Record<string, Record<string, any>>
}

// Wishlist item
export interface WishlistItem extends BaseEntity {
  user_id: string
  product_id: string
  product: ProductListItem
  variant_id?: string
  notes?: string
  priority: 'low' | 'medium' | 'high'
  added_at: string
}

// Recently viewed product
export interface RecentlyViewedProduct {
  product_id: string
  product: ProductListItem
  viewed_at: string
}

// Product recommendation
export interface ProductRecommendation {
  type: 'related' | 'similar' | 'frequently_bought' | 'trending' | 'personalized'
  products: ProductListItem[]
  reason?: string
  confidence_score?: number
}

// Product analytics
export interface ProductAnalytics {
  product_id: string
  views: number
  unique_views: number
  add_to_cart: number
  purchases: number
  conversion_rate: number
  bounce_rate: number
  average_time_on_page: number
  revenue: number
  profit: number
  return_rate: number
  review_score: number
  search_ranking: number
}

// Bulk product operations
export interface BulkProductOperation {
  product_ids: string[]
  action: 'activate' | 'deactivate' | 'delete' | 'update_category' | 'update_price' | 'update_inventory'
  parameters?: Record<string, any>
}

// Product import/export
export interface ProductImportData {
  name: string
  sku: string
  description?: string
  price: number
  category: string
  brand?: string
  stock_quantity: number
  images?: string[]
  attributes?: Record<string, any>
}

// Product export options
export interface ProductExportOptions {
  format: 'csv' | 'xlsx' | 'json'
  include_variants: boolean
  include_images: boolean
  include_reviews: boolean
  include_analytics: boolean
  filters?: ProductFilters
}

// Product statistics
export interface ProductStats {
  total_products: number
  active_products: number
  out_of_stock_products: number
  low_stock_products: number
  featured_products: number
  products_by_category: Array<{ category: string; count: number }>
  products_by_brand: Array<{ brand: string; count: number }>
  top_selling_products: ProductListItem[]
  most_viewed_products: ProductListItem[]
  highest_rated_products: ProductListItem[]
  revenue_by_product: Array<{ product_id: string; product_name: string; revenue: number }>
}
