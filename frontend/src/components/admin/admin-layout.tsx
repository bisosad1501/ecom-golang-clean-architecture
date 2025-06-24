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
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background flex">
      {/* Mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/60 backdrop-blur-sm lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Enhanced Sidebar */}
      <div className={cn(
        'fixed inset-y-0 left-0 z-50 w-72 bg-gradient-to-b from-white via-white to-muted/30 shadow-2xl transform transition-transform duration-300 ease-in-out lg:translate-x-0 lg:static lg:inset-0 lg:transform-none flex flex-col border-r border-border/50',
        sidebarOpen ? 'translate-x-0' : '-translate-x-full'
      )}>
        <div className="flex items-center justify-between h-20 px-8 border-b border-border/30 flex-shrink-0">
          <Link href="/admin" className="flex items-center space-x-3 group">
            <div className="h-12 w-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large group-hover:scale-105 transition-transform duration-200">
              <span className="text-white font-bold text-2xl">A</span>
            </div>
            <div>
              <span className="text-2xl font-bold text-foreground">Admin</span>
              <p className="text-xs text-muted-foreground">Dashboard</p>
            </div>
          </Link>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setSidebarOpen(false)}
            className="lg:hidden h-10 w-10 rounded-xl hover:bg-muted"
          >
            <X className="h-5 w-5" />
          </Button>
        </div>

        <nav className="flex-1 mt-8 px-6 overflow-y-auto">
          {/* Back to Store */}
          <Link
            href="/"
            className="flex items-center px-4 py-3 text-sm font-medium text-muted-foreground rounded-2xl hover:bg-muted hover:text-foreground mb-8 transition-all duration-200 group"
          >
            <div className="w-8 h-8 rounded-xl bg-muted flex items-center justify-center mr-3 group-hover:bg-primary group-hover:text-white transition-all duration-200">
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
                      ? 'bg-gradient-to-r from-primary-500 to-violet-600 text-white shadow-large'
                      : 'text-muted-foreground hover:bg-muted hover:text-foreground'
                  )}
                >
                  <div className={cn(
                    'w-8 h-8 rounded-xl flex items-center justify-center mr-3 transition-all duration-200',
                    isActive
                      ? 'bg-white/20 text-white'
                      : 'bg-muted group-hover:bg-primary group-hover:text-white'
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
        <div className="sticky top-0 z-30 bg-white/80 backdrop-blur-xl border-b border-border/30">
          <div className="flex items-center justify-between h-20 px-8">
            <div className="flex items-center space-x-6">
              <Button
                variant="ghost"
                size="icon"
                onClick={() => setSidebarOpen(true)}
                className="lg:hidden h-10 w-10 rounded-xl hover:bg-muted"
              >
                <Menu className="h-5 w-5" />
              </Button>

              <div>
                <h1 className="text-2xl font-bold text-foreground">
                  {pathname === '/admin' ? 'Dashboard' :
                   pathname.split('/').pop()?.replace('-', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                </h1>
                <p className="text-sm text-muted-foreground">
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
                    className="w-80 h-12 pl-12 pr-4 rounded-2xl border-2 border-border/30 focus:border-primary bg-muted/30 backdrop-blur-sm"
                    leftIcon={<Search className="h-5 w-5 text-muted-foreground" />}
                  />
                  <div className="absolute right-3 top-1/2 -translate-y-1/2">
                    <kbd className="px-2 py-1 text-xs bg-muted border border-border rounded-lg text-muted-foreground">⌘K</kbd>
                  </div>
                </div>
              </div>

              {/* Enhanced Notifications */}
              <Button variant="ghost" size="icon" className="relative h-12 w-12 rounded-2xl hover:bg-muted transition-colors">
                <Bell className="h-6 w-6" />
                <div className="absolute -top-1 -right-1 h-6 w-6 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center shadow-large">
                  <span className="text-white text-xs font-bold">3</span>
                </div>
              </Button>

              {/* Enhanced User Menu */}
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" className="flex items-center space-x-3 h-12 px-4 rounded-2xl hover:bg-muted transition-all duration-200 group">
                    <div className="w-10 h-10 bg-gradient-to-br from-primary-500 to-violet-600 rounded-2xl flex items-center justify-center shadow-medium group-hover:scale-105 transition-transform duration-200">
                      <span className="text-white text-sm font-bold">
                        {user?.first_name?.[0]}{user?.last_name?.[0]}
                      </span>
                    </div>
                    <div className="hidden lg:block text-left">
                      <div className="text-sm font-semibold text-foreground">
                        {user?.first_name} {user?.last_name}
                      </div>
                      <div className="text-xs text-muted-foreground capitalize flex items-center gap-1">
                        <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                        {user?.role.replace('_', ' ')}
                      </div>
                    </div>
                    <ChevronDown className="h-4 w-4 text-muted-foreground group-hover:text-foreground transition-colors" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-72 p-2 border-2 border-border/30 shadow-2xl rounded-2xl">
                  <div className="p-4 border-b border-border/30">
                    <div className="flex items-center gap-3">
                      <div className="w-12 h-12 bg-gradient-to-br from-primary-500 to-violet-600 rounded-2xl flex items-center justify-center">
                        <span className="text-white font-bold">
                          {user?.first_name?.[0]}{user?.last_name?.[0]}
                        </span>
                      </div>
                      <div>
                        <p className="font-semibold text-foreground">
                          {user?.first_name} {user?.last_name}
                        </p>
                        <p className="text-sm text-muted-foreground">{user?.email}</p>
                        <div className="flex items-center gap-1 mt-1">
                          <div className="w-2 h-2 bg-emerald-500 rounded-full"></div>
                          <span className="text-xs text-muted-foreground capitalize">
                            {user?.role.replace('_', ' ')}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div className="py-2">
                    <DropdownMenuItem asChild>
                      <Link href="/admin/profile" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-muted transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-blue-100 flex items-center justify-center">
                          <User className="h-4 w-4 text-blue-600" />
                        </div>
                        <div>
                          <p className="font-medium">Profile</p>
                          <p className="text-xs text-muted-foreground">Manage your account</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>

                    <DropdownMenuItem asChild>
                      <Link href="/admin/settings" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-muted transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-purple-100 flex items-center justify-center">
                          <Settings className="h-4 w-4 text-purple-600" />
                        </div>
                        <div>
                          <p className="font-medium">Settings</p>
                          <p className="text-xs text-muted-foreground">Preferences & config</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>

                    <DropdownMenuItem asChild>
                      <Link href="/" className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-muted transition-colors">
                        <div className="w-8 h-8 rounded-xl bg-emerald-100 flex items-center justify-center">
                          <Home className="h-4 w-4 text-emerald-600" />
                        </div>
                        <div>
                          <p className="font-medium">Back to Store</p>
                          <p className="text-xs text-muted-foreground">Visit your store</p>
                        </div>
                      </Link>
                    </DropdownMenuItem>
                  </div>

                  <div className="pt-2 border-t border-border/30">
                    <DropdownMenuItem
                      className="cursor-pointer flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-destructive/10 text-destructive focus:text-destructive transition-colors"
                      onClick={() => logout()}
                    >
                      <div className="w-8 h-8 rounded-xl bg-red-100 flex items-center justify-center">
                        <LogOut className="h-4 w-4 text-red-600" />
                      </div>
                      <div>
                        <p className="font-medium">Sign out</p>
                        <p className="text-xs text-muted-foreground">End your session</p>
                      </div>
                    </DropdownMenuItem>
                  </div>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>

        {/* Enhanced Page content */}
        <main className="p-8 bg-gradient-to-br from-background via-muted/10 to-background min-h-screen">
          {children}
        </main>
      </div>
    </div>
  )
}
