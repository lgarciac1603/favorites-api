package handlers

import (
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lgarciac1603/favorites-api/models"
)

// FavoritesHandler encapsula los handlers de favoritos
type FavoritesHandler struct {
	DB *sql.DB
}

// NewFavoritesHandler crea un nuevo handler
func NewFavoritesHandler(db *sql.DB) *FavoritesHandler {
	return &FavoritesHandler{DB: db}
}

// GetFavorites obtiene todos los favoritos del usuario
func (h *FavoritesHandler) GetFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Usuario no autenticado"})
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
		c.JSON(500, gin.H{"error": "Error consultando BD"})
		return
	}
	defer rows.Close()

	var favorites []models.Favorite

	for rows.Next() {
		var fav models.Favorite
		err := rows.Scan(&fav.ID, &fav.UserID, &fav.CryptoID, &fav.CryptoName, &fav.CreatedAt)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error leyendo datos"})
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

// Adds crypto to favorite
func (h *FavoritesHandler) PostFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDInt := userID.(int)

	var requestBody struct {
		CryptoID   string `json:"cryptoId" binding:"required"`
		CryptoName string `json:"cryptoName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "CryptoID y CryptoName son requeridos"})
		return
	}

	checkQuery := `SELECT id FROM user_favorites WHERE user_id = $1 AND crypto_id = $2`
	var existingID int
	err := h.DB.QueryRow(checkQuery, userIDInt, requestBody.CryptoID).Scan(&existingID)

	if err == nil {
		c.JSON(409, gin.H{"error": "Esta crypto ya está en favoritos"})
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
		c.JSON(500, gin.H{"error": "Error insertando en BD"})
		return
	}

	c.JSON(201, gin.H{
		"message": "Crypto añadida a favoritos",
		"data":    newFavorite,
	})
}

// Delete from favorites
func (h *FavoritesHandler) DeleteFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Usuario no autenticado"})
		return
	}

	userIDInt := userID.(int)
	cryptoID := c.Param("cryptoId")

	if cryptoID == "" {
		c.JSON(400, gin.H{"error": "CryptoId es requerido"})
		return
	}

	deleteQuery := `
		DELETE FROM user_favorites
		WHERE user_id = $1 AND crypto_id = $2
	`

	result, err := h.DB.Exec(deleteQuery, userIDInt, cryptoID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error eliminando de BD"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(404, gin.H{"error": "Crypto no encontrada en favoritos"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Crypto eliminada de favoritos",
	})
}
