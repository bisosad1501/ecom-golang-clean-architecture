import { Metadata } from 'next'
import AdminCustomersPage from '@/components/admin/admin-customers-page'

export const metadata: Metadata = {
  title: 'Customers | Admin Dashboard',
  description: 'Manage customer accounts, segments, and analytics.',
}

export default function Page() {
  return <AdminCustomersPage />
}
