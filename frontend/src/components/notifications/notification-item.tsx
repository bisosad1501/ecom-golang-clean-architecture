'use client';

import React from 'react';
import { formatDistanceToNow } from 'date-fns';
import { vi } from 'date-fns/locale';
import { 
  ShoppingCart, 
  CreditCard, 
  Truck, 
  Star, 
  AlertCircle, 
  Info,
  CheckCircle,
  Package
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Notification, NotificationType, NotificationCategory } from '@/types/notification';

interface NotificationItemProps {
  notification: Notification;
  onMarkAsRead: (notificationId: string) => void;
}

export function NotificationItem({ notification, onMarkAsRead }: NotificationItemProps) {
  const getNotificationIcon = (type: NotificationType, category: NotificationCategory) => {
    const iconClass = "h-5 w-5";
    
    switch (category) {
      case 'order':
        return <ShoppingCart className={`${iconClass} text-blue-500`} />;
      case 'payment':
        return <CreditCard className={`${iconClass} text-green-500`} />;
      case 'shipping':
        return <Truck className={`${iconClass} text-purple-500`} />;
      case 'review':
        return <Star className={`${iconClass} text-yellow-500`} />;
      case 'system':
        return <AlertCircle className={`${iconClass} text-red-500`} />;
      default:
        return <Info className={`${iconClass} text-gray-500`} />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
      case 'normal':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
      case 'low':
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-200';
    }
  };

  const handleClick = () => {
    if (!notification.is_read) {
      onMarkAsRead(notification.id);
    }
    
    // Navigate to related page if reference exists
    if (notification.reference_type && notification.reference_id) {
      switch (notification.reference_type) {
        case 'order':
          window.location.href = `/orders/${notification.reference_id}`;
          break;
        case 'payment':
          window.location.href = `/orders/${notification.reference_id}`;
          break;
        case 'product':
          window.location.href = `/products/${notification.reference_id}`;
          break;
        default:
          break;
      }
    }
  };

  const formatTime = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return formatDistanceToNow(date, { 
        addSuffix: true, 
        locale: vi 
      });
    } catch (error) {
      return 'Vừa xong';
    }
  };

  return (
    <div
      className={`p-4 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer transition-colors ${
        !notification.is_read ? 'bg-blue-50 dark:bg-blue-950' : ''
      }`}
      onClick={handleClick}
    >
      <div className="flex items-start gap-3">
        {/* Icon */}
        <div className="flex-shrink-0 mt-0.5">
          {getNotificationIcon(notification.type, notification.category)}
        </div>

        {/* Content */}
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1">
              <h4 className={`text-sm font-medium ${
                !notification.is_read ? 'text-gray-900 dark:text-gray-100' : 'text-gray-700 dark:text-gray-300'
              }`}>
                {notification.title}
              </h4>
              <p className={`text-sm mt-1 ${
                !notification.is_read ? 'text-gray-700 dark:text-gray-300' : 'text-gray-500 dark:text-gray-400'
              }`}>
                {notification.message}
              </p>
            </div>

            {/* Priority Badge */}
            {notification.priority === 'high' && (
              <Badge 
                variant="secondary" 
                className={`text-xs ${getPriorityColor(notification.priority)}`}
              >
                Quan trọng
              </Badge>
            )}
          </div>

          {/* Footer */}
          <div className="flex items-center justify-between mt-2">
            <span className="text-xs text-gray-500 dark:text-gray-400">
              {formatTime(notification.created_at)}
            </span>

            {!notification.is_read && (
              <div className="flex items-center gap-2">
                <div className="w-2 h-2 bg-blue-500 rounded-full"></div>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 px-2 text-xs"
                  onClick={(e) => {
                    e.stopPropagation();
                    onMarkAsRead(notification.id);
                  }}
                >
                  <CheckCircle className="h-3 w-3 mr-1" />
                  Đánh dấu đã đọc
                </Button>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
