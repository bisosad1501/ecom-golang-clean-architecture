'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Mail, ArrowLeft, CheckCircle } from 'lucide-react'
import { authApi } from '@/lib/api'
import { toast } from 'sonner'
import { useGuestOnly } from '@/hooks/use-auth-guard'

const resendSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
})

type ResendFormData = z.infer<typeof resendSchema>

export default function ResendVerificationPage() {
  useGuestOnly()

  const router = useRouter()
  const [isLoading, setIsLoading] = useState(false)
  const [isSuccess, setIsSuccess] = useState(false)
  const [sentEmail, setSentEmail] = useState('')

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ResendFormData>({
    resolver: zodResolver(resendSchema),
    defaultValues: {
      email: '',
    },
  })

  const onSubmit = async (data: ResendFormData) => {
    try {
      setIsLoading(true)

      await authApi.resendVerification(data.email)

      setIsSuccess(true)
      setSentEmail(data.email)
      toast.success('Verification email sent successfully!')

    } catch (error: any) {
      toast.error(error.message || 'Failed to send verification email')
    } finally {
      setIsLoading(false)
    }
  }

  const handleBackToLogin = () => {
    router.push('/auth/login')
  }

  if (isSuccess) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 py-12 px-4 sm:px-6 lg:px-8">
        <div className="max-w-md w-full space-y-8">
          <Card className="bg-gray-800/90 border-gray-700/50 backdrop-blur-sm shadow-2xl">
            <CardHeader className="text-center">
              <CardTitle className="text-2xl font-bold text-white">Check Your Email</CardTitle>
              <CardDescription className="text-gray-300">
                We've sent a verification link to your email
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="text-center space-y-4">
                <CheckCircle className="h-12 w-12 mx-auto text-green-400" />
                <h2 className="text-xl font-semibold text-green-400">Email Sent Successfully!</h2>
                <p className="text-gray-300">
                  We've sent a verification link to:
                </p>
                <p className="font-medium text-white">{sentEmail}</p>
                <div className="bg-blue-900/20 border border-blue-500/30 rounded-lg p-4 backdrop-blur-sm">
                  <p className="text-sm text-blue-200 font-medium">
                    <strong>Next steps:</strong>
                  </p>
                  <ol className="text-sm text-blue-300 mt-2 space-y-1">
                    <li>1. Check your email inbox</li>
                    <li>2. Click the verification link in the email</li>
                    <li>3. Return here to log in</li>
                  </ol>
                </div>
                <p className="text-xs text-gray-400">
                  Didn't receive the email? Check your spam folder or try again.
                </p>
                <Button
                  onClick={handleBackToLogin}
                  className="w-full bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold border-0"
                >
                  Back to Login
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <Card className="bg-gray-800/90 border-gray-700/50 backdrop-blur-sm shadow-2xl">
          <CardHeader className="text-center">
            <CardTitle className="text-2xl font-bold text-white">Resend Verification Email</CardTitle>
            <CardDescription className="text-gray-300">
              Enter your email address to receive a new verification link
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Input
                  {...register('email')}
                  type="email"
                  label="Email Address"
                  placeholder="Enter your email address"
                  error={errors.email?.message}
                  required
                  autoComplete="email"
                  size="lg"
                  className="transition-all duration-300 focus:scale-[1.01] bg-gray-800/90 border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 hover:border-gray-500 h-10 text-sm rounded-lg backdrop-blur-sm"
                />
              </div>

              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold border-0 transition-all duration-300 transform hover:scale-[1.02] hover:shadow-xl shadow-lg"
                disabled={isLoading}
                size="lg"
              >
                {isLoading ? (
                  <>
                    <Mail className="h-4 w-4 mr-2 animate-pulse" />
                    Sending...
                  </>
                ) : (
                  <>
                    <Mail className="h-4 w-4 mr-2" />
                    Send Verification Email
                  </>
                )}
              </Button>
            </form>

            <div className="text-center">
              <Button
                variant="ghost"
                onClick={handleBackToLogin}
                className="text-sm text-gray-400 hover:text-white hover:bg-gray-700/50"
              >
                <ArrowLeft className="h-4 w-4 mr-2" />
                Back to Login
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
