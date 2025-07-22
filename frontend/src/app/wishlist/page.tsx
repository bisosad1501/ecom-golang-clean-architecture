import { Metadata } from 'next'
import { WishlistPage } from '@/components/pages/wishlist-page'

export const metadata: Metadata = {
  title: 'My Wishlist | BiHub',
  description: 'View and manage your saved products',
}

export default function Wishlist() {
  return <WishlistPage />
}
