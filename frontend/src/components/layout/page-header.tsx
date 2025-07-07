/**
 * Page Header Component
 * Provides consistent page header structure with title, subtitle, and breadcrumbs
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { 
  getPageHeaderClasses, 
  getContainerClasses, 
  getPageTitleClasses,
  PAGE_LAYOUTS 
} from '@/constants/page-layouts'

interface PageHeaderProps {
  title: string
  subtitle?: string
  breadcrumbs?: ReactNode
  actions?: ReactNode
  size?: 'sm' | 'base' | 'lg'
  className?: string
}

export function PageHeader({ 
  title, 
  subtitle, 
  breadcrumbs, 
  actions, 
  size = 'base',
  className 
}: PageHeaderProps) {
  return (
    <div className={cn(getPageHeaderClasses(), className)}>
      <div className={getContainerClasses()}>
        {/* Breadcrumbs */}
        {breadcrumbs && (
          <div className={PAGE_LAYOUTS.PAGE_HEADER.breadcrumbs}>
            {breadcrumbs}
          </div>
        )}

        {/* Header content */}
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
          <div>
            {/* Title */}
            <h1 className={getPageTitleClasses(size)}>
              {title}
            </h1>
            
            {/* Subtitle */}
            {subtitle && (
              <p className={PAGE_LAYOUTS.PAGE_HEADER.subtitle}>
                {subtitle}
              </p>
            )}
          </div>

          {/* Actions */}
          {actions && (
            <div className="flex-shrink-0">
              {actions}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
