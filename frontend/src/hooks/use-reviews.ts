'use client';

import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { 
  reviewService, 
  Review, 
  ReviewsResponse, 
  ProductRatingSummary, 
  CreateReviewRequest, 
  UpdateReviewRequest,
  GetReviewsParams 
} from '@/services/review';

// Hook for product reviews
export function useProductReviews(productId: string, params?: GetReviewsParams) {
  return useQuery({
    queryKey: ['product-reviews', productId, params],
    queryFn: () => reviewService.getProductReviews(productId, params),
    enabled: !!productId,
  });
}

// Hook for product rating summary
export function useProductRatingSummary(productId: string) {
  return useQuery({
    queryKey: ['product-rating-summary', productId],
    queryFn: () => reviewService.getProductRatingSummary(productId),
    enabled: !!productId,
  });
}

// Hook for user reviews
export function useUserReviews(params?: GetReviewsParams) {
  return useQuery({
    queryKey: ['user-reviews', params],
    queryFn: () => reviewService.getUserReviews(params),
  });
}

// Hook for creating reviews
export function useCreateReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateReviewRequest) => reviewService.createReview(data),
    onSuccess: (review) => {
      toast.success('Review submitted successfully!');
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['product-reviews', review.product.id] });
      queryClient.invalidateQueries({ queryKey: ['product-rating-summary', review.product.id] });
      queryClient.invalidateQueries({ queryKey: ['user-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to submit review');
    },
  });
}

// Hook for updating reviews
export function useUpdateReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, data }: { reviewId: string; data: UpdateReviewRequest }) =>
      reviewService.updateReview(reviewId, data),
    onSuccess: (review) => {
      toast.success('Review updated successfully!');
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['product-reviews', review.product.id] });
      queryClient.invalidateQueries({ queryKey: ['product-rating-summary', review.product.id] });
      queryClient.invalidateQueries({ queryKey: ['user-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update review');
    },
  });
}

// Hook for deleting reviews
export function useDeleteReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (reviewId: string) => reviewService.deleteReview(reviewId),
    onSuccess: () => {
      toast.success('Review deleted successfully!');
      // Invalidate all review queries
      queryClient.invalidateQueries({ queryKey: ['product-reviews'] });
      queryClient.invalidateQueries({ queryKey: ['product-rating-summary'] });
      queryClient.invalidateQueries({ queryKey: ['user-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to delete review');
    },
  });
}

// Hook for voting on reviews
export function useVoteReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, isHelpful }: { reviewId: string; isHelpful: boolean }) =>
      reviewService.voteReview(reviewId, isHelpful),
    onSuccess: () => {
      // Invalidate review queries to refresh vote counts
      queryClient.invalidateQueries({ queryKey: ['product-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to vote on review');
    },
  });
}

// Hook for removing votes
export function useRemoveVote() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (reviewId: string) => reviewService.removeVote(reviewId),
    onSuccess: () => {
      // Invalidate review queries to refresh vote counts
      queryClient.invalidateQueries({ queryKey: ['product-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to remove vote');
    },
  });
}

// Hook for admin review management
export function useAdminReviews(params?: GetReviewsParams & { status?: string }) {
  return useQuery({
    queryKey: ['admin-reviews', params],
    queryFn: () => reviewService.getAdminReviews(params),
  });
}

// Hook for updating review status (admin)
export function useUpdateReviewStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, status }: { reviewId: string; status: 'approved' | 'hidden' | 'rejected' }) =>
      reviewService.updateReviewStatus(reviewId, status),
    onSuccess: () => {
      toast.success('Review status updated successfully!');
      // Invalidate admin and public review queries
      queryClient.invalidateQueries({ queryKey: ['admin-reviews'] });
      queryClient.invalidateQueries({ queryKey: ['product-reviews'] });
      queryClient.invalidateQueries({ queryKey: ['product-rating-summary'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to update review status');
    },
  });
}

// Hook for admin reply to review
export function useReplyToReview() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, reply }: { reviewId: string; reply: string }) =>
      reviewService.replyToReview(reviewId, reply),
    onSuccess: () => {
      toast.success('Reply added successfully!');
      // Invalidate review queries to show the reply
      queryClient.invalidateQueries({ queryKey: ['admin-reviews'] });
      queryClient.invalidateQueries({ queryKey: ['product-reviews'] });
    },
    onError: (error: any) => {
      toast.error(error.response?.data?.error || 'Failed to add reply');
    },
  });
}

// Custom hook for managing review list with pagination and filtering
export function useReviewList(productId: string) {
  const [params, setParams] = useState<GetReviewsParams>({
    limit: 10,
    offset: 0,
    sort_by: 'created_at',
    sort_order: 'desc',
  });

  const { data, isLoading, error } = useProductReviews(productId, params);
  const [allReviews, setAllReviews] = useState<Review[]>([]);

  // Reset reviews when productId changes
  useEffect(() => {
    setAllReviews([]);
    setParams(prev => ({ ...prev, offset: 0 }));
  }, [productId]);

  // Append new reviews when data changes
  useEffect(() => {
    if (data?.reviews) {
      if (params.offset === 0) {
        // First page or filter change - replace all reviews
        setAllReviews(data.reviews);
      } else {
        // Subsequent pages - append reviews
        setAllReviews(prev => [...prev, ...data.reviews]);
      }
    }
  }, [data, params.offset]);

  const loadMore = () => {
    if (data && allReviews.length < data.total_count) {
      setParams(prev => ({
        ...prev,
        offset: prev.offset + (prev.limit || 10),
      }));
    }
  };

  const updateParams = (newParams: GetReviewsParams) => {
    setParams(newParams);
  };

  const hasMore = data ? allReviews.length < data.total_count : false;

  return {
    reviews: allReviews,
    totalCount: data?.total_count || 0,
    params,
    updateParams,
    loadMore,
    hasMore,
    isLoading,
    error,
  };
}
