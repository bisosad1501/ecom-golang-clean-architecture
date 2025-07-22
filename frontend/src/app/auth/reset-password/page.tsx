import { Metadata } from 'next'
import { ResetPasswordForm } from '@/components/auth/reset-password-form'
import { AuthLayout } from '@/components/auth/auth-layout'
import { Lock, Shield, CheckCircle, Key } from 'lucide-react'

export const metadata: Metadata = {
  title: 'Reset Password | BiHub',
  description: 'Create a new password for your BiHub account',
}

export default function ResetPasswordPage() {
  return (
    <AuthLayout
      title="Create New Password"
      description="Enter a strong new password to secure your account"
      heroTitle="Almost There!"
      heroSubtitle="New Password"
      heroDescription="You're just one step away from regaining access to your account. Create a strong, secure password below."
      features={[
        {
          icon: <Lock className="h-5 w-5 text-white" />,
          title: "Secure Password",
          description: "Create a strong password with 8+ characters"
        },
        {
          icon: <Shield className="h-5 w-5 text-white" />,
          title: "Account Protection",
          description: "Your account will be secured immediately"
        },
        {
          icon: <CheckCircle className="h-5 w-5 text-white" />,
          title: "Instant Access",
          description: "Sign in immediately after password reset"
        },
        {
          icon: <Key className="h-5 w-5 text-white" />,
          title: "One-Time Link",
          description: "This reset link expires after use for security"
        }
      ]}
    >
      <ResetPasswordForm />
    </AuthLayout>
  )
}
