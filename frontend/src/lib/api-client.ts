import { AuthResponse, LoginRequest, RegisterRequest } from '@/types'

// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

// API Response wrapper
interface ApiResponse<T = any> {
  message?: string
  data?: T
  error?: string
  errors?: Record<string, string>
}

// HTTP Client class
class ApiClient {
  private baseURL: string
  private token: string | null = null

  constructor(baseURL: string) {
    this.baseURL = baseURL
    
    // Load token from localStorage on initialization
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('auth_token')
    }
  }

  // Set authentication token
  setToken(token: string | null) {
    this.token = token
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem('auth_token', token)
      } else {
        localStorage.removeItem('auth_token')
      }
    }
  }

  // Get authentication token
  getToken(): string | null {
    return this.token
  }

  // Make HTTP request
  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<{ data: T; response: Response }> {
    const url = `${this.baseURL}${endpoint}`
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string> || {}),
    }

    // Add authorization header if token exists
    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`
      console.log('ApiClient: Using token for request:', this.token.substring(0, 20) + '...')
    } else {
      console.log('ApiClient: No token available for request to:', endpoint)
    }

    const config: RequestInit = {
      ...options,
      headers,
    }

    console.log('ApiClient: Making request to:', url, 'with method:', options.method || 'GET')

    try {
      const response = await fetch(url, config)
      
      // Parse response
      let data: ApiResponse<T>
      try {
        data = await response.json()
      } catch {
        data = {} as ApiResponse<T>
      }

      console.log('ApiClient: Response status:', response.status, 'data:', data)

      // Handle HTTP errors
      if (!response.ok) {
        const errorMessage = data.error || data.message || `HTTP ${response.status}: ${response.statusText}`
        
        // Handle validation errors
        if (response.status === 400 && data.errors) {
          const error = new Error(errorMessage) as any
          error.code = 'VALIDATION_ERROR'
          error.details = data.errors
          throw error
        }
        
        // Handle authentication errors
        if (response.status === 401) {
          this.setToken(null) // Clear invalid token
          const error = new Error('Authentication required') as any
          error.code = 'UNAUTHORIZED'
          throw error
        }
        
        throw new Error(errorMessage)
      }

      return { data: (data.data || data) as T, response }
    } catch (error) {
      console.error(`API Error [${options.method || 'GET'} ${url}]:`, error)
      throw error
    }
  }

  // HTTP Methods
  async get<T>(endpoint: string, params?: Record<string, any>): Promise<{ data: T; response: Response }> {
    const url = new URL(`${this.baseURL}${endpoint}`)
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          url.searchParams.append(key, String(value))
        }
      })
    }
    
    return this.request<T>(endpoint + url.search, { method: 'GET' })
  }

  async post<T>(endpoint: string, data?: any): Promise<{ data: T; response: Response }> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T>(endpoint: string, data?: any): Promise<{ data: T; response: Response }> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async patch<T>(endpoint: string, data?: any): Promise<{ data: T; response: Response }> {
    return this.request<T>(endpoint, {
      method: 'PATCH',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T>(endpoint: string): Promise<{ data: T; response: Response }> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }
}

// Create API client instance
export const apiClient = new ApiClient(API_BASE_URL)

// Authentication API functions
export const authApi = {
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/login', credentials)
    return data
  },

  register: async (userData: RegisterRequest): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/register', userData)
    return data
  },

  logout: async (): Promise<void> => {
    try {
      await apiClient.post('/auth/logout')
    } catch (error) {
      // Ignore logout errors, just clear local token
      console.warn('Logout API call failed:', error)
    } finally {
      apiClient.setToken(null)
    }
  },

  refreshToken: async (): Promise<AuthResponse> => {
    const { data } = await apiClient.post<AuthResponse>('/auth/refresh')
    return data
  },

  getProfile: async (): Promise<any> => {
    const { data } = await apiClient.get('/users/profile')
    return data
  },
}

// Export default client
export default apiClient
