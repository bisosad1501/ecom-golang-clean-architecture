import { Metadata } from 'next'
import { UserPreferencesForm } from '@/components/profile/user-preferences-form'

export const metadata: Metadata = {
  title: 'Preferences | BiHub',
  description: 'Manage your account preferences and settings',
}

export default function UserPreferencesPage() {
  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Account Preferences</h1>
          <p className="mt-2 text-gray-600">
            Customize your experience and manage your notification settings
          </p>
        </div>

        {/* Preferences Form */}
        <UserPreferencesForm />
      </div>
    </div>
  )
}
