'use client'

import React from 'react'
import { cn } from '@/lib/utils'

interface ResponsiveGridProps {
  children: React.ReactNode
  variant?: 'categories' | 'products' | 'compact' | 'wide'
  className?: string
  gap?: 'sm' | 'md' | 'lg' | 'xl'
}

const gridVariants = {
  categories: {
    base: 'grid-cols-1',
    sm: 'sm:grid-cols-2',
    md: 'md:grid-cols-3',
    lg: 'lg:grid-cols-4',
    xl: 'xl:grid-cols-4',
    '2xl': '2xl:grid-cols-4'
  },
  products: {
    base: 'grid-cols-1',
    sm: 'sm:grid-cols-2',
    md: 'md:grid-cols-3',
    lg: 'lg:grid-cols-4',
    xl: 'xl:grid-cols-4',
    '2xl': '2xl:grid-cols-4'
  },
  compact: {
    base: 'grid-cols-2',
    sm: 'sm:grid-cols-3',
    md: 'md:grid-cols-4',
    lg: 'lg:grid-cols-5',
    xl: 'xl:grid-cols-6',
    '2xl': '2xl:grid-cols-8'
  },
  wide: {
    base: 'grid-cols-1',
    sm: 'sm:grid-cols-2',
    md: 'md:grid-cols-2',
    lg: 'lg:grid-cols-3',
    xl: 'xl:grid-cols-4',
    '2xl': '2xl:grid-cols-4'
  }
}

const gapSizes = {
  sm: 'gap-4',
  md: 'gap-6',
  lg: 'gap-8',
  xl: 'gap-10'
}

export function ResponsiveGrid({
  children,
  variant = 'categories',
  className,
  gap = 'lg'
}: ResponsiveGridProps) {
  const gridClasses = gridVariants[variant]
  
  return (
    <div className={cn(
      'grid',
      gridClasses.base,
      gridClasses.sm,
      gridClasses.md,
      gridClasses.lg,
      gridClasses.xl,
      gridClasses['2xl'],
      gapSizes[gap],
      className
    )}>
      {children}
    </div>
  )
}

// Grid item with animation support
interface GridItemProps {
  children: React.ReactNode
  index?: number
  className?: string
  animationDelay?: number
}

export function GridItem({
  children,
  index = 0,
  className,
  animationDelay = 50
}: GridItemProps) {
  return (
    <div
      className={cn(
        'animate-in fade-in slide-in-from-bottom-4 duration-500 h-full',
        className
      )}
      style={{
        animationDelay: `${index * animationDelay}ms`,
        animationFillMode: 'both'
      }}
    >
      {children}
    </div>
  )
}

// Masonry grid for varied content heights
interface MasonryGridProps {
  children: React.ReactNode
  columns?: {
    default: number
    sm?: number
    md?: number
    lg?: number
    xl?: number
  }
  gap?: string
  className?: string
}

export function MasonryGrid({
  children,
  columns = { default: 1, sm: 2, md: 3, lg: 4 },
  gap = '1.5rem',
  className
}: MasonryGridProps) {
  const columnClasses = [
    `columns-${columns.default}`,
    columns.sm && `sm:columns-${columns.sm}`,
    columns.md && `md:columns-${columns.md}`,
    columns.lg && `lg:columns-${columns.lg}`,
    columns.xl && `xl:columns-${columns.xl}`
  ].filter(Boolean).join(' ')

  return (
    <div 
      className={cn(columnClasses, className)}
      style={{ columnGap: gap }}
    >
      {React.Children.map(children, (child, index) => (
        <div 
          className="break-inside-avoid mb-6"
          style={{ 
            animationDelay: `${index * 100}ms`,
            animationFillMode: 'both'
          }}
        >
          {child}
        </div>
      ))}
    </div>
  )
}

// Auto-fit grid that adjusts based on content
interface AutoFitGridProps {
  children: React.ReactNode
  minItemWidth?: string
  gap?: string
  className?: string
}

export function AutoFitGrid({
  children,
  minItemWidth = '280px',
  gap = '1.5rem',
  className
}: AutoFitGridProps) {
  return (
    <div 
      className={cn('grid', className)}
      style={{
        gridTemplateColumns: `repeat(auto-fit, minmax(${minItemWidth}, 1fr))`,
        gap
      }}
    >
      {children}
    </div>
  )
}

// Staggered animation container
interface StaggeredContainerProps {
  children: React.ReactNode
  staggerDelay?: number
  className?: string
}

export function StaggeredContainer({
  children,
  staggerDelay = 100,
  className
}: StaggeredContainerProps) {
  return (
    <div className={className}>
      {React.Children.map(children, (child, index) => (
        <div
          key={index}
          className="animate-in fade-in slide-in-from-bottom-4 duration-500"
          style={{ 
            animationDelay: `${index * staggerDelay}ms`,
            animationFillMode: 'both'
          }}
        >
          {child}
        </div>
      ))}
    </div>
  )
}

// Grid skeleton loader
interface GridSkeletonProps {
  count?: number
  variant?: 'categories' | 'products' | 'compact'
  gap?: 'sm' | 'md' | 'lg' | 'xl'
  className?: string
}

export function GridSkeleton({
  count = 8,
  variant = 'categories',
  gap = 'lg',
  className
}: GridSkeletonProps) {
  return (
    <ResponsiveGrid variant={variant} gap={gap} className={className}>
      {Array.from({ length: count }, (_, i) => (
        <GridItem key={i} index={i}>
          <div className="bg-white/[0.08] backdrop-blur-sm border border-white/15 rounded-xl overflow-hidden">
            <div className="aspect-[4/3] bg-gradient-to-br from-gray-700/50 to-gray-800/50 animate-pulse" />
            <div className="p-4 space-y-3">
              <div className="h-5 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded animate-pulse" />
              <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-3/4 animate-pulse" />
              <div className="flex justify-between items-center">
                <div className="h-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded w-1/3 animate-pulse" />
                <div className="h-4 w-4 bg-gradient-to-r from-gray-600/50 to-gray-700/50 rounded animate-pulse" />
              </div>
            </div>
          </div>
        </GridItem>
      ))}
    </ResponsiveGrid>
  )
}

// Responsive container with max width constraints
interface ResponsiveContainerProps {
  children: React.ReactNode
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  className?: string
}

const containerSizes = {
  sm: 'max-w-3xl',
  md: 'max-w-5xl',
  lg: 'max-w-7xl',
  xl: 'max-w-[1400px]',
  full: 'max-w-none'
}

export function ResponsiveContainer({
  children,
  size = 'xl',
  className
}: ResponsiveContainerProps) {
  return (
    <div className={cn(
      'container mx-auto px-4 sm:px-6 lg:px-8',
      containerSizes[size],
      className
    )}>
      {children}
    </div>
  )
}
