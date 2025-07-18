import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { API_BASE_URL, AUTH_TOKEN_KEY, ERROR_MESSAGES } from '@/constants'
import { ApiResponse, ApiError, AuthResponse, LoginRequest, RegisterRequest, User } from '@/types'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor to add auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = this.getToken()
        console.log('API Client - Request URL:', config.url)
        console.log('API Client - Token found:', token ? 'Yes' : 'No')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
          console.log('API Client - Authorization header set')
        }
        return config
      },
      (error) => {
        console.error('API Client - Request error:', error)
        return Promise.reject(error)
      }
    )

    // Response interceptor to handle errors
    this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        console.log('API Client - Response status:', response.status)
        console.log('API Client - Response data:', response.data)
        return response
      },
      (error) => {
        console.error('API Client - Response error:', error)
        const apiError = this.handleError(error)
        return Promise.reject(apiError)
      }
    )
  }

  private getToken(): string | null {
    if (typeof window !== 'undefined') {
      return localStorage.getItem(AUTH_TOKEN_KEY)
    }
    return null
  }

  private handleError(error: any): ApiError {
    if (error.response) {
      // Server responded with error status
      const { status, data } = error.response

      switch (status) {
        case 400:
          return new ApiError(
            data.message || ERROR_MESSAGES.VALIDATION_ERROR,
            status,
            'VALIDATION_ERROR',
            data.details
          )
        case 401:
          this.handleUnauthorized()
          return new ApiError(
            ERROR_MESSAGES.UNAUTHORIZED,
            status,
            'UNAUTHORIZED'
          )
        case 403:
          return new ApiError(
            ERROR_MESSAGES.FORBIDDEN,
            status,
            'FORBIDDEN'
          )
        case 404:
          return new ApiError(
            ERROR_MESSAGES.NOT_FOUND,
            status,
            'NOT_FOUND'
          )
        case 409:
          return new ApiError(
            data.message || 'Conflict occurred',
            status,
            'CONFLICT',
            data.details
          )
        case 422:
          return new ApiError(
            data.message || 'Unprocessable entity',
            status,
            'UNPROCESSABLE_ENTITY',
            data.details
          )
        case 429:
          return new ApiError(
            ERROR_MESSAGES.RATE_LIMIT,
            status,
            'RATE_LIMIT'
          )
        case 500:
        default:
          return new ApiError(
            data.message || ERROR_MESSAGES.SERVER_ERROR,
            status,
            'SERVER_ERROR',
            data.details
          )
      }
    } else if (error.request) {
      // Network error
      return new ApiError(
        ERROR_MESSAGES.NETWORK_ERROR,
        0,
        'NETWORK_ERROR'
      )
    } else {
      // Other error
      return new ApiError(
        error.message || ERROR_MESSAGES.SERVER_ERROR,
        0,
        'UNKNOWN_ERROR'
      )
    }
  }

  private handleUnauthorized() {
    // Clear auth data but don't redirect automatically
    // Let the auth store handle the redirect logic
    if (typeof window !== 'undefined') {
      localStorage.removeItem(AUTH_TOKEN_KEY)
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')

      // Clear auth store state
      import('@/store/auth').then(({ useAuthStore }) => {
        const authStore = useAuthStore.getState()
        authStore.logout()
      })
    }
  }

  // Generic request methods
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.get(url, config)
    return response.data
  }

  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.post(url, data, config)
    return response.data
  }

  async put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.put(url, data, config)
    return response.data
  }

  async patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.patch(url, data, config)
    return response.data
  }

  async delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    const response = await this.client.delete(url, config)
    return response.data
  }

  // File upload method
  async upload<T = any>(url: string, file: File, onProgress?: (progress: number) => void): Promise<ApiResponse<T>> {
    const formData = new FormData()
    formData.append('file', file)

    const config: AxiosRequestConfig = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      },
    }

    const response = await this.client.post(url, formData, config)
    return response.data
  }

  // Multiple file upload
  async uploadMultiple<T = any>(
    url: string, 
    files: File[], 
    onProgress?: (progress: number) => void
  ): Promise<ApiResponse<T>> {
    const formData = new FormData()
    files.forEach((file, index) => {
      formData.append(`files[${index}]`, file)
    })

    const config: AxiosRequestConfig = {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          onProgress(progress)
        }
      },
    }

    const response = await this.client.post(url, formData, config)
    return response.data
  }

  // Download file
  async download(url: string, filename?: string): Promise<void> {
    const response = await this.client.get(url, {
      responseType: 'blob',
    })

    const blob = new Blob([response.data])
    const downloadUrl = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = downloadUrl
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(downloadUrl)
  }

  // Set auth token
  setToken(token: string) {
    if (typeof window !== 'undefined') {
      localStorage.setItem(AUTH_TOKEN_KEY, token)
    }
  }

  // Clear auth token
  clearToken() {
    if (typeof window !== 'undefined') {
      localStorage.removeItem(AUTH_TOKEN_KEY)
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
    }
  }

  // Payment-specific methods
  async createCheckoutSession(data: {
    order_id: string;
    amount: number;
    currency: string;
    description?: string;
    success_url: string;
    cancel_url: string;
    metadata?: Record<string, any>;
  }): Promise<ApiResponse<{
    success: boolean;
    session_id: string;
    session_url: string;
    message: string;
  }>> {
    return this.post('/payments/checkout-session', data)
  }

  async processPayment(data: {
    order_id: string;
    amount: number;
    currency: string;
    payment_method: string;
    payment_token?: string;
    metadata?: Record<string, any>;
  }): Promise<ApiResponse<any>> {
    return this.post('/payments', data)
  }

  async getPayment(paymentId: string): Promise<ApiResponse<any>> {
    return this.get(`/payments/${paymentId}`)
  }

  async processRefund(paymentId: string, data: {
    amount?: number;
    reason?: string;
  }): Promise<ApiResponse<any>> {
    return this.post(`/payments/${paymentId}/refund`, data)
  }

  async getPaymentMethods(): Promise<ApiResponse<any[]>> {
    return this.get('/payments/methods')
  }

  async savePaymentMethod(data: {
    type: string;
    provider: string;
    token: string;
    metadata?: Record<string, any>;
  }): Promise<ApiResponse<any>> {
    return this.post('/payments/methods', data)
  }

  async deletePaymentMethod(methodId: string): Promise<ApiResponse<any>> {
    return this.delete(`/payments/methods/${methodId}`)
  }

  async setDefaultPaymentMethod(methodId: string): Promise<ApiResponse<any>> {
    return this.put(`/payments/methods/${methodId}/default`, {})
  }

  // Order-specific methods
  async createOrder(data: {
    shipping_address: {
      street: string;
      city: string;
      state: string;
      zip_code: string;
      country: string;
    };
    billing_address: {
      street: string;
      city: string;
      state: string;
      zip_code: string;
      country: string;
    };
    coupon_code?: string;
    notes?: string;
  }): Promise<ApiResponse<any>> {
    return this.post('/orders', data)
  }

  async getOrders(params?: {
    page?: number;
    limit?: number;
    status?: string;
  }): Promise<ApiResponse<any[]>> {
    return this.get('/orders', { params })
  }

  async getOrder(orderId: string): Promise<ApiResponse<any>> {
    return this.get(`/orders/${orderId}`)
  }

  async cancelOrder(orderId: string, reason?: string): Promise<ApiResponse<any>> {
    return this.post(`/orders/${orderId}/cancel`, { reason })
  }

  // Cart-specific methods
  async getCart(): Promise<ApiResponse<any>> {
    return this.get('/cart')
  }

  async addToCart(data: {
    product_id: string;
    quantity: number;
    variant_id?: string;
  }): Promise<ApiResponse<any>> {
    return this.post('/cart/items', data)
  }

  async updateCartItem(data: {
    product_id: string;
    quantity: number;
    variant_id?: string;
  }): Promise<ApiResponse<any>> {
    return this.put('/cart/items', data)
  }

  async removeFromCart(productId: string): Promise<ApiResponse<any>> {
    return this.delete(`/cart/items/${productId}`)
  }

  async clearCart(): Promise<ApiResponse<any>> {
    return this.delete('/cart')
  }

  // Authentication methods
  async login(credentials: LoginRequest): Promise<AuthResponse> {
    const response = await this.post<AuthResponse>('/auth/login', credentials)
    return response.data!
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    const response = await this.post<AuthResponse>('/auth/register', userData)
    return response.data!
  }

  async logout(): Promise<void> {
    try {
      await this.post('/auth/logout')
    } catch (error) {
      // Ignore logout errors, just clear local token
      console.warn('Logout API call failed:', error)
    } finally {
      this.clearToken()
    }
  }

  // Compatibility methods for api-client.ts replacement
  // These methods return the same format as api-client.ts for seamless migration
  async getCompat<T>(endpoint: string, params?: Record<string, any>): Promise<{ data: T; response: any }> {
    const response = await this.get<T>(endpoint, { params })
    return {
      data: response.data!,
      response: { status: 200 } // Mock response object for compatibility
    }
  }

  async postCompat<T>(endpoint: string, data?: any): Promise<{ data: T; response: any }> {
    const response = await this.post<T>(endpoint, data)
    return {
      data: response.data!,
      response: { status: 200 }
    }
  }

  // Additional compatibility methods for layout components
  async putCompat<T>(endpoint: string, data?: any): Promise<{ data: T; response: any }> {
    const response = await this.put<T>(endpoint, data)
    return {
      data: response.data!,
      response: { status: 200 }
    }
  }

  async patchCompat<T>(endpoint: string, data?: any): Promise<{ data: T; response: any }> {
    const response = await this.patch<T>(endpoint, data)
    return {
      data: response.data!,
      response: { status: 200 }
    }
  }

  async deleteCompat<T>(endpoint: string): Promise<{ data: T; response: any }> {
    const response = await this.delete<T>(endpoint)
    return {
      data: response.data!,
      response: { status: 200 }
    }
  }

  async refreshToken(refreshToken: string): Promise<{ token: string; refresh_token: string }> {
    const response = await this.post<{ token: string; refresh_token: string }>('/auth/refresh', {
      refresh_token: refreshToken,
    })
    return response.data!
  }

  async getUserProfile(): Promise<User> {
    const response = await this.get<{ data: User }>('/users/profile')
    // Handle nested data structure from backend
    return response.data?.data || response.data
  }

  // Get raw axios instance for custom requests
  getClient(): AxiosInstance {
    return this.client
  }
}

// Create singleton instance
export const apiClient = new ApiClient()

// Compatibility layer: authApi that matches api-client.ts interface
export const authApi = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    console.log('authApi.login - Using Axios client')
    const response = await apiClient.post<any>('/auth/login', credentials)
    console.log('authApi.login - Raw response:', response)

    // Extract auth data from Axios response
    const authData = response.data
    console.log('authApi.login - Auth data:', authData)

    if (!authData) {
      throw new Error('No auth data in response')
    }

    if (!authData.user) {
      console.error('Missing user in authData:', authData)
      throw new Error('No user data in response')
    }

    if (!authData.token) {
      console.error('Missing token in authData:', authData)
      throw new Error('No token in response')
    }

    return authData
  },

  register: async (userData: RegisterRequest): Promise<any> => {
    console.log('authApi.register - Using Axios client')
    const response = await apiClient.post<any>('/auth/register', userData)
    console.log('authApi.register - Raw response:', response)

    const userData_response = response.data
    if (!userData_response) {
      throw new Error('Invalid register response format')
    }

    console.log('authApi.register - User data:', userData_response)
    return userData_response
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      console.warn('Logout API call failed:', error)
    } finally {
      apiClient.clearToken()
    }
  },

  refreshToken: async (): Promise<AuthResponse> => {
    const response = await apiClient.post<AuthResponse>('/auth/refresh')
    return response.data!
  },

  getProfile: async (): Promise<any> => {
    console.log('authApi.getProfile - Using Axios client')
    const response = await apiClient.get<any>('/users/profile')
    console.log('authApi.getProfile - Raw response:', response)

    const userData = response.data
    if (!userData) {
      throw new Error('Invalid profile response format')
    }

    console.log('authApi.getProfile - User data:', userData)
    return userData
  },

  oauthLogin: async (data: import('@/types/auth').OAuthLoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<any>(`/auth/${data.provider}/callback`, data)
    const authData = response.data

    if (!authData || !authData.user || !authData.token) {
      throw new Error('Invalid OAuth response format')
    }

    return authData
  },

  // Email verification methods
  verifyEmail: async (token: string): Promise<any> => {
    console.log('authApi.verifyEmail - Using Axios client')
    const response = await apiClient.get<any>(`/auth/verify-email?token=${token}`)
    console.log('authApi.verifyEmail - Raw response:', response)

    const verificationData = response.data
    if (!verificationData) {
      throw new Error('Invalid verification response format')
    }

    console.log('authApi.verifyEmail - Verification data:', verificationData)
    return verificationData
  },

  sendEmailVerification: async (): Promise<any> => {
    console.log('authApi.sendEmailVerification - Using Axios client')
    const response = await apiClient.post<any>('/users/verification/email/send', {})
    console.log('authApi.sendEmailVerification - Raw response:', response)

    const result = response.data
    if (!result) {
      throw new Error('Invalid send verification response format')
    }

    console.log('authApi.sendEmailVerification - Result:', result)
    return result
  },

  resendVerification: async (email: string): Promise<any> => {
    console.log('authApi.resendVerification - Using Axios client')
    const response = await apiClient.post<any>('/auth/resend-verification', { email })
    console.log('authApi.resendVerification - Raw response:', response)

    const result = response.data
    if (!result) {
      throw new Error('Invalid resend verification response format')
    }

    console.log('authApi.resendVerification - Result:', result)
    return result
  },
}

// Helper function to build query string
export function buildQueryString(params: Record<string, any>): string {
  const searchParams = new URLSearchParams()
  
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      if (Array.isArray(value)) {
        value.forEach((item) => searchParams.append(key, item.toString()))
      } else {
        searchParams.append(key, value.toString())
      }
    }
  })
  
  return searchParams.toString()
}

// Helper function for paginated requests
export async function getPaginated<T>(
  url: string,
  params: {
    page?: number
    limit?: number
    [key: string]: any
  } = {}
) {
  const queryString = buildQueryString(params)
  const fullUrl = queryString ? `${url}?${queryString}` : url
  return apiClient.get<T>(fullUrl)
}

// Export default instance
export default apiClient

// Compatibility exports for api-client.ts replacement
// This allows seamless migration from api-client.ts to api.ts
export { apiClient as apiClientCompat }

// Additional compatibility exports
export const productsApi = {
  getAll: async (params?: any): Promise<any> => {
    const queryParams = new URLSearchParams()
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.offset) queryParams.append('offset', params.offset.toString())
    // Backend supports category_id for search (uses ProductCategory many-to-many internally)
    if (params?.category_id) queryParams.append('category_id', params.category_id)
    if (params?.search) queryParams.append('search', params.search)
    if (params?.sort) queryParams.append('sort', params.sort)
    if (params?.order) queryParams.append('order', params.order)

    const url = `/products${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    const response = await apiClient.get<any>(url)
    return response.data
  },

  getById: async (id: string): Promise<any> => {
    const response = await apiClient.get<any>(`/products/${id}`)
    return response.data
  },

  create: async (productData: any): Promise<any> => {
    const response = await apiClient.post<any>('/admin/products', productData)
    return response.data
  },

  update: async (id: string, productData: any): Promise<any> => {
    const response = await apiClient.put<any>(`/admin/products/${id}`, productData)
    return response.data
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/products/${id}`)
  },
}

export const categoriesApi = {
  getAll: async (): Promise<any> => {
    const response = await apiClient.get<any>('/categories')
    return response.data
  },

  getById: async (id: string): Promise<any> => {
    const response = await apiClient.get<any>(`/categories/${id}`)
    return response.data
  },
}
