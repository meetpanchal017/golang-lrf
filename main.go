package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	routes "golang-jwt-demo/routes"
)

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get the port from environment variables or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.PublicRoutes(router)
	routes.UserRoutes(router)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
