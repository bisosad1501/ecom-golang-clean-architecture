export type NotificationType = 'in_app' | 'email' | 'sms' | 'push';

export type NotificationCategory = 
  | 'order' 
  | 'payment' 
  | 'shipping' 
  | 'review' 
  | 'system' 
  | 'promotion' 
  | 'security';

export type NotificationPriority = 'low' | 'normal' | 'high';

export type NotificationStatus = 'pending' | 'processing' | 'sent' | 'delivered' | 'failed' | 'read';

export interface Notification {
  id: string;
  user_id?: string;
  type: NotificationType;
  category: NotificationCategory;
  priority: NotificationPriority;
  status: NotificationStatus;
  title: string;
  message: string;
  data?: string;
  is_read: boolean;
  read_at?: string;
  recipient?: string;
  subject?: string;
  template?: string;
  reference_type?: string;
  reference_id?: string;
  retry_count: number;
  max_retries: number;
  next_retry_at?: string;
  last_error?: string;
  sent_at?: string;
  delivered_at?: string;
  created_at: string;
  updated_at: string;
}

export interface NotificationPreferences {
  id: string;
  user_id: string;
  email_enabled: boolean;
  sms_enabled: boolean;
  push_enabled: boolean;
  in_app_enabled: boolean;
  order_updates: boolean;
  promotional_emails: boolean;
  security_alerts: boolean;
  newsletter_enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface NotificationFilters {
  type?: NotificationType;
  category?: NotificationCategory;
  is_read?: boolean;
  priority?: NotificationPriority;
  date_from?: string;
  date_to?: string;
  limit?: number;
  offset?: number;
}

export interface NotificationResponse {
  notifications: Notification[];
  total: number;
  unread_count: number;
  has_more: boolean;
}

export interface CreateNotificationRequest {
  type: NotificationType;
  category: NotificationCategory;
  priority?: NotificationPriority;
  title: string;
  message: string;
  data?: Record<string, any>;
  recipient?: string;
  subject?: string;
  template?: string;
  reference_type?: string;
  reference_id?: string;
}

export interface UpdateNotificationPreferencesRequest {
  in_app_enabled?: boolean;
  email_enabled?: boolean;
  sms_enabled?: boolean;
  push_enabled?: boolean;
  email_order_updates?: boolean;
  email_payment_updates?: boolean;
  email_shipping_updates?: boolean;
  email_promotions?: boolean;
  email_newsletter?: boolean;
  sms_order_updates?: boolean;
  sms_payment_updates?: boolean;
  sms_shipping_updates?: boolean;
  sms_security_alerts?: boolean;
  push_order_updates?: boolean;
  push_payment_updates?: boolean;
  push_shipping_updates?: boolean;
  push_promotions?: boolean;
  in_app_order_updates?: boolean;
  in_app_payment_updates?: boolean;
  in_app_shipping_updates?: boolean;
  in_app_promotions?: boolean;
  in_app_system_updates?: boolean;
}
