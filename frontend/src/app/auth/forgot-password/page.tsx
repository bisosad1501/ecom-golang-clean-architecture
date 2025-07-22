import { Metadata } from 'next'
import { ForgotPasswordForm } from '@/components/auth/forgot-password-form'
import { AuthLayout } from '@/components/auth/auth-layout'
import { Mail, Shield, Clock, ArrowLeft } from 'lucide-react'

export const metadata: Metadata = {
  title: 'Forgot Password | BiHub',
  description: 'Reset your password to regain access to your BiHub account',
}

export default function ForgotPasswordPage() {
  return (
    <AuthLayout
      title="Reset Your Password"
      description="Enter your email address and we'll send you a link to reset your password"
      heroTitle="Forgot Your Password?"
      heroSubtitle="No Worries!"
      heroDescription="We'll help you get back into your account quickly and securely. Just enter your email address below."
      features={[
        {
          icon: <Mail className="h-5 w-5 text-white" />,
          title: "Email Recovery",
          description: "Secure password reset via email verification"
        },
        {
          icon: <Shield className="h-5 w-5 text-white" />,
          title: "Secure Process",
          description: "Bank-level security for your account protection"
        },
        {
          icon: <Clock className="h-5 w-5 text-white" />,
          title: "Quick Reset",
          description: "Get back to shopping in just a few minutes"
        },
        {
          icon: <ArrowLeft className="h-5 w-5 text-white" />,
          title: "Easy Return",
          description: "Remember your password? Sign in anytime"
        }
      ]}
    >
      <ForgotPasswordForm />
    </AuthLayout>
  )
}
