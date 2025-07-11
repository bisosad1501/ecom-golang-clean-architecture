package entities

import (
	"time"

	"github.com/google/uuid"
)

// StockReservationType represents the type of stock reservation
type StockReservationType string

const (
	ReservationTypeOrder    StockReservationType = "order"
	ReservationTypeCart     StockReservationType = "cart"
	ReservationTypePromotion StockReservationType = "promotion"
)

// StockReservationStatus represents the status of a stock reservation
type StockReservationStatus string

const (
	ReservationStatusActive    StockReservationStatus = "active"
	ReservationStatusConfirmed StockReservationStatus = "confirmed"
	ReservationStatusReleased  StockReservationStatus = "released"
	ReservationStatusExpired   StockReservationStatus = "expired"
)

// StockReservation represents a temporary stock reservation
type StockReservation struct {
	ID          uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID   uuid.UUID              `json:"product_id" gorm:"type:uuid;not null;index"`
	Product     Product                `json:"product" gorm:"foreignKey:ProductID"`
	OrderID     *uuid.UUID             `json:"order_id" gorm:"type:uuid;index"`
	Order       *Order                 `json:"order" gorm:"foreignKey:OrderID"`
	UserID      uuid.UUID              `json:"user_id" gorm:"type:uuid;not null;index"`
	User        User                   `json:"user" gorm:"foreignKey:UserID"`
	Quantity    int                    `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Type        StockReservationType   `json:"type" gorm:"not null"`
	Status      StockReservationStatus `json:"status" gorm:"default:'active'"`
	ReservedAt  time.Time              `json:"reserved_at" gorm:"autoCreateTime"`
	ExpiresAt   time.Time              `json:"expires_at" gorm:"not null;index"`
	ConfirmedAt *time.Time             `json:"confirmed_at"`
	ReleasedAt  *time.Time             `json:"released_at"`
	Notes       string                 `json:"notes" gorm:"type:text"`
	CreatedAt   time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for StockReservation entity
func (StockReservation) TableName() string {
	return "stock_reservations"
}

// IsExpired checks if the reservation has expired
func (sr *StockReservation) IsExpired() bool {
	return time.Now().After(sr.ExpiresAt)
}

// IsActive checks if the reservation is currently active
func (sr *StockReservation) IsActive() bool {
	return sr.Status == ReservationStatusActive && !sr.IsExpired()
}

// CanBeConfirmed checks if the reservation can be confirmed
func (sr *StockReservation) CanBeConfirmed() bool {
	return sr.Status == ReservationStatusActive && !sr.IsExpired()
}

// CanBeReleased checks if the reservation can be released
func (sr *StockReservation) CanBeReleased() bool {
	return sr.Status == ReservationStatusActive || sr.Status == ReservationStatusConfirmed
}

// Confirm confirms the reservation (converts to actual stock reduction)
func (sr *StockReservation) Confirm() {
	sr.Status = ReservationStatusConfirmed
	now := time.Now()
	sr.ConfirmedAt = &now
	sr.UpdatedAt = now
}

// Release releases the reservation
func (sr *StockReservation) Release() {
	sr.Status = ReservationStatusReleased
	now := time.Now()
	sr.ReleasedAt = &now
	sr.UpdatedAt = now
}

// MarkExpired marks the reservation as expired
func (sr *StockReservation) MarkExpired() {
	sr.Status = ReservationStatusExpired
	now := time.Now()
	sr.ReleasedAt = &now
	sr.UpdatedAt = now
}

// SetExpiration sets the expiration time
func (sr *StockReservation) SetExpiration(minutes int) {
	if minutes <= 0 {
		minutes = 30 // Default 30 minutes
	}
	sr.ExpiresAt = time.Now().Add(time.Duration(minutes) * time.Minute)
}

// ExtendExpiration extends the expiration time
func (sr *StockReservation) ExtendExpiration(minutes int) {
	if minutes <= 0 {
		minutes = 30 // Default 30 minutes
	}
	sr.ExpiresAt = sr.ExpiresAt.Add(time.Duration(minutes) * time.Minute)
	sr.UpdatedAt = time.Now()
}

// GetRemainingTime returns the remaining time before expiration
func (sr *StockReservation) GetRemainingTime() time.Duration {
	if sr.IsExpired() {
		return 0
	}
	return time.Until(sr.ExpiresAt)
}
