'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { CartConflictInfo, MergeStrategy } from '@/types'
import { formatPrice } from '@/lib/utils'
import { ShoppingCart, AlertTriangle, Merge, Replace, UserCheck, Plus } from 'lucide-react'

interface CartConflictModalProps {
  isOpen: boolean
  onClose: () => void
  conflictInfo: CartConflictInfo
  onMerge: (strategy: MergeStrategy) => Promise<void>
  isLoading?: boolean
}

export function CartConflictModal({
  isOpen,
  onClose,
  conflictInfo,
  onMerge,
  isLoading = false
}: CartConflictModalProps) {
  const [selectedStrategy, setSelectedStrategy] = useState<MergeStrategy>('merge')

  const handleMerge = async () => {
    await onMerge(selectedStrategy)
    onClose()
  }

  const strategyOptions = [
    {
      value: 'merge' as MergeStrategy,
      label: conflictInfo.has_conflict ? 'Merge Items' : 'Add to Cart',
      description: conflictInfo.has_conflict
        ? 'Add quantities together for duplicate items'
        : 'Add guest items to your current cart',
      icon: <Plus className="h-5 w-5" />,
      color: 'bg-green-500/20 text-green-300 border-green-500/30'
    },
    {
      value: 'replace' as MergeStrategy,
      label: 'Replace Cart',
      description: 'Replace your current cart with guest cart',
      icon: <Replace className="h-5 w-5" />,
      color: 'bg-orange-500/20 text-orange-300 border-orange-500/30'
    },
    {
      value: 'keep_user' as MergeStrategy,
      label: 'Ignore Guest Items',
      description: 'Keep your current cart, discard guest cart items',
      icon: <UserCheck className="h-5 w-5" />,
      color: 'bg-purple-500/20 text-purple-300 border-purple-500/30'
    }
  ]

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl bg-gradient-to-br from-slate-900/95 via-gray-900/95 to-slate-800/95 border-gray-700/50 text-white">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-3 text-xl">
            <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-orange-500 to-red-500 flex items-center justify-center">
              <AlertTriangle className="h-5 w-5 text-white" />
            </div>
            Cart Merge Conflict Detected
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-6">
          {/* Conflict Summary */}
          <Card className={conflictInfo.has_conflict ? "bg-red-500/10 border-red-500/30" : "bg-blue-500/10 border-blue-500/30"}>
            <CardContent className="p-4">
              <div className="flex items-center gap-3 mb-3">
                <ShoppingCart className={`h-5 w-5 ${conflictInfo.has_conflict ? 'text-red-400' : 'text-blue-400'}`} />
                <h3 className={`font-semibold ${conflictInfo.has_conflict ? 'text-red-300' : 'text-blue-300'}`}>
                  {conflictInfo.has_conflict ? 'Conflicting Items Found' : 'Cart Merge Available'}
                </h3>
              </div>
              <p className={`text-sm mb-4 ${conflictInfo.has_conflict ? 'text-red-200' : 'text-blue-200'}`}>
                {conflictInfo.has_conflict
                  ? 'You have items in both your current cart and guest cart. Choose how to handle the conflicts:'
                  : 'You have items from your browsing session. Choose how to handle them:'
                }
              </p>

              {conflictInfo.conflicting_items && conflictInfo.conflicting_items.length > 0 && (
                <div className="space-y-2">
                  {conflictInfo.conflicting_items.map((item, index) => (
                    <div key={index} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <div>
                        <div className="font-medium text-white">{item.product_name}</div>
                        <div className="text-sm text-gray-400">
                          Current: {item.user_quantity} × {formatPrice(item.user_price)} |
                          Guest: {item.guest_quantity} × {formatPrice(item.guest_price)}
                        </div>
                      </div>
                      {item.price_difference !== 0 && (
                        <Badge className="bg-yellow-500/20 text-yellow-300">
                          Price Diff: {formatPrice(Math.abs(item.price_difference))}
                        </Badge>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Cart Comparison */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Current Cart */}
            <Card className="bg-gray-800/30 border-gray-600/30">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg flex items-center gap-2">
                  <UserCheck className="h-5 w-5 text-blue-400" />
                  Your Current Cart
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-400">Items:</span>
                    <span className="text-white">{conflictInfo.user_cart?.item_count || 0}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">Total:</span>
                    <span className="text-white font-semibold">
                      {formatPrice(conflictInfo.user_cart?.total || 0)}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Guest Cart */}
            <Card className="bg-gray-800/30 border-gray-600/30">
              <CardHeader className="pb-3">
                <CardTitle className="text-lg flex items-center gap-2">
                  <ShoppingCart className="h-5 w-5 text-green-400" />
                  Guest Cart
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-gray-400">Items:</span>
                    <span className="text-white">{conflictInfo.guest_cart?.item_count || 0}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-gray-400">Total:</span>
                    <span className="text-white font-semibold">
                      {formatPrice(conflictInfo.guest_cart?.total || 0)}
                    </span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Strategy Selection */}
          <div>
            <h3 className="font-semibold text-white mb-4">Choose Merge Strategy:</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              {strategyOptions.map((option) => (
                <Card
                  key={option.value}
                  className={`cursor-pointer transition-all duration-200 ${
                    selectedStrategy === option.value
                      ? 'ring-2 ring-[#ff9000] bg-[#ff9000]/10 border-[#ff9000]/30'
                      : 'bg-gray-800/30 border-gray-600/30 hover:border-gray-500/50'
                  }`}
                  onClick={() => setSelectedStrategy(option.value)}
                >
                  <CardContent className="p-4">
                    <div className="flex items-start gap-3">
                      <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${option.color}`}>
                        {option.icon}
                      </div>
                      <div className="flex-1">
                        <div className="font-semibold text-white mb-1">{option.label}</div>
                        <div className="text-sm text-gray-400">{option.description}</div>
                      </div>
                      {selectedStrategy === option.value && (
                        <div className="w-5 h-5 rounded-full bg-[#ff9000] flex items-center justify-center">
                          <div className="w-2 h-2 rounded-full bg-white"></div>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          </div>

          {/* Recommendations */}
          {conflictInfo.recommendations && conflictInfo.recommendations.length > 0 && (
            <Card className="bg-blue-500/10 border-blue-500/30">
              <CardContent className="p-4">
                <h4 className="font-semibold text-blue-300 mb-2">Recommendations:</h4>
                <ul className="text-sm text-blue-200 space-y-1">
                  {conflictInfo.recommendations.map((rec, index) => (
                    <li key={index}>• {rec}</li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}

          {/* Actions */}
          <div className="flex gap-4 pt-4">
            <Button
              onClick={handleMerge}
              disabled={isLoading}
              className="flex-1 bg-[#ff9000] hover:bg-[#e68100] text-white"
            >
              {isLoading ? 'Merging...' : `Merge with ${strategyOptions.find(o => o.value === selectedStrategy)?.label}`}
            </Button>
            <Button
              onClick={onClose}
              variant="outline"
              disabled={isLoading}
              className="border-gray-600 text-gray-300 hover:bg-gray-800"
            >
              Cancel
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
