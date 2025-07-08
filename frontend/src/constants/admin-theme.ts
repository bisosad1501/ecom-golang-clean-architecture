// BiHub Admin Theme System
export const BIHUB_ADMIN_THEME = {
  // Brand Colors
  colors: {
    primary: '#FF9000',
    primaryHover: '#e67e00',
    primaryDark: '#cc6600',
    primaryLight: '#ffb84d',
    
    // Background Colors
    background: {
      main: 'from-slate-950 via-slate-900 to-slate-950',
      card: 'bg-gray-800/50',
      cardHover: 'bg-gray-800/70',
      sidebar: 'from-gray-900 via-gray-900 to-gray-800',
      header: 'bg-gray-900/90',
    },
    
    // Text Colors
    text: {
      primary: 'text-white',
      secondary: 'text-gray-300',
      muted: 'text-gray-400',
      disabled: 'text-gray-500',
    },
    
    // Border Colors
    border: {
      default: 'border-gray-700/50',
      hover: 'border-[#FF9000]',
      focus: 'border-[#FF9000]',
    },
    
    // Status Colors
    status: {
      success: {
        bg: 'bg-emerald-900/30',
        text: 'text-emerald-400',
        border: 'border-emerald-500/30',
      },
      warning: {
        bg: 'bg-yellow-900/30',
        text: 'text-yellow-400',
        border: 'border-yellow-500/30',
      },
      error: {
        bg: 'bg-red-900/30',
        text: 'text-red-400',
        border: 'border-red-500/30',
      },
      info: {
        bg: 'bg-blue-900/30',
        text: 'text-blue-400',
        border: 'border-blue-500/30',
      },
    },
  },
  
  // Component Styles
  components: {
    card: {
      base: 'bg-gray-800/50 border border-gray-700/50 rounded-2xl shadow-xl',
      hover: 'hover:bg-gray-800/70 hover:shadow-2xl transition-all duration-300',
      padding: 'p-6',
    },
    
    button: {
      primary: 'bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white font-semibold rounded-xl transition-all duration-300 transform hover:scale-[1.02] shadow-lg',
      secondary: 'border-2 border-gray-600 hover:border-[#FF9000] hover:bg-[#FF9000]/5 text-gray-300 hover:text-white rounded-xl transition-all duration-300',
      ghost: 'hover:bg-gray-800 text-gray-400 hover:text-white rounded-xl transition-all duration-300',
    },
    
    input: {
      base: 'bg-gray-800/90 border border-gray-600/80 text-white placeholder:text-gray-400 focus:border-[#FF9000] focus:ring-2 focus:ring-[#FF9000]/20 rounded-lg backdrop-blur-sm transition-all duration-300',
    },
    
    badge: {
      primary: 'bg-[#FF9000]/20 text-[#FF9000] border border-[#FF9000]/30 rounded-full px-3 py-1 text-xs font-semibold',
      success: 'bg-emerald-900/30 text-emerald-400 border border-emerald-500/30 rounded-full px-3 py-1 text-xs font-semibold',
      warning: 'bg-yellow-900/30 text-yellow-400 border border-yellow-500/30 rounded-full px-3 py-1 text-xs font-semibold',
      error: 'bg-red-900/30 text-red-400 border border-red-500/30 rounded-full px-3 py-1 text-xs font-semibold',
    },
  },
  
  // Typography
  typography: {
    heading: {
      h1: 'text-3xl font-bold text-white',
      h2: 'text-2xl font-bold text-white',
      h3: 'text-xl font-semibold text-white',
      h4: 'text-lg font-semibold text-white',
    },
    body: {
      large: 'text-base text-gray-300',
      medium: 'text-sm text-gray-300',
      small: 'text-xs text-gray-400',
    },
  },
  
  // Spacing
  spacing: {
    section: 'space-y-8',
    card: 'space-y-6',
    form: 'space-y-4',
    tight: 'space-y-2',
  },
  
  // Animations
  animations: {
    fadeIn: 'animate-in fade-in duration-300',
    slideIn: 'animate-in slide-in-from-bottom-4 duration-300',
    scaleIn: 'animate-in zoom-in-95 duration-200',
  },
} as const

// BiHub Brand Elements
export const BIHUB_BRAND = {
  name: 'BiHub',
  tagline: 'Your Premium E-commerce Destination',
  logo: {
    text: 'Bi',
    accent: 'hub',
    colors: {
      text: 'text-white',
      accent: 'bg-[#FF9000] text-black',
    },
  },
  
  // Admin-specific branding
  admin: {
    title: 'BiHub Admin',
    subtitle: 'Store Management Dashboard',
    features: [
      'Advanced Analytics',
      'Real-time Monitoring', 
      'Smart Inventory',
      'Customer Insights',
    ],
  },
} as const

// Utility functions for theme
export const getStatusColor = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'completed':
    case 'success':
      return BIHUB_ADMIN_THEME.colors.status.success
    case 'pending':
    case 'processing':
    case 'warning':
      return BIHUB_ADMIN_THEME.colors.status.warning
    case 'inactive':
    case 'failed':
    case 'error':
      return BIHUB_ADMIN_THEME.colors.status.error
    case 'draft':
    case 'info':
    default:
      return BIHUB_ADMIN_THEME.colors.status.info
  }
}

export const getBadgeVariant = (status: string) => {
  switch (status.toLowerCase()) {
    case 'active':
    case 'completed':
    case 'success':
      return 'success'
    case 'pending':
    case 'processing':
    case 'warning':
      return 'warning'
    case 'inactive':
    case 'failed':
    case 'error':
      return 'error'
    case 'draft':
    case 'info':
    default:
      return 'primary'
  }
}
