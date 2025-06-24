import { Metadata } from 'next'
import CategoriesPage from '@/components/pages/categories-page'

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

export default function Page() {
  return <CategoriesPage />
}
