import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { categoryService, CreateCategoryRequest, UpdateCategoryRequest } from '@/lib/services/categories'
import { Category } from '@/types'
import { toast } from 'sonner'

export const useCategories = () => {
  return useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryService.getCategories(),
  })
}

export const useCategory = (id: string) => {
  return useQuery({
    queryKey: ['categories', id],
    queryFn: () => categoryService.getCategory(id),
    enabled: !!id,
  })
}

export const useCreateCategory = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: (data: CreateCategoryRequest) => categoryService.createCategory(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      toast.success('Category created successfully!')
    },
    onError: (error: any) => {
      const message = error?.message || error?.error || 'Failed to create category'
      toast.error(message)
    },
  })
}

export const useUpdateCategory = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateCategoryRequest }) => {
      console.log('Updating category:', { id, data })
      return categoryService.updateCategory(id, data)
    },
    onSuccess: (updatedCategory) => {
      console.log('Category updated successfully:', updatedCategory)
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      // Remove toast from here, let component handle it
    },
    onError: (error: any) => {
      console.error('Update category error:', error)
      // Let the component handle the error toast for better UX
      throw error // Re-throw để component có thể handle
    },
  })
}

export const useDeleteCategory = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string) => {
      console.log('Deleting category with ID:', id)
      await categoryService.deleteCategory(id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] })
      // Let component handle success toast
    },
    onError: (error: any) => {
      console.error('Delete category error:', error)
      // Let the component handle the error toast for better UX
      throw error // Re-throw để component có thể handle
    },
  })
}
