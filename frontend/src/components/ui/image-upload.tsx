'use client'

import { useState, useRef } from 'react'
import { Upload, Link2, X, Image as ImageIcon } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'

interface ImageUploadProps {
  value?: string
  onChange: (url: string) => void
  onRemove: () => void
  className?: string
  label?: string
  placeholder?: string
}

export function ImageUpload({
  value,
  onChange,
  onRemove,
  className,
  label = "Image",
  placeholder = "Enter image URL or upload file"
}: ImageUploadProps) {
  const [isDragOver, setIsDragOver] = useState(false)
  const [urlInput, setUrlInput] = useState('')
  const [showUrlInput, setShowUrlInput] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleFileUpload = (file: File) => {
    if (!file.type.startsWith('image/')) {
      alert('Please select an image file')
      return
    }

    // Check file size (max 5MB)
    if (file.size > 5 * 1024 * 1024) {
      alert('File size must be less than 5MB')
      return
    }

    // Create object URL for preview (in real app, upload to server)
    const objectUrl = URL.createObjectURL(file)
    onChange(objectUrl)
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
    
    const file = e.dataTransfer.files[0]
    if (file) {
      handleFileUpload(file)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    setIsDragOver(false)
  }

  const handleUrlSubmit = () => {
    if (urlInput.trim()) {
      onChange(urlInput.trim())
      setUrlInput('')
      setShowUrlInput(false)
    }
  }

  const handleRemove = () => {
    if (value && value.startsWith('blob:')) {
      URL.revokeObjectURL(value)
    }
    onRemove()
  }

  return (
    <div className={cn("space-y-4", className)}>
      <Label>{label}</Label>
      
      {value ? (
        // Preview existing image
        <Card>
          <CardContent className="p-4">
            <div className="relative">
              <img
                src={value}
                alt="Preview"
                className="w-full h-40 object-cover rounded-lg"
                onError={(e) => {
                  e.currentTarget.src = '/placeholder-product.svg'
                }}
              />
              <Button
                type="button"
                variant="destructive"
                size="sm"
                className="absolute top-2 right-2"
                onClick={handleRemove}
              >
                <X className="h-4 w-4" />
              </Button>
            </div>
            <p className="text-xs text-gray-500 mt-2 break-all">{value}</p>
          </CardContent>
        </Card>
      ) : (
        // Upload area
        <Card>
          <CardContent className="p-6">
            <div
              className={cn(
                "border-2 border-dashed rounded-lg p-8 text-center transition-colors",
                isDragOver 
                  ? "border-blue-500 bg-blue-50" 
                  : "border-gray-300 hover:border-gray-400"
              )}
              onDrop={handleDrop}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
            >
              <div className="flex flex-col items-center space-y-4">
                <div className="p-4 bg-gray-100 rounded-full">
                  <ImageIcon className="h-8 w-8 text-gray-400" />
                </div>
                
                <div>
                  <h3 className="text-lg font-medium text-gray-900">
                    {isDragOver ? 'Drop image here!' : 'Upload Image'}
                  </h3>
                  <p className="text-sm text-gray-500 mt-1">
                    Drag and drop an image file, or click to browse
                  </p>
                  <p className="text-xs text-gray-400 mt-1">
                    Supports: JPG, PNG, GIF (max 5MB)
                  </p>
                </div>

                <div className="flex items-center space-x-4">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => fileInputRef.current?.click()}
                  >
                    <Upload className="mr-2 h-4 w-4" />
                    Choose File
                  </Button>
                  
                  <span className="text-gray-400">or</span>
                  
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setShowUrlInput(!showUrlInput)}
                  >
                    <Link2 className="mr-2 h-4 w-4" />
                    Enter URL
                  </Button>
                </div>

                <input
                  ref={fileInputRef}
                  type="file"
                  accept="image/*"
                  onChange={handleFileChange}
                  className="hidden"
                />
              </div>
            </div>

            {/* URL Input */}
            {showUrlInput && (
              <div className="mt-4 space-y-3">
                <div className="flex space-x-2">
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
                  />
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
