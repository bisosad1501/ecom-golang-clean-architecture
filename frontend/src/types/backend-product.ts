// Backend Product Response Types
// These types match the actual response structure from the Go backend

export interface BackendProductResponse {
  id: string
  name: string
  description: string
  short_description?: string
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
  category?: {
    id: string
    name: string
    description: string
    slug: string
    image: string
  }
  category_id?: string
  images?: Array<{
    id: string
    product_id: string
    url: string
    alt_text: string
    position: number
    created_at: string
  }>
  tags?: Array<{
    id: string
    name: string
    slug: string
    color?: string
  }>
  status: 'active' | 'draft' | 'archived' | 'inactive'
  is_digital: boolean
  is_available: boolean
  has_discount: boolean
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
      currency: 'USD',
      tax_inclusive: false,
    },
    inventory: {
      sku: backendProduct.sku,
      stock_quantity: backendProduct.stock,
      reserved_quantity: 0,
      available_quantity: backendProduct.stock,
      stock_status: backendProduct.stock > 0 ? 'in_stock' : 'out_of_stock',
      track_inventory: true,
      allow_backorders: false,
    },
    shipping: {
      weight: backendProduct.weight ? { value: backendProduct.weight, unit: 'kg' } : undefined,
      dimensions: backendProduct.dimensions,
      free_shipping: false,
    },
    // Keep original fields for backward compatibility
    compare_price: backendProduct.compare_price,
    cost_price: backendProduct.cost_price,
    stock: backendProduct.stock,
    weight: backendProduct.weight,
    is_digital: backendProduct.is_digital,
    short_description: backendProduct.short_description,
  }
}
