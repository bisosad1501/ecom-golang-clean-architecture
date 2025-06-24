'use client'

import { useEffect } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { useAuthStore } from '@/store/auth'
import { canAccessRoute, PROTECTED_ROUTES } from '@/lib/permissions'
import { toast } from 'sonner'

interface UseAuthGuardOptions {
  redirectTo?: string
  requireAuth?: boolean
  requireAdmin?: boolean
  allowedRoles?: string[]
  onUnauthorized?: () => void
}

export function useAuthGuard(options: UseAuthGuardOptions = {}) {
  const router = useRouter()
  const pathname = usePathname()
  const { user, isAuthenticated, isLoading, isHydrated } = useAuthStore()
  
  const {
    redirectTo = '/auth/login',
    requireAuth = false,
    requireAdmin = false,
    allowedRoles = [],
    onUnauthorized,
  } = options

  useEffect(() => {
    // Don't check while loading or before hydration
    if (isLoading || !isHydrated) return

    // Check if route requires authentication
    if (requireAuth && !isAuthenticated) {
      toast.error('Please sign in to access this page')
      router.push(`${redirectTo}?redirect=${encodeURIComponent(pathname)}`)
      return
    }

    // Check if route requires admin access
    if (requireAdmin && (!user || !canAccessRoute(user.role, pathname))) {
      toast.error('You do not have permission to access this page')
      onUnauthorized?.()
      router.push('/')
      return
    }

    // Check if user role is allowed
    if (allowedRoles.length > 0 && user && !allowedRoles.includes(user.role)) {
      toast.error('You do not have permission to access this page')
      onUnauthorized?.()
      router.push('/')
      return
    }

    // Check general route access
    if (!canAccessRoute(user?.role || null, pathname)) {
      if (!isAuthenticated) {
        toast.error('Please sign in to access this page')
        router.push(`${redirectTo}?redirect=${encodeURIComponent(pathname)}`)
      } else {
        toast.error('You do not have permission to access this page')
        onUnauthorized?.()
        router.push('/')
      }
      return
    }

    // Redirect authenticated users away from guest-only pages
    if (isAuthenticated && PROTECTED_ROUTES.GUEST_ONLY.includes(pathname as any)) {
      const redirectUrl = new URLSearchParams(window.location.search).get('redirect')
      router.push(redirectUrl || '/')
      return
    }
  }, [
    isLoading,
    isHydrated,
    isAuthenticated,
    user,
    pathname,
    requireAuth,
    requireAdmin,
    allowedRoles,
    redirectTo,
    onUnauthorized,
    router,
  ])

  return {
    isLoading: isLoading || !isHydrated,
    isAuthenticated,
    user,
    canAccess: canAccessRoute(user?.role || null, pathname),
  }
}

// Specific hooks for common use cases
export function useRequireAuth(redirectTo?: string) {
  return useAuthGuard({ requireAuth: true, redirectTo })
}

export function useRequireAdmin(redirectTo?: string) {
  return useAuthGuard({ requireAdmin: true, redirectTo })
}

export function useRequireRole(roles: string[], redirectTo?: string) {
  return useAuthGuard({ allowedRoles: roles, redirectTo })
}

export function useGuestOnly(redirectTo?: string) {
  const { isAuthenticated } = useAuthStore()
  const router = useRouter()
  
  useEffect(() => {
    if (isAuthenticated) {
      const redirectUrl = new URLSearchParams(window.location.search).get('redirect')
      router.push(redirectUrl || redirectTo || '/')
    }
  }, [isAuthenticated, redirectTo, router])
  
  return { isAuthenticated }
}
