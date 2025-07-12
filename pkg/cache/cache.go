package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	// Set stores a value with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	
	// Get retrieves a value by key
	Get(ctx context.Context, key string, dest interface{}) error
	
	// Delete removes a value by key
	Delete(ctx context.Context, key string) error
	
	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)
	
	// Clear removes all cached values
	Clear(ctx context.Context) error
	
	// SetMultiple stores multiple key-value pairs
	SetMultiple(ctx context.Context, items map[string]interface{}, expiration time.Duration) error
	
	// GetMultiple retrieves multiple values by keys
	GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error)
	
	// DeleteMultiple removes multiple values by keys
	DeleteMultiple(ctx context.Context, keys []string) error
	
	// Increment increments a numeric value
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	
	// Decrement decrements a numeric value
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
}

// MemoryCache implements an in-memory cache
type MemoryCache struct {
	data   map[string]*cacheItem
	mutex  sync.RWMutex
	ticker *time.Ticker
	done   chan bool
}

type cacheItem struct {
	Value      interface{}
	Expiration time.Time
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		data: make(map[string]*cacheItem),
		done: make(chan bool),
	}
	
	// Start cleanup goroutine
	cache.ticker = time.NewTicker(5 * time.Minute)
	go cache.cleanup()
	
	return cache
}

// Set stores a value with expiration
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}
	
	c.data[key] = &cacheItem{
		Value:      value,
		Expiration: exp,
	}
	
	return nil
}

// Get retrieves a value by key
func (c *MemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return fmt.Errorf("key not found: %s", key)
	}
	
	// Check expiration
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		delete(c.data, key)
		return fmt.Errorf("key expired: %s", key)
	}
	
	// Convert value to destination type
	return c.convertValue(item.Value, dest)
}

// Delete removes a value by key
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	delete(c.data, key)
	return nil
}

// Exists checks if a key exists
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	item, exists := c.data[key]
	if !exists {
		return false, nil
	}
	
	// Check expiration
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		delete(c.data, key)
		return false, nil
	}
	
	return true, nil
}

// Clear removes all cached values
func (c *MemoryCache) Clear(ctx context.Context) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.data = make(map[string]*cacheItem)
	return nil
}

// SetMultiple stores multiple key-value pairs
func (c *MemoryCache) SetMultiple(ctx context.Context, items map[string]interface{}, expiration time.Duration) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}
	
	for key, value := range items {
		c.data[key] = &cacheItem{
			Value:      value,
			Expiration: exp,
		}
	}
	
	return nil
}

// GetMultiple retrieves multiple values by keys
func (c *MemoryCache) GetMultiple(ctx context.Context, keys []string) (map[string]interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	
	result := make(map[string]interface{})
	now := time.Now()
	
	for _, key := range keys {
		if item, exists := c.data[key]; exists {
			// Check expiration
			if item.Expiration.IsZero() || now.Before(item.Expiration) {
				result[key] = item.Value
			} else {
				delete(c.data, key)
			}
		}
	}
	
	return result, nil
}

// DeleteMultiple removes multiple values by keys
func (c *MemoryCache) DeleteMultiple(ctx context.Context, keys []string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	for _, key := range keys {
		delete(c.data, key)
	}
	
	return nil
}

// Increment increments a numeric value
func (c *MemoryCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	item, exists := c.data[key]
	if !exists {
		c.data[key] = &cacheItem{Value: delta}
		return delta, nil
	}
	
	// Check expiration
	if !item.Expiration.IsZero() && time.Now().After(item.Expiration) {
		c.data[key] = &cacheItem{Value: delta}
		return delta, nil
	}
	
	// Convert to int64 and increment
	if val, ok := item.Value.(int64); ok {
		newVal := val + delta
		item.Value = newVal
		return newVal, nil
	}
	
	return 0, fmt.Errorf("value is not numeric")
}

// Decrement decrements a numeric value
func (c *MemoryCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return c.Increment(ctx, key, -delta)
}

// cleanup removes expired items
func (c *MemoryCache) cleanup() {
	for {
		select {
		case <-c.ticker.C:
			c.mutex.Lock()
			now := time.Now()
			for key, item := range c.data {
				if !item.Expiration.IsZero() && now.After(item.Expiration) {
					delete(c.data, key)
				}
			}
			c.mutex.Unlock()
		case <-c.done:
			return
		}
	}
}

// Close stops the cleanup goroutine
func (c *MemoryCache) Close() {
	c.ticker.Stop()
	close(c.done)
}

// convertValue converts a cached value to the destination type
func (c *MemoryCache) convertValue(value interface{}, dest interface{}) error {
	// If value is already the correct type, assign directly
	if reflect.TypeOf(value) == reflect.TypeOf(dest).Elem() {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(value))
		return nil
	}
	
	// Try JSON marshaling/unmarshaling for complex types
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cached value: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}
	
	return nil
}

// CacheKeyBuilder helps build consistent cache keys
type CacheKeyBuilder struct {
	prefix string
}

// NewCacheKeyBuilder creates a new cache key builder
func NewCacheKeyBuilder(prefix string) *CacheKeyBuilder {
	return &CacheKeyBuilder{prefix: prefix}
}

// Build builds a cache key with the given parts
func (b *CacheKeyBuilder) Build(parts ...string) string {
	key := b.prefix
	for _, part := range parts {
		key += ":" + part
	}
	return key
}

// ProductKey builds a product cache key
func (b *CacheKeyBuilder) ProductKey(productID string) string {
	return b.Build("product", productID)
}

// UserKey builds a user cache key
func (b *CacheKeyBuilder) UserKey(userID string) string {
	return b.Build("user", userID)
}

// CartKey builds a cart cache key
func (b *CacheKeyBuilder) CartKey(userID string) string {
	return b.Build("cart", userID)
}

// OrderKey builds an order cache key
func (b *CacheKeyBuilder) OrderKey(orderID string) string {
	return b.Build("order", orderID)
}
