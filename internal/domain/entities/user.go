package entities

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents user account status
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusInactive  UserStatus = "inactive"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusPending   UserStatus = "pending"
)

// UserRole represents the role of a user
type UserRole string

const (
	UserRoleCustomer UserRole = "customer"
	UserRoleAdmin    UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
	FirstName string    `json:"first_name" gorm:"not null" validate:"required"`
	LastName  string    `json:"last_name" gorm:"not null" validate:"required"`
	Phone     string     `json:"phone" gorm:"index"`
	Role      UserRole   `json:"role" gorm:"default:'customer'" validate:"required"`
	Status    UserStatus `json:"status" gorm:"default:'active'" validate:"required"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relationships
	Profile   *UserProfile `json:"profile,omitempty" gorm:"foreignKey:UserID"`
	Addresses []Address    `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Wishlist  []Product    `json:"wishlist,omitempty" gorm:"many2many:user_wishlists;"`
}

// TableName returns the table name for User entity
func (User) TableName() string {
	return "users"
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsModerator checks if the user is a moderator
func (u *User) IsModerator() bool {
	return u.Role == UserRoleModerator
}

// CanManageProducts checks if the user can manage products
func (u *User) CanManageProducts() bool {
	return u.Role == UserRoleAdmin || u.Role == UserRoleModerator
}

// UserProfile represents additional user profile information
type UserProfile struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Avatar      string    `json:"avatar"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Country     string    `json:"country"`
	ZipCode     string    `json:"zip_code"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for UserProfile entity
func (UserProfile) TableName() string {
	return "user_profiles"
}
