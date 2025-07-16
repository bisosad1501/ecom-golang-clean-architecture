import { apiClient } from './client';
import { 
  Notification, 
  NotificationFilters, 
  NotificationResponse,
  NotificationPreferences,
  UpdateNotificationPreferencesRequest,
  CreateNotificationRequest 
} from '@/types/notification';

export const notificationApi = {
  // Get user notifications
  async getUserNotifications(filters?: NotificationFilters): Promise<NotificationResponse> {
    const params = new URLSearchParams();
    
    if (filters?.type) params.append('type', filters.type);
    if (filters?.category) params.append('category', filters.category);
    if (filters?.is_read !== undefined) params.append('is_read', filters.is_read.toString());
    if (filters?.priority) params.append('priority', filters.priority);
    if (filters?.date_from) params.append('date_from', filters.date_from);
    if (filters?.date_to) params.append('date_to', filters.date_to);
    if (filters?.limit) params.append('limit', filters.limit.toString());
    if (filters?.offset) params.append('offset', filters.offset.toString());

    const queryString = params.toString();
    const url = `/notifications${queryString ? `?${queryString}` : ''}`;
    
    const response = await apiClient.get(url);
    return response.data;
  },

  // Get unread notification count
  async getUnreadCount(): Promise<number> {
    const response = await apiClient.get('/notifications/count');
    return response.data.count;
  },

  // Mark notification as read
  async markAsRead(notificationId: string): Promise<void> {
    await apiClient.put(`/notifications/${notificationId}/read`);
  },

  // Mark all notifications as read
  async markAllAsRead(): Promise<void> {
    await apiClient.put('/notifications/read-all');
  },

  // Get user notification preferences
  async getUserPreferences(): Promise<NotificationPreferences> {
    const response = await apiClient.get('/notifications/preferences');
    return response.data;
  },

  // Update user notification preferences
  async updateUserPreferences(updates: UpdateNotificationPreferencesRequest): Promise<NotificationPreferences> {
    const response = await apiClient.put('/notifications/preferences', updates);
    return response.data;
  },

  // Create notification (admin only)
  async createNotification(notification: CreateNotificationRequest): Promise<Notification> {
    const response = await apiClient.post('/admin/notifications', notification);
    return response.data;
  },

  // Get all notifications (admin only)
  async getAllNotifications(filters?: NotificationFilters): Promise<NotificationResponse> {
    const params = new URLSearchParams();
    
    if (filters?.type) params.append('type', filters.type);
    if (filters?.category) params.append('category', filters.category);
    if (filters?.is_read !== undefined) params.append('is_read', filters.is_read.toString());
    if (filters?.priority) params.append('priority', filters.priority);
    if (filters?.date_from) params.append('date_from', filters.date_from);
    if (filters?.date_to) params.append('date_to', filters.date_to);
    if (filters?.limit) params.append('limit', filters.limit.toString());
    if (filters?.offset) params.append('offset', filters.offset.toString());

    const queryString = params.toString();
    const url = `/admin/notifications${queryString ? `?${queryString}` : ''}`;
    
    const response = await apiClient.get(url);
    return response.data;
  },

  // Delete notification (admin only)
  async deleteNotification(notificationId: string): Promise<void> {
    await apiClient.delete(`/admin/notifications/${notificationId}`);
  },

  // Get notification statistics (admin only)
  async getNotificationStats(dateFrom?: string, dateTo?: string): Promise<any> {
    const params = new URLSearchParams();
    if (dateFrom) params.append('date_from', dateFrom);
    if (dateTo) params.append('date_to', dateTo);

    const queryString = params.toString();
    const url = `/admin/notifications/stats${queryString ? `?${queryString}` : ''}`;
    
    const response = await apiClient.get(url);
    return response.data;
  },

  // Retry failed notifications (admin only)
  async retryFailedNotifications(): Promise<void> {
    await apiClient.post('/admin/notifications/retry-failed');
  },
};
