// Permission constants for role-based access control
export const PERMISSIONS = {
  // Admin permissions
  ADMIN_ACCESS: 'admin.access',
  ADMIN_USERS: 'admin.users',
  ADMIN_ANALYTICS: 'admin.analytics',
  ADMIN_SYSTEM: 'admin.system',
  
  // Moderator permissions
  MODERATOR_PRODUCTS: 'moderator.products',
  MODERATOR_CATEGORIES: 'moderator.categories',
  MODERATOR_ORDERS: 'moderator.orders',
  
  // User permissions
  USER_PROFILE: 'user.profile',
  USER_ORDERS: 'user.orders',
  USER_CART: 'user.cart',
  USER_WISHLIST: 'user.wishlist',
} as const

export type Permission = typeof PERMISSIONS[keyof typeof PERMISSIONS]

// Role-based permission mapping (matches backend roles)
export const ROLE_PERMISSIONS = {
  admin: [
    PERMISSIONS.ADMIN_ACCESS,
    PERMISSIONS.ADMIN_USERS,
    PERMISSIONS.ADMIN_ANALYTICS,
    PERMISSIONS.ADMIN_SYSTEM,
    PERMISSIONS.MODERATOR_PRODUCTS,
    PERMISSIONS.MODERATOR_CATEGORIES,
    PERMISSIONS.MODERATOR_ORDERS,
    PERMISSIONS.USER_PROFILE,
    PERMISSIONS.USER_ORDERS,
    PERMISSIONS.USER_CART,
    PERMISSIONS.USER_WISHLIST,
  ],
  moderator: [
    PERMISSIONS.MODERATOR_PRODUCTS,
    PERMISSIONS.MODERATOR_CATEGORIES,
    PERMISSIONS.MODERATOR_ORDERS,
    PERMISSIONS.USER_PROFILE,
    PERMISSIONS.USER_ORDERS,
    PERMISSIONS.USER_CART,
    PERMISSIONS.USER_WISHLIST,
  ],
  customer: [
    PERMISSIONS.USER_PROFILE,
    PERMISSIONS.USER_ORDERS,
    PERMISSIONS.USER_CART,
    PERMISSIONS.USER_WISHLIST,
  ],
} as const

export type UserRole = keyof typeof ROLE_PERMISSIONS
