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

// Helper function to get badge variant based on product status
function getStatusVariant(status: string) {
  switch (status) {
    case 'active':
      return 'default'
    case 'draft':
      return 'secondary'
    case 'archived':
      return 'outline'
    default:
      return 'secondary'
  }
}

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
    <div className="space-y-8">
      {/* Enhanced Header */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
        <div>
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-emerald-500 to-emerald-600 flex items-center justify-center shadow-large">
              <Plus className="h-6 w-6 text-white" />
            </div>
            <span className="text-primary font-semibold">PRODUCT MANAGEMENT</span>
          </div>

          <h1 className="text-4xl lg:text-5xl font-bold text-foreground mb-4">
            Product <span className="text-gradient">Catalog</span>
          </h1>
          <p className="text-xl text-muted-foreground">
            Manage your product inventory, pricing, and availability
          </p>

          {/* Quick stats */}
          <div className="flex items-center gap-6 mt-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-primary">{pagination?.total || 0}</div>
              <div className="text-sm text-muted-foreground">Total Products</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-emerald-600">
                {products.filter(p => p.stock > 0).length}
              </div>
              <div className="text-sm text-muted-foreground">In Stock</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-orange-600">
                {products.filter(p => p.stock <= 5).length}
              </div>
              <div className="text-sm text-muted-foreground">Low Stock</div>
            </div>
          </div>
        </div>

        <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
          <Button
            onClick={() => setShowAddForm(true)}
            size="xl"
            variant="gradient"
            className="shadow-large hover:shadow-xl transition-all duration-200"
          >
            <Plus className="mr-2 h-5 w-5" />
            Add New Product
          </Button>
        </RequirePermission>
      </div>

      {/* Enhanced Filters */}
      <Card variant="elevated" className="border-0 shadow-large">
        <CardContent className="p-8">
          <div className="flex flex-col lg:flex-row items-center gap-6">
            <div className="flex-1 w-full">
              <div className="relative">
                <Input
                  placeholder="Search products by name, SKU, or category..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  leftIcon={<Search className="h-5 w-5" />}
                  size="lg"
                  className="w-full pr-12"
                />
                {searchQuery && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSearchQuery('')}
                    className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0"
                  >
                    ×
                  </Button>
                )}
              </div>
            </div>

            <div className="flex items-center gap-4">
              <Button variant="outline" size="lg" className="border-2 hover:border-primary transition-colors">
                <Filter className="mr-2 h-5 w-5" />
                Advanced Filters
              </Button>

              <Button variant="outline" size="lg" className="border-2 hover:border-primary transition-colors">
                Export
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Enhanced Products Grid */}
      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {[...Array(6)].map((_, i) => (
            <Card key={i} variant="elevated" className="border-0 shadow-large">
              <CardContent className="p-6">
                <div className="animate-pulse">
                  <div className="w-full h-48 bg-muted rounded-2xl mb-4"></div>
                  <div className="space-y-3">
                    <div className="h-6 bg-muted rounded w-3/4"></div>
                    <div className="h-4 bg-muted rounded w-1/2"></div>
                    <div className="h-4 bg-muted rounded w-1/4"></div>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : products.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {products.map((product) => (
            <Card
              key={product.id}
              variant="elevated"
              className="border-0 shadow-large hover:shadow-xl transition-all duration-300 group overflow-hidden"
            >
              <CardContent className="p-0">
                {/* Product Image */}
                <div className="relative h-64 overflow-hidden bg-muted">
                  {product.images?.[0]?.url ? (
                    <Image
                      src={product.images[0].url}
                      alt={product.name}
                      fill
                      className="object-cover group-hover:scale-105 transition-transform duration-300"
                      onError={(e) => {
                        const target = e.target as HTMLImageElement;
                        target.src = '/placeholder-product.svg';
                      }}
                    />
                  ) : (
                    <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-muted">
                      <Plus className="w-16 h-16" />
                    </div>
                  )}

                  {/* Status Badge */}
                  <div className="absolute top-4 left-4">
                    <Badge
                      variant={getStatusVariant(product.status)}
                      className="font-semibold shadow-medium"
                    >
                      {product.status.charAt(0).toUpperCase() + product.status.slice(1)}
                    </Badge>
                  </div>

                  {/* Stock Badge */}
                  <div className="absolute top-4 right-4">
                    <Badge
                      variant={product.stock > 5 ? "default" : product.stock > 0 ? "secondary" : "destructive"}
                      className="font-semibold shadow-medium"
                    >
                      {product.stock} in stock
                    </Badge>
                  </div>
                </div>

                {/* Product Info */}
                <div className="p-6">
                  <div className="mb-4">
                    <h3 className="text-xl font-bold text-foreground mb-2 line-clamp-2 group-hover:text-primary transition-colors">
                      {product.name}
                    </h3>
                    <p className="text-sm text-muted-foreground mb-2">
                      SKU: {product.sku}
                    </p>
                    {product.category && (
                      <Badge variant="outline" className="text-xs">
                        {product.category.name}
                      </Badge>
                    )}
                  </div>

                  {/* Pricing */}
                  <div className="flex items-center justify-between mb-6">
                    <div>
                      <p className="text-2xl font-bold text-primary">
                        {formatPrice(product.price)}
                      </p>
                      {product.cost_price && (
                        <p className="text-sm text-muted-foreground">
                          Cost: {formatPrice(product.cost_price)}
                        </p>
                      )}
                    </div>

                    {product.cost_price && (
                      <div className="text-right">
                        <p className="text-sm font-semibold text-emerald-600">
                          {Math.round(((product.price - product.cost_price) / product.price) * 100)}% margin
                        </p>
                      </div>
                    )}
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleViewProduct(product)}
                      className="flex-1 border-2 hover:border-primary transition-colors"
                    >
                      <Eye className="h-4 w-4 mr-2" />
                      View
                    </Button>

                    <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => handleEditProduct(product)}
                        className="flex-1 border-2 hover:border-primary transition-colors"
                      >
                        <Edit className="h-4 w-4 mr-2" />
                        Edit
                      </Button>
                    </RequirePermission>

                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="outline" size="sm" className="border-2 hover:border-primary transition-colors">
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
                        <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            onClick={() => handleDeleteProduct(product.id)}
                            className="text-destructive focus:text-destructive"
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            Delete
                          </DropdownMenuItem>
                        </RequirePermission>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      ) : (
        <div className="text-center py-24">
          <div className="relative mb-8">
            <div className="w-32 h-32 mx-auto rounded-full bg-gradient-to-br from-muted to-muted/50 flex items-center justify-center shadow-large">
              <Plus className="h-16 w-16 text-muted-foreground" />
            </div>
            <div className="absolute -top-2 -right-2 w-8 h-8 bg-gradient-to-br from-primary to-violet-600 rounded-full flex items-center justify-center shadow-medium">
              <span className="text-white text-sm font-bold">0</span>
            </div>
          </div>

          <h3 className="text-3xl font-bold text-foreground mb-4">
            {searchQuery ? 'No products found' : 'No products yet'}
          </h3>
          <p className="text-xl text-muted-foreground mb-12 max-w-md mx-auto leading-relaxed">
            {searchQuery
              ? `No products found matching "${searchQuery}". Try adjusting your search terms or browse all products.`
              : 'Start building your product catalog by adding your first product.'
            }
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
              <Button
                onClick={() => setShowAddForm(true)}
                size="xl"
                variant="gradient"
              >
                <Plus className="mr-2 h-5 w-5" />
                Add Your First Product
              </Button>
            </RequirePermission>

            {searchQuery && (
              <Button
                onClick={() => setSearchQuery('')}
                size="xl"
                variant="outline"
              >
                Clear Search
              </Button>
            )}
          </div>

          {/* Help section */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-16 max-w-4xl mx-auto">
            <div className="text-center">
              <div className="w-12 h-12 rounded-2xl bg-emerald-100 flex items-center justify-center mx-auto mb-3">
                <span className="text-emerald-600 text-lg">📦</span>
              </div>
              <h4 className="font-semibold text-foreground mb-1">Add Products</h4>
              <p className="text-sm text-muted-foreground">Upload images, set prices, and manage inventory</p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 rounded-2xl bg-blue-100 flex items-center justify-center mx-auto mb-3">
                <span className="text-blue-600 text-lg">🏷️</span>
              </div>
              <h4 className="font-semibold text-foreground mb-1">Organize</h4>
              <p className="text-sm text-muted-foreground">Create categories and manage product variants</p>
            </div>
            <div className="text-center">
              <div className="w-12 h-12 rounded-2xl bg-purple-100 flex items-center justify-center mx-auto mb-3">
                <span className="text-purple-600 text-lg">📊</span>
              </div>
              <h4 className="font-semibold text-foreground mb-1">Track Sales</h4>
              <p className="text-sm text-muted-foreground">Monitor performance and optimize pricing</p>
            </div>
          </div>
        </div>
      )}

      {/* Enhanced Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <Card variant="elevated" className="border-0 shadow-large">
          <CardContent className="p-6">
            <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
              <p className="text-muted-foreground">
                Showing <span className="font-semibold">{((pagination.page - 1) * pagination.limit) + 1}</span> to{' '}
                <span className="font-semibold">{Math.min(pagination.page * pagination.limit, pagination.total)}</span> of{' '}
                <span className="font-semibold">{pagination.total}</span> products
              </p>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="lg"
                  onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                  disabled={!pagination.has_prev}
                  className="border-2 hover:border-primary transition-colors"
                >
                  Previous
                </Button>

                <div className="flex items-center gap-1">
                  {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                    const pageNum = i + 1;
                    return (
                      <Button
                        key={pageNum}
                        variant={pageNum === pagination.page ? "default" : "ghost"}
                        size="lg"
                        onClick={() => setCurrentPage(pageNum)}
                        className="w-12 h-12"
                      >
                        {pageNum}
                      </Button>
                    );
                  })}
                </div>

                <Button
                  variant="outline"
                  size="lg"
                  onClick={() => setCurrentPage(prev => prev + 1)}
                  disabled={!pagination.has_next}
                  className="border-2 hover:border-primary transition-colors"
                >
                  Next
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>
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
                      {selectedProduct.cost_price && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Cost Price:</span>
                          <span className="font-medium">{formatPrice(selectedProduct.cost_price)}</span>
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
                      {selectedProduct.category && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Category:</span>
                          <span className="font-medium">{selectedProduct.category.name}</span>
                        </div>
                      )}
                      <div className="flex justify-between">
                        <span className="text-gray-500">Product Type:</span>
                        <Badge variant="outline" className="text-xs">
                          {selectedProduct.is_digital ? 'Digital' : 'Physical'}
                        </Badge>
                      </div>
                      {selectedProduct.weight && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Weight:</span>
                          <span className="text-gray-700">{selectedProduct.weight} kg</span>
                        </div>
                      )}
                      {selectedProduct.dimensions && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Dimensions:</span>
                          <span className="text-gray-700">
                            {selectedProduct.dimensions.length} × {selectedProduct.dimensions.width} × {selectedProduct.dimensions.height} cm
                          </span>
                        </div>
                      )}
                      {selectedProduct.rating && (
                        <div className="flex justify-between">
                          <span className="text-gray-500">Rating:</span>
                          <span className="text-gray-700">
                            ⭐ {selectedProduct.rating.average.toFixed(1)} ({selectedProduct.rating.count} reviews)
                          </span>
                        </div>
                      )}
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