'use client';

import React, { useState, useEffect } from 'react';
import {
  Truck,
  Search,
  MoreVertical,
  Eye,
  MapPin,
  Download,
  Clock,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';

// BiHub Components
import {
  BiHubAdminCard,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
  BiHubStatusBadge,
} from './bihub-admin-components';
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme';
import { RequirePermission } from '@/components/auth/permission-guard';
import { PERMISSIONS } from '@/lib/permissions';

interface Shipment {
  id: string;
  order_id: string;
  tracking_number: string;
  carrier: string;
  status: 'pending' | 'shipped' | 'in_transit' | 'delivered' | 'failed';
  shipping_address: {
    street: string;
    city: string;
    state: string;
    zip_code: string;
    country: string;
  };
  estimated_delivery: string;
  actual_delivery?: string;
  created_at: string;
  updated_at: string;
  customer_name: string;
  customer_email: string;
  items_count: number;
  total_weight: number;
}

export function AdminShipmentsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [carrierFilter, setCarrierFilter] = useState<string>('all');
  const [shipments, setShipments] = useState<Shipment[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Mock data for demonstration
  useEffect(() => {
    const mockShipments: Shipment[] = [
      {
        id: '1',
        order_id: 'ORD-001',
        tracking_number: 'TRK123456789',
        carrier: 'FedEx',
        status: 'in_transit',
        shipping_address: {
          street: '123 Main St',
          city: 'New York',
          state: 'NY',
          zip_code: '10001',
          country: 'USA'
        },
        estimated_delivery: '2024-01-20T18:00:00Z',
        created_at: '2024-01-15T10:30:00Z',
        updated_at: '2024-01-18T14:20:00Z',
        customer_name: 'John Doe',
        customer_email: 'john@example.com',
        items_count: 3,
        total_weight: 2.5
      },
      {
        id: '2',
        order_id: 'ORD-002',
        tracking_number: 'TRK987654321',
        carrier: 'UPS',
        status: 'delivered',
        shipping_address: {
          street: '456 Oak Ave',
          city: 'Los Angeles',
          state: 'CA',
          zip_code: '90210',
          country: 'USA'
        },
        estimated_delivery: '2024-01-18T17:00:00Z',
        actual_delivery: '2024-01-18T16:45:00Z',
        created_at: '2024-01-14T09:15:00Z',
        updated_at: '2024-01-18T16:45:00Z',
        customer_name: 'Jane Smith',
        customer_email: 'jane@example.com',
        items_count: 1,
        total_weight: 1.2
      },
      {
        id: '3',
        order_id: 'ORD-003',
        tracking_number: 'TRK555666777',
        carrier: 'DHL',
        status: 'pending',
        shipping_address: {
          street: '789 Pine St',
          city: 'Chicago',
          state: 'IL',
          zip_code: '60601',
          country: 'USA'
        },
        estimated_delivery: '2024-01-22T16:00:00Z',
        created_at: '2024-01-19T11:00:00Z',
        updated_at: '2024-01-19T11:00:00Z',
        customer_name: 'Mike Johnson',
        customer_email: 'mike@example.com',
        items_count: 2,
        total_weight: 3.1
      }
    ];
    
    setTimeout(() => {
      setShipments(mockShipments);
      setIsLoading(false);
    }, 1000);
  }, []);

  const filteredShipments = shipments.filter(shipment => {
    const matchesSearch = 
      shipment.tracking_number.toLowerCase().includes(searchTerm.toLowerCase()) ||
      shipment.order_id.toLowerCase().includes(searchTerm.toLowerCase()) ||
      shipment.customer_name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      shipment.customer_email.toLowerCase().includes(searchTerm.toLowerCase());
    
    const matchesStatus = statusFilter === 'all' || shipment.status === statusFilter;
    const matchesCarrier = carrierFilter === 'all' || shipment.carrier === carrierFilter;
    
    return matchesSearch && matchesStatus && matchesCarrier;
  });



  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatAddress = (address: Shipment['shipping_address']) => {
    return `${address.street}, ${address.city}, ${address.state} ${address.zip_code}, ${address.country}`;
  };

  const carriers = ['FedEx', 'UPS', 'DHL', 'USPS'];

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Shipment Management"
        subtitle="Track and manage order shipments across all carriers"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Shipments' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
            <Button className={BIHUB_ADMIN_THEME.components.button.secondary}>
              <Download className="mr-2 h-5 w-5" />
              Export Shipments
            </Button>
          </RequirePermission>
        }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <BiHubStatCard
          title="Pending"
          value={shipments.filter(s => s.status === 'pending').length.toString()}
          icon={<Clock className="h-5 w-5 text-white" />}
          change={-5}
          color="warning"
        />

        <BiHubStatCard
          title="In Transit"
          value={shipments.filter(s => s.status === 'in_transit' || s.status === 'shipped').length.toString()}
          icon={<Truck className="h-5 w-5 text-white" />}
          change={12}
          color="info"
        />

        <BiHubStatCard
          title="Delivered"
          value={shipments.filter(s => s.status === 'delivered').length.toString()}
          icon={<CheckCircle className="h-5 w-5 text-white" />}
          change={8}
          color="success"
        />

        <BiHubStatCard
          title="Failed"
          value={shipments.filter(s => s.status === 'failed').length.toString()}
          icon={<AlertCircle className="h-5 w-5 text-white" />}
          change={-2}
          color="error"
        />
      </div>

      {/* Shipment Management */}
      <BiHubAdminCard
        title="Shipment Tracking"
        subtitle="Monitor and manage all shipments across carriers"
        icon={<Truck className="h-5 w-5 text-white" />}
        headerAction={
          <div className="flex flex-col md:flex-row gap-4">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
              <Input
                placeholder="Search shipments..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-10 w-64"
              />
            </div>

            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="shipped">Shipped</SelectItem>
                <SelectItem value="in_transit">In Transit</SelectItem>
                <SelectItem value="delivered">Delivered</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>

            <Select value={carrierFilter} onValueChange={setCarrierFilter}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Filter by carrier" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Carriers</SelectItem>
                {carriers.map(carrier => (
                  <SelectItem key={carrier} value={carrier}>{carrier}</SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        }
      >
        <div className="space-y-4">
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 5 }, (_, i) => (
                <div key={i} className="animate-pulse">
                  <div className="h-32 bg-gray-700 rounded-lg"></div>
                </div>
            ))}
          </div>
        ) : filteredShipments.length === 0 ? (
          <BiHubEmptyState
            icon={<Truck className="h-12 w-12 text-gray-400" />}
            title="No shipments found"
            description={
              searchTerm || statusFilter !== 'all' || carrierFilter !== 'all'
                ? 'Try adjusting your filters'
                : 'No shipments have been created yet'
            }
          />
        ) : (
          filteredShipments.map((shipment) => (
            <div key={shipment.id} className="p-6 bg-gray-800 rounded-lg border border-gray-700 hover:border-gray-600 transition-colors">
              <div className="flex items-start justify-between mb-4">
                <div className="flex items-start gap-4">
                  <div className="w-12 h-12 bg-gradient-to-br from-blue-400 to-blue-600 rounded-xl flex items-center justify-center text-white">
                    <Truck className="w-6 h-6" />
                  </div>

                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-2">
                      <h4 className="font-semibold text-white">
                        {shipment.tracking_number}
                      </h4>
                      <BiHubStatusBadge status={shipment.status}>
                        {shipment.status.charAt(0).toUpperCase() + shipment.status.slice(1).replace('_', ' ')}
                      </BiHubStatusBadge>
                      <Badge variant="outline" className="text-xs border-gray-600 text-gray-300">
                        {shipment.carrier}
                      </Badge>
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-gray-300">
                      <div>
                        <p><strong>Order:</strong> {shipment.order_id}</p>
                        <p><strong>Customer:</strong> {shipment.customer_name}</p>
                        <p><strong>Email:</strong> {shipment.customer_email}</p>
                      </div>
                      <div>
                        <p><strong>Items:</strong> {shipment.items_count}</p>
                        <p><strong>Weight:</strong> {shipment.total_weight} kg</p>
                        <p><strong>Created:</strong> {formatDate(shipment.created_at)}</p>
                      </div>
                    </div>
                  </div>
                </div>

                <RequirePermission permission={PERMISSIONS.ORDERS_VIEW_ALL}>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>
                        <Eye className="w-4 h-4 mr-2" />
                        View Details
                      </DropdownMenuItem>
                      <DropdownMenuItem>
                        <Truck className="w-4 h-4 mr-2" />
                        Track Package
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </RequirePermission>
              </div>

              {/* Shipping Address */}
              <div className="bg-gray-700 rounded-lg p-4 mb-4">
                <div className="flex items-start gap-2">
                  <MapPin className="w-4 h-4 text-gray-400 mt-0.5" />
                  <div>
                    <p className="text-sm font-medium text-white">Shipping Address</p>
                    <p className="text-sm text-gray-300">{formatAddress(shipment.shipping_address)}</p>
                  </div>
                </div>
              </div>

              {/* Delivery Info */}
              <div className="flex items-center justify-between text-sm text-gray-300">
                <div>
                  <span className="font-medium">Estimated Delivery: </span>
                  {formatDate(shipment.estimated_delivery)}
                </div>
                {shipment.actual_delivery && (
                  <div>
                    <span className="font-medium">Delivered: </span>
                    {formatDate(shipment.actual_delivery)}
                  </div>
                )}
              </div>
            </div>
          ))
        )}
        </div>
      </BiHubAdminCard>
    </div>
  );
}
