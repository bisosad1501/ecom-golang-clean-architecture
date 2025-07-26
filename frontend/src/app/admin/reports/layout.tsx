import { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Reports | Admin Dashboard',
  description: 'Generate and download detailed business reports.',
}

export default function ReportsLayout({ children }: { children: React.ReactNode }) {
  return <>{children}</>;
}
