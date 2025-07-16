import { AuthResponse, LoginRequest, RegisterRequest } from '@/types'
import { AUTH_TOKEN_KEY } from '@/constants'

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
      this.token = localStorage.getItem(AUTH_TOKEN_KEY)
    }
  }

  // Set authentication token
  setToken(token: string | null) {
    this.token = token
    if (typeof window !== 'undefined') {
      if (token) {
        localStorage.setItem(AUTH_TOKEN_KEY, token)
      } else {
        localStorage.removeItem(AUTH_TOKEN_KEY)
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
          error.status = response.status
          throw error
        }

        // Handle authentication errors
        if (response.status === 401) {
          this.setToken(null) // Clear invalid token
          // Redirect to login if in browser
          if (typeof window !== 'undefined') {
            window.location.href = '/auth/login'
          }
          const error = new Error('Authentication required') as any
          error.code = 'UNAUTHORIZED'
          error.status = response.status
          throw error
        }

        // Handle authorization errors
        if (response.status === 403) {
          const error = new Error('Access denied') as any
          error.code = 'FORBIDDEN'
          error.status = response.status
          throw error
        }

        // Handle not found errors
        if (response.status === 404) {
          const error = new Error('Resource not found') as any
          error.code = 'NOT_FOUND'
          error.status = response.status
          throw error
        }

        // Handle rate limiting
        if (response.status === 429) {
          const error = new Error('Too many requests. Please try again later.') as any
          error.code = 'RATE_LIMITED'
          error.status = response.status
          throw error
        }

        // Handle server errors
        if (response.status >= 500) {
          const error = new Error('Server error. Please try again later.') as any
          error.code = 'SERVER_ERROR'
          error.status = response.status
          throw error
        }

        // Generic error
        const error = new Error(errorMessage) as any
        error.code = 'API_ERROR'
        error.status = response.status
        throw error
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
  oauthLogin: async (data: import('@/types/auth').OAuthLoginRequest): Promise<AuthResponse> => {
    // Note: Backend OAuth flow is different - it uses callbacks, not direct API calls
    // This method might need to be updated based on actual OAuth implementation
    const response = await apiClient.post<any>(`/auth/${data.provider}/callback`, data)

    // Backend returns: { "message": "...", "data": { "user": {...}, "token": "..." } }
    const authData = response.data.data
    if (!authData || !authData.user || !authData.token) {
      throw new Error('Invalid OAuth response format')
    }

    return authData
  },
  login: async (credentials: LoginRequest): Promise<AuthResponse> => {
    const response = await apiClient.post<any>('/auth/login', credentials)
    console.log('authApi.login - Raw response:', response)

    // ApiClient already extracts the 'data' field from backend response
    // Backend returns: { "message": "Login successful", "data": { "user": {...}, "token": "..." } }
    // ApiClient extracts data.data and returns: { data: { "user": {...}, "token": "..." }, response: Response }
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
    const response = await apiClient.post<any>('/auth/register', userData)
    console.log('authApi.register - Raw response:', response)

    // Backend returns: { "message": "User registered successfully", "data": { "id": "...", "email": "...", ... } }
    // ApiClient extracts data.data and returns: { data: { user data }, response: Response }
    // Note: Register doesn't return token, just user data
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
    const response = await apiClient.get<any>('/users/profile')
    console.log('authApi.getProfile - Raw response:', response)

    // Backend returns: { "message": "...", "data": { user data } }
    // ApiClient extracts data.data and returns: { data: { user data }, response: Response }
    const userData = response.data
    if (!userData) {
      throw new Error('Invalid profile response format')
    }

    console.log('authApi.getProfile - User data:', userData)
    return userData
  },
}

// Products API
export const productsApi = {
  getAll: async (params?: any): Promise<any> => {
    const queryParams = new URLSearchParams()
    if (params?.limit) queryParams.append('limit', params.limit.toString())
    if (params?.offset) queryParams.append('offset', params.offset.toString())
    if (params?.category_id) queryParams.append('category_id', params.category_id)
    if (params?.search) queryParams.append('search', params.search)
    if (params?.sort) queryParams.append('sort', params.sort)
    if (params?.order) queryParams.append('order', params.order)

    const url = `/products${queryParams.toString() ? '?' + queryParams.toString() : ''}`
    const response = await apiClient.get<any>(url)

    // Backend returns: { "message": "...", "data": { "products": [...], "total": 0, "limit": 10, "offset": 0 } }
    return response.data.data
  },

  getById: async (id: string): Promise<any> => {
    const response = await apiClient.get<any>(`/products/${id}`)
    // Backend returns: { "message": "...", "data": { product data } }
    return response.data.data
  },

  create: async (productData: any): Promise<any> => {
    const response = await apiClient.post<any>('/admin/products', productData)
    // Backend returns: { "message": "...", "data": { product data } }
    return response.data.data
  },

  update: async (id: string, productData: any): Promise<any> => {
    const response = await apiClient.put<any>(`/admin/products/${id}`, productData)
    // Backend returns: { "message": "...", "data": { product data } }
    return response.data.data
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/admin/products/${id}`)
  },
}

// Categories API
export const categoriesApi = {
  getAll: async (): Promise<any> => {
    const response = await apiClient.get<any>('/categories')
    // Backend returns: { "message": "...", "data": [...] }
    return response.data.data
  },

  getById: async (id: string): Promise<any> => {
    const response = await apiClient.get<any>(`/categories/${id}`)
    // Backend returns: { "message": "...", "data": { category data } }
    return response.data.data
  },
}

// Export default client
export default apiClient
