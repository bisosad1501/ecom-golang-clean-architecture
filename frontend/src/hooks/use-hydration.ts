'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/store/auth'

/**
 * Hook to handle client-side hydration properly
 * Prevents hydration mismatch issues with auth state
 */
export function useHydration() {
  const [isHydrated, setIsHydrated] = useState(false)
  const { isHydrated: storeHydrated } = useAuthStore()

  useEffect(() => {
    setIsHydrated(true)
  }, [])

  return isHydrated && storeHydrated
}
