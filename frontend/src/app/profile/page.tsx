import { Metadata } from 'next'
import ProfilePage from '@/components/pages/profile-page'

export const metadata: Metadata = {
  title: 'My Profile | EcomStore',
  description: 'Manage your account settings, view order history, and update your personal information.',
  robots: 'noindex, nofollow', // Profile pages should not be indexed
}

export default function Page() {
  return <ProfilePage />
}
