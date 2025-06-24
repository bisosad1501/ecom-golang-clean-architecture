import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient, authApi } from '@/lib/api-client'
import { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  isHydrated: boolean
  error: string | null
}

interface AuthActions {
  login: (credentials: LoginRequest) => Promise<void>
  register: (data: RegisterRequest) => Promise<void>
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

          const authResponse = await authApi.login(credentials)
          const { user, token } = authResponse

          // Store token in API client
          apiClient.setToken(token)
          
          // Also store in localStorage for API client to access
          if (typeof window !== 'undefined') {
            localStorage.setItem('auth_token', token)
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

          const authResponse = await authApi.register(data)
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

      logout: () => {
        // Clear token from API client
        apiClient.setToken(null)
        
        // Clear token from localStorage
        if (typeof window !== 'undefined') {
          localStorage.removeItem('auth_token')
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
          const response = await apiClient.get<User>('/users/profile')
          const user = response.data

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
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
        if (!token) return

        try {
          set({ isLoading: true })
          
          // Set token in API client and localStorage
          apiClient.setToken(token)
          if (typeof window !== 'undefined') {
            localStorage.setItem('auth_token', token)
          }
          
          // Verify token by fetching user profile
          const response = await apiClient.get<User>('/users/profile')
          const user = response.data

          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          })
        } catch (error: any) {
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


