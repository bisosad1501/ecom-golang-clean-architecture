'use client'

import React, { useState, useCallback } from 'react'
import { FormField } from '@/components/ui/form-field'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Upload, Link, X, GripVertical, Image as ImageIcon } from 'lucide-react'
import { DragDropContext, Droppable, Draggable, type DropResult } from '@hello-pangea/dnd'
import Image from 'next/image'
import { cn } from '@/lib/utils'
import { uploadMultipleImageFiles } from '@/lib/utils/image-upload'

interface ImageItem {
  url: string
  alt_text?: string
  position: number
}

interface MultiImageUploadProps {
  images: ImageItem[]
  onImagesChange: (images: ImageItem[]) => void
  maxImages?: number
  error?: string
  className?: string
  endpoint?: string
  uploadLabel?: string
}

export function MultiImageUpload({
  images,
  onImagesChange,
  maxImages = 10,
  error,
  className,
  endpoint = 'admin',
  uploadLabel = 'Upload Images',
}: MultiImageUploadProps) {
  const [isUploading, setIsUploading] = useState(false)
  const [urlInput, setUrlInput] = useState('')
  const [draggedOver, setDraggedOver] = useState(false)

  const handleFileUpload = async (files: FileList) => {
    if (files.length === 0 || images.length >= maxImages) return

    const fileArray = Array.from(files).slice(0, maxImages - images.length)
    
    try {
      setIsUploading(true)
      const uploadedUrls = await uploadMultipleImageFiles(fileArray, endpoint)
      
      const newImages: ImageItem[] = uploadedUrls.map((url, index) => ({
        url,
        alt_text: '',
        position: images.length + index,
      }))

      onImagesChange([...images, ...newImages])
    } catch (error) {
      console.error('Failed to upload images:', error)
    } finally {
      setIsUploading(false)
    }
  }

  const handleUrlAdd = () => {
    const trimmedUrl = urlInput.trim()
    if (trimmedUrl && images.length < maxImages) {
      const newImage: ImageItem = {
        url: trimmedUrl,
        alt_text: '',
        position: images.length,
      }
      onImagesChange([...images, newImage])
      setUrlInput('')
    }
  }

  const removeImage = (index: number) => {
    const newImages = images
      .filter((_, i) => i !== index)
      .map((img, i) => ({ ...img, position: i }))
    onImagesChange(newImages)
  }

  const handleDragEnd = (result: DropResult) => {
    if (!result.destination) return

    const items = Array.from(images)
    const [reorderedItem] = items.splice(result.source.index, 1)
    items.splice(result.destination.index, 0, reorderedItem)

    // Update positions
    const updatedItems = items.map((item, index) => ({
      ...item,
      position: index,
    }))

    onImagesChange(updatedItems)
  }

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDraggedOver(true)
  }, [])

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDraggedOver(false)
  }, [])

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    setDraggedOver(false)
    const files = e.dataTransfer.files
    if (files) {
      handleFileUpload(files)
    }
  }, [images.length, maxImages, endpoint])

  return (
    <FormField
      label="Product Images"
      error={error}
      hint={`Upload up to ${maxImages} images. Drag to reorder.`}
      className={className}
    >
      <div className="space-y-4">
        {/* Upload Area */}
        {images.length < maxImages && (
          <div className="space-y-3">
            {/* File Upload */}
            <div
              className={cn(
                "border-2 border-dashed rounded-lg p-6 text-center transition-colors",
                draggedOver ? "border-blue-500 bg-blue-50" : "border-gray-300",
                "hover:border-gray-400"
              )}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              <div className="flex flex-col items-center justify-center space-y-2">
                <Upload className="h-8 w-8 text-gray-400" />
                <div>
                  <label className="cursor-pointer">
                    <input
                      type="file"
                      multiple
                      accept="image/*"
                      onChange={(e) => e.target.files && handleFileUpload(e.target.files)}
                      className="hidden"
                      disabled={isUploading}
                    />
                    <Button
                      type="button"
                      variant="outline"
                      disabled={isUploading}
                      asChild
                    >
                      <span>
                        {isUploading ? 'Uploading...' : uploadLabel}
                      </span>
                    </Button>
                  </label>
                </div>
                <p className="text-sm text-gray-500">
                  or drag and drop images here
                </p>
              </div>
            </div>

            {/* URL Input */}
            <div className="flex gap-2">
              <Input
                placeholder="Or paste image URL"
                value={urlInput}
                onChange={(e) => setUrlInput(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleUrlAdd()}
              />
              <Button
                type="button"
                onClick={handleUrlAdd}
                disabled={!urlInput.trim()}
                size="sm"
                variant="outline"
              >
                <Link className="h-4 w-4" />
              </Button>
            </div>
          </div>
        )}

        {/* Images List */}
        {images.length > 0 && (
          <DragDropContext onDragEnd={handleDragEnd}>
            <Droppable droppableId="images">
              {(provided) => (
                <div
                  {...provided.droppableProps}
                  ref={provided.innerRef}
                  className="space-y-2"
                >
                  {images.map((image, index) => (
                    <Draggable
                      key={`${image.url}-${index}`}
                      draggableId={`${image.url}-${index}`}
                      index={index}
                    >
                      {(provided, snapshot) => (
                        <div
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          className={cn(
                            "flex items-center gap-3 p-3 border rounded-lg bg-white",
                            snapshot.isDragging && "shadow-lg"
                          )}
                        >
                          <div
                            {...provided.dragHandleProps}
                            className="text-gray-400 hover:text-gray-600 cursor-grab"
                          >
                            <GripVertical className="h-4 w-4" />
                          </div>

                          {/* Image Preview */}
                          <div className="relative w-16 h-16 rounded border overflow-hidden bg-gray-100">
                            {image.url ? (
                              <Image
                                src={image.url}
                                alt={image.alt_text || `Image ${index + 1}`}
                                fill
                                className="object-cover"
                                onError={(e) => {
                                  const target = e.target as HTMLImageElement
                                  target.style.display = 'none'
                                  target.nextElementSibling?.classList.remove('hidden')
                                }}
                              />
                            ) : null}
                            <div className="absolute inset-0 flex items-center justify-center text-gray-400">
                              <ImageIcon className="h-6 w-6" />
                            </div>
                          </div>

                          {/* Image Info */}
                          <div className="flex-1 min-w-0">
                            <p className="text-sm font-medium truncate">
                              Image {index + 1}
                            </p>
                            <p className="text-xs text-gray-500 truncate">
                              {image.url}
                            </p>
                            {index === 0 && (
                              <Badge variant="secondary" className="text-xs mt-1">
                                Primary
                              </Badge>
                            )}
                          </div>

                          {/* Remove Button */}
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => removeImage(index)}
                            className="text-red-600 hover:text-red-700 hover:bg-red-50"
                          >
                            <X className="h-4 w-4" />
                          </Button>
                        </div>
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </div>
              )}
            </Droppable>
          </DragDropContext>
        )}

        {/* Status Info */}
        <div className="flex justify-between text-sm text-gray-500">
          <span>{images.length} of {maxImages} images</span>
          {images.length >= maxImages && (
            <span className="text-amber-600">Maximum images reached</span>
          )}
        </div>
      </div>
    </FormField>
  )
}
