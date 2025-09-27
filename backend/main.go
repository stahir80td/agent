package main

import (
	"log"
	"os"
	"trading-platform/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// API routes FIRST (before static files!)
	api := r.Group("/api")
	{
		api.POST("/chat", handlers.HandleChat)
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Static files AFTER API routes
	r.Static("/assets", "./frontend/dist/assets")

	// Serve index.html for all other routes (SPA routing)
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	r.Run(":" + port)
}
