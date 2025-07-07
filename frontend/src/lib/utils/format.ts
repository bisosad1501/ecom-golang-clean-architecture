// ===== FORMATTING UTILITIES =====

/**
 * Format price with currency and locale support
 */
export function formatPrice(
  price: number,
  options: {
    currency?: string
    locale?: string
    notation?: Intl.NumberFormatOptions['notation']
    minimumFractionDigits?: number
    maximumFractionDigits?: number
  } = {}
) {
  const {
    currency = 'USD',
    locale = 'en-US',
    notation = 'standard',
    minimumFractionDigits = 0,
    maximumFractionDigits = 2,
  } = options

  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency,
    notation,
    minimumFractionDigits,
    maximumFractionDigits,
  }).format(price)
}

/**
 * Format number with locale support
 */
export function formatNumber(
  number: number,
  options: {
    locale?: string
    notation?: Intl.NumberFormatOptions['notation']
    minimumFractionDigits?: number
    maximumFractionDigits?: number
    useGrouping?: boolean
  } = {}
) {
  const {
    locale = 'en-US',
    notation = 'standard',
    minimumFractionDigits = 0,
    maximumFractionDigits = 2,
    useGrouping = true,
  } = options

  return new Intl.NumberFormat(locale, {
    notation,
    minimumFractionDigits,
    maximumFractionDigits,
    useGrouping,
  }).format(number)
}

/**
 * Format percentage
 */
export function formatPercentage(
  value: number,
  options: {
    locale?: string
    minimumFractionDigits?: number
    maximumFractionDigits?: number
  } = {}
) {
  const {
    locale = 'en-US',
    minimumFractionDigits = 0,
    maximumFractionDigits = 1,
  } = options

  return new Intl.NumberFormat(locale, {
    style: 'percent',
    minimumFractionDigits,
    maximumFractionDigits,
  }).format(value / 100)
}

/**
 * Format date with locale and options support
 */
export function formatDate(
  date: Date | string | number,
  options: Intl.DateTimeFormatOptions & { locale?: string } = {}
) {
  const { locale = 'en-US', ...formatOptions } = options
  
  const defaultOptions: Intl.DateTimeFormatOptions = {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    ...formatOptions,
  }

  return new Intl.DateTimeFormat(locale, defaultOptions).format(new Date(date))
}

/**
 * Format time
 */
export function formatTime(
  date: Date | string | number,
  options: {
    locale?: string
    hour12?: boolean
    includeSeconds?: boolean
  } = {}
) {
  const { locale = 'en-US', hour12 = true, includeSeconds = false } = options

  return new Intl.DateTimeFormat(locale, {
    hour: 'numeric',
    minute: '2-digit',
    ...(includeSeconds && { second: '2-digit' }),
    hour12,
  }).format(new Date(date))
}

/**
 * Format datetime
 */
export function formatDateTime(
  date: Date | string | number,
  options: {
    locale?: string
    dateStyle?: 'full' | 'long' | 'medium' | 'short'
    timeStyle?: 'full' | 'long' | 'medium' | 'short'
    hour12?: boolean
  } = {}
) {
  const {
    locale = 'en-US',
    dateStyle = 'medium',
    timeStyle = 'short',
    hour12 = true,
  } = options

  return new Intl.DateTimeFormat(locale, {
    dateStyle,
    timeStyle,
    hour12,
  }).format(new Date(date))
}

/**
 * Format relative time (e.g., "2 hours ago", "in 3 days")
 */
export function formatRelativeTime(
  date: Date | string | number,
  options: {
    locale?: string
    numeric?: 'always' | 'auto'
  } = {}
) {
  const { locale = 'en-US', numeric = 'auto' } = options
  const now = new Date()
  const targetDate = new Date(date)
  const diffInSeconds = Math.floor((targetDate.getTime() - now.getTime()) / 1000)

  const rtf = new Intl.RelativeTimeFormat(locale, { numeric })

  const intervals = [
    { unit: 'year' as const, seconds: 31536000 },
    { unit: 'month' as const, seconds: 2592000 },
    { unit: 'week' as const, seconds: 604800 },
    { unit: 'day' as const, seconds: 86400 },
    { unit: 'hour' as const, seconds: 3600 },
    { unit: 'minute' as const, seconds: 60 },
    { unit: 'second' as const, seconds: 1 },
  ]

  for (const { unit, seconds } of intervals) {
    const interval = Math.floor(Math.abs(diffInSeconds) / seconds)
    if (interval >= 1) {
      return rtf.format(diffInSeconds < 0 ? -interval : interval, unit)
    }
  }

  return rtf.format(0, 'second')
}

/**
 * Format file size
 */
export function formatFileSize(
  bytes: number,
  options: {
    locale?: string
    unit?: 'byte' | 'bit'
    notation?: 'standard' | 'scientific' | 'engineering' | 'compact'
  } = {}
) {
  const { locale = 'en-US', unit = 'byte', notation = 'standard' } = options

  return new Intl.NumberFormat(locale, {
    style: 'unit',
    unit,
    notation,
    unitDisplay: 'short',
  }).format(bytes)
}

/**
 * Format bytes to human readable format
 */
export function formatBytes(
  bytes: number,
  options: {
    decimals?: number
    binary?: boolean
  } = {}
) {
  const { decimals = 2, binary = false } = options

  if (bytes === 0) return '0 Bytes'

  const k = binary ? 1024 : 1000
  const sizes = binary
    ? ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']
    : ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB']

  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(decimals))} ${sizes[i]}`
}

/**
 * Format duration in milliseconds to human readable format
 */
export function formatDuration(
  milliseconds: number,
  options: {
    format?: 'long' | 'short' | 'narrow'
    units?: ('year' | 'month' | 'week' | 'day' | 'hour' | 'minute' | 'second')[]
    maxUnits?: number
  } = {}
) {
  const { format = 'long', units = ['hour', 'minute', 'second'], maxUnits = 2 } = options

  const durations = {
    year: 31536000000,
    month: 2592000000,
    week: 604800000,
    day: 86400000,
    hour: 3600000,
    minute: 60000,
    second: 1000,
  }

  const parts: string[] = []
  let remaining = milliseconds

  for (const unit of units) {
    if (parts.length >= maxUnits) break
    
    const duration = durations[unit]
    const value = Math.floor(remaining / duration)
    
    if (value > 0) {
      const formatter = new Intl.NumberFormat('en-US', {
        style: 'unit',
        unit,
        unitDisplay: format,
      })
      parts.push(formatter.format(value))
      remaining -= value * duration
    }
  }

  return parts.length > 0 ? parts.join(', ') : '0 seconds'
}

/**
 * Format phone number
 */
export function formatPhoneNumber(
  phoneNumber: string,
  options: {
    country?: string
    format?: 'national' | 'international' | 'e164' | 'rfc3966'
  } = {}
) {
  // This is a simple implementation. For production, use a library like libphonenumber-js
  const { format = 'national' } = options
  const cleaned = phoneNumber.replace(/\D/g, '')

  if (cleaned.length === 10) {
    const match = cleaned.match(/^(\d{3})(\d{3})(\d{4})$/)
    if (match) {
      switch (format) {
        case 'international':
          return `+1 (${match[1]}) ${match[2]}-${match[3]}`
        case 'e164':
          return `+1${cleaned}`
        case 'rfc3966':
          return `tel:+1-${match[1]}-${match[2]}-${match[3]}`
        default:
          return `(${match[1]}) ${match[2]}-${match[3]}`
      }
    }
  }

  return phoneNumber
}

/**
 * Format credit card number
 */
export function formatCreditCard(cardNumber: string, mask: boolean = true) {
  const cleaned = cardNumber.replace(/\D/g, '')
  const groups = cleaned.match(/.{1,4}/g) || []
  
  if (mask && groups.length > 1) {
    const maskedGroups = groups.map((group, index) => {
      if (index === groups.length - 1) return group // Last group unmasked
      return '****'
    })
    return maskedGroups.join(' ')
  }
  
  return groups.join(' ')
}

/**
 * Truncate text with ellipsis
 */
export function truncateText(
  text: string,
  options: {
    length?: number
    suffix?: string
    preserveWords?: boolean
  } = {}
) {
  const { length = 100, suffix = '...', preserveWords = true } = options

  if (text.length <= length) return text

  let truncated = text.slice(0, length)

  if (preserveWords) {
    const lastSpace = truncated.lastIndexOf(' ')
    if (lastSpace > 0) {
      truncated = truncated.slice(0, lastSpace)
    }
  }

  return truncated + suffix
}

/**
 * Format initials from name
 */
export function formatInitials(name: string, maxInitials: number = 2) {
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase())
    .slice(0, maxInitials)
    .join('')
}

/**
 * Format slug from text
 */
export function formatSlug(text: string) {
  return text
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '') // Remove special characters
    .replace(/[\s_-]+/g, '-') // Replace spaces and underscores with hyphens
    .replace(/^-+|-+$/g, '') // Remove leading/trailing hyphens
}

/**
 * Format title case
 */
export function formatTitleCase(text: string) {
  return text
    .toLowerCase()
    .split(' ')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

/**
 * Format sentence case
 */
export function formatSentenceCase(text: string) {
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase()
}
