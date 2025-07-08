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

  const getTotalStats = () => {
    const totalUsers = pagination?.total || 0
    const activeUsers = users.filter(user => user.is_active).length
    const adminUsers = users.filter(user => user.role === 'admin').length
    const customerUsers = users.filter(user => user.role === 'customer').length
    
    return { totalUsers, activeUsers, adminUsers, customerUsers }
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
      {/* BiHub Page Header */}
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

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <BiHubStatCard
          title="Total Users"
          value={stats.totalUsers}
          icon={<Users className="h-8 w-8 text-white" />}
          color="primary"
        />
        <BiHubStatCard
          title="Active Users"
          value={stats.activeUsers}
          icon={<UserCheck className="h-8 w-8 text-white" />}
          color="success"
        />
        <BiHubStatCard
          title="Admins"
          value={stats.adminUsers}
          icon={<Crown className="h-8 w-8 text-white" />}
          color="error"
        />
        <BiHubStatCard
          title="Customers"
          value={stats.customerUsers}
          icon={<User className="h-8 w-8 text-white" />}
          color="info"
        />
      </div>

      {/* Search & Filters */}
      <BiHubAdminCard
        title="Search & Filter Users"
        subtitle="Find and filter BiHub users by role and status"
        icon={<Search className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setViewMode(viewMode === 'grid' ? 'list' : 'grid')}
              className={cn(
                BIHUB_ADMIN_THEME.components.button.ghost,
                'h-10 w-10 p-0'
              )}
            >
              {viewMode === 'grid' ? (
                <List className="h-4 w-4" />
              ) : (
                <Grid className="h-4 w-4" />
              )}
            </Button>
          </div>
        }
      >
        <div className="flex flex-col lg:flex-row items-center gap-4">
          <div className="flex-1 w-full">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
              <Input
                placeholder="Search users by name, email..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={cn(
                  BIHUB_ADMIN_THEME.components.input.base,
                  'pl-10 pr-12 h-12'
                )}
              />
              {searchQuery && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setSearchQuery('')}
                  className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 text-gray-400 hover:text-white"
                >
                  ×
                </Button>
              )}
            </div>
          </div>

          <div className="flex items-center gap-4">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.secondary}>
                  <Filter className="mr-2 h-4 w-4" />
                  Role: {roleFilter || 'All'}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48 bg-gray-900 border-gray-700">
                <DropdownMenuItem
                  onClick={() => setRoleFilter('')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  All Roles
                </DropdownMenuItem>
                <DropdownMenuSeparator className="bg-gray-700" />
                <DropdownMenuItem
                  onClick={() => setRoleFilter('admin')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  Admin
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setRoleFilter('moderator')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  Moderator
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setRoleFilter('customer')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  Customer
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.secondary}>
                  <Filter className="mr-2 h-4 w-4" />
                  Status: {statusFilter || 'All'}
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" className="w-48 bg-gray-900 border-gray-700">
                <DropdownMenuItem
                  onClick={() => setStatusFilter('')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  All Status
                </DropdownMenuItem>
                <DropdownMenuSeparator className="bg-gray-700" />
                <DropdownMenuItem
                  onClick={() => setStatusFilter('active')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  Active
                </DropdownMenuItem>
                <DropdownMenuItem
                  onClick={() => setStatusFilter('inactive')}
                  className="text-gray-300 hover:text-white hover:bg-gray-800"
                >
                  Inactive
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </BiHubAdminCard>

      {/* Users List */}
      {isLoading ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {[...Array(6)].map((_, i) => (
            <div key={i} className={cn(
              BIHUB_ADMIN_THEME.components.card.base,
              'p-6 animate-pulse'
            )}>
              <div className="flex items-center gap-4 mb-4">
                <div className="w-12 h-12 bg-gray-700 rounded-full"></div>
                <div className="space-y-2 flex-1">
                  <div className="h-4 bg-gray-700 rounded w-1/3"></div>
                  <div className="h-3 bg-gray-700 rounded w-1/2"></div>
                </div>
              </div>
              <div className="space-y-2">
                <div className="h-3 bg-gray-700 rounded w-3/4"></div>
                <div className="h-3 bg-gray-700 rounded w-1/2"></div>
              </div>
            </div>
          ))}
        </div>
      ) : users.length > 0 ? (
        <div className={cn(
          viewMode === 'grid'
            ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
            : 'space-y-4'
        )}>
          {users.map((user) => (
            <div
              key={user.id}
              className={cn(
                BIHUB_ADMIN_THEME.components.card.base,
                BIHUB_ADMIN_THEME.components.card.hover,
                'group',
                viewMode === 'list' && 'flex items-center gap-6 p-6'
              )}
            >
              {/* User Avatar & Status */}
              <div className={cn(
                'relative flex-shrink-0',
                viewMode === 'grid' ? 'mb-4' : ''
              )}>
                <div className={cn(
                  'rounded-full overflow-hidden bg-gray-700 flex items-center justify-center',
                  viewMode === 'grid' ? 'w-16 h-16 mx-auto' : 'w-12 h-12'
                )}>
                  {user.profile?.avatar_url ? (
                    <Image
                      src={user.profile.avatar_url}
                      alt={`${user.first_name} ${user.last_name}`}
                      width={viewMode === 'grid' ? 64 : 48}
                      height={viewMode === 'grid' ? 64 : 48}
                      className="object-cover"
                    />
                  ) : (
                    <User className={cn(
                      'text-gray-400',
                      viewMode === 'grid' ? 'w-8 h-8' : 'w-6 h-6'
                    )} />
                  )}
                </div>

                {/* Status indicator */}
                <div className={cn(
                  'absolute -bottom-1 -right-1 w-4 h-4 rounded-full border-2 border-gray-800',
                  user.is_active ? 'bg-emerald-500' : 'bg-gray-400'
                )} />
              </div>

              {/* User Details */}
              <div className={cn(
                'flex-1',
                viewMode === 'grid' ? 'text-center p-6' : ''
              )}>
                <div className={cn(
                  'flex items-center gap-3 mb-3',
                  viewMode === 'grid' && 'flex-col gap-2'
                )}>
                  <h3 className={cn(
                    BIHUB_ADMIN_THEME.typography.heading.h4,
                    'group-hover:text-[#FF9000] transition-colors'
                  )}>
                    {user.first_name} {user.last_name}
                  </h3>

                  <div className="flex items-center gap-2 flex-wrap">
                    {/* Role Badge */}
                    <div className={cn(
                      'flex items-center gap-1 px-3 py-1 rounded-full text-xs font-semibold',
                      `bg-gradient-to-r ${getRoleColor(user.role)} text-white`
                    )}>
                      {getRoleIcon(user.role)}
                      <span className="ml-1">{user.role.charAt(0).toUpperCase() + user.role.slice(1)}</span>
                    </div>

                    {/* Status Badge */}
                    <BiHubStatusBadge status={user.is_active ? 'success' : 'error'}>
                      {user.is_active ? 'Active' : 'Inactive'}
                    </BiHubStatusBadge>
                  </div>
                </div>

                <div className={cn(
                  'space-y-2 text-sm',
                  viewMode === 'grid' ? 'text-center' : 'grid grid-cols-1 md:grid-cols-2 gap-2'
                )}>
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4 text-gray-400" />
                    <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                      {user.email}
                    </span>
                  </div>

                  {user.profile?.phone && (
                    <div className="flex items-center gap-2">
                      <Phone className="h-4 w-4 text-gray-400" />
                      <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                        {user.profile.phone}
                      </span>
                    </div>
                  )}

                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4 text-gray-400" />
                    <span className={BIHUB_ADMIN_THEME.typography.body.small}>
                      {formatDate(user.created_at)}
                    </span>
                  </div>
                </div>
              </div>

              {/* Actions */}
              <div className={cn(
                'flex items-center gap-2',
                viewMode === 'list' && 'flex-shrink-0'
              )}>
                <Button
                  variant="outline"
                  size="sm"
                  className={BIHUB_ADMIN_THEME.components.button.ghost}
                >
                  <Eye className="h-4 w-4 mr-2" />
                  {viewMode === 'grid' ? 'View' : ''}
                </Button>

                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="outline"
                      size="sm"
                      className={BIHUB_ADMIN_THEME.components.button.ghost}
                    >
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="bg-gray-900 border-gray-700">
                    <DropdownMenuItem className="text-gray-300 hover:text-white hover:bg-gray-800">
                      <Eye className="h-4 w-4 mr-2" />
                      View Details
                    </DropdownMenuItem>
                    <DropdownMenuSeparator className="bg-gray-700" />
                    {user.is_active ? (
                      <DropdownMenuItem className="text-red-400 hover:text-red-300 hover:bg-red-900/20">
                        <UserX className="h-4 w-4 mr-2" />
                        Deactivate User
                      </DropdownMenuItem>
                    ) : (
                      <DropdownMenuItem className="text-green-400 hover:text-green-300 hover:bg-green-900/20">
                        <UserCheck className="h-4 w-4 mr-2" />
                        Activate User
                      </DropdownMenuItem>
                    )}
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <BiHubEmptyState
          icon={<Users className="h-8 w-8 text-gray-400" />}
          title="No users found"
          description={
            searchQuery || roleFilter || statusFilter
              ? 'Try adjusting your search or filters to find users.'
              : 'There are no users in the BiHub system yet.'
          }
        />
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
