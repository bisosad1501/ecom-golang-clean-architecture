import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { categoryService, CreateCategoryRequest, UpdateCategoryRequest } from '@/lib/services/categories'
import { Category } from '@/types'
import { toast } from 'sonner'
import { useState, useCallback, useEffect } from 'react'

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

// Category SEO Management Hook
export function useCategorySEO(categoryId: string) {
  const [categorySEO, setCategorySEO] = useState<any>(null)
  const [seoInsights, setSeoInsights] = useState<any>(null)
  const [competitorAnalysis, setCompetitorAnalysis] = useState<any>(null)
  const [slugSuggestions, setSlugSuggestions] = useState<any>(null)
  const [isLoadingSEO, setIsLoadingSEO] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch category SEO data
  const fetchCategorySEO = useCallback(async () => {
    if (!categoryId) return

    setIsLoadingSEO(true)
    setError(null)

    try {
      const response = await fetch(`/api/v1/categories/${categoryId}/seo`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to fetch category SEO')
      }

      const data = await response.json()
      setCategorySEO(data.data)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setIsLoadingSEO(false)
    }
  }, [categoryId])

  // Fetch SEO insights
  const fetchSEOInsights = useCallback(async () => {
    if (!categoryId) return

    try {
      const response = await fetch(`/api/v1/admin/categories/${categoryId}/seo/insights`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        setSeoInsights(data.data)
      }
    } catch (err) {
      console.error('Failed to fetch SEO insights:', err)
    }
  }, [categoryId])

  // Fetch competitor analysis
  const fetchCompetitorAnalysis = useCallback(async () => {
    if (!categoryId) return

    try {
      const response = await fetch(`/api/v1/admin/categories/${categoryId}/seo/competitor-analysis`, {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })

      if (response.ok) {
        const data = await response.json()
        setCompetitorAnalysis(data.data)
      }
    } catch (err) {
      console.error('Failed to fetch competitor analysis:', err)
    }
  }, [categoryId])

  // Update category SEO
  const updateCategorySEO = useCallback(async (seoData: any) => {
    if (!categoryId) return

    const response = await fetch(`/api/v1/admin/categories/${categoryId}/seo`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(seoData)
    })

    if (!response.ok) {
      throw new Error('Failed to update category SEO')
    }

    const data = await response.json()
    setCategorySEO(data.data.seo)

    // Refresh insights after update
    fetchSEOInsights()

    return data.data
  }, [categoryId, fetchSEOInsights])

  // Generate SEO metadata
  const generateSEO = useCallback(async () => {
    if (!categoryId) return

    const response = await fetch(`/api/v1/admin/categories/${categoryId}/seo/generate`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to generate SEO metadata')
    }

    const data = await response.json()
    return data.data
  }, [categoryId])

  // Validate SEO
  const validateSEO = useCallback(async () => {
    if (!categoryId) return

    const response = await fetch(`/api/v1/admin/categories/${categoryId}/seo/validate`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to validate SEO')
    }

    const data = await response.json()
    return data.data
  }, [categoryId])

  // Optimize slug
  const optimizeSlug = useCallback(async (slugData: any) => {
    if (!categoryId) return

    const response = await fetch(`/api/v1/admin/categories/${categoryId}/slug/optimize`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify(slugData)
    })

    if (!response.ok) {
      throw new Error('Failed to optimize slug')
    }

    const data = await response.json()
    return data.data
  }, [categoryId])

  // Generate slug suggestions
  const generateSlugSuggestions = useCallback(async () => {
    if (!categoryId) return

    const response = await fetch(`/api/v1/admin/categories/${categoryId}/slug/suggestions`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) {
      throw new Error('Failed to generate slug suggestions')
    }

    const data = await response.json()
    setSlugSuggestions(data.data)
    return data.data
  }, [categoryId])

  // Validate slug availability
  const validateSlugAvailability = useCallback(async (slug: string, excludeId?: string) => {
    const params = new URLSearchParams({ slug })
    if (excludeId) {
      params.append('exclude_id', excludeId)
    }

    const response = await fetch(`/api/v1/categories/slug/validate?${params}`)

    if (!response.ok) {
      throw new Error('Failed to validate slug')
    }

    const data = await response.json()
    return data.data
  }, [])

  // Initialize data
  useEffect(() => {
    if (categoryId) {
      fetchCategorySEO()
      fetchSEOInsights()
      fetchCompetitorAnalysis()
    }
  }, [categoryId, fetchCategorySEO, fetchSEOInsights, fetchCompetitorAnalysis])

  return {
    categorySEO,
    seoInsights,
    competitorAnalysis,
    slugSuggestions,
    isLoadingSEO,
    error,
    updateCategorySEO,
    generateSEO,
    validateSEO,
    optimizeSlug,
    generateSlugSuggestions,
    validateSlugAvailability,
    refetch: () => {
      fetchCategorySEO()
      fetchSEOInsights()
      fetchCompetitorAnalysis()
    }
  }
}
