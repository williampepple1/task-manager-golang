package routes

import (
	"os"
	"task-manager/handlers"
	"task-manager/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Retrieve the secret key from environment variables or configuration
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		// Handle the missing secret key appropriately
		panic("JWT secret key must be set")
	}

	// User routes
	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))

	// Use the authorization middleware for the following routes
	authorized := r.Group("/")
	authorized.Use(middleware.Authorize())
	{
		authorized.GET("/tasks", handlers.ListTasks(db))
		authorized.POST("/tasks", handlers.CreateTask(db, secretKey))
		authorized.GET("/tasks/:id", handlers.GetTask(db))
		authorized.PUT("/tasks/:id", handlers.UpdateTask(db))
		authorized.DELETE("/tasks/:id", handlers.DeleteTask(db))
	}
}
