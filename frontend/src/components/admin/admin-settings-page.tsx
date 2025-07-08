'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '@/components/ui/tabs'
import {
  Settings,
  Store,
  Mail,
  Shield,
  Palette,
  Globe,
  CreditCard,
  Bell,
  Database,
  Key,
  Users,
  Package,
  Save,
  Upload,
  Download,
  RefreshCw,
} from 'lucide-react'
import { RequirePermission } from '@/components/auth/require-permission'
import { PERMISSIONS } from '@/constants/permissions'
import { toast } from 'sonner'

// BiHub Components
import {
  BiHubAdminCard,
  BiHubPageHeader,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

export default function AdminSettingsPage() {
  const [isLoading, setIsLoading] = useState(false)
  const [activeTab, setActiveTab] = useState('general')

  // Mock settings data - replace with real API calls
  const [settings, setSettings] = useState({
    general: {
      storeName: 'BiHub',
      storeDescription: 'Your premium e-commerce destination',
      storeEmail: 'admin@bihub.com',
      storePhone: '+1 (555) 123-4567',
      storeAddress: '123 Commerce Street, Business City, BC 12345',
      currency: 'USD',
      timezone: 'America/New_York',
      language: 'en',
    },
    appearance: {
      primaryColor: '#FF9000',
      secondaryColor: '#1F2937',
      logoUrl: '',
      faviconUrl: '',
      darkMode: true,
      customCSS: '',
    },
    email: {
      smtpHost: 'smtp.gmail.com',
      smtpPort: '587',
      smtpUsername: '',
      smtpPassword: '',
      fromEmail: 'noreply@bihub.com',
      fromName: 'BiHub Store',
    },
    notifications: {
      orderNotifications: true,
      customerNotifications: true,
      inventoryAlerts: true,
      securityAlerts: true,
      marketingEmails: false,
    },
    security: {
      twoFactorAuth: false,
      sessionTimeout: '24',
      passwordPolicy: 'medium',
      loginAttempts: '5',
    },
    payment: {
      stripeEnabled: true,
      paypalEnabled: false,
      stripePublicKey: '',
      stripeSecretKey: '',
      paypalClientId: '',
      paypalClientSecret: '',
    }
  })

  const handleSave = async (section: string) => {
    setIsLoading(true)
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000))
      toast.success(`${section} settings saved successfully!`)
    } catch (error) {
      toast.error('Failed to save settings')
    } finally {
      setIsLoading(false)
    }
  }

  const handleInputChange = (section: string, field: string, value: any) => {
    setSettings(prev => ({
      ...prev,
      [section]: {
        ...prev[section as keyof typeof prev],
        [field]: value
      }
    }))
  }

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Store Settings"
        subtitle="Configure your BiHub store settings and preferences"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Settings' }
        ]}
      />

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-6 bg-gray-800 border-gray-700">
          <TabsTrigger value="general" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <Store className="h-4 w-4 mr-2" />
            General
          </TabsTrigger>
          <TabsTrigger value="appearance" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <Palette className="h-4 w-4 mr-2" />
            Appearance
          </TabsTrigger>
          <TabsTrigger value="email" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <Mail className="h-4 w-4 mr-2" />
            Email
          </TabsTrigger>
          <TabsTrigger value="notifications" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <Bell className="h-4 w-4 mr-2" />
            Notifications
          </TabsTrigger>
          <TabsTrigger value="security" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <Shield className="h-4 w-4 mr-2" />
            Security
          </TabsTrigger>
          <TabsTrigger value="payment" className="data-[state=active]:bg-[#FF9000] data-[state=active]:text-white">
            <CreditCard className="h-4 w-4 mr-2" />
            Payment
          </TabsTrigger>
        </TabsList>

        {/* General Settings */}
        <TabsContent value="general">
          <BiHubAdminCard
            title="General Store Settings"
            subtitle="Basic information about your BiHub store"
            icon={<Store className="h-5 w-5 text-white" />}
            headerAction={
              <Button
                onClick={() => handleSave('general')}
                disabled={isLoading}
                className={BIHUB_ADMIN_THEME.components.button.primary}
              >
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </Button>
            }
          >
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="storeName" className="text-gray-300">Store Name</Label>
                <Input
                  id="storeName"
                  value={settings.general.storeName}
                  onChange={(e) => handleInputChange('general', 'storeName', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>
              
              <div className="space-y-2">
                <Label htmlFor="storeEmail" className="text-gray-300">Store Email</Label>
                <Input
                  id="storeEmail"
                  type="email"
                  value={settings.general.storeEmail}
                  onChange={(e) => handleInputChange('general', 'storeEmail', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="storePhone" className="text-gray-300">Store Phone</Label>
                <Input
                  id="storePhone"
                  value={settings.general.storePhone}
                  onChange={(e) => handleInputChange('general', 'storePhone', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="currency" className="text-gray-300">Currency</Label>
                <Select
                  value={settings.general.currency}
                  onValueChange={(value) => handleInputChange('general', 'currency', value)}
                >
                  <SelectTrigger className={BIHUB_ADMIN_THEME.components.input.base}>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-900 border-gray-700">
                    <SelectItem value="USD">USD - US Dollar</SelectItem>
                    <SelectItem value="EUR">EUR - Euro</SelectItem>
                    <SelectItem value="GBP">GBP - British Pound</SelectItem>
                    <SelectItem value="JPY">JPY - Japanese Yen</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="md:col-span-2 space-y-2">
                <Label htmlFor="storeDescription" className="text-gray-300">Store Description</Label>
                <Textarea
                  id="storeDescription"
                  value={settings.general.storeDescription}
                  onChange={(e) => handleInputChange('general', 'storeDescription', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                  rows={3}
                />
              </div>

              <div className="md:col-span-2 space-y-2">
                <Label htmlFor="storeAddress" className="text-gray-300">Store Address</Label>
                <Textarea
                  id="storeAddress"
                  value={settings.general.storeAddress}
                  onChange={(e) => handleInputChange('general', 'storeAddress', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                  rows={2}
                />
              </div>
            </div>
          </BiHubAdminCard>
        </TabsContent>

        {/* Appearance Settings */}
        <TabsContent value="appearance">
          <BiHubAdminCard
            title="Appearance Settings"
            subtitle="Customize your BiHub store's look and feel"
            icon={<Palette className="h-5 w-5 text-white" />}
            headerAction={
              <Button
                onClick={() => handleSave('appearance')}
                disabled={isLoading}
                className={BIHUB_ADMIN_THEME.components.button.primary}
              >
                <Save className="h-4 w-4 mr-2" />
                Save Changes
              </Button>
            }
          >
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="primaryColor" className="text-gray-300">Primary Color</Label>
                <div className="flex items-center gap-3">
                  <Input
                    id="primaryColor"
                    type="color"
                    value={settings.appearance.primaryColor}
                    onChange={(e) => handleInputChange('appearance', 'primaryColor', e.target.value)}
                    className="w-16 h-12 p-1 border-gray-600"
                  />
                  <Input
                    value={settings.appearance.primaryColor}
                    onChange={(e) => handleInputChange('appearance', 'primaryColor', e.target.value)}
                    className={cn(BIHUB_ADMIN_THEME.components.input.base, 'flex-1')}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="secondaryColor" className="text-gray-300">Secondary Color</Label>
                <div className="flex items-center gap-3">
                  <Input
                    id="secondaryColor"
                    type="color"
                    value={settings.appearance.secondaryColor}
                    onChange={(e) => handleInputChange('appearance', 'secondaryColor', e.target.value)}
                    className="w-16 h-12 p-1 border-gray-600"
                  />
                  <Input
                    value={settings.appearance.secondaryColor}
                    onChange={(e) => handleInputChange('appearance', 'secondaryColor', e.target.value)}
                    className={cn(BIHUB_ADMIN_THEME.components.input.base, 'flex-1')}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="logoUrl" className="text-gray-300">Logo URL</Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="logoUrl"
                    value={settings.appearance.logoUrl}
                    onChange={(e) => handleInputChange('appearance', 'logoUrl', e.target.value)}
                    className={cn(BIHUB_ADMIN_THEME.components.input.base, 'flex-1')}
                    placeholder="https://example.com/logo.png"
                  />
                  <Button variant="outline" size="sm" className={BIHUB_ADMIN_THEME.components.button.ghost}>
                    <Upload className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="faviconUrl" className="text-gray-300">Favicon URL</Label>
                <div className="flex items-center gap-2">
                  <Input
                    id="faviconUrl"
                    value={settings.appearance.faviconUrl}
                    onChange={(e) => handleInputChange('appearance', 'faviconUrl', e.target.value)}
                    className={cn(BIHUB_ADMIN_THEME.components.input.base, 'flex-1')}
                    placeholder="https://example.com/favicon.ico"
                  />
                  <Button variant="outline" size="sm" className={BIHUB_ADMIN_THEME.components.button.ghost}>
                    <Upload className="h-4 w-4" />
                  </Button>
                </div>
              </div>

              <div className="md:col-span-2 space-y-2">
                <Label htmlFor="customCSS" className="text-gray-300">Custom CSS</Label>
                <Textarea
                  id="customCSS"
                  value={settings.appearance.customCSS}
                  onChange={(e) => handleInputChange('appearance', 'customCSS', e.target.value)}
                  className={cn(BIHUB_ADMIN_THEME.components.input.base, 'font-mono')}
                  rows={6}
                  placeholder="/* Add your custom CSS here */"
                />
              </div>
            </div>
          </BiHubAdminCard>
        </TabsContent>

        {/* Email Settings */}
        <TabsContent value="email">
          <BiHubAdminCard
            title="Email Configuration"
            subtitle="Configure SMTP settings for sending emails"
            icon={<Mail className="h-5 w-5 text-white" />}
            headerAction={
              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  className={BIHUB_ADMIN_THEME.components.button.ghost}
                >
                  <RefreshCw className="h-4 w-4 mr-2" />
                  Test Connection
                </Button>
                <Button
                  onClick={() => handleSave('email')}
                  disabled={isLoading}
                  className={BIHUB_ADMIN_THEME.components.button.primary}
                >
                  <Save className="h-4 w-4 mr-2" />
                  Save Changes
                </Button>
              </div>
            }
          >
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="smtpHost" className="text-gray-300">SMTP Host</Label>
                <Input
                  id="smtpHost"
                  value={settings.email.smtpHost}
                  onChange={(e) => handleInputChange('email', 'smtpHost', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="smtpPort" className="text-gray-300">SMTP Port</Label>
                <Input
                  id="smtpPort"
                  value={settings.email.smtpPort}
                  onChange={(e) => handleInputChange('email', 'smtpPort', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="smtpUsername" className="text-gray-300">SMTP Username</Label>
                <Input
                  id="smtpUsername"
                  value={settings.email.smtpUsername}
                  onChange={(e) => handleInputChange('email', 'smtpUsername', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="smtpPassword" className="text-gray-300">SMTP Password</Label>
                <Input
                  id="smtpPassword"
                  type="password"
                  value={settings.email.smtpPassword}
                  onChange={(e) => handleInputChange('email', 'smtpPassword', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="fromEmail" className="text-gray-300">From Email</Label>
                <Input
                  id="fromEmail"
                  type="email"
                  value={settings.email.fromEmail}
                  onChange={(e) => handleInputChange('email', 'fromEmail', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="fromName" className="text-gray-300">From Name</Label>
                <Input
                  id="fromName"
                  value={settings.email.fromName}
                  onChange={(e) => handleInputChange('email', 'fromName', e.target.value)}
                  className={BIHUB_ADMIN_THEME.components.input.base}
                />
              </div>
            </div>
          </BiHubAdminCard>
        </TabsContent>

        {/* Add other tab contents here... */}
      </Tabs>
    </div>
  )
}
