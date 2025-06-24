import { Metadata } from 'next'
import AdminUsersPage from '@/components/admin/admin-users-page'

export const metadata: Metadata = {
  title: 'Users | Admin Dashboard',
  description: 'Manage user accounts, roles, and permissions.',
}

export default function Page() {
  return <AdminUsersPage />
}
