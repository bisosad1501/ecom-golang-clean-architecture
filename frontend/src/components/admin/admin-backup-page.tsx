'use client';

import React, { useState, useEffect } from 'react';
import {
  Database,
  Download,
  Upload,
  MoreVertical,
  Calendar,
  Clock,
  CheckCircle,
  AlertCircle,
  XCircle,
  RefreshCw,
  HardDrive,
  FileText
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';

// BiHub Components
import {
  BiHubAdminCard,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
  BiHubStatusBadge,
} from './bihub-admin-components';
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme';
import { RequirePermission } from '@/components/auth/permission-guard';
import { PERMISSIONS } from '@/lib/permissions';

interface Backup {
  id: string;
  name: string;
  description?: string;
  type: 'full' | 'incremental' | 'differential';
  status: 'completed' | 'in_progress' | 'failed' | 'scheduled';
  size: number;
  created_at: string;
  duration?: number;
  file_path?: string;
  tables_count: number;
  records_count: number;
}

export function AdminBackupPage() {
  const [backups, setBackups] = useState<Backup[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    type: 'full' as Backup['type']
  });

  // Mock data for demonstration
  useEffect(() => {
    const mockBackups: Backup[] = [
      {
        id: '1',
        name: 'Daily Backup - Jan 19',
        description: 'Automated daily backup',
        type: 'full',
        status: 'completed',
        size: 2.5 * 1024 * 1024 * 1024, // 2.5 GB
        created_at: '2024-01-19T02:00:00Z',
        duration: 1800, // 30 minutes
        file_path: '/backups/daily_backup_20240119.sql',
        tables_count: 25,
        records_count: 150000
      },
      {
        id: '2',
        name: 'Weekly Backup - Jan 14',
        description: 'Weekly full system backup',
        type: 'full',
        status: 'completed',
        size: 3.2 * 1024 * 1024 * 1024, // 3.2 GB
        created_at: '2024-01-14T01:00:00Z',
        duration: 2400, // 40 minutes
        file_path: '/backups/weekly_backup_20240114.sql',
        tables_count: 25,
        records_count: 145000
      },
      {
        id: '3',
        name: 'Manual Backup - Before Update',
        description: 'Manual backup before system update',
        type: 'full',
        status: 'completed',
        size: 2.8 * 1024 * 1024 * 1024, // 2.8 GB
        created_at: '2024-01-12T10:30:00Z',
        duration: 2100, // 35 minutes
        file_path: '/backups/manual_backup_20240112.sql',
        tables_count: 25,
        records_count: 142000
      },
      {
        id: '4',
        name: 'Incremental Backup - Jan 18',
        description: 'Incremental backup',
        type: 'incremental',
        status: 'failed',
        size: 0,
        created_at: '2024-01-18T02:00:00Z',
        tables_count: 0,
        records_count: 0
      }
    ];
    
    setTimeout(() => {
      setBackups(mockBackups);
      setIsLoading(false);
    }, 1000);
  }, []);

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'completed':
        return <Badge className="bg-green-100 text-green-800"><CheckCircle className="w-3 h-3 mr-1" />Completed</Badge>;
      case 'in_progress':
        return <Badge className="bg-blue-100 text-blue-800"><RefreshCw className="w-3 h-3 mr-1 animate-spin" />In Progress</Badge>;
      case 'failed':
        return <Badge className="bg-red-100 text-red-800"><XCircle className="w-3 h-3 mr-1" />Failed</Badge>;
      case 'scheduled':
        return <Badge className="bg-yellow-100 text-yellow-800"><Clock className="w-3 h-3 mr-1" />Scheduled</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
    }
  };

  const getTypeBadge = (type: string) => {
    switch (type) {
      case 'full':
        return <Badge variant="outline" className="bg-blue-50 text-blue-700">Full</Badge>;
      case 'incremental':
        return <Badge variant="outline" className="bg-green-50 text-green-700">Incremental</Badge>;
      case 'differential':
        return <Badge variant="outline" className="bg-purple-50 text-purple-700">Differential</Badge>;
      default:
        return <Badge variant="outline">{type}</Badge>;
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const formatDuration = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    if (hours > 0) {
      return `${hours}h ${minutes}m ${secs}s`;
    } else if (minutes > 0) {
      return `${minutes}m ${secs}s`;
    } else {
      return `${secs}s`;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const handleCreateBackup = async () => {
    if (!formData.name.trim()) {
      toast.error('Please enter a backup name');
      return;
    }

    setIsCreating(true);
    try {
      // TODO: Implement API call
      const newBackup: Backup = {
        id: Date.now().toString(),
        name: formData.name,
        description: formData.description,
        type: formData.type,
        status: 'in_progress',
        size: 0,
        created_at: new Date().toISOString(),
        tables_count: 0,
        records_count: 0
      };

      setBackups(prev => [newBackup, ...prev]);
      toast.success('Backup started successfully');
      setIsCreateDialogOpen(false);
      setFormData({ name: '', description: '', type: 'full' });

      // Simulate backup completion
      setTimeout(() => {
        setBackups(prev => prev.map(backup => 
          backup.id === newBackup.id 
            ? {
                ...backup,
                status: 'completed' as const,
                size: 2.1 * 1024 * 1024 * 1024,
                duration: 1500,
                tables_count: 25,
                records_count: 148000,
                file_path: `/backups/${formData.name.toLowerCase().replace(/\s+/g, '_')}.sql`
              }
            : backup
        ));
        toast.success('Backup completed successfully');
      }, 5000);
    } catch (error) {
      toast.error('Failed to create backup');
    } finally {
      setIsCreating(false);
    }
  };

  const handleDownloadBackup = (backup: Backup) => {
    if (backup.status !== 'completed' || !backup.file_path) {
      toast.error('Backup file is not available for download');
      return;
    }
    
    // TODO: Implement actual download
    toast.success(`Downloading ${backup.name}...`);
  };

  const handleDeleteBackup = async (backupId: string) => {
    try {
      // TODO: Implement API call
      setBackups(prev => prev.filter(b => b.id !== backupId));
      toast.success('Backup deleted successfully');
    } catch (error) {
      toast.error('Failed to delete backup');
    }
  };

  const totalBackupSize = backups.reduce((sum, backup) => sum + backup.size, 0);
  const completedBackups = backups.filter(b => b.status === 'completed').length;
  const failedBackups = backups.filter(b => b.status === 'failed').length;

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Backup Management"
        subtitle="Create and manage database backups for your BiHub system"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'System' },
          { label: 'Backup' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.SYSTEM_BACKUP}>
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
              <DialogTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.primary}>
                  <Database className="w-4 h-4 mr-2" />
                  Create Backup
                </Button>
              </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Create New Backup</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="name">Backup Name</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter backup name"
                />
              </div>
              <div>
                <Label htmlFor="description">Description (Optional)</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Enter backup description"
                  rows={3}
                />
              </div>
              <div>
                <Label>Backup Type</Label>
                <Select value={formData.type} onValueChange={(value: Backup['type']) => setFormData({ ...formData, type: value })}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="full">Full Backup</SelectItem>
                    <SelectItem value="incremental">Incremental Backup</SelectItem>
                    <SelectItem value="differential">Differential Backup</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={handleCreateBackup} disabled={isCreating || !formData.name.trim()}>
                  {isCreating ? 'Creating...' : 'Create Backup'}
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
        </RequirePermission>
      }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <BiHubStatCard
          title="Total Backups"
          value={backups.length.toString()}
          icon={<Database className="h-5 w-5 text-white" />}
          change={3}
          color="primary"
        />

        <BiHubStatCard
          title="Completed"
          value={completedBackups.toString()}
          icon={<CheckCircle className="h-5 w-5 text-white" />}
          change={2}
          color="success"
        />

        <BiHubStatCard
          title="Failed"
          value={failedBackups.toString()}
          icon={<XCircle className="h-5 w-5 text-white" />}
          change={-1}
          color="error"
        />

        <BiHubStatCard
          title="Total Size"
          value={formatFileSize(totalBackupSize)}
          icon={<HardDrive className="h-5 w-5 text-white" />}
          change={15}
          color="info"
        />
      </div>

      {/* Backup Schedule */}
      <BiHubAdminCard
        title="Backup Schedule"
        subtitle="Automated backup schedule configuration"
        icon={<Calendar className="h-5 w-5 text-white" />}
      >
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="bg-gray-700 p-4 rounded-lg border border-gray-600">
            <h4 className="font-medium text-white mb-2">Daily Backup</h4>
            <p className="text-sm text-gray-300">Every day at 2:00 AM</p>
            <p className="text-xs text-gray-400 mt-1">Next: Tomorrow 2:00 AM</p>
          </div>
          <div className="bg-gray-700 p-4 rounded-lg border border-gray-600">
            <h4 className="font-medium text-white mb-2">Weekly Backup</h4>
            <p className="text-sm text-gray-300">Every Sunday at 1:00 AM</p>
            <p className="text-xs text-gray-400 mt-1">Next: Sunday 1:00 AM</p>
          </div>
          <div className="bg-gray-700 p-4 rounded-lg border border-gray-600">
            <h4 className="font-medium text-white mb-2">Monthly Backup</h4>
            <p className="text-sm text-gray-300">1st of every month at 12:00 AM</p>
            <p className="text-xs text-gray-400 mt-1">Next: Feb 1st 12:00 AM</p>
          </div>
        </div>
      </BiHubAdminCard>

      {/* Backup History */}
      <BiHubAdminCard
        title="Backup History"
        subtitle="View and manage all database backups"
        icon={<Database className="h-5 w-5 text-white" />}
      >
        <div className="space-y-4">
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 4 }, (_, i) => (
                <div key={i} className="animate-pulse">
                  <div className="h-32 bg-gray-700 rounded-lg"></div>
                </div>
              ))}
            </div>
          ) : backups.length === 0 ? (
            <BiHubEmptyState
              icon={<Database className="h-12 w-12 text-gray-400" />}
              title="No backups found"
              description="Create your first backup to get started"
            />
          ) : (
          backups.map((backup) => (
            <div key={backup.id} className="p-6 bg-gray-800 rounded-lg border border-gray-700 hover:border-gray-600 transition-colors">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center text-white">
                    <Database className="w-6 h-6" />
                  </div>

                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h4 className="font-semibold text-white">{backup.name}</h4>
                      <BiHubStatusBadge status={backup.status}>
                        {backup.status.charAt(0).toUpperCase() + backup.status.slice(1).replace('_', ' ')}
                      </BiHubStatusBadge>
                      {getTypeBadge(backup.type)}
                    </div>

                    {backup.description && (
                      <p className="text-gray-300 mb-2">{backup.description}</p>
                    )}

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-300">
                      <div>
                        <p><strong>Created:</strong> {formatDate(backup.created_at)}</p>
                        <p><strong>Size:</strong> {formatFileSize(backup.size)}</p>
                        {backup.duration && (
                          <p><strong>Duration:</strong> {formatDuration(backup.duration)}</p>
                        )}
                      </div>
                      <div>
                        <p><strong>Tables:</strong> {backup.tables_count}</p>
                        <p><strong>Records:</strong> {backup.records_count.toLocaleString()}</p>
                        {backup.file_path && (
                          <p className="text-xs text-gray-400 truncate">
                            <strong>Path:</strong> {backup.file_path}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>

                <RequirePermission permission={PERMISSIONS.SYSTEM_BACKUP}>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      {backup.status === 'completed' && (
                        <DropdownMenuItem onClick={() => handleDownloadBackup(backup)}>
                          <Download className="w-4 h-4 mr-2" />
                          Download
                        </DropdownMenuItem>
                      )}
                      <DropdownMenuItem
                        onClick={() => handleDeleteBackup(backup.id)}
                        className="text-red-400"
                      >
                        <XCircle className="w-4 h-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </RequirePermission>
              </div>

              {backup.status === 'in_progress' && (
                <div className="w-full bg-gray-700 rounded-full h-2">
                  <div className="bg-[#FF9000] h-2 rounded-full animate-pulse" style={{ width: '45%' }} />
                </div>
              )}
            </div>
          ))
        )}
        </div>
      </BiHubAdminCard>
    </div>
  );
}
