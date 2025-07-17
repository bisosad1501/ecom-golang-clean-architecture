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
import { Product } from '@/types/product'
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

      {/* Enhanced Quick Stats with Modern Design */}
      <div className="space-y-6">
        {/* Primary Stats Row - Modern Glass Design */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          {/* Total Products */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-blue-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-blue-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-indigo-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-blue-400/80 text-sm font-medium uppercase tracking-wide">Total Products</p>
                <p className="text-2xl font-bold text-blue-100 mt-1">{pagination?.total || 0}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Package className="h-5 w-5 text-blue-400" />
              </div>
            </div>
          </div>

          {/* In Stock */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-emerald-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-emerald-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-emerald-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-green-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-emerald-400/80 text-sm font-medium uppercase tracking-wide">In Stock</p>
                <p className="text-2xl font-bold text-emerald-100 mt-1">
                  {products.filter((p: Product) => {
                    const stockStatus = (p as any).stock_status || 'in_stock'
                    const stock = (p as any).stock ?? 0
                    return stockStatus === 'in_stock' && stock > 0
                  }).length}
                </p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Package className="h-5 w-5 text-emerald-400" />
              </div>
            </div>
          </div>

          {/* Low Stock */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-amber-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-amber-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-amber-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 to-orange-500/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-amber-400/80 text-sm font-medium uppercase tracking-wide">Low Stock</p>
                <p className="text-2xl font-bold text-amber-100 mt-1">
                  {products.filter((p: Product) => (p as any).is_low_stock === true).length}
                </p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Package className="h-5 w-5 text-amber-400" />
              </div>
            </div>
          </div>

          {/* On Sale */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-orange-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-orange-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-orange-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-orange-500/5 to-red-500/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-orange-400/80 text-sm font-medium uppercase tracking-wide">On Sale</p>
                <p className="text-2xl font-bold text-orange-100 mt-1">
                  {products.filter((p: Product) => (p as any).is_on_sale === true).length}
                </p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-orange-500/20 to-red-500/20 border border-orange-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Package className="h-5 w-5 text-orange-400" />
              </div>
            </div>
          </div>
        </div>
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
                <h3 className="text-lg font-semibold text-white">Search & Filter Products</h3>
                <p className="text-sm text-gray-400">Find and manage your BiHub product catalog</p>
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

          {/* Search and Filter Controls */}
          <div className="flex flex-col lg:flex-row items-center gap-4">
            <div className="flex-1 w-full">
              <div className="relative group">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 group-focus-within:text-blue-400 transition-colors" />
                <Input
                  placeholder="Search products by name, SKU, or category..."
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
                    Clear
                  </Button>
                )}
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Button className="group relative bg-white/5 border border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white px-4 py-2 h-12 rounded-xl transition-all duration-200">
                <Filter className="mr-2 h-4 w-4" />
                <span className="font-medium">Advanced Filters</span>
                <div className="ml-2 w-2 h-2 rounded-full bg-blue-400 animate-pulse opacity-60"></div>
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Products List */}
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
      ) : products.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {products.map((product) => (
            <div
              key={product.id}
              className={cn(
                'group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200 shadow-lg hover:shadow-xl',
                viewMode === 'grid' 
                  ? 'p-5 hover:scale-[1.02]' 
                  : 'p-4 flex items-center gap-6'
              )}
            >
              {/* Gradient Background */}
              <div className={cn(
                'absolute inset-0 bg-gradient-to-br opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-2xl',
                product.status === 'active' ? 'from-emerald-500/5 to-green-600/5' :
                product.status === 'draft' ? 'from-amber-500/5 to-orange-500/5' :
                'from-red-500/5 to-rose-600/5'
              )} />
              {/* Product Header */}
              <div className={cn(
                'relative flex items-center gap-4',
                viewMode === 'grid' 
                  ? 'justify-between mb-4' 
                  : 'flex-shrink-0 w-80'
              )}>
                <div className="flex items-center gap-4">
                  {/* Product Image or Icon */}
                  <div className={cn(
                    'relative rounded-xl overflow-hidden shadow-lg border border-white/10',
                    viewMode === 'grid' ? 'w-12 h-12' : 'w-16 h-16',
                    product.status === 'active' ? 'border-emerald-500/30' :
                    product.status === 'draft' ? 'border-amber-500/30' :
                    'border-red-500/30'
                  )}>
                    {(product as any).images?.[0]?.url ? (
                      <Image
                        src={(product as any).images[0].url}
                        alt={product.name}
                        fill
                        className="object-cover"
                        onError={(e) => {
                          const target = e.target as HTMLImageElement;
                          target.src = '/placeholder-product.svg';
                        }}
                      />
                    ) : (
                      <div className={cn(
                        'w-full h-full bg-gradient-to-br flex items-center justify-center',
                        product.status === 'active' ? 'from-emerald-500/20 to-green-600/20' :
                        product.status === 'draft' ? 'from-amber-500/20 to-orange-500/20' :
                        'from-red-500/20 to-rose-600/20'
                      )}>
                        <Package className="h-5 w-5 text-white" />
                        <div className="absolute inset-0 bg-white/10 rounded-xl blur-sm"></div>
                      </div>
                    )}
                  </div>

                  <div>
                    <h3 className={cn(
                      "font-bold text-white group-hover:text-[#FF9000] transition-colors",
                      viewMode === 'grid' ? 'text-lg' : 'text-base'
                    )}>
                      {product.name}
                    </h3>
                    {/* Modern Status Badge */}
                    <div className={cn(
                      'inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border border-white/10 backdrop-blur-sm',
                      viewMode === 'grid' ? 'mt-2' : 'mt-1',
                      product.status === 'active' ? 'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-800' :
                      product.status === 'draft' ? 'bg-amber-100 text-amber-800 border-amber-200 dark:bg-amber-950/30 dark:text-amber-300 dark:border-amber-800' :
                      'bg-red-100 text-red-800 border-red-200 dark:bg-red-950/30 dark:text-red-300 dark:border-red-800'
                    )}>
                      <div className={cn(
                        'w-2 h-2 rounded-full mr-2',
                        product.status === 'active' ? 'bg-emerald-500' :
                        product.status === 'draft' ? 'bg-amber-500' :
                        'bg-red-500'
                      )} />
                      {product.status.charAt(0).toUpperCase() + product.status.slice(1)}
                    </div>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <div className="flex flex-col items-end gap-1">
                      <p className="text-2xl font-bold text-[#FF9000] group-hover:text-[#FF9000]/80 transition-colors">
                        {formatPrice((product as any).current_price || (product as any).price || 0)}
                      </p>
                      {(product as any).has_discount && (product as any).original_price && (
                        <div className="flex items-center gap-2">
                          <p className="text-sm text-gray-400 line-through">
                            {formatPrice((product as any).original_price || 0)}
                          </p>
                          <Badge className="bg-[#FF9000]/20 text-[#FF9000] text-xs px-1.5 py-0.5">
                            -{Math.round((product as any).discount_percentage || 0)}%
                          </Badge>
                        </div>
                      )}
                    </div>
                    <p className="text-xs text-gray-400 mt-1">
                      {(product as any).has_discount ? 'Current Price' : 'Price'}
                    </p>
                  </div>
                )}
              </div>

              {/* Product Details - Simplified */}
              <div className={cn(
                'relative',
                viewMode === 'grid' 
                  ? 'space-y-3' 
                  : 'flex-1 flex items-center gap-8'
              )}>
                {/* Grid view content */}
                {viewMode === 'grid' && (
                  <>
                    {/* SKU */}
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                        <Package className="h-4 w-4 text-gray-400" />
                      </div>
                      <div>
                        <p className="text-xs text-gray-400 uppercase tracking-wide">SKU</p>
                        <p className="text-sm text-white font-medium">
                          {product.sku || 'N/A'}
                        </p>
                      </div>
                    </div>

                    {/* Enhanced Stock for grid view */}
                    <div className="flex items-center gap-3">
                      <div className={cn(
                        "w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10",
                        (product as any).stock_status === 'in_stock' ? "text-emerald-400" :
                        (product as any).stock_status === 'out_of_stock' ? "text-red-400" : "text-amber-400"
                      )}>
                        <Package className="h-4 w-4" />
                      </div>
                      <div>
                        <p className="text-xs text-gray-400 uppercase tracking-wide">Stock</p>
                        <div className="flex items-center gap-2">
                          <p className={cn(
                            "text-sm font-medium",
                            (product as any).stock_status === 'in_stock' ? "text-emerald-400" :
                            (product as any).stock_status === 'out_of_stock' ? "text-red-400" : "text-amber-400"
                          )}>
                            {(product as any).stock ?? 0} items
                          </p>
                          {(product as any).is_low_stock && (
                            <Badge variant="outline" className="text-xs px-1.5 py-0 bg-amber-500/10 text-amber-400 border-amber-500/30">
                              Low Stock
                            </Badge>
                          )}
                          {(product as any).is_on_sale && (
                            <Badge variant="outline" className="text-xs px-1.5 py-0 bg-[#FF9000]/10 text-[#FF9000] border-[#FF9000]/30">
                              Sale
                            </Badge>
                          )}
                          {(product as any).featured && (
                            <Badge variant="outline" className="text-xs px-1.5 py-0 bg-purple-500/10 text-purple-400 border-purple-500/30">
                              Featured
                            </Badge>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Category for grid view */}
                    {product.category && (
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                          <Package className="h-4 w-4 text-purple-400" />
                        </div>
                        <div>
                          <p className="text-xs text-gray-400 uppercase tracking-wide">Category</p>
                          <p className="text-sm text-white font-medium">
                            {product.category.name}
                          </p>
                        </div>
                      </div>
                    )}
                  </>
                )}

                {/* List view layout - Horizontal */}
                {viewMode === 'list' && (
                  <>
                    {/* SKU Column */}
                    <div className="flex flex-col min-w-0 w-32">
                      <p className="text-xs text-gray-400 uppercase tracking-wide font-medium mb-1">SKU</p>
                      <p className="text-sm text-white font-medium truncate">
                        {product.sku || 'N/A'}
                      </p>
                    </div>

                    {/* Enhanced Price Column */}
                    <div className="flex flex-col min-w-0 w-32">
                      <p className="text-xs text-gray-400 uppercase tracking-wide font-medium mb-1">PRICE</p>
                      <div className="flex flex-col gap-1">
                        <p className="text-sm font-bold text-[#FF9000] truncate">
                          {formatPrice((product as any).current_price || (product as any).price || 0)}
                        </p>
                        {(product as any).has_discount && (product as any).original_price && (
                          <div className="flex items-center gap-1">
                            <p className="text-xs text-gray-400 line-through">
                              {formatPrice((product as any).original_price || 0)}
                            </p>
                            <Badge className="bg-[#FF9000]/20 text-[#FF9000] text-xs px-1 py-0">
                              -{Math.round((product as any).discount_percentage || 0)}%
                            </Badge>
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Enhanced Stock Column */}
                    <div className="flex flex-col min-w-0 w-32">
                      <p className="text-xs text-gray-400 uppercase tracking-wide font-medium mb-1">STOCK</p>
                      <div className="flex flex-col gap-1">
                        <div className="flex items-center gap-2">
                          <p className={cn(
                            "text-sm font-medium truncate",
                            (product as any).stock_status === 'in_stock' ? "text-emerald-400" :
                            (product as any).stock_status === 'out_of_stock' ? "text-red-400" : "text-amber-400"
                          )}>
                            {(product as any).stock ?? 0}
                          </p>
                          <Badge variant="outline" className={cn(
                            "text-xs px-1.5 py-0 flex-shrink-0",
                            (product as any).stock_status === 'in_stock' ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/30" :
                            (product as any).stock_status === 'out_of_stock' ? "bg-red-500/10 text-red-400 border-red-500/30" :
                            "bg-amber-500/10 text-amber-400 border-amber-500/30"
                          )}>
                            {(product as any).stock_status === 'in_stock' ? 'In Stock' :
                             (product as any).stock_status === 'out_of_stock' ? 'Out' :
                             (product as any).stock_status === 'on_backorder' ? 'Backorder' : 'Unknown'}
                          </Badge>
                        </div>
                        <div className="flex items-center gap-1">
                          {(product as any).is_low_stock && (
                            <Badge variant="outline" className="text-xs px-1 py-0 bg-amber-500/10 text-amber-400 border-amber-500/30">
                              Low
                            </Badge>
                          )}
                          {(product as any).featured && (
                            <Badge variant="outline" className="text-xs px-1 py-0 bg-purple-500/10 text-purple-400 border-purple-500/30">
                              Featured
                            </Badge>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Category Column */}
                    {product.category && (
                      <div className="flex flex-col min-w-0 w-32">
                        <p className="text-xs text-gray-400 uppercase tracking-wide font-medium mb-1">CATEGORY</p>
                        <p className="text-sm text-white font-medium truncate">
                          {product.category.name}
                        </p>
                      </div>
                    )}
                  </>
                )}
              </div>

              {/* Modern Action Buttons */}
              <div className={cn(
                'relative flex items-center gap-3 mt-6 pt-4 border-t border-white/10',
                viewMode === 'list' && 'flex-shrink-0 mt-0 pt-0 border-t-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleViewProduct(product)}
                  className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-blue-500/10 hover:border-blue-500/50 hover:text-blue-400 transition-all duration-200"
                >
                  <Eye className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'View Details' : 'View'}
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
                      onClick={() => handleViewProduct(product)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Eye className="mr-3 h-4 w-4" />
                      View Details
                    </DropdownMenuItem>
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>

                    <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                      <DropdownMenuItem
                        onClick={() => handleEditProduct(product)}
                        className="text-blue-400 hover:text-blue-300 hover:bg-blue-900/20 rounded-lg m-1 p-3"
                      >
                        <Edit className="mr-3 h-4 w-4" />
                        Edit Product
                      </DropdownMenuItem>
                    </RequirePermission>

                    <DropdownMenuItem
                      onClick={() => handleCopyProductUrl(product)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Copy className="mr-3 h-4 w-4" />
                      Copy URL
                    </DropdownMenuItem>

                    <DropdownMenuItem
                      onClick={() => handleDuplicateProduct(product)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Copy className="mr-3 h-4 w-4" />
                      Duplicate
                    </DropdownMenuItem>

                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    <DropdownMenuItem
                      onClick={() => handleArchiveProduct(product)}
                      className="text-amber-400 hover:text-amber-300 hover:bg-amber-900/20 rounded-lg m-1 p-3"
                    >
                      <Archive className="mr-3 h-4 w-4" />
                      Archive
                    </DropdownMenuItem>

                    <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
                      <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                      <DropdownMenuItem
                        onClick={() => handleDeleteProduct(product.id)}
                        className="text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded-lg m-1 p-3"
                      >
                        <Trash2 className="mr-3 h-4 w-4" />
                        Delete Product
                      </DropdownMenuItem>
                    </RequirePermission>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-12 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-slate-500/5 rounded-2xl"></div>
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-500/20 to-slate-500/20 border border-gray-400/30 flex items-center justify-center mx-auto mb-6">
              <Package className="h-8 w-8 text-gray-400" />
            </div>
            <h3 className="text-xl font-bold text-white mb-2">
              {searchQuery ? 'No products found' : 'No products yet'}
            </h3>
            <p className="text-gray-400 mb-6 max-w-md mx-auto">
              {searchQuery
                ? `No products found matching "${searchQuery}". Try adjusting your search terms.`
                : 'Start building your BiHub product catalog by adding your first product.'
              }
            </p>
            <div className="flex flex-col sm:flex-row gap-4 justify-center">
              <RequirePermission permission={PERMISSIONS.PRODUCTS_CREATE}>
                <Button
                  onClick={() => setShowAddForm(true)}
                  className="bg-blue-500/20 border border-blue-400/30 text-blue-400 hover:bg-blue-500/30 hover:border-blue-400/50 hover:text-blue-300 transition-all duration-200"
                >
                  <Plus className="mr-2 h-5 w-5" />
                  Add Your First Product
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
          <AlertDialogContent className="max-w-6xl max-h-[95vh] overflow-hidden bg-gray-900 border border-gray-700 shadow-2xl">
            <AlertDialogHeader className="border-b border-gray-700 pb-6 bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900">
              <AlertDialogTitle className="flex items-start gap-6">
                <div className="w-24 h-24 relative overflow-hidden rounded-2xl border-2 border-gray-600 shadow-lg bg-gray-800">
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
                    <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gradient-to-br from-gray-800 to-gray-700">
                      <Package className="w-10 h-10" />
                    </div>
                  )}
                </div>
                <div className="flex-1 space-y-3">
                  <div>
                    <h2 className="text-3xl font-bold text-white mb-2">{selectedProduct.name}</h2>
                    <div className="flex items-center gap-4 flex-wrap">
                      <span className="font-mono bg-gray-800 border border-gray-600 px-3 py-1.5 rounded-lg text-sm font-medium text-gray-300">
                        SKU: {selectedProduct.sku}
                      </span>
                      <Badge
                        variant={selectedProduct.status === 'active' ? 'default' : 'secondary'}
                        className={`px-3 py-1 text-sm font-medium ${
                          selectedProduct.status === 'active'
                            ? 'bg-emerald-900/30 text-emerald-400 border-emerald-500/30'
                            : 'bg-gray-700/50 text-gray-300 border-gray-600'
                        }`}
                      >
                        {selectedProduct.status?.toUpperCase()}
                      </Badge>
                      <Badge
                        variant={(selectedProduct as any).is_available ? 'default' : 'destructive'}
                        className={`px-3 py-1 text-sm font-medium ${
                          (selectedProduct as any).is_available
                            ? 'bg-emerald-900/30 text-emerald-400 border-emerald-500/30'
                            : 'bg-red-900/30 text-red-400 border-red-500/30'
                        }`}
                      >
                        {(selectedProduct as any).is_available ? 'Available' : 'Out of Stock'}
                      </Badge>
                    </div>
                  </div>
                  <div className="flex items-baseline gap-3">
                    <span className="text-4xl font-bold text-[#FF9000]">
                      {formatPrice((selectedProduct as any).price)}
                    </span>
                    {(selectedProduct as any).compare_price && (selectedProduct as any).compare_price > (selectedProduct as any).price && (
                      <span className="text-xl text-gray-400 line-through">
                        {formatPrice((selectedProduct as any).compare_price)}
                      </span>
                    )}
                  </div>
                </div>
              </AlertDialogTitle>
            </AlertDialogHeader>
            
            <div className="overflow-y-auto max-h-[60vh] px-6 py-4">
              {/* Product Details Grid */}
              <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
                {/* Left Column - Product Information */}
                <div className="lg:col-span-2 space-y-6">
                  {/* Pricing Information */}
                  <div className="bg-gradient-to-r from-gray-800 to-gray-700 rounded-xl p-6 border border-gray-600">
                    <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                      <span className="w-2 h-2 bg-[#FF9000] rounded-full"></span>
                      Pricing Information
                    </h3>
                    <div className="grid grid-cols-2 gap-4">
                      <div className="bg-gray-900 rounded-lg p-4 border border-gray-600">
                        <div className="text-sm text-gray-400 mb-1">Current Price</div>
                        <div className="text-2xl font-bold text-[#FF9000]">{formatPrice((selectedProduct as any).current_price)}</div>
                      </div>
                      {(selectedProduct as any).original_price && (
                        <div className="bg-gray-900 rounded-lg p-4 border border-gray-600">
                          <div className="text-sm text-gray-400 mb-1">Original Price</div>
                          <div className="text-xl font-semibold text-gray-500 line-through">{formatPrice((selectedProduct as any).original_price)}</div>
                        </div>
                      )}
                      {(selectedProduct as any).cost_price && (
                        <div className="bg-gray-900 rounded-lg p-4 border border-gray-600">
                          <div className="text-sm text-gray-400 mb-1">Cost Price</div>
                          <div className="text-xl font-semibold text-gray-300">{formatPrice((selectedProduct as any).cost_price)}</div>
                        </div>
                      )}
                      <div className="bg-gray-900 rounded-lg p-4 border border-gray-600">
                        <div className="text-sm text-gray-400 mb-1">Stock Quantity</div>
                        <div className={`text-xl font-semibold ${((selectedProduct as any).stock ?? 0) > 0 ? 'text-emerald-400' : 'text-red-400'}`}>
                          {(selectedProduct as any).stock ?? 0} items
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Product Details */}
                  <div className="bg-gray-800 rounded-xl p-6 border border-gray-600">
                    <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                      <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                      Product Details
                    </h3>
                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div className="space-y-3">
                        <div className="flex justify-between py-2 border-b border-gray-600">
                          <span className="text-gray-400 font-medium">Category:</span>
                          <span className="text-gray-200">{selectedProduct.category?.name || 'N/A'}</span>
                        </div>
                        <div className="flex justify-between py-2 border-b border-gray-600">
                          <span className="text-gray-400 font-medium">Product Type:</span>
                          <Badge variant="outline" className="text-xs bg-blue-900/30 text-blue-400 border-blue-500/30">
                            {(selectedProduct as any).is_digital ? 'Digital' : 'Physical'}
                          </Badge>
                        </div>
                        {(selectedProduct as any).weight && (
                          <div className="flex justify-between py-2 border-b border-gray-600">
                            <span className="text-gray-400 font-medium">Weight:</span>
                            <span className="text-gray-200">{(selectedProduct as any).weight} kg</span>
                          </div>
                        )}
                      </div>
                      <div className="space-y-3">
                        <div className="flex justify-between py-2 border-b border-gray-600">
                          <span className="text-gray-400 font-medium">Status:</span>
                          <Badge
                            variant={selectedProduct.status === 'active' ? 'default' : 'secondary'}
                            className={selectedProduct.status === 'active'
                              ? 'bg-emerald-900/30 text-emerald-400 border-emerald-500/30'
                              : 'bg-gray-700/50 text-gray-300 border-gray-600'
                            }
                          >
                            {selectedProduct.status?.toUpperCase()}
                          </Badge>
                        </div>
                        <div className="flex justify-between py-2 border-b border-gray-600">
                          <span className="text-gray-400 font-medium">Availability:</span>
                          <Badge
                            variant={(selectedProduct as any).is_available ? 'default' : 'destructive'}
                            className={`text-xs ${
                              (selectedProduct as any).is_available
                                ? 'bg-emerald-900/30 text-emerald-400 border-emerald-500/30'
                                : 'bg-red-900/30 text-red-400 border-red-500/30'
                            }`}
                          >
                            {(selectedProduct as any).is_available ? 'Available' : 'Not Available'}
                          </Badge>
                        </div>
                        {(selectedProduct as any).has_discount && (
                          <div className="flex justify-between py-2 border-b border-gray-600">
                            <span className="text-gray-400 font-medium">Has Discount:</span>
                            <Badge variant="secondary" className="text-xs bg-purple-900/30 text-purple-400 border-purple-500/30">Yes</Badge>
                          </div>
                        )}
                      </div>
                    </div>

                    {(selectedProduct as any).dimensions && (
                      <div className="mt-4 p-3 bg-gray-900 rounded-lg border border-gray-600">
                        <div className="text-sm text-gray-400 font-medium mb-1">Dimensions (L × W × H)</div>
                        <div className="text-lg font-semibold text-gray-200">
                          {(selectedProduct as any).dimensions.length} × {(selectedProduct as any).dimensions.width} × {(selectedProduct as any).dimensions.height} cm
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Description */}
                  <div className="bg-gray-800 rounded-xl p-6 border border-gray-600">
                    <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                      <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                      Description
                    </h3>
                    <p className="text-gray-300 text-sm leading-relaxed">{selectedProduct.description}</p>

                    {(selectedProduct as any).short_description && (
                      <div className="mt-4 pt-4 border-t border-gray-600">
                        <h4 className="font-medium text-white mb-2">Short Description</h4>
                        <p className="text-gray-400 text-sm">{(selectedProduct as any).short_description}</p>
                      </div>
                    )}
                  </div>

                  {/* Tags */}
                  {(selectedProduct as any).tags && (selectedProduct as any).tags.length > 0 && (
                    <div className="bg-gray-800 rounded-xl p-6 border border-gray-600">
                      <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                        <span className="w-2 h-2 bg-purple-500 rounded-full"></span>
                        Tags
                      </h3>
                      <div className="flex flex-wrap gap-2">
                        {(selectedProduct as any).tags.map((tag: any, index: number) => (
                          <Badge key={tag.id || index} variant="outline" className="text-sm px-3 py-1 bg-purple-900/30 border-purple-500/30 text-purple-400">
                            {tag.name || tag}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                </div>

                {/* Right Column - Images and Additional Info */}
                <div className="space-y-6">
                  {/* Images Gallery */}
                  {(selectedProduct as any).images && (selectedProduct as any).images.length > 0 && (
                    <div className="bg-gray-800 rounded-xl p-6 border border-gray-600">
                      <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                        <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                        Images ({(selectedProduct as any).images.length})
                      </h3>
                      <div className="grid grid-cols-2 gap-3">
                        {(selectedProduct as any).images.slice(0, 4).map((image: any, index: number) => (
                          <div key={index} className="aspect-square relative overflow-hidden rounded-lg border-2 border-gray-600 shadow-sm">
                            <Image
                              src={image.url}
                              alt={image.alt_text || `${selectedProduct.name} - Image ${index + 1}`}
                              fill
                              className="object-cover hover:scale-105 transition-transform duration-200"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.src = '/placeholder-product.svg';
                              }}
                            />
                          </div>
                        ))}
                        {(selectedProduct as any).images.length > 4 && (
                          <div className="aspect-square flex items-center justify-center bg-gradient-to-br from-gray-700 to-gray-600 rounded-lg border-2 border-gray-600 text-gray-300 text-sm font-medium">
                            +{(selectedProduct as any).images.length - 4} more
                          </div>
                        )}
                      </div>
                    </div>
                  )}

                  {/* Additional Information */}
                  <div className="bg-gray-800 rounded-xl p-6 border border-gray-600">
                    <h3 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                      <span className="w-2 h-2 bg-gray-400 rounded-full"></span>
                      Additional Information
                    </h3>
                    <div className="space-y-3 text-sm">
                      <div className="flex justify-between py-2 border-b border-gray-600">
                        <span className="text-gray-400 font-medium">Created:</span>
                        <span className="text-gray-200">{new Date((selectedProduct as any).created_at).toLocaleDateString()}</span>
                      </div>
                      <div className="flex justify-between py-2 border-b border-gray-600">
                        <span className="text-gray-400 font-medium">Updated:</span>
                        <span className="text-gray-200">{new Date((selectedProduct as any).updated_at).toLocaleDateString()}</span>
                      </div>
                      {/* Category ID removed - Backend uses ProductCategory many-to-many */}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <AlertDialogFooter className="border-t border-gray-700 pt-6 bg-gray-900">
              <div className="flex items-center justify-between w-full">
                <div className="text-sm text-gray-400">
                  Last updated: {new Date((selectedProduct as any).updated_at).toLocaleDateString()}
                </div>
                <div className="flex gap-3">
                  <AlertDialogCancel className="px-6 py-2 bg-gray-800 border-gray-600 text-gray-300 hover:bg-gray-700 hover:border-gray-500 hover:text-white transition-colors">
                    Close
                  </AlertDialogCancel>
                  <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
                    <AlertDialogAction
                      className="px-6 py-2 bg-[#FF9000] hover:bg-[#FF9000]/90 text-white border-[#FF9000] shadow-lg hover:shadow-xl transition-all duration-200"
                      onClick={() => {
                        setShowProductModal(false)
                        handleEditProduct(selectedProduct)
                      }}
                    >
                      <Edit className="w-4 h-4 mr-2" />
                      Edit Product
                    </AlertDialogAction>
                  </RequirePermission>
                </div>
              </div>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      )}

      {/* Edit Product Modal */}
      {editProduct && (
        <Dialog open={showEditModal} onOpenChange={setShowEditModal}>
          <DialogContent className="max-w-6xl max-h-[95vh] overflow-hidden bg-gray-900 border border-gray-700 shadow-2xl">
            <DialogHeader className="border-b border-gray-700 pb-6 bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900">
              <DialogTitle className="flex items-center gap-4">
                <div className="w-12 h-12 bg-gradient-to-br from-[#FF9000] to-orange-600 rounded-xl flex items-center justify-center shadow-lg">
                  <Edit className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h2 className="text-2xl font-bold text-white">Edit Product</h2>
                  <p className="text-sm text-gray-400 mt-1">{editProduct.name}</p>
                </div>
              </DialogTitle>
            </DialogHeader>
            <div className="overflow-y-auto max-h-[75vh] px-6 py-4">
              <div className="bg-gradient-to-br from-gray-800 to-gray-900 rounded-xl p-6 border border-gray-700">
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
              </div>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </div>
  )
}