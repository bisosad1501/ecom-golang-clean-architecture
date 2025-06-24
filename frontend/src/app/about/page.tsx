import { Metadata } from 'next'
import AboutPage from '@/components/pages/about-page'

export const metadata: Metadata = {
  title: 'About Us | EcomStore',
  description: 'Learn about our mission, values, and the passionate team behind EcomStore. Discover our journey from startup to global e-commerce platform.',
  keywords: ['about us', 'company', 'mission', 'values', 'team', 'ecommerce', 'story'],
  openGraph: {
    title: 'About EcomStore | Our Story & Mission',
    description: 'Learn about our mission, values, and the passionate team behind EcomStore. Discover our journey from startup to global e-commerce platform.',
    type: 'website',
  },
}

export default function Page() {
  return <AboutPage />
}
