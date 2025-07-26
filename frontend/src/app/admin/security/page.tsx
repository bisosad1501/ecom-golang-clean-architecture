import { Metadata } from 'next'
import AdminSecurityPage from '@/components/admin/admin-security-page'

export const metadata: Metadata = {
  title: 'Security | Admin Dashboard',
  description: 'Security monitoring, suspicious activity, and access logs.',
}

export default function Page() {
  return <AdminSecurityPage />
}
