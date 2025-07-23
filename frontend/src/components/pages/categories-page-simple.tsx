'use client'

import React, { useState, useMemo, useEffect } from 'react'
import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Search, Grid, List, FolderTree, ChevronRight, Home, X, ArrowRight, Eye } from 'lucide-react'
import Image from 'next/image'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { Category } from '@/types'
import { categoryService } from '@/lib/services/categories'
import { useQuery } from '@tanstack/react-query'

export function CategoriesPageSimple() {
  const searchParams = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get('search') || '')
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  const [sortBy, setSortBy] = useState('name')
  const [selectedParent, setSelectedParent] = useState('')
  const [showActiveOnly, setShowActiveOnly] = useState(true)
  const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc')
  const [currentPage, setCurrentPage] = useState(1)
  const [itemsPerPage, setItemsPerPage] = useState(12)

  // Fetch category tree
  const {
    data: categoryTree,
    isLoading,
    error
  } = useQuery({
    queryKey: ['categories', 'tree'],
    queryFn: () => categoryService.getCategoryTree(),
  })

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    const url = new URL(window.location.href)
    if (query) {
      url.searchParams.set('search', query)
    } else {
      url.searchParams.delete('search')
    }
    window.history.replaceState({}, '', url.toString())
  }

  if (error) {
    return (
      <div className="min-h-screen bg-black">
        <div className="container mx-auto px-4 py-8">
          <Card className="p-8 text-center bg-gray-900 border-gray-700">
            <h2 className="text-2xl font-bold text-white mb-4">
              Oops! Something went wrong
            </h2>
            <p className="text-gray-300 mb-6">
              We couldn't load the categories. Please try again later.
            </p>
            <Button onClick={() => window.location.reload()} variant="outline" className="bg-[#ff9000] text-white border-[#ff9000] hover:bg-[#ff9000]/90">
              Try Again
            </Button>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white relative overflow-hidden">
      <AnimatedBackground className="opacity-30" />
      
      <div className="container mx-auto px-4 lg:px-6 xl:px-8 py-4 lg:py-8 relative z-10">
        <div className="flex flex-col lg:flex-row gap-6 lg:gap-8">
          {/* Sidebar */}
          <div className="lg:w-80 flex-shrink-0">
            <div className="lg:sticky lg:top-8">
              <h1 className="text-3xl font-bold text-white mb-8">Categories</h1>
              
              {/* Search */}
              <div className="relative mb-6">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
                <Input
                  type="text"
                  placeholder="Search categories..."
                  value={searchQuery}
                  onChange={(e) => handleSearch(e.target.value)}
                  className="pl-10 bg-gray-900/50 border-gray-700 text-white placeholder-gray-400 focus:border-[#ff9000]"
                />
              </div>
              
              {/* Stats */}
              <div className="p-4 bg-gray-900/30 rounded-lg">
                <h3 className="text-lg font-semibold text-white mb-4">Statistics</h3>
                <div className="text-sm text-gray-300">
                  <div>Total: {categoryTree?.length || 0}</div>
                  <div>Loading: {isLoading ? 'Yes' : 'No'}</div>
                </div>
              </div>
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1 min-w-0">
            {isLoading ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {[...Array(6)].map((_, i) => (
                  <Card key={i} className="bg-gray-900/50 border-gray-700 overflow-hidden">
                    <CardContent className="p-0">
                      <div className="w-full h-48 bg-gray-800 animate-pulse" />
                      <div className="p-4">
                        <div className="h-5 w-3/4 mb-2 bg-gray-800 animate-pulse rounded" />
                        <div className="h-4 w-full mb-3 bg-gray-800 animate-pulse rounded" />
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            ) : categoryTree && categoryTree.length > 0 ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {categoryTree.map((category) => (
                  <Link key={category.id} href={`/categories/${category.slug || category.id}`}>
                    <Card className="bg-gray-900/50 border-gray-700 hover:border-[#ff9000]/50 transition-all duration-300 group overflow-hidden">
                      <CardContent className="p-0">
                        <div className="relative w-full h-48 bg-gradient-to-br from-[#ff9000]/20 to-[#ff9000]/10 flex items-center justify-center">
                          {category.image ? (
                            <Image
                              src={category.image}
                              alt={category.name}
                              width={120}
                              height={120}
                              className="object-cover rounded-lg group-hover:scale-110 transition-transform duration-500"
                            />
                          ) : (
                            <FolderTree className="w-16 h-16 text-[#ff9000]" />
                          )}
                        </div>
                        
                        <div className="p-4">
                          <h3 className="text-lg font-semibold text-white mb-2 group-hover:text-[#ff9000] transition-colors line-clamp-1">
                            {category.name}
                          </h3>
                          {category.description && (
                            <p className="text-gray-400 mb-3 text-sm line-clamp-2">
                              {category.description}
                            </p>
                          )}
                          <div className="flex items-center justify-between">
                            <Badge variant="secondary" className="bg-gray-800 text-gray-300 text-xs">
                              {category.product_count || 0} items
                            </Badge>
                            <ArrowRight className="h-4 w-4 text-[#ff9000] group-hover:translate-x-1 transition-transform duration-300" />
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </Link>
                ))}
              </div>
            ) : (
              <div className="text-center py-16">
                <FolderTree className="h-20 w-20 text-gray-600 mx-auto mb-4" />
                <h3 className="text-2xl font-semibold text-white mb-3">No categories found</h3>
                <p className="text-gray-400">No categories are available at the moment.</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
