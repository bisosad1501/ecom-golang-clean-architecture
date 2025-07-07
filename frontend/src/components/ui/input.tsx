import * as React from 'react'
import { cva, type VariantProps } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const inputVariants = cva(
  'flex w-full rounded-xl border bg-gray-800 px-4 py-3 text-sm text-white ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-gray-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50 transition-all duration-200 shadow-sm hover:shadow-medium',
  {
    variants: {
      variant: {
        default: 'border-gray-600 focus-visible:ring-[#FF9000]/30 focus-visible:border-[#FF9000]',
        error: 'border-red-500 focus-visible:ring-red-500/30 focus-visible:border-red-500',
        success: 'border-green-500 focus-visible:ring-green-500/30 focus-visible:border-green-500',
      },
      size: {
        default: 'h-12 px-4 py-3',
        sm: 'h-10 px-3 py-2 text-sm rounded-lg',
        lg: 'h-14 px-5 py-4 text-base',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
)

export interface InputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
    VariantProps<typeof inputVariants> {
  label?: string
  error?: string
  helperText?: string
  leftIcon?: React.ReactNode
  rightIcon?: React.ReactNode
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  (
    {
      className,
      variant,
      size,
      type,
      label,
      error,
      helperText,
      leftIcon,
      rightIcon,
      id,
      ...props
    },
    ref
  ) => {
    const inputId = id || React.useId()
    const hasError = !!error
    const finalVariant = hasError ? 'error' : variant

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="mb-3 block text-sm font-semibold text-white"
          >
            {label}
            {props.required && <span className="ml-1 text-red-400">*</span>}
          </label>
        )}
        <div className="relative">
          {leftIcon && (
            <div className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400">
              {leftIcon}
            </div>
          )}
          <input
            type={type}
            className={cn(
              inputVariants({ variant: finalVariant, size }),
              leftIcon && 'pl-12',
              rightIcon && 'pr-12',
              className
            )}
            ref={ref}
            id={inputId}
            {...props}
          />
          {rightIcon && (
            <div className="absolute right-4 top-1/2 -translate-y-1/2 text-gray-400">
              {rightIcon}
            </div>
          )}
        </div>
        {(error || helperText) && (
          <p
            className={cn(
              'mt-2 text-sm font-medium',
              hasError ? 'text-red-400' : 'text-gray-400'
            )}
          >
            {error || helperText}
          </p>
        )}
      </div>
    )
  }
)
Input.displayName = 'Input'

export { Input, inputVariants }
