'use client'

import { AdminLayout } from '@/components/admin/admin-layout'
import { useRequireAdmin } from '@/hooks/use-auth-guard'
import { useHydration } from '@/hooks/use-hydration'

export default function AdminLayoutWrapper({
  children,
}: {
  children: React.ReactNode
}) {
  const { isLoading, canAccess } = useRequireAdmin()
  const isHydrated = useHydration()

  // Show loading while hydrating or checking auth
  if (!isHydrated || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading admin panel...</p>
        </div>
      </div>
    )
  }

  if (!canAccess) {
    return null // useRequireAdmin will handle redirect
  }

  return <AdminLayout>{children}</AdminLayout>
}
