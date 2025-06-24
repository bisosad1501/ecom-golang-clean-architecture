import { toast } from 'sonner'
import api from '@/lib/api'

export interface UploadImageResponse {
  url: string
}

/**
 * Upload an image file to the server and return the URL
 * Authentication và authorization được xử lý ở middleware level
 */
export async function uploadImageFile(file: File, customEndpoint?: string): Promise<string> {
  // Validate file type
  if (!file.type.startsWith('image/')) {
    throw new Error('Please select an image file')
  }

  // Check file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    throw new Error('File size must be less than 5MB')
  }

  try {
    // Auto-detect endpoint based on context if not provided
    let endpoint = customEndpoint
    if (!endpoint) {
      // Check if we're in admin context (simple heuristic)
      const currentPath = typeof window !== 'undefined' ? window.location.pathname : ''
      if (currentPath.includes('/admin')) {
        endpoint = '/admin/upload/image'
      } else {
        endpoint = '/upload/image'
      }
    }

    // Sử dụng endpoint được xác định, authentication xử lý ở middleware
    const response = await api.upload<{ url: string }>(endpoint, file)
    
    // Handle the response structure (same logic as SingleImageUpload)
    let serverUrl: string
    
    // Cast response to any to access url property directly
    const responseAny = response as any
    
    // The backend returns directly {url: '...', message: '...'} not wrapped in data
    if (responseAny.url) {
      serverUrl = responseAny.url
    } else if (response.data && response.data.url) {
      serverUrl = response.data.url
    } else if (response.data && typeof response.data === 'string') {
      serverUrl = response.data
    } else {
      console.error('Unexpected response structure:', response)
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
    
    return fullUrl
    
  } catch (error: any) {
    console.error('Upload failed:', error)
    
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
    
    throw new Error(errorMessage)
  }
}

/**
 * Upload multiple image files and return their URLs
 */
export async function uploadMultipleImageFiles(files: File[], customEndpoint?: string): Promise<string[]> {
  const results: string[] = []
  const errors: string[] = []
  
  for (const file of files) {
    try {
      const url = await uploadImageFile(file, customEndpoint)
      results.push(url)
    } catch (error: any) {
      errors.push(`${file.name}: ${error.message}`)
    }
  }
  
  if (errors.length > 0) {
    toast.error(`Some uploads failed: ${errors.join(', ')}`)
  }
  
  return results
}

/**
 * Upload an image file for public use (no authentication required)
 */
export async function uploadPublicImageFile(file: File): Promise<string> {
  return uploadImageFile(file, '/public/upload/image')
}

/**
 * Upload an image file for admin use (admin authentication required)
 */
export async function uploadAdminImageFile(file: File): Promise<string> {
  return uploadImageFile(file, '/admin/upload/image')
}
