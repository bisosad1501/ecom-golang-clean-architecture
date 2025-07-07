// Unified Design System Tokens for ShopHub
export const DESIGN_TOKENS = {
  // ===== COLOR SYSTEM =====
  COLORS: {
    // Primary Brand Colors
    PRIMARY: {
      50: '#FFF8F0',
      100: '#FFEFDB',
      200: '#FFDFB7',
      300: '#FFCF93',
      400: '#FFBF6F',
      500: '#FF9000',  // Main brand orange
      600: '#E6820E',
      700: '#CC7300',
      800: '#B26400',
      900: '#995500',
      950: '#663800',
    },

    // Neutral Colors
    NEUTRAL: {
      0: '#FFFFFF',
      50: '#F8FAFC',
      100: '#F1F5F9',
      200: '#E2E8F0',
      300: '#CBD5E1',
      400: '#94A3B8',
      500: '#64748B',
      600: '#475569',
      700: '#334155',
      800: '#1E293B',
      900: '#0F172A',
      950: '#000000',
    },

    // Semantic Colors
    SUCCESS: { 50: '#F0FDF4', 500: '#22C55E', 600: '#16A34A', 900: '#14532D' },
    WARNING: { 50: '#FFFBEB', 500: '#F59E0B', 600: '#D97706', 900: '#92400E' },
    ERROR: { 50: '#FEF2F2', 500: '#EF4444', 600: '#DC2626', 900: '#7F1D1D' },
    INFO: { 50: '#EFF6FF', 500: '#3B82F6', 600: '#2563EB', 900: '#1E3A8A' },
  },

  // ===== TYPOGRAPHY SYSTEM =====
  TYPOGRAPHY: {
    // Font Families
    FAMILIES: {
      SANS: ['Inter', 'system-ui', 'sans-serif'],
      MONO: ['JetBrains Mono', 'Consolas', 'monospace'],
    },

    // Font Sizes & Line Heights
    SCALE: {
      XS: { size: '0.75rem', lineHeight: '1rem' },     // 12px
      SM: { size: '0.875rem', lineHeight: '1.25rem' }, // 14px
      BASE: { size: '1rem', lineHeight: '1.5rem' },    // 16px
      LG: { size: '1.125rem', lineHeight: '1.75rem' }, // 18px
      XL: { size: '1.25rem', lineHeight: '1.75rem' },  // 20px
      '2XL': { size: '1.5rem', lineHeight: '2rem' },   // 24px
      '3XL': { size: '1.875rem', lineHeight: '2.25rem' }, // 30px
      '4XL': { size: '2.25rem', lineHeight: '2.5rem' },   // 36px
      '5XL': { size: '3rem', lineHeight: '1' },           // 48px
      '6XL': { size: '3.75rem', lineHeight: '1' },        // 60px
    },

    // Font Weights
    WEIGHTS: {
      LIGHT: 300,
      NORMAL: 400,
      MEDIUM: 500,
      SEMIBOLD: 600,
      BOLD: 700,
      EXTRABOLD: 800,
    },

    // Tailwind Classes (for convenience)
    CLASSES: {
      DISPLAY_LARGE: 'text-4xl lg:text-5xl font-bold',
      DISPLAY_MEDIUM: 'text-3xl lg:text-4xl font-bold',
      HEADING_1: 'text-2xl lg:text-3xl font-semibold',
      HEADING_2: 'text-xl lg:text-2xl font-semibold',
      HEADING_3: 'text-lg lg:text-xl font-medium',
      BODY_LARGE: 'text-base',
      BODY_DEFAULT: 'text-sm',
      BODY_SMALL: 'text-xs',
      LABEL: 'text-xs font-medium',
      CAPTION: 'text-xs text-muted-foreground',
    },
  },

  // ===== SPACING SYSTEM =====
  SPACING: {
    // Base spacing scale (rem)
    SCALE: {
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
      32: '8rem',     // 128px
    },

    // Tailwind Classes (for convenience)
    CLASSES: {
      SECTION_LARGE: 'py-16',
      SECTION_DEFAULT: 'py-12',
      SECTION_SMALL: 'py-8',
      COMPONENT_LARGE: 'p-6',
      COMPONENT_DEFAULT: 'p-4',
      COMPONENT_SMALL: 'p-3',
      GAP_LARGE: 'gap-8',
      GAP_DEFAULT: 'gap-6',
      GAP_SMALL: 'gap-4',
      GAP_TINY: 'gap-2',
    },
  },

  // ===== BORDER RADIUS =====
  RADIUS: {
    NONE: 'rounded-none',
    SM: 'rounded-sm',
    BASE: 'rounded',
    MD: 'rounded-md',
    LG: 'rounded-lg',
    XL: 'rounded-xl',
    '2XL': 'rounded-2xl',
    '3XL': 'rounded-3xl',
    FULL: 'rounded-full',
    // Legacy support
    LARGE: 'rounded-lg',
    DEFAULT: 'rounded-md',
    SMALL: 'rounded-sm',
  },

  // ===== COMPONENT TOKENS =====
  COMPONENTS: {
    // Button sizes
    BUTTON: {
      SM: 'h-8 px-3 text-xs',
      MD: 'h-10 px-4 text-sm',
      LG: 'h-12 px-6 text-base',
      ICON_SM: 'h-8 w-8',
      ICON_MD: 'h-10 w-10',
      ICON_LG: 'h-12 w-12',
    },

    // Input sizes
    INPUT: {
      SM: 'h-8 px-3 text-sm',
      MD: 'h-10 px-4 text-sm',
      LG: 'h-12 px-4 text-base',
    },

    // Card padding
    CARD: {
      SM: 'p-4',
      MD: 'p-6',
      LG: 'p-8',
    },

    // Icon sizes
    ICON: {
      XS: 'h-3 w-3',
      SM: 'h-4 w-4',
      MD: 'h-5 w-5',
      LG: 'h-6 w-6',
      XL: 'h-8 w-8',
    },
  },

  // ===== LEGACY SUPPORT =====
  // These are for backward compatibility with existing components
  BUTTONS: {
    ICON_DEFAULT: 'h-10 w-10',
    ICON_SMALL: 'h-8 w-8',
    ICON_LARGE: 'h-12 w-12',
  },

  ICONS: {
    TINY: 'h-3 w-3',
    SMALL: 'h-4 w-4',
    DEFAULT: 'h-5 w-5',
    LARGE: 'h-6 w-6',
    EXTRA_LARGE: 'h-8 w-8',
  },

  CONTAINERS: {
    CARD_PADDING: 'p-4',
    FORM_PADDING: 'p-6',
    MODAL_PADDING: 'p-8',
  },

  // ===== ANIMATION SYSTEM =====
  ANIMATIONS: {
    DURATION: {
      FAST: '150ms',
      DEFAULT: '200ms',
      SLOW: '300ms',
      SLOWER: '500ms',
    },

    EASING: {
      DEFAULT: 'cubic-bezier(0.4, 0, 0.2, 1)',
      IN: 'cubic-bezier(0.4, 0, 1, 1)',
      OUT: 'cubic-bezier(0, 0, 0.2, 1)',
      IN_OUT: 'cubic-bezier(0.4, 0, 0.2, 1)',
    },

    CLASSES: {
      TRANSITION: 'transition-all duration-200 ease-out',
      HOVER_SCALE: 'hover:scale-105',
      HOVER_LIFT: 'hover:-translate-y-1',
      PULSE: 'animate-pulse',
      SPIN: 'animate-spin',
      FLOAT: 'animate-float',
    },
  },

  // ===== BREAKPOINTS =====
  BREAKPOINTS: {
    SM: '640px',
    MD: '768px',
    LG: '1024px',
    XL: '1280px',
    '2XL': '1536px',
  },
} as const

// Type helpers
export type DesignTokens = typeof DESIGN_TOKENS
export type ColorScale = keyof typeof DESIGN_TOKENS.COLORS.PRIMARY
export type SpacingScale = keyof typeof DESIGN_TOKENS.SPACING.SCALE
export type TypographyScale = keyof typeof DESIGN_TOKENS.TYPOGRAPHY.SCALE
