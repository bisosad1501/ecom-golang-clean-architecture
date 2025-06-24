'use client'

import { useState } from 'react'
import { Plus, Search, Filter, MoreHorizontal, Edit, Trash2, Eye, Copy, Archive } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useProducts, useDeleteProduct } from '@/hooks/use-products'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { formatPrice } from '@/lib/utils/price'
import { AdminPriceDisplay } from '@/components/ui/price-display'
import { AddProductForm } from '@/components/forms/add-product-form'
import { EditProductForm } from '@/components/forms/edit-product-form'
import { Product } from '@/types'
import { toast } from 'sonner'
import Image from 'next/image'
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

export function AdminProductsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [showAddForm, setShowAddForm] = useState(false)
  const [deleteProductId, setDeleteProductId] = useState<string | null>(null)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [showProductModal, setShowProductModal] = useState(false)
  const [editProduct, setEditProduct] = useState<Product | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  
  const { data, isLoading, refetch } = useProducts({
    page: currentPage,
    limit: 10,
    search: searchQuery || undefined,
  })

  const deleteProductMutation = useDeleteProduct()

  const products = data?.data || []
  const pagination = data?.pagination

  // Action handlers
  const handleViewProduct = (product: Product) => {
    setSelectedProduct(product)
    setShowProductModal(true)
  }

  const handleEditProduct = (product: Product) => {
    setEditProduct(product)
    setShowEditModal(true)
  }

  const handleDeleteProduct = (productId: string) => {
    setDeleteProductId(productId)
  }

  const confirmDeleteProduct = async () => {
    if (!deleteProductId) return
    
    try {
      await deleteProductMutation.mutateAsync(deleteProductId)
      setDeleteProductId(null)
      refetch()
      toast.success('Product deleted successfully!')
    } catch (error: any) {
      console.error('Failed to delete product:', error)
      const errorMessage = error?.message || error?.error || 'Failed to delete product. Please try again.'
      toast.error(errorMessage)
    }
  }

  const handleDuplicateProduct = (product: Product) => {
    // TODO: Implement duplicate functionality
    toast.info('Duplicate functionality coming soon!')
  }

  const handleArchiveProduct = (product: Product) => {
    // TODO: Implement archive functionality
    toast.info('Archive functionality coming soon!')
  }

  const handleCopyProductUrl = async (product: Product) => {
    try {
      const url = `${window.location.origin}/products/${product.id}`
      await navigator.clipboard.writeText(url)
      toast.success('Product URL copied to clipboard!')
    } catch (error) {
      // Fallback for browsers that don't support clipboard API
      const textArea = document.createElement('textarea')
      textArea.value = `${window.location.origin}/products/${product.id}`
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand('copy')
      document.body.removeChild(textArea)
      toast.success('Product URL copied to clipboard!')
    }
  }

  if (showAddForm) {
    return (
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Add New Product</h1>
            <p className="text-gray-600 mt-2">Create a new product in your catalog</p>
          </div>
          
          <Button 
            variant="outline" 
            onClick={() => setShowAddForm(false)}
          >
            Back to Products
          </Button>
        </div>

        <AddProductForm 
          onSuccess={() => setShowAddForm(false)}
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
          <h1 className="text-3xl font-bold text-gray-900">Products</h1>
          <p className="text-gray-600 mt-2">Manage your product catalog</p>
        </div>
        
        <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
          <Button onClick={() => setShowAddForm(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Product
          </Button>
        </RequirePermission>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-center space-x-4">
            <div className="flex-1">
              <Input
                placeholder="Search products..."
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
        </CardContent>
      </Card>

      {/* Products Table */}
      <Card>
        <CardHeader>
          <CardTitle>All Products</CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="p-6">
              <div className="space-y-4">
                {[...Array(5)].map((_, i) => (
                  <div key={i} className="animate-pulse flex items-center space-x-4 p-4">
                    <div className="w-16 h-16 bg-gray-200 rounded"></div>
                    <div className="flex-1 space-y-2">
                      <div className="h-4 bg-gray-200 rounded w-1/4"></div>
                      <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                    </div>
                    <div className="w-20 h-4 bg-gray-200 rounded"></div>
                  </div>
                ))}
              </div>
            </div>
          ) : products.length > 0 ? (
            <div className="overflow-x-auto">
              {/* Table Header - Desktop */}
              <div className="hidden md:grid grid-cols-[80px_1fr_140px_120px_100px_120px_140px] gap-4 p-4 bg-gray-50 border-b text-sm font-medium text-gray-600 min-w-[900px]">
                <div className="text-center">Image</div>
                <div>Product</div>
                <div className="text-center">Price</div>
                <div className="text-center">Cost</div>
                <div className="text-center">Stock</div>
                <div className="text-center">Status</div>
                <div className="text-center">Actions</div>
              </div>
              
              {/* Table Body - Desktop */}
              <div className="hidden md:block">
                {products.map((product) => (
                  <div key={product.id} className="grid grid-cols-[80px_1fr_140px_120px_100px_120px_140px] gap-4 p-4 border-b hover:bg-gray-50 items-center min-w-[900px]">
                    {/* Product Image */}
                    <div className="w-16 h-16 relative overflow-hidden rounded-md border bg-gray-100 mx-auto">
                      {product.images?.[0]?.url ? (
                        <Image
                          src={product.images[0].url}
                          alt={product.name}
                          fill
                          className="object-cover"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement;
                            target.src = '/placeholder-product.svg';
                          }}
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center text-gray-400">
                          <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                          </svg>
                        </div>
                      )}
                    </div>

                    {/* Product Info */}
                    <div className="min-w-0">
                      <h3 className="font-medium text-gray-900 truncate">{product.name}</h3>
                      <p className="text-sm text-gray-500 truncate">SKU: {product.sku}</p>
                      {product.tags && product.tags.length > 0 && (
                        <div className="flex items-center gap-1 mt-2 flex-wrap">
                          {product.tags.slice(0, 2).map((tag) => (
                            <Badge key={tag.id} variant="outline" className="text-xs">
                              {tag.name}
                            </Badge>
                          ))}
                          {product.tags.length > 2 && (
                            <Badge variant="outline" className="text-xs">
                              +{product.tags.length - 2}
                            </Badge>
                          )}
                        </div>
                      )}
                    </div>

                    {/* Price */}
                    <div className="text-center space-y-1">
                      <div className="font-semibold text-lg text-gray-900">
                        {formatPrice(product.price)}
                      </div>
                      {product.compare_price && product.compare_price > product.price && (
                        <>
                          <div className="text-xs text-gray-400 line-through">
                            {formatPrice(product.compare_price)}
                          </div>
                          <div className="text-xs text-green-600 font-medium">
                            {Math.round(((product.compare_price - product.price) / product.compare_price) * 100)}% OFF
                          </div>
                        </>
                      )}
                    </div>

                    {/* Cost */}
                    <div className="text-center space-y-1">
                      <div className="font-semibold text-lg text-gray-900">
                        {product.cost_price ? formatPrice(product.cost_price) : '-'}
                      </div>
                      {product.cost_price && product.price > product.cost_price && (
                        <div className="text-xs text-green-600 font-medium">
                          Profit: {formatPrice(product.price - product.cost_price)}
                        </div>
                      )}
                    </div>

                    {/* Stock */}
                    <div className="text-center space-y-1">
                      <div className="font-semibold text-lg text-gray-900">{product.stock}</div>
                      <div className="text-xs text-gray-500">units</div>
                    </div>

                    {/* Status */}
                    <div className="text-center space-y-1">
                      <div>
                        <Badge 
                          variant={product.stock > 0 ? 'default' : 'destructive'} 
                          className="text-xs"
                        >
                          {product.stock > 0 ? 'In Stock' : 'Out'}
                        </Badge>
                      </div>
                      <div>
                        <Badge 
                          variant={product.status === 'active' ? 'default' : 'secondary'}
                          className="text-xs"
                        >
                          {product.status}
                        </Badge>
                      </div>
                      {product.compare_price && product.compare_price > product.price && (
                        <div>
                          <Badge variant="secondary" className="text-xs">
                            Sale
                          </Badge>
                        </div>
                      )}
                    </div>

                    {/* Actions */}
                    <div className="flex items-center justify-center gap-1">
                      <Button 
                        variant="ghost" 
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => handleViewProduct(product)}
                        title="View product"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                      
                      <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                        <Button 
                          variant="ghost" 
                          size="icon"
                          className="h-8 w-8"
                          onClick={() => handleEditProduct(product)}
                          title="Edit product"
                        >
                          <Edit className="h-4 w-4" />
                        </Button>
                      </RequirePermission>
                      
                      <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
                        <Button 
                          variant="ghost" 
                          size="icon" 
                          className="h-8 w-8 text-red-600 hover:text-red-700"
                          onClick={() => handleDeleteProduct(product.id)}
                          title="Delete product"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </RequirePermission>
                      
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8" title="More actions">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-48">
                          <DropdownMenuItem onClick={() => handleCopyProductUrl(product)}>
                            <Copy className="mr-2 h-4 w-4" />
                            Copy URL
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleDuplicateProduct(product)}>
                            <Copy className="mr-2 h-4 w-4" />
                            Duplicate
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem onClick={() => handleArchiveProduct(product)}>
                            <Archive className="mr-2 h-4 w-4" />
                            Archive
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </div>
                ))}
              </div>

              {/* Mobile View */}
              <div className="md:hidden">
                {products.map((product) => (
                  <div key={product.id} className="border-b p-4 space-y-3">
                    <div className="flex items-start space-x-3">
                      {/* Product Image */}
                      <div className="w-16 h-16 relative overflow-hidden rounded-md border bg-gray-100 flex-shrink-0">
                        {product.images?.[0]?.url ? (
                          <Image
                            src={product.images[0].url}
                            alt={product.name}
                            fill
                            className="object-cover"
                            onError={(e) => {
                              const target = e.target as HTMLImageElement;
                              target.src = '/placeholder-product.svg';
                            }}
                          />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center text-gray-400">
                            <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                          </div>
                        )}
                      </div>

                      {/* Product Info */}
                      <div className="flex-1 min-w-0">
                        <h3 className="font-medium text-gray-900 truncate">{product.name}</h3>
                        <p className="text-sm text-gray-500">SKU: {product.sku}</p>
                        <div className="mt-2">
                          <AdminPriceDisplay product={product} />
                        </div>
                      </div>

                      {/* Actions */}
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" title="More actions">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-48">
                          <DropdownMenuItem onClick={() => handleViewProduct(product)}>
                            <Eye className="mr-2 h-4 w-4" />
                            View Details
                          </DropdownMenuItem>
                          <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                            <DropdownMenuItem onClick={() => handleEditProduct(product)}>
                              <Edit className="mr-2 h-4 w-4" />
                              Edit Product
                            </DropdownMenuItem>
                          </RequirePermission>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem onClick={() => handleCopyProductUrl(product)}>
                            <Copy className="mr-2 h-4 w-4" />
                            Copy URL
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleDuplicateProduct(product)}>
                            <Copy className="mr-2 h-4 w-4" />
                            Duplicate
                          </DropdownMenuItem>
                          <DropdownMenuItem onClick={() => handleArchiveProduct(product)}>
                            <Archive className="mr-2 h-4 w-4" />
                            Archive
                          </DropdownMenuItem>
                          <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
                            <DropdownMenuSeparator />
                            <DropdownMenuItem 
                              onClick={() => handleDeleteProduct(product.id)}
                              className="text-red-600"
                            >
                              <Trash2 className="mr-2 h-4 w-4" />
                              Delete Product
                            </DropdownMenuItem>
                          </RequirePermission>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>

                    {/* Bottom row with status badges and stock info */}
                    <div className="flex items-center justify-between pt-2 border-t">
                      <div className="flex items-center gap-2">
                        <Badge variant={product.stock > 0 ? 'default' : 'destructive'} className="text-xs">
                          {product.stock > 0 ? 'In Stock' : 'Out'}
                        </Badge>
                        <Badge variant={product.status === 'active' ? 'default' : 'secondary'} className="text-xs">
                          {product.status}
                        </Badge>
                        {product.compare_price && product.compare_price > product.price && (
                          <Badge variant="secondary" className="text-xs">Sale</Badge>
                        )}
                      </div>
                      <div className="text-sm text-gray-500">
                        <span className="font-medium">{product.stock}</span> units
                      </div>
                    </div>

                    {/* Tags */}
                    {product.tags && product.tags.length > 0 && (
                      <div className="flex items-center gap-1 flex-wrap">
                        {product.tags.slice(0, 3).map((tag) => (
                          <Badge key={tag.id} variant="outline" className="text-xs">
                            {tag.name}
                          </Badge>
                        ))}
                        {product.tags.length > 3 && (
                          <Badge variant="outline" className="text-xs">
                            +{product.tags.length - 3}
                          </Badge>
                        )}
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </div>
          ) : (
            <div className="text-center py-12 p-6">
              <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Plus className="h-8 w-8 text-gray-400" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">No products found</h3>
              <p className="text-gray-600 mb-6">
                {searchQuery 
                  ? `No products found matching "${searchQuery}". Try adjusting your search terms.` 
                  : 'Get started by adding your first product.'
                }
              </p>
              <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
                <Button onClick={() => setShowAddForm(true)}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Product
                </Button>
              </RequirePermission>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-gray-600">
            Showing {((pagination.page - 1) * pagination.limit) + 1} to {Math.min(pagination.page * pagination.limit, pagination.total)} of {pagination.total} products
          </p>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={!pagination.has_prev}
            >
              Previous
            </Button>
            <span className="text-sm text-gray-600">
              Page {pagination.page} of {pagination.total_pages}
            </span>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage(prev => prev + 1)}
              disabled={!pagination.has_next}
            >
              Next
            </Button>
          </div>
        </div>
      )}
      
      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteProductId} onOpenChange={() => setDeleteProductId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Product</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this product? This action cannot be undone.
              All associated data including images, reviews, and sales history will be permanently removed.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteProductMutation.isPending}>
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteProduct}
              disabled={deleteProductMutation.isPending}
              className="bg-red-600 hover:bg-red-700"
            >
              {deleteProductMutation.isPending ? 'Deleting...' : 'Delete Product'}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* Product View Modal */}
      {selectedProduct && (
        <AlertDialog open={showProductModal} onOpenChange={setShowProductModal}>
          <AlertDialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
            <AlertDialogHeader>
              <AlertDialogTitle className="flex items-center gap-3">
                <div className="w-12 h-12 relative overflow-hidden rounded-lg border">
                  {selectedProduct.images?.[0]?.url ? (
                    <Image
                      src={selectedProduct.images[0].url}
                      alt={selectedProduct.name}
                      fill
                      className="object-cover"
                      onError={(e) => {
                        const target = e.target as HTMLImageElement;
                        target.src = '/placeholder-product.svg';
                      }}
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-100">
                      <Eye className="w-6 h-6" />
                    </div>
                  )}
                </div>
                <div>
                  <h2 className="text-xl font-semibold">{selectedProduct.name}</h2>
                  <p className="text-sm text-gray-500">SKU: {selectedProduct.sku}</p>
                </div>
              </AlertDialogTitle>
            </AlertDialogHeader>
            
            <div className="space-y-6">
              {/* Product Details Grid */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <h3 className="font-medium text-gray-900 mb-2">Product Information</h3>
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-gray-500">Price:</span>
                        <span className="font-medium">{formatPrice(selectedProduct.price)}</span>
                      </div>
                      {selectedProduct.compare_price && selectedProduct.compare_price > selectedProduct.price && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Compare Price:</span>
                          <span className="line-through text-gray-400">{formatPrice(selectedProduct.compare_price)}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span className="text-gray-500">Stock:</span>
                        <span className={selectedProduct.stock > 0 ? 'text-green-600' : 'text-red-600'}>
                          {selectedProduct.stock} {selectedProduct.stock === 1 ? 'item' : 'items'}
                        </span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-gray-500">Status:</span>
                        <Badge variant={selectedProduct.status === 'active' ? 'default' : 'secondary'}>
                          {selectedProduct.status}
                        </Badge>
                      </div>
                    </div>
                  </div>
                  
                  {selectedProduct.tags && selectedProduct.tags.length > 0 && (
                    <div>
                      <h3 className="font-medium text-gray-900 mb-2">Tags</h3>
                      <div className="flex flex-wrap gap-1">
                        {selectedProduct.tags.map((tag) => (
                          <Badge key={tag.id} variant="outline" className="text-xs">
                            {tag.name}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
                
                <div className="space-y-4">
                  {selectedProduct.images && selectedProduct.images.length > 0 && (
                    <div>
                      <h3 className="font-medium text-gray-900 mb-2">Images ({selectedProduct.images.length})</h3>
                      <div className="grid grid-cols-3 gap-2">
                        {selectedProduct.images.slice(0, 6).map((image, index) => (
                          <div key={index} className="aspect-square relative overflow-hidden rounded-lg border">
                            <Image
                              src={image.url}
                              alt={image.alt_text || `${selectedProduct.name} - Image ${index + 1}`}
                              fill
                              className="object-cover"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.src = '/placeholder-product.svg';
                              }}
                            />
                          </div>
                        ))}
                        {selectedProduct.images.length > 6 && (
                          <div className="aspect-square flex items-center justify-center bg-gray-100 rounded-lg border text-gray-500 text-sm">
                            +{selectedProduct.images.length - 6} more
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </div>
              
              {/* Description */}
              <div>
                <h3 className="font-medium text-gray-900 mb-2">Description</h3>
                <p className="text-gray-700 text-sm leading-relaxed">{selectedProduct.description}</p>
              </div>
              
              {/* Short Description */}
              {selectedProduct.short_description && (
                <div>
                  <h3 className="font-medium text-gray-900 mb-2">Short Description</h3>
                  <p className="text-gray-700 text-sm">{selectedProduct.short_description}</p>
                </div>
              )}
            </div>
            
            <AlertDialogFooter>
              <AlertDialogCancel>Close</AlertDialogCancel>
              <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                <AlertDialogAction onClick={() => {
                  setShowProductModal(false)
                  handleEditProduct(selectedProduct)
                }}>
                  Edit Product
                </AlertDialogAction>
              </RequirePermission>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}

      {/* Edit Product Modal */}
      {editProduct && (
        <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
          <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>Edit Product: {editProduct.name}</DialogTitle>
            </DialogHeader>
            <EditProductForm
              product={editProduct}
              onSuccess={() => {
                setShowEditModal(false)
                setEditProduct(null)
                refetch()
              }}
              onCancel={() => {
                setShowEditModal(false)
                setEditProduct(null)
              }}
            />
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}