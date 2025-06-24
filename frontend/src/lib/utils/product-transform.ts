/**
 * Data transformation utility để đảm bảo consistency giữa FE và BE
 */

import { UpdateProductRequest } from '@/lib/services/products'

export function transformUpdateProductData(formData: any): UpdateProductRequest {
  const transformed: any = {}
  
  // Only include fields that have values (not undefined/null/empty)
  if (formData.name && formData.name.trim()) {
    transformed.name = formData.name.trim()
  }
  
  if (formData.description && formData.description.trim()) {
    transformed.description = formData.description.trim()
  }
  
  if (formData.price !== undefined && formData.price !== null) {
    transformed.price = Number(formData.price)
  }
  
  if (formData.compare_price !== undefined && formData.compare_price !== null && formData.compare_price > 0) {
    transformed.compare_price = Number(formData.compare_price)
  }
  
  if (formData.cost_price !== undefined && formData.cost_price !== null && formData.cost_price > 0) {
    transformed.cost_price = Number(formData.cost_price)
  }
  
  if (formData.stock !== undefined && formData.stock !== null) {
    transformed.stock = Number(formData.stock)
  }
  
  if (formData.weight !== undefined && formData.weight !== null && formData.weight > 0) {
    transformed.weight = Number(formData.weight)
  }
  
  if (formData.category_id && formData.category_id.length > 0) {
    console.log('=== transformUpdateProductData: Processing category_id ===')
    console.log('formData.category_id:', formData.category_id)
    console.log('formData.category_id type:', typeof formData.category_id)
    console.log('formData.category_id length:', formData.category_id.length)
    transformed.category_id = formData.category_id
    console.log('transformed.category_id set to:', transformed.category_id)
  } else {
    console.log('=== transformUpdateProductData: category_id NOT included ===')
    console.log('formData.category_id:', formData.category_id)
    console.log('formData.category_id type:', typeof formData.category_id)
    if (formData.category_id) {
      console.log('formData.category_id length:', formData.category_id.length)
    }
  }
  
  if (formData.status) {
    transformed.status = formData.status
  }
  
  if (formData.is_digital !== undefined) {
    transformed.is_digital = Boolean(formData.is_digital)
  }
  
  if (formData.images && Array.isArray(formData.images)) {
    transformed.images = formData.images
  }

  if (formData.image_changes) {
    transformed.image_changes = formData.image_changes
  }

  if (formData.tags && Array.isArray(formData.tags)) {
    transformed.tags = formData.tags
  }

  console.log('transformUpdateProductData:', { original: formData, transformed })
  return transformed
}
