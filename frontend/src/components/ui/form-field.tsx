'use client'

import React from 'react'
import { Label } from '@/components/ui/label'
import { AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

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
      <Label htmlFor={fieldId}>
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </Label>
      
      {React.cloneElement(children as React.ReactElement, {
        ...((children as any)?.props || {}),
        id: fieldId,
        className: cn(
          (children as any)?.props?.className,
          error && 'border-red-500'
        ),
      } as any)}
      
      {error && (
        <div className="flex items-center text-sm text-red-600">
          <AlertCircle className="h-4 w-4 mr-1 flex-shrink-0" />
          <span>{error}</span>
        </div>
      )}
      
      {hint && !error && (
        <p className="text-xs text-gray-500">{hint}</p>
      )}
    </div>
  )
}
