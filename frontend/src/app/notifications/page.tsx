'use client';

import React, { useState, useEffect } from 'react';
import {
  Bell,
  Search,
  ShoppingCart,
  CreditCard,
  Truck,
  Star,
  Gift,
  Shield,
  Clock,
  Eye,
  EyeOff,
  Trash2,
  Settings,
  CheckCircle,
  AlertCircle,
  Package,
  Heart,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useWebSocketNotifications } from '@/hooks/use-websocket-notifications';
import { formatDistanceToNow } from 'date-fns';
import { vi } from 'date-fns/locale';
import { cn } from '@/lib/utils';

export default function CustomerNotificationsPage() {
  const {
    isConnected,
    notifications,
    unreadCount,
    markAsRead,
    clearNotifications,
  } = useWebSocketNotifications();

  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [showUnreadOnly, setShowUnreadOnly] = useState(false);

  // Filter notifications
  const filteredNotifications = notifications.filter(notification => {
    const matchesSearch = notification.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         notification.message.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = selectedCategory === 'all' || notification.category === selectedCategory;
    const matchesReadStatus = !showUnreadOnly || !notification.is_read;
    
    return matchesSearch && matchesCategory && matchesReadStatus;
  });

  // Group notifications by category
  const groupedNotifications = filteredNotifications.reduce((acc, notification) => {
    const category = notification.category;
    if (!acc[category]) {
      acc[category] = [];
    }
    acc[category].push(notification);
    return acc;
  }, {} as Record<string, typeof notifications>);

  // Get notification stats
  const stats = {
    total: notifications.length,
    unread: unreadCount,
    orders: notifications.filter(n => n.category === 'order').length,
    promotions: notifications.filter(n => n.category === 'promotion').length,
  };

  // Get notification icon and color
  const getCategoryInfo = (category: string) => {
    const categoryInfo: Record<string, { icon: React.ReactNode; color: string; label: string }> = {
      order: {
        icon: <ShoppingCart className="h-5 w-5" />,
        color: 'text-blue-600 bg-blue-100',
        label: 'Orders'
      },
      payment: {
        icon: <CreditCard className="h-5 w-5" />,
        color: 'text-green-600 bg-green-100',
        label: 'Payments'
      },
      shipping: {
        icon: <Truck className="h-5 w-5" />,
        color: 'text-orange-600 bg-orange-100',
        label: 'Shipping'
      },
      promotion: {
        icon: <Gift className="h-5 w-5" />,
        color: 'text-pink-600 bg-pink-100',
        label: 'Promotions'
      },
      review: {
        icon: <Star className="h-5 w-5" />,
        color: 'text-yellow-600 bg-yellow-100',
        label: 'Reviews'
      },
      security: {
        icon: <Shield className="h-5 w-5" />,
        color: 'text-red-600 bg-red-100',
        label: 'Security'
      },
    };
    return categoryInfo[category] || {
      icon: <Bell className="h-5 w-5" />,
      color: 'text-gray-600 bg-gray-100',
      label: 'Other'
    };
  };

  // Get priority color
  const getPriorityColor = (priority: string) => {
    const colors = {
      high: 'border-l-red-500 bg-red-50',
      normal: 'border-l-blue-500 bg-blue-50',
      low: 'border-l-gray-500 bg-gray-50',
    };
    return colors[priority as keyof typeof colors] || colors.normal;
  };

  return (
    <div className="container mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 flex items-center gap-3">
            <Bell className="h-8 w-8 text-blue-600" />
            Your Notifications
          </h1>
          <p className="text-gray-600 mt-1">
            Stay updated with your orders, promotions, and account activities
          </p>
        </div>
        
        <div className="flex items-center gap-3">
          <div className={cn(
            "flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium",
            isConnected ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
          )}>
            <div className={cn(
              "h-2 w-2 rounded-full",
              isConnected ? "bg-green-500" : "bg-red-500"
            )} />
            {isConnected ? "Live Updates" : "Disconnected"}
          </div>
          
          <Button variant="outline" size="sm">
            <Settings className="h-4 w-4 mr-2" />
            Preferences
          </Button>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card className="bg-gradient-to-r from-blue-50 to-blue-100 border-blue-200">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-blue-700">Total</p>
                <p className="text-2xl font-bold text-blue-900">{stats.total}</p>
              </div>
              <Bell className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-gradient-to-r from-orange-50 to-orange-100 border-orange-200">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-orange-700">Unread</p>
                <p className="text-2xl font-bold text-orange-900">{stats.unread}</p>
              </div>
              <EyeOff className="h-8 w-8 text-orange-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-gradient-to-r from-green-50 to-green-100 border-green-200">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-green-700">Orders</p>
                <p className="text-2xl font-bold text-green-900">{stats.orders}</p>
              </div>
              <ShoppingCart className="h-8 w-8 text-green-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-gradient-to-r from-pink-50 to-pink-100 border-pink-200">
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-pink-700">Promotions</p>
                <p className="text-2xl font-bold text-pink-900">{stats.promotions}</p>
              </div>
              <Gift className="h-8 w-8 text-pink-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card>
        <CardContent className="p-4">
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search your notifications..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10"
                />
              </div>
            </div>
            
            <Select value={selectedCategory} onValueChange={setSelectedCategory}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                <SelectItem value="order">🛒 Orders</SelectItem>
                <SelectItem value="payment">💳 Payments</SelectItem>
                <SelectItem value="shipping">🚚 Shipping</SelectItem>
                <SelectItem value="promotion">🎁 Promotions</SelectItem>
                <SelectItem value="review">⭐ Reviews</SelectItem>
                <SelectItem value="security">🔒 Security</SelectItem>
              </SelectContent>
            </Select>
            
            <Button
              variant={showUnreadOnly ? "default" : "outline"}
              onClick={() => setShowUnreadOnly(!showUnreadOnly)}
              className="whitespace-nowrap"
            >
              <EyeOff className="h-4 w-4 mr-2" />
              Unread Only
            </Button>
            
            <Button
              variant="outline"
              onClick={() => clearNotifications()}
              className="text-red-600 hover:text-red-700"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Clear All
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Notifications */}
      <div className="space-y-6">
        {Object.keys(groupedNotifications).length === 0 ? (
          <Card>
            <CardContent className="p-12">
              <div className="flex flex-col items-center justify-center text-center">
                <div className="relative mb-6">
                  <Bell className="h-20 w-20 text-gray-200" />
                  <CheckCircle className="absolute -bottom-2 -right-2 h-8 w-8 text-green-500 bg-white rounded-full" />
                </div>
                <h3 className="text-xl font-semibold text-gray-900 mb-2">All caught up!</h3>
                <p className="text-gray-600 mb-4">
                  {searchTerm || selectedCategory !== 'all' || showUnreadOnly
                    ? "No notifications match your current filters"
                    : "You're all up to date. New notifications will appear here."}
                </p>
                {(searchTerm || selectedCategory !== 'all' || showUnreadOnly) && (
                  <Button
                    variant="outline"
                    onClick={() => {
                      setSearchTerm('');
                      setSelectedCategory('all');
                      setShowUnreadOnly(false);
                    }}
                  >
                    Clear Filters
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ) : (
          Object.entries(groupedNotifications).map(([category, categoryNotifications]) => {
            const categoryInfo = getCategoryInfo(category);
            return (
              <Card key={category}>
                <CardHeader className="pb-3">
                  <CardTitle className="flex items-center gap-3">
                    <div className={cn("p-2 rounded-lg", categoryInfo.color)}>
                      {categoryInfo.icon}
                    </div>
                    <span>{categoryInfo.label}</span>
                    <Badge variant="secondary" className="ml-auto">
                      {categoryNotifications.length}
                    </Badge>
                  </CardTitle>
                </CardHeader>
                <CardContent className="p-0">
                  <div className="space-y-1 p-4">
                    {categoryNotifications.map((notification) => (
                      <div
                        key={notification.id}
                        className={cn(
                          "p-4 rounded-lg border-l-4 cursor-pointer transition-all hover:shadow-md",
                          getPriorityColor(notification.priority),
                          !notification.is_read && "shadow-sm",
                          notification.is_read && "opacity-75"
                        )}
                        onClick={() => markAsRead(notification.id)}
                      >
                        <div className="flex items-start gap-4">
                          <div className="flex-1 min-w-0">
                            <div className="flex items-start justify-between mb-2">
                              <h3 className="text-sm font-semibold text-gray-900 line-clamp-1 pr-4">
                                {notification.title}
                              </h3>
                              {!notification.is_read && (
                                <div className="h-2 w-2 bg-blue-600 rounded-full flex-shrink-0 mt-1" />
                              )}
                            </div>
                            <p className="text-sm text-gray-600 line-clamp-2 mb-3 leading-relaxed">
                              {notification.message}
                            </p>
                            <div className="flex items-center justify-between">
                              <span className="text-xs text-gray-400 flex items-center gap-1">
                                <Clock className="h-3 w-3" />
                                {formatDistanceToNow(new Date(notification.created_at), {
                                  addSuffix: true,
                                  locale: vi,
                                })}
                              </span>
                              {notification.priority === 'high' && (
                                <Badge variant="destructive" className="text-xs">
                                  High Priority
                                </Badge>
                              )}
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            );
          })
        )}
      </div>
    </div>
  );
}
