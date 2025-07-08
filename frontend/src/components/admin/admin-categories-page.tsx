'use client'

import { useState } from 'react'
import { Plus, Search, Filter, MoreHorizontal, Edit, Trash2, Eye, FolderTree, ChevronDown, ChevronRight, Settings, Expand, Minimize, Tag, Folder, FolderOpen, Grid, List, Move } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
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
import { useCategories, useDeleteCategory } from '@/hooks/use-categories'
import { formatDate } from '@/lib/utils'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export function AdminCategoriesPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [showAddForm, setShowAddForm] = useState(false)
  const [deleteCategoryId, setDeleteCategoryId] = useState<string | null>(null)
  const [selectedCategory, setSelectedCategory] = useState<Category | null>(null)
  const [showCategoryModal, setShowCategoryModal] = useState(false)
  const [editCategory, setEditCategory] = useState<Category | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  
  const { data: categories, isLoading, refetch } = useCategories()
  const deleteCategory = useDeleteCategory()

  // Filter categories based on search
  const filteredCategories = categories?.filter(category =>
    category.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    category.description?.toLowerCase().includes(searchQuery.toLowerCase())
  ) || []

  // Mock stats - replace with real data
  const stats = {
    totalCategories: filteredCategories.length,
    activeCategories: filteredCategories.filter(c => c.is_active).length,
    parentCategories: filteredCategories.filter(c => !c.parent_id).length,
    subCategories: filteredCategories.filter(c => c.parent_id).length,
  }

  const getCategoryIcon = (category: Category) => {
    if (category.parent_id) {
      return <Folder className="h-4 w-4 text-white" />
    }
    return <FolderOpen className="h-4 w-4 text-white" />
  }

  const getCategoryColor = (category: Category) => {
    if (category.parent_id) {
      return 'from-blue-500 to-blue-600'
    }
    return 'from-purple-500 to-purple-600'
  }

  // Action handlers
  const handleEditCategory = (category: Category) => {
    setEditCategory(category)
    setShowEditModal(true)
  }

  const handleDeleteCategory = (categoryId: string) => {
    setDeleteCategoryId(categoryId)
  }

  const confirmDeleteCategory = async () => {
    if (!deleteCategoryId) return

    try {
      await deleteCategory.mutateAsync(deleteCategoryId)
      setDeleteCategoryId(null)
      refetch()
      toast.success('Category deleted successfully!')
    } catch (error: any) {
      console.error('Failed to delete category:', error)
      toast.error('Failed to delete category')
      setDeleteCategoryId(null)
    }
  }

  const getCategoryLevel = (category: Category) => {
    return category.level || 0
  }

  const getChildrenCount = (parentId: string) => {
    return categories?.filter((cat: Category) => cat.parent_id === parentId).length || 0
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
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Category Management"
        subtitle="Organize and manage BiHub product categories and subcategories"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Categories' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
            <Button
              onClick={() => setShowAddForm(true)}
              className={BIHUB_ADMIN_THEME.components.button.primary}
            >
              <Plus className="mr-2 h-5 w-5" />
              Add Category
            </Button>
          </RequirePermission>
        }
      />

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <BiHubStatCard
          title="Total Categories"
          value={stats.totalCategories}
          icon={<Tag className="h-8 w-8 text-white" />}
          color="primary"
        />
        <BiHubStatCard
          title="Active Categories"
          value={stats.activeCategories}
          icon={<FolderOpen className="h-8 w-8 text-white" />}
          color="success"
        />
        <BiHubStatCard
          title="Parent Categories"
          value={stats.parentCategories}
          icon={<Folder className="h-8 w-8 text-white" />}
          color="info"
        />
        <BiHubStatCard
          title="Sub Categories"
          value={stats.subCategories}
          icon={<Folder className="h-8 w-8 text-white" />}
          color="warning"
        />
      </div>

      {/* Search & Filters */}
      <BiHubAdminCard
        title="Search & Filter Categories"
        subtitle="Find and organize BiHub categories"
        icon={<Search className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className={cn(
                BIHUB_ADMIN_THEME.components.button.ghost,
                'h-10 w-10 p-0'
              )}
            >
              {viewMode === 'grid' ? (
                <List className="h-4 w-4" />
              ) : (
                <Grid className="h-4 w-4" />
              )}
            </Button>
          </div>
        }
      >
        <div className="flex flex-col lg:flex-row items-center gap-4">
          <div className="flex-1 w-full">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                placeholder="Search categories by name or description..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={cn(
                  BIHUB_ADMIN_THEME.components.input.base,
                  'pl-10 pr-12 h-12'
                )}
              />
              {searchQuery && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSearchQuery('')}
                  className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 text-gray-400 hover:text-white"
                >
                  ×
                </Button>
              )}
            </div>
          </div>
        </div>
      </BiHubAdminCard>

      {/* Categories List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className={cn(
              BIHUB_ADMIN_THEME.components.card.base,
              'p-6 animate-pulse'
            )}>
              <div className="flex items-center gap-4 mb-4">
                <div className="w-12 h-12 bg-gray-700 rounded-xl"></div>
                <div className="space-y-2 flex-1">
                  <div className="h-4 bg-gray-700 rounded w-1/3"></div>
                  <div className="h-3 bg-gray-700 rounded w-1/2"></div>
                </div>
              </div>
              <div className="space-y-2">
                <div className="h-3 bg-gray-700 rounded w-3/4"></div>
                <div className="h-3 bg-gray-700 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      ) : filteredCategories.length === 0 ? (
        <BiHubEmptyState
          icon={<Tag className="h-8 w-8 text-gray-400" />}
          title={searchQuery ? 'No categories found' : 'No categories yet'}
          description={
            searchQuery
              ? `No categories found matching "${searchQuery}". Try adjusting your search terms.`
              : 'Start organizing your BiHub products by creating categories.'
          }
          action={
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                <Button
                  onClick={() => setShowAddForm(true)}
                  className={BIHUB_ADMIN_THEME.components.button.primary}
                >
                  <Plus className="mr-2 h-5 w-5" />
                  Create First Category
                </Button>
              </RequirePermission>

              {searchQuery && (
                <Button
                  onClick={() => setSearchQuery('')}
                  className={BIHUB_ADMIN_THEME.components.button.secondary}
                >
                  Clear Search
                </Button>
              )}
            </div>
          }
        />
      ) : (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {filteredCategories.map((category) => (
            <div
              key={category.id}
              className={cn(
                BIHUB_ADMIN_THEME.components.card.base,
                BIHUB_ADMIN_THEME.components.card.hover,
                'group',
                viewMode === 'list' && 'flex items-center gap-6 p-6'
              )}
            >
              {/* Category Icon & Info */}
              <div className={cn(
                'flex items-center gap-4',
                viewMode === 'grid' ? 'mb-4' : 'flex-1'
              )}>
                <div className={cn(
                  'w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center',
                  getCategoryColor(category)
                )}>
                  {getCategoryIcon(category)}
                </div>

                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className={cn(
                      BIHUB_ADMIN_THEME.typography.heading.h4,
                      'group-hover:text-[#FF9000] transition-colors'
                    )}>
                      {category.name}
                    </h3>

                    <BiHubStatusBadge status={category.is_active ? 'success' : 'error'}>
                      {category.is_active ? 'Active' : 'Inactive'}
                    </BiHubStatusBadge>
                  </div>

                  {category.description && (
                    <p className={cn(
                      BIHUB_ADMIN_THEME.typography.body.small,
                      'line-clamp-2'
                    )}>
                      {category.description}
                    </p>
                  )}

                  <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
                    <span>Created: {formatDate(category.created_at)}</span>
                    {category.parent_id && (
                      <span>Subcategory</span>
                    )}
                  </div>
                </div>
              </div>
              {/* Actions */}
              <div className={cn(
                'flex items-center gap-2 mt-4',
                viewMode === 'list' && 'flex-shrink-0 mt-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleEditCategory(category)}
                  className={BIHUB_ADMIN_THEME.components.button.ghost}
                >
                  <Edit className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'Edit' : ''}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className={BIHUB_ADMIN_THEME.components.button.ghost}
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-48 bg-gray-900 border-gray-700">
                    <DropdownMenuItem
                      onClick={() => handleEditCategory(category)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <Edit className="mr-2 h-4 w-4" />
                      Edit Category
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-gray-300 hover:text-white hover:bg-gray-800"
                    >
                      <Move className="mr-2 h-4 w-4" />
                      Move Category
                    </DropdownMenuItem>
                    <DropdownMenuSeparator className="bg-gray-700" />
                    <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                      <DropdownMenuItem
                        onClick={() => handleDeleteCategory(category.id)}
                        className="text-red-400 hover:text-red-300 hover:bg-red-900/20"
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Delete Category
                      </DropdownMenuItem>
                    </RequirePermission>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Add Category Form */}
      {showAddForm && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className={cn(BIHUB_ADMIN_THEME.components.card.base, 'max-w-2xl w-full mx-4')}>
            <AddCategoryForm
              onSuccess={() => {
                setShowAddForm(false)
                refetch()
              }}
              onCancel={() => setShowAddForm(false)}
            />
          </div>
        </div>
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteCategoryId} onOpenChange={() => setDeleteCategoryId(null)}>
        <AlertDialogContent className="bg-gray-900 border-gray-700">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-white">Delete Category</AlertDialogTitle>
            <AlertDialogDescription className="text-gray-400">
              Are you sure you want to delete this category? This action cannot be undone.
              All products in this category will need to be reassigned.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel className={BIHUB_ADMIN_THEME.components.button.secondary}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteCategory}
              className="bg-red-600 hover:bg-red-700 text-white"
            >
              Delete Category
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Edit Category Dialog */}
      <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
        <DialogContent className="max-w-2xl bg-gray-900 border-gray-700">
          <DialogHeader>
            <DialogTitle className="text-white">Edit Category</DialogTitle>
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
    </div>
  )
}
