'use client'

import { useState } from 'react'
import {
  Server,
  Database,
  HardDrive,
  Cpu,
  MemoryStick,
  Activity,
  AlertTriangle,
  CheckCircle,
  Download,
  Upload,
  Trash2,
  RefreshCw,
  Settings,
  Monitor,
  Wifi,
  Clock,
  FileText,
  Shield,
  Zap
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
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

interface SystemMetric {
  name: string
  value: number
  unit: string
  status: 'good' | 'warning' | 'critical'
  icon: React.ReactNode
}

interface SystemLog {
  id: string
  level: 'info' | 'warning' | 'error' | 'debug'
  message: string
  timestamp: string
  source: string
}

interface BackupInfo {
  id: string
  name: string
  size: string
  created_at: string
  status: 'completed' | 'failed' | 'in_progress'
  type: 'manual' | 'automatic'
}

export default function AdminSystemPage() {
  const [activeTab, setActiveTab] = useState<'overview' | 'logs' | 'backups' | 'maintenance'>('overview')
  const [isBackingUp, setIsBackingUp] = useState(false)

  // Mock system metrics
  const systemMetrics: SystemMetric[] = [
    {
      name: 'CPU Usage',
      value: 45,
      unit: '%',
      status: 'good',
      icon: <Cpu className="h-5 w-5" />
    },
    {
      name: 'Memory Usage',
      value: 68,
      unit: '%',
      status: 'warning',
      icon: <MemoryStick className="h-5 w-5" />
    },
    {
      name: 'Disk Usage',
      value: 82,
      unit: '%',
      status: 'critical',
      icon: <HardDrive className="h-5 w-5" />
    },
    {
      name: 'Network I/O',
      value: 23,
      unit: 'MB/s',
      status: 'good',
      icon: <Wifi className="h-5 w-5" />
    }
  ]

  // Mock system logs
  const systemLogs: SystemLog[] = [
    {
      id: '1',
      level: 'info',
      message: 'Database backup completed successfully',
      timestamp: '2024-01-31T10:30:00Z',
      source: 'backup-service'
    },
    {
      id: '2',
      level: 'warning',
      message: 'High memory usage detected (85%)',
      timestamp: '2024-01-31T10:25:00Z',
      source: 'monitoring'
    },
    {
      id: '3',
      level: 'error',
      message: 'Failed to connect to external payment service',
      timestamp: '2024-01-31T10:20:00Z',
      source: 'payment-gateway'
    }
  ]

  // Mock backup data
  const backups: BackupInfo[] = [
    {
      id: '1',
      name: 'Daily Backup - 2024-01-31',
      size: '2.4 GB',
      created_at: '2024-01-31T02:00:00Z',
      status: 'completed',
      type: 'automatic'
    },
    {
      id: '2',
      name: 'Manual Backup - Pre-Update',
      size: '2.3 GB',
      created_at: '2024-01-30T14:30:00Z',
      status: 'completed',
      type: 'manual'
    }
  ]

  const getMetricColor = (status: string) => {
    switch (status) {
      case 'good':
        return 'text-green-400'
      case 'warning':
        return 'text-yellow-400'
      case 'critical':
        return 'text-red-400'
      default:
        return 'text-gray-400'
    }
  }

  const getLogLevelColor = (level: string) => {
    switch (level) {
      case 'error':
        return 'text-red-400 bg-red-500/20'
      case 'warning':
        return 'text-yellow-400 bg-yellow-500/20'
      case 'info':
        return 'text-blue-400 bg-blue-500/20'
      case 'debug':
        return 'text-gray-400 bg-gray-500/20'
      default:
        return 'text-gray-400 bg-gray-500/20'
    }
  }

  const handleBackup = async () => {
    setIsBackingUp(true)
    // Simulate backup process
    setTimeout(() => {
      setIsBackingUp(false)
    }, 3000)
  }

  const handleCleanup = () => {
    console.log('Starting system cleanup...')
  }

  const tabs = [
    { id: 'overview', label: 'System Overview', icon: <Monitor className="h-4 w-4" /> },
    { id: 'logs', label: 'System Logs', icon: <FileText className="h-4 w-4" /> },
    { id: 'backups', label: 'Backups', icon: <Database className="h-4 w-4" /> },
    { id: 'maintenance', label: 'Maintenance', icon: <Settings className="h-4 w-4" /> }
  ]

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="System Management"
        subtitle="Monitor system health, manage backups, and perform maintenance"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'System' }
        ]}
        action={
          <div className="flex items-center gap-3">
            <Button variant="outline" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Refresh
            </Button>
            <RequirePermission permission={PERMISSIONS.SYSTEM_BACKUP}>
              <Button 
                onClick={handleBackup}
                disabled={isBackingUp}
                className={BIHUB_ADMIN_THEME.components.button.primary}
              >
                {isBackingUp ? (
                  <RefreshCw className="mr-2 h-5 w-5 animate-spin" />
                ) : (
                  <Download className="mr-2 h-5 w-5" />
                )}
                {isBackingUp ? 'Backing Up...' : 'Create Backup'}
              </Button>
            </RequirePermission>
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

      {/* System Overview Tab */}
      {activeTab === 'overview' && (
        <div className="space-y-6">
          {/* System Metrics */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            {systemMetrics.map((metric, index) => (
              <div key={index} className="bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4">
                <div className="flex items-center justify-between mb-3">
                  <div className={cn('w-10 h-10 rounded-xl flex items-center justify-center', getMetricColor(metric.status))}>
                    {metric.icon}
                  </div>
                  <BiHubStatusBadge status={metric.status}>
                    {metric.status}
                  </BiHubStatusBadge>
                </div>
                <h3 className="text-white font-semibold">{metric.name}</h3>
                <div className="flex items-center gap-2 mt-2">
                  <span className={cn('text-2xl font-bold', getMetricColor(metric.status))}>
                    {metric.value}
                  </span>
                  <span className="text-gray-400 text-sm">{metric.unit}</span>
                </div>
                {metric.unit === '%' && (
                  <Progress 
                    value={metric.value} 
                    className="mt-3"
                    // className={cn('mt-3', getMetricColor(metric.status))}
                  />
                )}
              </div>
            ))}
          </div>

          {/* System Status */}
          <BiHubAdminCard
            title="System Status"
            subtitle="Current system health and services"
            icon={<Activity className="h-5 w-5 text-white" />}
          >
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[
                { name: 'Database', status: 'online', uptime: '99.9%' },
                { name: 'API Server', status: 'online', uptime: '99.8%' },
                { name: 'File Storage', status: 'online', uptime: '100%' },
                { name: 'Payment Gateway', status: 'warning', uptime: '98.5%' },
                { name: 'Email Service', status: 'online', uptime: '99.7%' },
                { name: 'CDN', status: 'online', uptime: '99.9%' }
              ].map((service, index) => (
                <div key={index} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <div className={cn(
                      'w-3 h-3 rounded-full',
                      service.status === 'online' ? 'bg-green-500' : 
                      service.status === 'warning' ? 'bg-yellow-500' : 'bg-red-500'
                    )}></div>
                    <span className="text-white font-medium">{service.name}</span>
                  </div>
                  <span className="text-gray-400 text-sm">{service.uptime}</span>
                </div>
              ))}
            </div>
          </BiHubAdminCard>
        </div>
      )}

      {/* System Logs Tab */}
      {activeTab === 'logs' && (
        <BiHubAdminCard
          title="System Logs"
          subtitle="Recent system events and messages"
          icon={<FileText className="h-5 w-5 text-white" />}
          headerAction={
            <Button variant="outline" size="sm">
              <Download className="h-4 w-4 mr-2" />
              Export Logs
            </Button>
          }
        >
          <div className="space-y-3">
            {systemLogs.map((log) => (
              <div key={log.id} className="flex items-start gap-3 p-3 bg-gray-800/50 rounded-lg">
                <Badge className={cn('text-xs', getLogLevelColor(log.level))}>
                  {log.level.toUpperCase()}
                </Badge>
                <div className="flex-1">
                  <p className="text-white">{log.message}</p>
                  <div className="flex items-center gap-2 mt-1 text-sm text-gray-400">
                    <Clock className="h-3 w-3" />
                    <span>{formatDate(log.timestamp)}</span>
                    <span>•</span>
                    <span>{log.source}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      )}

      {/* Backups Tab */}
      {activeTab === 'backups' && (
        <BiHubAdminCard
          title="Database Backups"
          subtitle="Manage system backups and restore points"
          icon={<Database className="h-5 w-5 text-white" />}
          headerAction={
            <RequirePermission permission={PERMISSIONS.SYSTEM_BACKUP}>
              <Button 
                onClick={handleBackup}
                disabled={isBackingUp}
                className={BIHUB_ADMIN_THEME.components.button.primary}
              >
                {isBackingUp ? (
                  <RefreshCw className="mr-2 h-5 w-5 animate-spin" />
                ) : (
                  <Download className="mr-2 h-5 w-5" />
                )}
                {isBackingUp ? 'Creating...' : 'Create Backup'}
              </Button>
            </RequirePermission>
          }
        >
          <div className="space-y-4">
            {backups.map((backup) => (
              <div key={backup.id} className="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
                    <Database className="h-5 w-5 text-blue-400" />
                  </div>
                  <div>
                    <h3 className="text-white font-medium">{backup.name}</h3>
                    <div className="flex items-center gap-2 text-sm text-gray-400">
                      <span>{backup.size}</span>
                      <span>•</span>
                      <span>{formatDate(backup.created_at)}</span>
                      <span>•</span>
                      <BiHubStatusBadge status={backup.type}>
                        {backup.type}
                      </BiHubStatusBadge>
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="sm">
                    <Download className="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" className="text-red-400 hover:text-red-300">
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </BiHubAdminCard>
      )}

      {/* Maintenance Tab */}
      {activeTab === 'maintenance' && (
        <div className="space-y-6">
          <BiHubAdminCard
            title="System Maintenance"
            subtitle="Perform system cleanup and optimization tasks"
            icon={<Settings className="h-5 w-5 text-white" />}
          >
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="p-4 bg-gray-800/50 rounded-lg">
                <div className="flex items-center gap-3 mb-3">
                  <Trash2 className="h-5 w-5 text-red-400" />
                  <h3 className="text-white font-medium">System Cleanup</h3>
                </div>
                <p className="text-gray-400 text-sm mb-4">
                  Remove temporary files, clear caches, and optimize database
                </p>
                <Button onClick={handleCleanup} variant="outline" size="sm">
                  <Zap className="h-4 w-4 mr-2" />
                  Run Cleanup
                </Button>
              </div>

              <div className="p-4 bg-gray-800/50 rounded-lg">
                <div className="flex items-center gap-3 mb-3">
                  <Shield className="h-5 w-5 text-green-400" />
                  <h3 className="text-white font-medium">Security Scan</h3>
                </div>
                <p className="text-gray-400 text-sm mb-4">
                  Scan for security vulnerabilities and malware
                </p>
                <Button variant="outline" size="sm">
                  <Shield className="h-4 w-4 mr-2" />
                  Start Scan
                </Button>
              </div>
            </div>
          </BiHubAdminCard>
        </div>
      )}
    </div>
  )
}
