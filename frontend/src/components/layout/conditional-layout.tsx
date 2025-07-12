'use client'

import { usePathname } from 'next/navigation'
import { Header } from '@/components/layout/header'
import { Footer } from '@/components/layout/footer'
import { GlobalCartConflictModal } from '@/components/cart/global-cart-conflict-modal'

interface ConditionalLayoutProps {
  children: React.ReactNode
}

export function ConditionalLayout({ children }: ConditionalLayoutProps) {
  const pathname = usePathname()
  
  // Don't show header/footer for admin pages
  const isAdminPage = pathname.startsWith('/admin')
  
  if (isAdminPage) {
    return <>{children}</>
  }
  
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main className="flex-1">
        {children}
      </main>
      <Footer />

      {/* Global Cart Conflict Modal */}
      <GlobalCartConflictModal />
    </div>
  )
}
