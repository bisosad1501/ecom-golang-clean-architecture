// ===== AUTHENTICATION SERVICE =====

import { apiClient } from '@/lib/api'
import type {
  AuthResponse,
  LoginRequest,
  RegisterRequest,
  User,
} from '@/types'

export class AuthService {
  // Authentication methods
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/login', credentials)
    return response.data!
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/register', userData)
    return response.data!
  }

  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      // Ignore logout errors, just clear local token
      console.warn('Logout API call failed:', error)
    } finally {
      apiClient.clearToken()
    }
  }

  async refreshToken(): Promise<AuthResponse> {
    const response = await apiClient.post<AuthResponse>('/auth/refresh')
    return response.data!
  }

  async getProfile(): Promise<User> {
    const response = await apiClient.get<User>('/users/profile')
    return response.data!
  }

  async updateProfile(userData: Partial<User>): Promise<User> {
    const response = await apiClient.put<User>('/users/profile', userData)
    return response.data!
  }

  async changePassword(data: { current_password: string; new_password: string }): Promise<void> {
    await apiClient.post('/auth/change-password', data)
  }

  async forgotPassword(email: string): Promise<void> {
    await apiClient.post('/auth/forgot-password', { email })
  }

  async resetPassword(data: { token: string; password: string }): Promise<void> {
    await apiClient.post('/auth/reset-password', data)
  }

  async verifyEmail(token: string): Promise<void> {
    await apiClient.post('/auth/verify-email', { token })
  }

  async resendVerification(email: string): Promise<void> {
    await apiClient.post('/auth/resend-verification', { email })
  }
}

// Create and export singleton instance
export const authService = new AuthService()
export default authService
