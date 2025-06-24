'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { 
  User, 
  Package, 
  Heart, 
  Settings, 
  MapPin, 
  CreditCard,
  Bell,
  Shield,
  LogOut
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { useAuthStore } from '@/store/auth'
import { useRequireAuth } from '@/hooks/use-auth-guard'
import { formatPrice } from '@/lib/utils'
import { cn } from '@/lib/utils'

export function AccountPage() {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('overview')
  const { user, logout } = useAuthStore()
  const { isLoading, canAccess } = useRequireAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      </div>
    )
  }

  if (!canAccess || !user) {
    return null // useRequireAuth will handle redirect
  }

  const handleLogout = () => {
    logout()
    router.push('/')
  }

  const tabs = [
    { id: 'overview', label: 'Overview', icon: User },
    { id: 'orders', label: 'Orders', icon: Package },
    { id: 'wishlist', label: 'Wishlist', icon: Heart },
    { id: 'addresses', label: 'Addresses', icon: MapPin },
    { id: 'payment', label: 'Payment Methods', icon: CreditCard },
    { id: 'notifications', label: 'Notifications', icon: Bell },
    { id: 'security', label: 'Security', icon: Shield },
    { id: 'settings', label: 'Settings', icon: Settings },
  ]

  // Mock data - would come from API
  const recentOrders = [
    {
      id: '1',
      order_number: 'ORD-001',
      date: '2024-01-15',
      status: 'delivered',
      total: 129.99,
      items: 3,
    },
    {
      id: '2',
      order_number: 'ORD-002',
      date: '2024-01-10',
      status: 'shipped',
      total: 89.50,
      items: 2,
    },
    {
      id: '3',
      order_number: 'ORD-003',
      date: '2024-01-05',
      status: 'processing',
      total: 199.99,
      items: 1,
    },
  ]

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'delivered':
        return 'success'
      case 'shipped':
        return 'default'
      case 'processing':
        return 'warning'
      case 'cancelled':
        return 'destructive'
      default:
        return 'secondary'
    }
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">My Account</h1>
          <p className="text-gray-600 mt-2">
            Welcome back, {user.first_name}! Manage your account and view your orders.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <Card>
              <CardHeader className="text-center">
                <div className="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <User className="h-10 w-10 text-primary-600" />
                </div>
                <CardTitle>{user.first_name} {user.last_name}</CardTitle>
                <p className="text-sm text-gray-600">{user.email}</p>
              </CardHeader>
              <CardContent className="p-0">
                <nav className="space-y-1">
                  {tabs.map((tab) => {
                    const Icon = tab.icon
                    return (
                      <button
                        key={tab.id}
                        onClick={() => setActiveTab(tab.id)}
                        className={cn(
                          'w-full flex items-center space-x-3 px-6 py-3 text-left text-sm font-medium transition-colors',
                          activeTab === tab.id
                            ? 'bg-primary-50 text-primary-700 border-r-2 border-primary-600'
                            : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                        )}
                      >
                        <Icon className="h-5 w-5" />
                        <span>{tab.label}</span>
                      </button>
                    )
                  })}
                  
                  <button
                    onClick={handleLogout}
                    className="w-full flex items-center space-x-3 px-6 py-3 text-left text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
                  >
                    <LogOut className="h-5 w-5" />
                    <span>Sign Out</span>
                  </button>
                </nav>
              </CardContent>
            </Card>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-3">
            {activeTab === 'overview' && (
              <div className="space-y-6">
                {/* Stats Cards */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <Card>
                    <CardContent className="p-6">
                      <div className="flex items-center">
                        <Package className="h-8 w-8 text-primary-600" />
                        <div className="ml-4">
                          <p className="text-sm font-medium text-gray-600">Total Orders</p>
                          <p className="text-2xl font-bold text-gray-900">12</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                  
                  <Card>
                    <CardContent className="p-6">
                      <div className="flex items-center">
                        <Heart className="h-8 w-8 text-red-500" />
                        <div className="ml-4">
                          <p className="text-sm font-medium text-gray-600">Wishlist Items</p>
                          <p className="text-2xl font-bold text-gray-900">8</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                  
                  <Card>
                    <CardContent className="p-6">
                      <div className="flex items-center">
                        <CreditCard className="h-8 w-8 text-green-600" />
                        <div className="ml-4">
                          <p className="text-sm font-medium text-gray-600">Total Spent</p>
                          <p className="text-2xl font-bold text-gray-900">{formatPrice(1249.99)}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </div>

                {/* Recent Orders */}
                <Card>
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle>Recent Orders</CardTitle>
                      <Button 
                        variant="outline" 
                        size="sm"
                        onClick={() => setActiveTab('orders')}
                      >
                        View All
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent>
                    <div className="space-y-4">
                      {recentOrders.map((order) => (
                        <div key={order.id} className="flex items-center justify-between p-4 border rounded-lg">
                          <div className="flex items-center space-x-4">
                            <div>
                              <p className="font-medium">{order.order_number}</p>
                              <p className="text-sm text-gray-600">{order.date}</p>
                            </div>
                            <Badge variant={getStatusColor(order.status) as any}>
                              {order.status}
                            </Badge>
                          </div>
                          <div className="text-right">
                            <p className="font-medium">{formatPrice(order.total)}</p>
                            <p className="text-sm text-gray-600">{order.items} items</p>
                          </div>
                        </div>
                      ))}
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {activeTab === 'orders' && (
              <Card>
                <CardHeader>
                  <CardTitle>Order History</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    {recentOrders.map((order) => (
                      <div key={order.id} className="border rounded-lg p-6">
                        <div className="flex items-center justify-between mb-4">
                          <div>
                            <h3 className="font-semibold">{order.order_number}</h3>
                            <p className="text-sm text-gray-600">Placed on {order.date}</p>
                          </div>
                          <Badge variant={getStatusColor(order.status) as any}>
                            {order.status}
                          </Badge>
                        </div>
                        <div className="flex items-center justify-between">
                          <p className="text-sm text-gray-600">{order.items} items</p>
                          <div className="flex items-center space-x-4">
                            <p className="font-semibold">{formatPrice(order.total)}</p>
                            <Button variant="outline" size="sm">
                              View Details
                            </Button>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            )}

            {activeTab === 'wishlist' && (
              <Card>
                <CardHeader>
                  <CardTitle>My Wishlist</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-center py-12">
                    <Heart className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">
                      Your wishlist is empty
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Save items you love to your wishlist and shop them later.
                    </p>
                    <Button onClick={() => router.push('/products')}>
                      Browse Products
                    </Button>
                  </div>
                </CardContent>
              </Card>
            )}

            {activeTab === 'addresses' && (
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle>Saved Addresses</CardTitle>
                    <Button>Add New Address</Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-center py-12">
                    <MapPin className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">
                      No addresses saved
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Add your addresses to make checkout faster.
                    </p>
                    <Button>Add Your First Address</Button>
                  </div>
                </CardContent>
              </Card>
            )}

            {activeTab === 'payment' && (
              <Card>
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <CardTitle>Payment Methods</CardTitle>
                    <Button>Add Payment Method</Button>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="text-center py-12">
                    <CreditCard className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                    <h3 className="text-lg font-semibold text-gray-900 mb-2">
                      No payment methods saved
                    </h3>
                    <p className="text-gray-600 mb-6">
                      Add your payment methods for faster checkout.
                    </p>
                    <Button>Add Payment Method</Button>
                  </div>
                </CardContent>
              </Card>
            )}

            {activeTab === 'notifications' && (
              <Card>
                <CardHeader>
                  <CardTitle>Notification Preferences</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">Order Updates</h4>
                        <p className="text-sm text-gray-600">Get notified about your order status</p>
                      </div>
                      <input type="checkbox" defaultChecked className="rounded" />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">Promotions</h4>
                        <p className="text-sm text-gray-600">Receive emails about sales and promotions</p>
                      </div>
                      <input type="checkbox" defaultChecked className="rounded" />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">New Products</h4>
                        <p className="text-sm text-gray-600">Be the first to know about new arrivals</p>
                      </div>
                      <input type="checkbox" className="rounded" />
                    </div>
                  </div>
                  
                  <Button>Save Preferences</Button>
                </CardContent>
              </Card>
            )}

            {activeTab === 'security' && (
              <Card>
                <CardHeader>
                  <CardTitle>Security Settings</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">Change Password</h4>
                        <p className="text-sm text-gray-600">Update your account password</p>
                      </div>
                      <Button variant="outline">Change</Button>
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">Two-Factor Authentication</h4>
                        <p className="text-sm text-gray-600">Add an extra layer of security</p>
                      </div>
                      <Button variant="outline">Enable</Button>
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="font-medium">Login Sessions</h4>
                        <p className="text-sm text-gray-600">Manage your active sessions</p>
                      </div>
                      <Button variant="outline">View</Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            )}

            {activeTab === 'settings' && (
              <Card>
                <CardHeader>
                  <CardTitle>Account Settings</CardTitle>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-4">
                    <div>
                      <h4 className="font-medium mb-2">Personal Information</h4>
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            First Name
                          </label>
                          <input
                            type="text"
                            defaultValue={user.first_name}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-700 mb-1">
                            Last Name
                          </label>
                          <input
                            type="text"
                            defaultValue={user.last_name}
                            className="w-full px-3 py-2 border border-gray-300 rounded-md"
                          />
                        </div>
                      </div>
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Email Address
                      </label>
                      <input
                        type="email"
                        defaultValue={user.email}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md"
                      />
                    </div>
                    
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Phone Number
                      </label>
                      <input
                        type="tel"
                        defaultValue={user.phone || ''}
                        className="w-full px-3 py-2 border border-gray-300 rounded-md"
                      />
                    </div>
                  </div>
                  
                  <Button>Save Changes</Button>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
