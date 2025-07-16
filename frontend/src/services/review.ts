import { apiClient } from '@/lib/api';

// Types
export interface ReviewUser {
  id: string;
  first_name: string;
  last_name: string;
  avatar?: string;
}

export interface ReviewProduct {
  id: string;
  name: string;
  image?: string;
}

export interface Review {
  id: string;
  user: ReviewUser;
  product: ReviewProduct;
  rating: number;
  title: string;
  comment: string;
  status: 'pending' | 'approved' | 'hidden' | 'rejected';
  is_verified: boolean;
  admin_reply?: string;
  admin_reply_at?: string;
  helpful_count: number;
  not_helpful_count: number;
  helpful_percentage: number;
  images?: string[];
  user_vote?: 'helpful' | 'not_helpful';
  created_at: string;
  updated_at: string;
}

export interface ReviewsResponse {
  reviews: Review[];
  total_count: number;
  limit: number;
  offset: number;
}

export interface ProductRatingSummary {
  product_id: string;
  average_rating: number;
  total_reviews: number;
  rating_counts: {
    [key: string]: number; // "1": 5, "2": 10, etc.
  };
}

export interface CreateReviewRequest {
  product_id: string;
  order_id?: string;
  rating: number;
  title: string;
  comment: string;
  images?: string[];
}

export interface UpdateReviewRequest {
  rating?: number;
  title?: string;
  comment?: string;
  images?: string[];
}

export interface GetReviewsParams {
  rating?: number;
  verified?: boolean;
  sort_by?: 'created_at' | 'rating' | 'helpful_count';
  sort_order?: 'asc' | 'desc';
  limit?: number;
  offset?: number;
}

// API Service Class
class ReviewService {
  // Public APIs (no auth required)
  async getProductReviews(productId: string, params?: GetReviewsParams): Promise<ReviewsResponse> {
    const response = await apiClient.get(`/public/reviews/product/${productId}`, { params });
    return response.data;
  }

  async getProductRatingSummary(productId: string): Promise<ProductRatingSummary> {
    const response = await apiClient.get(`/public/reviews/product/${productId}/summary`);
    return response.data;
  }

  // Protected APIs (auth required)
  async createReview(data: CreateReviewRequest): Promise<Review> {
    const response = await apiClient.post(`/reviews`, data);
    return response.data;
  }

  async updateReview(reviewId: string, data: UpdateReviewRequest): Promise<Review> {
    const response = await apiClient.put(`/reviews/${reviewId}`, data);
    return response.data;
  }

  async deleteReview(reviewId: string): Promise<void> {
    await apiClient.delete(`/reviews/${reviewId}`);
  }

  async getUserReviews(params?: GetReviewsParams): Promise<ReviewsResponse> {
    const response = await apiClient.get(`/reviews/user`, { params });
    return response.data;
  }

  async voteReview(reviewId: string, isHelpful: boolean): Promise<void> {
    await apiClient.post(`/reviews/${reviewId}/vote`, { is_helpful: isHelpful });
  }

  async removeVote(reviewId: string): Promise<void> {
    await apiClient.delete(`/reviews/${reviewId}/vote`);
  }

  // Admin APIs
  async getAdminReviews(params?: GetReviewsParams & { status?: string }): Promise<ReviewsResponse> {
    const response = await apiClient.get(`/admin/reviews`, { params });
    return response.data;
  }

  async updateReviewStatus(reviewId: string, status: 'approved' | 'hidden' | 'rejected'): Promise<void> {
    await apiClient.put(`/admin/reviews/${reviewId}/status`, { status });
  }

  async replyToReview(reviewId: string, reply: string): Promise<void> {
    await apiClient.post(`/admin/reviews/${reviewId}/reply`, { reply });
  }
}

export const reviewService = new ReviewService();
export default reviewService;
