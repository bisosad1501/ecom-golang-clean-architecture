'use client'

import React, { useState, useEffect } from 'react'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import {
  User,
  Mail,
  Phone,
  Calendar,
  Edit,
  Save,
  Shield,
  X,
  Camera,
  ShoppingBag,
  Lock,
  Package,
  Clock,
  CheckCircle,
  Upload,
  Globe,
  Award,
  Activity,
  DollarSign,
  Truck,
} from 'lucide-react'
import { useProfile, useUpdateProfile, useChangePassword } from '@/hooks/use-users'
import { useOrders } from '@/hooks/use-orders'
import { formatDate, formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubStatCard,
} from '@/components/admin/bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'

export default function ProfilePage() {
  const [isEditing, setIsEditing] = useState(false)
  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [showAvatarUpload, setShowAvatarUpload] = useState(false)
  const [activeTab, setActiveTab] = useState('overview')

  const [profileData, setProfileData] = useState({
    first_name: '',
    last_name: '',
    profile: {
      phone: '',
      date_of_birth: '',
      gender: '',
      bio: '',
      website: '',
      avatar_url: '',
    }
  })

  const [passwordData, setPasswordData] = useState({
    current_password: '',
    new_password: '',
    confirm_password: '',
  })

  const { data: user, isLoading: userLoading } = useProfile()
  const { data: ordersData, isLoading: ordersLoading } = useOrders({ limit: 10 })
  const updateProfile = useUpdateProfile()
  const changePassword = useChangePassword()

  const orders = ordersData?.data || []

  // Use backend user metrics as single source of truth (matches UserMetricsService)
  const userStats = {
    totalOrders: user?.total_orders || 0,
    completedOrders: orders.filter(order => order.status === 'delivered').length,
    totalSpent: user?.total_spent || 0,
    loyaltyPoints: user?.loyalty_points || 0,
    membershipTier: user?.membership_tier || 'bronze',
    memberSince: user?.created_at ? new Date(user.created_at).getFullYear() : new Date().getFullYear(),
  }

  // Initialize form data when user data loads
  useEffect(() => {
    if (user) {
      setProfileData({
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        profile: {
          phone: user.profile?.phone || '',
          date_of_birth: user.profile?.date_of_birth || '',
          gender: user.profile?.gender || '',
          bio: user.profile?.bio || '',
          website: user.profile?.website || '',
          avatar_url: user.profile?.avatar_url || '',
        }
      })
    }
  }, [user])

  const handleProfileUpdate = async () => {
    try {
      // Only send data that backend supports
      const updateData = {
        first_name: profileData.first_name,
        last_name: profileData.last_name,
        phone: profileData.profile.phone,
      }

      await updateProfile.mutateAsync(updateData)
      setIsEditing(false)
      toast.success('Profile updated successfully!')
    } catch (error) {
      console.error('Failed to update profile:', error)
      toast.error('Failed to update profile')
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
      toast.success('Password changed successfully!')
    } catch (error) {
      console.error('Failed to change password:', error)
      toast.error('Failed to change password')
    }
  }

  const getOrderStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'success'
      case 'shipped':
        return 'info'
      case 'processing':
        return 'warning'
      case 'pending':
        return 'error'
      case 'cancelled':
        return 'error'
      default:
        return 'default'
    }
  }

  const getOrderStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return <CheckCircle className="h-4 w-4" />
      case 'shipped':
        return <Truck className="h-4 w-4" />
      case 'processing':
        return <Package className="h-4 w-4" />
      case 'pending':
        return <Clock className="h-4 w-4" />
      case 'cancelled':
        return <X className="h-4 w-4" />
      default:
        return <ShoppingBag className="h-4 w-4" />
    }
  }

  if (userLoading) {
    return (
      <div className="min-h-screen bg-gray-950 py-8">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
            {/* Profile Card Skeleton */}
            <div className="lg:col-span-1">
              <div className={cn(BIHUB_ADMIN_THEME.components.card.base, 'p-6 animate-pulse')}>
                <div className="w-24 h-24 bg-gray-700 rounded-full mx-auto mb-4"></div>
                <div className="space-y-2">
                  <div className="h-6 bg-gray-700 rounded w-3/4 mx-auto"></div>
                  <div className="h-4 bg-gray-700 rounded w-1/2 mx-auto"></div>
                </div>
              </div>
            </div>

            {/* Content Skeleton */}
            <div className="lg:col-span-3">
              <div className={cn(BIHUB_ADMIN_THEME.components.card.base, 'p-6 animate-pulse')}>
                <div className="space-y-4">
                  {[...Array(4)].map((_, i) => (
                    <div key={i} className="space-y-2">
                      <div className="h-4 bg-gray-700 rounded w-1/4"></div>
                      <div className="h-10 bg-gray-700 rounded"></div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-gray-950 flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <Lock className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white mb-4">Access Denied</h1>
          <p className="text-gray-400 mb-8">
            You need to be logged in to view your profile.
          </p>
          <Button className="bg-[#FF9000] hover:bg-[#e67e00] text-white">
            <a href="/auth/login">Login to BiHub</a>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-950 py-8">
      <div className="container mx-auto px-4">
        {/* BiHub Profile Header */}
        <div className="relative overflow-hidden bg-gradient-to-br from-[#FF9000] via-[#e67e00] to-[#cc6600] rounded-2xl shadow-2xl mb-8">
          <div className="absolute inset-0 bg-black/10"></div>
          <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full -translate-y-32 translate-x-32"></div>
          <div className="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full translate-y-24 -translate-x-24"></div>

          <div className="relative p-8">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
              <div className="text-white">
                <div className="flex items-center gap-3 mb-4">
                  <div className="w-12 h-12 rounded-2xl bg-white/20 backdrop-blur-sm flex items-center justify-center">
                    <User className="h-6 w-6 text-white" />
                  </div>
                  <span className="text-white/90 font-semibold">BIHUB PROFILE</span>
                </div>

                <h1 className="text-4xl lg:text-5xl font-bold mb-4">
                  Welcome back, <span className="text-white/90">{user.first_name}!</span>
                </h1>
                <p className="text-xl text-white/80 leading-relaxed">
                  Manage your BiHub account settings and view your order history
                </p>
              </div>

              <div className="flex flex-col sm:flex-row gap-4">
                <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6 text-center">
                  <p className="text-white/70 text-sm font-medium">Member Since</p>
                  <p className="text-white font-semibold text-lg">
                    {userStats.memberSince}
                  </p>
                </div>

                <div className="bg-white/10 backdrop-blur-sm rounded-2xl p-6 text-center">
                  <p className="text-white/70 text-sm font-medium">Total Orders</p>
                  <p className="text-white font-semibold text-lg">
                    {userStats.totalOrders}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-6 gap-6 mb-8">
          <BiHubStatCard
            title="Total Orders"
            value={userStats.totalOrders}
            icon={<ShoppingBag className="h-8 w-8 text-white" />}
            color="primary"
          />
          <BiHubStatCard
            title="Completed Orders"
            value={userStats.completedOrders}
            icon={<CheckCircle className="h-8 w-8 text-white" />}
            color="success"
          />
          <BiHubStatCard
            title="Total Spent"
            value={formatPrice(userStats.totalSpent)}
            icon={<DollarSign className="h-8 w-8 text-white" />}
            color="info"
          />
          <BiHubStatCard
            title="Loyalty Points"
            value={userStats.loyaltyPoints.toLocaleString()}
            icon={<Award className="h-8 w-8 text-white" />}
            color="warning"
          />
          <BiHubStatCard
            title="Membership Tier"
            value={userStats.membershipTier.charAt(0).toUpperCase() + userStats.membershipTier.slice(1)}
            icon={<Shield className="h-8 w-8 text-white" />}
            color="secondary"
          />
          <BiHubStatCard
            title="Member Since"
            value={userStats.memberSince}
            icon={<Calendar className="h-8 w-8 text-white" />}
            color="accent"
          />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Profile Sidebar */}
          <div className="lg:col-span-1">
            <BiHubAdminCard
              title="Profile"
              subtitle="Your BiHub account information"
              icon={<User className="h-5 w-5 text-white" />}
            >
              <div className="text-center">
                {/* Avatar */}
                <div className="relative mb-6">
                  <div className="w-24 h-24 rounded-full overflow-hidden mx-auto bg-gray-700 flex items-center justify-center">
                    {user.profile?.avatar_url ? (
                      <Image
                        src={user.profile.avatar_url}
                        alt={`${user.first_name} ${user.last_name}`}
                        width={96}
                        height={96}
                        className="object-cover"
                      />
                    ) : (
                      <User className="w-12 h-12 text-gray-400" />
                    )}
                  </div>

                  <Dialog open={showAvatarUpload} onOpenChange={setShowAvatarUpload}>
                    <DialogTrigger asChild>
                      <Button
                        size="sm"
                        variant="outline"
                        className="absolute bottom-0 right-0 rounded-full w-8 h-8 p-0 bg-[#FF9000] hover:bg-[#e67e00] border-[#FF9000] text-white"
                      >
                        <Camera className="h-4 w-4" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="bg-gray-900 border-gray-700">
                      <DialogHeader>
                        <DialogTitle className="text-white">Update Profile Picture</DialogTitle>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div className="flex items-center justify-center w-full">
                          <label className="flex flex-col items-center justify-center w-full h-32 border-2 border-gray-600 border-dashed rounded-lg cursor-pointer bg-gray-800 hover:bg-gray-700">
                            <div className="flex flex-col items-center justify-center pt-5 pb-6">
                              <Upload className="w-8 h-8 mb-4 text-gray-400" />
                              <p className="mb-2 text-sm text-gray-400">
                                <span className="font-semibold">Click to upload</span> or drag and drop
                              </p>
                              <p className="text-xs text-gray-400">PNG, JPG or GIF (MAX. 800x400px)</p>
                            </div>
                            <input type="file" className="hidden" accept="image/*" />
                          </label>
                        </div>
                      </div>
                    </DialogContent>
                  </Dialog>
                </div>

                {/* User Info */}
                <h2 className="text-xl font-bold text-white mb-2">
                  {user.first_name} {user.last_name}
                </h2>
                <p className="text-gray-400 mb-4">{user.email}</p>

                <div className="flex flex-col gap-2 mb-6">
                  <BiHubStatusBadge status="info">
                    <Shield className="h-3 w-3 mr-1" />
                    {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                  </BiHubStatusBadge>
                  <BiHubStatusBadge status={user.is_active ? "success" : "error"}>
                    {user.is_active ? 'Active Account' : 'Inactive Account'}
                  </BiHubStatusBadge>
                </div>

                {/* Quick Actions */}
                <div className="space-y-2">
                  <Button
                    onClick={() => setIsEditing(true)}
                    className="w-full bg-[#FF9000] hover:bg-[#e67e00] text-white"
                  >
                    <Edit className="h-4 w-4 mr-2" />
                    Edit Profile
                  </Button>
                  <Button
                    onClick={() => setShowPasswordForm(true)}
                    variant="outline"
                    className="w-full border-gray-600 text-gray-300 hover:text-white hover:bg-gray-800"
                  >
                    <Lock className="h-4 w-4 mr-2" />
                    Change Password
                  </Button>
                </div>
              </div>
            </BiHubAdminCard>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-3">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
              <TabsList className="grid w-full grid-cols-4 bg-gray-800 border-gray-700">
                <TabsTrigger
                  value="overview"
                  className="flex items-center gap-2 data-[state=active]:bg-[#FF9000] data-[state=active]:text-white"
                >
                  <User className="h-4 w-4" />
                  Overview
                </TabsTrigger>
                <TabsTrigger
                  value="orders"
                  className="flex items-center gap-2 data-[state=active]:bg-[#FF9000] data-[state=active]:text-white"
                >
                  <ShoppingBag className="h-4 w-4" />
                  Orders
                </TabsTrigger>
                <TabsTrigger
                  value="profile"
                  className="flex items-center gap-2 data-[state=active]:bg-[#FF9000] data-[state=active]:text-white"
                >
                  <Edit className="h-4 w-4" />
                  Edit Profile
                </TabsTrigger>
                <TabsTrigger
                  value="security"
                  className="flex items-center gap-2 data-[state=active]:bg-[#FF9000] data-[state=active]:text-white"
                >
                  <Lock className="h-4 w-4" />
                  Security
                </TabsTrigger>
              </TabsList>

              {/* Overview Tab */}
              <TabsContent value="overview">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <BiHubAdminCard
                    title="Account Information"
                    subtitle="Your basic account details"
                    icon={<User className="h-5 w-5 text-white" />}
                  >
                    <div className="space-y-4">
                      <div className="flex items-center gap-3">
                        <Mail className="h-5 w-5 text-gray-400" />
                        <div>
                          <p className="text-sm text-gray-400">Email</p>
                          <p className="text-white font-medium">{user.email}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-3">
                        <Phone className="h-5 w-5 text-gray-400" />
                        <div>
                          <p className="text-sm text-gray-400">Phone</p>
                          <p className="text-white font-medium">{user.profile?.phone || 'Not provided'}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-3">
                        <Calendar className="h-5 w-5 text-gray-400" />
                        <div>
                          <p className="text-sm text-gray-400">Member Since</p>
                          <p className="text-white font-medium">{formatDate(user.created_at)}</p>
                        </div>
                      </div>
                    </div>
                  </BiHubAdminCard>

                  <BiHubAdminCard
                    title="Recent Activity"
                    subtitle="Your latest BiHub activities"
                    icon={<Activity className="h-5 w-5 text-white" />}
                  >
                    <div className="space-y-3">
                      {orders.slice(0, 3).map((order, index) => (
                        <div key={index} className="flex items-center gap-3 p-3 bg-gray-800 rounded-lg">
                          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                            {getOrderStatusIcon(order.status)}
                          </div>
                          <div className="flex-1">
                            <p className="text-white text-sm font-medium">
                              Order #{order.order_number || `ORD-${index + 1}`}
                            </p>
                            <p className="text-gray-400 text-xs">
                              {formatDate(new Date())} • {formatPrice(order.total_amount || 0)}
                            </p>
                          </div>
                          <BiHubStatusBadge status={getOrderStatusColor(order.status)}>
                            {order.status}
                          </BiHubStatusBadge>
                        </div>
                      ))}

                      {orders.length === 0 && (
                        <div className="text-center py-6">
                          <ShoppingBag className="h-8 w-8 text-gray-600 mx-auto mb-2" />
                          <p className="text-gray-400 text-sm">No recent activity</p>
                        </div>
                      )}
                    </div>
                  </BiHubAdminCard>

                  {/* Membership Tier Card */}
                  <BiHubAdminCard
                    title="Membership Tier"
                    subtitle="Your current membership benefits"
                    icon={<Award className="h-5 w-5 text-white" />}
                  >
                    <div className="space-y-4">
                      <div className="text-center">
                        <div className={cn(
                          "inline-flex items-center px-4 py-2 rounded-full text-sm font-medium",
                          userStats.membershipTier === 'bronze' && "bg-amber-100 text-amber-800",
                          userStats.membershipTier === 'silver' && "bg-gray-100 text-gray-800",
                          userStats.membershipTier === 'gold' && "bg-yellow-100 text-yellow-800",
                          userStats.membershipTier === 'platinum' && "bg-purple-100 text-purple-800"
                        )}>
                          <Shield className="h-4 w-4 mr-2" />
                          {userStats.membershipTier.charAt(0).toUpperCase() + userStats.membershipTier.slice(1)} Member
                        </div>
                      </div>

                      <div className="space-y-3">
                        <div className="flex items-center justify-between">
                          <span className="text-gray-400 text-sm">Loyalty Points</span>
                          <span className="text-white font-medium">{userStats.loyaltyPoints.toLocaleString()}</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-gray-400 text-sm">Total Spent</span>
                          <span className="text-white font-medium">{formatPrice(userStats.totalSpent)}</span>
                        </div>

                        {/* Membership Benefits */}
                        <div className="mt-4 p-3 bg-gray-800 rounded-lg">
                          <p className="text-white text-sm font-medium mb-2">Current Benefits:</p>
                          <ul className="text-gray-400 text-xs space-y-1">
                            {userStats.membershipTier === 'bronze' && (
                              <li>• Basic member benefits</li>
                            )}
                            {userStats.membershipTier === 'silver' && (
                              <>
                                <li>• 5% discount on orders</li>
                                <li>• Priority customer support</li>
                              </>
                            )}
                            {userStats.membershipTier === 'gold' && (
                              <>
                                <li>• 10% discount on orders</li>
                                <li>• Free shipping on all orders</li>
                                <li>• Priority customer support</li>
                              </>
                            )}
                            {userStats.membershipTier === 'platinum' && (
                              <>
                                <li>• 15% discount on orders</li>
                                <li>• Free shipping on all orders</li>
                                <li>• Priority customer support</li>
                                <li>• Exclusive early access to sales</li>
                              </>
                            )}
                          </ul>
                        </div>
                      </div>
                    </div>
                  </BiHubAdminCard>
                </div>
              </TabsContent>

              {/* Edit Profile Tab */}
              <TabsContent value="profile">
                <BiHubAdminCard
                  title="Edit Profile"
                  subtitle="Update your personal information"
                  icon={<Edit className="h-5 w-5 text-white" />}
                  headerAction={
                    <Button
                      onClick={() => setIsEditing(!isEditing)}
                      variant="outline"
                      className={cn(
                        BIHUB_ADMIN_THEME.components.button.ghost,
                        isEditing && 'bg-red-600 hover:bg-red-700 text-white'
                      )}
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
                  }
                >
                  <div className="space-y-6">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <Label htmlFor="firstName" className="text-gray-300">First Name</Label>
                        <Input
                          id="firstName"
                          value={profileData.first_name}
                          onChange={(e) => setProfileData(prev => ({ ...prev, first_name: e.target.value }))}
                          disabled={!isEditing}
                          className={BIHUB_ADMIN_THEME.components.input.base}
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="lastName" className="text-gray-300">Last Name</Label>
                        <Input
                          id="lastName"
                          value={profileData.last_name}
                          onChange={(e) => setProfileData(prev => ({ ...prev, last_name: e.target.value }))}
                          disabled={!isEditing}
                          className={BIHUB_ADMIN_THEME.components.input.base}
                        />
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="email" className="text-gray-300">Email Address</Label>
                      <div className="relative">
                        <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <Input
                          id="email"
                          value={user.email}
                          disabled
                          className={cn(BIHUB_ADMIN_THEME.components.input.base, 'pl-10 opacity-60')}
                        />
                      </div>
                      <p className="text-xs text-gray-500">
                        Email cannot be changed. Contact BiHub support if needed.
                      </p>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="phone" className="text-gray-300">Phone Number</Label>
                      <div className="relative">
                        <Phone className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <Input
                          id="phone"
                          value={profileData.profile.phone}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, phone: e.target.value }
                          }))}
                          disabled={!isEditing}
                          className={cn(BIHUB_ADMIN_THEME.components.input.base, 'pl-10')}
                        />
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <Label htmlFor="dateOfBirth" className="text-gray-300">Date of Birth</Label>
                        <Input
                          id="dateOfBirth"
                          type="date"
                          value={profileData.profile.date_of_birth}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, date_of_birth: e.target.value }
                          }))}
                          disabled={!isEditing}
                          className={BIHUB_ADMIN_THEME.components.input.base}
                        />
                      </div>

                      <div className="space-y-2">
                        <Label htmlFor="gender" className="text-gray-300">Gender</Label>
                        <Select
                          value={profileData.profile.gender}
                          onValueChange={(value) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, gender: value }
                          }))}
                          disabled={!isEditing}
                        >
                          <SelectTrigger className={BIHUB_ADMIN_THEME.components.input.base}>
                            <SelectValue placeholder="Select gender" />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-gray-700">
                            <SelectItem value="male">Male</SelectItem>
                            <SelectItem value="female">Female</SelectItem>
                            <SelectItem value="other">Other</SelectItem>
                            <SelectItem value="prefer_not_to_say">Prefer not to say</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="bio" className="text-gray-300">Bio</Label>
                      <Textarea
                        id="bio"
                        value={profileData.profile.bio}
                        onChange={(e) => setProfileData(prev => ({
                          ...prev,
                          profile: { ...prev.profile, bio: e.target.value }
                        }))}
                        disabled={!isEditing}
                        className={BIHUB_ADMIN_THEME.components.input.base}
                        rows={3}
                        placeholder="Tell us about yourself..."
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="website" className="text-gray-300">Website</Label>
                      <div className="relative">
                        <Globe className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                        <Input
                          id="website"
                          value={profileData.profile.website}
                          onChange={(e) => setProfileData(prev => ({
                            ...prev,
                            profile: { ...prev.profile, website: e.target.value }
                          }))}
                          disabled={!isEditing}
                          className={cn(BIHUB_ADMIN_THEME.components.input.base, 'pl-10')}
                          placeholder="https://your-website.com"
                        />
                      </div>
                    </div>

                    {isEditing && (
                      <div className="flex items-center gap-4 pt-6 border-t border-gray-700">
                        <Button
                          onClick={handleProfileUpdate}
                          disabled={updateProfile.isPending}
                          className="bg-[#FF9000] hover:bg-[#e67e00] text-white"
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
                          className="border-gray-600 text-gray-300 hover:text-white hover:bg-gray-800"
                        >
                          Cancel
                        </Button>
                      </div>
                    )}
                  </div>
                </BiHubAdminCard>
              </TabsContent>

              {/* Orders Tab */}
              <TabsContent value="orders">
                <BiHubAdminCard
                  title="Order History"
                  subtitle="Your BiHub order history and tracking"
                  icon={<ShoppingBag className="h-5 w-5 text-white" />}
                  headerAction={
                    <Button
                      variant="outline"
                      className={BIHUB_ADMIN_THEME.components.button.ghost}
                    >
                      View All Orders
                    </Button>
                  }
                >
                  {ordersLoading ? (
                    <div className="space-y-4">
                      {[...Array(3)].map((_, i) => (
                        <div key={i} className="animate-pulse">
                          <div className="h-20 bg-gray-700 rounded-xl"></div>
                        </div>
                      ))}
                    </div>
                  ) : orders.length > 0 ? (
                    <div className="space-y-4">
                      {orders.map((order, index) => (
                        <div key={index} className="p-4 bg-gray-800 rounded-lg hover:bg-gray-750 transition-colors">
                          <div className="flex items-center justify-between mb-3">
                            <div className="flex items-center gap-3">
                              <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
                                {getOrderStatusIcon(order.status)}
                              </div>
                              <div>
                                <h3 className="font-semibold text-white">
                                  Order #{order.order_number || `ORD-${index + 1}`}
                                </h3>
                                <p className="text-sm text-gray-400">
                                  {formatDate(new Date())}
                                </p>
                              </div>
                            </div>
                            <BiHubStatusBadge status={getOrderStatusColor(order.status)}>
                              {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                            </BiHubStatusBadge>
                          </div>

                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-4">
                              <div className="flex items-center gap-2">
                                <Package className="h-4 w-4 text-gray-400" />
                                <span className="text-sm text-gray-300">{order.items?.length || 1} items</span>
                              </div>
                              <div className="flex items-center gap-2">
                                <span className="text-lg font-bold text-[#FF9000]">
                                  {formatPrice(order.total_amount)}
                                </span>
                              </div>
                            </div>

                            <Button
                              variant="outline"
                              size="sm"
                              className={BIHUB_ADMIN_THEME.components.button.ghost}
                            >
                              View Details
                            </Button>
                          </div>
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="text-center py-12">
                      <ShoppingBag className="h-16 w-16 text-gray-600 mx-auto mb-4" />
                      <h3 className="text-lg font-semibold text-white mb-2">No orders yet</h3>
                      <p className="text-gray-400 mb-6">
                        Start shopping on BiHub to see your orders here.
                      </p>
                      <Button className="bg-[#FF9000] hover:bg-[#e67e00] text-white">
                        <a href="/products">Browse Products</a>
                      </Button>
                    </div>
                  )}
                </BiHubAdminCard>
              </TabsContent>

              {/* Security Tab */}
              <TabsContent value="security">
                <BiHubAdminCard
                  title="Security Settings"
                  subtitle="Manage your BiHub account security"
                  icon={<Lock className="h-5 w-5 text-white" />}
                >
                  <div className="space-y-6">
                    <div className="p-6 bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between mb-4">
                        <div>
                          <h3 className="font-semibold text-white">Password</h3>
                          <p className="text-sm text-gray-400">
                            Last changed {formatDate(user.updated_at)}
                          </p>
                        </div>
                        <Button
                          variant="outline"
                          onClick={() => setShowPasswordForm(!showPasswordForm)}
                          className={BIHUB_ADMIN_THEME.components.button.ghost}
                        >
                          <Lock className="h-4 w-4 mr-2" />
                          Change Password
                        </Button>
                      </div>

                      {showPasswordForm && (
                        <div className="space-y-4 pt-4 border-t border-gray-700">
                          <div className="space-y-2">
                            <Label htmlFor="currentPassword" className="text-gray-300">Current Password</Label>
                            <Input
                              id="currentPassword"
                              type="password"
                              value={passwordData.current_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, current_password: e.target.value }))}
                              className={BIHUB_ADMIN_THEME.components.input.base}
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="newPassword" className="text-gray-300">New Password</Label>
                            <Input
                              id="newPassword"
                              type="password"
                              value={passwordData.new_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, new_password: e.target.value }))}
                              className={BIHUB_ADMIN_THEME.components.input.base}
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="confirmPassword" className="text-gray-300">Confirm New Password</Label>
                            <Input
                              id="confirmPassword"
                              type="password"
                              value={passwordData.confirm_password}
                              onChange={(e) => setPasswordData(prev => ({ ...prev, confirm_password: e.target.value }))}
                              className={BIHUB_ADMIN_THEME.components.input.base}
                            />
                          </div>

                          <div className="flex items-center gap-4">
                            <Button
                              onClick={handlePasswordChange}
                              disabled={changePassword.isPending}
                              className="bg-[#FF9000] hover:bg-[#e67e00] text-white"
                            >
                              {changePassword.isPending ? 'Changing...' : 'Change Password'}
                            </Button>
                            <Button
                              variant="outline"
                              onClick={() => setShowPasswordForm(false)}
                              className="border-gray-600 text-gray-300 hover:text-white hover:bg-gray-800"
                            >
                              Cancel
                            </Button>
                          </div>
                        </div>
                      )}
                    </div>

                    <div className="p-6 bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-white">Two-Factor Authentication</h3>
                          <p className="text-sm text-gray-400">
                            Add an extra layer of security to your BiHub account
                          </p>
                        </div>
                        <Button
                          variant="outline"
                          className={BIHUB_ADMIN_THEME.components.button.ghost}
                        >
                          <Shield className="h-4 w-4 mr-2" />
                          Enable 2FA
                        </Button>
                      </div>
                    </div>

                    <div className="p-6 bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold text-white">Login Sessions</h3>
                          <p className="text-sm text-gray-400">
                            Manage your active login sessions
                          </p>
                        </div>
                        <Button
                          variant="outline"
                          className={BIHUB_ADMIN_THEME.components.button.ghost}
                        >
                          <Activity className="h-4 w-4 mr-2" />
                          View Sessions
                        </Button>
                      </div>
                    </div>
                  </div>
                </BiHubAdminCard>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </div>
    </div>
  )
}
