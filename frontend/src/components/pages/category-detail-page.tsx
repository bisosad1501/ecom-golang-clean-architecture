'use client'

import { useState, useMemo } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ChevronRight, Home } from 'lucide-react'
import { ArrowLeft, Grid, List, SlidersHorizontal, Package, FolderTree } from 'lucide-react'
import { AnimatedBackground } from '@/components/ui/animated-background'
import { categoryService } from '@/lib/services/categories'
import { productService } from '@/lib/services/products'
import { Category, Product } from '@/types'

interface Props {
  slug: string
  searchParams: {
    page?: string
    limit?: string
    sort_by?: string
    sort_order?: string
  }
}

export function CategoryDetailPage({ slug, searchParams }: Props) {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid')
  
  // Parse search params
  const page = parseInt(searchParams.page || '1')
  const limit = parseInt(searchParams.limit || '20')
  const sortBy = searchParams.sort_by || 'name'
  const sortOrder = searchParams.sort_order || 'asc'

  // Fetch category by slug
  const { 
    data: category, 
    isLoading: isLoadingCategory,
    error: categoryError 
  } = useQuery({
    queryKey: ['category', slug],
    queryFn: () => categoryService.getCategoryBySlug(slug),
  })

  // Fetch category products
  const { 
    data: productsData, 
    isLoading: isLoadingProducts 
  } = useQuery({
    queryKey: ['category-products', category?.id, page, limit, sortBy, sortOrder],
    queryFn: () => {
      if (!category?.id) return null
      return productService.getProducts({
        category_id: category.id,
        page,
        limit,
        sort_by: sortBy,
        sort_order: sortOrder as 'asc' | 'desc'
      })
    },
    enabled: !!category?.id,
  })

  // Fetch subcategories
  const { 
    data: subcategories 
  } = useQuery({
    queryKey: ['subcategories', category?.id],
    queryFn: () => {
      if (!category?.id) return []
      return categoryService.getSubcategories(category.id)
    },
    enabled: !!category?.id,
  })

  if (categoryError) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black text-white">
        <div className="container mx-auto px-4 py-8">
          <Card className="bg-gray-900/50 border-gray-700 max-w-md mx-auto">
            <CardContent className="p-8 text-center">
              <h2 className="text-2xl font-bold text-white mb-4">Category Not Found</h2>
              <p className="text-gray-400 mb-6">
                The category you're looking for doesn't exist or has been removed.
              </p>
              <Button asChild variant="outline" className="bg-gray-800 border-gray-600 text-white hover:bg-gray-700">
                <Link href="/categories">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Back to Categories
                </Link>
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Nike-style Breadcrumb */}
        <div className="py-4 border-b border-gray-200">
          <nav className="flex items-center space-x-2 text-sm">
            <Link href="/" className="text-gray-600 hover:text-black transition-colors flex items-center">
              <Home className="h-4 w-4" />
              <span className="ml-1">Home</span>
            </Link>
            <ChevronRight className="h-4 w-4 text-gray-400" />
            <Link href="/categories" className="text-gray-600 hover:text-black transition-colors">
              Categories
            </Link>
            <ChevronRight className="h-4 w-4 text-gray-400" />
            <span className="text-black font-medium">
              {category?.name || slug}
            </span>
          </nav>
        </div>

        {/* Category Header */}
        {isLoadingCategory ? (
          <div className="mb-12">
            <div className="h-8 w-48 bg-gray-800 animate-pulse rounded mb-4" />
            <div className="h-4 w-96 bg-gray-800 animate-pulse rounded" />
          </div>
        ) : category && (
          <div className="mb-12">
            <div className="flex items-center gap-4 mb-6">
              <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-[#ff9000] to-[#ff9000]/80 flex items-center justify-center shadow-lg shadow-[#ff9000]/20">
                {category.image ? (
                  <Image
                    src={category.image}
                    alt={category.name}
                    width={64}
                    height={64}
                    className="rounded-2xl object-cover"
                  />
                ) : (
                  <FolderTree className="h-8 w-8 text-white" />
                )}
              </div>
              <div>
                <h1 className="text-4xl lg:text-5xl font-bold bg-gradient-to-r from-white via-gray-200 to-[#ff9000] bg-clip-text text-transparent">
                  {category.name}
                </h1>
                {category.description && (
                  <p className="text-xl text-gray-400 mt-2">
                    {category.description}
                  </p>
                )}
              </div>
            </div>

            {/* Category Stats */}
            <div className="flex items-center gap-4 mb-8">
              <Badge className="bg-[#ff9000]/10 text-[#ff9000] border-[#ff9000]/20">
                {productsData?.total || 0} Products
              </Badge>
              {category.level !== undefined && (
                <Badge className="bg-gray-800 text-gray-300 border-gray-600">
                  Level {category.level}
                </Badge>
              )}
              {category.path && (
                <Badge variant="outline" className="bg-gray-900/50 border-gray-700 text-gray-300">
                  {category.path}
                </Badge>
              )}
            </div>
          </div>
        )}

        {/* Subcategories */}
        {subcategories && subcategories.length > 0 && (
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-white mb-6">Subcategories</h2>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
              {subcategories.map((subcat) => (
                <Link key={subcat.id} href={`/categories/${subcat.slug}`}>
                  <Card className="bg-gray-900/50 border-gray-700 hover:border-[#ff9000]/50 transition-all duration-300 group">
                    <CardContent className="p-4 text-center">
                      <div className="w-12 h-12 rounded-xl bg-gray-800 flex items-center justify-center mx-auto mb-3 group-hover:bg-[#ff9000]/20 transition-colors">
                        <Package className="h-6 w-6 text-gray-400 group-hover:text-[#ff9000]" />
                      </div>
                      <h3 className="text-sm font-medium text-white group-hover:text-[#ff9000] transition-colors">
                        {subcat.name}
                      </h3>
                    </CardContent>
                  </Card>
                </Link>
              ))}
            </div>
          </div>
        )}

        {/* Products Section */}
        <div className="mb-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-white">Products</h2>
            
            {/* View Toggle */}
            <div className="flex items-center gap-2">
              <div className="flex items-center gap-1 bg-gray-900/50 border border-gray-700 rounded-xl p-1">
                <Button
                  variant={viewMode === 'grid' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('grid')}
                  className={`flex items-center gap-2 rounded-lg transition-all duration-300 ${
                    viewMode === 'grid' 
                      ? 'bg-[#ff9000] text-white shadow-lg shadow-[#ff9000]/20' 
                      : 'text-gray-400 hover:text-white hover:bg-gray-800'
                  }`}
                >
                  <Grid className="h-4 w-4" />
                  Grid
                </Button>
                <Button
                  variant={viewMode === 'list' ? 'default' : 'ghost'}
                  size="sm"
                  onClick={() => setViewMode('list')}
                  className={`flex items-center gap-2 rounded-lg transition-all duration-300 ${
                    viewMode === 'list' 
                      ? 'bg-[#ff9000] text-white shadow-lg shadow-[#ff9000]/20' 
                      : 'text-gray-400 hover:text-white hover:bg-gray-800'
                  }`}
                >
                  <List className="h-4 w-4" />
                  List
                </Button>
              </div>
            </div>
          </div>

          {/* Products Grid/List */}
          {isLoadingProducts ? (
            <div className={viewMode === 'grid' 
              ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
              : "space-y-4"
            }>
              {[...Array(8)].map((_, i) => (
                <Card key={i} className="bg-gray-900/50 border-gray-700">
                  <CardContent className="p-0">
                    <div className="w-full h-48 bg-gray-800 animate-pulse" />
                    <div className="p-6">
                      <div className="h-6 w-3/4 mb-3 bg-gray-800 animate-pulse rounded" />
                      <div className="h-4 w-1/2 mb-4 bg-gray-800 animate-pulse rounded" />
                      <div className="h-8 w-1/3 bg-gray-800 animate-pulse rounded" />
                    </div>
                  </CardContent>
                </Card>
              ))}
            </div>
          ) : productsData?.products && productsData.products.length > 0 ? (
            <div className={viewMode === 'grid' 
              ? "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
              : "space-y-4"
            }>
              {/* Products will be rendered here - you can import ProductCard component */}
              {productsData.products.map((product) => (
                <div key={product.id} className="text-white">
                  {/* Placeholder for product card */}
                  Product: {product.name}
                </div>
              ))}
            </div>
          ) : (
            <Card className="bg-gray-900/50 border-gray-700">
              <CardContent className="p-8 text-center">
                <Package className="h-16 w-16 text-gray-400 mx-auto mb-4" />
                <h3 className="text-xl font-semibold text-white mb-2">No Products Found</h3>
                <p className="text-gray-400">
                  This category doesn't have any products yet.
                </p>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  )
}
