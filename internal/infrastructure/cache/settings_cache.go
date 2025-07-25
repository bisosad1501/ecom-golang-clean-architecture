package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"

	"github.com/go-redis/redis/v8"
)

// SettingsCache handles caching for settings
type SettingsCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewSettingsCache creates a new settings cache
func NewSettingsCache(client *redis.Client) *SettingsCache {
	return &SettingsCache{
		client: client,
		ttl:    15 * time.Minute, // Cache for 15 minutes
	}
}

// GetSetting gets a setting from cache
func (c *SettingsCache) GetSetting(ctx context.Context, key string) (*entities.Setting, error) {
	cacheKey := fmt.Sprintf("setting:%s", key)

	data, err := c.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var setting entities.Setting
	if err := json.Unmarshal([]byte(data), &setting); err != nil {
		return nil, err
	}

	return &setting, nil
}

// SetSetting sets a setting in cache
func (c *SettingsCache) SetSetting(ctx context.Context, setting *entities.Setting) error {
	cacheKey := fmt.Sprintf("setting:%s", setting.Key)

	data, err := json.Marshal(setting)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey, data, c.ttl).Err()
}

// GetSettingsByCategory gets settings by category from cache
func (c *SettingsCache) GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error) {
	cacheKey := fmt.Sprintf("settings:category:%s", category)

	data, err := c.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var settings []*entities.Setting
	if err := json.Unmarshal([]byte(data), &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

// SetSettingsByCategory sets settings by category in cache
func (c *SettingsCache) SetSettingsByCategory(ctx context.Context, category string, settings []*entities.Setting) error {
	cacheKey := fmt.Sprintf("settings:category:%s", category)

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, cacheKey, data, c.ttl).Err()
}

// InvalidateSetting removes a setting from cache
func (c *SettingsCache) InvalidateSetting(ctx context.Context, key string) error {
	cacheKey := fmt.Sprintf("setting:%s", key)
	return c.client.Del(ctx, cacheKey).Err()
}

// InvalidateCategory removes all settings of a category from cache
func (c *SettingsCache) InvalidateCategory(ctx context.Context, category string) error {
	cacheKey := fmt.Sprintf("settings:category:%s", category)
	return c.client.Del(ctx, cacheKey).Err()
}

// InvalidateAll removes all settings from cache
func (c *SettingsCache) InvalidateAll(ctx context.Context) error {
	// Get all setting keys
	keys, err := c.client.Keys(ctx, "setting:*").Result()
	if err != nil {
		return err
	}

	// Get all category keys
	categoryKeys, err := c.client.Keys(ctx, "settings:category:*").Result()
	if err != nil {
		return err
	}

	// Combine all keys
	allKeys := append(keys, categoryKeys...)

	if len(allKeys) > 0 {
		return c.client.Del(ctx, allKeys...).Err()
	}

	return nil
}

// WarmupCache preloads frequently accessed settings
func (c *SettingsCache) WarmupCache(ctx context.Context, settingsRepo interface {
	GetSettingsByCategory(ctx context.Context, category string) ([]*entities.Setting, error)
}) error {
	categories := []string{"general", "store", "payment", "email", "tax", "shipping", "seo", "security", "notifications", "integrations"}

	for _, category := range categories {
		settings, err := settingsRepo.GetSettingsByCategory(ctx, category)
		if err != nil {
			continue // Skip on error, don't fail the entire warmup
		}

		// Cache the category
		c.SetSettingsByCategory(ctx, category, settings)

		// Cache individual settings
		for _, setting := range settings {
			c.SetSetting(ctx, setting)
		}
	}

	return nil
}

// GetCacheStats returns cache statistics
func (c *SettingsCache) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.client.Info(ctx, "memory").Result()
	if err != nil {
		return nil, err
	}

	// Count setting keys
	settingKeys, err := c.client.Keys(ctx, "setting:*").Result()
	if err != nil {
		return nil, err
	}

	categoryKeys, err := c.client.Keys(ctx, "settings:category:*").Result()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"setting_keys_count":  len(settingKeys),
		"category_keys_count": len(categoryKeys),
		"total_keys":          len(settingKeys) + len(categoryKeys),
		"ttl_minutes":         int(c.ttl.Minutes()),
		"memory_info":         info,
	}

	return stats, nil
}
