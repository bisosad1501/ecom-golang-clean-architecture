package entities

import (
	"time"

	"github.com/google/uuid"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// PaymentStatus represents the payment status of an order
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Order represents an order in the system
type Order struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderNumber     string          `json:"order_number" gorm:"uniqueIndex;not null"`
	UserID          uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	User            User            `json:"user" gorm:"foreignKey:UserID"`
	Items           []OrderItem     `json:"items" gorm:"foreignKey:OrderID"`
	Status          OrderStatus     `json:"status" gorm:"default:'pending'"`
	PaymentStatus   PaymentStatus   `json:"payment_status" gorm:"default:'pending'"`
	Subtotal        float64         `json:"subtotal" gorm:"not null"`
	TaxAmount       float64         `json:"tax_amount" gorm:"default:0"`
	ShippingAmount  float64         `json:"shipping_amount" gorm:"default:0"`
	DiscountAmount  float64         `json:"discount_amount" gorm:"default:0"`
	Total           float64         `json:"total" gorm:"not null"`
	Currency        string          `json:"currency" gorm:"default:'USD'"`
	ShippingAddress *OrderAddress   `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  *OrderAddress   `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	Notes           string          `json:"notes" gorm:"type:text"`
	Payment         *Payment        `json:"payment" gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Order entity
func (Order) TableName() string {
	return "orders"
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;not null;index"`
	ProductID   uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product     Product   `json:"product" gorm:"foreignKey:ProductID"`
	ProductName string    `json:"product_name" gorm:"not null"`
	ProductSKU  string    `json:"product_sku" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Price       float64   `json:"price" gorm:"not null"`
	Total       float64   `json:"total" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for OrderItem entity
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderAddress represents an address for orders
type OrderAddress struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Company   string `json:"company"`
	Address1  string `json:"address1" validate:"required"`
	Address2  string `json:"address2"`
	City      string `json:"city" validate:"required"`
	State     string `json:"state" validate:"required"`
	ZipCode   string `json:"zip_code" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Phone     string `json:"phone"`
}

// GetFullName returns the full name from the address
func (a *OrderAddress) GetFullName() string {
	return a.FirstName + " " + a.LastName
}

// GetFullAddress returns the formatted full address
func (a *OrderAddress) GetFullAddress() string {
	address := a.Address1
	if a.Address2 != "" {
		address += ", " + a.Address2
	}
	address += ", " + a.City + ", " + a.State + " " + a.ZipCode + ", " + a.Country
	return address
}

// CanBeCancelled checks if the order can be cancelled
func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending || o.Status == OrderStatusConfirmed
}

// CanBeRefunded checks if the order can be refunded
func (o *Order) CanBeRefunded() bool {
	return o.PaymentStatus == PaymentStatusPaid && 
		   (o.Status == OrderStatusDelivered || o.Status == OrderStatusShipped)
}

// IsCompleted checks if the order is completed
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusDelivered
}

// IsPaid checks if the order is paid
func (o *Order) IsPaid() bool {
	return o.PaymentStatus == PaymentStatusPaid
}

// GetItemCount returns the total number of items in the order
func (o *Order) GetItemCount() int {
	count := 0
	for _, item := range o.Items {
		count += item.Quantity
	}
	return count
}

// CalculateTotal calculates the total amount of the order
func (o *Order) CalculateTotal() {
	o.Total = o.Subtotal + o.TaxAmount + o.ShippingAmount - o.DiscountAmount
}
