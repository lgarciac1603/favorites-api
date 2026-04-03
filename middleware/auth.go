package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// validate user's token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get auth header
		authHeader := c.GetHeader("Authorization")

		// validate header
		if authHeader == "" {
			c.JSON(401, gin.H{ "error": "Token not found" })
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{ "error": "Invalid token format" })
			c.Abort()
			return
		}

		token := parts[1]

		/* TODO: validate with users API:
		*	Simulates that all tokens are valid
		*/
		userID, err := ValidateToken(token)
		if err != nil {
			c.JSON(401, gin.H{ "error": "Invalid token" })
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

/* TODO: Connect with users backend
* ValidateToken to return userID
*/
func ValidateToken(token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("empty token")
	}

	return token, nil
}
