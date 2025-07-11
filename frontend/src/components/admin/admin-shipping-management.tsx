'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { 
  Dialog, 
  DialogContent, 
  DialogDescription, 
  DialogHeader, 
  DialogTitle,
  DialogFooter 
} from '@/components/ui/dialog'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Truck, 
  Package, 
  MapPin, 
  Clock, 
  CheckCircle, 
  MessageCircle,
  Edit,
  Save,
  X
} from 'lucide-react'
import { Order, OrderEvent } from '@/types'
import { formatDate, formatPrice } from '@/lib/utils'
import { OrderTimeline } from '@/components/ui/order-timeline'
import { toast } from 'sonner'

interface AdminShippingManagementProps {
  order: Order
  events?: OrderEvent[]
  onUpdateShipping?: (data: any) => Promise<void>
  onUpdateDeliveryStatus?: (status: string) => Promise<void>
  onAddNote?: (note: string, isPublic: boolean) => Promise<void>
}

export function AdminShippingManagement({ 
  order, 
  events = [],
  onUpdateShipping,
  onUpdateDeliveryStatus,
  onAddNote
}: AdminShippingManagementProps) {
  const [showShippingModal, setShowShippingModal] = useState(false)
  const [showNoteModal, setShowNoteModal] = useState(false)
  const [isUpdating, setIsUpdating] = useState(false)
  
  // Shipping form state
  const [shippingData, setShippingData] = useState({
    tracking_number: order.tracking_number || '',
    carrier: order.carrier || '',
    shipping_method: order.shipping_method || '',
    tracking_url: order.tracking_url || ''
  })
  
  // Note form state
  const [noteData, setNoteData] = useState({
    note: '',
    is_public: false
  })

  const handleUpdateShipping = async () => {
    if (!onUpdateShipping) return
    
    setIsUpdating(true)
    try {
      await onUpdateShipping(shippingData)
      setShowShippingModal(false)
      toast.success('Shipping info updated successfully')
    } catch (error) {
      toast.error('Failed to update shipping info')
    } finally {
      setIsUpdating(false)
    }
  }

  const handleUpdateDeliveryStatus = async (status: string) => {
    if (!onUpdateDeliveryStatus) return
    
    setIsUpdating(true)
    try {
      await onUpdateDeliveryStatus(status)
      toast.success('Delivery status updated successfully')
    } catch (error) {
      toast.error('Failed to update delivery status')
    } finally {
      setIsUpdating(false)
    }
  }

  const handleAddNote = async () => {
    if (!onAddNote || !noteData.note.trim()) return
    
    setIsUpdating(true)
    try {
      await onAddNote(noteData.note, noteData.is_public)
      setShowNoteModal(false)
      setNoteData({ note: '', is_public: false })
      toast.success('Note added successfully')
    } catch (error) {
      toast.error('Failed to add note')
    } finally {
      setIsUpdating(false)
    }
  }

  return (
    <div className="space-y-6">
      {/* Shipping Status Overview */}
      <Card className="bg-gray-900/50 border-gray-800">
        <CardHeader>
          <CardTitle className="text-white flex items-center gap-2">
            <Truck className="h-5 w-5" />
            Shipping & Delivery Management
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Current Status */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="space-y-2">
              <Label className="text-gray-400">Order Status</Label>
              <Badge variant="outline" className="w-fit">
                {order.status}
              </Badge>
            </div>
            <div className="space-y-2">
              <Label className="text-gray-400">Fulfillment Status</Label>
              <Badge variant="outline" className="w-fit">
                {order.fulfillment_status}
              </Badge>
            </div>
            <div className="space-y-2">
              <Label className="text-gray-400">Priority</Label>
              <Badge variant="outline" className="w-fit">
                {order.priority}
              </Badge>
            </div>
          </div>

          <Separator className="bg-gray-800" />

          {/* Shipping Information */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h4 className="text-white font-medium">Shipping Information</h4>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setShowShippingModal(true)}
                className="border-gray-700 text-gray-300 hover:bg-gray-800"
              >
                <Edit className="h-4 w-4 mr-2" />
                Edit Shipping
              </Button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label className="text-gray-400">Tracking Number</Label>
                <div className="text-white">
                  {order.tracking_number || 'Not assigned'}
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-400">Carrier</Label>
                <div className="text-white">
                  {order.carrier || 'Not assigned'}
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-400">Shipping Method</Label>
                <div className="text-white">
                  {order.shipping_method || 'Not assigned'}
                </div>
              </div>
              <div className="space-y-2">
                <Label className="text-gray-400">Delivery Status</Label>
                <div className="flex items-center gap-2">
                  <Badge variant="outline">
                    {order.is_delivered ? 'Delivered' : order.is_shipped ? 'Shipped' : 'Pending'}
                  </Badge>
                  {order.can_be_delivered && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => handleUpdateDeliveryStatus('delivered')}
                      disabled={isUpdating}
                      className="border-green-700 text-green-300 hover:bg-green-800"
                    >
                      <CheckCircle className="h-4 w-4 mr-2" />
                      Mark Delivered
                    </Button>
                  )}
                </div>
              </div>
            </div>
          </div>

          <Separator className="bg-gray-800" />

          {/* Quick Actions */}
          <div className="flex flex-wrap gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowNoteModal(true)}
              className="border-gray-700 text-gray-300 hover:bg-gray-800"
            >
              <MessageCircle className="h-4 w-4 mr-2" />
              Add Note
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Order Timeline */}
      <OrderTimeline 
        orderId={order.id} 
        events={events}
        showPrivateEvents={true}
      />

      {/* Shipping Update Modal */}
      <Dialog open={showShippingModal} onOpenChange={setShowShippingModal}>
        <DialogContent className="bg-gray-900 border-gray-800 text-white">
          <DialogHeader>
            <DialogTitle>Update Shipping Information</DialogTitle>
            <DialogDescription className="text-gray-400">
              Update tracking and shipping details for order {order.order_number}
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="tracking_number">Tracking Number</Label>
              <Input
                id="tracking_number"
                value={shippingData.tracking_number}
                onChange={(e) => setShippingData(prev => ({ ...prev, tracking_number: e.target.value }))}
                className="bg-gray-800 border-gray-700 text-white"
                placeholder="Enter tracking number"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="carrier">Carrier</Label>
              <Select
                value={shippingData.carrier}
                onValueChange={(value) => setShippingData(prev => ({ ...prev, carrier: value }))}
              >
                <SelectTrigger className="bg-gray-800 border-gray-700 text-white">
                  <SelectValue placeholder="Select carrier" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="FedEx">FedEx</SelectItem>
                  <SelectItem value="UPS">UPS</SelectItem>
                  <SelectItem value="DHL">DHL</SelectItem>
                  <SelectItem value="USPS">USPS</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="shipping_method">Shipping Method</Label>
              <Input
                id="shipping_method"
                value={shippingData.shipping_method}
                onChange={(e) => setShippingData(prev => ({ ...prev, shipping_method: e.target.value }))}
                className="bg-gray-800 border-gray-700 text-white"
                placeholder="e.g., Express Delivery"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="tracking_url">Tracking URL</Label>
              <Input
                id="tracking_url"
                value={shippingData.tracking_url}
                onChange={(e) => setShippingData(prev => ({ ...prev, tracking_url: e.target.value }))}
                className="bg-gray-800 border-gray-700 text-white"
                placeholder="https://..."
              />
            </div>
          </div>
          
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowShippingModal(false)}
              className="border-gray-700 text-gray-300 hover:bg-gray-800"
            >
              <X className="h-4 w-4 mr-2" />
              Cancel
            </Button>
            <Button
              onClick={handleUpdateShipping}
              disabled={isUpdating}
              className="bg-blue-600 hover:bg-blue-700 text-white"
            >
              <Save className="h-4 w-4 mr-2" />
              {isUpdating ? 'Updating...' : 'Update Shipping'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Add Note Modal */}
      <Dialog open={showNoteModal} onOpenChange={setShowNoteModal}>
        <DialogContent className="bg-gray-900 border-gray-800 text-white">
          <DialogHeader>
            <DialogTitle>Add Order Note</DialogTitle>
            <DialogDescription className="text-gray-400">
              Add a note to order {order.order_number}
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="note">Note</Label>
              <Textarea
                id="note"
                value={noteData.note}
                onChange={(e) => setNoteData(prev => ({ ...prev, note: e.target.value }))}
                className="bg-gray-800 border-gray-700 text-white"
                placeholder="Enter your note..."
                rows={4}
              />
            </div>
            
            <div className="flex items-center space-x-2">
              <input
                type="checkbox"
                id="is_public"
                checked={noteData.is_public}
                onChange={(e) => setNoteData(prev => ({ ...prev, is_public: e.target.checked }))}
                className="rounded border-gray-700"
              />
              <Label htmlFor="is_public" className="text-sm">
                Make this note visible to customer
              </Label>
            </div>
          </div>
          
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => setShowNoteModal(false)}
              className="border-gray-700 text-gray-300 hover:bg-gray-800"
            >
              <X className="h-4 w-4 mr-2" />
              Cancel
            </Button>
            <Button
              onClick={handleAddNote}
              disabled={isUpdating || !noteData.note.trim()}
              className="bg-blue-600 hover:bg-blue-700 text-white"
            >
              <MessageCircle className="h-4 w-4 mr-2" />
              {isUpdating ? 'Adding...' : 'Add Note'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
