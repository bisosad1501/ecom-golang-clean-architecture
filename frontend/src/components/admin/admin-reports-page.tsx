'use client'

import { useState } from 'react'
import {
  FileText,
  Download,
  Calendar,
  Filter,
  BarChart3,
  TrendingUp,
  DollarSign,
  ShoppingCart,
  Users,
  Package,
  Clock,
  CheckCircle,
  AlertCircle,
  RefreshCw,
  Plus,
  Eye
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
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

interface Report {
  id: string
  name: string
  type: 'sales' | 'inventory' | 'customers' | 'analytics' | 'financial'
  format: 'pdf' | 'excel' | 'csv'
  status: 'pending' | 'processing' | 'completed' | 'failed'
  progress: number
  created_at: string
  completed_at?: string
  download_url?: string
  file_size?: string
  created_by: string
}

export default function AdminReportsPage() {
  const [selectedType, setSelectedType] = useState<string>('all')
  const [selectedStatus, setSelectedStatus] = useState<string>('all')
  const [showGenerateForm, setShowGenerateForm] = useState(false)

  // Mock data - replace with real API calls
  const reports: Report[] = [
    {
      id: '1',
      name: 'Monthly Sales Report - January 2024',
      type: 'sales',
      format: 'pdf',
      status: 'completed',
      progress: 100,
      created_at: '2024-01-31T10:30:00Z',
      completed_at: '2024-01-31T10:35:00Z',
      download_url: '/reports/sales-jan-2024.pdf',
      file_size: '2.4 MB',
      created_by: 'Admin User'
    },
    {
      id: '2',
      name: 'Inventory Status Report',
      type: 'inventory',
      format: 'excel',
      status: 'processing',
      progress: 65,
      created_at: '2024-01-31T14:20:00Z',
      created_by: 'Admin User'
    },
    {
      id: '3',
      name: 'Customer Analytics Q4 2023',
      type: 'customers',
      format: 'pdf',
      status: 'failed',
      progress: 0,
      created_at: '2024-01-30T09:15:00Z',
      created_by: 'Admin User'
    }
  ]

  const filteredReports = reports.filter(report => {
    const matchesType = selectedType === 'all' || report.type === selectedType
    const matchesStatus = selectedStatus === 'all' || report.status === selectedStatus
    return matchesType && matchesStatus
  })

  const stats = {
    total: reports.length,
    completed: reports.filter(r => r.status === 'completed').length,
    processing: reports.filter(r => r.status === 'processing').length,
    failed: reports.filter(r => r.status === 'failed').length
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'completed':
        return <CheckCircle className="h-4 w-4 text-green-400" />
      case 'processing':
        return <RefreshCw className="h-4 w-4 text-blue-400 animate-spin" />
      case 'pending':
        return <Clock className="h-4 w-4 text-yellow-400" />
      case 'failed':
        return <AlertCircle className="h-4 w-4 text-red-400" />
      default:
        return <Clock className="h-4 w-4 text-gray-400" />
    }
  }

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'sales':
        return <DollarSign className="h-5 w-5 text-green-400" />
      case 'inventory':
        return <Package className="h-5 w-5 text-blue-400" />
      case 'customers':
        return <Users className="h-5 w-5 text-purple-400" />
      case 'analytics':
        return <BarChart3 className="h-5 w-5 text-orange-400" />
      case 'financial':
        return <TrendingUp className="h-5 w-5 text-emerald-400" />
      default:
        return <FileText className="h-5 w-5 text-gray-400" />
    }
  }

  const handleDownload = (report: Report) => {
    if (report.download_url) {
      window.open(report.download_url, '_blank')
    }
  }

  const handleGenerateReport = () => {
    setShowGenerateForm(true)
  }

  const reportTypes = [
    { value: 'sales', label: 'Sales Reports' },
    { value: 'inventory', label: 'Inventory Reports' },
    { value: 'customers', label: 'Customer Reports' },
    { value: 'analytics', label: 'Analytics Reports' },
    { value: 'financial', label: 'Financial Reports' }
  ]

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="Reports & Analytics"
        subtitle="Generate and download detailed business reports"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Reports' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
            <Button
              onClick={handleGenerateReport}
              className={BIHUB_ADMIN_THEME.components.button.primary}
            >
              <Plus className="mr-2 h-5 w-5" />
              Generate Report
            </Button>
          </RequirePermission>
        }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div className="bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-400/80 text-sm font-medium">Total Reports</p>
              <p className="text-2xl font-bold text-blue-100">{stats.total}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
              <FileText className="h-5 w-5 text-blue-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-400/80 text-sm font-medium">Completed</p>
              <p className="text-2xl font-bold text-green-100">{stats.completed}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-green-500/20 flex items-center justify-center">
              <CheckCircle className="h-5 w-5 text-green-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-yellow-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-yellow-400/80 text-sm font-medium">Processing</p>
              <p className="text-2xl font-bold text-yellow-100">{stats.processing}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-yellow-500/20 flex items-center justify-center">
              <RefreshCw className="h-5 w-5 text-yellow-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-red-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-red-400/80 text-sm font-medium">Failed</p>
              <p className="text-2xl font-bold text-red-100">{stats.failed}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-red-500/20 flex items-center justify-center">
              <AlertCircle className="h-5 w-5 text-red-400" />
            </div>
          </div>
        </div>
      </div>

      {/* Quick Generate Section */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        {reportTypes.map((type) => (
          <div
            key={type.value}
            className="group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 hover:scale-[1.02] transition-all duration-200 cursor-pointer"
            onClick={handleGenerateReport}
          >
            <div className="flex items-center gap-3">
              <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                {getTypeIcon(type.value)}
              </div>
              <div>
                <h3 className="text-lg font-semibold text-white">{type.label}</h3>
                <p className="text-gray-400 text-sm">Generate new {type.label.toLowerCase()}</p>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Reports List */}
      <BiHubAdminCard
        title="Report History"
        subtitle="View and download generated reports"
        icon={<FileText className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-3">
            <Select value={selectedType} onValueChange={setSelectedType}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Filter by type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                {reportTypes.map((type) => (
                  <SelectItem key={type.value} value={type.value}>
                    {type.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            
            <Select value={selectedStatus} onValueChange={setSelectedStatus}>
              <SelectTrigger className="w-40">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="processing">Processing</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>
          </div>
        }
      >
        <div className="space-y-4">
          {filteredReports.map((report) => (
            <div
              key={report.id}
              className="group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-[#FF9000] to-[#e67e00] flex items-center justify-center">
                    {getTypeIcon(report.type)}
                  </div>
                  
                  <div className="flex-1">
                    <div className="flex items-center gap-3">
                      <h3 className="text-lg font-semibold text-white">{report.name}</h3>
                      <BiHubStatusBadge status={report.status}>
                        {report.status}
                      </BiHubStatusBadge>
                      <Badge variant="outline" className="text-xs">
                        {report.format.toUpperCase()}
                      </Badge>
                    </div>
                    
                    <div className="flex items-center gap-4 mt-2 text-sm text-gray-400">
                      <div className="flex items-center gap-1">
                        {getStatusIcon(report.status)}
                        <span>
                          {report.status === 'processing' ? `${report.progress}% complete` : 
                           report.status === 'completed' ? `Completed ${formatDate(report.completed_at!)}` :
                           report.status === 'failed' ? 'Generation failed' :
                           'Pending'}
                        </span>
                      </div>
                      <span>•</span>
                      <span>Created {formatDate(report.created_at)}</span>
                      <span>•</span>
                      <span>By {report.created_by}</span>
                      {report.file_size && (
                        <>
                          <span>•</span>
                          <span>{report.file_size}</span>
                        </>
                      )}
                    </div>

                    {report.status === 'processing' && (
                      <div className="mt-3">
                        <div className="w-full bg-gray-700 rounded-full h-2">
                          <div 
                            className="bg-blue-500 h-2 rounded-full transition-all duration-300"
                            style={{ width: `${report.progress}%` }}
                          ></div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {report.status === 'completed' && report.download_url && (
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={() => handleDownload(report)}
                      className="opacity-0 group-hover:opacity-100 transition-opacity"
                    >
                      <Download className="h-4 w-4" />
                    </Button>
                  )}
                  
                  <Button variant="ghost" size="sm">
                    <Eye className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </div>
          ))}

          {filteredReports.length === 0 && (
            <div className="text-center py-12">
              <FileText className="h-12 w-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-300 mb-2">No reports found</h3>
              <p className="text-gray-400 mb-6">
                Generate your first report to get started with analytics.
              </p>
              <RequirePermission permission={PERMISSIONS.ANALYTICS_VIEW}>
                <Button
                  onClick={handleGenerateReport}
                  className={BIHUB_ADMIN_THEME.components.button.primary}
                >
                  <Plus className="mr-2 h-5 w-5" />
                  Generate First Report
                </Button>
              </RequirePermission>
            </div>
          )}
        </div>
      </BiHubAdminCard>
    </div>
  )
}
