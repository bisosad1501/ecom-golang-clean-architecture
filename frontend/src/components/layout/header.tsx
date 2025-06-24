'use client'

import { useState } from 'react'
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
  Package
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { useCartStore, getCartItemCount } from '@/store/cart'
import { APP_NAME, MAIN_NAV, USER_NAV } from '@/constants'
import { getVisibleNavItems, canAccessAdminPanel } from '@/lib/permissions'
import { RequireAuth, RequireGuest, RequireAdmin } from '@/components/auth/permission-guard'
import { cn } from '@/lib/utils'

export function Header() {
  const router = useRouter()
  const [isMenuOpen, setIsMenuOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false)

  const { user, isAuthenticated, logout, refreshUser } = useAuthStore()
  const { cart, openCart } = useCartStore()

  const cartItemCount = getCartItemCount(cart)
  const visibleNavItems = getVisibleNavItems(user?.role || null)

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
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 shadow-sm">
      {/* Top bar */}
      <div className="border-b border-border bg-muted py-2">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <div className="hidden md:block">
              Free shipping on orders over $50
            </div>
            <div className="flex items-center space-x-4">
              <Link href="/help" className="hover:text-foreground">
                Help
              </Link>
              <Link href="/contact" className="hover:text-foreground">
                Contact
              </Link>
              <Link href="/track-order" className="hover:text-foreground">
                Track Order
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Main header */}
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-2">
              <div className="h-8 w-8 rounded-lg bg-primary flex items-center justify-center">
                <span className="text-primary-foreground font-bold text-lg">E</span>
              </div>
              <span className="text-xl font-bold text-foreground">{APP_NAME}</span>
            </Link>
          </div>

          {/* Search bar */}
          <div className="hidden md:flex flex-1 max-w-lg mx-8">
            <form onSubmit={handleSearch} className="flex w-full">
              <Input
                type="search"
                placeholder="Search products..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="rounded-r-none border-r-0"
                leftIcon={<Search className="h-4 w-4" />}
              />
              <Button 
                type="submit" 
                className="rounded-l-none"
                variant="default"
              >
                Search
              </Button>
            </form>
          </div>

          {/* Right side actions */}
          <div className="flex items-center space-x-4">
            {/* Mobile search */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => {/* TODO: Open mobile search */}}
            >
              <Search className="h-5 w-5" />
            </Button>

            {/* Wishlist */}
            <RequireAuth>
              <Link href="/wishlist">
                <Button variant="ghost" size="icon" className="relative">
                  <Heart className="h-5 w-5" />
                </Button>
              </Link>
            </RequireAuth>

            {/* Cart */}
            <Button
              variant="ghost"
              size="icon"
              className="relative"
              onClick={openCart}
            >
              <ShoppingCart className="h-5 w-5" />
              {cartItemCount > 0 && (
                <Badge 
                  variant="destructive" 
                  className="absolute -top-2 -right-2 h-5 w-5 rounded-full p-0 text-xs"
                >
                  {cartItemCount}
                </Badge>
              )}
            </Button>

            {/* User menu */}
            {isAuthenticated ? (
              <div className="relative">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                  className="relative"
                >
                  <User className="h-5 w-5" />
                </Button>

                {/* User dropdown */}
                {isUserMenuOpen && (
                  <div className="absolute right-0 mt-2 w-48 rounded-md border bg-white py-1 shadow-lg">
                    <div className="px-4 py-2 border-b">
                      <p className="text-sm font-medium text-gray-900">
                        {user?.first_name} {user?.last_name}
                      </p>
                      <p className="text-sm text-gray-500">{user?.email}</p>
                    </div>
                    
                    {USER_NAV.map((item) => (
                      <Link
                        key={item.href}
                        href={item.href}
                        className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
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
                    {user?.role === 'admin' && (
                      <Link
                        href="/admin"
                        className="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 border-t"
                        onClick={() => setIsUserMenuOpen(false)}
                      >
                        <Settings className="mr-3 h-4 w-4" />
                        Admin Panel
                      </Link>
                    )}
                    
                    <button
                      onClick={handleLogout}
                      className="flex w-full items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
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
              className="md:hidden"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
            >
              {isMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </Button>
          </div>
        </div>
      </div>

      {/* Navigation */}
      <nav className="border-t border-gray-100">
        <div className="container mx-auto px-4">
          <div className="hidden md:flex h-12 items-center space-x-8">
            {MAIN_NAV.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-sm font-medium text-gray-700 hover:text-gray-900 transition-colors"
              >
                {item.title}
              </Link>
            ))}
          </div>
        </div>
      </nav>

      {/* Mobile menu */}
      {isMenuOpen && (
        <div className="md:hidden border-t border-gray-100 bg-white">
          <div className="container mx-auto px-4 py-4">
            {/* Mobile search */}
            <form onSubmit={handleSearch} className="mb-4">
              <div className="flex">
                <Input
                  type="search"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="rounded-r-none border-r-0"
                />
                <Button type="submit" className="rounded-l-none">
                  <Search className="h-4 w-4" />
                </Button>
              </div>
            </form>

            {/* Mobile navigation */}
            <div className="space-y-2">
              {MAIN_NAV.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="block py-2 text-base font-medium text-gray-700 hover:text-gray-900"
                  onClick={() => setIsMenuOpen(false)}
                >
                  {item.title}
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
