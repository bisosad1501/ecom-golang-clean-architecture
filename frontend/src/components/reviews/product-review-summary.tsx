'use client';

import React from 'react';
import { Star, TrendingUp } from 'lucide-react';
import { ProductRatingSummary } from '@/services/review';
import { RatingDisplay, RatingBreakdown } from './rating-stars';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { cn } from '@/lib/utils';

interface ProductReviewSummaryProps {
  summary: ProductRatingSummary;
  onWriteReview?: () => void;
  canWriteReview?: boolean;
  className?: string;
}

export function ProductReviewSummary({
  summary,
  onWriteReview,
  canWriteReview = false,
  className,
}: ProductReviewSummaryProps) {
  const { average_rating, total_reviews, rating_counts } = summary;

  // Calculate recommendation percentage (4-5 star reviews)
  const positiveReviews = (rating_counts['4'] || 0) + (rating_counts['5'] || 0);
  const recommendationPercentage = total_reviews > 0 ? (positiveReviews / total_reviews) * 100 : 0;

  return (
    <Card className={cn('w-full', className)}>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>Customer Reviews</span>
          {canWriteReview && onWriteReview && (
            <Button onClick={onWriteReview} size="sm">
              Add Your Rating
            </Button>
          )}
        </CardTitle>
      </CardHeader>
      
      <CardContent className="space-y-6">
        {total_reviews === 0 ? (
          // No reviews state
          <div className="text-center py-8">
            <Star className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              No reviews yet
            </h3>
            <p className="text-gray-600 mb-4">
              Share your experience with this product - even just a star rating helps!
            </p>
            {canWriteReview && onWriteReview && (
              <Button onClick={onWriteReview}>
                Rate This Product
              </Button>
            )}
          </div>
        ) : (
          <>
            {/* Overall Rating */}
            <div className="flex items-start gap-6">
              <div className="text-center">
                <div className="text-4xl font-bold text-gray-900 mb-1">
                  {average_rating.toFixed(1)}
                </div>
                <RatingDisplay 
                  rating={average_rating} 
                  totalReviews={total_reviews}
                  size="md"
                  className="justify-center"
                />
              </div>
              
              <div className="flex-1">
                <RatingBreakdown 
                  ratingCounts={rating_counts}
                  totalReviews={total_reviews}
                />
              </div>
            </div>

            {/* Recommendation */}
            {recommendationPercentage > 0 && (
              <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                <div className="flex items-center gap-2 mb-2">
                  <TrendingUp className="w-5 h-5 text-green-600" />
                  <span className="font-semibold text-green-800">
                    {recommendationPercentage.toFixed(0)}% of customers recommend this product
                  </span>
                </div>
                <Progress 
                  value={recommendationPercentage} 
                  className="h-2 bg-green-100"
                />
              </div>
            )}

            {/* Quick Stats */}
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t">
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-900">
                  {average_rating.toFixed(1)}
                </div>
                <div className="text-sm text-gray-600">Average Rating</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-900">
                  {total_reviews}
                </div>
                <div className="text-sm text-gray-600">Total Reviews</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-900">
                  {rating_counts['5'] || 0}
                </div>
                <div className="text-sm text-gray-600">5-Star Reviews</div>
              </div>
              
              <div className="text-center">
                <div className="text-2xl font-bold text-gray-900">
                  {recommendationPercentage.toFixed(0)}%
                </div>
                <div className="text-sm text-gray-600">Recommend</div>
              </div>
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

// Compact version for product cards
interface CompactReviewSummaryProps {
  summary: ProductRatingSummary;
  className?: string;
}

export function CompactReviewSummary({ summary, className }: CompactReviewSummaryProps) {
  const { average_rating, total_reviews } = summary;

  if (total_reviews === 0) {
    return (
      <div className={cn('flex items-center gap-2 text-sm text-gray-500', className)}>
        <div className="flex">
          {Array.from({ length: 5 }, (_, i) => (
            <Star key={i} className="w-4 h-4 text-gray-300" />
          ))}
        </div>
        <span>No reviews</span>
      </div>
    );
  }

  return (
    <div className={cn('flex items-center gap-2', className)}>
      <RatingDisplay 
        rating={average_rating}
        totalReviews={total_reviews}
        size="sm"
        showCount={true}
      />
    </div>
  );
}
