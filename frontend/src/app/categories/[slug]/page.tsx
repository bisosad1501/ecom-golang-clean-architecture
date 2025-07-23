import { Metadata } from 'next'
import { Suspense } from 'react'
import { notFound } from 'next/navigation'
import { CategoryDetailPage } from '@/components/pages/category-detail-page'

interface Props {
  params: Promise<{
    slug: string
  }>
  searchParams: Promise<{
    page?: string
    limit?: string
    sort_by?: string
    sort_order?: string
  }>
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  // In a real app, you'd fetch the category data here
  const { slug } = await params
  const categoryName = slug.replace(/-/g, ' ').replace(/\b\w/g, l => l.toUpperCase())

  return {
    title: `${categoryName} | EcomStore`,
    description: `Browse ${categoryName} products. Find the best deals and latest items in our ${categoryName} category.`,
    keywords: [categoryName, 'products', 'shopping', 'ecommerce'],
    openGraph: {
      title: `${categoryName} Products | EcomStore`,
      description: `Browse ${categoryName} products. Find the best deals and latest items in our ${categoryName} category.`,
      type: 'website',
    },
  }
}

async function CategoryDetailContent({ params, searchParams }: Props) {
  const { slug } = await params
  const resolvedSearchParams = await searchParams

  return (
    <CategoryDetailPage
      slug={slug}
      searchParams={resolvedSearchParams}
    />
  )
}

export default function Page({ params, searchParams }: Props) {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 rounded-full bg-[#ff9000] animate-spin mx-auto mb-4"></div>
          <p className="text-white">Loading category...</p>
        </div>
      </div>
    }>
      <CategoryDetailContent params={params} searchParams={searchParams} />
    </Suspense>
  )
}
