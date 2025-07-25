package services

import (
	"context"

	"ecom-golang-clean-architecture/internal/domain/entities"
)

// SettingsCache interface for settings caching
type SettingsCache interface {
	// GetSetting gets a setting from cache
	GetSetting(ctx context.Context, key string) (*entities.Setting, error)
	
	// SetSetting sets a setting in cache
	SetSetting(ctx context.Context, setting *entities.Setting) error
	
	// GetSettingsByCategory gets settings by category from cache
	GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error)
	
	// SetSettingsByCategory sets settings by category in cache
	SetSettingsByCategory(ctx context.Context, category string, settings []*entities.Setting) error
	
	// InvalidateSetting removes a setting from cache
	InvalidateSetting(ctx context.Context, key string) error
	
	// InvalidateCategory removes all settings of a category from cache
	InvalidateCategory(ctx context.Context, category string) error
	
	// InvalidateAll removes all settings from cache
	InvalidateAll(ctx context.Context) error
	
	// WarmupCache preloads frequently accessed settings
	WarmupCache(ctx context.Context, settingsRepo interface {
		GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error)
	}) error
	
	// GetCacheStats returns cache statistics
	GetCacheStats(ctx context.Context) (map[string]interface{}, error)
}
