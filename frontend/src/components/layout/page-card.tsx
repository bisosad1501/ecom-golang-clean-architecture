/**
 * Page Card Component
 * Provides consistent card styling across all pages
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { getCardClasses } from '@/constants/page-layouts'

interface PageCardProps {
  children: ReactNode
  padding?: 'sm' | 'base' | 'lg'
  shadow?: 'none' | 'sm' | 'base' | 'lg'
  className?: string
}

export function PageCard({ 
  children, 
  padding = 'base', 
  shadow = 'none',
  className 
}: PageCardProps) {
  return (
    <div className={cn(getCardClasses(padding, shadow), className)}>
      {children}
    </div>
  )
}
