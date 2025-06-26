'use client'

import { RegisterForm } from '@/components/auth/register-form'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import Link from 'next/link'
import { APP_NAME } from '@/constants'
import { useGuestOnly } from '@/hooks/use-auth-guard'

export default function RegisterPage() {
  useGuestOnly()

  return (
    <div className="min-h-screen bg-gradient-to-br from-violet-50 via-background to-primary-50 relative overflow-hidden pt-[140px]">{/* Add padding top for header */}
      {/* Background Pattern */}
      <div className="absolute inset-0 opacity-5">
        <svg className="w-full h-full" viewBox="0 0 100 100" fill="none">
          <pattern id="registerGrid" width="10" height="10" patternUnits="userSpaceOnUse">
            <path d="M 10 0 L 0 0 0 10" fill="none" stroke="currentColor" strokeWidth="0.5"/>
          </pattern>
          <rect width="100%" height="100%" fill="url(#registerGrid)" />
        </svg>
      </div>

      {/* Floating Elements */}
      <div className="absolute top-20 right-20 w-32 h-32 bg-gradient-to-br from-violet-400/20 to-primary-500/20 rounded-full blur-3xl animate-float"></div>
      <div className="absolute bottom-20 left-20 w-24 h-24 bg-gradient-to-br from-primary-400/20 to-violet-500/20 rounded-full blur-2xl animate-pulse-slow"></div>

      <div className="relative z-10 min-h-[calc(100vh-140px)] flex">{/* Adjust inner container height */}
        {/* Left Side - Registration Form */}
        <div className="flex-1 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-6">{/* Reduced padding */}
          <div className="max-w-lg w-full space-y-6">{/* Reduced spacing */}
            {/* Mobile Logo */}
            <div className="lg:hidden text-center">
              <Link href="/" className="flex items-center justify-center space-x-3 mb-6 group">{/* Reduced margin */}
                <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large group-hover:shadow-xl transition-all duration-300">
                  <span className="text-white font-bold text-2xl">E</span>
                </div>
                <span className="text-3xl font-bold text-gradient">{APP_NAME}</span>
              </Link>
            </div>

            <Card variant="elevated" className="border-0 shadow-2xl">
              <CardHeader className="text-center pb-6">{/* Reduced padding */}
                <CardTitle className="text-2xl font-bold text-foreground mb-2">{/* Smaller title and margin */}
                  Join our community
                </CardTitle>
                <CardDescription className="text-base text-muted-foreground">{/* Smaller text */}
                  Create your account and start your shopping journey
                </CardDescription>
              </CardHeader>
              <CardContent className="px-8 pb-6">{/* Reduced padding */}
                <RegisterForm />

                <div className="mt-6 text-center">{/* Reduced margin */}
                  <p className="text-muted-foreground">
                    Already have an account?{' '}
                    <Link
                      href="/auth/login"
                      className="font-semibold text-primary hover:text-primary-600 transition-colors"
                    >
                      Sign in instead
                    </Link>
                  </p>
                </div>
              </CardContent>
            </Card>

            <div className="text-center text-sm text-muted-foreground">
              By creating an account, you agree to our{' '}
              <Link href="/terms" className="text-primary hover:text-primary-600 transition-colors font-medium">
                Terms of Service
              </Link>{' '}
              and{' '}
              <Link href="/privacy" className="text-primary hover:text-primary-600 transition-colors font-medium">
                Privacy Policy
              </Link>
            </div>
          </div>
        </div>

        {/* Right Side - Benefits */}
        <div className="hidden lg:flex lg:w-1/2 flex-col justify-center px-12 xl:px-20">
          <div className="max-w-md">
            <Link href="/" className="flex items-center space-x-3 mb-12 group">
              <div className="h-16 w-16 rounded-3xl bg-gradient-to-br from-violet-500 via-violet-600 to-primary-600 flex items-center justify-center shadow-2xl group-hover:shadow-3xl group-hover:scale-105 transition-all duration-300">
                <span className="text-white font-bold text-3xl">E</span>
              </div>
              <span className="text-4xl font-bold text-gradient bg-gradient-to-r from-violet-600 via-violet-500 to-primary-500 bg-clip-text text-transparent">{APP_NAME}</span>
            </Link>

            <h1 className="text-5xl font-bold text-foreground mb-6 leading-tight">
              Start your amazing
              <span className="text-gradient"> shopping experience</span>
            </h1>

            <p className="text-xl text-muted-foreground leading-relaxed mb-8">
              Join thousands of satisfied customers and unlock exclusive benefits, personalized recommendations, and seamless shopping.
            </p>

            {/* Benefits */}
            <div className="space-y-6">
              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center shadow-large">
                  <span className="text-white text-lg">🎁</span>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground mb-1">Exclusive Member Deals</h3>
                  <p className="text-muted-foreground">Access special discounts and early access to sales</p>
                </div>
              </div>

              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center shadow-large">
                  <span className="text-white text-lg">🚀</span>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground mb-1">Fast & Free Shipping</h3>
                  <p className="text-muted-foreground">Enjoy free shipping on orders over $50</p>
                </div>
              </div>

              <div className="flex items-start gap-4">
                <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-purple-500 to-purple-600 flex items-center justify-center shadow-large">
                  <span className="text-white text-lg">💝</span>
                </div>
                <div>
                  <h3 className="font-semibold text-foreground mb-1">Personalized Experience</h3>
                  <p className="text-muted-foreground">Get product recommendations tailored to your preferences</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}