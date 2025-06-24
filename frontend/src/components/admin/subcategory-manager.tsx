'use client'

import { useState } from 'react'
import { Plus, Edit, Trash2, Eye, MoreHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
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
import { ErrorDialog } from '@/components/ui/error-dialog'
import { useDeleteCategory } from '@/hooks/use-categories'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { getSubcategoryDeleteErrorMessage, CategoryDeleteError } from '@/lib/utils/category-delete-errors'

interface SubcategoryManagerProps {
  parentCategory: Category
  subcategories: Category[]
  onRefresh: () => void
}

export function SubcategoryManager({ 
  parentCategory, 
  subcategories, 
  onRefresh 
}: SubcategoryManagerProps) {
  const [showAddForm, setShowAddForm] = useState(false)
  const [deleteSubcategoryId, setDeleteSubcategoryId] = useState<string | null>(null)
  const [selectedSubcategory, setSelectedSubcategory] = useState<Category | null>(null)
  const [showViewModal, setShowViewModal] = useState(false)
  const [editSubcategory, setEditSubcategory] = useState<Category | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [deleteError, setDeleteError] = useState<CategoryDeleteError | null>(null)
  const [showDeleteErrorModal, setShowDeleteErrorModal] = useState(false)

  const deleteCategory = useDeleteCategory()

  const handleAddSubcategory = () => {
    setShowAddForm(true)
  }

  const handleViewSubcategory = (subcategory: Category) => {
    setSelectedSubcategory(subcategory)
    setShowViewModal(true)
  }

  const handleEditSubcategory = (subcategory: Category) => {
    setEditSubcategory(subcategory)
    setShowEditModal(true)
  }

  const handleDeleteSubcategory = (subcategoryId: string) => {
    setDeleteSubcategoryId(subcategoryId)
  }

  const confirmDeleteSubcategory = async () => {
    if (!deleteSubcategoryId) return
    
    try {
      await deleteCategory.mutateAsync(deleteSubcategoryId)
      setDeleteSubcategoryId(null)
      onRefresh()
      toast.success('Subcategory deleted successfully!')
    } catch (error: any) {
      console.error('Failed to delete subcategory:', error)
      
      // Parse the error and show detailed error dialog
      const deleteError = getSubcategoryDeleteErrorMessage(error)
      setDeleteError(deleteError)
      setShowDeleteErrorModal(true)
      setDeleteSubcategoryId(null) // Close the delete confirmation dialog
    }
  }

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle className="text-lg">
              Subcategories of "{parentCategory.name}"
            </CardTitle>
            <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
              <Button onClick={handleAddSubcategory} size="sm">
                <Plus className="mr-2 h-4 w-4" />
                Add Subcategory
              </Button>
            </RequirePermission>
          </div>
        </CardHeader>
        <CardContent>
          {subcategories.length === 0 ? (
            <div className="text-center py-8">
              <div className="text-gray-400 mb-2">
                <Plus className="mx-auto h-12 w-12" />
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">No subcategories</h3>
              <p className="text-gray-500 mb-4">
                Add subcategories to organize your products better.
              </p>
              <RequirePermission permission={PERMISSIONS.CATEGORIES_CREATE}>
                <Button onClick={handleAddSubcategory}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add First Subcategory
                </Button>
              </RequirePermission>
            </div>
          ) : (
            <div className="space-y-3">
              {subcategories.map((subcategory) => (
                <div
                  key={subcategory.id}
                  className="flex items-center justify-between p-4 border rounded-lg hover:bg-gray-50"
                >
                  <div className="flex items-center space-x-4">
                    <div className="w-8 h-8 bg-gradient-to-br from-gray-400 to-gray-600 rounded-md flex items-center justify-center">
                      <span className="text-xs text-white">└─</span>
                    </div>
                    <div>
                      <div className="flex items-center space-x-2">
                        <h4 className="font-medium text-gray-900">{subcategory.name}</h4>
                        <Badge variant="outline" className="text-xs">
                          Subcategory
                        </Badge>
                      </div>
                      <p className="text-sm text-gray-500 mt-1">
                        {subcategory.description || 'No description'}
                      </p>
                      <div className="flex items-center space-x-4 mt-2 text-xs text-gray-400">
                        <span>ID: {subcategory.id}</span>
                        {subcategory.created_at && (
                          <span>Created: {new Date(subcategory.created_at).toLocaleDateString()}</span>
                        )}
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center space-x-2">
                    <Badge variant={subcategory.is_active ? "default" : "secondary"}>
                      {subcategory.is_active ? 'Active' : 'Inactive'}
                    </Badge>

                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => handleViewSubcategory(subcategory)}>
                          <Eye className="mr-2 h-4 w-4" />
                          View Details
                        </DropdownMenuItem>
                        <RequirePermission permission={PERMISSIONS.CATEGORIES_UPDATE}>
                          <DropdownMenuItem onClick={() => handleEditSubcategory(subcategory)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Edit
                          </DropdownMenuItem>
                        </RequirePermission>
                        <DropdownMenuSeparator />
                        <RequirePermission permission={PERMISSIONS.CATEGORIES_DELETE}>
                          <DropdownMenuItem 
                            onClick={() => handleDeleteSubcategory(subcategory.id)}
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
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Add Subcategory Form */}
      {showAddForm && (
        <Card>
          <CardHeader>
            <CardTitle>Add New Subcategory</CardTitle>
          </CardHeader>
          <CardContent>
            <AddCategoryForm
              parentCategory={parentCategory}
              onSuccess={() => {
                setShowAddForm(false)
                onRefresh()
              }}
              onCancel={() => setShowAddForm(false)}
            />
          </CardContent>
        </Card>
      )}

      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteSubcategoryId} onOpenChange={() => setDeleteSubcategoryId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the subcategory
              and may affect products assigned to this subcategory.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction 
              onClick={confirmDeleteSubcategory}
              className="bg-red-600 hover:bg-red-700 focus:ring-red-600"
            >
              Delete Subcategory
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* View Subcategory Dialog */}
      <Dialog open={showViewModal} onOpenChange={setShowViewModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Subcategory Details</DialogTitle>
          </DialogHeader>
          {selectedSubcategory && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="text-sm font-medium text-gray-500">Name</label>
                  <p className="text-sm text-gray-900 mt-1">{selectedSubcategory.name}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">ID</label>
                  <p className="text-sm text-gray-900 mt-1">{selectedSubcategory.id}</p>
                </div>
                <div className="col-span-2">
                  <label className="text-sm font-medium text-gray-500">Description</label>
                  <p className="text-sm text-gray-900 mt-1">
                    {selectedSubcategory.description || 'No description provided'}
                  </p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Parent Category</label>
                  <p className="text-sm text-gray-900 mt-1">{parentCategory.name}</p>
                </div>
                <div>
                  <label className="text-sm font-medium text-gray-500">Status</label>
                  <div className="mt-1">
                    <Badge variant={selectedSubcategory.is_active ? "default" : "secondary"}>
                      {selectedSubcategory.is_active ? 'Active' : 'Inactive'}
                    </Badge>
                  </div>
                </div>
                {selectedSubcategory.created_at && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">Created</label>
                    <p className="text-sm text-gray-900 mt-1">
                      {new Date(selectedSubcategory.created_at).toLocaleString()}
                    </p>
                  </div>
                )}
                {selectedSubcategory.updated_at && (
                  <div>
                    <label className="text-sm font-medium text-gray-500">Last Updated</label>
                    <p className="text-sm text-gray-900 mt-1">
                      {new Date(selectedSubcategory.updated_at).toLocaleString()}
                    </p>
                  </div>
                )}
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Subcategory Dialog */}
      <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Edit Subcategory</DialogTitle>
          </DialogHeader>
          {editSubcategory && (
            <EditCategoryForm
              category={editSubcategory}
              onSuccess={() => {
                setShowEditModal(false)
                setEditSubcategory(null)
                onRefresh()
              }}
              onCancel={() => {
                setShowEditModal(false)
                setEditSubcategory(null)
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
            if (deleteSubcategoryId) {
              confirmDeleteSubcategory()
            }
          }}
        />
      )}
    </div>
  )
}
