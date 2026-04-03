// middleware/auth.go
package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates user's token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")

		// Validate header exists
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}

		// Validate Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token := parts[1]

		// TODO: validate with users API
		// Currently simulates that all tokens are valid
		userID, err := ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// TODO: Connect with users backend
// ValidateToken returns userID from token
func ValidateToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}
