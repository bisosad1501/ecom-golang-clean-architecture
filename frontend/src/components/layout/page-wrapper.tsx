/**
 * Page Wrapper Component
 * Provides consistent page structure and styling
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { getPageWrapperClasses } from '@/constants/page-layouts'

interface PageWrapperProps {
  children: ReactNode
  className?: string
}

export function PageWrapper({ children, className }: PageWrapperProps) {
  return (
    <div className={cn(getPageWrapperClasses(), className)}>
      {children}
    </div>
  )
}
