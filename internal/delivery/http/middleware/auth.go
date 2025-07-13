package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"ecom-golang-clean-architecture/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddlewareStruct holds the auth middleware configuration
type AuthMiddlewareStruct struct {
	jwtSecret string
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddlewareStruct {
	return &AuthMiddlewareStruct{
		jwtSecret: cfg.JWT.Secret,
	}
}

// RequireAuth returns a middleware that requires authentication
func (a *AuthMiddlewareStruct) RequireAuth() gin.HandlerFunc {
	return AuthMiddleware(a.jwtSecret)
}

// AuthMiddleware creates JWT authentication middleware
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method is HMAC and specifically HS256
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			} else if method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", method.Alg())
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract and validate claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check expiration
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Token has expired",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token missing expiration",
				})
				c.Abort()
				return
			}

			// Validate required claims
			userID, hasUserID := claims["user_id"]
			email, hasEmail := claims["email"]
			role, hasRole := claims["role"]

			if !hasUserID || !hasEmail || !hasRole {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token missing required claims",
				})
				c.Abort()
				return
			}

			// Validate user_id format
			if userIDStr, ok := userID.(string); ok {
				if len(userIDStr) == 0 {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid user ID in token",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid user ID format in token",
				})
				c.Abort()
				return
			}

			// Validate email format
			if emailStr, ok := email.(string); ok {
				if len(emailStr) == 0 || !strings.Contains(emailStr, "@") {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid email in token",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid email format in token",
				})
				c.Abort()
				return
			}

			// Validate role
			if roleStr, ok := role.(string); ok {
				validRoles := []string{string(entities.UserRoleCustomer), string(entities.UserRoleModerator), string(entities.UserRoleAdmin)}
				isValidRole := false
				for _, validRole := range validRoles {
					if roleStr == validRole {
						isValidRole = true
						break
					}
				}
				if !isValidRole {
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": "Invalid role in token",
					})
					c.Abort()
					return
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid role format in token",
				})
				c.Abort()
				return
			}

			c.Set("user_id", userID)
			c.Set("email", email)
			c.Set("role", role)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if user has admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found",
			})
			c.Abort()
			return
		}

		// Add debug logging
		fmt.Printf("AdminMiddleware: role=%v, expectedRole=%s\n", role, string(entities.UserRoleAdmin))

		if role != string(entities.UserRoleAdmin) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ModeratorMiddleware checks if user has moderator or admin role
func ModeratorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User role not found",
			})
			c.Abort()
			return
		}

		userRole := role.(string)
		
		// Add debug logging
		fmt.Printf("ModeratorMiddleware: role=%s, checking against admin=%s or moderator=%s\n", 
			userRole, string(entities.UserRoleAdmin), string(entities.UserRoleModerator))
		
		if userRole != string(entities.UserRoleAdmin) && userRole != string(entities.UserRoleModerator) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Moderator or admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
