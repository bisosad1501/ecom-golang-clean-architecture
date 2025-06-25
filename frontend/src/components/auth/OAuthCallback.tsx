'use client';

import React, { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { oauthService } from '../../services/oauth';
import { useAuthStore } from '../../store/auth';

interface OAuthCallbackProps {
  provider: 'google' | 'facebook';
}

export const OAuthCallback: React.FC<OAuthCallbackProps> = ({ provider }) => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { setAuthData } = useAuthStore();
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [error, setError] = useState<string>('');

  useEffect(() => {
    const handleCallback = async () => {
      try {
        console.log('=== OAuth Callback Debug ===');
        console.log('Current URL:', window.location.href);
        console.log('Hash:', window.location.hash);
        console.log('Search params:', window.location.search);
        
        // Check for OAuth errors in query params first
        const errorParam = searchParams.get('error');
        if (errorParam) {
          throw new Error(`OAuth error: ${decodeURIComponent(errorParam)}`);
        }

        // Check for token in URL fragment (set by backend redirect) - PRIORITY
        const fragment = window.location.hash.substring(1);
        const fragmentParams = new URLSearchParams(fragment);
        const token = fragmentParams.get('token');
        const success = fragmentParams.get('success');

        console.log('Fragment params:', { token: token ? 'present' : 'missing', success });

        if (token && success === 'true') {
          // Token received from backend redirect
          console.log('Processing token from fragment...');
          setStatus('loading');
          
          // We need to get user info with this token
          // Set token in API client first
          const { apiClient } = await import('../../lib/api-client');
          apiClient.setToken(token);
          
          // Get user profile to complete login (correct endpoint)
          const userResponse = await apiClient.get('/users/profile');
          
          // Store authentication data (userResponse.data contains the actual user data)
          setAuthData(token, userResponse.data as any);
          
          setStatus('success');
          
          // Clear URL fragment
          window.history.replaceState(null, '', window.location.pathname);
          
          // Redirect to homepage or intended page
          const redirectTo = sessionStorage.getItem('oauth_redirect') || '/';
          sessionStorage.removeItem('oauth_redirect');
          
          setTimeout(() => {
            router.push(redirectTo);
          }, 1500);
          
          return; // Exit early, don't try fallback
        }

        // Fallback: Check for error in fragment
        const fragmentError = fragmentParams.get('error');
        if (fragmentError) {
          throw new Error(`OAuth error: ${decodeURIComponent(fragmentError)}`);
        }

        // Fallback: handle OAuth callback via API (old method) - ONLY if no fragment data
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        
        console.log('Query params:', { code: code ? 'present' : 'missing', state: state ? 'present' : 'missing' });

        if (code && state) {
          console.log('No token in fragment, falling back to API callback...');
          setStatus('loading');

          // Handle OAuth callback
          const response = await oauthService.handleOAuthCallback(provider, code, state);

          // Store authentication data
          setAuthData(response.token, response.user as any);

          setStatus('success');

          // Redirect to homepage or intended page
          const redirectTo = sessionStorage.getItem('oauth_redirect') || '/';
          sessionStorage.removeItem('oauth_redirect');
          
          setTimeout(() => {
            router.push(redirectTo);
          }, 1500);
        } else {
          // Check if user is already authenticated before throwing error
          const { useAuthStore } = await import('../../store/auth');
          const authStore = useAuthStore.getState();
          
          console.log('Auth state check:', authStore);
          
          if (authStore?.isAuthenticated) {
            console.log('User already authenticated, redirecting to home...');
            // User is already authenticated, just redirect
            setTimeout(() => {
              router.push('/');
            }, 500);
            return;
          }
          
          console.warn('No authentication data found in URL');
          throw new Error('Missing authentication data - no token in fragment and no code/state in query');
        }

      } catch (error) {
        console.error(`${provider} OAuth callback error:`, error);
        setError(error instanceof Error ? error.message : `${provider} authentication failed`);
        setStatus('error');

        // Redirect to login page after error
        setTimeout(() => {
          router.push('/auth/login');
        }, 3000);
      }
    };

    handleCallback();
  }, [searchParams, router, setAuthData, provider]);

  const getProviderName = () => {
    return provider.charAt(0).toUpperCase() + provider.slice(1);
  };

  const getProviderColor = () => {
    return provider === 'google' ? 'text-blue-600' : 'text-blue-800';
  };

  const getProviderIcon = () => {
    if (provider === 'google') {
      return (
        <svg className="w-8 h-8" viewBox="0 0 24 24">
          <path
            fill="#4285F4"
            d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
          />
          <path
            fill="#34A853"
            d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
          />
          <path
            fill="#FBBC05"
            d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
          />
          <path
            fill="#EA4335"
            d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
          />
        </svg>
      );
    } else {
      return (
        <svg className="w-8 h-8 text-[#1877F2]" fill="currentColor" viewBox="0 0 24 24">
          <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
        </svg>
      );
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <div className="flex justify-center mb-4">
            {getProviderIcon()}
          </div>
          
          {status === 'loading' && (
            <>
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900 mx-auto mb-4"></div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Completing {getProviderName()} Sign In
              </h2>
              <p className="text-gray-600">
                Please wait while we verify your account...
              </p>
            </>
          )}

          {status === 'success' && (
            <>
              <div className="flex justify-center mb-4">
                <div className="rounded-full bg-green-100 p-2">
                  <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Welcome!
              </h2>
              <p className="text-gray-600">
                Successfully signed in with {getProviderName()}. Redirecting you to your dashboard...
              </p>
            </>
          )}

          {status === 'error' && (
            <>
              <div className="flex justify-center mb-4">
                <div className="rounded-full bg-red-100 p-2">
                  <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </div>
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Authentication Failed
              </h2>
              <p className="text-gray-600 mb-4">
                {error || `Failed to sign in with ${getProviderName()}`}
              </p>
              <p className="text-sm text-gray-500">
                Redirecting you back to the login page...
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
};
