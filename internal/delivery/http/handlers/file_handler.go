package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/usecases"
)

// FileHandler handles file upload operations
type FileHandler struct {
	fileUseCase usecases.FileUseCase
}

// NewFileHandler creates a new file handler
func NewFileHandler(fileUseCase usecases.FileUseCase) *FileHandler {
	return &FileHandler{
		fileUseCase: fileUseCase,
	}
}

// UploadImage handles general image file uploads for authenticated users
// @Summary Upload an image file
// @Description Upload an image file and return the URL (authentication required)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload/image [post]
func (h *FileHandler) UploadImage(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDStr := userID.(string)
	h.uploadImageHandler(c, entities.FileUploadTypeUser, &userIDStr)
}

// UploadImageAdmin handles admin-specific image uploads  
// @Summary Upload an image file (admin)
// @Description Upload an image file and return the URL (admin only)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Image file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/upload/image [post]
func (h *FileHandler) UploadImageAdmin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDStr := userID.(string)
	h.uploadImageHandler(c, entities.FileUploadTypeAdmin, &userIDStr)
}

// UploadImagePublic handles public image uploads (no authentication required)
// @Summary Upload an image file (public)
// @Description Upload an image file and return the URL (no authentication required)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /public/upload/image [post]
func (h *FileHandler) UploadImagePublic(c *gin.Context) {
	h.uploadImageHandler(c, entities.FileUploadTypePublic, nil)
}

// uploadImageHandler is the common handler logic for image uploads
func (h *FileHandler) uploadImageHandler(c *gin.Context, uploadType entities.FileUploadType, userID *string) {
	// Debug logging
	fmt.Printf("=== Upload Image Handler ===\n")
	fmt.Printf("Upload Type: %s\n", uploadType)
	fmt.Printf("User ID: %v\n", userID)
	fmt.Printf("Content Type: %s\n", c.GetHeader("Content-Type"))
	
	// Get the file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Printf("Error getting file from form: %v\n", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "No file provided",
		})
		return
	}
	defer file.Close()

	fmt.Printf("File received: %s, Size: %d bytes\n", header.Filename, header.Size)

	// Upload image using use case
	response, err := h.fileUseCase.UploadImage(c.Request.Context(), file, header, uploadType, userID)
	if err != nil {
		fmt.Printf("Upload failed: %v\n", err)
		// Handle specific error types
		switch {
		case err.Error() == "file size exceeds maximum allowed size":
			c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{
				Error: err.Error(),
			})
		case err.Error() == "file extension is not allowed" || err.Error() == "content type is not allowed":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to upload file: " + err.Error(),
			})
		}
		return
	}

	fmt.Printf("Upload successful: %+v\n", response)
	// Return success response
	c.JSON(http.StatusOK, response)
}

// UploadDocument handles document file uploads for authenticated users
// @Summary Upload a document file
// @Description Upload a document file and return the URL (authentication required)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Document file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /upload/document [post]
func (h *FileHandler) UploadDocument(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDStr := userID.(string)
	h.uploadDocumentHandler(c, entities.FileUploadTypeUser, &userIDStr)
}

// UploadDocumentAdmin handles admin-specific document uploads  
// @Summary Upload a document file (admin)
// @Description Upload a document file and return the URL (admin only)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Document file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /admin/upload/document [post]
func (h *FileHandler) UploadDocumentAdmin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userIDStr := userID.(string)
	h.uploadDocumentHandler(c, entities.FileUploadTypeAdmin, &userIDStr)
}

// UploadDocumentPublic handles public document uploads (no authentication required)
// @Summary Upload a document file (public)
// @Description Upload a document file and return the URL (no authentication required)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file to upload"
// @Success 200 {object} entities.FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 413 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /public/upload/document [post]
func (h *FileHandler) UploadDocumentPublic(c *gin.Context) {
	h.uploadDocumentHandler(c, entities.FileUploadTypePublic, nil)
}

// uploadDocumentHandler is the common handler logic for document uploads
func (h *FileHandler) uploadDocumentHandler(c *gin.Context, uploadType entities.FileUploadType, userID *string) {
	// Get the file from form data
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "No file provided",
		})
		return
	}
	defer file.Close()

	// Upload document using use case
	response, err := h.fileUseCase.UploadDocument(c.Request.Context(), file, header, uploadType, userID)
	if err != nil {
		// Handle specific error types
		switch {
		case err.Error() == "file size exceeds maximum allowed size":
			c.JSON(http.StatusRequestEntityTooLarge, ErrorResponse{
				Error: err.Error(),
			})
		case err.Error() == "file extension is not allowed" || err.Error() == "content type is not allowed":
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Failed to upload file: " + err.Error(),
			})
		}
		return
	}

	// Return success response
	c.JSON(http.StatusOK, response)
}

// DeleteFile handles file deletion
// @Summary Delete a file
// @Description Delete a file by ID (authentication required)
// @Tags files
// @Produce json
// @Security BearerAuth
// @Param id path string true "File ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /files/{id} [delete]
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "File ID is required",
		})
		return
	}

	// Delete file using use case
	if err := h.fileUseCase.DeleteFile(c.Request.Context(), fileID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

// GetFileUpload handles getting file upload info
// @Summary Get file upload information
// @Description Get file upload information by ID
// @Tags files
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} entities.FileUpload
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /files/{id} [get]
func (h *FileHandler) GetFileUpload(c *gin.Context) {
	fileID := c.Param("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "File ID is required",
		})
		return
	}

	// Get file upload using use case
	fileUpload, err := h.fileUseCase.GetFileUpload(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "File not found",
		})
		return
	}

	c.JSON(http.StatusOK, fileUpload)
}

// GetFileUploads handles getting list of file uploads
// @Summary Get list of file uploads
// @Description Get list of file uploads with filtering and pagination
// @Tags files
// @Produce json
// @Param upload_type query string false "Upload type (admin, user, public)"
// @Param category query string false "File category (images, documents)"
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} entities.FileUpload
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /files [get]
func (h *FileHandler) GetFileUploads(c *gin.Context) {
	// Parse query parameters
	uploadTypeStr := c.Query("upload_type")
	category := c.Query("category")
	
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// Convert upload type string to enum
	var uploadType entities.FileUploadType
	switch uploadTypeStr {
	case "admin":
		uploadType = entities.FileUploadTypeAdmin
	case "user":
		uploadType = entities.FileUploadTypeUser
	case "public":
		uploadType = entities.FileUploadTypePublic
	default:
		uploadType = "" // Empty means get all types
	}

	// Get file uploads using use case
	fileUploads, err := h.fileUseCase.GetFileUploads(c.Request.Context(), uploadType, category, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get file uploads: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": fileUploads,
		"pagination": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(fileUploads),
		},
	})
}
