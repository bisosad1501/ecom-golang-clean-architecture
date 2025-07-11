package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"ecom-golang-clean-architecture/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

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

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
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
