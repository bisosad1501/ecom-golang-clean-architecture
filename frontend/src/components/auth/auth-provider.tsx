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

  // Remove excessive logging in development

  return <>{children}</>
}
