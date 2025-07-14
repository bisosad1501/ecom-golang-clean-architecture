package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)



// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    *uuid.UUID `json:"user_id" gorm:"type:uuid;index"` // Nullable for guest carts
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SessionID *string    `json:"session_id" gorm:"index"` // For guest users
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`

	// Calculated fields (stored for performance)
	Subtotal     float64 `json:"subtotal" gorm:"default:0"`
	TaxAmount    float64 `json:"tax_amount" gorm:"default:0"`
	ShippingAmount float64 `json:"shipping_amount" gorm:"default:0"`
	Total        float64 `json:"total" gorm:"default:0"`
	ItemCount    int     `json:"item_count" gorm:"default:0"`

	// Cart lifecycle
	Status    string     `json:"status" gorm:"default:'active'"`
	ExpiresAt *time.Time `json:"expires_at" gorm:"index"` // For cart abandonment

	// Metadata
	Currency string `json:"currency" gorm:"default:'USD'"`
	Notes    string `json:"notes" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Cart entity
func (Cart) TableName() string {
	return "carts"
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CartID    uuid.UUID `json:"cart_id" gorm:"type:uuid;not null;index"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null;index"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity" gorm:"not null" validate:"required,gt=0"`
	Price     float64   `json:"price" gorm:"not null" validate:"required,gte=0"`
	Total     float64   `json:"total" gorm:"not null" validate:"required,gte=0"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for CartItem entity
func (CartItem) TableName() string {
	return "cart_items"
}

// GetTotal calculates the total amount of the cart
func (c *Cart) GetTotal() float64 {
	total := 0.0
	for _, item := range c.Items {
		total += item.GetSubtotal()
	}
	return total
}

// GetItemCount returns the total number of items in the cart
func (c *Cart) GetItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

// UpdateCalculatedFields updates the calculated fields (subtotal, total, item_count)
func (c *Cart) UpdateCalculatedFields() {
	newSubtotal := c.GetTotal()
	newItemCount := c.GetItemCount()
	newTotal := newSubtotal + c.TaxAmount + c.ShippingAmount

	// Only update if values have changed to avoid unnecessary database writes
	if c.Subtotal != newSubtotal || c.ItemCount != newItemCount || c.Total != newTotal {
		c.Subtotal = newSubtotal
		c.Total = newTotal
		c.ItemCount = newItemCount
		c.UpdatedAt = time.Now()
	}
}

// UpdateCalculatedFieldsForce forces update of calculated fields regardless of changes
func (c *Cart) UpdateCalculatedFieldsForce() {
	c.Subtotal = c.GetTotal()
	// Tax and shipping will be calculated when converting to order
	// For now, cart total = subtotal (tax and shipping added at checkout)
	c.Total = c.Subtotal + c.TaxAmount + c.ShippingAmount
	c.ItemCount = c.GetItemCount()
	c.UpdatedAt = time.Now()
}

// SetExpiration sets cart expiration (default 7 days for logged in, 1 day for guest)
func (c *Cart) SetExpiration() {
	var hours int
	if c.UserID != nil {
		hours = 24 * 7 // 7 days for logged in users
	} else {
		hours = 24 // 1 day for guest users
	}
	expiry := time.Now().Add(time.Duration(hours) * time.Hour)
	c.ExpiresAt = &expiry
}

// IsExpired checks if the cart has expired
func (c *Cart) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

// MarkAsAbandoned marks the cart as abandoned
func (c *Cart) MarkAsAbandoned() {
	c.Status = "abandoned"
	c.UpdatedAt = time.Now()
}

// MarkAsConverted marks the cart as converted to order
func (c *Cart) MarkAsConverted() {
	c.Status = "converted"
	c.UpdatedAt = time.Now()
}

// SetTaxAmount sets the tax amount for the cart
func (c *Cart) SetTaxAmount(taxAmount float64) {
	if taxAmount < 0 {
		taxAmount = 0
	}
	c.TaxAmount = taxAmount
	c.UpdateCalculatedFields()
}

// SetShippingAmount sets the shipping amount for the cart
func (c *Cart) SetShippingAmount(shippingAmount float64) {
	if shippingAmount < 0 {
		shippingAmount = 0
	}
	c.ShippingAmount = shippingAmount
	c.UpdateCalculatedFields()
}

// Validate validates cart data
func (c *Cart) Validate() error {
	if c.UserID == nil && (c.SessionID == nil || *c.SessionID == "") {
		return fmt.Errorf("cart must have either user_id or session_id")
	}

	// Validate that cart doesn't have both UserID and SessionID
	if c.UserID != nil && c.SessionID != nil && *c.SessionID != "" {
		return fmt.Errorf("cart cannot have both user_id and session_id")
	}

	// Validate SessionID format if present
	if c.SessionID != nil && *c.SessionID != "" {
		if len(*c.SessionID) < 8 || len(*c.SessionID) > 128 {
			return fmt.Errorf("session_id must be between 8 and 128 characters")
		}
		// Basic format validation - should contain only alphanumeric and safe characters
		for _, char := range *c.SessionID {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				 (char >= '0' && char <= '9') || char == '-' || char == '_' || char == '.') {
				return fmt.Errorf("session_id contains invalid characters")
			}
		}
	}

	// Validate currency
	if c.Currency == "" {
		c.Currency = "USD"
	}

	// Validate financial fields
	if c.Subtotal < 0 {
		return fmt.Errorf("subtotal cannot be negative")
	}
	if c.Total < 0 {
		return fmt.Errorf("total cannot be negative")
	}
	if c.ItemCount < 0 {
		return fmt.Errorf("item_count cannot be negative")
	}

	// Validate items
	for i, item := range c.Items {
		if err := item.Validate(); err != nil {
			return fmt.Errorf("item %d validation failed: %w", i, err)
		}
	}

	// Validate calculated fields consistency
	expectedSubtotal := c.GetTotal()
	expectedItemCount := c.GetItemCount()
	expectedTotal := c.Subtotal + c.TaxAmount + c.ShippingAmount

	if c.Subtotal != expectedSubtotal {
		return fmt.Errorf("subtotal %.2f does not match calculated subtotal %.2f", c.Subtotal, expectedSubtotal)
	}
	if c.ItemCount != expectedItemCount {
		return fmt.Errorf("item_count %d does not match calculated item_count %d", c.ItemCount, expectedItemCount)
	}
	if c.Total != expectedTotal {
		return fmt.Errorf("total %.2f does not match calculated total %.2f (subtotal %.2f + tax %.2f + shipping %.2f)",
			c.Total, expectedTotal, c.Subtotal, c.TaxAmount, c.ShippingAmount)
	}
	validCurrencies := []string{"USD", "EUR", "GBP", "JPY", "VND"}
	isValidCurrency := false
	for _, currency := range validCurrencies {
		if c.Currency == currency {
			isValidCurrency = true
			break
		}
	}
	if !isValidCurrency {
		return fmt.Errorf("invalid currency: %s", c.Currency)
	}

	// Validate status
	if c.Status == "" {
		c.Status = "active"
	}
	validStatuses := []string{"active", "abandoned", "converted", "expired"}
	isValidStatus := false
	for _, status := range validStatuses {
		if c.Status == status {
			isValidStatus = true
			break
		}
	}
	if !isValidStatus {
		return fmt.Errorf("invalid status: %s", c.Status)
	}



	return nil
}

// IsGuest checks if this is a guest cart
func (c *Cart) IsGuest() bool {
	return c.UserID == nil && c.SessionID != nil
}

// IsUserCart checks if this is a user cart
func (c *Cart) IsUserCart() bool {
	return c.UserID != nil
}



// IsEmpty checks if the cart is empty
func (c *Cart) IsEmpty() bool {
	return len(c.Items) == 0
}

// HasItem checks if the cart contains a specific product
func (c *Cart) HasItem(productID uuid.UUID) bool {
	for _, item := range c.Items {
		if item.ProductID == productID {
			return true
		}
	}
	return false
}

// GetItem returns a cart item by product ID
func (c *Cart) GetItem(productID uuid.UUID) *CartItem {
	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			return &c.Items[i]
		}
	}
	return nil
}

// AddItem adds an item to the cart or updates quantity if it exists
func (c *Cart) AddItem(productID uuid.UUID, quantity int, price float64) {
	if existingItem := c.GetItem(productID); existingItem != nil {
		existingItem.Quantity += quantity
		existingItem.Price = price // Update price to current price
		existingItem.Total = float64(existingItem.Quantity) * existingItem.Price // Calculate total
		existingItem.UpdatedAt = time.Now()
	} else {
		newItem := CartItem{
			ID:        uuid.New(),
			CartID:    c.ID,
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
			Total:     float64(quantity) * price, // Calculate total for new item
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		c.Items = append(c.Items, newItem)
	}
	c.UpdateCalculatedFields()
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID uuid.UUID) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.UpdateCalculatedFields()
			break
		}
	}
}

// UpdateItemQuantity updates the quantity of a cart item
func (c *Cart) UpdateItemQuantity(productID uuid.UUID, quantity int) {
	if item := c.GetItem(productID); item != nil {
		if quantity <= 0 {
			c.RemoveItem(productID)
		} else {
			item.Quantity = quantity
			item.Total = float64(item.Quantity) * item.Price // Recalculate total
			item.UpdatedAt = time.Now()
			c.UpdateCalculatedFields()
		}
	}
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.UpdateCalculatedFields()
}

// GetTotal returns the stored total for a cart item
func (ci *CartItem) GetTotal() float64 {
	return ci.Total
}

// GetSubtotal is an alias for GetTotal for backward compatibility
func (ci *CartItem) GetSubtotal() float64 {
	return ci.GetTotal()
}

// CalculateTotal calculates and updates the total for a cart item
func (ci *CartItem) CalculateTotal() {
	ci.Total = ci.Price * float64(ci.Quantity)
	ci.UpdatedAt = time.Now()
}

// Validate validates cart item data
func (ci *CartItem) Validate() error {
	if ci.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	if ci.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}
	if ci.Total < 0 {
		return fmt.Errorf("total cannot be negative")
	}
	// Verify that total matches price * quantity
	expectedTotal := ci.Price * float64(ci.Quantity)
	if ci.Total != expectedTotal {
		return fmt.Errorf("total %.2f does not match price %.2f * quantity %d = %.2f",
			ci.Total, ci.Price, ci.Quantity, expectedTotal)
	}
	return nil
}
