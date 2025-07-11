import { ProductDetailPage } from '@/components/pages/product-detail-page'

interface ProductPageProps {
  params: Promise<{
    id: string
  }>
}

export default async function ProductPage({ params }: ProductPageProps) {
  const { id } = await params
  return <ProductDetailPage productId={id} />
}
