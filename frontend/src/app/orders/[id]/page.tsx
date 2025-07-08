import { OrderDetailPage } from '@/components/pages/OrderDetailPage'

interface Props {
  params: {
    id: string
  }
}

export default function OrderDetail({ params }: Props) {
  return <OrderDetailPage orderId={params.id} />
}
