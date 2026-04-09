package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type meResponse struct {
	ID int `json:"id"`
}

// AuthMiddleware validates user's token
func AuthMiddleware(authAPIBaseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Token not found"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		token := parts[1]

		userID, err := ValidateToken(authAPIBaseURL, token)
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// ValidateToken calls the primary API to validate token and return userId
func ValidateToken(authAPIBaseURL, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("empty token")
	}
	if authAPIBaseURL == "" {
		return "", fmt.Errorf("missing AUTH_API_URL")
	}

	base := strings.TrimRight(authAPIBaseURL, "/")
	req, err := http.NewRequest(http.MethodGet, base+"/me", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token validation failed: %s", resp.Status)
	}

	var me meResponse
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return "", err
	}
	if me.ID == 0 {
		return "", fmt.Errorf("missing user id")
	}

	return strconv.Itoa(me.ID), nil
}
