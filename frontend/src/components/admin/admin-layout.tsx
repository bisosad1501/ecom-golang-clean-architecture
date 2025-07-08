'use client'

import { useState } from 'react'
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
  ChevronDown
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { 
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useAuthStore } from '@/store/auth'
import { getAdminSidebarItems } from '@/lib/permissions'
import { APP_NAME } from '@/constants'
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
}

export function AdminLayout({ children }: { children: React.ReactNode }) {
  console.log('=== AdminLayout RENDERING ===')
  
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const pathname = usePathname()
  const { user, logout } = useAuthStore()

  console.log('AdminLayout - user:', user?.role, 'pathname:', pathname)

  const sidebarItems = user ? getAdminSidebarItems(user.role) : []

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-slate-950 to-black flex">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Enhanced Sidebar */}
      <div className={cn(
        'fixed inset-y-0 left-0 z-50 w-72 bg-gradient-to-b from-gray-900 via-gray-900 to-gray-800 shadow-2xl transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 lg:transform-none flex flex-col border-r border-gray-700/50',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      )}>
        <div className="flex items-center justify-between h-20 px-8 border-b border-gray-700/30 flex-shrink-0">
          <Link href="/admin" className="flex items-center space-x-3 group">
            <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center shadow-large group-hover:scale-105 transition-transform duration-200">
              <span className="text-white font-bold text-2xl">A</span>
            </div>
            <div>
              <span className="text-2xl font-bold text-white">Admin</span>
              <p className="text-xs text-gray-400">Dashboard</p>
            </div>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setSidebarOpen(false)}
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
          <div className="space-y-3">
            {sidebarItems.map((item) => {
              const Icon = iconMap[item.icon as keyof typeof iconMap]

              // Fix active logic to avoid conflicts
              let isActive = false
              if (item.href === '/admin') {
                // Dashboard should only be active on exact /admin path
                isActive = pathname === '/admin'
              } else {
                // Other items should be active when pathname starts with their href
                isActive = pathname === item.href || pathname.startsWith(`${item.href}/`)
              }

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    'flex items-center px-4 py-3 text-sm font-semibold rounded-2xl transition-all duration-200 group',
                    isActive
                      ? 'bg-gradient-to-r from-[#FF9000] to-[#e67e00] text-white shadow-large'
                      : 'text-gray-400 hover:bg-gray-800 hover:text-white'
                  )}
                >
                  <div className={cn(
                    'w-8 h-8 rounded-xl flex items-center justify-center mr-3 transition-all duration-200',
                    isActive
                      ? 'bg-white/20 text-white'
                      : 'bg-gray-800 group-hover:bg-[#FF9000] group-hover:text-white'
                  )}>
                    <Icon className="h-4 w-4" />
                  </div>
                  {item.label}
                  {isActive && (
                    <div className="ml-auto w-2 h-2 bg-white rounded-full"></div>
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
        <div className="sticky top-0 z-30 bg-gray-900/90 backdrop-blur-xl border-b border-gray-700/30">
          <div className="flex items-center justify-between h-20 px-8">
            <div className="flex items-center space-x-6">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden h-10 w-10 rounded-xl hover:bg-gray-800 text-gray-400 hover:text-white"
              >
                <Menu className="h-5 w-5" />
              </Button>

              <div>
                <h1 className="text-2xl font-bold text-white">
                  {pathname === '/admin' ? 'Dashboard' :
                   pathname.split('/').pop()?.replace('-', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                </h1>
                <p className="text-sm text-gray-400">
                  {pathname === '/admin' ? 'Overview of your store performance' : 'Manage your store'}
                </p>
              </div>
            </div>

            <div className="flex items-center space-x-6">
              {/* Enhanced Search */}
              <div className="hidden md:block">
                <div className="relative">
                  <Input
                    type="search"
                    placeholder="Search products, orders, customers..."
                    className="w-80 h-12 pl-12 pr-4 rounded-2xl border-2 border-gray-600/30 focus:border-[#FF9000] bg-gray-800/30 backdrop-blur-sm text-white placeholder:text-gray-400"
                    leftIcon={<Search className="h-5 w-5 text-gray-400" />}
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    <kbd className="px-2 py-1 text-xs bg-gray-800 border border-gray-600 rounded-lg text-gray-400">⌘K</kbd>
                  </div>
                </div>
              </div>

              {/* Enhanced Notifications */}
              <Button variant="ghost" size="icon" className="relative h-12 w-12 rounded-2xl hover:bg-gray-800 transition-colors text-gray-400 hover:text-white">
                <Bell className="h-6 w-6" />
                <div className="absolute -top-1 -right-1 h-6 w-6 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-full flex items-center justify-center shadow-large">
                  <span className="text-white text-xs font-bold">3</span>
                </div>
              </Button>

              {/* Enhanced User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="flex items-center space-x-3 h-12 px-4 rounded-2xl hover:bg-gray-800 transition-all duration-200 group">
                    <div className="w-10 h-10 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-2xl flex items-center justify-center shadow-medium group-hover:scale-105 transition-transform duration-200">
                      <span className="text-white text-sm font-bold">
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
                    <ChevronDown className="h-4 w-4 text-gray-400 group-hover:text-white transition-colors" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-72 p-2 border-2 border-gray-600/30 shadow-2xl rounded-2xl bg-gray-900">
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

        {/* Enhanced Page content */}
        <main className="p-8 bg-gradient-to-br from-slate-950 via-slate-900/10 to-slate-950 min-h-screen">
          {children}
        </main>
      </div>
    </div>
  )
}
