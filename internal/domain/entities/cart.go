package entities

import (
	"time"

	"github.com/google/uuid"
)

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
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
	Price     float64   `json:"price" gorm:"not null"`
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
		existingItem.UpdatedAt = time.Now()
	} else {
		newItem := CartItem{
			ID:        uuid.New(),
			CartID:    c.ID,
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		c.Items = append(c.Items, newItem)
	}
	c.UpdatedAt = time.Now()
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(productID uuid.UUID) {
	for i, item := range c.Items {
		if item.ProductID == productID {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			c.UpdatedAt = time.Now()
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
			item.UpdatedAt = time.Now()
			c.UpdatedAt = time.Now()
		}
	}
}

// Clear removes all items from the cart
func (c *Cart) Clear() {
	c.Items = []CartItem{}
	c.UpdatedAt = time.Now()
}

// GetSubtotal calculates the subtotal for a cart item
func (ci *CartItem) GetSubtotal() float64 {
	return ci.Price * float64(ci.Quantity)
}
