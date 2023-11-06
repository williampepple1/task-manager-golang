package routes

import (
	"task-manager/handlers"
	"task-manager/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// User routes
	r.POST("/register", handlers.RegisterUser(db))
	r.POST("/login", handlers.LoginUser(db))

	// Use the authorization middleware for the following routes
	authorized := r.Group("/")
	authorized.Use(middleware.Authorize())
	{
		authorized.GET("/tasks", handlers.ListTasks(db))
		authorized.POST("/tasks", handlers.CreateTask(db))
		authorized.GET("/tasks/:id", handlers.GetTask(db))
		authorized.PUT("/tasks/:id", handlers.UpdateTask(db))
		authorized.DELETE("/tasks/:id", handlers.DeleteTask(db))
	}
}
