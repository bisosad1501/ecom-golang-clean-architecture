'use client'

import React from 'react'
import { cn } from '@/lib/utils'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

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
    <div className={cn(BIHUB_ADMIN_THEME.spacing.section, className)}>
      {/* BiHub Form Header */}
      <div className="space-y-2 mb-8">
        <h2 className={cn(BIHUB_ADMIN_THEME.typography.heading.h2, "text-2xl")}>
          {title}
        </h2>
        {description && (
          <p className={BIHUB_ADMIN_THEME.typography.body.large}>
            {description}
          </p>
        )}
      </div>

      {/* BiHub Form Content */}
      <div className={cn('overflow-y-auto', maxHeight)}>
        <div className={BIHUB_ADMIN_THEME.spacing.section}>
          {children}
        </div>
      </div>
    </div>
  )
}
