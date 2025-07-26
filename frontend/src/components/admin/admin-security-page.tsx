'use client'

import { useState } from 'react'
import {
  Shield,
  AlertTriangle,
  Eye,
  Lock,
  Unlock,
  UserX,
  Activity,
  MapPin,
  Clock,
  Monitor,
  Smartphone,
  Globe,
  Ban,
  CheckCircle,
  XCircle,
  RefreshCw,
  Download,
  Filter
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { RequirePermission } from '@/components/auth/permission-guard'
import { PERMISSIONS } from '@/lib/permissions'
import { formatDate } from '@/lib/utils'
import {
  BiHubAdminCard,
  BiHubStatusBadge,
  BiHubPageHeader,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

interface SecurityEvent {
  id: string
  type: 'login_attempt' | 'failed_login' | 'suspicious_activity' | 'blocked_ip' | 'password_change'
  severity: 'low' | 'medium' | 'high' | 'critical'
  user_email?: string
  ip_address: string
  location: string
  device: string
  timestamp: string
  details: string
  status: 'resolved' | 'investigating' | 'open'
}

interface BlockedIP {
  id: string
  ip_address: string
  reason: string
  blocked_at: string
  expires_at?: string
  is_permanent: boolean
  blocked_by: string
}

export default function AdminSecurityPage() {
  const [activeTab, setActiveTab] = useState<'overview' | 'events' | 'blocked' | 'settings'>('overview')
  const [searchQuery, setSearchQuery] = useState('')
  const [filterSeverity, setFilterSeverity] = useState<string>('all')

  // Mock security events
  const securityEvents: SecurityEvent[] = [
    {
      id: '1',
      type: 'failed_login',
      severity: 'medium',
      user_email: 'admin@example.com',
      ip_address: '192.168.1.100',
      location: 'New York, US',
      device: 'Chrome on Windows',
      timestamp: '2024-01-31T10:30:00Z',
      details: 'Multiple failed login attempts (5 attempts)',
      status: 'investigating'
    },
    {
      id: '2',
      type: 'suspicious_activity',
      severity: 'high',
      ip_address: '45.123.45.67',
      location: 'Unknown',
      device: 'Bot/Crawler',
      timestamp: '2024-01-31T09:15:00Z',
      details: 'Rapid API requests from suspicious IP',
      status: 'open'
    },
    {
      id: '3',
      type: 'login_attempt',
      severity: 'low',
      user_email: 'user@example.com',
      ip_address: '10.0.0.1',
      location: 'San Francisco, US',
      device: 'Safari on iPhone',
      timestamp: '2024-01-31T08:45:00Z',
      details: 'Successful login from new device',
      status: 'resolved'
    }
  ]

  // Mock blocked IPs
  const blockedIPs: BlockedIP[] = [
    {
      id: '1',
      ip_address: '45.123.45.67',
      reason: 'Suspicious activity - Bot traffic',
      blocked_at: '2024-01-31T09:20:00Z',
      expires_at: '2024-02-07T09:20:00Z',
      is_permanent: false,
      blocked_by: 'Auto-Security System'
    },
    {
      id: '2',
      ip_address: '192.168.1.200',
      reason: 'Brute force attack attempts',
      blocked_at: '2024-01-30T15:30:00Z',
      is_permanent: true,
      blocked_by: 'Admin User'
    }
  ]

  const filteredEvents = securityEvents.filter(event => {
    const matchesSearch = 
      event.ip_address.includes(searchQuery) ||
      event.user_email?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      event.details.toLowerCase().includes(searchQuery.toLowerCase())
    
    const matchesSeverity = filterSeverity === 'all' || event.severity === filterSeverity
    
    return matchesSearch && matchesSeverity
  })

  const stats = {
    totalEvents: securityEvents.length,
    criticalEvents: securityEvents.filter(e => e.severity === 'critical').length,
    blockedIPs: blockedIPs.length,
    activeThreats: securityEvents.filter(e => e.status === 'open').length
  }

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'text-red-400 bg-red-500/20'
      case 'high':
        return 'text-orange-400 bg-orange-500/20'
      case 'medium':
        return 'text-yellow-400 bg-yellow-500/20'
      case 'low':
        return 'text-blue-400 bg-blue-500/20'
      default:
        return 'text-gray-400 bg-gray-500/20'
    }
  }

  const getEventIcon = (type: string) => {
    switch (type) {
      case 'failed_login':
        return <XCircle className="h-4 w-4 text-red-400" />
      case 'login_attempt':
        return <CheckCircle className="h-4 w-4 text-green-400" />
      case 'suspicious_activity':
        return <AlertTriangle className="h-4 w-4 text-orange-400" />
      case 'blocked_ip':
        return <Ban className="h-4 w-4 text-red-400" />
      case 'password_change':
        return <Lock className="h-4 w-4 text-blue-400" />
      default:
        return <Activity className="h-4 w-4 text-gray-400" />
    }
  }

  const handleBlockIP = (ip: string) => {
    console.log('Blocking IP:', ip)
  }

  const handleUnblockIP = (id: string) => {
    console.log('Unblocking IP:', id)
  }

  const tabs = [
    { id: 'overview', label: 'Security Overview', icon: <Shield className="h-4 w-4" /> },
    { id: 'events', label: 'Security Events', icon: <Activity className="h-4 w-4" /> },
    { id: 'blocked', label: 'Blocked IPs', icon: <Ban className="h-4 w-4" /> },
    { id: 'settings', label: 'Security Settings', icon: <Lock className="h-4 w-4" /> }
  ]

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="Security Management"
        subtitle="Monitor security events, manage threats, and configure security settings"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Security' }
        ]}
        action={
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export Report
            </Button>
          </div>
        }
      />

      {/* Tab Navigation */}
      <div className="flex space-x-1 mb-6 bg-gray-800/50 p-1 rounded-xl">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={cn(
              'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
              activeTab === tab.id
                ? 'bg-[#FF9000] text-white shadow-lg'
                : 'text-gray-400 hover:text-white hover:bg-gray-700/50'
            )}
          >
            {tab.icon}
            {tab.label}
          </button>
        ))}
      </div>

      {/* Security Overview Tab */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* Security Stats */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-blue-400/80 text-sm font-medium">Total Events</p>
                  <p className="text-2xl font-bold text-blue-100">{stats.totalEvents}</p>
                </div>
                <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
                  <Activity className="h-5 w-5 text-blue-400" />
                </div>
              </div>
            </div>

            <div className="bg-white/5 backdrop-blur-sm border border-red-300/20 rounded-xl p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-red-400/80 text-sm font-medium">Active Threats</p>
                  <p className="text-2xl font-bold text-red-100">{stats.activeThreats}</p>
                </div>
                <div className="w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center">
                  <AlertTriangle className="h-5 w-5 text-red-400" />
                </div>
              </div>
            </div>

            <div className="bg-white/5 backdrop-blur-sm border border-orange-300/20 rounded-xl p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-orange-400/80 text-sm font-medium">Blocked IPs</p>
                  <p className="text-2xl font-bold text-orange-100">{stats.blockedIPs}</p>
                </div>
                <div className="w-10 h-10 rounded-xl bg-orange-500/20 flex items-center justify-center">
                  <Ban className="h-5 w-5 text-orange-400" />
                </div>
              </div>
            </div>

            <div className="bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-green-400/80 text-sm font-medium">Security Score</p>
                  <p className="text-2xl font-bold text-green-100">85%</p>
                </div>
                <div className="w-10 h-10 rounded-xl bg-green-500/20 flex items-center justify-center">
                  <Shield className="h-5 w-5 text-green-400" />
                </div>
              </div>
            </div>
          </div>

          {/* Recent Security Events */}
          <BiHubAdminCard
            title="Recent Security Events"
            subtitle="Latest security events and alerts"
            icon={<Activity className="h-5 w-5 text-white" />}
            headerAction={
              <Button variant="outline" size="sm" onClick={() => setActiveTab('events')}>
                <Eye className="h-4 w-4 mr-2" />
                View All
              </Button>
            }
          >
            <div className="space-y-3">
              {securityEvents.slice(0, 5).map((event) => (
                <div key={event.id} className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
                  <div className="flex items-center gap-2">
                    {getEventIcon(event.type)}
                    <Badge className={cn('text-xs', getSeverityColor(event.severity))}>
                      {event.severity.toUpperCase()}
                    </Badge>
                  </div>
                  <div className="flex-1">
                    <p className="text-white">{event.details}</p>
                    <div className="flex items-center gap-2 mt-1 text-sm text-gray-400">
                      <MapPin className="h-3 w-3" />
                      <span>{event.ip_address} - {event.location}</span>
                      <span>•</span>
                      <Clock className="h-3 w-3" />
                      <span>{formatDate(event.timestamp)}</span>
                    </div>
                  </div>
                  <BiHubStatusBadge status={event.status}>
                    {event.status}
                  </BiHubStatusBadge>
                </div>
              ))}
            </div>
          </BiHubAdminCard>
        </div>
      )}

      {/* Security Events Tab */}
      {activeTab === 'events' && (
        <BiHubAdminCard
          title="Security Events"
          subtitle="All security events and monitoring alerts"
          icon={<Activity className="h-5 w-5 text-white" />}
          headerAction={
            <div className="flex items-center gap-3">
              <div className="relative">
                <Input
                  placeholder="Search events..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="w-64"
                />
              </div>
              <Select value={filterSeverity} onValueChange={setFilterSeverity}>
                <SelectTrigger className="w-32">
                  <SelectValue placeholder="Severity" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All</SelectItem>
                  <SelectItem value="critical">Critical</SelectItem>
                  <SelectItem value="high">High</SelectItem>
                  <SelectItem value="medium">Medium</SelectItem>
                  <SelectItem value="low">Low</SelectItem>
                </SelectContent>
              </Select>
            </div>
          }
        >
          <div className="space-y-4">
            {filteredEvents.map((event) => (
              <div key={event.id} className="flex items-start gap-4 p-4 bg-gray-800/50 rounded-lg">
                <div className="flex items-center gap-2">
                  {getEventIcon(event.type)}
                  <Badge className={cn('text-xs', getSeverityColor(event.severity))}>
                    {event.severity.toUpperCase()}
                  </Badge>
                </div>
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="text-white font-medium">{event.details}</h3>
                    <BiHubStatusBadge status={event.status}>
                      {event.status}
                    </BiHubStatusBadge>
                  </div>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm text-gray-400">
                    <div className="flex items-center gap-1">
                      <Globe className="h-3 w-3" />
                      <span>{event.ip_address}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <MapPin className="h-3 w-3" />
                      <span>{event.location}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Monitor className="h-3 w-3" />
                      <span>{event.device}</span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Clock className="h-3 w-3" />
                      <span>{formatDate(event.timestamp)}</span>
                    </div>
                  </div>
                  {event.user_email && (
                    <div className="mt-2 text-sm text-gray-400">
                      User: {event.user_email}
                    </div>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="sm">
                    <Eye className="h-4 w-4" />
                  </Button>
                  <Button 
                    variant="ghost" 
                    size="sm"
                    onClick={() => handleBlockIP(event.ip_address)}
                    className="text-red-400 hover:text-red-300"
                  >
                    <Ban className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      )}

      {/* Blocked IPs Tab */}
      {activeTab === 'blocked' && (
        <BiHubAdminCard
          title="Blocked IP Addresses"
          subtitle="Manage blocked and restricted IP addresses"
          icon={<Ban className="h-5 w-5 text-white" />}
        >
          <div className="space-y-4">
            {blockedIPs.map((blockedIP) => (
              <div key={blockedIP.id} className="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg">
                <div>
                  <div className="flex items-center gap-3">
                    <h3 className="text-white font-medium">{blockedIP.ip_address}</h3>
                    <Badge variant={blockedIP.is_permanent ? 'destructive' : 'secondary'}>
                      {blockedIP.is_permanent ? 'Permanent' : 'Temporary'}
                    </Badge>
                  </div>
                  <p className="text-gray-400 text-sm mt-1">{blockedIP.reason}</p>
                  <div className="flex items-center gap-4 mt-2 text-sm text-gray-400">
                    <span>Blocked: {formatDate(blockedIP.blocked_at)}</span>
                    {blockedIP.expires_at && (
                      <>
                        <span>•</span>
                        <span>Expires: {formatDate(blockedIP.expires_at)}</span>
                      </>
                    )}
                    <span>•</span>
                    <span>By: {blockedIP.blocked_by}</span>
                  </div>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleUnblockIP(blockedIP.id)}
                  className="text-green-400 hover:text-green-300"
                >
                  <Unlock className="h-4 w-4 mr-2" />
                  Unblock
                </Button>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      )}

      {/* Security Settings Tab */}
      {activeTab === 'settings' && (
        <div className="space-y-6">
          <BiHubAdminCard
            title="Security Configuration"
            subtitle="Configure security policies and settings"
            icon={<Lock className="h-5 w-5 text-white" />}
          >
            <div className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <h3 className="text-white font-medium">Login Security</h3>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">Two-Factor Authentication</span>
                      <Badge variant="secondary">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">Password Complexity</span>
                      <Badge variant="secondary">Strong</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">Session Timeout</span>
                      <Badge variant="secondary">30 minutes</Badge>
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  <h3 className="text-white font-medium">Threat Protection</h3>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">Rate Limiting</span>
                      <Badge variant="secondary">Active</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">DDoS Protection</span>
                      <Badge variant="secondary">Enabled</Badge>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                      <span className="text-gray-300">IP Blocking</span>
                      <Badge variant="secondary">Automatic</Badge>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </BiHubAdminCard>
        </div>
      )}
    </div>
  )
}
