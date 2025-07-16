'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { oauthService } from '@/lib/services/oauth';

interface OAuthButtonsProps {
  onError?: (error: string) => void;
  className?: string;
}

export const OAuthButtons: React.FC<OAuthButtonsProps> = ({ 
  onError, 
  className = '' 
}) => {
  const [isGoogleLoading, setIsGoogleLoading] = useState(false);
  const [isFacebookLoading, setIsFacebookLoading] = useState(false);

  const handleGoogleLogin = async () => {
    try {
      setIsGoogleLoading(true);
      await oauthService.loginWithGoogle();
    } catch (error) {
      setIsGoogleLoading(false);
      const errorMessage = error instanceof Error ? error.message : 'Google login failed';
      onError?.(errorMessage);
    }
  };

  const handleFacebookLogin = async () => {
    try {
      setIsFacebookLoading(true);
      await oauthService.loginWithFacebook();
    } catch (error) {
      setIsFacebookLoading(false);
      const errorMessage = error instanceof Error ? error.message : 'Facebook login failed';
      onError?.(errorMessage);
    }
  };

  return (
    <div className={`grid grid-cols-2 gap-3 ${className}`}>
      {/* Google Login Button */}
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={handleGoogleLogin}
        disabled={isGoogleLoading || isFacebookLoading}
        className="h-9 bg-gray-800/90 border-gray-600/80 hover:border-[#FF9000] hover:bg-gray-700/90 text-white transition-all duration-300 transform hover:scale-[1.02] rounded-lg backdrop-blur-sm font-medium"
      >
        {isGoogleLoading ? (
          <>
            <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-[#FF9000] mr-1"></div>
            <span className="text-xs">Connecting...</span>
          </>
        ) : (
          <>
            <svg className="w-4 h-4 mr-1" viewBox="0 0 24 24">
              <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
              <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
              <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
              <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
            </svg>
            <span className="text-xs">Google</span>
          </>
        )}
      </Button>

      {/* Facebook Login Button */}
      <Button
        type="button"
        variant="outline"
        size="sm"
        onClick={handleFacebookLogin}
        disabled={isGoogleLoading || isFacebookLoading}
        className="h-9 bg-gray-800/90 border-gray-600/80 hover:border-[#FF9000] hover:bg-gray-700/90 text-white transition-all duration-300 transform hover:scale-[1.02] rounded-lg backdrop-blur-sm font-medium"
      >
        {isFacebookLoading ? (
          <>
            <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-[#FF9000] mr-1"></div>
            <span className="text-xs">Connecting...</span>
          </>
        ) : (
          <>
            <svg className="w-4 h-4 mr-1" fill="#1877F2" viewBox="0 0 24 24">
              <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
            </svg>
            <span className="text-xs">Facebook</span>
          </>
        )}
      </Button>
    </div>
  );
};
