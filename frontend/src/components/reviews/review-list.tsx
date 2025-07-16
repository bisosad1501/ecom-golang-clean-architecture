'use client';

import React from 'react';
import { Filter, SortAsc, SortDesc, Star } from 'lucide-react';
import { Review, GetReviewsParams } from '@/lib/services/review';
import { ReviewCard } from './review-card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';

interface ReviewListProps {
  reviews: Review[];
  totalCount: number;
  currentParams: GetReviewsParams;
  onParamsChange: (params: GetReviewsParams) => void;
  onVote?: (reviewId: string, isHelpful: boolean) => void;
  onRemoveVote?: (reviewId: string) => void;
  onEdit?: (review: Review) => void;
  onDelete?: (reviewId: string) => void;
  onReply?: (reviewId: string) => void;
  onLoadMore?: () => void;
  isLoading?: boolean;
  hasMore?: boolean;
  currentUserId?: string;
  isAdmin?: boolean;
  className?: string;
}

export function ReviewList({
  reviews,
  totalCount,
  currentParams,
  onParamsChange,
  onVote,
  onRemoveVote,
  onEdit,
  onDelete,
  onReply,
  onLoadMore,
  isLoading = false,
  hasMore = false,
  currentUserId,
  isAdmin = false,
  className,
}: ReviewListProps) {
  const handleSortChange = (sortBy: string) => {
    const newSortOrder = 
      currentParams.sort_by === sortBy && currentParams.sort_order === 'desc' 
        ? 'asc' 
        : 'desc';
    
    onParamsChange({
      ...currentParams,
      sort_by: sortBy as any,
      sort_order: newSortOrder,
    });
  };

  const handleRatingFilter = (rating: number | null) => {
    onParamsChange({
      ...currentParams,
      rating: rating || undefined,
      offset: 0, // Reset to first page
    });
  };

  const handleVerifiedFilter = (verified: boolean | null) => {
    onParamsChange({
      ...currentParams,
      verified: verified || undefined,
      offset: 0, // Reset to first page
    });
  };

  const clearFilters = () => {
    onParamsChange({
      limit: currentParams.limit,
      offset: 0,
    });
  };

  const hasActiveFilters = currentParams.rating || currentParams.verified !== undefined;

  return (
    <div className={cn('space-y-6', className)}>
      {/* Header with filters and sorting */}
      <div className="flex items-center justify-between flex-wrap gap-4">
        <div className="flex items-center gap-4">
          <h3 className="text-lg font-semibold">
            Reviews ({totalCount})
          </h3>
          
          {hasActiveFilters && (
            <Button variant="ghost" size="sm" onClick={clearFilters}>
              Clear Filters
            </Button>
          )}
        </div>

        <div className="flex items-center gap-2">
          {/* Rating Filter */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="flex items-center gap-2">
                <Filter className="w-4 h-4" />
                Rating
                {currentParams.rating && (
                  <Badge variant="secondary" className="ml-1">
                    {currentParams.rating}★
                  </Badge>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => handleRatingFilter(null)}>
                All Ratings
              </DropdownMenuItem>
              {[5, 4, 3, 2, 1].map((rating) => (
                <DropdownMenuItem key={rating} onClick={() => handleRatingFilter(rating)}>
                  <div className="flex items-center gap-2">
                    <span>{rating}</span>
                    <Star className="w-4 h-4 fill-yellow-400 text-yellow-400" />
                  </div>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Verified Filter */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="sm" className="flex items-center gap-2">
                <Filter className="w-4 h-4" />
                Type
                {currentParams.verified !== undefined && (
                  <Badge variant="secondary" className="ml-1">
                    {currentParams.verified ? 'Verified' : 'All'}
                  </Badge>
                )}
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={() => handleVerifiedFilter(null)}>
                All Reviews
              </DropdownMenuItem>
              <DropdownMenuItem onClick={() => handleVerifiedFilter(true)}>
                Verified Purchases Only
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          {/* Sort */}
          <Select
            value={currentParams.sort_by || 'created_at'}
            onValueChange={handleSortChange}
          >
            <SelectTrigger className="w-40">
              <SelectValue placeholder="Sort by" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="created_at">
                <div className="flex items-center gap-2">
                  Newest First
                  {currentParams.sort_by === 'created_at' && (
                    currentParams.sort_order === 'desc' ? 
                      <SortDesc className="w-4 h-4" /> : 
                      <SortAsc className="w-4 h-4" />
                  )}
                </div>
              </SelectItem>
              <SelectItem value="rating">
                <div className="flex items-center gap-2">
                  Rating
                  {currentParams.sort_by === 'rating' && (
                    currentParams.sort_order === 'desc' ? 
                      <SortDesc className="w-4 h-4" /> : 
                      <SortAsc className="w-4 h-4" />
                  )}
                </div>
              </SelectItem>
              <SelectItem value="helpful_count">
                <div className="flex items-center gap-2">
                  Most Helpful
                  {currentParams.sort_by === 'helpful_count' && (
                    currentParams.sort_order === 'desc' ? 
                      <SortDesc className="w-4 h-4" /> : 
                      <SortAsc className="w-4 h-4" />
                  )}
                </div>
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Reviews */}
      <div className="space-y-4">
        {reviews.length === 0 ? (
          <div className="text-center py-12">
            <Star className="w-12 h-12 text-gray-300 mx-auto mb-4" />
            <h4 className="text-lg font-medium text-gray-900 mb-2">No reviews yet</h4>
            <p className="text-gray-600">
              {hasActiveFilters 
                ? 'No reviews match your current filters.' 
                : 'Be the first to review this product!'
              }
            </p>
          </div>
        ) : (
          reviews.map((review) => (
            <ReviewCard
              key={review.id}
              review={review}
              onVote={onVote}
              onRemoveVote={onRemoveVote}
              onEdit={onEdit}
              onDelete={onDelete}
              onReply={onReply}
              isOwner={currentUserId === review.user.id}
              isAdmin={isAdmin}
            />
          ))
        )}
      </div>

      {/* Load More */}
      {hasMore && (
        <div className="text-center pt-6">
          <Button
            variant="outline"
            onClick={onLoadMore}
            disabled={isLoading}
            className="min-w-32"
          >
            {isLoading ? 'Loading...' : 'Load More Reviews'}
          </Button>
        </div>
      )}
    </div>
  );
}
