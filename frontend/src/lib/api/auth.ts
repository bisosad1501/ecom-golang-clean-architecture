// ===== AUTHENTICATION API SERVICE =====

import { apiClient } from './client'
import { API_CONFIG } from '@/constants'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  ForgotPasswordRequest,
  ResetPasswordRequest,
  ChangePasswordRequest,
  VerifyEmailRequest,
  ResendVerificationRequest,
  TwoFactorSetupResponse,
  TwoFactorVerifyRequest,
  OAuthLoginRequest,
  User,
  UserProfile,
} from '@/types'

export class AuthService {
  // Login with email and password
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.ENDPOINTS.AUTH.LOGIN,
      credentials,
      { skipAuth: true }
    )
    
    // Set token in client
    if (response.token) {
      apiClient.setToken(response.token)
    }
    
    return response
  }

  // Register new user
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.ENDPOINTS.AUTH.REGISTER,
      data,
      { skipAuth: true }
    )
    
    // Set token in client if auto-login after registration
    if (response.token) {
      apiClient.setToken(response.token)
    }
    
    return response
  }

  // Logout user
  async logout(): Promise<void> {
    try {
      await apiClient.post(API_CONFIG.ENDPOINTS.AUTH.LOGOUT)
    } finally {
      // Always clear token, even if request fails
      apiClient.clearToken()
    }
  }

  // Refresh authentication token
  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      API_CONFIG.ENDPOINTS.AUTH.REFRESH
    )
    
    if (response.token) {
      apiClient.setToken(response.token)
    }
    
    return response
  }

  // Get current user profile
  async getProfile(): Promise<User> {
    return apiClient.get<User>(API_CONFIG.ENDPOINTS.AUTH.PROFILE)
  }

  // Update user profile
  async updateProfile(data: Partial<UserProfile>): Promise<User> {
    return apiClient.put<User>(API_CONFIG.ENDPOINTS.AUTH.PROFILE, data)
  }

  // Change password
  async changePassword(data: ChangePasswordRequest): Promise<void> {
    return apiClient.post(API_CONFIG.ENDPOINTS.USERS.CHANGE_PASSWORD, data)
  }

  // Forgot password - send reset email
  async forgotPassword(data: ForgotPasswordRequest): Promise<void> {
    return apiClient.post(
      API_CONFIG.ENDPOINTS.AUTH.FORGOT_PASSWORD,
      data,
      { skipAuth: true }
    )
  }

  // Reset password with token
  async resetPassword(data: ResetPasswordRequest): Promise<void> {
    return apiClient.post(
      API_CONFIG.ENDPOINTS.AUTH.RESET_PASSWORD,
      data,
      { skipAuth: true }
    )
  }

  // Verify email address
  async verifyEmail(data: VerifyEmailRequest): Promise<void> {
    return apiClient.post(
      API_CONFIG.ENDPOINTS.AUTH.VERIFY_EMAIL,
      data,
      { skipAuth: true }
    )
  }

  // Resend email verification
  async resendVerification(data: ResendVerificationRequest): Promise<void> {
    return apiClient.post(
      API_CONFIG.ENDPOINTS.AUTH.VERIFY_EMAIL,
      data,
      { skipAuth: true }
    )
  }

  // Two-factor authentication setup
  async setupTwoFactor(): Promise<TwoFactorSetupResponse> {
    return apiClient.post<TwoFactorSetupResponse>('/auth/2fa/setup')
  }

  // Verify two-factor authentication
  async verifyTwoFactor(data: TwoFactorVerifyRequest): Promise<void> {
    return apiClient.post('/auth/2fa/verify', data)
  }

  // Disable two-factor authentication
  async disableTwoFactor(data: { password: string }): Promise<void> {
    return apiClient.post('/auth/2fa/disable', data)
  }

  // OAuth login (Google, Facebook, etc.)
  async oauthLogin(data: OAuthLoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>(
      `/auth/oauth/${data.provider}`,
      data,
      { skipAuth: true }
    )
    
    if (response.token) {
      apiClient.setToken(response.token)
    }
    
    return response
  }

  // Get OAuth authorization URL
  async getOAuthUrl(provider: string, redirectUri: string): Promise<{ url: string }> {
    return apiClient.get<{ url: string }>(
      `/auth/oauth/${provider}/url?redirect_uri=${encodeURIComponent(redirectUri)}`,
      { skipAuth: true }
    )
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return !!apiClient.getToken()
  }

  // Get current token
  getToken(): string | null {
    return apiClient.getToken()
  }

  // Validate token (check if still valid)
  async validateToken(): Promise<boolean> {
    try {
      await this.getProfile()
      return true
    } catch {
      return false
    }
  }

  // Get user sessions
  async getSessions(): Promise<any[]> {
    return apiClient.get('/auth/sessions')
  }

  // Revoke session
  async revokeSession(sessionId: string): Promise<void> {
    return apiClient.delete(`/auth/sessions/${sessionId}`)
  }

  // Revoke all sessions except current
  async revokeAllSessions(): Promise<void> {
    return apiClient.post('/auth/sessions/revoke-all')
  }

  // Account deletion
  async deleteAccount(data: { password: string; reason?: string }): Promise<void> {
    return apiClient.post('/auth/delete-account', data)
  }

  // Export user data (GDPR compliance)
  async exportData(format: 'json' | 'csv' = 'json'): Promise<void> {
    return apiClient.download(`/auth/export-data?format=${format}`, `user-data.${format}`)
  }

  // Check email availability
  async checkEmailAvailability(email: string): Promise<{ available: boolean }> {
    return apiClient.get<{ available: boolean }>(
      `/auth/check-email?email=${encodeURIComponent(email)}`,
      { skipAuth: true }
    )
  }

  // Check username availability
  async checkUsernameAvailability(username: string): Promise<{ available: boolean }> {
    return apiClient.get<{ available: boolean }>(
      `/auth/check-username?username=${encodeURIComponent(username)}`,
      { skipAuth: true }
    )
  }

  // Send test email (for development)
  async sendTestEmail(email: string): Promise<void> {
    return apiClient.post(
      '/auth/send-test-email',
      { email },
      { skipAuth: true }
    )
  }
}

// Create and export service instance
export const authService = new AuthService()

// Export for testing and custom instances
export default AuthService
