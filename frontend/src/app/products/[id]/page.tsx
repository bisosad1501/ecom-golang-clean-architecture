import { ProductDetailPage } from '@/components/pages/product-detail-page'

interface ProductPageProps {
  params: {
    id: string
  }
}

export default function ProductPage({ params }: ProductPageProps) {
  return <ProductDetailPage productId={params.id} />
}
