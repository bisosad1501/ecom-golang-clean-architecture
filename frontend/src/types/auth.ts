// ===== AUTHENTICATION & USER TYPES =====

import { BaseEntity, Address, ContactInfo } from './common'

// User role types
export type UserRole = 'customer' | 'admin' | 'moderator' | 'super_admin'

// User status types
export type UserStatus = 'active' | 'inactive' | 'suspended' | 'pending_verification'

// User preferences
export interface UserPreferences {
  newsletter_subscribed: boolean
  marketing_emails: boolean
  order_updates: boolean
  sms_notifications: boolean
  push_notifications: boolean
  theme: 'light' | 'dark' | 'system'
  language: string
  currency: string
  timezone: string
}

// User profile
export interface UserProfile extends BaseEntity {
  user_id: string
  phone?: string
  date_of_birth?: string
  gender?: 'male' | 'female' | 'other' | 'prefer_not_to_say'
  avatar_url?: string
  bio?: string
  website?: string
  social_links?: Record<string, string>
  preferences: UserPreferences
  contact_info?: ContactInfo
}

// Main user interface
export interface User extends BaseEntity {
  email: string
  first_name: string
  last_name: string
  username?: string
  role: UserRole
  status: UserStatus
  is_active: boolean
  email_verified: boolean
  phone_verified: boolean
  two_factor_enabled: boolean
  last_login_at?: string
  login_count: number
  profile?: UserProfile
  addresses?: Address[]
}

// Authentication response (matches backend LoginResponse)
export interface AuthResponse {
  user: User
  token: string
}

// Backend API response wrapper
export interface ApiResponse<T = any> {
  message?: string
  data?: T
  error?: string
  details?: string
}

// Login request
export interface LoginRequest {
  email: string
  password: string
  remember_me?: boolean
  device_info?: DeviceInfo
}

// Register request (matches backend RegisterRequest)
export interface RegisterRequest {
  email: string
  password: string
  first_name: string
  last_name: string
  phone?: string
}

// Extended register request for frontend forms
export interface RegisterFormRequest extends RegisterRequest {
  password_confirmation: string
  terms_accepted: boolean
  marketing_consent?: boolean
}

// Password reset request
export interface ForgotPasswordRequest {
  email: string
}

// Password reset confirmation
export interface ResetPasswordRequest {
  token: string
  email: string
  password: string
  password_confirmation: string
}

// Change password request
export interface ChangePasswordRequest {
  current_password: string
  new_password: string
  new_password_confirmation: string
}

// Email verification request
export interface VerifyEmailRequest {
  token: string
  email: string
}

// Resend verification email
export interface ResendVerificationRequest {
  email: string
}

// Two-factor authentication setup
export interface TwoFactorSetupResponse {
  secret: string
  qr_code: string
  backup_codes: string[]
}

// Two-factor authentication verification
export interface TwoFactorVerifyRequest {
  code: string
  backup_code?: string
}

// Device information
export interface DeviceInfo {
  device_type: 'mobile' | 'tablet' | 'desktop'
  device_name?: string
  browser: string
  browser_version?: string
  os: string
  os_version?: string
  ip_address?: string
  user_agent: string
  screen_resolution?: string
  timezone?: string
  language?: string
}

// Login session
export interface LoginSession extends BaseEntity {
  user_id: string
  device_info: DeviceInfo
  is_active: boolean
  last_activity_at: string
  expires_at: string
  revoked_at?: string
}

// OAuth provider types
export type OAuthProvider = 'google' | 'facebook' | 'twitter' | 'github' | 'apple'

// OAuth login request
export interface OAuthLoginRequest {
  provider: OAuthProvider
  code: string
  state?: string
  redirect_uri: string
}

// OAuth account linking
export interface OAuthAccount extends BaseEntity {
  user_id: string
  provider: OAuthProvider
  provider_user_id: string
  provider_username?: string
  provider_email?: string
  access_token?: string
  refresh_token?: string
  expires_at?: string
  scope?: string
}

// Permission types
export interface Permission {
  id: string
  name: string
  description: string
  resource: string
  action: string
}

// Role with permissions
export interface Role {
  id: string
  name: string
  description: string
  permissions: Permission[]
  is_default: boolean
}

// User with full role details
export interface UserWithRole extends User {
  role_details: Role
}

// Authentication context
export interface AuthContext {
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
  logout: () => Promise<void>
  refreshToken: () => Promise<void>
  updateProfile: (data: Partial<UserProfile>) => Promise<void>
  changePassword: (data: ChangePasswordRequest) => Promise<void>
  forgotPassword: (data: ForgotPasswordRequest) => Promise<void>
  resetPassword: (data: ResetPasswordRequest) => Promise<void>
  verifyEmail: (data: VerifyEmailRequest) => Promise<void>
  resendVerification: (data: ResendVerificationRequest) => Promise<void>
}

// JWT token payload
export interface JWTPayload {
  sub: string // user id
  email: string
  role: UserRole
  iat: number // issued at
  exp: number // expires at
  aud: string // audience
  iss: string // issuer
}

// Password strength
export interface PasswordStrength {
  score: number // 0-4
  feedback: string[]
  warning?: string
  suggestions?: string[]
}

// Account verification status
export interface VerificationStatus {
  email_verified: boolean
  phone_verified: boolean
  identity_verified: boolean
  address_verified: boolean
  payment_verified: boolean
}

// User activity log
export interface UserActivity extends BaseEntity {
  user_id: string
  action: string
  description: string
  ip_address: string
  user_agent: string
  metadata?: Record<string, any>
}

// Account deletion request
export interface AccountDeletionRequest {
  password: string
  reason?: string
  feedback?: string
}

// Export request (GDPR compliance)
export interface DataExportRequest {
  format: 'json' | 'csv' | 'pdf'
  include_orders: boolean
  include_reviews: boolean
  include_activity: boolean
}

// Privacy settings
export interface PrivacySettings {
  profile_visibility: 'public' | 'private' | 'friends'
  show_email: boolean
  show_phone: boolean
  show_activity: boolean
  allow_friend_requests: boolean
  allow_messages: boolean
  data_processing_consent: boolean
  marketing_consent: boolean
  analytics_consent: boolean
}

// Security settings
export interface SecuritySettings {
  two_factor_enabled: boolean
  login_notifications: boolean
  suspicious_activity_alerts: boolean
  session_timeout: number // minutes
  allowed_devices_limit: number
  require_password_change: boolean
  password_change_interval: number // days
}

// User invitation
export interface UserInvitation extends BaseEntity {
  email: string
  role: UserRole
  invited_by: string
  expires_at: string
  accepted_at?: string
  token: string
  message?: string
}

// Bulk user operations
export interface BulkUserOperation {
  user_ids: string[]
  action: 'activate' | 'deactivate' | 'delete' | 'change_role' | 'send_notification'
  parameters?: Record<string, any>
}

// User statistics
export interface UserStats {
  total_users: number
  active_users: number
  new_users_today: number
  new_users_this_week: number
  new_users_this_month: number
  users_by_role: Record<UserRole, number>
  users_by_status: Record<UserStatus, number>
  top_countries: Array<{ country: string; count: number }>
  registration_trend: Array<{ date: string; count: number }>
}
