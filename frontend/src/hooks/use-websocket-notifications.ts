'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { toast } from 'sonner';
import { Notification as NotificationData } from '@/types/notification';
import { useAuthStore } from '@/store/auth';

interface WebSocketMessage {
  type: 'notification' | 'system' | 'error';
  event: string;
  notification?: NotificationData;
  data?: Record<string, unknown>;
  timestamp: string;
}

interface WebSocketStats {
  connected_users: number;
  total_clients: number;
  users_with_clients: Record<string, unknown>;
}

export function useWebSocketNotifications() {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionTime, setConnectionTime] = useState<Date | null>(null);
  const [lastMessage, setLastMessage] = useState<Date | null>(null);
  const [notifications, setNotifications] = useState<NotificationData[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const reconnectAttempts = useRef(0);
  const maxReconnectAttempts = 5;
  const { token, isAuthenticated } = useAuthStore();

  // WebSocket URL
  const getWebSocketUrl = useCallback(() => {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const host = process.env.NODE_ENV === 'production' 
      ? window.location.host 
      : 'localhost:8080';
    return `${protocol}//${host}/api/v1/ws/notifications?token=${token}`;
  }, [token]);

  // Handle incoming WebSocket messages
  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const message: WebSocketMessage = JSON.parse(event.data);
      setLastMessage(new Date());
      
      console.log('📨 WebSocket message received:', message);

      switch (message.type) {
        case 'notification':
          if (message.notification) {
            handleNewNotification(message.notification);
          }
          break;
          
        case 'system':
          console.log('🔧 System message:', message.data?.message);
          // if (message.event === 'connected') {
          //   toast.success('🔔 Connected to real-time notifications');
          // }
          break;
          
        case 'error':
          console.error('❌ WebSocket error:', message.data);
          const errorMsg = message.data && typeof message.data === 'object' && 'message' in message.data
            ? String(message.data.message)
            : 'Unknown error';
          toast.error(`WebSocket error: ${errorMsg}`);
          break;
          
        default:
          console.log('📨 Unknown message type:', message);
      }
    } catch (err) {
      console.error('❌ Failed to parse WebSocket message:', err);
    }
  }, []);

  // Handle new notification
  const handleNewNotification = useCallback((notification: NotificationData) => {
    console.log('🔔 New real-time notification:', notification);
    
    // Add to notifications list
    setNotifications(prev => [notification, ...prev]);
    
    // Update unread count if notification is unread
    if (!notification.is_read) {
      setUnreadCount(prev => prev + 1);
    }

    // Show toast notification
    const icon = getNotificationIcon(notification.category);
    const title = `${icon} ${notification.title}`;
    
    toast(title, {
      description: notification.message,
      duration: 5000,
      action: notification.category === 'order' ? {
        label: 'View Order',
        onClick: () => {
          // Navigate to order details
          if (notification.reference_id) {
            window.location.href = `/orders/${notification.reference_id}`;
          }
        }
      } : undefined,
    });

    // Show browser notification if permission granted
    if (typeof window !== 'undefined' && 'Notification' in window && window.Notification.permission === 'granted') {
      new window.Notification(notification.title, {
        body: notification.message,
        icon: '/logo-Bihub.png',
        tag: notification.id,
        requireInteraction: notification.priority === 'high',
      });
    }
  }, []);

  // Get notification icon based on category
  const getNotificationIcon = (category: string): string => {
    const icons: Record<string, string> = {
      order: '🛒',
      payment: '💳',
      shipping: '🚚',
      promotion: '🎉',
      account: '🔐',
      system: '⚙️',
      review: '⭐',
      inventory: '📦',
      security: '🔒',
      cart: '🛍️',
    };
    return icons[category] || '🔔';
  };

  // Connect to WebSocket
  const connect = useCallback(() => {
    if (!isAuthenticated || !token) {
      console.log('❌ Cannot connect: not authenticated');
      return;
    }

    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.log('✅ WebSocket already connected');
      return;
    }

    try {
      const wsUrl = getWebSocketUrl();
      console.log('🔌 Connecting to WebSocket:', wsUrl);
      
      wsRef.current = new WebSocket(wsUrl);

      wsRef.current.onopen = () => {
        console.log('✅ WebSocket connected');
        setIsConnected(true);
        setConnectionTime(new Date());
        setError(null);
        reconnectAttempts.current = 0;
      };

      wsRef.current.onmessage = handleMessage;

      wsRef.current.onclose = (event) => {
        console.log('❌ WebSocket disconnected:', event.code, event.reason);
        setIsConnected(false);
        setConnectionTime(null);
        
        // Attempt to reconnect if not a normal closure
        if (event.code !== 1000 && reconnectAttempts.current < maxReconnectAttempts) {
          const delay = Math.pow(2, reconnectAttempts.current) * 1000; // Exponential backoff
          console.log(`🔄 Reconnecting in ${delay}ms (attempt ${reconnectAttempts.current + 1})`);
          
          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttempts.current++;
            connect();
          }, delay);
        }
      };

      wsRef.current.onerror = (error) => {
        console.error('❌ WebSocket error:', error);
        setError('WebSocket connection error');
        toast.error('🔌 Connection error - notifications may be delayed');
      };

    } catch (err) {
      console.error('❌ Failed to create WebSocket connection:', err);
      setError('Failed to connect to notification service');
    }
  }, [isAuthenticated, token, getWebSocketUrl, handleMessage]);

  // Disconnect from WebSocket
  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }

    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnected');
      wsRef.current = null;
    }

    setIsConnected(false);
    setConnectionTime(null);
    reconnectAttempts.current = 0;
  }, []);

  // Send test notification
  const sendTestNotification = useCallback(async (title: string, message: string) => {
    if (!token) return;

    try {
      // Get user info from token to get user ID
      const userResponse = await fetch('http://localhost:8080/api/v1/auth/me', {
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (!userResponse.ok) {
        toast.error('❌ Failed to get user info');
        return;
      }

      const userData = await userResponse.json();
      const userId = userData.data?.id;

      if (!userId) {
        toast.error('❌ User ID not found');
        return;
      }

      const response = await fetch(`http://localhost:8080/api/v1/ws/test/${userId}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          title,
          message,
          category: 'test',
          priority: 'normal',
        }),
      });

      if (response.ok) {
        toast.success('✅ Test notification sent');
      } else {
        toast.error('❌ Failed to send test notification');
      }
    } catch (error) {
      console.error('❌ Error sending test notification:', error);
      toast.error('❌ Error sending test notification');
    }
  }, [token]);

  // Get WebSocket stats
  const getStats = useCallback(async (): Promise<WebSocketStats | null> => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/ws/stats', {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        return data.data;
      }
    } catch (error) {
      console.error('❌ Error fetching WebSocket stats:', error);
    }
    return null;
  }, [token]);

  // Request notification permission
  const requestNotificationPermission = useCallback(async () => {
    if (typeof window !== 'undefined' && 'Notification' in window && window.Notification.permission === 'default') {
      const permission = await window.Notification.requestPermission();
      if (permission === 'granted') {
        toast.success('🔔 Browser notifications enabled');
      } else {
        toast.info('🔕 Browser notifications disabled');
      }
      return permission;
    }
    return typeof window !== 'undefined' && 'Notification' in window ? window.Notification.permission : 'default';
  }, []);

  // Auto-connect when authenticated
  useEffect(() => {
    if (isAuthenticated && token) {
      // Small delay to ensure token is ready
      const timer = setTimeout(connect, 1000);
      return () => clearTimeout(timer);
    } else {
      disconnect();
    }
  }, [isAuthenticated, token, connect, disconnect]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      disconnect();
    };
  }, [disconnect]);

  // Request notification permission on mount
  useEffect(() => {
    requestNotificationPermission();
  }, [requestNotificationPermission]);

  return {
    // Connection state
    isConnected,
    connectionTime,
    lastMessage,
    error,
    
    // Notifications
    notifications,
    unreadCount,
    
    // Actions
    connect,
    disconnect,
    sendTestNotification,
    getStats,
    requestNotificationPermission,
    
    // Utils
    clearNotifications: () => setNotifications([]),
    markAsRead: (notificationId: string) => {
      setNotifications(prev => 
        prev.map(n => 
          n.id === notificationId 
            ? { ...n, is_read: true, read_at: new Date().toISOString() }
            : n
        )
      );
      setUnreadCount(prev => Math.max(0, prev - 1));
    },
  };
}
