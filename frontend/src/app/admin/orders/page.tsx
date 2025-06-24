import { Metadata } from 'next'
import AdminOrdersPage from '@/components/admin/admin-orders-page'

export const metadata: Metadata = {
  title: 'Orders | Admin Dashboard',
  description: 'Manage customer orders, track fulfillment, and process payments.',
}

export default function Page() {
  return <AdminOrdersPage />
}
