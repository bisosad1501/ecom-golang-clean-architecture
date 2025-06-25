package handlers

import (
	"net/http"

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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to process Google callback",
			Details: err.Error(),
		})
		return
	}

	// For web applications, you might want to redirect to frontend with token
	// For API-only, return JSON response
	if c.Query("redirect") == "web" {
		// Redirect to frontend with token in URL fragment (more secure than query param)
		frontendURL := "http://localhost:3000/auth/callback"
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"#token="+response.Token)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Google authentication successful",
		Data:    response,
	})
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to process Facebook callback",
			Details: err.Error(),
		})
		return
	}

	// For web applications, you might want to redirect to frontend with token
	// For API-only, return JSON response
	if c.Query("redirect") == "web" {
		// Redirect to frontend with token in URL fragment (more secure than query param)
		frontendURL := "http://localhost:3000/auth/callback"
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"#token="+response.Token)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Facebook authentication successful",
		Data:    response,
	})
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
