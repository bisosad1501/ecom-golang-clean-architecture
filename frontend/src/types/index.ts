// Base Types
export interface BaseEntity {
  id: string
  created_at: string
  updated_at: string
}

// User Types
export interface User extends BaseEntity {
  email: string
  first_name: string
  last_name: string
  role: UserRole
  is_active: boolean
  email_verified: boolean
  profile?: UserProfile
}

export interface UserProfile extends BaseEntity {
  user_id: string
  phone?: string
  date_of_birth?: string
  gender?: 'male' | 'female' | 'other'
  avatar_url?: string
  bio?: string
  preferences?: UserPreferences
}

export interface UserPreferences {
  newsletter_subscribed: boolean
  marketing_emails: boolean
  order_updates: boolean
  theme: 'light' | 'dark' | 'system'
  language: string
  currency: string
}

export type UserRole = 'customer' | 'admin' | 'moderator' | 'super_admin'

// Authentication Types
export interface AuthResponse {
  user: User
  token: string
  refresh_token: string
  expires_at: string
}

export interface LoginRequest {
  email: string
  password: string
  remember_me?: boolean
}

export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  terms_accepted: boolean
}

export interface ResetPasswordRequest {
  email: string
}

export interface ChangePasswordRequest {
  current_password: string
  new_password: string
  confirm_password: string
}

// Product Types
export interface Product extends BaseEntity {
  name: string
  description: string
  short_description?: string
  sku: string
  price: number
  compare_price?: number
  cost_price?: number
  stock: number
  status: ProductStatus
  is_digital: boolean
  is_available: boolean
  has_discount: boolean
  weight?: number
  dimensions?: Dimensions
  category_id: string
  category?: Category
  images: ProductImage[]
  tags: ProductTag[]
  reviews?: Review[]
  rating?: ProductRating
  inventory?: Inventory
}

export interface Dimensions {
  length: number
  width: number
  height: number
}

export interface ProductImage extends BaseEntity {
  product_id?: string
  url: string
  alt_text?: string
  position: number
}

export interface ProductTag extends BaseEntity {
  name: string
  slug: string
  color?: string
}

export interface ProductRating {
  average: number
  count: number
  distribution: {
    1: number
    2: number
    3: number
    4: number
    5: number
  }
}

export type ProductStatus = 'draft' | 'active' | 'inactive' | 'archived'

// Category Types
export interface Category extends BaseEntity {
  name: string
  slug: string
  description?: string
  image_url?: string
  parent_id?: string
  parent?: Category
  children?: Category[]
  is_active: boolean
  sort_order: number
  level?: number
  path?: string
  product_count?: number
}

// Cart Types
export interface Cart extends BaseEntity {
  user_id: string
  items: CartItem[]
  total_amount: number
  item_count: number
}

export interface CartItem extends BaseEntity {
  cart_id: string
  product_id: string
  product: Product
  quantity: number
  unit_price: number
  total_price: number
}

// Order Types
export interface Order extends BaseEntity {
  order_number: string
  user_id: string
  user?: User
  status: OrderStatus
  payment_status: PaymentStatus
  items: OrderItem[]
  shipping_address: Address
  billing_address: Address
  subtotal: number
  tax_amount: number
  shipping_amount: number
  discount_amount: number
  total_amount: number
  notes?: string
  shipped_at?: string
  delivered_at?: string
}

export interface OrderItem extends BaseEntity {
  order_id: string
  product_id: string
  product: Product
  quantity: number
  unit_price: number
  total_price: number
}

export type OrderStatus = 
  | 'pending'
  | 'confirmed'
  | 'processing'
  | 'shipped'
  | 'delivered'
  | 'cancelled'
  | 'refunded'

export type PaymentStatus = 
  | 'pending'
  | 'processing'
  | 'completed'
  | 'failed'
  | 'cancelled'
  | 'refunded'

// Address Types
export interface Address extends BaseEntity {
  user_id: string
  type: 'shipping' | 'billing'
  first_name: string
  last_name: string
  company?: string
  address1: string
  address2?: string
  city: string
  state: string
  zip_code: string
  country: string
  phone?: string
  is_default: boolean
}

// Payment Types
export interface Payment extends BaseEntity {
  order_id: string
  method: PaymentMethod
  status: PaymentStatus
  amount: number
  currency: string
  transaction_id?: string
  gateway_response?: string
  processed_at?: string
}

export type PaymentMethod = 'credit_card' | 'debit_card' | 'paypal' | 'stripe' | 'bank_transfer'

// Review Types
export interface Review extends BaseEntity {
  product_id: string
  user_id: string
  user?: User
  rating: number
  title?: string
  comment?: string
  status: ReviewStatus
  is_verified_purchase: boolean
  helpful_count: number
  images?: ReviewImage[]
}

export interface ReviewImage extends BaseEntity {
  review_id: string
  url: string
  alt_text?: string
}

export type ReviewStatus = 'pending' | 'approved' | 'rejected'

// Wishlist Types
export interface Wishlist extends BaseEntity {
  user_id: string
  items: WishlistItem[]
}

export interface WishlistItem extends BaseEntity {
  wishlist_id: string
  product_id: string
  product: Product
}

// Coupon Types
export interface Coupon extends BaseEntity {
  code: string
  type: CouponType
  value: number
  minimum_amount?: number
  maximum_discount?: number
  usage_limit?: number
  used_count: number
  expires_at?: string
  is_active: boolean
}

export type CouponType = 'percentage' | 'fixed_amount' | 'free_shipping'

// Inventory Types
export interface Inventory extends BaseEntity {
  product_id: string
  warehouse_id?: string
  quantity: number
  reserved_quantity: number
  available_quantity: number
  reorder_level: number
  reorder_quantity: number
}

// Shipping Types
export interface ShippingMethod extends BaseEntity {
  name: string
  description?: string
  carrier: string
  base_cost: number
  cost_per_kg: number
  min_delivery_days: number
  max_delivery_days: number
  is_active: boolean
}

export interface Shipment extends BaseEntity {
  order_id: string
  tracking_number: string
  carrier: string
  status: ShipmentStatus
  shipped_at?: string
  estimated_delivery?: string
  actual_delivery?: string
}

export type ShipmentStatus = 
  | 'pending'
  | 'processing'
  | 'shipped'
  | 'in_transit'
  | 'out_for_delivery'
  | 'delivered'
  | 'failed'
  | 'returned'

// API Response Types
export interface ApiResponse<T = any> {
  data: T
  message?: string
  success: boolean
}

export interface PaginatedResponse<T = any> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean
  }
}

export interface ApiError {
  message: string
  code?: string
  details?: Record<string, any>
}

// Filter and Search Types
export interface ProductFilters {
  category_id?: string
  min_price?: number
  max_price?: number
  in_stock?: boolean
  rating?: number
  tags?: string[]
  search?: string
}

export interface SortOption {
  field: string
  direction: 'asc' | 'desc'
}

export interface PaginationParams {
  page?: number
  limit?: number
}

// Form Types
export interface ContactForm {
  name: string
  email: string
  subject: string
  message: string
}

export interface NewsletterForm {
  email: string
}

// UI State Types
export interface LoadingState {
  isLoading: boolean
  error?: string | null
}

export interface ModalState {
  isOpen: boolean
  data?: any
}

// Analytics Types
export interface AnalyticsEvent {
  event_type: string
  event_name: string
  properties?: Record<string, any>
  user_id?: string
  session_id: string
}

// Notification Types
export interface Notification extends BaseEntity {
  user_id: string
  type: NotificationType
  title: string
  message: string
  is_read: boolean
  action_url?: string
}

export type NotificationType = 
  | 'order_update'
  | 'payment_success'
  | 'payment_failed'
  | 'product_back_in_stock'
  | 'promotion'
  | 'system'

// Search Types
export interface SearchResult {
  products: Product[]
  categories: Category[]
  total_results: number
  search_time: number
  suggestions?: string[]
}

export interface SearchFilters extends ProductFilters {
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

// Theme Types
export type Theme = 'light' | 'dark' | 'system'

// Component Props Types
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive'
  size?: 'sm' | 'md' | 'lg'
  isLoading?: boolean
  disabled?: boolean
  children: React.ReactNode
  onClick?: () => void
  type?: 'button' | 'submit' | 'reset'
  className?: string
}

export interface InputProps {
  label?: string
  placeholder?: string
  error?: string
  required?: boolean
  disabled?: boolean
  type?: string
  value?: string
  onChange?: (value: string) => void
  className?: string
}

// Store Types (for Zustand)
export interface AuthStore {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => void
  updateProfile: (data: Partial<UserProfile>) => Promise<void>
}

export interface CartStore {
  cart: Cart | null
  isLoading: boolean
  addItem: (productId: string, quantity: number) => Promise<void>
  updateItem: (itemId: string, quantity: number) => Promise<void>
  removeItem: (itemId: string) => Promise<void>
  clearCart: () => Promise<void>
  fetchCart: () => Promise<void>
}

export interface WishlistStore {
  wishlist: Wishlist | null
  isLoading: boolean
  addItem: (productId: string) => Promise<void>
  removeItem: (productId: string) => Promise<void>
  fetchWishlist: () => Promise<void>
}
