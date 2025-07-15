// ===== COMMON BASE TYPES =====

// Base entity interface for all database entities
export interface BaseEntity {
  id: string
  created_at: string
  updated_at: string
}

// API Response wrapper
export interface ApiResponse<T = any> {
  data: T
  message?: string
  success: boolean
  errors?: string[]
}

// Enhanced paginated response (matching backend)
export interface PaginatedResponse<T = any> {
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    total_pages: number
    has_next: boolean
    has_prev: boolean

    // Enhanced ecommerce fields
    start_index: number
    end_index: number
    next_page?: number
    prev_page?: number

    // SEO and UX fields
    canonical_url?: string
    page_sizes?: number[]

    // Performance and caching fields
    use_cursor?: boolean
    cache_key?: string
  }
}

// Ecommerce pagination context
export interface EcommercePaginationContext {
  entity_type: string
  user_id?: string
  category_id?: string
  search_query?: string
  sort_by?: string
  filter_applied?: boolean
}

// Generic list response
export interface ListResponse<T = any> {
  items: T[]
  total: number
  page: number
  limit: number
}

// Error response
export interface ErrorResponse {
  error: string
  message: string
  status_code: number
  details?: Record<string, any>
}

// Loading states
export type LoadingState = 'idle' | 'loading' | 'success' | 'error'

// Sort order
export type SortOrder = 'asc' | 'desc'

// Generic sort option
export interface SortOption {
  field: string
  order: SortOrder
  label: string
}

// Generic filter option
export interface FilterOption {
  key: string
  value: any
  label: string
  type: 'text' | 'number' | 'boolean' | 'date' | 'select' | 'range'
}

// Search parameters
export interface SearchParams {
  query?: string
  page?: number
  limit?: number
  sort_by?: string
  sort_order?: SortOrder
  filters?: Record<string, any>
}

// File upload types
export interface FileUpload {
  file: File
  preview?: string
  progress?: number
  status: 'pending' | 'uploading' | 'success' | 'error'
  error?: string
}

// Image types
export interface Image {
  id: string
  url: string
  alt?: string
  width?: number
  height?: number
  size?: number
  format?: string
}

// Address interface
export interface Address {
  id?: string
  type: 'billing' | 'shipping' | 'both'
  first_name: string
  last_name: string
  company?: string
  address_line_1: string
  address_line_2?: string
  city: string
  state: string
  postal_code: string
  country: string
  phone?: string
  is_default: boolean
}

// Contact information
export interface ContactInfo {
  email?: string
  phone?: string
  website?: string
  social_links?: Record<string, string>
}

// SEO metadata
export interface SeoMetadata {
  title?: string
  description?: string
  keywords?: string[]
  og_image?: string
  canonical_url?: string
}

// Audit fields
export interface AuditFields {
  created_by?: string
  updated_by?: string
  deleted_at?: string
  deleted_by?: string
}

// Status types
export type Status = 'active' | 'inactive' | 'pending' | 'archived' | 'deleted'

// Visibility types
export type Visibility = 'public' | 'private' | 'draft' | 'published'

// Priority levels
export type Priority = 'low' | 'medium' | 'high' | 'urgent'

// Generic key-value pair
export interface KeyValuePair {
  key: string
  value: any
  label?: string
}

// Coordinates for maps
export interface Coordinates {
  latitude: number
  longitude: number
}

// Time range
export interface TimeRange {
  start: string
  end: string
}

// Date range
export interface DateRange {
  start_date: string
  end_date: string
}

// Price range
export interface PriceRange {
  min: number
  max: number
  currency?: string
}

// Dimensions
export interface Dimensions {
  width: number
  height: number
  depth?: number
  unit: 'cm' | 'in' | 'mm' | 'm'
}

// Weight
export interface Weight {
  value: number
  unit: 'g' | 'kg' | 'lb' | 'oz'
}

// Currency
export interface Currency {
  code: string
  symbol: string
  name: string
  decimal_places: number
}

// Language
export interface Language {
  code: string
  name: string
  native_name: string
  flag?: string
}

// Country
export interface Country {
  code: string
  name: string
  flag?: string
  currency?: string
  phone_code?: string
}

// Notification types
export type NotificationType = 'info' | 'success' | 'warning' | 'error'

// Theme types
export type Theme = 'light' | 'dark' | 'system'

// Device types
export type DeviceType = 'mobile' | 'tablet' | 'desktop'

// Browser types
export type BrowserType = 'chrome' | 'firefox' | 'safari' | 'edge' | 'other'

// Generic callback function
export type Callback<T = void> = (data?: T) => void

// Generic async callback function
export type AsyncCallback<T = void> = (data?: T) => Promise<void>

// Event handler types
export type EventHandler<T = any> = (event: T) => void
export type AsyncEventHandler<T = any> = (event: T) => Promise<void>

// Form field types
export type FormFieldType = 
  | 'text' 
  | 'email' 
  | 'password' 
  | 'number' 
  | 'tel' 
  | 'url' 
  | 'search'
  | 'textarea'
  | 'select'
  | 'checkbox'
  | 'radio'
  | 'file'
  | 'date'
  | 'datetime-local'
  | 'time'
  | 'color'
  | 'range'

// Validation rule types
export interface ValidationRule {
  type: 'required' | 'email' | 'min' | 'max' | 'pattern' | 'custom'
  value?: any
  message: string
}

// Form field definition
export interface FormField {
  name: string
  label: string
  type: FormFieldType
  placeholder?: string
  default_value?: any
  required?: boolean
  disabled?: boolean
  readonly?: boolean
  validation_rules?: ValidationRule[]
  options?: Array<{ label: string; value: any }>
  help_text?: string
}

// Generic form data
export type FormData = Record<string, any>

// Generic component props
export interface ComponentProps {
  className?: string
  children?: React.ReactNode
  id?: string
  'data-testid'?: string
}

// Modal props
export interface ModalProps extends ComponentProps {
  isOpen: boolean
  onClose: () => void
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closeOnOverlayClick?: boolean
  closeOnEscape?: boolean
}

// Toast notification
export interface Toast {
  id: string
  type: NotificationType
  title: string
  message?: string
  duration?: number
  action?: {
    label: string
    onClick: () => void
  }
}

// Menu item
export interface MenuItem {
  id: string
  label: string
  href?: string
  icon?: string
  badge?: string
  children?: MenuItem[]
  onClick?: () => void
  disabled?: boolean
  hidden?: boolean
}

// Breadcrumb item
export interface BreadcrumbItem {
  label: string
  href?: string
  active?: boolean
}

// Tab item
export interface TabItem {
  id: string
  label: string
  content: React.ReactNode
  disabled?: boolean
  badge?: string
}

// Generic option type
export interface Option<T = any> {
  label: string
  value: T
  disabled?: boolean
  description?: string
  icon?: string
}

// Generic tree node
export interface TreeNode<T = any> {
  id: string
  label: string
  data?: T
  children?: TreeNode<T>[]
  parent_id?: string
  expanded?: boolean
  selected?: boolean
  disabled?: boolean
}

// Analytics event
export interface AnalyticsEvent {
  name: string
  properties?: Record<string, any>
  timestamp?: string
  user_id?: string
  session_id?: string
}

// Feature flag
export interface FeatureFlag {
  key: string
  enabled: boolean
  description?: string
  rollout_percentage?: number
  conditions?: Record<string, any>
}

// ===== API ERROR CLASS =====
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public code?: string,
    public details?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}
