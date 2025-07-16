'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { ChevronRight, Home, Package } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { apiClient } from '@/lib/api'
import { Category } from '@/types'
import { cn } from '@/lib/utils'

interface CategoryBreadcrumbsProps {
  categoryId?: string
  showProductCount?: boolean
  className?: string
}

interface BreadcrumbItem {
  id: string
  name: string
  href: string
  productCount?: number
}

export function CategoryBreadcrumbs({ 
  categoryId, 
  showProductCount = false,
  className 
}: CategoryBreadcrumbsProps) {
  const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbItem[]>([])
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    if (!categoryId) {
      setBreadcrumbs([])
      return
    }

    const fetchCategoryPath = async () => {
      try {
        setIsLoading(true)
        
        // Fetch category path from our new API endpoint
        const response = await apiClient.get(`/categories/${categoryId}/path`)
        const categoryPath: Category[] = response.data || []
        
        const breadcrumbItems: BreadcrumbItem[] = categoryPath.map((category) => ({
          id: category.id,
          name: category.name,
          href: `/categories/${category.id}`,
          productCount: showProductCount ? category.product_count : undefined
        }))
        
        setBreadcrumbs(breadcrumbItems)
      } catch (error) {
        console.error('Error fetching category path:', error)
        setBreadcrumbs([])
      } finally {
        setIsLoading(false)
      }
    }

    fetchCategoryPath()
  }, [categoryId, showProductCount])

  if (!categoryId && breadcrumbs.length === 0) {
    return null
  }

  return (
    <nav className={cn("flex items-center space-x-1 text-sm", className)} aria-label="Breadcrumb">
      {/* Home Link */}
      <Link 
        href="/" 
        className="flex items-center text-gray-500 hover:text-primary-600 transition-colors group"
      >
        <Home className="h-4 w-4 group-hover:scale-110 transition-transform duration-200" />
        <span className="sr-only">Home</span>
      </Link>
      
      {breadcrumbs.length > 0 && (
        <>
          <ChevronRight className="h-4 w-4 text-gray-400" />
          
          {/* Categories Link */}
          <Link 
            href="/categories" 
            className="flex items-center text-gray-500 hover:text-primary-600 transition-colors group"
          >
            <Package className="h-4 w-4 mr-1 group-hover:scale-110 transition-transform duration-200" />
            <span>Categories</span>
          </Link>
        </>
      )}

      {/* Category Path */}
      {breadcrumbs.map((item, index) => {
        const isLast = index === breadcrumbs.length - 1
        
        return (
          <div key={item.id} className="flex items-center">
            <ChevronRight className="h-4 w-4 text-gray-400 mx-1" />
            
            {isLast ? (
              <div className="flex items-center">
                <span className="font-medium text-gray-900 flex items-center">
                  {item.name}
                  {item.productCount !== undefined && (
                    <Badge variant="secondary" className="ml-2 text-xs">
                      {item.productCount}
                    </Badge>
                  )}
                </span>
              </div>
            ) : (
              <Link
                href={item.href}
                className="text-gray-500 hover:text-primary-600 transition-colors hover:underline flex items-center"
              >
                {item.name}
                {item.productCount !== undefined && (
                  <Badge variant="outline" className="ml-2 text-xs">
                    {item.productCount}
                  </Badge>
                )}
              </Link>
            )}
          </div>
        )
      })}
      
      {isLoading && (
        <div className="flex items-center ml-2">
          <div className="h-4 w-16 bg-gray-200 rounded animate-pulse"></div>
        </div>
      )}
    </nav>
  )
}

// Structured Data Component for SEO
interface CategoryStructuredDataProps {
  breadcrumbs: BreadcrumbItem[]
}

export function CategoryStructuredData({ breadcrumbs }: CategoryStructuredDataProps) {
  const structuredData = {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    "itemListElement": [
      {
        "@type": "ListItem",
        "position": 1,
        "name": "Home",
        "item": process.env.NEXT_PUBLIC_SITE_URL
      },
      {
        "@type": "ListItem", 
        "position": 2,
        "name": "Categories",
        "item": `${process.env.NEXT_PUBLIC_SITE_URL}/categories`
      },
      ...breadcrumbs.map((item, index) => ({
        "@type": "ListItem",
        "position": index + 3,
        "name": item.name,
        "item": `${process.env.NEXT_PUBLIC_SITE_URL}${item.href}`
      }))
    ]
  }

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
    />
  )
}
