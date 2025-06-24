/**
 * Price utilities for consistent price handling across the application
 * These utilities mirror the backend business logic for pricing
 */

import { Product } from '@/types'

export interface PriceInfo {
  price: number
  comparePrice?: number
  costPrice?: number
  hasDiscount: boolean
  discountPercentage: number
  isAvailable: boolean
}

/**
 * Check if a product has a discount
 * Mirrors backend logic: ComparePrice != nil && *ComparePrice > Price
 */
export function hasDiscount(price: number, comparePrice?: number): boolean {
  return comparePrice !== undefined && comparePrice !== null && comparePrice > price
}

/**
 * Calculate discount percentage
 * Mirrors backend logic: (ComparePrice - Price) / ComparePrice * 100
 */
export function getDiscountPercentage(price: number, comparePrice?: number): number {
  if (!hasDiscount(price, comparePrice) || !comparePrice) {
    return 0
  }
  return ((comparePrice - price) / comparePrice) * 100
}

/**
 * Check if a product is available
 * Mirrors backend logic: Status == "active" && Stock > 0
 */
export function isProductAvailable(status: string, stock: number): boolean {
  return status === 'active' && stock > 0
}

/**
 * Get comprehensive price information for a product
 */
export function getPriceInfo(product: Product): PriceInfo {
  return {
    price: product.price,
    comparePrice: product.compare_price,
    costPrice: product.cost_price,
    hasDiscount: hasDiscount(product.price, product.compare_price),
    discountPercentage: getDiscountPercentage(product.price, product.compare_price),
    isAvailable: isProductAvailable(product.status, product.stock),
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
 * Validate price inputs
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
