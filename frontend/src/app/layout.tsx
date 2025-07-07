import type { Metadata } from "next";
import "./globals.css";
import { Providers } from '@/components/providers'
import { ConditionalLayout } from '@/components/layout/conditional-layout'
import { CartSidebar } from '@/components/cart/cart-sidebar'
import { Toaster } from 'sonner'
import { DEFAULT_SEO } from '@/constants'

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'),
  title: DEFAULT_SEO.title,
  description: DEFAULT_SEO.description,
  keywords: DEFAULT_SEO.keywords,
  openGraph: {
    title: DEFAULT_SEO.title,
    description: DEFAULT_SEO.description,
    images: [DEFAULT_SEO.ogImage],
  },
  twitter: {
    card: 'summary_large_image',
    title: DEFAULT_SEO.title,
    description: DEFAULT_SEO.description,
    images: [DEFAULT_SEO.ogImage],
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning className="dark">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@100;200;300;400;500;600;700;800;900&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="bg-black text-white font-inter" suppressHydrationWarning>
        <Providers>
          <ConditionalLayout>
            {children}
          </ConditionalLayout>
          <CartSidebar />
          <Toaster position="top-right" richColors />
        </Providers>
      </body>
    </html>
  );
}
