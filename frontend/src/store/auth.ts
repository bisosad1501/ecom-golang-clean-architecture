import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient } from '@/lib/api'
import { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types'
import { AUTH_TOKEN_KEY } from '@/constants'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  isHydrated: boolean
  error: string | null
}

interface AuthActions {
  login: (credentials: LoginRequest) => Promise<User>
  register: (data: RegisterRequest) => Promise<void>
  setAuthData: (token: string, user: User) => void
  logout: () => void
  updateProfile: (data: Partial<User>) => Promise<void>
  clearError: () => void
  setLoading: (loading: boolean) => void
  setHydrated: (hydrated: boolean) => void
  checkAuth: () => Promise<void>
  refreshUser: () => Promise<void>
}

type AuthStore = AuthState & AuthActions

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // Initial state
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      isHydrated: false,
      error: null,

      // Actions
      login: async (credentials: LoginRequest) => {
        try {
          set({ isLoading: true, error: null })

          const authResponse = await apiClient.login(credentials)
          const { user, token } = authResponse

          // Store token in API client
          apiClient.setToken(token)
          
          // Also store in localStorage for API client to access
          if (typeof window !== 'undefined') {
            localStorage.setItem(AUTH_TOKEN_KEY, token)
          }

          // Log role for debugging
          console.log('User logged in with role:', user.role)

          set({
            user,
            token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })

          // Return user data for immediate use
          return user
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Login failed',
          })
          throw error
        }
      },

      register: async (data: RegisterRequest) => {
        try {
          set({ isLoading: true, error: null })

          const authResponse = await apiClient.register(data)
          const { user, token } = authResponse

          // Store token in API client
          apiClient.setToken(token)

          set({
            user,
            token,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Registration failed',
          })
          throw error
        }
      },

      setAuthData: (token: string, user: User) => {
        // Store token in API client
        apiClient.setToken(token)
        
        // Also store in localStorage for API client to access
        if (typeof window !== 'undefined') {
          localStorage.setItem(AUTH_TOKEN_KEY, token)
        }

        set({
          user,
          token,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        })
      },

      logout: () => {
        // Clear token from API client
        apiClient.clearToken()
        
        // Clear token from localStorage
        if (typeof window !== 'undefined') {
          localStorage.removeItem(AUTH_TOKEN_KEY)
        }

        set({
          user: null,
          token: null,
          isAuthenticated: false,
          isHydrated: true, // Keep hydrated state
          error: null,
        })
      },

      refreshUser: async () => {
        const { token } = get()
        if (!token) return

        try {
          set({ isLoading: true })
          
          // Set token in API client
          apiClient.setToken(token)
          
          // Fetch fresh user data
          const user = await apiClient.getUserProfile()
          console.log('refreshUser response:', user)

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          console.error('refreshUser error:', error)
          // Token is invalid, clear auth state
          get().logout()
          set({ isLoading: false })
        }
      },

      updateProfile: async (data: Partial<User>) => {
        try {
          set({ isLoading: true, error: null })
          
          const response = await apiClient.put<User>('/users/profile', data)
          const updatedUser = response.data

          set({
            user: updatedUser,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'Profile update failed',
          })
          throw error
        }
      },

      checkAuth: async () => {
        const { token } = get()
        console.log('checkAuth - starting with token:', !!token)
        if (!token) return

        try {
          set({ isLoading: true })
          
          // Set token in API client and localStorage
          apiClient.setToken(token)
          if (typeof window !== 'undefined') {
            localStorage.setItem(AUTH_TOKEN_KEY, token)
          }
          
          // Verify token by fetching user profile
          const user = await apiClient.getUserProfile()
          console.log('checkAuth extracted user:', user)

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
          
          console.log('checkAuth - auth state updated successfully')
        } catch (error: any) {
          console.error('checkAuth error:', error)
          // Token is invalid, clear auth state
          get().logout()
          set({ isLoading: false })
        }
      },

      clearError: () => set({ error: null }),
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      setHydrated: (hydrated: boolean) => set({ isHydrated: hydrated }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated,
      }),
      onRehydrateStorage: () => (state) => {
        // Set hydrated flag and check auth after rehydration
        if (state) {
          state.isHydrated = true
          state.isLoading = false
          
          // Auto-check auth after hydration if we have a token
          if (state.token) {
            state.checkAuth()
          }
        }
      },
    }
  )
)


