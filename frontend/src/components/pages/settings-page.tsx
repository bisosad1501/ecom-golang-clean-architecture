'use client'

import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AnimatedBackground } from '@/components/ui/animated-background'
import {
  Settings,
  Globe,
  Bell,
  Shield,
  Eye,
  Moon,
  Sun,
  Monitor,
  Mail,
  Smartphone,
  Save,
  RefreshCw,
} from 'lucide-react'
import { useProfile, useUserPreferences, useUpdateUserPreferences } from '@/hooks/use-users'
import { formatDate, cn } from '@/lib/utils'
import { toast } from 'sonner'

export function SettingsPage() {
  const [activeTab, setActiveTab] = useState('general')
  const [preferences, setPreferences] = useState({
    language: 'en',
    currency: 'USD',
    timezone: 'UTC',
    theme: 'dark',
    email_notifications: true,
    push_notifications: false,
    marketing_emails: false,
    price_alerts: true,
    stock_alerts: true,
    new_item_alerts: false,
  })

  const { data: user, isLoading: userLoading } = useProfile()
  const { data: userPrefs, isLoading: prefsLoading } = useUserPreferences()
  const updatePreferences = useUpdateUserPreferences()

  // Initialize preferences when data loads
  useEffect(() => {
    if (user && userPrefs) {
      setPreferences({
        language: user.language || userPrefs.language || 'en',
        currency: user.currency || userPrefs.currency || 'USD',
        timezone: user.timezone || userPrefs.timezone || 'UTC',
        theme: userPrefs.theme || 'dark',
        email_notifications: userPrefs.email_notifications ?? true,
        push_notifications: userPrefs.push_notifications ?? false,
        marketing_emails: userPrefs.marketing_emails ?? false,
        price_alerts: userPrefs.price_alerts ?? true,
        stock_alerts: userPrefs.stock_alerts ?? true,
        new_item_alerts: userPrefs.new_item_alerts ?? false,
      })
    }
  }, [user, userPrefs])

  const handleSavePreferences = async () => {
    try {
      await updatePreferences.mutateAsync(preferences)
      toast.success('Settings saved successfully!')
    } catch (error) {
      console.error('Failed to save settings:', error)
      toast.error('Failed to save settings')
    }
  }

  if (userLoading || prefsLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden py-8">
        <AnimatedBackground className="opacity-30" />
        <div className="container mx-auto px-4 relative z-10">
          <div className="text-center">
            <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-[#ff9000] mx-auto"></div>
            <p className="text-white mt-4">Loading settings...</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden py-8">
      {/* Enhanced Background Pattern - Matching Products/Cart Pages */}
      <AnimatedBackground className="opacity-30" />

      <div className="container mx-auto px-4 max-w-6xl relative z-10">
        {/* Header */}
        <div className="relative overflow-hidden bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl shadow-2xl shadow-black/20 mb-8">
          <div className="absolute inset-0 bg-gradient-to-br from-[#ff9000]/5 to-orange-600/5"></div>
          <div className="absolute top-0 right-0 w-64 h-64 bg-gradient-to-br from-[#ff9000]/10 to-transparent rounded-full -translate-y-32 translate-x-32"></div>

          <div className="relative p-8">
            <div className="flex items-center gap-4 mb-6">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 backdrop-blur-sm border border-[#ff9000]/30 flex items-center justify-center">
                <Settings className="h-8 w-8 text-[#ff9000]" />
              </div>
              <div>
                <h1 className="text-4xl font-bold bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
                  Settings
                </h1>
                <p className="text-xl text-white/70 leading-relaxed">
                  Customize your BiHub experience and preferences
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Settings Tabs */}
        <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-2 shadow-xl shadow-black/10 mb-8">
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <TabsList className="grid w-full grid-cols-4 bg-transparent gap-2">
              <TabsTrigger
                value="general"
                className="flex items-center gap-2 px-6 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Globe className="h-4 w-4" />
                General
              </TabsTrigger>
              <TabsTrigger
                value="notifications"
                className="flex items-center gap-2 px-6 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Bell className="h-4 w-4" />
                Notifications
              </TabsTrigger>
              <TabsTrigger
                value="privacy"
                className="flex items-center gap-2 px-6 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Shield className="h-4 w-4" />
                Privacy
              </TabsTrigger>
              <TabsTrigger
                value="appearance"
                className="flex items-center gap-2 px-6 py-3 rounded-xl font-semibold transition-all duration-300 data-[state=active]:bg-gradient-to-r data-[state=active]:from-[#ff9000] data-[state=active]:to-orange-600 data-[state=active]:text-white data-[state=active]:shadow-lg text-white/70 hover:text-white hover:bg-white/10"
              >
                <Eye className="h-4 w-4" />
                Appearance
              </TabsTrigger>
            </TabsList>

            {/* General Tab */}
            <TabsContent value="general" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-blue-500/20 to-cyan-500/20 flex items-center justify-center">
                    <Globe className="h-5 w-5 text-blue-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">General Settings</h3>
                    <p className="text-white/60 text-sm">Configure your basic preferences</p>
                  </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div className="space-y-2">
                    <Label className="text-white/80">Language</Label>
                    <Select 
                      value={preferences.language} 
                      onValueChange={(value) => setPreferences(prev => ({ ...prev, language: value }))}
                    >
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
                    <Select 
                      value={preferences.currency} 
                      onValueChange={(value) => setPreferences(prev => ({ ...prev, currency: value }))}
                    >
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
                    <Select 
                      value={preferences.timezone} 
                      onValueChange={(value) => setPreferences(prev => ({ ...prev, timezone: value }))}
                    >
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
                </div>

                <div className="flex justify-end pt-6 border-t border-white/10 mt-8">
                  <Button 
                    onClick={handleSavePreferences}
                    disabled={updatePreferences.isPending}
                    className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                  >
                    {updatePreferences.isPending ? (
                      <>
                        <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Save Settings
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </TabsContent>

            {/* Notifications Tab */}
            <TabsContent value="notifications" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-green-500/20 to-emerald-500/20 flex items-center justify-center">
                    <Bell className="h-5 w-5 text-green-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Notification Settings</h3>
                    <p className="text-white/60 text-sm">Control how you receive notifications</p>
                  </div>
                </div>

                <div className="space-y-6">
                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div className="flex items-center gap-3">
                      <Mail className="h-5 w-5 text-blue-400" />
                      <div>
                        <h4 className="text-white font-semibold">Email Notifications</h4>
                        <p className="text-white/60 text-sm">Receive notifications via email</p>
                      </div>
                    </div>
                    <Switch
                      checked={preferences.email_notifications}
                      onCheckedChange={(checked) => setPreferences(prev => ({ ...prev, email_notifications: checked }))}
                    />
                  </div>

                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div className="flex items-center gap-3">
                      <Smartphone className="h-5 w-5 text-purple-400" />
                      <div>
                        <h4 className="text-white font-semibold">Push Notifications</h4>
                        <p className="text-white/60 text-sm">Receive push notifications on your device</p>
                      </div>
                    </div>
                    <Switch
                      checked={preferences.push_notifications}
                      onCheckedChange={(checked) => setPreferences(prev => ({ ...prev, push_notifications: checked }))}
                    />
                  </div>

                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div className="flex items-center gap-3">
                      <Bell className="h-5 w-5 text-[#ff9000]" />
                      <div>
                        <h4 className="text-white font-semibold">Price Alerts</h4>
                        <p className="text-white/60 text-sm">Get notified when prices drop</p>
                      </div>
                    </div>
                    <Switch
                      checked={preferences.price_alerts}
                      onCheckedChange={(checked) => setPreferences(prev => ({ ...prev, price_alerts: checked }))}
                    />
                  </div>

                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div className="flex items-center gap-3">
                      <Bell className="h-5 w-5 text-green-400" />
                      <div>
                        <h4 className="text-white font-semibold">Stock Alerts</h4>
                        <p className="text-white/60 text-sm">Get notified when items are back in stock</p>
                      </div>
                    </div>
                    <Switch
                      checked={preferences.stock_alerts}
                      onCheckedChange={(checked) => setPreferences(prev => ({ ...prev, stock_alerts: checked }))}
                    />
                  </div>
                </div>

                <div className="flex justify-end pt-6 border-t border-white/10 mt-8">
                  <Button
                    onClick={handleSavePreferences}
                    disabled={updatePreferences.isPending}
                    className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                  >
                    {updatePreferences.isPending ? (
                      <>
                        <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Save Settings
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </TabsContent>

            {/* Appearance Tab */}
            <TabsContent value="appearance" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                    <Eye className="h-5 w-5 text-purple-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Appearance Settings</h3>
                    <p className="text-white/60 text-sm">Customize how BiHub looks</p>
                  </div>
                </div>

                <div className="space-y-6">
                  <div className="space-y-4">
                    <Label className="text-white/80 text-lg">Theme</Label>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                      {[
                        { value: 'light', label: 'Light', icon: Sun },
                        { value: 'dark', label: 'Dark', icon: Moon },
                        { value: 'auto', label: 'Auto', icon: Monitor },
                      ].map((theme) => {
                        const Icon = theme.icon
                        return (
                          <div
                            key={theme.value}
                            onClick={() => setPreferences(prev => ({ ...prev, theme: theme.value }))}
                            className={cn(
                              "p-4 rounded-xl border cursor-pointer transition-all duration-300",
                              preferences.theme === theme.value
                                ? "bg-gradient-to-br from-[#ff9000]/20 to-orange-600/20 border-[#ff9000]/40"
                                : "bg-white/[0.02] border-white/10 hover:border-white/20"
                            )}
                          >
                            <div className="flex items-center gap-3">
                              <div className={cn(
                                "w-10 h-10 rounded-lg flex items-center justify-center",
                                preferences.theme === theme.value
                                  ? "bg-[#ff9000]/20"
                                  : "bg-white/[0.05]"
                              )}>
                                <Icon className={cn(
                                  "h-5 w-5",
                                  preferences.theme === theme.value ? "text-[#ff9000]" : "text-white/60"
                                )} />
                              </div>
                              <div>
                                <h4 className="text-white font-semibold">{theme.label}</h4>
                                <p className="text-white/60 text-sm">
                                  {theme.value === 'auto' ? 'System preference' : `${theme.label} mode`}
                                </p>
                              </div>
                            </div>
                          </div>
                        )
                      })}
                    </div>
                  </div>
                </div>

                <div className="flex justify-end pt-6 border-t border-white/10 mt-8">
                  <Button
                    onClick={handleSavePreferences}
                    disabled={updatePreferences.isPending}
                    className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                  >
                    {updatePreferences.isPending ? (
                      <>
                        <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Save Settings
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </TabsContent>

            {/* Privacy Tab */}
            <TabsContent value="privacy" className="mt-8">
              <div className="bg-gradient-to-br from-white/[0.08] to-white/[0.02] backdrop-blur-xl border border-white/20 rounded-2xl p-8 shadow-xl shadow-black/10">
                <div className="flex items-center gap-3 mb-6">
                  <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-red-500/20 to-pink-500/20 flex items-center justify-center">
                    <Shield className="h-5 w-5 text-red-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-bold text-lg">Privacy Settings</h3>
                    <p className="text-white/60 text-sm">Control your privacy and data sharing</p>
                  </div>
                </div>

                <div className="space-y-6">
                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div>
                      <h4 className="text-white font-semibold">Marketing Emails</h4>
                      <p className="text-white/60 text-sm">Receive promotional emails and offers</p>
                    </div>
                    <Switch
                      checked={preferences.marketing_emails}
                      onCheckedChange={(checked) => setPreferences(prev => ({ ...prev, marketing_emails: checked }))}
                    />
                  </div>

                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div>
                      <h4 className="text-white font-semibold">Data Collection</h4>
                      <p className="text-white/60 text-sm">Allow BiHub to collect usage data for improvements</p>
                    </div>
                    <Switch defaultChecked />
                  </div>

                  <div className="flex items-center justify-between p-4 bg-white/[0.02] rounded-xl border border-white/10">
                    <div>
                      <h4 className="text-white font-semibold">Profile Visibility</h4>
                      <p className="text-white/60 text-sm">Control who can see your profile information</p>
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
                </div>

                <div className="flex justify-end pt-6 border-t border-white/10 mt-8">
                  <Button
                    onClick={handleSavePreferences}
                    disabled={updatePreferences.isPending}
                    className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white font-semibold px-6 py-2 rounded-xl"
                  >
                    {updatePreferences.isPending ? (
                      <>
                        <RefreshCw className="h-4 w-4 mr-2 animate-spin" />
                        Saving...
                      </>
                    ) : (
                      <>
                        <Save className="h-4 w-4 mr-2" />
                        Save Settings
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </div>
  )
}
