import { Metadata } from 'next'
import SettingsPage from '@/components/pages/settings-page'

export const metadata: Metadata = {
  title: 'Settings | BiHub',
  description: 'Manage your BiHub account settings, preferences, and privacy options.',
  keywords: 'settings, preferences, account, privacy, notifications, bihub',
  robots: 'noindex, nofollow',
}

export default function Settings() {
  return <SettingsPage />
}
