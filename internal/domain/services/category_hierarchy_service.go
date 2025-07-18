package services

import (
	"context"
	"sync"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// CategoryHierarchyService provides optimized category hierarchy operations
// Caches category relationships to avoid expensive recursive queries
type CategoryHierarchyService interface {
	// Get all descendant category IDs (including self)
	GetDescendantCategoryIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error)
	
	// Get all ancestor category IDs (including self)
	GetAncestorCategoryIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error)
	
	// Refresh cache (call when categories are modified)
	RefreshCache(ctx context.Context) error
	
	// Get category path (breadcrumb)
	GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error)
}

type categoryHierarchyService struct {
	categoryRepo repositories.CategoryRepository
	
	// Cache structures
	mu                sync.RWMutex
	descendantsCache  map[uuid.UUID][]uuid.UUID // categoryID -> descendant IDs
	ancestorsCache    map[uuid.UUID][]uuid.UUID // categoryID -> ancestor IDs
	categoriesCache   map[uuid.UUID]*entities.Category // categoryID -> category
	lastRefresh       time.Time
	cacheExpiry       time.Duration
}

// NewCategoryHierarchyService creates a new category hierarchy service
func NewCategoryHierarchyService(categoryRepo repositories.CategoryRepository) CategoryHierarchyService {
	service := &categoryHierarchyService{
		categoryRepo:     categoryRepo,
		descendantsCache: make(map[uuid.UUID][]uuid.UUID),
		ancestorsCache:   make(map[uuid.UUID][]uuid.UUID),
		categoriesCache:  make(map[uuid.UUID]*entities.Category),
		cacheExpiry:      15 * time.Minute, // Cache for 15 minutes
	}
	
	// Initialize cache
	go func() {
		ctx := context.Background()
		if err := service.RefreshCache(ctx); err != nil {
			// Log error but don't fail
			println("Warning: Failed to initialize category hierarchy cache:", err.Error())
		}
	}()
	
	return service
}

// GetDescendantCategoryIDs returns all descendant category IDs (optimized with cache)
func (s *categoryHierarchyService) GetDescendantCategoryIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	// Check if cache needs refresh
	if err := s.ensureCacheValid(ctx); err != nil {
		return nil, err
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return cached descendants
	if descendants, exists := s.descendantsCache[categoryID]; exists {
		return descendants, nil
	}

	// FIXED: If not in cache, fallback to database query instead of returning single category
	return s.queryDescendantsFromDB(ctx, categoryID)
}

// GetAncestorCategoryIDs returns all ancestor category IDs (optimized with cache)
func (s *categoryHierarchyService) GetAncestorCategoryIDs(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	// Check if cache needs refresh
	if err := s.ensureCacheValid(ctx); err != nil {
		return nil, err
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Return cached ancestors
	if ancestors, exists := s.ancestorsCache[categoryID]; exists {
		return ancestors, nil
	}

	// FIXED: If not in cache, fallback to database query instead of returning single category
	return s.queryAncestorsFromDB(ctx, categoryID)
}

// GetCategoryPath returns the full path from root to category (breadcrumb)
func (s *categoryHierarchyService) GetCategoryPath(ctx context.Context, categoryID uuid.UUID) ([]*entities.Category, error) {
	ancestorIDs, err := s.GetAncestorCategoryIDs(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var path []*entities.Category
	for _, id := range ancestorIDs {
		if category, exists := s.categoriesCache[id]; exists {
			path = append(path, category)
		}
	}
	
	return path, nil
}

// RefreshCache rebuilds the category hierarchy cache
func (s *categoryHierarchyService) RefreshCache(ctx context.Context) error {
	// Get all active categories (using List with large limit)
	categories, err := s.categoryRepo.List(ctx, 10000, 0)
	if err != nil {
		return err
	}
	
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Clear existing cache
	s.descendantsCache = make(map[uuid.UUID][]uuid.UUID)
	s.ancestorsCache = make(map[uuid.UUID][]uuid.UUID)
	s.categoriesCache = make(map[uuid.UUID]*entities.Category)
	
	// Build categories cache
	for _, category := range categories {
		if category.IsActive {
			s.categoriesCache[category.ID] = category
		}
	}
	
	// Build hierarchy maps
	for _, category := range categories {
		if !category.IsActive {
			continue
		}
		
		// Build descendants cache
		descendants := s.findDescendants(category.ID, categories)
		s.descendantsCache[category.ID] = descendants
		
		// Build ancestors cache
		ancestors := s.findAncestors(category.ID, categories)
		s.ancestorsCache[category.ID] = ancestors
	}
	
	s.lastRefresh = time.Now()
	return nil
}

// ensureCacheValid checks if cache needs refresh
func (s *categoryHierarchyService) ensureCacheValid(ctx context.Context) error {
	s.mu.RLock()
	needsRefresh := time.Since(s.lastRefresh) > s.cacheExpiry
	s.mu.RUnlock()
	
	if needsRefresh {
		return s.RefreshCache(ctx)
	}
	
	return nil
}

// findDescendants recursively finds all descendant category IDs
func (s *categoryHierarchyService) findDescendants(categoryID uuid.UUID, allCategories []*entities.Category) []uuid.UUID {
	var descendants []uuid.UUID
	descendants = append(descendants, categoryID) // Include self
	
	// Find direct children
	for _, category := range allCategories {
		if category.ParentID != nil && *category.ParentID == categoryID && category.IsActive {
			// Recursively get descendants of child
			childDescendants := s.findDescendants(category.ID, allCategories)
			descendants = append(descendants, childDescendants...)
		}
	}
	
	return descendants
}

// findAncestors recursively finds all ancestor category IDs
func (s *categoryHierarchyService) findAncestors(categoryID uuid.UUID, allCategories []*entities.Category) []uuid.UUID {
	var ancestors []uuid.UUID
	
	// Find the category
	var currentCategory *entities.Category
	for _, category := range allCategories {
		if category.ID == categoryID && category.IsActive {
			currentCategory = category
			break
		}
	}
	
	if currentCategory == nil {
		return ancestors
	}
	
	// Add current category
	ancestors = append(ancestors, currentCategory.ID)
	
	// If has parent, recursively get ancestors
	if currentCategory.ParentID != nil {
		parentAncestors := s.findAncestors(*currentCategory.ParentID, allCategories)
		ancestors = append(parentAncestors, ancestors...) // Prepend parent ancestors
	}
	
	return ancestors
}

// FIXED: Add database fallback methods for cache misses
// queryDescendantsFromDB queries descendants directly from database when cache misses
func (s *categoryHierarchyService) queryDescendantsFromDB(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	// Use the category repository's GetCategoryTree method
	descendants, err := s.categoryRepo.GetCategoryTree(ctx, categoryID)
	if err != nil {
		// If database query fails, return just the category itself as fallback
		return []uuid.UUID{categoryID}, nil
	}
	return descendants, nil
}

// queryAncestorsFromDB queries ancestors directly from database when cache misses
func (s *categoryHierarchyService) queryAncestorsFromDB(ctx context.Context, categoryID uuid.UUID) ([]uuid.UUID, error) {
	// Get all categories to build ancestor chain
	categories, err := s.categoryRepo.List(ctx, 10000, 0)
	if err != nil {
		// If database query fails, return just the category itself as fallback
		return []uuid.UUID{categoryID}, nil
	}

	// Find ancestors by traversing up the parent chain
	ancestors := s.findAncestors(categoryID, categories)
	return ancestors, nil
}
