'use client'

import { useEffect, useRef } from 'react'
import { useAuthStore } from '@/store/auth'

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { setHydrated, token, user } = useAuthStore()
  const hasHydrated = useRef(false)

  useEffect(() => {
    // Mark as hydrated after first render (client-side)
    if (!hasHydrated.current) {
      hasHydrated.current = true
      setHydrated(true)
      
      if (process.env.NODE_ENV === 'development') {
        console.log('Auth Provider: Hydration complete', { 
          hasToken: !!token, 
          hasUser: !!user,
          userRole: user?.role 
        })
      }
    }
  }, [setHydrated, token, user])

  return <>{children}</>
}
