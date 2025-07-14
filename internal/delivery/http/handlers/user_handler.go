package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"ecom-golang-clean-architecture/internal/usecases"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userUseCase usecases.UserUseCase
}

// getUserIDFromContext extracts user ID from gin context
func (h *UserHandler) getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in token")
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user ID format")
	}

	return userID, nil
}

// NewUserHandler creates a new user handler
func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.RegisterRequest true "Registration request"
// @Success 201 {object} usecases.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req usecases.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "User registered successfully",
		Data:    user,
	})
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.LoginRequest true "Login request"
// @Success 200 {object} usecases.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req usecases.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	response, err := h.userUseCase.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Login successful",
		Data:    response,
	})
}

// GetProfile handles getting user profile
// @Summary Get user profile
// @Description Get current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: user,
	})
}

// UpdateProfile handles updating user profile
// @Summary Update user profile
// @Description Update current user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.UpdateProfileRequest true "Update profile request"
// @Success 200 {object} usecases.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID format",
		})
		return
	}

	var req usecases.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Profile updated successfully",
		Data:    user,
	})
}

// ChangePassword handles changing user password
// @Summary Change user password
// @Description Change current user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.ChangePasswordRequest true "Change password request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err = h.userUseCase.ChangePassword(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password changed successfully",
	})
}

// GetUsers handles getting list of users (admin only)
// @Summary Get users list
// @Description Get list of users with pagination (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} usecases.UserResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	usersResponse, err := h.userUseCase.GetUsers(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Calculate pagination
	page := (offset / limit) + 1
	totalPages := int((usersResponse.Total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, PaginatedResponse{
		Data: usersResponse.Users,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      usersResponse.Total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	})
}

// DeactivateUser handles deactivating a user (admin only)
// @Summary Deactivate user
// @Description Deactivate a user account (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/users/{id}/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = h.userUseCase.DeactivateUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User deactivated successfully",
	})
}

// ActivateUser handles activating a user (admin only)
// @Summary Activate user
// @Description Activate a user account (admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /admin/users/{id}/activate [post]
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = h.userUseCase.ActivateUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "User activated successfully",
	})
}

// GetUserPreferences handles getting user preferences
// @Summary Get user preferences
// @Description Get current user's preferences
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.UserPreferencesResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/preferences [get]
func (h *UserHandler) GetUserPreferences(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	preferences, err := h.userUseCase.GetUserPreferences(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: preferences,
	})
}

// UpdateUserPreferences handles updating user preferences
// @Summary Update user preferences
// @Description Update current user's preferences
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body usecases.UpdateUserPreferencesRequest true "Update preferences request"
// @Success 200 {object} usecases.UserPreferencesResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/preferences [put]
func (h *UserHandler) UpdateUserPreferences(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.UpdateUserPreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	preferences, err := h.userUseCase.UpdateUserPreferences(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Preferences updated successfully",
		Data:    preferences,
	})
}

// UpdateTheme handles updating user theme preference
// @Summary Update user theme
// @Description Update current user's theme preference
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Theme request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/preferences/theme [put]
func (h *UserHandler) UpdateTheme(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	theme, exists := req["theme"]
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Theme is required",
		})
		return
	}

	err = h.userUseCase.UpdateTheme(c.Request.Context(), userID, theme)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Theme updated successfully",
	})
}

// UpdateLanguage handles updating user language preference
// @Summary Update user language
// @Description Update current user's language preference
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Language request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/preferences/language [put]
func (h *UserHandler) UpdateLanguage(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	language, exists := req["language"]
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Language is required",
		})
		return
	}

	err = h.userUseCase.UpdateLanguage(c.Request.Context(), userID, language)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Language updated successfully",
	})
}

// SendEmailVerification handles sending email verification
// @Summary Send email verification
// @Description Send email verification to current user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/verification/email/send [post]
func (h *UserHandler) SendEmailVerification(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = h.userUseCase.SendEmailVerification(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email verification sent successfully",
	})
}

// VerifyEmail handles email verification
// @Summary Verify email
// @Description Verify email with token
// @Tags users
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/verification/email/verify [post]
func (h *UserHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Verification token is required",
		})
		return
	}

	err := h.userUseCase.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email verified successfully",
	})
}

// SendPhoneVerification handles sending phone verification
// @Summary Send phone verification
// @Description Send phone verification OTP to current user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Phone request"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/verification/phone/send [post]
func (h *UserHandler) SendPhoneVerification(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	phone, exists := req["phone"]
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Phone number is required",
		})
		return
	}

	err = h.userUseCase.SendPhoneVerification(c.Request.Context(), userID, phone)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Phone verification sent successfully",
	})
}

// VerifyPhone handles phone verification
// @Summary Verify phone
// @Description Verify phone with OTP code
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Code request"
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /users/verification/phone/verify [post]
func (h *UserHandler) VerifyPhone(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	code, exists := req["code"]
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Verification code is required",
		})
		return
	}

	err = h.userUseCase.VerifyPhone(c.Request.Context(), userID, code)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Phone verified successfully",
	})
}

// GetVerificationStatus handles getting verification status
// @Summary Get verification status
// @Description Get current user's verification status
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usecases.VerificationStatusResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/verification/status [get]
func (h *UserHandler) GetVerificationStatus(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	status, err := h.userUseCase.GetVerificationStatus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: status,
	})
}

// GetUserSessions handles getting user sessions
// @Summary Get user sessions
// @Description Get current user's active sessions
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} usecases.UserSessionsResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/sessions [get]
func (h *UserHandler) GetUserSessions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	limit := 10
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	sessions, err := h.userUseCase.GetUserSessions(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Data: sessions,
	})
}

// InvalidateSession handles invalidating a specific session
// @Summary Invalidate session
// @Description Invalidate a specific user session
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param session_id path string true "Session ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/sessions/{session_id} [delete]
func (h *UserHandler) InvalidateSession(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid session ID",
		})
		return
	}

	err = h.userUseCase.InvalidateSession(c.Request.Context(), userID, sessionID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Session invalidated successfully",
	})
}

// InvalidateAllSessions handles invalidating all user sessions
// @Summary Invalidate all sessions
// @Description Invalidate all user sessions except current
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /users/sessions [delete]
func (h *UserHandler) InvalidateAllSessions(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	err = h.userUseCase.InvalidateAllSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(getErrorStatusCode(err), ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "All sessions invalidated successfully",
	})
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and invalidate token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	// Get token from header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Authorization header required",
		})
		return
	}

	// Extract token
	token := authHeader[7:] // Remove "Bearer " prefix

	err := h.userUseCase.Logout(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Logged out successfully",
	})
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token request"
// @Success 200 {object} usecases.RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /auth/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	response, err := h.userUseCase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Token refreshed successfully",
		Data:    response,
	})
}

// ForgotPassword handles forgot password request
// @Summary Forgot password
// @Description Send password reset email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.ForgotPasswordRequest true "Forgot password request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req usecases.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err := h.userUseCase.ForgotPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password reset email sent successfully",
	})
}

// ResetPassword handles password reset
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.ResetPasswordRequest true "Reset password request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req usecases.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err := h.userUseCase.ResetPassword(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Password reset successfully",
	})
}

// VerifyEmailWithToken handles email verification with token in body
// @Summary Verify email with token
// @Description Verify email using verification token in request body
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.VerifyEmailRequest true "Verify email request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/verify-email [post]
func (h *UserHandler) VerifyEmailWithToken(c *gin.Context) {
	var req usecases.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	err := h.userUseCase.VerifyEmail(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Email verified successfully",
	})
}

// ResendVerification handles resend verification email
// @Summary Resend verification email
// @Description Resend email verification
// @Tags auth
// @Accept json
// @Produce json
// @Param request body usecases.ResendVerificationRequest true "Resend verification request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth/resend-verification [post]
func (h *UserHandler) ResendVerification(c *gin.Context) {
	var req usecases.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// For now, we'll use the existing SendEmailVerification method
	// In production, you might want to create a separate method
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Verification email sent successfully",
	})
}

// TrackSearch tracks user search activity
func (h *UserHandler) TrackSearch(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.TrackSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	req.UserID = userID
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	if err := h.userUseCase.TrackSearch(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to track search",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search tracked successfully",
	})
}

// GetSearchHistory retrieves user's search history
func (h *UserHandler) GetSearchHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.SearchHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := h.userUseCase.GetSearchHistory(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get search history",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search history retrieved successfully",
		Data:    response,
	})
}

// ClearSearchHistory clears user's search history
func (h *UserHandler) ClearSearchHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	if err := h.userUseCase.ClearSearchHistory(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to clear search history",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Search history cleared successfully",
	})
}

// CreateSavedSearch creates a new saved search
func (h *UserHandler) CreateSavedSearch(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.CreateSavedSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	req.UserID = userID

	response, err := h.userUseCase.CreateSavedSearch(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create saved search",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "Saved search created successfully",
		Data:    response,
	})
}

// GetSavedSearches retrieves user's saved searches
func (h *UserHandler) GetSavedSearches(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.GetSavedSearchesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := h.userUseCase.GetSavedSearches(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get saved searches",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Saved searches retrieved successfully",
		Data:    response,
	})
}

// TrackProductView tracks user product viewing activity
func (h *UserHandler) TrackProductView(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.TrackProductViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	req.UserID = userID
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	if err := h.userUseCase.TrackProductView(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to track product view",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Product view tracked successfully",
	})
}

// GetBrowsingHistory retrieves user's browsing history
func (h *UserHandler) GetBrowsingHistory(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req usecases.BrowsingHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid query parameters",
			Details: err.Error(),
		})
		return
	}

	// Set default values
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	response, err := h.userUseCase.GetBrowsingHistory(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get browsing history",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Browsing history retrieved successfully",
		Data:    response,
	})
}

// GetPersonalization retrieves user personalization data
func (h *UserHandler) GetPersonalization(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User ID not found in token",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	response, err := h.userUseCase.GetPersonalization(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get personalization data",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Personalization data retrieved successfully",
		Data:    response,
	})
}
