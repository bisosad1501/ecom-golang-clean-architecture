'use client'

import { useState, useRef, useCallback, useEffect } from 'react'
import { DragDropContext, Droppable, Draggable, type DropResult } from '@hello-pangea/dnd'
import { createPortal } from 'react-dom'
import Image from 'next/image'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { 
  X, 
  Upload, 
  Link, 
  Star, 
  Image as ImageIcon, 
  Edit2,
  Plus,
  Loader2
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { toast } from 'sonner'

export interface ImageUploadItem {
  url: string
  alt_text?: string
  position: number
}

interface ImageUploadGridProps {
  images: ImageUploadItem[]
  onImagesChange: (images: ImageUploadItem[]) => void
  onUploadFiles: (files: File[]) => Promise<string[]>
  maxImages?: number
  isUploading?: boolean
  className?: string
}

export function ImageUploadGrid({
  images,
  onImagesChange,
  onUploadFiles,
  maxImages = 10,
  isUploading = false,
  className
}: ImageUploadGridProps) {
  const [isDragOver, setIsDragOver] = useState(false)
  const [imageUrl, setImageUrl] = useState('')
  const [isDragging, setIsDragging] = useState(false)
  const [portalContainer, setPortalContainer] = useState<HTMLElement | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    // Create or get portal container for drag previews
    let container = document.getElementById('drag-portal')
    if (!container) {
      container = document.createElement('div')
      container.id = 'drag-portal'
      container.style.position = 'fixed'
      container.style.top = '0'
      container.style.left = '0'
      container.style.zIndex = '9999'
      container.style.pointerEvents = 'none'
      document.body.appendChild(container)
    }
    setPortalContainer(container)

    return () => {
      // Clean up on unmount
      const existingContainer = document.getElementById('drag-portal')
      if (existingContainer && existingContainer.children.length === 0) {
        document.body.removeChild(existingContainer)
      }
    }
  }, [])

  const handleFileSelect = async (files: File[]) => {
    if (images.length + files.length > maxImages) {
      toast.error(`Maximum ${maxImages} images allowed`)
      return
    }

    try {
      const uploadedUrls = await onUploadFiles(files)
      
      // Add successfully uploaded images
      const newImages = uploadedUrls.map((url, index) => ({
        url,
        alt_text: files[index]?.name?.split('.')[0] || '',
        position: images.length + index,
      }))

      onImagesChange([...images, ...newImages])
      
      if (uploadedUrls.length > 0) {
        toast.success(`Successfully uploaded ${uploadedUrls.length} image(s)`)
      }
    } catch (error: any) {
      toast.error(error.message || 'Failed to upload images')
    }
  }

  const handleFileInputChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (files && files.length > 0) {
      handleFileSelect(Array.from(files))
    }
    // Clear the input
    event.target.value = ''
  }

  const handleDrop = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    setIsDragOver(false)

    const files = Array.from(event.dataTransfer.files).filter(file => 
      file.type.startsWith('image/')
    )

    if (files.length === 0) {
      toast.error('Please drop image files only')
      return
    }

    handleFileSelect(files)
  }, [images.length, maxImages, onUploadFiles, onImagesChange])

  const handleDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    setIsDragOver(true)
  }, [])

  const handleDragLeave = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    // Only set isDragOver to false if we're actually leaving the drop zone
    if (!event.currentTarget.contains(event.relatedTarget as Node)) {
      setIsDragOver(false)
    }
  }, [])

  const removeImage = (slotIndex: number) => {
    const newImages = [...images]
    newImages.splice(slotIndex, 1)
    onImagesChange(newImages)
  }

  const updateImageAltText = (slotIndex: number, newAltText: string) => {
    const newImages = [...images]
    if (newImages[slotIndex]) {
      newImages[slotIndex] = { ...newImages[slotIndex], alt_text: newAltText }
      onImagesChange(newImages)
    }
  }

  const handleDragEnd = (result: DropResult) => {
    setIsDragging(false)
    
    const { destination, source } = result
    
    if (!destination || destination.index === source.index) {
      return
    }
    
    const newImages = [...images]
    const [movedImage] = newImages.splice(source.index, 1)
    newImages.splice(destination.index, 0, movedImage)
    
    // Update positions to match array indices
    const reorderedImages = newImages.map((img, i) => ({ ...img, position: i }))
    onImagesChange(reorderedImages)
  }

  const handleDragStart = () => {
    setIsDragging(true)
  }

  const addImageFromUrl = () => {
    if (!imageUrl.trim()) return
    
    try {
      new URL(imageUrl) // Validate URL
      const newImage: ImageUploadItem = {
        url: imageUrl,
        alt_text: '',
        position: images.length,
      }
      onImagesChange([...images, newImage])
      setImageUrl('')
    } catch {
      toast.error('Please enter a valid image URL')
    }
  }

  return (
    <div 
      className={cn("space-y-6", className)}
      style={{ 
        overflow: 'visible',
        // Ensure container doesn't interfere with drag positioning
        position: 'relative',
        isolation: 'isolate'
      }}
    >
      {/* Upload Section */}
      <div className="space-y-4">
        {/* File Upload Drop Zone */}
        <div
          className={cn(
            "relative border-2 border-dashed rounded-lg p-6 text-center transition-all duration-200",
            isDragOver 
              ? "border-blue-500 bg-blue-50 dark:bg-blue-950/20" 
              : "border-gray-300 hover:border-gray-400",
            isUploading && "pointer-events-none opacity-60"
          )}
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
        >
          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept="image/*"
            onChange={handleFileInputChange}
            className="hidden"
            disabled={isUploading}
          />
          
          {isUploading ? (
            <div className="flex flex-col items-center gap-3">
              <Loader2 className="h-8 w-8 animate-spin text-blue-500" />
              <div>
                <p className="font-medium">Uploading images...</p>
                <p className="text-sm text-gray-500">Please wait</p>
              </div>
            </div>
          ) : (
            <div className="flex flex-col items-center gap-3">
              <div className={cn(
                "rounded-full p-3 transition-colors",
                isDragOver 
                  ? "bg-blue-100 dark:bg-blue-900/30" 
                  : "bg-gray-100 dark:bg-gray-800"
              )}>
                <Upload className={cn(
                  "h-6 w-6",
                  isDragOver ? "text-blue-600" : "text-gray-600"
                )} />
              </div>
              
              <div>
                <p className="font-medium">
                  {isDragOver ? "Drop images here" : "Drag and drop images"}
                </p>
                <p className="text-sm text-gray-500">
                  or{" "}
                  <button
                    type="button"
                    onClick={() => fileInputRef.current?.click()}
                    className="text-blue-600 hover:text-blue-700 font-medium underline"
                  >
                    browse files
                  </button>
                </p>
              </div>
              
              <div className="flex gap-2 text-xs text-gray-500">
                <span>PNG, JPG up to 10MB</span>
                <span>•</span>
                <span>Max {maxImages} images</span>
              </div>
            </div>
          )}
        </div>

        {/* URL Upload */}
        <div className="flex gap-2">
          <div className="flex-1">
            <Label htmlFor="image-url" className="sr-only">Image URL</Label>
            <Input
              id="image-url"
              type="url"
              placeholder="Or paste an image URL..."
              value={imageUrl}
              onChange={(e) => setImageUrl(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && addImageFromUrl()}
              disabled={isUploading}
            />
          </div>
          <Button
            type="button"
            variant="outline"
            onClick={addImageFromUrl}
            disabled={!imageUrl.trim() || isUploading}
            className="shrink-0"
          >
            <Link className="h-4 w-4 mr-2" />
            Add URL
          </Button>
        </div>
      </div>

      {/* Images Grid - Always show slots */}
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <h4 className="font-medium">
            Product Images ({images.length}/{maxImages})
          </h4>
          {images.length > 1 && (
            <p className="text-sm text-gray-500">Drag images to reorder slots</p>
          )}
        </div>

        {/* Container with proper overflow handling for drag & drop */}
        <div className="relative" style={{ overflow: 'visible' }}>
          <DragDropContext 
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
          >
            <Droppable droppableId="images-grid" direction="horizontal">
              {(provided, snapshot) => (
                <div
                  ref={provided.innerRef}
                  {...provided.droppableProps}
                  className={cn(
                    "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 p-3 min-h-[120px] rounded-lg transition-all duration-200 relative",
                    snapshot.isDraggingOver && "bg-blue-50 border-2 border-dashed border-blue-400"
                  )}
                  style={{
                    overflow: 'visible',
                    position: 'relative'
                  }}
                >
                  {/* Render draggable images */}
                  {images.map((image, imageIndex) => (
                    <Draggable
                      key={`image-${image.url}-${imageIndex}`}
                      draggableId={`image-${image.url}-${imageIndex}`}
                      index={imageIndex}
                    >
                      {(provided, snapshot) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          {...provided.dragHandleProps}
                          className={cn(
                            "group relative aspect-square rounded-lg overflow-hidden border-2 bg-white cursor-grab active:cursor-grabbing transition-shadow duration-200",
                            snapshot.isDragging 
                              ? "border-blue-500 shadow-xl ring-2 ring-blue-200 z-50" 
                              : "border-gray-200 hover:border-blue-300 hover:shadow-md"
                          )}
                          style={{
                            ...provided.draggableProps.style,
                            // Simple, clean positioning for drag
                            zIndex: snapshot.isDragging ? 9999 : 'auto',
                          }}
                        >
                          {/* Slot Number */}
                          <div className={cn(
                            "absolute top-2 left-2 z-20 text-xs px-1.5 py-0.5 rounded font-medium transition-colors",
                            snapshot.isDragging 
                              ? "bg-blue-600 text-white" 
                              : "bg-white/90 text-gray-600"
                          )}>
                            #{imageIndex + 1}
                            {imageIndex === 0 && (
                              <Star className="inline h-3 w-3 ml-1 text-current" />
                            )}
                          </div>

                          {/* Dragging Indicator */}
                          {snapshot.isDragging && (
                            <div className="absolute inset-0 bg-blue-500/20 flex items-center justify-center">
                              <div className="bg-blue-600 text-white px-3 py-1 rounded-full text-sm font-medium shadow-lg">
                                Moving...
                              </div>
                            </div>
                          )}

                          {/* Image */}
                          <Image
                            src={image.url}
                            alt={image.alt_text || `Product image ${imageIndex + 1}`}
                            fill
                            className="object-cover select-none"
                            unoptimized={true}
                            onError={(e) => {
                              const target = e.target as HTMLImageElement;
                              target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgZmlsbD0iI0Y3RjhGOSIvPjx0ZXh0IHg9IjUwJSIgeT0iNTAlIiBkb21pbmFudC1iYXNlbGluZT0ibWlkZGxlIiB0ZXh0LWFuY2hvcj0ibWlkZGxlIiBmb250LWZhbWlseT0ic3lzdGVtLXVpIiBmb250LXNpemU9IjE2IiBmaWxsPSIjOUNBM0FGIj5FcnJvcjwvdGV4dD48L3N2Zz4=';
                              target.onerror = null;
                            }}
                          />

                          {/* Overlay Controls - Hide when dragging */}
                          {!snapshot.isDragging && (
                            <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors flex items-center justify-center">
                              <Button
                                type="button"
                                size="sm"
                                variant="secondary"
                                className="opacity-0 group-hover:opacity-100 transition-opacity h-8 w-8 p-0 bg-white/90 hover:bg-white"
                                onClick={(e) => {
                                  e.stopPropagation()
                                  const newAltText = prompt('Enter alt text for this image:', image.alt_text || '')
                                  if (newAltText !== null) {
                                    updateImageAltText(imageIndex, newAltText)
                                  }
                                }}
                                title="Edit alt text"
                              >
                                <Edit2 className="h-3 w-3" />
                              </Button>
                            </div>
                          )}

                          {/* Remove Button - Hide when dragging */}
                          {!snapshot.isDragging && (
                            <Button
                              type="button"
                              size="sm"
                              variant="destructive"
                              className="absolute top-1 right-1 z-30 opacity-0 group-hover:opacity-100 transition-opacity h-6 w-6 p-0 rounded-full bg-red-500 hover:bg-red-600 text-white"
                              onClick={(e) => {
                                e.stopPropagation()
                                removeImage(imageIndex)
                              }}
                              title="Remove image"
                            >
                              <X className="h-3 w-3" />
                            </Button>
                          )}

                          {/* Alt Text Display */}
                          {image.alt_text && (
                            <div className="absolute bottom-1 left-1 right-1 bg-black/70 text-white text-xs px-2 py-1 rounded truncate">
                              {image.alt_text}
                            </div>
                          )}
                        </div>
                      )}
                    </Draggable>
                  ))}

                  {/* Render empty slots */}
                  {Array.from({ length: maxImages - images.length }, (_, emptyIndex) => (
                    <div
                      key={`empty-${emptyIndex}`}
                      className={cn(
                        "relative aspect-square rounded-lg border-2 border-dashed bg-gray-50/50 transition-all duration-200",
                        snapshot.isDraggingOver 
                          ? "border-blue-400 bg-blue-50 scale-105" 
                          : "border-gray-300 hover:border-gray-400"
                      )}
                    >
                      {/* Slot Number */}
                      <div className="absolute top-2 left-2 z-10 bg-white/90 text-gray-400 text-xs px-1.5 py-0.5 rounded font-medium">
                        #{images.length + emptyIndex + 1}
                      </div>

                      {/* Empty Indicator */}
                      <div className="absolute inset-0 flex flex-col items-center justify-center text-gray-400">
                        <div className="text-center">
                          <div className="w-8 h-8 rounded-full border-2 border-dashed border-gray-300 flex items-center justify-center mb-2">
                            <Plus className="h-4 w-4" />
                          </div>
                          <div className="text-xs">Empty</div>
                        </div>
                      </div>
                    </div>
                  ))}
                  {provided.placeholder}
                </div>
              )}
            </Droppable>
          </DragDropContext>
        </div>
      </div>

      {/* Empty State */}
      {images.length === 0 && !isUploading && (
        <div className="text-center py-6 text-gray-500">
          <p className="text-sm">Upload images to fill the slots above</p>
        </div>
      )}
    </div>
  )
}
