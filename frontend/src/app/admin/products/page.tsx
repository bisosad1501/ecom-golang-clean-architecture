'use client'

import { AdminProductsPage } from '@/components/admin/admin-products-page'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'

export default function AdminProducts() {
  return (
    <RequirePermission 
      permissions={[PERMISSIONS.PRODUCTS_VIEW, PERMISSIONS.PRODUCTS_CREATE]}
      fallback={
        <div className="text-center py-12">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Access Denied</h2>
          <p className="text-gray-600">You don't have permission to manage products.</p>
        </div>
      }
    >
      <AdminProductsPage />
    </RequirePermission>
  )
}
