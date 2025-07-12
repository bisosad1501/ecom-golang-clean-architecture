'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/store/auth'
import { useCartStore } from '@/store/cart'
import { CartConflictModal } from './cart-conflict-modal'
import { MergeStrategy } from '@/types'
import { toast } from 'sonner'

export function GlobalCartConflictModal() {
  const { pendingCartConflict, clearPendingCartConflict } = useAuthStore()
  const { mergeGuestCart, fetchCart } = useCartStore()
  const [isLoading, setIsLoading] = useState(false)

  const handleMerge = async (strategy: MergeStrategy) => {
    if (!pendingCartConflict) return

    setIsLoading(true)
    try {
      // Merge with selected strategy (mergeGuestCart will get session ID from localStorage)
      await mergeGuestCart(strategy)
      
      // Fetch fresh cart after merge
      await fetchCart()
      
      // Clear conflict state
      clearPendingCartConflict()
      
      toast.success(`Cart merged successfully using ${strategy} strategy!`)
    } catch (error: any) {
      console.error('Failed to merge cart:', error)
      toast.error('Failed to merge cart. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleClose = () => {
    // If user closes without merging, use auto strategy as fallback
    if (pendingCartConflict) {
      handleMerge('auto')
    }
  }

  if (!pendingCartConflict) {
    return null
  }

  return (
    <CartConflictModal
      isOpen={!!pendingCartConflict}
      onClose={handleClose}
      conflictInfo={pendingCartConflict}
      onMerge={handleMerge}
      isLoading={isLoading}
    />
  )
}
