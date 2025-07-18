'use client';

import React, { useState, useEffect } from 'react';
import {
  Bell,
  Filter,
  Search,
  Calendar,
  TrendingUp,
  Users,
  ShoppingCart,
  Star,
  Package,
  AlertTriangle,
  CheckCircle,
  Clock,
  Eye,
  EyeOff,
  Trash2,
  Settings,
  Download,
  RefreshCw,
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

export default function AdminNotificationsPage() {
  const {
    isConnected,
    notifications,
    unreadCount,
    markAsRead,
    clearNotifications,
  } = useWebSocketNotifications();

  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [selectedPriority, setSelectedPriority] = useState('all');
  const [showUnreadOnly, setShowUnreadOnly] = useState(false);

  // Filter notifications
  const filteredNotifications = notifications.filter(notification => {
    const matchesSearch = notification.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         notification.message.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = selectedCategory === 'all' || notification.category === selectedCategory;
    const matchesPriority = selectedPriority === 'all' || notification.priority === selectedPriority;
    const matchesReadStatus = !showUnreadOnly || !notification.is_read;
    
    return matchesSearch && matchesCategory && matchesPriority && matchesReadStatus;
  });

  // Get notification stats
  const stats = {
    total: notifications.length,
    unread: unreadCount,
    high_priority: notifications.filter(n => n.priority === 'high').length,
    today: notifications.filter(n => {
      const today = new Date();
      const notificationDate = new Date(n.created_at);
      return notificationDate.toDateString() === today.toDateString();
    }).length,
  };

  // Get category counts
  const categoryStats = notifications.reduce((acc, notification) => {
    acc[notification.category] = (acc[notification.category] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  // Get notification icon
  const getNotificationIcon = (category: string) => {
    const icons: Record<string, React.ReactNode> = {
      order: <ShoppingCart className="h-4 w-4 text-blue-600" />,
      system: <AlertTriangle className="h-4 w-4 text-orange-600" />,
      user: <Users className="h-4 w-4 text-green-600" />,
      review: <Star className="h-4 w-4 text-yellow-600" />,
      inventory: <Package className="h-4 w-4 text-purple-600" />,
      revenue: <TrendingUp className="h-4 w-4 text-emerald-600" />,
    };
    return icons[category] || <Bell className="h-4 w-4 text-gray-600" />;
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
            <TrendingUp className="h-8 w-8 text-blue-600" />
            Admin Notifications
          </h1>
          <p className="text-gray-600 mt-1">
            Monitor business activities and system alerts in real-time
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
            <Download className="h-4 w-4 mr-2" />
            Export
          </Button>
          
          <Button variant="outline" size="sm">
            <Settings className="h-4 w-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Total Notifications</p>
                <p className="text-2xl font-bold text-gray-900">{stats.total}</p>
              </div>
              <Bell className="h-8 w-8 text-blue-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Unread</p>
                <p className="text-2xl font-bold text-orange-600">{stats.unread}</p>
              </div>
              <EyeOff className="h-8 w-8 text-orange-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">High Priority</p>
                <p className="text-2xl font-bold text-red-600">{stats.high_priority}</p>
              </div>
              <AlertTriangle className="h-8 w-8 text-red-600" />
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600">Today</p>
                <p className="text-2xl font-bold text-green-600">{stats.today}</p>
              </div>
              <Calendar className="h-8 w-8 text-green-600" />
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
                  placeholder="Search notifications..."
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
                <SelectItem value="order">Orders ({categoryStats.order || 0})</SelectItem>
                <SelectItem value="user">Users ({categoryStats.user || 0})</SelectItem>
                <SelectItem value="system">System ({categoryStats.system || 0})</SelectItem>
                <SelectItem value="review">Reviews ({categoryStats.review || 0})</SelectItem>
                <SelectItem value="inventory">Inventory ({categoryStats.inventory || 0})</SelectItem>
              </SelectContent>
            </Select>
            
            <Select value={selectedPriority} onValueChange={setSelectedPriority}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All</SelectItem>
                <SelectItem value="high">High</SelectItem>
                <SelectItem value="normal">Normal</SelectItem>
                <SelectItem value="low">Low</SelectItem>
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

      {/* Notifications List */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center justify-between">
            <span>Business Notifications ({filteredNotifications.length})</span>
            <Button variant="ghost" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          {filteredNotifications.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-center">
              <Bell className="h-16 w-16 text-gray-200 mb-4" />
              <p className="text-lg font-medium text-gray-600 mb-2">No notifications found</p>
              <p className="text-sm text-gray-400">
                {searchTerm || selectedCategory !== 'all' || selectedPriority !== 'all' || showUnreadOnly
                  ? "Try adjusting your filters"
                  : "Business alerts will appear here"}
              </p>
            </div>
          ) : (
            <ScrollArea className="h-96">
              <div className="space-y-1 p-4">
                {filteredNotifications.map((notification) => (
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
                      <div className="flex-shrink-0 mt-1">
                        {getNotificationIcon(notification.category)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between mb-2">
                          <h3 className="text-sm font-semibold text-gray-900 line-clamp-1 pr-4">
                            {notification.title}
                          </h3>
                          <div className="flex items-center gap-2 flex-shrink-0">
                            {!notification.is_read && (
                              <div className="h-2 w-2 bg-blue-600 rounded-full" />
                            )}
                            <Badge
                              variant={notification.priority === 'high' ? 'destructive' : 'outline'}
                              className="text-xs capitalize"
                            >
                              {notification.priority}
                            </Badge>
                          </div>
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
                          <Badge variant="outline" className="text-xs capitalize">
                            {notification.category}
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </ScrollArea>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
