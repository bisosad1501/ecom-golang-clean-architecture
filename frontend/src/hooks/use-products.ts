import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { productService, ProductsParams, CreateProductRequest, UpdateProductRequest } from '@/lib/services/products'
import { Product } from '@/types/product'
import { toast } from 'sonner'

// Query keys
export const productKeys = {
  all: ['products'] as const,
  lists: () => [...productKeys.all, 'list'] as const,
  list: (params: ProductsParams) => [...productKeys.lists(), params] as const,
  details: () => [...productKeys.all, 'detail'] as const,
  detail: (id: string) => [...productKeys.details(), id] as const,
  featured: () => [...productKeys.all, 'featured'] as const,
  related: (id: string) => [...productKeys.all, 'related', id] as const,
  search: (query: string, params: ProductsParams) => [...productKeys.all, 'search', query, params] as const,
  suggestions: (query: string) => [...productKeys.all, 'suggestions', query] as const,
  reviews: (id: string) => [...productKeys.all, 'reviews', id] as const,
  analytics: (id: string, period: string) => [...productKeys.all, 'analytics', id, period] as const,
  admin: () => [...productKeys.all, 'admin'] as const,
  adminList: (params: ProductsParams) => [...productKeys.admin(), 'list', params] as const,
}

// Hooks for fetching products
export function useProducts(params: ProductsParams = {}) {
  return useQuery({
    queryKey: productKeys.list(params),
    queryFn: async () => {
      try {
        const result = await productService.getProducts(params)
        return result
      } catch (error) {
        console.error('useProducts: Error fetching products:', error)
        throw error
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export function useAdminProducts(params: ProductsParams = {}) {
  return useQuery({
    queryKey: productKeys.adminList(params),
    queryFn: async () => {
      console.log('useAdminProducts: Fetching admin products with params:', params)
      try {
        const result = await productService.getAdminProducts(params)
        console.log('useAdminProducts: Successfully fetched admin products:', result)
        return result
      } catch (error) {
        console.error('useAdminProducts: Error fetching admin products:', error)
        throw error
      }
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: productKeys.detail(id),
    queryFn: () => productService.getProduct(id),
    enabled: !!id,
    staleTime: 10 * 60 * 1000, // 10 minutes
  })
}

export function useProductBySlug(slug: string) {
  return useQuery({
    queryKey: [...productKeys.details(), 'slug', slug],
    queryFn: () => productService.getProductBySlug(slug),
    enabled: !!slug,
    staleTime: 10 * 60 * 1000,
  })
}

export function useFeaturedProducts(limit = 8) {
  return useQuery({
    queryKey: [...productKeys.featured(), limit],
    queryFn: () => productService.getFeaturedProducts(limit),
    staleTime: 15 * 60 * 1000, // 15 minutes
  })
}

export function useRelatedProducts(productId: string, limit = 4) {
  return useQuery({
    queryKey: [...productKeys.related(productId), limit],
    queryFn: () => productService.getRelatedProducts(productId, limit),
    enabled: !!productId,
    staleTime: 10 * 60 * 1000,
  })
}

export function useProductsByCategory(categoryId: string, params: ProductsParams = {}) {
  return useQuery({
    queryKey: [...productKeys.all, 'category', categoryId, params],
    queryFn: () => productService.getProductsByCategory(categoryId, params),
    enabled: !!categoryId,
    staleTime: 5 * 60 * 1000,
  })
}

export function useSearchProducts(query: string, params: ProductsParams = {}) {
  return useQuery({
    queryKey: productKeys.search(query, params),
    queryFn: () => productService.searchProducts(query, params),
    enabled: !!query && query.length > 2,
    staleTime: 2 * 60 * 1000, // 2 minutes for search results
  })
}

export function useProductSuggestions(query: string, limit = 5) {
  return useQuery({
    queryKey: [...productKeys.suggestions(query), limit],
    queryFn: () => productService.getProductSuggestions(query, limit),
    enabled: !!query && query.length > 1,
    staleTime: 1 * 60 * 1000, // 1 minute for suggestions
  })
}

export function useProductReviews(productId: string, params: { page?: number; limit?: number } = {}) {
  return useQuery({
    queryKey: [...productKeys.reviews(productId), params],
    queryFn: () => productService.getProductReviews(productId, params),
    enabled: !!productId,
    staleTime: 5 * 60 * 1000,
  })
}

export function useProductAnalytics(productId: string, period = '30d') {
  return useQuery({
    queryKey: productKeys.analytics(productId, period),
    queryFn: () => productService.getProductAnalytics(productId, period),
    enabled: !!productId,
    staleTime: 10 * 60 * 1000,
  })
}

// Mutations for modifying products
export function useCreateProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateProductRequest) => productService.createProduct(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      toast.success('Product created successfully!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create product')
    },
  })
}

export function useUpdateProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProductRequest }) => {
      console.log('=== useUpdateProduct mutationFn CALLED ===')
      console.log('useUpdateProduct mutation:', { id, data })
      console.log('Timestamp:', new Date().toISOString())
      return productService.updateProduct(id, data)
    },
    onSuccess: (updatedProduct, { id }) => {
      console.log('=== useUpdateProduct onSuccess CALLED ===')
      console.log('useUpdateProduct success:', { updatedProduct, id })
      console.log('Timestamp:', new Date().toISOString())
      
      // Only invalidate queries, don't set cache directly to avoid conflicts
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      queryClient.invalidateQueries({ queryKey: productKeys.detail(id) })
      
      // Comment out direct cache update to prevent conflicts
      // queryClient.setQueryData(productKeys.detail(id), updatedProduct)
    },
    onError: (error: any) => {
      console.error('useUpdateProduct error:', error)
      toast.error(error.message || 'Failed to update product')
    },
  })
}

export function useDeleteProduct() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => productService.deleteProduct(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      // Don't show toast here, let the component handle it
    },
    onError: (error: any) => {
      console.error('Delete product error:', error)
      // Don't show toast here, let the component handle it
    },
  })
}

export function useBulkDeleteProducts() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (productIds: string[]) => productService.bulkDeleteProducts(productIds),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: productKeys.lists() })
      toast.success('Products deleted successfully!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to delete products')
    },
  })
}

export function useUploadProductImage() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ productId, file }: { productId: string; file: File }) =>
      productService.uploadProductImage(productId, file),
    onSuccess: (_, { productId }) => {
      // Invalidate product details to refetch with new image
      queryClient.invalidateQueries({ queryKey: productKeys.detail(productId) })
      toast.success('Image uploaded successfully!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to upload image')
    },
  })
}

export function useAddToWishlist() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (productId: string) => productService.addToWishlist(productId),
    onSuccess: () => {
      // Invalidate wishlist queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Added to wishlist!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to add to wishlist')
    },
  })
}

export function useRemoveFromWishlist() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (productId: string) => productService.removeFromWishlist(productId),
    onSuccess: () => {
      // Invalidate wishlist queries
      queryClient.invalidateQueries({ queryKey: ['wishlist'] })
      toast.success('Removed from wishlist!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to remove from wishlist')
    },
  })
}

export function useCreateReview() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ productId, data }: { productId: string; data: { rating: number; title?: string; comment?: string } }) =>
      productService.createReview(productId, data),
    onSuccess: (_, { productId }) => {
      // Invalidate product reviews and product details
      queryClient.invalidateQueries({ queryKey: productKeys.reviews(productId) })
      queryClient.invalidateQueries({ queryKey: productKeys.detail(productId) })
      toast.success('Review submitted successfully!')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to submit review')
    },
  })
}
