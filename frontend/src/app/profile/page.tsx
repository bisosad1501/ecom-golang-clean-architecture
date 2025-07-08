import { Metadata } from 'next'
import ProfilePage from '@/components/pages/profile-page'

export const metadata: Metadata = {
  title: 'My Profile | BiHub',
  description: 'Manage your BiHub account settings, view order history, and update your personal information.',
  robots: 'noindex, nofollow', // Profile pages should not be indexed
}

export default function Page() {
  return <ProfilePage />
}
