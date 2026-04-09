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
	fmt.Printf("Connecting to PostgreSQL: %s:%s\n", cfg.Host, cfg.Port)

	err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.CloseDB()

	router := gin.Default()

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	favHandler := handlers.NewFavoritesHandler(database.DB)

	protected := router.Group("/favorites")
	protected.Use(middleware.AuthMiddleware(cfg.AuthAPI))
	{
		protected.GET("", favHandler.GetFavorites)
		protected.POST("", favHandler.PostFavorite)
		protected.DELETE("/:cryptoId", favHandler.DeleteFavorite)
	}

	fmt.Printf("Server listening on :%s\n", cfg.AppPort)
	router.Run(":" + cfg.AppPort)
}