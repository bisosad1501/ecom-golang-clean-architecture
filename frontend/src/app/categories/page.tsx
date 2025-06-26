import { Metadata } from 'next'
import { Suspense } from 'react'
import { SimpleCategoryPage } from '@/components/layout/simple-category-page'

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
  return <SimpleCategoryPage />
}

export default function Page() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <CategoriesContent />
    </Suspense>
  )
}
