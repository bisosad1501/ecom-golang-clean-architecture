'use client'

import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { 
  Settings, 
  Palette, 
  Globe, 
  DollarSign, 
  Clock, 
  Bell, 
  Mail, 
  MessageSquare, 
  Shield, 
  Eye, 
  ShoppingCart,
  CreditCard,
  Heart,
  Save,
  Loader2
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { authService } from '@/lib/services/auth'
import { toast } from 'sonner'
import type { UserPreferences, UpdateUserPreferencesRequest } from '@/types'

const preferencesSchema = z.object({
  // Display preferences
  theme: z.enum(['light', 'dark', 'system']),
  language: z.string().min(2),
  currency: z.string().min(3),
  timezone: z.string().min(3),
  
  // Notification preferences
  email_notifications: z.boolean(),
  sms_notifications: z.boolean(),
  push_notifications: z.boolean(),
  marketing_emails: z.boolean(),
  order_updates: z.boolean(),
  security_alerts: z.boolean(),
  newsletter_enabled: z.boolean(),
  promotional_emails: z.boolean(),
  
  // Privacy preferences
  profile_visibility: z.enum(['public', 'private']),
  activity_visibility: z.enum(['public', 'private']),
  
  // Shopping preferences
  save_payment_methods: z.boolean(),
  auto_apply_coupons: z.boolean(),
  wishlist_visibility: z.enum(['public', 'private']),
})

type PreferencesFormData = z.infer<typeof preferencesSchema>

export function UserPreferencesForm() {
  const [isLoading, setIsLoading] = useState(false)
  const [isFetching, setIsFetching] = useState(true)
  const [preferences, setPreferences] = useState<UserPreferences | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isDirty },
    reset,
    watch,
  } = useForm<PreferencesFormData>({
    resolver: zodResolver(preferencesSchema),
  })

  const watchedTheme = watch('theme')

  // Fetch user preferences on component mount
  useEffect(() => {
    const fetchPreferences = async () => {
      try {
        setIsFetching(true)
        const userPrefs = await authService.getUserPreferences()
        setPreferences(userPrefs)
        
        // Reset form with fetched data
        reset({
          theme: userPrefs.theme as 'light' | 'dark' | 'system',
          language: userPrefs.language,
          currency: userPrefs.currency,
          timezone: userPrefs.timezone,
          email_notifications: userPrefs.email_notifications,
          sms_notifications: userPrefs.sms_notifications,
          push_notifications: userPrefs.push_notifications,
          marketing_emails: userPrefs.marketing_emails,
          order_updates: userPrefs.order_updates,
          security_alerts: userPrefs.security_alerts,
          newsletter_enabled: userPrefs.newsletter_enabled,
          promotional_emails: userPrefs.promotional_emails,
          profile_visibility: userPrefs.profile_visibility as 'public' | 'private',
          activity_visibility: userPrefs.activity_visibility as 'public' | 'private',
          save_payment_methods: userPrefs.save_payment_methods,
          auto_apply_coupons: userPrefs.auto_apply_coupons,
          wishlist_visibility: userPrefs.wishlist_visibility as 'public' | 'private',
        })
      } catch (error: any) {
        console.error('Failed to fetch preferences:', error)
        toast.error('Failed to load preferences')
      } finally {
        setIsFetching(false)
      }
    }

    fetchPreferences()
  }, [reset])

  const onSubmit = async (data: PreferencesFormData) => {
    try {
      setIsLoading(true)
      console.log('🔄 [PREFERENCES] Updating preferences:', data)
      
      const updateData: UpdateUserPreferencesRequest = {
        theme: data.theme,
        language: data.language,
        currency: data.currency,
        timezone: data.timezone,
        email_notifications: data.email_notifications,
        sms_notifications: data.sms_notifications,
        push_notifications: data.push_notifications,
        marketing_emails: data.marketing_emails,
        order_updates: data.order_updates,
        security_alerts: data.security_alerts,
        newsletter_enabled: data.newsletter_enabled,
        promotional_emails: data.promotional_emails,
        profile_visibility: data.profile_visibility,
        activity_visibility: data.activity_visibility,
        save_payment_methods: data.save_payment_methods,
        auto_apply_coupons: data.auto_apply_coupons,
        wishlist_visibility: data.wishlist_visibility,
      }
      
      const updatedPrefs = await authService.updateUserPreferences(updateData)
      setPreferences(updatedPrefs)
      
      console.log('✅ [PREFERENCES] Preferences updated successfully')
      toast.success('Preferences updated successfully!')
      
    } catch (error: any) {
      console.error('❌ [PREFERENCES] Error:', error)
      toast.error(error.message || 'Failed to update preferences')
    } finally {
      setIsLoading(false)
    }
  }

  if (isFetching) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="flex items-center space-x-2">
          <Loader2 className="h-6 w-6 animate-spin text-[#FF9000]" />
          <span className="text-gray-600">Loading preferences...</span>
        </div>
      </div>
    )
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
      {/* Display Preferences */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Palette className="h-5 w-5 text-[#FF9000]" />
            <span>Display Preferences</span>
          </CardTitle>
          <CardDescription>
            Customize how the interface looks and feels
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Theme */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Palette className="h-4 w-4" />
                <span>Theme</span>
              </label>
              <select
                {...register('theme')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="light">Light</option>
                <option value="dark">Dark</option>
                <option value="system">System</option>
              </select>
              {errors.theme && (
                <p className="text-sm text-red-600">{errors.theme.message}</p>
              )}
            </div>

            {/* Language */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Globe className="h-4 w-4" />
                <span>Language</span>
              </label>
              <select
                {...register('language')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="en">English</option>
                <option value="vi">Tiếng Việt</option>
                <option value="zh">中文</option>
                <option value="ja">日本語</option>
                <option value="ko">한국어</option>
              </select>
              {errors.language && (
                <p className="text-sm text-red-600">{errors.language.message}</p>
              )}
            </div>

            {/* Currency */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <DollarSign className="h-4 w-4" />
                <span>Currency</span>
              </label>
              <select
                {...register('currency')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="USD">USD - US Dollar</option>
                <option value="VND">VND - Vietnamese Dong</option>
                <option value="EUR">EUR - Euro</option>
                <option value="GBP">GBP - British Pound</option>
                <option value="JPY">JPY - Japanese Yen</option>
              </select>
              {errors.currency && (
                <p className="text-sm text-red-600">{errors.currency.message}</p>
              )}
            </div>

            {/* Timezone */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Clock className="h-4 w-4" />
                <span>Timezone</span>
              </label>
              <select
                {...register('timezone')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="UTC">UTC</option>
                <option value="Asia/Ho_Chi_Minh">Asia/Ho Chi Minh</option>
                <option value="America/New_York">America/New York</option>
                <option value="Europe/London">Europe/London</option>
                <option value="Asia/Tokyo">Asia/Tokyo</option>
              </select>
              {errors.timezone && (
                <p className="text-sm text-red-600">{errors.timezone.message}</p>
              )}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Notification Preferences */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Bell className="h-5 w-5 text-[#FF9000]" />
            <span>Notification Preferences</span>
          </CardTitle>
          <CardDescription>
            Choose how you want to receive notifications
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Email Notifications */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Mail className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Email Notifications</label>
                  <p className="text-xs text-gray-500">Receive notifications via email</p>
                </div>
              </div>
              <input
                {...register('email_notifications')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* SMS Notifications */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <MessageSquare className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">SMS Notifications</label>
                  <p className="text-xs text-gray-500">Receive notifications via SMS</p>
                </div>
              </div>
              <input
                {...register('sms_notifications')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* Push Notifications */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Bell className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Push Notifications</label>
                  <p className="text-xs text-gray-500">Receive browser push notifications</p>
                </div>
              </div>
              <input
                {...register('push_notifications')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* Order Updates */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <ShoppingCart className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Order Updates</label>
                  <p className="text-xs text-gray-500">Get notified about order status</p>
                </div>
              </div>
              <input
                {...register('order_updates')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* Marketing Emails */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Mail className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Marketing Emails</label>
                  <p className="text-xs text-gray-500">Receive promotional emails</p>
                </div>
              </div>
              <input
                {...register('marketing_emails')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* Security Alerts */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Shield className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Security Alerts</label>
                  <p className="text-xs text-gray-500">Important security notifications</p>
                </div>
              </div>
              <input
                {...register('security_alerts')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Privacy Preferences */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <Eye className="h-5 w-5 text-[#FF9000]" />
            <span>Privacy Preferences</span>
          </CardTitle>
          <CardDescription>
            Control who can see your information
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Profile Visibility */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Eye className="h-4 w-4" />
                <span>Profile Visibility</span>
              </label>
              <select
                {...register('profile_visibility')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
            </div>

            {/* Activity Visibility */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Settings className="h-4 w-4" />
                <span>Activity Visibility</span>
              </label>
              <select
                {...register('activity_visibility')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
            </div>

            {/* Wishlist Visibility */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700 flex items-center space-x-2">
                <Heart className="h-4 w-4" />
                <span>Wishlist Visibility</span>
              </label>
              <select
                {...register('wishlist_visibility')}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-[#FF9000] focus:border-[#FF9000] transition-colors"
              >
                <option value="private">Private</option>
                <option value="public">Public</option>
              </select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Shopping Preferences */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center space-x-2">
            <ShoppingCart className="h-5 w-5 text-[#FF9000]" />
            <span>Shopping Preferences</span>
          </CardTitle>
          <CardDescription>
            Customize your shopping experience
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="space-y-4">
            {/* Save Payment Methods */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <CreditCard className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Save Payment Methods</label>
                  <p className="text-xs text-gray-500">Securely save payment methods for faster checkout</p>
                </div>
              </div>
              <input
                {...register('save_payment_methods')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>

            {/* Auto Apply Coupons */}
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Settings className="h-4 w-4 text-gray-500" />
                <div>
                  <label className="text-sm font-medium text-gray-700">Auto Apply Coupons</label>
                  <p className="text-xs text-gray-500">Automatically apply available coupons at checkout</p>
                </div>
              </div>
              <input
                {...register('auto_apply_coupons')}
                type="checkbox"
                className="h-4 w-4 text-[#FF9000] focus:ring-[#FF9000] border-gray-300 rounded"
              />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Save Button */}
      <div className="flex justify-end">
        <Button
          type="submit"
          disabled={!isDirty || isLoading}
          className="bg-[#FF9000] hover:bg-[#e67e00] text-white px-8 py-2 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          isLoading={isLoading}
          loadingText="Saving..."
        >
          <Save className="h-4 w-4 mr-2" />
          Save Preferences
        </Button>
      </div>
    </form>
  )
}
