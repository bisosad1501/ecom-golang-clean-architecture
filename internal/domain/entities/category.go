package entities

import (
	"time"

	"github.com/google/uuid"
)

// Category represents a product category
type Category struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string     `json:"name" gorm:"not null;index" validate:"required"`
	Description string     `json:"description" gorm:"type:text"`
	Slug        string     `json:"slug" gorm:"uniqueIndex;not null" validate:"required"`
	Image       string     `json:"image"`
	ParentID    *uuid.UUID `json:"parent_id" gorm:"type:uuid;index"`
	Parent      *Category  `json:"parent" gorm:"foreignKey:ParentID"`
	Children    []Category `json:"children" gorm:"foreignKey:ParentID"`
	// Products relationship removed - use ProductCategory many-to-many as single source of truth
	IsActive    bool       `json:"is_active" gorm:"default:true"`

	// SEO fields
	MetaTitle       string `json:"meta_title" gorm:"type:varchar(255)"`
	MetaDescription string `json:"meta_description" gorm:"type:text"`
	MetaKeywords    string `json:"meta_keywords" gorm:"type:text"`
	CanonicalURL    string `json:"canonical_url" gorm:"type:varchar(500)"`
	OGTitle         string `json:"og_title" gorm:"type:varchar(255)"`
	OGDescription   string `json:"og_description" gorm:"type:text"`
	OGImage         string `json:"og_image" gorm:"type:varchar(500)"`
	TwitterTitle    string `json:"twitter_title" gorm:"type:varchar(255)"`
	TwitterDescription string `json:"twitter_description" gorm:"type:text"`
	TwitterImage    string `json:"twitter_image" gorm:"type:varchar(500)"`
	SchemaMarkup    string `json:"schema_markup" gorm:"type:text"` // JSON string for structured data
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Category entity
func (Category) TableName() string {
	return "categories"
}

// IsRootCategory checks if the category is a root category
func (c *Category) IsRootCategory() bool {
	return c.ParentID == nil
}

// HasChildren checks if the category has children
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// GetLevel returns the level of the category in the hierarchy
func (c *Category) GetLevel() int {
	if c.IsRootCategory() {
		return 0
	}
	if c.Parent != nil {
		return c.Parent.GetLevel() + 1
	}
	return 1
}

// GetPath returns the full path of the category
func (c *Category) GetPath() string {
	if c.IsRootCategory() {
		return c.Name
	}
	if c.Parent != nil {
		return c.Parent.GetPath() + " > " + c.Name
	}
	return c.Name
}
