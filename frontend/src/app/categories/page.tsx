import { Metadata } from 'next'
import { Suspense } from 'react'
import { CategoriesPage } from '@/components/pages/categories-page'

export const metadata: Metadata = {
  title: 'Categories | EcomStore',
  description: 'Browse products by category. Find exactly what you\'re looking for with our organized product categories.',
  keywords: ['categories', 'products', 'shopping', 'browse', 'ecommerce'],
  openGraph: {
    title: 'Product Categories | EcomStore',
    description: 'Browse products by category. Find exactly what you\'re looking for with our organized product categories.',
    type: 'website',
  },
}

function CategoriesContent() {
  return <CategoriesPage />
}

export default function Page() {
  return (
    <Suspense fallback={
      <div className="min-h-screen bg-gradient-to-br from-black via-gray-900 to-black flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 rounded-full bg-[#ff9000] animate-spin mx-auto mb-4"></div>
          <p className="text-white">Loading categories...</p>
        </div>
      </div>
    }>
      <CategoriesContent />
    </Suspense>
  )
}
