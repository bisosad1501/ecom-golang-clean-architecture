'use client'

import { useState, useEffect } from 'react'
import { useRouter, useSearchParams } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff, Lock, CheckCircle, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { authService } from '@/lib/services/auth'
import { toast } from 'sonner'
import type { ResetPasswordFormRequest } from '@/types'

const resetPasswordSchema = z.object({
  new_password: z
    .string()
    .min(6, 'Password must be at least 6 characters')
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, 'Password must contain at least one uppercase letter, one lowercase letter, and one number'),
  password_confirmation: z.string(),
}).refine((data) => data.new_password === data.password_confirmation, {
  message: "Passwords don't match",
  path: ["password_confirmation"],
})

type ResetPasswordFormData = z.infer<typeof resetPasswordSchema>

export function ResetPasswordForm() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const [isLoading, setIsLoading] = useState(false)
  const [isSuccess, setIsSuccess] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [token, setToken] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
    watch,
  } = useForm<ResetPasswordFormData>({
    resolver: zodResolver(resetPasswordSchema),
    defaultValues: {
      new_password: '',
      password_confirmation: '',
    },
  })

  const password = watch('new_password')

  useEffect(() => {
    const tokenParam = searchParams.get('token')
    if (!tokenParam) {
      toast.error('Invalid reset link. Please request a new password reset.')
      router.push('/auth/forgot-password')
      return
    }
    setToken(tokenParam)
  }, [searchParams, router])

  const onSubmit = async (data: ResetPasswordFormData) => {
    if (!token) {
      toast.error('Invalid reset token. Please request a new password reset.')
      return
    }

    try {
      setIsLoading(true)
      console.log('🔄 [RESET PASSWORD] Resetting password with token:', token)
      
      await authService.resetPassword({
        token,
        new_password: data.new_password,
      })
      
      console.log('✅ [RESET PASSWORD] Password reset successful')
      setIsSuccess(true)
      toast.success('Password reset successful! You can now sign in with your new password.')
      
    } catch (error: any) {
      console.error('❌ [RESET PASSWORD] Error:', error)
      
      if (error.message?.includes('expired')) {
        toast.error('Reset link has expired. Please request a new password reset.', {
          action: {
            label: 'Request New Reset',
            onClick: () => router.push('/auth/forgot-password'),
          },
        })
      } else if (error.message?.includes('already used')) {
        toast.error('Reset link has already been used. Please request a new password reset.', {
          action: {
            label: 'Request New Reset',
            onClick: () => router.push('/auth/forgot-password'),
          },
        })
      } else if (error.code === 'VALIDATION_ERROR' && error.details) {
        // Handle field-specific validation errors
        Object.entries(error.details).forEach(([field, message]) => {
          setError(field as keyof ResetPasswordFormData, {
            type: 'server',
            message: message as string,
          })
        })
      } else {
        toast.error(error.message || 'Failed to reset password. Please try again.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  // Password strength indicator
  const getPasswordStrength = (password: string) => {
    if (!password) return { strength: 0, label: '', color: '' }
    
    let strength = 0
    if (password.length >= 6) strength++
    if (password.length >= 8) strength++
    if (/[a-z]/.test(password)) strength++
    if (/[A-Z]/.test(password)) strength++
    if (/\d/.test(password)) strength++
    if (/[^a-zA-Z\d]/.test(password)) strength++

    if (strength <= 2) return { strength, label: 'Weak', color: 'bg-red-500' }
    if (strength <= 4) return { strength, label: 'Medium', color: 'bg-yellow-500' }
    return { strength, label: 'Strong', color: 'bg-green-500' }
  }

  const passwordStrength = getPasswordStrength(password)

  if (!token) {
    return (
      <div className="space-y-4 text-center">
        <div className="flex justify-center">
          <div className="w-16 h-16 bg-red-500/20 rounded-full flex items-center justify-center">
            <AlertCircle className="h-8 w-8 text-red-500" />
          </div>
        </div>
        <div className="space-y-2">
          <h3 className="text-lg font-semibold text-white">Invalid Reset Link</h3>
          <p className="text-sm text-gray-300">
            This password reset link is invalid or has expired.
          </p>
        </div>
        <Link
          href="/auth/forgot-password"
          className="inline-flex items-center justify-center w-full px-4 py-2 bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold rounded-lg transition-all duration-300"
        >
          Request New Reset Link
        </Link>
      </div>
    )
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
          <h3 className="text-lg font-semibold text-white">Password Reset Successful!</h3>
          <p className="text-sm text-gray-300 leading-relaxed">
            Your password has been successfully updated. You can now sign in with your new password.
          </p>
        </div>

        {/* Sign In Button */}
        <Button
          onClick={() => router.push('/auth/login')}
          className="w-full bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold py-3 text-sm rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg border-0"
          size="lg"
        >
          Sign In Now
        </Button>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {/* New Password */}
      <div className="space-y-1">
        <Input
          {...register('new_password')}
          type={showPassword ? 'text' : 'password'}
          label="New Password"
          placeholder="Enter your new password"
          error={errors.new_password?.message}
          required
          autoComplete="new-password"
          size="lg"
          className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
          leftIcon={<Lock className="h-4 w-4 text-gray-400" />}
          rightIcon={
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="text-gray-400 hover:text-[#FF9000] transition-all duration-200 p-1 rounded-lg hover:bg-gray-700/50"
            >
              {showPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          }
        />
        
        {/* Password Strength Indicator */}
        {password && (
          <div className="space-y-1">
            <div className="flex items-center justify-between">
              <span className="text-xs text-gray-400">Password strength:</span>
              <span className={`text-xs font-medium ${
                passwordStrength.label === 'Strong' ? 'text-green-400' :
                passwordStrength.label === 'Medium' ? 'text-yellow-400' : 'text-red-400'
              }`}>
                {passwordStrength.label}
              </span>
            </div>
            <div className="w-full bg-gray-700 rounded-full h-1">
              <div 
                className={`h-1 rounded-full transition-all duration-300 ${passwordStrength.color}`}
                style={{ width: `${(passwordStrength.strength / 6) * 100}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {/* Confirm Password */}
      <div className="space-y-1">
        <Input
          {...register('password_confirmation')}
          type={showConfirmPassword ? 'text' : 'password'}
          label="Confirm New Password"
          placeholder="Confirm your new password"
          error={errors.password_confirmation?.message}
          required
          autoComplete="new-password"
          size="lg"
          className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
          leftIcon={<Lock className="h-4 w-4 text-gray-400" />}
          rightIcon={
            <button
              type="button"
              onClick={() => setShowConfirmPassword(!showConfirmPassword)}
              className="text-gray-400 hover:text-[#FF9000] transition-all duration-200 p-1 rounded-lg hover:bg-gray-700/50"
            >
              {showConfirmPassword ? (
                <EyeOff className="h-4 w-4" />
              ) : (
                <Eye className="h-4 w-4" />
              )}
            </button>
          }
        />
      </div>

      {/* Password Requirements */}
      <div className="bg-gray-800/50 rounded-lg p-3">
        <h4 className="text-xs font-medium text-white mb-2">Password Requirements:</h4>
        <ul className="text-xs text-gray-400 space-y-1">
          <li className={password?.length >= 6 ? 'text-green-400' : ''}>
            • At least 6 characters long
          </li>
          <li className={/[a-z]/.test(password || '') ? 'text-green-400' : ''}>
            • Contains lowercase letter
          </li>
          <li className={/[A-Z]/.test(password || '') ? 'text-green-400' : ''}>
            • Contains uppercase letter
          </li>
          <li className={/\d/.test(password || '') ? 'text-green-400' : ''}>
            • Contains number
          </li>
        </ul>
      </div>

      {/* Submit Button */}
      <Button
        type="submit"
        className="w-full mt-4 bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold py-3 text-sm rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg border-0"
        size="lg"
        isLoading={isLoading}
        loadingText="Resetting password..."
      >
        <Lock className="h-4 w-4 mr-2" />
        Reset Password
      </Button>

      {/* Back to Login */}
      <div className="text-center pt-2">
        <Link
          href="/auth/login"
          className="text-sm font-medium text-gray-400 hover:text-[#FF9000] transition-all duration-200 hover:underline"
        >
          Back to Sign In
        </Link>
      </div>
    </form>
  )
}
