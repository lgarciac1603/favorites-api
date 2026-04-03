package handlers

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lgarciac1603/favorites-api/models"
)

var favoritesDB = make(map[string][]models.Favorite)

func GetFavorites(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{ "error": "User not found" })
		return
	}

	userIDStr := userID.(string)

	favorites := favoritesDB[userIDStr]

	if favorites == nil {
		favorites = []models.Favorite{}
	}

	c.JSON(200, gin.H{
		"data": favorites,
		"total": len(favorites),
	})
}

// Add favorite crypto
func PostFavorite(c *gin.Context) {
	userID, exists := c.Get("iserID")
	if !exists {
		c.JSON(401, gin.H{ "error": "Unauthenticated user" })
		return
	}

	userIDStr := userID.(string)

	var requestBody struct {
		CryptoID 		string `json:"cryptoId" binding:"required"`
		CryptoName 	string `json:"cryptoName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{ "error": "CryptoID and CryptoName are required" })
		return
	}

	for _, fav := range favoritesDB[userIDStr] {
		if fav.CryptoId == requestBody.CryptoID {
			c.JSON(409, gin.H{ "error": "Crypto already on favorites" })
			return
		}
	}

	newFavorite := models.Favorite{
		ID: generateID(), // auxiliar function
		UserID:     userIDStr,
		CryptoId:   requestBody.CryptoID,
		CryptoName: requestBody.CryptoName,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	favoritesDB[userIDStr] = append(favoritesDB[userIDStr], newFavorite)

	c.JSON(201, gin.H{
		"message": "Crypto added to users favorites",
		"data": newFavorite,
	})
}

// Delete favorite crypto
func DeleteFavorite(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{ "error": "Unauthenticated user" })
		return
	}

	userIDStr := userID.(string)
	
	// Obtener cryptoId de la URL
	cryptoID := c.Param("cryptoId")

	if cryptoID == "" {
		c.JSON(400, gin.H{"error": "CryptoId is required"})
		return
	}

	favorites := favoritesDB[userIDStr]
	found := false

	for i, fav := range favorites {
		if fav.CryptoId == cryptoID {
			// Eliminar del slice
			favoritesDB[userIDStr] = append(favorites[:i], favorites[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(404, gin.H{"error": "Crypto no encontrada en favoritos"})
		return
	}
	
	c.JSON(200, gin.H{
		"message": "Crypto eliminada de favoritos",
	})
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
