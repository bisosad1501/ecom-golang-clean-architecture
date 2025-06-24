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

      {/* Social login buttons */}
      <div className="grid grid-cols-2 gap-4">
        <Button
          type="button"
          variant="outline"
          size="lg"
          className="border-2 hover:border-primary transition-all duration-200"
          onClick={() => toast.info('Google login coming soon!')}
        >
          <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
            <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
            <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
            <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
            <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
          </svg>
          Google
        </Button>
        <Button
          type="button"
          variant="outline"
          size="lg"
          className="border-2 hover:border-primary transition-all duration-200"
          onClick={() => toast.info('Facebook login coming soon!')}
        >
          <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
            <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
          </svg>
          Facebook
        </Button>
      </div>
    </form>
  )
}
