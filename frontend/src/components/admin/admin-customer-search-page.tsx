'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Search,
  Filter,
  Users,
  TrendingUp,
  Crown,
  Shield,
  User,
  Mail,
  Phone,
  Calendar,
  DollarSign,
  ShoppingBag,
  Star,
  Eye,
  Download,
  BarChart3,
} from 'lucide-react'
import { formatDate, formatCurrency, cn } from '@/lib/utils'
import {
  useCustomerSearch,
  useCustomerSegments,
  useCustomerAnalytics,
  useHighValueCustomers,
  CustomerSearchFilters,
} from '@/hooks/use-users'
import { toast } from 'sonner'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

export default function AdminCustomerSearchPage() {
  const [searchFilters, setSearchFilters] = useState<CustomerSearchFilters>({
    limit: 20,
    offset: 0,
  })
  const [activeTab, setActiveTab] = useState<'search' | 'segments' | 'analytics' | 'high-value'>('search')

  // Hooks for data fetching
  const { data: searchResults, isLoading: searchLoading, error: searchError } = useCustomerSearch(searchFilters)
  const { data: segments, isLoading: segmentsLoading } = useCustomerSegments()
  const { data: analytics, isLoading: analyticsLoading } = useCustomerAnalytics()
  const { data: highValueCustomers, isLoading: highValueLoading } = useHighValueCustomers()

  const handleSearchChange = (field: keyof CustomerSearchFilters, value: any) => {
    setSearchFilters(prev => ({
      ...prev,
      [field]: value,
      offset: 0, // Reset to first page when filters change
    }))
  }

  const handleExportCustomers = () => {
    toast.success('Customer export started. You will receive an email when ready.')
  }

  const getRoleIcon = (role: string) => {
    switch (role.toLowerCase()) {
      case 'admin':
        return <Crown className="h-4 w-4" />
      case 'moderator':
        return <Shield className="h-4 w-4" />
      case 'customer':
        return <User className="h-4 w-4" />
      default:
        return <User className="h-4 w-4" />
    }
  }

  const getSegmentColor = (segment: string) => {
    switch (segment.toLowerCase()) {
      case 'new':
        return 'bg-blue-100 text-blue-800'
      case 'occasional':
        return 'bg-yellow-100 text-yellow-800'
      case 'regular':
        return 'bg-green-100 text-green-800'
      case 'loyal':
        return 'bg-purple-100 text-purple-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getTierColor = (tier: string) => {
    switch (tier.toLowerCase()) {
      case 'bronze':
        return 'bg-orange-100 text-orange-800'
      case 'silver':
        return 'bg-gray-100 text-gray-800'
      case 'gold':
        return 'bg-yellow-100 text-yellow-800'
      case 'platinum':
        return 'bg-purple-100 text-purple-800'
      case 'diamond':
        return 'bg-blue-100 text-blue-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const renderSearchTab = () => (
    <div className="space-y-6">
      {/* Search Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Customer Search Filters
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <label className="text-sm font-medium mb-2 block">Search Query</label>
              <Input
                placeholder="Search customers..."
                value={searchFilters.query || ''}
                onChange={(e) => handleSearchChange('query', e.target.value)}
                className="w-full"
              />
            </div>
            
            <div>
              <label className="text-sm font-medium mb-2 block">Role</label>
              <Select
                value={searchFilters.role || ''}
                onValueChange={(value) => handleSearchChange('role', value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All roles" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All roles</SelectItem>
                  <SelectItem value="customer">Customer</SelectItem>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="moderator">Moderator</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <label className="text-sm font-medium mb-2 block">Segment</label>
              <Select
                value={searchFilters.customer_segment || ''}
                onValueChange={(value) => handleSearchChange('customer_segment', value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All segments" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All segments</SelectItem>
                  <SelectItem value="new">New</SelectItem>
                  <SelectItem value="occasional">Occasional</SelectItem>
                  <SelectItem value="regular">Regular</SelectItem>
                  <SelectItem value="loyal">Loyal</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <label className="text-sm font-medium mb-2 block">Membership Tier</label>
              <Select
                value={searchFilters.membership_tier || ''}
                onValueChange={(value) => handleSearchChange('membership_tier', value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All tiers" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">All tiers</SelectItem>
                  <SelectItem value="bronze">Bronze</SelectItem>
                  <SelectItem value="silver">Silver</SelectItem>
                  <SelectItem value="gold">Gold</SelectItem>
                  <SelectItem value="platinum">Platinum</SelectItem>
                  <SelectItem value="diamond">Diamond</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
            <div>
              <label className="text-sm font-medium mb-2 block">Min Total Spent</label>
              <Input
                type="number"
                placeholder="0"
                value={searchFilters.min_total_spent || ''}
                onChange={(e) => handleSearchChange('min_total_spent', parseFloat(e.target.value) || undefined)}
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-2 block">Max Total Spent</label>
              <Input
                type="number"
                placeholder="No limit"
                value={searchFilters.max_total_spent || ''}
                onChange={(e) => handleSearchChange('max_total_spent', parseFloat(e.target.value) || undefined)}
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-2 block">Min Orders</label>
              <Input
                type="number"
                placeholder="0"
                value={searchFilters.min_total_orders || ''}
                onChange={(e) => handleSearchChange('min_total_orders', parseInt(e.target.value) || undefined)}
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-2 block">Max Orders</label>
              <Input
                type="number"
                placeholder="No limit"
                value={searchFilters.max_total_orders || ''}
                onChange={(e) => handleSearchChange('max_total_orders', parseInt(e.target.value) || undefined)}
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Search Results */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Search Results
              {searchResults && (
                <Badge variant="secondary">
                  {searchResults.total} customers
                </Badge>
              )}
            </div>
            <Button onClick={handleExportCustomers} variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {searchLoading ? (
            <div className="text-center py-8">Loading customers...</div>
          ) : searchError ? (
            <div className="text-center py-8 text-red-600">
              Error loading customers: {searchError.message}
            </div>
          ) : searchResults?.customers.length === 0 ? (
            <BiHubEmptyState
              icon={<Users className="h-8 w-8 text-gray-400" />}
              title="No customers found"
              description="Try adjusting your search filters"
            />
          ) : (
            <div className="space-y-4">
              {searchResults?.customers.map((customer) => (
                <div
                  key={customer.id}
                  className="border rounded-lg p-4 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <div className="flex items-center gap-2">
                          {getRoleIcon(customer.role)}
                          <h3 className="font-semibold">
                            {customer.first_name} {customer.last_name}
                          </h3>
                        </div>
                        <Badge className={getSegmentColor(customer.customer_segment)}>
                          {customer.customer_segment}
                        </Badge>
                        <Badge className={getTierColor(customer.membership_tier)}>
                          {customer.membership_tier}
                        </Badge>
                        {customer.is_vip && (
                          <Badge className="bg-gold-100 text-gold-800">
                            <Crown className="h-3 w-3 mr-1" />
                            VIP
                          </Badge>
                        )}
                      </div>
                      
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-gray-600">
                        <div className="flex items-center gap-1">
                          <Mail className="h-4 w-4" />
                          {customer.email}
                        </div>
                        <div className="flex items-center gap-1">
                          <ShoppingBag className="h-4 w-4" />
                          {customer.order_count} orders
                        </div>
                        <div className="flex items-center gap-1">
                          <DollarSign className="h-4 w-4" />
                          {formatCurrency(customer.total_spent)}
                        </div>
                        <div className="flex items-center gap-1">
                          <Star className="h-4 w-4" />
                          {customer.loyalty_points} points
                        </div>
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      <BiHubStatusBadge
                        status={customer.is_active ? 'active' : 'inactive'}
                        variant={customer.is_active ? 'success' : 'secondary'}
                      />
                      <Button variant="ghost" size="sm">
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )

  const renderSegmentsTab = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            Customer Segments
          </CardTitle>
        </CardHeader>
        <CardContent>
          {segmentsLoading ? (
            <div className="text-center py-8">Loading segments...</div>
          ) : segments ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              {segments.segments.map((segment) => (
                <BiHubStatCard
                  key={segment.segment}
                  title={segment.segment.charAt(0).toUpperCase() + segment.segment.slice(1)}
                  value={segment.count.toString()}
                  subtitle={`${segment.percentage.toFixed(1)}% of customers`}
                  icon={<Users className="h-6 w-6" />}
                  trend={{
                    value: segment.avg_spent,
                    label: `Avg: ${formatCurrency(segment.avg_spent)}`,
                    isPositive: segment.avg_spent > 0
                  }}
                  className={getSegmentColor(segment.segment)}
                />
              ))}
            </div>
          ) : (
            <BiHubEmptyState
              icon={<BarChart3 className="h-8 w-8 text-gray-400" />}
              title="No segment data available"
              description="Customer segmentation data will appear here"
            />
          )}
        </CardContent>
      </Card>
    </div>
  )

  const renderAnalyticsTab = () => (
    <div className="space-y-6">
      {analyticsLoading ? (
        <div className="text-center py-8">Loading analytics...</div>
      ) : analytics ? (
        <>
          {/* Overview Stats */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <BiHubStatCard
              title="Total Customers"
              value={analytics.overview.total_customers.toString()}
              icon={<Users className="h-6 w-6" />}
              trend={{
                value: analytics.overview.new_customers,
                label: `${analytics.overview.new_customers} new`,
                isPositive: true
              }}
            />
            <BiHubStatCard
              title="Active Customers"
              value={analytics.overview.active_customers.toString()}
              icon={<Users className="h-6 w-6" />}
              trend={{
                value: analytics.overview.churn_rate,
                label: `${analytics.overview.churn_rate.toFixed(1)}% churn`,
                isPositive: analytics.overview.churn_rate < 10
              }}
            />
            <BiHubStatCard
              title="Avg Lifetime Value"
              value={formatCurrency(analytics.overview.avg_lifetime_value)}
              icon={<DollarSign className="h-6 w-6" />}
              trend={{
                value: analytics.overview.avg_order_value,
                label: `${formatCurrency(analytics.overview.avg_order_value)} AOV`,
                isPositive: true
              }}
            />
            <BiHubStatCard
              title="Repeat Purchase Rate"
              value={`${analytics.retention_metrics.repeat_purchase_rate.toFixed(1)}%`}
              icon={<TrendingUp className="h-6 w-6" />}
              trend={{
                value: analytics.retention_metrics.day_30_retention,
                label: `${analytics.retention_metrics.day_30_retention.toFixed(1)}% 30-day retention`,
                isPositive: analytics.retention_metrics.day_30_retention > 50
              }}
            />
          </div>

          {/* Tier Distribution */}
          <Card>
            <CardHeader>
              <CardTitle>Membership Tier Distribution</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-5 gap-4">
                {analytics.tier_distribution.map((tier) => (
                  <div key={tier.tier} className="text-center p-4 border rounded-lg">
                    <div className={cn("inline-flex items-center justify-center w-12 h-12 rounded-full mb-2", getTierColor(tier.tier))}>
                      <Crown className="h-6 w-6" />
                    </div>
                    <h3 className="font-semibold capitalize">{tier.tier}</h3>
                    <p className="text-2xl font-bold">{tier.count}</p>
                    <p className="text-sm text-gray-600">{tier.percentage.toFixed(1)}%</p>
                    <p className="text-sm text-gray-600">{formatCurrency(tier.revenue)} revenue</p>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </>
      ) : (
        <BiHubEmptyState
          icon={<TrendingUp className="h-8 w-8 text-gray-400" />}
          title="No analytics data available"
          description="Customer analytics will appear here"
        />
      )}
    </div>
  )

  const renderHighValueTab = () => (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Crown className="h-5 w-5" />
              High Value Customers
              {highValueCustomers && (
                <Badge variant="secondary">
                  {highValueCustomers.total} customers
                </Badge>
              )}
            </div>
            <Button onClick={handleExportCustomers} variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent>
          {highValueLoading ? (
            <div className="text-center py-8">Loading high value customers...</div>
          ) : highValueCustomers?.customers.length === 0 ? (
            <BiHubEmptyState
              icon={<Crown className="h-8 w-8 text-gray-400" />}
              title="No high value customers found"
              description="High value customers will appear here"
            />
          ) : (
            <>
              {/* Criteria Info */}
              {highValueCustomers?.criteria && (
                <div className="mb-6 p-4 bg-blue-50 rounded-lg">
                  <h4 className="font-semibold mb-2">High Value Criteria:</h4>
                  <div className="flex gap-4 text-sm">
                    <span>Min Total Spent: {formatCurrency(highValueCustomers.criteria.min_total_spent)}</span>
                    <span>Min Total Orders: {highValueCustomers.criteria.min_total_orders}</span>
                  </div>
                </div>
              )}

              <div className="space-y-4">
                {highValueCustomers?.customers.map((customer) => (
                  <div
                    key={customer.id}
                    className="border rounded-lg p-4 hover:bg-gray-50 transition-colors"
                  >
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="flex items-center gap-3 mb-2">
                          <div className="flex items-center gap-2">
                            <Crown className="h-5 w-5 text-yellow-600" />
                            <h3 className="font-semibold">
                              {customer.first_name} {customer.last_name}
                            </h3>
                          </div>
                          <Badge className={getSegmentColor(customer.customer_segment)}>
                            {customer.customer_segment}
                          </Badge>
                          <Badge className={getTierColor(customer.membership_tier)}>
                            {customer.membership_tier}
                          </Badge>
                          {customer.is_vip && (
                            <Badge className="bg-purple-100 text-purple-800">
                              VIP
                            </Badge>
                          )}
                        </div>

                        <div className="grid grid-cols-2 md:grid-cols-5 gap-4 text-sm text-gray-600">
                          <div className="flex items-center gap-1">
                            <Mail className="h-4 w-4" />
                            {customer.email}
                          </div>
                          <div className="flex items-center gap-1">
                            <ShoppingBag className="h-4 w-4" />
                            {customer.order_count} orders
                          </div>
                          <div className="flex items-center gap-1">
                            <DollarSign className="h-4 w-4" />
                            {formatCurrency(customer.total_spent)}
                          </div>
                          <div className="flex items-center gap-1">
                            <Star className="h-4 w-4" />
                            {customer.loyalty_points} points
                          </div>
                          <div className="flex items-center gap-1">
                            <Calendar className="h-4 w-4" />
                            {formatDate(customer.created_at)}
                          </div>
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <BiHubStatusBadge
                          status={customer.is_active ? 'active' : 'inactive'}
                          variant={customer.is_active ? 'success' : 'secondary'}
                        />
                        <Button variant="ghost" size="sm">
                          <Eye className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="Customer Search & Segmentation"
        subtitle="Advanced customer search, segmentation, and analytics"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Customers' }
        ]}
      />

      {/* Tab Navigation */}
      <div className="flex space-x-1 mb-6">
        {[
          { id: 'search', label: 'Search', icon: Search },
          { id: 'segments', label: 'Segments', icon: BarChart3 },
          { id: 'analytics', label: 'Analytics', icon: TrendingUp },
          { id: 'high-value', label: 'High Value', icon: Crown },
        ].map((tab) => (
          <Button
            key={tab.id}
            variant={activeTab === tab.id ? 'default' : 'ghost'}
            onClick={() => setActiveTab(tab.id as any)}
            className="flex items-center gap-2"
          >
            <tab.icon className="h-4 w-4" />
            {tab.label}
          </Button>
        ))}
      </div>

      {/* Tab Content */}
      {activeTab === 'search' && renderSearchTab()}
      {activeTab === 'segments' && renderSegmentsTab()}
      {activeTab === 'analytics' && renderAnalyticsTab()}
      {activeTab === 'high-value' && renderHighValueTab()}
    </div>
  )
}
