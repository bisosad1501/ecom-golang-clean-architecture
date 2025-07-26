'use client'

import { useState } from 'react'
import {
  Search,
  Filter,
  MoreHorizontal,
  Edit,
  Trash2,
  Eye,
  Users,
  UserCheck,
  Crown,
  Shield,
  Mail,
  Phone,
  Calendar,
  DollarSign,
  ShoppingBag,
  TrendingUp,
  Download,
  UserPlus
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { formatPrice, formatDate } from '@/lib/utils'
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

interface Customer {
  id: string
  first_name: string
  last_name: string
  email: string
  phone?: string
  role: 'customer' | 'moderator' | 'admin'
  is_active: boolean
  total_orders: number
  total_spent: number
  last_order_date?: string
  created_at: string
  avatar?: string
}

export default function AdminCustomersPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCustomers, setSelectedCustomers] = useState<string[]>([])
  const [filterRole, setFilterRole] = useState<string>('all')
  const [filterStatus, setFilterStatus] = useState<string>('all')

  // Mock data - replace with real API calls
  const customers: Customer[] = [
    {
      id: '1',
      first_name: 'John',
      last_name: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      role: 'customer',
      is_active: true,
      total_orders: 15,
      total_spent: 2450.50,
      last_order_date: '2024-01-15T10:30:00Z',
      created_at: '2023-06-15T09:00:00Z'
    },
    {
      id: '2',
      first_name: 'Jane',
      last_name: 'Smith',
      email: 'jane.smith@example.com',
      role: 'customer',
      is_active: true,
      total_orders: 8,
      total_spent: 1200.00,
      last_order_date: '2024-01-10T14:20:00Z',
      created_at: '2023-08-20T11:15:00Z'
    },
    {
      id: '3',
      first_name: 'Admin',
      last_name: 'User',
      email: 'admin@bihub.com',
      role: 'admin',
      is_active: true,
      total_orders: 0,
      total_spent: 0,
      created_at: '2023-01-01T00:00:00Z'
    }
  ]

  const filteredCustomers = customers.filter(customer => {
    const matchesSearch = 
      customer.first_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      customer.last_name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      customer.email.toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesRole = filterRole === 'all' || customer.role === filterRole
    const matchesStatus = filterStatus === 'all' || 
      (filterStatus === 'active' && customer.is_active) ||
      (filterStatus === 'inactive' && !customer.is_active)
    
    return matchesSearch && matchesRole && matchesStatus
  })

  const stats = {
    total: customers.length,
    active: customers.filter(c => c.is_active).length,
    totalSpent: customers.reduce((sum, c) => sum + c.total_spent, 0),
    avgOrderValue: customers.reduce((sum, c) => sum + c.total_spent, 0) / Math.max(customers.reduce((sum, c) => sum + c.total_orders, 0), 1)
  }

  const getRoleIcon = (role: string) => {
    switch (role) {
      case 'admin':
        return <Crown className="h-4 w-4 text-yellow-400" />
      case 'moderator':
        return <Shield className="h-4 w-4 text-blue-400" />
      default:
        return <UserCheck className="h-4 w-4 text-green-400" />
    }
  }

  const getRoleBadgeVariant = (role: string) => {
    switch (role) {
      case 'admin':
        return 'warning'
      case 'moderator':
        return 'info'
      default:
        return 'success'
    }
  }

  const handleDeleteCustomer = (customerId: string) => {
    console.log('Delete customer:', customerId)
  }

  const handleToggleStatus = (customerId: string) => {
    console.log('Toggle status:', customerId)
  }

  const handleViewCustomer = (customerId: string) => {
    console.log('View customer:', customerId)
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="Customer Management"
        subtitle="Manage customer accounts, roles, and analytics"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Customers' }
        ]}
        action={
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
            <RequirePermission permission={PERMISSIONS.USERS_CREATE}>
              <Button className={BIHUB_ADMIN_THEME.components.button.primary}>
                <UserPlus className="mr-2 h-5 w-5" />
                Add Customer
              </Button>
            </RequirePermission>
          </div>
        }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div className="bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-400/80 text-sm font-medium">Total Customers</p>
              <p className="text-2xl font-bold text-blue-100">{stats.total}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
              <Users className="h-5 w-5 text-blue-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-400/80 text-sm font-medium">Active Customers</p>
              <p className="text-2xl font-bold text-green-100">{stats.active}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-green-500/20 flex items-center justify-center">
              <UserCheck className="h-5 w-5 text-green-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-purple-400/80 text-sm font-medium">Total Revenue</p>
              <p className="text-2xl font-bold text-purple-100">{formatPrice(stats.totalSpent)}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-purple-500/20 flex items-center justify-center">
              <DollarSign className="h-5 w-5 text-purple-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-orange-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-orange-400/80 text-sm font-medium">Avg Order Value</p>
              <p className="text-2xl font-bold text-orange-100">{formatPrice(stats.avgOrderValue)}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-orange-500/20 flex items-center justify-center">
              <TrendingUp className="h-5 w-5 text-orange-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Customer List */}
      <BiHubAdminCard
        title="Customer List"
        subtitle="Manage customer accounts and permissions"
        icon={<Users className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search customers..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-64"
              />
            </div>
            <Button variant="outline" size="sm">
              <Filter className="h-4 w-4 mr-2" />
              Filter
            </Button>
          </div>
        }
      >
        <div className="space-y-4">
          {filteredCustomers.map((customer) => (
            <div
              key={customer.id}
              className="group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                    <span className="text-white font-bold text-lg">
                      {customer.first_name[0]}{customer.last_name[0]}
                    </span>
                  </div>
                  
                  <div>
                    <div className="flex items-center gap-3">
                      <h3 className="text-lg font-semibold text-white">
                        {customer.first_name} {customer.last_name}
                      </h3>
                      <div className="flex items-center gap-1">
                        {getRoleIcon(customer.role)}
                        <BiHubStatusBadge status={customer.role}>
                          {customer.role}
                        </BiHubStatusBadge>
                      </div>
                      <BiHubStatusBadge status={customer.is_active ? 'active' : 'inactive'}>
                        {customer.is_active ? 'Active' : 'Inactive'}
                      </BiHubStatusBadge>
                    </div>
                    <div className="flex items-center gap-2 mt-1">
                      <Mail className="h-4 w-4 text-gray-400" />
                      <span className="text-gray-400 text-sm">{customer.email}</span>
                      {customer.phone && (
                        <>
                          <span className="text-gray-600">•</span>
                          <Phone className="h-4 w-4 text-gray-400" />
                          <span className="text-gray-400 text-sm">{customer.phone}</span>
                        </>
                      )}
                    </div>
                    <div className="flex items-center gap-4 mt-2 text-sm text-gray-400">
                      <div className="flex items-center gap-1">
                        <ShoppingBag className="h-4 w-4" />
                        <span>{customer.total_orders} orders</span>
                      </div>
                      <span>•</span>
                      <div className="flex items-center gap-1">
                        <DollarSign className="h-4 w-4" />
                        <span>{formatPrice(customer.total_spent)} spent</span>
                      </div>
                      <span>•</span>
                      <div className="flex items-center gap-1">
                        <Calendar className="h-4 w-4" />
                        <span>Joined {formatDate(customer.created_at)}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleViewCustomer(customer.id)}
                    className="opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <Eye className="h-4 w-4" />
                  </Button>
                  
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => handleViewCustomer(customer.id)}>
                        <Eye className="h-4 w-4 mr-2" />
                        View Details
                      </DropdownMenuItem>
                      <RequirePermission permission={PERMISSIONS.USERS_UPDATE}>
                        <DropdownMenuItem>
                          <Edit className="h-4 w-4 mr-2" />
                          Edit Customer
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => handleToggleStatus(customer.id)}>
                          <UserCheck className="h-4 w-4 mr-2" />
                          {customer.is_active ? 'Deactivate' : 'Activate'}
                        </DropdownMenuItem>
                      </RequirePermission>
                      <RequirePermission permission={PERMISSIONS.USERS_DELETE}>
                        <DropdownMenuItem 
                          onClick={() => handleDeleteCustomer(customer.id)}
                          className="text-red-600"
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </RequirePermission>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            </div>
          ))}

          {filteredCustomers.length === 0 && (
            <div className="text-center py-12">
              <Users className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-300 mb-2">No customers found</h3>
              <p className="text-gray-400 mb-6">
                {searchQuery ? 'Try adjusting your search terms.' : 'No customers registered yet.'}
              </p>
            </div>
          )}
        </div>
      </BiHubAdminCard>
    </div>
  )
}
