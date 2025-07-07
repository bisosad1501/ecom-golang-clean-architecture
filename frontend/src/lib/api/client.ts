// ===== UNIFIED API CLIENT =====

import { API_CONFIG, HTTP_STATUS, STORAGE_KEYS } from '@/constants'
import type { ApiResponse, ErrorResponse } from '@/types'

// Request configuration interface
export interface RequestConfig extends RequestInit {
  timeout?: number
  retries?: number
  retryDelay?: number
  skipAuth?: boolean
}

// API Error class
export class ApiError extends Error {
  public status: number
  public code?: string
  public details?: Record<string, any>

  constructor(
    message: string,
    status: number,
    code?: string,
    details?: Record<string, any>
  ) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.code = code
    this.details = details
  }
}

// Main API Client class
export class ApiClient {
  private baseURL: string
  private defaultTimeout: number
  private defaultRetries: number
  private token: string | null = null

  constructor(
    baseURL: string = API_CONFIG.BASE_URL + API_CONFIG.API_VERSION,
    options: {
      timeout?: number
      retries?: number
    } = {}
  ) {
    this.baseURL = baseURL
    this.defaultTimeout = options.timeout || API_CONFIG.TIMEOUT.DEFAULT
    this.defaultRetries = options.retries || API_CONFIG.RETRY.ATTEMPTS
    
    // Load token from storage on initialization
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem(STORAGE_KEYS.AUTH_TOKEN)
    }
  }

  // Token management
  setToken(token: string | null) {
    this.token = token
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, token)
      } else {
        localStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
      }
    }
  }

  getToken(): string | null {
    return this.token
  }

  clearToken() {
    this.setToken(null)
  }

  // Build request headers
  private buildHeaders(config: RequestConfig = {}): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      ...config.headers,
    }

    // Add authorization header if token exists and not skipped
    if (this.token && !config.skipAuth) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    return headers
  }

  // Handle request timeout
  private createTimeoutPromise(timeout: number): Promise<never> {
    return new Promise((_, reject) => {
      setTimeout(() => {
        reject(new ApiError('Request timeout', 408, 'TIMEOUT'))
      }, timeout)
    })
  }

  // Retry logic with exponential backoff
  private async retryRequest<T>(
    requestFn: () => Promise<T>,
    retries: number,
    delay: number = API_CONFIG.RETRY.DELAY
  ): Promise<T> {
    try {
      return await requestFn()
    } catch (error) {
      if (retries > 0 && this.shouldRetry(error)) {
        await new Promise(resolve => setTimeout(resolve, delay))
        return this.retryRequest(
          requestFn,
          retries - 1,
          delay * API_CONFIG.RETRY.BACKOFF
        )
      }
      throw error
    }
  }

  // Determine if request should be retried
  private shouldRetry(error: any): boolean {
    if (error instanceof ApiError) {
      // Retry on server errors and timeout
      return error.status >= 500 || error.status === 408
    }
    return false
  }

  // Parse error response
  private async parseErrorResponse(response: Response): Promise<ApiError> {
    let errorData: ErrorResponse | null = null
    
    try {
      errorData = await response.json()
    } catch {
      // If JSON parsing fails, use status text
    }

    const message = errorData?.message || response.statusText || 'An error occurred'
    const code = errorData?.details?.code
    const details = errorData?.details

    return new ApiError(message, response.status, code, details)
  }

  // Main request method
  private async makeRequest<T>(
    endpoint: string,
    config: RequestConfig = {}
  ): Promise<T> {
    const {
      timeout = this.defaultTimeout,
      retries = this.defaultRetries,
      skipAuth = false,
      ...fetchConfig
    } = config

    const url = `${this.baseURL}${endpoint}`
    const headers = this.buildHeaders({ ...config, skipAuth })

    const requestFn = async (): Promise<T> => {
      const fetchPromise = fetch(url, {
        ...fetchConfig,
        headers,
      })

      const response = await Promise.race([
        fetchPromise,
        this.createTimeoutPromise(timeout),
      ])

      // Handle non-ok responses
      if (!response.ok) {
        throw await this.parseErrorResponse(response)
      }

      // Handle empty responses
      const contentType = response.headers.get('content-type')
      if (!contentType || !contentType.includes('application/json')) {
        return null as T
      }

      const data: ApiResponse<T> = await response.json()
      
      // Handle API response format
      if (data.success === false) {
        throw new ApiError(
          data.errors?.[0] || 'API request failed',
          response.status,
          'API_ERROR',
          data
        )
      }

      return data.data || (data as T)
    }

    return this.retryRequest(requestFn, retries)
  }

  // HTTP Methods
  async get<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'GET',
    })
  }

  async post<T>(
    endpoint: string,
    data?: any,
    config?: RequestConfig
  ): Promise<T> {
    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T>(
    endpoint: string,
    data?: any,
    config?: RequestConfig
  ): Promise<T> {
    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async patch<T>(
    endpoint: string,
    data?: any,
    config?: RequestConfig
  ): Promise<T> {
    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T>(endpoint: string, config?: RequestConfig): Promise<T> {
    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'DELETE',
    })
  }

  // File upload method
  async upload<T>(
    endpoint: string,
    formData: FormData,
    config?: Omit<RequestConfig, 'headers'>
  ): Promise<T> {
    const headers: HeadersInit = {}
    
    // Add authorization header if token exists
    if (this.token && !config?.skipAuth) {
      headers['Authorization'] = `Bearer ${this.token}`
    }

    return this.makeRequest<T>(endpoint, {
      ...config,
      method: 'POST',
      headers,
      body: formData,
    })
  }

  // Download file method
  async download(
    endpoint: string,
    filename?: string,
    config?: RequestConfig
  ): Promise<void> {
    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...config,
      headers: this.buildHeaders(config),
    })

    if (!response.ok) {
      throw await this.parseErrorResponse(response)
    }

    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    
    link.href = url
    link.download = filename || 'download'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    
    window.URL.revokeObjectURL(url)
  }

  // Health check
  async healthCheck(): Promise<{ status: string; timestamp: string }> {
    return this.get('/health', { skipAuth: true })
  }
}

// Create and export default instance
export const apiClient = new ApiClient()

// Export for testing and custom instances
export default ApiClient
