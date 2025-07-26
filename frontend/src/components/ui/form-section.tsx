'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

interface FormSectionProps {
  title: string
  description?: string
  children: React.ReactNode
  className?: string
}

export function FormSection({ title, description, children, className }: FormSectionProps) {
  return (
    <Card className={cn(
      'bg-gray-800/50 border-gray-700/50 rounded-xl shadow-lg',
      className
    )}>
      <CardHeader className="pb-4">
        <CardTitle className={cn(
          BIHUB_ADMIN_THEME.typography.heading.h4,
          "text-lg"
        )}>
          {title}
        </CardTitle>
        {description && (
          <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
            {description}
          </p>
        )}
      </CardHeader>
      <CardContent className={BIHUB_ADMIN_THEME.spacing.form}>
        {children}
      </CardContent>
    </Card>
  )
}
