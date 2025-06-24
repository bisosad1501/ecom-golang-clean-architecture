import { UserRole } from '@/types'

// Permission constants
export const PERMISSIONS = {
  // Product permissions
  PRODUCTS_VIEW: 'products:view',
  PRODUCTS_CREATE: 'products:create',
  PRODUCTS_UPDATE: 'products:update',
  PRODUCTS_DELETE: 'products:delete',
  PRODUCTS_MANAGE_INVENTORY: 'products:manage_inventory',
  
  // Category permissions
  CATEGORIES_VIEW: 'categories:view',
  CATEGORIES_CREATE: 'categories:create',
  CATEGORIES_UPDATE: 'categories:update',
  CATEGORIES_DELETE: 'categories:delete',
  
  // Order permissions
  ORDERS_VIEW_ALL: 'orders:view_all',
  ORDERS_VIEW_OWN: 'orders:view_own',
  ORDERS_UPDATE: 'orders:update',
  ORDERS_CANCEL: 'orders:cancel',
  ORDERS_REFUND: 'orders:refund',
  
  // User permissions
  USERS_VIEW_ALL: 'users:view_all',
  USERS_VIEW_OWN: 'users:view_own',
  USERS_UPDATE_ALL: 'users:update_all',
  USERS_UPDATE_OWN: 'users:update_own',
  USERS_DELETE: 'users:delete',
  USERS_MANAGE_ROLES: 'users:manage_roles',
  
  // Review permissions
  REVIEWS_VIEW_ALL: 'reviews:view_all',
  REVIEWS_CREATE: 'reviews:create',
  REVIEWS_UPDATE_OWN: 'reviews:update_own',
  REVIEWS_DELETE_OWN: 'reviews:delete_own',
  REVIEWS_MODERATE: 'reviews:moderate',
  
  // Coupon permissions
  COUPONS_VIEW: 'coupons:view',
  COUPONS_CREATE: 'coupons:create',
  COUPONS_UPDATE: 'coupons:update',
  COUPONS_DELETE: 'coupons:delete',
  
  // Analytics permissions
  ANALYTICS_VIEW: 'analytics:view',
  ANALYTICS_EXPORT: 'analytics:export',
  
  // System permissions
  SYSTEM_SETTINGS: 'system:settings',
  SYSTEM_LOGS: 'system:logs',
  SYSTEM_BACKUP: 'system:backup',
} as const

// Role-based permissions mapping
export const ROLE_PERMISSIONS: Record<UserRole, string[]> = {
  customer: [
    PERMISSIONS.PRODUCTS_VIEW,
    PERMISSIONS.CATEGORIES_VIEW,
    PERMISSIONS.ORDERS_VIEW_OWN,
    PERMISSIONS.ORDERS_CANCEL,
    PERMISSIONS.USERS_VIEW_OWN,
    PERMISSIONS.USERS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_CREATE,
    PERMISSIONS.REVIEWS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_DELETE_OWN,
  ],

  moderator: [
    // Customer permissions
    PERMISSIONS.PRODUCTS_VIEW,
    PERMISSIONS.CATEGORIES_VIEW,
    PERMISSIONS.ORDERS_VIEW_OWN,
    PERMISSIONS.ORDERS_CANCEL,
    PERMISSIONS.USERS_VIEW_OWN,
    PERMISSIONS.USERS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_CREATE,
    PERMISSIONS.REVIEWS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_DELETE_OWN,
    // Additional moderator permissions
    PERMISSIONS.REVIEWS_VIEW_ALL,
    PERMISSIONS.REVIEWS_MODERATE,
    PERMISSIONS.ORDERS_VIEW_ALL,
    PERMISSIONS.USERS_VIEW_ALL,
  ],

  admin: [
    // Customer permissions
    PERMISSIONS.PRODUCTS_VIEW,
    PERMISSIONS.CATEGORIES_VIEW,
    PERMISSIONS.ORDERS_VIEW_OWN,
    PERMISSIONS.ORDERS_CANCEL,
    PERMISSIONS.USERS_VIEW_OWN,
    PERMISSIONS.USERS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_CREATE,
    PERMISSIONS.REVIEWS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_DELETE_OWN,
    // Moderator permissions
    PERMISSIONS.REVIEWS_VIEW_ALL,
    PERMISSIONS.REVIEWS_MODERATE,
    PERMISSIONS.ORDERS_VIEW_ALL,
    PERMISSIONS.USERS_VIEW_ALL,
    // Additional admin permissions
    PERMISSIONS.PRODUCTS_CREATE,
    PERMISSIONS.PRODUCTS_UPDATE,
    PERMISSIONS.PRODUCTS_DELETE,
    PERMISSIONS.PRODUCTS_MANAGE_INVENTORY,
    PERMISSIONS.CATEGORIES_CREATE,
    PERMISSIONS.CATEGORIES_UPDATE,
    PERMISSIONS.CATEGORIES_DELETE,
    PERMISSIONS.ORDERS_UPDATE,
    PERMISSIONS.ORDERS_REFUND,
    PERMISSIONS.USERS_UPDATE_ALL,
    PERMISSIONS.COUPONS_VIEW,
    PERMISSIONS.COUPONS_CREATE,
    PERMISSIONS.COUPONS_UPDATE,
    PERMISSIONS.COUPONS_DELETE,
    PERMISSIONS.ANALYTICS_VIEW,
    PERMISSIONS.ANALYTICS_EXPORT,
  ],

  super_admin: [
    // Customer permissions
    PERMISSIONS.PRODUCTS_VIEW,
    PERMISSIONS.CATEGORIES_VIEW,
    PERMISSIONS.ORDERS_VIEW_OWN,
    PERMISSIONS.ORDERS_CANCEL,
    PERMISSIONS.USERS_VIEW_OWN,
    PERMISSIONS.USERS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_CREATE,
    PERMISSIONS.REVIEWS_UPDATE_OWN,
    PERMISSIONS.REVIEWS_DELETE_OWN,
    // Moderator permissions
    PERMISSIONS.REVIEWS_VIEW_ALL,
    PERMISSIONS.REVIEWS_MODERATE,
    PERMISSIONS.ORDERS_VIEW_ALL,
    PERMISSIONS.USERS_VIEW_ALL,
    // Admin permissions
    PERMISSIONS.PRODUCTS_CREATE,
    PERMISSIONS.PRODUCTS_UPDATE,
    PERMISSIONS.PRODUCTS_DELETE,
    PERMISSIONS.PRODUCTS_MANAGE_INVENTORY,
    PERMISSIONS.CATEGORIES_CREATE,
    PERMISSIONS.CATEGORIES_UPDATE,
    PERMISSIONS.CATEGORIES_DELETE,
    PERMISSIONS.ORDERS_UPDATE,
    PERMISSIONS.ORDERS_REFUND,
    PERMISSIONS.USERS_UPDATE_ALL,
    PERMISSIONS.COUPONS_VIEW,
    PERMISSIONS.COUPONS_CREATE,
    PERMISSIONS.COUPONS_UPDATE,
    PERMISSIONS.COUPONS_DELETE,
    PERMISSIONS.ANALYTICS_VIEW,
    PERMISSIONS.ANALYTICS_EXPORT,
    // Additional super admin permissions
    PERMISSIONS.USERS_DELETE,
    PERMISSIONS.USERS_MANAGE_ROLES,
    PERMISSIONS.SYSTEM_SETTINGS,
    PERMISSIONS.SYSTEM_LOGS,
    PERMISSIONS.SYSTEM_BACKUP,
  ],
}

// Helper functions
export function hasPermission(userRole: UserRole, permission: string): boolean {
  const rolePermissions = ROLE_PERMISSIONS[userRole] || []
  return rolePermissions.includes(permission)
}

export function hasAnyPermission(userRole: UserRole, permissions: string[]): boolean {
  return permissions.some(permission => hasPermission(userRole, permission))
}

export function hasAllPermissions(userRole: UserRole, permissions: string[]): boolean {
  return permissions.every(permission => hasPermission(userRole, permission))
}

export function canAccessAdminPanel(userRole: UserRole): boolean {
  return hasAnyPermission(userRole, [
    PERMISSIONS.PRODUCTS_CREATE,
    PERMISSIONS.ORDERS_VIEW_ALL,
    PERMISSIONS.USERS_VIEW_ALL,
    PERMISSIONS.ANALYTICS_VIEW,
  ])
}

export function canManageProducts(userRole: UserRole): boolean {
  return hasAnyPermission(userRole, [
    PERMISSIONS.PRODUCTS_CREATE,
    PERMISSIONS.PRODUCTS_UPDATE,
    PERMISSIONS.PRODUCTS_DELETE,
  ])
}

export function canManageOrders(userRole: UserRole): boolean {
  return hasAnyPermission(userRole, [
    PERMISSIONS.ORDERS_VIEW_ALL,
    PERMISSIONS.ORDERS_UPDATE,
    PERMISSIONS.ORDERS_REFUND,
  ])
}

export function canManageUsers(userRole: UserRole): boolean {
  return hasAnyPermission(userRole, [
    PERMISSIONS.USERS_VIEW_ALL,
    PERMISSIONS.USERS_UPDATE_ALL,
    PERMISSIONS.USERS_DELETE,
  ])
}

export function canModerateReviews(userRole: UserRole): boolean {
  return hasPermission(userRole, PERMISSIONS.REVIEWS_MODERATE)
}

export function canViewAnalytics(userRole: UserRole): boolean {
  return hasPermission(userRole, PERMISSIONS.ANALYTICS_VIEW)
}

// Route protection helpers
export const PROTECTED_ROUTES = {
  ADMIN_ONLY: ['/admin', '/admin/*'],
  AUTHENTICATED_ONLY: ['/account', '/account/*', '/checkout', '/orders', '/orders/*'],
  GUEST_ONLY: ['/auth/login', '/auth/register'],
} as const

export function canAccessRoute(userRole: UserRole | null, path: string): boolean {
  // Check if route requires admin access
  if (PROTECTED_ROUTES.ADMIN_ONLY.some(route => 
    path.startsWith(route.replace('/*', ''))
  )) {
    return userRole ? canAccessAdminPanel(userRole) : false
  }
  
  // Check if route requires authentication
  if (PROTECTED_ROUTES.AUTHENTICATED_ONLY.some(route => 
    path.startsWith(route.replace('/*', ''))
  )) {
    return userRole !== null
  }
  
  // Check if route is for guests only
  if (PROTECTED_ROUTES.GUEST_ONLY.includes(path)) {
    return userRole === null
  }
  
  // Public routes
  return true
}

// UI permission helpers
export function getVisibleNavItems(userRole: UserRole | null) {
  const items = [
    { href: '/', label: 'Home', public: true },
    { href: '/products', label: 'Products', public: true },
    { href: '/categories', label: 'Categories', public: true },
  ]
  
  if (userRole) {
    items.push(
      { href: '/account', label: 'My Account', public: false },
      { href: '/orders', label: 'My Orders', public: false },
      { href: '/wishlist', label: 'Wishlist', public: false }
    )
    
    if (canAccessAdminPanel(userRole)) {
      items.push({ href: '/admin', label: 'Admin Panel', public: false })
    }
  }
  
  return items
}

export function getAdminSidebarItems(userRole: UserRole) {
  const items = []
  
  if (hasPermission(userRole, PERMISSIONS.ANALYTICS_VIEW)) {
    items.push({ href: '/admin', label: 'Dashboard', icon: 'BarChart3' })
  }
  
  if (canManageProducts(userRole)) {
    items.push({ href: '/admin/products', label: 'Products', icon: 'Package' })
  }
  
  if (hasAnyPermission(userRole, [PERMISSIONS.CATEGORIES_CREATE, PERMISSIONS.CATEGORIES_UPDATE])) {
    items.push({ href: '/admin/categories', label: 'Categories', icon: 'Folder' })
  }
  
  if (canManageOrders(userRole)) {
    items.push({ href: '/admin/orders', label: 'Orders', icon: 'ShoppingCart' })
  }
  
  if (canManageUsers(userRole)) {
    items.push({ href: '/admin/users', label: 'Users', icon: 'Users' })
  }
  
  if (canModerateReviews(userRole)) {
    items.push({ href: '/admin/reviews', label: 'Reviews', icon: 'Star' })
  }
  
  if (hasAnyPermission(userRole, [PERMISSIONS.COUPONS_VIEW, PERMISSIONS.COUPONS_CREATE])) {
    items.push({ href: '/admin/coupons', label: 'Coupons', icon: 'Ticket' })
  }
  
  if (hasPermission(userRole, PERMISSIONS.SYSTEM_SETTINGS)) {
    items.push({ href: '/admin/settings', label: 'Settings', icon: 'Settings' })
  }
  
  return items
}
