'use client'

import React from 'react'
import { Label } from '@/components/ui/label'
import { AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

interface FormFieldProps {
  label: string
  required?: boolean
  error?: string
  hint?: string
  children: React.ReactNode
  className?: string
}

export function FormField({
  label,
  required = false,
  error,
  hint,
  children,
  className
}: FormFieldProps) {
  const fieldId = React.useId()

  return (
    <div className={cn('space-y-2', className)}>
      <Label
        htmlFor={fieldId}
        className={cn(
          BIHUB_ADMIN_THEME.typography.body.medium,
          "font-medium text-white"
        )}
      >
        {label}
        {required && <span className="text-red-400 ml-1">*</span>}
      </Label>

      {React.cloneElement(children as React.ReactElement, {
        ...((children as any)?.props || {}),
        id: fieldId,
        className: cn(
          'bg-gray-700/50 border-gray-600/50 text-white placeholder:text-gray-400',
          'focus:border-[#FF9000] focus:ring-[#FF9000]/20',
          'rounded-lg transition-colors duration-200',
          (children as any)?.props?.className,
          error && 'border-red-400 focus:border-red-400'
        ),
      } as any)}

      {error && (
        <div className="flex items-center text-sm text-red-400">
          <AlertCircle className="h-4 w-4 mr-1 flex-shrink-0" />
          <span>{error}</span>
        </div>
      )}

      {hint && !error && (
        <p className={cn(BIHUB_ADMIN_THEME.typography.body.small, "text-gray-400")}>
          {hint}
        </p>
      )}
    </div>
  )
}
