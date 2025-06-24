package usecases

import (
	"context"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"
	"github.com/google/uuid"
)

// AddressUseCase defines address use cases
type AddressUseCase interface {
	CreateAddress(ctx context.Context, userID uuid.UUID, req CreateAddressRequest) (*AddressResponse, error)
	GetUserAddresses(ctx context.Context, userID uuid.UUID) ([]*AddressResponse, error)
	GetAddress(ctx context.Context, userID, addressID uuid.UUID) (*AddressResponse, error)
	UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req UpdateAddressRequest) (*AddressResponse, error)
	DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error
	SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID, addressType entities.AddressType) error
	GetDefaultAddress(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) (*AddressResponse, error)
}

type addressUseCase struct {
	addressRepo repositories.AddressRepository
}

// NewAddressUseCase creates a new address use case
func NewAddressUseCase(addressRepo repositories.AddressRepository) AddressUseCase {
	return &addressUseCase{
		addressRepo: addressRepo,
	}
}

// CreateAddressRequest represents create address request
type CreateAddressRequest struct {
	Type      entities.AddressType `json:"type" validate:"required,oneof=shipping billing both"`
	FirstName string               `json:"first_name" validate:"required"`
	LastName  string               `json:"last_name" validate:"required"`
	Company   string               `json:"company"`
	Address1  string               `json:"address1" validate:"required"`
	Address2  string               `json:"address2"`
	City      string               `json:"city" validate:"required"`
	State     string               `json:"state" validate:"required"`
	ZipCode   string               `json:"zip_code" validate:"required"`
	Country   string               `json:"country" validate:"required"`
	Phone     string               `json:"phone"`
	IsDefault bool                 `json:"is_default"`
}

// UpdateAddressRequest represents update address request
type UpdateAddressRequest struct {
	Type      *entities.AddressType `json:"type"`
	FirstName *string               `json:"first_name"`
	LastName  *string               `json:"last_name"`
	Company   *string               `json:"company"`
	Address1  *string               `json:"address1"`
	Address2  *string               `json:"address2"`
	City      *string               `json:"city"`
	State     *string               `json:"state"`
	ZipCode   *string               `json:"zip_code"`
	Country   *string               `json:"country"`
	Phone     *string               `json:"phone"`
	IsDefault *bool                 `json:"is_default"`
}

// AddressResponse represents address response
type AddressResponse struct {
	ID          uuid.UUID            `json:"id"`
	Type        entities.AddressType `json:"type"`
	FirstName   string               `json:"first_name"`
	LastName    string               `json:"last_name"`
	Company     string               `json:"company"`
	Address1    string               `json:"address1"`
	Address2    string               `json:"address2"`
	City        string               `json:"city"`
	State       string               `json:"state"`
	ZipCode     string               `json:"zip_code"`
	Country     string               `json:"country"`
	Phone       string               `json:"phone"`
	IsDefault   bool                 `json:"is_default"`
	FullName    string               `json:"full_name"`
	FullAddress string               `json:"full_address"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

// CreateAddress creates a new address for user
func (uc *addressUseCase) CreateAddress(ctx context.Context, userID uuid.UUID, req CreateAddressRequest) (*AddressResponse, error) {
	address := &entities.Address{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      req.Type,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Company:   req.Company,
		Address1:  req.Address1,
		Address2:  req.Address2,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		Country:   req.Country,
		Phone:     req.Phone,
		IsDefault: req.IsDefault,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.addressRepo.Create(ctx, address); err != nil {
		return nil, err
	}

	// If this is set as default, update other addresses
	if req.IsDefault {
		if err := uc.addressRepo.SetAsDefault(ctx, userID, address.ID, req.Type); err != nil {
			return nil, err
		}
	}

	return uc.toAddressResponse(address), nil
}

// GetUserAddresses gets all addresses for a user
func (uc *addressUseCase) GetUserAddresses(ctx context.Context, userID uuid.UUID) ([]*AddressResponse, error) {
	addresses, err := uc.addressRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*AddressResponse, len(addresses))
	for i, address := range addresses {
		responses[i] = uc.toAddressResponse(address)
	}

	return responses, nil
}

// GetAddress gets a specific address for user
func (uc *addressUseCase) GetAddress(ctx context.Context, userID, addressID uuid.UUID) (*AddressResponse, error) {
	// Verify address belongs to user
	exists, err := uc.addressRepo.ExistsByUserIDAndID(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrAddressNotFound
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	return uc.toAddressResponse(address), nil
}

// UpdateAddress updates an existing address
func (uc *addressUseCase) UpdateAddress(ctx context.Context, userID, addressID uuid.UUID, req UpdateAddressRequest) (*AddressResponse, error) {
	// Verify address belongs to user
	exists, err := uc.addressRepo.ExistsByUserIDAndID(ctx, userID, addressID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, entities.ErrAddressNotFound
	}

	address, err := uc.addressRepo.GetByID(ctx, addressID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Type != nil {
		address.Type = *req.Type
	}
	if req.FirstName != nil {
		address.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		address.LastName = *req.LastName
	}
	if req.Company != nil {
		address.Company = *req.Company
	}
	if req.Address1 != nil {
		address.Address1 = *req.Address1
	}
	if req.Address2 != nil {
		address.Address2 = *req.Address2
	}
	if req.City != nil {
		address.City = *req.City
	}
	if req.State != nil {
		address.State = *req.State
	}
	if req.ZipCode != nil {
		address.ZipCode = *req.ZipCode
	}
	if req.Country != nil {
		address.Country = *req.Country
	}
	if req.Phone != nil {
		address.Phone = *req.Phone
	}
	if req.IsDefault != nil {
		address.IsDefault = *req.IsDefault
	}

	address.UpdatedAt = time.Now()

	if err := uc.addressRepo.Update(ctx, address); err != nil {
		return nil, err
	}

	// If this is set as default, update other addresses
	if req.IsDefault != nil && *req.IsDefault {
		if err := uc.addressRepo.SetAsDefault(ctx, userID, addressID, address.Type); err != nil {
			return nil, err
		}
	}

	return uc.toAddressResponse(address), nil
}

// DeleteAddress deletes an address
func (uc *addressUseCase) DeleteAddress(ctx context.Context, userID, addressID uuid.UUID) error {
	// Verify address belongs to user
	exists, err := uc.addressRepo.ExistsByUserIDAndID(ctx, userID, addressID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrAddressNotFound
	}

	return uc.addressRepo.Delete(ctx, addressID)
}

// SetDefaultAddress sets an address as default
func (uc *addressUseCase) SetDefaultAddress(ctx context.Context, userID, addressID uuid.UUID, addressType entities.AddressType) error {
	// Verify address belongs to user
	exists, err := uc.addressRepo.ExistsByUserIDAndID(ctx, userID, addressID)
	if err != nil {
		return err
	}
	if !exists {
		return entities.ErrAddressNotFound
	}

	return uc.addressRepo.SetAsDefault(ctx, userID, addressID, addressType)
}

// GetDefaultAddress gets the default address for a user
func (uc *addressUseCase) GetDefaultAddress(ctx context.Context, userID uuid.UUID, addressType entities.AddressType) (*AddressResponse, error) {
	address, err := uc.addressRepo.GetDefaultByUserID(ctx, userID, addressType)
	if err != nil {
		return nil, err
	}

	return uc.toAddressResponse(address), nil
}

// toAddressResponse converts address entity to response
func (uc *addressUseCase) toAddressResponse(address *entities.Address) *AddressResponse {
	return &AddressResponse{
		ID:          address.ID,
		Type:        address.Type,
		FirstName:   address.FirstName,
		LastName:    address.LastName,
		Company:     address.Company,
		Address1:    address.Address1,
		Address2:    address.Address2,
		City:        address.City,
		State:       address.State,
		ZipCode:     address.ZipCode,
		Country:     address.Country,
		Phone:       address.Phone,
		IsDefault:   address.IsDefault,
		FullName:    address.GetFullName(),
		FullAddress: address.GetFullAddress(),
		CreatedAt:   address.CreatedAt,
		UpdatedAt:   address.UpdatedAt,
	}
}
