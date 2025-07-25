package database

import (
	"context"
	"fmt"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type settingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *gorm.DB) repositories.SettingsRepository {
	return &settingsRepository{db: db}
}

// Settings CRUD
func (r *settingsRepository) GetSetting(ctx context.Context, key string) (*entities.Setting, error) {
	var setting entities.Setting
	err := r.db.WithContext(ctx).Where("key = ?", key).First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("setting not found")
		}
		return nil, err
	}
	return &setting, nil
}

func (r *settingsRepository) GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error) {
	var settings []*entities.Setting
	err := r.db.WithContext(ctx).Where("category = ?", category).Find(&settings).Error
	return settings, err
}

func (r *settingsRepository) GetAllSettings(ctx context.Context) ([]*entities.Setting, error) {
	var settings []*entities.Setting
	err := r.db.WithContext(ctx).Find(&settings).Error
	return settings, err
}

func (r *settingsRepository) SetSetting(ctx context.Context, setting *entities.Setting) error {
	return r.db.WithContext(ctx).Save(setting).Error
}

func (r *settingsRepository) SetSettings(ctx context.Context, settings []*entities.Setting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, setting := range settings {
			// Use ON CONFLICT to handle duplicate keys
			result := tx.Exec(`
				INSERT INTO settings (key, value, type, category, description, is_public, is_encrypted, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
				ON CONFLICT (key) DO UPDATE SET
					value = EXCLUDED.value,
					type = EXCLUDED.type,
					category = EXCLUDED.category,
					description = EXCLUDED.description,
					is_public = EXCLUDED.is_public,
					is_encrypted = EXCLUDED.is_encrypted,
					updated_at = NOW()
			`, setting.Key, setting.Value, setting.Type, setting.Category, setting.Description, setting.IsPublic, setting.IsEncrypted)

			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

func (r *settingsRepository) DeleteSetting(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&entities.Setting{}).Error
}

// Tax Management
func (r *settingsRepository) GetTaxRates(ctx context.Context, filters repositories.TaxRateFilters) ([]*entities.TaxRate, error) {
	query := r.db.WithContext(ctx).Model(&entities.TaxRate{})

	if filters.Country != "" {
		query = query.Where("country = ?", filters.Country)
	}
	if filters.State != "" {
		query = query.Where("state = ?", filters.State)
	}
	if filters.City != "" {
		query = query.Where("city = ?", filters.City)
	}
	if filters.ZipCode != "" {
		query = query.Where("zip_code = ?", filters.ZipCode)
	}
	if filters.TaxClass != "" {
		query = query.Where("tax_class = ?", filters.TaxClass)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var taxRates []*entities.TaxRate
	err := query.Order("priority ASC, created_at DESC").Find(&taxRates).Error
	return taxRates, err
}

func (r *settingsRepository) GetTaxRate(ctx context.Context, id uuid.UUID) (*entities.TaxRate, error) {
	var taxRate entities.TaxRate
	err := r.db.WithContext(ctx).First(&taxRate, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tax rate not found")
		}
		return nil, err
	}
	return &taxRate, nil
}

func (r *settingsRepository) CreateTaxRate(ctx context.Context, taxRate *entities.TaxRate) error {
	return r.db.WithContext(ctx).Create(taxRate).Error
}

func (r *settingsRepository) UpdateTaxRate(ctx context.Context, taxRate *entities.TaxRate) error {
	return r.db.WithContext(ctx).Save(taxRate).Error
}

func (r *settingsRepository) DeleteTaxRate(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.TaxRate{}, "id = ?", id).Error
}

func (r *settingsRepository) GetTaxClasses(ctx context.Context) ([]*entities.TaxClass, error) {
	var taxClasses []*entities.TaxClass
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&taxClasses).Error
	return taxClasses, err
}

func (r *settingsRepository) GetTaxClass(ctx context.Context, id uuid.UUID) (*entities.TaxClass, error) {
	var taxClass entities.TaxClass
	err := r.db.WithContext(ctx).First(&taxClass, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tax class not found")
		}
		return nil, err
	}
	return &taxClass, nil
}

func (r *settingsRepository) CreateTaxClass(ctx context.Context, taxClass *entities.TaxClass) error {
	return r.db.WithContext(ctx).Create(taxClass).Error
}

func (r *settingsRepository) UpdateTaxClass(ctx context.Context, taxClass *entities.TaxClass) error {
	return r.db.WithContext(ctx).Save(taxClass).Error
}

func (r *settingsRepository) DeleteTaxClass(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.TaxClass{}, "id = ?", id).Error
}

// Shipping Management
func (r *settingsRepository) GetShippingZones(ctx context.Context) ([]*entities.ShippingZone, error) {
	var zones []*entities.ShippingZone
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("priority ASC, name ASC").Find(&zones).Error
	return zones, err
}

func (r *settingsRepository) GetShippingZone(ctx context.Context, id uuid.UUID) (*entities.ShippingZone, error) {
	var zone entities.ShippingZone
	err := r.db.WithContext(ctx).First(&zone, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shipping zone not found")
		}
		return nil, err
	}
	return &zone, nil
}

func (r *settingsRepository) CreateShippingZone(ctx context.Context, zone *entities.ShippingZone) error {
	return r.db.WithContext(ctx).Create(zone).Error
}

func (r *settingsRepository) UpdateShippingZone(ctx context.Context, zone *entities.ShippingZone) error {
	return r.db.WithContext(ctx).Save(zone).Error
}

func (r *settingsRepository) DeleteShippingZone(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ShippingZone{}, "id = ?", id).Error
}

func (r *settingsRepository) GetShippingMethods(ctx context.Context, zoneID *uuid.UUID) ([]*entities.ShippingMethod, error) {
	query := r.db.WithContext(ctx).Preload("Zone")

	if zoneID != nil {
		query = query.Where("zone_id = ?", *zoneID)
	}

	var methods []*entities.ShippingMethod
	err := query.Where("is_active = ?", true).Order("name ASC").Find(&methods).Error
	return methods, err
}

func (r *settingsRepository) GetShippingMethod(ctx context.Context, id uuid.UUID) (*entities.ShippingMethod, error) {
	var method entities.ShippingMethod
	err := r.db.WithContext(ctx).Preload("Zone").First(&method, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shipping method not found")
		}
		return nil, err
	}
	return &method, nil
}

func (r *settingsRepository) CreateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error {
	return r.db.WithContext(ctx).Create(method).Error
}

func (r *settingsRepository) UpdateShippingMethod(ctx context.Context, method *entities.ShippingMethod) error {
	return r.db.WithContext(ctx).Save(method).Error
}

func (r *settingsRepository) DeleteShippingMethod(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ShippingMethod{}, "id = ?", id).Error
}

func (r *settingsRepository) GetShippingClasses(ctx context.Context) ([]*entities.ShippingClass, error) {
	var classes []*entities.ShippingClass
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&classes).Error
	return classes, err
}

func (r *settingsRepository) GetShippingClass(ctx context.Context, id uuid.UUID) (*entities.ShippingClass, error) {
	var class entities.ShippingClass
	err := r.db.WithContext(ctx).First(&class, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("shipping class not found")
		}
		return nil, err
	}
	return &class, nil
}

func (r *settingsRepository) CreateShippingClass(ctx context.Context, class *entities.ShippingClass) error {
	return r.db.WithContext(ctx).Create(class).Error
}

func (r *settingsRepository) UpdateShippingClass(ctx context.Context, class *entities.ShippingClass) error {
	return r.db.WithContext(ctx).Save(class).Error
}

func (r *settingsRepository) DeleteShippingClass(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.ShippingClass{}, "id = ?", id).Error
}

func (r *settingsRepository) GetPackagingOptions(ctx context.Context) ([]*entities.PackagingOption, error) {
	var options []*entities.PackagingOption
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("name ASC").Find(&options).Error
	return options, err
}

func (r *settingsRepository) GetPackagingOption(ctx context.Context, id uuid.UUID) (*entities.PackagingOption, error) {
	var option entities.PackagingOption
	err := r.db.WithContext(ctx).First(&option, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("packaging option not found")
		}
		return nil, err
	}
	return &option, nil
}

func (r *settingsRepository) CreatePackagingOption(ctx context.Context, option *entities.PackagingOption) error {
	return r.db.WithContext(ctx).Create(option).Error
}

func (r *settingsRepository) UpdatePackagingOption(ctx context.Context, option *entities.PackagingOption) error {
	return r.db.WithContext(ctx).Save(option).Error
}

func (r *settingsRepository) DeletePackagingOption(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.PackagingOption{}, "id = ?", id).Error
}
