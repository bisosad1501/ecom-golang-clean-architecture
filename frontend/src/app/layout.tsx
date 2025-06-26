import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Providers } from '@/components/providers'
import { ConditionalLayout } from '@/components/layout/conditional-layout'
import { CartSidebar } from '@/components/cart/cart-sidebar'
import { Toaster } from 'sonner'
import { APP_NAME, DEFAULT_SEO } from '@/constants'

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
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
      <body className={`${inter.className} bg-black text-white`}>
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
