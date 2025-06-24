'use client'

import React from 'react'
import { cn } from '@/lib/utils'

interface AdminFormLayoutProps {
  title: string
  description?: string
  children: React.ReactNode
  className?: string
  maxHeight?: string
}

export function AdminFormLayout({
  title,
  description,
  children,
  className,
  maxHeight = 'max-h-[75vh]',
}: AdminFormLayoutProps) {
  return (
    <div className={cn('space-y-6', className)}>
      {/* Form Header */}
      <div className="space-y-2">
        <h2 className="text-lg font-semibold text-gray-900">{title}</h2>
        {description && (
          <p className="text-sm text-gray-600">{description}</p>
        )}
      </div>

      {/* Form Content */}
      <div className={cn('overflow-y-auto', maxHeight)}>
        <div className="space-y-6">
          {children}
        </div>
      </div>
    </div>
  )
}
