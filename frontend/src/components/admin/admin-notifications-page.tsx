'use client'

import React, { useState } from 'react'
import {
  Bell,
  Mail,
  MessageSquare,
  AlertTriangle,
  CheckCircle,
  Info,
  X,
  Eye,
  Trash2,
  Filter,
  Search,
  Settings,
  Send,
  Users,
  Calendar,
  Clock
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
  BiHubActionButton,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME, getBadgeVariant } from '@/constants/admin-theme'

// Mock data for notifications
const mockNotifications = [
  {
    id: '1',
    type: 'order',
    title: 'New Order Received',
    message: 'Order #12345 has been placed by John Doe',
    timestamp: '2 minutes ago',
    read: false,
    priority: 'high',
    icon: 'ShoppingCart'
  },
  {
    id: '2',
    type: 'user',
    title: 'New User Registration',
    message: 'Jane Smith has created a new account',
    timestamp: '15 minutes ago',
    read: false,
    priority: 'medium',
    icon: 'Users'
  },
  {
    id: '3',
    type: 'system',
    title: 'System Backup Completed',
    message: 'Daily backup has been completed successfully',
    timestamp: '1 hour ago',
    read: true,
    priority: 'low',
    icon: 'CheckCircle'
  },
  {
    id: '4',
    type: 'security',
    title: 'Security Alert',
    message: 'Multiple failed login attempts detected',
    timestamp: '2 hours ago',
    read: false,
    priority: 'high',
    icon: 'AlertTriangle'
  },
  {
    id: '5',
    type: 'review',
    title: 'New Product Review',
    message: 'A new 5-star review has been submitted',
    timestamp: '3 hours ago',
    read: true,
    priority: 'medium',
    icon: 'MessageSquare'
  }
]

const notificationStats = [
  { label: 'Total Notifications', value: '156', change: '+12%', trend: 'up' },
  { label: 'Unread', value: '23', change: '+5%', trend: 'up' },
  { label: 'High Priority', value: '8', change: '-2%', trend: 'down' },
  { label: 'This Week', value: '45', change: '+18%', trend: 'up' }
]

export default function AdminNotificationsPage() {
  const [searchTerm, setSearchTerm] = useState('')
  const [filterType, setFilterType] = useState('all')
  const [filterPriority, setFilterPriority] = useState('all')
  const [notifications, setNotifications] = useState(mockNotifications)

  const handleMarkAsRead = (id: string) => {
    setNotifications(prev =>
      prev.map(notif =>
        notif.id === id ? { ...notif, read: true } : notif
      )
    )
  }

  const handleMarkAllAsRead = () => {
    setNotifications(prev =>
      prev.map(notif => ({ ...notif, read: true }))
    )
  }

  const handleDelete = (id: string) => {
    setNotifications(prev => prev.filter(notif => notif.id !== id))
  }

  const filteredNotifications = notifications.filter(notif => {
    const matchesSearch = notif.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         notif.message.toLowerCase().includes(searchTerm.toLowerCase())
    const matchesType = filterType === 'all' || notif.type === filterType
    const matchesPriority = filterPriority === 'all' || notif.priority === filterPriority
    
    return matchesSearch && matchesType && matchesPriority
  })

  const getNotificationIcon = (type: string) => {
    switch (type) {
      case 'order': return <Bell className="h-5 w-5" />
      case 'user': return <Users className="h-5 w-5" />
      case 'system': return <Settings className="h-5 w-5" />
      case 'security': return <AlertTriangle className="h-5 w-5" />
      case 'review': return <MessageSquare className="h-5 w-5" />
      default: return <Info className="h-5 w-5" />
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'text-red-400 bg-red-500/20'
      case 'medium': return 'text-yellow-400 bg-yellow-500/20'
      case 'low': return 'text-green-400 bg-green-500/20'
      default: return 'text-gray-400 bg-gray-500/20'
    }
  }

  return (
    <div className="space-y-8">
      <BiHubPageHeader
        title="Notifications"
        subtitle="Manage system notifications and alerts"
        icon={<Bell className="h-6 w-6 text-white" />}
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {notificationStats.map((stat, index) => (
          <BiHubAdminCard
            key={index}
            title={stat.label}
            subtitle={stat.value}
            icon={<Bell className="h-5 w-5 text-white" />}
          >
            <div className="flex items-center justify-between">
              <span className="text-2xl font-bold text-white">{stat.value}</span>
              <BiHubStatusBadge
                status={stat.trend === 'up' ? 'success' : 'warning'}
                text={stat.change}
              />
            </div>
          </BiHubAdminCard>
        ))}
      </div>

      {/* Filters and Actions */}
      <BiHubAdminCard
        title="Notification Management"
        subtitle="Filter and manage your notifications"
        icon={<Filter className="h-5 w-5 text-white" />}
      >
        <div className="space-y-6">
          {/* Search and Filters */}
          <div className="flex flex-col lg:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  placeholder="Search notifications..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 bg-gray-800 border-gray-600 text-white"
                />
              </div>
            </div>
            
            <Select value={filterType} onValueChange={setFilterType}>
              <SelectTrigger className="w-full lg:w-48 bg-gray-800 border-gray-600 text-white">
                <SelectValue placeholder="Filter by type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                <SelectItem value="order">Orders</SelectItem>
                <SelectItem value="user">Users</SelectItem>
                <SelectItem value="system">System</SelectItem>
                <SelectItem value="security">Security</SelectItem>
                <SelectItem value="review">Reviews</SelectItem>
              </SelectContent>
            </Select>

            <Select value={filterPriority} onValueChange={setFilterPriority}>
              <SelectTrigger className="w-full lg:w-48 bg-gray-800 border-gray-600 text-white">
                <SelectValue placeholder="Filter by priority" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Priorities</SelectItem>
                <SelectItem value="high">High</SelectItem>
                <SelectItem value="medium">Medium</SelectItem>
                <SelectItem value="low">Low</SelectItem>
              </SelectContent>
            </Select>

            <Button
              onClick={handleMarkAllAsRead}
              className="bg-[#FF9000] hover:bg-[#e67e00] text-white"
            >
              <CheckCircle className="h-4 w-4 mr-2" />
              Mark All Read
            </Button>
          </div>

          {/* Notifications List */}
          <div className="space-y-4">
            {filteredNotifications.map((notification) => (
              <div
                key={notification.id}
                className={`p-4 rounded-xl border transition-all duration-200 ${
                  notification.read
                    ? 'bg-gray-800/50 border-gray-600/30'
                    : 'bg-gray-800 border-[#FF9000]/30 shadow-lg'
                }`}
              >
                <div className="flex items-start gap-4">
                  <div className={`p-2 rounded-lg ${getPriorityColor(notification.priority)}`}>
                    {getNotificationIcon(notification.type)}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1">
                        <h3 className={`font-semibold ${notification.read ? 'text-gray-300' : 'text-white'}`}>
                          {notification.title}
                          {!notification.read && (
                            <span className="ml-2 w-2 h-2 bg-[#FF9000] rounded-full inline-block"></span>
                          )}
                        </h3>
                        <p className="text-gray-400 text-sm mt-1">{notification.message}</p>
                        <div className="flex items-center gap-4 mt-2">
                          <span className="text-xs text-gray-500 flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {notification.timestamp}
                          </span>
                          <Badge variant="outline" className={`text-xs ${getPriorityColor(notification.priority)}`}>
                            {notification.priority}
                          </Badge>
                        </div>
                      </div>
                      
                      <div className="flex items-center gap-2">
                        {!notification.read && (
                          <Button
                            size="sm"
                            variant="ghost"
                            onClick={() => handleMarkAsRead(notification.id)}
                            className="text-gray-400 hover:text-white"
                          >
                            <Eye className="h-4 w-4" />
                          </Button>
                        )}
                        <Button
                          size="sm"
                          variant="ghost"
                          onClick={() => handleDelete(notification.id)}
                          className="text-gray-400 hover:text-red-400"
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {filteredNotifications.length === 0 && (
            <div className="text-center py-12">
              <Bell className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-400 mb-2">No notifications found</h3>
              <p className="text-gray-500">Try adjusting your search or filter criteria.</p>
            </div>
          )}
        </div>
      </BiHubAdminCard>
    </div>
  )
}
