import { Metadata } from 'next'
import AdminSettingsPage from '@/components/admin/admin-settings-page'

export const metadata: Metadata = {
  title: 'Settings | Admin Dashboard',
  description: 'Manage system settings, configurations, and preferences.',
}

export default function Page() {
  return <AdminSettingsPage />
}
