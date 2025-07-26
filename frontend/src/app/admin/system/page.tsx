import { Metadata } from 'next'
import AdminSystemPage from '@/components/admin/admin-system-page'

export const metadata: Metadata = {
  title: 'System | Admin Dashboard',
  description: 'System management, logs, backups, and maintenance.',
}

export default function Page() {
  return <AdminSystemPage />
}
