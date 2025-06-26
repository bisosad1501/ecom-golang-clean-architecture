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
  Package,
  Truck,
  ArrowRight
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { CategoryMegaMenu } from '@/components/layout/category-mega-menu'
import { useAuthStore } from '@/store/auth'
import { useCartStore, getCartItemCount } from '@/store/cart'
import { APP_NAME, MAIN_NAV, USER_NAV } from '@/constants'
import { getVisibleNavItems } from '@/lib/permissions'
import { RequireAuth, RequireGuest } from '@/components/auth/permission-guard'
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

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/30 bg-background/95 backdrop-blur-xl supports-[backdrop-filter]:bg-background/90 shadow-sm">
      {/* Single Unified Header */}
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Left: Logo + Category Menu */}
          <div className="flex items-center space-x-6">
            {/* Logo */}
            <Link href="/" className="flex items-center space-x-2 group">
              <div className="h-10 w-10 rounded-xl bg-gradient-to-br from-primary-500 via-primary-600 to-violet-600 flex items-center justify-center shadow-lg group-hover:shadow-xl group-hover:scale-105 transition-all duration-300">
                <span className="text-primary-foreground font-bold text-xl">E</span>
              </div>
              <span className="hidden sm:block text-xl font-bold text-gradient bg-gradient-to-r from-primary-600 via-primary-500 to-violet-500 bg-clip-text text-transparent">{APP_NAME}</span>
            </Link>

            {/* Category Mega Menu */}
            <div className="hidden lg:block">
              <CategoryMegaMenu />
            </div>

            {/* Navigation Links */}
            <nav className="hidden lg:flex items-center space-x-1">
              {visibleNavItems.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600 hover:bg-primary-50 rounded-lg transition-all duration-200"
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </div>

          {/* Center: Search Bar */}
          <div className="flex-1 max-w-lg mx-6">
            <form onSubmit={handleSearch} className="flex w-full relative group">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
                <Input
                  type="search"
                  placeholder="Search products..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 pr-3 h-10 rounded-l-lg border-r-0 focus:ring-2 focus:ring-primary/30 focus:border-primary text-sm"
                />
              </div>
              <Button
                type="submit"
                className="rounded-l-none rounded-r-lg h-10 px-4 text-sm"
                variant="default"
              >
                Search
              </Button>
            </form>
          </div>

          {/* Right: Actions + User Menu */}
          <div className="flex items-center space-x-2">
            {/* Quick Actions */}
            <div className="hidden md:flex items-center space-x-1">
              <Button variant="ghost" size="sm" asChild>
                <Link href="/track-order" className="text-xs">
                  <Truck className="h-4 w-4 mr-1" />
                  Track
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/help" className="text-xs">
                  Help
                </Link>
              </Button>
            </div>

            {/* Wishlist */}
            <Button variant="ghost" size="sm" asChild className="relative">
              <Link href="/wishlist">
                <Heart className="h-5 w-5" />
              </Link>
            </Button>

            {/* Cart */}
            <Button 
              variant="ghost" 
              size="sm" 
              onClick={openCart}
              className="relative"
            >
              <ShoppingCart className="h-5 w-5" />
              {cartItemCount > 0 && (
                <Badge className="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 text-xs">
                  {cartItemCount}
                </Badge>
              )}
            </Button>

            {/* User Menu */}
            <RequireGuest>
              <div className="flex items-center space-x-2">
                <Button variant="ghost" size="sm" asChild>
                  <Link href="/auth/login">Sign In</Link>
                </Button>
                <Button size="sm" asChild>
                  <Link href="/auth/register">Sign Up</Link>
                </Button>
              </div>
            </RequireGuest>

            <RequireAuth>
              <div className="relative">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                  className="flex items-center space-x-2"
                >
                  <User className="h-5 w-5" />
                  <span className="hidden sm:block text-sm">{user?.first_name || user?.email || 'Account'}</span>
                </Button>

                {isUserMenuOpen && (
                  <div className="absolute right-0 top-full mt-2 w-48 bg-background border border-border rounded-lg shadow-lg z-50">
                    <div className="p-2">
                      <div className="px-3 py-2 text-sm text-gray-600 border-b">
                        {user?.email}
                      </div>
                      
                      {USER_NAV.filter(item => getVisibleNavItems(user?.role || null).some(nav => nav.href === item.href))
                        .map((item) => (
                        <Link
                          key={item.href}
                          href={item.href}
                          className="flex items-center px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded"
                          onClick={() => setIsUserMenuOpen(false)}
                        >
                          <span className="h-4 w-4 mr-2">•</span>
                          {item.title}
                        </Link>
                      ))}
                      
                      <button
                        onClick={handleLogout}
                        className="flex items-center w-full px-3 py-2 text-sm text-red-600 hover:bg-red-50 rounded"
                      >
                        <LogOut className="h-4 w-4 mr-2" />
                        Sign Out
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </RequireAuth>

            {/* Mobile Menu */}
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="lg:hidden"
            >
              {isMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </Button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="lg:hidden border-t border-border py-4">
            <div className="space-y-2">
              {/* Mobile Search */}
              <div className="mb-4">
                <form onSubmit={handleSearch} className="flex">
                  <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                      type="search"
                      placeholder="Search products..."
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      className="pl-10 rounded-r-none"
                    />
                  </div>
                  <Button type="submit" className="rounded-l-none">
                    Search
                  </Button>
                </form>
              </div>

              {/* Mobile Navigation */}
              {visibleNavItems.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="block px-3 py-2 text-sm font-medium text-gray-700 hover:text-primary-600 hover:bg-primary-50 rounded-lg"
                  onClick={() => setIsMenuOpen(false)}
                >
                  {item.label}
                </Link>
              ))}

              {/* Mobile Quick Actions */}
              <div className="pt-2 border-t">
                <Link href="/track-order" className="block px-3 py-2 text-sm text-gray-600">
                  Track Order
                </Link>
                <Link href="/help" className="block px-3 py-2 text-sm text-gray-600">
                  Help
                </Link>
                <Link href="/contact" className="block px-3 py-2 text-sm text-gray-600">
                  Contact
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Promotional Banner (Optional) */}
      <div className="bg-gradient-to-r from-primary-600 to-violet-600 text-white py-2">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-center text-sm font-medium">
            <Truck className="h-4 w-4 mr-2" />
            <span>Free shipping on orders over $50 • New arrivals daily</span>
          </div>
        </div>
      </div>
    </header>
  )
}
