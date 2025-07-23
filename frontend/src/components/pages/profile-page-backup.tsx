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
import { useProfile, useUpdateProfile, useChangePassword, useUserPreferences, useUpdateUserPreferences, useUserSessions } from '@/hooks/use-users'
import { useOrders } from '@/hooks/use-orders'
import { formatDate, formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'

// Remove unused BiHub components

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
          phone: user.phone || '', // Phone is at user level in BE
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
      // Backend UpdateProfileRequest only supports: first_name, last_name, phone
      const updateData = {
        first_name: profileData.first_name,
        last_name: profileData.last_name,
        phone: profileData.profile.phone, // Note: phone is at user level in BE, not profile level
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
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="text-center">
          <div className="w-16 h-16 bg-gradient-to-br from-red-500 to-red-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <Lock className="h-8 w-8 text-white" />
          </div>
          <h1 className="text-3xl font-bold text-white mb-4">Access Denied</h1>
          <p className="text-gray-400 mb-8">
            You need to be logged in to view your profile.
          </p>
          <Button className="bg-[#ff9000] hover:bg-orange-600 text-white">
            <a href="/auth/login">Login to BiHub</a>
          </Button>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-black py-8">
      <div className="container mx-auto px-4 max-w-7xl">
        {/* Modern Profile Header */}
        <div className="relative overflow-hidden bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl shadow-2xl shadow-black/20 mb-8">
          <div className="absolute inset-0 bg-gradient-to-br from-[#ff9000]/5 to-orange-600/5"></div>
          <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-to-br from-[#ff9000]/10 to-transparent rounded-full -translate-y-48 translate-x-48"></div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-gradient-to-tr from-blue-500/10 to-transparent rounded-full translate-y-32 -translate-x-32"></div>

          <div className="relative p-8">
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-8">
              <div className="text-white">
                <div className="flex items-center gap-4 mb-6">
                  <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 backdrop-blur-sm border border-[#ff9000]/30 flex items-center justify-center">
                    <User className="h-8 w-8 text-[#ff9000]" />
                  </div>
                  <div>
                    <span className="text-[#ff9000] font-bold text-sm tracking-wider">PROFILE DASHBOARD</span>
                    <p className="text-white/60 text-sm">BiHub Account Management</p>
                  </div>
                </div>

                <h1 className="text-4xl lg:text-5xl font-bold mb-4 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                  Welcome back, <span className="text-[#ff9000]">{user.first_name}!</span>
                </h1>
                <p className="text-xl text-white/70 leading-relaxed max-w-2xl">
                  Manage your account settings, view order history, and customize your BiHub experience
                </p>
              </div>

              <div className="flex flex-col sm:flex-row gap-4">
                <div className="bg-gradient-to-br from-white/[0.12] to-white/[0.04] backdrop-blur-sm border border-white/20 rounded-2xl p-6 text-center shadow-lg">
                  <p className="text-white/60 text-sm font-medium mb-1">Member Since</p>
                  <p className="text-white font-bold text-2xl">
                    {userStats.memberSince}
                  </p>
                </div>

                <div className="bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 backdrop-blur-sm border border-[#ff9000]/30 rounded-2xl p-6 text-center shadow-lg">
                  <p className="text-white/60 text-sm font-medium mb-1">Total Orders</p>
                  <p className="text-[#ff9000] font-bold text-2xl">
                    {userStats.totalOrders}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Modern Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 flex items-center justify-center">
                <ShoppingBag className="h-6 w-6 text-[#ff9000]" />
              </div>
              <span className="text-2xl font-bold text-white">{userStats.totalOrders}</span>
            </div>
            <h3 className="text-white/80 font-semibold">Total Orders</h3>
            <p className="text-white/50 text-sm">All time purchases</p>
          </div>

          <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-600/20 flex items-center justify-center">
                <CheckCircle className="h-6 w-6 text-green-400" />
              </div>
              <span className="text-2xl font-bold text-white">{userStats.completedOrders}</span>
            </div>
            <h3 className="text-white/80 font-semibold">Completed</h3>
            <p className="text-white/50 text-sm">Successfully delivered</p>
          </div>

          <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-600/20 flex items-center justify-center">
                <DollarSign className="h-6 w-6 text-blue-400" />
              </div>
              <span className="text-2xl font-bold text-white">{formatPrice(userStats.totalSpent)}</span>
            </div>
            <h3 className="text-white/80 font-semibold">Total Spent</h3>
            <p className="text-white/50 text-sm">Lifetime value</p>
          </div>

          <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-600/20 flex items-center justify-center">
                <Award className="h-6 w-6 text-purple-400" />
              </div>
              <span className="text-2xl font-bold text-white">{userStats.loyaltyPoints.toLocaleString()}</span>
            </div>
            <h3 className="text-white/80 font-semibold">Loyalty Points</h3>
            <p className="text-white/50 text-sm">Rewards earned</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-8">
          {/* Profile Sidebar */}
          <div className="lg:col-span-1">
            <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl shadow-xl shadow-black/10 p-6">
              <div className="text-center">
                {/* Avatar */}
                <div className="relative mb-6">
                  <div className="w-32 h-32 rounded-2xl overflow-hidden mx-auto bg-gradient-to-br from-gray-700 to-gray-800 flex items-center justify-center border-4 border-white/10">
                    {user.profile?.avatar_url ? (
                      <Image
                        src={user.profile.avatar_url}
                        alt={`${user.first_name} ${user.last_name}`}
                        width={128}
                        height={128}
                        className="object-cover"
                      />
                    ) : (
                      <User className="w-16 h-16 text-gray-400" />
                    )}
                  </div>

                  <Dialog open={showAvatarUpload} onOpenChange={setShowAvatarUpload}>
                    <DialogTrigger asChild>
                      <Button
                        size="sm"
                        variant="outline"
                        className="absolute -bottom-2 -right-2 rounded-full w-12 h-12 p-0 bg-gradient-to-br from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 border-2 border-white/20 text-white shadow-lg"
                      >
                        <Camera className="h-5 w-5" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="bg-gradient-to-br from-gray-900 to-gray-800 border border-white/20 backdrop-blur-xl">
                      <DialogHeader>
                        <DialogTitle className="text-white text-xl font-bold">Update Profile Picture</DialogTitle>
                      </DialogHeader>
                      <div className="space-y-6">
                        <div className="flex items-center justify-center w-full">
                          <label className="flex flex-col items-center justify-center w-full h-40 border-2 border-white/20 border-dashed rounded-2xl cursor-pointer bg-white/[0.02] hover:bg-white/[0.05] transition-all duration-300">
                            <div className="flex flex-col items-center justify-center pt-5 pb-6">
                              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 flex items-center justify-center mb-4">
                                <Upload className="w-6 h-6 text-[#ff9000]" />
                              </div>
                              <p className="mb-2 text-sm text-white">
                                <span className="font-semibold">Click to upload</span> or drag and drop
                              </p>
                              <p className="text-xs text-white/60">PNG, JPG or GIF (MAX. 2MB)</p>
                            </div>
                            <input type="file" className="hidden" accept="image/*" />
                          </label>
                        </div>
                      </div>
                    </DialogContent>
                  </Dialog>
                </div>

                {/* User Info */}
                <h2 className="text-2xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent mb-2">
                  {user.first_name} {user.last_name}
                </h2>
                <p className="text-white/60 mb-6">{user.email}</p>

                <div className="flex flex-col gap-3 mb-8">
                  <div className="inline-flex items-center px-3 py-2 rounded-xl bg-gradient-to-r from-blue-500/20 to-cyan-500/20 border border-blue-500/30">
                    <Shield className="h-4 w-4 mr-2 text-blue-400" />
                    <span className="text-blue-400 font-semibold text-sm">
                      {user.role.charAt(0).toUpperCase() + user.role.slice(1)}
                    </span>
                  </div>
                  <div className={cn(
                    "inline-flex items-center px-3 py-2 rounded-xl border",
                    user.is_active
                      ? "bg-gradient-to-r from-green-500/20 to-emerald-500/20 border-green-500/30 text-green-400"
                      : "bg-gradient-to-r from-red-500/20 to-pink-500/20 border-red-500/30 text-red-400"
                  )}>
                    <div className={cn(
                      "w-2 h-2 rounded-full mr-2",
                      user.is_active ? "bg-green-400" : "bg-red-400"
                    )}></div>
                    <span className="font-semibold text-sm">
                      {user.is_active ? 'Active Account' : 'Inactive Account'}
                    </span>
                  </div>
                </div>

                {/* Quick Actions */}
                <div className="space-y-3">
                  <Button
                    onClick={() => setActiveTab('profile')}
                    className="w-full bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold py-3 rounded-xl shadow-lg"
                  >
                    <Edit className="h-4 w-4 mr-2" />
                    Edit Profile
                  </Button>
                  <Button
                    onClick={() => setActiveTab('security')}
                    variant="outline"
                    className="w-full border-white/30 text-white hover:text-white hover:bg-white/10 font-semibold py-3 rounded-xl backdrop-blur-sm"
                  >
                    <Lock className="h-4 w-4 mr-2" />
                    Security Settings
                  </Button>
                </div>
              </div>
            </div>
          </div>

          {/* Main Content */}
          <div className="lg:col-span-3">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-2 shadow-xl shadow-black/10">
                <TabsList className="grid w-full grid-cols-5 bg-transparent gap-2">
                  <TabsTrigger
                    value="overview"
                    className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
                  >
                    <User className="h-4 w-4" />
                    Overview
                  </TabsTrigger>
                  <TabsTrigger
                    value="orders"
                    className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
                  >
                    <ShoppingBag className="h-4 w-4" />
                    Orders
                  </TabsTrigger>
                  <TabsTrigger
                    value="profile"
                    className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
                  >
                    <Edit className="h-4 w-4" />
                    Profile
                  </TabsTrigger>
                  <TabsTrigger
                    value="preferences"
                    className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
                  >
                    <Globe className="h-4 w-4" />
                    Settings
                  </TabsTrigger>
                  <TabsTrigger
                    value="security"
                    className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
                  >
                    <Lock className="h-4 w-4" />
                    Security
                  </TabsTrigger>
                </TabsList>
              </div>

              {/* Overview Tab */}
              <TabsContent value="overview">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                  <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                        <User className="h-5 w-5 text-blue-400" />
                      </div>
                      <div>
                        <h3 className="text-white font-bold text-lg">Account Information</h3>
                        <p className="text-white/60 text-sm">Your basic account details</p>
                      </div>
                    </div>

                    <div className="space-y-6">
                      <div className="flex items-center gap-4 p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                          <Mail className="h-5 w-5 text-green-400" />
                        </div>
                        <div className="flex-1">
                          <p className="text-white/60 text-sm font-medium">Email Address</p>
                          <p className="text-white font-semibold">{user.email}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                          <Phone className="h-5 w-5 text-purple-400" />
                        </div>
                        <div className="flex-1">
                          <p className="text-white/60 text-sm font-medium">Phone Number</p>
                          <p className="text-white font-semibold">{user.phone || 'Not provided'}</p>
                        </div>
                      </div>

                      <div className="flex items-center gap-4 p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 flex items-center justify-center">
                          <Calendar className="h-5 w-5 text-[#ff9000]" />
                        </div>
                        <div className="flex-1">
                          <p className="text-white/60 text-sm font-medium">Member Since</p>
                          <p className="text-white font-semibold">{formatDate(user.created_at)}</p>
                        </div>
                      </div>
                    </div>
                  </div>

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

              {/* Preferences Tab */}
              <TabsContent value="preferences">
                <div className="space-y-8">
                  <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                        <Globe className="h-5 w-5 text-purple-400" />
                      </div>
                      <div>
                        <h3 className="text-white font-bold text-lg">Account Preferences</h3>
                        <p className="text-white/60 text-sm">Customize your BiHub experience</p>
                      </div>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <Label className="text-white/80">Language</Label>
                        <Select defaultValue={user.language || 'en'}>
                          <SelectTrigger className="bg-white/[0.02] border-white/20 text-white">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-white/20">
                            <SelectItem value="en">English</SelectItem>
                            <SelectItem value="vi">Tiếng Việt</SelectItem>
                            <SelectItem value="fr">Français</SelectItem>
                            <SelectItem value="es">Español</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="space-y-2">
                        <Label className="text-white/80">Currency</Label>
                        <Select defaultValue={user.currency || 'USD'}>
                          <SelectTrigger className="bg-white/[0.02] border-white/20 text-white">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-white/20">
                            <SelectItem value="USD">USD ($)</SelectItem>
                            <SelectItem value="VND">VND (₫)</SelectItem>
                            <SelectItem value="EUR">EUR (€)</SelectItem>
                            <SelectItem value="GBP">GBP (£)</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="space-y-2">
                        <Label className="text-white/80">Timezone</Label>
                        <Select defaultValue={user.timezone || 'UTC'}>
                          <SelectTrigger className="bg-white/[0.02] border-white/20 text-white">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-white/20">
                            <SelectItem value="UTC">UTC</SelectItem>
                            <SelectItem value="Asia/Ho_Chi_Minh">Asia/Ho Chi Minh</SelectItem>
                            <SelectItem value="America/New_York">America/New York</SelectItem>
                            <SelectItem value="Europe/London">Europe/London</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="space-y-2">
                        <Label className="text-white/80">Theme</Label>
                        <Select defaultValue="dark">
                          <SelectTrigger className="bg-white/[0.02] border-white/20 text-white">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-white/20">
                            <SelectItem value="dark">Dark</SelectItem>
                            <SelectItem value="light">Light</SelectItem>
                            <SelectItem value="auto">Auto</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <div className="flex justify-end pt-6 border-t border-white/10 mt-6">
                      <Button className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl">
                        Save Preferences
                      </Button>
                    </div>
                  </div>

                  <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                        <Shield className="h-5 w-5 text-green-400" />
                      </div>
                      <div>
                        <h3 className="text-white font-bold text-lg">Privacy Settings</h3>
                        <p className="text-white/60 text-sm">Control your privacy and data sharing</p>
                      </div>
                    </div>

                    <div className="space-y-6">
                      <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div>
                          <h4 className="text-white font-semibold">Profile Visibility</h4>
                          <p className="text-white/60 text-sm">Control who can see your profile</p>
                        </div>
                        <Select defaultValue="private">
                          <SelectTrigger className="w-32 bg-white/[0.02] border-white/20 text-white">
                            <SelectValue />
                          </SelectTrigger>
                          <SelectContent className="bg-gray-900 border-white/20">
                            <SelectItem value="public">Public</SelectItem>
                            <SelectItem value="private">Private</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>

                      <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div>
                          <h4 className="text-white font-semibold">Data Collection</h4>
                          <p className="text-white/60 text-sm">Allow BiHub to collect usage data</p>
                        </div>
                        <Button variant="outline" size="sm" className="border-white/30 text-white hover:bg-white/10">
                          Enabled
                        </Button>
                      </div>

                      <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                        <div>
                          <h4 className="text-white font-semibold">Marketing Emails</h4>
                          <p className="text-white/60 text-sm">Receive promotional emails and offers</p>
                        </div>
                        <Button variant="outline" size="sm" className="border-white/30 text-white hover:bg-white/10">
                          Disabled
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              </TabsContent>

              {/* Security Tab */}
              <TabsContent value="security">
                <div className="space-y-8">
                  <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                    <div className="flex items-center gap-3 mb-6">
                      <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-red-500/20 to-pink-500/20 flex items-center justify-center">
                        <Lock className="h-5 w-5 text-red-400" />
                      </div>
                      <div>
                        <h3 className="text-white font-bold text-lg">Security Settings</h3>
                        <p className="text-white/60 text-sm">Manage your account security</p>
                      </div>
                    </div>

                    <div className="space-y-6">
                      <div className="p-6 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="flex items-center justify-between mb-4">
                          <div>
                            <h4 className="font-semibold text-white text-lg">Password</h4>
                            <p className="text-white/60 text-sm">
                              Last changed {formatDate(user.updated_at)}
                            </p>
                          </div>
                          <Button
                            variant="outline"
                            onClick={() => setShowPasswordForm(!showPasswordForm)}
                            className="border-white/30 text-white hover:bg-white/10 font-semibold px-4 py-2 rounded-xl"
                          >
                            <Lock className="h-4 w-4 mr-2" />
                            Change Password
                          </Button>
                        </div>

                        {showPasswordForm && (
                          <div className="space-y-6 pt-6 border-t border-white/10">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                              <div className="space-y-2">
                                <Label htmlFor="currentPassword" className="text-white/80">Current Password</Label>
                                <Input
                                  id="currentPassword"
                                  type="password"
                                  value={passwordData.current_password}
                                  onChange={(e) => setPasswordData(prev => ({ ...prev, current_password: e.target.value }))}
                                  className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                                  placeholder="Enter current password"
                                />
                              </div>

                              <div className="space-y-2">
                                <Label htmlFor="newPassword" className="text-white/80">New Password</Label>
                                <Input
                                  id="newPassword"
                                  type="password"
                                  value={passwordData.new_password}
                                  onChange={(e) => setPasswordData(prev => ({ ...prev, new_password: e.target.value }))}
                                  className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                                  placeholder="Enter new password"
                                />
                              </div>
                            </div>

                            <div className="space-y-2">
                              <Label htmlFor="confirmPassword" className="text-white/80">Confirm New Password</Label>
                              <Input
                                id="confirmPassword"
                                type="password"
                                value={passwordData.confirm_password}
                                onChange={(e) => setPasswordData(prev => ({ ...prev, confirm_password: e.target.value }))}
                                className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                                placeholder="Confirm new password"
                              />
                            </div>

                            <div className="flex items-center gap-4">
                              <Button
                                onClick={handlePasswordChange}
                                disabled={changePassword.isPending}
                                className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                              >
                                {changePassword.isPending ? (
                                  <>
                                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                                    Changing...
                                  </>
                                ) : (
                                  'Change Password'
                                )}
                              </Button>
                              <Button
                                variant="outline"
                                onClick={() => setShowPasswordForm(false)}
                                className="border-white/30 text-white hover:bg-white/10 font-semibold px-6 py-2 rounded-xl"
                              >
                                Cancel
                              </Button>
                            </div>
                          </div>
                        )}
                    </div>

                      </div>

                      <div className="p-6 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-semibold text-white text-lg">Two-Factor Authentication</h4>
                            <p className="text-white/60 text-sm">
                              Add an extra layer of security to your BiHub account
                            </p>
                          </div>
                          <Button
                            variant="outline"
                            className="border-white/30 text-white hover:bg-white/10 font-semibold px-4 py-2 rounded-xl"
                          >
                            <Shield className="h-4 w-4 mr-2" />
                            Enable 2FA
                          </Button>
                        </div>
                      </div>

                      <div className="p-6 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-semibold text-white text-lg">Active Sessions</h4>
                            <p className="text-white/60 text-sm">
                              Manage your active login sessions across devices
                            </p>
                          </div>
                          <Button
                            variant="outline"
                            className="border-white/30 text-white hover:bg-white/10 font-semibold px-4 py-2 rounded-xl"
                          >
                            <Activity className="h-4 w-4 mr-2" />
                            View Sessions
                          </Button>
                        </div>
                      </div>

                      <div className="p-6 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="flex items-center justify-between">
                          <div>
                            <h4 className="font-semibold text-white text-lg">Account Deletion</h4>
                            <p className="text-white/60 text-sm">
                              Permanently delete your BiHub account and all data
                            </p>
                          </div>
                          <Button
                            variant="outline"
                            className="border-red-500/30 text-red-400 hover:bg-red-500/10 font-semibold px-4 py-2 rounded-xl"
                          >
                            <X className="h-4 w-4 mr-2" />
                            Delete Account
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </div>
    </div>
  )
}
