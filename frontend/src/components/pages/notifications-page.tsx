'use client';

import React, { useState, useEffect } from 'react';
import { Bell, Filter, CheckCheck, Trash2, Settings } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { useNotifications } from '@/hooks/use-notifications';
import { NotificationItem } from '@/components/notifications/notification-item';
import { NotificationFilters, NotificationType, NotificationCategory } from '@/types/notification';

export function NotificationsPage() {
  const [activeTab, setActiveTab] = useState('all');
  const [filters, setFilters] = useState<NotificationFilters>({
    limit: 20,
    offset: 0,
  });

  const { 
    notifications, 
    unreadCount, 
    loading, 
    markAsRead, 
    markAllAsRead,
    fetchNotifications 
  } = useNotifications();

  // Fetch notifications on mount and when filters change
  useEffect(() => {
    const currentFilters = { ...filters };
    
    // Apply tab filters
    if (activeTab === 'unread') {
      currentFilters.is_read = false;
    } else if (activeTab === 'read') {
      currentFilters.is_read = true;
    }
    
    fetchNotifications(currentFilters);
  }, [activeTab, filters, fetchNotifications]);

  const handleTabChange = (value: string) => {
    setActiveTab(value);
    setFilters(prev => ({ ...prev, offset: 0 })); // Reset pagination
  };

  const handleFilterChange = (key: keyof NotificationFilters, value: any) => {
    setFilters(prev => ({
      ...prev,
      [key]: value,
      offset: 0, // Reset pagination when filters change
    }));
  };

  const handleLoadMore = () => {
    setFilters(prev => ({
      ...prev,
      offset: (prev.offset || 0) + (prev.limit || 20),
    }));
  };

  const getTabCount = (tab: string) => {
    switch (tab) {
      case 'unread':
        return unreadCount;
      case 'read':
        return notifications.filter(n => n.is_read).length;
      default:
        return notifications.length;
    }
  };

  const filteredNotifications = notifications.filter(notification => {
    if (activeTab === 'unread' && notification.is_read) return false;
    if (activeTab === 'read' && !notification.is_read) return false;
    return true;
  });

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="flex items-center justify-between mb-8">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-orange-100 rounded-lg">
            <Bell className="h-6 w-6 text-orange-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              Thông báo
            </h1>
            <p className="text-gray-600 dark:text-gray-400">
              Quản lý và xem tất cả thông báo của bạn
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {unreadCount > 0 && (
            <Button
              variant="outline"
              size="sm"
              onClick={markAllAsRead}
              className="gap-2"
            >
              <CheckCheck className="h-4 w-4" />
              Đánh dấu tất cả đã đọc
            </Button>
          )}
          
          <Button
            variant="outline"
            size="sm"
            onClick={() => window.location.href = '/settings/notifications'}
            className="gap-2"
          >
            <Settings className="h-4 w-4" />
            Cài đặt
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card className="mb-6">
        <CardHeader className="pb-4">
          <CardTitle className="text-lg flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Bộ lọc
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-4">
            <div className="flex-1 min-w-[200px]">
              <label className="text-sm font-medium mb-2 block">Loại thông báo</label>
              <Select
                value={filters.type || 'all'}
                onValueChange={(value) => 
                  handleFilterChange('type', value === 'all' ? undefined : value as NotificationType)
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Tất cả loại" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Tất cả loại</SelectItem>
                  <SelectItem value="in_app">Trong ứng dụng</SelectItem>
                  <SelectItem value="email">Email</SelectItem>
                  <SelectItem value="sms">SMS</SelectItem>
                  <SelectItem value="push">Push</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex-1 min-w-[200px]">
              <label className="text-sm font-medium mb-2 block">Danh mục</label>
              <Select
                value={filters.category || 'all'}
                onValueChange={(value) => 
                  handleFilterChange('category', value === 'all' ? undefined : value as NotificationCategory)
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Tất cả danh mục" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Tất cả danh mục</SelectItem>
                  <SelectItem value="order">Đơn hàng</SelectItem>
                  <SelectItem value="payment">Thanh toán</SelectItem>
                  <SelectItem value="shipping">Vận chuyển</SelectItem>
                  <SelectItem value="review">Đánh giá</SelectItem>
                  <SelectItem value="system">Hệ thống</SelectItem>
                  <SelectItem value="promotion">Khuyến mãi</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="flex-1 min-w-[200px]">
              <label className="text-sm font-medium mb-2 block">Mức độ ưu tiên</label>
              <Select
                value={filters.priority || 'all'}
                onValueChange={(value) => 
                  handleFilterChange('priority', value === 'all' ? undefined : value)
                }
              >
                <SelectTrigger>
                  <SelectValue placeholder="Tất cả mức độ" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">Tất cả mức độ</SelectItem>
                  <SelectItem value="high">Cao</SelectItem>
                  <SelectItem value="normal">Bình thường</SelectItem>
                  <SelectItem value="low">Thấp</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={handleTabChange} className="mb-6">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="all" className="gap-2">
            Tất cả
            <Badge variant="secondary" className="ml-1">
              {getTabCount('all')}
            </Badge>
          </TabsTrigger>
          <TabsTrigger value="unread" className="gap-2">
            Chưa đọc
            <Badge variant="destructive" className="ml-1">
              {getTabCount('unread')}
            </Badge>
          </TabsTrigger>
          <TabsTrigger value="read" className="gap-2">
            Đã đọc
            <Badge variant="secondary" className="ml-1">
              {getTabCount('read')}
            </Badge>
          </TabsTrigger>
        </TabsList>

        <TabsContent value={activeTab} className="mt-6">
          <Card>
            <CardContent className="p-0">
              {loading && notifications.length === 0 ? (
                <div className="p-8 text-center">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-orange-500 mx-auto mb-4"></div>
                  <p className="text-gray-500">Đang tải thông báo...</p>
                </div>
              ) : filteredNotifications.length === 0 ? (
                <div className="p-8 text-center">
                  <Bell className="h-16 w-16 text-gray-300 mx-auto mb-4" />
                  <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100 mb-2">
                    Không có thông báo
                  </h3>
                  <p className="text-gray-500">
                    {activeTab === 'unread' 
                      ? 'Bạn đã đọc hết tất cả thông báo' 
                      : 'Chưa có thông báo nào trong danh mục này'
                    }
                  </p>
                </div>
              ) : (
                <div className="divide-y divide-gray-100 dark:divide-gray-800">
                  {filteredNotifications.map((notification) => (
                    <NotificationItem
                      key={notification.id}
                      notification={notification}
                      onMarkAsRead={markAsRead}
                    />
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Load More Button */}
          {filteredNotifications.length > 0 && filteredNotifications.length >= (filters.limit || 20) && (
            <div className="mt-6 text-center">
              <Button
                variant="outline"
                onClick={handleLoadMore}
                disabled={loading}
                className="gap-2"
              >
                {loading ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
                    Đang tải...
                  </>
                ) : (
                  'Tải thêm'
                )}
              </Button>
            </div>
          )}
        </TabsContent>
      </Tabs>
    </div>
  );
}
