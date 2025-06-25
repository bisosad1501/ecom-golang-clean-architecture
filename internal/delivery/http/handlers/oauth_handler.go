package handlers

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"ecom-golang-clean-architecture/internal/usecases"
)

// OAuthHandler handles OAuth-related HTTP requests
type OAuthHandler struct {
	oauthUseCase usecases.OAuthUseCase
}

// NewOAuthHandler creates a new OAuth handler
func NewOAuthHandler(oauthUseCase usecases.OAuthUseCase) *OAuthHandler {
	return &OAuthHandler{
		oauthUseCase: oauthUseCase,
	}
}

// GetGoogleAuthURL generates Google OAuth URL
func (h *OAuthHandler) GetGoogleAuthURL(c *gin.Context) {
	response, err := h.oauthUseCase.GetGoogleAuthURL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate Google auth URL",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Google auth URL generated successfully",
		Data:    response,
	})
}

// GetFacebookAuthURL generates Facebook OAuth URL
func (h *OAuthHandler) GetFacebookAuthURL(c *gin.Context) {
	response, err := h.oauthUseCase.GetFacebookAuthURL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to generate Facebook auth URL",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Facebook auth URL generated successfully",
		Data:    response,
	})
}

// GoogleCallback handles Google OAuth callback
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Authorization code is required",
		})
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "State parameter is required",
		})
		return
	}

	req := &usecases.OAuthCallbackRequest{
		Code:  code,
		State: state,
	}

	response, err := h.oauthUseCase.HandleGoogleCallback(c.Request.Context(), req)
	if err != nil {
		// Redirect to frontend with error (URL encode the error message)
		frontendURL := "http://localhost:3000/auth/google/callback?error=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, frontendURL)
		return
	}

	// Always redirect to frontend callback page with success token
	// Use URL fragment to pass token (more secure than query params)
	frontendURL := "http://localhost:3000/auth/google/callback"
	redirectURL := frontendURL + "#token=" + response.Token + "&success=true"
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// FacebookCallback handles Facebook OAuth callback
func (h *OAuthHandler) FacebookCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Authorization code is required",
		})
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "State parameter is required",
		})
		return
	}

	req := &usecases.OAuthCallbackRequest{
		Code:  code,
		State: state,
	}

	response, err := h.oauthUseCase.HandleFacebookCallback(c.Request.Context(), req)
	if err != nil {
		// Redirect to frontend with error (URL encode the error message)
		frontendURL := "http://localhost:3000/auth/facebook/callback?error=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, frontendURL)
		return
	}

	// Always redirect to frontend callback page with success token
	// Use URL fragment to pass token (more secure than query params)
	frontendURL := "http://localhost:3000/auth/facebook/callback"
	redirectURL := frontendURL + "#token=" + response.Token + "&success=true"
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GoogleLogin initiates Google OAuth flow (alternative endpoint)
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	response, err := h.oauthUseCase.GetGoogleAuthURL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to initiate Google login",
			Details: err.Error(),
		})
		return
	}

	// Redirect directly to Google OAuth
	c.Redirect(http.StatusTemporaryRedirect, response.URL)
}

// FacebookLogin initiates Facebook OAuth flow (alternative endpoint)
func (h *OAuthHandler) FacebookLogin(c *gin.Context) {
	response, err := h.oauthUseCase.GetFacebookAuthURL(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to initiate Facebook login",
			Details: err.Error(),
		})
		return
	}

	// Redirect directly to Facebook OAuth
	c.Redirect(http.StatusTemporaryRedirect, response.URL)
}
