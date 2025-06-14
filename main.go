package main

import (
	"Go_Gin_To-Do_List_API/config"
	"Go_Gin_To-Do_List_API/database"
	"Go_Gin_To-Do_List_API/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// @title Go Gin To-Do List API
// @version 1.0
// @description This is a secure to-do list API with user authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load config and generate JWT secret if needed
	config.LoadAndInitConfig()

	// Set Gin mode based on an environment variable.
	// Defaults to debug mode if GIN_MODE is not set.
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	database.ConnectDatabase()

	// Set up router
	r := router.SetupRouter()

	// Start the server
	log.Println("Starting server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
