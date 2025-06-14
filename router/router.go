package router

import (
	"Go_Gin_To-Do_List_API/auth"
	"Go_Gin_To-Do_List_API/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin router for a pure JSON API.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Set a trusted proxy to prevent header spoofing.
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}

	// Simple health check endpoint at the root.
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "API is running"})
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Public routes for authentication
		public := api.Group("/")
		{
			public.POST("/register", handlers.Register)
			public.POST("/login", handlers.Login)
			public.GET("/verify-email", handlers.VerifyEmail)
			public.POST("/forgot-password", handlers.ForgotPassword)
			public.POST("/reset-password", handlers.ResetPassword)
		}

		// Protected routes that require JWT authentication
		protected := api.Group("/")
		protected.Use(auth.Middleware())
		{
			// Task routes
			protected.POST("/tasks", handlers.CreateTask)
			protected.GET("/tasks", handlers.GetTasks)
			protected.GET("/tasks/:id", handlers.GetTask)
			protected.PUT("/tasks/:id", handlers.UpdateTask)
			protected.DELETE("/tasks/:id", handlers.DeleteTask)
		}
	}
	return r
}
