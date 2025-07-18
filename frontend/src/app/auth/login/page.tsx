'use client'

import { useEffect, useState } from 'react'
import { useSearchParams } from 'next/navigation'
import { LoginForm } from '@/components/auth/login-form'
import { AuthLayout } from '@/components/auth/auth-layout'
import Link from 'next/link'
import { useGuestOnly } from '@/hooks/use-auth-guard'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Mail, CheckCircle } from 'lucide-react'
import { ResendButton } from '@/components/auth/resend-button'

export default function LoginPage() {
  useGuestOnly()

  const searchParams = useSearchParams()
  const [showAlert, setShowAlert] = useState(false)
  const [alertType, setAlertType] = useState<'verify-email' | 'verified'>('verify-email')
  const [userEmail, setUserEmail] = useState('')

  useEffect(() => {
    const message = searchParams.get('message')
    const verified = searchParams.get('verified')
    const email = searchParams.get('email')

    if (email) {
      setUserEmail(decodeURIComponent(email))
    }

    if (message === 'verify-email') {
      setAlertType('verify-email')
      setShowAlert(true)
    } else if (verified === 'true') {
      setAlertType('verified')
      setShowAlert(true)
      // Auto-hide verified message after 5 seconds
      setTimeout(() => setShowAlert(false), 5000)
    }
  }, [searchParams])

  const features = [
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
      ),
      title: "Bank-level Security",
      description: "256-bit encryption, multi-factor authentication, and secure data protection"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      ),
      title: "Smart Analytics",
      description: "AI-powered insights, personalized recommendations, and detailed reports"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      ),
      title: "Premium Experience",
      description: "Exclusive features, priority support, and enhanced user experience"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
      title: "24/7 Availability",
      description: "Round-the-clock access to your account and instant customer support"
    }
  ]

  const bottomContent = (
    <div className="mt-4 text-center">
      <p className="text-gray-300 text-xs">
        Don't have an account?{' '}
        <Link
          href="/auth/register"
          className="font-semibold text-[#FF9000] hover:text-[#FF9000]/80 transition-all duration-200 hover:underline"
        >
          Create one now
        </Link>
      </p>
    </div>
  )

  return (
    <AuthLayout
      title="Welcome back!"
      description="Please sign in to your account to continue your amazing shopping experience"
      heroTitle="Welcome back to your"
      heroSubtitle="digital experience"
      heroDescription="Seamless access to your personalized dashboard, exclusive features, and premium shopping experience"
      features={features}
      bottomContent={bottomContent}
    >
      {showAlert && (
        <div className="mb-4">
          {alertType === 'verify-email' ? (
            <div className="bg-blue-900/20 border border-blue-500/30 rounded-lg p-3 backdrop-blur-sm">
              <div className="flex items-start space-x-3">
                <Mail className="h-5 w-5 text-blue-400 mt-0.5 flex-shrink-0" />
                <div className="flex-1">
                  <p className="text-blue-200 text-sm font-medium">Email verification required</p>
                  <p className="text-blue-300/80 text-xs mt-1">
                    Check your inbox and click the verification link before logging in.
                  </p>
                  {userEmail && (
                    <div className="mt-2">
                      <ResendButton
                        email={userEmail}
                        variant="ghost"
                        size="sm"
                        className="text-xs text-blue-400 hover:text-blue-300 p-0 h-auto font-normal underline hover:no-underline"
                      />
                    </div>
                  )}
                  {!userEmail && (
                    <Link
                      href="/auth/resend-verification"
                      className="text-blue-400 hover:text-blue-300 text-xs underline hover:no-underline mt-1 inline-block"
                    >
                      Resend verification email
                    </Link>
                  )}
                </div>
                <button
                  onClick={() => setShowAlert(false)}
                  className="text-blue-400/60 hover:text-blue-300 text-lg leading-none"
                >
                  ×
                </button>
              </div>
            </div>
          ) : (
            <div className="bg-green-900/20 border border-green-500/30 rounded-lg p-3 backdrop-blur-sm">
              <div className="flex items-start space-x-3">
                <CheckCircle className="h-5 w-5 text-green-400 mt-0.5 flex-shrink-0" />
                <div className="flex-1">
                  <p className="text-green-200 text-sm font-medium">Email verified!</p>
                  <p className="text-green-300/80 text-xs mt-1">You can now log in to your account.</p>
                </div>
                <button
                  onClick={() => setShowAlert(false)}
                  className="text-green-400/60 hover:text-green-300 text-lg leading-none"
                >
                  ×
                </button>
              </div>
            </div>
          )}
        </div>
      )}
      <LoginForm />
    </AuthLayout>
  )
}
