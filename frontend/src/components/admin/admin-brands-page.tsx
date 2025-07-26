'use client';

import React, { useState, useEffect } from 'react';
import {
  Plus,
  Search,
  MoreVertical,
  Edit,
  Trash2,
  Tag,
  Image as ImageIcon
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';

// BiHub Components
import {
  BiHubAdminCard,
  BiHubPageHeader,
  BiHubEmptyState,
  BiHubStatCard,
} from './bihub-admin-components';
import { BIHUB_ADMIN_THEME } from '@/constants/admin-theme';
import { RequirePermission } from '@/components/auth/permission-guard';
import { PERMISSIONS } from '@/lib/permissions';

interface Brand {
  id: string;
  name: string;
  description?: string;
  logo?: string;
  website?: string;
  is_active: boolean;
  product_count: number;
  created_at: string;
  updated_at: string;
}

export function AdminBrandsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [brands, setBrands] = useState<Brand[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingBrand, setEditingBrand] = useState<Brand | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    website: '',
    logo: ''
  });

  // Mock data for demonstration
  useEffect(() => {
    const mockBrands: Brand[] = [
      {
        id: '1',
        name: 'Apple',
        description: 'Technology company known for innovative products',
        logo: '/api/placeholder/60/60',
        website: 'https://apple.com',
        is_active: true,
        product_count: 25,
        created_at: '2024-01-01T00:00:00Z',
        updated_at: '2024-01-01T00:00:00Z'
      },
      {
        id: '2',
        name: 'Samsung',
        description: 'Electronics and technology conglomerate',
        logo: '/api/placeholder/60/60',
        website: 'https://samsung.com',
        is_active: true,
        product_count: 18,
        created_at: '2024-01-02T00:00:00Z',
        updated_at: '2024-01-02T00:00:00Z'
      },
      {
        id: '3',
        name: 'Nike',
        description: 'Athletic footwear and apparel',
        logo: '/api/placeholder/60/60',
        website: 'https://nike.com',
        is_active: true,
        product_count: 42,
        created_at: '2024-01-03T00:00:00Z',
        updated_at: '2024-01-03T00:00:00Z'
      }
    ];
    
    setTimeout(() => {
      setBrands(mockBrands);
      setIsLoading(false);
    }, 1000);
  }, []);

  const filteredBrands = brands.filter(brand =>
    brand.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    brand.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleCreateBrand = async () => {
    try {
      // TODO: Implement API call
      toast.success('Brand created successfully');
      setIsCreateDialogOpen(false);
      setFormData({ name: '', description: '', website: '', logo: '' });
    } catch (error) {
      toast.error('Failed to create brand');
    }
  };

  const handleUpdateBrand = async () => {
    try {
      // TODO: Implement API call
      toast.success('Brand updated successfully');
      setEditingBrand(null);
      setFormData({ name: '', description: '', website: '', logo: '' });
    } catch (error) {
      toast.error('Failed to update brand');
    }
  };

  const handleDeleteBrand = async (brandId: string) => {
    try {
      // TODO: Implement API call
      setBrands(brands.filter(b => b.id !== brandId));
      toast.success('Brand deleted successfully');
    } catch (error) {
      toast.error('Failed to delete brand');
    }
  };

  const openEditDialog = (brand: Brand) => {
    setEditingBrand(brand);
    setFormData({
      name: brand.name,
      description: brand.description || '',
      website: brand.website || '',
      logo: brand.logo || ''
    });
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className={BIHUB_ADMIN_THEME.spacing.section}>
      {/* BiHub Page Header */}
      <BiHubPageHeader
        title="Brand Management"
        subtitle="Manage product brands and manufacturers in your BiHub catalog"
        breadcrumbs={[
          { label: 'Admin' },
          { label: 'Brands' }
        ]}
        action={
          <RequirePermission permission={PERMISSIONS.SYSTEM_SETTINGS}>
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
              <DialogTrigger asChild>
                <Button className={BIHUB_ADMIN_THEME.components.button.primary}>
                  <Plus className="w-4 h-4 mr-2" />
                  Add Brand
                </Button>
              </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>Create New Brand</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <div>
                <Label htmlFor="name">Brand Name</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="Enter brand name"
                />
              </div>
              <div>
                <Label htmlFor="description">Description</Label>
                <Textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  placeholder="Enter brand description"
                  rows={3}
                />
              </div>
              <div>
                <Label htmlFor="website">Website</Label>
                <Input
                  id="website"
                  value={formData.website}
                  onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                  placeholder="https://example.com"
                />
              </div>
              <div>
                <Label htmlFor="logo">Logo URL</Label>
                <Input
                  id="logo"
                  value={formData.logo}
                  onChange={(e) => setFormData({ ...formData, logo: e.target.value })}
                  placeholder="https://example.com/logo.png"
                />
              </div>
              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                  Cancel
                </Button>
                <Button onClick={handleCreateBrand} disabled={!formData.name.trim()}>
                  Create Brand
                </Button>
              </div>
            </div>
          </DialogContent>
        </Dialog>
        </RequirePermission>
      }
      />

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <BiHubStatCard
          title="Total Brands"
          value={brands.length.toString()}
          icon={<Tag className="h-5 w-5 text-white" />}
          trend={{ value: 12, isPositive: true }}
          color="blue"
        />

        <BiHubStatCard
          title="Active Brands"
          value={brands.filter(b => b.is_active).length.toString()}
          icon={<Tag className="h-5 w-5 text-white" />}
          trend={{ value: 8, isPositive: true }}
          color="green"
        />

        <BiHubStatCard
          title="Total Products"
          value={brands.reduce((sum, b) => sum + b.product_count, 0).toString()}
          icon={<Tag className="h-5 w-5 text-white" />}
          trend={{ value: 24, isPositive: true }}
          color="orange"
        />
      </div>

      {/* Brands Management */}
      <BiHubAdminCard
        title="Brand Catalog"
        subtitle="Manage your product brands and manufacturers"
        icon={<Tag className="h-5 w-5 text-white" />}
        headerAction={
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
            <Input
              placeholder="Search brands..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-10 w-64"
            />
          </div>
        }
      >
        <div className="space-y-4">
          {isLoading ? (
            <div className="space-y-4">
              {Array.from({ length: 3 }, (_, i) => (
                <div key={i} className="animate-pulse">
                  <div className="h-20 bg-gray-700 rounded-lg"></div>
                </div>
              ))}
            </div>
          ) : filteredBrands.length === 0 ? (
            <BiHubEmptyState
              icon={<Tag className="h-12 w-12 text-gray-400" />}
              title="No brands found"
              description={searchTerm ? 'Try adjusting your search' : 'Create your first brand to get started'}
            />
          ) : (
          filteredBrands.map((brand) => (
            <div key={brand.id} className="p-4 bg-gray-800 rounded-lg border border-gray-700 hover:border-gray-600 transition-colors">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <div className="w-16 h-16 bg-gray-700 rounded-xl flex items-center justify-center overflow-hidden">
                    {brand.logo ? (
                      <img src={brand.logo} alt={brand.name} className="w-full h-full object-cover" />
                    ) : (
                      <ImageIcon className="w-8 h-8 text-gray-400" />
                    )}
                  </div>

                  <div className="flex-1">
                    <div className="flex items-center gap-2 mb-1">
                      <h3 className="text-lg font-semibold text-white">{brand.name}</h3>
                      <Badge variant={brand.is_active ? "default" : "secondary"}>
                        {brand.is_active ? 'Active' : 'Inactive'}
                      </Badge>
                    </div>

                    {brand.description && (
                      <p className="text-gray-300 mb-2">{brand.description}</p>
                    )}

                    <div className="flex items-center gap-4 text-sm text-gray-400">
                      <span>{brand.product_count} products</span>
                      <span>Created {formatDate(brand.created_at)}</span>
                      {brand.website && (
                        <a
                          href={brand.website}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-[#FF9000] hover:underline"
                        >
                          Visit Website
                        </a>
                      )}
                    </div>
                  </div>
                </div>

                <RequirePermission permission={PERMISSIONS.SYSTEM_SETTINGS}>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => openEditDialog(brand)}>
                        <Edit className="w-4 h-4 mr-2" />
                        Edit
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => handleDeleteBrand(brand.id)}
                        className="text-red-400"
                      >
                        <Trash2 className="w-4 h-4 mr-2" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </RequirePermission>
              </div>
            </div>
          ))
        )}
        </div>
      </BiHubAdminCard>

      {/* Edit Dialog */}
      <Dialog open={!!editingBrand} onOpenChange={() => setEditingBrand(null)}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Edit Brand</DialogTitle>
          </DialogHeader>
          <div className="space-y-4">
            <div>
              <Label htmlFor="edit-name">Brand Name</Label>
              <Input
                id="edit-name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="Enter brand name"
              />
            </div>
            <div>
              <Label htmlFor="edit-description">Description</Label>
              <Textarea
                id="edit-description"
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="Enter brand description"
                rows={3}
              />
            </div>
            <div>
              <Label htmlFor="edit-website">Website</Label>
              <Input
                id="edit-website"
                value={formData.website}
                onChange={(e) => setFormData({ ...formData, website: e.target.value })}
                placeholder="https://example.com"
              />
            </div>
            <div>
              <Label htmlFor="edit-logo">Logo URL</Label>
              <Input
                id="edit-logo"
                value={formData.logo}
                onChange={(e) => setFormData({ ...formData, logo: e.target.value })}
                placeholder="https://example.com/logo.png"
              />
            </div>
            <div className="flex justify-end gap-2">
              <Button variant="outline" onClick={() => setEditingBrand(null)}>
                Cancel
              </Button>
              <Button onClick={handleUpdateBrand} disabled={!formData.name.trim()}>
                Update Brand
              </Button>
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
