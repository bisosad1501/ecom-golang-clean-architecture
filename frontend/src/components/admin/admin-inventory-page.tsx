'use client';

import React, { useState, useEffect } from 'react';
import {
  Warehouse,
  Search,
  MoreVertical,
  Edit,
  AlertTriangle,
  TrendingUp,
  Package,
  Plus,
  Minus,
  XCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';
import { inventoryService, type InventoryItem, type AdjustStockRequest } from '@/lib/services/inventory';

export function AdminInventoryPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [inventory, setInventory] = useState<InventoryItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [adjustingItem, setAdjustingItem] = useState<InventoryItem | null>(null);
  const [adjustmentType, setAdjustmentType] = useState<'add' | 'remove'>('add');
  const [adjustmentQuantity, setAdjustmentQuantity] = useState('');
  const [adjustmentReason, setAdjustmentReason] = useState('');
  const [totalCount, setTotalCount] = useState(0);
  const [lowStockCount, setLowStockCount] = useState(0);
  const [outOfStockCount, setOutOfStockCount] = useState(0);
  const [totalValue, setTotalValue] = useState(0);

  // Load inventory data from API
  const loadInventory = async () => {
    try {
      setIsLoading(true);
      const response = await inventoryService.getInventories({
        page: 1,
        limit: 100,
        search: searchTerm || undefined,
        status: statusFilter !== 'all' ? statusFilter : undefined,
      });

      setInventory(response.inventories || []);
      setTotalCount(response.total || 0);

      // Calculate stats
      const items = response.inventories || [];
      setLowStockCount(items.filter(item => item.quantity_available <= item.reorder_level && item.quantity_available > 0).length);
      setOutOfStockCount(items.filter(item => item.quantity_available <= 0).length);
      setTotalValue(items.reduce((sum, item) => sum + (item.quantity_on_hand * item.last_cost), 0));

    } catch (error) {
      console.error('Failed to load inventory:', error);
      toast.error('Failed to load inventory data');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadInventory();
  }, [searchTerm, statusFilter]);

  // Remove mock data section completely
  const filteredInventory = inventory.filter(item => {
    const matchesSearch =
      (item.product?.name || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.product?.sku || '').toLowerCase().includes(searchTerm.toLowerCase()) ||
      (item.warehouse?.name || '').toLowerCase().includes(searchTerm.toLowerCase());

    let matchesStatus = true;
    if (statusFilter !== 'all') {
      switch (statusFilter) {
        case 'in_stock':
          matchesStatus = item.quantity_available > item.reorder_level;
          break;
        case 'low_stock':
          matchesStatus = item.quantity_available <= item.reorder_level && item.quantity_available > 0;
          break;
        case 'out_of_stock':
          matchesStatus = item.quantity_available <= 0;
          break;
        case 'overstocked':
          matchesStatus = item.quantity_on_hand > item.max_stock_level;
          break;
      }
    }

    return matchesSearch && matchesStatus;
  });

  const getStatusBadge = (item: InventoryItem) => {
    if (item.quantity_available <= 0) {
      return <Badge className="bg-red-900/30 text-red-400 border-red-700">Out of Stock</Badge>;
    } else if (item.quantity_available <= item.reorder_level) {
      return <Badge className="bg-yellow-900/30 text-yellow-400 border-yellow-700"><AlertTriangle className="w-3 h-3 mr-1" />Low Stock</Badge>;
    } else if (item.quantity_on_hand > item.max_stock_level) {
      return <Badge className="bg-blue-900/30 text-blue-400 border-blue-700">Overstocked</Badge>;
    } else {
      return <Badge className="bg-green-900/30 text-green-400 border-green-700">In Stock</Badge>;
    }
  };

  const handleStockAdjustment = async () => {
    if (!adjustingItem || !adjustmentQuantity || !adjustmentReason) {
      toast.error('Please fill in all fields');
      return;
    }

    try {
      const quantity = parseInt(adjustmentQuantity);
      const quantityDelta = adjustmentType === 'add' ? quantity : -quantity;

      if (adjustingItem.quantity_on_hand + quantityDelta < 0) {
        toast.error('Stock cannot be negative');
        return;
      }

      const adjustRequest: AdjustStockRequest = {
        product_id: adjustingItem.product_id,
        warehouse_id: adjustingItem.warehouse_id,
        quantity_delta: quantityDelta,
        reason: adjustmentReason,
        notes: `${adjustmentType === 'add' ? 'Added' : 'Removed'} ${quantity} units via admin panel`
      };

      await inventoryService.adjustStock(adjustRequest);

      toast.success(`Stock ${adjustmentType === 'add' ? 'added' : 'removed'} successfully`);
      setAdjustingItem(null);
      setAdjustmentQuantity('');
      setAdjustmentReason('');

      // Reload inventory data
      await loadInventory();
    } catch (error) {
      console.error('Failed to adjust stock:', error);
      toast.error('Failed to adjust stock');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Inventory Management</h1>
          <p className="text-gray-300">Monitor and manage product inventory levels</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-sm border-[#FF9000] text-[#FF9000]">
            {totalCount} Products
          </Badge>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-xl flex items-center justify-center">
                <Package className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Total Items</p>
                <p className="text-2xl font-bold text-white">
                  {inventory.reduce((sum, item) => sum + item.quantity_on_hand, 0)}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center">
                <TrendingUp className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Total Value</p>
                <p className="text-2xl font-bold text-white">
                  ${totalValue.toLocaleString()}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
        
        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-xl flex items-center justify-center">
                <AlertTriangle className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Low Stock</p>
                <p className="text-2xl font-bold text-white">{lowStockCount}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-red-500 to-red-600 rounded-xl flex items-center justify-center">
                <XCircle className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Out of Stock</p>
                <p className="text-2xl font-bold text-white">{outOfStockCount}</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
        <CardHeader>
          <CardTitle className="text-lg text-white">Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  placeholder="Search by product name, SKU, or warehouse..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />
              </div>
            </div>

            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-48 bg-gray-700 border-gray-600 text-white">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent className="bg-gray-800 border-gray-600">
                <SelectItem value="all" className="text-white hover:bg-gray-700">All Status</SelectItem>
                <SelectItem value="in_stock" className="text-white hover:bg-gray-700">In Stock</SelectItem>
                <SelectItem value="low_stock" className="text-white hover:bg-gray-700">Low Stock</SelectItem>
                <SelectItem value="out_of_stock" className="text-white hover:bg-gray-700">Out of Stock</SelectItem>
                <SelectItem value="overstocked" className="text-white hover:bg-gray-700">Overstocked</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Inventory List */}
      <div className="space-y-4">
        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 5 }, (_, i) => (
              <Card key={i} className="animate-pulse bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
                <CardContent className="p-6">
                  <div className="h-32 bg-gray-700 rounded"></div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : filteredInventory.length === 0 ? (
          <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
            <CardContent className="text-center py-12">
              <Warehouse className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-white mb-2">No inventory found</h3>
              <p className="text-gray-400">
                {searchTerm || statusFilter !== 'all'
                  ? 'Try adjusting your filters'
                  : 'No inventory items have been added yet'
                }
              </p>
            </CardContent>
          </Card>
        ) : (
          filteredInventory.map((item) => (
            <Card key={item.id} className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
              <CardContent className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-start gap-4">
                    <div className="w-16 h-16 bg-gray-700 rounded-xl flex items-center justify-center overflow-hidden">
                      {item.product?.images?.[0]?.url ? (
                        <img src={item.product.images[0].url} alt={item.product.name} className="w-full h-full object-cover" />
                      ) : (
                        <Package className="w-8 h-8 text-gray-400" />
                      )}
                    </div>

                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-2">
                        <h4 className="font-semibold text-white">{item.product?.name || 'Unknown Product'}</h4>
                        {getStatusBadge(item)}
                      </div>

                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm text-gray-300">
                        <div>
                          <p><strong className="text-gray-200">SKU:</strong> {item.product?.sku || 'N/A'}</p>
                          <p><strong className="text-gray-200">Warehouse:</strong> {item.warehouse?.name || 'N/A'}</p>
                          <p><strong className="text-gray-200">Cost:</strong> ${item.last_cost?.toFixed(2) || '0.00'}</p>
                        </div>
                        <div>
                          <p><strong className="text-gray-200">On Hand:</strong> {item.quantity_on_hand}</p>
                          <p><strong className="text-gray-200">Available:</strong> {item.quantity_available}</p>
                          <p><strong className="text-gray-200">Reserved:</strong> {item.quantity_reserved}</p>
                        </div>
                        <div>
                          <p><strong className="text-gray-200">Reorder Level:</strong> {item.reorder_level}</p>
                          <p><strong className="text-gray-200">Max Stock:</strong> {item.max_stock_level}</p>
                          <p><strong className="text-gray-200">Updated:</strong> {formatDate(item.updated_at)}</p>
                        </div>
                      </div>
                    </div>
                  </div>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="text-gray-300 hover:text-white hover:bg-gray-700">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="bg-gray-800 border-gray-600">
                      <DropdownMenuItem onClick={() => setAdjustingItem(item)} className="text-white hover:bg-gray-700">
                        <Edit className="w-4 h-4 mr-2" />
                        Adjust Stock
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                {/* Stock Level Bar */}
                <div className="w-full bg-gray-600 rounded-full h-2 mb-2">
                  <div
                    className={cn(
                      "h-2 rounded-full transition-all duration-300",
                      item.quantity_available <= 0 ? 'bg-red-500' :
                      item.quantity_available <= item.reorder_level ? 'bg-yellow-500' :
                      item.quantity_on_hand > item.max_stock_level ? 'bg-blue-500' :
                      'bg-green-500'
                    )}
                    style={{
                      width: `${Math.min((item.quantity_on_hand / item.max_stock_level) * 100, 100)}%`
                    }}
                  />
                </div>
                <div className="flex justify-between text-xs text-gray-400">
                  <span>0</span>
                  <span>Reorder: {item.reorder_level}</span>
                  <span>Max: {item.max_stock_level}</span>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {/* Stock Adjustment Dialog */}
      <Dialog open={!!adjustingItem} onOpenChange={() => setAdjustingItem(null)}>
        <DialogContent className="max-w-md bg-gray-800 border-gray-600">
          <DialogHeader>
            <DialogTitle className="text-white">Adjust Stock</DialogTitle>
          </DialogHeader>

          {adjustingItem && (
            <div className="space-y-4">
              <div className="bg-gray-700 p-4 rounded-lg">
                <h4 className="font-medium text-white">{adjustingItem.product?.name || 'Unknown Product'}</h4>
                <p className="text-sm text-gray-300">Current Stock: {adjustingItem.quantity_on_hand}</p>
                <p className="text-sm text-gray-300">SKU: {adjustingItem.product?.sku || 'N/A'}</p>
              </div>

              <div>
                <Label className="text-gray-200">Adjustment Type</Label>
                <Select value={adjustmentType} onValueChange={(value: 'add' | 'remove') => setAdjustmentType(value)}>
                  <SelectTrigger className="bg-gray-700 border-gray-600 text-white">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-800 border-gray-600">
                    <SelectItem value="add" className="text-white hover:bg-gray-700">
                      <div className="flex items-center">
                        <Plus className="w-4 h-4 mr-2 text-green-400" />
                        Add Stock
                      </div>
                    </SelectItem>
                    <SelectItem value="remove" className="text-white hover:bg-gray-700">
                      <div className="flex items-center">
                        <Minus className="w-4 h-4 mr-2 text-red-400" />
                        Remove Stock
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="quantity" className="text-gray-200">Quantity</Label>
                <Input
                  id="quantity"
                  type="number"
                  min="1"
                  value={adjustmentQuantity}
                  onChange={(e) => setAdjustmentQuantity(e.target.value)}
                  placeholder="Enter quantity"
                  className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />
              </div>

              <div>
                <Label htmlFor="reason" className="text-gray-200">Reason</Label>
                <Select value={adjustmentReason} onValueChange={setAdjustmentReason}>
                  <SelectTrigger className="bg-gray-700 border-gray-600 text-white">
                    <SelectValue placeholder="Select reason" />
                  </SelectTrigger>
                  <SelectContent className="bg-gray-800 border-gray-600">
                    <SelectItem value="restock" className="text-white hover:bg-gray-700">Restock</SelectItem>
                    <SelectItem value="damaged" className="text-white hover:bg-gray-700">Damaged</SelectItem>
                    <SelectItem value="lost" className="text-white hover:bg-gray-700">Lost</SelectItem>
                    <SelectItem value="returned" className="text-white hover:bg-gray-700">Returned</SelectItem>
                    <SelectItem value="correction" className="text-white hover:bg-gray-700">Inventory Correction</SelectItem>
                    <SelectItem value="other" className="text-white hover:bg-gray-700">Other</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setAdjustingItem(null)} className="border-gray-600 text-gray-300 hover:bg-gray-700">
                  Cancel
                </Button>
                <Button
                  onClick={handleStockAdjustment}
                  disabled={!adjustmentQuantity || !adjustmentReason}
                  className="bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600] text-white"
                >
                  Adjust Stock
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
