'use client'

import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { Download, X, FileText, Building2, MapPin, Calendar, Hash, CreditCard } from 'lucide-react'
import { formatPrice, formatDate, cn } from '@/lib/utils'
import { Order, OrderItem } from '@/types'

interface InvoicePreviewModalProps {
  isOpen: boolean
  onClose: () => void
  order: Order
}

export function InvoicePreviewModal({ isOpen, onClose, order }: InvoicePreviewModalProps) {
  const [isDownloading, setIsDownloading] = useState(false)

  const handleDownload = async () => {
    setIsDownloading(true)
    try {
      // Simulate download delay
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      // Generate PDF content using jsPDF or similar library
      // For now, we'll create a simple HTML that can be converted to PDF
      const invoiceHTML = generateInvoiceHTML(order)
      
      // Create a new window for printing/PDF generation
      const printWindow = window.open('', '_blank')
      if (printWindow) {
        printWindow.document.write(invoiceHTML)
        printWindow.document.close()
        
        // Wait for content to load then trigger print dialog
        printWindow.onload = () => {
          printWindow.print()
          printWindow.close()
        }
      }
      
      onClose()
    } catch (error) {
      console.error('Error generating PDF:', error)
    } finally {
      setIsDownloading(false)
    }
  }

  const generateInvoiceHTML = (order: Order) => {
    return `
<!DOCTYPE html>
<html>
<head>
    <title>Invoice ${order.order_number}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: white; }
        .header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 30px; border-bottom: 2px solid #ff9000; padding-bottom: 20px; }
        .company-info { display: flex; align-items: center; gap: 15px; }
        .logo { width: 60px; height: 60px; border-radius: 8px; }
        .company-name { font-size: 28px; font-weight: bold; margin: 0; display: flex; align-items: center; }
        .company-subtitle { color: #666; font-size: 14px; margin: 5px 0 0 0; }
        .invoice-title { font-size: 24px; margin: 0; color: #333; font-weight: bold; }
        .invoice-number { color: #666; font-size: 14px; margin: 5px 0 0 0; }
        .details { margin: 20px 0; display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
        .section-title { font-size: 16px; font-weight: bold; color: #333; margin-bottom: 10px; border-bottom: 1px solid #eee; padding-bottom: 5px; }
        .table { width: 100%; border-collapse: collapse; margin: 20px 0; }
        .table th, .table td { border: 1px solid #ddd; padding: 12px; text-align: left; }
        .table th { background-color: #f8f9fa; font-weight: bold; color: #333; }
        .table .total-row { font-weight: bold; background-color: #fff5e6; }
        .table .grand-total { background-color: #ff9000; color: white; font-size: 18px; }
        .footer { text-align: center; margin-top: 40px; padding-top: 20px; border-top: 1px solid #eee; color: #666; font-size: 12px; }
        @media print {
            body { margin: 0; }
            .no-print { display: none; }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="company-info">
            <div>
                <h1 class="company-name">
                    <span style="color: #333; font-size: 28px;">Bi</span><span style="background-color: #FF9000; color: #000; padding: 0 2px; margin-left: 2px; border-radius: 3px; font-weight: bold; letter-spacing: 0.3px; line-height: 1.1; display: inline-block;">hub</span>
                </h1>
                <p class="company-subtitle">Your Premium E-commerce Store</p>
            </div>
        </div>
        <div>
            <h2 class="invoice-title">INVOICE</h2>
            <p class="invoice-number">#${order.order_number}</p>
            <p class="invoice-number">${formatDate(new Date(order.created_at))}</p>
        </div>
    </div>
    
    <div class="details">
        <div>
            <h3 class="section-title">Bill To</h3>
            <div>
                <strong>${order.shipping_address?.first_name} ${order.shipping_address?.last_name}</strong><br>
                ${order.shipping_address?.company ? order.shipping_address.company + '<br>' : ''}
                ${order.shipping_address?.address1}<br>
                ${order.shipping_address?.address2 ? order.shipping_address.address2 + '<br>' : ''}
                ${order.shipping_address?.city}, ${order.shipping_address?.state} ${order.shipping_address?.zip_code}<br>
                ${order.shipping_address?.country}<br>
                ${order.shipping_address?.phone ? '📞 ' + order.shipping_address.phone : ''}
            </div>
        </div>
        <div>
            <h3 class="section-title">Payment Info</h3>
            <div>
                <strong>Status:</strong> ${order.payment_status?.charAt(0).toUpperCase() + (order.payment_status?.slice(1) || '')}<br>
                <strong>Method:</strong> ${order.payment?.method || 'Stripe'}<br>
                <strong>Currency:</strong> ${order.currency || 'USD'}
            </div>
        </div>
    </div>
    
    <table class="table">
        <thead>
            <tr>
                <th>Item</th>
                <th style="text-align: center;">Quantity</th>
                <th style="text-align: right;">Price</th>
                <th style="text-align: right;">Total</th>
            </tr>
        </thead>
        <tbody>
            ${order.items?.map((item: OrderItem) => `
                <tr>
                    <td>
                        <strong>${item.product?.name || item.product_name}</strong>
                    </td>
                    <td style="text-align: center;">${item.quantity}</td>
                    <td style="text-align: right;">${formatPrice(item.price)}</td>
                    <td style="text-align: right;"><strong>${formatPrice(item.price * item.quantity)}</strong></td>
                </tr>
            `).join('')}
        </tbody>
        <tfoot>
            <tr>
                <td colspan="3"><strong>Subtotal</strong></td>
                <td style="text-align: right;"><strong>${formatPrice(order.subtotal || 0)}</strong></td>
            </tr>
            <tr>
                <td colspan="3"><strong>Shipping</strong></td>
                <td style="text-align: right;"><strong>${order.shipping_amount === 0 ? 'FREE' : formatPrice(order.shipping_amount || 0)}</strong></td>
            </tr>
            <tr>
                <td colspan="3"><strong>Tax</strong></td>
                <td style="text-align: right;"><strong>${(order.tax_amount || 0) === 0 ? '$0' : formatPrice(order.tax_amount)}</strong></td>
            </tr>
            ${order.discount_amount && order.discount_amount > 0 ? `
            <tr>
                <td colspan="3"><strong>Discount</strong></td>
                <td style="text-align: right; color: #28a745;"><strong>-${formatPrice(order.discount_amount)}</strong></td>
            </tr>
            ` : ''}
            <tr class="grand-total">
                <td colspan="3"><strong>TOTAL</strong></td>
                <td style="text-align: right;"><strong>${formatPrice(order.total || 0)}</strong></td>
            </tr>
        </tfoot>
    </table>
    
    <div class="footer">
        <p>Thank you for your business! If you have any questions about this invoice, please contact our support team.</p>
        <p>BiHub - Your trusted e-commerce partner</p>
    </div>
</body>
</html>
    `
  }

  return (
    <Dialog open={isOpen} onOpenChange={onClose}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto bg-gradient-to-br from-slate-900/95 via-gray-900/98 to-slate-800/95 border-gray-700/50 text-white">
        <DialogHeader className="border-b border-gray-700/50 pb-4">
          <DialogTitle className="flex items-center gap-3 text-xl font-semibold">
            <div className="w-10 h-10 bg-gradient-to-br from-[#ff9000]/20 to-orange-600/30 rounded-lg flex items-center justify-center">
              <FileText className="h-5 w-5 text-[#ff9000]" />
            </div>
            Invoice Preview
          </DialogTitle>
        </DialogHeader>

        {/* Invoice Preview Content */}
        <div className="space-y-6">
          {/* Invoice Header */}
          <Card className="bg-gradient-to-br from-white/[0.02] via-white/[0.03] to-white/[0.02] border-gray-700/40">
            <CardContent className="p-6">
              <div className="flex justify-between items-start mb-6">
                <div>
                  <div className="flex items-center gap-3 mb-2">
                    <span className="text-2xl font-bold flex items-center">
                      <span className="text-white">Bi</span>
                      <span className="ml-0.5 px-0.5 py-0 rounded-[3px] text-black font-bold" style={{letterSpacing: '0.3px', backgroundColor: '#FF9000', lineHeight: '1.1'}}>hub</span>
                    </span>
                  </div>
                  <p className="text-gray-400 text-sm">Your trusted online marketplace</p>
                </div>
                <div className="text-right">
                  <h3 className="text-xl font-bold text-white mb-2">INVOICE</h3>
                  <div className="space-y-1 text-sm">
                    <div className="flex items-center gap-2">
                      <Hash className="h-4 w-4 text-gray-400" />
                      <span className="text-gray-400">Number:</span>
                      <span className="text-white font-semibold">{order.order_number}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <Calendar className="h-4 w-4 text-gray-400" />
                      <span className="text-gray-400">Date:</span>
                      <span className="text-white">{formatDate(new Date(order.created_at))}</span>
                    </div>
                  </div>
                </div>
              </div>

              <Separator className="bg-gray-700/50 my-6" />

              {/* Billing Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <MapPin className="h-4 w-4 text-blue-300" />
                    <h4 className="font-semibold text-white">Bill To</h4>
                  </div>
                  {order.shipping_address ? (
                    <div className="space-y-1 text-sm text-gray-300 pl-6">
                      <p className="font-medium text-white">
                        {order.shipping_address.first_name} {order.shipping_address.last_name}
                      </p>
                      {order.shipping_address.company && <p>{order.shipping_address.company}</p>}
                      <p>{order.shipping_address.address1}</p>
                      {order.shipping_address.address2 && <p>{order.shipping_address.address2}</p>}
                      <p>{order.shipping_address.city}, {order.shipping_address.state} {order.shipping_address.zip_code}</p>
                      <p>{order.shipping_address.country}</p>
                      {order.shipping_address.phone && <p>📞 {order.shipping_address.phone}</p>}
                    </div>
                  ) : (
                    <p className="text-gray-400 text-sm pl-6">No billing address available</p>
                  )}
                </div>

                <div>
                  <div className="flex items-center gap-2 mb-3">
                    <CreditCard className="h-4 w-4 text-green-300" />
                    <h4 className="font-semibold text-white">Payment Info</h4>
                  </div>
                  <div className="space-y-2 text-sm pl-6">
                    <div className="flex justify-between">
                      <span className="text-gray-400">Status:</span>
                      <span className={cn(
                        "font-semibold",
                        order.payment_status === 'paid' ? 'text-green-400' : 'text-yellow-400'
                      )}>
                        {order.payment_status?.charAt(0).toUpperCase() + (order.payment_status?.slice(1) || '')}
                      </span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">Method:</span>
                      <span className="text-white">{order.payment?.method || 'Stripe'}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-400">Currency:</span>
                      <span className="text-white">{order.currency || 'USD'}</span>
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Items Table */}
          <Card className="bg-gradient-to-br from-white/[0.02] via-white/[0.03] to-white/[0.02] border-gray-700/40">
            <CardContent className="p-6">
              <h4 className="font-semibold text-white mb-4">Order Items</h4>
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-gray-700/50">
                      <th className="text-left py-3 text-gray-400 font-medium">Item</th>
                      <th className="text-center py-3 text-gray-400 font-medium">Qty</th>
                      <th className="text-right py-3 text-gray-400 font-medium">Price</th>
                      <th className="text-right py-3 text-gray-400 font-medium">Total</th>
                    </tr>
                  </thead>
                  <tbody>
                    {order.items?.map((item: OrderItem, index: number) => (
                      <tr key={item.id} className="border-b border-gray-700/30">
                        <td className="py-4">
                          <div className="space-y-1">
                            <p className="text-white font-medium">
                              {item.product?.name || item.product_name}
                            </p>
                          </div>
                        </td>
                        <td className="py-4 text-center text-white">{item.quantity}</td>
                        <td className="py-4 text-right text-white">{formatPrice(item.price)}</td>
                        <td className="py-4 text-right text-white font-semibold">
                          {formatPrice(item.price * item.quantity)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <Separator className="bg-gray-700/50 my-6" />

              {/* Totals */}
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Subtotal</span>
                  <span className="text-white font-semibold">{formatPrice(order.subtotal || 0)}</span>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Shipping</span>
                  <span className="text-white font-semibold">
                    {(order.shipping_amount || 0) === 0 ? (
                      <span className="text-green-400 font-bold">FREE</span>
                    ) : (
                      formatPrice(order.shipping_amount || 0)
                    )}
                  </span>
                </div>
                
                <div className="flex justify-between items-center">
                  <span className="text-gray-400">Tax</span>
                  <span className="text-white font-semibold">
                    {(order.tax_amount || 0) === 0 ? (
                      <span className="text-gray-400">$0</span>
                    ) : (
                      formatPrice(order.tax_amount)
                    )}
                  </span>
                </div>
                
                {order.discount_amount && order.discount_amount > 0 && (
                  <div className="flex justify-between items-center">
                    <span className="text-gray-400">Discount</span>
                    <span className="text-green-400 font-semibold">-{formatPrice(order.discount_amount)}</span>
                  </div>
                )}
                
                <Separator className="bg-gray-700/50 my-3" />
                
                <div className="flex justify-between items-center py-3 bg-[#ff9000]/10 border border-[#ff9000]/30 rounded-lg px-4">
                  <span className="text-white font-bold text-lg">TOTAL</span>
                  <span className="text-[#ff9000] font-bold text-xl">{formatPrice(order.total || 0)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Footer Note */}
          <Card className="bg-gradient-to-br from-blue-500/5 via-blue-500/10 to-blue-600/5 border-blue-500/30">
            <CardContent className="p-4">
              <p className="text-blue-300 text-sm text-center">
                Thank you for your business! If you have any questions about this invoice, please contact our support team.
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Action Buttons */}
        <div className="flex justify-end gap-3 pt-4 border-t border-gray-700/50">
          <Button
            variant="outline"
            onClick={onClose}
            className="border-gray-600/50 text-gray-300 hover:bg-gray-800/50 hover:border-gray-500 hover:text-white"
          >
            <X className="h-4 w-4 mr-2" />
            Cancel
          </Button>
          <Button
            onClick={handleDownload}
            disabled={isDownloading}
            className="bg-gradient-to-r from-[#ff9000] to-orange-600 hover:from-orange-600 hover:to-orange-700 text-white border-0 shadow-lg disabled:opacity-50"
          >
            <Download className="h-4 w-4 mr-2" />
            {isDownloading ? 'Generating PDF...' : 'Download PDF'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  )
}
