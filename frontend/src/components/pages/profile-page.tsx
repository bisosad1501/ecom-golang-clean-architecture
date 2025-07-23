'use client'

import React, { useState, useEffect } from 'react'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AnimatedBackground } from '@/components/ui/animated-background'
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
  CheckCircle,
  Upload,
  Globe,
  Award,
  Activity,
  DollarSign,
} from 'lucide-react'
import { useProfile, useUpdateProfile, useChangePassword, useUserStats, useUserSessions, useUserActivity } from '@/hooks/use-users'
import { useOrders } from '@/hooks/use-orders'
import { formatDate, formatPrice, cn } from '@/lib/utils'
import { toast } from 'sonner'

export default function ProfilePage() {
  const [isEditing, setIsEditing] = useState(false)
  const [activeTab, setActiveTab] = useState('overview')
  const [showPasswordForm, setShowPasswordForm] = useState(false)
  const [passwordData, setPasswordData] = useState({
    current_password: '',
    new_password: '',
    confirm_password: '',
  })

  const [profileData, setProfileData] = useState({
    first_name: '',
    last_name: '',
    profile: {
      phone: '',
    }
  })

  const { data: user, isLoading: userLoading } = useProfile()
  const { data: ordersData, isLoading: ordersLoading } = useOrders({ limit: 10 })
  const { data: userStatsData } = useUserStats()
  const { data: userSessionsData } = useUserSessions(5, 0)
  const { data: userActivityData } = useUserActivity(10, 0)
  const updateProfile = useUpdateProfile()
  const changePassword = useChangePassword()

  const orders = ordersData?.data || []

  // Use backend user stats as single source of truth, fallback to user entity
  const userStats = {
    totalOrders: userStatsData?.total_orders || user?.total_orders || 0,
    completedOrders: userStatsData?.completed_orders || orders.filter(order => order.status === 'delivered').length,
    totalSpent: userStatsData?.total_spent || user?.total_spent || 0,
    loyaltyPoints: userStatsData?.loyalty_points || user?.loyalty_points || 0,
    membershipTier: userStatsData?.membership_tier || user?.membership_tier || 'bronze',
    memberSince: user?.created_at ? new Date(user.created_at).getFullYear() : new Date().getFullYear(),
  }

  // Initialize form data when user data loads
  useEffect(() => {
    if (user) {
      setProfileData({
        first_name: user.first_name || '',
        last_name: user.last_name || '',
        profile: {
          phone: user.profile?.phone || '', // Phone is in profile
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
        phone: profileData.profile.phone, // Send phone to BE user level
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
      await changePassword.mutateAsync({
        current_password: passwordData.current_password,
        new_password: passwordData.new_password,
      })
      setShowPasswordForm(false)
      setPasswordData({
        current_password: '',
        new_password: '',
        confirm_password: '',
      })
      toast.success('Password changed successfully!')
    } catch (error) {
      console.error('Failed to change password:', error)
      toast.error('Failed to change password')
    }
  }

  if (userLoading) {
    return (
      <div className="min-h-screen bg-black py-8">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-[#ff9000] mx-auto"></div>
            <p className="text-white mt-4">Loading your profile...</p>
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
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden py-8">
      {/* Enhanced Background Pattern - Matching Products/Cart Pages */}
      <AnimatedBackground className="opacity-30" />

      <div className="container mx-auto px-4 max-w-7xl relative z-10">
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

        {/* Tabs Navigation */}
        <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-2 shadow-xl shadow-black/10 mb-8">
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList className="grid w-full grid-cols-5 bg-transparent gap-2">
              <TabsTrigger
                value="overview"
                className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <User className="h-4 w-4" />
                Overview
              </TabsTrigger>
              <TabsTrigger
                value="profile"
                className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Edit className="h-4 w-4" />
                Profile
              </TabsTrigger>
              <TabsTrigger
                value="orders"
                className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <ShoppingBag className="h-4 w-4" />
                Orders
              </TabsTrigger>
              <TabsTrigger
                value="activity"
                className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Activity className="h-4 w-4" />
                Activity
              </TabsTrigger>
              <TabsTrigger
                value="security"
                className="flex items-center gap-2 px-4 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Lock className="h-4 w-4" />
                Security
              </TabsTrigger>
            </TabsList>

            {/* Overview Tab */}
            <TabsContent value="overview" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                    <User className="h-5 w-5 text-blue-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Account Overview</h3>
                    <p className="text-white/60 text-sm">Your account summary and recent activity</p>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-4">
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
                        <p className="text-white font-semibold">{user.profile?.phone || 'Not provided'}</p>
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

                  <div className="space-y-4">
                    <h4 className="text-white font-semibold mb-4">Recent Orders</h4>
                    {orders.length > 0 ? (
                      orders.slice(0, 3).map((order, index) => (
                        <div key={index} className="flex items-center gap-3 p-3 bg-white/[0.02] rounded-xl border border-white/10">
                          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                            <Package className="h-4 w-4 text-blue-400" />
                          </div>
                          <div className="flex-1">
                            <p className="text-white text-sm font-medium">
                              Order #{order.order_number || `ORD-${index + 1}`}
                            </p>
                            <p className="text-white/60 text-xs">
                              {formatDate(new Date())} • {formatPrice(order.total || 0)}
                            </p>
                          </div>
                          <span className="text-xs px-2 py-1 rounded-full bg-green-500/20 text-green-400">
                            {order.status}
                          </span>
                        </div>
                      ))
                    ) : (
                      <div className="text-center py-6">
                        <ShoppingBag className="h-8 w-8 text-white/30 mx-auto mb-2" />
                        <p className="text-white/60 text-sm">No recent orders</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </TabsContent>

            {/* Profile Tab */}
            <TabsContent value="profile" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                    <Edit className="h-5 w-5 text-blue-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Edit Profile</h3>
                    <p className="text-white/60 text-sm">Update your personal information</p>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-2">
                    <Label className="text-white/80">First Name</Label>
                    <Input
                      value={profileData.first_name}
                      onChange={(e) => setProfileData(prev => ({ ...prev, first_name: e.target.value }))}
                      disabled={!isEditing}
                      className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-white/80">Last Name</Label>
                    <Input
                      value={profileData.last_name}
                      onChange={(e) => setProfileData(prev => ({ ...prev, last_name: e.target.value }))}
                      disabled={!isEditing}
                      className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label className="text-white/80">Email</Label>
                    <Input
                      value={user.email}
                      disabled
                      className="bg-white/[0.02] border-white/20 text-white/60 placeholder:text-white/40"
                    />
                    <p className="text-xs text-white/50">Email cannot be changed</p>
                  </div>

                  <div className="space-y-2">
                    <Label className="text-white/80">Phone</Label>
                    <Input
                      value={profileData.profile.phone}
                      onChange={(e) => setProfileData(prev => ({
                        ...prev,
                        profile: { ...prev.profile, phone: e.target.value }
                      }))}
                      disabled={!isEditing}
                      className="bg-white/[0.02] border-white/20 text-white placeholder:text-white/40"
                    />
                  </div>
                </div>

                <div className="flex items-center gap-4 pt-6 border-t border-white/10 mt-6">
                  {isEditing ? (
                    <>
                      <Button
                        onClick={handleProfileUpdate}
                        disabled={updateProfile.isPending}
                        className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
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
                        className="border-white/30 text-white hover:bg-white/10 font-semibold px-6 py-2 rounded-xl"
                      >
                        Cancel
                      </Button>
                    </>
                  ) : (
                    <Button
                      onClick={() => setIsEditing(true)}
                      className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                    >
                      <Edit className="h-4 w-4 mr-2" />
                      Edit Profile
                    </Button>
                  )}
                </div>
              </div>
            </TabsContent>

            {/* Orders Tab */}
            <TabsContent value="orders" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                    <ShoppingBag className="h-5 w-5 text-green-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Order History</h3>
                    <p className="text-white/60 text-sm">View and track your orders</p>
                  </div>
                </div>

                {ordersLoading ? (
                  <div className="text-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-[#ff9000] mx-auto mb-4"></div>
                    <p className="text-white/60">Loading orders...</p>
                  </div>
                ) : orders.length > 0 ? (
                  <div className="space-y-4">
                    {orders.map((order, index) => (
                      <div key={index} className="p-6 bg-white/[0.02] rounded-xl border border-white/10">
                        <div className="flex items-center justify-between mb-4">
                          <div className="flex items-center gap-3">
                            <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                              <Package className="h-5 w-5 text-blue-400" />
                            </div>
                            <div>
                              <h4 className="font-semibold text-white">
                                Order #{order.order_number || `ORD-${index + 1}`}
                              </h4>
                              <p className="text-sm text-white/60">
                                {formatDate(new Date())}
                              </p>
                            </div>
                          </div>
                          <span className="px-3 py-1 rounded-full text-sm font-medium bg-green-500/20 text-green-400">
                            {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                          </span>
                        </div>

                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-4">
                            <div className="flex items-center gap-2">
                              <Package className="h-4 w-4 text-white/60" />
                              <span className="text-sm text-white/80">{order.items?.length || 1} items</span>
                            </div>
                            <div className="flex items-center gap-2">
                              <span className="text-lg font-bold text-[#ff9000]">
                                {formatPrice(order.total)}
                              </span>
                            </div>
                          </div>

                          <Button
                            variant="outline"
                            size="sm"
                            className="border-white/30 text-white hover:bg-white/10"
                          >
                            View Details
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <ShoppingBag className="h-16 w-16 text-white/30 mx-auto mb-4" />
                    <h4 className="text-lg font-semibold text-white mb-2">No orders yet</h4>
                    <p className="text-white/60 mb-6">
                      Start shopping on BiHub to see your orders here.
                    </p>
                    <Button className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white">
                      <a href="/products">Browse Products</a>
                    </Button>
                  </div>
                )}
              </div>
            </TabsContent>

            {/* Activity Tab */}
            <TabsContent value="activity" className="mt-8">
              <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
                {/* User Activity */}
                <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                  <div className="flex items-center gap-3 mb-6">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                      <Activity className="h-5 w-5 text-green-400" />
                    </div>
                    <div>
                      <h3 className="text-white font-bold text-lg">Recent Activity</h3>
                      <p className="text-white/60 text-sm">Your recent actions on BiHub</p>
                    </div>
                  </div>

                  <div className="space-y-4">
                    {userActivityData?.activities?.length > 0 ? (
                      userActivityData.activities.slice(0, 5).map((activity: any, index: number) => (
                        <div key={index} className="flex items-center gap-3 p-3 bg-white/[0.02] rounded-xl border border-white/10">
                          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                            <Package className="h-4 w-4 text-blue-400" />
                          </div>
                          <div className="flex-1">
                            <p className="text-white text-sm font-medium">
                              {activity.description || 'Activity'}
                            </p>
                            <p className="text-white/60 text-xs">
                              {formatDate(activity.created_at || new Date())}
                            </p>
                          </div>
                        </div>
                      ))
                    ) : (
                      <div className="text-center py-6">
                        <Activity className="h-8 w-8 text-white/30 mx-auto mb-2" />
                        <p className="text-white/60 text-sm">No recent activity</p>
                      </div>
                    )}
                  </div>
                </div>

                {/* Active Sessions */}
                <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-6 shadow-xl shadow-black/10">
                  <div className="flex items-center gap-3 mb-6">
                    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                      <Shield className="h-5 w-5 text-purple-400" />
                    </div>
                    <div>
                      <h3 className="text-white font-bold text-lg">Active Sessions</h3>
                      <p className="text-white/60 text-sm">Your login sessions across devices</p>
                    </div>
                  </div>

                  <div className="space-y-4">
                    {userSessionsData?.sessions?.length > 0 ? (
                      userSessionsData.sessions.map((session: any, index: number) => (
                        <div key={index} className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                          <div className="flex items-center gap-3">
                            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                              <Shield className="h-4 w-4 text-green-400" />
                            </div>
                            <div>
                              <p className="text-white text-sm font-medium">
                                {session.device_info || 'Unknown Device'}
                              </p>
                              <p className="text-white/60 text-xs">
                                {session.ip_address} • {formatDate(session.created_at)}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-2">
                            {session.is_current && (
                              <span className="text-xs px-2 py-1 rounded-full bg-green-500/20 text-green-400">
                                Current
                              </span>
                            )}
                            <Button
                              variant="outline"
                              size="sm"
                              className="border-red-500/30 text-red-400 hover:bg-red-500/10 text-xs px-2 py-1"
                            >
                              Revoke
                            </Button>
                          </div>
                        </div>
                      ))
                    ) : (
                      <div className="text-center py-6">
                        <Shield className="h-8 w-8 text-white/30 mx-auto mb-2" />
                        <p className="text-white/60 text-sm">No active sessions</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </TabsContent>

            {/* Security Tab */}
            <TabsContent value="security" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
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
                        <h4 className="font-semibold text-white text-lg">Session Management</h4>
                        <p className="text-white/60 text-sm">
                          View and manage your active login sessions
                        </p>
                      </div>
                      <Button
                        variant="outline"
                        onClick={() => setActiveTab('activity')}
                        className="border-white/30 text-white hover:bg-white/10 font-semibold px-4 py-2 rounded-xl"
                      >
                        <Activity className="h-4 w-4 mr-2" />
                        View Sessions
                      </Button>
                    </div>
                  </div>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}
