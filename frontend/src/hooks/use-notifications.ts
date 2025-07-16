'use client';

import { useState, useCallback, useEffect } from 'react';
import { toast } from 'sonner';
import { 
  Notification, 
  NotificationFilters, 
  NotificationResponse,
  NotificationPreferences,
  UpdateNotificationPreferencesRequest 
} from '@/types/notification';
import { notificationApi } from '@/lib/api/notification';

export function useNotifications() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch notifications
  const fetchNotifications = useCallback(async (filters?: NotificationFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await notificationApi.getUserNotifications(filters);
      setNotifications(response.notifications);
      setUnreadCount(response.unread_count);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Không thể tải thông báo';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  // Fetch unread count only
  const fetchUnreadCount = useCallback(async () => {
    try {
      const count = await notificationApi.getUnreadCount();
      setUnreadCount(count);
    } catch (err) {
      console.error('Failed to fetch unread count:', err);
    }
  }, []);

  // Mark notification as read
  const markAsRead = useCallback(async (notificationId: string) => {
    try {
      await notificationApi.markAsRead(notificationId);
      
      // Update local state
      setNotifications(prev => 
        prev.map(notification => 
          notification.id === notificationId 
            ? { ...notification, is_read: true, read_at: new Date().toISOString() }
            : notification
        )
      );
      
      // Update unread count
      setUnreadCount(prev => Math.max(0, prev - 1));
      
      toast.success('Đã đánh dấu thông báo là đã đọc');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Không thể đánh dấu thông báo';
      toast.error(errorMessage);
    }
  }, []);

  // Mark all notifications as read
  const markAllAsRead = useCallback(async () => {
    try {
      await notificationApi.markAllAsRead();
      
      // Update local state
      setNotifications(prev => 
        prev.map(notification => ({ 
          ...notification, 
          is_read: true, 
          read_at: new Date().toISOString() 
        }))
      );
      
      setUnreadCount(0);
      toast.success('Đã đánh dấu tất cả thông báo là đã đọc');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Không thể đánh dấu tất cả thông báo';
      toast.error(errorMessage);
    }
  }, []);

  return {
    notifications,
    unreadCount,
    loading,
    error,
    fetchNotifications,
    fetchUnreadCount,
    markAsRead,
    markAllAsRead,
  };
}

export function useNotificationPreferences() {
  const [preferences, setPreferences] = useState<NotificationPreferences | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Fetch user preferences
  const fetchPreferences = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const prefs = await notificationApi.getUserPreferences();
      setPreferences(prefs);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Không thể tải cài đặt thông báo';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  // Update user preferences
  const updatePreferences = useCallback(async (updates: UpdateNotificationPreferencesRequest) => {
    try {
      setLoading(true);
      setError(null);
      
      const updatedPrefs = await notificationApi.updateUserPreferences(updates);
      setPreferences(updatedPrefs);
      
      toast.success('Đã cập nhật cài đặt thông báo');
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Không thể cập nhật cài đặt thông báo';
      setError(errorMessage);
      toast.error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  // Load preferences on mount
  useEffect(() => {
    fetchPreferences();
  }, [fetchPreferences]);

  return {
    preferences,
    loading,
    error,
    fetchPreferences,
    updatePreferences,
  };
}
