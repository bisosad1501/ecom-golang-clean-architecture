import { ApiError } from '@/types'

export interface CategoryDeleteError {
  title: string
  message: string
  suggestions?: string[]
}

export function getCategoryDeleteErrorMessage(error: ApiError): CategoryDeleteError {
  // Check if error contains specific details about why deletion failed
  const errorMessage = error.message?.toLowerCase() || ''
  const errorDetails = error.details || []

  // Products dependency
  if (
    errorMessage.includes('product') || 
    errorMessage.includes('foreign key') ||
    errorMessage.includes('constraint') ||
    errorDetails.some((detail: any) => 
      detail.includes?.('product') || detail.includes?.('reference')
    )
  ) {
    return {
      title: 'Cannot Delete Category',
      message: 'This category cannot be deleted because it is currently being used by one or more products.',
      suggestions: [
        'Move products to another category first',
        'Delete all products in this category',
        'Or disable the category instead of deleting it'
      ]
    }
  }

  // Subcategories dependency
  if (
    errorMessage.includes('subcategory') || 
    errorMessage.includes('child') ||
    errorMessage.includes('parent') ||
    errorDetails.some((detail: any) => 
      detail.includes?.('subcategory') || detail.includes?.('child')
    )
  ) {
    return {
      title: 'Cannot Delete Parent Category',
      message: 'This category cannot be deleted because it has subcategories.',
      suggestions: [
        'Delete all subcategories first',
        'Move subcategories to another parent category',
        'Or disable the category instead of deleting it'
      ]
    }
  }

  // Orders dependency
  if (
    errorMessage.includes('order') || 
    errorMessage.includes('purchase') ||
    errorDetails.some((detail: any) => 
      detail.includes?.('order') || detail.includes?.('purchase')
    )
  ) {
    return {
      title: 'Cannot Delete Category',
      message: 'This category cannot be deleted because it is referenced in existing orders.',
      suggestions: [
        'Categories used in orders cannot be deleted for data integrity',
        'Consider disabling the category instead',
        'Contact administrator if deletion is absolutely necessary'
      ]
    }
  }

  // Reviews dependency
  if (
    errorMessage.includes('review') || 
    errorMessage.includes('rating') ||
    errorDetails.some((detail: any) => 
      detail.includes?.('review') || detail.includes?.('rating')
    )
  ) {
    return {
      title: 'Cannot Delete Category',
      message: 'This category cannot be deleted because it has product reviews associated with it.',
      suggestions: [
        'Move products with reviews to another category',
        'Or disable the category instead of deleting it'
      ]
    }
  }

  // Permission denied
  if (error.code === 'FORBIDDEN' || error.code === 'UNAUTHORIZED') {
    return {
      title: 'Permission Denied',
      message: 'You do not have permission to delete this category.',
      suggestions: [
        'Contact your administrator for proper permissions',
        'Ensure you are logged in with the correct account'
      ]
    }
  }

  // Category not found
  if (error.code === 'NOT_FOUND') {
    return {
      title: 'Category Not Found',
      message: 'The category you are trying to delete no longer exists.',
      suggestions: [
        'The category may have already been deleted',
        'Refresh the page to see the current list'
      ]
    }
  }

  // Validation error
  if (error.code === 'VALIDATION_ERROR' || error.code === 'UNPROCESSABLE_ENTITY') {
    return {
      title: 'Invalid Request',
      message: error.message || 'The delete request is invalid.',
      suggestions: [
        'Ensure the category ID is valid',
        'Try refreshing the page and attempting again'
      ]
    }
  }

  // Conflict error
  if (error.code === 'CONFLICT') {
    return {
      title: 'Deletion Conflict',
      message: error.message || 'This category is currently being used and cannot be deleted.',
      suggestions: [
        'Check if the category is being used by products, orders, or other data',
        'Remove all dependencies before attempting to delete',
        'Consider disabling the category instead'
      ]
    }
  }

  // Network error
  if (error.code === 'NETWORK_ERROR') {
    return {
      title: 'Connection Error',
      message: 'Unable to connect to the server. Please check your internet connection.',
      suggestions: [
        'Check your internet connection',
        'Try again in a few moments',
        'Contact support if the problem persists'
      ]
    }
  }

  // Generic server error
  return {
    title: 'Server Error',
    message: error.message || 'An unexpected error occurred while deleting the category.',
    suggestions: [
      'This category might be in use by other data',
      'Try again in a few moments',
      'Contact support if the problem persists',
      'Consider disabling the category instead of deleting it'
    ]
  }
}

export function getSubcategoryDeleteErrorMessage(error: ApiError): CategoryDeleteError {
  const baseError = getCategoryDeleteErrorMessage(error)
  
  // Customize messages for subcategories
  if (baseError.title === 'Cannot Delete Category') {
    return {
      ...baseError,
      title: 'Cannot Delete Subcategory',
      message: baseError.message.replace('category', 'subcategory')
    }
  }
  
  return baseError
}
