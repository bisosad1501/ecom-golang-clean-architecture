// ===== UNIFIED COMPONENTS EXPORT =====
// Atomic Design Pattern Implementation

// Export all component levels
export * from './atoms'
export * from './molecules'
export * from './organisms'
export * from './templates'
export * from './pages'

// Export providers and utilities
export * from './providers'

// Re-export commonly used components for convenience
// Atoms
export { Button, Input, Label, Badge, Avatar, Icon, Spinner } from './atoms'

// Molecules  
export { FormField, SearchBar, ProductCard, Pagination, DataTable } from './molecules'

// Organisms
export { Header, Footer, ProductGrid, ShoppingCart, LoginForm } from './organisms'

// Templates
export { MainLayout, AdminLayout, AuthLayout, ProductPageTemplate } from './templates'

// Pages
export { HomePage, ProductsPage, ProductDetailPage, LoginPage, AdminDashboardPage } from './pages'

// Legacy exports for backward compatibility
// These will be deprecated in future versions
export { Button as UIButton } from './ui/button'
export { Input as UIInput } from './ui/input'
export { Card as UICard } from './ui/card'
export { Badge as UIBadge } from './ui/badge'

// Component categories for easier imports
export const Atoms = {
  Button: () => import('./atoms/Button'),
  Input: () => import('./atoms/Input'),
  Label: () => import('./atoms/Label'),
  Badge: () => import('./atoms/Badge'),
  Avatar: () => import('./atoms/Avatar'),
  Icon: () => import('./atoms/Icon'),
  Spinner: () => import('./atoms/Spinner'),
  Skeleton: () => import('./atoms/Skeleton'),
  Separator: () => import('./atoms/Separator'),
  // Add more atoms as needed
}

export const Molecules = {
  FormField: () => import('./molecules/FormField'),
  SearchBar: () => import('./molecules/SearchBar'),
  ProductCard: () => import('./molecules/ProductCard'),
  CategoryCard: () => import('./molecules/CategoryCard'),
  Pagination: () => import('./molecules/Pagination'),
  DataTable: () => import('./molecules/DataTable'),
  ImageGallery: () => import('./molecules/ImageGallery'),
  // Add more molecules as needed
}

export const Organisms = {
  Header: () => import('./organisms/Header'),
  Footer: () => import('./organisms/Footer'),
  Navigation: () => import('./organisms/Navigation'),
  HeroSection: () => import('./organisms/HeroSection'),
  ProductGrid: () => import('./organisms/ProductGrid'),
  ProductDetails: () => import('./organisms/ProductDetails'),
  ShoppingCart: () => import('./organisms/ShoppingCart'),
  CheckoutForm: () => import('./organisms/CheckoutForm'),
  LoginForm: () => import('./organisms/LoginForm'),
  SearchResults: () => import('./organisms/SearchResults'),
  // Add more organisms as needed
}

export const Templates = {
  MainLayout: () => import('./templates/MainLayout'),
  AdminLayout: () => import('./templates/AdminLayout'),
  AuthLayout: () => import('./templates/AuthLayout'),
  ProductPageTemplate: () => import('./templates/ProductPageTemplate'),
  CategoryPageTemplate: () => import('./templates/CategoryPageTemplate'),
  CheckoutPageTemplate: () => import('./templates/CheckoutPageTemplate'),
  // Add more templates as needed
}

export const Pages = {
  HomePage: () => import('./pages/HomePage'),
  ProductsPage: () => import('./pages/ProductsPage'),
  ProductDetailPage: () => import('./pages/ProductDetailPage'),
  CategoryPage: () => import('./pages/CategoryPage'),
  SearchPage: () => import('./pages/SearchPage'),
  CartPage: () => import('./pages/CartPage'),
  CheckoutPage: () => import('./pages/CheckoutPage'),
  LoginPage: () => import('./pages/LoginPage'),
  RegisterPage: () => import('./pages/RegisterPage'),
  ProfilePage: () => import('./pages/ProfilePage'),
  AdminDashboardPage: () => import('./pages/AdminDashboardPage'),
  // Add more pages as needed
}

// Component registry for dynamic loading
export const ComponentRegistry = {
  // Atoms
  'atoms/Button': () => import('./atoms/Button'),
  'atoms/Input': () => import('./atoms/Input'),
  'atoms/Badge': () => import('./atoms/Badge'),
  'atoms/Avatar': () => import('./atoms/Avatar'),
  'atoms/Icon': () => import('./atoms/Icon'),
  
  // Molecules
  'molecules/ProductCard': () => import('./molecules/ProductCard'),
  'molecules/SearchBar': () => import('./molecules/SearchBar'),
  'molecules/FormField': () => import('./molecules/FormField'),
  'molecules/Pagination': () => import('./molecules/Pagination'),
  
  // Organisms
  'organisms/Header': () => import('./organisms/Header'),
  'organisms/Footer': () => import('./organisms/Footer'),
  'organisms/ProductGrid': () => import('./organisms/ProductGrid'),
  'organisms/ShoppingCart': () => import('./organisms/ShoppingCart'),
  
  // Templates
  'templates/MainLayout': () => import('./templates/MainLayout'),
  'templates/AdminLayout': () => import('./templates/AdminLayout'),
  'templates/ProductPageTemplate': () => import('./templates/ProductPageTemplate'),
  
  // Pages
  'pages/HomePage': () => import('./pages/HomePage'),
  'pages/ProductsPage': () => import('./pages/ProductsPage'),
  'pages/LoginPage': () => import('./pages/LoginPage'),
}

// Utility function to dynamically load components
export async function loadComponent(componentPath: string) {
  const loader = ComponentRegistry[componentPath as keyof typeof ComponentRegistry]
  if (!loader) {
    throw new Error(`Component not found: ${componentPath}`)
  }
  return await loader()
}

// Component metadata for documentation and tooling
export const ComponentMetadata = {
  atoms: {
    count: Object.keys(Atoms).length,
    description: 'Basic building blocks - smallest components',
    examples: ['Button', 'Input', 'Badge', 'Icon']
  },
  molecules: {
    count: Object.keys(Molecules).length,
    description: 'Combinations of atoms that function together as a unit',
    examples: ['ProductCard', 'SearchBar', 'FormField', 'Pagination']
  },
  organisms: {
    count: Object.keys(Organisms).length,
    description: 'Complex components made of molecules and atoms',
    examples: ['Header', 'ProductGrid', 'ShoppingCart', 'CheckoutForm']
  },
  templates: {
    count: Object.keys(Templates).length,
    description: 'Page-level layouts that combine organisms, molecules, and atoms',
    examples: ['MainLayout', 'AdminLayout', 'ProductPageTemplate']
  },
  pages: {
    count: Object.keys(Pages).length,
    description: 'Complete page components that use templates and handle routing',
    examples: ['HomePage', 'ProductsPage', 'LoginPage', 'AdminDashboardPage']
  }
}

// Design system information
export const DesignSystem = {
  name: 'ShopHub Design System',
  version: '1.0.0',
  pattern: 'Atomic Design',
  description: 'A comprehensive design system built with Atomic Design principles for e-commerce applications',
  levels: ['atoms', 'molecules', 'organisms', 'templates', 'pages'],
  theme: {
    colors: 'Unified color system with primary orange theme',
    typography: 'Inter font family with consistent scale',
    spacing: 'Consistent spacing scale based on 4px grid',
    components: 'Reusable components following design tokens'
  },
  features: [
    'TypeScript support',
    'Tailwind CSS integration',
    'Dark/Light theme support',
    'Responsive design',
    'Accessibility compliance',
    'Performance optimized',
    'Storybook documentation',
    'Unit test coverage'
  ]
}

// Export design system info
export { DesignSystem as designSystem }

// Type exports for better TypeScript support
export type {
  // Atom types
  ButtonProps,
  InputProps,
  BadgeProps,
  AvatarProps,
  IconProps,
  
  // Molecule types
  FormFieldProps,
  SearchBarProps,
  ProductCardProps,
  PaginationProps,
  
  // Organism types
  HeaderProps,
  FooterProps,
  ProductGridProps,
  ShoppingCartProps,
  
  // Template types
  MainLayoutProps,
  AdminLayoutProps,
  ProductPageTemplateProps,
  
  // Page types
  HomePageProps,
  ProductsPageProps,
  LoginPageProps,
} from './atoms'

// Component factory for creating custom instances
export class ComponentFactory {
  static createButton(props: any) {
    return import('./atoms/Button').then(({ Button }) => Button)
  }
  
  static createProductCard(props: any) {
    return import('./molecules/ProductCard').then(({ ProductCard }) => ProductCard)
  }
  
  static createLayout(type: 'main' | 'admin' | 'auth') {
    switch (type) {
      case 'main':
        return import('./templates/MainLayout').then(({ MainLayout }) => MainLayout)
      case 'admin':
        return import('./templates/AdminLayout').then(({ AdminLayout }) => AdminLayout)
      case 'auth':
        return import('./templates/AuthLayout').then(({ AuthLayout }) => AuthLayout)
      default:
        throw new Error(`Unknown layout type: ${type}`)
    }
  }
}

// Export component factory
export { ComponentFactory as componentFactory }
