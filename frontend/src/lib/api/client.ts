import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { toast } from 'sonner';

// Create axios instance
const createApiClient = (): AxiosInstance => {
  const client = axios.create({
    baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1',
    timeout: 30000,
    headers: {
      'Content-Type': 'application/json',
    },
  });

  // Request interceptor to add auth token
  client.interceptors.request.use(
    (config) => {
      // Get token from localStorage or cookies
      const token = typeof window !== 'undefined' ? localStorage.getItem('auth_token') : null;
      
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      
      return config;
    },
    (error) => {
      return Promise.reject(error);
    }
  );

  // Response interceptor for error handling
  client.interceptors.response.use(
    (response: AxiosResponse) => {
      return response;
    },
    (error) => {
      // Handle common errors
      if (error.response) {
        const { status, data } = error.response;
        
        switch (status) {
          case 401:
            // Unauthorized - redirect to login
            if (typeof window !== 'undefined') {
              localStorage.removeItem('auth_token');
              window.location.href = '/auth/login';
            }
            break;
          case 403:
            toast.error('Bạn không có quyền thực hiện hành động này');
            break;
          case 404:
            toast.error('Không tìm thấy tài nguyên');
            break;
          case 422:
            // Validation errors
            if (data.errors) {
              const errorMessages = Object.values(data.errors).flat();
              errorMessages.forEach((message: any) => toast.error(message));
            } else {
              toast.error(data.message || 'Dữ liệu không hợp lệ');
            }
            break;
          case 429:
            toast.error('Quá nhiều yêu cầu. Vui lòng thử lại sau');
            break;
          case 500:
            toast.error('Lỗi máy chủ. Vui lòng thử lại sau');
            break;
          default:
            toast.error(data.message || 'Đã xảy ra lỗi');
        }
      } else if (error.request) {
        // Network error
        toast.error('Không thể kết nối đến máy chủ');
      } else {
        // Other errors
        toast.error('Đã xảy ra lỗi không xác định');
      }
      
      return Promise.reject(error);
    }
  );

  return client;
};

export const apiClient = createApiClient();

// Helper functions for common HTTP methods
export const api = {
  get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> =>
    apiClient.get(url, config),
    
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> =>
    apiClient.post(url, data, config),
    
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> =>
    apiClient.put(url, data, config),
    
  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> =>
    apiClient.patch(url, data, config),
    
  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<AxiosResponse<T>> =>
    apiClient.delete(url, config),
};

export default apiClient;
