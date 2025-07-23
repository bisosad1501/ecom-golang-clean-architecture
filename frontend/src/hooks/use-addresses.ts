import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import apiClient from '@/lib/api'

// Types
export interface Address {
  id: string
  user_id: string
  type: 'shipping' | 'billing' | 'both'
  first_name: string
  last_name: string
  company?: string
  address1: string
  address2?: string
  city: string
  state: string
  zip_code: string
  country: string
  phone?: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CreateAddressRequest {
  type: 'shipping' | 'billing' | 'both'
  first_name: string
  last_name: string
  company?: string
  address1: string
  address2?: string
  city: string
  state: string
  zip_code: string
  country: string
  phone?: string
  is_default?: boolean
}

export interface UpdateAddressRequest extends Partial<CreateAddressRequest> {}

// Query keys
export const addressKeys = {
  all: ['addresses'] as const,
  lists: () => [...addressKeys.all, 'list'] as const,
  list: (filters: Record<string, any>) => [...addressKeys.lists(), filters] as const,
  details: () => [...addressKeys.all, 'detail'] as const,
  detail: (id: string) => [...addressKeys.details(), id] as const,
}

// Get user addresses
export function useAddresses() {
  return useQuery({
    queryKey: addressKeys.lists(),
    queryFn: async (): Promise<Address[]> => {
      const response = await apiClient.get<Address[]>('/users/addresses')
      return response.data
    },
  })
}

// Get single address
export function useAddress(id: string) {
  return useQuery({
    queryKey: addressKeys.detail(id),
    queryFn: async (): Promise<Address> => {
      const response = await apiClient.get<Address>(`/users/addresses/${id}`)
      return response.data
    },
    enabled: !!id,
  })
}

// Create address
export function useCreateAddress() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (data: CreateAddressRequest): Promise<Address> => {
      const response = await apiClient.post<Address>('/users/addresses', data)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: addressKeys.lists() })
      toast.success('Address created successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to create address')
    },
  })
}

// Update address
export function useUpdateAddress() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async ({ id, data }: { id: string; data: UpdateAddressRequest }): Promise<Address> => {
      const response = await apiClient.put<Address>(`/users/addresses/${id}`, data)
      return response.data
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: addressKeys.lists() })
      queryClient.setQueryData(addressKeys.detail(data.id), data)
      toast.success('Address updated successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to update address')
    },
  })
}

// Delete address
export function useDeleteAddress() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string): Promise<void> => {
      await apiClient.delete(`/users/addresses/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: addressKeys.lists() })
      toast.success('Address deleted successfully')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to delete address')
    },
  })
}

// Set default address
export function useSetDefaultAddress() {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: async (id: string): Promise<Address> => {
      const response = await apiClient.put<Address>(`/users/addresses/${id}/default`)
      return response.data
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: addressKeys.lists() })
      toast.success('Default address updated')
    },
    onError: (error: any) => {
      toast.error(error.message || 'Failed to set default address')
    },
  })
}
