'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Eye, EyeOff } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useAuthStore } from '@/store/auth'
import { RegisterRequest } from '@/types'
import { toast } from 'sonner'

const registerSchema = z.object({
  first_name: z.string().min(2, 'First name must be at least 2 characters'),
  last_name: z.string().min(2, 'Last name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address'),
  password: z.string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, 'Password must contain at least one uppercase letter, one lowercase letter, and one number'),
  confirm_password: z.string(),
  terms_accepted: z.boolean().refine(val => val === true, 'You must accept the terms and conditions'),
}).refine((data) => data.password === data.confirm_password, {
  message: "Passwords don't match",
  path: ["confirm_password"],
})

type RegisterFormData = z.infer<typeof registerSchema>

export function RegisterForm() {
  const router = useRouter()
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const { register: registerUser, isLoading, error, clearError } = useAuthStore()

  const {
    register,
    handleSubmit,
    formState: { errors },
    setError,
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      first_name: '',
      last_name: '',
      email: '',
      password: '',
      confirm_password: '',
      terms_accepted: false,
    },
  })

  const onSubmit = async (data: RegisterFormData) => {
    try {
      clearError()
      const { confirm_password, ...registerData } = data
      await registerUser(registerData as RegisterRequest)
      toast.success('Account created successfully! Welcome to our store!')
      router.push('/')
    } catch (error: any) {
      if (error.code === 'VALIDATION_ERROR' && error.details) {
        // Handle field-specific validation errors
        Object.entries(error.details).forEach(([field, message]) => {
          setError(field as keyof RegisterFormData, {
            type: 'server',
            message: message as string,
          })
        })
      } else {
        toast.error(error.message || 'Registration failed')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      {/* Name fields */}
      <div className="grid grid-cols-2 gap-4">
        <Input
          {...register('first_name')}
          type="text"
          label="First name"
          placeholder="John"
          error={errors.first_name?.message}
          required
          autoComplete="given-name"
        />
        <Input
          {...register('last_name')}
          type="text"
          label="Last name"
          placeholder="Doe"
          error={errors.last_name?.message}
          required
          autoComplete="family-name"
        />
      </div>

      {/* Email */}
      <Input
        {...register('email')}
        type="email"
        label="Email address"
        placeholder="john@example.com"
        error={errors.email?.message}
        required
        autoComplete="email"
      />

      {/* Password */}
      <Input
        {...register('password')}
        type={showPassword ? 'text' : 'password'}
        label="Password"
        placeholder="Create a strong password"
        error={errors.password?.message}
        required
        autoComplete="new-password"
        helperText="Must contain at least 8 characters with uppercase, lowercase, and number"
        rightIcon={
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="text-gray-400 hover:text-gray-600"
          >
            {showPassword ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
          </button>
        }
      />

      {/* Confirm Password */}
      <Input
        {...register('confirm_password')}
        type={showConfirmPassword ? 'text' : 'password'}
        label="Confirm password"
        placeholder="Confirm your password"
        error={errors.confirm_password?.message}
        required
        autoComplete="new-password"
        rightIcon={
          <button
            type="button"
            onClick={() => setShowConfirmPassword(!showConfirmPassword)}
            className="text-gray-400 hover:text-gray-600"
          >
            {showConfirmPassword ? (
              <EyeOff className="h-4 w-4" />
            ) : (
              <Eye className="h-4 w-4" />
            )}
          </button>
        }
      />

      {/* Terms and conditions */}
      <div className="flex items-start">
        <input
          {...register('terms_accepted')}
          id="terms"
          type="checkbox"
          className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded mt-1"
        />
        <label htmlFor="terms" className="ml-2 block text-sm text-gray-900">
          I agree to the{' '}
          <a href="/terms" className="text-primary-600 hover:text-primary-500" target="_blank">
            Terms of Service
          </a>{' '}
          and{' '}
          <a href="/privacy" className="text-primary-600 hover:text-primary-500" target="_blank">
            Privacy Policy
          </a>
        </label>
      </div>
      {errors.terms_accepted && (
        <p className="text-sm text-error-600">{errors.terms_accepted.message}</p>
      )}

      {/* Submit button */}
      <div className="mt-6">
        <Button
          type="submit"
          className="w-full"
          size="lg"
          isLoading={isLoading}
          loadingText="Creating account..."
        >
          Create account
        </Button>
      </div>

      {/* Social registration */}
      <div className="mt-6">
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-white text-gray-500">Or sign up with</span>
          </div>
        </div>

        <div className="mt-6 grid grid-cols-2 gap-3">
          <Button
            type="button"
            variant="outline"
            className="w-full"
            onClick={() => toast.info('Google signup coming soon!')}
          >
            <svg className="h-5 w-5 mr-2" viewBox="0 0 24 24">
              <path
                fill="currentColor"
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
              />
              <path
                fill="currentColor"
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              />
              <path
                fill="currentColor"
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
              />
              <path
                fill="currentColor"
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
              />
            </svg>
            Google
          </Button>

          <Button
            type="button"
            variant="outline"
            className="w-full"
            onClick={() => toast.info('Facebook signup coming soon!')}
          >
            <svg className="h-5 w-5 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
            </svg>
            Facebook
          </Button>
        </div>
      </div>
    </form>
  )
}
