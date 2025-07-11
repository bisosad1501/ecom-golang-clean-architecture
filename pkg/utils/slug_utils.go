package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateSlug creates a URL-friendly slug from a name
func GenerateSlug(name string) string {
	if name == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing hyphens
	slug = strings.Trim(slug, "-")

	// Ensure it's not empty
	if slug == "" {
		slug = fmt.Sprintf("product-%d", time.Now().Unix())
	}

	return slug
}

// GenerateUniqueSlug creates a unique slug by appending a suffix if needed
func GenerateUniqueSlug(baseSlug string, existingSlugs []string) string {
	if baseSlug == "" {
		baseSlug = fmt.Sprintf("product-%d", time.Now().Unix())
	}

	// Check if base slug is already unique
	if !contains(existingSlugs, baseSlug) {
		return baseSlug
	}

	// Try appending numbers until we find a unique slug
	for i := 1; i <= 1000; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		if !contains(existingSlugs, candidateSlug) {
			return candidateSlug
		}
	}

	// Fallback to timestamp-based slug
	return fmt.Sprintf("%s-%d", baseSlug, time.Now().Unix())
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ValidateSlug validates that a slug meets requirements
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug cannot be empty")
	}

	if len(slug) > 255 {
		return fmt.Errorf("slug cannot be longer than 255 characters")
	}

	// Check if slug contains only valid characters
	reg := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !reg.MatchString(slug) {
		return fmt.Errorf("slug can only contain lowercase letters, numbers, and hyphens")
	}

	// Check if slug starts or ends with hyphen
	if strings.HasPrefix(slug, "-") || strings.HasSuffix(slug, "-") {
		return fmt.Errorf("slug cannot start or end with a hyphen")
	}

	// Check for consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug cannot contain consecutive hyphens")
	}

	return nil
}
