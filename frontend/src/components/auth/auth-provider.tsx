'use client'

import { useEffect, useRef } from 'react'
import { useAuthStore } from '@/store/auth'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { setHydrated, token, user, isAuthenticated } = useAuthStore()
  const hasHydrated = useRef(false)

  useEffect(() => {
    // Mark as hydrated after first render (client-side)
    if (!hasHydrated.current) {
      hasHydrated.current = true
      setHydrated(true)
    }
  }, [setHydrated])

  useEffect(() => {
    // Log auth state changes for debugging
    if (process.env.NODE_ENV === 'development') {
      console.log('Auth Provider: Auth state changed', { 
        hasToken: !!token, 
        hasUser: !!user,
        isAuthenticated,
        userRole: user?.role 
      })
    }
  }, [token, user, isAuthenticated])

  return <>{children}</>
}
