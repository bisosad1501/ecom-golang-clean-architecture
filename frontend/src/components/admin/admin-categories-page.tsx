'use client'

import { useState } from 'react'
import { Plus, Search, Filter, MoreHorizontal, Edit, Trash2, Eye, FolderTree, ChevronDown, ChevronRight, Settings, Expand, Minimize, Tag, Folder, FolderOpen, Grid, List, Move, Calendar } from 'lucide-react'
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
      <div className={BIHUB_ADMIN_THEME.spacing.section}>
        <BiHubPageHeader
          title="Add New Category"
          subtitle="Create a new category to organize your BiHub products"
          breadcrumbs={[
            { label: 'Admin' },
            { label: 'Categories' },
            { label: 'Add Category' }
          ]}
          action={
            <Button
              variant="outline"
              onClick={() => setShowAddForm(false)}
              className={BIHUB_ADMIN_THEME.components.button.secondary}
            >
              Back to Categories
            </Button>
          }
        />

        <BiHubAdminCard
          title="Category Information"
          subtitle="Fill in the details for your new category"
          icon={<FolderTree className="h-5 w-5 text-white" />}
        >
          <AddCategoryForm
            onSuccess={() => {
              setShowAddForm(false)
              refetch()
            }}
            onCancel={() => setShowAddForm(false)}
          />
        </BiHubAdminCard>
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

      {/* Modern Search & Filters */}
      <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-6 shadow-lg">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-purple-500/5 rounded-2xl"></div>
        <div className="relative">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 border border-blue-400/30 flex items-center justify-center">
                <Search className="h-5 w-5 text-blue-400" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">Search & Filter Categories</h3>
                <p className="text-sm text-gray-400">Find and organize your BiHub categories</p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className="group relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white transition-all duration-200"
            >
              <div className="flex items-center gap-2">
                {viewMode === 'grid' ? (
                  <List className="h-4 w-4" />
                ) : (
                  <Grid className="h-4 w-4" />
                )}
                <span className="text-sm font-medium">
                  {viewMode === 'grid' ? 'List View' : 'Grid View'}
                </span>
              </div>
            </Button>
          </div>

          {/* Search Controls */}
          <div className="flex flex-col lg:flex-row items-center gap-4">
            <div className="flex-1 w-full">
              <div className="relative group">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 group-focus-within:text-blue-400 transition-colors" />
                <Input
                  placeholder="Search categories by name or description..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full h-12 pl-12 pr-12 bg-white/5 border-gray-600/50 rounded-xl text-white placeholder:text-gray-400 focus:bg-white/10 focus:border-blue-500/50 focus:ring-2 focus:ring-blue-500/20 transition-all duration-200"
                />
                {searchQuery && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSearchQuery('')}
                    className="absolute right-3 top-1/2 -translate-y-1/2 h-6 w-6 p-0 text-gray-400 hover:text-white hover:bg-white/10 rounded-md transition-all duration-200"
                  >
                    ×
                  </Button>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Categories List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 animate-pulse">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gray-700/50 rounded-xl"></div>
                  <div>
                    <div className="h-5 bg-gray-700/50 rounded w-24 mb-2"></div>
                    <div className="h-4 bg-gray-700/50 rounded w-16"></div>
                  </div>
                </div>
                {viewMode === 'grid' && (
                  <div className="h-6 bg-gray-700/50 rounded w-20"></div>
                )}
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-700/50 rounded w-3/4"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/2"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/4"></div>
              </div>
              <div className="flex items-center gap-3 mt-6 pt-4 border-t border-gray-700/50">
                <div className="h-8 bg-gray-700/50 rounded w-20"></div>
                <div className="h-8 bg-gray-700/50 rounded w-8"></div>
              </div>
            </div>
          ))}
        </div>
      ) : filteredCategories.length === 0 ? (
        <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-12 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-slate-500/5 rounded-2xl"></div>
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-500/20 to-slate-500/20 border border-gray-400/30 flex items-center justify-center mx-auto mb-6">
              <Tag className="h-8 w-8 text-gray-400" />
            </div>
            <h3 className="text-xl font-bold text-white mb-2">
              {searchQuery ? 'No categories found' : 'No categories yet'}
            </h3>
            <p className="text-gray-400 mb-6 max-w-md mx-auto">
              {searchQuery
                ? `No categories found matching "${searchQuery}". Try adjusting your search terms.`
                : 'Start organizing your BiHub products by creating categories.'
              }
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                <Button
                  onClick={() => setShowAddForm(true)}
                  className="bg-blue-500/20 border border-blue-400/30 text-blue-400 hover:bg-blue-500/30 hover:border-blue-400/50 hover:text-blue-300 transition-all duration-200"
                >
                  <Plus className="mr-2 h-5 w-5" />
                  Create First Category
                </Button>
              </RequirePermission>

              {searchQuery && (
                <Button
                  onClick={() => setSearchQuery('')}
                  className="bg-gray-500/20 border border-gray-400/30 text-gray-400 hover:bg-gray-500/30 hover:border-gray-400/50 hover:text-gray-300 transition-all duration-200"
                >
                  Clear Search
                </Button>
              )}
            </div>
          </div>
        </div>
      ) : (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {filteredCategories.map((category) => (
            <div
              key={category.id}
              className={cn(
                'group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 hover:bg-white/10 hover:border-gray-600/50 hover:scale-[1.02] transition-all duration-200 shadow-lg hover:shadow-xl',
                viewMode === 'list' && 'flex items-center gap-6'
              )}
            >
              {/* Gradient Background */}
              <div className={cn(
                'absolute inset-0 bg-gradient-to-br opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-2xl',
                category.is_active ? 'from-emerald-500/5 to-green-600/5' : 'from-red-500/5 to-rose-600/5'
              )} />

              {/* Category Header */}
              <div className={cn(
                'relative flex items-center justify-between mb-4',
                viewMode === 'list' && 'flex-1 mb-0'
              )}>
                <div className="flex items-center gap-4">
                  {/* Modern Category Icon */}
                  <div className={cn(
                    'relative w-12 h-12 rounded-xl bg-gradient-to-br flex items-center justify-center shadow-lg border border-white/10',
                    getCategoryColor(category)
                  )}>
                    <div className="relative z-10 text-white">
                      {getCategoryIcon(category)}
                    </div>
                    <div className="absolute inset-0 bg-white/10 rounded-xl blur-sm"></div>
                  </div>

                  <div>
                    <h3 className="text-lg font-bold text-white group-hover:text-[#FF9000] transition-colors">
                      {category.name}
                    </h3>
                    {/* Modern Status Badge */}
                    <div className={cn(
                      'inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold mt-2 border border-white/10 backdrop-blur-sm',
                      category.is_active
                        ? 'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-800'
                        : 'bg-red-100 text-red-800 border-red-200 dark:bg-red-950/30 dark:text-red-300 dark:border-red-800'
                    )}>
                      <div className={cn(
                        'w-2 h-2 rounded-full mr-2',
                        category.is_active ? 'bg-emerald-500' : 'bg-red-500'
                      )} />
                      {category.is_active ? 'Active' : 'Inactive'}
                    </div>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <p className="text-2xl font-bold text-[#FF9000] group-hover:text-[#FF9000]/80 transition-colors">
                      {getChildrenCount(category.id)}
                    </p>
                    <p className="text-xs text-gray-400 mt-1">Subcategories</p>
                  </div>
                )}
              </div>

              {/* Category Details */}
              <div className={cn(
                'relative space-y-3',
                viewMode === 'list' && 'flex-1 grid grid-cols-2 md:grid-cols-4 gap-4 space-y-0'
              )}>
                {category.description && (
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                      <FolderTree className="h-4 w-4 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 uppercase tracking-wide">Description</p>
                      <p className="text-sm text-white font-medium line-clamp-1">
                        {category.description}
                      </p>
                    </div>
                  </div>
                )}

                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                    <Calendar className="h-4 w-4 text-gray-400" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 uppercase tracking-wide">Created</p>
                    <p className="text-sm text-white font-medium">
                      {formatDate(category.created_at)}
                    </p>
                  </div>
                </div>

                {category.parent_id && (
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                      <FolderTree className="h-4 w-4 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 uppercase tracking-wide">Type</p>
                      <p className="text-sm text-white font-medium">Subcategory</p>
                    </div>
                  </div>
                )}

                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                    <FolderTree className="h-4 w-4 text-gray-400" />
                  </div>
                  <div>
                    <p className="text-xs text-gray-400 uppercase tracking-wide">Level</p>
                    <p className="text-sm text-white font-medium">
                      {getCategoryLevel(category)}
                    </p>
                  </div>
                </div>
              </div>

              {/* Modern Action Buttons */}
              <div className={cn(
                'relative flex items-center gap-3 mt-6 pt-4 border-t border-white/10',
                viewMode === 'list' && 'flex-shrink-0 mt-0 pt-0 border-t-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleEditCategory(category)}
                  className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-blue-500/10 hover:border-blue-500/50 hover:text-blue-400 transition-all duration-200"
                >
                  <Edit className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'Edit Category' : 'Edit'}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-purple-500/10 hover:border-purple-500/50 hover:text-purple-400 transition-all duration-200"
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                    <DropdownMenuItem
                      onClick={() => handleEditCategory(category)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Edit className="mr-3 h-4 w-4" />
                      Edit Category
                    </DropdownMenuItem>
                    <DropdownMenuItem
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Move className="mr-3 h-4 w-4" />
                      Move Category
                    </DropdownMenuItem>
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                      <DropdownMenuItem
                        onClick={() => handleDeleteCategory(category.id)}
                        className="text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded-lg m-1 p-3"
                      >
                        <Trash2 className="mr-3 h-4 w-4" />
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
