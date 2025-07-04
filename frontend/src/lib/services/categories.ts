import { apiClient } from '@/lib/api'
import { Category, ApiResponse } from '@/types'

export interface CreateCategoryRequest {
  name: string
  slug?: string
  description?: string
  parent_id?: string
  image?: string
  is_active?: boolean
  sort_order?: number
}

export interface UpdateCategoryRequest extends Partial<CreateCategoryRequest> {}

class CategoryService {
  // Get all categories
  async getCategories(): Promise<Category[]> {
    const response = await apiClient.get('/categories')
    return response.data
  }

  // Get root categories (no parent)
  async getRootCategories(): Promise<Category[]> {
    const response = await apiClient.get<Category[]>('/categories/root')
    return response.data || []
  }

  // Get category tree (hierarchical structure)
  async getCategoryTree(): Promise<Category[]> {
    const response = await apiClient.get<Category[]>('/categories/tree')
    return response.data || []
  }

  // Get single category by ID
  async getCategory(id: string): Promise<Category> {
    const response = await apiClient.get(`/categories/${id}`)
    return response.data
  }

  // Get category by slug
  async getCategoryBySlug(slug: string): Promise<Category> {
    const response = await apiClient.get(`/categories/slug/${slug}`)
    return response.data
  }

  // Get category path (breadcrumbs)
  async getCategoryPath(categoryId: string): Promise<Category[]> {
    const response = await apiClient.get<Category[]>(`/categories/${categoryId}/path`)
    return response.data || []
  }

  // Get category children
  async getCategoryChildren(categoryId: string): Promise<Category[]> {
    const response = await apiClient.get<Category[]>(`/categories/${categoryId}/children`)
    return response.data || []
  }

  // Get category product count (including subcategories)
  async getCategoryProductCount(categoryId: string): Promise<number> {
    const response = await apiClient.get<{ product_count: number }>(`/categories/${categoryId}/product-count`)
    return response.data?.product_count || 0
  }

  // Get category breadcrumb
  async getCategoryBreadcrumb(id: string): Promise<Category[]> {
    const response = await apiClient.get(`/categories/${id}/breadcrumb`)
    return response.data
  }

  // Admin methods
  async createCategory(data: CreateCategoryRequest): Promise<Category> {
    console.log('CategoryService: Creating category with data:', data)
    const response = await apiClient.post('/admin/categories', data)
    console.log('CategoryService: Create response:', response.data)
    return response.data
  }

  async updateCategory(id: string, data: UpdateCategoryRequest): Promise<Category> {
    console.log('CategoryService: Original data:', data)
    console.log('CategoryService: Sending to backend:', data)
    
    const response = await apiClient.put(`/admin/categories/${id}`, data)
    console.log('CategoryService: Backend response:', response.data)
    return response.data
  }

  async deleteCategory(id: string): Promise<void> {
    console.log('CategoryService: Deleting category with ID:', id)
    try {
      const response = await apiClient.delete(`/admin/categories/${id}`)
      console.log('CategoryService: Delete response:', response)
    } catch (error) {
      console.error('CategoryService: Delete error:', error)
      throw error
    }
  }

  async uploadCategoryImage(categoryId: string, file: File): Promise<{ url: string }> {
    const response = await apiClient.upload(`/admin/categories/${categoryId}/image`, file)
    return response.data
  }

  async reorderCategories(categoryIds: string[]): Promise<void> {
    await apiClient.patch('/admin/categories/reorder', { category_ids: categoryIds })
  }

  // Bulk operations
  async bulkUpdateCategories(categoryIds: string[], data: UpdateCategoryRequest): Promise<void> {
    await apiClient.patch('/admin/categories/bulk', { category_ids: categoryIds, ...data })
  }

  async bulkDeleteCategories(categoryIds: string[]): Promise<void> {
    await apiClient.delete('/admin/categories/bulk', { data: { category_ids: categoryIds } })
  }

  // Category analytics
  async getCategoryAnalytics(categoryId: string, period = '30d'): Promise<any> {
    const response = await apiClient.get(`/admin/categories/${categoryId}/analytics?period=${period}`)
    return response.data
  }

  // Popular categories
  async getPopularCategories(limit = 10): Promise<Category[]> {
    const response = await apiClient.get(`/categories/popular?limit=${limit}`)
    return response.data
  }

  // Search categories
  async searchCategories(query: string): Promise<Category[]> {
    const response = await apiClient.get(`/categories/search?q=${query}`)
    return response.data
  }
}

export const categoryService = new CategoryService()
