'use client';

import React, { useState } from 'react';
import {
  Bell,
  BellRing,
  Check,
  X,
  Settings,
  Trash2,
  Clock,
  ShoppingCart,
  CreditCard,
  Truck,
  Star,
  Gift,
  Shield,
  Users,
  Package,
  TrendingUp,
  AlertCircle,
  CheckCircle,
  Info,
  Zap,
  Wifi,
  WifiOff,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
  DropdownMenuItem,
} from '@/components/ui/dropdown-menu';
import { useWebSocketNotifications } from '@/hooks/use-websocket-notifications';
import { useAuthStore } from '@/store/auth';
import { formatDistanceToNow } from 'date-fns';
import { vi } from 'date-fns/locale';
import { cn } from '@/lib/utils';

interface ModernNotificationBellProps {
  className?: string;
}

export function ModernNotificationBell({ className }: ModernNotificationBellProps) {
  const { user } = useAuthStore();
  const {
    isConnected,
    connectionTime,
    lastMessage,
    notifications,
    unreadCount,
    error,
    markAsRead,
    clearNotifications,
    sendTestNotification,
  } = useWebSocketNotifications();

  const [isOpen, setIsOpen] = useState(false);

  // Get notification icon based on category and role
  const getNotificationIcon = (category: string, isAdmin: boolean = false) => {
    if (isAdmin) {
      const adminIcons: Record<string, React.ReactNode> = {
        order: <ShoppingCart className="h-4 w-4 text-blue-600" />,
        system: <AlertCircle className="h-4 w-4 text-orange-600" />,
        user: <Users className="h-4 w-4 text-green-600" />,
        review: <Star className="h-4 w-4 text-yellow-600" />,
        inventory: <Package className="h-4 w-4 text-purple-600" />,
        revenue: <TrendingUp className="h-4 w-4 text-emerald-600" />,
      };
      return adminIcons[category] || <Bell className="h-4 w-4 text-gray-600" />;
    }

    const customerIcons: Record<string, React.ReactNode> = {
      order: <ShoppingCart className="h-4 w-4 text-blue-600" />,
      payment: <CreditCard className="h-4 w-4 text-green-600" />,
      shipping: <Truck className="h-4 w-4 text-orange-600" />,
      promotion: <Gift className="h-4 w-4 text-pink-600" />,
      review: <Star className="h-4 w-4 text-yellow-600" />,
      security: <Shield className="h-4 w-4 text-red-600" />,
    };
    return customerIcons[category] || <Bell className="h-4 w-4 text-gray-600" />;
  };

  // Get priority color
  const getPriorityColor = (priority: string) => {
    const colors = {
      high: 'border-l-red-500 bg-red-50 hover:bg-red-100',
      normal: 'border-l-blue-500 bg-blue-50 hover:bg-blue-100',
      low: 'border-l-gray-500 bg-gray-50 hover:bg-gray-100',
    };
    return colors[priority as keyof typeof colors] || colors.normal;
  };

  // Get notification title based on role
  const getNotificationTitle = (isAdmin: boolean) => {
    if (isAdmin) {
      return (
        <div className="flex items-center gap-2">
          <TrendingUp className="h-4 w-4 text-blue-600" />
          <span className="font-semibold text-gray-900">Admin Dashboard</span>
          <Badge variant="secondary" className="text-xs bg-blue-100 text-blue-700">
            Business
          </Badge>
        </div>
      );
    }
    return (
      <div className="flex items-center gap-2">
        <Bell className="h-4 w-4 text-blue-600" />
        <span className="font-semibold text-gray-900">Notifications</span>
      </div>
    );
  };

  const isAdmin = user?.role === 'admin';

  return (
    <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className={cn(
            "relative h-10 w-10 rounded-full hover:bg-gray-100 transition-all duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2",
            className
          )}
        >
          {isConnected ? (
            <BellRing className={cn(
              "h-5 w-5 transition-colors",
              unreadCount > 0 ? "text-blue-600" : "text-gray-600"
            )} />
          ) : (
            <Bell className="h-5 w-5 text-gray-400" />
          )}
          
          {/* Unread count badge */}
          {unreadCount > 0 && (
            <Badge
              variant="destructive"
              className="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs font-bold animate-pulse"
            >
              {unreadCount > 99 ? '99+' : unreadCount}
            </Badge>
          )}
          
          {/* Connection status indicator */}
          <div
            className={cn(
              "absolute bottom-0 right-0 h-2.5 w-2.5 rounded-full border-2 border-white transition-colors",
              isConnected ? "bg-green-500" : "bg-red-500"
            )}
          />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent
        align="end"
        className="w-96 p-0 shadow-xl border-0 bg-white"
        sideOffset={8}
      >
        <Card className="border-0 shadow-none">
          <CardHeader className="pb-3 bg-gradient-to-r from-blue-50 to-indigo-50">
            <div className="flex items-center justify-between">
              {getNotificationTitle(isAdmin)}
              <div className="flex items-center gap-2">
                {/* Connection status */}
                <div
                  className={cn(
                    "flex items-center gap-1.5 text-xs px-2 py-1 rounded-full font-medium",
                    isConnected
                      ? "bg-green-100 text-green-700"
                      : "bg-red-100 text-red-700"
                  )}
                >
                  {isConnected ? (
                    <Wifi className="h-3 w-3" />
                  ) : (
                    <WifiOff className="h-3 w-3" />
                  )}
                  {isConnected ? "Live" : "Offline"}
                </div>
                
                {/* Settings menu */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="sm" className="h-7 w-7 p-0 hover:bg-white/50">
                      <Settings className="h-3.5 w-3.5" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-48">
                    <DropdownMenuItem 
                      onClick={() => clearNotifications()}
                      className="text-red-600 focus:text-red-600"
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      Clear All
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem>
                      <Settings className="h-4 w-4 mr-2" />
                      Preferences
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
            
            {/* Error state */}
            {error && (
              <div className="flex items-center gap-2 text-sm text-red-600 bg-red-100 p-2 rounded-lg mt-2">
                <AlertCircle className="h-4 w-4 flex-shrink-0" />
                <span className="text-xs">{error}</span>
              </div>
            )}
          </CardHeader>

          <CardContent className="p-0">
            {notifications.length === 0 ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <div className="relative mb-4">
                  <Bell className="h-16 w-16 text-gray-200" />
                  <div className="absolute -top-1 -right-1 h-4 w-4 bg-gray-100 rounded-full flex items-center justify-center">
                    <Check className="h-2.5 w-2.5 text-gray-400" />
                  </div>
                </div>
                <p className="text-sm text-gray-600 font-medium mb-1">All caught up!</p>
                <p className="text-xs text-gray-400">
                  {isAdmin ? "Business alerts will appear here" : "Your updates will appear here"}
                </p>
              </div>
            ) : (
              <ScrollArea className="h-96">
                <div className="space-y-0.5 p-2">
                  {notifications.map((notification, index) => (
                    <div
                      key={notification.id}
                      className={cn(
                        "p-3 rounded-lg border-l-4 cursor-pointer transition-all duration-200",
                        getPriorityColor(notification.priority),
                        !notification.is_read && "shadow-sm",
                        notification.is_read && "opacity-75"
                      )}
                      onClick={() => markAsRead(notification.id)}
                    >
                      <div className="flex items-start gap-3">
                        <div className="flex-shrink-0 mt-0.5">
                          {getNotificationIcon(notification.category, isAdmin)}
                        </div>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-start justify-between mb-1">
                            <h4 className="text-sm font-medium text-gray-900 line-clamp-1 pr-2">
                              {notification.title}
                            </h4>
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
                            <Badge
                              variant="outline"
                              className="text-xs capitalize border-gray-200 text-gray-600"
                            >
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

            {/* Footer */}
            {notifications.length > 0 && (
              <>
                <Separator />
                <div className="p-3 bg-gray-50">
                  <Button
                    variant="ghost"
                    size="sm"
                    className="w-full text-sm font-medium text-gray-700 hover:text-gray-900 hover:bg-white"
                    onClick={() => {
                      setIsOpen(false);
                      // Navigate to full notifications page
                      window.location.href = isAdmin ? '/admin/notifications' : '/notifications';
                    }}
                  >
                    View All Notifications
                  </Button>
                </div>
              </>
            )}
          </CardContent>
        </Card>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
