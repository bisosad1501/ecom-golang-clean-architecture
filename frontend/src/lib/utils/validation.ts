// ===== VALIDATION UTILITIES =====

/**
 * Email validation
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Phone number validation (basic)
 */
export function isValidPhone(phone: string): boolean {
  const phoneRegex = /^\+?[\d\s\-\(\)]+$/
  return phoneRegex.test(phone) && phone.replace(/\D/g, '').length >= 10
}

/**
 * Password strength validation
 */
export function validatePassword(password: string): {
  isValid: boolean
  score: number
  feedback: string[]
  requirements: {
    minLength: boolean
    hasUppercase: boolean
    hasLowercase: boolean
    hasNumber: boolean
    hasSpecialChar: boolean
  }
} {
  const requirements = {
    minLength: password.length >= 8,
    hasUppercase: /[A-Z]/.test(password),
    hasLowercase: /[a-z]/.test(password),
    hasNumber: /\d/.test(password),
    hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password),
  }

  const score = Object.values(requirements).filter(Boolean).length
  const feedback: string[] = []

  if (!requirements.minLength) {
    feedback.push('Password must be at least 8 characters long')
  }
  if (!requirements.hasUppercase) {
    feedback.push('Password must contain at least one uppercase letter')
  }
  if (!requirements.hasLowercase) {
    feedback.push('Password must contain at least one lowercase letter')
  }
  if (!requirements.hasNumber) {
    feedback.push('Password must contain at least one number')
  }
  if (!requirements.hasSpecialChar) {
    feedback.push('Password must contain at least one special character')
  }

  return {
    isValid: score === 5,
    score,
    feedback,
    requirements,
  }
}

/**
 * URL validation
 */
export function isValidUrl(url: string): boolean {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

/**
 * Credit card number validation (Luhn algorithm)
 */
export function isValidCreditCard(cardNumber: string): boolean {
  const cleaned = cardNumber.replace(/\D/g, '')
  
  if (cleaned.length < 13 || cleaned.length > 19) {
    return false
  }

  let sum = 0
  let isEven = false

  for (let i = cleaned.length - 1; i >= 0; i--) {
    let digit = parseInt(cleaned[i])

    if (isEven) {
      digit *= 2
      if (digit > 9) {
        digit -= 9
      }
    }

    sum += digit
    isEven = !isEven
  }

  return sum % 10 === 0
}

/**
 * Credit card type detection
 */
export function getCreditCardType(cardNumber: string): string {
  const cleaned = cardNumber.replace(/\D/g, '')
  
  const patterns = {
    visa: /^4/,
    mastercard: /^5[1-5]|^2[2-7]/,
    amex: /^3[47]/,
    discover: /^6(?:011|5)/,
    dinersclub: /^3[068]/,
    jcb: /^35/,
  }

  for (const [type, pattern] of Object.entries(patterns)) {
    if (pattern.test(cleaned)) {
      return type
    }
  }

  return 'unknown'
}

/**
 * Postal code validation by country
 */
export function isValidPostalCode(postalCode: string, country: string = 'US'): boolean {
  const patterns: Record<string, RegExp> = {
    US: /^\d{5}(-\d{4})?$/,
    CA: /^[A-Za-z]\d[A-Za-z] \d[A-Za-z]\d$/,
    UK: /^[A-Za-z]{1,2}\d[A-Za-z\d]? \d[A-Za-z]{2}$/,
    DE: /^\d{5}$/,
    FR: /^\d{5}$/,
    JP: /^\d{3}-\d{4}$/,
    AU: /^\d{4}$/,
  }

  const pattern = patterns[country.toUpperCase()]
  return pattern ? pattern.test(postalCode) : true
}

/**
 * Social Security Number validation (US)
 */
export function isValidSSN(ssn: string): boolean {
  const cleaned = ssn.replace(/\D/g, '')
  
  if (cleaned.length !== 9) {
    return false
  }

  // Check for invalid patterns
  const invalidPatterns = [
    /^000/, // First 3 digits cannot be 000
    /^666/, // First 3 digits cannot be 666
    /^9/, // First digit cannot be 9
    /^\d{3}00/, // Middle 2 digits cannot be 00
    /^\d{5}0000$/, // Last 4 digits cannot be 0000
  ]

  return !invalidPatterns.some(pattern => pattern.test(cleaned))
}

/**
 * Tax ID validation (EIN - Employer Identification Number)
 */
export function isValidEIN(ein: string): boolean {
  const cleaned = ein.replace(/\D/g, '')
  
  if (cleaned.length !== 9) {
    return false
  }

  // Valid prefixes for EIN
  const validPrefixes = [
    '01', '02', '03', '04', '05', '06', '10', '11', '12', '13', '14', '15', '16',
    '20', '21', '22', '23', '24', '25', '26', '27', '30', '31', '32', '33', '34',
    '35', '36', '37', '38', '39', '40', '41', '42', '43', '44', '45', '46', '47',
    '48', '50', '51', '52', '53', '54', '55', '56', '57', '58', '59', '60', '61',
    '62', '63', '64', '65', '66', '67', '68', '71', '72', '73', '74', '75', '76',
    '77', '80', '81', '82', '83', '84', '85', '86', '87', '88', '90', '91', '92',
    '93', '94', '95', '98', '99'
  ]

  const prefix = cleaned.substring(0, 2)
  return validPrefixes.includes(prefix)
}

/**
 * File type validation
 */
export function isValidFileType(file: File, allowedTypes: string[]): boolean {
  return allowedTypes.includes(file.type)
}

/**
 * File size validation
 */
export function isValidFileSize(file: File, maxSizeInBytes: number): boolean {
  return file.size <= maxSizeInBytes
}

/**
 * Image file validation
 */
export function isValidImageFile(file: File): boolean {
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  return isValidFileType(file, allowedTypes)
}

/**
 * Date validation
 */
export function isValidDate(date: string | Date): boolean {
  const dateObj = new Date(date)
  return dateObj instanceof Date && !isNaN(dateObj.getTime())
}

/**
 * Age validation
 */
export function isValidAge(birthDate: string | Date, minAge: number = 0, maxAge: number = 150): boolean {
  if (!isValidDate(birthDate)) {
    return false
  }

  const today = new Date()
  const birth = new Date(birthDate)
  const age = today.getFullYear() - birth.getFullYear()
  const monthDiff = today.getMonth() - birth.getMonth()

  const actualAge = monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())
    ? age - 1
    : age

  return actualAge >= minAge && actualAge <= maxAge
}

/**
 * Username validation
 */
export function isValidUsername(username: string): boolean {
  // Username must be 3-20 characters, alphanumeric and underscores only
  const usernameRegex = /^[a-zA-Z0-9_]{3,20}$/
  return usernameRegex.test(username)
}

/**
 * Slug validation
 */
export function isValidSlug(slug: string): boolean {
  const slugRegex = /^[a-z0-9]+(?:-[a-z0-9]+)*$/
  return slugRegex.test(slug)
}

/**
 * Hex color validation
 */
export function isValidHexColor(color: string): boolean {
  const hexRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/
  return hexRegex.test(color)
}

/**
 * IP address validation
 */
export function isValidIPAddress(ip: string): boolean {
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  const ipv6Regex = /^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/
  
  return ipv4Regex.test(ip) || ipv6Regex.test(ip)
}

/**
 * MAC address validation
 */
export function isValidMACAddress(mac: string): boolean {
  const macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/
  return macRegex.test(mac)
}

/**
 * JSON validation
 */
export function isValidJSON(jsonString: string): boolean {
  try {
    JSON.parse(jsonString)
    return true
  } catch {
    return false
  }
}

/**
 * Required field validation
 */
export function isRequired(value: any): boolean {
  if (value === null || value === undefined) {
    return false
  }
  
  if (typeof value === 'string') {
    return value.trim().length > 0
  }
  
  if (Array.isArray(value)) {
    return value.length > 0
  }
  
  return true
}

/**
 * Minimum length validation
 */
export function hasMinLength(value: string, minLength: number): boolean {
  return value.length >= minLength
}

/**
 * Maximum length validation
 */
export function hasMaxLength(value: string, maxLength: number): boolean {
  return value.length <= maxLength
}

/**
 * Range validation for numbers
 */
export function isInRange(value: number, min: number, max: number): boolean {
  return value >= min && value <= max
}

/**
 * Pattern validation
 */
export function matchesPattern(value: string, pattern: RegExp): boolean {
  return pattern.test(value)
}

/**
 * Custom validation function type
 */
export type ValidationFunction<T = any> = (value: T) => boolean | string

/**
 * Validation rule interface
 */
export interface ValidationRule<T = any> {
  validator: ValidationFunction<T>
  message: string
}

/**
 * Validate value against multiple rules
 */
export function validateValue<T>(value: T, rules: ValidationRule<T>[]): {
  isValid: boolean
  errors: string[]
} {
  const errors: string[] = []
  
  for (const rule of rules) {
    const result = rule.validator(value)
    if (result === false || typeof result === 'string') {
      errors.push(typeof result === 'string' ? result : rule.message)
    }
  }
  
  return {
    isValid: errors.length === 0,
    errors,
  }
}
