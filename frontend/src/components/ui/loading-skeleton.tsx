import { cn } from '@/lib/utils'

interface LoadingSkeletonProps {
  className?: string
  variant?: 'default' | 'card' | 'text' | 'avatar' | 'button'
  lines?: number
}

export function LoadingSkeleton({ 
  className, 
  variant = 'default',
  lines = 1 
}: LoadingSkeletonProps) {
  const baseClasses = "skeleton rounded-lg bg-gray-700"
  
  const variants = {
    default: "h-4 w-full",
    card: "h-32 w-full",
    text: "h-4",
    avatar: "h-12 w-12 rounded-full",
    button: "h-10 w-24"
  }

  if (variant === 'text' && lines > 1) {
    return (
      <div className="space-y-2">
        {Array.from({ length: lines }).map((_, i) => (
          <div
            key={i}
            className={cn(
              baseClasses,
              variants.text,
              i === lines - 1 ? "w-3/4" : "w-full",
              className
            )}
          />
        ))}
      </div>
    )
  }

  return (
    <div
      className={cn(
        baseClasses,
        variants[variant],
        className
      )}
    />
  )
}

// Specialized loading components
export function ProductCardSkeleton() {
  return (
    <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-xl p-6 space-y-4">
      <LoadingSkeleton variant="card" />
      <LoadingSkeleton variant="text" lines={2} />
      <div className="flex justify-between items-center">
        <LoadingSkeleton className="h-6 w-20" />
        <LoadingSkeleton variant="button" />
      </div>
    </div>
  )
}

export function OrderCardSkeleton() {
  return (
    <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-xl p-6 space-y-4">
      <div className="flex items-center gap-4">
        <LoadingSkeleton variant="avatar" />
        <div className="flex-1 space-y-2">
          <LoadingSkeleton className="h-6 w-32" />
          <LoadingSkeleton className="h-4 w-24" />
        </div>
        <LoadingSkeleton className="h-6 w-16" />
      </div>
      <div className="grid grid-cols-3 gap-4">
        <LoadingSkeleton className="h-4 w-full" />
        <LoadingSkeleton className="h-4 w-full" />
        <LoadingSkeleton className="h-4 w-full" />
      </div>
    </div>
  )
}

export function CartItemSkeleton() {
  return (
    <div className="bg-gray-900/50 backdrop-blur-sm border border-gray-700 rounded-xl p-6">
      <div className="flex flex-col sm:flex-row gap-6">
        <LoadingSkeleton className="h-36 w-36 rounded-2xl" />
        <div className="flex-1 space-y-4">
          <LoadingSkeleton variant="text" lines={3} />
          <div className="flex justify-between items-center">
            <LoadingSkeleton className="h-8 w-24" />
            <LoadingSkeleton className="h-10 w-32" />
          </div>
        </div>
      </div>
    </div>
  )
}

export function PageLoadingSkeleton() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      <div className="container mx-auto px-4 py-12">
        <div className="space-y-8">
          {/* Header skeleton */}
          <div className="space-y-4">
            <LoadingSkeleton className="h-8 w-64" />
            <LoadingSkeleton className="h-4 w-96" />
          </div>
          
          {/* Content skeleton */}
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-6">
              {Array.from({ length: 3 }).map((_, i) => (
                <LoadingSkeleton key={i} variant="card" />
              ))}
            </div>
            <div className="space-y-6">
              {Array.from({ length: 2 }).map((_, i) => (
                <LoadingSkeleton key={i} className="h-48 w-full" />
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
