package handlers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lgarciac1603/favorites-api/models"
)

// FavoritesHandler encapsulates favorite handlers
type FavoritesHandler struct {
	DB *sql.DB
}

// NewFavoritesHandler creates a new handler
func NewFavoritesHandler(db *sql.DB) *FavoritesHandler {
	return &FavoritesHandler{DB: db}
}

// GetFavorites retrieves all favorites for the authenticated user
func (h *FavoritesHandler) GetFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userIDInt := userID.(int)

	query := `
		SELECT id, user_id, crypto_id, crypto_name, created_at
		FROM user_favorites
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := h.DB.Query(query, userIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error querying database"})
		return
	}
	defer rows.Close()

	var favorites []models.Favorite

	for rows.Next() {
		var fav models.Favorite
		err := rows.Scan(&fav.ID, &fav.UserID, &fav.CryptoID, &fav.CryptoName, &fav.CreatedAt)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error reading data"})
			return
		}
		favorites = append(favorites, fav)
	}

	if favorites == nil {
		favorites = []models.Favorite{}
	}

	c.JSON(200, gin.H{
		"data":  favorites,
		"total": len(favorites),
	})
}

// PostFavorite adds a cryptocurrency to user's favorites
func (h *FavoritesHandler) PostFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userIDInt := userID.(int)

	var requestBody struct {
		CryptoID   string `json:"cryptoId" binding:"required"`
		CryptoName string `json:"cryptoName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "CryptoID and CryptoName are required"})
		return
	}

	checkQuery := `SELECT id FROM user_favorites WHERE user_id = $1 AND crypto_id = $2`
	var existingID int
	err := h.DB.QueryRow(checkQuery, userIDInt, requestBody.CryptoID).Scan(&existingID)

	if err == nil {
		c.JSON(409, gin.H{"error": "This crypto is already in favorites"})
		return
	}

	insertQuery := `
		INSERT INTO user_favorites (user_id, crypto_id, crypto_name, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, crypto_id, crypto_name, created_at
	`

	var newFavorite models.Favorite
	createdAt := time.Now().Format(time.RFC3339)

	err = h.DB.QueryRow(
		insertQuery,
		userIDInt,
		requestBody.CryptoID,
		requestBody.CryptoName,
		createdAt,
	).Scan(
		&newFavorite.ID,
		&newFavorite.UserID,
		&newFavorite.CryptoID,
		&newFavorite.CryptoName,
		&newFavorite.CreatedAt,
	)

	if err != nil {
		c.JSON(500, gin.H{"error": "Error inserting into database"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Crypto added to favorites",
		"data":    newFavorite,
	})
}

// DeleteFavorite removes a cryptocurrency from user's favorites
func (h *FavoritesHandler) DeleteFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	userIDInt := userID.(int)
	cryptoID := c.Param("cryptoId")

	if cryptoID == "" {
		c.JSON(400, gin.H{"error": "CryptoId is required"})
		return
	}

	deleteQuery := `
		DELETE FROM user_favorites
		WHERE user_id = $1 AND crypto_id = $2
	`

	result, err := h.DB.Exec(deleteQuery, userIDInt, cryptoID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error deleting from database"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Crypto not found in favorites"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Crypto removed from favorites",
	})
}
