/**
 * Page Container Component
 * Provides consistent container with max-width and padding
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { getContainerClasses } from '@/constants/page-layouts'

interface PageContainerProps {
  children: ReactNode
  className?: string
}

export function PageContainer({ children, className }: PageContainerProps) {
  return (
    <div className={cn(getContainerClasses(), className)}>
      {children}
    </div>
  )
}
