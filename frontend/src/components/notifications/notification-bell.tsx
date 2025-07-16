'use client';

import React, { useState, useEffect } from 'react';
import { Bell, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { useNotifications } from '@/hooks/use-notifications';
import { useAuthStore } from '@/store/auth';
import { NotificationItem } from './notification-item';
import { NotificationDropdown } from './notification-dropdown';

interface NotificationBellProps {
  className?: string;
}

export function NotificationBell({ className }: NotificationBellProps) {
  const [isOpen, setIsOpen] = useState(false);
  const { isAuthenticated, isHydrated } = useAuthStore();
  const {
    notifications,
    unreadCount,
    loading,
    markAsRead,
    markAllAsRead,
    fetchNotifications
  } = useNotifications();

  // Fetch notifications on mount only if authenticated
  useEffect(() => {
    if (isAuthenticated && isHydrated) {
      fetchNotifications();
    }
  }, [fetchNotifications, isAuthenticated, isHydrated]);

  // Auto-refresh notifications every 60 seconds (reduced frequency) with silent mode
  useEffect(() => {
    if (!isAuthenticated || !isHydrated) return;

    const interval = setInterval(() => {
      // Use silent mode to avoid showing errors during auto-refresh
      fetchNotifications(undefined, true);
    }, 60000); // Increased from 30s to 60s

    return () => clearInterval(interval);
  }, [fetchNotifications, isAuthenticated, isHydrated]);

  const handleBellClick = () => {
    setIsOpen(!isOpen);
    if (!isOpen) {
      fetchNotifications();
    }
  };

  const handleMarkAsRead = async (notificationId: string) => {
    await markAsRead(notificationId);
  };

  const handleMarkAllAsRead = async () => {
    await markAllAsRead();
  };

  return (
    <div className={`relative ${className}`}>
      {/* Notification Bell Button */}
      <Button
        variant="ghost"
        size="sm"
        className="relative p-2 hover:bg-gray-100 dark:hover:bg-gray-800"
        onClick={handleBellClick}
      >
        <Bell className="h-5 w-5 text-gray-600 dark:text-gray-300" />
        {unreadCount > 0 && (
          <Badge 
            variant="destructive" 
            className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs"
          >
            {unreadCount > 99 ? '99+' : unreadCount}
          </Badge>
        )}
      </Button>

      {/* Notification Dropdown */}
      {isOpen && (
        <NotificationDropdown
          notifications={notifications}
          unreadCount={unreadCount}
          loading={loading}
          onClose={() => setIsOpen(false)}
          onMarkAsRead={handleMarkAsRead}
          onMarkAllAsRead={handleMarkAllAsRead}
        />
      )}
    </div>
  );
}
