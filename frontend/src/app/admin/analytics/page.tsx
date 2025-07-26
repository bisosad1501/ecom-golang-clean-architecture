import { Metadata } from 'next'
import AdminAnalyticsPage from '@/components/admin/admin-analytics-page'

export const metadata: Metadata = {
  title: 'Analytics | Admin Dashboard',
  description: 'View detailed analytics and insights about your store performance.',
}

export default function Page() {
  return <AdminAnalyticsPage />
}
