'use client';

import React from 'react';
import { Star } from 'lucide-react';
import { cn } from '@/lib/utils';

interface RatingStarsProps {
  rating: number;
  maxRating?: number;
  size?: 'sm' | 'md' | 'lg';
  interactive?: boolean;
  onRatingChange?: (rating: number) => void;
  className?: string;
  showValue?: boolean;
}

const sizeClasses = {
  sm: 'w-4 h-4',
  md: 'w-5 h-5',
  lg: 'w-6 h-6',
};

export function RatingStars({
  rating,
  maxRating = 5,
  size = 'md',
  interactive = false,
  onRatingChange,
  className,
  showValue = false,
}: RatingStarsProps) {
  const [hoverRating, setHoverRating] = React.useState(0);

  const handleStarClick = (starRating: number) => {
    if (interactive && onRatingChange) {
      onRatingChange(starRating);
    }
  };

  const handleStarHover = (starRating: number) => {
    if (interactive) {
      setHoverRating(starRating);
    }
  };

  const handleMouseLeave = () => {
    if (interactive) {
      setHoverRating(0);
    }
  };

  const displayRating = hoverRating || rating;

  return (
    <div className={cn('flex items-center gap-1', className)}>
      <div className="flex items-center" onMouseLeave={handleMouseLeave}>
        {Array.from({ length: maxRating }, (_, index) => {
          const starRating = index + 1;
          const isFilled = starRating <= displayRating;
          const isPartial = starRating - 0.5 <= displayRating && displayRating < starRating;

          return (
            <button
              key={index}
              type="button"
              className={cn(
                'relative transition-colors',
                interactive && 'hover:scale-110 cursor-pointer',
                !interactive && 'cursor-default'
              )}
              onClick={() => handleStarClick(starRating)}
              onMouseEnter={() => handleStarHover(starRating)}
              disabled={!interactive}
            >
              <Star
                className={cn(
                  sizeClasses[size],
                  'transition-colors',
                  isFilled
                    ? 'fill-yellow-400 text-yellow-400'
                    : isPartial
                    ? 'fill-yellow-200 text-yellow-400'
                    : 'fill-gray-200 text-gray-300'
                )}
              />
            </button>
          );
        })}
      </div>
      
      {showValue && (
        <span className="text-sm text-gray-600 ml-2">
          {rating.toFixed(1)} / {maxRating}
        </span>
      )}
    </div>
  );
}

// Rating Display Component (read-only with better styling)
interface RatingDisplayProps {
  rating: number;
  totalReviews?: number;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
  showCount?: boolean;
}

export function RatingDisplay({
  rating,
  totalReviews,
  size = 'md',
  className,
  showCount = true,
}: RatingDisplayProps) {
  return (
    <div className={cn('flex items-center gap-2', className)}>
      <RatingStars rating={rating} size={size} />
      <div className="flex items-center gap-1 text-sm text-gray-600">
        <span className="font-medium">{rating.toFixed(1)}</span>
        {showCount && totalReviews !== undefined && (
          <span>({totalReviews} review{totalReviews !== 1 ? 's' : ''})</span>
        )}
      </div>
    </div>
  );
}

// Rating Breakdown Component
interface RatingBreakdownProps {
  ratingCounts: { [key: string]: number };
  totalReviews: number;
  className?: string;
}

export function RatingBreakdown({ ratingCounts, totalReviews, className }: RatingBreakdownProps) {
  const ratings = [5, 4, 3, 2, 1];

  return (
    <div className={cn('space-y-2', className)}>
      {ratings.map((rating) => {
        const count = ratingCounts[rating.toString()] || 0;
        const percentage = totalReviews > 0 ? (count / totalReviews) * 100 : 0;

        return (
          <div key={rating} className="flex items-center gap-3">
            <div className="flex items-center gap-1 w-12">
              <span className="text-sm font-medium">{rating}</span>
              <Star className="w-3 h-3 fill-yellow-400 text-yellow-400" />
            </div>
            <div className="flex-1 bg-gray-200 rounded-full h-2">
              <div
                className="bg-yellow-400 h-2 rounded-full transition-all duration-300"
                style={{ width: `${percentage}%` }}
              />
            </div>
            <span className="text-sm text-gray-600 w-8 text-right">{count}</span>
          </div>
        );
      })}
    </div>
  );
}
