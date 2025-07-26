'use client';

import React, { useState, useEffect } from 'react';
import {
  MessageSquare,
  EyeOff,
  Reply,
  Search,
  MoreVertical,
  CheckCircle,
  XCircle,
  Star,
  Clock
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';

import { reviewService, Review, GetReviewsParams } from '@/lib/services/review';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';

// Simple RatingStars component
function RatingStars({ rating, size = 'sm' }: { rating: number; size?: 'sm' | 'md' | 'lg' }) {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6'
  }

  return (
    <div className="flex items-center gap-1">
      {Array.from({ length: 5 }, (_, i) => (
        <Star
          key={i}
          className={cn(
            sizeClasses[size],
            i < rating ? 'fill-yellow-400 text-yellow-400' : 'text-gray-300'
          )}
        />
      ))}
    </div>
  )
}

export function AdminReviewsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [replyingToReview, setReplyingToReview] = useState<Review | null>(null);
  const [replyText, setReplyText] = useState('');
  const [reviews, setReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [totalCount, setTotalCount] = useState(0);

  const [params, setParams] = useState<GetReviewsParams & { status?: string }>({
    limit: 20,
    offset: 0,
    sort_by: 'created_at',
    sort_order: 'desc',
  });

  // Load reviews
  useEffect(() => {
    loadReviews();
  }, [params, statusFilter]);

  const loadReviews = async () => {
    try {
      setIsLoading(true);
      const response = await reviewService.getAdminReviews({
        ...params,
        status: statusFilter === 'all' ? undefined : statusFilter,
      });
      setReviews(response.reviews);
      setTotalCount(response.total_count);
    } catch (error) {
      toast.error('Failed to load reviews');
      console.error('Error loading reviews:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleStatusChange = async (reviewId: string, status: 'approved' | 'hidden' | 'rejected') => {
    try {
      await reviewService.updateReviewStatus(reviewId, status);
      toast.success(`Review ${status} successfully`);
      loadReviews(); // Reload reviews
    } catch (error) {
      toast.error(`Failed to ${status} review`);
      console.error('Error updating review status:', error);
    }
  };

  const handleReply = async () => {
    if (!replyingToReview || !replyText.trim()) return;

    try {
      await reviewService.replyToReview(replyingToReview.id, replyText.trim());
      toast.success('Reply sent successfully');
      setReplyingToReview(null);
      setReplyText('');
      loadReviews(); // Reload reviews
    } catch (error) {
      toast.error('Failed to send reply');
      console.error('Error sending reply:', error);
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'approved':
        return <Badge className="bg-green-900/30 text-green-400 border-green-700">Approved</Badge>;
      case 'hidden':
        return <Badge className="bg-yellow-900/30 text-yellow-400 border-yellow-700">Hidden</Badge>;
      case 'rejected':
        return <Badge className="bg-red-900/30 text-red-400 border-red-700">Rejected</Badge>;
      case 'pending':
        return <Badge className="bg-gray-700/50 text-gray-300 border-gray-600">Pending</Badge>;
      default:
        return <Badge variant="outline" className="border-gray-600 text-gray-400">{status}</Badge>;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const filteredReviews = reviews.filter((review: Review) =>
    review.comment?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    review.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    (review.user && `${review.user.first_name ?? ''} ${review.user.last_name ?? ''}`.toLowerCase().includes(searchTerm.toLowerCase()))
  );

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Review Management</h1>
          <p className="text-gray-300">Manage customer reviews and feedback</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-sm border-[#FF9000] text-[#FF9000]">
            {totalCount} Total Reviews
          </Badge>
        </div>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-xl flex items-center justify-center">
                <MessageSquare className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Total Reviews</p>
                <p className="text-2xl font-bold text-white">{totalCount}</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-xl flex items-center justify-center">
                <Clock className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Pending</p>
                <p className="text-2xl font-bold text-white">
                  {reviews.filter(r => r.status === 'pending').length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center">
                <CheckCircle className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Approved</p>
                <p className="text-2xl font-bold text-white">
                  {reviews.filter(r => r.status === 'approved').length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
          <CardContent className="p-6">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-gradient-to-br from-red-500 to-red-600 rounded-xl flex items-center justify-center">
                <XCircle className="w-6 h-6 text-white" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Rejected</p>
                <p className="text-2xl font-bold text-white">
                  {reviews.filter(r => r.status === 'rejected').length}
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Filters */}
      <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
        <CardHeader>
          <CardTitle className="text-lg text-white">Filters</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col md:flex-row gap-4">
            <div className="flex-1">
              <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                <Input
                  placeholder="Search reviews, users, or products..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-10 bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />
              </div>
            </div>
            
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-48 bg-gray-700 border-gray-600 text-white">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent className="bg-gray-800 border-gray-600">
                <SelectItem value="all" className="text-white hover:bg-gray-700">All Reviews</SelectItem>
                <SelectItem value="pending" className="text-white hover:bg-gray-700">Pending</SelectItem>
                <SelectItem value="approved" className="text-white hover:bg-gray-700">Approved</SelectItem>
                <SelectItem value="hidden" className="text-white hover:bg-gray-700">Hidden</SelectItem>
                <SelectItem value="rejected" className="text-white hover:bg-gray-700">Rejected</SelectItem>
              </SelectContent>
            </Select>

            <Select
              value={`${params.sort_by}-${params.sort_order}`}
              onValueChange={(value) => {
                const [sort_by, sort_order] = value.split('-');
                setParams(prev => ({ ...prev, sort_by: sort_by as any, sort_order: sort_order as any }));
              }}
            >
              <SelectTrigger className="w-48 bg-gray-700 border-gray-600 text-white">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent className="bg-gray-800 border-gray-600">
                <SelectItem value="created_at-desc" className="text-white hover:bg-gray-700">Newest First</SelectItem>
                <SelectItem value="created_at-asc" className="text-white hover:bg-gray-700">Oldest First</SelectItem>
                <SelectItem value="rating-desc" className="text-white hover:bg-gray-700">Highest Rating</SelectItem>
                <SelectItem value="rating-asc" className="text-white hover:bg-gray-700">Lowest Rating</SelectItem>
                <SelectItem value="helpful_count-desc" className="text-white hover:bg-gray-700">Most Helpful</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
      </Card>

      {/* Reviews List */}
      <div className="space-y-4">
        {isLoading ? (
          <div className="space-y-4">
            {Array.from({ length: 5 }, (_, i) => (
              <Card key={i} className="animate-pulse bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
                <CardContent className="p-6">
                  <div className="h-32 bg-gray-700 rounded"></div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : filteredReviews.length === 0 ? (
          <Card className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
            <CardContent className="text-center py-12">
              <MessageSquare className="w-12 h-12 text-gray-500 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-white mb-2">No reviews found</h3>
              <p className="text-gray-400">
                {searchTerm || statusFilter !== 'all'
                  ? 'Try adjusting your filters'
                  : 'No reviews have been submitted yet'
                }
              </p>
            </CardContent>
          </Card>
        ) : (
          filteredReviews.map((review) => (
            <Card key={review.id} className="bg-gradient-to-br from-gray-800 to-gray-900 border-gray-700">
              <CardContent className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-start gap-4">
                    <div className="w-10 h-10 bg-gradient-to-br from-[#FF9000] to-[#e67e00] rounded-full flex items-center justify-center text-white font-semibold">
                      {review.user ? `${review.user.first_name?.[0] ?? ''}${review.user.last_name?.[0] ?? ''}` : ''}
                    </div>

                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-semibold text-white">
                          {review.user ? `${review.user.first_name ?? ''} ${review.user.last_name ?? ''}` : ''}
                        </h4>
                        {review.is_verified && (
                          <Badge variant="secondary" className="text-xs bg-green-900/30 text-green-400 border-green-700">Verified</Badge>
                        )}
                        {getStatusBadge(review.status)}
                      </div>

                      <div className="flex items-center gap-2 mb-2">
                        <RatingStars rating={review.rating} size="sm" />
                        <span className="text-sm text-gray-400">
                          {formatDate(review.created_at)}
                        </span>
                      </div>

                      <p className="text-sm text-gray-400 mb-2">
                        Product: {review.product ? review.product.name : ''}
                      </p>
                    </div>
                  </div>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm" className="text-gray-400 hover:text-white hover:bg-gray-700">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="bg-gray-800 border-gray-600">
                      <DropdownMenuItem
                        onClick={() => handleStatusChange(review.id, 'approved')}
                        disabled={review.status === 'approved'}
                        className="text-white hover:bg-gray-700 focus:bg-gray-700"
                      >
                        <CheckCircle className="w-4 h-4 mr-2 text-green-400" />
                        Approve
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => handleStatusChange(review.id, 'hidden')}
                        disabled={review.status === 'hidden'}
                        className="text-white hover:bg-gray-700 focus:bg-gray-700"
                      >
                        <EyeOff className="w-4 h-4 mr-2 text-yellow-400" />
                        Hide
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => handleStatusChange(review.id, 'rejected')}
                        disabled={review.status === 'rejected'}
                        className="text-red-400 hover:bg-gray-700 focus:bg-gray-700"
                      >
                        <XCircle className="w-4 h-4 mr-2" />
                        Reject
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        onClick={() => setReplyingToReview(review)}
                        className="text-white hover:bg-gray-700 focus:bg-gray-700"
                      >
                        <Reply className="w-4 h-4 mr-2" />
                        Reply
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                {/* Review Content */}
                <div className="mb-4">
                  {review.title && (
                    <h5 className="font-semibold text-white mb-2">{review.title}</h5>
                  )}
                  <p className="text-gray-300 leading-relaxed">{review.comment}</p>
                </div>

                {/* Admin Reply */}
                {review.admin_reply && (
                  <div className="bg-gradient-to-br from-[#FF9000]/10 to-[#e67e00]/10 border border-[#FF9000]/30 rounded-lg p-4 mt-4">
                    <div className="flex items-center gap-2 mb-2">
                      <Badge variant="outline" className="text-[#FF9000] border-[#FF9000]/50 bg-[#FF9000]/10">
                        Admin Response
                      </Badge>
                      {review.admin_reply_at && (
                        <span className="text-sm text-gray-400">
                          {formatDate(review.admin_reply_at)}
                        </span>
                      )}
                    </div>
                    <p className="text-gray-300">{review.admin_reply}</p>
                  </div>
                )}

                {/* Review Stats */}
                <div className="flex items-center gap-4 pt-4 border-t border-gray-700 text-sm text-gray-400">
                  <span>👍 {review.helpful_count} helpful</span>
                  <span>👎 {review.not_helpful_count} not helpful</span>
                  <span>{typeof review.helpful_percentage === 'number' ? review.helpful_percentage.toFixed(0) : '0'}% helpful rate</span>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {/* Reply Dialog */}
      <Dialog open={!!replyingToReview} onOpenChange={() => setReplyingToReview(null)}>
        <DialogContent className="max-w-2xl bg-gray-800 border-gray-700">
          <DialogHeader>
            <DialogTitle className="text-white">Reply to Review</DialogTitle>
          </DialogHeader>

          {replyingToReview && (
            <div className="space-y-4">
              {/* Review Preview */}
              <div className="bg-gray-700 p-4 rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <RatingStars rating={replyingToReview.rating} size="sm" />
                  <span className="font-medium text-white">
                    {replyingToReview.user ? `${replyingToReview.user.first_name ?? ''} ${replyingToReview.user.last_name ?? ''}` : ''}
                  </span>
                </div>
                <p className="text-gray-300">{replyingToReview.comment}</p>
              </div>

              {/* Reply Form */}
              <div className="space-y-4">
                <Textarea
                  placeholder="Write your response to this review..."
                  value={replyText}
                  onChange={(e) => setReplyText(e.target.value)}
                  rows={4}
                  className="bg-gray-700 border-gray-600 text-white placeholder-gray-400"
                />

                <div className="flex items-center gap-2 justify-end">
                  <Button
                    variant="outline"
                    onClick={() => setReplyingToReview(null)}
                    className="border-gray-600 text-gray-300 hover:bg-gray-700"
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={handleReply}
                    disabled={!replyText.trim()}
                    className="bg-gradient-to-r from-[#FF9000] to-[#e67e00] hover:from-[#e67e00] hover:to-[#cc6600]"
                  >
                    Send Reply
                  </Button>
                </div>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
