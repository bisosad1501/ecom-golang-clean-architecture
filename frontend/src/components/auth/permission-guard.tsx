'use client'

import { ReactNode } from 'react'
import { useAuthStore } from '@/store/auth'
import { hasPermission, hasAnyPermission, hasAllPermissions, canAccessAdminPanel } from '@/lib/permissions'
import { UserRole } from '@/types'

interface PermissionGuardProps {
  children: ReactNode
  permission?: string
  permissions?: string[]
  requireAll?: boolean
  role?: UserRole
  roles?: UserRole[]
  requireAuth?: boolean
  requireAdmin?: boolean
  fallback?: ReactNode
  inverse?: boolean
}

export function PermissionGuard({
  children,
  permission,
  permissions = [],
  requireAll = false,
  role,
  roles = [],
  requireAuth = false,
  requireAdmin = false,
  fallback = null,
  inverse = false,
}: PermissionGuardProps) {
  const { user, isAuthenticated } = useAuthStore()

  // Check authentication requirement
  if (requireAuth && !isAuthenticated) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // If not authenticated or user is null, show based on inverse
  if (!isAuthenticated || !user) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // Check admin requirement
  if (requireAdmin && !canAccessAdminPanel(user.role)) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // Check specific role
  if (role && user.role !== role) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // Check multiple roles
  if (roles.length > 0 && !roles.includes(user.role)) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // Check single permission
  if (permission && !hasPermission(user.role, permission)) {
    return inverse ? <>{children}</> : <>{fallback}</>
  }

  // Check multiple permissions
  if (permissions.length > 0) {
    const hasAccess = requireAll
      ? hasAllPermissions(user.role, permissions)
      : hasAnyPermission(user.role, permissions)
    
    if (!hasAccess) {
      return inverse ? <>{children}</> : <>{fallback}</>
    }
  }

  // If all checks pass, show children (or fallback if inverse)
  return inverse ? <>{fallback}</> : <>{children}</>
}

// Convenience components for common use cases
export function RequireAuth({ children, fallback = null }: { children: ReactNode; fallback?: ReactNode }) {
  return (
    <PermissionGuard requireAuth fallback={fallback}>
      {children}
    </PermissionGuard>
  )
}

export function RequireGuest({ children, fallback = null }: { children: ReactNode; fallback?: ReactNode }) {
  return (
    <PermissionGuard requireAuth inverse fallback={fallback}>
      {children}
    </PermissionGuard>
  )
}

export function RequireAdmin({ children, fallback = null }: { children: ReactNode; fallback?: ReactNode }) {
  return (
    <PermissionGuard requireAdmin fallback={fallback}>
      {children}
    </PermissionGuard>
  )
}

export function RequireRole({ 
  role, 
  roles, 
  children, 
  fallback = null 
}: { 
  role?: UserRole
  roles?: UserRole[]
  children: ReactNode
  fallback?: ReactNode 
}) {
  return (
    <PermissionGuard role={role} roles={roles} fallback={fallback}>
      {children}
    </PermissionGuard>
  )
}

export function RequirePermission({ 
  permission, 
  permissions, 
  requireAll = false,
  children, 
  fallback = null 
}: { 
  permission?: string
  permissions?: string[]
  requireAll?: boolean
  children: ReactNode
  fallback?: ReactNode 
}) {
  return (
    <PermissionGuard 
      permission={permission} 
      permissions={permissions} 
      requireAll={requireAll}
      fallback={fallback}
    >
      {children}
    </PermissionGuard>
  )
}
