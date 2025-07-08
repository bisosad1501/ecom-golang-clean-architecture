'use client'

import { useState } from 'react'
import { Plus, Search, Filter, MoreHorizontal, Edit, Trash2, Eye, Copy, Archive, Package, Grid, List } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
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
import { Badge } from '@/components/ui/badge'
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

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubActionButton,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export function AdminProductsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [showAddForm, setShowAddForm] = useState(false)
  const [deleteProductId, setDeleteProductId] = useState<string | null>(null)
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null)
  const [showProductModal, setShowProductModal] = useState(false)
  const [editProduct, setEditProduct] = useState<Product | null>(null)
  const [showEditModal, setShowEditModal] = useState(false)
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')

  const { data, isLoading, refetch } = useProducts({
    page: currentPage,
    limit: 12,
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
      <div className={BIHUB_ADMIN_THEME.spacing.section}>
        <BiHubPageHeader
          title="Add New Product"
          subtitle="Create a new product in your BiHub catalog"
          breadcrumbs={[
            { label: 'Admin' },
            { label: 'Products' },
            { label: 'Add Product' }
          ]}
          action={
            <Button
              variant="outline"
              onClick={() => setShowAddForm(false)}
              className={BIHUB_ADMIN_THEME.components.button.secondary}
            >
              Back to Products
            </Button>
          }
        />

        <BiHubAdminCard
          title="Product Information"
          subtitle="Fill in the details for your new product"
          icon={<Package className="h-5 w-5 text-white" />}
        >
          <AddProductForm
            onSuccess={() => {
              setShowAddForm(false)
              refetch()
              toast.success('Product created successfully!')
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
        title="Product Catalog"
        subtitle="Manage your BiHub product inventory, pricing, and availability"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Products' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
            <Button
              onClick={() => setShowAddForm(true)}
              className={BIHUB_ADMIN_THEME.components.button.primary}
            >
              <Plus className="mr-2 h-5 w-5" />
              Add New Product
            </Button>
          </RequirePermission>
        }
      />

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className={cn(
          BIHUB_ADMIN_THEME.components.card.base,
          'p-6 text-center'
        )}>
          <div className="text-3xl font-bold text-[#FF9000] mb-2">
            {pagination?.total || 0}
          </div>
          <div className={BIHUB_ADMIN_THEME.typography.body.medium}>
            Total Products
          </div>
        </div>

        <div className={cn(
          BIHUB_ADMIN_THEME.components.card.base,
          'p-6 text-center'
        )}>
          <div className="text-3xl font-bold text-emerald-400 mb-2">
            {products.filter(p => p.stock > 0).length}
          </div>
          <div className={BIHUB_ADMIN_THEME.typography.body.medium}>
            In Stock
          </div>
        </div>

        <div className={cn(
          BIHUB_ADMIN_THEME.components.card.base,
          'p-6 text-center'
        )}>
          <div className="text-3xl font-bold text-yellow-400 mb-2">
            {products.filter(p => p.stock <= 5).length}
          </div>
          <div className={BIHUB_ADMIN_THEME.typography.body.medium}>
            Low Stock
          </div>
        </div>
      </div>

      {/* Enhanced Filters */}
      <BiHubAdminCard
        title="Search & Filter"
        subtitle="Find products quickly with advanced search"
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
            <Button
              variant="outline"
              size="sm"
              className={BIHUB_ADMIN_THEME.components.button.secondary}
            >
              <Filter className="mr-2 h-4 w-4" />
              Filters
            </Button>
          </div>
        }
      >
        <div className="flex flex-col lg:flex-row items-center gap-4">
          <div className="flex-1 w-full">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                placeholder="Search products by name, SKU, or category..."
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

      {/* Products Listing */}
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
              <div className="w-full h-48 bg-gray-700 rounded-xl mb-4"></div>
              <div className="space-y-3">
                <div className="h-6 bg-gray-700 rounded w-3/4"></div>
                <div className="h-4 bg-gray-700 rounded w-1/2"></div>
                <div className="h-4 bg-gray-700 rounded w-1/4"></div>
              </div>
            </div>
          ))}
        </div>
      ) : products.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {products.map((product) => (
            <div
              key={product.id}
              className={cn(
                BIHUB_ADMIN_THEME.components.card.base,
                BIHUB_ADMIN_THEME.components.card.hover,
                'group overflow-hidden',
                viewMode === 'list' && 'flex items-center gap-6 p-6'
              )}
            >
              {/* Product Image */}
              <div className={cn(
                'relative overflow-hidden bg-gray-700 rounded-xl',
                viewMode === 'grid' ? 'h-48 w-full' : 'h-24 w-24 flex-shrink-0'
              )}>
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
                  <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-700">
                    <Package className="w-8 h-8" />
                  </div>
                )}

                {/* Status Badge */}
                {viewMode === 'grid' && (
                  <div className="absolute top-3 left-3">
                    <BiHubStatusBadge status={product.status}>
                      {product.status.charAt(0).toUpperCase() + product.status.slice(1)}
                    </BiHubStatusBadge>
                  </div>
                )}

                {/* Stock Badge */}
                {viewMode === 'grid' && (
                  <div className="absolute top-3 right-3">
                    <BiHubStatusBadge
                      status={product.stock > 5 ? "success" : product.stock > 0 ? "warning" : "error"}
                    >
                      {product.stock} in stock
                    </BiHubStatusBadge>
                  </div>
                )}
              </div>

              {/* Product Info */}
              <div className={cn(
                viewMode === 'grid' ? 'p-6' : 'flex-1',
                viewMode === 'list' && 'flex items-center justify-between'
              )}>
                <div className={viewMode === 'list' ? 'flex-1' : ''}>
                  <div className={cn(viewMode === 'grid' ? 'mb-4' : 'mb-2')}>
                    <h3 className={cn(
                      BIHUB_ADMIN_THEME.typography.heading.h4,
                      'mb-2 line-clamp-2 group-hover:text-[#FF9000] transition-colors',
                      viewMode === 'list' && 'text-lg'
                    )}>
                      {product.name}
                    </h3>
                    <p className={cn(BIHUB_ADMIN_THEME.typography.body.small, 'mb-2')}>
                      SKU: {product.sku}
                    </p>

                    <div className="flex items-center gap-2 flex-wrap">
                      {viewMode === 'list' && (
                        <>
                          <BiHubStatusBadge status={product.status}>
                            {product.status.charAt(0).toUpperCase() + product.status.slice(1)}
                          </BiHubStatusBadge>
                          <BiHubStatusBadge
                            status={product.stock > 5 ? "success" : product.stock > 0 ? "warning" : "error"}
                          >
                            {product.stock} in stock
                          </BiHubStatusBadge>
                        </>
                      )}
                      {product.category && (
                        <span className="px-2 py-1 bg-gray-700/50 text-gray-300 text-xs rounded-md">
                          {product.category.name}
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Pricing */}
                  <div className={cn(
                    'flex items-center justify-between',
                    viewMode === 'grid' ? 'mb-6' : 'mb-2'
                  )}>
                    <div>
                      <p className={cn(
                        'font-bold text-[#FF9000]',
                        viewMode === 'grid' ? 'text-2xl' : 'text-xl'
                      )}>
                        {formatPrice(product.price)}
                      </p>
                      {product.cost_price && (
                        <p className={BIHUB_ADMIN_THEME.typography.body.small}>
                          Cost: {formatPrice(product.cost_price)}
                        </p>
                      )}
                    </div>

                    {product.cost_price && (
                      <div className="text-right">
                        <p className="text-sm font-semibold text-emerald-400">
                          {Math.round(((product.price - product.cost_price) / product.price) * 100)}% margin
                        </p>
                      </div>
                    )}
                  </div>
                </div>

                {/* Actions */}
                <div className={cn(
                  'flex items-center gap-2',
                  viewMode === 'list' && 'flex-shrink-0'
                )}>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleViewProduct(product)}
                    className={cn(
                      BIHUB_ADMIN_THEME.components.button.ghost,
                      viewMode === 'grid' && 'flex-1'
                    )}
                  >
                    <Eye className="h-4 w-4 mr-2" />
                    {viewMode === 'grid' ? 'View' : ''}
                  </Button>

                  <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleEditProduct(product)}
                      className={cn(
                        BIHUB_ADMIN_THEME.components.button.secondary,
                        viewMode === 'grid' && 'flex-1'
                      )}
                    >
                      <Edit className="h-4 w-4 mr-2" />
                      {viewMode === 'grid' ? 'Edit' : ''}
                    </Button>
                  </RequirePermission>

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
                        onClick={() => handleCopyProductUrl(product)}
                        className="text-gray-300 hover:text-white hover:bg-gray-800"
                      >
                        <Copy className="mr-2 h-4 w-4" />
                        Copy URL
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => handleDuplicateProduct(product)}
                        className="text-gray-300 hover:text-white hover:bg-gray-800"
                      >
                        <Copy className="mr-2 h-4 w-4" />
                        Duplicate
                      </DropdownMenuItem>
                      <DropdownMenuSeparator className="bg-gray-700" />
                      <DropdownMenuItem
                        onClick={() => handleArchiveProduct(product)}
                        className="text-gray-300 hover:text-white hover:bg-gray-800"
                      >
                        <Archive className="mr-2 h-4 w-4" />
                        Archive
                      </DropdownMenuItem>
                      <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
                        <DropdownMenuSeparator className="bg-gray-700" />
                        <DropdownMenuItem
                          onClick={() => handleDeleteProduct(product.id)}
                          className="text-red-400 hover:text-red-300 hover:bg-red-900/20"
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
      ) : (
        <BiHubEmptyState
          icon={<Package className="h-8 w-8 text-gray-400" />}
          title={searchQuery ? 'No products found' : 'No products yet'}
          description={
            searchQuery
              ? `No products found matching "${searchQuery}". Try adjusting your search terms or browse all products.`
              : 'Start building your BiHub product catalog by adding your first product.'
          }
          action={
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
                <Button
                  onClick={() => setShowAddForm(true)}
                  className={BIHUB_ADMIN_THEME.components.button.primary}
                >
                  <Plus className="mr-2 h-5 w-5" />
                  Add Your First Product
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

      )}

      {/* Enhanced Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <BiHubAdminCard
          title="Pagination"
          subtitle={`Showing ${((pagination.page - 1) * pagination.limit) + 1} to ${Math.min(pagination.page * pagination.limit, pagination.total)} of ${pagination.total} products`}
        >
          <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
            <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
              Page {pagination.page} of {pagination.total_pages}
            </p>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={!pagination.has_prev}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
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
                      onClick={() => setCurrentPage(pageNum)}
                      className={cn(
                        'w-10 h-10',
                        pageNum === pagination.page
                          ? BIHUB_ADMIN_THEME.components.button.primary
                          : BIHUB_ADMIN_THEME.components.button.ghost
                      )}
                    >
                      {pageNum}
                    </Button>
                  );
                })}
              </div>

              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => prev + 1)}
                disabled={!pagination.has_next}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Next
              </Button>
            </div>
          </div>
        </BiHubAdminCard>
      )}
      
      {/* Delete Confirmation Dialog */}
      <AlertDialog open={!!deleteProductId} onOpenChange={() => setDeleteProductId(null)}>
        <AlertDialogContent className="bg-gray-900 border-gray-700">
          <AlertDialogHeader>
            <AlertDialogTitle className="text-white">Delete Product</AlertDialogTitle>
            <AlertDialogDescription className="text-gray-400">
              Are you sure you want to delete this product? This action cannot be undone.
              All associated data including images, reviews, and sales history will be permanently removed.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel
              disabled={deleteProductMutation.isPending}
              className={BIHUB_ADMIN_THEME.components.button.secondary}
            >
              Cancel
            </AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteProduct}
              disabled={deleteProductMutation.isPending}
              className="bg-red-600 hover:bg-red-700 text-white"
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
          <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto bg-gray-900 border-gray-700">
            <DialogHeader>
              <DialogTitle className="text-white">Edit Product: {editProduct.name}</DialogTitle>
            </DialogHeader>
            <EditProductForm
              product={editProduct}
              onSuccess={() => {
                setShowEditModal(false)
                setEditProduct(null)
                refetch()
                toast.success('Product updated successfully!')
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