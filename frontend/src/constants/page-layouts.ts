/**
 * Centralized Page Layout System for BiHub E-commerce
 * All page sizing, spacing, and layout configurations in one place
 * Easy to maintain and update across the entire application
 */

import { DESIGN_SYSTEM } from './design-system'

// ===== PAGE LAYOUT CONSTANTS =====

export const PAGE_LAYOUTS = {
  // Container and spacing
  CONTAINER: {
    className: 'container mx-auto px-4',
    maxWidth: '1280px', // xl breakpoint
  },

  // Section spacing (consistent across all pages)
  SECTION: {
    // Vertical padding for page sections
    padding: {
      sm: 'py-6',      // 24px top/bottom - for compact sections
      base: 'py-8',    // 32px top/bottom - standard sections  
      lg: 'py-12',     // 48px top/bottom - major sections
      xl: 'py-16',     // 64px top/bottom - hero sections
    },
    
    // Gap between sections
    gap: {
      sm: 'space-y-6',   // 24px between sections
      base: 'space-y-8', // 32px between sections
      lg: 'space-y-12',  // 48px between sections
    },
  },

  // Page headers (consistent across all pages)
  PAGE_HEADER: {
    // Container
    container: 'bg-gradient-to-r from-gray-900 to-black border-b border-gray-800',
    padding: 'py-6 lg:py-8',
    
    // Title styling
    title: {
      sm: 'text-xl lg:text-2xl font-bold text-white mb-2',
      base: 'text-2xl lg:text-3xl font-bold text-white mb-2', 
      lg: 'text-3xl lg:text-4xl font-bold text-white mb-3',
    },
    
    // Subtitle/description
    subtitle: 'text-sm lg:text-base text-gray-300 mb-4',
    
    // Breadcrumbs
    breadcrumbs: 'text-xs text-gray-400 mb-3',
  },

  // Content areas
  CONTENT: {
    // Main content wrapper
    wrapper: 'min-h-screen bg-black text-white',
    
    // Content sections
    section: 'py-8',
    
    // Grid layouts
    grid: {
      // Product grids
      products: 'grid gap-4 lg:gap-6 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
      
      // Category grids  
      categories: 'grid gap-4 grid-cols-2 md:grid-cols-4 lg:grid-cols-6',
      
      // Feature grids (3 columns)
      features: 'grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
      
      // Two column layout
      twoColumn: 'grid gap-8 grid-cols-1 lg:grid-cols-2',
      
      // Three column layout  
      threeColumn: 'grid gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3',
    },
  },

  // Card styling (consistent across all pages)
  CARD: {
    // Base card styles
    base: 'bg-gray-800 border border-gray-700 rounded-lg hover:border-gray-600 transition-all duration-300',
    
    // Card padding options
    padding: {
      sm: 'p-4',
      base: 'p-6', 
      lg: 'p-8',
    },
    
    // Card shadows
    shadow: {
      none: '',
      sm: 'shadow-sm',
      base: 'shadow-md',
      lg: 'shadow-lg',
    },
  },

  // Form layouts
  FORM: {
    // Form container
    container: 'space-y-6',
    
    // Form sections
    section: 'space-y-4',
    
    // Field groups
    fieldGroup: 'space-y-2',
    
    // Button groups
    buttonGroup: 'flex gap-3 pt-4',
    
    // Two column form
    twoColumn: 'grid gap-4 grid-cols-1 md:grid-cols-2',
  },

  // Typography scales for different page types
  TYPOGRAPHY: {
    // Hero sections
    hero: {
      title: 'text-3xl lg:text-4xl xl:text-5xl font-bold leading-tight',
      subtitle: 'text-base lg:text-lg text-gray-300 leading-relaxed',
    },
    
    // Page titles
    pageTitle: {
      sm: 'text-xl lg:text-2xl font-bold text-white',
      base: 'text-2xl lg:text-3xl font-bold text-white',
      lg: 'text-3xl lg:text-4xl font-bold text-white',
    },
    
    // Section titles
    sectionTitle: {
      sm: 'text-lg lg:text-xl font-bold text-white',
      base: 'text-xl lg:text-2xl font-bold text-white', 
      lg: 'text-2xl lg:text-3xl font-bold text-white',
    },
    
    // Card titles
    cardTitle: 'text-lg font-semibold text-white',
    
    // Body text
    body: {
      sm: 'text-sm text-gray-300',
      base: 'text-base text-gray-300',
      lg: 'text-lg text-gray-300',
    },
    
    // Labels
    label: 'text-sm font-medium text-gray-200',
    
    // Captions
    caption: 'text-xs text-gray-400',
  },

  // Button configurations
  BUTTON: {
    // Size mappings to our design system
    sizes: {
      sm: 'h-8 px-3 text-sm',
      base: 'h-10 px-4 text-base',
      lg: 'h-12 px-6 text-lg',
    },
    
    // Common button groups
    group: 'flex gap-3',
    groupVertical: 'flex flex-col gap-2',
  },

  // Loading states
  LOADING: {
    // Skeleton cards
    skeletonCard: 'animate-pulse bg-gray-800 border border-gray-700 rounded-lg',
    skeletonContent: 'space-y-3 p-4',
    skeletonLine: 'h-4 bg-gray-600 rounded',
    skeletonImage: 'aspect-square bg-gray-700 rounded-lg',
  },

  // Responsive breakpoints for consistent usage
  BREAKPOINTS: {
    sm: '640px',
    md: '768px', 
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
  },
} as const

// ===== UTILITY FUNCTIONS =====

/**
 * Get consistent page wrapper classes
 */
export const getPageWrapperClasses = () => {
  return PAGE_LAYOUTS.CONTENT.wrapper
}

/**
 * Get consistent container classes
 */
export const getContainerClasses = () => {
  return PAGE_LAYOUTS.CONTAINER.className
}

/**
 * Get consistent section classes
 */
export const getSectionClasses = (size: 'sm' | 'base' | 'lg' | 'xl' = 'base') => {
  return PAGE_LAYOUTS.SECTION.padding[size]
}

/**
 * Get consistent page header classes
 */
export const getPageHeaderClasses = () => {
  return `${PAGE_LAYOUTS.PAGE_HEADER.container} ${PAGE_LAYOUTS.PAGE_HEADER.padding}`
}

/**
 * Get consistent page title classes
 */
export const getPageTitleClasses = (size: 'sm' | 'base' | 'lg' = 'base') => {
  return PAGE_LAYOUTS.PAGE_HEADER.title[size]
}

/**
 * Get consistent card classes
 */
export const getCardClasses = (padding: 'sm' | 'base' | 'lg' = 'base', shadow: 'none' | 'sm' | 'base' | 'lg' = 'none') => {
  return `${PAGE_LAYOUTS.CARD.base} ${PAGE_LAYOUTS.CARD.padding[padding]} ${PAGE_LAYOUTS.CARD.shadow[shadow]}`
}

/**
 * Get consistent grid classes
 */
export const getGridClasses = (type: 'products' | 'categories' | 'features' | 'twoColumn' | 'threeColumn') => {
  return PAGE_LAYOUTS.CONTENT.grid[type]
}

/**
 * Get consistent form classes
 */
export const getFormClasses = (type: 'container' | 'section' | 'fieldGroup' | 'buttonGroup' | 'twoColumn' = 'container') => {
  return PAGE_LAYOUTS.FORM[type]
}

/**
 * Get consistent typography classes
 */
export const getTypographyClasses = (
  category: 'hero' | 'pageTitle' | 'sectionTitle' | 'cardTitle' | 'body' | 'label' | 'caption',
  size?: 'sm' | 'base' | 'lg'
) => {
  const typography = PAGE_LAYOUTS.TYPOGRAPHY[category]
  
  if (typeof typography === 'string') {
    return typography
  }
  
  if (category === 'hero') {
    return size === 'subtitle' ? typography.subtitle : typography.title
  }
  
  return typography[size || 'base']
}

// ===== COMMON PAGE PATTERNS =====

export const PAGE_PATTERNS = {
  // Standard page with header
  standardPage: {
    wrapper: getPageWrapperClasses(),
    header: getPageHeaderClasses(),
    container: getContainerClasses(),
    content: getSectionClasses('base'),
  },
  
  // Landing page pattern
  landingPage: {
    wrapper: getPageWrapperClasses(),
    hero: getSectionClasses('xl'),
    section: getSectionClasses('lg'),
    container: getContainerClasses(),
  },
  
  // Form page pattern
  formPage: {
    wrapper: getPageWrapperClasses(),
    header: getPageHeaderClasses(),
    container: getContainerClasses(),
    form: getFormClasses('container'),
  },
  
  // Grid page pattern (products, categories, etc.)
  gridPage: {
    wrapper: getPageWrapperClasses(),
    header: getPageHeaderClasses(),
    container: getContainerClasses(),
    grid: getGridClasses('products'),
  },
} as const
