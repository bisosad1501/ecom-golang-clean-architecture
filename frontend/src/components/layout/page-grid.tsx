/**
 * Page Grid Component
 * Provides consistent grid layouts for different content types
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { getGridClasses } from '@/constants/page-layouts'

interface PageGridProps {
  children: ReactNode
  type?: 'products' | 'categories' | 'features' | 'twoColumn' | 'threeColumn'
  className?: string
}

export function PageGrid({ 
  children, 
  type = 'products', 
  className 
}: PageGridProps) {
  return (
    <div className={cn(getGridClasses(type), className)}>
      {children}
    </div>
  )
}
