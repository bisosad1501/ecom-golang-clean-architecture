'use client'

import { useState } from 'react'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Search,
  Filter,
  MoreHorizontal,
  Eye,
  Edit,
  Users,
  UserCheck,
  UserX,
  Mail,
  Phone,
  Calendar,
  Shield,
  Crown,
  User,
  Download,
  UserPlus,
  Grid,
  List,
} from 'lucide-react'
import { User as UserType } from '@/types'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { formatDate } from '@/lib/utils'
import { useUsers } from '@/hooks/use-users'
import { toast } from 'sonner'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export default function AdminUsersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [currentPage, setCurrentPage] = useState(1)
  const [roleFilter, setRoleFilter] = useState<string>('')
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')

  const {
    data: usersData,
    isLoading,
    error
  } = useUsers({
    limit: 20,
    page: currentPage,
    search: searchQuery,
    role: roleFilter,
    is_active: statusFilter ? statusFilter === 'active' : undefined
  })

  const users = usersData?.data || []
  const pagination = usersData?.pagination

  const getRoleIcon = (role: string) => {
    switch (role.toLowerCase()) {
      case 'admin':
        return <Crown className="h-4 w-4 text-white" />
      case 'moderator':
        return <Shield className="h-4 w-4 text-white" />
      case 'customer':
        return <User className="h-4 w-4 text-white" />
      default:
        return <User className="h-4 w-4 text-white" />
    }
  }

  const getRoleColor = (role: string) => {
    switch (role.toLowerCase()) {
      case 'admin':
        return 'from-red-500 to-red-600'
      case 'moderator':
        return 'from-blue-500 to-blue-600'
      case 'customer':
        return 'from-green-500 to-green-600'
      default:
        return 'from-gray-500 to-gray-600'
    }
  }

  const handleExportUsers = () => {
    console.log('Export users')
  }

  const handleViewUser = (user: UserType) => {
    console.log('Viewing user:', user)
    // TODO: Implement user view modal or navigate to user detail page
    toast.info(`Viewing user: ${user.first_name} ${user.last_name}`, {
      description: `Email: ${user.email} | Role: ${user.role}`,
      duration: 3000,
    })
  }

  const handleEditUser = (user: UserType) => {
    console.log('Editing user:', user)
    // TODO: Implement user edit modal or navigate to user edit page
    toast.info(`Edit User Feature`, {
      description: `Edit functionality for ${user.first_name} ${user.last_name} will be available soon!`,
      duration: 3000,
    })
  }

  const handleDeleteUser = (user: UserType) => {
    console.log('Deleting user:', user)
    // TODO: Implement user deletion
    toast.warning(`Delete User`, {
      description: `Delete functionality for ${user.first_name} ${user.last_name} will be available soon!`,
      duration: 3000,
    })
  }

  const handleChangeUserRole = (user: UserType, newRole: string) => {
    console.log('Changing user role:', user, newRole)
    // TODO: Implement role change
    toast.info(`Role Change`, {
      description: `Role change for ${user.first_name} ${user.last_name} to ${newRole} will be available soon!`,
      duration: 3000,
    })
  }

  const handleToggleUserStatus = (user: UserType) => {
    console.log('Toggling user status:', user)
    const newStatus = !user.is_active
    // TODO: Implement status toggle
    toast.info(`Status Toggle`, {
      description: `Status change for ${user.first_name} ${user.last_name} to ${newStatus ? 'Active' : 'Inactive'} will be available soon!`,
      duration: 3000,
    })
  }

  const getTotalStats = () => {
    // Use pagination.total for total users, with fallback to users array length
    const totalUsers = pagination?.total || users?.length || 0
    const usersArray = Array.isArray(users) ? users : []
    
    // Debug logging
    console.log('=== USER STATS DEBUG ===')
    console.log('pagination:', pagination)
    console.log('users:', users)
    console.log('users.length:', users?.length)
    console.log('pagination.total:', pagination?.total)
    console.log('calculated totalUsers:', totalUsers)
    
    // For stats based on current page data (since we don't have full dataset)
    // In a real app, you'd want separate API calls for these stats
    const activeUsers = usersArray.filter(user => user.is_active).length
    const adminUsers = usersArray.filter(user => user.role === 'admin').length
    const customerUsers = usersArray.filter(user => user.role === 'customer').length
    
    return { 
      totalUsers, 
      activeUsers, 
      adminUsers, 
      customerUsers 
    }
  }

  const stats = getTotalStats()

  if (error) {
    return (
      <BiHubEmptyState
        icon={<UserX className="h-8 w-8 text-red-400" />}
        title="Error loading users"
        description={error.message}
      />
    )
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="User Management"
        subtitle="Manage BiHub user accounts, roles, and permissions"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Users' }
        ]}
        action={
          <div className="flex items-center gap-4">
            <Button
              onClick={handleExportUsers}
              className={BIHUB_ADMIN_THEME.components.button.secondary}
            >
              <Download className="mr-2 h-4 w-4" />
              Export Users
            </Button>

            <RequirePermission permission={PERMISSIONS.USERS_MANAGE_ROLES}>
              <Button className={BIHUB_ADMIN_THEME.components.button.primary}>
                <UserPlus className="mr-2 h-4 w-4" />
                Add User
              </Button>
            </RequirePermission>
          </div>
        }
      />

      {/* Enhanced Quick Stats with Modern Design */}
      <div className="space-y-6">
        {/* Primary Stats Row - Modern Glass Design */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          {/* Total Users */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-blue-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-blue-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-indigo-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-blue-400/80 text-sm font-medium uppercase tracking-wide">Total Users</p>
                <p className="text-2xl font-bold text-blue-100 mt-1">{stats.totalUsers}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Users className="h-5 w-5 text-blue-400" />
              </div>
            </div>
          </div>

          {/* Active Users */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-emerald-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-emerald-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-emerald-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/5 to-green-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-emerald-400/80 text-sm font-medium uppercase tracking-wide">Active Users</p>
                <p className="text-2xl font-bold text-emerald-100 mt-1">{stats.activeUsers}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <UserCheck className="h-5 w-5 text-emerald-400" />
              </div>
            </div>
          </div>

          {/* Admins */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-purple-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-purple-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-purple-500/5 to-violet-600/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-purple-400/80 text-sm font-medium uppercase tracking-wide">Admins</p>
                <p className="text-2xl font-bold text-purple-100 mt-1">{stats.adminUsers}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-purple-500/20 to-violet-600/20 border border-purple-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <Crown className="h-5 w-5 text-purple-400" />
              </div>
            </div>
          </div>

          {/* Customers */}
          <div className="group relative bg-white/5 backdrop-blur-sm border border-amber-300/20 rounded-xl p-4 hover:bg-white/10 hover:border-amber-400/40 hover:scale-[1.02] transition-all duration-200 shadow-sm hover:shadow-lg hover:shadow-amber-500/10">
            <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 to-orange-500/5 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-200"></div>
            <div className="relative flex items-center justify-between">
              <div>
                <p className="text-amber-400/80 text-sm font-medium uppercase tracking-wide">Customers</p>
                <p className="text-2xl font-bold text-amber-100 mt-1">{stats.customerUsers}</p>
              </div>
              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-amber-500/20 to-orange-500/20 border border-amber-400/30 flex items-center justify-center group-hover:scale-110 transition-transform">
                <User className="h-5 w-5 text-amber-400" />
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Search & Filters */}
      <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-6 shadow-lg">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-purple-500/5 rounded-2xl"></div>
        <div className="relative">
          {/* Header */}
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-purple-500/20 border border-blue-400/30 flex items-center justify-center">
                <Search className="h-5 w-5 text-blue-400" />
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">Search & Filter Users</h3>
                <p className="text-sm text-gray-400">Find and filter BiHub users by role and status</p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className="group relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white transition-all duration-200"
            >
              <div className="flex items-center gap-2">
                {viewMode === 'grid' ? (
                  <List className="h-4 w-4" />
                ) : (
                  <Grid className="h-4 w-4" />
                )}
                <span className="text-sm font-medium">
                  {viewMode === 'grid' ? 'List View' : 'Grid View'}
                </span>
              </div>
            </Button>
          </div>

          {/* Search and Filter Controls */}
          <div className="flex flex-col lg:flex-row items-center gap-4">
            <div className="flex-1 w-full">
              <div className="relative group">
                <Search className="absolute left-4 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400 group-focus-within:text-blue-400 transition-colors" />
                <Input
                  placeholder="Search users by name, email, or role..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-full h-12 pl-12 pr-12 bg-white/5 border-gray-600/50 rounded-xl text-white placeholder:text-gray-400 focus:bg-white/10 focus:border-blue-500/50 focus:ring-2 focus:ring-blue-500/20 transition-all duration-200"
                />
                {searchQuery && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setSearchQuery('')}
                    className="absolute right-3 top-1/2 -translate-y-1/2 h-6 w-6 p-0 text-gray-400 hover:text-white hover:bg-white/10 rounded-md transition-all duration-200"
                  >
                    ×
                  </Button>
                )}
              </div>
            </div>

            <div className="flex items-center gap-3">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button className="group relative bg-white/5 border border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white px-4 py-2 h-12 rounded-xl transition-all duration-200">
                    <Filter className="mr-2 h-4 w-4" />
                    <span className="font-medium">Role: {roleFilter || 'All'}</span>
                    <div className="ml-2 w-2 h-2 rounded-full bg-blue-400 animate-pulse opacity-60"></div>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                  <DropdownMenuItem
                    onClick={() => setRoleFilter('')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gray-600/50 flex items-center justify-center">
                        <Filter className="h-3 w-3 text-white" />
                      </div>
                      <span className="font-medium">All Roles</span>
                    </div>
                  </DropdownMenuItem>
                  <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                  <DropdownMenuItem
                    onClick={() => setRoleFilter('admin')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-purple-500/20 to-violet-600/20 border border-purple-400/30 flex items-center justify-center">
                        <Crown className="h-3 w-3 text-purple-400" />
                      </div>
                      <span className="font-medium">Admin</span>
                    </div>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setRoleFilter('moderator')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-blue-500/20 to-indigo-600/20 border border-blue-400/30 flex items-center justify-center">
                        <Shield className="h-3 w-3 text-blue-400" />
                      </div>
                      <span className="font-medium">Moderator</span>
                    </div>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setRoleFilter('customer')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center">
                        <User className="h-3 w-3 text-emerald-400" />
                      </div>
                      <span className="font-medium">Customer</span>
                    </div>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>

              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button className="group relative bg-white/5 border border-gray-600/50 text-gray-300 hover:bg-white/10 hover:border-gray-500 hover:text-white px-4 py-2 h-12 rounded-xl transition-all duration-200">
                    <Filter className="mr-2 h-4 w-4" />
                    <span className="font-medium">Status: {statusFilter || 'All'}</span>
                    <div className="ml-2 w-2 h-2 rounded-full bg-green-400 animate-pulse opacity-60"></div>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                  <DropdownMenuItem
                    onClick={() => setStatusFilter('')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gray-600/50 flex items-center justify-center">
                        <Filter className="h-3 w-3 text-white" />
                      </div>
                      <span className="font-medium">All Status</span>
                    </div>
                  </DropdownMenuItem>
                  <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                  <DropdownMenuItem
                    onClick={() => setStatusFilter('active')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-emerald-500/20 to-green-600/20 border border-emerald-400/30 flex items-center justify-center">
                        <UserCheck className="h-3 w-3 text-emerald-400" />
                      </div>
                      <span className="font-medium">Active</span>
                    </div>
                  </DropdownMenuItem>
                  <DropdownMenuItem
                    onClick={() => setStatusFilter('inactive')}
                    className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                  >
                    <div className="flex items-center gap-3">
                      <div className="w-6 h-6 rounded-lg bg-gradient-to-br from-red-500/20 to-rose-600/20 border border-red-400/30 flex items-center justify-center">
                        <UserX className="h-3 w-3 text-red-400" />
                      </div>
                      <span className="font-medium">Inactive</span>
                    </div>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          </div>
        </div>
      </div>

      {/* Modern Users List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 animate-pulse">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 bg-gray-700/50 rounded-xl"></div>
                  <div>
                    <div className="h-5 bg-gray-700/50 rounded w-24 mb-2"></div>
                    <div className="h-4 bg-gray-700/50 rounded w-16"></div>
                  </div>
                </div>
                {viewMode === 'grid' && (
                  <div className="h-6 bg-gray-700/50 rounded w-20"></div>
                )}
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-700/50 rounded w-3/4"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/2"></div>
                <div className="h-4 bg-gray-700/50 rounded w-1/4"></div>
              </div>
              <div className="flex items-center gap-3 mt-6 pt-4 border-t border-gray-700/50">
                <div className="h-8 bg-gray-700/50 rounded w-20"></div>
                <div className="h-8 bg-gray-700/50 rounded w-8"></div>
              </div>
            </div>
          ))}
        </div>
      ) : Array.isArray(users) && users.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4'
            : 'space-y-3'
        )}>
          {Array.isArray(users) && users.map((user) => (
            <div
              key={user.id}
              className={cn(
                'group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-5 hover:bg-white/10 hover:border-gray-600/50 hover:scale-[1.02] transition-all duration-200 shadow-lg hover:shadow-xl',
                viewMode === 'list' && 'flex items-center gap-6'
              )}
            >
              {/* Gradient Background */}
              <div className={cn(
                'absolute inset-0 bg-gradient-to-br opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-2xl',
                user.role === 'admin' ? 'from-purple-500/5 to-violet-600/5' :
                user.role === 'moderator' ? 'from-blue-500/5 to-indigo-600/5' :
                'from-emerald-500/5 to-green-600/5'
              )} />

              {/* User Header */}
              <div className={cn(
                'relative flex items-center justify-between mb-4',
                viewMode === 'list' && 'flex-shrink-0 w-80 mb-0'
              )}>
                <div className="flex items-center gap-4">
                  {/* Modern User Avatar */}
                  <div className={cn(
                    'relative rounded-xl bg-gradient-to-br flex items-center justify-center shadow-lg border border-white/10',
                    viewMode === 'list' ? 'w-10 h-10' : 'w-12 h-12',
                    user.role === 'admin' ? 'from-purple-500/20 to-violet-600/20' :
                    user.role === 'moderator' ? 'from-blue-500/20 to-indigo-600/20' :
                    'from-emerald-500/20 to-green-600/20'
                  )}>
                    {user.profile?.avatar_url ? (
                      <Image
                        src={user.profile.avatar_url}
                        alt={`${user.first_name} ${user.last_name}`}
                        width={viewMode === 'list' ? 40 : 48}
                        height={viewMode === 'list' ? 40 : 48}
                        className="object-cover rounded-xl"
                      />
                    ) : (
                      <>
                        <div className="relative z-10 text-white">
                          <User className="h-4 w-4" />
                        </div>
                        <div className="absolute inset-0 bg-white/10 rounded-xl blur-sm"></div>
                      </>
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <h3 className={cn(
                      "font-bold text-white group-hover:text-[#FF9000] transition-colors truncate",
                      viewMode === 'list' ? 'text-base' : 'text-lg'
                    )}>
                      {user.first_name} {user.last_name}
                    </h3>
                    {/* Modern Role Badge */}
                    <div className={cn(
                      'inline-flex items-center px-2 py-1 rounded-full text-xs font-semibold border border-white/10 backdrop-blur-sm',
                      viewMode === 'list' ? 'mt-1' : 'mt-2',
                      user.role === 'admin' ? 'bg-purple-100 text-purple-800 border-purple-200 dark:bg-purple-950/30 dark:text-purple-300 dark:border-purple-800' :
                      user.role === 'moderator' ? 'bg-blue-100 text-blue-800 border-blue-200 dark:bg-blue-950/30 dark:text-blue-300 dark:border-blue-800' :
                      'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-800'
                    )}>
                      <div className={cn(
                        'w-2 h-2 rounded-full mr-1',
                        user.role === 'admin' ? 'bg-purple-500' :
                        user.role === 'moderator' ? 'bg-blue-500' :
                        'bg-emerald-500'
                      )} />
                      {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                    </div>
                  </div>
                </div>

                {viewMode === 'grid' && (
                  <div className="text-right">
                    <div className={cn(
                      'inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border border-white/10 backdrop-blur-sm',
                      user.is_active
                        ? 'bg-emerald-100 text-emerald-800 border-emerald-200 dark:bg-emerald-950/30 dark:text-emerald-300 dark:border-emerald-800'
                        : 'bg-red-100 text-red-800 border-red-200 dark:bg-red-950/30 dark:text-red-300 dark:border-red-800'
                    )}>
                      <div className={cn(
                        'w-2 h-2 rounded-full mr-2',
                        user.is_active ? 'bg-emerald-500' : 'bg-red-500'
                      )} />
                      {user.is_active ? 'Active' : 'Inactive'}
                    </div>
                  </div>
                )}
              </div>

              {/* User Details */}
              <div className={cn(
                'relative space-y-3',
                viewMode === 'list' && 'flex-1 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 space-y-0'
              )}>
                {/* Email - Always show */}
                <div className="flex items-center gap-3 min-w-0">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10 flex-shrink-0">
                    <Mail className="h-4 w-4 text-gray-400" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs text-gray-400 uppercase tracking-wide font-medium">Email</p>
                    <p className="text-sm text-white font-medium truncate">
                      {user.email}
                    </p>
                  </div>
                </div>

                {/* Status - Show in list view */}
                {viewMode === 'list' && (
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10 flex-shrink-0">
                      <User className="h-4 w-4 text-gray-400" />
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="text-xs text-gray-400 uppercase tracking-wide font-medium">Status</p>
                      <div className="flex items-center gap-2">
                        <div className={cn(
                          'w-2 h-2 rounded-full',
                          user.is_active ? 'bg-emerald-500' : 'bg-red-500'
                        )} />
                        <p className={cn(
                          "text-sm font-medium",
                          user.is_active ? 'text-emerald-400' : 'text-red-400'
                        )}>
                          {user.is_active ? 'Active' : 'Inactive'}
                        </p>
                      </div>
                    </div>
                  </div>
                )}

                {/* Joined Date */}
                <div className="flex items-center gap-3 min-w-0">
                  <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10 flex-shrink-0">
                    <Calendar className="h-4 w-4 text-gray-400" />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="text-xs text-gray-400 uppercase tracking-wide font-medium">Joined</p>
                    <p className="text-sm text-white font-medium">
                      {formatDate(user.created_at)}
                    </p>
                  </div>
                </div>

                {/* Phone - Show if available */}
                {user.profile?.phone && (
                  <div className="flex items-center gap-3 min-w-0">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10 flex-shrink-0">
                      <Phone className="h-4 w-4 text-gray-400" />
                    </div>
                    <div className="min-w-0 flex-1">
                      <p className="text-xs text-gray-400 uppercase tracking-wide font-medium">Phone</p>
                      <p className="text-sm text-white font-medium truncate">
                        {user.profile.phone}
                      </p>
                    </div>
                  </div>
                )}

                {/* Show Status in grid view only */}
                {viewMode === 'grid' && (
                  <div className="flex items-center gap-3">
                    <div className="w-8 h-8 rounded-lg bg-white/5 flex items-center justify-center border border-white/10">
                      <User className="h-4 w-4 text-gray-400" />
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 uppercase tracking-wide">Status</p>
                      <p className="text-sm text-white font-medium">
                        {user.is_active ? 'Active User' : 'Inactive User'}
                      </p>
                    </div>
                  </div>
                )}
              </div>

              {/* Modern Action Buttons */}
              <div className={cn(
                'relative flex items-center gap-3 mt-6 pt-4 border-t border-white/10',
                viewMode === 'list' && 'flex-shrink-0 mt-0 pt-0 border-t-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleViewUser(user)}
                  className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-blue-500/10 hover:border-blue-500/50 hover:text-blue-400 transition-all duration-200"
                >
                  <Eye className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'View Profile' : 'View'}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className="group/btn relative bg-white/5 border-gray-600/50 text-gray-300 hover:bg-purple-500/10 hover:border-purple-500/50 hover:text-purple-400 transition-all duration-200"
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56 bg-gray-900/95 backdrop-blur-sm border-gray-700/50 rounded-xl shadow-xl">
                    <DropdownMenuItem 
                      onClick={() => handleViewUser(user)}
                      className="text-gray-300 hover:text-white hover:bg-gray-800/50 rounded-lg m-1 p-3"
                    >
                      <Eye className="mr-3 h-4 w-4" />
                      View Profile
                    </DropdownMenuItem>
                    <DropdownMenuItem 
                      onClick={() => handleEditUser(user)}
                      className="text-blue-400 hover:text-blue-300 hover:bg-blue-900/20 rounded-lg m-1 p-3"
                    >
                      <Edit className="mr-3 h-4 w-4" />
                      Edit User
                    </DropdownMenuItem>
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    {user.is_active ? (
                      <DropdownMenuItem 
                        onClick={() => handleToggleUserStatus(user)}
                        className="text-red-400 hover:text-red-300 hover:bg-red-900/20 rounded-lg m-1 p-3"
                      >
                        <UserX className="mr-3 h-4 w-4" />
                        Deactivate User
                      </DropdownMenuItem>
                    ) : (
                      <DropdownMenuItem 
                        onClick={() => handleToggleUserStatus(user)}
                        className="text-green-400 hover:text-green-300 hover:bg-green-900/20 rounded-lg m-1 p-3"
                      >
                        <UserCheck className="mr-3 h-4 w-4" />
                        Activate User
                      </DropdownMenuItem>
                    )}
                    <RequirePermission permission={PERMISSIONS.USERS_MANAGE_ROLES}>
                      <DropdownMenuItem 
                        onClick={() => {
                          const roles = ['admin', 'moderator', 'customer']
                          const roleOptions = roles.filter(role => role !== user.role).join(', ')
                          const newRole = prompt(
                            `Change role for ${user.first_name} ${user.last_name}\n\nCurrent role: ${user.role}\nAvailable roles: ${roleOptions}\n\nEnter new role:`, 
                            user.role
                          )
                          if (newRole && newRole !== user.role && roles.includes(newRole.toLowerCase())) {
                            handleChangeUserRole(user, newRole.toLowerCase())
                          } else if (newRole && !roles.includes(newRole.toLowerCase())) {
                            toast.error('Invalid role', {
                              description: `Valid roles are: ${roles.join(', ')}`,
                              duration: 3000,
                            })
                          }
                        }}
                        className="text-purple-400 hover:text-purple-300 hover:bg-purple-900/20 rounded-lg m-1 p-3"
                      >
                        <Shield className="mr-3 h-4 w-4" />
                        Change Role
                      </DropdownMenuItem>
                    </RequirePermission>
                    <div className="h-px bg-gray-700/50 mx-2 my-1"></div>
                    <RequirePermission permission={PERMISSIONS.USERS_DELETE}>
                      <DropdownMenuItem 
                        onClick={() => handleDeleteUser(user)}
                        className="text-red-500 hover:text-red-400 hover:bg-red-900/20 rounded-lg m-1 p-3"
                      >
                        <UserX className="mr-3 h-4 w-4" />
                        Delete User
                      </DropdownMenuItem>
                    </RequirePermission>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-2xl p-12 text-center">
          <div className="absolute inset-0 bg-gradient-to-br from-gray-500/5 to-slate-500/5 rounded-2xl"></div>
          <div className="relative">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-gray-500/20 to-slate-500/20 border border-gray-400/30 flex items-center justify-center mx-auto mb-6">
              <Users className="h-8 w-8 text-gray-400" />
            </div>
            <h3 className="text-xl font-bold text-white mb-2">No users found</h3>
            <p className="text-gray-400 mb-6 max-w-md mx-auto">
              {searchQuery || roleFilter || statusFilter
                ? 'Try adjusting your search or filters to find users.'
                : 'There are no users in the BiHub system yet.'
              }
            </p>
          </div>
        </div>
      )}

      {/* Pagination */}
      {pagination && pagination.total_pages > 1 && (
        <BiHubAdminCard
          title="Pagination"
          subtitle={`Showing ${((currentPage - 1) * (pagination.limit || 20)) + 1} to ${Math.min(currentPage * (pagination.limit || 20), pagination.total)} of ${pagination.total} users`}
        >
          <div className="flex items-center justify-between">
            <p className={BIHUB_ADMIN_THEME.typography.body.medium}>
              Page {currentPage} of {pagination.total_pages}
            </p>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Previous
              </Button>

              <div className="flex items-center gap-1">
                {Array.from({ length: Math.min(5, pagination.total_pages) }, (_, i) => {
                  const page = i + 1
                  return (
                    <Button
                      key={page}
                      variant={currentPage === page ? "default" : "ghost"}
                      onClick={() => setCurrentPage(page)}
                      className={cn(
                        'w-10 h-10',
                        currentPage === page
                          ? BIHUB_ADMIN_THEME.components.button.primary
                          : BIHUB_ADMIN_THEME.components.button.ghost
                      )}
                    >
                      {page}
                    </Button>
                  )
                })}
              </div>

              <Button
                variant="outline"
                onClick={() => setCurrentPage(prev => Math.min(pagination.total_pages, prev + 1))}
                disabled={currentPage === pagination.total_pages}
                className={BIHUB_ADMIN_THEME.components.button.secondary}
              >
                Next
              </Button>
            </div>
          </div>
        </BiHubAdminCard>
      )}
    </div>
  )
}
