'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { useSearchParams } from 'next/navigation'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Search, Grid, FolderTree, ArrowRight, Package } from 'lucide-react'
import { useCategories } from '@/hooks/use-categories'
import { Category } from '@/types'

export default function CategoriesPage() {
  const searchParams = useSearchParams()
  const [searchQuery, setSearchQuery] = useState(searchParams.get('search') || '')
  const [viewMode, setViewMode] = useState<'grid' | 'tree'>('grid')
  
  const { 
    data: categories, 
    isLoading, 
    error 
  } = useCategories()

  // Filter categories based on search query
  const filteredCategories = categories?.filter(category =>
    category.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    category.description?.toLowerCase().includes(searchQuery.toLowerCase())
  ) || []

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    // Update URL without page reload
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
      <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background py-16">
        <div className="container mx-auto px-4">
          <div className="text-center">
            <h1 className="text-3xl font-bold text-destructive mb-4">Error Loading Categories</h1>
            <p className="text-muted-foreground mb-8">
              We're having trouble loading the categories. Please try again later.
            </p>
            <Button onClick={() => window.location.reload()}>
              Try Again
            </Button>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-muted/20 to-background">
      <div className="container mx-auto px-4 py-12">
        {/* Enhanced Header */}
        <div className="mb-12 text-center">
          <div className="flex items-center justify-center gap-3 mb-6">
            <div className="w-12 h-12 rounded-2xl bg-gradient-to-br from-primary-500 to-violet-600 flex items-center justify-center shadow-large">
              <FolderTree className="h-6 w-6 text-white" />
            </div>
            <span className="text-primary font-semibold">PRODUCT CATEGORIES</span>
          </div>
          
          <h1 className="text-4xl lg:text-6xl font-bold text-foreground mb-6">
            Browse by <span className="text-gradient">Category</span>
          </h1>
          <p className="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
            Discover products organized by category. Find exactly what you're looking for 
            with our intuitive category navigation.
          </p>
        </div>

        {/* Search and View Controls */}
        <div className="mb-12">
          <Card variant="elevated" className="border-0 shadow-large">
            <CardContent className="p-8">
              <div className="flex flex-col lg:flex-row items-center gap-6">
                <div className="flex-1 w-full">
                  <div className="relative">
                    <Input
                      placeholder="Search categories..."
                      value={searchQuery}
                      onChange={(e) => handleSearch(e.target.value)}
                      leftIcon={<Search className="h-5 w-5" />}
                      size="lg"
                      className="w-full pr-12"
                    />
                    {searchQuery && (
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleSearch('')}
                        className="absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0"
                      >
                        ×
                      </Button>
                    )}
                  </div>
                </div>
                
                <div className="flex items-center bg-background border-2 border-border rounded-2xl p-1 shadow-medium">
                  <Button
                    variant={viewMode === 'grid' ? 'default' : 'ghost'}
                    size="sm"
                    onClick={() => setViewMode('grid')}
                    className="rounded-xl h-10 px-4"
                  >
                    <Grid className="h-4 w-4 mr-2" />
                    Grid
                  </Button>
                  <Button
                    variant={viewMode === 'tree' ? 'default' : 'ghost'}
                    size="sm"
                    onClick={() => setViewMode('tree')}
                    className="rounded-xl h-10 px-4"
                  >
                    <FolderTree className="h-4 w-4 mr-2" />
                    Tree
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Categories Content */}
        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
            {[...Array(8)].map((_, i) => (
              <Card key={i} variant="elevated" className="border-0 shadow-large">
                <CardContent className="p-6">
                  <div className="w-full h-48 rounded-2xl mb-4 bg-muted animate-pulse" />
                  <div className="h-6 w-3/4 mb-2 bg-muted animate-pulse rounded" />
                  <div className="h-4 w-1/2 mb-4 bg-muted animate-pulse rounded" />
                  <div className="h-4 w-1/4 bg-muted animate-pulse rounded" />
                </CardContent>
              </Card>
            ))}
          </div>
        ) : filteredCategories.length > 0 ? (
          viewMode === 'grid' ? (
            <CategoryHierarchyGrid categories={filteredCategories} />
          ) : (
            <CategoryTree categories={filteredCategories} />
          )
        ) : (
          <EmptyState searchQuery={searchQuery} onClearSearch={() => handleSearch('')} />
        )}
      </div>
    </div>
  )
}

// Category Card Component
function CategoryCard({ 
  category, 
  isRoot = false, 
  isOrphaned = false 
}: { 
  category: Category
  isRoot?: boolean
  isOrphaned?: boolean
}) {
  return (
    <Card variant="elevated" className={`border-0 shadow-large hover:shadow-xl transition-all duration-300 group overflow-hidden ${
      isOrphaned ? 'ring-2 ring-yellow-400 ring-opacity-50' : ''
    }`}>
      <CardContent className="p-0">
        {/* Category Image */}
        <div className="relative h-48 overflow-hidden bg-muted">
          {category.image ? (
            <Image
              src={category.image}
              alt={category.name}
              fill
              className="object-cover group-hover:scale-105 transition-transform duration-300"
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center text-muted-foreground bg-gradient-to-br from-muted to-muted/50">
              {isRoot ? (
                <FolderTree className="w-16 h-16" />
              ) : (
                <Package className="w-16 h-16" />
              )}
            </div>
          )}
          
          {/* Category Type Badge */}
          <div className="absolute top-3 left-3">
            {isOrphaned ? (
              <Badge variant="destructive" className="text-xs">
                ⚠ Orphaned
              </Badge>
            ) : isRoot ? (
              <Badge variant="default" className="text-xs bg-blue-600">
                📁 Main Category
              </Badge>
            ) : (
              <Badge variant="secondary" className="text-xs">
                📂 Subcategory
              </Badge>
            )}
          </div>
          
          {/* Level indicator */}
          <div className="absolute top-3 right-3">
            <Badge variant="outline" className="text-xs bg-white/90">
              Level {category.level || 0}
            </Badge>
          </div>
          
          {/* Overlay */}
          <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
          
          {/* View Products Button */}
          <div className="absolute bottom-4 left-4 right-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <Button asChild size="sm" variant="secondary" className="w-full">
              <Link href={`/products?category=${category.id}`}>
                View Products
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
        
        {/* Category Info */}
        <div className="p-6">
          <div className="mb-3">
            <h3 className="text-xl font-bold text-foreground mb-1 group-hover:text-primary transition-colors">
              {category.name}
            </h3>
            
            {/* Path breadcrumb for non-root categories */}
            {!isRoot && category.path && (
              <p className="text-xs text-muted-foreground mb-2 font-mono bg-muted/50 px-2 py-1 rounded">
                {category.path}
              </p>
            )}
          </div>
          
          {category.description && (
            <p className="text-muted-foreground text-sm mb-4 line-clamp-2">
              {category.description}
            </p>
          )}
          
          <div className="flex items-center justify-between">
            <Badge variant="outline" className="text-xs">
              {category.product_count || 0} products
            </Badge>
            
            <Button asChild variant="ghost" size="sm">
              <Link href={`/products?category=${category.id}`}>
                Browse
                <ArrowRight className="ml-1 h-3 w-3" />
              </Link>
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

// Category Tree Component
function CategoryTree({ categories }: { categories: Category[] }) {
  // Group categories by parent
  const rootCategories = categories.filter(cat => !cat.parent_id)
  
  return (
    <div className="space-y-6">
      {rootCategories.map((category) => (
        <CategoryTreeNode key={category.id} category={category} allCategories={categories} />
      ))}
    </div>
  )
}

// Category Tree Node Component
function CategoryTreeNode({ category, allCategories, level = 0 }: { 
  category: Category
  allCategories: Category[]
  level?: number 
}) {
  const children = allCategories.filter(cat => cat.parent_id === category.id)
  const [isExpanded, setIsExpanded] = useState(level < 2)
  
  return (
    <div className={`${level > 0 ? 'ml-8' : ''}`}>
      <Card variant={level === 0 ? "elevated" : "outlined"} className="border-0 shadow-medium">
        <CardContent className="p-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              {category.image && (
                <div className="w-12 h-12 rounded-xl overflow-hidden">
                  <Image
                    src={category.image}
                    alt={category.name}
                    width={48}
                    height={48}
                    className="object-cover"
                  />
                </div>
              )}
              
              <div>
                <h3 className="text-lg font-semibold text-foreground">
                  {category.name}
                </h3>
                {category.description && (
                  <p className="text-sm text-muted-foreground">
                    {category.description}
                  </p>
                )}
              </div>
            </div>
            
            <div className="flex items-center gap-3">
              <Badge variant="outline">
                {category.product_count || 0} products
              </Badge>
              
              <Button asChild variant="outline" size="sm">
                <Link href={`/products?category=${category.id}`}>
                  View Products
                </Link>
              </Button>
              
              {children.length > 0 && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsExpanded(!isExpanded)}
                >
                  {isExpanded ? 'Collapse' : 'Expand'} ({children.length})
                </Button>
              )}
            </div>
          </div>
        </CardContent>
      </Card>
      
      {isExpanded && children.length > 0 && (
        <div className="mt-4 space-y-4">
          {children.map((child) => (
            <CategoryTreeNode
              key={child.id}
              category={child}
              allCategories={allCategories}
              level={level + 1}
            />
          ))}
        </div>
      )}
    </div>
  )
}

// Category Hierarchy Grid Component
function CategoryHierarchyGrid({ categories }: { categories: Category[] }) {
  // Group categories by hierarchy
  const rootCategories = categories.filter(cat => !cat.parent_id)
  const childCategories = categories.filter(cat => cat.parent_id)
  
  return (
    <div className="space-y-12">
      {/* Root Categories Section */}
      <div>
        <div className="flex items-center gap-3 mb-8">
          <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center">
            <FolderTree className="h-4 w-4 text-white" />
          </div>
          <h2 className="text-2xl font-bold text-foreground">Main Categories</h2>
          <Badge variant="outline" className="ml-2">
            {rootCategories.length} categories
          </Badge>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {rootCategories.map((category) => (
            <CategoryCard key={category.id} category={category} isRoot={true} />
          ))}
        </div>
      </div>

      {/* Root Categories with their Children */}
      {rootCategories.map((rootCategory) => {
        const children = childCategories.filter(cat => cat.parent_id === rootCategory.id)
        
        if (children.length === 0) return null
        
        return (
          <div key={`${rootCategory.id}-children`}>
            <div className="flex items-center gap-3 mb-8">
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-green-500 to-green-600 flex items-center justify-center">
                <span className="text-white text-xs font-bold">└─</span>
              </div>
              <h2 className="text-2xl font-bold text-foreground">
                <span className="text-muted-foreground">{rootCategory.name}</span> Subcategories
              </h2>
              <Badge variant="outline" className="ml-2">
                {children.length} subcategories
              </Badge>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {children.map((category) => (
                <CategoryCard key={category.id} category={category} isRoot={false} />
              ))}
            </div>
          </div>
        )
      })}
      
      {/* Orphaned Categories (if any) */}
      {childCategories.filter(child => 
        !rootCategories.some(root => root.id === child.parent_id)
      ).length > 0 && (
        <div>
          <div className="flex items-center gap-3 mb-8">
            <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-yellow-500 to-orange-500 flex items-center justify-center">
              <span className="text-white text-xs">⚠</span>
            </div>
            <h2 className="text-2xl font-bold text-foreground">Uncategorized Items</h2>
            <Badge variant="destructive" className="ml-2">
              Needs attention
            </Badge>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {childCategories.filter(child => 
              !rootCategories.some(root => root.id === child.parent_id)
            ).map((category) => (
              <CategoryCard key={category.id} category={category} isRoot={false} isOrphaned={true} />
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

// Empty State Component
function EmptyState({ searchQuery, onClearSearch }: { 
  searchQuery: string
  onClearSearch: () => void 
}) {
  return (
    <div className="text-center py-24">
      <div className="relative mb-8">
        <div className="w-32 h-32 mx-auto rounded-full bg-gradient-to-br from-muted to-muted/50 flex items-center justify-center shadow-large">
          <FolderTree className="h-16 w-16 text-muted-foreground" />
        </div>
      </div>
      
      <h3 className="text-3xl font-bold text-foreground mb-4">
        {searchQuery ? 'No categories found' : 'No categories available'}
      </h3>
      <p className="text-xl text-muted-foreground mb-12 max-w-md mx-auto leading-relaxed">
        {searchQuery 
          ? `No categories found matching "${searchQuery}". Try adjusting your search terms.` 
          : 'Categories will appear here once they are added to the system.'
        }
      </p>
      
      {searchQuery && (
        <Button onClick={onClearSearch} size="xl" variant="outline">
          Clear Search
        </Button>
      )}
    </div>
  )
}
