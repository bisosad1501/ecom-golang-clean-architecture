/**
 * Price utilities for consistent price handling across the application
 * Updated to work with backend ProductResponse structure
 */

import { Product } from '@/types'

export interface PriceInfo {
  // Current selling price (backend computed)
  currentPrice: number
  // Original price
  price: number
  // Compare price for discount calculation
  comparePrice?: number
  // Cost price
  costPrice?: number
  // Sale price if on sale
  salePrice?: number
  // Backend computed discount flags
  hasDiscount: boolean
  isOnSale: boolean
  // Backend computed discount percentage
  discountPercentage: number
  // Availability
  isAvailable: boolean
  // Stock info
  stock: number
  stockStatus: string
  isLowStock: boolean
}

/**
 * Check if a product has a discount (using backend computed field)
 */
export function hasDiscount(product: Product): boolean {
  return product.has_discount || product.is_on_sale
}

/**
 * Get discount percentage (using backend computed field)
 */
export function getDiscountPercentage(product: Product): number {
  return product.sale_discount_percentage || 0
}

/**
 * Check if a product is available (using backend computed field)
 */
export function isProductAvailable(product: Product): boolean {
  return product.is_available
}

/**
 * Get the current selling price (using backend computed field)
 */
export function getCurrentPrice(product: Product): number {
  return product.current_price || product.price
}

/**
 * Get the display price for comparison (original price when on sale)
 */
export function getComparePrice(product: Product): number | undefined {
  if (product.is_on_sale && product.price !== product.current_price) {
    return product.price
  }
  return product.compare_price
}

/**
 * Get comprehensive price information for a product
 * Uses backend computed fields for accuracy
 */
export function getPriceInfo(product: Product): PriceInfo {
  return {
    currentPrice: getCurrentPrice(product),
    price: product.price,
    comparePrice: product.compare_price,
    costPrice: product.cost_price,
    salePrice: product.sale_price,
    hasDiscount: product.has_discount,
    isOnSale: product.is_on_sale,
    discountPercentage: product.sale_discount_percentage,
    isAvailable: product.is_available,
    stock: product.stock,
    stockStatus: product.stock_status,
    isLowStock: product.is_low_stock,
  }
}

/**
 * Format price for display
 */
export function formatPrice(price: number, currency = '$'): string {
  return `${currency}${price.toFixed(2)}`
}

/**
 * Format discount percentage for display
 */
export function formatDiscountPercentage(percentage: number): string {
  return `${Math.round(percentage)}% OFF`
}

/**
 * Clean price input (convert 0 to undefined for optional price fields)
 */
export function cleanPriceInput(value: number | undefined): number | undefined {
  if (value === 0 || value === null) {
    return undefined
  }
  return value
}

/**
 * Validate price inputs (for admin forms)
 */
export function validatePriceInputs(price: number, comparePrice?: number, costPrice?: number): string[] {
  const errors: string[] = []

  if (price <= 0) {
    errors.push('Price must be greater than 0')
  }

  if (comparePrice !== undefined && comparePrice <= 0) {
    errors.push('Compare price must be greater than 0')
  }

  if (costPrice !== undefined && costPrice <= 0) {
    errors.push('Cost price must be greater than 0')
  }

  if (comparePrice !== undefined && comparePrice <= price) {
    errors.push('Compare price should be greater than price for discount')
  }

  if (costPrice !== undefined && costPrice >= price) {
    errors.push('Cost price should be less than selling price for profit')
  }

  return errors
}

/**
 * Legacy function for backward compatibility
 * @deprecated Use getPriceInfo(product).hasDiscount instead
 */
export function hasDiscountLegacy(price: number, comparePrice?: number): boolean {
  return comparePrice !== undefined && comparePrice !== null && comparePrice > price
}

/**
 * Legacy function for backward compatibility
 * @deprecated Use getPriceInfo(product).discountPercentage instead
 */
export function getDiscountPercentageLegacy(price: number, comparePrice?: number): number {
  if (!hasDiscountLegacy(price, comparePrice) || !comparePrice) {
    return 0
  }
  return ((comparePrice - price) / comparePrice) * 100
}
