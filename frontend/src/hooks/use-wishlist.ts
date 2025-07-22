// ===== WISHLIST HOOKS =====

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { wishlistService } from '@/lib/services/wishlist'
import { useAuthStore } from '@/store/auth'
import { toast } from 'sonner'
import type { GetWishlistRequest } from '@/types/wishlist'

// Get user's wishlist
export function useWishlist(params?: GetWishlistRequest) {
  const { isAuthenticated, isHydrated, user } = useAuthStore()

  const isEnabled = isAuthenticated && isHydrated && !!user

  const query = useQuery({
    queryKey: ['wishlist', user?.id, params], // Include user ID in query key
    queryFn: async () => {
      try {
        const result = await wishlistService.getWishlist(params)
        return result
      } catch (error) {
        console.error('useWishlist queryFn error:', error)
        throw error
      }
    },
    enabled: isEnabled,
    staleTime: 5 * 60 * 1000, // 5 minutes
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    retry: (failureCount, error: any) => {
      // Don't retry on 401/403 errors
      if (error?.status === 401 || error?.status === 403) {
        return false
      }
      return failureCount < 2
    },
  })

  return query
}

// Get wishlist count
export function useWishlistCount() {
  const { isAuthenticated } = useAuthStore()

  return useQuery({
    queryKey: ['wishlist', 'count'],
    queryFn: () => wishlistService.getWishlistCount(),
    enabled: isAuthenticated,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Check if product is in wishlist
export function useIsInWishlist(productId: string) {
  const { isAuthenticated } = useAuthStore()

  return useQuery({
    queryKey: ['wishlist', 'status', productId],
    queryFn: () => wishlistService.isInWishlist(productId),
    enabled: isAuthenticated && !!productId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

// Add to wishlist mutation
export function useAddToWishlist() {
  const queryClient = useQueryClient()
  const { isAuthenticated } = useAuthStore()

  return useMutation({
    mutationFn: (productId: string) => {
      if (!isAuthenticated) {
        throw new Error('Please sign in to add items to your wishlist')
      }
      return wishlistService.addToWishlist(productId)
    },
    onSuccess: (_, productId) => {
      // Invalidate wishlist queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      // Update specific product status
      queryClient.setQueryData(['wishlist', 'status', productId], true)
      toast.success('Added to wishlist!')
    },
    onError: (error: any) => {
      if (error.message?.includes('already in wishlist')) {
        toast.info('Product is already in your wishlist')
      } else {
        toast.error(error.message || 'Failed to add to wishlist')
      }
    },
  })
}

// Remove from wishlist mutation
export function useRemoveFromWishlist() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (productId: string) => wishlistService.removeFromWishlist(productId),
    onSuccess: (_, productId) => {
      // Invalidate wishlist queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      // Update specific product status
      queryClient.setQueryData(['wishlist', 'status', productId], false)
      toast.success('Removed from wishlist!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to remove from wishlist')
    },
  })
}

// Clear wishlist mutation
export function useClearWishlist() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => wishlistService.clearWishlist(),
    onSuccess: () => {
      // Invalidate all wishlist queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Wishlist cleared!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to clear wishlist')
    },
  })
}

// Move to cart mutation
export function useMoveToCart() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (productId: string) => wishlistService.moveToCart(productId),
    onSuccess: (_, productId) => {
      // Invalidate wishlist and cart queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      queryClient.invalidateQueries({ queryKey: ['cart'] })
      // Update specific product status
      queryClient.setQueryData(['wishlist', 'status', productId], false)
      toast.success('Moved to cart!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to move to cart')
    },
  })
}
