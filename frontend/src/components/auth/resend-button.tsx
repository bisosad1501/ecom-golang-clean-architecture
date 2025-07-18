'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Mail, Clock } from 'lucide-react'
import { authApi } from '@/lib/api'
import { toast } from 'sonner'

interface ResendButtonProps {
  email: string
  onSuccess?: () => void
  variant?: 'default' | 'outline' | 'ghost'
  className?: string
  size?: 'sm' | 'default' | 'lg'
}

export function ResendButton({
  email,
  onSuccess,
  variant = 'outline',
  className = '',
  size = 'default'
}: ResendButtonProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const [canResend, setCanResend] = useState(true)

  useEffect(() => {
    let timer: NodeJS.Timeout
    if (countdown > 0) {
      timer = setTimeout(() => {
        setCountdown(countdown - 1)
      }, 1000)
    } else if (countdown === 0 && !canResend) {
      setCanResend(true)
    }
    return () => clearTimeout(timer)
  }, [countdown, canResend])

  const handleResend = async () => {
    if (!email) {
      toast.error('Email address is required')
      return
    }

    try {
      setIsLoading(true)

      await authApi.resendVerification(email)

      toast.success('Verification email sent successfully!')

      // Start countdown
      setCountdown(60)
      setCanResend(false)

      onSuccess?.()

    } catch (error: any) {
      toast.error(error.message || 'Failed to send verification email')
    } finally {
      setIsLoading(false)
    }
  }

  const isDisabled = isLoading || !canResend || !email

  return (
    <Button
      onClick={handleResend}
      disabled={isDisabled}
      variant={variant}
      size={size}
      className={`${className} ${
        variant === 'outline' 
          ? 'border-gray-600 text-gray-300 hover:bg-gray-700 hover:text-white disabled:opacity-50' 
          : ''
      }`}
    >
      {isLoading ? (
        <>
          <Mail className="h-4 w-4 mr-2 animate-pulse" />
          Sending...
        </>
      ) : !canResend ? (
        <>
          <Clock className="h-4 w-4 mr-2" />
          Resend in {countdown}s
        </>
      ) : (
        <>
          <Mail className="h-4 w-4 mr-2" />
          Resend Verification Email
        </>
      )}
    </Button>
  )
}
