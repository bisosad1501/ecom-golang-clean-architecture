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
          return {
            message: data.message || ERROR_MESSAGES.VALIDATION_ERROR,
            code: 'VALIDATION_ERROR',
            details: data.details,
          }
        case 401:
          this.handleUnauthorized()
          return {
            message: ERROR_MESSAGES.UNAUTHORIZED,
            code: 'UNAUTHORIZED',
          }
        case 403:
          return {
            message: ERROR_MESSAGES.FORBIDDEN,
            code: 'FORBIDDEN',
          }
        case 404:
          return {
            message: ERROR_MESSAGES.NOT_FOUND,
            code: 'NOT_FOUND',
          }
        case 409:
          return {
            message: data.message || 'Conflict occurred',
            code: 'CONFLICT',
            details: data.details,
          }
        case 422:
          return {
            message: data.message || 'Unprocessable entity',
            code: 'UNPROCESSABLE_ENTITY',
            details: data.details,
          }
        case 429:
          return {
            message: ERROR_MESSAGES.RATE_LIMIT,
            code: 'RATE_LIMIT',
          }
        case 500:
        default:
          return {
            message: data.message || ERROR_MESSAGES.SERVER_ERROR,
            code: 'SERVER_ERROR',
            details: data.details,
          }
      }
    } else if (error.request) {
      // Network error
      return {
        message: ERROR_MESSAGES.NETWORK_ERROR,
        code: 'NETWORK_ERROR',
      }
    } else {
      // Other error
      return {
        message: error.message || ERROR_MESSAGES.SERVER_ERROR,
        code: 'UNKNOWN_ERROR',
      }
    }
  }

  private handleUnauthorized() {
    // Clear auth data and redirect to login
    if (typeof window !== 'undefined') {
      localStorage.removeItem(AUTH_TOKEN_KEY)
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('user')
      window.location.href = '/auth/login'
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
