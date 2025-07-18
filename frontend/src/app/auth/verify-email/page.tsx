'use client'

import { useEffect, useState } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { CheckCircle, XCircle, Loader2, Mail } from 'lucide-react'
import { authApi } from '@/lib/api'
import { toast } from 'sonner'
import { ResendButton } from '@/components/auth/resend-button'

export default function VerifyEmailPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const token = searchParams.get('token')
  
  const [status, setStatus] = useState<'loading' | 'success' | 'error' | 'no-token'>('loading')
  const [message, setMessage] = useState('')
  const [userEmail, setUserEmail] = useState('')

  useEffect(() => {
    const success = searchParams.get('success')
    const error = searchParams.get('error')
    const email = searchParams.get('email')

    // Try to get email from URL params or localStorage
    const savedEmail = email || localStorage.getItem('pendingVerificationEmail') || ''
    if (savedEmail) {
      setUserEmail(savedEmail)
    }

    if (success === 'true') {
      setStatus('success')
      setMessage('Email verified successfully! You can now log in to your account.')
      setUserEmail(email || savedEmail)

      // Clear saved email
      localStorage.removeItem('pendingVerificationEmail')

      // Redirect to login after 3 seconds
      setTimeout(() => {
        router.push('/auth/login?verified=true')
      }, 3000)
      return
    }

    if (error) {
      setStatus('error')
      if (error === 'missing-token') {
        setMessage('No verification token provided. Please check your email for the correct link.')
      } else if (error === 'verification-failed') {
        setMessage('Failed to verify email. The token may be invalid or expired.')
      } else {
        setMessage('Verification failed. Please try again.')
      }
      return
    }

    if (!token) {
      setStatus('no-token')
      setMessage('No verification token provided')
      return
    }

    verifyEmail(token)
  }, [token, searchParams, router])

  const verifyEmail = async (verificationToken: string) => {
    try {
      setStatus('loading')
      
      const response = await authApi.verifyEmail(verificationToken)
      
      setStatus('success')
      setMessage('Email verified successfully! You can now log in to your account.')
      setUserEmail(response.user?.email || '')
      
      toast.success('Email verified successfully!')
      
      // Redirect to login after 3 seconds
      setTimeout(() => {
        router.push('/auth/login?verified=true')
      }, 3000)
      
    } catch (error: any) {
      setStatus('error')
      setMessage(error.message || 'Failed to verify email. The token may be invalid or expired.')
      toast.error('Email verification failed')
    }
  }

  const handleResendSuccess = () => {
    toast.success('Verification email sent! Please check your inbox.')
  }

  const handleGoToLogin = () => {
    router.push('/auth/login')
  }

  const renderContent = () => {
    switch (status) {
      case 'loading':
        return (
          <div className="text-center space-y-4">
            <Loader2 className="h-12 w-12 animate-spin mx-auto text-orange-500" />
            <h2 className="text-xl font-semibold text-white">Verifying your email...</h2>
            <p className="text-gray-300">Please wait while we verify your email address.</p>
          </div>
        )

      case 'success':
        return (
          <div className="text-center space-y-4">
            <CheckCircle className="h-12 w-12 mx-auto text-green-500" />
            <h2 className="text-xl font-semibold text-green-400">Email Verified Successfully!</h2>
            <p className="text-gray-300">{message}</p>
            {userEmail && (
              <p className="text-sm text-gray-400">
                Email: <span className="font-medium text-white">{userEmail}</span>
              </p>
            )}
            <div className="space-y-2">
              <p className="text-sm text-gray-400">Redirecting to login page in 3 seconds...</p>
              <Button onClick={handleGoToLogin} className="w-full bg-gradient-to-r from-orange-500 to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0">
                Go to Login Now
              </Button>
            </div>
          </div>
        )

      case 'error':
        return (
          <div className="text-center space-y-4">
            <XCircle className="h-12 w-12 mx-auto text-red-500" />
            <h2 className="text-xl font-semibold text-red-400">Verification Failed</h2>
            <p className="text-gray-300">{message}</p>
            <div className="space-y-2">
              <ResendButton
                email={userEmail}
                onSuccess={handleResendSuccess}
                variant="outline"
                className="w-full"
              />
              <Button onClick={handleGoToLogin} variant="ghost" className="w-full text-gray-400 hover:text-white hover:bg-gray-700">
                Back to Login
              </Button>
            </div>
          </div>
        )

      case 'no-token':
        return (
          <div className="text-center space-y-4">
            <XCircle className="h-12 w-12 mx-auto text-red-500" />
            <h2 className="text-xl font-semibold text-red-400">Invalid Verification Link</h2>
            <p className="text-gray-300">
              The verification link is invalid or missing. Please check your email for the correct link.
            </p>
            <div className="space-y-2">
              <ResendButton
                email={userEmail}
                onSuccess={handleResendSuccess}
                variant="outline"
                className="w-full"
              />
              <Button onClick={handleGoToLogin} variant="ghost" className="w-full text-gray-400 hover:text-white hover:bg-gray-700">
                Back to Login
              </Button>
            </div>
          </div>
        )

      default:
        return null
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <Card className="bg-gray-800/50 border-gray-700/50 backdrop-blur-sm shadow-2xl">
          <CardHeader className="text-center space-y-4">
            <div className="mx-auto w-16 h-16 bg-gradient-to-br from-orange-500 to-orange-600 rounded-full flex items-center justify-center shadow-lg">
              <Mail className="h-8 w-8 text-white" />
            </div>
            <CardTitle className="text-2xl font-bold text-white">Email Verification</CardTitle>
            <CardDescription className="text-gray-300">
              Verify your email address to complete your account setup
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {renderContent()}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
