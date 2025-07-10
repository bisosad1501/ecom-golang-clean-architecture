'use client'

import { ReactNode } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { BIHUB_ADMIN_THEME, BIHUB_BRAND, getStatusColor, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'
import { TrendingUp, Package, Users, ShoppingCart, Star, Crown } from 'lucide-react'

// BiHub Admin Logo Component
export function BiHubAdminLogo({ size = 'default' }: { size?: 'small' | 'default' | 'large' }) {
  const sizes = {
    small: { container: 'h-8 w-8', text: 'text-lg' },
    default: { container: 'h-12 w-12', text: 'text-2xl' },
    large: { container: 'h-16 w-16', text: 'text-3xl' },
  }
  
  return (
    <div className="flex items-center space-x-3 group">
      <div className={cn(
        'rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center shadow-large group-hover:scale-105 transition-transform duration-200',
        sizes[size].container
      )}>
        <span className={cn('text-white font-bold', sizes[size].text)}>A</span>
      </div>
      <div>
        <span className={cn(BIHUB_ADMIN_THEME.typography.heading.h3, 'flex items-center')}>
          <span className={BIHUB_BRAND.logo.colors.text}>{BIHUB_BRAND.logo.text}</span>
          <span className={cn(
            'ml-0.5 px-1.5 py-0.5 rounded-[2px] font-bold',
            BIHUB_BRAND.logo.colors.accent
          )} style={{letterSpacing: '0.3px'}}>
            {BIHUB_BRAND.logo.accent}
          </span>
          <span className="ml-2 text-gray-400 text-sm font-normal">Admin</span>
        </span>
      </div>
    </div>
  )
}

// BiHub Admin Card Component
interface BiHubAdminCardProps {
  title: string
  subtitle?: string
  icon?: ReactNode
  children: ReactNode
  className?: string
  headerAction?: ReactNode
}

export function BiHubAdminCard({ 
  title, 
  subtitle, 
  icon, 
  children, 
  className,
  headerAction 
}: BiHubAdminCardProps) {
  return (
    <Card className={cn(
      BIHUB_ADMIN_THEME.components.card.base,
      BIHUB_ADMIN_THEME.components.card.hover,
      className
    )}>
      <CardHeader className="pb-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {icon && (
              <div className="w-10 h-10 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                {icon}
              </div>
            )}
            <div>
              <CardTitle className={BIHUB_ADMIN_THEME.typography.heading.h3}>
                {title}
              </CardTitle>
              {subtitle && (
                <p className={BIHUB_ADMIN_THEME.typography.body.small}>
                  {subtitle}
                </p>
              )}
            </div>
          </div>
          {headerAction}
        </div>
      </CardHeader>
      <CardContent className={BIHUB_ADMIN_THEME.components.card.padding}>
        {children}
      </CardContent>
    </Card>
  )
}

// BiHub Status Badge Component
interface BiHubStatusBadgeProps {
  status: string
  children: ReactNode
}

export function BiHubStatusBadge({ status, children }: BiHubStatusBadgeProps) {
  const variant = getBadgeVariant(status)
  const badgeClass = BIHUB_ADMIN_THEME.components.badge[variant as keyof typeof BIHUB_ADMIN_THEME.components.badge]
  
  return (
    <Badge className={badgeClass}>
      {children}
    </Badge>
  )
}

// BiHub Stat Card Component
interface BiHubStatCardProps {
  title: string
  value: string | number
  change?: number
  icon: ReactNode
  color: 'primary' | 'success' | 'warning' | 'error' | 'info'
}

export function BiHubStatCard({ title, value, change, icon, color }: BiHubStatCardProps) {
  const colorClasses = {
    primary: 'from-[#FF9000] to-[#e67e00]',
    success: 'from-emerald-500 to-emerald-600',
    warning: 'from-yellow-500 to-yellow-600',
    error: 'from-red-500 to-red-600',
    info: 'from-blue-500 to-blue-600',
  }
  
  const changeColorClasses = {
    primary: 'bg-orange-900/30 text-[#FF9000]',
    success: 'bg-emerald-900/30 text-emerald-400',
    warning: 'bg-yellow-900/30 text-yellow-400',
    error: 'bg-red-900/30 text-red-400',
    info: 'bg-blue-900/30 text-blue-400',
  }
  
  return (
    <Card className={cn(
      BIHUB_ADMIN_THEME.components.card.base,
      'hover:shadow-xl transition-all duration-300 group'
    )}>
      <CardContent className="p-8">
        <div className="flex items-center justify-between mb-6">
          <div className={cn(
            'w-16 h-16 rounded-3xl bg-gradient-to-br flex items-center justify-center shadow-large group-hover:scale-110 transition-transform duration-300',
            colorClasses[color]
          )}>
            {icon}
          </div>
          {change !== undefined && (
            <div className={cn(
              'flex items-center gap-2 px-3 py-1 rounded-full',
              changeColorClasses[color]
            )}>
              <TrendingUp className="h-4 w-4" />
              <span className="text-sm font-semibold">+{change}%</span>
            </div>
          )}
        </div>

        <div>
          <p className={cn(BIHUB_ADMIN_THEME.typography.body.small, 'mb-2')}>
            {title}
          </p>
          <p className={cn(BIHUB_ADMIN_THEME.typography.heading.h1, 'text-3xl')}>
            {value}
          </p>
        </div>
      </CardContent>
    </Card>
  )
}

// BiHub Action Button Component
interface BiHubActionButtonProps {
  icon: ReactNode
  title: string
  description?: string
  onClick?: () => void
  variant?: 'primary' | 'secondary'
}

export function BiHubActionButton({ 
  icon, 
  title, 
  description, 
  onClick,
  variant = 'secondary' 
}: BiHubActionButtonProps) {
  return (
    <Button
      variant="outline"
      onClick={onClick}
      className={cn(
        'h-24 flex-col gap-3 transition-all duration-200 group',
        variant === 'primary' 
          ? BIHUB_ADMIN_THEME.components.button.primary
          : BIHUB_ADMIN_THEME.components.button.secondary
      )}
    >
      <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center group-hover:scale-110 transition-transform duration-200">
        {icon}
      </div>
      <div className="text-center">
        <span className="font-semibold block">{title}</span>
        {description && (
          <span className="text-xs opacity-70">{description}</span>
        )}
      </div>
    </Button>
  )
}

// BiHub Page Header Component
interface BiHubPageHeaderProps {
  title: string
  subtitle?: string
  action?: ReactNode
  breadcrumbs?: Array<{ label: string; href?: string }>
}

export function BiHubPageHeader({ title, subtitle, action, breadcrumbs }: BiHubPageHeaderProps) {
  return (
    <div className="mb-8">
      {breadcrumbs && (
        <nav className="flex items-center space-x-2 text-sm text-gray-400 mb-4">
          {breadcrumbs.map((crumb, index) => (
            <div key={index} className="flex items-center">
              {index > 0 && <span className="mx-2">/</span>}
              <span className={index === breadcrumbs.length - 1 ? 'text-[#FF9000]' : ''}>
                {crumb.label}
              </span>
            </div>
          ))}
        </nav>
      )}
      
      <div className="flex items-center justify-between">
        <div>
          <h1 className={BIHUB_ADMIN_THEME.typography.heading.h1}>
            {title}
          </h1>
          {subtitle && (
            <p className={cn(BIHUB_ADMIN_THEME.typography.body.large, 'mt-2')}>
              {subtitle}
            </p>
          )}
        </div>
        {action}
      </div>
    </div>
  )
}

// BiHub Empty State Component
interface BiHubEmptyStateProps {
  icon: ReactNode
  title: string
  description: string
  action?: ReactNode
}

export function BiHubEmptyState({ icon, title, description, action }: BiHubEmptyStateProps) {
  return (
    <div className="text-center py-12">
      <div className="w-16 h-16 rounded-full bg-gray-700 flex items-center justify-center mx-auto mb-4">
        {icon}
      </div>
      <h3 className={cn(BIHUB_ADMIN_THEME.typography.heading.h3, 'mb-2')}>
        {title}
      </h3>
      <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
        {description}
      </p>
      {action && (
        <div className="mt-6">
          {action}
        </div>
      )}
    </div>
  )
}
