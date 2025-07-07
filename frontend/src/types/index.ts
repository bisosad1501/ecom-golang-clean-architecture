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

// Order types
export interface Order extends BaseEntity {
  order_number: string
  user_id: string
  status: 'pending' | 'confirmed' | 'processing' | 'shipped' | 'delivered' | 'cancelled' | 'refunded'
  payment_status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled' | 'refunded'
  items: OrderItem[]
  shipping_address: Address
  billing_address: Address
  subtotal: number
  tax_amount: number
  shipping_amount: number
  discount_amount: number
  total_amount: number
}

export interface OrderItem {
  id: string
  product_id: string
  product_name: string
  product_image?: string
  quantity: number
  price: number
  total: number
}

export interface CreateOrderRequest {
  items: Array<{
    product_id: string
    quantity: number
    price: number
  }>
  shipping_address: Address
  billing_address: Address
  payment_method: string
}

// Cart types
export interface Cart {
  id: string
  user_id?: string
  items: CartItem[]
  subtotal: number
  total: number
  created_at: string
  updated_at: string
}

export interface CartItem {
  id: string
  product_id: string
  product: Product
  quantity: number
  price: number
  total: number
}

// Payment types
export interface PaymentStore {
  paymentMethods: PaymentMethod[]
  currentPayment: Payment | null
  isLoading: boolean
  error: string | null
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
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled' | 'refunded'
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



