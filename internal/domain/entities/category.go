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
	Products    []Product  `json:"products" gorm:"foreignKey:CategoryID"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
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
