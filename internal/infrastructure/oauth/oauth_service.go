package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"

	"ecom-golang-clean-architecture/internal/infrastructure/config"
)

// Service handles OAuth operations
type Service struct {
	config *config.OAuthConfig
}

// NewService creates a new OAuth service
func NewService(cfg *config.OAuthConfig) *Service {
	return &Service{
		config: cfg,
	}
}

// GetGoogleAuthURL returns the Google OAuth authorization URL
func (s *Service) GetGoogleAuthURL(state string) string {
	return s.config.Google.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetFacebookAuthURL returns the Facebook OAuth authorization URL
func (s *Service) GetFacebookAuthURL(state string) string {
	return s.config.Facebook.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeGoogleCode exchanges authorization code for user info
func (s *Service) ExchangeGoogleCode(ctx context.Context, code string) (*config.OAuthUserInfo, error) {
	token, err := s.config.Google.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange Google code: %w", err)
	}

	client := s.config.Google.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Google response: %w", err)
	}

	var googleUser config.GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Google user info: %w", err)
	}

	return googleUser.ToStandardUserInfo(), nil
}

// ExchangeFacebookCode exchanges authorization code for user info
func (s *Service) ExchangeFacebookCode(ctx context.Context, code string) (*config.OAuthUserInfo, error) {
	token, err := s.config.Facebook.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange Facebook code: %w", err)
	}

	client := s.config.Facebook.Client(ctx, token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		return nil, fmt.Errorf("failed to get Facebook user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read Facebook response: %w", err)
	}

	var facebookUser config.FacebookUserInfo
	if err := json.Unmarshal(body, &facebookUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Facebook user info: %w", err)
	}

	return facebookUser.ToStandardUserInfo(), nil
}

// ValidateState validates OAuth state parameter
func (s *Service) ValidateState(receivedState, expectedState string) bool {
	return receivedState == expectedState
}
