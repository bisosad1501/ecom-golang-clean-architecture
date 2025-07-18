'use client';

import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { ScrollArea } from '@/components/ui/scroll-area';
import { 
  Bell, 
  Wifi, 
  WifiOff, 
  Users, 
  Activity, 
  Clock,
  Send,
  Trash2,
  RefreshCw
} from 'lucide-react';
import { useWebSocketNotifications } from '@/hooks/use-websocket-notifications';
import { formatDistanceToNow } from 'date-fns';
import { vi } from 'date-fns/locale';
import { toast } from 'sonner';

export default function RealTimeNotificationsPage() {
  const {
    isConnected,
    connectionTime,
    lastMessage,
    notifications,
    unreadCount,
    error,
    connect,
    disconnect,
    sendTestNotification,
    getStats,
    clearNotifications,
    markAsRead,
  } = useWebSocketNotifications();

  const [stats, setStats] = useState<any>(null);
  const [testTitle, setTestTitle] = useState('🛒 Test Order Notification');
  const [testMessage, setTestMessage] = useState('Your test order has been placed successfully!');

  // Fetch WebSocket stats
  const fetchStats = async () => {
    const statsData = await getStats();
    setStats(statsData);
  };

  // Auto-fetch stats every 10 seconds
  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 10000);
    return () => clearInterval(interval);
  }, []);

  // Get notification icon
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

  // Send test notification
  const handleSendTest = () => {
    if (!testTitle.trim() || !testMessage.trim()) {
      toast.error('Please enter both title and message');
      return;
    }
    sendTestNotification(testTitle, testMessage);
  };

  // Quick test notifications
  const quickTests = [
    {
      title: '🛒 Order Confirmed!',
      message: 'Your order #ORD-12345 has been placed successfully. Total: $1,199.99',
    },
    {
      title: '💳 Payment Received',
      message: 'Payment for order #ORD-12345 has been confirmed.',
    },
    {
      title: '🚚 Order Shipped!',
      message: 'Your order is on its way! Track: TRK123456789',
    },
    {
      title: '📦 Order Delivered!',
      message: 'Your order has been delivered! Please leave a review.',
    },
    {
      title: '🎉 Special Offer!',
      message: 'Flash Sale: 20% off all smartphones! Use code FLASH20.',
    },
    {
      title: '🔐 Security Alert',
      message: 'New device login detected from MacBook Pro in San Francisco.',
    },
  ];

  return (
    <div className="container mx-auto py-8 space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">🔔 Real-time Notifications</h1>
          <p className="text-muted-foreground">
            Test and monitor WebSocket-based real-time notifications
          </p>
        </div>
        <Button onClick={fetchStats} variant="outline" size="sm">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh Stats
        </Button>
      </div>

      {/* Connection Status */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            {isConnected ? (
              <Wifi className="h-5 w-5 text-green-500" />
            ) : (
              <WifiOff className="h-5 w-5 text-red-500" />
            )}
            Connection Status
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="text-center">
              <div className={`text-2xl font-bold ${isConnected ? 'text-green-600' : 'text-red-600'}`}>
                {isConnected ? 'Connected' : 'Disconnected'}
              </div>
              <div className="text-sm text-muted-foreground">Status</div>
            </div>
            
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600">
                {connectionTime ? formatDistanceToNow(connectionTime, { addSuffix: true, locale: vi }) : '--'}
              </div>
              <div className="text-sm text-muted-foreground">Connected Since</div>
            </div>
            
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-600">
                {lastMessage ? formatDistanceToNow(lastMessage, { addSuffix: true, locale: vi }) : '--'}
              </div>
              <div className="text-sm text-muted-foreground">Last Message</div>
            </div>
            
            <div className="text-center">
              <div className="text-2xl font-bold text-orange-600">
                {unreadCount}
              </div>
              <div className="text-sm text-muted-foreground">Unread Count</div>
            </div>
          </div>

          {error && (
            <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
              <p className="text-red-600 text-sm">❌ {error}</p>
            </div>
          )}

          <div className="mt-4 flex gap-2">
            <Button 
              onClick={connect} 
              disabled={isConnected}
              variant="outline"
              size="sm"
            >
              <Wifi className="h-4 w-4 mr-2" />
              Connect
            </Button>
            <Button 
              onClick={disconnect} 
              disabled={!isConnected}
              variant="outline"
              size="sm"
            >
              <WifiOff className="h-4 w-4 mr-2" />
              Disconnect
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* WebSocket Stats */}
      {stats && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              WebSocket Statistics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {stats.total_clients || 0}
                </div>
                <div className="text-sm text-muted-foreground">Total Clients</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-green-600">
                  {stats.connected_users || 0}
                </div>
                <div className="text-sm text-muted-foreground">Connected Users</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-purple-600">
                  {Object.keys(stats.users_with_clients || {}).length}
                </div>
                <div className="text-sm text-muted-foreground">Active Sessions</div>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Test Notifications */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Send className="h-5 w-5" />
              Test Notifications
            </CardTitle>
            <CardDescription>
              Send test notifications to verify real-time delivery
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            {/* Custom test */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Custom Test</label>
              <input
                type="text"
                placeholder="Notification title..."
                value={testTitle}
                onChange={(e) => setTestTitle(e.target.value)}
                className="w-full p-2 border rounded-md text-sm"
              />
              <textarea
                placeholder="Notification message..."
                value={testMessage}
                onChange={(e) => setTestMessage(e.target.value)}
                className="w-full p-2 border rounded-md text-sm h-20 resize-none"
              />
              <Button onClick={handleSendTest} className="w-full">
                <Send className="h-4 w-4 mr-2" />
                Send Test Notification
              </Button>
            </div>

            <Separator />

            {/* Quick tests */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Quick Tests</label>
              <div className="grid grid-cols-1 gap-2">
                {quickTests.map((test, index) => (
                  <Button
                    key={index}
                    variant="outline"
                    size="sm"
                    onClick={() => sendTestNotification(test.title, test.message)}
                    className="justify-start text-left h-auto p-3"
                  >
                    <div>
                      <div className="font-medium text-sm">{test.title}</div>
                      <div className="text-xs text-muted-foreground truncate">
                        {test.message}
                      </div>
                    </div>
                  </Button>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Real-time Notifications */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <Bell className="h-5 w-5" />
                Real-time Notifications
                {unreadCount > 0 && (
                  <Badge variant="destructive">{unreadCount}</Badge>
                )}
              </div>
              {notifications.length > 0 && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={clearNotifications}
                >
                  <Trash2 className="h-4 w-4 mr-2" />
                  Clear All
                </Button>
              )}
            </CardTitle>
            <CardDescription>
              Notifications received in real-time via WebSocket
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ScrollArea className="h-96">
              {notifications.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  <Bell className="h-12 w-12 mx-auto mb-4 opacity-50" />
                  <p>No notifications yet</p>
                  <p className="text-sm">Send a test notification to see it appear here instantly!</p>
                </div>
              ) : (
                <div className="space-y-3">
                  {notifications.map((notification) => (
                    <div
                      key={notification.id}
                      onClick={() => markAsRead(notification.id)}
                      className={`p-3 border rounded-lg cursor-pointer transition-colors hover:bg-muted ${
                        notification.is_read ? 'opacity-70' : 'border-blue-200 bg-blue-50'
                      }`}
                    >
                      <div className="flex items-start gap-3">
                        <span className="text-xl">
                          {getNotificationIcon(notification.category)}
                        </span>
                        <div className="flex-1 min-w-0">
                          <div className="flex items-center gap-2 mb-1">
                            <h4 className="font-medium text-sm">{notification.title}</h4>
                            <Badge variant="outline" className="text-xs">
                              {notification.priority}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground mb-2">
                            {notification.message}
                          </p>
                          <div className="flex items-center justify-between text-xs text-muted-foreground">
                            <span>
                              {formatDistanceToNow(new Date(notification.created_at), { 
                                addSuffix: true, 
                                locale: vi 
                              })}
                            </span>
                            <Badge variant="secondary" className="text-xs">
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
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
