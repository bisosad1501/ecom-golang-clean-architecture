'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { ArrowLeft, Mail, CheckCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { authService } from '@/lib/services/auth'
import { toast } from 'sonner'

const forgotPasswordSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
})

type ForgotPasswordFormData = z.infer<typeof forgotPasswordSchema>

export function ForgotPasswordForm() {
  const [isLoading, setIsLoading] = useState(false)
  const [isSuccess, setIsSuccess] = useState(false)
  const [submittedEmail, setSubmittedEmail] = useState('')

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<ForgotPasswordFormData>({
    resolver: zodResolver(forgotPasswordSchema),
    defaultValues: {
      email: '',
    },
  })

  const onSubmit = async (data: ForgotPasswordFormData) => {
    try {
      setIsLoading(true)
      console.log('🔄 [FORGOT PASSWORD] Sending reset email to:', data.email)
      
      await authService.forgotPassword(data.email)
      
      console.log('✅ [FORGOT PASSWORD] Reset email sent successfully')
      setSubmittedEmail(data.email)
      setIsSuccess(true)
      toast.success('Password reset email sent! Check your inbox.')
      
    } catch (error: any) {
      console.error('❌ [FORGOT PASSWORD] Error:', error)
      
      if (error.code === 'VALIDATION_ERROR' && error.details) {
        // Handle field-specific validation errors
        Object.entries(error.details).forEach(([field, message]) => {
          setError(field as keyof ForgotPasswordFormData, {
            type: 'server',
            message: message as string,
          })
        })
      } else {
        toast.error(error.message || 'Failed to send reset email. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  const handleResendEmail = async () => {
    if (!submittedEmail) return
    
    try {
      setIsLoading(true)
      await authService.forgotPassword(submittedEmail)
      toast.success('Reset email sent again! Check your inbox.')
    } catch (error: any) {
      toast.error(error.message || 'Failed to resend email. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  if (isSuccess) {
    return (
      <div className="space-y-6 text-center">
        {/* Success Icon */}
        <div className="flex justify-center">
          <div className="w-16 h-16 bg-green-500/20 rounded-full flex items-center justify-center">
            <CheckCircle className="h-8 w-8 text-green-500" />
          </div>
        </div>

        {/* Success Message */}
        <div className="space-y-2">
          <h3 className="text-lg font-semibold text-white">Check Your Email</h3>
          <p className="text-sm text-gray-300 leading-relaxed">
            We've sent a password reset link to{' '}
            <span className="font-medium text-[#FF9000]">{submittedEmail}</span>
          </p>
          <p className="text-xs text-gray-400">
            The link will expire in 1 hour for security reasons.
          </p>
        </div>

        {/* Instructions */}
        <div className="bg-gray-800/50 rounded-lg p-4 text-left">
          <h4 className="text-sm font-medium text-white mb-2">Next Steps:</h4>
          <ol className="text-xs text-gray-300 space-y-1 list-decimal list-inside">
            <li>Check your email inbox (and spam folder)</li>
            <li>Click the "Reset Password" link in the email</li>
            <li>Enter your new password</li>
            <li>Sign in with your new password</li>
          </ol>
        </div>

        {/* Action Buttons */}
        <div className="space-y-3">
          <Button
            onClick={handleResendEmail}
            variant="outline"
            className="w-full bg-gray-800/50 border-gray-600 text-gray-300 hover:bg-gray-700/50 hover:text-white transition-all duration-200"
            isLoading={isLoading}
            loadingText="Sending..."
          >
            <Mail className="h-4 w-4 mr-2" />
            Resend Email
          </Button>

          <Link
            href="/auth/login"
            className="inline-flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-[#FF9000] hover:text-[#e67e00] transition-all duration-200 hover:underline"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Back to Sign In
          </Link>
        </div>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* Email Input */}
      <div className="space-y-1">
        <Input
          {...register('email')}
          type="email"
          label="Email address"
          placeholder="Enter your email address"
          error={errors.email?.message}
          required
          autoComplete="email"
          autoFocus
          size="lg"
          className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
          leftIcon={<Mail className="h-4 w-4 text-gray-400" />}
        />
      </div>

      {/* Help Text */}
      <div className="bg-blue-500/10 border border-blue-500/20 rounded-lg p-3">
        <p className="text-xs text-blue-300 leading-relaxed">
          <strong>Don't worry!</strong> Enter the email address associated with your account and we'll send you a secure link to reset your password.
        </p>
      </div>

      {/* Submit Button */}
      <Button
        type="submit"
        className="w-full mt-4 bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold py-3 text-sm rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg border-0"
        size="lg"
        isLoading={isLoading}
        loadingText="Sending reset email..."
      >
        <Mail className="h-4 w-4 mr-2" />
        Send Reset Email
      </Button>

      {/* Back to Login */}
      <div className="text-center pt-2">
        <Link
          href="/auth/login"
          className="inline-flex items-center text-sm font-medium text-gray-400 hover:text-[#FF9000] transition-all duration-200 hover:underline"
        >
          <ArrowLeft className="h-4 w-4 mr-1" />
          Remember your password? Sign in
        </Link>
      </div>
    </form>
  )
}
