package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// OAuthConfig holds OAuth configuration for different providers
type OAuthConfig struct {
	Google   *oauth2.Config
	Facebook *oauth2.Config
}

// NewOAuthConfig creates a new OAuth configuration
func NewOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		Google: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		Facebook: &oauth2.Config{
			ClientID:     os.Getenv("FACEBOOK_APP_ID"),
			ClientSecret: os.Getenv("FACEBOOK_APP_SECRET"),
			RedirectURL:  os.Getenv("FACEBOOK_REDIRECT_URL"),
			Scopes: []string{
				"email",
				"public_profile",
			},
			Endpoint: facebook.Endpoint,
		},
	}
}

// GoogleUserInfo represents user info from Google
type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

// FacebookUserInfo represents user info from Facebook
type FacebookUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			Height       int    `json:"height"`
			IsSilhouette bool   `json:"is_silhouette"`
			URL          string `json:"url"`
			Width        int    `json:"width"`
		} `json:"data"`
	} `json:"picture"`
}

// OAuthProvider represents OAuth provider type
type OAuthProvider string

const (
	ProviderGoogle   OAuthProvider = "google"
	ProviderFacebook OAuthProvider = "facebook"
)

// OAuthUserInfo represents standardized user info from OAuth providers
type OAuthUserInfo struct {
	Provider     OAuthProvider `json:"provider"`
	ProviderID   string        `json:"provider_id"`
	Email        string        `json:"email"`
	Name         string        `json:"name"`
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Picture      string        `json:"picture"`
	Verified     bool          `json:"verified"`
}

// ToStandardUserInfo converts Google user info to standard format
func (g *GoogleUserInfo) ToStandardUserInfo() *OAuthUserInfo {
	return &OAuthUserInfo{
		Provider:   ProviderGoogle,
		ProviderID: g.ID,
		Email:      g.Email,
		Name:       g.Name,
		FirstName:  g.GivenName,
		LastName:   g.FamilyName,
		Picture:    g.Picture,
		Verified:   g.VerifiedEmail,
	}
}

// ToStandardUserInfo converts Facebook user info to standard format
func (f *FacebookUserInfo) ToStandardUserInfo() *OAuthUserInfo {
	return &OAuthUserInfo{
		Provider:   ProviderFacebook,
		ProviderID: f.ID,
		Email:      f.Email,
		Name:       f.Name,
		FirstName:  "", // Facebook doesn't provide separate first/last names in basic profile
		LastName:   "",
		Picture:    f.Picture.Data.URL,
		Verified:   true, // Facebook emails are generally verified
	}
}
