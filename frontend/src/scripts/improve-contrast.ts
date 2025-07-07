/**
 * Script to improve contrast across all pages
 * Applies high contrast design system to all components
 */

export const CONTRAST_IMPROVEMENTS = {
  // Common text color improvements
  TEXT_REPLACEMENTS: [
    { from: 'text-gray-300', to: 'text-gray-200' },
    { from: 'text-gray-400', to: 'text-gray-300' },
    { from: 'text-gray-500', to: 'text-gray-400' },
    { from: 'text-muted-foreground', to: 'text-gray-200' },
    { from: 'text-foreground', to: 'text-white' },
  ],

  // Background improvements
  BG_REPLACEMENTS: [
    { from: 'bg-gray-800', to: 'bg-gray-900' },
    { from: 'bg-gray-700', to: 'bg-gray-800' },
    { from: 'bg-muted', to: 'bg-gray-900' },
    { from: 'bg-background', to: 'bg-black' },
  ],

  // Border improvements
  BORDER_REPLACEMENTS: [
    { from: 'border-gray-800', to: 'border-gray-600' },
    { from: 'border-gray-700', to: 'border-gray-600' },
    { from: 'border-muted', to: 'border-gray-600' },
  ],

  // Button improvements
  BUTTON_IMPROVEMENTS: {
    primary: 'bg-orange-500 hover:bg-orange-400 text-white border-orange-500',
    secondary: 'bg-gray-700 hover:bg-gray-600 text-white border-gray-600',
    outline: 'bg-transparent hover:bg-gray-800 text-white border-gray-500 hover:border-gray-400',
  },

  // Form improvements
  FORM_IMPROVEMENTS: {
    input: 'bg-gray-800 border-gray-600 text-white placeholder:text-gray-400 focus:border-orange-500',
    label: 'text-gray-200',
    error: 'text-red-400',
  },

  // Card improvements
  CARD_IMPROVEMENTS: {
    default: 'bg-gray-900 border-gray-600 text-white',
    elevated: 'bg-gray-800 border-gray-600 text-white',
    interactive: 'bg-gray-900 hover:bg-gray-800 border-gray-600 hover:border-gray-500 text-white',
  },
}

// Pages that need contrast improvements
export const PAGES_TO_IMPROVE = [
  'frontend/src/components/pages/search-page.tsx',
  'frontend/src/components/pages/categories-page.tsx',
  'frontend/src/components/pages/product-detail-page.tsx',
  'frontend/src/app/profile/page.tsx',
  'frontend/src/app/admin/dashboard/page.tsx',
  'frontend/src/components/admin/**/*.tsx',
  'frontend/src/components/profile/**/*.tsx',
]

// Component files that need improvements
export const COMPONENTS_TO_IMPROVE = [
  'frontend/src/components/ui/input.tsx',
  'frontend/src/components/ui/button.tsx',
  'frontend/src/components/ui/card.tsx',
  'frontend/src/components/ui/badge.tsx',
  'frontend/src/components/navigation/**/*.tsx',
  'frontend/src/components/products/**/*.tsx',
]

export function applyContrastImprovements(content: string): string {
  let improvedContent = content

  // Apply text improvements
  CONTRAST_IMPROVEMENTS.TEXT_REPLACEMENTS.forEach(({ from, to }) => {
    improvedContent = improvedContent.replace(new RegExp(from, 'g'), to)
  })

  // Apply background improvements
  CONTRAST_IMPROVEMENTS.BG_REPLACEMENTS.forEach(({ from, to }) => {
    improvedContent = improvedContent.replace(new RegExp(from, 'g'), to)
  })

  // Apply border improvements
  CONTRAST_IMPROVEMENTS.BORDER_REPLACEMENTS.forEach(({ from, to }) => {
    improvedContent = improvedContent.replace(new RegExp(from, 'g'), to)
  })

  return improvedContent
}

// Specific improvements for different page types
export const PAGE_SPECIFIC_IMPROVEMENTS = {
  // Product pages
  PRODUCT_PAGES: {
    cardBg: 'bg-gray-900 border-gray-600 hover:border-gray-500',
    priceText: 'text-orange-400',
    titleText: 'text-white',
    descriptionText: 'text-gray-200',
  },

  // Auth pages
  AUTH_PAGES: {
    formBg: 'bg-gray-900 border-gray-700',
    inputBg: 'bg-gray-800 border-gray-600 text-white placeholder:text-gray-400',
    labelText: 'text-gray-200',
    linkText: 'text-orange-400 hover:text-orange-300',
  },

  // Admin pages
  ADMIN_PAGES: {
    tableBg: 'bg-gray-900 border-gray-700',
    tableHeaderBg: 'bg-gray-800',
    tableText: 'text-white',
    sidebarBg: 'bg-gray-950 border-gray-800',
  },

  // Profile pages
  PROFILE_PAGES: {
    sectionBg: 'bg-gray-900 border-gray-700',
    cardBg: 'bg-gray-800 border-gray-600',
    text: 'text-white',
    mutedText: 'text-gray-200',
  },
}

// High contrast utility classes
export const HIGH_CONTRAST_CLASSES = {
  // Text utilities
  TEXT: {
    primary: 'text-white',
    secondary: 'text-gray-200',
    muted: 'text-gray-300',
    brand: 'text-orange-400',
    success: 'text-green-400',
    warning: 'text-yellow-400',
    error: 'text-red-400',
  },

  // Background utilities
  BG: {
    primary: 'bg-black',
    secondary: 'bg-gray-950',
    card: 'bg-gray-900',
    elevated: 'bg-gray-800',
    muted: 'bg-gray-950/50',
  },

  // Border utilities
  BORDER: {
    primary: 'border-gray-600',
    secondary: 'border-gray-700',
    accent: 'border-orange-500',
    muted: 'border-gray-800',
  },

  // Interactive states
  INTERACTIVE: {
    hover: 'hover:bg-gray-800',
    active: 'active:bg-gray-700',
    focus: 'focus:border-orange-500 focus:ring-2 focus:ring-orange-500/20',
    disabled: 'disabled:bg-gray-800 disabled:text-gray-500',
  },
}

// Function to generate high contrast component classes
export function getHighContrastComponentClasses(component: string, variant?: string) {
  const baseClasses = {
    button: {
      primary: `${HIGH_CONTRAST_CLASSES.BG.primary} ${HIGH_CONTRAST_CLASSES.TEXT.primary} border-orange-500 hover:bg-orange-400`,
      secondary: `${HIGH_CONTRAST_CLASSES.BG.elevated} ${HIGH_CONTRAST_CLASSES.TEXT.primary} ${HIGH_CONTRAST_CLASSES.BORDER.primary} hover:bg-gray-700`,
      outline: `bg-transparent ${HIGH_CONTRAST_CLASSES.TEXT.primary} ${HIGH_CONTRAST_CLASSES.BORDER.primary} hover:bg-gray-800`,
    },
    card: {
      default: `${HIGH_CONTRAST_CLASSES.BG.card} ${HIGH_CONTRAST_CLASSES.BORDER.primary} ${HIGH_CONTRAST_CLASSES.TEXT.primary}`,
      elevated: `${HIGH_CONTRAST_CLASSES.BG.elevated} ${HIGH_CONTRAST_CLASSES.BORDER.primary} ${HIGH_CONTRAST_CLASSES.TEXT.primary}`,
      interactive: `${HIGH_CONTRAST_CLASSES.BG.card} ${HIGH_CONTRAST_CLASSES.BORDER.primary} ${HIGH_CONTRAST_CLASSES.TEXT.primary} hover:border-gray-500 transition-colors`,
    },
    input: {
      default: `${HIGH_CONTRAST_CLASSES.BG.elevated} ${HIGH_CONTRAST_CLASSES.BORDER.primary} ${HIGH_CONTRAST_CLASSES.TEXT.primary} placeholder:text-gray-400 focus:border-orange-500`,
    },
    badge: {
      default: `${HIGH_CONTRAST_CLASSES.BG.elevated} ${HIGH_CONTRAST_CLASSES.TEXT.secondary} ${HIGH_CONTRAST_CLASSES.BORDER.primary}`,
      brand: `bg-orange-500 ${HIGH_CONTRAST_CLASSES.TEXT.primary} border-orange-500`,
    },
  }

  return baseClasses[component as keyof typeof baseClasses]?.[variant || 'default'] || ''
}
