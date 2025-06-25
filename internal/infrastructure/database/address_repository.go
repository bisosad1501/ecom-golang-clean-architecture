package database

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository creates a new address repository
func NewAddressRepository(db *gorm.DB) repositories.AddressRepository {
	return &addressRepository{db: db}
}

// Create creates a new address
func (r *addressRepository) Create(ctx context.Context, address *entities.Address) error {
	return r.db.WithContext(ctx).Create(address).Error
}

// GetByID gets an address by ID
func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Address, error) {
	var address entities.Address
	err := r.db.WithContext(ctx).First(&address, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// GetByUser gets all addresses for a user
func (r *addressRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error) {
	var addresses []*entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// GetDefaultByUser gets the default address for a user
func (r *addressRepository) GetDefaultByUser(ctx context.Context, userID uuid.UUID) (*entities.Address, error) {
	var address entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_default = ?", userID, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Update updates an address
func (r *addressRepository) Update(ctx context.Context, address *entities.Address) error {
	address.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Save(address).Error
}

// Delete deletes an address
func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Address{}, "id = ?", id).Error
}



// List lists addresses with filters
func (r *addressRepository) List(ctx context.Context, filters repositories.AddressFilters) ([]*entities.Address, error) {
	var addresses []*entities.Address
	query := r.db.WithContext(ctx)

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}

	if filters.Country != "" {
		query = query.Where("country = ?", filters.Country)
	}

	if filters.State != "" {
		query = query.Where("state = ?", filters.State)
	}

	if filters.City != "" {
		query = query.Where("city = ?", filters.City)
	}

	// Apply sorting
	switch filters.SortBy {
	case "created_at":
		if filters.SortOrder == "desc" {
			query = query.Order("created_at DESC")
		} else {
			query = query.Order("created_at ASC")
		}
	case "updated_at":
		if filters.SortOrder == "desc" {
			query = query.Order("updated_at DESC")
		} else {
			query = query.Order("updated_at ASC")
		}
	case "type":
		if filters.SortOrder == "desc" {
			query = query.Order("type DESC")
		} else {
			query = query.Order("type ASC")
		}
	default:
		query = query.Order("is_default DESC, created_at DESC")
	}

	// Apply pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	err := query.Find(&addresses).Error
	return addresses, err
}

// Count counts addresses with filters
func (r *addressRepository) Count(ctx context.Context, filters repositories.AddressFilters) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entities.Address{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.IsDefault != nil {
		query = query.Where("is_default = ?", *filters.IsDefault)
	}

	if filters.Country != "" {
		query = query.Where("country = ?", filters.Country)
	}

	if filters.State != "" {
		query = query.Where("state = ?", filters.State)
	}

	if filters.City != "" {
		query = query.Where("city = ?", filters.City)
	}

	err := query.Count(&count).Error
	return count, err
}

// GetByType gets addresses by type for a user
func (r *addressRepository) GetByType(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) ([]*entities.Address, error) {
	var addresses []*entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, addressType).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// ValidateAddress validates if an address belongs to a user
func (r *addressRepository) ValidateAddress(ctx context.Context, addressID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetShippingAddresses gets shipping addresses for a user
func (r *addressRepository) GetShippingAddresses(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error) {
	var addresses []*entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND (type = ? OR type = ?)", userID, entities.AddressTypeShipping, entities.AddressTypeBoth).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// GetBillingAddresses gets billing addresses for a user
func (r *addressRepository) GetBillingAddresses(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error) {
	var addresses []*entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND (type = ? OR type = ?)", userID, entities.AddressTypeBilling, entities.AddressTypeBoth).
		Order("is_default DESC, created_at DESC").
		Find(&addresses).Error
	return addresses, err
}

// CountByUser counts addresses for a user
func (r *addressRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Address{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

// DeleteByUser deletes all addresses for a user
func (r *addressRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Address{}, "user_id = ?", userID).Error
}

// ExistsByUserIDAndID checks if an address exists for a user
func (r *addressRepository) ExistsByUserIDAndID(ctx context.Context, userID, addressID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetByUserID gets all addresses for a user
func (r *addressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Address, error) {
	return r.GetByUser(ctx, userID)
}

// GetDefaultByUserID gets the default address for a user by type
func (r *addressRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) (*entities.Address, error) {
	var address entities.Address
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ? AND is_default = ?", userID, addressType, true).
		First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// SetAsDefault sets an address as default for a specific type
func (r *addressRepository) SetAsDefault(ctx context.Context, userID, addressID uuid.UUID, addressType entities.AddressType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Unset all other addresses of this type as default
		err := tx.Model(&entities.Address{}).
			Where("user_id = ? AND type = ?", userID, addressType).
			Update("is_default", false).Error
		if err != nil {
			return err
		}

		// Set the specified address as default
		return tx.Model(&entities.Address{}).
			Where("id = ? AND user_id = ? AND type = ?", addressID, userID, addressType).
			Update("is_default", true).Error
	})
}

// GetByUserIDAndType gets addresses by user and type
func (r *addressRepository) GetByUserIDAndType(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) ([]*entities.Address, error) {
	return r.GetByType(ctx, userID, addressType)
}
