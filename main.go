// main.go
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lgarciac1603/favorites-api/config"
	"github.com/lgarciac1603/favorites-api/database"
	"github.com/lgarciac1603/favorites-api/handlers"
	"github.com/lgarciac1603/favorites-api/middleware"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Printf("Conectando a PostgreSQL: %s:%s\n", cfg.Host, cfg.Port)

	err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error inicializando BD: %v", err)
	}
	defer database.CloseDB()

	router := gin.Default()

	// Create handler injecting DB
	favHandler := handlers.NewFavoritesHandler(database.DB)

	protected := router.Group("/favorites")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("", favHandler.GetFavorites)
		protected.POST("", favHandler.PostFavorite)
		protected.DELETE("/:cryptoId", favHandler.DeleteFavorite)
	}

	fmt.Println("Servidor escuchando en :8080")
	router.Run(":8080")
}