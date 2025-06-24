import { Metadata } from 'next'
import ContactPage from '@/components/pages/contact-page'

export const metadata: Metadata = {
  title: 'Contact Us | EcomStore',
  description: 'Get in touch with our support team. We\'re here to help with orders, shipping, returns, and any questions you may have.',
  keywords: ['contact', 'support', 'help', 'customer service', 'phone', 'email', 'chat'],
  openGraph: {
    title: 'Contact EcomStore | Customer Support',
    description: 'Get in touch with our support team. We\'re here to help with orders, shipping, returns, and any questions you may have.',
    type: 'website',
  },
}

export default function Page() {
  return <ContactPage />
}
