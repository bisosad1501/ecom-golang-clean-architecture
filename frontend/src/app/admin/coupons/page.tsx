import { Metadata } from 'next'
import AdminCouponsPage from '@/components/admin/admin-coupons-page'

export const metadata: Metadata = {
  title: 'Coupons | Admin Dashboard',
  description: 'Manage discount coupons, promotional codes, and marketing campaigns.',
}

export default function Page() {
  return <AdminCouponsPage />
}
