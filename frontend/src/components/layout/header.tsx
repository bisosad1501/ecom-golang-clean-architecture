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

  const handleRefreshUser = async () => {
    try {
      await refreshUser()
      console.log('User data refreshed')
    } catch (error) {
      console.error('Failed to refresh user:', error)
    }
  }

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/30 bg-background/95 backdrop-blur-xl supports-[backdrop-filter]:bg-background/90 shadow-medium">
      {/* Top bar */}
      <div className="border-b border-border/20 bg-gradient-to-r from-primary-50/80 to-violet-50/80 py-3">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between text-sm">
            <div className="hidden md:flex items-center space-x-3">
              <div className="flex items-center gap-2 bg-primary-100/50 rounded-full px-3 py-1">
                <Truck className="h-4 w-4 text-primary-600" />
                <span className="font-semibold text-primary-700">Free shipping on orders over $50</span>
              </div>
            </div>
            <div className="flex items-center space-x-8">
              <Link href="/help" className="hover:text-primary transition-all duration-200 font-medium hover:scale-105 flex items-center gap-1 group">
                <span>Help</span>
                <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all duration-200" />
              </Link>
              <Link href="/contact" className="hover:text-primary transition-all duration-200 font-medium hover:scale-105 flex items-center gap-1 group">
                <span>Contact</span>
                <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all duration-200" />
              </Link>
              <Link href="/track-order" className="hover:text-primary transition-all duration-200 font-medium hover:scale-105 flex items-center gap-1 group">
                <span>Track Order</span>
                <ArrowRight className="h-3 w-3 opacity-0 group-hover:opacity-100 group-hover:translate-x-1 transition-all duration-200" />
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Main header */}
      <div className="container mx-auto px-4">
        <div className="flex h-20 items-center justify-between">
          {/* Logo */}
          <div className="flex items-center">
            <Link href="/" className="flex items-center space-x-3 group">
              <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-primary-500 via-primary-600 to-violet-600 flex items-center justify-center shadow-large group-hover:shadow-xl group-hover:scale-105 transition-all duration-300">
                <span className="text-primary-foreground font-bold text-2xl">E</span>
              </div>
              <span className="text-3xl font-bold text-gradient bg-gradient-to-r from-primary-600 via-primary-500 to-violet-500 bg-clip-text text-transparent">{APP_NAME}</span>
            </Link>
          </div>

          {/* Enhanced Search bar */}
          <div className="hidden md:flex flex-1 max-w-3xl mx-12">
            <form onSubmit={handleSearch} className="flex w-full relative group">
              <div className="relative flex-1">
                <Search className="absolute left-5 top-1/2 transform -translate-y-1/2 h-5 w-5 text-muted-foreground group-focus-within:text-primary transition-colors" />
                <Input
                  type="search"
                  placeholder="Search for products, brands, categories..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-14 pr-6 h-14 rounded-l-2xl border-r-0 focus:ring-2 focus:ring-primary/30 focus:border-primary text-base shadow-medium hover:shadow-large transition-all duration-200"
                />

                {/* Search suggestions dropdown */}
                {searchQuery && (
                  <div className="absolute top-full left-0 right-0 bg-background border border-border rounded-b-2xl shadow-2xl z-50 max-h-96 overflow-y-auto">
                    <div className="p-4">
                      <div className="text-sm text-muted-foreground mb-3 font-medium">Popular searches</div>
                      <div className="space-y-2">
                        {['Electronics', 'Fashion', 'Home & Garden', 'Sports'].map((suggestion) => (
                          <div key={suggestion} className="flex items-center gap-3 p-2 hover:bg-muted rounded-lg cursor-pointer transition-colors">
                            <Search className="h-4 w-4 text-muted-foreground" />
                            <span className="text-sm">{suggestion}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
              <Button
                type="submit"
                className="rounded-l-none rounded-r-2xl h-14 px-8 text-base font-semibold shadow-large hover:shadow-xl"
                variant="gradient"
              >
                Search
              </Button>
            </form>
          </div>

          {/* Right side actions */}
          <div className="flex items-center space-x-3">
            {/* Mobile search */}
            <Button
              variant="ghost"
              size="icon"
              className="md:hidden h-12 w-12 rounded-2xl hover:bg-primary/10 hover:scale-105 transition-all duration-200"
              onClick={() => {/* TODO: Open mobile search */}}
            >
              <Search className="h-5 w-5" />
            </Button>

            {/* Wishlist */}
            <RequireAuth>
              <Link href="/wishlist">
                <Button variant="ghost" size="icon" className="relative h-12 w-12 rounded-2xl hover:bg-primary/10 hover:scale-105 transition-all duration-200 group">
                  <Heart className="h-5 w-5 group-hover:text-red-500 transition-colors" />
                </Button>
              </Link>
            </RequireAuth>

            {/* Enhanced Cart with preview */}
            <div className="relative group">
              <Button
                variant="ghost"
                size="icon"
                className="relative h-12 w-12 rounded-2xl hover:bg-primary/10 hover:scale-105 transition-all duration-200"
                onClick={openCart}
              >
                <ShoppingCart className="h-5 w-5 group-hover:text-primary transition-colors" />
                {cartItemCount > 0 && (
                  <Badge
                    variant="default"
                    className="absolute -top-1 -right-1 h-6 w-6 rounded-full p-0 text-xs font-bold shadow-large animate-pulse"
                  >
                    {cartItemCount}
                  </Badge>
                )}
              </Button>

              {/* Cart preview on hover */}
              {cartItemCount > 0 && cart && (
                <div className="absolute top-full right-0 mt-2 w-80 bg-background border border-border rounded-2xl shadow-2xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-50">
                  <div className="p-6">
                    <div className="flex items-center justify-between mb-4">
                      <h3 className="font-semibold text-lg">Shopping Cart</h3>
                      <Badge variant="secondary">{cartItemCount} items</Badge>
                    </div>

                    <div className="space-y-3 max-h-64 overflow-y-auto">
                      {cart.items.slice(0, 3).map((item) => (
                        <div key={item.id} className="flex items-center gap-3 p-2 hover:bg-muted rounded-lg transition-colors">
                          <div className="w-12 h-12 bg-muted rounded-lg"></div>
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
                      <Button className="w-full" variant="gradient" onClick={openCart}>
                        View Cart
                      </Button>
                    </div>
                  </div>
                </div>
              )}
            </div>

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
