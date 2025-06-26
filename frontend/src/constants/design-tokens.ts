// Design Tokens for Consistent UI - BiHub PornHub-inspired Theme
export const DESIGN_TOKENS = {
  // Color System (PornHub inspired)
  COLORS: {
    // Primary colors - official PornHub palette
    PRIMARY: {
      BLACK: '#000000',        // Pure black
      ORANGE: '#FF9000',       // Official PornHub orange
      WHITE: '#FFFFFF',        // Pure white
    },
    
    // Background variants
    BACKGROUNDS: {
      PRIMARY: 'bg-black',
      SECONDARY: 'bg-gray-900',
      SURFACE: 'bg-white',
      MUTED: 'bg-gray-100',
    },
    
    // Text colors
    TEXT: {
      PRIMARY: 'text-white',           // White text on dark
      SECONDARY: 'text-gray-300',      // Light gray
      ACCENT: 'text-orange-500',       // Orange accent
      DARK: 'text-black',              // Black text on light
      MUTED: 'text-gray-500',          // Muted text
    },
    
    // Hover states
    HOVER: {
      ORANGE: 'hover:bg-orange-500',
      ORANGE_LIGHT: 'hover:bg-orange-500/10',
      BLACK: 'hover:bg-black/90',
    },
    
    // Border colors
    BORDERS: {
      DEFAULT: 'border-gray-800',
      LIGHT: 'border-gray-200',
      ORANGE: 'border-orange-500',
    }
  },

  // Typography Scale (following Tailwind scale but standardized)
  TYPOGRAPHY: {
    // Display text (hero sections)
    DISPLAY_LARGE: 'text-3xl lg:text-4xl font-bold',
    DISPLAY_MEDIUM: 'text-2xl lg:text-3xl font-bold',
    
    // Headings
    HEADING_1: 'text-xl lg:text-2xl font-bold',
    HEADING_2: 'text-lg lg:text-xl font-bold', 
    HEADING_3: 'text-base lg:text-lg font-semibold',
    
    // Body text
    BODY_LARGE: 'text-base',
    BODY_DEFAULT: 'text-sm',
    BODY_SMALL: 'text-xs',
    
    // Labels and captions
    LABEL: 'text-xs font-medium',
    CAPTION: 'text-xs text-muted-foreground',
  },

  // Spacing Scale
  SPACING: {
    // Section padding
    SECTION_LARGE: 'py-16',
    SECTION_DEFAULT: 'py-12', 
    SECTION_SMALL: 'py-8',
    
    // Component padding
    COMPONENT_LARGE: 'p-6',
    COMPONENT_DEFAULT: 'p-4',
    COMPONENT_SMALL: 'p-3',
    
    // Gaps
    GAP_LARGE: 'gap-8',
    GAP_DEFAULT: 'gap-6',
    GAP_SMALL: 'gap-4',
    GAP_TINY: 'gap-2',
    
    // Margins
    MARGIN_LARGE: 'mb-8',
    MARGIN_DEFAULT: 'mb-6',
    MARGIN_SMALL: 'mb-4',
    MARGIN_TINY: 'mb-2',
  },

  // Icon Sizes  
  ICONS: {
    LARGE: 'h-6 w-6',
    DEFAULT: 'h-4 w-4', 
    SMALL: 'h-3.5 w-3.5',
    TINY: 'h-3 w-3',
  },

  // Border Radius
  RADIUS: {
    LARGE: 'rounded-xl',
    DEFAULT: 'rounded-lg',
    SMALL: 'rounded-md',
    FULL: 'rounded-full',
  },

  // Button Sizes (height consistency)
  BUTTONS: {
    LARGE: 'h-12 px-6 text-base',
    DEFAULT: 'h-10 px-4 text-sm',
    SMALL: 'h-8 px-3 text-xs',
    ICON_LARGE: 'h-12 w-12',
    ICON_DEFAULT: 'h-10 w-10', 
    ICON_SMALL: 'h-8 w-8',
  },

  // Container sizes
  CONTAINERS: {
    CARD_PADDING: 'p-4',
    FORM_PADDING: 'p-6',
    MODAL_PADDING: 'p-8',
  }
} as const
