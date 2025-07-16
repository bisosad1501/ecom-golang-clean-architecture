'use client';

import React from 'react';
import { Bell, Mail, MessageSquare, Smartphone, Save } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Separator } from '@/components/ui/separator';
import { useNotificationPreferences } from '@/hooks/use-notifications';
import { UpdateNotificationPreferencesRequest } from '@/types/notification';

export function NotificationSettingsPage() {
  const { preferences, loading, updatePreferences } = useNotificationPreferences();

  const handleToggle = async (key: keyof UpdateNotificationPreferencesRequest, value: boolean) => {
    await updatePreferences({ [key]: value });
  };

  if (loading && !preferences) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="animate-pulse">
          <div className="h-8 bg-gray-200 rounded w-1/3 mb-4"></div>
          <div className="h-4 bg-gray-200 rounded w-2/3 mb-8"></div>
          <div className="space-y-4">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-32 bg-gray-200 rounded"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="flex items-center gap-3 mb-8">
        <div className="p-2 bg-orange-100 rounded-lg">
          <Bell className="h-6 w-6 text-orange-600" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
            Cài đặt thông báo
          </h1>
          <p className="text-gray-600 dark:text-gray-400">
            Quản lý cách bạn nhận thông báo từ Bihub
          </p>
        </div>
      </div>

      <div className="space-y-6">
        {/* General Notification Settings */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Bell className="h-5 w-5" />
              Cài đặt chung
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Thông báo trong ứng dụng</h4>
                <p className="text-sm text-gray-500">Nhận thông báo khi đang sử dụng ứng dụng</p>
              </div>
              <Switch
                checked={preferences?.in_app_enabled || false}
                onCheckedChange={(checked) => handleToggle('in_app_enabled', checked)}
                disabled={loading}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Thông báo email</h4>
                <p className="text-sm text-gray-500">Nhận thông báo qua email</p>
              </div>
              <Switch
                checked={preferences?.email_enabled || false}
                onCheckedChange={(checked) => handleToggle('email_enabled', checked)}
                disabled={loading}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Thông báo SMS</h4>
                <p className="text-sm text-gray-500">Nhận thông báo qua tin nhắn SMS</p>
              </div>
              <Switch
                checked={preferences?.sms_enabled || false}
                onCheckedChange={(checked) => handleToggle('sms_enabled', checked)}
                disabled={loading}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Thông báo đẩy</h4>
                <p className="text-sm text-gray-500">Nhận thông báo đẩy trên thiết bị</p>
              </div>
              <Switch
                checked={preferences?.push_enabled || false}
                onCheckedChange={(checked) => handleToggle('push_enabled', checked)}
                disabled={loading}
              />
            </div>
          </CardContent>
        </Card>

        {/* Email Notifications */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Mail className="h-5 w-5" />
              Thông báo email
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cập nhật đơn hàng</h4>
                <p className="text-sm text-gray-500">Thông báo về trạng thái đơn hàng</p>
              </div>
              <Switch
                checked={preferences?.email_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('email_order_updates', checked)}
                disabled={loading || !preferences?.email_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cập nhật thanh toán</h4>
                <p className="text-sm text-gray-500">Thông báo về giao dịch thanh toán</p>
              </div>
              <Switch
                checked={preferences?.email_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('email_payment_updates', checked)}
                disabled={loading || !preferences?.email_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cập nhật vận chuyển</h4>
                <p className="text-sm text-gray-500">Thông báo về tình trạng giao hàng</p>
              </div>
              <Switch
                checked={preferences?.email_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('email_shipping_updates', checked)}
                disabled={loading || !preferences?.email_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Khuyến mãi và ưu đãi</h4>
                <p className="text-sm text-gray-500">Thông báo về các chương trình khuyến mãi</p>
              </div>
              <Switch
                checked={preferences?.email_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('email_promotions', checked)}
                disabled={loading || !preferences?.email_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Bản tin</h4>
                <p className="text-sm text-gray-500">Nhận bản tin định kỳ từ Bihub</p>
              </div>
              <Switch
                checked={preferences?.email_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('email_newsletter', checked)}
                disabled={loading || !preferences?.email_enabled}
              />
            </div>
          </CardContent>
        </Card>

        {/* SMS Notifications */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <MessageSquare className="h-5 w-5" />
              Thông báo SMS
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cập nhật đơn hàng</h4>
                <p className="text-sm text-gray-500">SMS về trạng thái đơn hàng quan trọng</p>
              </div>
              <Switch
                checked={preferences?.sms_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('sms_order_updates', checked)}
                disabled={loading || !preferences?.sms_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cảnh báo bảo mật</h4>
                <p className="text-sm text-gray-500">SMS về các hoạt động bảo mật</p>
              </div>
              <Switch
                checked={preferences?.sms_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('sms_security_alerts', checked)}
                disabled={loading || !preferences?.sms_enabled}
              />
            </div>
          </CardContent>
        </Card>

        {/* Push Notifications */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Smartphone className="h-5 w-5" />
              Thông báo đẩy
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Cập nhật đơn hàng</h4>
                <p className="text-sm text-gray-500">Thông báo đẩy về đơn hàng</p>
              </div>
              <Switch
                checked={preferences?.push_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('push_order_updates', checked)}
                disabled={loading || !preferences?.push_enabled}
              />
            </div>

            <Separator />

            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium">Khuyến mãi</h4>
                <p className="text-sm text-gray-500">Thông báo đẩy về ưu đãi đặc biệt</p>
              </div>
              <Switch
                checked={preferences?.push_enabled && true} // Simplified for demo
                onCheckedChange={(checked) => handleToggle('push_promotions', checked)}
                disabled={loading || !preferences?.push_enabled}
              />
            </div>
          </CardContent>
        </Card>

        {/* Save Button */}
        <div className="flex justify-end">
          <Button
            onClick={() => window.location.href = '/notifications'}
            className="gap-2"
          >
            <Save className="h-4 w-4" />
            Quay lại thông báo
          </Button>
        </div>
      </div>
    </div>
  );
}
