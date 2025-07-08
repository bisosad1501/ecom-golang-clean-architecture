'use client'

import { useState, useEffect, useRef } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import {
  Search,
  ShoppingCart,
  User,
  Heart,
  Menu,
  X,
  LogOut,
  Settings,
  Package,
  Truck,
  ArrowRight,
  Shield,
  ShoppingBag
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { useCartStore, getCartItemCount } from '@/store/cart'
import { APP_CONFIG } from '@/constants/app'
import { USER_NAV } from '@/constants'
import { getVisibleNavItems } from '@/lib/permissions'
import { RequireAuth, RequireGuest } from '@/components/auth/permission-guard'
import { cn } from '@/lib/utils'
import { CategoryMegaMenu } from './category-mega-menu'
import { DESIGN_TOKENS } from '@/constants/design-tokens'

export function Header() {
  const router = useRouter()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false)
  const [isShoppingMode, setIsShoppingMode] = useState(true) // Toggle between admin and shopping mode
  const userMenuRef = useRef<HTMLDivElement>(null)

  const { user, isAuthenticated, logout, refreshUser } = useAuthStore()
  const { cart, openCart } = useCartStore()

  const cartItemCount = getCartItemCount(cart)
  const visibleNavItems = getVisibleNavItems(user?.role || null)
  const isAdmin = user?.role === 'admin'

  // Close user menu when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (userMenuRef.current && !userMenuRef.current.contains(event.target as Node)) {
        setIsUserMenuOpen(false)
      }
    }

    if (isUserMenuOpen) {
      document.addEventListener('mousedown', handleClickOutside)
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside)
    }
  }, [isUserMenuOpen])

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault()
    if (searchQuery.trim()) {
      router.push(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
      setSearchQuery('')
    }
  }

  const handleLogout = () => {
    logout()
    setIsUserMenuOpen(false)
    router.push('/')
  }

  const handleRefreshUser = async () => {
    try {
      await refreshUser()
      console.log('User data refreshed')
    } catch (error) {
      console.error('Failed to refresh user:', error)
    }
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b border-black/20 bg-black backdrop-blur-lg shadow-lg">
      {/* Top bar */}
      <div className="border-b border-black/30 bg-black py-1.5">{/* Pure black background */}
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between text-xs text-white">{/* White text on dark background */}
            <div className="hidden md:flex items-center space-x-2">{/* Reduced spacing */}
              <div className="flex items-center gap-1.5 rounded-full px-2.5 py-0.5" style={{backgroundColor: 'rgba(255, 144, 0, 0.2)'}}>{/* BiHub orange accent */}
                <Truck className="h-3.5 w-3.5" style={{color: '#FF9000'}} />{/* BiHub orange */}
                <span className="font-semibold" style={{color: '#FF9000'}}>Free shipping on orders over $50</span>
              </div>
            </div>
            <div className="flex items-center space-x-6">{/* Reduced spacing */}
              <Link
                href="/help"
                className="transition-all duration-200 font-medium hover:scale-105 flex items-center gap-1 group"
                onMouseEnter={(e) => e.currentTarget.style.color = '#FF9000'}
                onMouseLeave={(e) => e.currentTarget.style.color = 'white'}
              >{/* BiHub orange hover */}
                <span>Help</span>
                <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all duration-200" />
              </Link>
              <Link
                href="/track-order"
                className="transition-all duration-200 font-medium hover:scale-105 flex items-center gap-1 group"
                onMouseEnter={(e) => e.currentTarget.style.color = '#FF9000'}
                onMouseLeave={(e) => e.currentTarget.style.color = 'white'}
              >
                <span>Track Order</span>
                <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all duration-200" />
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Main header */}
      <div className="container mx-auto px-4">
        <div className="flex h-14 items-center justify-between">{/* Reduced height */}
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-2.5 group">
              <span className="text-2xl font-bold flex items-center">
                <span className="text-white">Bi</span>
                <span className="ml-1 px-2 py-0.5 rounded-[2px] text-black font-bold" style={{letterSpacing: '0.5px', backgroundColor: '#FF9000'}}>hub</span>
              </span>
            </Link>
          </div>

          {/* Enhanced Search bar */}
          <div className="hidden md:flex flex-1 max-w-2xl mx-8">
            <form onSubmit={handleSearch} className="flex w-full relative group">
              <div className="relative flex-1">
                <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-5 w-5 text-gray-400 transition-colors z-10"
                  style={{'--focus-color': '#FF9000'} as React.CSSProperties}
                />
                <Input
                  type="search"
                  placeholder="Search for products, brands, categories..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-12 pr-4 h-11 rounded-l-full rounded-r-none border-0 bg-gray-900 text-white placeholder:text-gray-400 focus:outline-none shadow-lg transition-all duration-200 text-base"
                  style={{
                    '--tw-ring-color': 'rgba(255, 144, 0, 0.4)',
                  } as React.CSSProperties}
                  onFocus={(e) => {
                    e.currentTarget.style.boxShadow = '0 0 0 2px rgba(255, 144, 0, 0.4)';
                    const searchIcon = e.currentTarget.parentElement?.querySelector('.lucide-search') as HTMLElement;
                    if (searchIcon) searchIcon.style.color = '#FF9000';
                  }}
                  onBlur={(e) => {
                    e.currentTarget.style.boxShadow = '';
                    const searchIcon = e.currentTarget.parentElement?.querySelector('.lucide-search') as HTMLElement;
                    if (searchIcon) searchIcon.style.color = '#9ca3af';
                  }}
                />
              </div>
              <Button
                type="submit"
                className="rounded-r-full rounded-l-none h-11 px-6 text-white text-base font-semibold shadow-lg transition-all duration-200 border-0"
                style={{
                  backgroundColor: '#FF9000',
                  boxShadow: '0 2px 8px 0 rgba(255,144,0,0.08)'
                }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = '#e67e00'}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = '#FF9000'}
              >
                Search
              </Button>
            </form>
          </div>

          {/* Right side actions */}
          <div className={`flex items-center ${DESIGN_TOKENS.SPACING.CLASSES.GAP_SMALL}`}>{/* Smaller spacing */}
            {/* Mobile search */}
            <Button
              variant="ghost"
              size="icon"
              className={`md:hidden h-10 w-10 rounded-xl hover:bg-orange-500/10 hover:scale-105 transition-all duration-200 text-white`}
              onClick={() => {/* TODO: Open mobile search */}}
            >
              <Search className="h-4 w-4" />
            </Button>

            {/* Wishlist with preview - Show for all authenticated users or admin in shopping mode */}
            <RequireAuth>
              {(!isAdmin || isShoppingMode) && (
                <div className="relative group">
                  <Button variant="ghost" size="icon" className="relative h-10 w-10 rounded-xl hover:bg-orange-500/10 hover:scale-105 transition-all duration-200 text-white">
                    <Heart className="h-4 w-4 group-hover:text-orange-500 transition-colors text-white" />
                    {/* Add wishlist count badge if needed */}
                    {/* <Badge variant="default" className="absolute -top-1 -right-1 h-6 w-6 rounded-full p-0 text-xs font-bold">0</Badge> */}
                  </Button>

                  {/* Wishlist preview on hover */}
                  <div className="absolute top-full right-0 mt-2 w-80 bg-background border border-border rounded-2xl shadow-2xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-50">
                    <div className="p-6">
                      <div className="flex items-center justify-between mb-4">
                        <h3 className="font-semibold text-lg">Wishlist</h3>
                        <Badge variant="secondary">0 items</Badge>
                      </div>

                      {/* Empty wishlist state */}
                      <div className="text-center py-8">
                        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                          <Heart className="h-8 w-8 text-gray-400" />
                        </div>
                        <h4 className="font-medium text-gray-900 mb-2">Your wishlist is empty</h4>
                        <p className="text-sm text-gray-500 mb-4">Save items you love for later</p>
                        <Button variant="gradient" onClick={() => router.push('/products')}>
                          Discover Products
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </RequireAuth>

            {/* Enhanced Cart with preview - Show for all authenticated users or admin in shopping mode */}
            {(!isAdmin || isShoppingMode) && (
              <div className="relative group">
                <Button
                  variant="ghost"
                  size="icon"
                  className="relative h-10 w-10 rounded-xl hover:bg-orange-500/10 hover:scale-105 transition-all duration-200 text-white"
                  onClick={() => router.push('/cart')}
                >
                  <ShoppingCart className="h-4 w-4 group-hover:text-orange-500 transition-colors text-white" />
                  {cartItemCount > 0 && (
                    <Badge
                      variant="default"
                      className="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 text-xs font-bold shadow-large animate-pulse flex items-center justify-center bg-orange-500 text-white border-0"
                    >
                      {cartItemCount}
                    </Badge>
                  )}
                </Button>

                {/* Cart preview on hover */}
                <div className="absolute top-full right-0 mt-2 w-80 bg-background border border-border rounded-2xl shadow-2xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-50">
                  <div className="p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="font-semibold text-lg">Shopping Cart</h3>
                      <Badge variant="secondary">{cartItemCount} items</Badge>
                    </div>

                    {cartItemCount > 0 && cart ? (
                      <>
                        <div className="space-y-3 max-h-64 overflow-y-auto">
                          {cart.items.slice(0, 3).map((item) => (
                            <div key={item.id} className="flex items-center gap-3 p-2 hover:bg-muted rounded-lg transition-colors">
                              <div className="w-12 h-12 bg-muted rounded-lg overflow-hidden">
                                {item.product.images && item.product.images.length > 0 ? (
                                  <img 
                                    src={item.product.images[0].url || item.product.images[0].image_url || '/placeholder-product.svg'} 
                                    alt={item.product.name}
                                    className="w-full h-full object-cover"
                                    onError={(e) => {
                                      const target = e.target as HTMLImageElement;
                                      target.src = '/placeholder-product.svg';
                                    }}
                                  />
                                ) : (
                                  <div className="w-full h-full bg-gray-200 flex items-center justify-center">
                                    <Package className="h-6 w-6 text-gray-400" />
                                  </div>
                                )}
                              </div>
                              <div className="flex-1 min-w-0">
                                <div className="font-medium text-sm truncate">{item.product.name}</div>
                                <div className="text-xs text-muted-foreground">Qty: {item.quantity}</div>
                              </div>
                              <div className="text-sm font-semibold">${item.product.price}</div>
                            </div>
                          ))}
                        </div>

                        {cart.items.length > 3 && (
                          <div className="text-center text-sm text-muted-foreground mt-3">
                            +{cart.items.length - 3} more items
                          </div>
                        )}

                        <div className="border-t border-border mt-4 pt-4">
                          <div className="flex items-center justify-between mb-3">
                            <span className="font-semibold">Total:</span>
                            <span className="font-bold text-lg text-primary">
                              ${cart.items.reduce((sum, item) => sum + (item.product.price * item.quantity), 0).toFixed(2)}
                            </span>
                          </div>
                          <Button className="w-full" variant="gradient" onClick={() => router.push('/cart')}>
                            View Cart
                          </Button>
                        </div>
                      </>
                    ) : (
                      <div className="text-center py-8">
                        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
                          <ShoppingCart className="h-8 w-8 text-gray-400" />
                        </div>
                        <h4 className="font-medium text-gray-900 mb-2">Your cart is empty</h4>
                        <p className="text-sm text-gray-500 mb-4">Add some products to get started</p>
                        <Button variant="gradient" onClick={() => router.push('/products')}>
                          Continue Shopping
                        </Button>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}

            {/* Admin Mode Toggle - Only show for admin users */}
            {isAdmin && (
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setIsShoppingMode(!isShoppingMode)}
                className={cn(
                  "relative h-10 w-10 rounded-xl transition-all duration-200",
                  isShoppingMode 
                    ? "bg-emerald-100 hover:bg-emerald-200 text-emerald-700" 
                    : "bg-blue-100 hover:bg-blue-200 text-blue-700"
                )}
                title={isShoppingMode ? "Switch to Admin Mode" : "Switch to Shopping Mode"}
              >
                {isShoppingMode ? (
                  <ShoppingBag className="h-4 w-4" />
                ) : (
                  <Shield className="h-4 w-4" />
                )}
              </Button>
            )}

            {/* User menu */}
            {isAuthenticated ? (
              <div className="relative" ref={userMenuRef}>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                  className="relative h-10 w-10 rounded-xl hover:bg-orange-500/10 hover:scale-105 transition-all duration-200 text-white"
                >
                  <User className="h-4 w-4" />
                </Button>

                {/* User dropdown */}
                {isUserMenuOpen && (
                  <div className="absolute right-0 mt-2 w-48 rounded-2xl border bg-white py-1 shadow-2xl z-[60] backdrop-blur-sm">
                    <div className="px-4 py-3 border-b border-gray-100">
                      <p className="text-sm font-semibold text-gray-900">
                        {user?.first_name} {user?.last_name}
                      </p>
                      <p className="text-xs text-gray-500">{user?.email}</p>
                      {isAdmin && (
                        <div className={cn(
                          "inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium mt-2",
                          isShoppingMode 
                            ? "bg-emerald-100 text-emerald-700" 
                            : "bg-blue-100 text-blue-700"
                        )}>
                          {isShoppingMode ? (
                            <>
                              <ShoppingBag className="h-3 w-3" />
                              Shopping Mode
                            </>
                          ) : (
                            <>
                              <Shield className="h-3 w-3" />
                              Admin Mode
                            </>
                          )}
                        </div>
                      )}
                    </div>
                    
                    {USER_NAV.map((item) => (
                      <Link
                        key={item.href}
                        href={item.href}
                        className="flex items-center px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-primary-600 transition-colors"
                        onClick={() => setIsUserMenuOpen(false)}
                      >
                        {item.icon === 'User' && <User className="mr-3 h-4 w-4" />}
                        {item.icon === 'Package' && <Package className="mr-3 h-4 w-4" />}
                        {item.icon === 'Heart' && <Heart className="mr-3 h-4 w-4" />}
                        {item.icon === 'Settings' && <Settings className="mr-3 h-4 w-4" />}
                        {item.title}
                      </Link>
                    ))}

                    {/* Admin Panel Link */}
                    {isAdmin && (
                      <>
                        <Link
                          href="/admin"
                          className="flex items-center px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-primary-600 border-t border-gray-100 transition-colors"
                          onClick={() => setIsUserMenuOpen(false)}
                        >
                          <Settings className="mr-3 h-4 w-4" />
                          Admin Panel
                        </Link>
                        
                        <button
                          onClick={() => {
                            setIsShoppingMode(!isShoppingMode)
                            setIsUserMenuOpen(false)
                          }}
                          className="flex w-full items-center px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-primary-600 transition-colors"
                        >
                          {isShoppingMode ? (
                            <>
                              <Shield className="mr-3 h-4 w-4" />
                              Switch to Admin Mode
                            </>
                          ) : (
                            <>
                              <ShoppingBag className="mr-3 h-4 w-4" />
                              Switch to Shopping Mode
                            </>
                          )}
                        </button>
                      </>
                    )}
                    
                    <button
                      onClick={handleLogout}
                      className="flex w-full items-center px-4 py-2.5 text-sm text-gray-700 hover:bg-gray-50 hover:text-red-600 border-t border-gray-100 transition-colors"
                    >
                      <LogOut className="mr-3 h-4 w-4" />
                      Sign out
                    </button>
                  </div>
                )}
              </div>
            ) : (
              <RequireGuest>
                <div className="flex items-center space-x-2">
                  <Link href="/auth/login">
                    <Button variant="ghost" size="sm">
                      Sign in
                    </Button>
                  </Link>
                  <Link href="/auth/register">
                    <Button size="sm">
                      Sign up
                    </Button>
                  </Link>
                </div>
              </RequireGuest>
            )}

            {/* Mobile menu button */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden h-10 w-10 rounded-xl text-white hover:bg-orange-500/10"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              {isMenuOpen ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
            </Button>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <div className="border-t border-gray-700/30 bg-black">
        <div className="container mx-auto px-4">
          <div className="flex items-center h-14">
            {/* Category Mega Menu */}
            <CategoryMegaMenu />
            
            {/* Navigation Links */}
            <nav className="hidden md:flex items-center space-x-1 flex-1 ml-6">
              <Link
                href="/"
                className="h-14 flex items-center px-4 text-sm font-medium text-white hover:text-orange-500 hover:bg-orange-500/10 border-b-2 border-transparent hover:border-orange-500 transition-all duration-200"
              >
                Home
              </Link>
              <Link
                href="/products"
                className="h-14 flex items-center px-4 text-sm font-medium text-white hover:text-orange-500 hover:bg-orange-500/10 border-b-2 border-transparent hover:border-orange-500 transition-all duration-200"
              >
                Products
              </Link>
              <Link
                href="/about"
                className="h-14 flex items-center px-4 text-sm font-medium text-white hover:text-orange-500 hover:bg-orange-500/10 border-b-2 border-transparent hover:border-orange-500 transition-all duration-200"
              >
                About
              </Link>
              <Link
                href="/contact"
                className="h-14 flex items-center px-4 text-sm font-medium text-white hover:text-orange-500 hover:bg-orange-500/10 border-b-2 border-transparent hover:border-orange-500 transition-all duration-200"
              >
                Contact
              </Link>
            </nav>

            
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {isMenuOpen && (
        <div className="md:hidden border-t border-gray-700/60 bg-black shadow-lg">
          <div className="container mx-auto px-4 py-4">
            {/* Mobile search */}
            <form onSubmit={handleSearch} className="mb-4">
              <div className="flex">
                <Input
                  type="search"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="rounded-r-none border-r-0 bg-gray-900 border-gray-600 text-white placeholder:text-gray-400"
                />
                <Button type="submit" className="rounded-l-none bg-orange-500 hover:bg-orange-600">
                  <Search className="h-4 w-4" />
                </Button>
              </div>
            </form>

            {/* Mobile navigation */}
            <div className="space-y-2">
              <Link
                href="/"
                className="block py-2 text-base font-medium text-gray-300 hover:text-orange-500"
                onClick={() => setIsMenuOpen(false)}
              >
                Home
              </Link>
              <Link
                href="/products"
                className="block py-2 text-base font-medium text-gray-300 hover:text-orange-500"
                onClick={() => setIsMenuOpen(false)}
              >
                Products
              </Link>
              <Link
                href="/about"
                className="block py-2 text-base font-medium text-gray-300 hover:text-orange-500"
                onClick={() => setIsMenuOpen(false)}
              >
                About
              </Link>
              <Link
                href="/contact"
                className="block py-2 text-base font-medium text-gray-300 hover:text-orange-500"
                onClick={() => setIsMenuOpen(false)}
              >
                Contact
              </Link>
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
