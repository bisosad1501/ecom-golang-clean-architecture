'use client'

import { Product } from '@/types'
import { formatPrice } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  Dialog, 
  DialogContent, 
  DialogHeader, 
  DialogTitle,
  DialogFooter 
} from '@/components/ui/dialog'
import { Eye, Edit, Trash2, Copy, Star } from 'lucide-react'
import Image from 'next/image'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'

interface ProductQuickViewProps {
  product: Product | null
  open: boolean
  onOpenChange: (open: boolean) => void
  onEdit?: (product: Product) => void
  onDelete?: (product: Product) => void
  onCopyUrl?: (product: Product) => void
}

export function ProductQuickView({ 
  product, 
  open, 
  onOpenChange, 
  onEdit, 
  onDelete, 
  onCopyUrl 
}: ProductQuickViewProps) {
  if (!product) return null

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3">
            <div className="w-12 h-12 relative overflow-hidden rounded-lg border">
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
                <div className="w-full h-full flex items-center justify-center text-gray-400 bg-gray-100">
                  <Eye className="w-6 h-6" />
                </div>
              )}
            </div>
            <div>
              <h2 className="text-xl font-semibold">{product.name}</h2>
              <p className="text-sm text-gray-500">SKU: {product.sku}</p>
            </div>
          </DialogTitle>
        </DialogHeader>
        
        <div className="space-y-6">
          {/* Product Stats Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-blue-600">Price</p>
                  <p className="text-2xl font-bold text-blue-900">{formatPrice(product.price)}</p>
                  {product.compare_price && product.compare_price > product.price && (
                    <p className="text-sm text-blue-600 line-through">{formatPrice(product.compare_price)}</p>
                  )}
                </div>
                <div className="text-blue-500">
                  <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
                  </svg>
                </div>
              </div>
            </div>
            
            <div className={`border rounded-lg p-4 ${product.stock > 0 ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
              <div className="flex items-center justify-between">
                <div>
                  <p className={`text-sm font-medium ${product.stock > 0 ? 'text-green-600' : 'text-red-600'}`}>Stock</p>
                  <p className={`text-2xl font-bold ${product.stock > 0 ? 'text-green-900' : 'text-red-900'}`}>{product.stock}</p>
                  <p className={`text-sm ${product.stock > 0 ? 'text-green-600' : 'text-red-600'}`}>
                    {product.stock > 0 ? 'In Stock' : 'Out of Stock'}
                  </p>
                </div>
                <div className={product.stock > 0 ? 'text-green-500' : 'text-red-500'}>
                  <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                </div>
              </div>
            </div>
            
            <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-purple-600">Status</p>
                  <div className="mt-1">
                    <Badge variant={product.status === 'active' ? 'default' : 'secondary'}>
                      {product.status}
                    </Badge>
                  </div>
                  {product.compare_price && product.compare_price > product.price && (
                    <p className="text-sm text-purple-600 mt-1">On Sale</p>
                  )}
                </div>
                <div className="text-purple-500">
                  <Star className="w-8 h-8" fill="currentColor" />
                </div>
              </div>
            </div>
          </div>

          {/* Product Details Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div className="space-y-4">
              {/* Description */}
              <div>
                <h3 className="font-medium text-gray-900 mb-2">Description</h3>
                <div className="bg-gray-50 rounded-lg p-4">
                  <p className="text-gray-700 text-sm leading-relaxed">{product.description}</p>
                </div>
              </div>
              
              {/* Short Description */}
              {product.short_description && (
                <div>
                  <h3 className="font-medium text-gray-900 mb-2">Short Description</h3>
                  <div className="bg-gray-50 rounded-lg p-4">
                    <p className="text-gray-700 text-sm">{product.short_description}</p>
                  </div>
                </div>
              )}
              
              {/* Tags */}
              {product.tags && product.tags.length > 0 && (
                <div>
                  <h3 className="font-medium text-gray-900 mb-2">Tags</h3>
                  <div className="flex flex-wrap gap-1">
                    {product.tags.map((tag) => (
                      <Badge key={tag.id} variant="outline" className="text-xs">
                        {tag.name}
                      </Badge>
                    ))}
                  </div>
                </div>
              )}
            </div>
            
            <div className="space-y-4">
              {/* Images */}
              {product.images && product.images.length > 0 && (
                <div>
                  <h3 className="font-medium text-gray-900 mb-2">Images ({product.images.length})</h3>
                  <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                    {product.images.slice(0, 6).map((image, index) => (
                      <div key={index} className="aspect-square relative overflow-hidden rounded-lg border group">
                        <Image
                          src={image.url}
                          alt={image.alt_text || `${product.name} - Image ${index + 1}`}
                          fill
                          className="object-cover transition-transform group-hover:scale-110"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement;
                            target.src = '/placeholder-product.svg';
                          }}
                        />
                        {index === 0 && (
                          <div className="absolute top-1 left-1 bg-blue-500 text-white text-xs px-1.5 py-0.5 rounded">
                            Main
                          </div>
                        )}
                      </div>
                    ))}
                    {product.images.length > 6 && (
                      <div className="aspect-square flex items-center justify-center bg-gray-100 rounded-lg border text-gray-500 text-sm font-medium">
                        +{product.images.length - 6} more
                      </div>
                    )}
                  </div>
                </div>
              )}
              
              {/* Additional Info */}
              <div>
                <h3 className="font-medium text-gray-900 mb-2">Product Information</h3>
                <div className="bg-gray-50 rounded-lg p-4 space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-gray-500">Category:</span>
                    <span className="font-medium">{product.category?.name || 'Uncategorized'}</span>
                  </div>
                  {product.weight && (
                    <div className="flex justify-between">
                      <span className="text-gray-500">Weight:</span>
                      <span className="font-medium">{product.weight}g</span>
                    </div>
                  )}
                  {product.is_digital && (
                    <div className="flex justify-between">
                      <span className="text-gray-500">Type:</span>
                      <Badge variant="outline" className="text-xs">Digital Product</Badge>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <span className="text-gray-500">Created:</span>
                    <span className="font-medium">
                      {product.created_at ? new Date(product.created_at).toLocaleDateString() : 'Unknown'}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <DialogFooter className="flex justify-between">
          <div className="flex gap-2">
            {onCopyUrl && (
              <Button variant="outline" size="sm" onClick={() => onCopyUrl(product)}>
                <Copy className="mr-2 h-4 w-4" />
                Copy URL
              </Button>
            )}
          </div>
          
          <div className="flex gap-2">
            <RequirePermission permission={PERMISSIONS.PRODUCTS_UPDATE}>
              {onEdit && (
                <Button variant="outline" onClick={() => onEdit(product)}>
                  <Edit className="mr-2 h-4 w-4" />
                  Edit
                </Button>
              )}
            </RequirePermission>
            
            <RequirePermission permission={PERMISSIONS.PRODUCTS_DELETE}>
              {onDelete && (
                <Button variant="destructive" onClick={() => onDelete(product)}>
                  <Trash2 className="mr-2 h-4 w-4" />
                  Delete
                </Button>
              )}
            </RequirePermission>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
