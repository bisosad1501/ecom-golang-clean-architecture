/**
 * Page Section Component
 * Provides consistent section structure with optional title and container
 */

import { ReactNode } from 'react'
import { cn } from '@/lib/utils'
import { 
  getSectionClasses, 
  getContainerClasses, 
  getTypographyClasses,
  PAGE_LAYOUTS 
} from '@/constants/page-layouts'

interface PageSectionProps {
  children: ReactNode
  title?: string
  subtitle?: string
  size?: 'sm' | 'base' | 'lg' | 'xl'
  container?: boolean
  className?: string
  titleClassName?: string
  contentClassName?: string
}

export function PageSection({ 
  children, 
  title, 
  subtitle, 
  size = 'base',
  container = true,
  className,
  titleClassName,
  contentClassName
}: PageSectionProps) {
  const sectionContent = (
    <>
      {/* Section header */}
      {(title || subtitle) && (
        <div className="text-center mb-6 lg:mb-8">
          {title && (
            <h2 className={cn(
              getTypographyClasses('sectionTitle', size === 'xl' ? 'lg' : size),
              titleClassName
            )}>
              {title}
            </h2>
          )}
          {subtitle && (
            <p className={cn(
              getTypographyClasses('body', 'base'),
              'mt-2 max-w-2xl mx-auto',
              titleClassName
            )}>
              {subtitle}
            </p>
          )}
        </div>
      )}

      {/* Section content */}
      <div className={cn(contentClassName)}>
        {children}
      </div>
    </>
  )

  return (
    <section className={cn(getSectionClasses(size), className)}>
      {container ? (
        <div className={getContainerClasses()}>
          {sectionContent}
        </div>
      ) : (
        sectionContent
      )}
    </section>
  )
}
