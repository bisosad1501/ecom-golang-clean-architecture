import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient, authApi } from '@/lib/api'
import { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types'
import type { OAuthLoginRequest } from '@/types/auth'
import { AUTH_TOKEN_KEY } from '@/constants'

interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
  isHydrated: boolean
  error: string | null
  pendingCartConflict: any | null  // Store conflict info for modal
}

interface AuthActions {
  login: (credentials: LoginRequest) => Promise<User>
  oauthLogin: (data: OAuthLoginRequest) => Promise<User>
  register: (data: RegisterRequest) => Promise<void>
  setAuthData: (token: string, user: User) => void
  logout: () => void
  updateProfile: (data: Partial<User>) => Promise<void>
  clearError: () => void
  setLoading: (loading: boolean) => void
  setHydrated: (hydrated: boolean) => void
  checkAuth: () => Promise<void>
  refreshUser: () => Promise<void>
  setPendingCartConflict: (conflict: any) => void
  clearPendingCartConflict: () => void
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
      pendingCartConflict: null,

      // Actions
      oauthLogin: async (data: OAuthLoginRequest) => {
        try {
          set({ isLoading: true, error: null })

          const authResponse = await authApi.oauthLogin(data)
          const { user, token } = authResponse

          apiClient.setToken(token)
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

          // Check for cart conflicts before merging
          try {
            const { useCartStore } = await import('@/store/cart')
            const cartStore = useCartStore.getState()
            const conflict = await cartStore.checkMergeConflict()
            console.log('🔍 Conflict check result:', conflict)
            if (conflict && (conflict.guest_cart_exists || conflict.user_cart_exists)) {
              console.log('🔍 Cart merge needed, showing modal for user choice')
              set({ pendingCartConflict: conflict })
            } else {
              console.log('ℹ️ No guest cart to merge')
            }
          } catch (cartError) {
            console.error('❌ Failed to handle cart merge:', cartError)
          }
          return user
        } catch (error: any) {
          set({
            isLoading: false,
            error: error.message || 'OAuth login failed',
          })
          throw error
        }
      },
      login: async (credentials: LoginRequest) => {
        try {
          set({ isLoading: true, error: null })

          const authResponse = await authApi.login(credentials)
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

          // Check for cart conflicts before merging
          try {
            const { useCartStore } = await import('@/store/cart')
            const cartStore = useCartStore.getState()

            console.log('🔍 Checking for cart conflicts...')

            // Check if there are conflicts
            const conflict = await cartStore.checkMergeConflict()
            console.log('🔍 Conflict check result:', conflict)

            if (conflict && (conflict.guest_cart_exists || conflict.user_cart_exists)) {
              console.log('🔍 Cart merge needed, showing modal for user choice')
              // Always show modal when there's any cart to merge, let user decide
              set({ pendingCartConflict: conflict })
            } else {
              console.log('ℹ️ No guest cart to merge')
            }
          } catch (cartError) {
            console.error('❌ Failed to handle cart merge:', cartError)
            // Don't fail login if cart merge fails
          }

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

          // Backend register only returns user data, not token
          const userData = await authApi.register(data)

          set({
            isLoading: false,
            error: null,
          })

          // Registration successful, but user needs to login separately
          // or we could auto-login them here by calling login
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
        apiClient.setToken(null)
        
        // Clear token from localStorage
        if (typeof window !== 'undefined') {
          localStorage.removeItem(AUTH_TOKEN_KEY)
        }

        // IMPORTANT: Clear cart when user logs out to prevent cart sharing
        // Import cart store and clear it
        import('@/store/cart').then(({ useCartStore }) => {
          const cartStore = useCartStore.getState()
          cartStore.clearCartLocal() // Use local clear to avoid API call during logout
          // Also clear the cart from localStorage immediately
          if (typeof window !== 'undefined') {
            localStorage.removeItem('cart-storage')
          }
        })

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
          const user = await authApi.getProfile()
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
        if (!token) {
          set({ isLoading: false })
          return
        }

        try {
          set({ isLoading: true })

          // Set token in API client and localStorage
          apiClient.setToken(token)
          if (typeof window !== 'undefined') {
            localStorage.setItem(AUTH_TOKEN_KEY, token)
          }

          // Verify token by fetching user profile
          const user = await authApi.getProfile()
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
          // Token is invalid, clear auth state silently (don't call logout to avoid redirect)
          apiClient.setToken(null)
          if (typeof window !== 'undefined') {
            localStorage.removeItem(AUTH_TOKEN_KEY)
          }

          set({
            user: null,
            token: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          })
        }
      },

      clearError: () => set({ error: null }),
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      setHydrated: (hydrated: boolean) => set({ isHydrated: hydrated }),

      setPendingCartConflict: (conflict: any) => set({ pendingCartConflict: conflict }),
      clearPendingCartConflict: () => set({ pendingCartConflict: null }),
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
          // Use setTimeout to avoid blocking the hydration process
          if (state.token) {
            setTimeout(() => {
              state.checkAuth()
            }, 100)
          }
        }
      },
    }
  )
)


