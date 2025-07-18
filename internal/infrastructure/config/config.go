package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Email    EmailConfig
	Payment  PaymentConfig
	Upload   UploadConfig
	Log      LogConfig
	CORS     CORSConfig
}

// AppConfig holds application configuration
type AppConfig struct {
	Name string
	Env  string
	Host string
	Port string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	Timezone string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

// EmailConfig holds Gmail SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	ReplyToEmail string
}



// PaymentConfig holds payment configuration
type PaymentConfig struct {
	StripeSecretKey      string
	StripePublishableKey string
	StripeWebhookSecret  string
	PayPalClientID       string
	PayPalClientSecret   string
	PayPalSandbox        bool
}

// UploadConfig holds file upload configuration
type UploadConfig struct {
	Path        string
	MaxFileSize int64
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string
	Format string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Có thể log hoặc bỏ qua nếu không cần thiết
	}

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "ecom-api"),
			Env:  getEnv("APP_ENV", "development"),
			Host: getEnv("APP_HOST", "localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "ecommerce_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Timezone: getEnv("DB_TIMEZONE", "UTC"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:      getEnvOrPanic("JWT_SECRET", "JWT_SECRET is required for security"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		Email: EmailConfig{
			SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			SMTPPort:     getEnv("SMTP_PORT", "587"),
			SMTPUsername: getEnv("SMTP_USERNAME", ""),
			SMTPPassword: getEnv("SMTP_PASSWORD", ""),
			FromEmail:    getEnv("FROM_EMAIL", ""),
			FromName:     getEnv("FROM_NAME", "BiHub Store"),
			ReplyToEmail: getEnv("REPLY_TO_EMAIL", ""),
		},
		Payment: PaymentConfig{
			StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
			StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PayPalClientID:       getEnv("PAYPAL_CLIENT_ID", ""),
			PayPalClientSecret:   getEnv("PAYPAL_CLIENT_SECRET", ""),
			PayPalSandbox:        getEnvAsBool("PAYPAL_SANDBOX", true),
		},
		Upload: UploadConfig{
			Path:        getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:8080"}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Content-Type", "Authorization", "X-Session-ID"}),
		},
	}

	return config, nil
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.Name +
		" sslmode=" + c.SSLMode +
		" TimeZone=" + c.Timezone
}

// GetRedisAddr returns Redis address
func (c *RedisConfig) GetRedisAddr() string {
	return c.Host + ":" + c.Port
}

// GetJWTExpireDuration returns JWT expire duration
func (c *JWTConfig) GetJWTExpireDuration() time.Duration {
	return time.Duration(c.ExpireHours) * time.Hour
}

// IsProduction checks if the environment is production
func (c *AppConfig) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment checks if the environment is development
func (c *AppConfig) IsDevelopment() bool {
	return c.Env == "development"
}

// GetAddress returns the full address
func (c *AppConfig) GetAddress() string {
	return c.Host + ":" + c.Port
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrPanic(key, message string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	panic(message + ": " + key + " environment variable is required")
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		var result []string
		for _, item := range strings.Split(value, ",") {
			if trimmed := strings.TrimSpace(item); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
