'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '@/store/auth'
import { useCartStore, isGuestCart } from '@/store/cart'
import { LoginRequest } from '@/types'
import { toast } from 'sonner'
import { OAuthButtons } from './OAuthButtons'
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'
import { canAccessAdminPanel } from '@/lib/permissions'

const loginSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  remember_me: z.boolean().optional(),
})

type LoginFormData = z.infer<typeof loginSchema>

export function LoginForm() {
  const router = useRouter()
  const [showPassword, setShowPassword] = useState(false)
  const { login, isLoading, error, clearError, pendingCartConflict } = useAuthStore()
  const { cart } = useCartStore()

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
      remember_me: false,
    },
  })

  const onSubmit = async (data: LoginFormData) => {
    try {
      clearError()
      console.log('🚀 [NORMAL LOGIN] Starting login process...')
      console.log('🚀 [NORMAL LOGIN] Guest session ID:', localStorage.getItem('guest_session_id'))

      const user = await login(data as LoginRequest)
      console.log('✅ [NORMAL LOGIN] Login successful!')
      console.log('✅ [NORMAL LOGIN] User logged in with role:', user.role)
      console.log('✅ [NORMAL LOGIN] Full user data:', user)
      
      // Check if there's a pending cart conflict
      if (pendingCartConflict) {
        toast.success('Login successful! Please choose how to merge your carts.')
      } else {
        toast.success('Welcome back!')
      }
      
      // Check for redirect URL from query params
      const redirectUrl = new URLSearchParams(window.location.search).get('redirect')

      // Check admin access using permission system
      const canAccessAdmin = canAccessAdminPanel(user.role)
      console.log('Can access admin panel:', canAccessAdmin)

      // Determine redirect destination
      let destination = '/'
      if (redirectUrl) {
        destination = redirectUrl
      } else if (canAccessAdmin) {
        destination = '/admin'
      }

      console.log('Redirecting to:', destination)

      // Use setTimeout to ensure auth state is fully updated before redirect
      setTimeout(() => {
        router.replace(destination)
      }, 100)
      
    } catch (error: any) {
      console.error('Login error:', error)

      // Check if error is related to email verification
      if (error.message && error.message.includes('email not verified')) {
        toast.error('Please verify your email before logging in', {
          action: {
            label: 'Resend Email',
            onClick: () => router.push('/auth/resend-verification'),
          },
        })
        return
      }

      if (error.code === 'VALIDATION_ERROR' && error.details) {
        // Handle field-specific validation errors
        Object.entries(error.details).forEach(([field, message]) => {
          setError(field as keyof LoginFormData, {
            type: 'server',
            message: message as string,
          })
        })
      } else {
        toast.error(error.message || 'Login failed')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
      {/* Email */}
      <div className="space-y-1">
        <Input
          {...register('email')}
          type="email"
          label="Email address"
          placeholder="Enter your email address"
          error={errors.email?.message}
          required
          autoComplete="email"
          size="lg"
          className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
        />
      </div>

      {/* Password */}
      <div className="space-y-1">
        <Input
          {...register('password')}
          type={showPassword ? 'text' : 'password'}
          label="Password"
          placeholder="Enter your password"
          error={errors.password?.message}
          required
          autoComplete="current-password"
          size="lg"
          className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
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
      </div>

      {/* Remember me & Forgot password */}
      <div className="flex items-center justify-between pt-1">
        <div className="flex items-center">
          <input
            {...register('remember_me')}
            id="remember-me"
            type="checkbox"
            className="h-3 w-3 border-gray-600 bg-gray-800 rounded transition-colors focus:ring-2 focus:ring-[#FF9000]/20"
            style={{accentColor: '#FF9000'}}
          />
          <label htmlFor="remember-me" className="ml-2 block text-xs font-medium text-gray-200">
            Remember me
          </label>
        </div>

        <Link
          href="/auth/forgot-password"
          className="text-xs font-medium text-[#FF9000] hover:text-[#e67e00] transition-all duration-200 hover:underline"
        >
          Forgot password?
        </Link>
      </div>

      {/* Submit button */}
      <Button
        type="submit"
        className="w-full mt-4 bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold py-3 text-sm rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg border-0"
        size="lg"
        isLoading={isLoading}
        loadingText="Signing in..."
      >
        Sign in to your account
      </Button>

      {/* Social login divider */}
      <div className="relative mt-4">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-gray-600/60"></div>
        </div>
        <div className="relative flex justify-center text-xs">
          <span className="px-4 bg-gray-900/90 text-gray-400 font-medium backdrop-blur-sm">Or continue with</span>
        </div>
      </div>

      {/* OAuth login buttons */}
      <div className="mt-4">
        <OAuthButtons
          onError={(error) => toast.error(error)}
        />
      </div>
    </form>
  )
}
