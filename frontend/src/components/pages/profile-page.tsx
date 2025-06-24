'use client'

import { useState } from 'react'
import Image from 'next/image'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  User,
  Mail,
  Phone,
  MapPin,
  Calendar,
  Edit,
  Save,
  X,
  Camera,
  Shield,
  ShoppingBag,
  Heart,
  Settings,
  CreditCard,
  Truck,
  Bell,
  Lock,
} from 'lucide-react'
import { useProfile, useUpdateProfile, useChangePassword } from '@/hooks/use-users'
import { useOrders } from '@/hooks/use-orders'
import { formatDate, formatPrice } from '@/lib/utils'
import { toast } from 'sonner'

export default function ProfilePage() {
  const [isEditing, setIsEditing] = useState(false)
  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [profileData, setProfileData] = useState({
    first_name: '',
    last_name: '',
    phone: '',
    profile: {
      date_of_birth: '',
      gender: '',
      address: '',
      city: '',
      country: '',
      postal_code: '',
    }
  })
  const [passwordData, setPasswordData] = useState({
    current_password: '',
    new_password: '',
    confirm_password: '',
  })

  const { data: user, isLoading: userLoading } = useProfile()
  const { data: ordersData, isLoading: ordersLoading } = useOrders({ limit: 5 })
  const updateProfile = useUpdateProfile()
  const changePassword = useChangePassword()

  const orders = ordersData?.data || []

  // Initialize form data when user data loads
  React.useEffect(() => {
    if (user) {
      setProfileData({
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        phone: user.phone || '',
        profile: {
          date_of_birth: user.profile?.date_of_birth || '',
          gender: user.profile?.gender || '',
          address: user.profile?.address || '',
          city: user.profile?.city || '',
          country: user.profile?.country || '',
          postal_code: user.profile?.postal_code || '',
        }
      })
    }
  }, [user])

  const handleProfileUpdate = async () => {
    try {
      await updateProfile.mutateAsync(profileData)
      setIsEditing(false)
    } catch (error) {
      console.error('Failed to update profile:', error)
    }
  }

  const handlePasswordChange = async () => {
    if (passwordData.new_password !== passwordData.confirm_password) {
      toast.error('New passwords do not match')
      return
    }

    try {
      await changePassword.mutateAsync(passwordData)
      setPasswordData({
        current_password: '',
        new_password: '',
        confirm_password: '',
      })
      setShowPasswordForm(false)
    } catch (error) {
      console.error('Failed to change password:', error)
    }
  }

  const getOrderStatusVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'default'
      case 'shipped':
        return 'default'
      case 'processing':
        return 'secondary'
      case 'pending':
        return 'secondary'
      case 'cancelled':
        return 'destructive'
      default:
        return 'outline'
    }
  }

  if (userLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-12">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <div className="lg:col-span-1">
              <Card variant="elevated" className="border-0 shadow-large">
                <CardContent className="p-8">
                  <div className="animate-pulse">
                    <div className="w-32 h-32 bg-muted rounded-full mx-auto mb-6"></div>
                    <div className="space-y-3">
                      <div className="h-6 bg-muted rounded w-3/4 mx-auto"></div>
                      <div className="h-4 bg-muted rounded w-1/2 mx-auto"></div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            </div>
            <div className="lg:col-span-2">
              <Card variant="elevated" className="border-0 shadow-large">
                <CardContent className="p-8">
                  <div className="animate-pulse space-y-6">
                    {[...Array(4)].map((_, i) => (
                      <div key={i} className="space-y-3">
                        <div className="h-4 bg-muted rounded w-1/4"></div>
                        <div className="h-10 bg-muted rounded"></div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-12">
        <div className="container mx-auto px-4 text-center">
          <h1 className="text-3xl font-bold text-destructive mb-4">Access Denied</h1>
          <p className="text-muted-foreground mb-8">
            You need to be logged in to view your profile.
          </p>
          <Button asChild>
            <a href="/auth/login">Login</a>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-12">
      <div className="container mx-auto px-4">
        {/* Header */}
        <div className="mb-12">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large">
              <User className="h-6 w-6 text-white" />
            </div>
            <span className="text-primary font-semibold">MY PROFILE</span>
          </div>
          
          <h1 className="text-4xl lg:text-5xl font-bold text-foreground mb-4">
            Welcome back, <span className="text-gradient">{user.first_name}</span>
          </h1>
          <p className="text-xl text-muted-foreground">
            Manage your account settings and view your order history
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Profile Sidebar */}
          <div className="lg:col-span-1">
            <Card variant="elevated" className="border-0 shadow-large">
              <CardContent className="p-8 text-center">
                {/* Avatar */}
                <div className="relative mb-6">
                  <div className="w-32 h-32 rounded-full overflow-hidden mx-auto bg-muted flex items-center justify-center">
                    {user.avatar ? (
                      <Image
                        src={user.avatar}
                        alt={`${user.first_name} ${user.last_name}`}
                        width={128}
                        height={128}
                        className="object-cover"
                      />
                    ) : (
                      <User className="w-16 h-16 text-muted-foreground" />
                    )}
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    className="absolute bottom-0 right-0 rounded-full w-10 h-10 p-0"
                  >
                    <Camera className="h-4 w-4" />
                  </Button>
                </div>

                {/* User Info */}
                <h2 className="text-2xl font-bold text-foreground mb-2">
                  {user.first_name} {user.last_name}
                </h2>
                <p className="text-muted-foreground mb-4">{user.email}</p>
                
                <div className="flex items-center justify-center gap-2 mb-6">
                  <Badge variant="default" className="font-semibold">
                    <Shield className="h-3 w-3 mr-1" />
                    {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                  </Badge>
                  <Badge variant={user.is_active ? "default" : "secondary"}>
                    {user.is_active ? 'Active' : 'Inactive'}
                  </Badge>
                </div>

                {/* Quick Stats */}
                <div className="grid grid-cols-2 gap-4 text-center">
                  <div>
                    <div className="text-2xl font-bold text-primary">{orders.length}</div>
                    <div className="text-sm text-muted-foreground">Orders</div>
                  </div>
                  <div>
                    <div className="text-2xl font-bold text-emerald-600">
                      {formatDate(user.created_at, { month: 'short', year: 'numeric' })}
                    </div>
                    <div className="text-sm text-muted-foreground">Member Since</div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-2">
            <Tabs defaultValue="profile" className="space-y-8">
              <TabsList className="grid w-full grid-cols-4">
                <TabsTrigger value="profile" className="flex items-center gap-2">
                  <User className="h-4 w-4" />
                  Profile
                </TabsTrigger>
                <TabsTrigger value="orders" className="flex items-center gap-2">
                  <ShoppingBag className="h-4 w-4" />
                  Orders
                </TabsTrigger>
                <TabsTrigger value="security" className="flex items-center gap-2">
                  <Lock className="h-4 w-4" />
                  Security
                </TabsTrigger>
                <TabsTrigger value="settings" className="flex items-center gap-2">
                  <Settings className="h-4 w-4" />
                  Settings
                </TabsTrigger>
              </TabsList>

              {/* Profile Tab */}
              <TabsContent value="profile">
                <Card variant="elevated" className="border-0 shadow-large">
                  <CardHeader>
                    <div className="flex items-center justify-between">
                      <CardTitle className="text-2xl">Personal Information</CardTitle>
                      <Button
                        variant={isEditing ? "outline" : "default"}
                        onClick={() => {
                          if (isEditing) {
                            setIsEditing(false)
                            // Reset form data
                            setProfileData({
                              first_name: user.first_name || '',
                              last_name: user.last_name || '',
                              phone: user.phone || '',
                              profile: {
                                date_of_birth: user.profile?.date_of_birth || '',
                                gender: user.profile?.gender || '',
                                address: user.profile?.address || '',
                                city: user.profile?.city || '',
                                country: user.profile?.country || '',
                                postal_code: user.profile?.postal_code || '',
                              }
                            })
                          } else {
                            setIsEditing(true)
                          }
                        }}
                      >
                        {isEditing ? (
                          <>
                            <X className="h-4 w-4 mr-2" />
                            Cancel
                          </>
                        ) : (
                          <>
                            <Edit className="h-4 w-4 mr-2" />
                            Edit
                          </>
                        )}
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent className="p-8 space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          First Name
                        </label>
                        <Input
                          value={profileData.first_name}
                          onChange={(e) => setProfileData(prev => ({ ...prev, first_name: e.target.value }))}
                          disabled={!isEditing}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Last Name
                        </label>
                        <Input
                          value={profileData.last_name}
                          onChange={(e) => setProfileData(prev => ({ ...prev, last_name: e.target.value }))}
                          disabled={!isEditing}
                        />
                      </div>
                    </div>

                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Email Address
                      </label>
                      <Input
                        value={user.email}
                        disabled
                        leftIcon={<Mail className="h-4 w-4" />}
                      />
                      <p className="text-xs text-muted-foreground mt-1">
                        Email cannot be changed. Contact support if needed.
                      </p>
                    </div>

                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Phone Number
                      </label>
                      <Input
                        value={profileData.phone}
                        onChange={(e) => setProfileData(prev => ({ ...prev, phone: e.target.value }))}
                        disabled={!isEditing}
                        leftIcon={<Phone className="h-4 w-4" />}
                      />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Date of Birth
                        </label>
                        <Input
                          type="date"
                          value={profileData.profile.date_of_birth}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, date_of_birth: e.target.value }
                          }))}
                          disabled={!isEditing}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Gender
                        </label>
                        <Select
                          value={profileData.profile.gender}
                          onValueChange={(value) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, gender: value }
                          }))}
                          disabled={!isEditing}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Select gender" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="male">Male</SelectItem>
                            <SelectItem value="female">Female</SelectItem>
                            <SelectItem value="other">Other</SelectItem>
                            <SelectItem value="prefer_not_to_say">Prefer not to say</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <div>
                      <label className="block text-sm font-semibold text-foreground mb-2">
                        Address
                      </label>
                      <Input
                        value={profileData.profile.address}
                        onChange={(e) => setProfileData(prev => ({
                          ...prev,
                          profile: { ...prev.profile, address: e.target.value }
                        }))}
                        disabled={!isEditing}
                        leftIcon={<MapPin className="h-4 w-4" />}
                      />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          City
                        </label>
                        <Input
                          value={profileData.profile.city}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, city: e.target.value }
                          }))}
                          disabled={!isEditing}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Country
                        </label>
                        <Input
                          value={profileData.profile.country}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, country: e.target.value }
                          }))}
                          disabled={!isEditing}
                        />
                      </div>
                      
                      <div>
                        <label className="block text-sm font-semibold text-foreground mb-2">
                          Postal Code
                        </label>
                        <Input
                          value={profileData.profile.postal_code}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, postal_code: e.target.value }
                          }))}
                          disabled={!isEditing}
                        />
                      </div>
                    </div>

                    {isEditing && (
                      <div className="flex items-center gap-4 pt-6 border-t">
                        <Button
                          onClick={handleProfileUpdate}
                          disabled={updateProfile.isPending}
                          variant="gradient"
                        >
                          {updateProfile.isPending ? (
                            <>
                              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                              Saving...
                            </>
                          ) : (
                            <>
                              <Save className="h-4 w-4 mr-2" />
                              Save Changes
                            </>
                          )}
                        </Button>
                        <Button
                          variant="outline"
                          onClick={() => setIsEditing(false)}
                        >
                          Cancel
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Orders Tab */}
              <TabsContent value="orders">
                <Card variant="elevated" className="border-0 shadow-large">
                  <CardHeader>
                    <CardTitle className="text-2xl">Recent Orders</CardTitle>
                  </CardHeader>
                  <CardContent className="p-8">
                    {ordersLoading ? (
                      <div className="space-y-4">
                        {[...Array(3)].map((_, i) => (
                          <div key={i} className="animate-pulse">
                            <div className="h-20 bg-muted rounded-xl"></div>
                          </div>
                        ))}
                      </div>
                    ) : orders.length > 0 ? (
                      <div className="space-y-6">
                        {orders.map((order) => (
                          <div key={order.id} className="border border-border rounded-xl p-6 hover:shadow-medium transition-shadow">
                            <div className="flex items-center justify-between mb-4">
                              <div>
                                <h3 className="font-semibold text-foreground">
                                  Order #{order.order_number || order.id.slice(0, 8)}
                                </h3>
                                <p className="text-sm text-muted-foreground">
                                  {formatDate(order.created_at)}
                                </p>
                              </div>
                              <Badge variant={getOrderStatusVariant(order.status)}>
                                {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                              </Badge>
                            </div>
                            
                            <div className="flex items-center justify-between">
                              <div className="flex items-center gap-4">
                                <div className="flex items-center gap-2">
                                  <ShoppingBag className="h-4 w-4 text-muted-foreground" />
                                  <span className="text-sm">{order.items?.length || 0} items</span>
                                </div>
                                <div className="flex items-center gap-2">
                                  <span className="text-lg font-bold text-primary">
                                    {formatPrice(order.total_amount)}
                                  </span>
                                </div>
                              </div>
                              
                              <Button variant="outline" size="sm">
                                View Details
                              </Button>
                            </div>
                          </div>
                        ))}
                        
                        <div className="text-center pt-6">
                          <Button variant="outline">
                            View All Orders
                          </Button>
                        </div>
                      </div>
                    ) : (
                      <div className="text-center py-12">
                        <ShoppingBag className="h-16 w-16 text-muted-foreground mx-auto mb-4" />
                        <h3 className="text-lg font-semibold text-foreground mb-2">No orders yet</h3>
                        <p className="text-muted-foreground mb-6">
                          Start shopping to see your orders here.
                        </p>
                        <Button asChild>
                          <a href="/products">Browse Products</a>
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Security Tab */}
              <TabsContent value="security">
                <Card variant="elevated" className="border-0 shadow-large">
                  <CardHeader>
                    <CardTitle className="text-2xl">Security Settings</CardTitle>
                  </CardHeader>
                  <CardContent className="p-8 space-y-6">
                    <div className="border border-border rounded-xl p-6">
                      <div className="flex items-center justify-between mb-4">
                        <div>
                          <h3 className="font-semibold text-foreground">Password</h3>
                          <p className="text-sm text-muted-foreground">
                            Last changed {formatDate(user.updated_at)}
                          </p>
                        </div>
                        <Button
                          variant="outline"
                          onClick={() => setShowPasswordForm(!showPasswordForm)}
                        >
                          Change Password
                        </Button>
                      </div>
                      
                      {showPasswordForm && (
                        <div className="space-y-4 pt-4 border-t">
                          <div>
                            <label className="block text-sm font-semibold text-foreground mb-2">
                              Current Password
                            </label>
                            <Input
                              type="password"
                              value={passwordData.current_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, current_password: e.target.value }))}
                            />
                          </div>
                          
                          <div>
                            <label className="block text-sm font-semibold text-foreground mb-2">
                              New Password
                            </label>
                            <Input
                              type="password"
                              value={passwordData.new_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, new_password: e.target.value }))}
                            />
                          </div>
                          
                          <div>
                            <label className="block text-sm font-semibold text-foreground mb-2">
                              Confirm New Password
                            </label>
                            <Input
                              type="password"
                              value={passwordData.confirm_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, confirm_password: e.target.value }))}
                            />
                          </div>
                          
                          <div className="flex items-center gap-4">
                            <Button
                              onClick={handlePasswordChange}
                              disabled={changePassword.isPending}
                              variant="gradient"
                            >
                              {changePassword.isPending ? 'Changing...' : 'Change Password'}
                            </Button>
                            <Button
                              variant="outline"
                              onClick={() => setShowPasswordForm(false)}
                            >
                              Cancel
                            </Button>
                          </div>
                        </div>
                      )}
                    </div>

                    <div className="border border-border rounded-xl p-6">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-foreground">Two-Factor Authentication</h3>
                          <p className="text-sm text-muted-foreground">
                            Add an extra layer of security to your account
                          </p>
                        </div>
                        <Button variant="outline">
                          Enable 2FA
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Settings Tab */}
              <TabsContent value="settings">
                <Card variant="elevated" className="border-0 shadow-large">
                  <CardHeader>
                    <CardTitle className="text-2xl">Account Settings</CardTitle>
                  </CardHeader>
                  <CardContent className="p-8 space-y-6">
                    <div className="border border-border rounded-xl p-6">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-foreground">Email Notifications</h3>
                          <p className="text-sm text-muted-foreground">
                            Receive updates about your orders and account
                          </p>
                        </div>
                        <Button variant="outline">
                          <Bell className="h-4 w-4 mr-2" />
                          Manage
                        </Button>
                      </div>
                    </div>

                    <div className="border border-border rounded-xl p-6">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-foreground">Payment Methods</h3>
                          <p className="text-sm text-muted-foreground">
                            Manage your saved payment methods
                          </p>
                        </div>
                        <Button variant="outline">
                          <CreditCard className="h-4 w-4 mr-2" />
                          Manage
                        </Button>
                      </div>
                    </div>

                    <div className="border border-border rounded-xl p-6">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-foreground">Shipping Addresses</h3>
                          <p className="text-sm text-muted-foreground">
                            Manage your delivery addresses
                          </p>
                        </div>
                        <Button variant="outline">
                          <Truck className="h-4 w-4 mr-2" />
                          Manage
                        </Button>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </div>
    </div>
  )
}
