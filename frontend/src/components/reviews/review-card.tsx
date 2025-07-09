'use client';

import React from 'react';
import { ThumbsUp, ThumbsDown, Shield, MoreVertical, Reply } from 'lucide-react';
import { Review } from '@/services/review';
import { RatingStars } from './rating-stars';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';

interface ReviewCardProps {
  review: Review;
  onVote?: (reviewId: string, isHelpful: boolean) => void;
  onRemoveVote?: (reviewId: string) => void;
  onEdit?: (review: Review) => void;
  onDelete?: (reviewId: string) => void;
  onReply?: (reviewId: string) => void;
  showActions?: boolean;
  isOwner?: boolean;
  isAdmin?: boolean;
  className?: string;
}

export function ReviewCard({
  review,
  onVote,
  onRemoveVote,
  onEdit,
  onDelete,
  onReply,
  showActions = true,
  isOwner = false,
  isAdmin = false,
  className,
}: ReviewCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const handleVote = (isHelpful: boolean) => {
    if (review.user_vote === (isHelpful ? 'helpful' : 'not_helpful')) {
      // Remove vote if clicking the same vote
      onRemoveVote?.(review.id);
    } else {
      // Add or change vote
      onVote?.(review.id, isHelpful);
    }
  };

  return (
    <Card className={cn('w-full', className)}>
      <CardHeader className="pb-4">
        <div className="flex items-start justify-between">
          <div className="flex items-start gap-3">
            {/* User Avatar */}
            <div className="w-10 h-10 bg-gradient-to-br from-orange-400 to-orange-600 rounded-full flex items-center justify-center text-white font-semibold">
              {review.user.first_name[0]}{review.user.last_name[0]}
            </div>
            
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <h4 className="font-semibold text-gray-900">
                  {review.user.first_name} {review.user.last_name}
                </h4>
                {review.is_verified && (
                  <Badge variant="secondary" className="text-xs">
                    <Shield className="w-3 h-3 mr-1" />
                    Verified Purchase
                  </Badge>
                )}
              </div>
              
              <div className="flex items-center gap-2 mb-2">
                <RatingStars rating={review.rating} size="sm" />
                <span className="text-sm text-gray-500">
                  {formatDate(review.created_at)}
                </span>
              </div>
            </div>
          </div>

          {/* Actions Menu */}
          {showActions && (isOwner || isAdmin) && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm">
                  <MoreVertical className="w-4 h-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {isOwner && (
                  <>
                    <DropdownMenuItem onClick={() => onEdit?.(review)}>
                      Edit Review
                    </DropdownMenuItem>
                    <DropdownMenuItem 
                      onClick={() => onDelete?.(review.id)}
                      className="text-red-600"
                    >
                      Delete Review
                    </DropdownMenuItem>
                  </>
                )}
                {isAdmin && (
                  <DropdownMenuItem onClick={() => onReply?.(review.id)}>
                    <Reply className="w-4 h-4 mr-2" />
                    Reply as Admin
                  </DropdownMenuItem>
                )}
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
      </CardHeader>

      <CardContent className="pt-0">
        {/* Review Title */}
        {review.title && (
          <h5 className="font-semibold text-gray-900 mb-2">{review.title}</h5>
        )}

        {/* Review Comment */}
        {review.comment ? (
          <p className="text-gray-700 mb-4 leading-relaxed">{review.comment}</p>
        ) : (
          <p className="text-gray-500 italic mb-4">Customer rated this product {review.rating} star{review.rating !== 1 ? 's' : ''}</p>
        )}

        {/* Review Images */}
        {review.images && review.images.length > 0 && (
          <div className="flex gap-2 mb-4 overflow-x-auto">
            {review.images.map((image, index) => (
              <img
                key={index}
                src={image}
                alt={`Review image ${index + 1}`}
                className="w-20 h-20 object-cover rounded-lg border"
              />
            ))}
          </div>
        )}

        {/* Admin Reply */}
        {review.admin_reply && (
          <div className="bg-orange-50 border border-orange-200 rounded-lg p-4 mb-4">
            <div className="flex items-center gap-2 mb-2">
              <Badge variant="outline" className="text-orange-600 border-orange-300">
                Admin Response
              </Badge>
              {review.admin_reply_at && (
                <span className="text-sm text-gray-500">
                  {formatDate(review.admin_reply_at)}
                </span>
              )}
            </div>
            <p className="text-gray-700">{review.admin_reply}</p>
          </div>
        )}

        {/* Helpful Votes */}
        {showActions && (
          <div className="flex items-center justify-between pt-4 border-t">
            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleVote(true)}
                className={cn(
                  'flex items-center gap-2',
                  review.user_vote === 'helpful' && 'text-green-600 bg-green-50'
                )}
              >
                <ThumbsUp className="w-4 h-4" />
                Helpful ({review.helpful_count})
              </Button>
              
              <Button
                variant="ghost"
                size="sm"
                onClick={() => handleVote(false)}
                className={cn(
                  'flex items-center gap-2',
                  review.user_vote === 'not_helpful' && 'text-red-600 bg-red-50'
                )}
              >
                <ThumbsDown className="w-4 h-4" />
                Not Helpful ({review.not_helpful_count})
              </Button>
            </div>

            {review.helpful_count + review.not_helpful_count > 0 && (
              <span className="text-sm text-gray-500">
                {review.helpful_percentage.toFixed(0)}% found this helpful
              </span>
            )}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
