'use client'

import { LoginForm } from '@/components/auth/login-form'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'
import { APP_NAME } from '@/constants'
import { useGuestOnly } from '@/hooks/use-auth-guard'

export default function LoginPage() {
  useGuestOnly()

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <Link href="/" className="flex items-center justify-center space-x-2 mb-8">
            <div className="h-10 w-10 rounded-lg bg-primary-600 flex items-center justify-center">
              <span className="text-white font-bold text-xl">E</span>
            </div>
            <span className="text-2xl font-bold text-gray-900">{APP_NAME}</span>
          </Link>
        </div>

        <Card>
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold">Sign in to your account</CardTitle>
            <CardDescription>
              Welcome back! Please enter your details to continue.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <LoginForm />
            
            <div className="mt-6 text-center">
              <p className="text-sm text-gray-600">
                Don't have an account?{' '}
                <Link 
                  href="/auth/register" 
                  className="font-medium text-primary-600 hover:text-primary-500"
                >
                  Sign up here
                </Link>
              </p>
            </div>
          </CardContent>
        </Card>

        <div className="text-center">
          <Link 
            href="/auth/forgot-password" 
            className="text-sm text-primary-600 hover:text-primary-500"
          >
            Forgot your password?
          </Link>
        </div>
      </div>
    </div>
  )
}
