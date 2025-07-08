'use client'

import { AdminDashboard } from '@/components/admin/admin-dashboard'
import { useRequireAdmin } from '@/hooks/use-auth-guard'

export default function AdminPage() {
  // Require admin authentication
  const { isLoading } = useRequireAdmin()
  
  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-gray-600">Loading...</p>
        </div>
      </div>
    )
  }

  return <AdminDashboard />
}
