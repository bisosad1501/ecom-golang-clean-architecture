// ===== UNIFIED TYPES EXPORT =====
// Export all types from organized files
export * from './common'
export * from './auth'
export * from './product'

// Re-export commonly used types for convenience
export type {
  BaseEntity,
  ApiResponse,
  PaginatedResponse,
  LoadingState,
  SortOrder,
  Address,
  Image,
  Theme,
  NotificationType
} from './common'

export type {
  User,
  UserRole,
  UserProfile,
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  AuthContext
} from './auth'

export type {
  Product,
  ProductListItem,
  Category,
  Brand,
  ProductStatus,
  ProductSearchParams,
  ProductFilters,
  WishlistItem
} from './product'

// ===== ADDITIONAL MISSING TYPES =====

// Order types (matching backend response structure)
export interface Order extends BaseEntity {
  order_number: string
  user?: {
    id: string
    email: string
    first_name: string
    last_name: string
  }

  // Order Status & Management
  status: 'pending' | 'confirmed' | 'processing' | 'ready_to_ship' | 'shipped' | 'out_for_delivery' | 'delivered' | 'cancelled' | 'refunded' | 'returned' | 'exchanged'
  fulfillment_status: 'pending' | 'processing' | 'packed' | 'shipped' | 'delivered' | 'returned' | 'cancelled'
  payment_status: 'pending' | 'processing' | 'paid' | 'failed' | 'cancelled' | 'refunded'
  priority: 'low' | 'normal' | 'high' | 'urgent' | 'critical'
  source: 'web' | 'mobile' | 'admin' | 'api' | 'phone' | 'email' | 'social'
  customer_type: 'guest' | 'registered' | 'vip' | 'wholesale' | 'corporate'

  // Financial Information
  items: OrderItem[]
  shipping_address?: OrderAddress
  billing_address?: OrderAddress
  subtotal: number
  tax_amount: number
  shipping_amount: number
  discount_amount: number
  tip_amount: number
  total: number
  currency: string

  // Shipping & Delivery
  shipping_method?: string
  tracking_number?: string
  tracking_url?: string
  carrier?: string
  estimated_delivery?: string
  actual_delivery?: string
  delivery_instructions?: string
  delivery_attempts: number

  // Customer Information
  customer_notes?: string
  admin_notes?: string
  internal_notes?: string

  // Gift Options
  is_gift: boolean
  gift_message?: string
  gift_wrap: boolean

  // Business Information
  sales_channel?: string
  referral_source?: string
  coupon_codes?: string
  tags?: string

  // Fulfillment Information
  warehouse_id?: string
  packed_at?: string
  shipped_at?: string
  processed_at?: string

  // Order Capabilities
  payment?: any
  item_count: number
  can_be_cancelled: boolean
  can_be_refunded: boolean
  can_be_shipped: boolean
  can_be_delivered: boolean
  is_shipped: boolean
  is_delivered: boolean
  has_tracking: boolean

  created_at: string
  updated_at: string
}

export interface OrderItem {
  id: string
  product?: any  // Use any for now to avoid circular import
  product_name: string
  product_sku: string
  quantity: number
  price: number
  total: number
}

export interface OrderAddress {
  id?: string
  type?: string
  first_name: string
  last_name: string
  company?: string
  address1: string // Backend uses address1, not address_line_1
  address2?: string // Backend uses address2, not address_line_2
  city: string
  state: string
  zip_code: string // Backend uses zip_code, not postal_code
  country: string
  phone?: string
  is_default?: boolean
}

// Order Event types
export interface OrderEvent {
  id: string
  order_id: string
  event_type: 'created' | 'status_changed' | 'payment_received' | 'payment_failed' | 'shipped' | 'delivered' | 'cancelled' | 'refunded' | 'returned' | 'note_added' | 'tracking_updated' | 'inventory_reserved' | 'inventory_released' | 'custom'
  title: string
  description: string
  data?: string
  user_id?: string
  user?: {
    id: string
    first_name: string
    last_name: string
    email: string
  }
  is_public: boolean
  created_at: string
}

export interface CreateOrderRequest {
  items: Array<{
    product_id: string
    quantity: number
    price: number
  }>
  shipping_address: any
  billing_address: any
  payment_method: string
}

// Cart types - matching backend CartResponse exactly
export interface Cart {
  id: string
  user_id?: string
  session_id?: string  // For guest carts
  items: CartItem[]
  item_count: number   // Calculated field from backend
  subtotal: number
  tax_amount: number   // Added from backend
  shipping_amount: number // Added from backend
  total: number
  status: string       // active, abandoned, converted
  currency: string     // USD, EUR, etc.
  notes?: string       // Added from backend
  expires_at?: string  // Cart expiration time
  is_guest: boolean    // Added from backend
  created_at: string
  updated_at: string
}

export interface CartItem {
  id: string
  product: Product     // Use proper Product type with all computed fields
  quantity: number
  price: number        // Price at time of adding to cart
  subtotal: number     // Backend calculated subtotal (price * quantity)
  created_at: string
  updated_at: string
}

// Payment types
export interface PaymentStore {
  paymentMethods: PaymentMethod[]
  currentPayment: Payment | null
  isLoading: boolean
  error: string | null

  // Methods
  fetchPaymentMethods: () => Promise<void>
  savePaymentMethod: (data: any) => Promise<void>
  deletePaymentMethod: (methodId: string) => Promise<void>
  setDefaultPaymentMethod: (methodId: string) => Promise<void>
  createCheckoutSession: (data: CreateCheckoutSessionRequest) => Promise<CheckoutSession>
  processPayment: (data: ProcessPaymentRequest) => Promise<void>
  fetchPayment: (paymentId: string) => Promise<void>
  processRefund: (paymentId: string, data: RefundRequest) => Promise<void>
}

export interface PaymentMethod {
  id: string
  type: 'card' | 'paypal' | 'bank_transfer' | 'cash_on_delivery'
  name: string
  details: Record<string, any>
  is_default: boolean
}

export interface Payment {
  id: string
  order_id: string
  amount: number
  currency: string
  status: 'pending' | 'processing' | 'paid' | 'completed' | 'failed' | 'cancelled' | 'refunded'
  method: PaymentMethod
  created_at: string
}

export interface CheckoutSession {
  id: string
  url: string
  expires_at: string
}

export interface CreateCheckoutSessionRequest {
  order_id: string
  success_url: string
  cancel_url: string
}

export interface ProcessPaymentRequest {
  order_id: string
  payment_method_id: string
  amount: number
}

export interface RefundRequest {
  payment_id: string
  amount?: number
  reason?: string
}

// Cart conflict detection types
export interface ConflictingItem {
  product_id: string
  product_name: string
  user_quantity: number
  guest_quantity: number
  user_price: number
  guest_price: number
  price_difference: number
}

export interface CartConflictInfo {
  has_conflict: boolean
  user_cart_exists: boolean
  guest_cart_exists: boolean
  conflicting_items?: ConflictingItem[]
  user_cart?: Cart
  guest_cart?: Cart
  recommendations?: string[]
}

// Merge strategy types
export type MergeStrategy = 'auto' | 'replace' | 'keep_user' | 'merge'