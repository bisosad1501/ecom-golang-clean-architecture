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
import { LoginRequest } from '@/types'
import { toast } from 'sonner'
import { OAuthButtons } from './OAuthButtons'

const loginSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  remember_me: z.boolean().optional(),
})

type LoginFormData = z.infer<typeof loginSchema>

export function LoginForm() {
  const router = useRouter()
  const [showPassword, setShowPassword] = useState(false)
  const { login, isLoading, error, clearError } = useAuthStore()

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
      await login(data as LoginRequest)
      
      // Get fresh user data after login
      const userRole = useAuthStore.getState().user?.role
      console.log('User logged in with role:', userRole)
      
      toast.success('Welcome back!')
      
      // Redirect based on user role immediately
      if (userRole === 'admin') {
        console.log('Redirecting to admin panel')
        window.location.href = '/admin' // Force full page navigation for admin
      } else {
        console.log('Redirecting to home')
        router.replace('/') // Use replace for cleaner history
      }
    } catch (error: any) {
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
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
      {/* Email */}
      <div className="space-y-2">
        <Input
          {...register('email')}
          type="email"
          label="Email address"
          placeholder="Enter your email address"
          error={errors.email?.message}
          required
          autoComplete="email"
          size="lg"
          className="transition-all duration-200 focus:scale-[1.02]"
        />
      </div>

      {/* Password */}
      <div className="space-y-2">
        <Input
          {...register('password')}
          type={showPassword ? 'text' : 'password'}
          label="Password"
          placeholder="Enter your password"
          error={errors.password?.message}
          required
          autoComplete="current-password"
          size="lg"
          className="transition-all duration-200 focus:scale-[1.02]"
          rightIcon={
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="text-muted-foreground hover:text-primary transition-colors p-1 rounded-lg hover:bg-muted"
            >
              {showPassword ? (
                <EyeOff className="h-5 w-5" />
              ) : (
                <Eye className="h-5 w-5" />
              )}
            </button>
          }
        />
      </div>

      {/* Remember me & Forgot password */}
      <div className="flex items-center justify-between">
        <div className="flex items-center">
          <input
            {...register('remember_me')}
            id="remember-me"
            type="checkbox"
            className="h-4 w-4 text-primary focus:ring-primary/30 border-border rounded transition-colors"
          />
          <label htmlFor="remember-me" className="ml-3 block text-sm font-medium text-foreground">
            Remember me
          </label>
        </div>

        <Link
          href="/auth/forgot-password"
          className="text-sm font-medium text-primary hover:text-primary-600 transition-colors"
        >
          Forgot password?
        </Link>
      </div>

      {/* Submit button */}
      <Button
        type="submit"
        className="w-full"
        size="xl"
        variant="gradient"
        isLoading={isLoading}
        loadingText="Signing you in..."
      >
        Sign in to your account
      </Button>

      {/* Social login divider */}
      <div className="relative">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-border"></div>
        </div>
        <div className="relative flex justify-center text-sm">
          <span className="px-4 bg-background text-muted-foreground font-medium">Or continue with</span>
        </div>
      </div>

      {/* OAuth login buttons */}
      <OAuthButtons
        onError={(error) => toast.error(error)}
      />
    </form>
  )
}
