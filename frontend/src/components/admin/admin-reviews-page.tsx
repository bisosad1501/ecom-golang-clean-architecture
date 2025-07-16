'use client';

import React, { useState } from 'react';
import { 
  MessageSquare, 
  Eye, 
  EyeOff, 
  Trash2, 
  Reply, 
  Filter,
  Search,
  MoreVertical,
  CheckCircle,
  XCircle,
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
import { RatingStars } from '@/components/reviews/rating-stars';
import { 
  useAdminReviews, 
  useUpdateReviewStatus, 
  useReplyToReview 
} from '@/hooks/use-reviews';
import { Review, GetReviewsParams } from '@/lib/services/review';
import { cn } from '@/lib/utils';

export function AdminReviewsPage() {
  const [searchTerm, setSearchTerm] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [replyingToReview, setReplyingToReview] = useState<Review | null>(null);
  const [replyText, setReplyText] = useState('');

  const [params, setParams] = useState<GetReviewsParams & { status?: string }>({
    limit: 20,
    offset: 0,
    sort_by: 'created_at',
    sort_order: 'desc',
  });

  // Queries and mutations
  const { data: reviewsData, isLoading } = useAdminReviews({
    ...params,
    status: statusFilter === 'all' ? undefined : statusFilter,
  });

  const updateStatusMutation = useUpdateReviewStatus();
  const replyMutation = useReplyToReview();

  const handleStatusChange = (reviewId: string, status: 'approved' | 'hidden' | 'rejected') => {
    updateStatusMutation.mutate({ reviewId, status });
  };

  const handleReply = async () => {
    if (!replyingToReview || !replyText.trim()) return;

    try {
      await replyMutation.mutateAsync({
        reviewId: replyingToReview.id,
        reply: replyText.trim(),
      });
      setReplyingToReview(null);
      setReplyText('');
    } catch (error) {
      // Error handled by mutation
    }
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'approved':
        return <Badge className="bg-green-100 text-green-800">Approved</Badge>;
      case 'hidden':
        return <Badge className="bg-yellow-100 text-yellow-800">Hidden</Badge>;
      case 'rejected':
        return <Badge className="bg-red-100 text-red-800">Rejected</Badge>;
      case 'pending':
        return <Badge className="bg-gray-100 text-gray-800">Pending</Badge>;
      default:
        return <Badge variant="outline">{status}</Badge>;
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

  const filteredReviews = reviewsData?.reviews.filter(review =>
    review.comment.toLowerCase().includes(searchTerm.toLowerCase()) ||
    review.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    `${review.user.first_name} ${review.user.last_name}`.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Review Management</h1>
          <p className="text-gray-600">Manage customer reviews and feedback</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline" className="text-sm">
            {reviewsData?.total_count || 0} Total Reviews
          </Badge>
        </div>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Filters</CardTitle>
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
                  className="pl-10"
                />
              </div>
            </div>
            
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Filter by status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Reviews</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="approved">Approved</SelectItem>
                <SelectItem value="hidden">Hidden</SelectItem>
                <SelectItem value="rejected">Rejected</SelectItem>
              </SelectContent>
            </Select>

            <Select 
              value={`${params.sort_by}-${params.sort_order}`}
              onValueChange={(value) => {
                const [sort_by, sort_order] = value.split('-');
                setParams(prev => ({ ...prev, sort_by: sort_by as any, sort_order: sort_order as any }));
              }}
            >
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Sort by" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="created_at-desc">Newest First</SelectItem>
                <SelectItem value="created_at-asc">Oldest First</SelectItem>
                <SelectItem value="rating-desc">Highest Rating</SelectItem>
                <SelectItem value="rating-asc">Lowest Rating</SelectItem>
                <SelectItem value="helpful_count-desc">Most Helpful</SelectItem>
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
              <Card key={i} className="animate-pulse">
                <CardContent className="p-6">
                  <div className="h-32 bg-gray-200 rounded"></div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : filteredReviews.length === 0 ? (
          <Card>
            <CardContent className="text-center py-12">
              <MessageSquare className="w-12 h-12 text-gray-300 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">No reviews found</h3>
              <p className="text-gray-600">
                {searchTerm || statusFilter !== 'all' 
                  ? 'Try adjusting your filters' 
                  : 'No reviews have been submitted yet'
                }
              </p>
            </CardContent>
          </Card>
        ) : (
          filteredReviews.map((review) => (
            <Card key={review.id}>
              <CardContent className="p-6">
                <div className="flex items-start justify-between mb-4">
                  <div className="flex items-start gap-4">
                    <div className="w-10 h-10 bg-gradient-to-br from-orange-400 to-orange-600 rounded-full flex items-center justify-center text-white font-semibold">
                      {review.user.first_name[0]}{review.user.last_name[0]}
                    </div>
                    
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-semibold text-gray-900">
                          {review.user.first_name} {review.user.last_name}
                        </h4>
                        {review.is_verified && (
                          <Badge variant="secondary" className="text-xs">Verified</Badge>
                        )}
                        {getStatusBadge(review.status)}
                      </div>
                      
                      <div className="flex items-center gap-2 mb-2">
                        <RatingStars rating={review.rating} size="sm" />
                        <span className="text-sm text-gray-500">
                          {formatDate(review.created_at)}
                        </span>
                      </div>
                      
                      <p className="text-sm text-gray-600 mb-2">
                        Product: {review.product.name}
                      </p>
                    </div>
                  </div>

                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <MoreVertical className="w-4 h-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem 
                        onClick={() => handleStatusChange(review.id, 'approved')}
                        disabled={review.status === 'approved'}
                      >
                        <CheckCircle className="w-4 h-4 mr-2 text-green-600" />
                        Approve
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => handleStatusChange(review.id, 'hidden')}
                        disabled={review.status === 'hidden'}
                      >
                        <EyeOff className="w-4 h-4 mr-2 text-yellow-600" />
                        Hide
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => handleStatusChange(review.id, 'rejected')}
                        disabled={review.status === 'rejected'}
                        className="text-red-600"
                      >
                        <XCircle className="w-4 h-4 mr-2" />
                        Reject
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => setReplyingToReview(review)}>
                        <Reply className="w-4 h-4 mr-2" />
                        Reply
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>

                {/* Review Content */}
                <div className="mb-4">
                  {review.title && (
                    <h5 className="font-semibold text-gray-900 mb-2">{review.title}</h5>
                  )}
                  <p className="text-gray-700 leading-relaxed">{review.comment}</p>
                </div>

                {/* Admin Reply */}
                {review.admin_reply && (
                  <div className="bg-orange-50 border border-orange-200 rounded-lg p-4 mt-4">
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

                {/* Review Stats */}
                <div className="flex items-center gap-4 pt-4 border-t text-sm text-gray-600">
                  <span>👍 {review.helpful_count} helpful</span>
                  <span>👎 {review.not_helpful_count} not helpful</span>
                  <span>{review.helpful_percentage.toFixed(0)}% helpful rate</span>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>

      {/* Reply Dialog */}
      <Dialog open={!!replyingToReview} onOpenChange={() => setReplyingToReview(null)}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Reply to Review</DialogTitle>
          </DialogHeader>
          
          {replyingToReview && (
            <div className="space-y-4">
              {/* Review Preview */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <div className="flex items-center gap-2 mb-2">
                  <RatingStars rating={replyingToReview.rating} size="sm" />
                  <span className="font-medium">
                    {replyingToReview.user.first_name} {replyingToReview.user.last_name}
                  </span>
                </div>
                <p className="text-gray-700">{replyingToReview.comment}</p>
              </div>

              {/* Reply Form */}
              <div className="space-y-4">
                <Textarea
                  placeholder="Write your response to this review..."
                  value={replyText}
                  onChange={(e) => setReplyText(e.target.value)}
                  rows={4}
                />
                
                <div className="flex items-center gap-2 justify-end">
                  <Button 
                    variant="outline" 
                    onClick={() => setReplyingToReview(null)}
                  >
                    Cancel
                  </Button>
                  <Button 
                    onClick={handleReply}
                    disabled={!replyText.trim() || replyMutation.isPending}
                  >
                    {replyMutation.isPending ? 'Sending...' : 'Send Reply'}
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
