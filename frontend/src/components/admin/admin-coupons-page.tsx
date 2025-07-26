'use client'

import { useState } from 'react'
import {
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Edit,
  Trash2,
  Eye,
  Copy,
  Ticket,
  Calendar,
  DollarSign,
  Percent,
  Users,
  ShoppingCart,
  TrendingUp,
  Download
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

interface Coupon {
  id: string
  code: string
  name: string
  type: 'percentage' | 'fixed'
  value: number
  min_order_amount?: number
  max_discount_amount?: number
  usage_limit?: number
  used_count: number
  start_date: string
  end_date: string
  is_active: boolean
  created_at: string
}

export default function AdminCouponsPage() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedCoupons, setSelectedCoupons] = useState<string[]>([])
  const [showAddForm, setShowAddForm] = useState(false)
  const [editingCoupon, setEditingCoupon] = useState<Coupon | null>(null)

  // Mock data - replace with real API calls
  const coupons: Coupon[] = [
    {
      id: '1',
      code: 'WELCOME10',
      name: 'Welcome Discount',
      type: 'percentage',
      value: 10,
      min_order_amount: 50,
      usage_limit: 1000,
      used_count: 245,
      start_date: '2024-01-01',
      end_date: '2024-12-31',
      is_active: true,
      created_at: '2024-01-01T00:00:00Z'
    },
    {
      id: '2',
      code: 'SAVE20',
      name: 'Fixed Discount',
      type: 'fixed',
      value: 20,
      min_order_amount: 100,
      max_discount_amount: 50,
      usage_limit: 500,
      used_count: 89,
      start_date: '2024-02-01',
      end_date: '2024-06-30',
      is_active: true,
      created_at: '2024-02-01T00:00:00Z'
    }
  ]

  const filteredCoupons = coupons.filter(coupon =>
    coupon.code.toLowerCase().includes(searchQuery.toLowerCase()) ||
    coupon.name.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const stats = {
    total: coupons.length,
    active: coupons.filter(c => c.is_active).length,
    totalUsage: coupons.reduce((sum, c) => sum + c.used_count, 0),
    totalSavings: 15420 // Mock value
  }

  const handleDeleteCoupon = (couponId: string) => {
    // Implement delete logic
    console.log('Delete coupon:', couponId)
  }

  const handleToggleStatus = (couponId: string) => {
    // Implement toggle status logic
    console.log('Toggle status:', couponId)
  }

  const handleCopyCoupon = (couponId: string) => {
    // Implement copy logic
    console.log('Copy coupon:', couponId)
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="Coupon Management"
        subtitle="Create and manage discount coupons and promotional codes"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Coupons' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.COUPONS_CREATE}>
            <Button
              onClick={() => setShowAddForm(true)}
              className={BIHUB_ADMIN_THEME.components.button.primary}
            >
              <Plus className="mr-2 h-5 w-5" />
              Create Coupon
            </Button>
          </RequirePermission>
        }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div className="bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-400/80 text-sm font-medium">Total Coupons</p>
              <p className="text-2xl font-bold text-blue-100">{stats.total}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
              <Ticket className="h-5 w-5 text-blue-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-400/80 text-sm font-medium">Active Coupons</p>
              <p className="text-2xl font-bold text-green-100">{stats.active}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-green-500/20 flex items-center justify-center">
              <TrendingUp className="h-5 w-5 text-green-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-purple-400/80 text-sm font-medium">Total Usage</p>
              <p className="text-2xl font-bold text-purple-100">{stats.totalUsage.toLocaleString()}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-purple-500/20 flex items-center justify-center">
              <Users className="h-5 w-5 text-purple-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-orange-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-orange-400/80 text-sm font-medium">Total Savings</p>
              <p className="text-2xl font-bold text-orange-100">{formatPrice(stats.totalSavings)}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-orange-500/20 flex items-center justify-center">
              <DollarSign className="h-5 w-5 text-orange-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Filters and Search */}
      <BiHubAdminCard
        title="Coupon List"
        subtitle="Manage your promotional codes and discounts"
        icon={<Ticket className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search coupons..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-64"
              />
            </div>
            <Button variant="outline" size="sm">
              <Filter className="h-4 w-4 mr-2" />
              Filter
            </Button>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          </div>
        }
      >
        <div className="space-y-4">
          {filteredCoupons.map((coupon) => (
            <div
              key={coupon.id}
              className="group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                    {coupon.type === 'percentage' ? (
                      <Percent className="h-6 w-6 text-white" />
                    ) : (
                      <DollarSign className="h-6 w-6 text-white" />
                    )}
                  </div>
                  
                  <div>
                    <div className="flex items-center gap-3">
                      <h3 className="text-lg font-semibold text-white">{coupon.code}</h3>
                      <BiHubStatusBadge status={coupon.is_active ? 'active' : 'inactive'}>
                        {coupon.is_active ? 'Active' : 'Inactive'}
                      </BiHubStatusBadge>
                    </div>
                    <p className="text-gray-400 text-sm">{coupon.name}</p>
                    <div className="flex items-center gap-4 mt-2 text-sm text-gray-400">
                      <span>
                        {coupon.type === 'percentage' ? `${coupon.value}% off` : `${formatPrice(coupon.value)} off`}
                      </span>
                      <span>•</span>
                      <span>{coupon.used_count}/{coupon.usage_limit || '∞'} used</span>
                      <span>•</span>
                      <span>Expires {formatDate(coupon.end_date)}</span>
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleCopyCoupon(coupon.id)}
                    className="opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <Copy className="h-4 w-4" />
                  </Button>
                  
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => setEditingCoupon(coupon)}>
                        <Edit className="h-4 w-4 mr-2" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => handleToggleStatus(coupon.id)}>
                        <Eye className="h-4 w-4 mr-2" />
                        {coupon.is_active ? 'Deactivate' : 'Activate'}
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => handleDeleteCoupon(coupon.id)}
                        className="text-red-600"
                      >
                        <Trash2 className="h-4 w-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            </div>
          ))}

          {filteredCoupons.length === 0 && (
            <div className="text-center py-12">
              <Ticket className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-300 mb-2">No coupons found</h3>
              <p className="text-gray-400 mb-6">
                {searchQuery ? 'Try adjusting your search terms.' : 'Create your first coupon to get started.'}
              </p>
              {!searchQuery && (
                <RequirePermission permission={PERMISSIONS.COUPONS_CREATE}>
                  <Button
                    onClick={() => setShowAddForm(true)}
                    className={BIHUB_ADMIN_THEME.components.button.primary}
                  >
                    <Plus className="mr-2 h-5 w-5" />
                    Create First Coupon
                  </Button>
                </RequirePermission>
              )}
            </div>
          )}
        </div>
      </BiHubAdminCard>
    </div>
  )
}
