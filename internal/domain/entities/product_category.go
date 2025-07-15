package entities

import (
	"time"

	"github.com/google/uuid"
)

// ProductCategory represents the many-to-many relationship between products and categories
type ProductCategory struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID  uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	CategoryID uuid.UUID `json:"category_id" gorm:"type:uuid;not null;index"`
	IsPrimary  bool      `json:"is_primary" gorm:"default:false"` // One category can be marked as primary
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Product  *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// TableName returns the table name for ProductCategory
func (ProductCategory) TableName() string {
	return "product_categories"
}

// ProductCategoryFilters represents filters for querying product categories
type ProductCategoryFilters struct {
	ProductID  *uuid.UUID
	CategoryID *uuid.UUID
	IsPrimary  *bool
	Limit      int
	Offset     int
}

// ProductWithCategories represents a product with its categories
type ProductWithCategories struct {
	*Product
	Categories       []*Category `json:"categories"`
	PrimaryCategory  *Category   `json:"primary_category"`
	CategoryIDs      []uuid.UUID `json:"category_ids"`
	PrimaryCategoryID *uuid.UUID `json:"primary_category_id"`
}

// CategoryWithProducts represents a category with its products
type CategoryWithProducts struct {
	*Category
	Products   []*Product `json:"products"`
	ProductIDs []uuid.UUID `json:"product_ids"`
}

// AddCategory adds a category to a product
func (p *ProductWithCategories) AddCategory(categoryID uuid.UUID, isPrimary bool) {
	// Check if category already exists
	for _, id := range p.CategoryIDs {
		if id == categoryID {
			return // Already exists
		}
	}

	p.CategoryIDs = append(p.CategoryIDs, categoryID)
	
	if isPrimary {
		p.PrimaryCategoryID = &categoryID
	}
}

// RemoveCategory removes a category from a product
func (p *ProductWithCategories) RemoveCategory(categoryID uuid.UUID) {
	// Remove from CategoryIDs
	for i, id := range p.CategoryIDs {
		if id == categoryID {
			p.CategoryIDs = append(p.CategoryIDs[:i], p.CategoryIDs[i+1:]...)
			break
		}
	}

	// Clear primary if it was the primary category
	if p.PrimaryCategoryID != nil && *p.PrimaryCategoryID == categoryID {
		p.PrimaryCategoryID = nil
	}
}

// SetPrimaryCategory sets a category as primary
func (p *ProductWithCategories) SetPrimaryCategory(categoryID uuid.UUID) error {
	// Check if category exists in the list
	found := false
	for _, id := range p.CategoryIDs {
		if id == categoryID {
			found = true
			break
		}
	}

	if !found {
		// Add the category first
		p.AddCategory(categoryID, true)
	} else {
		p.PrimaryCategoryID = &categoryID
	}

	return nil
}

// GetPrimaryCategory returns the primary category
func (p *ProductWithCategories) GetPrimaryCategory() *Category {
	return p.PrimaryCategory
}

// HasCategory checks if product belongs to a category
func (p *ProductWithCategories) HasCategory(categoryID uuid.UUID) bool {
	for _, id := range p.CategoryIDs {
		if id == categoryID {
			return true
		}
	}
	return false
}

// GetCategoryCount returns the number of categories
func (p *ProductWithCategories) GetCategoryCount() int {
	return len(p.CategoryIDs)
}
