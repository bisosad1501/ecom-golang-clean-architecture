// Route: /orders/[id] (Individual order details)
// Component: OrderDetailPage from order-detail-page.tsx
import { OrderDetailPage } from '@/components/pages/order-detail-page'

interface Props {
  params: Promise<{
    id: string
  }>
}

export default async function OrderDetail({ params }: Props) {
  const { id } = await params
  return <OrderDetailPage orderId={id} />
}
