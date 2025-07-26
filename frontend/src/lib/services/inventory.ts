import { apiClient } from '@/lib/api';

export interface InventoryItem {
  id: string;
  product_id: string;
  warehouse_id: string;
  product?: {
    id: string;
    name: string;
    sku: string;
    images?: Array<{ url: string; alt_text?: string }>;
  };
  warehouse?: {
    id: string;
    name: string;
    code: string;
  };
  quantity_on_hand: number;
  quantity_reserved: number;
  quantity_available: number;
  reorder_level: number;
  max_stock_level: number;
  min_stock_level: number;
  average_cost: number;
  last_cost: number;
  last_movement_at?: string;
  last_count_at?: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface InventoryMovement {
  id: string;
  inventory_id: string;
  type: string;
  reason: string;
  quantity: number;
  quantity_before: number;
  quantity_after: number;
  unit_cost?: number;
  total_cost?: number;
  reference_type?: string;
  reference_id?: string;
  batch_number?: string;
  expiry_date?: string;
  notes?: string;
  created_by: string;
  created_at: string;
}

export interface StockAlert {
  id: string;
  inventory_id: string;
  type: string;
  message: string;
  severity: string;
  is_resolved: boolean;
  resolved_at?: string;
  resolved_by?: string;
  created_at: string;
}

export interface GetInventoriesParams {
  page?: number;
  limit?: number;
  search?: string;
  warehouse_id?: string;
  status?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

export interface GetInventoriesResponse {
  inventories: InventoryItem[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface UpdateInventoryRequest {
  product_id: string;
  warehouse_id: string;
  quantity_on_hand?: number;
  reorder_level?: number;
  max_stock_level?: number;
  min_stock_level?: number;
}

export interface AdjustStockRequest {
  product_id: string;
  warehouse_id: string;
  quantity_delta: number;
  reason: string;
  notes?: string;
}

export interface RecordMovementRequest {
  product_id: string;
  warehouse_id: string;
  type: string;
  reason: string;
  quantity: number;
  unit_cost?: number;
  reference_type?: string;
  reference_id?: string;
  batch_number?: string;
  expiry_date?: string;
  notes?: string;
}

export interface TransferStockRequest {
  product_id: string;
  from_warehouse_id: string;
  to_warehouse_id: string;
  quantity: number;
  reason: string;
  notes?: string;
}

export interface GetMovementsParams {
  page?: number;
  limit?: number;
  inventory_id?: string;
  type?: string;
  from_date?: string;
  to_date?: string;
}

export interface GetMovementsResponse {
  movements: InventoryMovement[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface GetAlertsResponse {
  alerts: StockAlert[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export class InventoryService {
  private baseUrl = '/admin/inventory';

  async getInventories(params: GetInventoriesParams = {}): Promise<GetInventoriesResponse> {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.limit) searchParams.append('limit', params.limit.toString());
    if (params.search) searchParams.append('search', params.search);
    if (params.warehouse_id) searchParams.append('warehouse_id', params.warehouse_id);
    if (params.status) searchParams.append('status', params.status);
    if (params.sort_by) searchParams.append('sort_by', params.sort_by);
    if (params.sort_order) searchParams.append('sort_order', params.sort_order);

    const url = `${this.baseUrl}?${searchParams.toString()}`;
    const response = await apiClient.get(url);
    return response.data;
  }

  async getInventory(id: string): Promise<InventoryItem> {
    const response = await apiClient.get(`${this.baseUrl}/${id}`);
    return response.data;
  }

  async updateInventory(id: string, data: UpdateInventoryRequest): Promise<InventoryItem> {
    const response = await apiClient.put(`${this.baseUrl}/${id}`, data);
    return response.data;
  }

  async adjustStock(data: AdjustStockRequest): Promise<InventoryMovement> {
    const response = await apiClient.post(`${this.baseUrl}/adjust`, data);
    return response.data;
  }

  async recordMovement(data: RecordMovementRequest): Promise<InventoryMovement> {
    const response = await apiClient.post(`${this.baseUrl}/movements`, data);
    return response.data;
  }

  async transferStock(data: TransferStockRequest): Promise<void> {
    await apiClient.post(`${this.baseUrl}/transfer`, data);
  }

  async getMovements(params: GetMovementsParams = {}): Promise<GetMovementsResponse> {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.limit) searchParams.append('limit', params.limit.toString());
    if (params.inventory_id) searchParams.append('inventory_id', params.inventory_id);
    if (params.type) searchParams.append('type', params.type);
    if (params.from_date) searchParams.append('from_date', params.from_date);
    if (params.to_date) searchParams.append('to_date', params.to_date);

    const url = `${this.baseUrl}/movements?${searchParams.toString()}`;
    const response = await apiClient.get(url);
    return response.data;
  }

  async getStockAlerts(params: { page?: number; limit?: number } = {}): Promise<GetAlertsResponse> {
    const searchParams = new URLSearchParams();
    
    if (params.page) searchParams.append('page', params.page.toString());
    if (params.limit) searchParams.append('limit', params.limit.toString());

    const url = `${this.baseUrl}/alerts?${searchParams.toString()}`;
    const response = await apiClient.get(url);
    return response.data;
  }

  async resolveAlert(alertId: string): Promise<void> {
    await apiClient.put(`${this.baseUrl}/alerts/${alertId}/resolve`);
  }

  async getLowStockItems(): Promise<InventoryItem[]> {
    const response = await apiClient.get(`${this.baseUrl}/low-stock`);
    return response.data;
  }

  async getOutOfStockItems(): Promise<InventoryItem[]> {
    const response = await apiClient.get(`${this.baseUrl}/out-of-stock`);
    return response.data;
  }
}

export const inventoryService = new InventoryService();
