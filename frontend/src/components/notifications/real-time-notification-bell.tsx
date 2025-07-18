'use client';

import { useState } from 'react';
import { Bell, Wifi, WifiOff, Circle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuHeader,
  DropdownMenuTrigger,
  DropdownMenuSeparator,
} from '@/components/ui/dropdown-menu';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useWebSocketNotifications } from '@/hooks/use-websocket-notifications';
import { formatDistanceToNow } from 'date-fns';
import { vi } from 'date-fns/locale';

export function RealTimeNotificationBell() {
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

  // Get priority color
  const getPriorityColor = (priority: string): string => {
    switch (priority) {
      case 'high': return 'text-red-500';
      case 'normal': return 'text-blue-500';
      case 'low': return 'text-gray-500';
      default: return 'text-gray-500';
    }
  };

  // Handle notification click
  const handleNotificationClick = (notification: any) => {
    if (!notification.is_read) {
      markAsRead(notification.id);
    }

    // Navigate based on notification type
    if (notification.category === 'order' && notification.reference_id) {
      window.location.href = `/orders/${notification.reference_id}`;
    } else if (notification.category === 'payment' && notification.reference_id) {
      window.location.href = `/orders/${notification.reference_id}`;
    }
  };

  // Send test notification
  const handleTestNotification = () => {
    const testMessages = [
      {
        title: '🛒 Order Confirmed!',
        message: 'Your order #ORD-12345 has been placed successfully. Total: $1,199.99'
      },
      {
        title: '💳 Payment Received',
        message: 'Payment for order #ORD-12345 has been confirmed.'
      },
      {
        title: '🚚 Order Shipped!',
        message: 'Your order is on its way! Track: TRK123456789'
      },
      {
        title: '📦 Order Delivered!',
        message: 'Your order has been delivered! Please leave a review.'
      }
    ];

    const randomMessage = testMessages[Math.floor(Math.random() * testMessages.length)];
    sendTestNotification(randomMessage.title, randomMessage.message);
  };

  return (
    <DropdownMenu open={isOpen} onOpenChange={setIsOpen}>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="relative">
          <Bell className="h-5 w-5" />
          
          {/* Connection status indicator */}
          <div className="absolute -top-1 -right-1">
            {isConnected ? (
              <Wifi className="h-3 w-3 text-green-500" />
            ) : (
              <WifiOff className="h-3 w-3 text-red-500" />
            )}
          </div>
          
          {/* Unread count badge */}
          {unreadCount > 0 && (
            <Badge 
              variant="destructive" 
              className="absolute -top-2 -right-2 h-5 w-5 rounded-full p-0 text-xs flex items-center justify-center"
            >
              {unreadCount > 99 ? '99+' : unreadCount}
            </Badge>
          )}
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-80">
        <DropdownMenuHeader className="pb-2">
          <div className="flex items-center justify-between">
            <h3 className="font-semibold">🔔 Real-time Notifications</h3>
            <div className="flex items-center gap-2">
              {/* Connection status */}
              <div className="flex items-center gap-1 text-xs">
                <Circle 
                  className={`h-2 w-2 fill-current ${
                    isConnected ? 'text-green-500' : 'text-red-500'
                  }`} 
                />
                <span className={isConnected ? 'text-green-600' : 'text-red-600'}>
                  {isConnected ? 'Connected' : 'Disconnected'}
                </span>
              </div>
            </div>
          </div>
          
          {/* Connection info */}
          <div className="text-xs text-muted-foreground mt-1">
            {isConnected && connectionTime && (
              <div>Connected: {formatDistanceToNow(connectionTime, { addSuffix: true, locale: vi })}</div>
            )}
            {lastMessage && (
              <div>Last message: {formatDistanceToNow(lastMessage, { addSuffix: true, locale: vi })}</div>
            )}
            {error && (
              <div className="text-red-500">Error: {error}</div>
            )}
          </div>
        </DropdownMenuHeader>

        <DropdownMenuSeparator />

        {/* Test notification button */}
        <div className="p-2">
          <Button 
            variant="outline" 
            size="sm" 
            onClick={handleTestNotification}
            className="w-full text-xs"
          >
            🧪 Send Test Notification
          </Button>
        </div>

        <DropdownMenuSeparator />

        {/* Notifications list */}
        <ScrollArea className="h-96">
          {notifications.length === 0 ? (
            <div className="p-4 text-center text-muted-foreground">
              <Bell className="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">No notifications yet</p>
              <p className="text-xs">Real-time notifications will appear here</p>
            </div>
          ) : (
            <div className="space-y-1">
              {notifications.map((notification) => (
                <div
                  key={notification.id}
                  onClick={() => handleNotificationClick(notification)}
                  className={`p-3 hover:bg-muted cursor-pointer border-l-2 ${
                    notification.is_read 
                      ? 'border-l-gray-200 opacity-70' 
                      : 'border-l-blue-500'
                  }`}
                >
                  <div className="flex items-start gap-2">
                    <span className="text-lg">
                      {getNotificationIcon(notification.category)}
                    </span>
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h4 className="font-medium text-sm truncate">
                          {notification.title}
                        </h4>
                        <Circle 
                          className={`h-2 w-2 fill-current ${getPriorityColor(notification.priority)}`} 
                        />
                      </div>
                      <p className="text-xs text-muted-foreground mt-1 line-clamp-2">
                        {notification.message}
                      </p>
                      <div className="flex items-center justify-between mt-2">
                        <span className="text-xs text-muted-foreground">
                          {formatDistanceToNow(new Date(notification.created_at), { 
                            addSuffix: true, 
                            locale: vi 
                          })}
                        </span>
                        <Badge variant="outline" className="text-xs">
                          {notification.category}
                        </Badge>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </ScrollArea>

        {/* Footer actions */}
        {notifications.length > 0 && (
          <>
            <DropdownMenuSeparator />
            <div className="p-2 flex gap-2">
              <Button 
                variant="outline" 
                size="sm" 
                onClick={clearNotifications}
                className="flex-1 text-xs"
              >
                Clear All
              </Button>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={() => window.location.href = '/notifications'}
                className="flex-1 text-xs"
              >
                View All
              </Button>
            </div>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
