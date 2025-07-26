'use client'

import { useState, useMemo, useCallback } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  Menu,
  X,
  BarChart3,
  Package,
  Folder,
  ShoppingCart,
  Users,
  Star,
  Ticket,
  Settings,
  Home,
  Bell,
  Search,
  LogOut,
  User,
  ChevronDown,
  TrendingUp,
  UserCheck,
  FileText,
  FileImage,
  Shield,
  Server,
  Tag,
  Truck,
  Warehouse,
  Database,
  PieChart,
  CreditCard,
  Layers,
  Archive,
  Lock,
  Cog
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useAuthStore } from '@/store/auth'
import { getAdminSidebarItems } from '@/lib/permissions'

import { cn } from '@/lib/utils'

const iconMap = {
  BarChart3,
  Package,
  Folder,
  ShoppingCart,
  Users,
  Star,
  Ticket,
  Settings,
  TrendingUp,
  UserCheck,
  FileText,
  FileImage,
  Shield,
  Server,
  Bell,
  Tag,
  Truck,
  Warehouse,
  Database,
  PieChart,
  CreditCard,
  Layers,
  Archive,
  Lock,
  Cog,
}

// Color mapping for sidebar items - harmonious with BiHub orange (#FF9000)
const getItemColors = (color: string) => {
  const colorMap = {
    blue: {
      bg: 'bg-blue-500/20',
      text: 'text-blue-400',
      hover: 'group-hover:bg-blue-500',
      border: 'border-blue-400/30'
    },
    green: {
      bg: 'bg-green-500/20',
      text: 'text-green-400',
      hover: 'group-hover:bg-green-500',
      border: 'border-green-400/30'
    },
    purple: {
      bg: 'bg-purple-500/20',
      text: 'text-purple-400',
      hover: 'group-hover:bg-purple-500',
      border: 'border-purple-400/30'
    },
    orange: {
      bg: 'bg-orange-500/20',
      text: 'text-orange-400',
      hover: 'group-hover:bg-orange-500',
      border: 'border-orange-400/30'
    },
    teal: {
      bg: 'bg-teal-500/20',
      text: 'text-teal-400',
      hover: 'group-hover:bg-teal-500',
      border: 'border-teal-400/30'
    },
    indigo: {
      bg: 'bg-indigo-500/20',
      text: 'text-indigo-400',
      hover: 'group-hover:bg-indigo-500',
      border: 'border-indigo-400/30'
    },
    cyan: {
      bg: 'bg-cyan-500/20',
      text: 'text-cyan-400',
      hover: 'group-hover:bg-cyan-500',
      border: 'border-cyan-400/30'
    },
    pink: {
      bg: 'bg-pink-500/20',
      text: 'text-pink-400',
      hover: 'group-hover:bg-pink-500',
      border: 'border-pink-400/30'
    },
    violet: {
      bg: 'bg-violet-500/20',
      text: 'text-violet-400',
      hover: 'group-hover:bg-violet-500',
      border: 'border-violet-400/30'
    },
    emerald: {
      bg: 'bg-emerald-500/20',
      text: 'text-emerald-400',
      hover: 'group-hover:bg-emerald-500',
      border: 'border-emerald-400/30'
    },
    amber: {
      bg: 'bg-amber-500/20',
      text: 'text-amber-400',
      hover: 'group-hover:bg-amber-500',
      border: 'border-amber-400/30'
    },
    lime: {
      bg: 'bg-lime-500/20',
      text: 'text-lime-400',
      hover: 'group-hover:bg-lime-500',
      border: 'border-lime-400/30'
    },
    red: {
      bg: 'bg-red-500/20',
      text: 'text-red-400',
      hover: 'group-hover:bg-red-500',
      border: 'border-red-400/30'
    },
    slate: {
      bg: 'bg-slate-500/20',
      text: 'text-slate-400',
      hover: 'group-hover:bg-slate-500',
      border: 'border-slate-400/30'
    },
    rose: {
      bg: 'bg-rose-500/20',
      text: 'text-rose-400',
      hover: 'group-hover:bg-rose-500',
      border: 'border-rose-400/30'
    },
    stone: {
      bg: 'bg-stone-500/20',
      text: 'text-stone-400',
      hover: 'group-hover:bg-stone-500',
      border: 'border-stone-400/30'
    },
    yellow: {
      bg: 'bg-yellow-500/20',
      text: 'text-yellow-400',
      hover: 'group-hover:bg-yellow-500',
      border: 'border-yellow-400/30'
    },
    gray: {
      bg: 'bg-gray-500/20',
      text: 'text-gray-400',
      hover: 'group-hover:bg-gray-500',
      border: 'border-gray-400/30'
    }
  }

  return colorMap[color as keyof typeof colorMap] || colorMap.gray
}

export function AdminLayout({ children }: { children: React.ReactNode }) {
  console.log('=== AdminLayout RENDERING ===')
  
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const pathname = usePathname()
  const { user, logout } = useAuthStore()

  console.log('AdminLayout - user:', user?.role, 'pathname:', pathname)

  // Memoize sidebar items to prevent unnecessary recalculations
  const sidebarItems = useMemo(() => {
    return user ? getAdminSidebarItems(user.role) : getAdminSidebarItems('admin')
  }, [user?.role])

  // Memoize sidebar close handler to prevent unnecessary re-renders
  const handleSidebarClose = useCallback(() => {
    setSidebarOpen(false)
  }, [])

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-950 to-black flex">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm lg:hidden"
          onClick={handleSidebarClose}
        />
      )}

      {/* Enhanced Sidebar */}
      <div className={cn(
        'fixed inset-y-0 left-0 z-50 w-72 bg-gradient-to-b from-gray-900 via-gray-900 to-gray-800 shadow-2xl transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 lg:transform-none flex flex-col border-r border-gray-700/50',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      )}>
        <div className="flex items-center justify-between h-20 px-8 border-b border-gray-700/30 flex-shrink-0">
          <Link href="/admin" className="flex items-center space-x-3 group">
            <div className="flex items-center">
              <span className="text-2xl font-bold flex items-center">
                <span className="text-white">Bi</span>
                <span className="ml-0.5 px-0.5 py-0 rounded-[3px] text-black font-bold" style={{letterSpacing: '0.3px', backgroundColor: '#FF9000', lineHeight: '1.1'}}>hub</span>
              </span>
            </div>
            <div>
              <span className="text-lg font-bold text-white">Admin</span>
              <p className="text-xs text-gray-400">Dashboard</p>
            </div>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={handleSidebarClose}
            className="lg:hidden h-10 w-10 rounded-xl hover:bg-gray-800 text-gray-400 hover:text-white"
          >
            <X className="h-5 w-5" />
          </Button>
        </div>

        <nav className="flex-1 mt-8 px-6 overflow-y-auto">
          {/* Back to Store */}
          <Link
            href="/"
            className="flex items-center px-4 py-3 text-sm font-medium text-gray-400 rounded-2xl hover:bg-gray-800 hover:text-white mb-8 transition-all duration-200 group"
          >
            <div className="w-8 h-8 rounded-xl bg-gray-800 flex items-center justify-center mr-3 group-hover:bg-[#FF9000] group-hover:text-white transition-all duration-200">
              <Home className="h-4 w-4" />
            </div>
            Back to Store
          </Link>

          {/* Navigation Items */}
          <div className="space-y-2">
            {sidebarItems.map((item) => {
              const Icon = iconMap[item.icon as keyof typeof iconMap] || BarChart3
              const colors = getItemColors(item.color || 'gray')

              // Optimized active logic to avoid conflicts
              const isActive = useMemo(() => {
                if (item.href === '/admin') {
                  // Dashboard should only be active on exact /admin path
                  return pathname === '/admin'
                } else {
                  // Other items should be active when pathname starts with their href
                  return pathname === item.href || pathname.startsWith(`${item.href}/`)
                }
              }, [pathname, item.href])

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex items-center px-4 py-3 text-sm font-semibold rounded-xl transition-all duration-150 group will-change-transform',
                    isActive
                      ? 'bg-gradient-to-r from-[#FF9000] to-[#e67e00] text-white shadow-lg transform scale-[1.02]'
                      : 'text-gray-400 hover:bg-gray-800/70 hover:text-white hover:transform hover:scale-[1.01]'
                  )}
                >
                  <div className={cn(
                    'w-8 h-8 rounded-lg flex items-center justify-center mr-3 transition-all duration-150',
                    isActive
                      ? 'bg-white/20 text-white'
                      : `${colors.bg} ${colors.text} ${colors.hover} group-hover:text-white`
                  )}>
                    <Icon className="h-4 w-4" />
                  </div>
                  <span className="truncate">{item.label}</span>
                  {isActive && (
                    <div className="ml-auto w-2 h-2 bg-white rounded-full animate-pulse"></div>
                  )}
                </Link>
              )
            })}
          </div>
        </nav>
      </div>

      {/* Main content */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Enhanced Top bar */}
        <div className="sticky top-0 z-30 bg-gray-900/95 backdrop-blur-xl border-b border-gray-700/30 will-change-transform">
          <div className="flex items-center justify-between h-20 px-8">
            <div className="flex items-center space-x-6">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden h-10 w-10 rounded-xl hover:bg-gray-800 text-gray-400 hover:text-white transition-colors duration-150"
              >
                <Menu className="h-5 w-5" />
              </Button>

              {/* Memoized page title to prevent unnecessary recalculations */}
              {useMemo(() => (
                <div className="min-w-0 flex-1">
                  <h1 className="text-xl lg:text-2xl font-bold text-white truncate">
                    {pathname === '/admin' ? 'Dashboard' :
                     pathname.split('/').pop()?.replace('-', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                  </h1>
                  <p className="text-xs lg:text-sm text-gray-400 hidden sm:block">
                    {pathname === '/admin' ? 'Overview of your store performance' : 'Manage your store'}
                  </p>
                </div>
              ), [pathname])}
            </div>

            <div className="flex items-center space-x-3 lg:space-x-6">
              {/* Mobile Search Button */}
              <Button
                variant="ghost"
                size="icon"
                className="md:hidden h-10 w-10 rounded-xl hover:bg-gray-800 text-gray-400 hover:text-white"
              >
                <Search className="h-5 w-5" />
              </Button>

              {/* Enhanced Search - Desktop */}
              <div className="hidden md:block">
                <div className="relative">
                  <Input
                    type="search"
                    placeholder="Search products, orders, customers..."
                    className="w-64 lg:w-80 h-10 lg:h-12 pl-10 lg:pl-12 pr-4 rounded-xl lg:rounded-2xl border-2 border-gray-600/30 focus:border-[#FF9000] bg-gray-800/30 backdrop-blur-sm text-white placeholder:text-gray-400 text-sm lg:text-base"
                    leftIcon={<Search className="h-4 w-4 lg:h-5 lg:w-5 text-gray-400" />}
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2 hidden lg:block">
                    <kbd className="px-2 py-1 text-xs bg-gray-800 border border-gray-600 rounded-lg text-gray-400">⌘K</kbd>
                  </div>
                </div>
              </div>

              {/* Enhanced Notifications */}
              <Button variant="ghost" size="icon" className="relative h-10 w-10 lg:h-12 lg:w-12 rounded-xl lg:rounded-2xl hover:bg-gray-800 transition-colors text-gray-400 hover:text-white">
                <Bell className="h-5 w-5 lg:h-6 lg:w-6" />
                <div className="absolute -top-1 -right-1 h-5 w-5 lg:h-6 lg:w-6 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-full flex items-center justify-center shadow-large">
                  <span className="text-white text-xs font-bold">3</span>
                </div>
              </Button>

              {/* Enhanced User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="flex items-center space-x-2 lg:space-x-3 h-10 lg:h-12 px-2 lg:px-4 rounded-xl lg:rounded-2xl hover:bg-gray-800 transition-all duration-200 group">
                    <div className="w-8 h-8 lg:w-10 lg:h-10 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-xl lg:rounded-2xl flex items-center justify-center shadow-medium group-hover:scale-105 transition-transform duration-200">
                      <span className="text-white text-xs lg:text-sm font-bold">
                        {user?.first_name?.[0]}{user?.last_name?.[0]}
                      </span>
                    </div>
                    <div className="hidden lg:block text-left">
                      <div className="text-sm font-semibold text-white">
                        {user?.first_name} {user?.last_name}
                      </div>
                      <div className="text-xs text-gray-400 capitalize flex items-center gap-1">
                        <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                        {user?.role.replace('_', ' ')}
                      </div>
                    </div>
                    <ChevronDown className="h-3 w-3 lg:h-4 lg:w-4 text-gray-400 group-hover:text-white transition-colors hidden sm:block" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-64 lg:w-72 p-2 border-2 border-gray-600/30 shadow-2xl rounded-xl lg:rounded-2xl bg-gray-900">
                  <div className="p-4 border-b border-gray-700/30">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-2xl flex items-center justify-center">
                        <span className="text-white font-bold">
                          {user?.first_name?.[0]}{user?.last_name?.[0]}
                        </span>
                      </div>
                      <div>
                        <p className="font-semibold text-white">
                          {user?.first_name} {user?.last_name}
                        </p>
                        <p className="text-sm text-gray-400">{user?.email}</p>
                        <div className="flex items-center gap-1 mt-1">
                          <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                          <span className="text-xs text-gray-400 capitalize">
                            {user?.role.replace('_', ' ')}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="py-2">
                    <DropdownMenuItem asChild>
                      <Link href="/admin/profile" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-800 transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-blue-900/30 flex items-center justify-center">
                          <User className="h-4 w-4 text-blue-400" />
                        </div>
                        <div>
                          <p className="font-medium text-white">Profile</p>
                          <p className="text-xs text-gray-400">Manage your account</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>

                    <DropdownMenuItem asChild>
                      <Link href="/admin/settings" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-800 transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-purple-900/30 flex items-center justify-center">
                          <Settings className="h-4 w-4 text-purple-400" />
                        </div>
                        <div>
                          <p className="font-medium text-white">Settings</p>
                          <p className="text-xs text-gray-400">Preferences & config</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>

                    <DropdownMenuItem asChild>
                      <Link href="/" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-800 transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-emerald-900/30 flex items-center justify-center">
                          <Home className="h-4 w-4 text-emerald-400" />
                        </div>
                        <div>
                          <p className="font-medium text-white">Back to Store</p>
                          <p className="text-xs text-gray-400">Visit your store</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>
                  </div>

                  <div className="pt-2 border-t border-gray-700/30">
                    <DropdownMenuItem
                      className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-red-900/20 text-red-400 focus:text-red-400 transition-colors"
                      onClick={() => logout()}
                    >
                      <div className="w-8 h-8 rounded-xl bg-red-900/30 flex items-center justify-center">
                        <LogOut className="h-4 w-4 text-red-400" />
                      </div>
                      <div>
                        <p className="font-medium">Sign out</p>
                        <p className="text-xs text-gray-400">End your session</p>
                      </div>
                    </DropdownMenuItem>
                  </div>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>

        {/* Enhanced Page content with optimized rendering */}
        <main className="p-8 bg-gradient-to-br from-slate-950 via-slate-900/10 to-slate-950 min-h-screen will-change-transform">
          <div className="transform-gpu">
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
