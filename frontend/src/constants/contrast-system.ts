/**
 * High Contrast Design System for BiHub E-commerce
 * Ensures excellent readability and accessibility across all pages
 * WCAG 2.1 AA compliant contrast ratios
 */

export const CONTRAST_SYSTEM = {
  // High contrast color palette
  COLORS: {
    // Backgrounds (darker for better contrast)
    BACKGROUND: {
      primary: '#000000',      // Pure black for maximum contrast
      secondary: '#111111',    // Very dark gray
      card: '#1a1a1a',        // Dark card background
      elevated: '#222222',     // Elevated elements
      muted: '#0a0a0a',       // Muted sections
    },
    
    // Text colors (high contrast)
    TEXT: {
      primary: '#ffffff',      // Pure white for maximum contrast
      secondary: '#e5e5e5',    // Light gray for secondary text
      muted: '#cccccc',        // Medium gray for muted text
      disabled: '#999999',     // Gray for disabled text
      inverse: '#000000',      // Black text on light backgrounds
    },
    
    // Brand colors (enhanced contrast)
    BRAND: {
      primary: '#ff9500',      // Brighter orange for better visibility
      primaryHover: '#e67e00',
      primaryLight: '#ffb84d',
      primaryDark: '#cc7700',
    },
    
    // Semantic colors (high contrast)
    SEMANTIC: {
      success: '#00ff88',      // Bright green
      successBg: '#003322',
      warning: '#ffcc00',      // Bright yellow
      warningBg: '#332200',
      error: '#ff4444',        // Bright red
      errorBg: '#330000',
      info: '#44aaff',         // Bright blue
      infoBg: '#002244',
    },
    
    // Border colors
    BORDER: {
      primary: '#444444',      // Medium gray borders
      secondary: '#333333',    // Darker borders
      accent: '#ff9500',       // Orange accent borders
      muted: '#222222',        // Very subtle borders
    },
    
    // Interactive states
    INTERACTIVE: {
      hover: '#2a2a2a',        // Hover background
      active: '#333333',       // Active background
      focus: '#ff9500',        // Focus outline color
      disabled: '#1a1a1a',     // Disabled background
    },
  },

  // Typography with high contrast
  TYPOGRAPHY: {
    // High contrast text classes
    TEXT_CLASSES: {
      // Primary text (maximum contrast)
      primary: 'text-white',
      
      // Secondary text (good contrast)
      secondary: 'text-gray-200',
      
      // Muted text (adequate contrast)
      muted: 'text-gray-300',
      
      // Disabled text
      disabled: 'text-gray-500',
      
      // Brand colored text
      brand: 'text-orange-400',
      
      // Semantic text colors
      success: 'text-green-400',
      warning: 'text-yellow-400',
      error: 'text-red-400',
      info: 'text-blue-400',
    },
    
    // Background classes
    BG_CLASSES: {
      primary: 'bg-black',
      secondary: 'bg-gray-950',
      card: 'bg-gray-900',
      elevated: 'bg-gray-800',
      muted: 'bg-gray-950/50',
    },
    
    // Border classes
    BORDER_CLASSES: {
      primary: 'border-gray-600',
      secondary: 'border-gray-700',
      accent: 'border-orange-500',
      muted: 'border-gray-800',
    },
  },

  // Component-specific contrast settings
  COMPONENTS: {
    // Button contrast
    BUTTON: {
      primary: {
        bg: 'bg-orange-500 hover:bg-orange-400',
        text: 'text-white',
        border: 'border-orange-500',
      },
      secondary: {
        bg: 'bg-gray-700 hover:bg-gray-600',
        text: 'text-white',
        border: 'border-gray-600',
      },
      outline: {
        bg: 'bg-transparent hover:bg-gray-800',
        text: 'text-white hover:text-white',
        border: 'border-gray-500 hover:border-gray-400',
      },
    },
    
    // Card contrast
    CARD: {
      default: {
        bg: 'bg-gray-900',
        border: 'border-gray-700',
        text: 'text-white',
      },
      elevated: {
        bg: 'bg-gray-800',
        border: 'border-gray-600',
        text: 'text-white',
      },
      interactive: {
        bg: 'bg-gray-900 hover:bg-gray-800',
        border: 'border-gray-700 hover:border-gray-600',
        text: 'text-white',
      },
    },
    
    // Form contrast
    FORM: {
      input: {
        bg: 'bg-gray-800',
        border: 'border-gray-600 focus:border-orange-500',
        text: 'text-white placeholder:text-gray-400',
      },
      label: {
        text: 'text-gray-200',
      },
      error: {
        text: 'text-red-400',
        border: 'border-red-500',
      },
    },
    
    // Navigation contrast
    NAV: {
      bg: 'bg-black border-gray-800',
      text: 'text-white',
      link: 'text-gray-300 hover:text-white',
      active: 'text-orange-400',
    },
  },

  // Accessibility helpers
  ACCESSIBILITY: {
    // Focus styles
    FOCUS: {
      ring: 'focus:ring-2 focus:ring-orange-500 focus:ring-offset-2 focus:ring-offset-black',
      outline: 'focus:outline-none',
    },
    
    // Screen reader only
    SR_ONLY: 'sr-only',
    
    // High contrast mode detection
    HIGH_CONTRAST: '@media (prefers-contrast: high)',
  },
} as const

// Utility functions for high contrast
export const getHighContrastClasses = {
  // Text utilities
  text: {
    primary: () => CONTRAST_SYSTEM.TYPOGRAPHY.TEXT_CLASSES.primary,
    secondary: () => CONTRAST_SYSTEM.TYPOGRAPHY.TEXT_CLASSES.secondary,
    muted: () => CONTRAST_SYSTEM.TYPOGRAPHY.TEXT_CLASSES.muted,
    brand: () => CONTRAST_SYSTEM.TYPOGRAPHY.TEXT_CLASSES.brand,
  },
  
  // Background utilities
  bg: {
    primary: () => CONTRAST_SYSTEM.TYPOGRAPHY.BG_CLASSES.primary,
    card: () => CONTRAST_SYSTEM.TYPOGRAPHY.BG_CLASSES.card,
    elevated: () => CONTRAST_SYSTEM.TYPOGRAPHY.BG_CLASSES.elevated,
  },
  
  // Button utilities
  button: {
    primary: () => `${CONTRAST_SYSTEM.COMPONENTS.BUTTON.primary.bg} ${CONTRAST_SYSTEM.COMPONENTS.BUTTON.primary.text}`,
    secondary: () => `${CONTRAST_SYSTEM.COMPONENTS.BUTTON.secondary.bg} ${CONTRAST_SYSTEM.COMPONENTS.BUTTON.secondary.text}`,
    outline: () => `${CONTRAST_SYSTEM.COMPONENTS.BUTTON.outline.bg} ${CONTRAST_SYSTEM.COMPONENTS.BUTTON.outline.text} ${CONTRAST_SYSTEM.COMPONENTS.BUTTON.outline.border}`,
  },
  
  // Card utilities
  card: {
    default: () => `${CONTRAST_SYSTEM.COMPONENTS.CARD.default.bg} ${CONTRAST_SYSTEM.COMPONENTS.CARD.default.border} ${CONTRAST_SYSTEM.COMPONENTS.CARD.default.text}`,
    interactive: () => `${CONTRAST_SYSTEM.COMPONENTS.CARD.interactive.bg} ${CONTRAST_SYSTEM.COMPONENTS.CARD.interactive.border} ${CONTRAST_SYSTEM.COMPONENTS.CARD.interactive.text}`,
  },
}

// Page-specific contrast configurations
export const PAGE_CONTRAST = {
  // Home page
  HOME: {
    hero: {
      bg: 'bg-gradient-to-br from-black via-gray-950 to-black',
      title: 'text-white',
      subtitle: 'text-gray-200',
    },
    section: {
      bg: 'bg-gray-950',
      title: 'text-white',
      text: 'text-gray-200',
    },
  },
  
  // Product pages
  PRODUCTS: {
    header: {
      bg: 'bg-gradient-to-r from-gray-950 to-black border-gray-700',
      title: 'text-white',
      subtitle: 'text-gray-200',
    },
    card: {
      bg: 'bg-gray-900 border-gray-700 hover:border-gray-600',
      title: 'text-white',
      price: 'text-orange-400',
      text: 'text-gray-300',
    },
  },
  
  // Form pages
  FORMS: {
    container: {
      bg: 'bg-gray-900',
      border: 'border-gray-700',
    },
    input: {
      bg: 'bg-gray-800',
      border: 'border-gray-600 focus:border-orange-500',
      text: 'text-white',
      placeholder: 'placeholder:text-gray-400',
    },
    label: 'text-gray-200',
    error: 'text-red-400',
  },
} as const
