'use client'

import { RegisterForm } from '@/components/auth/register-form'
import { AuthLayout } from '@/components/auth/auth-layout'
import Link from 'next/link'
import { useGuestOnly } from '@/hooks/use-auth-guard'

export default function RegisterPage() {
  useGuestOnly()

  const features = [
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
        </svg>
      ),
      title: "Exclusive Member Deals",
      description: "Access special discounts and early sales with member-only pricing"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
        </svg>
      ),
      title: "Fast & Free Shipping",
      description: "Enjoy free shipping on orders over $50 with express delivery options"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      ),
      title: "Personalized Experience",
      description: "AI-powered product recommendations tailored to your preferences"
    },
    {
      icon: (
        <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
      ),
      title: "Secure Shopping",
      description: "Bank-level security with encrypted transactions and data protection"
    }
  ]

  const bottomContent = (
    <>
      <div className="mt-4 text-center">
        <p className="text-gray-300 text-xs">
          Already have an account?{' '}
          <Link
            href="/auth/login"
            className="font-semibold text-[#FF9000] hover:text-[#FF9000]/80 transition-all duration-200 hover:underline"
          >
            Sign in instead
          </Link>
        </p>
      </div>

      <div className="mt-2 text-center text-xs text-gray-400 leading-relaxed">
        By creating an account, you agree to our{' '}
        <Link
          href="/terms"
          className="text-[#FF9000] hover:text-[#FF9000]/80 transition-colors font-medium hover:underline"
        >
          Terms
        </Link>{' '}
        and{' '}
        <Link
          href="/privacy"
          className="text-[#FF9000] hover:text-[#FF9000]/80 transition-colors font-medium hover:underline"
        >
          Privacy Policy
        </Link>
      </div>
    </>
  )

  return (
    <AuthLayout
      title="Join our community"
      description="Create your account and start your amazing shopping journey with exclusive benefits"
      heroTitle="Start your amazing"
      heroSubtitle="shopping experience"
      heroDescription="Join thousands of satisfied customers and unlock exclusive benefits, personalized recommendations, and premium features"
      features={features}
      bottomContent={bottomContent}
    >
      <RegisterForm />
    </AuthLayout>
  )
}