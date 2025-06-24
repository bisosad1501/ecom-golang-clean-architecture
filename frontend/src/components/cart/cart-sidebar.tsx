'use client'

import { Fragment } from 'react'
import Image from 'next/image'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { X, Plus, Minus, ShoppingBag, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { useCartStore, getCartTotal, getCartItemCount } from '@/store/cart'
import { formatPrice } from '@/lib/utils'
import { cn } from '@/lib/utils'

export function CartSidebar() {
  const pathname = usePathname()
  const { 
    cart, 
    isOpen, 
    isLoading, 
    closeCart, 
    updateItem, 
    removeItem 
  } = useCartStore()

  // Don't show cart sidebar on admin pages
  const isAdminPage = pathname.startsWith('/admin')
  if (isAdminPage) {
    return null
  }

  const cartTotal = getCartTotal(cart)
  const cartItemCount = getCartItemCount(cart)

  const handleUpdateQuantity = async (itemId: string, newQuantity: number) => {
    try {
      await updateItem(itemId, newQuantity)
    } catch (error) {
      console.error('Failed to update cart item:', error)
    }
  }

  const handleRemoveItem = async (itemId: string) => {
    try {
      await removeItem(itemId)
    } catch (error) {
      console.error('Failed to remove cart item:', error)
    }
  }

  if (!isOpen) return null

  return (
    <>
      {/* Backdrop */}
      <div 
        className="fixed inset-0 z-50 bg-black bg-opacity-50 transition-opacity"
        onClick={closeCart}
      />

      {/* Sidebar */}
      <div className="fixed right-0 top-0 z-50 h-full w-full max-w-md bg-white shadow-xl">
        <div className="flex h-full flex-col">
          {/* Header */}
          <div className="flex items-center justify-between border-b px-6 py-4">
            <h2 className="text-lg font-semibold">
              Shopping Cart ({cartItemCount})
            </h2>
            <Button
              variant="ghost"
              size="icon"
              onClick={closeCart}
              className="h-8 w-8"
            >
              <X className="h-4 w-4" />
            </Button>
          </div>

          {/* Cart content */}
          <div className="flex-1 overflow-y-auto">
            {isLoading ? (
              <div className="flex items-center justify-center h-32">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              </div>
            ) : !cart || cart.items.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full px-6 text-center">
                <ShoppingBag className="h-16 w-16 text-gray-300 mb-4" />
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  Your cart is empty
                </h3>
                <p className="text-gray-500 mb-6">
                  Add some products to get started
                </p>
                <Button onClick={closeCart} asChild>
                  <Link href="/products">
                    Continue Shopping
                  </Link>
                </Button>
              </div>
            ) : (
              <div className="px-6 py-4">
                <div className="space-y-4">
                  {cart.items.map((item) => (
                    <div key={item.id} className="flex items-center space-x-4">
                      {/* Product image */}
                      <div className="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-md border">
                        <Image
                          src={item.product.images?.[0]?.url || '/placeholder-product.jpg'}
                          alt={item.product.name}
                          fill
                          className="object-cover"
                        />
                      </div>

                      {/* Product details */}
                      <div className="flex-1 min-w-0">
                        <h4 className="text-sm font-medium text-gray-900 truncate">
                          {item.product.name}
                        </h4>
                        <p className="text-sm text-gray-500">
                          {formatPrice(item.unit_price)}
                        </p>
                        
                        {/* Quantity controls */}
                        <div className="flex items-center space-x-2 mt-2">
                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity - 1)}
                            disabled={isLoading || item.quantity <= 1}
                          >
                            <Minus className="h-3 w-3" />
                          </Button>
                          
                          <span className="text-sm font-medium w-8 text-center">
                            {item.quantity}
                          </span>
                          
                          <Button
                            variant="outline"
                            size="icon"
                            className="h-8 w-8"
                            onClick={() => handleUpdateQuantity(item.id, item.quantity + 1)}
                            disabled={isLoading}
                          >
                            <Plus className="h-3 w-3" />
                          </Button>
                        </div>
                      </div>

                      {/* Price and remove */}
                      <div className="flex flex-col items-end space-y-2">
                        <p className="text-sm font-medium text-gray-900">
                          {formatPrice(item.total_price)}
                        </p>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-gray-400 hover:text-red-500"
                          onClick={() => handleRemoveItem(item.id)}
                          disabled={isLoading}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Footer */}
          {cart && cart.items.length > 0 && (
            <div className="border-t bg-gray-50 px-6 py-4">
              {/* Subtotal */}
              <div className="flex items-center justify-between mb-4">
                <span className="text-base font-medium text-gray-900">
                  Subtotal
                </span>
                <span className="text-lg font-semibold text-gray-900">
                  {formatPrice(cartTotal)}
                </span>
              </div>

              {/* Shipping notice */}
              <p className="text-sm text-gray-500 mb-4">
                Shipping and taxes calculated at checkout.
              </p>

              {/* Action buttons */}
              <div className="space-y-2">
                <Button 
                  className="w-full" 
                  size="lg"
                  onClick={closeCart}
                  asChild
                >
                  <Link href="/checkout">
                    Checkout
                  </Link>
                </Button>
                
                <Button 
                  variant="outline" 
                  className="w-full"
                  onClick={closeCart}
                  asChild
                >
                  <Link href="/cart">
                    View Cart
                  </Link>
                </Button>
              </div>

              {/* Continue shopping */}
              <div className="mt-4 text-center">
                <button
                  onClick={closeCart}
                  className="text-sm text-primary-600 hover:text-primary-700 font-medium"
                >
                  Continue Shopping
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  )
}
