'use client'

import { LoginForm } from '@/components/auth/login-form'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'
import { APP_NAME } from '@/constants'
import { useGuestOnly } from '@/hooks/use-auth-guard'

export default function LoginPage() {
  useGuestOnly()

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-50 via-background to-violet-50 relative overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <svg className="w-full h-full" viewBox="0 0 100 100" fill="none">
          <pattern id="authGrid" width="10" height="10" patternUnits="userSpaceOnUse">
            <path d="M 10 0 L 0 0 0 10" fill="none" stroke="currentColor" strokeWidth="0.5"/>
          </pattern>
          <rect width="100%" height="100%" fill="url(#authGrid)" />
        </svg>
      </div>

      {/* Floating Elements */}
      <div className="absolute top-20 left-20 w-32 h-32 bg-gradient-to-br from-primary-400/20 to-violet-500/20 rounded-full blur-3xl animate-float"></div>
      <div className="absolute bottom-20 right-20 w-24 h-24 bg-gradient-to-br from-violet-400/20 to-primary-500/20 rounded-full blur-2xl animate-pulse-slow"></div>

      <div className="relative z-10 min-h-screen flex">
        {/* Left Side - Branding */}
        <div className="hidden lg:flex lg:w-1/2 flex-col justify-center px-12 xl:px-20">
          <div className="max-w-md">
            <Link href="/" className="flex items-center space-x-3 mb-12 group">
              <div className="h-16 w-16 rounded-3xl bg-gradient-to-br from-primary-500 via-primary-600 to-violet-600 flex items-center justify-center shadow-2xl group-hover:shadow-3xl group-hover:scale-105 transition-all duration-300">
                <span className="text-white font-bold text-3xl">E</span>
              </div>
              <span className="text-4xl font-bold text-gradient bg-gradient-to-r from-primary-600 via-primary-500 to-violet-500 bg-clip-text text-transparent">{APP_NAME}</span>
            </Link>

            <h1 className="text-5xl font-bold text-foreground mb-6 leading-tight">
              Welcome back to your
              <span className="text-gradient"> shopping journey</span>
            </h1>

            <p className="text-xl text-muted-foreground leading-relaxed mb-8">
              Sign in to access your account, track orders, and discover amazing products tailored just for you.
            </p>

            {/* Trust Indicators */}
            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-emerald-100 flex items-center justify-center">
                  <span className="text-emerald-600 text-sm">✓</span>
                </div>
                <span className="text-muted-foreground">Secure & encrypted login</span>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-blue-100 flex items-center justify-center">
                  <span className="text-blue-600 text-sm">✓</span>
                </div>
                <span className="text-muted-foreground">Access to exclusive deals</span>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-8 h-8 rounded-lg bg-purple-100 flex items-center justify-center">
                  <span className="text-purple-600 text-sm">✓</span>
                </div>
                <span className="text-muted-foreground">Personalized recommendations</span>
              </div>
            </div>
          </div>
        </div>

        {/* Right Side - Login Form */}
        <div className="flex-1 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-12">
          <div className="max-w-md w-full space-y-8">
            {/* Mobile Logo */}
            <div className="lg:hidden text-center">
              <Link href="/" className="flex items-center justify-center space-x-3 mb-8 group">
                <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large group-hover:shadow-xl transition-all duration-300">
                  <span className="text-white font-bold text-2xl">E</span>
                </div>
                <span className="text-3xl font-bold text-gradient">{APP_NAME}</span>
              </Link>
            </div>

            <Card variant="elevated" className="border-0 shadow-2xl">
              <CardHeader className="text-center pb-8">
                <CardTitle className="text-3xl font-bold text-foreground mb-3">Welcome back!</CardTitle>
                <CardDescription className="text-lg text-muted-foreground">
                  Please sign in to your account to continue
                </CardDescription>
              </CardHeader>
              <CardContent className="px-8 pb-8">
                <LoginForm />

                <div className="mt-8 text-center">
                  <p className="text-muted-foreground">
                    Don't have an account?{' '}
                    <Link
                      href="/auth/register"
                      className="font-semibold text-primary hover:text-primary-600 transition-colors"
                    >
                      Create one now
                    </Link>
                  </p>
                </div>
              </CardContent>
            </Card>

            <div className="text-center">
              <Link
                href="/auth/forgot-password"
                className="text-muted-foreground hover:text-primary transition-colors font-medium"
              >
                Forgot your password?
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
