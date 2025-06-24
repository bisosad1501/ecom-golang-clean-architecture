'use client'

import { useState } from 'react'
import { Plus, Search, Filter, MoreHorizontal, Edit, Trash2, Eye, FolderTree, ChevronDown, ChevronRight, Settings, Expand, Minimize } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { Category } from '@/types'
import { toast } from 'sonner'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { AddCategoryForm } from '@/components/forms/add-category-form'
import { EditCategoryForm } from '@/components/forms/edit-category-form'
import { SubcategoryManager } from '@/components/admin/subcategory-manager'
import { CategoryStatusToggle } from '@/components/admin/category-status-toggle'
import { ErrorDialog } from '@/components/ui/error-dialog'
import { useCategories, useDeleteCategory } from '@/hooks/use-categories'
import { getCategoryDeleteErrorMessage, CategoryDeleteError } from '@/lib/utils/category-delete-errors'

export function AdminCategoriesPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [showAddForm, setShowAddForm] = useState(false)
  const [deleteCategoryId, setDeleteCategoryId] = useState<string | null>(null)
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)
  const [showCategoryModal, setShowCategoryModal] = useState(false)
  const [editCategory, setEditCategory] = useState<Category | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [addSubcategoryParent, setAddSubcategoryParent] = useState<Category | null>(null)
  const [showAddSubcategoryModal, setShowAddSubcategoryModal] = useState(false)
  const [manageSubcategoriesParent, setManageSubcategoriesParent] = useState<Category | null>(null)
  const [showManageSubcategoriesModal, setShowManageSubcategoriesModal] = useState(false)
  const [collapsedCategories, setCollapsedCategories] = useState<Set<string>>(new Set())
  const [deleteError, setDeleteError] = useState<CategoryDeleteError | null>(null)
  const [showDeleteErrorModal, setShowDeleteErrorModal] = useState(false)
  
  const { data: categories, isLoading, refetch } = useCategories()
  const deleteCategory = useDeleteCategory()

  // Filter categories based on search
  const filteredCategories = categories?.filter(category =>
    category.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    category.description?.toLowerCase().includes(searchQuery.toLowerCase())
  ) || []

  // Organize categories into tree structure
  const parentCategories = filteredCategories.filter(cat => !cat.parent_id)
  const childCategories = filteredCategories.filter(cat => cat.parent_id)

  // Action handlers
  const handleViewCategory = (category: Category) => {
    setSelectedCategory(category)
    setShowCategoryModal(true)
  }

  const handleEditCategory = (category: Category) => {
    setEditCategory(category)
    setShowEditModal(true)
  }

  const handleDeleteCategory = (categoryId: string) => {
    setDeleteCategoryId(categoryId)
  }

  const handleAddSubcategory = (parent: Category) => {
    setAddSubcategoryParent(parent)
    setShowAddSubcategoryModal(true)
  }

  const handleManageSubcategories = (parent: Category) => {
    setManageSubcategoriesParent(parent)
    setShowManageSubcategoriesModal(true)
  }

  const toggleCategoryCollapse = (categoryId: string) => {
    const newCollapsed = new Set(collapsedCategories)
    if (newCollapsed.has(categoryId)) {
      newCollapsed.delete(categoryId)
    } else {
      newCollapsed.add(categoryId)
    }
    setCollapsedCategories(newCollapsed)
  }

  const expandAll = () => {
    setCollapsedCategories(new Set())
  }

  const collapseAll = () => {
    const allParentIds = parentCategories.map(cat => cat.id)
    setCollapsedCategories(new Set(allParentIds))
  }

  const confirmDeleteCategory = async () => {
    if (!deleteCategoryId) return
    
    try {
      // Check if category has subcategories
      const hasSubcategories = childCategories.some(child => child.parent_id === deleteCategoryId)
      
      if (hasSubcategories) {
        const deleteError: CategoryDeleteError = {
          title: 'Cannot Delete Parent Category',
          message: 'This category cannot be deleted because it has subcategories.',
          suggestions: [
            'Delete all subcategories first',
            'Move subcategories to another parent category',
            'Or disable the category instead of deleting it'
          ]
        }
        setDeleteError(deleteError)
        setShowDeleteErrorModal(true)
        setDeleteCategoryId(null)
        return
      }

      await deleteCategory.mutateAsync(deleteCategoryId)
      setDeleteCategoryId(null)
      refetch()
      toast.success('Category deleted successfully!')
    } catch (error: any) {
      console.error('Failed to delete category:', error)
      
      // Parse the error and show detailed error dialog
      const deleteError = getCategoryDeleteErrorMessage(error)
      setDeleteError(deleteError)
      setShowDeleteErrorModal(true)
      setDeleteCategoryId(null) // Close the delete confirmation dialog
    }
  }

  const getCategoryLevel = (category: Category) => {
    return category.level || 0
  }

  const getChildrenCount = (parentId: string) => {
    return childCategories.filter(cat => cat.parent_id === parentId).length
  }

  if (showAddForm) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Add New Category</h1>
            <p className="text-gray-600 mt-2">Create a new category for your products</p>
          </div>
          
          <Button 
            variant="outline" 
            onClick={() => setShowAddForm(false)}
          >
            Back to Categories
          </Button>
        </div>

        <AddCategoryForm 
          onSuccess={() => {
            setShowAddForm(false)
            refetch()
          }}
          onCancel={() => setShowAddForm(false)}
        />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Categories</h1>
          <p className="text-gray-600 mt-2">Manage your product categories</p>
        </div>
        
        <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
          <Button onClick={() => setShowAddForm(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Category
          </Button>
        </RequirePermission>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <FolderTree className="h-8 w-8 text-blue-600" />
              <div>
                <p className="text-sm font-medium text-gray-600">Total Categories</p>
                <p className="text-2xl font-bold text-gray-900">{filteredCategories.length}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <div className="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
                <span className="text-blue-600 font-bold">P</span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-600">Parent Categories</p>
                <p className="text-2xl font-bold text-gray-900">{parentCategories.length}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <div className="h-8 w-8 bg-gray-100 rounded-lg flex items-center justify-center">
                <span className="text-gray-600 font-bold">S</span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-600">Subcategories</p>
                <p className="text-2xl font-bold text-gray-900">{childCategories.length}</p>
              </div>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center space-x-2">
              <div className="h-8 w-8 bg-green-100 rounded-lg flex items-center justify-center">
                <span className="text-green-600 font-bold">✓</span>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-600">Active Categories</p>
                <p className="text-2xl font-bold text-gray-900">
                  {filteredCategories.filter(cat => cat.is_active).length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4 flex-1">
              <div className="flex-1">
                <Input
                  placeholder="Search categories..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  leftIcon={<Search className="h-4 w-4" />}
                />
              </div>
              <Button variant="outline">
                <Filter className="mr-2 h-4 w-4" />
                Filters
              </Button>
            </div>
            
            {/* Expand/Collapse Controls */}
            <div className="flex items-center space-x-2 ml-4">
              <Button variant="outline" size="sm" onClick={expandAll}>
                <Expand className="mr-2 h-4 w-4" />
                Expand All
              </Button>
              <Button variant="outline" size="sm" onClick={collapseAll}>
                <Minimize className="mr-2 h-4 w-4" />
                Collapse All
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Categories Table */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FolderTree className="h-5 w-5" />
            All Categories
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-6">
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="animate-pulse flex items-center space-x-4 p-4">
                    <div className="w-12 h-12 bg-gray-200 rounded"></div>
                    <div className="flex-1 space-y-2">
                      <div className="h-4 bg-gray-200 rounded w-1/4"></div>
                      <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                    </div>
                    <div className="w-20 h-4 bg-gray-200 rounded"></div>
                  </div>
                ))}
              </div>
            </div>
          ) : filteredCategories.length === 0 ? (
            <div className="p-6 text-center">
              <FolderTree className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-sm font-medium text-gray-900">No categories</h3>
              <p className="mt-1 text-sm text-gray-500">
                {searchQuery ? 'No categories match your search.' : 'Get started by creating a new category.'}
              </p>
              {!searchQuery && (
                <div className="mt-6">
                  <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                    <Button onClick={() => setShowAddForm(true)}>
                      <Plus className="mr-2 h-4 w-4" />
                      Add Category
                    </Button>
                  </RequirePermission>
                </div>
              )}
            </div>
          ) : (
            <div className="divide-y divide-gray-200">
              {/* Parent Categories */}
              {parentCategories.map((category) => {
                const children = childCategories.filter(cat => cat.parent_id === category.id)
                const isCollapsed = collapsedCategories.has(category.id)
                
                return (
                  <div key={category.id}>
                    {/* Parent Category Row */}
                    <div className="p-4 hover:bg-gray-50">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-4 flex-1">
                          <div className="flex items-center space-x-3">
                            {/* Collapse/Expand Button */}
                            {children.length > 0 && (
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={() => toggleCategoryCollapse(category.id)}
                                className="w-6 h-6 p-0"
                              >
                                {isCollapsed ? (
                                  <ChevronRight className="h-4 w-4" />
                                ) : (
                                  <ChevronDown className="h-4 w-4" />
                                )}
                              </Button>
                            )}
                            {children.length === 0 && <div className="w-6" />}
                            
                            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center overflow-hidden">
                              {category.image ? (
                                <img 
                                  src={category.image} 
                                  alt={category.name}
                                  className="w-full h-full object-cover"
                                  onError={(e) => {
                                    const target = e.target as HTMLImageElement
                                    target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjQwIiBoZWlnaHQ9IjQwIiByeD0iOCIgZmlsbD0idXJsKCNncmFkaWVudCkiLz4KPHN2ZyB4PSIxMiIgeT0iMTIiIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiPgo8cGF0aCBkPSJNMjIgMTlhMiAyIDAgMCAxLTIgMkg0YTIgMiAwIDAgMS0yLTJWNWEyIDIgMCAwIDEgMi0yaDE2YTIgMiAwIDAgMSAyIDJ2MTR6Ii8+CjxwYXRoIGQ9Im0yMiAxMyAyLTIiLz4KPHN2ZyB4PSI5IiB5PSI5IiB3aWR0aD0iNiIgaGVpZ2h0PSI2IiB2aWV3Qm94PSIwIDAgMjQgMjQiIGZpbGw9Im5vbmUiIHN0cm9rZT0id2hpdGUiIHN0cm9rZS13aWR0aD0iMiI+CjxjaXJjbGUgY3g9IjkiIGN5PSI5IiByPSIyIi8+Cjwvc3ZnPgo8L3N2Zz4KPHA+CjxkZWZzPgo8bGluZWFyR3JhZGllbnQgaWQ9ImdyYWRpZW50IiB4MT0iMCUiIHkxPSIwJSIgeDI9IjEwMCUiIHkyPSIxMDAlIj4KPHN0b3Agb2Zmc2V0PSIwJSIgc3R5bGU9InN0b3AtY29sb3I6IzM5ODNmNjtzdG9wLW9wYWNpdHk6MSIgLz4KPHN0b3Agb2Zmc2V0PSIxMDAlIiBzdHlsZT0ic3RvcC1jb2xvcjojOWMzNGZiO3N0b3Atb3BhY2l0eToxIiAvPgo8L2xpbmVhckdyYWRpZW50Pgo8L2RlZnM+Cjwvc3ZnPgo='
                                  }}
                                />
                              ) : (
                                <FolderTree className="h-5 w-5 text-white" />
                              )}
                            </div>
                            <div>
                              <div className="flex items-center space-x-2">
                                <h3 className="text-sm font-medium text-gray-900">
                                  {category.name}
                                </h3>
                                <Badge variant="secondary">
                                  Parent
                                </Badge>
                                {children.length > 0 && (
                                  <Badge variant="outline">
                                    {children.length} subcategories
                                  </Badge>
                                )}
                              </div>
                              <p className="text-sm text-gray-500 mt-1">
                                {category.description || 'No description'}
                              </p>
                              <div className="flex items-center space-x-4 mt-2 text-xs text-gray-400">
                                <span>ID: {category.id}</span>
                                <span>Level: {getCategoryLevel(category)}</span>
                                {category.created_at && (
                                  <span>Created: {new Date(category.created_at).toLocaleDateString()}</span>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>

                        <div className="flex items-center space-x-2">
                          <Badge 
                            variant={category.is_active ? "default" : "secondary"}
                          >
                            {category.is_active ? 'Active' : 'Inactive'}
                          </Badge>

                          {/* Add Subcategory Button */}
                          <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() => handleAddSubcategory(category)}
                              className="text-xs"
                            >
                              <Plus className="mr-1 h-3 w-3" />
                              Add Sub
                            </Button>
                          </RequirePermission>

                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" className="h-8 w-8 p-0">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end">
                              <DropdownMenuItem onClick={() => handleViewCategory(category)}>
                                <Eye className="mr-2 h-4 w-4" />
                                View Details
                              </DropdownMenuItem>
                              <RequirePermission permission={PERMISSIONS.CATEGORIES_UPDATE}>
                                <DropdownMenuItem onClick={() => handleEditCategory(category)}>
                                  <Edit className="mr-2 h-4 w-4" />
                                  Edit
                                </DropdownMenuItem>
                              </RequirePermission>
                              <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                                <DropdownMenuItem onClick={() => handleAddSubcategory(category)}>
                                  <Plus className="mr-2 h-4 w-4" />
                                  Add Subcategory
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => handleManageSubcategories(category)}>
                                  <Settings className="mr-2 h-4 w-4" />
                                  Manage Subcategories
                                </DropdownMenuItem>
                              </RequirePermission>
                              <DropdownMenuSeparator />
                              <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                                <DropdownMenuItem 
                                  onClick={() => handleDeleteCategory(category.id)}
                                  className="text-red-600 focus:text-red-600"
                                >
                                  <Trash2 className="mr-2 h-4 w-4" />
                                  Delete
                                </DropdownMenuItem>
                              </RequirePermission>
                            </DropdownMenuContent>
                          </DropdownMenu>
                        </div>
                      </div>
                    </div>

                    {/* Child Categories */}
                    {!isCollapsed && children.map((child) => (
                      <div key={child.id} className="pl-8 pr-4 py-3 bg-gray-50 border-l-2 border-gray-200 hover:bg-gray-100">
                        <div className="flex items-center justify-between">
                          <div className="flex items-center space-x-4 flex-1">
                            <div className="flex items-center space-x-3">
                              <div className="w-8 h-8 bg-gradient-to-br from-gray-400 to-gray-600 rounded-md flex items-center justify-center overflow-hidden">
                                {child.image ? (
                                  <img 
                                    src={child.image} 
                                    alt={child.name}
                                    className="w-full h-full object-cover"
                                    onError={(e) => {
                                      const target = e.target as HTMLImageElement
                                      target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzIiIGhlaWdodD0iMzIiIHZpZXdCb3g9IjAgMCAzMiAzMiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjMyIiBoZWlnaHQ9IjMyIiByeD0iNiIgZmlsbD0idXJsKCNncmFkaWVudCkiLz4KPHN2ZyB4PSI4IiB5PSI4IiB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgc3Ryb2tlPSJ3aGl0ZSIgc3Ryb2tlLXdpZHRoPSIyIj4KPHBhdGggZD0iTTIyIDE5YTIgMiAwIDAgMS0yIDJINGEyIDIgMCAwIDEtMi0yVjVhMiAyIDAgMCAxIDItMmgxNmEyIDIgMCAwIDEgMiAydjE0eiIvPgo8cGF0aCBkPSJtMjIgMTMgMi0yIi8+CjxzdmcgeD0iOSIgeT0iOSIgd2lkdGg9IjYiIGhlaWdodD0iNiIgdmlld0JveD0iMCAwIDI0IDI0IiBmaWxsPSJub25lIiBzdHJva2U9IndoaXRlIiBzdHJva2Utd2lkdGg9IjIiPgo8Y2lyY2xlIGN4PSI5IiBjeT0iOSIgcj0iMiIvPgo8L3N2Zz4KPC9zdmc+CjwvcD4KPGRlZnM+CjxsaW5lYXJHcmFkaWVudCBpZD0iZ3JhZGllbnQiIHgxPSIwJSIgeTE9IjAlIiB4Mj0iMTAwJSIgeTI9IjEwMCUiPgo8c3RvcCBvZmZzZXQ9IjAlIiBzdHlsZT0ic3RvcC1jb2xvcjojOWNhM2FmO3N0b3Atb3BhY2l0eToxIiAvPgo8c3RvcCBvZmZzZXQ9IjEwMCUiIHN0eWxlPSJzdG9wLWNvbG9yOiM0Yjc0OGU7c3RvcC1vcGFjaXR5OjEiIC8+CjwvbGluZWFyR3JhZGllbnQ+CjwvZGVmcz4KPC9zdmc+Cg=='
                                    }}
                                  />
                                ) : (
                                  <span className="text-xs text-white">└─</span>
                                )}
                              </div>
                              <div>
                                <div className="flex items-center space-x-2">
                                  <h3 className="text-sm font-medium text-gray-900">
                                    {child.name}
                                  </h3>
                                  <Badge variant="outline" className="text-xs">
                                    Sub-category
                                  </Badge>
                                </div>
                                <p className="text-sm text-gray-500 mt-1">
                                  {child.description || 'No description'}
                                </p>
                                <div className="flex items-center space-x-4 mt-2 text-xs text-gray-400">
                                  <span>ID: {child.id}</span>
                                  <span>Level: {getCategoryLevel(child)}</span>
                                  {child.created_at && (
                                    <span>Created: {new Date(child.created_at).toLocaleDateString()}</span>
                                  )}
                                </div>
                              </div>
                            </div>
                          </div>

                          <div className="flex items-center space-x-2">
                            <Badge 
                              variant={child.is_active ? "default" : "secondary"}
                            >
                              {child.is_active ? 'Active' : 'Inactive'}
                            </Badge>

                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" className="h-8 w-8 p-0">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end">
                                <DropdownMenuItem onClick={() => handleViewCategory(child)}>
                                  <Eye className="mr-2 h-4 w-4" />
                                  View Details
                                </DropdownMenuItem>
                                <RequirePermission permission={PERMISSIONS.CATEGORIES_UPDATE}>
                                  <DropdownMenuItem onClick={() => handleEditCategory(child)}>
                                    <Edit className="mr-2 h-4 w-4" />
                                    Edit
                                  </DropdownMenuItem>
                                </RequirePermission>
                                <DropdownMenuSeparator />
                                <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                                  <DropdownMenuItem 
                                    onClick={() => handleDeleteCategory(child.id)}
                                    className="text-red-600 focus:text-red-600"
                                  >
                                    <Trash2 className="mr-2 h-4 w-4" />
                                    Delete
                                  </DropdownMenuItem>
                                </RequirePermission>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                )
              })}

              {/* Orphaned Child Categories (categories with parent_id but no existing parent) */}
              {childCategories.filter(child => 
                !parentCategories.some(parent => parent.id === child.parent_id)
              ).map((orphan) => (
                <div key={orphan.id} className="p-4 hover:bg-gray-50 border-l-4 border-yellow-400">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4 flex-1">
                      <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 bg-gradient-to-br from-yellow-500 to-orange-600 rounded-lg flex items-center justify-center">
                          <span className="text-white text-xs">⚠</span>
                        </div>
                        <div>
                          <div className="flex items-center space-x-2">
                            <h3 className="text-sm font-medium text-gray-900">
                              {orphan.name}
                            </h3>
                            <Badge variant="destructive">
                              Orphaned
                            </Badge>
                          </div>
                          <p className="text-sm text-gray-500 mt-1">
                            {orphan.description || 'No description'}
                          </p>
                          <p className="text-xs text-yellow-600 mt-1">
                            Parent category not found (ID: {orphan.parent_id})
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="flex items-center space-x-2">
                      <Badge 
                        variant={orphan.is_active ? "default" : "secondary"}
                      >
                        {orphan.is_active ? 'Active' : 'Inactive'}
                      </Badge>

                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" className="h-8 w-8 p-0">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => handleViewCategory(orphan)}>
                            <Eye className="mr-2 h-4 w-4" />
                            View Details
                          </DropdownMenuItem>
                          <RequirePermission permission={PERMISSIONS.CATEGORIES_UPDATE}>
                            <DropdownMenuItem onClick={() => handleEditCategory(orphan)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit
                            </DropdownMenuItem>
                          </RequirePermission>
                          <DropdownMenuSeparator />
                          <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                            <DropdownMenuItem 
                              onClick={() => handleDeleteCategory(orphan.id)}
                              className="text-red-600 focus:text-red-600"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete
                            </DropdownMenuItem>
                          </RequirePermission>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteCategoryId} onOpenChange={() => setDeleteCategoryId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the category
              and may affect products assigned to this category.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction 
              onClick={confirmDeleteCategory}
              className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
            >
              Delete Category
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* View Category Dialog */}
      <Dialog open={showCategoryModal} onOpenChange={setShowCategoryModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Category Details</DialogTitle>
          </DialogHeader>
          {selectedCategory && (
            <div className="space-y-6">
              {/* Category Image */}
              {selectedCategory.image && (
                <div>
                  <label className="text-sm font-medium text-gray-500">Category Image</label>
                  <div className="mt-2">
                    <div className="aspect-video w-full max-w-md mx-auto bg-gray-100 rounded-lg overflow-hidden">
                      <img
                        src={selectedCategory.image}
                        alt={selectedCategory.name}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          const target = e.target as HTMLImageElement
                          target.src = '/placeholder-product.svg'
                        }}
                      />
                    </div>
                    <p className="text-xs text-gray-400 mt-1 break-all">{selectedCategory.image}</p>
                  </div>
                </div>
              )}
              
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-500">Name</label>
                  <p className="text-sm text-gray-900 mt-1">{selectedCategory.name}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">ID</label>
                  <p className="text-sm text-gray-900 mt-1">{selectedCategory.id}</p>
                </div>
                <div className="col-span-2">
                  <label className="text-sm font-medium text-gray-500">Description</label>
                  <p className="text-sm text-gray-900 mt-1">
                    {selectedCategory.description || 'No description provided'}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Level</label>
                  <p className="text-sm text-gray-900 mt-1">{getCategoryLevel(selectedCategory)}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Status</label>
                  <div className="mt-1">
                    <Badge variant={selectedCategory.is_active ? "default" : "secondary"}>
                      {selectedCategory.is_active ? 'Active' : 'Inactive'}
                    </Badge>
                  </div>
                </div>
                {selectedCategory.parent_id && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">Parent Category</label>
                    <p className="text-sm text-gray-900 mt-1">{selectedCategory.parent_id}</p>
                  </div>
                )}
                <div>
                  <label className="text-sm font-medium text-gray-500">Sub-categories</label>
                  <p className="text-sm text-gray-900 mt-1">
                    {getChildrenCount(selectedCategory.id)} sub-categories
                  </p>
                </div>
                {selectedCategory.created_at && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">Created</label>
                    <p className="text-sm text-gray-900 mt-1">
                      {new Date(selectedCategory.created_at).toLocaleString()}
                    </p>
                  </div>
                )}
                {selectedCategory.updated_at && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">Last Updated</label>
                    <p className="text-sm text-gray-900 mt-1">
                      {new Date(selectedCategory.updated_at).toLocaleString()}
                    </p>
                  </div>
                )}
              </div>
              
              {/* Category Status Toggle */}
              <CategoryStatusToggle 
                category={selectedCategory} 
                onUpdate={() => {
                  refetch()
                  // Update the selected category state to reflect changes
                  setSelectedCategory(prev => prev ? { ...prev, is_active: !prev.is_active } : null)
                }}
              />
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Category Dialog */}
      <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-hidden">
          <DialogHeader>
            <DialogTitle>Edit Category</DialogTitle>
          </DialogHeader>
          {editCategory && (
            <EditCategoryForm
              category={editCategory}
              onSuccess={() => {
                setShowEditModal(false)
                setEditCategory(null)
                refetch()
              }}
              onCancel={() => {
                setShowEditModal(false)
                setEditCategory(null)
              }}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Add Subcategory Dialog */}
      <Dialog open={showAddSubcategoryModal} onOpenChange={setShowAddSubcategoryModal}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-hidden">
          <DialogHeader>
            <DialogTitle>
              Add Subcategory to "{addSubcategoryParent?.name}"
            </DialogTitle>
          </DialogHeader>
          {addSubcategoryParent && (
            <AddCategoryForm
              parentCategory={addSubcategoryParent}
              onSuccess={() => {
                setShowAddSubcategoryModal(false)
                setAddSubcategoryParent(null)
                refetch()
              }}
              onCancel={() => {
                setShowAddSubcategoryModal(false)
                setAddSubcategoryParent(null)
              }}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Manage Subcategories Dialog */}
      <Dialog open={showManageSubcategoriesModal} onOpenChange={setShowManageSubcategoriesModal}>
        <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              Manage Subcategories
            </DialogTitle>
          </DialogHeader>
          {manageSubcategoriesParent && (
            <SubcategoryManager
              parentCategory={manageSubcategoriesParent}
              subcategories={childCategories.filter(cat => cat.parent_id === manageSubcategoriesParent.id)}
              onRefresh={() => {
                refetch()
              }}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Error Dialog */}
      {deleteError && (
        <ErrorDialog
          open={showDeleteErrorModal}
          onOpenChange={setShowDeleteErrorModal}
          error={deleteError}
          onRetry={() => {
            // Optionally allow retry
            if (deleteCategoryId) {
              confirmDeleteCategory()
            }
          }}
        />
      )}
    </div>
  )
}
