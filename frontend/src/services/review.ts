import axios from 'axios';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

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
  private getAuthHeaders() {
    const token = localStorage.getItem('token');
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  // Public APIs (no auth required)
  async getProductReviews(productId: string, params?: GetReviewsParams): Promise<ReviewsResponse> {
    const response = await axios.get(`${API_BASE_URL}/public/reviews/product/${productId}`, {
      params,
    });
    return response.data.data;
  }

  async getProductRatingSummary(productId: string): Promise<ProductRatingSummary> {
    const response = await axios.get(`${API_BASE_URL}/public/reviews/product/${productId}/summary`);
    return response.data.data;
  }

  // Protected APIs (auth required)
  async createReview(data: CreateReviewRequest): Promise<Review> {
    const response = await axios.post(`${API_BASE_URL}/reviews`, data, {
      headers: this.getAuthHeaders(),
    });
    return response.data.data;
  }

  async updateReview(reviewId: string, data: UpdateReviewRequest): Promise<Review> {
    const response = await axios.put(`${API_BASE_URL}/reviews/${reviewId}`, data, {
      headers: this.getAuthHeaders(),
    });
    return response.data.data;
  }

  async deleteReview(reviewId: string): Promise<void> {
    await axios.delete(`${API_BASE_URL}/reviews/${reviewId}`, {
      headers: this.getAuthHeaders(),
    });
  }

  async getUserReviews(params?: GetReviewsParams): Promise<ReviewsResponse> {
    const response = await axios.get(`${API_BASE_URL}/reviews/user`, {
      headers: this.getAuthHeaders(),
      params,
    });
    return response.data.data;
  }

  async voteReview(reviewId: string, isHelpful: boolean): Promise<void> {
    await axios.post(`${API_BASE_URL}/reviews/${reviewId}/vote`, 
      { is_helpful: isHelpful },
      { headers: this.getAuthHeaders() }
    );
  }

  async removeVote(reviewId: string): Promise<void> {
    await axios.delete(`${API_BASE_URL}/reviews/${reviewId}/vote`, {
      headers: this.getAuthHeaders(),
    });
  }

  // Admin APIs
  async getAdminReviews(params?: GetReviewsParams & { status?: string }): Promise<ReviewsResponse> {
    const response = await axios.get(`${API_BASE_URL}/admin/reviews`, {
      headers: this.getAuthHeaders(),
      params,
    });
    return response.data.data;
  }

  async updateReviewStatus(reviewId: string, status: 'approved' | 'hidden' | 'rejected'): Promise<void> {
    await axios.put(`${API_BASE_URL}/admin/reviews/${reviewId}/status`, 
      { status },
      { headers: this.getAuthHeaders() }
    );
  }

  async replyToReview(reviewId: string, reply: string): Promise<void> {
    await axios.post(`${API_BASE_URL}/admin/reviews/${reviewId}/reply`, 
      { reply },
      { headers: this.getAuthHeaders() }
    );
  }
}

export const reviewService = new ReviewService();
export default reviewService;
