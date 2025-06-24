package utils

import (
	"os"
	"strings"
)

// DeleteImageFile deletes an image file from the uploads directory
func DeleteImageFile(imageURL string) error {
	println("DEBUG: DeleteImageFile called with URL:", imageURL)
	
	if imageURL == "" {
		println("DEBUG: Image URL is empty, nothing to delete")
		return nil
	}

	// Extract the relative path from URL
	relativePath := ExtractRelativePathFromURL(imageURL)
	println("DEBUG: Extracted relative path:", relativePath)
	
	if relativePath == "" {
		println("DEBUG: No valid relative path extracted, not a local upload")
		return nil // Not a local upload or invalid URL
	}
	
	// Check if file exists
	if _, err := os.Stat(relativePath); os.IsNotExist(err) {
		println("DEBUG: File does not exist:", relativePath)
		return nil // File doesn't exist, nothing to delete
	}

	println("DEBUG: Attempting to delete file:", relativePath)
	// Delete the file
	err := os.Remove(relativePath)
	if err != nil {
		println("DEBUG: Failed to delete file:", err.Error())
	} else {
		println("DEBUG: Successfully deleted file:", relativePath)
	}
	return err
}

// ExtractRelativePathFromURL extracts the relative file path from an image URL
func ExtractRelativePathFromURL(imageURL string) string {
	if imageURL == "" {
		return ""
	}

	// Handle full URLs with domain
	if strings.Contains(imageURL, "://") {
		// Extract path part after domain
		parts := strings.Split(imageURL, "/")
		// Find index of "uploads"
		for i, part := range parts {
			if part == "uploads" && i+1 < len(parts) {
				return strings.Join(parts[i:], "/")
			}
		}
		return ""
	}

	// Handle relative URLs starting with /uploads/
	if strings.HasPrefix(imageURL, "/uploads/") {
		return strings.TrimPrefix(imageURL, "/")
	}

	return ""
}

// ExtractFilePathFromURL extracts the relative file path from an image URL (alias for ExtractRelativePathFromURL)
func ExtractFilePathFromURL(imageURL string) string {
	return ExtractRelativePathFromURL(imageURL)
}
