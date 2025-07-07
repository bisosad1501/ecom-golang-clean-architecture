/**
 * Unified Design System for BiHub E-commerce
 * Consistent spacing, sizing, colors, and typography across the application
 */

export const DESIGN_SYSTEM = {
  // Color Palette
  COLORS: {
    // Primary Brand Colors
    PRIMARY: {
      50: '#FFF7ED',
      100: '#FFEDD5',
      200: '#FED7AA',
      300: '#FDBA74',
      400: '#FB923C',
      500: '#FF9000', // Main brand color
      600: '#E67E00',
      700: '#C2410C',
      800: '#9A3412',
      900: '#7C2D12',
    },
    
    // Neutral Colors (Dark Theme)
    NEUTRAL: {
      50: '#F9FAFB',
      100: '#F3F4F6',
      200: '#E5E7EB',
      300: '#D1D5DB',
      400: '#9CA3AF',
      500: '#6B7280',
      600: '#4B5563',
      700: '#374151',
      800: '#1F2937',
      900: '#111827',
      950: '#0A0A0A',
    },
    
    // Semantic Colors
    SUCCESS: '#10B981',
    WARNING: '#F59E0B',
    ERROR: '#EF4444',
    INFO: '#3B82F6',
  },

  // Typography Scale
  TYPOGRAPHY: {
    FONT_SIZES: {
      xs: '0.75rem',    // 12px
      sm: '0.875rem',   // 14px
      base: '1rem',     // 16px
      lg: '1.125rem',   // 18px
      xl: '1.25rem',    // 20px
      '2xl': '1.5rem',  // 24px
      '3xl': '1.875rem', // 30px
      '4xl': '2.25rem', // 36px
      '5xl': '3rem',    // 48px
    },
    
    FONT_WEIGHTS: {
      normal: '400',
      medium: '500',
      semibold: '600',
      bold: '700',
    },
    
    LINE_HEIGHTS: {
      tight: '1.25',
      normal: '1.5',
      relaxed: '1.75',
    },
  },

  // Spacing Scale (consistent with Tailwind)
  SPACING: {
    0: '0',
    1: '0.25rem',   // 4px
    2: '0.5rem',    // 8px
    3: '0.75rem',   // 12px
    4: '1rem',      // 16px
    5: '1.25rem',   // 20px
    6: '1.5rem',    // 24px
    8: '2rem',      // 32px
    10: '2.5rem',   // 40px
    12: '3rem',     // 48px
    16: '4rem',     // 64px
    20: '5rem',     // 80px
    24: '6rem',     // 96px
  },

  // Border Radius
  RADIUS: {
    none: '0',
    sm: '0.125rem',   // 2px
    base: '0.25rem',  // 4px
    md: '0.375rem',   // 6px
    lg: '0.5rem',     // 8px
    xl: '0.75rem',    // 12px
    '2xl': '1rem',    // 16px
    full: '9999px',
  },

  // Shadows
  SHADOWS: {
    sm: '0 1px 2px 0 rgb(0 0 0 / 0.05)',
    base: '0 1px 3px 0 rgb(0 0 0 / 0.1), 0 1px 2px -1px rgb(0 0 0 / 0.1)',
    md: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
    lg: '0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)',
    xl: '0 20px 25px -5px rgb(0 0 0 / 0.1), 0 8px 10px -6px rgb(0 0 0 / 0.1)',
  },

  // Component Sizes
  COMPONENT_SIZES: {
    // Button sizes
    BUTTON: {
      sm: {
        height: '2rem',      // 32px
        padding: '0.5rem 0.75rem',
        fontSize: '0.875rem',
      },
      base: {
        height: '2.5rem',    // 40px
        padding: '0.625rem 1rem',
        fontSize: '1rem',
      },
      lg: {
        height: '3rem',      // 48px
        padding: '0.75rem 1.5rem',
        fontSize: '1.125rem',
      },
    },
    
    // Input sizes
    INPUT: {
      sm: {
        height: '2rem',
        padding: '0.5rem 0.75rem',
        fontSize: '0.875rem',
      },
      base: {
        height: '2.5rem',
        padding: '0.625rem 0.75rem',
        fontSize: '1rem',
      },
      lg: {
        height: '3rem',
        padding: '0.75rem 1rem',
        fontSize: '1.125rem',
      },
    },
    
    // Card padding
    CARD: {
      sm: '1rem',      // 16px
      base: '1.5rem',  // 24px
      lg: '2rem',      // 32px
    },
  },

  // Layout Breakpoints
  BREAKPOINTS: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
  },

  // Grid Systems
  GRID: {
    PRODUCT_GRID: {
      mobile: 'grid-cols-1',
      tablet: 'sm:grid-cols-2',
      desktop: 'lg:grid-cols-3 xl:grid-cols-4',
    },
    
    CATEGORY_GRID: {
      mobile: 'grid-cols-2',
      tablet: 'md:grid-cols-4',
      desktop: 'lg:grid-cols-6',
    },
  },

  // Animation Durations
  ANIMATION: {
    fast: '150ms',
    normal: '300ms',
    slow: '500ms',
  },

  // Z-Index Scale
  Z_INDEX: {
    dropdown: 1000,
    sticky: 1020,
    fixed: 1030,
    modal_backdrop: 1040,
    modal: 1050,
    popover: 1060,
    tooltip: 1070,
    toast: 1080,
  },
} as const

// Utility functions for consistent styling
export const getButtonClasses = (size: 'sm' | 'base' | 'lg' = 'base', variant: 'primary' | 'secondary' | 'outline' = 'primary') => {
  const sizeClasses = {
    sm: 'h-8 px-3 text-sm',
    base: 'h-10 px-4 text-base',
    lg: 'h-12 px-6 text-lg',
  }
  
  const variantClasses = {
    primary: 'bg-orange-500 hover:bg-orange-600 text-white',
    secondary: 'bg-gray-600 hover:bg-gray-700 text-white',
    outline: 'border border-gray-600 text-white hover:bg-gray-800',
  }
  
  return `${sizeClasses[size]} ${variantClasses[variant]} rounded-lg transition-all duration-300 font-medium`
}

export const getCardClasses = (padding: 'sm' | 'base' | 'lg' = 'base') => {
  const paddingClasses = {
    sm: 'p-4',
    base: 'p-6',
    lg: 'p-8',
  }
  
  return `bg-gray-800 border border-gray-700 rounded-lg ${paddingClasses[padding]} hover:border-gray-600 transition-all duration-300`
}

export const getTextClasses = (size: keyof typeof DESIGN_SYSTEM.TYPOGRAPHY.FONT_SIZES, weight: keyof typeof DESIGN_SYSTEM.TYPOGRAPHY.FONT_WEIGHTS = 'normal') => {
  return `text-${size} font-${weight}`
}

// Consistent spacing utilities
export const SPACING_CLASSES = {
  section: 'py-8 lg:py-12',
  container: 'container mx-auto px-4',
  cardGap: 'gap-4 lg:gap-6',
  elementGap: 'gap-2 lg:gap-3',
} as const
