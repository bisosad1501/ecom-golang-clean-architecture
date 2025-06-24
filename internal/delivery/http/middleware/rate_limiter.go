package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

// RateLimiter interface
type RateLimiter interface {
	Allow(key string) bool
	Reset(key string) error
}

// InMemoryRateLimiter implements rate limiting using in-memory storage
type InMemoryRateLimiter struct {
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter(rps int, burst int) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// Allow checks if request is allowed
func (rl *InMemoryRateLimiter) Allow(key string) bool {
	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}
	return limiter.Allow()
}

// Reset resets the rate limiter for a key
func (rl *InMemoryRateLimiter) Reset(key string) error {
	delete(rl.limiters, key)
	return nil
}

// RedisRateLimiter implements rate limiting using Redis
type RedisRateLimiter struct {
	client *redis.Client
	window time.Duration
	limit  int
}

// NewRedisRateLimiter creates a new Redis-based rate limiter
func NewRedisRateLimiter(client *redis.Client, window time.Duration, limit int) *RedisRateLimiter {
	return &RedisRateLimiter{
		client: client,
		window: window,
		limit:  limit,
	}
}

// Allow checks if request is allowed using sliding window
func (rl *RedisRateLimiter) Allow(key string) bool {
	ctx := rl.client.Context()
	now := time.Now().Unix()
	windowStart := now - int64(rl.window.Seconds())

	pipe := rl.client.Pipeline()
	
	// Remove expired entries
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
	
	// Count current requests in window
	pipe.ZCard(ctx, key)
	
	// Add current request
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now), Member: now})
	
	// Set expiry
	pipe.Expire(ctx, key, rl.window)
	
	results, err := pipe.Exec(ctx)
	if err != nil {
		return false
	}

	// Get count from second command
	count := results[1].(*redis.IntCmd).Val()
	
	return count < int64(rl.limit)
}

// Reset resets the rate limiter for a key
func (rl *RedisRateLimiter) Reset(key string) error {
	return rl.client.Del(rl.client.Context(), key).Err()
}

// RateLimitConfig represents rate limit configuration
type RateLimitConfig struct {
	RPS        int           // Requests per second
	Burst      int           // Burst capacity
	Window     time.Duration // Time window for Redis limiter
	KeyFunc    func(*gin.Context) string // Function to generate key
	SkipFunc   func(*gin.Context) bool   // Function to skip rate limiting
	OnLimitHit func(*gin.Context)        // Function called when limit is hit
}

// DefaultKeyFunc generates key based on client IP
func DefaultKeyFunc(c *gin.Context) string {
	return c.ClientIP()
}

// UserKeyFunc generates key based on user ID
func UserKeyFunc(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists {
		return c.ClientIP()
	}
	return fmt.Sprintf("user:%v", userID)
}

// APIKeyFunc generates key based on API key
func APIKeyFunc(c *gin.Context) string {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		return c.ClientIP()
	}
	return fmt.Sprintf("api:%s", apiKey)
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter RateLimiter, config RateLimitConfig) gin.HandlerFunc {
	if config.KeyFunc == nil {
		config.KeyFunc = DefaultKeyFunc
	}

	return func(c *gin.Context) {
		// Skip rate limiting if skip function returns true
		if config.SkipFunc != nil && config.SkipFunc(c) {
			c.Next()
			return
		}

		key := config.KeyFunc(c)
		
		if !limiter.Allow(key) {
			// Rate limit exceeded
			if config.OnLimitHit != nil {
				config.OnLimitHit(c)
			}

			c.Header("X-RateLimit-Limit", strconv.Itoa(config.RPS))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CreateRateLimiters creates different rate limiters for different endpoints
func CreateRateLimiters(redisClient *redis.Client) map[string]gin.HandlerFunc {
	limiters := make(map[string]gin.HandlerFunc)

	// General API rate limiter (100 requests per minute)
	generalLimiter := NewRedisRateLimiter(redisClient, time.Minute, 100)
	limiters["general"] = RateLimitMiddleware(generalLimiter, RateLimitConfig{
		RPS:     100,
		Window:  time.Minute,
		KeyFunc: DefaultKeyFunc,
	})

	// Auth endpoints rate limiter (5 requests per minute)
	authLimiter := NewRedisRateLimiter(redisClient, time.Minute, 5)
	limiters["auth"] = RateLimitMiddleware(authLimiter, RateLimitConfig{
		RPS:     5,
		Window:  time.Minute,
		KeyFunc: DefaultKeyFunc,
		OnLimitHit: func(c *gin.Context) {
			// Log suspicious activity
			fmt.Printf("Rate limit hit for auth endpoint from IP: %s\n", c.ClientIP())
		},
	})

	// Search endpoints rate limiter (20 requests per minute)
	searchLimiter := NewRedisRateLimiter(redisClient, time.Minute, 20)
	limiters["search"] = RateLimitMiddleware(searchLimiter, RateLimitConfig{
		RPS:     20,
		Window:  time.Minute,
		KeyFunc: DefaultKeyFunc,
	})

	// Admin endpoints rate limiter (200 requests per minute)
	adminLimiter := NewRedisRateLimiter(redisClient, time.Minute, 200)
	limiters["admin"] = RateLimitMiddleware(adminLimiter, RateLimitConfig{
		RPS:     200,
		Window:  time.Minute,
		KeyFunc: UserKeyFunc,
		SkipFunc: func(c *gin.Context) bool {
			// Skip rate limiting for super admin
			role, exists := c.Get("user_role")
			return exists && role == "super_admin"
		},
	})

	// File upload rate limiter (10 uploads per hour)
	uploadLimiter := NewRedisRateLimiter(redisClient, time.Hour, 10)
	limiters["upload"] = RateLimitMiddleware(uploadLimiter, RateLimitConfig{
		RPS:     10,
		Window:  time.Hour,
		KeyFunc: UserKeyFunc,
	})

	return limiters
}

// Adaptive rate limiter that adjusts based on system load
type AdaptiveRateLimiter struct {
	baseLimiter RateLimiter
	baseLimit   int
	maxLimit    int
	minLimit    int
	loadFunc    func() float64 // Function to get current system load (0.0 to 1.0)
}

// NewAdaptiveRateLimiter creates a new adaptive rate limiter
func NewAdaptiveRateLimiter(baseLimiter RateLimiter, baseLimit, minLimit, maxLimit int, loadFunc func() float64) *AdaptiveRateLimiter {
	return &AdaptiveRateLimiter{
		baseLimiter: baseLimiter,
		baseLimit:   baseLimit,
		maxLimit:    maxLimit,
		minLimit:    minLimit,
		loadFunc:    loadFunc,
	}
}

// Allow checks if request is allowed with adaptive limiting
func (arl *AdaptiveRateLimiter) Allow(key string) bool {
	load := arl.loadFunc()
	
	// Adjust limit based on load
	_ = arl.baseLimit // Use base limit for now
	if load > 0.8 {
		// High load, reduce limit - would adjust limiter in real implementation
	} else if load < 0.3 {
		// Low load, increase limit - would adjust limiter in real implementation
	}

	// For simplicity, we'll use the base limiter
	// In a real implementation, you'd adjust the limiter's rate
	return arl.baseLimiter.Allow(key)
}

// Reset resets the adaptive rate limiter
func (arl *AdaptiveRateLimiter) Reset(key string) error {
	return arl.baseLimiter.Reset(key)
}

// Circuit breaker pattern for rate limiting
type CircuitBreakerRateLimiter struct {
	limiter       RateLimiter
	failureCount  int
	maxFailures   int
	resetTimeout  time.Duration
	lastFailTime  time.Time
	state         string // "closed", "open", "half-open"
}

// NewCircuitBreakerRateLimiter creates a new circuit breaker rate limiter
func NewCircuitBreakerRateLimiter(limiter RateLimiter, maxFailures int, resetTimeout time.Duration) *CircuitBreakerRateLimiter {
	return &CircuitBreakerRateLimiter{
		limiter:      limiter,
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        "closed",
	}
}

// Allow checks if request is allowed with circuit breaker
func (cb *CircuitBreakerRateLimiter) Allow(key string) bool {
	now := time.Now()

	switch cb.state {
	case "open":
		if now.Sub(cb.lastFailTime) > cb.resetTimeout {
			cb.state = "half-open"
			cb.failureCount = 0
		} else {
			return false
		}
	case "half-open":
		// Allow limited requests to test if service is back
		if cb.limiter.Allow(key) {
			cb.state = "closed"
			cb.failureCount = 0
			return true
		} else {
			cb.state = "open"
			cb.lastFailTime = now
			return false
		}
	}

	// Closed state
	if cb.limiter.Allow(key) {
		return true
	} else {
		cb.failureCount++
		if cb.failureCount >= cb.maxFailures {
			cb.state = "open"
			cb.lastFailTime = now
		}
		return false
	}
}

// Reset resets the circuit breaker
func (cb *CircuitBreakerRateLimiter) Reset(key string) error {
	cb.state = "closed"
	cb.failureCount = 0
	return cb.limiter.Reset(key)
}
