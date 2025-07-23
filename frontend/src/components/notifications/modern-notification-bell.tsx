'use client';

import React from 'react';
import {
  Bell,
  BellRing,
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
  Wifi,
  WifiOff,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { AnimatedBackground } from '@/components/ui/animated-background';
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
    notifications,
    unreadCount,
    markAsRead,
  } = useWebSocketNotifications();

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

  // Get priority color - Updated for dark theme
  const getPriorityColor = (priority: string) => {
    const colors = {
      high: 'border-l-red-500 bg-red-900/20 hover:bg-red-800/30',
      normal: 'border-l-blue-500 bg-blue-900/20 hover:bg-blue-800/30',
      low: 'border-l-gray-500 bg-gray-800/20 hover:bg-gray-700/30',
    };
    return colors[priority as keyof typeof colors] || colors.normal;
  };

  // Get notification title based on role - Updated for dark theme
  const getNotificationTitle = (isAdmin: boolean) => {
    if (isAdmin) {
      return (
        <div className="flex items-center gap-2">
          <TrendingUp className="h-4 w-4 text-blue-400" />
          <span className="font-semibold text-white">Admin Dashboard</span>
          <Badge variant="secondary" className="text-xs bg-blue-900/30 text-blue-300 border-blue-700">
            Business
          </Badge>
        </div>
      );
    }
    return (
      <div className="flex items-center gap-2">
        <Bell className="h-4 w-4 text-blue-400" />
        <span className="font-semibold text-white">Notifications</span>
      </div>
    );
  };

  const isAdmin = user?.role === 'admin';

  return (
    <div className="relative group">
      <Button
        variant="ghost"
        size="icon"
        className={cn(
          "relative h-10 w-10 rounded-xl hover:bg-orange-500/10 hover:scale-105 transition-all duration-200 text-white",
          className
        )}
        onClick={() => {
          // Navigate to notifications page on click
          window.location.href = isAdmin ? '/admin/notifications' : '/notifications';
        }}
      >
        {isConnected ? (
          <BellRing className={cn(
            "h-4 w-4 transition-colors group-hover:text-orange-500",
            unreadCount > 0 ? "text-orange-500" : "text-white"
          )} />
        ) : (
          <Bell className="h-4 w-4 text-gray-400" />
        )}
        {/* Unread count badge */}
        {unreadCount > 0 && (
          <Badge
            variant="default"
            className="absolute -top-1 -right-1 h-5 w-5 rounded-full p-0 text-xs font-bold shadow-large animate-pulse flex items-center justify-center bg-orange-500 text-white border-0"
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

      {/* Notification preview on hover */}
      <div className="absolute top-full right-0 mt-2 w-96 bg-gradient-to-br from-black via-gray-900 to-black border border-gray-700 rounded-xl shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-300 z-50 overflow-hidden">
        {/* Animated Background */}
        <AnimatedBackground className="opacity-20" />

        <div className="relative z-10">

          <div className="pb-3 bg-gradient-to-r from-gray-800/50 to-gray-700/50 backdrop-blur-sm p-4">
            <div className="flex items-center justify-between">
              {getNotificationTitle(isAdmin)}
              <div className="flex items-center gap-2">
                {/* Connection status */}
                <div
                  className={cn(
                    "flex items-center gap-1.5 text-xs px-2 py-1 rounded-full font-medium",
                    isConnected
                      ? "bg-green-900/30 text-green-300 border border-green-700"
                      : "bg-red-900/30 text-red-300 border border-red-700"
                  )}
                >
                  {isConnected ? (
                    <Wifi className="h-3 w-3" />
                  ) : (
                    <WifiOff className="h-3 w-3" />
                  )}
                  {isConnected ? "Live" : "Offline"}
                </div>
              </div>
            </div>
          </div>

          {/* Notifications content */}
          <div className="max-h-96 overflow-y-auto">
            <div className="space-y-1 p-2">
              {notifications.length === 0 ? (
                <div className="p-8 text-center">
                  <div className="w-16 h-16 bg-gray-800 rounded-full flex items-center justify-center mx-auto mb-4">
                    <Bell className="h-8 w-8 text-gray-400" />
                  </div>
                  <h4 className="font-medium text-white mb-2">No notifications yet</h4>
                  <p className="text-sm text-gray-400">Real-time notifications will appear here</p>
                </div>
              ) : (
                notifications.slice(0, 5).map((notification) => (
                  <div
                    key={notification.id}
                    className={cn(
                      "p-3 rounded-xl border-l-4 cursor-pointer transition-all duration-300 group bg-gray-800/30 backdrop-blur-sm hover:bg-gray-700/40",
                      getPriorityColor(notification.priority),
                      !notification.is_read ? "shadow-lg ring-2 ring-orange-500/30" : "opacity-70",
                      "hover:scale-[1.02]"
                    )}
                    onClick={() => markAsRead(notification.id)}
                  >
                    <div className="flex items-start gap-3">
                      <div className="flex-shrink-0 mt-0.5">
                        {getNotificationIcon(notification.category, isAdmin)}
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between mb-1">
                          <h4 className="text-sm font-semibold text-white line-clamp-1 pr-2">
                            {notification.title}
                          </h4>
                          {!notification.is_read && (
                            <div className="h-2 w-2 bg-orange-500 rounded-full flex-shrink-0 mt-1 animate-pulse" />
                          )}
                        </div>
                        <p className="text-sm text-gray-300 line-clamp-2 mb-3 leading-relaxed">
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
                            variant="default"
                            className="text-xs capitalize bg-orange-900/30 text-orange-300 border border-orange-700"
                          >
                            {notification.category}
                          </Badge>
                        </div>
                      </div>
                    </div>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Footer */}
          {notifications.length > 0 && (
            <>
              <div className="border-t border-gray-700"></div>
              <div className="p-3 bg-gray-800/30 backdrop-blur-sm">
                <Button
                  variant="ghost"
                  size="sm"
                  className="w-full text-sm font-medium text-gray-200 hover:text-white hover:bg-gray-700/50 transition-all duration-300 hover:scale-[1.02]"
                  onClick={() => {
                    // Navigate to full notifications page
                    window.location.href = isAdmin ? '/admin/notifications' : '/notifications';
                  }}
                >
                  View All Notifications
                </Button>
                {notifications.length > 5 && (
                  <p className="text-xs text-gray-400 text-center mt-2">
                    +{notifications.length - 5} more notifications
                  </p>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
