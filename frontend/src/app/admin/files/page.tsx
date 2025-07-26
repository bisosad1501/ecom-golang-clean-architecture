import { Metadata } from 'next'
import AdminFilesPage from '@/components/admin/admin-files-page'

export const metadata: Metadata = {
  title: 'Files | Admin Dashboard',
  description: 'Manage uploaded files, images, and documents.',
}

export default function Page() {
  return <AdminFilesPage />
}
