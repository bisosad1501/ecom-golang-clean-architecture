import { apiClient } from '../lib/api-client';

export interface OAuthURLResponse {
  url: string;
  state: string;
}

export interface OAuthUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  phone?: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface OAuthLoginResponse {
  token: string;
  user: OAuthUser;
}

export class OAuthService {
  /**
   * Get Google OAuth URL
   */
  async getGoogleAuthURL(): Promise<OAuthURLResponse> {
    const response = await apiClient.get<{
      message: string;
      data: OAuthURLResponse;
    }>('/auth/google/url');

    console.log('Google OAuth URL response:', response);
    console.log('Google OAuth URL response.data:', response.data);
    console.log('Google OAuth URL response.data.data:', (response.data as any)?.data);

    // The API client extracts data.data from the backend response
    // Backend returns: { message: string, data: OAuthURLResponse }
    // ApiClient extracts the 'data' field and returns it as response.data
    // So response.data should be OAuthURLResponse
    const oauthData = response.data as any;
    if (oauthData && typeof oauthData === 'object' && 'data' in oauthData) {
      return oauthData.data;
    }
    
    // If already extracted
    return response.data as unknown as OAuthURLResponse;
  }

  /**
   * Get Facebook OAuth URL
   */
  async getFacebookAuthURL(): Promise<OAuthURLResponse> {
    const response = await apiClient.get<{
      message: string;
      data: OAuthURLResponse;
    }>('/auth/facebook/url');

    console.log('Facebook OAuth URL response:', response);
    console.log('Facebook OAuth URL response.data:', response.data);
    console.log('Facebook OAuth URL response.data.data:', (response.data as any)?.data);

    // The API client extracts data.data from the backend response
    // Backend returns: { message: string, data: OAuthURLResponse }
    // ApiClient extracts the 'data' field and returns it as response.data
    // So response.data should be OAuthURLResponse
    const oauthData = response.data as any;
    if (oauthData && typeof oauthData === 'object' && 'data' in oauthData) {
      return oauthData.data;
    }
    
    // If already extracted
    return response.data as unknown as OAuthURLResponse;
  }

  /**
   * Handle OAuth callback (for both Google and Facebook)
   */
  async handleOAuthCallback(
    provider: 'google' | 'facebook',
    code: string,
    state: string
  ): Promise<OAuthLoginResponse> {
    const response = await apiClient.get<{
      message: string;
      data: OAuthLoginResponse;
    }>(`/auth/${provider}/callback?code=${code}&state=${state}`);


    return response.data.data;
  }

  /**
   * Initiate Google OAuth login (redirect)
   */
  initiateGoogleLogin(): void {
    window.location.href = `${process.env.REACT_APP_API_URL}/api/v1/auth/google/login`;
  }

  /**
   * Initiate Facebook OAuth login (redirect)
   */
  initiateFacebookLogin(): void {
    window.location.href = `${process.env.REACT_APP_API_URL}/api/v1/auth/facebook/login`;
  }

  /**
   * Get Google OAuth URL and redirect
   */
  async loginWithGoogle(): Promise<void> {
    try {
      const response = await this.getGoogleAuthURL();
      
      // Validate that we have a URL
      if (!response || !response.url) {
        throw new Error('Invalid OAuth URL response from server');
      }
      
      window.location.href = response.url;
    } catch (error) {
      console.error('Failed to get Google OAuth URL:', error);
      throw error;
    }
  }

  /**
   * Get Facebook OAuth URL and redirect
   */
  async loginWithFacebook(): Promise<void> {
    try {
      const response = await this.getFacebookAuthURL();
      
      // Validate that we have a URL
      if (!response || !response.url) {
        throw new Error('Invalid OAuth URL response from server');
      }
      
      window.location.href = response.url;
    } catch (error) {
      console.error('Failed to get Facebook OAuth URL:', error);
      throw error;
    }
  }
}

export const oauthService = new OAuthService();
