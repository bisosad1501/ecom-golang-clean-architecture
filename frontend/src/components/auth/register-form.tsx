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
import { getHighContrastClasses, PAGE_CONTRAST } from '@/constants/contrast-system'

const registerSchema = z.object({
  first_name: z.string().min(2, 'First name must be at least 2 characters'),
  last_name: z.string().min(2, 'Last name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address'),
  password: z.string()
    .min(8, 'Password must be at least 8 characters')
    .regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, 'Password must contain at least one uppercase letter, one lowercase letter, and one number'),
  confirm_password: z.string(),
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
    },
  })

  const onSubmit = async (data: RegisterFormData) => {
    try {
      clearError()
      const { confirm_password, ...registerData } = data
      await registerUser(registerData as RegisterRequest)
      toast.success('Account created successfully! Please login to continue.')
      // Redirect to login page since backend register doesn't auto-login
      router.push('/auth/login')
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
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-3">
      {/* Name fields */}
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <Input
            {...register('first_name')}
            type="text"
            label="First name"
            placeholder="John"
            error={errors.first_name?.message}
            required
            autoComplete="given-name"
            size="lg"
            className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
          />
        </div>
        <div className="space-y-1">
          <Input
            {...register('last_name')}
            type="text"
            label="Last name"
            placeholder="Doe"
            error={errors.last_name?.message}
            required
            autoComplete="family-name"
            size="lg"
            className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
          />
        </div>
      </div>

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

      {/* Password fields */}
      <div className="grid grid-cols-2 gap-3">
        <div className="space-y-1">
          <Input
            {...register('password')}
            type={showPassword ? 'text' : 'password'}
            label="Password"
            placeholder="Create password"
            error={errors.password?.message}
            required
            autoComplete="new-password"
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

        <div className="space-y-1">
          <Input
            {...register('confirm_password')}
            type={showConfirmPassword ? 'text' : 'password'}
            label="Confirm"
            placeholder="Confirm"
            error={errors.confirm_password?.message}
            required
            autoComplete="new-password"
            size="lg"
            className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
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
      </div>

      {/* Submit button */}
      <Button
        type="submit"
        className="w-full mt-4 bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold py-3 text-sm rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg border-0"
        size="lg"
        isLoading={isLoading}
        loadingText="Creating account..."
      >
        Create your account
      </Button>
    </form>
  )
}
