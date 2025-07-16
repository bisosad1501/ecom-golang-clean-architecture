'use client';

import React, { useState } from 'react';
import { Plus, MessageSquare } from 'lucide-react';
import { 
  ProductReviewSummary, 
  ReviewList, 
  ReviewForm 
} from '@/components/reviews';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { 
  useProductRatingSummary, 
  useReviewList, 
  useCreateReview, 
  useUpdateReview, 
  useDeleteReview, 
  useVoteReview, 
  useRemoveVote 
} from '@/hooks/use-reviews';
import { useAuthStore } from '@/store/auth';
import { Review, CreateReviewRequest, UpdateReviewRequest } from '@/lib/services/review';
import { cn } from '@/lib/utils';

interface ProductReviewsSectionProps {
  productId: string;
  className?: string;
}

export function ProductReviewsSection({ productId, className }: ProductReviewsSectionProps) {
  const [showReviewForm, setShowReviewForm] = useState(false);
  const [editingReview, setEditingReview] = useState<Review | null>(null);
  
  const { user, isAuthenticated } = useAuthStore();
  
  // Queries
  const { data: ratingSummary, isLoading: summaryLoading } = useProductRatingSummary(productId);
  const {
    reviews,
    totalCount,
    params,
    updateParams,
    loadMore,
    hasMore,
    isLoading: reviewsLoading,
  } = useReviewList(productId);

  // Mutations
  const createReviewMutation = useCreateReview();
  const updateReviewMutation = useUpdateReview();
  const deleteReviewMutation = useDeleteReview();
  const voteReviewMutation = useVoteReview();
  const removeVoteMutation = useRemoveVote();

  // Allow multiple reviews per user
  const canWriteReview = isAuthenticated;

  const handleCreateReview = async (data: CreateReviewRequest) => {
    try {
      await createReviewMutation.mutateAsync(data);
      setShowReviewForm(false);
    } catch (error) {
      // Error is handled by the mutation
    }
  };

  const handleUpdateReview = async (data: UpdateReviewRequest) => {
    if (!editingReview) return;
    
    try {
      await updateReviewMutation.mutateAsync({
        reviewId: editingReview.id,
        data,
      });
      setEditingReview(null);
    } catch (error) {
      // Error is handled by the mutation
    }
  };

  const handleDeleteReview = async (reviewId: string) => {
    if (confirm('Are you sure you want to delete this review?')) {
      try {
        await deleteReviewMutation.mutateAsync(reviewId);
      } catch (error) {
        // Error is handled by the mutation
      }
    }
  };

  const handleVoteReview = async (reviewId: string, isHelpful: boolean) => {
    if (!isAuthenticated) {
      // Redirect to login or show login modal
      return;
    }

    try {
      await voteReviewMutation.mutateAsync({ reviewId, isHelpful });
    } catch (error) {
      // Error is handled by the mutation
    }
  };

  const handleRemoveVote = async (reviewId: string) => {
    try {
      await removeVoteMutation.mutateAsync(reviewId);
    } catch (error) {
      // Error is handled by the mutation
    }
  };

  const handleEditReview = (review: Review) => {
    setEditingReview(review);
  };

  if (summaryLoading) {
    return (
      <div className={cn('space-y-6', className)}>
        <div className="animate-pulse">
          <div className="h-48 bg-gray-200 rounded-lg mb-6"></div>
          <div className="space-y-4">
            {Array.from({ length: 3 }, (_, i) => (
              <div key={i} className="h-32 bg-gray-200 rounded-lg"></div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (!ratingSummary) {
    return (
      <div className={cn('text-center py-8', className)}>
        <MessageSquare className="w-12 h-12 text-gray-300 mx-auto mb-4" />
        <p className="text-gray-600">Unable to load reviews</p>
      </div>
    );
  }

  return (
    <div className={cn('space-y-8', className)}>
      {/* Review Summary */}
      <ProductReviewSummary
        summary={ratingSummary}
        canWriteReview={canWriteReview}
        onWriteReview={() => setShowReviewForm(true)}
      />

      {/* Review List */}
      <ReviewList
        reviews={reviews}
        totalCount={totalCount}
        currentParams={params}
        onParamsChange={updateParams}
        onVote={handleVoteReview}
        onRemoveVote={handleRemoveVote}
        onEdit={handleEditReview}
        onDelete={handleDeleteReview}
        onLoadMore={loadMore}
        isLoading={reviewsLoading}
        hasMore={hasMore}
        currentUserId={user?.id}
      />

      {/* Write Review Dialog */}
      <Dialog open={showReviewForm} onOpenChange={setShowReviewForm}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Rate This Product</DialogTitle>
          </DialogHeader>
          <ReviewForm
            productId={productId}
            onSubmit={handleCreateReview}
            onCancel={() => setShowReviewForm(false)}
            isLoading={createReviewMutation.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Review Dialog */}
      <Dialog open={!!editingReview} onOpenChange={() => setEditingReview(null)}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Your Review</DialogTitle>
          </DialogHeader>
          {editingReview && (
            <ReviewForm
              productId={productId}
              existingReview={editingReview}
              onSubmit={handleUpdateReview}
              onCancel={() => setEditingReview(null)}
              isLoading={updateReviewMutation.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Floating Write Review Button (Mobile) */}
      {canWriteReview && (
        <div className="fixed bottom-6 right-6 md:hidden z-50">
          <Button
            onClick={() => setShowReviewForm(true)}
            size="lg"
            className="rounded-full shadow-lg"
          >
            <Plus className="w-5 h-5 mr-2" />
            Rate
          </Button>
        </div>
      )}
    </div>
  );
}
