'use client'

import { useState } from 'react'
import {
  Upload,
  FileImage,
  FileText,
  File,
  Trash2,
  Download,
  Eye,
  Search,
  Filter,
  Grid,
  List,
  FolderOpen,
  Image,
  Video,
  Music,
  Archive,
  MoreHorizontal,
  Copy,
  Share,
  Edit
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
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
  BiHubPageHeader,
} from './bihub-admin-components'
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme'
import { cn } from '@/lib/utils'

interface FileItem {
  id: string
  name: string
  type: 'image' | 'document' | 'video' | 'audio' | 'archive' | 'other'
  size: string
  url: string
  thumbnail?: string
  uploaded_at: string
  uploaded_by: string
  folder?: string
  mime_type: string
}

export default function AdminFilesPage() {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [searchQuery, setSearchQuery] = useState('')
  const [filterType, setFilterType] = useState<string>('all')
  const [selectedFiles, setSelectedFiles] = useState<string[]>([])

  // Mock file data
  const files: FileItem[] = [
    {
      id: '1',
      name: 'product-hero-banner.jpg',
      type: 'image',
      size: '2.4 MB',
      url: '/uploads/images/product-hero-banner.jpg',
      thumbnail: '/uploads/thumbnails/product-hero-banner.jpg',
      uploaded_at: '2024-01-31T10:30:00Z',
      uploaded_by: 'Admin User',
      folder: 'banners',
      mime_type: 'image/jpeg'
    },
    {
      id: '2',
      name: 'user-manual.pdf',
      type: 'document',
      size: '1.8 MB',
      url: '/uploads/documents/user-manual.pdf',
      uploaded_at: '2024-01-30T14:20:00Z',
      uploaded_by: 'Content Manager',
      folder: 'documents',
      mime_type: 'application/pdf'
    },
    {
      id: '3',
      name: 'product-demo.mp4',
      type: 'video',
      size: '15.6 MB',
      url: '/uploads/videos/product-demo.mp4',
      thumbnail: '/uploads/thumbnails/product-demo.jpg',
      uploaded_at: '2024-01-29T09:15:00Z',
      uploaded_by: 'Marketing Team',
      folder: 'videos',
      mime_type: 'video/mp4'
    },
    {
      id: '4',
      name: 'logo-variations.zip',
      type: 'archive',
      size: '5.2 MB',
      url: '/uploads/archives/logo-variations.zip',
      uploaded_at: '2024-01-28T16:45:00Z',
      uploaded_by: 'Design Team',
      folder: 'assets',
      mime_type: 'application/zip'
    }
  ]

  const filteredFiles = files.filter(file => {
    const matchesSearch = file.name.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesType = filterType === 'all' || file.type === filterType
    return matchesSearch && matchesType
  })

  const stats = {
    totalFiles: files.length,
    totalSize: '24.8 MB',
    images: files.filter(f => f.type === 'image').length,
    documents: files.filter(f => f.type === 'document').length
  }

  const getFileIcon = (type: string, size: 'sm' | 'lg' = 'sm') => {
    const iconSize = size === 'lg' ? 'h-8 w-8' : 'h-5 w-5'
    
    switch (type) {
      case 'image':
        return <Image className={cn(iconSize, 'text-green-400')} />
      case 'document':
        return <FileText className={cn(iconSize, 'text-blue-400')} />
      case 'video':
        return <Video className={cn(iconSize, 'text-purple-400')} />
      case 'audio':
        return <Music className={cn(iconSize, 'text-orange-400')} />
      case 'archive':
        return <Archive className={cn(iconSize, 'text-yellow-400')} />
      default:
        return <File className={cn(iconSize, 'text-gray-400')} />
    }
  }

  const handleFileSelect = (fileId: string) => {
    setSelectedFiles(prev => 
      prev.includes(fileId) 
        ? prev.filter(id => id !== fileId)
        : [...prev, fileId]
    )
  }

  const handleDeleteFile = (fileId: string) => {
    console.log('Delete file:', fileId)
  }

  const handleDownloadFile = (file: FileItem) => {
    window.open(file.url, '_blank')
  }

  const handleUpload = () => {
    console.log('Upload new file')
  }

  const fileTypes = [
    { value: 'image', label: 'Images' },
    { value: 'document', label: 'Documents' },
    { value: 'video', label: 'Videos' },
    { value: 'audio', label: 'Audio' },
    { value: 'archive', label: 'Archives' }
  ]

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      <BiHubPageHeader
        title="File Management"
        subtitle="Upload, organize, and manage your media files and documents"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Files' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.FILES_UPLOAD}>
            <Button
              onClick={handleUpload}
              className={BIHUB_ADMIN_THEME.components.button.primary}
            >
              <Upload className="mr-2 h-5 w-5" />
              Upload Files
            </Button>
          </RequirePermission>
        }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <div className="bg-white/5 backdrop-blur-sm border border-blue-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-400/80 text-sm font-medium">Total Files</p>
              <p className="text-2xl font-bold text-blue-100">{stats.totalFiles}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-blue-500/20 flex items-center justify-center">
              <FileImage className="h-5 w-5 text-blue-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-green-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-400/80 text-sm font-medium">Images</p>
              <p className="text-2xl font-bold text-green-100">{stats.images}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-green-500/20 flex items-center justify-center">
              <Image className="h-5 w-5 text-green-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-purple-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-purple-400/80 text-sm font-medium">Documents</p>
              <p className="text-2xl font-bold text-purple-100">{stats.documents}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-purple-500/20 flex items-center justify-center">
              <FileText className="h-5 w-5 text-purple-400" />
            </div>
          </div>
        </div>

        <div className="bg-white/5 backdrop-blur-sm border border-orange-300/20 rounded-xl p-4">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-orange-400/80 text-sm font-medium">Total Size</p>
              <p className="text-2xl font-bold text-orange-100">{stats.totalSize}</p>
            </div>
            <div className="w-10 h-10 rounded-xl bg-orange-500/20 flex items-center justify-center">
              <FolderOpen className="h-5 w-5 text-orange-400" />
            </div>
          </div>
        </div>
      </div>

      {/* File Management */}
      <BiHubAdminCard
        title="File Library"
        subtitle="Browse and manage your uploaded files"
        icon={<FileImage className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex items-center gap-3">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
              <Input
                placeholder="Search files..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-64"
              />
            </div>
            
            <Select value={filterType} onValueChange={setFilterType}>
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Type" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Types</SelectItem>
                {fileTypes.map((type) => (
                  <SelectItem key={type.value} value={type.value}>
                    {type.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>

            <div className="flex items-center border border-gray-600 rounded-lg">
              <Button
                variant={viewMode === 'grid' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setViewMode('grid')}
                className="rounded-r-none"
              >
                <Grid className="h-4 w-4" />
              </Button>
              <Button
                variant={viewMode === 'list' ? 'default' : 'ghost'}
                size="sm"
                onClick={() => setViewMode('list')}
                className="rounded-l-none"
              >
                <List className="h-4 w-4" />
              </Button>
            </div>
          </div>
        }
      >
        {viewMode === 'grid' ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
            {filteredFiles.map((file) => (
              <div
                key={file.id}
                className={cn(
                  'group relative bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl p-4 hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200 cursor-pointer',
                  selectedFiles.includes(file.id) && 'border-[#FF9000] bg-[#FF9000]/10'
                )}
                onClick={() => handleFileSelect(file.id)}
              >
                <div className="aspect-square mb-3 rounded-lg bg-gray-800/50 flex items-center justify-center overflow-hidden">
                  {file.thumbnail ? (
                    <img 
                      src={file.thumbnail} 
                      alt={file.name}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    getFileIcon(file.type, 'lg')
                  )}
                </div>
                
                <div className="space-y-1">
                  <h3 className="text-white text-sm font-medium truncate" title={file.name}>
                    {file.name}
                  </h3>
                  <div className="flex items-center justify-between text-xs text-gray-400">
                    <span>{file.size}</span>
                    <Badge variant="outline" className="text-xs">
                      {file.type}
                    </Badge>
                  </div>
                </div>

                <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                        <MoreHorizontal className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => handleDownloadFile(file)}>
                        <Download className="h-4 w-4 mr-2" />
                        Download
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Eye className="h-4 w-4 mr-2" />
                        Preview
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Copy className="h-4 w-4 mr-2" />
                        Copy URL
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Edit className="h-4 w-4 mr-2" />
                        Rename
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => handleDeleteFile(file.id)}
                        className="text-red-600"
                      >
                        <Trash2 className="h-4 w-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="space-y-2">
            {filteredFiles.map((file) => (
              <div
                key={file.id}
                className={cn(
                  'group flex items-center gap-4 p-4 bg-white/5 backdrop-blur-sm border border-gray-700/50 rounded-xl hover:bg-white/10 hover:border-gray-600/50 transition-all duration-200 cursor-pointer',
                  selectedFiles.includes(file.id) && 'border-[#FF9000] bg-[#FF9000]/10'
                )}
                onClick={() => handleFileSelect(file.id)}
              >
                <div className="w-12 h-12 rounded-lg bg-gray-800/50 flex items-center justify-center flex-shrink-0">
                  {file.thumbnail ? (
                    <img 
                      src={file.thumbnail} 
                      alt={file.name}
                      className="w-full h-full object-cover rounded-lg"
                    />
                  ) : (
                    getFileIcon(file.type)
                  )}
                </div>
                
                <div className="flex-1 min-w-0">
                  <h3 className="text-white font-medium truncate">{file.name}</h3>
                  <div className="flex items-center gap-4 text-sm text-gray-400">
                    <span>{file.size}</span>
                    <span>•</span>
                    <span>{formatDate(file.uploaded_at)}</span>
                    <span>•</span>
                    <span>By {file.uploaded_by}</span>
                    {file.folder && (
                      <>
                        <span>•</span>
                        <span>{file.folder}</span>
                      </>
                    )}
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <Badge variant="outline" className="text-xs">
                    {file.type}
                  </Badge>
                  
                  <div className="opacity-0 group-hover:opacity-100 transition-opacity flex items-center gap-1">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleDownloadFile(file)
                      }}
                    >
                      <Download className="h-4 w-4" />
                    </Button>
                    
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="sm" onClick={(e) => e.stopPropagation()}>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem>
                          <Eye className="h-4 w-4 mr-2" />
                          Preview
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Copy className="h-4 w-4 mr-2" />
                          Copy URL
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Share className="h-4 w-4 mr-2" />
                          Share
                        </DropdownMenuItem>
                        <DropdownMenuItem>
                          <Edit className="h-4 w-4 mr-2" />
                          Rename
                        </DropdownMenuItem>
                        <DropdownMenuItem 
                          onClick={() => handleDeleteFile(file.id)}
                          className="text-red-600"
                        >
                          <Trash2 className="h-4 w-4 mr-2" />
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}

        {filteredFiles.length === 0 && (
          <div className="text-center py-12">
            <FileImage className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-300 mb-2">No files found</h3>
            <p className="text-gray-400 mb-6">
              {searchQuery ? 'Try adjusting your search terms.' : 'Upload your first file to get started.'}
            </p>
            {!searchQuery && (
              <RequirePermission permission={PERMISSIONS.FILES_UPLOAD}>
                <Button
                  onClick={handleUpload}
                  className={BIHUB_ADMIN_THEME.components.button.primary}
                >
                  <Upload className="mr-2 h-5 w-5" />
                  Upload First File
                </Button>
              </RequirePermission>
            )}
          </div>
        )}
      </BiHubAdminCard>
    </div>
  )
}
