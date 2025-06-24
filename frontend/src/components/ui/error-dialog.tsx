'use client'

import { AlertTriangle, Info, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { CategoryDeleteError } from '@/lib/utils/category-delete-errors'

interface ErrorDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  error: CategoryDeleteError
  onRetry?: () => void
}

export function ErrorDialog({ open, onOpenChange, error, onRetry }: ErrorDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <div className="flex items-center space-x-3">
            <div className="flex-shrink-0">
              <AlertTriangle className="h-6 w-6 text-red-600" />
            </div>
            <DialogTitle className="text-lg font-semibold text-gray-900">
              {error.title}
            </DialogTitle>
          </div>
        </DialogHeader>
        
        <div className="space-y-4">
          {/* Error Message */}
          <p className="text-sm text-gray-600 leading-relaxed">
            {error.message}
          </p>

          {/* Suggestions */}
          {error.suggestions && error.suggestions.length > 0 && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div className="flex items-start space-x-2">
                <Info className="h-4 w-4 text-blue-600 mt-0.5 flex-shrink-0" />
                <div>
                  <h4 className="text-sm font-medium text-blue-900 mb-2">
                    Suggestions:
                  </h4>
                  <ul className="text-sm text-blue-800 space-y-1">
                    {error.suggestions.map((suggestion, index) => (
                      <li key={index} className="flex items-start">
                        <span className="mr-2">•</span>
                        <span>{suggestion}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex items-center justify-end space-x-3 pt-4 border-t">
            {onRetry && (
              <Button
                variant="outline"
                onClick={() => {
                  onOpenChange(false)
                  onRetry()
                }}
              >
                Try Again
              </Button>
            )}
            <Button
              onClick={() => onOpenChange(false)}
              className="min-w-[80px]"
            >
              Close
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}
