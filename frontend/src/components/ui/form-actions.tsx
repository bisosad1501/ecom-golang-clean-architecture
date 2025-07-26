'use client'

import React from 'react'
import { Button } from '@/components/ui/button'
import { Loader2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

interface FormActionsProps {
  onCancel?: () => void
  onSubmit?: () => void
  submitLabel?: string
  cancelLabel?: string
  isSubmitting?: boolean
  disabled?: boolean
  className?: string
  submitVariant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
}

export function FormActions({
  onCancel,
  onSubmit,
  submitLabel = 'Save',
  cancelLabel = 'Cancel',
  isSubmitting = false,
  disabled = false,
  className,
  submitVariant = 'default',
}: FormActionsProps) {
  return (
    <div className={cn(
      'flex justify-end space-x-4 pt-6 border-t border-gray-700/50 mt-8',
      className
    )}>
      {onCancel && (
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting || disabled}
          className={cn(
            BIHUB_ADMIN_THEME.components.button.secondary,
            "px-6 py-2"
          )}
        >
          {cancelLabel}
        </Button>
      )}
      <Button
        type={onSubmit ? 'button' : 'submit'}
        variant={submitVariant}
        onClick={onSubmit}
        disabled={isSubmitting || disabled}
        className={cn(
          BIHUB_ADMIN_THEME.components.button.primary,
          "min-w-[120px] px-6 py-2"
        )}
      >
        {isSubmitting ? (
          <>
            <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            Saving...
          </>
        ) : (
          submitLabel
        )}
      </Button>
    </div>
  )
}
