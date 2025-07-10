// Backend Product Response Types
// These types match the actual response structure from the Go backend

export interface BackendProductResponse {
  id: string
  name: string
  description: string
  short_description: string
  sku: string

  // SEO and Metadata
  slug: string
  meta_title: string
  meta_description: string
  keywords: string
  featured: boolean
  visibility: 'visible' | 'hidden' | 'private'

  // Pricing
  price: number
  compare_price?: number
  cost_price?: number

  // Sale Pricing
  sale_price?: number
  sale_start_date?: string
  sale_end_date?: string
  current_price: number
  is_on_sale: boolean
  sale_discount_percentage: number

  // Inventory
  stock: number
  low_stock_threshold: number
  track_quantity: boolean
  allow_backorder: boolean
  stock_status: 'in_stock' | 'out_of_stock' | 'on_backorder' | 'discontinued'
  is_low_stock: boolean

  // Physical Properties
  weight?: number
  dimensions?: {
    length: number
    width: number
    height: number
  }

  // Shipping and Tax
  requires_shipping: boolean
  shipping_class: string
  tax_class: string
  country_of_origin: string

  // Categorization
  category?: {
    id: string
    name: string
    description: string
    slug: string
    image: string
  }
  brand?: {
    id: string
    name: string
    slug: string
    description: string
    logo: string
    website: string
    is_active: boolean
  }

  // Content
  images: Array<{
    id: string
    product_id: string
    url: string
    alt_text: string
    position: number
    created_at: string
  }>
  tags: Array<{
    id: string
    name: string
    slug: string
  }>
  attributes: Array<{
    id: string
    attribute_id: string
    term_id?: string
    name: string
    value: string
    position: number
  }>
  variants: Array<{
    id: string
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
    image: string
    position: number
    is_active: boolean
    attributes: Array<{
      id: string
      attribute_id: string
      attribute_name: string
      term_id: string
      term_name: string
      term_value: string
    }>
  }>

  // Status and Type
  status: 'active' | 'draft' | 'archived' | 'inactive'
  product_type: 'simple' | 'variable' | 'grouped' | 'external' | 'digital'
  is_digital: boolean
  is_available: boolean
  has_discount: boolean
  has_variants: boolean
  main_image: string

  created_at: string
  updated_at: string
}

// Type guard to check if a product is a backend response
export function isBackendProduct(product: any): product is BackendProductResponse {
  return product && typeof product.id === 'string' && typeof product.price === 'number'
}

// Helper function to convert backend product to frontend format
export function convertBackendProductToFrontend(backendProduct: BackendProductResponse): any {
  return {
    ...backendProduct,
    // Map backend fields to frontend expected structure
    pricing: {
      price: backendProduct.price,
      compare_price: backendProduct.compare_price,
      cost_price: backendProduct.cost_price,
      sale_price: backendProduct.sale_price,
      current_price: backendProduct.current_price,
      currency: 'USD',
      tax_inclusive: false,
      discount_percentage: backendProduct.sale_discount_percentage,
    },
    inventory: {
      sku: backendProduct.sku,
      stock_quantity: backendProduct.stock,
      reserved_quantity: 0,
      available_quantity: backendProduct.stock,
      stock_status: backendProduct.stock_status,
      low_stock_threshold: backendProduct.low_stock_threshold,
      track_inventory: backendProduct.track_quantity,
      allow_backorders: backendProduct.allow_backorder,
    },
    shipping: {
      weight: backendProduct.weight ? { value: backendProduct.weight, unit: 'kg' } : undefined,
      dimensions: backendProduct.dimensions,
      shipping_class: backendProduct.shipping_class,
      free_shipping: !backendProduct.requires_shipping,
      origin_country: backendProduct.country_of_origin,
    },
    seo: {
      meta_title: backendProduct.meta_title,
      meta_description: backendProduct.meta_description,
      keywords: backendProduct.keywords,
      slug: backendProduct.slug,
    },
    // Keep original fields for backward compatibility and direct access
    slug: backendProduct.slug,
    meta_title: backendProduct.meta_title,
    meta_description: backendProduct.meta_description,
    keywords: backendProduct.keywords,
    featured: backendProduct.featured,
    visibility: backendProduct.visibility,
    sale_price: backendProduct.sale_price,
    sale_start_date: backendProduct.sale_start_date,
    sale_end_date: backendProduct.sale_end_date,
    current_price: backendProduct.current_price,
    is_on_sale: backendProduct.is_on_sale,
    sale_discount_percentage: backendProduct.sale_discount_percentage,
    low_stock_threshold: backendProduct.low_stock_threshold,
    track_quantity: backendProduct.track_quantity,
    allow_backorder: backendProduct.allow_backorder,
    stock_status: backendProduct.stock_status,
    is_low_stock: backendProduct.is_low_stock,
    requires_shipping: backendProduct.requires_shipping,
    shipping_class: backendProduct.shipping_class,
    tax_class: backendProduct.tax_class,
    country_of_origin: backendProduct.country_of_origin,
    brand: backendProduct.brand,
    attributes: backendProduct.attributes,
    variants: backendProduct.variants,
    product_type: backendProduct.product_type,
    has_variants: backendProduct.has_variants,
    main_image: backendProduct.main_image,
    // Legacy fields
    compare_price: backendProduct.compare_price,
    cost_price: backendProduct.cost_price,
    stock: backendProduct.stock,
    weight: backendProduct.weight,
    is_digital: backendProduct.is_digital,
    short_description: backendProduct.short_description,
  }
}
