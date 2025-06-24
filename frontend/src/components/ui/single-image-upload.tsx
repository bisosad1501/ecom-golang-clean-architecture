'use client'

import { useState, useRef } from 'react'
import { Upload, Link2, X, ImageIcon, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'
import api from '@/lib/api'

interface SingleImageUploadProps {
  value?: string
  onChange: (url: string) => void
  onRemove: () => void
  className?: string
  label?: string
  placeholder?: string
  disabled?: boolean
}

export function SingleImageUpload({
  value,
  onChange,
  onRemove,
  className,
  label = "Image",
  placeholder = "Enter image URL",
  disabled = false
}: SingleImageUploadProps) {
  const [isDragOver, setIsDragOver] = useState(false)
  const [urlInput, setUrlInput] = useState('')
  const [showUrlInput, setShowUrlInput] = useState(false)
  const [isUploading, setIsUploading] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileUpload = async (file: File) => {
    if (!file.type.startsWith('image/')) {
      toast.error('Please select an image file')
      return
    }

    // Check file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      toast.error('File size must be less than 5MB')
      return
    }

    setIsUploading(true)
    try {
      // Upload file to server
      const response = await api.upload<{ url: string }>('/admin/upload/image', file)
      
      console.log('Upload response:', response)
      console.log('Response.data:', response.data)
      console.log('Type of response.data:', typeof response.data)
      
      // Check if response exists
      if (!response) {
        throw new Error('No response from server')
      }
      
      // Handle the response structure
      let serverUrl: string
      
      // Cast response to any to access url property directly
      const responseAny = response as any
      
      // The backend returns directly {url: '...', message: '...'} not wrapped in data
      if (responseAny.url) {
        // Direct response from backend
        serverUrl = responseAny.url
        console.log('Found URL directly in response:', serverUrl)
      } else if (response.data && response.data.url) {
        // Wrapped in data property
        serverUrl = response.data.url
        console.log('Found URL in response.data.url:', serverUrl)
      } else if (response.data && typeof response.data === 'string') {
        // Sometimes the response might be a direct string
        serverUrl = response.data
        console.log('Found URL as direct string:', serverUrl)
      } else {
        console.error('Unexpected response structure:', response)
        console.error('Response data:', response.data)
        throw new Error('Invalid response structure from server')
      }
      
      if (!serverUrl) {
        throw new Error('No URL returned from server')
      }
      
      // Convert relative URL to absolute URL for frontend display
      const baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'
      
      // Remove /api/v1 from baseUrl for static files
      const staticBaseUrl = baseUrl.replace('/api/v1', '')
      const fullUrl = serverUrl.startsWith('http') ? serverUrl : `${staticBaseUrl}${serverUrl}`
      
      onChange(fullUrl)
      toast.success('Image uploaded successfully!')
      
    } catch (error: any) {
      console.error('Upload failed:', error)
      console.error('Error details:', {
        message: error?.message,
        response: error?.response,
        status: error?.response?.status,
        data: error?.response?.data
      })
      
      let errorMessage = 'Failed to upload image'
      
      if (error?.response?.status === 401) {
        errorMessage = 'Unauthorized - please log in again'
      } else if (error?.response?.status === 413) {
        errorMessage = 'File too large - maximum size is 5MB'
      } else if (error?.response?.status === 415) {
        errorMessage = 'Unsupported file type - please use JPG, PNG, or GIF'
      } else if (error?.message) {
        errorMessage = error.message
      } else if (error?.response?.data?.message) {
        errorMessage = error.response.data.message
      }
      
      toast.error(errorMessage)
    } finally {
      setIsUploading(false)
    }
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      handleFileUpload(file)
    }
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
    
    if (disabled) return
    
    const file = e.dataTransfer.files[0]
    if (file) {
      handleFileUpload(file)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    if (!disabled) {
      setIsDragOver(true)
    }
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }

  const handleUrlSubmit = () => {
    const url = urlInput.trim()
    if (!url) return
    
    // Basic URL validation
    try {
      new URL(url)
      onChange(url)
      setUrlInput('')
      setShowUrlInput(false)
    } catch {
      toast.error('Please enter a valid URL')
    }
  }

  const handleRemove = () => {
    if (value && value.startsWith('blob:')) {
      URL.revokeObjectURL(value)
    }
    onRemove()
  }

  const handleChooseFile = () => {
    if (!disabled) {
      fileInputRef.current?.click()
    }
  }

  const handleShowUrlInput = () => {
    if (!disabled) {
      setShowUrlInput(!showUrlInput)
    }
  }

  return (
    <div className={cn("space-y-4", className)}>
      {label && <Label>{label}</Label>}
      
      {value ? (
        // Preview existing image
        <Card className="group">
          <CardContent className="p-4">
            <div className="relative">
              <div className="aspect-video relative overflow-hidden rounded-lg bg-gray-100">
                <img
                  src={value}
                  alt="Preview"
                  className="w-full h-full object-cover transition-transform hover:scale-105"
                  onError={(e) => {
                    const target = e.target as HTMLImageElement
                    target.src = '/placeholder-product.svg'
                  }}
                />
              </div>
              
              {!disabled && (
                <Button
                  type="button"
                  variant="destructive"
                  size="sm"
                  className="absolute top-2 right-2 opacity-80 hover:opacity-100 transition-opacity"
                  onClick={handleRemove}
                >
                  <X className="h-4 w-4" />
                </Button>
              )}
            </div>
            
            <div className="mt-3 space-y-2">
              <p className="text-xs text-gray-500 break-all">{value}</p>
              
              {!disabled && (
                <div className="flex gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleChooseFile}
                    disabled={isUploading}
                  >
                    <Upload className="mr-1 h-3 w-3" />
                    Replace
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={handleShowUrlInput}
                  >
                    <Link2 className="mr-1 h-3 w-3" />
                    Change URL
                  </Button>
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      ) : (
        // Upload area
        <Card>
          <CardContent className="p-6">
            <div
              className={cn(
                "border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200",
                disabled 
                  ? "border-gray-200 bg-gray-50 cursor-not-allowed"
                  : isDragOver 
                    ? "border-blue-500 bg-blue-50 scale-105" 
                    : "border-gray-300 hover:border-gray-400 hover:bg-gray-50 cursor-pointer"
              )}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onClick={handleChooseFile}
            >
              <div className="flex flex-col items-center space-y-4">
                <div className={cn(
                  "p-4 rounded-full transition-colors",
                  disabled 
                    ? "bg-gray-200"
                    : isDragOver 
                      ? "bg-blue-100" 
                      : "bg-gray-100"
                )}>
                  <ImageIcon className={cn(
                    "h-8 w-8 transition-colors",
                    disabled 
                      ? "text-gray-400"
                      : isDragOver 
                        ? "text-blue-500" 
                        : "text-gray-400"
                  )} />
                </div>
                
                <div>
                  <h3 className={cn(
                    "text-lg font-medium transition-colors",
                    disabled 
                      ? "text-gray-400"
                      : isDragOver 
                        ? "text-blue-600" 
                        : "text-gray-900"
                  )}>
                    {isUploading 
                      ? 'Uploading...' 
                      : isDragOver 
                        ? 'Drop image here!' 
                        : 'Upload Category Image'
                    }
                  </h3>
                  <p className="text-sm text-gray-500 mt-1">
                    {disabled 
                      ? 'Image upload disabled'
                      : 'Drag and drop an image file, or click to browse'
                    }
                  </p>
                  {!disabled && (
                    <p className="text-xs text-gray-400 mt-1">
                      Supports: JPG, PNG, GIF (max 5MB)
                    </p>
                  )}
                </div>

                {!disabled && (
                  <div className="flex items-center space-x-4">
                    <Button
                      type="button"
                      variant="outline"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleChooseFile()
                      }}
                      disabled={isUploading}
                    >
                      <Upload className="mr-2 h-4 w-4" />
                      {isUploading ? 'Uploading...' : 'Choose File'}
                    </Button>
                    
                    <span className="text-gray-400">or</span>
                    
                    <Button
                      type="button"
                      variant="outline"
                      onClick={(e) => {
                        e.stopPropagation()
                        handleShowUrlInput()
                      }}
                      disabled={isUploading}
                    >
                      <Link2 className="mr-2 h-4 w-4" />
                      Enter URL
                    </Button>
                  </div>
                )}

                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                  className="hidden"
                  disabled={disabled}
                />
              </div>
            </div>

            {/* URL Input */}
            {showUrlInput && !disabled && (
              <div className="mt-4 space-y-3">
                <div className="flex space-x-2">
                  <div className="flex-1 relative">
                    <Link2 className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
                    <Input
                      placeholder={placeholder}
                      value={urlInput}
                      onChange={(e) => setUrlInput(e.target.value)}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter') {
                          e.preventDefault()
                          handleUrlSubmit()
                        }
                      }}
                      className="pl-10"
                    />
                  </div>
                  <Button
                    type="button"
                    onClick={handleUrlSubmit}
                    disabled={!urlInput.trim()}
                  >
                    Add
                  </Button>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => {
                    setShowUrlInput(false)
                    setUrlInput('')
                  }}
                >
                  Cancel
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
